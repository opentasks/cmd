package config

import (
	"strings"
	"testing"
)

func TestProjectLister_List(t *testing.T) {
	// Create test config with projects
	globalConfig := &OpentaskGlobalConfigFile{
		Projects: []GlobalProjectConfig{
			{
				ID:   "proj1",
				Name: "Project One",
				Storage: &StorageConfig{
					Path: "/home/user/proj1",
				},
				Context: []ProjectContext{
					{Path: "/home/user/proj1"},
				},
			},
			{
				ID:   "proj2",
				Name: "Project Two",
				Storage: &StorageConfig{
					Path: "/home/user/proj2",
				},
				Context: []ProjectContext{},
			},
		},
	}

	lister := NewProjectLister(globalConfig)
	result := lister.List()

	// Verify result contains project information
	if !strings.Contains(result, "proj1") {
		t.Error("Result does not contain proj1 ID")
	}
	if !strings.Contains(result, "proj2") {
		t.Error("Result does not contain proj2 ID")
	}
	if !strings.Contains(result, "Project One") {
		t.Error("Result does not contain 'Project One'")
	}
	if !strings.Contains(result, "Project Two") {
		t.Error("Result does not contain 'Project Two'")
	}
	if !strings.Contains(result, "Storage:") {
		t.Error("Result does not contain storage information")
	}
	if !strings.Contains(result, "Contexts:") {
		t.Error("Result does not contain contexts information")
	}
}

func TestProjectLister_List_NoProjects(t *testing.T) {
	// Create empty config
	globalConfig := &OpentaskGlobalConfigFile{
		Projects: []GlobalProjectConfig{},
	}

	lister := NewProjectLister(globalConfig)
	result := lister.List()

	// Empty result expected
	if result != "" {
		t.Errorf("Expected empty result for no projects, got: %q", result)
	}
}

func TestProjectLister_List_NilConfig(t *testing.T) {
	lister := NewProjectLister(nil)
	result := lister.List()

	// Empty result expected
	if result != "" {
		t.Errorf("Expected empty result for nil config, got: %q", result)
	}
}

// Note: GetActive() and active project markers removed
