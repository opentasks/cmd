package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestDiscoverConfigFiles tests config file discovery walking up directories
func TestDiscoverConfigFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Create directory structure:
	// tmpDir/
	//   opentask.toml (root)
	//   subdir/
	//     opentask.toml (mid)
	//     deep/
	//       opentask.toml (leaf)

	rootConfig := filepath.Join(tmpDir, "opentask.toml")
	subdir := filepath.Join(tmpDir, "subdir")
	midConfig := filepath.Join(subdir, "opentask.toml")
	deepdir := filepath.Join(subdir, "deep")
	leafConfig := filepath.Join(deepdir, "opentask.toml")

	// Create directories
	if err := os.MkdirAll(deepdir, 0755); err != nil {
		t.Fatalf("Failed to create directories: %v", err)
	}

	// Create config files
	for _, path := range []string{rootConfig, midConfig, leafConfig} {
		if err := os.WriteFile(path, []byte("[project]\nname=\"test\"\n"), 0644); err != nil {
			t.Fatalf("Failed to create config file: %v", err)
		}
	}

	// Test discovery from leaf directory
	found, err := DiscoverConfigFiles(leafConfig)
	if err != nil {
		t.Errorf("DiscoverConfigFiles() error = %v", err)
	}

	if len(found) != 3 {
		t.Errorf("DiscoverConfigFiles() found %d files, want 3", len(found))
	}

	// Verify order: should be [leaf, mid, root] (closest first)
	if len(found) >= 1 && found[0] != leafConfig {
		t.Errorf("First file = %s, want %s", found[0], leafConfig)
	}
	if len(found) >= 2 && found[1] != midConfig {
		t.Errorf("Second file = %s, want %s", found[1], midConfig)
	}
	if len(found) >= 3 && found[2] != rootConfig {
		t.Errorf("Third file = %s, want %s", found[2], rootConfig)
	}
}

// TestDiscoverConfigStopsAtGitRoot tests that discovery stops at .git directory
func TestDiscoverConfigStopsAtGitRoot(t *testing.T) {
	tmpDir := t.TempDir()

	// Create structure:
	// tmpDir/
	//   .git/
	//   opentask.toml (should find this)
	//   subdir/
	//     opentask.toml (should find this)

	gitDir := filepath.Join(tmpDir, ".git")
	rootConfig := filepath.Join(tmpDir, "opentask.toml")
	subdir := filepath.Join(tmpDir, "subdir")
	subConfig := filepath.Join(subdir, "opentask.toml")

	// Create directories and files
	if err := os.MkdirAll(gitDir, 0755); err != nil {
		t.Fatalf("Failed to create .git directory: %v", err)
	}
	if err := os.MkdirAll(subdir, 0755); err != nil {
		t.Fatalf("Failed to create subdir: %v", err)
	}

	for _, path := range []string{rootConfig, subConfig} {
		if err := os.WriteFile(path, []byte("[project]\nname=\"test\"\n"), 0644); err != nil {
			t.Fatalf("Failed to create config: %v", err)
		}
	}

	// Start discovery from subdir
	found, err := DiscoverConfigFiles(subConfig)
	if err != nil {
		t.Errorf("DiscoverConfigFiles() error = %v", err)
	}

	// Should find subConfig and rootConfig, but stop at tmpDir (git root)
	if len(found) < 2 {
		t.Errorf("DiscoverConfigFiles() found %d files, want at least 2", len(found))
	}

	// Verify we found both configs
	hasRoot := false
	hasSub := false
	for _, f := range found {
		if f == rootConfig {
			hasRoot = true
		}
		if f == subConfig {
			hasSub = true
		}
	}
	if !hasRoot || !hasSub {
		t.Errorf("Missing expected config files. hasRoot=%v, hasSub=%v", hasRoot, hasSub)
	}

	// Verify we didn't go beyond git root
	for _, f := range found {
		if !strings.HasPrefix(f, tmpDir) {
			t.Errorf("Found config outside git root: %s", f)
		}
	}
}

// TestMergeConfigs tests merging multiple config files
func TestMergeConfigs(t *testing.T) {
	tmpDir := t.TempDir()

	// Create base config
	baseConfig := filepath.Join(tmpDir, "base.toml")
	baseContent := `
[project]
name = "Base Project"
owner = "base-owner"

[workflow]
statuses = ["todo", "done"]
initial = "todo"

[storage]
backend = "markdown-fs"
path = "./tasks"
`
	if err := os.WriteFile(baseConfig, []byte(baseContent), 0644); err != nil {
		t.Fatalf("Failed to create base config: %v", err)
	}

	// Create override config
	overrideConfig := filepath.Join(tmpDir, "override.toml")
	overrideContent := `
[project]
name = "Override Project"

[workflow]
statuses = ["todo", "in-progress", "done"]

[storage]
path = "./override-tasks"
`
	if err := os.WriteFile(overrideConfig, []byte(overrideContent), 0644); err != nil {
		t.Fatalf("Failed to create override config: %v", err)
	}

	// Merge configs: baseConfig first, then overrideConfig
	merged, err := MergeConfigs([]string{baseConfig, overrideConfig})
	if err != nil {
		t.Errorf("MergeConfigs() error = %v", err)
	}

	// Override should win for name
	if merged.Project.Name != "Override Project" {
		t.Errorf("Project name = %q, want 'Override Project' (override should win)", merged.Project.Name)
	}

	// But base values should be preserved if not overridden
	if merged.Project.Owner != "base-owner" {
		t.Errorf("Project owner = %q, want 'base-owner' (should be preserved from base)", merged.Project.Owner)
	}

	// Override should win for statuses
	if len(merged.Workflow.Statuses) != 3 {
		t.Errorf("Workflow statuses length = %d, want 3", len(merged.Workflow.Statuses))
	}

	// Storage path should be from override
	if !strings.HasSuffix(merged.Storage.Path, "override-tasks") {
		t.Errorf("Storage path = %q, want suffix 'override-tasks'", merged.Storage.Path)
	}

	// Backend should be preserved from base if not in override
	if merged.Storage.Backend != "markdown-fs" {
		t.Errorf("Storage backend = %q, want 'markdown-fs'", merged.Storage.Backend)
	}
}

// TestMergeConfigsEmpty tests merging with default when no configs found
func TestMergeConfigsEmpty(t *testing.T) {
	merged, err := MergeConfigs([]string{})
	if err != nil {
		t.Errorf("MergeConfigs(empty) error = %v", err)
	}

	if merged == nil {
		t.Error("MergeConfigs(empty) returned nil")
	}

	// Should have defaults
	if len(merged.Workflow.Statuses) == 0 {
		t.Error("MergeConfigs(empty) should have default workflow")
	}
}
