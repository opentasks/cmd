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
[global]
active_project = "work"

[workflow]
statuses = ["todo", "done"]
initial = "todo"

[templates]
epic = "~/.local/share/opentask/templates/epic.md"

[[global.projects]]
id = "work"
name = "Work Tasks"

[global.projects.storage]
backend = "markdown-fs"
path = "` + filepath.Join(tmpDir, "work", ".tasks") + `"

[[global.projects]]
id = "personal"
name = "Personal Tasks"

[global.projects.storage]
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
[project.project]
name = "Work"
owner = "team"

[project.storage]
backend = "markdown-fs"
path = "` + filepath.Join(workRoot, ".tasks") + `"

[project.workflow]
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
[project.project]
name = "Sub Project"
owner = "john"

[project.templates]
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
[project.project]
name = "My Project"
description = "Standalone project"

[project.storage]
backend = "markdown-fs"
path = "./.tasks"

[project.workflow]
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
[project.project]
name = "Root Project"

[project.workflow]
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
[project.project]
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
[project.project]
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

	// Note: Level 3 defines [project.project] which replaces the entire section,
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
[project.project]
name = "Test"

[project.storage]
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
