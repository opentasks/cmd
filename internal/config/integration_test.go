package config

import (
	"os"
	"path/filepath"
	"testing"
)

// TestIntegrationGlobalAndProjectConfigs tests the full end-to-end resolution with real files
func TestIntegrationGlobalAndProjectConfigs(t *testing.T) {
	tmpDir := t.TempDir()

	// Setup: Create global config
	globalDir := filepath.Join(tmpDir, ".config", "opentask")
	os.MkdirAll(globalDir, 0755)
	globalPath := filepath.Join(globalDir, "config.toml")

	globalConfig := `
active_project = "work"

[workflow]
statuses = ["todo", "done"]
initial = "todo"

[templates]
epic = "~/.local/share/opentask/templates/epic.md"

[[projects]]
id = "work"
name = "Work Tasks"

[projects.storage]
backend = "markdown-fs"
path = "` + filepath.Join(tmpDir, "work", ".tasks") + `"

[[projects]]
id = "personal"
name = "Personal Tasks"

[projects.storage]
backend = "markdown-fs"
path = "` + filepath.Join(tmpDir, "personal", ".tasks") + `"
`

	if err := os.WriteFile(globalPath, []byte(globalConfig), 0644); err != nil {
		t.Fatalf("Failed to write global config: %v", err)
	}

	// Setup: Create work project structure with nested configs
	workRoot := filepath.Join(tmpDir, "work")
	os.MkdirAll(workRoot, 0755)
	workProjectPath := filepath.Join(workRoot, ".opentask.toml")

	workProjectConfig := `
[project]
name = "Work"
owner = "team"

[storage]
backend = "markdown-fs"
path = "` + filepath.Join(workRoot, ".tasks") + `"

[workflow]
statuses = ["backlog", "todo", "in-progress", "done"]
initial = "todo"
`

	if err := os.WriteFile(workProjectPath, []byte(workProjectConfig), 0644); err != nil {
		t.Fatalf("Failed to write work project config: %v", err)
	}

	// Setup: Create nested sub-project
	subProjectDir := filepath.Join(workRoot, "sub-project")
	os.MkdirAll(subProjectDir, 0755)
	subProjectPath := filepath.Join(subProjectDir, ".opentask.toml")

	subProjectConfig := `
[project]
name = "Sub Project"
owner = "john"

[templates]
task = "templates/custom-task.md"
`

	if err := os.WriteFile(subProjectPath, []byte(subProjectConfig), 0644); err != nil {
		t.Fatalf("Failed to write sub-project config: %v", err)
	}

	// Set HOME to tmpDir so global config is found
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", oldHome)

	// Test Case 1: Resolve at work root - should get work project config + global defaults
	t.Run("resolve at work root", func(t *testing.T) {
		resolved, err := ResolveProjectConfig(workRoot)
		if err != nil {
			t.Fatalf("ResolveProjectConfig failed: %v", err)
		}

		if resolved.Project.Name != "Work" {
			t.Errorf("Project.Name = %q, want %q", resolved.Project.Name, "Work")
		}

		if resolved.Project.Owner != "team" {
			t.Errorf("Project.Owner = %q, want %q", resolved.Project.Owner, "team")
		}

		// Should have project-specific workflow
		if len(resolved.Workflow.Statuses) != 4 {
			t.Errorf("Workflow.Statuses = %v, want 4 statuses", resolved.Workflow.Statuses)
		}

		// Storage path should be absolute
		if !filepath.IsAbs(resolved.Storage.Path) {
			t.Errorf("Storage.Path = %q, expected absolute path", resolved.Storage.Path)
		}

		// Should have discovered work project config
		if len(resolved.DiscoveredFiles) == 0 {
			t.Fatal("No discovered files")
		}

		// Active project should be derived from global
		if resolved.ActiveProject != "work" {
			t.Errorf("ActiveProject = %q, want %q", resolved.ActiveProject, "work")
		}
	})

	// Test Case 2: Resolve at sub-project - should merge with parent + global
	t.Run("resolve at sub-project with parent override", func(t *testing.T) {
		resolved, err := ResolveProjectConfig(subProjectDir)
		if err != nil {
			t.Fatalf("ResolveProjectConfig failed: %v", err)
		}

		// Sub-project should override parent project name
		if resolved.Project.Name != "Sub Project" {
			t.Errorf("Project.Name = %q, want %q", resolved.Project.Name, "Sub Project")
		}

		// Owner should be from sub-project
		if resolved.Project.Owner != "john" {
			t.Errorf("Project.Owner = %q, want %q", resolved.Project.Owner, "john")
		}

		// Workflow should come from parent (work project config)
		if len(resolved.Workflow.Statuses) != 4 {
			t.Errorf("Workflow.Statuses = %v, want 4 from parent", resolved.Workflow.Statuses)
		}

		// Should have discovered both configs
		if len(resolved.DiscoveredFiles) < 2 {
			t.Errorf("Expected at least 2 discovered files, got %d", len(resolved.DiscoveredFiles))
		}

		// Storage should come from parent
		if !filepath.IsAbs(resolved.Storage.Path) {
			t.Errorf("Storage.Path = %q, expected absolute path", resolved.Storage.Path)
		}
	})
}

// TestIntegrationProjectWithoutGlobalConfig tests resolution when only project config exists
func TestIntegrationProjectWithoutGlobalConfig(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a simple project without global config
	projectDir := filepath.Join(tmpDir, "my-project")
	os.MkdirAll(projectDir, 0755)
	projectPath := filepath.Join(projectDir, ".opentask.toml")

	projectConfig := `
[project]
name = "My Project"
description = "Standalone project"

[storage]
backend = "markdown-fs"
path = "./.tasks"

[workflow]
statuses = ["todo", "in-progress", "done"]
initial = "todo"
`

	if err := os.WriteFile(projectPath, []byte(projectConfig), 0644); err != nil {
		t.Fatalf("Failed to write project config: %v", err)
	}

	// Set HOME to ensure no global config exists
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", oldHome)

	resolved, err := ResolveProjectConfig(projectDir)
	if err != nil {
		t.Fatalf("ResolveProjectConfig failed: %v", err)
	}

	if resolved.Project.Name != "My Project" {
		t.Errorf("Project.Name = %q, want %q", resolved.Project.Name, "My Project")
	}

	// Active project should be derived from directory name
	if resolved.ActiveProject != "my-project" {
		t.Errorf("ActiveProject = %q, want %q", resolved.ActiveProject, "my-project")
	}

	// Storage path should be absolute
	if !filepath.IsAbs(resolved.Storage.Path) {
		t.Errorf("Storage.Path = %q, expected absolute path", resolved.Storage.Path)
	}
}

// TestIntegrationMultipleProjectLevels tests hierarchical config resolution
func TestIntegrationMultipleProjectLevels(t *testing.T) {
	tmpDir := t.TempDir()

	// Level 1: Root project
	level1Dir := filepath.Join(tmpDir, "root-project")
	os.MkdirAll(level1Dir, 0755)
	level1Path := filepath.Join(level1Dir, ".opentask.toml")

	level1Config := `
[project]
name = "Root Project"

[workflow]
statuses = ["todo", "in-progress", "review", "done"]
initial = "todo"
`

	if err := os.WriteFile(level1Path, []byte(level1Config), 0644); err != nil {
		t.Fatalf("Failed to write level 1 config: %v", err)
	}

	// Level 2: Sub-project under root
	level2Dir := filepath.Join(level1Dir, "module-a")
	os.MkdirAll(level2Dir, 0755)
	level2Path := filepath.Join(level2Dir, ".opentask.toml")

	level2Config := `
[project]
name = "Module A"
owner = "team-a"
`

	if err := os.WriteFile(level2Path, []byte(level2Config), 0644); err != nil {
		t.Fatalf("Failed to write level 2 config: %v", err)
	}

	// Level 3: Nested project
	level3Dir := filepath.Join(level2Dir, "sub-module")
	os.MkdirAll(level3Dir, 0755)
	level3Path := filepath.Join(level3Dir, ".opentask.toml")

	level3Config := `
[project]
name = "Sub Module"
`

	if err := os.WriteFile(level3Path, []byte(level3Config), 0644); err != nil {
		t.Fatalf("Failed to write level 3 config: %v", err)
	}

	// Test resolution at level 3 - should merge all three levels
	resolved, err := ResolveProjectConfig(level3Dir)
	if err != nil {
		t.Fatalf("ResolveProjectConfig failed: %v", err)
	}

	// Should have name from level 3 (closest)
	if resolved.Project.Name != "Sub Module" {
		t.Errorf("Project.Name = %q, want %q", resolved.Project.Name, "Sub Module")
	}

	// Note: Level 3 defines [project] which replaces the entire section,
	// so owner is not inherited from level 2. This is correct behavior -
	// sections replace rather than merge. If we want inheritance, level 3 would need to
	// explicitly include owner or not define project section at all.
	// For this test, we'll verify level 3's config takes precedence
	if resolved.Project.Owner != "" {
		t.Errorf("Project.Owner = %q, want empty string (level 3 overrides)", resolved.Project.Owner)
	}

	// Should have workflow from level 1
	if len(resolved.Workflow.Statuses) != 4 {
		t.Errorf("Workflow.Statuses = %v, want 4 from level 1", resolved.Workflow.Statuses)
	}

	// Should have discovered all 3 config files + global if it exists
	if len(resolved.DiscoveredFiles) < 3 {
		t.Errorf("Expected at least 3 discovered files, got %d", len(resolved.DiscoveredFiles))
	}
}

// TestIntegrationStoragePathResolution tests correct path resolution
func TestIntegrationStoragePathResolution(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name            string
		configPath      string
		expectedAbsPath bool
	}{
		{
			name:            "relative path becomes absolute",
			configPath:      "./.tasks",
			expectedAbsPath: true,
		},
		{
			name:            "absolute path stays absolute",
			configPath:      filepath.Join(tmpDir, "custom", "tasks"),
			expectedAbsPath: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projectDir := filepath.Join(tmpDir, "proj-"+tt.name)
			os.MkdirAll(projectDir, 0755)
			projectPath := filepath.Join(projectDir, ".opentask.toml")

			projectConfig := `
[project]
name = "Test"

[storage]
backend = "markdown-fs"
path = "` + tt.configPath + `"
`

			if err := os.WriteFile(projectPath, []byte(projectConfig), 0644); err != nil {
				t.Fatalf("Failed to write project config: %v", err)
			}

			resolved, err := ResolveProjectConfig(projectDir)
			if err != nil {
				t.Fatalf("ResolveProjectConfig failed: %v", err)
			}

			if tt.expectedAbsPath && !filepath.IsAbs(resolved.Storage.Path) {
				t.Errorf("Storage.Path = %q, expected absolute path", resolved.Storage.Path)
			}
		})
	}
}

// TestIntegrationConfigInitWithTaskCreation tests the config init workflow:
// 1. Initialize a new project with config init
// 2. Resolve the config
// 3. Verify the storage path is correct
func TestIntegrationConfigInitWithTaskCreation(t *testing.T) {
	tmpDir := t.TempDir()

	// Setup: Create a git repository
	projectDir := filepath.Join(tmpDir, "test-project")
	gitDir := filepath.Join(projectDir, ".git")
	os.MkdirAll(gitDir, 0755)

	// Simulate config init by creating a project config file
	configPath := filepath.Join(projectDir, ".opentask.toml")
	configContent := `# opentask project configuration for test-project
# This file defines project-specific settings

# Project metadata
[project]
name = "test-project"
description = ""
owner = ""

# Storage configuration (project-specific)
[storage]
backend = "markdown-fs"
path = "./.tasks"
`

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	// Test 1: Verify config can be resolved
	resolved, err := ResolveProjectConfig(projectDir)
	if err != nil {
		t.Fatalf("ResolveProjectConfig failed: %v", err)
	}

	// Test 2: Verify storage path is absolute and correct
	expectedPath := filepath.Join(projectDir, ".tasks")
	if resolved.Storage.Path != expectedPath {
		t.Errorf("Storage.Path = %q, want %q", resolved.Storage.Path, expectedPath)
	}

	// Test 3: Verify discovered files include only local config and global config
	if len(resolved.DiscoveredFiles) < 1 {
		t.Errorf("DiscoveredFiles is empty, should have at least local config")
	}

	// First file should be the local config
	if resolved.DiscoveredFiles[0] != configPath {
		t.Errorf("First discovered file = %q, want %q", resolved.DiscoveredFiles[0], configPath)
	}

	// Test 4: Verify project metadata
	if resolved.Project.Name != "test-project" {
		t.Errorf("Project.Name = %q, want 'test-project'", resolved.Project.Name)
	}
}

// TestIntegrationConfigStopsAtGitRoot tests that config discovery stops at git root
// When .git is found at a directory level, that directory's config is included,
// but we don't walk up further past the git root.
func TestIntegrationConfigStopsAtGitRoot(t *testing.T) {
	tmpDir := t.TempDir()

	// Create structure:
	// tmpDir/
	//   .opentask.toml (root - SHOULD be discovered and then stop)
	//   .git/
	//   subdir/
	//     .opentask.toml (should be discovered - below git root)
	//     deep/
	//       (no config here - discovery starts here)

	rootConfig := filepath.Join(tmpDir, ".opentask.toml")
	gitDir := filepath.Join(tmpDir, ".git")
	subdir := filepath.Join(tmpDir, "subdir")
	subConfig := filepath.Join(subdir, ".opentask.toml")
	deepdir := filepath.Join(subdir, "deep")

	// Create directories
	os.MkdirAll(gitDir, 0755)
	os.MkdirAll(deepdir, 0755)

	// Create config files
	for _, path := range []string{rootConfig, subConfig} {
		content := `[project]
name = "test"

[storage]
path = "./.tasks"
`
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write config: %v", err)
		}
	}

	// Resolve config from deepdir
	resolved, err := ResolveProjectConfig(deepdir)
	if err != nil {
		t.Fatalf("ResolveProjectConfig failed: %v", err)
	}

	// Filter to project configs (not global config)
	home, _ := os.UserHomeDir()
	globalPath := filepath.Join(home, ".config", "opentask", "config.toml")
	var projectConfigs []string
	for _, f := range resolved.DiscoveredFiles {
		if f != globalPath {
			projectConfigs = append(projectConfigs, f)
		}
	}

	// Should find exactly 2 configs (subConfig and rootConfig)
	// Because discovery walks from deepdir -> subdir -> tmpdir (finds git) and includes tmpdir's config
	if len(projectConfigs) != 2 {
		t.Errorf("Found %d project configs, want 2", len(projectConfigs))
		for _, f := range projectConfigs {
			t.Logf("  - %s", f)
		}
	}

	// Verify both configs are found
	hasSubConfig := false
	hasRootConfig := false
	for _, f := range projectConfigs {
		if f == subConfig {
			hasSubConfig = true
		}
		if f == rootConfig {
			hasRootConfig = true
		}
	}
	if !hasSubConfig {
		t.Errorf("subConfig not found in discovered files")
	}
	if !hasRootConfig {
		t.Errorf("rootConfig not found in discovered files (should find config at git root)")
	}
}

// TestIntegrationConfigDoesntWalkPastGitRoot tests that discovery doesn't walk PAST git root
func TestIntegrationConfigDoesntWalkPastGitRoot(t *testing.T) {
	tmpDir := t.TempDir()

	// Create structure:
	// tmpDir/
	//   .opentask.toml (above git root - should NOT be discovered)
	//   subdir/
	//     .git/
	//     .opentask.toml (at git root - SHOULD be discovered)
	//     deep/
	//       (no config - discovery starts here)

	rootConfig := filepath.Join(tmpDir, ".opentask.toml")
	subdir := filepath.Join(tmpDir, "subdir")
	gitDir := filepath.Join(subdir, ".git")
	subConfig := filepath.Join(subdir, ".opentask.toml")
	deepdir := filepath.Join(subdir, "deep")

	// Create directories
	os.MkdirAll(gitDir, 0755)
	os.MkdirAll(deepdir, 0755)

	// Create config files
	for _, path := range []string{rootConfig, subConfig} {
		content := `[project]
name = "test"

[storage]
path = "./.tasks"
`
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write config: %v", err)
		}
	}

	// Resolve config from deepdir
	resolved, err := ResolveProjectConfig(deepdir)
	if err != nil {
		t.Fatalf("ResolveProjectConfig failed: %v", err)
	}

	// Filter to project configs (not global config)
	home, _ := os.UserHomeDir()
	globalPath := filepath.Join(home, ".config", "opentask", "config.toml")
	var projectConfigs []string
	for _, f := range resolved.DiscoveredFiles {
		if f != globalPath {
			projectConfigs = append(projectConfigs, f)
		}
	}

	// Should find only 1 config (subConfig) - NOT rootConfig because discovery stops at git root
	if len(projectConfigs) != 1 {
		t.Errorf("Found %d project configs, want 1 (should not walk past git root)", len(projectConfigs))
		for _, f := range projectConfigs {
			t.Logf("  - %s", f)
		}
	}

	// Verify only subConfig is found
	hasSubConfig := false
	hasRootConfig := false
	for _, f := range projectConfigs {
		if f == subConfig {
			hasSubConfig = true
		}
		if f == rootConfig {
			hasRootConfig = true
		}
	}
	if !hasSubConfig {
		t.Errorf("subConfig not found in discovered files")
	}
	if hasRootConfig {
		t.Errorf("rootConfig should NOT be discovered (discovery should stop at git root)")
	}
}
