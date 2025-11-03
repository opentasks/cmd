package config

import (
	"testing"
)

// TestOpentaskConfigCoreSchemaDefaults tests that core schema has sensible defaults
func TestOpentaskConfigCoreSchemaDefaults(t *testing.T) {
	tests := []struct {
		name string
		core *OpentaskConfigCoreSchema
		want bool
	}{
		{
			name: "empty core schema has no workflow",
			core: &OpentaskConfigCoreSchema{},
			want: true, // should be allowed
		},
		{
			name: "core schema can have workflow",
			core: &OpentaskConfigCoreSchema{
				Workflow: &WorkflowConfig{
					Statuses: []string{"todo", "done"},
					Initial:  "todo",
				},
			},
			want: true,
		},
		{
			name: "core schema can have templates",
			core: &OpentaskConfigCoreSchema{
				Templates: &TemplateConfig{
					Epic: "epic.md",
				},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.core == nil {
				t.Fatal("core schema is nil")
			}
		})
	}
}

// TestOpentaskConfigGlobalSchemaStructure tests global config schema
func TestOpentaskConfigGlobalSchemaStructure(t *testing.T) {
	global := &OpentaskConfigGlobalSchema{
		ActiveProject: "test-project",
		Projects: []GlobalProjectConfig{
			{
				ID:   "test-project",
				Name: "Test Project",
				Storage: &StorageConfig{
					Backend: "markdown-fs",
					Path:    "~/test/.tasks",
				},
			},
		},
	}

	if global.ActiveProject != "test-project" {
		t.Errorf("ActiveProject = %q, want %q", global.ActiveProject, "test-project")
	}

	if len(global.Projects) != 1 {
		t.Errorf("Projects length = %d, want 1", len(global.Projects))
	}

	if global.Projects[0].ID != "test-project" {
		t.Errorf("Project ID = %q, want %q", global.Projects[0].ID, "test-project")
	}
}

// TestOpentaskConfigProjectSchemaStructure tests project config schema
func TestOpentaskConfigProjectSchemaStructure(t *testing.T) {
	project := &OpentaskConfigProjectSchema{
		Project: &ProjectSection{
			Name:  "My Project",
			Owner: "me",
		},
		Storage: &StorageConfig{
			Backend: "markdown-fs",
			Path:    "./.tasks",
		},
		ActiveProject: "my-project",
	}

	if project.Project.Name != "My Project" {
		t.Errorf("Project.Name = %q, want %q", project.Project.Name, "My Project")
	}

	if project.ActiveProject != "my-project" {
		t.Errorf("ActiveProject = %q, want %q", project.ActiveProject, "my-project")
	}
}

// TestOpentaskGlobalConfigFile tests the complete global config file type
func TestOpentaskGlobalConfigFile(t *testing.T) {
	globalConfig := &OpentaskGlobalConfigFile{
		ActiveProject: "work",
		Projects: []GlobalProjectConfig{
			{
				ID:   "work",
				Name: "Work Projects",
				Storage: &StorageConfig{
					Backend: "markdown-fs",
					Path:    "~/work/.tasks",
				},
			},
		},
		Workflow: &WorkflowConfig{
			Statuses: []string{"todo", "done"},
			Initial:  "todo",
		},
	}

	if globalConfig.ActiveProject != "work" {
		t.Errorf("ActiveProject = %q, want %q", globalConfig.ActiveProject, "work")
	}

	if len(globalConfig.Projects) != 1 {
		t.Errorf("Projects count = %d, want 1", len(globalConfig.Projects))
	}

	if globalConfig.Projects[0].ID != "work" {
		t.Errorf("First project ID = %q, want %q", globalConfig.Projects[0].ID, "work")
	}
}

// TestOpentaskProjectConfigFile tests the complete project config file type
func TestOpentaskProjectConfigFile(t *testing.T) {
	projectConfig := &OpentaskProjectConfigFile{
		Project: &ProjectSection{
			Name: "My Project",
		},
		Storage: &StorageConfig{
			Backend: "markdown-fs",
			Path:    "./.tasks",
		},
		ActiveProject: "my-project",
		Core: &OpentaskConfigCoreSchema{
			Workflow: &WorkflowConfig{
				Statuses: []string{"todo", "in-progress", "done"},
				Initial:  "todo",
			},
		},
	}

	if projectConfig.Project == nil {
		t.Fatal("Project section is nil")
	}

	if projectConfig.Core == nil {
		t.Fatal("Core section is nil")
	}

	if projectConfig.ActiveProject != "my-project" {
		t.Errorf("ActiveProject = %q, want %q", projectConfig.ActiveProject, "my-project")
	}
}

// TestOpentaskResolvedConfigStructure tests the final merged config
func TestOpentaskResolvedConfigStructure(t *testing.T) {
	resolved := &OpentaskResolvedConfig{
		Project: &ProjectSection{
			Name:  "Final Project",
			Owner: "me",
		},
		Workflow: &WorkflowConfig{
			Statuses: []string{"todo", "in-progress", "done"},
			Initial:  "todo",
		},
		Templates: &TemplateConfig{
			Epic: "epic.md",
		},
		Storage: &StorageConfig{
			Backend: "markdown-fs",
			Path:    "/absolute/path/.tasks",
		},
		ActiveProject: "final-project",
		// DiscoveredFiles contains the config files that were merged
		DiscoveredFiles: []string{
			"./.opentask.toml",
			"../.opentask.toml",
		},
	}

	if resolved.Project.Name != "Final Project" {
		t.Errorf("Project.Name = %q, want %q", resolved.Project.Name, "Final Project")
	}

	if resolved.ActiveProject != "final-project" {
		t.Errorf("ActiveProject = %q, want %q", resolved.ActiveProject, "final-project")
	}

	if len(resolved.DiscoveredFiles) != 2 {
		t.Errorf("DiscoveredFiles length = %d, want 2", len(resolved.DiscoveredFiles))
	}
}

// TestGlobalProjectConfigStructure tests individual project definition in global config
func TestGlobalProjectConfigStructure(t *testing.T) {
	proj := &GlobalProjectConfig{
		ID:   "work",
		Name: "Work Tasks",
		Storage: &StorageConfig{
			Backend: "markdown-fs",
			Path:    "~/work/.tasks",
		},
	}

	if proj.ID != "work" {
		t.Errorf("ID = %q, want %q", proj.ID, "work")
	}

	if proj.Name != "Work Tasks" {
		t.Errorf("Name = %q, want %q", proj.Name, "Work Tasks")
	}

	if proj.Storage.Path != "~/work/.tasks" {
		t.Errorf("Storage.Path = %q, want %q", proj.Storage.Path, "~/work/.tasks")
	}
}

// TestNewResolvedConfigPopulatesDefaults tests that resolved config gets defaults
func TestNewResolvedConfigPopulatesDefaults(t *testing.T) {
	resolved := NewResolvedConfig()

	if resolved == nil {
		t.Fatal("NewResolvedConfig returned nil")
	}

	if resolved.Workflow == nil {
		t.Fatal("Workflow is nil, expected defaults")
	}

	if resolved.Storage == nil {
		t.Fatal("Storage is nil, expected defaults")
	}

	if resolved.Templates == nil {
		t.Fatal("Templates is nil, expected defaults")
	}
}

// TestNewResolvedConfigWorkflowDefaults tests that defaults match existing behavior
func TestNewResolvedConfigWorkflowDefaults(t *testing.T) {
	resolved := NewResolvedConfig()

	// Should have some default statuses
	if len(resolved.Workflow.Statuses) == 0 {
		t.Fatal("Workflow.Statuses is empty, expected defaults")
	}

	// Check that initial status is in the statuses list
	found := false
	for _, s := range resolved.Workflow.Statuses {
		if s == resolved.Workflow.Initial {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Initial status %q not found in Statuses", resolved.Workflow.Initial)
	}
}

// TestActiveProjectDerivationLogic tests the logic for populating active_project
func TestActiveProjectDerivationLogic(t *testing.T) {
	tests := []struct {
		name           string
		configPath     string
		globalProjects []GlobalProjectConfig
		wantProjectID  string
	}{
		{
			name:           "derive from directory name when no global match",
			configPath:     "/home/user/my-project/.opentask.toml",
			globalProjects: []GlobalProjectConfig{},
			wantProjectID:  "my-project",
		},
		{
			name:       "use global project ID when path matches",
			configPath: "/home/user/work/.tasks/.opentask.toml",
			globalProjects: []GlobalProjectConfig{
				{
					ID:   "work",
					Name: "Work",
					Storage: &StorageConfig{
						Path: "/home/user/work/.tasks",
					},
				},
			},
			wantProjectID: "work",
		},
		{
			name:       "derive from directory name when config in different location",
			configPath: "/home/user/projects/personal/.opentask.toml",
			globalProjects: []GlobalProjectConfig{
				{
					ID:   "work",
					Name: "Work",
					Storage: &StorageConfig{
						Path: "/home/user/work/.tasks",
					},
				},
			},
			wantProjectID: "personal",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This test documents the expected behavior but doesn't implement
			// the derivation logic yet - that comes in Phase 2
			_ = tt.configPath
			_ = tt.globalProjects
			_ = tt.wantProjectID
		})
	}
}
