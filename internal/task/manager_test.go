package task

import (
	"testing"

	"github.com/zenobi-us/opentask/internal/model"
)

func TestDetermineInitialStatus(t *testing.T) {
	manager := NewManager()

	tests := []struct {
		name        string
		taskType    string
		description string
		want        string
	}{
		{
			name:        "high-level task without description",
			taskType:    model.TypeEpic,
			description: "",
			want:        "backlog",
		},
		{
			name:        "high-level task with description",
			taskType:    model.TypePlan,
			description: "Some description",
			want:        "todo",
		},
		{
			name:        "regular task without description",
			taskType:    model.TypeTask,
			description: "",
			want:        "todo",
		},
		{
			name:        "regular task with description",
			taskType:    model.TypeTask,
			description: "Some description",
			want:        "todo",
		},
		{
			name:        "research type without description",
			taskType:    model.TypeResearch,
			description: "",
			want:        "backlog",
		},
		{
			name:        "story type without description",
			taskType:    model.TypeStory,
			description: "",
			want:        "backlog",
		},
		{
			name:        "decision type without description",
			taskType:    model.TypeDecision,
			description: "",
			want:        "backlog",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := manager.DetermineInitialStatus(tt.taskType, tt.description)
			if got != tt.want {
				t.Errorf("DetermineInitialStatus(%q, %q) = %q, want %q", tt.taskType, tt.description, got, tt.want)
			}
		})
	}
}

func TestCreateTask(t *testing.T) {
	manager := NewManager()

	task := manager.CreateTask(42, "Test Task", model.TypeTask, "todo", "Test description", []string{"tag1"}, 0)

	if task.ID != 42 {
		t.Errorf("task.ID = %d, want 42", task.ID)
	}
	if task.Title != "Test Task" {
		t.Errorf("task.Title = %q, want %q", task.Title, "Test Task")
	}
	if task.Type != model.TypeTask {
		t.Errorf("task.Type = %q, want %q", task.Type, model.TypeTask)
	}
	if task.Status != "todo" {
		t.Errorf("task.Status = %q, want %q", task.Status, "todo")
	}
	if task.Description != "Test description" {
		t.Errorf("task.Description = %q, want %q", task.Description, "Test description")
	}
	if len(task.Tags) != 1 || task.Tags[0] != "tag1" {
		t.Errorf("task.Tags = %v, want [tag1]", task.Tags)
	}
	if task.CreatedAt.IsZero() {
		t.Errorf("task.CreatedAt should be set")
	}
	if task.UpdatedAt.IsZero() {
		t.Errorf("task.UpdatedAt should be set")
	}
}

func TestCreateTaskWithParent(t *testing.T) {
	manager := NewManager()

	task := manager.CreateTask(42, "Child Task", model.TypeTask, "todo", "", []string{}, 10)

	if len(task.Relationships) != 1 {
		t.Errorf("task.Relationships should have 1 relationship, got %d", len(task.Relationships))
	}

	rel := task.Relationships[0]
	if rel.Type != model.RelParent {
		t.Errorf("relationship type = %q, want %q", rel.Type, model.RelParent)
	}
	if rel.TaskID != 10 {
		t.Errorf("relationship TaskID = %d, want 10", rel.TaskID)
	}
}

func TestMergeTags(t *testing.T) {
	manager := NewManager()

	existing := []string{"tag1", "tag2"}
	newTags := []string{"tag2", "tag3"}

	result := manager.MergeTags(existing, newTags)

	// Check that we have exactly 3 tags (union)
	if len(result) != 3 {
		t.Errorf("MergeTags returned %d tags, want 3", len(result))
	}

	// Check that all expected tags are present
	tagMap := make(map[string]bool)
	for _, tag := range result {
		tagMap[tag] = true
	}

	for _, expected := range []string{"tag1", "tag2", "tag3"} {
		if !tagMap[expected] {
			t.Errorf("MergeTags missing tag %q", expected)
		}
	}
}

func TestRemoveTags(t *testing.T) {
	manager := NewManager()

	existing := []string{"tag1", "tag2", "tag3"}
	toRemove := []string{"tag2"}

	result := manager.RemoveTags(existing, toRemove)

	if len(result) != 2 {
		t.Errorf("RemoveTags returned %d tags, want 2", len(result))
	}

	// Check that the remaining tags are correct
	tagMap := make(map[string]bool)
	for _, tag := range result {
		tagMap[tag] = true
	}

	if !tagMap["tag1"] || !tagMap["tag3"] {
		t.Errorf("RemoveTags didn't preserve correct tags")
	}

	if tagMap["tag2"] {
		t.Errorf("RemoveTags didn't remove tag2")
	}
}
