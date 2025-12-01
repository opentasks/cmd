package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/opentasks/cmd/internal/config"
	"github.com/spf13/cobra"
)

// configCmd represents the config command group
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
	Long:  "Commands for viewing and managing configuration",
}

// configViewCmd views the resolved configuration with discovery details
var configViewCmd = &cobra.Command{
	Use:     "view",
	Short:   "View resolved configuration",
	Long:    "Display the resolved project configuration after discovering and merging all config files",
	Aliases: []string{"show"},
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

		return config.Renderer(
			config.WithCwd(cwd),
			config.WithVerbose(verboseFlag),
			config.WithResolvedConfig(resolved),
		)
	},
}

// configInitCmd initializes a new opentask project
var configInitCmd = &cobra.Command{
	Use:     "init",
	Short:   "Initialize a new project",
	Long:    "Create a new .opentask.toml configuration file in the current directory",
	PreRunE: allowUnresolved, // Allow init to work without existing project
	RunE: func(cmd *cobra.Command, args []string) error {
		nameFlag, _ := cmd.Flags().GetString("name")
		storagePath, _ := cmd.Flags().GetString("storage")
		forceFlag, _ := cmd.Flags().GetBool("force")

		// Get current working directory
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current working directory: %w", err)
		}

		initializer := config.NewConfigInitializer(cwd)
		return initializer.Initialize(nameFlag, storagePath, forceFlag)
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

		if globalConfig == nil {
			fmt.Println("No global configuration found at " + globalPath)
			fmt.Println("\nTo create global config:")
			fmt.Println("  1. mkdir -p ~/.config/opentask")
			fmt.Println("  2. Create config.toml with project definitions")
			return nil
		}

		// List projects
		if len(globalConfig.Projects) == 0 {
			fmt.Println("No projects configured in global config")
			return nil
		}

		fmt.Println("Configured projects:")
		fmt.Println()

		lister := config.NewProjectLister(globalConfig)
		fmt.Print(lister.List())
		fmt.Println()

		// Show which project would be active for current directory
		cwd, err := os.Getwd()
		if err == nil {
			resolved, err := config.ResolveProjectConfig(cwd)
			if err == nil && resolved.IsResolved {
				fmt.Printf("\nNote: Current directory resolves to project: %s\n", resolved.ActiveProject)
			}
		}

		return nil
	},
}

func init() {
	configCmd.AddCommand(configInitCmd)
	configCmd.AddCommand(configViewCmd)
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
}
