package display

import (
	"strings"
	"testing"
	"time"

	"github.com/zenobi-us/opentask/internal/model"
)

func TestTruncate(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		maxLen int
		want   string
	}{
		{
			name:   "string shorter than max",
			input:  "hello",
			maxLen: 10,
			want:   "hello",
		},
		{
			name:   "string equal to max",
			input:  "hello",
			maxLen: 5,
			want:   "hello",
		},
		{
			name:   "string longer than max",
			input:  "hello world",
			maxLen: 8,
			want:   "hello...",
		},
		{
			name:   "empty string",
			input:  "",
			maxLen: 5,
			want:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Truncate(tt.input, tt.maxLen)
			if got != tt.want {
				t.Errorf("Truncate(%q, %d) = %q, want %q", tt.input, tt.maxLen, got, tt.want)
			}
		})
	}
}

func TestConfigFileTree(t *testing.T) {
	files := []string{"/path/to/config1.toml", "/path/to/config2.toml"}
	result := ConfigFileTree(files)

	// Check that result contains expected elements
	if !strings.Contains(result, "config1.toml") {
		t.Errorf("ConfigFileTree should contain config1.toml")
	}
	if !strings.Contains(result, "config2.toml") {
		t.Errorf("ConfigFileTree should contain config2.toml")
	}
	if !strings.Contains(result, "(builtin) defaults") {
		t.Errorf("ConfigFileTree should contain (builtin) defaults")
	}
	if !strings.Contains(result, "└──") || !strings.Contains(result, "├──") {
		t.Errorf("ConfigFileTree should contain tree characters")
	}
}

func TestTaskTable(t *testing.T) {
	tasks := []*model.Task{
		{
			ID:        1,
			Title:     "Short title",
			Type:      "task",
			Status:    "todo",
			CreatedAt: time.Now(),
		},
		{
			ID:        2,
			Title:     "A very long title that should be truncated in the table display",
			Type:      "epic",
			Status:    "backlog",
			CreatedAt: time.Now(),
		},
	}

	result := TaskTable(tasks)

	// Check that result contains header and expected content
	if !strings.Contains(result, "ID") || !strings.Contains(result, "Type") {
		t.Errorf("TaskTable should contain header")
	}
	if !strings.Contains(result, "1") || !strings.Contains(result, "Short title") {
		t.Errorf("TaskTable should contain first task")
	}
	if !strings.Contains(result, "2") {
		t.Errorf("TaskTable should contain second task ID")
	}
}

func TestTaskTableEmpty(t *testing.T) {
	result := TaskTable([]*model.Task{})
	if result != "" {
		t.Errorf("TaskTable with empty slice should return empty string, got %q", result)
	}
}

func TestTaskDetails(t *testing.T) {
	now := time.Now().UTC()
	task := &model.Task{
		ID:          42,
		Title:       "Test Task",
		Type:        "task",
		Status:      "in-progress",
		Description: "Test description",
		Tags:        []string{"tag1", "tag2"},
		CreatedAt:   now,
		UpdatedAt:   now,
		Relationships: []model.Relationship{
			{Type: model.RelParent, TaskID: 10},
		},
	}

	result := TaskDetails(task)

	// Check that result contains expected content
	if !strings.Contains(result, "42") {
		t.Errorf("TaskDetails should contain task ID")
	}
	if !strings.Contains(result, "Test Task") {
		t.Errorf("TaskDetails should contain task title")
	}
	if !strings.Contains(result, "task") {
		t.Errorf("TaskDetails should contain task type")
	}
	if !strings.Contains(result, "in-progress") {
		t.Errorf("TaskDetails should contain task status")
	}
	if !strings.Contains(result, "Test description") {
		t.Errorf("TaskDetails should contain task description")
	}
	if !strings.Contains(result, "tag1") || !strings.Contains(result, "tag2") {
		t.Errorf("TaskDetails should contain task tags")
	}
	if !strings.Contains(result, "parent") {
		t.Errorf("TaskDetails should contain relationship type")
	}
}
