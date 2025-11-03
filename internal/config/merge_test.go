package config

import (
	"os"
	"path/filepath"
	"testing"
)

// TestLoadGlobalConfigBasic tests loading a global config file
func TestLoadGlobalConfigBasic(t *testing.T) {
	// Create a temporary global config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "global.toml")

	globalConfig := `
[global]
active_project = "work"

[workflow]
statuses = ["todo", "done"]
initial = "todo"

[[global.projects]]
id = "work"
name = "Work Tasks"

[global.projects.storage]
backend = "markdown-fs"
path = "~/work/.tasks"

[[global.projects]]
id = "personal"
name = "Personal Tasks"

[global.projects.storage]
backend = "markdown-fs"
path = "~/personal/.tasks"
`

	if err := os.WriteFile(configPath, []byte(globalConfig), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	cfg, err := LoadGlobalConfig(configPath)
	if err != nil {
		t.Fatalf("LoadGlobalConfig failed: %v", err)
	}

	if cfg == nil {
		t.Fatal("LoadGlobalConfig returned nil")
	}

	if cfg.Global == nil {
		t.Fatal("Global section is nil")
	}

	if cfg.Global.ActiveProject != "work" {
		t.Errorf("ActiveProject = %q, want %q", cfg.Global.ActiveProject, "work")
	}

	if len(cfg.Global.Projects) != 2 {
		t.Errorf("Projects count = %d, want 2", len(cfg.Global.Projects))
	}

	if cfg.Global.Projects[0].ID != "work" {
		t.Errorf("First project ID = %q, want %q", cfg.Global.Projects[0].ID, "work")
	}
}

// TestLoadProjectConfigBasic tests loading a project config file
func TestLoadProjectConfigBasic(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".opentask.toml")

	projectConfig := `
[project.project]
name = "My Project"
description = "A test project"
owner = "me"

[project]
active_project = "my-project"

[project.storage]
backend = "markdown-fs"
path = "./.tasks"

[project.workflow]
statuses = ["todo", "in-progress", "done"]
initial = "todo"
`

	if err := os.WriteFile(configPath, []byte(projectConfig), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	cfg, err := LoadProjectConfig(configPath)
	if err != nil {
		t.Fatalf("LoadProjectConfig failed: %v", err)
	}

	if cfg == nil {
		t.Fatal("LoadProjectConfig returned nil")
	}

	if cfg.Project == nil {
		t.Fatal("Project section is nil")
	}

	if cfg.Project.Project.Name != "My Project" {
		t.Errorf("Project.Name = %q, want %q", cfg.Project.Project.Name, "My Project")
	}

	if cfg.Project.ActiveProject != "my-project" {
		t.Errorf("ActiveProject = %q, want %q", cfg.Project.ActiveProject, "my-project")
	}

	if cfg.Project.Storage.Path != "./.tasks" {
		t.Errorf("Storage.Path = %q, want %q", cfg.Project.Storage.Path, "./.tasks")
	}
}

// TestResolveProjectConfigSimpleGlobalOnly tests resolution with only global config
func TestResolveProjectConfigSimpleGlobalOnly(t *testing.T) {
	tmpDir := t.TempDir()

	// Create global config
	globalDir := filepath.Join(tmpDir, "config")
	os.MkdirAll(globalDir, 0755)
	globalPath := filepath.Join(globalDir, "config.toml")

	globalConfig := `
[global]
active_project = "test"

[[projects]]
id = "test"
name = "Test Project"

[projects.storage]
backend = "markdown-fs"
path = "` + tmpDir + `/.tasks"

[workflow]
statuses = ["todo", "done"]
initial = "todo"
`

	if err := os.WriteFile(globalPath, []byte(globalConfig), 0644); err != nil {
		t.Fatalf("Failed to write global config: %v", err)
	}

	// Set HOME to tmpDir so global config is found
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", oldHome)

	// Now resolve config at a project directory
	projectDir := filepath.Join(tmpDir, "project")
	os.MkdirAll(projectDir, 0755)

	resolved, err := ResolveProjectConfig(projectDir)
	if err != nil {
		t.Fatalf("ResolveProjectConfig failed: %v", err)
	}

	if resolved == nil {
		t.Fatal("ResolveProjectConfig returned nil")
	}

	if resolved.Workflow == nil {
		t.Fatal("Workflow is nil")
	}

	if len(resolved.Workflow.Statuses) == 0 {
		t.Fatal("Workflow.Statuses is empty")
	}
}

// TestResolveProjectConfigProjectOverridesGlobal tests that project config overrides global
func TestResolveProjectConfigProjectOverridesGlobal(t *testing.T) {
	tmpDir := t.TempDir()

	// Create global config
	globalDir := filepath.Join(tmpDir, ".config", "opentask")
	os.MkdirAll(globalDir, 0755)
	globalPath := filepath.Join(globalDir, "config.toml")

	globalConfig := `
[global]
active_project = "test"

[[global.projects]]
id = "test"
name = "Test Project"

[global.projects.storage]
backend = "markdown-fs"
path = "` + filepath.Join(tmpDir, ".tasks") + `"

[workflow]
statuses = ["todo", "done"]
initial = "todo"
`

	if err := os.WriteFile(globalPath, []byte(globalConfig), 0644); err != nil {
		t.Fatalf("Failed to write global config: %v", err)
	}

	// Create project config
	projectDir := filepath.Join(tmpDir, ".tasks")
	os.MkdirAll(projectDir, 0755)
	projectPath := filepath.Join(projectDir, ".opentask.toml")

	projectConfig := `
[project.project]
name = "Overridden Project"

[project.workflow]
statuses = ["backlog", "todo", "in-progress", "done"]
initial = "todo"
`

	if err := os.WriteFile(projectPath, []byte(projectConfig), 0644); err != nil {
		t.Fatalf("Failed to write project config: %v", err)
	}

	// Set HOME to tmpDir so global config is found
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", oldHome)

	resolved, err := ResolveProjectConfig(projectDir)
	if err != nil {
		t.Fatalf("ResolveProjectConfig failed: %v", err)
	}

	// Project workflow should override global
	if len(resolved.Workflow.Statuses) != 4 {
		t.Errorf("Workflow.Statuses count = %d, want 4 (from project), got %v", len(resolved.Workflow.Statuses), resolved.Workflow.Statuses)
	}

	// Project name should be used
	if resolved.Project.Name != "Overridden Project" {
		t.Errorf("Project.Name = %q, want %q", resolved.Project.Name, "Overridden Project")
	}
}

// TestResolveProjectConfigActiveProjectDerivation tests auto-populating active_project
func TestResolveProjectConfigActiveProjectDerivation(t *testing.T) {
	tmpDir := t.TempDir()

	// Create global config
	globalDir := filepath.Join(tmpDir, ".config", "opentask")
	os.MkdirAll(globalDir, 0755)
	globalPath := filepath.Join(globalDir, "config.toml")

	projectStoragePath := filepath.Join(tmpDir, "work", ".tasks")

	globalConfig := `
[global]
active_project = "work"

[[global.projects]]
id = "work"
name = "Work"

[global.projects.storage]
backend = "markdown-fs"
path = "` + projectStoragePath + `"
`

	if err := os.WriteFile(globalPath, []byte(globalConfig), 0644); err != nil {
		t.Fatalf("Failed to write global config: %v", err)
	}

	// Create project config WITHOUT active_project field
	projectDir := projectStoragePath
	os.MkdirAll(projectDir, 0755)
	projectPath := filepath.Join(projectDir, ".opentask.toml")

	projectConfig := `
[project.project]
name = "Work Project"
`

	if err := os.WriteFile(projectPath, []byte(projectConfig), 0644); err != nil {
		t.Fatalf("Failed to write project config: %v", err)
	}

	// Set HOME to tmpDir so global config is found
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", oldHome)

	resolved, err := ResolveProjectConfig(projectDir)
	if err != nil {
		t.Fatalf("ResolveProjectConfig failed: %v", err)
	}

	// active_project should be derived from matching global project
	if resolved.ActiveProject != "work" {
		t.Errorf("ActiveProject = %q, want %q (derived from global projects)", resolved.ActiveProject, "work")
	}
}

// TestResolveProjectConfigFallbackToDirectoryName tests deriving from directory name
func TestResolveProjectConfigFallbackToDirectoryName(t *testing.T) {
	tmpDir := t.TempDir()

	// Create project config WITHOUT active_project and NO global config match
	projectDir := filepath.Join(tmpDir, "my-project")
	os.MkdirAll(projectDir, 0755)
	projectPath := filepath.Join(projectDir, ".opentask.toml")

	projectConfig := `
[project]
name = "My Project"

[storage]
backend = "markdown-fs"
path = "./.tasks"
`

	if err := os.WriteFile(projectPath, []byte(projectConfig), 0644); err != nil {
		t.Fatalf("Failed to write project config: %v", err)
	}

	// Set HOME to tmpDir (no global config)
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", oldHome)

	resolved, err := ResolveProjectConfig(projectDir)
	if err != nil {
		t.Fatalf("ResolveProjectConfig failed: %v", err)
	}

	// active_project should be derived from directory name
	if resolved.ActiveProject != "my-project" {
		t.Errorf("ActiveProject = %q, want %q (derived from dir name)", resolved.ActiveProject, "my-project")
	}
}

// TestMergeConfigsSimple tests merging a single config
func TestMergeConfigsSimple(t *testing.T) {
	global := &OpentaskGlobalConfigFile{
		Global: &OpentaskConfigGlobalSchema{
			ActiveProject: "test",
		},
		Core: &OpentaskConfigCoreSchema{
			Workflow: &WorkflowConfig{
				Statuses: []string{"todo", "done"},
				Initial:  "todo",
			},
		},
	}

	resolved := NewResolvedConfig()
	resolved = mergeGlobalConfig(resolved, global)

	if len(resolved.Workflow.Statuses) != 2 {
		t.Errorf("Workflow.Statuses count = %d, want 2", len(resolved.Workflow.Statuses))
	}

	if resolved.Workflow.Statuses[0] != "todo" {
		t.Errorf("Workflow.Statuses[0] = %q, want %q", resolved.Workflow.Statuses[0], "todo")
	}
}

// TestMergeConfigsProjectOverridesGlobal tests correct override priority
func TestMergeConfigsProjectOverridesGlobal(t *testing.T) {
	global := &OpentaskGlobalConfigFile{
		Core: &OpentaskConfigCoreSchema{
			Workflow: &WorkflowConfig{
				Statuses: []string{"todo", "done"},
				Initial:  "todo",
			},
		},
	}

	project := &OpentaskProjectConfigFile{
		Core: &OpentaskConfigCoreSchema{
			Workflow: &WorkflowConfig{
				Statuses: []string{"todo", "in-progress", "done"},
				Initial:  "todo",
			},
		},
	}

	resolved := NewResolvedConfig()
	resolved = mergeGlobalConfig(resolved, global)
	resolved = mergeProjectConfig(resolved, project)

	if len(resolved.Workflow.Statuses) != 3 {
		t.Errorf("Workflow.Statuses count = %d, want 3 (from project)", len(resolved.Workflow.Statuses))
	}
}
