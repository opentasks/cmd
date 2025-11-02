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

	// Load configuration hierarchically (searches up from current path)
	var cfg *config.ProjectConfig

	if configPath != "" {
		// Explicit config path provided - load only that file
		cfg, err = config.LoadConfig(configPath)
		if err != nil {
			return fmt.Errorf("failed to load configuration from %s: %w", configPath, err)
		}
	} else {
		// Discover and load configs hierarchically from the project path
		cfg, _, err = config.LoadConfigHierarchical(absPath)
		if err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}
	}

	// Ensure storage path is set
	if cfg.Storage.Path == "" {
		cfg.Storage.Path = absPath
	}

	// Initialize storage
	storageConfig := storage.StorageConfig{
		Backend: cfg.Storage.Backend,
		Path:    cfg.Storage.Path,
		Options: cfg.Storage.Options,
	}

	Store, err = storage.NewStorage(storageConfig)
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}

	// Initialize query engine
	Engine = query.NewQueryEngine(Store)

	if verbose {
		fmt.Fprintf(os.Stderr, "Initialized storage: %s at %s\n", cfg.Storage.Backend, cfg.Storage.Path)
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
