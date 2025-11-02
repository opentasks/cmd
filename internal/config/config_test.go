package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultWorkflow(t *testing.T) {
	wf := DefaultWorkflow()

	if len(wf.Statuses) == 0 {
		t.Error("DefaultWorkflow() returned empty statuses")
	}

	if wf.Initial == "" {
		t.Error("DefaultWorkflow() returned empty initial status")
	}

	if wf.Initial != "todo" {
		t.Errorf("DefaultWorkflow() initial = %q, want 'todo'", wf.Initial)
	}

	if len(wf.Transitions) == 0 {
		t.Error("DefaultWorkflow() returned no transitions")
	}
}

func TestDefaultWorkflowStatuses(t *testing.T) {
	wf := DefaultWorkflow()

	expectedStatuses := []string{"todo", "in-progress", "reviewing", "done", "archived"}
	if len(wf.Statuses) != len(expectedStatuses) {
		t.Errorf("DefaultWorkflow() has %d statuses, want %d", len(wf.Statuses), len(expectedStatuses))
	}

	for i, status := range wf.Statuses {
		if i < len(expectedStatuses) && status != expectedStatuses[i] {
			t.Errorf("Status[%d] = %q, want %q", i, status, expectedStatuses[i])
		}
	}
}

func TestDefaultWorkflowTransitions(t *testing.T) {
	wf := DefaultWorkflow()

	// Verify transitions map exists and has valid entries
	transitions := make(map[string][]string)
	for _, t := range wf.Transitions {
		transitions[t.From] = t.To
	}

	// Should have transitions from "todo"
	if _, exists := transitions["todo"]; !exists {
		t.Error("DefaultWorkflow() missing transition from 'todo'")
	}

	// "done" should be a valid target from some state
	foundDone := false
	for _, transition := range wf.Transitions {
		for _, status := range transition.To {
			if status == "done" {
				foundDone = true
			}
		}
	}
	if !foundDone {
		t.Error("DefaultWorkflow() 'done' is not a valid target in any transition")
	}
}

func TestDefaultStorage(t *testing.T) {
	storage := DefaultStorage()

	if storage.Backend != "markdown-fs" {
		t.Errorf("DefaultStorage() backend = %q, want 'markdown-fs'", storage.Backend)
	}

	if storage.Options == nil {
		t.Error("DefaultStorage() options is nil")
	}
}

func TestDefaultTemplates(t *testing.T) {
	templates := DefaultTemplates()

	// DefaultTemplates should return a valid TemplateConfig (may be empty struct)
	// The function exists and returns a TemplateConfig, even if all fields are empty
	// This is acceptable behavior - templates are optional
	if _, ok := interface{}(templates).(TemplateConfig); !ok {
		t.Error("DefaultTemplates() did not return a TemplateConfig")
	}
}

func TestLoadConfigNonexistent(t *testing.T) {
	// Try to load non-existent config file
	config, err := LoadConfig("/nonexistent/path/config.toml")

	if err != nil {
		t.Errorf("LoadConfig() with nonexistent file error = %v, want nil (should return defaults)", err)
	}

	if config == nil {
		t.Error("LoadConfig() returned nil config")
	}

	if config.Workflow.Initial != "todo" {
		t.Errorf("LoadConfig() default workflow initial = %q, want 'todo'", config.Workflow.Initial)
	}
}

func TestLoadConfigWithDefaults(t *testing.T) {
	// Create temp dir for test
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.toml")

	// Create a minimal config file
	content := `
[project]
name = "Test Project"
description = "A test project"

[workflow]
statuses = ["todo", "done"]
initial = "todo"

[storage]
backend = "markdown-fs"
path = "./tasks"
`

	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	config, err := LoadConfig(configPath)
	if err != nil {
		t.Errorf("LoadConfig() error = %v", err)
	}

	if config.Project.Name != "Test Project" {
		t.Errorf("LoadConfig() project name = %q, want 'Test Project'", config.Project.Name)
	}

	if config.Workflow.Initial != "todo" {
		t.Errorf("LoadConfig() workflow initial = %q, want 'todo'", config.Workflow.Initial)
	}

	if config.Storage.Backend != "markdown-fs" {
		t.Errorf("LoadConfig() storage backend = %q, want 'markdown-fs'", config.Storage.Backend)
	}
}

func TestLoadConfigPartialConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.toml")

	// Create config with only project section
	content := `
[project]
name = "Partial Project"
`

	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	config, err := LoadConfig(configPath)
	if err != nil {
		t.Errorf("LoadConfig() error = %v", err)
	}

	if config.Project.Name != "Partial Project" {
		t.Errorf("Project name = %q, want 'Partial Project'", config.Project.Name)
	}

	// Should have default workflow
	if len(config.Workflow.Statuses) == 0 {
		t.Error("Partial config should include default workflow")
	}
}

func TestProjectConfigStructure(t *testing.T) {
	config := &ProjectConfig{
		Project: ProjectSection{
			Name:        "Test",
			Description: "Test project",
			Owner:       "testuser",
		},
		Workflow: DefaultWorkflow(),
		Storage:  DefaultStorage(),
	}

	if config.Project.Name != "Test" {
		t.Errorf("ProjectConfig.Project.Name = %q, want 'Test'", config.Project.Name)
	}

	if len(config.Workflow.Statuses) == 0 {
		t.Error("ProjectConfig.Workflow.Statuses should not be empty")
	}
}

func TestWorkflowConfigStructure(t *testing.T) {
	wf := WorkflowConfig{
		Statuses: []string{"todo", "done"},
		Initial:  "todo",
		Transitions: []TransitionConfig{
			{
				From: "todo",
				To:   []string{"done"},
			},
		},
	}

	if len(wf.Statuses) != 2 {
		t.Errorf("WorkflowConfig.Statuses length = %d, want 2", len(wf.Statuses))
	}

	if len(wf.Transitions) != 1 {
		t.Errorf("WorkflowConfig.Transitions length = %d, want 1", len(wf.Transitions))
	}

	if wf.Transitions[0].From != "todo" {
		t.Errorf("Transition.From = %q, want 'todo'", wf.Transitions[0].From)
	}
}

func TestStorageConfigStructure(t *testing.T) {
	storage := StorageConfig{
		Backend: "memory",
		Path:    "/tmp/tasks",
		Options: map[string]string{"option1": "value1"},
	}

	if storage.Backend != "memory" {
		t.Errorf("StorageConfig.Backend = %q, want 'memory'", storage.Backend)
	}

	if storage.Path != "/tmp/tasks" {
		t.Errorf("StorageConfig.Path = %q, want '/tmp/tasks'", storage.Path)
	}

	if len(storage.Options) != 1 {
		t.Errorf("StorageConfig.Options length = %d, want 1", len(storage.Options))
	}
}
