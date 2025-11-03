package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

// OpentaskConfigCoreSchema defines shared configuration fields used in both global and project contexts.
// This includes workflow definitions and template configurations that can be overridden at different levels.
type OpentaskConfigCoreSchema struct {
	Workflow  *WorkflowConfig `toml:"workflow"`
	Templates *TemplateConfig `toml:"templates"`
}

// OpentaskConfigGlobalSchema defines user-level global configuration.
// This includes the active project selection and list of known projects.
type OpentaskConfigGlobalSchema struct {
	ActiveProject string                `toml:"active_project"`
	Projects      []GlobalProjectConfig `toml:"projects"`
}

// GlobalProjectConfig represents a project entry in the global configuration.
// Projects can be defined at the global level with their storage configuration.
type GlobalProjectConfig struct {
	ID        string          `toml:"id"`
	Name      string          `toml:"name"`
	Storage   *StorageConfig  `toml:"storage"`
	Workflow  *WorkflowConfig `toml:"workflow"`  // Optional project-specific workflow override
	Templates *TemplateConfig `toml:"templates"` // Optional project-specific templates override
}

// OpentaskConfigProjectSchema defines project-level configuration.
// This includes project metadata, storage path, and optional workflow/template overrides specific to this project.
type OpentaskConfigProjectSchema struct {
	Project       *ProjectSection `toml:"project"`
	Storage       *StorageConfig  `toml:"storage"`
	Workflow      *WorkflowConfig `toml:"workflow"`
	Templates     *TemplateConfig `toml:"templates"`
	ActiveProject string          `toml:"active_project"` // May be auto-populated if not specified
}

// OpentaskGlobalConfigFile represents the complete structure of a global config file (~/.config/opentask/config.toml).
// It combines the global schema with shared core settings.
type OpentaskGlobalConfigFile struct {
	Global *OpentaskConfigGlobalSchema `toml:"global"`
	Core   *OpentaskConfigCoreSchema   `toml:"core"` // Optional global defaults for workflow/templates
}

// OpentaskProjectConfigFile represents the complete structure of a project config file (.opentask.toml).
// It combines project-specific settings with core settings that can be overridden at project level.
type OpentaskProjectConfigFile struct {
	Project *OpentaskConfigProjectSchema `toml:"project"`
	Core    *OpentaskConfigCoreSchema    `toml:"core"` // Optional project-specific workflow/templates
}

// OpentaskResolvedConfig is the final merged configuration after resolving from all sources.
// This is what the application uses at runtime - a complete, validated configuration with no ambiguities.
// The DiscoveredFiles field tracks which config files were merged to produce this result.
type OpentaskResolvedConfig struct {
	Project         *ProjectSection `toml:"-"`
	Workflow        *WorkflowConfig `toml:"-"`
	Templates       *TemplateConfig `toml:"-"`
	Storage         *StorageConfig  `toml:"-"`
	ActiveProject   string          `toml:"-"`
	DiscoveredFiles []string        `toml:"-"` // Config files that were merged (for debugging/display)
}

// ProjectConfig holds all configuration for a project (legacy, kept for backwards compatibility during transition)
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

// NewResolvedConfig creates a new resolved config with sensible defaults.
// This is used as the starting point before merging from actual config files.
func NewResolvedConfig() *OpentaskResolvedConfig {
	return &OpentaskResolvedConfig{
		Project: &ProjectSection{},
		Workflow: &WorkflowConfig{
			Statuses: []string{"todo", "in-progress", "reviewing", "done", "archived"},
			Initial:  "todo",
			Transitions: []TransitionConfig{
				{From: "todo", To: []string{"in-progress", "archived"}},
				{From: "in-progress", To: []string{"reviewing", "todo", "archived"}},
				{From: "reviewing", To: []string{"done", "in-progress", "archived"}},
				{From: "done", To: []string{"archived"}},
			},
		},
		Templates:       &TemplateConfig{},
		Storage:         &StorageConfig{Backend: "markdown-fs"},
		ActiveProject:   "",
		DiscoveredFiles: []string{},
	}
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
