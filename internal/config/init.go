package config

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

//go:embed config.tmpl
var configTemplate string

// ConfigInitializer handles initialization of new opentask projects
type ConfigInitializer struct {
	cwd string
}

// NewConfigInitializer creates a new ConfigInitializer for the given working directory
func NewConfigInitializer(cwd string) *ConfigInitializer {
	return &ConfigInitializer{
		cwd: cwd,
	}
}

// Initialize creates a new .opentask.toml configuration file
// If force is true, it will overwrite an existing config file
func (ci *ConfigInitializer) Initialize(name, storagePath string, force bool) error {
	configPath := filepath.Join(ci.cwd, ".opentask.toml")

	// Check if config already exists
	if _, err := os.Stat(configPath); err == nil && !force {
		return fmt.Errorf(".opentask.toml already exists in %s. Use --force to overwrite", ci.cwd)
	}

	// Determine project name
	projectName := name
	if projectName == "" {
		// Use directory name as default
		projectName = filepath.Base(ci.cwd)
	}

	// Render template with values
	configContent, err := ci.renderTemplate(map[string]string{
		"ProjectName": projectName,
		"StoragePath": storagePath,
	})
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	// Write config file
	// #nosec G306 - File permissions 0600 are intentional for config files (sensitive data)
	// #nosec G304 - configPath constructed from user-provided path, validated by caller
	if err := os.WriteFile(configPath, []byte(configContent), 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	// Print success message
	ci.printSummary(configPath, projectName, storagePath)

	return nil
}

// renderTemplate renders the embedded template with the given data
func (ci *ConfigInitializer) renderTemplate(data map[string]string) (string, error) {
	tmpl, err := template.New("config").Funcs(template.FuncMap{
		"quote": func(s string) string {
			return `"` + s + `"`
		},
	}).Parse(configTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// printSummary prints a formatted summary of the initialization
func (ci *ConfigInitializer) printSummary(configPath, projectName, storagePath string) {
	fmt.Printf("Initialized opentask project in %s\n", ci.cwd)
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
}
