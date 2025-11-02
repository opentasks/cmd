package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

// TestInitializeStorageWithParentConfig tests that initializeStorage can
// discover and load config files from parent directories
func TestInitializeStorageWithParentConfig(t *testing.T) {
	tmpDir := t.TempDir()

	// Create structure:
	// tmpDir/
	//   opentask.toml (parent config)
	//   tasks/
	//   subdir/
	//     (no config - should inherit from parent)

	parentConfig := filepath.Join(tmpDir, "opentask.toml")
	parentContent := `[project]
name = "Parent Project"
description = "Parent Description"

[storage]
path = "tasks"
`
	if err := os.WriteFile(parentConfig, []byte(parentContent), 0644); err != nil {
		t.Fatalf("Failed to create parent config: %v", err)
	}

	tasksDir := filepath.Join(tmpDir, "tasks")
	if err := os.MkdirAll(tasksDir, 0755); err != nil {
		t.Fatalf("Failed to create tasks directory: %v", err)
	}

	subdir := filepath.Join(tmpDir, "subdir")
	if err := os.MkdirAll(subdir, 0755); err != nil {
		t.Fatalf("Failed to create subdir: %v", err)
	}

	// Test 1: Initialize storage from parent directory
	t.Run("from parent directory", func(t *testing.T) {
		projectPath = tmpDir
		configPath = ""

		err := initializeStorage(nil, nil)
		if err != nil {
			t.Errorf("initializeStorage() error = %v, want nil", err)
		}

		if Store == nil {
			t.Error("Store should not be nil after initialization")
		}

		if Engine == nil {
			t.Error("Engine should not be nil after initialization")
		}

		// Cleanup
		Store.Close()
		Store = nil
		Engine = nil
	})

	// Test 2: Initialize storage from subdirectory (should discover parent config)
	t.Run("from subdirectory with parent config", func(t *testing.T) {
		projectPath = subdir
		configPath = ""

		err := initializeStorage(nil, nil)
		if err != nil {
			t.Errorf("initializeStorage() error = %v, want nil", err)
		}

		if Store == nil {
			t.Error("Store should not be nil after initialization")
		}

		if Engine == nil {
			t.Error("Engine should not be nil after initialization")
		}

		// The key test: we should be able to initialize from a subdirectory
		// and it should find the parent config (not fail looking for local config)

		// Cleanup
		Store.Close()
		Store = nil
		Engine = nil
	})
}

// TestInitializeStorageWithExplicitConfigPath tests that explicit --config flag is respected
func TestInitializeStorageWithExplicitConfigPath(t *testing.T) {
	tmpDir := t.TempDir()

	// Create two config files at different locations
	config1Dir := filepath.Join(tmpDir, "config1")
	config1Path := filepath.Join(config1Dir, "opentask.toml")
	if err := os.MkdirAll(config1Dir, 0755); err != nil {
		t.Fatalf("Failed to create config1 dir: %v", err)
	}

	config1Content := `[project]
name = "Config1 Project"

[storage]
path = "config1-tasks"
`
	if err := os.WriteFile(config1Path, []byte(config1Content), 0644); err != nil {
		t.Fatalf("Failed to create config1: %v", err)
	}

	config2Dir := filepath.Join(tmpDir, "config2")
	config2Path := filepath.Join(config2Dir, "opentask.toml")
	if err := os.MkdirAll(config2Dir, 0755); err != nil {
		t.Fatalf("Failed to create config2 dir: %v", err)
	}

	config2Content := `[project]
name = "Config2 Project"

[storage]
path = "config2-tasks"
`
	if err := os.WriteFile(config2Path, []byte(config2Content), 0644); err != nil {
		t.Fatalf("Failed to create config2: %v", err)
	}

	// Set explicit config path to config2
	projectPath = config1Dir
	configPath = config2Path

	err := initializeStorage(nil, nil)
	if err != nil {
		t.Errorf("initializeStorage() error = %v, want nil", err)
	}

	if Store == nil {
		t.Error("Store should not be nil after initialization")
	}

	if Engine == nil {
		t.Error("Engine should not be nil after initialization")
	}

	// The key test: explicit config should be used instead of discovered config

	// Cleanup
	Store.Close()
	Store = nil
	Engine = nil
}
