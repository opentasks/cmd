package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/charmbracelet/glamour"
	"github.com/spf13/cobra"
	"github.com/zenobi-us/opentask/internal/config"
)

// buildConfigFileTree creates a list of resolved config files with merge flow indicators
// Files are shown in priority order (highest first)
func buildConfigFileTree(files []string) string {
	var result strings.Builder

	// Build list with all items (config files + defaults)
	allItems := make([]string, len(files)+1)
	for i, file := range files {
		// Get path relative to current directory for most readable display
		cwd, err := os.Getwd()
		var displayPath string

		// First, try to show as relative path from cwd
		if err == nil {
			relPath, err := filepath.Rel(cwd, file)
			if err == nil && !strings.HasPrefix(relPath, "..") {
				// Only use relative path if it doesn't go up many directories
				displayPath = relPath
			}
		}

		// If not a good relative path, try to use ~ for home directory
		if displayPath == "" {
			homeDir, err := os.UserHomeDir()
			if err == nil && strings.HasPrefix(file, homeDir) {
				displayPath = "~" + file[len(homeDir):]
			}
		}

		// Fall back to absolute path if nothing else worked
		if displayPath == "" {
			displayPath = file
		}

		allItems[i] = displayPath
	}
	allItems[len(files)] = "(builtin) defaults"

	// Render as vertical list
	for i, item := range allItems {
		if i == len(allItems)-1 {
			// Last item
			result.WriteString(fmt.Sprintf("└── %s\n", item))
		} else {
			// Not last item
			result.WriteString(fmt.Sprintf("├── %s\n", item))
			result.WriteString("│   ↓\n")
		}
	}

	return result.String()
}

// configCmd represents the config command group
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
	Long:  "Commands for viewing and managing configuration",
}

// configViewCmd views the resolved configuration with discovery details
var configViewCmd = &cobra.Command{
	Use:   "view",
	Short: "View resolved configuration",
	Long:  "Display the resolved project configuration after discovering and merging all config files",
	RunE: func(cmd *cobra.Command, args []string) error {
		pathFlag, _ := cmd.Flags().GetBool("path")
		jsonFlag, _ := cmd.Flags().GetBool("json")
		verboseFlag, _ := cmd.Flags().GetBool("verbose")

		// Get current directory to start resolution
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current working directory: %w", err)
		}

		// Resolve configuration using new schema
		resolved, err := config.ResolveProjectConfig(cwd)
		if err != nil {
			return fmt.Errorf("failed to resolve configuration: %w", err)
		}

		// If --path flag, just show the resolved storage path
		if pathFlag {
			fmt.Println(resolved.Storage.Path)
			return nil
		}

		// If --json flag, output JSON
		if jsonFlag {
			jsonOutput := map[string]interface{}{
				"discoveredFiles": resolved.DiscoveredFiles,
				"resolved":        resolved,
			}
			data, err := json.MarshalIndent(jsonOutput, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal JSON: %w", err)
			}
			fmt.Println(string(data))
			return nil
		}

		// Build file tree
		fileTree := buildConfigFileTree(resolved.DiscoveredFiles)

		// Convert resolved config to legacy ProjectConfig format for display
		displayConfig := &config.ProjectConfig{
			Project:   *resolved.Project,
			Workflow:  *resolved.Workflow,
			Templates: *resolved.Templates,
			Storage:   *resolved.Storage,
		}

		// Convert config to TOML string
		tomlStr, err := config.ConfigAsToml(displayConfig)
		if err != nil {
			return fmt.Errorf("failed to format config as TOML: %w", err)
		}

		// Prepare template data
		templateData := map[string]interface{}{
			"FoundFiles":                 resolved.DiscoveredFiles,
			"MergingOrder":               resolved.DiscoveredFiles,
			"ResolvedConfigAsTomlString": tomlStr,
			"StopReason":                 "reached filesystem root",
			"StopDir":                    cwd,
			"FileTree":                   fileTree,
		}

		// Execute template
		tmpl, err := template.New("configView").
			Funcs(template.FuncMap{
				"quote": func(s string) string {
					return `"` + s + `"`
				},
				"add": func(a, b int) int {
					return a + b
				},
			}).
			Parse(config.ViewTemplate)
		if err != nil {
			return fmt.Errorf("failed to parse template: %w", err)
		}

		// Render markdown to a buffer first
		var mdBuf bytes.Buffer
		if err := tmpl.Execute(&mdBuf, templateData); err != nil {
			return fmt.Errorf("failed to execute template: %w", err)
		}

		// Render markdown with syntax highlighting using glamour
		renderer, err := glamour.NewTermRenderer(
			glamour.WithAutoStyle(),
			glamour.WithWordWrap(120),
		)
		if err != nil {
			return fmt.Errorf("failed to create markdown renderer: %w", err)
		}

		output, err := renderer.Render(mdBuf.String())
		if err != nil {
			return fmt.Errorf("failed to render markdown: %w", err)
		}

		fmt.Print(output)

		// Verbose output
		if verboseFlag {
			fmt.Println("\n=== Verbose Mode: Merging Details ===")

			// Show discovered config files
			if len(resolved.DiscoveredFiles) > 0 {
				for i, file := range resolved.DiscoveredFiles {
					fmt.Printf("\n[Step %d] Applying: %s\n", i+1, file)
					// For resolved config, just show what was merged
					if resolved.Project.Name != "" {
						fmt.Printf("  - project.name: %q\n", resolved.Project.Name)
					}
					if len(resolved.Workflow.Statuses) > 0 {
						fmt.Printf("  - workflow.statuses: %v\n", resolved.Workflow.Statuses)
					}
					if resolved.Storage.Path != "" {
						fmt.Printf("  - storage.path: %q\n", resolved.Storage.Path)
					}
				}
			}

			// Show defaults as final virtual layer
			stepNum := len(resolved.DiscoveredFiles) + 1
			fmt.Printf("\n[Step %d] (Virtual) Default configuration\n", stepNum)
			defaults := config.ProjectConfig{
				Workflow:  config.DefaultWorkflow(),
				Templates: config.DefaultTemplates(),
				Storage:   config.DefaultStorage(),
			}
			fmt.Printf("  - workflow.statuses: %v\n", defaults.Workflow.Statuses)
			fmt.Printf("  - workflow.initial: %q\n", defaults.Workflow.Initial)
			fmt.Printf("  - storage.backend: %q\n", defaults.Storage.Backend)
		}

		return nil
	},
}

// configInitCmd initializes a new opentask project
var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new project",
	Long:  "Create a new .opentask.toml configuration file in the current directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		nameFlag, _ := cmd.Flags().GetString("name")
		storagePath, _ := cmd.Flags().GetString("storage")
		forceFlag, _ := cmd.Flags().GetBool("force")

		// Get current working directory
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current working directory: %w", err)
		}

		configPath := filepath.Join(cwd, ".opentask.toml")

		// Check if config already exists
		if _, err := os.Stat(configPath); err == nil && !forceFlag {
			return fmt.Errorf(".opentask.toml already exists in %s. Use --force to overwrite", cwd)
		}

		// Determine project name
		projectName := nameFlag
		if projectName == "" {
			// Use directory name as default
			projectName = filepath.Base(cwd)
		}

		// Create config content using new schema structure
		// Note: active_project will be auto-derived if not set
		configContent := fmt.Sprintf(`# opentask project configuration for %s
# This file defines project-specific settings

# Project metadata
[project.project]
name = %q
description = ""
owner = ""

# Storage configuration (project-specific)
[project.storage]
backend = "markdown-fs"
path = %q

# Project-specific workflow (optional - comment out to use global defaults)
# [project.workflow]
# statuses = ["todo", "in-progress", "reviewing", "done", "archived"]
# initial = "todo"

# [[project.workflow.transitions]]
# from = "todo"
# to = ["in-progress", "archived"]

# [[project.workflow.transitions]]
# from = "in-progress"
# to = ["reviewing", "todo", "archived"]

# [[project.workflow.transitions]]
# from = "reviewing"
# to = ["done", "in-progress", "archived"]

# [[project.workflow.transitions]]
# from = "done"
# to = ["archived"]

# Project-specific templates (optional - comment out to use global templates)
# [project.templates]
# epic = "templates/epic.md"
# plan = "templates/plan.md"
# research = "templates/research.md"
# story = "templates/story.md"
# decision = "templates/decision.md"
# task = "templates/task.md"
`, projectName, projectName, storagePath)

		// Write config file
		if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
			return fmt.Errorf("failed to write config file: %w", err)
		}

		// Print success message
		fmt.Printf("Initialized opentask project in %s\n", cwd)
		fmt.Printf("Created: %s\n", configPath)
		fmt.Printf("Storage: %s (local directory)\n\n", storagePath)
		fmt.Println("Configuration:")
		fmt.Printf("  - Project name: %s\n", projectName)
		fmt.Printf("  - Active project ID will be auto-derived from directory name or global config\n\n")
		fmt.Println("Next steps:")
		fmt.Println("  1. Create a task: opentask task new \"Your task title\"")
		fmt.Println("  2. List tasks: opentask task list")
		fmt.Println("  3. View config: opentask config view")
		fmt.Println("  4. Edit config: " + configPath)
		fmt.Println("\nTo set up global configuration:")
		fmt.Println("  1. Create ~/.config/opentask/config.toml with your global settings")
		fmt.Println("  2. Define projects at the global level for multi-project support")

		return nil
	},
}

// configShowCmd shows the current configuration
var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show configuration",
	Long:  "Display the current project configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Current configuration...")
		// TODO: Implement config display
		return nil
	},
}

// configGetCmd gets a config value
var configGetCmd = &cobra.Command{
	Use:   "get [key]",
	Short: "Get a config value",
	Long:  "Get a specific configuration value",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		fmt.Printf("Getting config value for: %s\n", key)
		// TODO: Implement config get
		return nil
	},
}

// configSetCmd sets a config value
var configSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Set a config value",
	Long:  "Set a configuration value",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		value := args[1]
		fmt.Printf("Setting %s = %s\n", key, value)
		// TODO: Implement config set
		return nil
	},
}

// configProjectsCmd manages global projects
var configProjectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "Manage projects",
	Long:  "List and manage projects defined in global configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		activeFlag, _ := cmd.Flags().GetString("active")

		// Load global config
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}

		globalPath := filepath.Join(home, ".config", "opentask", "config.toml")
		globalConfig, err := config.LoadGlobalConfig(globalPath)
		if err != nil {
			return fmt.Errorf("failed to load global config: %w", err)
		}

		if globalConfig == nil || globalConfig.Global == nil {
			fmt.Println("No global configuration found at " + globalPath)
			fmt.Println("\nTo create global config:")
			fmt.Println("  1. mkdir -p ~/.config/opentask")
			fmt.Println("  2. Create config.toml with project definitions")
			return nil
		}

		// Handle --active flag to set active project
		if activeFlag != "" {
			// Validate project exists
			found := false
			for _, proj := range globalConfig.Global.Projects {
				if proj.ID == activeFlag {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("project %q not found in global config", activeFlag)
			}
			globalConfig.Global.ActiveProject = activeFlag
			fmt.Printf("Set active project to: %s\n", activeFlag)
			// TODO: Persist change to global config file
			return nil
		}

		// List projects (default)
		if len(globalConfig.Global.Projects) == 0 {
			fmt.Println("No projects configured in global config")
			return nil
		}

		fmt.Println("Configured projects:")
		fmt.Println()
		for _, proj := range globalConfig.Global.Projects {
			prefix := "  "
			if proj.ID == globalConfig.Global.ActiveProject {
				prefix = "* "
			}
			fmt.Printf("%s%s (%s)\n", prefix, proj.Name, proj.ID)
			if proj.Storage != nil && proj.Storage.Path != "" {
				fmt.Printf("     Path: %s\n", proj.Storage.Path)
			}
		}

		if globalConfig.Global.ActiveProject != "" {
			fmt.Printf("\nActive project: %s\n", globalConfig.Global.ActiveProject)
		}

		return nil
	},
}

func init() {
	configCmd.AddCommand(configInitCmd)
	configCmd.AddCommand(configViewCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configProjectsCmd)

	// Flags for config init
	configInitCmd.Flags().StringP("name", "n", "", "Project name (default: current directory name)")
	configInitCmd.Flags().StringP("storage", "s", "./.tasks/", "Storage path (default: ./.tasks/)")
	configInitCmd.Flags().BoolP("force", "f", false, "Overwrite existing .opentask.toml")

	// Flags for config view
	configViewCmd.Flags().BoolP("path", "p", false, "Show only the resolved storage path")
	configViewCmd.Flags().BoolP("json", "j", false, "Output resolved config as JSON")
	configViewCmd.Flags().BoolP("verbose", "v", false, "Show each config file contents during merging")

	// Flags for config projects
	configProjectsCmd.Flags().StringP("active", "a", "", "Set the active project by ID")
}
