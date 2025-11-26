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
active_project = "work"

[workflow]
statuses = ["todo", "done"]
initial = "todo"

[[projects]]
id = "work"
name = "Work Tasks"

[projects.storage]
backend = "markdown-fs"
path = "~/work/.tasks"

[[projects]]
id = "personal"
name = "Personal Tasks"

[projects.storage]
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

	if cfg.ActiveProject != "work" {
		t.Errorf("ActiveProject = %q, want %q", cfg.ActiveProject, "work")
	}

	if len(cfg.Projects) != 2 {
		t.Errorf("Projects count = %d, want 2", len(cfg.Projects))
	}

	if cfg.Projects[0].ID != "work" {
		t.Errorf("First project ID = %q, want %q", cfg.Projects[0].ID, "work")
	}
}

// TestLoadProjectConfigBasic tests loading a project config file
func TestLoadProjectConfigBasic(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".opentask.toml")

	projectConfig := `
active_project = "my-project"

[project]
name = "My Project"
description = "A test project"
owner = "me"

[storage]
backend = "markdown-fs"
path = "./.tasks"

[workflow]
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

	if cfg.Project.Name != "My Project" {
		t.Errorf("Project.Name = %q, want %q", cfg.Project.Name, "My Project")
	}

	if cfg.ActiveProject != "my-project" {
		t.Errorf("ActiveProject = %q, want %q", cfg.ActiveProject, "my-project")
	}

	if cfg.Storage.Path != "./.tasks" {
		t.Errorf("Storage.Path = %q, want %q", cfg.Storage.Path, "./.tasks")
	}
}

// TestResolveProjectConfigSimpleGlobalOnly tests resolution with only global config
func TestResolveProjectConfigSimpleGlobalOnly(t *testing.T) {
	tmpDir := t.TempDir()

	// Create global config
	globalDir := filepath.Join(tmpDir, "config")
	_ = os.MkdirAll(globalDir, 0755)
	globalPath := filepath.Join(globalDir, "config.toml")

	globalConfig := `
active_project = "test"

[[projects]]
id = "test"
name = "Test Project"

[projects.storage]
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
	_ = os.MkdirAll(projectDir, 0755)
	projectPath := filepath.Join(projectDir, ".opentask.toml")

	projectConfig := `
[project]
name = "Overridden Project"

[workflow]
statuses = ["backlog", "todo", "in-progress", "done"]
initial = "todo"
`

	if err := os.WriteFile(projectPath, []byte(projectConfig), 0644); err != nil {
		t.Fatalf("Failed to write project config: %v", err)
	}

	// Set HOME to tmpDir so global config is found
	oldHome := os.Getenv("HOME")
	_ = os.Setenv("HOME", tmpDir)
	defer func() { _ = os.Setenv("HOME", oldHome) }()

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
	_ = os.MkdirAll(globalDir, 0755)
	globalPath := filepath.Join(globalDir, "config.toml")

	projectStoragePath := filepath.Join(tmpDir, "work", ".tasks")

	globalConfig := `
active_project = "work"

[[projects]]
id = "work"
name = "Work"

[projects.storage]
backend = "markdown-fs"
path = "` + projectStoragePath + `"
`

	if err := os.WriteFile(globalPath, []byte(globalConfig), 0644); err != nil {
		t.Fatalf("Failed to write global config: %v", err)
	}

	// Create project config WITHOUT active_project field
	projectDir := projectStoragePath
	_ = os.MkdirAll(projectDir, 0755)
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
	_ = os.Setenv("HOME", tmpDir)
	defer func() { _ = os.Setenv("HOME", oldHome) }()

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
	_ = os.MkdirAll(projectDir, 0755)
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
	_ = os.Setenv("HOME", tmpDir)
	defer func() { _ = os.Setenv("HOME", oldHome) }()

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
		ActiveProject: "test",
		Workflow: &WorkflowConfig{
			Statuses: []string{"todo", "done"},
			Initial:  "todo",
		},
	}

	resolved := NewResolvedConfig()
	resolved = MergeGlobalConfig(resolved, global)

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
		Workflow: &WorkflowConfig{
			Statuses: []string{"todo", "done"},
			Initial:  "todo",
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
	resolved = MergeGlobalConfig(resolved, global)
	resolved = MergeProjectConfig(resolved, project)

	if len(resolved.Workflow.Statuses) != 3 {
		t.Errorf("Workflow.Statuses count = %d, want 3 (from project)", len(resolved.Workflow.Statuses))
	}
}

// TestFindProjectByContext tests matching a cwd to a project by context
func TestFindProjectByContext(t *testing.T) {
	tests := []struct {
		name         string
		cwd          string
		contexts     map[string][]string // project id -> context paths
		expectedProj string              // expected matched project ID
		expectMatch  bool
	}{
		{
			name: "exact match",
			cwd:  "/mnt/repos/myproject",
			contexts: map[string][]string{
				"proj1": {"/mnt/repos/myproject"},
			},
			expectedProj: "proj1",
			expectMatch:  true,
		},
		{
			name: "subdirectory match",
			cwd:  "/mnt/repos/myproject/src/main",
			contexts: map[string][]string{
				"proj1": {"/mnt/repos/myproject"},
			},
			expectedProj: "proj1",
			expectMatch:  true,
		},
		{
			name: "longest match wins",
			cwd:  "/mnt/repos/myproject/src",
			contexts: map[string][]string{
				"proj1": {"/mnt/repos"},
				"proj2": {"/mnt/repos/myproject"},
				"proj3": {"/mnt/repos/myproject/src"},
			},
			expectedProj: "proj3",
			expectMatch:  true,
		},
		{
			name: "no match",
			cwd:  "/home/user/notes",
			contexts: map[string][]string{
				"proj1": {"/mnt/repos/myproject"},
			},
			expectMatch: false,
		},
		{
			name: "multiple contexts per project",
			cwd:  "/mnt/repos/myproject.worktrees/feature",
			contexts: map[string][]string{
				"proj1": {"/mnt/repos/myproject", "/mnt/repos/myproject.worktrees/feature"},
			},
			expectedProj: "proj1",
			expectMatch:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Build global projects
			var projects []GlobalProjectConfig
			for projID, contextPaths := range tt.contexts {
				proj := GlobalProjectConfig{
					ID:   projID,
					Name: projID,
					Storage: &StorageConfig{
						Backend: "markdown-fs",
						Path:    "/tmp/tasks",
					},
				}
				for _, ctxPath := range contextPaths {
					proj.Context = append(proj.Context, ProjectContext{Path: ctxPath})
				}
				projects = append(projects, proj)
			}

			// Test matching
			projID, proj := FindProjectByContext(tt.cwd, projects)

			if tt.expectMatch {
				if !tt.expectMatch && projID == "" {
					t.Errorf("expected match for %s, got no match", tt.cwd)
				}
				if projID != tt.expectedProj {
					t.Errorf("expected project %s, got %s", tt.expectedProj, projID)
				}
				if proj == nil {
					t.Errorf("expected matched project config, got nil")
				}
			} else {
				if projID != "" {
					t.Errorf("expected no match for %s, got project %s", tt.cwd, projID)
				}
				if proj != nil {
					t.Errorf("expected nil project, got %v", proj)
				}
			}
		})
	}
}

// TestAddContextPath tests adding context paths to a project
func TestAddContextPath(t *testing.T) {
	proj := &GlobalProjectConfig{
		ID:   "test",
		Name: "Test Project",
	}

	// Add a context path
	err := proj.AddContextPath("/mnt/repos/myproject")
	if err != nil {
		t.Errorf("failed to add context: %v", err)
	}

	// Verify it was added
	if len(proj.Context) != 1 {
		t.Errorf("expected 1 context, got %d", len(proj.Context))
	}

	// Try adding duplicate
	err = proj.AddContextPath("/mnt/repos/myproject")
	if err == nil {
		t.Errorf("should have failed adding duplicate context")
	}

	// Add another
	err = proj.AddContextPath("/mnt/repos/other")
	if err != nil {
		t.Errorf("failed to add second context: %v", err)
	}

	if len(proj.Context) != 2 {
		t.Errorf("expected 2 contexts, got %d", len(proj.Context))
	}
}

// TestRemoveContextPath tests removing context paths from a project
func TestRemoveContextPath(t *testing.T) {
	proj := &GlobalProjectConfig{
		ID:   "test",
		Name: "Test Project",
		Context: []ProjectContext{
			{Path: "/mnt/repos/myproject"},
			{Path: "/mnt/repos/other"},
		},
	}

	// Remove first context
	err := proj.RemoveContextPath("/mnt/repos/myproject")
	if err != nil {
		t.Errorf("failed to remove context: %v", err)
	}

	// Verify it was removed
	if len(proj.Context) != 1 {
		t.Errorf("expected 1 context, got %d", len(proj.Context))
	}

	// Try removing non-existent context
	err = proj.RemoveContextPath("/mnt/repos/nonexistent")
	if err == nil {
		t.Errorf("should have failed removing non-existent context")
	}
}

// TODO: TestResolveProjectConfigWithContext
// This test needs proper mocking of home directory to avoid finding real global config
// Core context matching is tested in TestFindProjectByContext
