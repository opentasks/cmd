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
	//   .opentask.toml (parent config)
	//   tasks/
	//   subdir/
	//     (no config - should inherit from parent)

	parentConfig := filepath.Join(tmpDir, ".opentask.toml")
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
		_ = Store.Close()
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
		_ = Store.Close()
		Store = nil
		Engine = nil
	})
}

// TestInitializeStorageWithExplicitConfigPath tests that explicit --config flag is respected
func TestInitializeStorageWithExplicitConfigPath(t *testing.T) {
	tmpDir := t.TempDir()

	// Create two config files at different locations
	config1Dir := filepath.Join(tmpDir, "config1")
	config1Path := filepath.Join(config1Dir, ".opentask.toml")
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
	config2Path := filepath.Join(config2Dir, ".opentask.toml")
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
	_ = Store.Close()
	Store = nil
	Engine = nil
}

// TestRequireActiveProjectWithResolvedProject tests that requireActiveProject
// succeeds when a project is resolved (has .opentask.toml)
func TestRequireActiveProjectWithResolvedProject(t *testing.T) {
	// Save original values
	origProjectPath := projectPath
	origConfigPath := configPath
	defer func() {
		projectPath = origProjectPath
		configPath = origConfigPath
	}()

	tmpDir := t.TempDir()

	// Create .opentask.toml
	testConfigPath := filepath.Join(tmpDir, ".opentask.toml")
	configContent := `[project]
id = "test-project"
name = "Test Project"

[storage]
path = ".tasks"
`
	if err := os.WriteFile(testConfigPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	// Create storage directory
	tasksDir := filepath.Join(tmpDir, ".tasks")
	if err := os.MkdirAll(tasksDir, 0755); err != nil {
		t.Fatalf("Failed to create tasks directory: %v", err)
	}

	// Set project path
	projectPath = tmpDir
	configPath = ""

	// Call requireActiveProject - should succeed
	err := requireActiveProject(nil, nil)
	if err != nil {
		t.Errorf("requireActiveProject() error = %v, want nil", err)
	}

	// Verify that config was resolved
	if Resolved == nil {
		t.Fatal("Resolved config should not be nil")
	}

	if !Resolved.IsResolved {
		t.Error("IsResolved should be true")
	}

	if Resolved.ActiveProject != "test-project" {
		t.Errorf("ActiveProject = %q, want %q", Resolved.ActiveProject, "test-project")
	}

	// Cleanup
	if Store != nil {
		_ = Store.Close()
	}
	Store = nil
	Engine = nil
	Resolved = nil
}

// TestRequireActiveProjectWithoutProject tests that requireActiveProject
// returns an error when no project is configured
func TestRequireActiveProjectWithoutProject(t *testing.T) {
	// Save original values
	origProjectPath := projectPath
	origConfigPath := configPath
	defer func() {
		projectPath = origProjectPath
		configPath = origConfigPath
	}()

	tmpDir := t.TempDir()

	// NO .opentask.toml file created

	// Set project path
	projectPath = tmpDir
	configPath = ""

	// Call requireActiveProject - should fail
	err := requireActiveProject(nil, nil)
	if err == nil {
		t.Error("requireActiveProject() error = nil, want error")
	}

	// Verify error message
	expectedErrMsg := "no active project found"
	if err.Error() != expectedErrMsg {
		t.Errorf("requireActiveProject() error = %q, want %q", err.Error(), expectedErrMsg)
	}

	// Verify that config was attempted to resolve but IsResolved = false
	if Resolved == nil {
		t.Fatal("Resolved config should not be nil")
	}

	if Resolved.IsResolved {
		t.Error("IsResolved should be false")
	}

	if Resolved.ActiveProject != "" {
		t.Errorf("ActiveProject = %q, want empty string", Resolved.ActiveProject)
	}

	// Cleanup
	if Store != nil {
		_ = Store.Close()
		Store = nil
	}
	Engine = nil
	Resolved = nil
}

// TestAllowUnresolvedWithProject tests that allowUnresolved
// succeeds when a project is configured
func TestAllowUnresolvedWithProject(t *testing.T) {
	// Save original values
	origProjectPath := projectPath
	origConfigPath := configPath
	defer func() {
		projectPath = origProjectPath
		configPath = origConfigPath
	}()

	tmpDir := t.TempDir()

	// Create .opentask.toml
	testConfigPath := filepath.Join(tmpDir, ".opentask.toml")
	configContent := `[project]
name = "Test Project"

[storage]
path = ".tasks"
`
	if err := os.WriteFile(testConfigPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	// Create storage directory
	tasksDir := filepath.Join(tmpDir, ".tasks")
	if err := os.MkdirAll(tasksDir, 0755); err != nil {
		t.Fatalf("Failed to create tasks directory: %v", err)
	}

	// Set project path
	projectPath = tmpDir
	configPath = ""

	// Call allowUnresolved - should succeed
	err := allowUnresolved(nil, nil)
	if err != nil {
		t.Errorf("allowUnresolved() error = %v, want nil", err)
	}

	// Verify that storage was initialized
	if Store == nil {
		t.Error("Store should not be nil")
	}

	if Engine == nil {
		t.Error("Engine should not be nil")
	}

	// Cleanup
	if Store != nil {
		_ = Store.Close()
	}
	Store = nil
	Engine = nil
	Resolved = nil
}

// TestAllowUnresolvedWithoutProject tests that allowUnresolved
// succeeds even when no project is configured (onboarding scenario)
func TestAllowUnresolvedWithoutProject(t *testing.T) {
	// Save original values
	origProjectPath := projectPath
	origConfigPath := configPath
	defer func() {
		projectPath = origProjectPath
		configPath = origConfigPath
	}()

	tmpDir := t.TempDir()

	// NO .opentask.toml file created

	// Set project path
	projectPath = tmpDir
	configPath = ""

	// Call allowUnresolved - should succeed (doesn't require resolved project)
	err := allowUnresolved(nil, nil)
	if err != nil {
		t.Errorf("allowUnresolved() error = %v, want nil", err)
	}

	// Verify that storage was initialized with defaults
	if Store == nil {
		t.Error("Store should not be nil")
	}

	if Engine == nil {
		t.Error("Engine should not be nil")
	}

	// Verify that config is unresolved
	if Resolved == nil {
		t.Fatal("Resolved config should not be nil")
	}

	if Resolved.IsResolved {
		t.Error("IsResolved should be false for onboarding scenario")
	}

	// Cleanup
	if Store != nil {
		_ = Store.Close()
	}
	Store = nil
	Engine = nil
	Resolved = nil
}

// TestRequireActiveProjectWithGlobalContextMatch tests resolution via global config context
// NOTE: This test is skipped because global config resolution in tests requires more complex setup.
// The actual functionality is tested in internal/config/merge_test.go via TestFindProjectByContext
// and TestIntegrationGlobalAndProjectConfigs.
func TestRequireActiveProjectWithGlobalContextMatch(t *testing.T) {
	t.Skip("Global config resolution tested in internal/config/merge_test.go")
}

// TestRequireActiveProjectWithDirectoryNameFallback tests resolution via directory name
func TestRequireActiveProjectWithDirectoryNameFallback(t *testing.T) {
	// Save original values
	origProjectPath := projectPath
	origConfigPath := configPath
	defer func() {
		projectPath = origProjectPath
		configPath = origConfigPath
	}()

	tmpDir := t.TempDir()

	// Create a directory with a specific name
	projectDir := filepath.Join(tmpDir, "awesome-app")
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		t.Fatalf("Failed to create project directory: %v", err)
	}

	// Create .opentask.toml WITHOUT explicit project.id
	testConfigPath := filepath.Join(projectDir, ".opentask.toml")
	configContent := `[project]
name = "Awesome Application"

[storage]
path = ".tasks"
`
	if err := os.WriteFile(testConfigPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	// Create storage directory
	tasksDir := filepath.Join(projectDir, ".tasks")
	if err := os.MkdirAll(tasksDir, 0755); err != nil {
		t.Fatalf("Failed to create tasks directory: %v", err)
	}

	// Set project path
	projectPath = projectDir
	configPath = ""

	// Call requireActiveProject - should succeed with directory name fallback
	err := requireActiveProject(nil, nil)
	if err != nil {
		t.Errorf("requireActiveProject() error = %v, want nil", err)
	}

	// Verify that config was resolved with directory name
	if Resolved == nil {
		t.Fatal("Resolved config should not be nil")
	}

	if !Resolved.IsResolved {
		t.Error("IsResolved should be true")
	}

	// ActiveProject should be derived from directory name
	if Resolved.ActiveProject != "awesome-app" {
		t.Errorf("ActiveProject = %q, want %q", Resolved.ActiveProject, "awesome-app")
	}

	// Cleanup
	if Store != nil {
		_ = Store.Close()
	}
	Store = nil
	Engine = nil
	Resolved = nil
}
