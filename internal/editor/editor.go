package editor

import (
	"fmt"
	"os"
	"os/exec"
)

// EditFile opens content in the user's preferred editor and returns the edited content
func EditFile(initialContent string) (string, error) {
	// Get the editor from environment variable
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi" // Fallback to vi if EDITOR is not set
	}

	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "opentask-*.md")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpFilePath := tmpFile.Name()
	defer func() { _ = os.Remove(tmpFilePath) }()

	// Write initial content to the temp file
	if _, err := tmpFile.WriteString(initialContent); err != nil {
		_ = tmpFile.Close()
		return "", fmt.Errorf("failed to write temp file: %w", err)
	}
	_ = tmpFile.Close()

	// Launch the editor
	cmd := exec.Command(editor, tmpFilePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("editor command failed: %w", err)
	}

	// Read the updated content
	updatedContent, err := os.ReadFile(tmpFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read updated content: %w", err)
	}

	return string(updatedContent), nil
}
