package editor

import (
	"os"
	"testing"
)

func TestEditFile(t *testing.T) {
	// Skip this test in CI or when EDITOR is not set
	editor := os.Getenv("EDITOR")
	if editor == "" {
		t.Skip("EDITOR environment variable not set")
	}

	initialContent := "test content\n"

	// This test would require an interactive editor, so we'll test the error handling
	// by using a non-existent editor
	t.Setenv("EDITOR", "/nonexistent/editor")

	_, err := EditFile(initialContent)
	if err == nil {
		t.Errorf("EditFile with non-existent editor should return error")
	}
}

func TestEditFileCreatesTempFile(t *testing.T) {
	// Test that EditFile handles temp file creation
	// We'll use 'cat' which will just output the file content without editing
	t.Setenv("EDITOR", "cat")

	initialContent := "test content\n"
	result, err := EditFile(initialContent)

	// 'cat' just outputs the file, so result should contain the initial content
	if err != nil {
		t.Errorf("EditFile with 'cat' command failed: %v", err)
	}

	if result != initialContent {
		t.Errorf("EditFile returned %q, expected %q", result, initialContent)
	}
}

func TestEditFileEmptyContent(t *testing.T) {
	t.Setenv("EDITOR", "cat")

	result, err := EditFile("")
	if err != nil {
		t.Errorf("EditFile with empty content failed: %v", err)
	}

	if result != "" {
		t.Errorf("EditFile with empty content returned %q, expected empty string", result)
	}
}
