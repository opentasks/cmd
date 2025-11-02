package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// DiscoverConfigFiles walks up from the given file path to find all .opentask.toml files
// Also checks for user global config at ~/.config/opentask/config.toml
// Stops at filesystem root
// Returns files in order from closest to furthest (leaf to root), with user config last
func DiscoverConfigFiles(startPath string) ([]string, error) {
	var found []string

	// Start from the directory containing the config file
	currentDir := filepath.Dir(startPath)

	// If startPath is a directory, use it directly
	info, err := os.Stat(startPath)
	if err == nil && info.IsDir() {
		currentDir = startPath
	}

	// Convert to absolute path
	currentDir, err = filepath.Abs(currentDir)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	for {
		// Check if .opentask.toml exists in current directory
		configPath := filepath.Join(currentDir, ".opentask.toml")
		if _, err := os.Stat(configPath); err == nil {
			found = append(found, configPath)
		}

		// Check if we're at filesystem root
		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			// We're at filesystem root
			break
		}

		// Move to parent directory
		currentDir = parent
	}

	// Check for user global config at ~/.config/opentask/config.toml
	homeDir, err := os.UserHomeDir()
	if err == nil {
		userConfigPath := filepath.Join(homeDir, ".config", "opentask", "config.toml")
		if _, err := os.Stat(userConfigPath); err == nil {
			found = append(found, userConfigPath)
		}
	}

	return found, nil
}

// MergeConfigs merges multiple config files
// Configs are merged in order: earliest (furthest) configs first, latest (closest) override
// Later configs override earlier configs
func MergeConfigs(configPaths []string) (*ProjectConfig, error) {
	// Start with defaults
	merged := &ProjectConfig{
		Project:   ProjectSection{},
		Workflow:  DefaultWorkflow(),
		Templates: DefaultTemplates(),
		Storage:   DefaultStorage(),
	}

	// If no configs provided, return defaults
	if len(configPaths) == 0 {
		return merged, nil
	}

	// Merge each config in order (furthest first, closest last)
	for _, configPath := range configPaths {
		config, err := LoadConfig(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load config %s: %w", configPath, err)
		}

		// Merge project section
		if config.Project.Name != "" {
			merged.Project.Name = config.Project.Name
		}
		if config.Project.Description != "" {
			merged.Project.Description = config.Project.Description
		}
		if config.Project.Owner != "" {
			merged.Project.Owner = config.Project.Owner
		}

		// Merge workflow (override if provided)
		if len(config.Workflow.Statuses) > 0 {
			merged.Workflow.Statuses = config.Workflow.Statuses
		}
		if config.Workflow.Initial != "" {
			merged.Workflow.Initial = config.Workflow.Initial
		}
		if len(config.Workflow.Transitions) > 0 {
			merged.Workflow.Transitions = config.Workflow.Transitions
		}

		// Merge templates (override if provided)
		if config.Templates.Epic != "" {
			merged.Templates.Epic = config.Templates.Epic
		}
		if config.Templates.Plan != "" {
			merged.Templates.Plan = config.Templates.Plan
		}
		if config.Templates.Research != "" {
			merged.Templates.Research = config.Templates.Research
		}
		if config.Templates.Story != "" {
			merged.Templates.Story = config.Templates.Story
		}
		if config.Templates.Decision != "" {
			merged.Templates.Decision = config.Templates.Decision
		}
		if config.Templates.Task != "" {
			merged.Templates.Task = config.Templates.Task
		}

		// Merge storage (override if provided)
		if config.Storage.Backend != "" {
			merged.Storage.Backend = config.Storage.Backend
		}
		if config.Storage.Path != "" {
			merged.Storage.Path = config.Storage.Path
		}
		if len(config.Storage.Options) > 0 {
			merged.Storage.Options = config.Storage.Options
		}
	}

	return merged, nil
}

// LoadConfigHierarchical loads and merges config files walking up from the given path
// Returns the merged configuration and the paths that were found
func LoadConfigHierarchical(startPath string) (*ProjectConfig, []string, error) {
	// Discover all config files
	configPaths, err := DiscoverConfigFiles(startPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to discover configs: %w", err)
	}

	// Reverse to get furthest first
	reversedPaths := make([]string, len(configPaths))
	for i, p := range configPaths {
		reversedPaths[len(configPaths)-1-i] = p
	}

	// Merge all configs
	merged, err := MergeConfigs(reversedPaths)
	if err != nil {
		return nil, configPaths, fmt.Errorf("failed to merge configs: %w", err)
	}

	return merged, configPaths, nil
}

// ConfigDiscoveryInfo provides details about the config discovery process
type ConfigDiscoveryInfo struct {
	FoundFiles    []string // Files found, closest first
	ReversedFiles []string // Same as FoundFiles but reversed (furthest first)
	ResolvedPath  string   // Final resolved storage path (absolute)
	StopReason    string   // Why discovery stopped ("git root" or "filesystem root")
	MergingOrder  []string // How files were merged (furthest first)
}

// DiscoverAndAnalyze performs config discovery and returns analysis info
func DiscoverAndAnalyze(startPath string) (*ConfigDiscoveryInfo, error) {
	// Discover configs
	found, err := DiscoverConfigFiles(startPath)
	if err != nil {
		return nil, err
	}

	// Always stops at filesystem root
	stopReason := "filesystem root"

	// Create reversed list for merging order
	reversed := make([]string, len(found))
	for i, p := range found {
		reversed[len(found)-1-i] = p
	}

	// Load merged config to get resolved path
	merged, err := MergeConfigs(reversed)
	if err != nil {
		return nil, err
	}

	resolvedPath := merged.Storage.Path
	if resolvedPath == "" {
		// Use current directory if not set
		abs, _ := filepath.Abs(startPath)
		dir := abs
		if info, err := os.Stat(abs); err == nil && !info.IsDir() {
			dir = filepath.Dir(abs)
		}
		resolvedPath = dir
	}

	// Ensure absolute path
	if !filepath.IsAbs(resolvedPath) {
		abs, _ := filepath.Abs(resolvedPath)
		resolvedPath = abs
	}

	return &ConfigDiscoveryInfo{
		FoundFiles:    found,
		ReversedFiles: reversed,
		ResolvedPath:  resolvedPath,
		StopReason:    stopReason,
		MergingOrder:  reversed,
	}, nil
}
