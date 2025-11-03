package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// LoadGlobalConfig loads and parses a global configuration file.
// Returns nil if file doesn't exist.
func LoadGlobalConfig(path string) (*OpentaskGlobalConfigFile, error) {
	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, nil
	}

	var cfg OpentaskGlobalConfigFile
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse global config: %w", err)
	}

	return &cfg, nil
}

// LoadProjectConfig loads and parses a project configuration file.
// Returns nil if file doesn't exist.
func LoadProjectConfig(path string) (*OpentaskProjectConfigFile, error) {
	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, nil
	}

	var cfg OpentaskProjectConfigFile
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse project config: %w", err)
	}

	return &cfg, nil
}

// MergeGlobalConfig merges global config into the resolved config.
// Global config provides defaults that can be overridden by project configs.
// Note: We start with defaults, so global config will override them if present.
func MergeGlobalConfig(resolved *OpentaskResolvedConfig, global *OpentaskGlobalConfigFile) *OpentaskResolvedConfig {
	if global == nil {
		return resolved
	}

	// Merge global schema if present
	if global.Global != nil {
		if global.Global.ActiveProject != "" {
			resolved.ActiveProject = global.Global.ActiveProject
		}
	}

	// Merge core config - always apply if present (overrides built-in defaults)
	if global.Core != nil {
		if global.Core.Workflow != nil {
			resolved.Workflow = global.Core.Workflow
		}
		if global.Core.Templates != nil {
			resolved.Templates = global.Core.Templates
		}
	}

	return resolved
}

// MergeProjectConfig merges project config into the resolved config.
// Project config has higher priority than global config.
func MergeProjectConfig(resolved *OpentaskResolvedConfig, project *OpentaskProjectConfigFile) *OpentaskResolvedConfig {
	if project == nil {
		return resolved
	}

	// Merge project schema
	if project.Project != nil {
		if project.Project.Project != nil {
			resolved.Project = project.Project.Project
		}

		if project.Project.Storage != nil {
			resolved.Storage = project.Project.Storage
		}

		if project.Project.Workflow != nil {
			resolved.Workflow = project.Project.Workflow
		}

		if project.Project.Templates != nil {
			resolved.Templates = project.Project.Templates
		}

		if project.Project.ActiveProject != "" {
			resolved.ActiveProject = project.Project.ActiveProject
		}
	}

	// Merge project core config
	// Priority: project.Project.X > project.Core.X > global.Core.X > built-in defaults
	if project.Core != nil {
		// If project.Project doesn't have workflow, use project.Core workflow
		if (project.Project == nil || project.Project.Workflow == nil) && project.Core.Workflow != nil {
			resolved.Workflow = project.Core.Workflow
		}
		// If project.Project doesn't have templates, use project.Core templates
		if (project.Project == nil || project.Project.Templates == nil) && project.Core.Templates != nil {
			resolved.Templates = project.Core.Templates
		}
	}

	return resolved
}

// deriveActiveProject derives the active_project from the config file path and global projects.
// Priority:
// 1. If active_project is already set, return it
// 2. Check if the config directory matches a global project storage path
// 3. Fall back to the config directory name
func deriveActiveProject(activeProject string, configDir string, globalProjects []GlobalProjectConfig) string {
	if activeProject != "" {
		return activeProject
	}

	// Try to match config directory with global project storage path
	// Use absolute path for comparison
	configDirAbs, err := filepath.Abs(configDir)
	if err != nil {
		configDirAbs = configDir
	}

	for _, proj := range globalProjects {
		if proj.Storage == nil {
			continue
		}

		projPath := proj.Storage.Path

		// Expand ~ to home directory
		if projPath[0] == '~' {
			home, err := os.UserHomeDir()
			if err == nil {
				projPath = filepath.Join(home, projPath[1:])
			}
		}

		projPathAbs, err := filepath.Abs(projPath)
		if err != nil {
			continue
		}

		// Check if paths match
		if configDirAbs == projPathAbs {
			return proj.ID
		}
	}

	// Fall back to directory name
	return filepath.Base(configDir)
}

// ResolveProjectConfig resolves the final merged configuration for a project.
// This is the main entry point for the resolution algorithm.
//
// Resolution order (highest to lowest priority):
// 1. Current directory .opentask.toml (project schema fields)
// 2. Parent .opentask.toml files (project schema fields)
// 3. Global config [[projects]] matching project ID (project schema fields)
// 4. Global config [global] section (core schema fields)
// 5. Built-in defaults
func ResolveProjectConfig(cwd string) (*OpentaskResolvedConfig, error) {
	// Start with defaults
	resolved := NewResolvedConfig()

	// Discover project config files
	projectFiles, err := DiscoverConfigFiles(cwd)
	if err != nil {
		return nil, fmt.Errorf("failed to discover config files: %w", err)
	}

	// Load and merge project configs (furthest to closest, so closest overrides)
	// Note: projectFiles may include the global config at the end
	var globalConfig *OpentaskGlobalConfigFile
	var globalPath string

	for i := len(projectFiles) - 1; i >= 0; i-- {
		projectFile := projectFiles[i]

		// Check if this is the global config file
		home, err := os.UserHomeDir()
		isGlobalConfig := false
		if err == nil {
			globalPath = filepath.Join(home, ".config", "opentask", "config.toml")
			isGlobalConfig = (projectFile == globalPath)
		}

		if isGlobalConfig {
			// Load global config separately
			var err error
			globalConfig, err = LoadGlobalConfig(projectFile)
			if err != nil {
				return nil, fmt.Errorf("failed to load global config: %w", err)
			}
			if globalConfig != nil {
				resolved = MergeGlobalConfig(resolved, globalConfig)
				resolved.DiscoveredFiles = append([]string{projectFile}, resolved.DiscoveredFiles...)
			}
		} else {
			// Load project config
			projectConfig, err := LoadProjectConfig(projectFile)
			if err != nil {
				return nil, fmt.Errorf("failed to load project config %s: %w", projectFile, err)
			}

			if projectConfig != nil {
				resolved = MergeProjectConfig(resolved, projectConfig)
				// Insert at beginning to maintain order (closest first)
				resolved.DiscoveredFiles = append([]string{projectFile}, resolved.DiscoveredFiles...)
			}
		}
	}

	// Resolve storage path (must be absolute)
	if resolved.Storage != nil && resolved.Storage.Path != "" {
		storagePath := resolved.Storage.Path

		// Expand ~ to home directory
		if len(storagePath) > 0 && storagePath[0] == '~' {
			home, err := os.UserHomeDir()
			if err == nil {
				storagePath = filepath.Join(home, storagePath[1:])
			}
		}

		// Make absolute relative to cwd
		if !filepath.IsAbs(storagePath) {
			storagePath = filepath.Join(cwd, storagePath)
		}

		resolved.Storage.Path = storagePath
	}

	// Derive active_project if not set
	if resolved.ActiveProject == "" {
		var globalProjects []GlobalProjectConfig
		if globalConfig != nil && globalConfig.Global != nil {
			globalProjects = globalConfig.Global.Projects
		}

		// Find the closest project config directory
		if len(projectFiles) > 0 {
			configDir := filepath.Dir(projectFiles[0])
			resolved.ActiveProject = deriveActiveProject("", configDir, globalProjects)
		} else {
			// No project config, derive from cwd
			resolved.ActiveProject = deriveActiveProject("", cwd, globalProjects)
		}
	}

	return resolved, nil
}
