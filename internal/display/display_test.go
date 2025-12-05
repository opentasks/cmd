package display

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/opentasks/cmd/internal/model"
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

// TestPrintOnboardingBox_BasicOutput tests basic output structure and content
func TestPrintOnboardingBox_BasicOutput(t *testing.T) {
	var buf bytes.Buffer
	cwd := "/home/user/my-project"

	PrintOnboardingBox(&buf, cwd)
	output := buf.String()

	// Check that output contains expected content
	if !strings.Contains(output, "NO OPENTASK PROJECT FOUND") {
		t.Errorf("PrintOnboardingBox should contain header")
	}
	if !strings.Contains(output, "opentask config init") {
		t.Errorf("PrintOnboardingBox should contain opentask config init command")
	}
	if !strings.Contains(strings.ToLower(output), "project") {
		t.Errorf("PrintOnboardingBox should mention project")
	}
}

// TestPrintOnboardingBox_PathTruncation tests long path truncation
func TestPrintOnboardingBox_PathTruncation(t *testing.T) {
	var buf bytes.Buffer
	// Very long path that should be truncated
	cwd := "/very/long/path/that/exceeds/the/maximum/path/length/displayed/in/onboarding/message/to/test/truncation"

	PrintOnboardingBox(&buf, cwd)
	output := buf.String()

	// Should contain the output and handle truncation gracefully
	if output == "" {
		t.Errorf("PrintOnboardingBox should produce output even with long paths")
	}
	// Should still contain the directory structure hint
	if !strings.Contains(output, "opentask") {
		t.Errorf("PrintOnboardingBox should mention opentask even with long paths")
	}
}

// TestPrintOnboardingBox_EmptyPath tests empty path handling
func TestPrintOnboardingBox_EmptyPath(t *testing.T) {
	var buf bytes.Buffer
	cwd := ""

	PrintOnboardingBox(&buf, cwd)
	output := buf.String()

	// Should produce output
	if output == "" {
		t.Errorf("PrintOnboardingBox should produce output even with empty path")
	}
	// Should still have the main message
	if !strings.Contains(output, "NO OPENTASK PROJECT FOUND") {
		t.Errorf("PrintOnboardingBox should contain main message with empty path")
	}
}

// TestPrintOnboardingBox_UnknownDirectory tests unknown directory handling
func TestPrintOnboardingBox_UnknownDirectory(t *testing.T) {
	var buf bytes.Buffer
	cwd := "(unknown directory)"

	PrintOnboardingBox(&buf, cwd)
	output := buf.String()

	// Should produce output
	if output == "" {
		t.Errorf("PrintOnboardingBox should produce output for unknown directory")
	}
	// Should still have the instructions
	if !strings.Contains(output, "opentask") {
		t.Errorf("PrintOnboardingBox should contain opentask reference for unknown directory")
	}
}

// TestPrintOnboardingBox_OutputNotEmpty tests that output is never empty and contains formatting
func TestPrintOnboardingBox_OutputNotEmpty(t *testing.T) {
	tests := []struct {
		name string
		cwd  string
	}{
		{"root directory", "/"},
		{"home directory", "/home/user"},
		{"deeply nested", "/a/b/c/d/e/f/g/h/i/j/k/l/m/n/o/p"},
		{"relative path", "relative/path"},
		{"single character", "x"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			PrintOnboardingBox(&buf, tt.cwd)
			output := buf.String()

			if output == "" {
				t.Errorf("PrintOnboardingBox should produce output for %q", tt.cwd)
			}
			// Verify key messages are present
			if !strings.Contains(strings.ToLower(output), "project") && !strings.Contains(strings.ToLower(output), "opentask") {
				t.Errorf("PrintOnboardingBox output should be meaningful for %q", tt.cwd)
			}
		})
	}
}
