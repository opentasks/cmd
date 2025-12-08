package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
// Note: ActiveProject is NOT merged here - it's derived at runtime via ResolveActiveProject().
func MergeGlobalConfig(resolved *OpentaskResolvedConfig, global *OpentaskGlobalConfigFile) *OpentaskResolvedConfig {
	if global == nil {
		return resolved
	}

	// Merge workflow if present (overrides built-in defaults)
	if global.Workflow != nil {
		resolved.Workflow = global.Workflow
	}

	// Merge templates if present (overrides built-in defaults)
	if global.Templates != nil {
		resolved.Templates = global.Templates
	}

	// Merge graph config if present (overrides built-in defaults)
	if global.Graph != nil {
		resolved.Graph = global.Graph
	}

	return resolved
}

// MergeProjectConfig merges project config into the resolved config.
// Project config has higher priority than global config.
// Priority: project.Project > project.Core > global config > built-in defaults
func MergeProjectConfig(resolved *OpentaskResolvedConfig, project *OpentaskProjectConfigFile) *OpentaskResolvedConfig {
	if project == nil {
		return resolved
	}

	// Merge top-level project fields (now flattened, no intermediate schema)
	if project.Project != nil {
		resolved.Project = project.Project
	}

	// Merge storage field-by-field to preserve defaults
	if project.Storage != nil {
		if project.Storage.Backend != "" {
			resolved.Storage.Backend = project.Storage.Backend
		}
		if project.Storage.Path != "" {
			resolved.Storage.Path = project.Storage.Path
		}
		if len(project.Storage.Options) > 0 {
			resolved.Storage.Options = project.Storage.Options
		}
	}

	if project.Workflow != nil {
		resolved.Workflow = project.Workflow
	}

	if project.Templates != nil {
		resolved.Templates = project.Templates
	}

	if project.Graph != nil {
		resolved.Graph = project.Graph
	}

	// Merge project core config (fallback if top-level fields not set)
	// Priority: project.X > project.Core.X > global.Core.X > built-in defaults
	// Note: ActiveProject is NOT merged here - it's derived at runtime via ResolveActiveProject()
	if project.Core != nil {
		// If project doesn't have workflow, use project.Core workflow
		if project.Workflow == nil && project.Core.Workflow != nil {
			resolved.Workflow = project.Core.Workflow
		}
		// If project doesn't have templates, use project.Core templates
		if project.Templates == nil && project.Core.Templates != nil {
			resolved.Templates = project.Core.Templates
		}
		// If project doesn't have graph config, use project.Core graph config
		if project.Graph == nil && project.Core.Graph != nil {
			resolved.Graph = project.Core.Graph
		}
	}

	return resolved
}

// ResolveActiveProject derives the active project ID from runtime context.
// This is a pure function: same inputs always produce same output.
//
// Resolution priority:
//  1. Explicit project.id in local .opentask.toml
//  2. Context match in global config (longest path wins)
//  3. Directory name fallback (if local config exists AND no global match)
//  4. Unresolved (empty string)
//
// Returns (projectID, isResolved) where isResolved indicates whether a project was successfully identified.
func ResolveActiveProject(cwd string, localConfig *OpentaskProjectConfigFile, globalConfig *OpentaskGlobalConfigFile) (string, bool) {
	// Priority 1: Explicit project.id in local .opentask.toml
	if localConfig != nil && localConfig.Project != nil && localConfig.Project.ID != "" {
		return localConfig.Project.ID, true
	}

	// Priority 2: Context match in global config
	if globalConfig != nil {
		projectID, matchedProject := FindProjectByContext(cwd, globalConfig.Projects)
		if matchedProject != nil {
			return projectID, true
		}
	}

	// Priority 3: Derive from local .opentask.toml directory name (fallback if no global match found)
	if localConfig != nil {
		return filepath.Base(cwd), true
	}

	// Priority 4: Unresolved - onboarding required
	return "", false
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
//
// ActiveProject is ALWAYS derived at runtime via ResolveActiveProject() - never read from files.
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
	var localConfig *OpentaskProjectConfigFile
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
				// Save the closest (first) local config for resolution
				if i == 0 {
					localConfig = projectConfig
				}
				resolved = MergeProjectConfig(resolved, projectConfig)
				// Insert at beginning to maintain order (closest first)
				resolved.DiscoveredFiles = append([]string{projectFile}, resolved.DiscoveredFiles...)
			}
		}
	}

	// ═══════════════════════════════════════════════════════════
	// DERIVE ACTIVE PROJECT (pure function - never read from file)
	// ═══════════════════════════════════════════════════════════
	resolved.ActiveProject, resolved.IsResolved = ResolveActiveProject(cwd, localConfig, globalConfig)

	// If resolved, merge settings from matching global project
	if resolved.IsResolved && globalConfig != nil {
		for _, globalProj := range globalConfig.Projects {
			if globalProj.ID == resolved.ActiveProject {
				// Found matching global project, merge its settings

				// Merge storage if not already set
				if resolved.Storage.Path == "" && globalProj.Storage != nil {
					if globalProj.Storage.Path != "" {
						resolved.Storage.Path = globalProj.Storage.Path
					}
				}

				// Merge templates if all are empty (only set if project config didn't provide them)
				if resolved.Templates != nil && resolved.Templates.Epic == "" && resolved.Templates.Plan == "" &&
					resolved.Templates.Research == "" && resolved.Templates.Story == "" &&
					resolved.Templates.Decision == "" && resolved.Templates.Task == "" &&
					globalProj.Templates != nil {
					resolved.Templates = globalProj.Templates
				}

				// Merge workflow if not already set (only if project config didn't provide it)
				if (resolved.Workflow == nil || len(resolved.Workflow.Statuses) == 0) && globalProj.Workflow != nil {
					resolved.Workflow = globalProj.Workflow
				}

				break
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

	return resolved, nil
}

// FindProjectByContext finds a project in global config that matches the given cwd.
// Returns the project ID and project config if found, empty string if no match.
// Matching algorithm:
// 1. Find longest matching context path (most specific wins)
// 2. Return that project's ID
// If multiple projects have overlapping contexts, longest match takes precedence.
func FindProjectByContext(cwd string, globalProjects []GlobalProjectConfig) (string, *GlobalProjectConfig) {
	// Convert cwd to absolute path
	cwdAbs, err := filepath.Abs(cwd)
	if err != nil {
		return "", nil
	}

	// Normalize path and expand ~
	if len(cwdAbs) > 0 && cwdAbs[0] == '~' {
		home, err := os.UserHomeDir()
		if err == nil {
			cwdAbs = filepath.Join(home, cwdAbs[1:])
		}
	}
	cwdAbs = filepath.Clean(cwdAbs)

	var longestMatchPath string
	var matchedProject *GlobalProjectConfig

	// Check all projects for context matches
	for i := range globalProjects {
		for _, ctx := range globalProjects[i].Context {
			// Normalize context path
			ctxPath := filepath.Clean(ctx.Path)

			// Expand ~ to home directory
			if len(ctxPath) > 0 && ctxPath[0] == '~' {
				home, err := os.UserHomeDir()
				if err == nil {
					ctxPath = filepath.Join(home, ctxPath[1:])
				}
			}

			// Convert to absolute
			ctxPath, err := filepath.Abs(ctxPath)
			if err != nil {
				continue
			}
			ctxPath = filepath.Clean(ctxPath)

			// Check if cwd is under this context path
			// Use HasPrefix to check if cwd starts with context path
			if cwdAbs == ctxPath || strings.HasPrefix(cwdAbs+string(filepath.Separator), ctxPath+string(filepath.Separator)) {
				// This context matches
				// Keep the longest match (most specific)
				if len(ctxPath) > len(longestMatchPath) {
					longestMatchPath = ctxPath
					matchedProject = &globalProjects[i]
				}
			}
		}
	}

	if matchedProject != nil {
		return matchedProject.ID, matchedProject
	}

	return "", nil
}

// AddContextPath adds a context path to a project.
// Resolves the path to absolute and normalizes it.
func (p *GlobalProjectConfig) AddContextPath(path string) error {
	// Expand ~ to home directory
	if len(path) > 0 && path[0] == '~' {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to expand home directory: %w", err)
		}
		path = filepath.Join(home, path[1:])
	}

	// Convert to absolute
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	absPath = filepath.Clean(absPath)

	// Check if path already exists
	for _, ctx := range p.Context {
		if filepath.Clean(ctx.Path) == absPath {
			return fmt.Errorf("context path already exists: %s", absPath)
		}
	}

	// Add the context
	p.Context = append(p.Context, ProjectContext{Path: absPath})
	return nil
}

// RemoveContextPath removes a context path from a project.
func (p *GlobalProjectConfig) RemoveContextPath(path string) error {
	// Expand ~ to home directory
	if len(path) > 0 && path[0] == '~' {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to expand home directory: %w", err)
		}
		path = filepath.Join(home, path[1:])
	}

	// Convert to absolute
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	absPath = filepath.Clean(absPath)

	// Find and remove the context
	for i, ctx := range p.Context {
		if filepath.Clean(ctx.Path) == absPath {
			p.Context = append(p.Context[:i], p.Context[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("context path not found: %s", absPath)
}
