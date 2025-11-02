package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

// ProjectConfig holds all configuration for a project
type ProjectConfig struct {
	Project   ProjectSection `toml:"project"`
	Workflow  WorkflowConfig `toml:"workflow"`
	Templates TemplateConfig `toml:"templates"`
	Storage   StorageConfig  `toml:"storage"`
}

// ProjectSection holds project metadata
type ProjectSection struct {
	Name        string `toml:"name"`
	Description string `toml:"description"`
	Owner       string `toml:"owner"`
}

// WorkflowConfig defines the task status workflow
type WorkflowConfig struct {
	Statuses    []string           `toml:"statuses"`
	Initial     string             `toml:"initial"`
	Transitions []TransitionConfig `toml:"transitions"`
}

// TransitionConfig defines allowed status transitions
type TransitionConfig struct {
	From string   `toml:"from"`
	To   []string `toml:"to"`
}

// TemplateConfig holds paths to template files
type TemplateConfig struct {
	Epic     string `toml:"epic"`
	Plan     string `toml:"plan"`
	Research string `toml:"research"`
	Story    string `toml:"story"`
	Decision string `toml:"decision"`
	Task     string `toml:"task"`
}

// StorageConfig holds storage backend configuration
type StorageConfig struct {
	Backend string            `toml:"backend"`
	Path    string            `toml:"path"`
	Options map[string]string `toml:"options"`
}

// DefaultWorkflow returns the default workflow configuration
func DefaultWorkflow() WorkflowConfig {
	return WorkflowConfig{
		Statuses: []string{"todo", "in-progress", "reviewing", "done", "archived"},
		Initial:  "todo",
		Transitions: []TransitionConfig{
			{From: "todo", To: []string{"in-progress", "archived"}},
			{From: "in-progress", To: []string{"reviewing", "todo", "archived"}},
			{From: "reviewing", To: []string{"done", "in-progress", "archived"}},
			{From: "done", To: []string{"archived"}},
		},
	}
}

// DefaultStorage returns the default storage configuration
func DefaultStorage() StorageConfig {
	return StorageConfig{
		Backend: "markdown-fs",
		Path:    "", // Will be resolved by caller to project root
		Options: make(map[string]string),
	}
}

// DefaultTemplates returns the default template configuration
func DefaultTemplates() TemplateConfig {
	return TemplateConfig{}
}

// LoadConfig loads configuration from a config.toml file
// If the file doesn't exist, returns default configuration
func LoadConfig(configPath string) (*ProjectConfig, error) {
	// Check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Return defaults if file doesn't exist
		return &ProjectConfig{
			Workflow:  DefaultWorkflow(),
			Storage:   DefaultStorage(),
			Templates: DefaultTemplates(),
		}, nil
	}

	// Read and parse config file
	var config ProjectConfig
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Apply defaults for missing sections (but preserve values that were set)
	if len(config.Workflow.Statuses) == 0 {
		config.Workflow = DefaultWorkflow()
	}
	if config.Storage.Backend == "" {
		defaultStorage := DefaultStorage()
		config.Storage.Backend = defaultStorage.Backend
		// Don't override path if it was explicitly set
		if config.Storage.Path == "" {
			config.Storage.Path = defaultStorage.Path
		}
		// Don't override options if they were set
		if len(config.Storage.Options) == 0 {
			config.Storage.Options = defaultStorage.Options
		}
	}

	// Resolve relative paths
	configDir := filepath.Dir(configPath)
	if config.Storage.Path != "" && !filepath.IsAbs(config.Storage.Path) {
		config.Storage.Path = filepath.Join(configDir, config.Storage.Path)
	} else if config.Storage.Path == "" {
		config.Storage.Path = configDir
	}

	return &config, nil
}

// ValidateWorkflow checks if the workflow configuration is valid
func (wf *WorkflowConfig) ValidateWorkflow() error {
	// Check that initial status exists in statuses
	initialExists := false
	for _, s := range wf.Statuses {
		if s == wf.Initial {
			initialExists = true
			break
		}
	}
	if !initialExists {
		return fmt.Errorf("initial status '%s' not found in statuses", wf.Initial)
	}

	// Check that all transition sources/destinations exist in statuses
	for _, t := range wf.Transitions {
		// Check 'from' status
		fromExists := false
		for _, s := range wf.Statuses {
			if s == t.From {
				fromExists = true
				break
			}
		}
		if !fromExists {
			return fmt.Errorf("transition source '%s' not found in statuses", t.From)
		}

		// Check 'to' statuses
		for _, toStatus := range t.To {
			toExists := false
			for _, s := range wf.Statuses {
				if s == toStatus {
					toExists = true
					break
				}
			}
			if !toExists {
				return fmt.Errorf("transition destination '%s' not found in statuses", toStatus)
			}
		}
	}

	return nil
}

// IsValidTransition checks if a status transition is allowed
func (wf *WorkflowConfig) IsValidTransition(from, to string) bool {
	for _, t := range wf.Transitions {
		if t.From == from {
			for _, toStatus := range t.To {
				if toStatus == to {
					return true
				}
			}
		}
	}
	return false
}

// ResolveTemplate resolves a template file path using the hierarchy
// Checks in order: config, project root, XDG_DATA_HOME, built-in
func ResolveTemplate(taskType, configPath string) (string, error) {
	// Map task type to template filename
	var defaultTemplate string

	switch strings.ToLower(taskType) {
	case "epic":
		defaultTemplate = "epic.md"
	case "plan":
		defaultTemplate = "plan.md"
	case "research":
		defaultTemplate = "research.md"
	case "story":
		defaultTemplate = "story.md"
	case "decision":
		defaultTemplate = "decision.md"
	case "task":
		defaultTemplate = "task.md"
	default:
		return "", fmt.Errorf("unknown task type: %s", taskType)
	}

	// Try loading config to get custom template path
	config, err := LoadConfig(configPath)
	if err == nil {
		// Check config-specified path
		configDir := filepath.Dir(configPath)
		var configTemplatePath string

		switch strings.ToLower(taskType) {
		case "epic":
			configTemplatePath = config.Templates.Epic
		case "plan":
			configTemplatePath = config.Templates.Plan
		case "research":
			configTemplatePath = config.Templates.Research
		case "story":
			configTemplatePath = config.Templates.Story
		case "decision":
			configTemplatePath = config.Templates.Decision
		case "task":
			configTemplatePath = config.Templates.Task
		}

		if configTemplatePath != "" {
			if filepath.IsAbs(configTemplatePath) {
				if _, err := os.Stat(configTemplatePath); err == nil {
					return configTemplatePath, nil
				}
			} else {
				resolvedPath := filepath.Join(configDir, configTemplatePath)
				if _, err := os.Stat(resolvedPath); err == nil {
					return resolvedPath, nil
				}
			}
		}
	}

	// Try project root templates directory
	configDir := filepath.Dir(configPath)
	projectTemplatePath := filepath.Join(configDir, "templates", defaultTemplate)
	if _, err := os.Stat(projectTemplatePath); err == nil {
		return projectTemplatePath, nil
	}

	// Try XDG_DATA_HOME
	xdgDataHome := os.Getenv("XDG_DATA_HOME")
	if xdgDataHome != "" {
		xdgTemplatePath := filepath.Join(xdgDataHome, "opentask", "templates", defaultTemplate)
		if _, err := os.Stat(xdgTemplatePath); err == nil {
			return xdgTemplatePath, nil
		}
	}

	// Fall back to HOME/.local/share
	homeDir, err := os.UserHomeDir()
	if err == nil {
		homeTemplatePath := filepath.Join(homeDir, ".local", "share", "opentask", "templates", defaultTemplate)
		if _, err := os.Stat(homeTemplatePath); err == nil {
			return homeTemplatePath, nil
		}
	}

	// No template found - return empty string (built-in will be used)
	return "", nil
}
