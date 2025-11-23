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
		ActiveProject: "",
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

func TestProjectLister_GetActive(t *testing.T) {
	// Create config with active project
	globalConfig := &OpentaskGlobalConfigFile{
		Projects: []GlobalProjectConfig{
			{
				ID:   "proj1",
				Name: "Project One",
			},
		},
		ActiveProject: "proj1",
	}

	lister := NewProjectLister(globalConfig)
	active := lister.GetActive()

	if active != "proj1" {
		t.Errorf("Expected active project 'proj1', got %q", active)
	}
}

func TestProjectLister_GetActive_NoActiveProject(t *testing.T) {
	globalConfig := &OpentaskGlobalConfigFile{
		Projects: []GlobalProjectConfig{
			{
				ID:   "proj1",
				Name: "Project One",
			},
		},
		ActiveProject: "",
	}

	lister := NewProjectLister(globalConfig)
	active := lister.GetActive()

	if active != "" {
		t.Errorf("Expected empty active project, got %q", active)
	}
}

func TestProjectLister_GetActive_NilConfig(t *testing.T) {
	lister := NewProjectLister(nil)
	active := lister.GetActive()

	if active != "" {
		t.Errorf("Expected empty string for nil config, got %q", active)
	}
}

func TestProjectLister_List_WithActiveMarker(t *testing.T) {
	globalConfig := &OpentaskGlobalConfigFile{
		Projects: []GlobalProjectConfig{
			{
				ID:   "proj1",
				Name: "Project One",
			},
			{
				ID:   "proj2",
				Name: "Project Two",
			},
		},
		ActiveProject: "proj1",
	}

	lister := NewProjectLister(globalConfig)
	result := lister.List()

	// proj1 should have active marker
	lines := strings.Split(result, "\n")
	found := false
	for _, line := range lines {
		if strings.Contains(line, "proj1") && strings.Contains(line, "*") {
			found = true
			break
		}
	}
	if !found {
		t.Error("Active project marker (*) not found for proj1")
	}

	// proj2 should not have active marker
	found = false
	for _, line := range lines {
		if strings.Contains(line, "proj2") && strings.Contains(line, "*") {
			found = true
			break
		}
	}
	if found {
		t.Error("Inactive project proj2 should not have active marker")
	}
}

func TestProjectLister_formatProjectEntry(t *testing.T) {
	globalConfig := &OpentaskGlobalConfigFile{
		Projects:      []GlobalProjectConfig{},
		ActiveProject: "",
	}

	lister := NewProjectLister(globalConfig)

	proj := GlobalProjectConfig{
		ID:   "testproj",
		Name: "Test Project",
		Storage: &StorageConfig{
			Path: "/home/user/test",
		},
		Context: []ProjectContext{
			{Path: "/home/user/test/dir1"},
			{Path: "/home/user/test/dir2"},
		},
	}

	result := lister.formatProjectEntry(proj)

	// Verify format
	if !strings.Contains(result, "testproj") {
		t.Error("Entry does not contain project ID")
	}
	if !strings.Contains(result, "Test Project") {
		t.Error("Entry does not contain project name")
	}
	if !strings.Contains(result, "Storage:") {
		t.Error("Entry does not contain storage label")
	}
	if !strings.Contains(result, "Contexts:") {
		t.Error("Entry does not contain contexts label")
	}
	if !strings.Contains(result, "dir1") {
		t.Error("Entry does not contain first context")
	}
	if !strings.Contains(result, "dir2") {
		t.Error("Entry does not contain second context")
	}
}

func TestProjectLister_formatProjectEntry_NoStorage(t *testing.T) {
	globalConfig := &OpentaskGlobalConfigFile{
		Projects: []GlobalProjectConfig{},
	}

	lister := NewProjectLister(globalConfig)

	proj := GlobalProjectConfig{
		ID:      "testproj",
		Name:    "Test Project",
		Storage: nil,
		Context: []ProjectContext{},
	}

	result := lister.formatProjectEntry(proj)

	// Verify basic format still works
	if !strings.Contains(result, "testproj") {
		t.Error("Entry does not contain project ID")
	}
	if !strings.Contains(result, "Test Project") {
		t.Error("Entry does not contain project name")
	}
	if !strings.Contains(result, "Contexts: (none)") {
		t.Error("Entry does not indicate no contexts")
	}
}

func TestProjectLister_formatProjectEntry_NoName_UseID(t *testing.T) {
	globalConfig := &OpentaskGlobalConfigFile{
		Projects: []GlobalProjectConfig{},
	}

	lister := NewProjectLister(globalConfig)

	proj := GlobalProjectConfig{
		ID:   "proj-id-123",
		Name: "", // No name, should use ID
	}

	result := lister.formatProjectEntry(proj)

	// Should display ID in name position when name is empty
	if !strings.Contains(result, "proj-id-123") {
		t.Error("Entry should display ID when name is empty")
	}
}
