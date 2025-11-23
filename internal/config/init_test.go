package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestConfigInitializer_Initialize_CreateFile(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()

	// Initialize a new config
	initializer := NewConfigInitializer(tmpDir)
	err := initializer.Initialize("TestProject", "./.tasks", false)
	if err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	// Verify file was created
	configPath := filepath.Join(tmpDir, ".opentask.toml")
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read created config file: %v", err)
	}

	// Verify content contains expected values
	contentStr := string(content)
	if !strings.Contains(contentStr, "TestProject") {
		t.Error("Config does not contain project name")
	}
	if !strings.Contains(contentStr, "./.tasks") {
		t.Error("Config does not contain storage path")
	}
	if !strings.Contains(contentStr, "[project]") {
		t.Error("Config does not contain [project] section")
	}
	if !strings.Contains(contentStr, "[storage]") {
		t.Error("Config does not contain [storage] section")
	}
}

func TestConfigInitializer_Initialize_ExistingFile_NoForce(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()

	// Create an existing config file
	configPath := filepath.Join(tmpDir, ".opentask.toml")
	if err := os.WriteFile(configPath, []byte("existing"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Try to initialize without force
	initializer := NewConfigInitializer(tmpDir)
	err := initializer.Initialize("TestProject", "./.tasks", false)

	// Should return error
	if err == nil {
		t.Error("Expected error when config already exists without --force")
	}
	if !strings.Contains(err.Error(), ".opentask.toml already exists") {
		t.Errorf("Expected 'already exists' error, got: %v", err)
	}
}

func TestConfigInitializer_Initialize_ExistingFile_Force(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()

	// Create an existing config file
	configPath := filepath.Join(tmpDir, ".opentask.toml")
	if err := os.WriteFile(configPath, []byte("existing"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Initialize with force flag
	initializer := NewConfigInitializer(tmpDir)
	err := initializer.Initialize("NewProject", "./.tasks", true)
	if err != nil {
		t.Fatalf("Initialize with force failed: %v", err)
	}

	// Verify file was overwritten
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	contentStr := string(content)
	if strings.Contains(contentStr, "existing") {
		t.Error("Old content was not overwritten")
	}
	if !strings.Contains(contentStr, "NewProject") {
		t.Error("New project name not in config")
	}
}

func TestConfigInitializer_Initialize_DefaultProjectName(t *testing.T) {
	// Create temporary directory with a specific name
	tmpDir := t.TempDir()
	dirName := filepath.Base(tmpDir)

	// Initialize without specifying project name
	initializer := NewConfigInitializer(tmpDir)
	err := initializer.Initialize("", "./.tasks", false)
	if err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	// Verify file uses directory name
	configPath := filepath.Join(tmpDir, ".opentask.toml")
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, dirName) {
		t.Errorf("Expected directory name %q in config, but not found", dirName)
	}
}

func TestConfigInitializer_renderTemplate(t *testing.T) {
	initializer := NewConfigInitializer("")

	data := map[string]string{
		"ProjectName": "MyProject",
		"StoragePath": "/path/to/storage",
	}

	result, err := initializer.renderTemplate(data)
	if err != nil {
		t.Fatalf("renderTemplate failed: %v", err)
	}

	// Verify template was rendered with correct values
	if !strings.Contains(result, "MyProject") {
		t.Error("ProjectName not rendered in template")
	}
	if !strings.Contains(result, "/path/to/storage") {
		t.Error("StoragePath not rendered in template")
	}
	if !strings.Contains(result, "[project]") {
		t.Error("Template structure missing [project] section")
	}
	if !strings.Contains(result, "[storage]") {
		t.Error("Template structure missing [storage] section")
	}
}

func TestConfigInitializer_renderTemplate_WithSpecialChars(t *testing.T) {
	initializer := NewConfigInitializer("")

	data := map[string]string{
		"ProjectName": "My-Project_123",
		"StoragePath": "./tasks-dir/storage",
	}

	result, err := initializer.renderTemplate(data)
	if err != nil {
		t.Fatalf("renderTemplate failed: %v", err)
	}

	// Verify special characters are handled correctly
	if !strings.Contains(result, "My-Project_123") {
		t.Error("Special characters in ProjectName not preserved")
	}
	if !strings.Contains(result, "./tasks-dir/storage") {
		t.Error("Special characters in StoragePath not preserved")
	}
}
