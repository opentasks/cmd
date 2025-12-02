package e2e

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/opentasks/cmd/internal/testutil"
)

// TestE2E_SmokeTest validates the basic E2E test infrastructure
// This test ensures that:
// - SetupE2EEnvironment creates isolated environment
// - Fixture builders work correctly
// - Storage persists tasks to filesystem
// - Query engine can retrieve stored tasks
func TestE2E_SmokeTest(t *testing.T) {
	// Setup isolated test environment with real file storage
	env := SetupE2EEnvironment(t)
	defer env.Cleanup()

	t.Run("environment setup creates temp directory", func(t *testing.T) {
		// Verify temp directory exists
		if _, err := os.Stat(env.TmpDir); os.IsNotExist(err) {
			t.Fatalf("Expected temp directory to exist at %s", env.TmpDir)
		}

		// Verify directory is writable
		testFile := filepath.Join(env.TmpDir, "test.txt")
		if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
			t.Fatalf("Temp directory is not writable: %v", err)
		}
	})

	t.Run("create epic and verify storage", func(t *testing.T) {
		// Create epic using fixture builder
		epic := testutil.CreateEpic(t, env.Store,
			testutil.WithTitle("Authentication System"),
			testutil.WithStatus("planning"),
		)

		// Verify epic was created with correct properties
		if epic.ID == 0 {
			t.Fatal("Expected epic to have non-zero ID")
		}
		if epic.Title != "Authentication System" {
			t.Errorf("Expected title 'Authentication System', got '%s'", epic.Title)
		}
		if epic.Status != "planning" {
			t.Errorf("Expected status 'planning', got '%s'", epic.Status)
		}

		// Verify epic was saved to storage
		loaded, err := env.Store.LoadTask(env.Ctx, epic.ID)
		if err != nil {
			t.Fatalf("Failed to load epic from storage: %v", err)
		}
		if loaded.Title != epic.Title {
			t.Errorf("Loaded epic title mismatch: expected '%s', got '%s'",
				epic.Title, loaded.Title)
		}

		// Verify file was created on disk
		files, err := os.ReadDir(env.TmpDir)
		if err != nil {
			t.Fatalf("Failed to read temp directory: %v", err)
		}
		if len(files) == 0 {
			t.Fatal("Expected at least one task file in temp directory")
		}
	})

	t.Run("create story with parent relationship", func(t *testing.T) {
		// Create epic first
		epic := testutil.CreateEpic(t, env.Store,
			testutil.WithTitle("User Management"),
		)

		// Create story linked to epic
		story := testutil.CreateStory(t, env.Store,
			testutil.WithTitle("User Login"),
			testutil.WithParent(epic.ID),
		)

		// Verify story has parent relationship
		if len(story.Relationships) == 0 {
			t.Fatal("Expected story to have parent relationship")
		}
		if story.Relationships[0].TaskID != epic.ID {
			t.Errorf("Expected parent ID %d, got %d",
				epic.ID, story.Relationships[0].TaskID)
		}

		// Verify query engine can find children
		children, err := env.Engine.FindChildren(env.Ctx, epic.ID)
		if err != nil {
			t.Fatalf("Failed to query children: %v", err)
		}
		if len(children) != 1 {
			t.Errorf("Expected 1 child, got %d", len(children))
		}
		if len(children) > 0 && children[0].ID != story.ID {
			t.Errorf("Expected child ID %d, got %d", story.ID, children[0].ID)
		}
	})

	t.Run("memory storage works for fast tests", func(t *testing.T) {
		// Create separate memory environment
		memEnv := SetupMemoryEnvironment(t)
		defer memEnv.Cleanup()

		// Create task in memory
		task := testutil.CreateTask(t, memEnv.Store,
			testutil.WithTitle("Write Tests"),
		)

		// Verify task exists in memory
		loaded, err := memEnv.Store.LoadTask(memEnv.Ctx, task.ID)
		if err != nil {
			t.Fatalf("Failed to load task from memory: %v", err)
		}
		if loaded.Title != task.Title {
			t.Errorf("Memory storage title mismatch")
		}

		// Verify no files were created (memory storage)
		if memEnv.TmpDir != "" {
			t.Error("Expected empty TmpDir for memory storage")
		}
	})
}

// TestE2E_FixtureBuilders validates all fixture builder functions
func TestE2E_FixtureBuilders(t *testing.T) {
	env := SetupMemoryEnvironment(t) // Use memory for speed
	defer env.Cleanup()

	t.Run("CreateEpic with all options", func(t *testing.T) {
		epic := testutil.CreateEpic(t, env.Store,
			testutil.WithTitle("Complex Epic"),
			testutil.WithStatus("active"),
			testutil.WithDescription("Detailed description"),
			testutil.WithTags("feature", "priority-high"),
		)

		if epic.Title != "Complex Epic" {
			t.Errorf("Title mismatch: got '%s'", epic.Title)
		}
		if epic.Status != "active" {
			t.Errorf("Status mismatch: got '%s'", epic.Status)
		}
		if epic.Description != "Detailed description" {
			t.Errorf("Description mismatch: got '%s'", epic.Description)
		}
		if len(epic.Tags) != 2 {
			t.Errorf("Expected 2 tags, got %d", len(epic.Tags))
		}
	})

	t.Run("CreateStory with parent and blocks", func(t *testing.T) {
		epic := testutil.CreateEpic(t, env.Store)
		blockedTask := testutil.CreateTask(t, env.Store)

		story := testutil.CreateStory(t, env.Store,
			testutil.WithTitle("Story with relationships"),
			testutil.WithParent(epic.ID),
			testutil.WithBlocks(blockedTask.ID),
		)

		if len(story.Relationships) != 2 {
			t.Fatalf("Expected 2 relationships, got %d", len(story.Relationships))
		}

		// Verify parent relationship
		hasParent := false
		hasBlocks := false
		for _, rel := range story.Relationships {
			if rel.Type == "parent" && rel.TaskID == epic.ID {
				hasParent = true
			}
			if rel.Type == "blocks" && rel.TaskID == blockedTask.ID {
				hasBlocks = true
			}
		}
		if !hasParent {
			t.Error("Missing parent relationship")
		}
		if !hasBlocks {
			t.Error("Missing blocks relationship")
		}
	})

	t.Run("CreateTask with minimal options", func(t *testing.T) {
		task := testutil.CreateTask(t, env.Store,
			testutil.WithTitle("Simple Task"),
		)

		if task.Title != "Simple Task" {
			t.Errorf("Title mismatch: got '%s'", task.Title)
		}
		if task.Status != "todo" {
			t.Errorf("Expected default status 'todo', got '%s'", task.Status)
		}
	})
}
