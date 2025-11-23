package project

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Manager provides pure business logic for project operations
type Manager struct{}

// NewManager creates a new project manager
func NewManager() *Manager {
	return &Manager{}
}

// ResolvePath resolves a path to an absolute path, expanding ~ if needed
func (m *Manager) ResolvePath(path string) (string, error) {
	// Expand ~ to home directory
	if len(path) > 0 && path[0] == '~' {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to expand home directory: %w", err)
		}
		path = filepath.Join(home, path[1:])
	}

	// Convert to absolute
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}

	return filepath.Clean(absPath), nil
}

// GlobalConfigPath returns the path to the global config file
func (m *Manager) GlobalConfigPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get config directory: %w", err)
	}

	opentaskDir := filepath.Join(configDir, "opentask")
	configPath := filepath.Join(opentaskDir, "config.toml")

	// Create directory if it doesn't exist
	if err := os.MkdirAll(opentaskDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create config directory: %w", err)
	}

	return configPath, nil
}

// ValidatePath validates that a path exists
func (m *Manager) ValidatePath(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("path does not exist: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("path is not a directory: %s", path)
	}
	return nil
}

// FormatPathForDisplay formats a path for display, using ~ for home directory if applicable
func (m *Manager) FormatPathForDisplay(path string) string {
	home, err := os.UserHomeDir()
	if err == nil && strings.HasPrefix(path, home) {
		return "~" + path[len(home):]
	}
	return path
}
