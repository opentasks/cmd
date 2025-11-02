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
	"github.com/xlab/treeprint"
	"github.com/zenobi-us/opentask/internal/config"
)

// buildConfigFileTree creates a tree structure of the resolved config files
func buildConfigFileTree(files []string) string {
	if len(files) == 0 {
		return ""
	}

	// Build tree starting from root
	tree := treeprint.New()
	root := findRootPath(files)

	for i, file := range files {
		// Make path relative to root for display
		relPath, err := filepath.Rel(root, file)
		if err != nil {
			relPath = file
		}

		// Mark the last file (highest priority) with a special indicator
		label := relPath
		if i == 0 {
			label = relPath + " ⭐ (merged last)"
		}

		tree.AddNode(label)
	}

	return tree.String()
}

// findRootPath finds the common root directory for all files
func findRootPath(files []string) string {
	if len(files) == 0 {
		return "/"
	}
	if len(files) == 1 {
		return filepath.Dir(files[0])
	}

	// Start with first file's directory
	parts := strings.Split(filepath.Clean(filepath.Dir(files[0])), string(filepath.Separator))

	// Find common path parts with other files
	for _, file := range files[1:] {
		otherParts := strings.Split(filepath.Clean(filepath.Dir(file)), string(filepath.Separator))
		// Trim to common prefix length
		minLen := len(parts)
		if len(otherParts) < minLen {
			minLen = len(otherParts)
		}
		parts = parts[:minLen]

		// Check which parts match
		for i := 0; i < len(parts); i++ {
			if i >= len(otherParts) || parts[i] != otherParts[i] {
				parts = parts[:i]
				break
			}
		}
	}

	if len(parts) == 0 {
		return "/"
	}
	return filepath.Join(parts...)
}

// commonPath finds the common path between two paths
func commonPath(a, b string) string {
	aParts := strings.Split(filepath.Clean(a), string(filepath.Separator))
	bParts := strings.Split(filepath.Clean(b), string(filepath.Separator))

	var common []string
	for i := 0; i < len(aParts) && i < len(bParts); i++ {
		if aParts[i] == bParts[i] {
			common = append(common, aParts[i])
		} else {
			break
		}
	}

	if len(common) == 0 {
		return "/"
	}
	return filepath.Join(common...)
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

		// Get current directory to start discovery
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current working directory: %w", err)
		}

		// Perform discovery and analysis
		info, err := config.DiscoverAndAnalyze(cwd)
		if err != nil {
			return fmt.Errorf("failed to discover configuration: %w", err)
		}

		// If --path flag, just show the resolved storage path
		if pathFlag {
			fmt.Println(info.ResolvedPath)
			return nil
		}

		// Load the merged configuration
		merged, _, err := config.LoadConfigHierarchical(cwd)
		if err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}

		// If --json flag, output JSON
		if jsonFlag {
			jsonOutput := map[string]interface{}{
				"foundFiles":    info.FoundFiles,
				"mergingOrder":  info.MergingOrder,
				"stopReason":    info.StopReason,
				"resolvedPath":  info.ResolvedPath,
				"configuration": merged,
			}
			data, err := json.MarshalIndent(jsonOutput, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal JSON: %w", err)
			}
			fmt.Println(string(data))
			return nil
		}

		// Prepare template data
		stopDir := ""
		if len(info.MergingOrder) > 0 {
			stopDir = filepath.Dir(filepath.Join(info.MergingOrder[len(info.MergingOrder)-1]))
		}

		resolvedPath := info.ResolvedPath
		if resolvedPath == "" {
			resolvedPath = cwd
		}

		// Update merged config storage path to show resolved absolute path
		merged.Storage.Path = resolvedPath

		// Convert config to TOML string
		tomlStr, err := config.ConfigAsToml(merged)
		if err != nil {
			return fmt.Errorf("failed to format config as TOML: %w", err)
		}

		// Build file tree
		fileTree := buildConfigFileTree(info.FoundFiles)

		templateData := map[string]interface{}{
			"FoundFiles":                 info.FoundFiles,
			"MergingOrder":               info.MergingOrder,
			"ResolvedConfigAsTomlString": tomlStr,
			"StopReason":                 info.StopReason,
			"StopDir":                    stopDir,
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
		if verboseFlag && len(info.FoundFiles) > 0 {
			fmt.Println("\n=== Verbose Mode: Merging Details ===")
			for i, file := range info.MergingOrder {
				fmt.Printf("\n[Step %d] Applying: %s\n", i+1, file)
				cfg, err := config.LoadConfig(file)
				if err != nil {
					fmt.Printf("  Error loading: %v\n", err)
					continue
				}
				if cfg.Project.Name != "" {
					fmt.Printf("  - project.name: %q\n", cfg.Project.Name)
				}
				if len(cfg.Workflow.Statuses) > 0 {
					fmt.Printf("  - workflow.statuses: %v\n", cfg.Workflow.Statuses)
				}
				if cfg.Storage.Path != "" {
					fmt.Printf("  - storage.path: %q\n", cfg.Storage.Path)
				}
			}
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

		// Create config content
		configContent := fmt.Sprintf(`# opentask configuration for %s
[project]
name = %q
description = ""
owner = ""

[workflow]
statuses = ["todo", "in-progress", "reviewing", "done", "archived"]
initial = "todo"

# Transitions define allowed state changes
[[workflow.transitions]]
from = "todo"
to = ["in-progress", "archived"]

[[workflow.transitions]]
from = "in-progress"
to = ["reviewing", "todo", "archived"]

[[workflow.transitions]]
from = "reviewing"
to = ["done", "in-progress", "archived"]

[[workflow.transitions]]
from = "done"
to = ["archived"]

[storage]
backend = "markdown-fs"
path = %q

[templates]
# Leave empty to use built-in templates
# Or specify custom template paths
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
		fmt.Println("Next steps:")
		fmt.Println("  1. Create a task: opentask task new \"Your task title\"")
		fmt.Println("  2. List tasks: opentask task list")
		fmt.Println("  3. View config: opentask config view")
		fmt.Println("  4. Edit config: " + configPath)

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

func init() {
	configCmd.AddCommand(configInitCmd)
	configCmd.AddCommand(configViewCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)

	// Flags for config init
	configInitCmd.Flags().StringP("name", "n", "", "Project name (default: current directory name)")
	configInitCmd.Flags().StringP("storage", "s", "./.tasks/", "Storage path (default: ./.tasks/)")
	configInitCmd.Flags().BoolP("force", "f", false, "Overwrite existing .opentask.toml")

	// Flags for config view
	configViewCmd.Flags().BoolP("path", "p", false, "Show only the resolved storage path")
	configViewCmd.Flags().BoolP("json", "j", false, "Output resolved config as JSON")
	configViewCmd.Flags().BoolP("verbose", "v", false, "Show each config file contents during merging")
}
