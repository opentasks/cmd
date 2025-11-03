package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zenobi-us/opentask/internal/config"
	"github.com/zenobi-us/opentask/internal/query"
	"github.com/zenobi-us/opentask/internal/storage"
)

var (
	// Global flags
	projectPath string
	configPath  string
	verbose     bool

	// Global state
	Engine *query.QueryEngine
	Store  storage.BaseStorage
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "opentask",
	Short: "opentask - A simple task management system",
	Long: `opentask is a task management system that stores tasks as markdown files
with YAML frontmatter. It provides a command-line interface for managing tasks,
organizing them hierarchically, and tracking relationships.`,
	PersistentPreRunE:  initializeStorage,
	PersistentPostRunE: cleanupStorage,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

// initializeStorage initializes the storage backend and query engine
func initializeStorage(cmd *cobra.Command, args []string) error {
	// Determine project path
	path := projectPath
	if path == "" {
		// Try environment variable
		path = os.Getenv("opentask_PROJECT_PATH")
	}
	if path == "" {
		// Use current directory
		path = "."
	}

	// Resolve to absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to resolve project path: %w", err)
	}

	// Resolve configuration (searches hierarchically from current path)
	var resolved *config.OpentaskResolvedConfig

	if configPath != "" {
		// Explicit config path provided - load only that file
		projectConfig, err := config.LoadProjectConfig(configPath)
		if err != nil {
			return fmt.Errorf("failed to load configuration from %s: %w", configPath, err)
		}
		if projectConfig == nil {
			return fmt.Errorf("configuration file not found: %s", configPath)
		}
		// Convert to resolved config
		resolved = config.NewResolvedConfig()
		resolved = config.MergeProjectConfig(resolved, projectConfig)
	} else {
		// Discover and resolve configs hierarchically from the project path
		resolved, err = config.ResolveProjectConfig(absPath)
		if err != nil {
			return fmt.Errorf("failed to resolve configuration: %w", err)
		}
	}

	// Ensure storage path is set
	if resolved.Storage.Path == "" {
		resolved.Storage.Path = absPath
	}

	// Initialize storage
	storageConfig := storage.StorageConfig{
		Backend: resolved.Storage.Backend,
		Path:    resolved.Storage.Path,
		Options: resolved.Storage.Options,
	}

	Store, err = storage.NewStorage(storageConfig)
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}

	// Initialize query engine
	Engine = query.NewQueryEngine(Store)

	if verbose {
		fmt.Fprintf(os.Stderr, "Initialized storage: %s at %s\n", resolved.Storage.Backend, resolved.Storage.Path)
		if len(resolved.DiscoveredFiles) > 0 {
			fmt.Fprintf(os.Stderr, "Configuration files:\n")
			for _, f := range resolved.DiscoveredFiles {
				fmt.Fprintf(os.Stderr, "  - %s\n", f)
			}
		}
	}

	return nil
}

// cleanupStorage performs cleanup after command execution
func cleanupStorage(cmd *cobra.Command, args []string) error {
	if Store != nil {
		return Store.Close()
	}
	return nil
}

// GetContext returns a context for the command
func GetContext() context.Context {
	return context.Background()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&projectPath, "path", "", "Project path (default: current directory)")
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "", "Config file path")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Enable verbose output")

	// Bind flags to viper
	viper.BindPFlag("project.path", rootCmd.PersistentFlags().Lookup("path"))
	viper.BindPFlag("config.path", rootCmd.PersistentFlags().Lookup("config"))

	// Add subcommands
	rootCmd.AddCommand(taskCmd)
	rootCmd.AddCommand(projectCmd)
	rootCmd.AddCommand(configCmd)
}
