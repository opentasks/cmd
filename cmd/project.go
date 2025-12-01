package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/opentasks/cmd/internal/config"
	"github.com/opentasks/cmd/internal/project"
	"github.com/spf13/cobra"
)

// projectCmd represents the project command group
var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Manage projects",
	Long:  "Commands for managing projects and their working directory contexts",
}

// projectAttachCmd attaches a working directory to a project
var projectAttachCmd = &cobra.Command{
	Use:   "attach [PATH]",
	Short: "Attach a working directory to a project",
	Long: `Attach a working directory to a project so that tasks are automatically
found in the project's storage when working from that directory.

If PATH is not provided, the current working directory is used.

Example:
  opentask project attach --project opentask
  opentask project attach /path/to/worktree --project opentask`,
	RunE: func(cmd *cobra.Command, args []string) error {
		projectID, _ := cmd.Flags().GetString("project")
		if projectID == "" {
			return fmt.Errorf("--project flag is required")
		}

		// Determine path
		var targetPath string
		if len(args) > 0 {
			targetPath = args[0]
		} else {
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current working directory: %w", err)
			}
			targetPath = cwd
		}

		// Resolve path
		pm := project.NewManager()
		resolvedPath, err := pm.ResolvePath(targetPath)
		if err != nil {
			return err
		}

		// Load global config
		globalPath, err := pm.GlobalConfigPath()
		if err != nil {
			return err
		}

		globalConfig, err := config.LoadGlobalConfig(globalPath)
		if err != nil {
			return fmt.Errorf("failed to load global config: %w", err)
		}

		if globalConfig == nil {
			globalConfig = &config.OpentaskGlobalConfigFile{
				Projects: []config.GlobalProjectConfig{},
			}
		}

		// Use service to attach context
		svc := config.NewProjectService(globalConfig)
		if err := svc.AttachContext(projectID, resolvedPath); err != nil {
			return fmt.Errorf("failed to add context path: %w", err)
		}

		// Save global config
		saver := config.NewGlobalConfigSaver()
		if err := saver.Save(globalPath, globalConfig); err != nil {
			return err
		}

		fmt.Printf("✓ Attached %s to project %s\n", resolvedPath, projectID)
		return nil
	},
}

// projectDetachCmd detaches a working directory from a project
var projectDetachCmd = &cobra.Command{
	Use:   "detach [PATH]",
	Short: "Detach a working directory from a project",
	Long: `Remove a working directory from a project.

If PATH is not provided, the current working directory is used.

Example:
  opentask project detach --project opentask
  opentask project detach /path/to/worktree --project opentask`,
	RunE: func(cmd *cobra.Command, args []string) error {
		projectID, _ := cmd.Flags().GetString("project")
		if projectID == "" {
			return fmt.Errorf("--project flag is required")
		}

		// Determine path
		var targetPath string
		if len(args) > 0 {
			targetPath = args[0]
		} else {
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current working directory: %w", err)
			}
			targetPath = cwd
		}

		// Resolve path
		pm := project.NewManager()
		resolvedPath, err := pm.ResolvePath(targetPath)
		if err != nil {
			return err
		}

		// Load global config
		globalPath, err := pm.GlobalConfigPath()
		if err != nil {
			return err
		}

		globalConfig, err := config.LoadGlobalConfig(globalPath)
		if err != nil {
			return fmt.Errorf("failed to load global config: %w", err)
		}

		if globalConfig == nil {
			return fmt.Errorf("global config not found at %s", globalPath)
		}

		// Use service to detach context
		svc := config.NewProjectService(globalConfig)
		if err := svc.DetachContext(projectID, resolvedPath); err != nil {
			return fmt.Errorf("failed to remove context path: %w", err)
		}

		// Save global config
		saver := config.NewGlobalConfigSaver()
		if err := saver.Save(globalPath, globalConfig); err != nil {
			return err
		}

		fmt.Printf("✓ Detached %s from project %s\n", resolvedPath, projectID)
		return nil
	},
}

// projectListCmd lists all projects and their contexts
var projectListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all projects",
	Long:  "Display all projects with their storage paths and working directory contexts",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load global config
		pm := project.NewManager()
		globalPath, err := pm.GlobalConfigPath()
		if err != nil {
			return err
		}

		globalConfig, err := config.LoadGlobalConfig(globalPath)
		if err != nil {
			return fmt.Errorf("failed to load global config: %w", err)
		}

		if globalConfig == nil || len(globalConfig.Projects) == 0 {
			fmt.Println("No projects configured")
			return nil
		}

		fmt.Println("Projects:")
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

// projectRemoveCmd removes a project from the global config
var projectRemoveCmd = &cobra.Command{
	Use:   "remove [PROJECT_ID]",
	Short: "Remove a project",
	Long: `Remove a project from the global configuration.

This will permanently delete the project configuration, though the storage
directory itself will not be affected.

Example:
  opentask project remove myproject`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectID := args[0]

		// Load global config
		pm := project.NewManager()
		globalPath, err := pm.GlobalConfigPath()
		if err != nil {
			return err
		}

		globalConfig, err := config.LoadGlobalConfig(globalPath)
		if err != nil {
			return fmt.Errorf("failed to load global config: %w", err)
		}

		if globalConfig == nil {
			return fmt.Errorf("global config not found at %s", globalPath)
		}

		// Use service to get project for confirmation
		svc := config.NewProjectService(globalConfig)
		proj := svc.GetProject(projectID)
		if proj == nil {
			return fmt.Errorf("project not found: %s", projectID)
		}

		// Display project info and ask for confirmation
		fmt.Printf("About to remove project: %s\n", projectID)
		if proj.Name != "" {
			fmt.Printf("Name: %s\n", proj.Name)
		}
		if proj.Storage != nil && proj.Storage.Path != "" {
			fmt.Printf("Storage path: %s\n", pm.FormatPathForDisplay(proj.Storage.Path))
		}
		fmt.Println()
		fmt.Print("Are you sure you want to remove this project? (yes/no): ")

		// Read user input
		scanner := bufio.NewScanner(os.Stdin)
		if !scanner.Scan() {
			return fmt.Errorf("failed to read confirmation")
		}

		response := strings.TrimSpace(strings.ToLower(scanner.Text()))
		if response != "yes" && response != "y" {
			fmt.Println("Cancelled.")
			return nil
		}

		// Use service to remove project
		if err := svc.RemoveProject(projectID); err != nil {
			return err
		}

		// Save global config
		saver := config.NewGlobalConfigSaver()
		if err := saver.Save(globalPath, globalConfig); err != nil {
			return err
		}

		fmt.Printf("✓ Removed project %s\n", projectID)
		return nil
	},
}

func init() {
	// Add subcommands
	projectCmd.AddCommand(projectAttachCmd)
	projectCmd.AddCommand(projectDetachCmd)
	projectCmd.AddCommand(projectListCmd)
	projectCmd.AddCommand(projectRemoveCmd)

	// Add flags for attach and detach
	projectAttachCmd.Flags().StringP("project", "p", "", "Project ID (required)")
	projectDetachCmd.Flags().StringP("project", "p", "", "Project ID (required)")
}
