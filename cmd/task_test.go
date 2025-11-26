package cmd

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/opentasks/cmd/internal/model"
	"github.com/opentasks/cmd/internal/query"
	"github.com/opentasks/cmd/internal/storage"
)

// setupTestEnvironment creates a test storage and engine with a sample task
func setupTestEnvironment(t *testing.T) (storage.BaseStorage, *query.QueryEngine, string) {
	tmpDir := t.TempDir()

	storageConfig := storage.StorageConfig{
		Backend: "markdown-fs",
		Path:    tmpDir,
	}

	store, err := storage.NewStorage(storageConfig)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	engine := query.NewQueryEngine(store)

	// Create a test task
	ctx := context.Background()
	testTask := &model.Task{
		ID:          1,
		Title:       "Test Task",
		Type:        "story",
		Status:      "todo",
		Tags:        []string{"test", "urgent"},
		Description: "This is a test task",
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	if err := store.SaveTask(ctx, testTask); err != nil {
		t.Fatalf("Failed to save test task: %v", err)
	}

	return store, engine, tmpDir
}

// TestTaskUpdateTitleChangesFilename verifies that updating a title creates a new file
func TestTaskUpdateTitleChangesFilename(t *testing.T) {
	store, _, tmpDir := setupTestEnvironment(t)
	defer func() { _ = store.Close() }()

	ctx := context.Background()

	// Get original file
	task, err := store.LoadTask(ctx, 1)
	if err != nil {
		t.Fatalf("Failed to load task: %v", err)
	}

	originalFiles, _ := filepath.Glob(filepath.Join(tmpDir, "s-1-*.md"))
	if len(originalFiles) != 1 {
		t.Fatalf("Expected 1 original file, got %d", len(originalFiles))
	}

	// Update title
	task.Title = "Updated Title"
	task.UpdatedAt = time.Now().UTC()

	if err := store.SaveTask(ctx, task); err != nil {
		t.Fatalf("SaveTask failed: %v", err)
	}

	// Verify old file is gone and new file exists
	allFiles, _ := filepath.Glob(filepath.Join(tmpDir, "s-1-*.md"))
	if len(allFiles) != 1 {
		t.Errorf("Expected 1 file after update, got %d: %v", len(allFiles), allFiles)
	}

	// Verify the new filename is correct
	expectedFile := filepath.Join(tmpDir, "s-1-updated-title.md")
	if allFiles[0] != expectedFile {
		t.Errorf("Expected file %s, got %s", expectedFile, allFiles[0])
	}
}

// TestTaskLoadAfterTitleUpdate verifies data persists correctly after title update
func TestTaskLoadAfterTitleUpdate(t *testing.T) {
	store, _, _ := setupTestEnvironment(t)
	defer func() { _ = store.Close() }()

	ctx := context.Background()

	// Load and update task
	task, _ := store.LoadTask(ctx, 1)
	task.Title = "New Title"
	task.Status = "in-progress"
	task.Tags = []string{"updated", "test"}
	task.UpdatedAt = time.Now().UTC()

	if err := store.SaveTask(ctx, task); err != nil {
		t.Fatalf("SaveTask failed: %v", err)
	}

	// Load again and verify
	loaded, err := store.LoadTask(ctx, 1)
	if err != nil {
		t.Fatalf("LoadTask failed: %v", err)
	}

	if loaded.Title != "New Title" {
		t.Errorf("Title: got %q, want %q", loaded.Title, "New Title")
	}

	if loaded.Status != "in-progress" {
		t.Errorf("Status: got %q, want %q", loaded.Status, "in-progress")
	}

	if len(loaded.Tags) != 2 {
		t.Errorf("Tags: expected 2, got %d", len(loaded.Tags))
	}
}

// TestTaskAddTagsPreservesExisting verifies that adding tags doesn't lose existing ones
func TestTaskAddTagsPreservesExisting(t *testing.T) {
	store, _, _ := setupTestEnvironment(t)
	defer func() { _ = store.Close() }()

	ctx := context.Background()

	task, _ := store.LoadTask(ctx, 1)
	originalTags := len(task.Tags)

	// Add new tag
	task.Tags = append(task.Tags, "new-tag")
	task.UpdatedAt = time.Now().UTC()

	if err := store.SaveTask(ctx, task); err != nil {
		t.Fatalf("SaveTask failed: %v", err)
	}

	// Load and verify
	loaded, _ := store.LoadTask(ctx, 1)
	if len(loaded.Tags) != originalTags+1 {
		t.Errorf("Expected %d tags, got %d", originalTags+1, len(loaded.Tags))
	}

	hasNewTag := false
	for _, tag := range loaded.Tags {
		if tag == "new-tag" {
			hasNewTag = true
			break
		}
	}

	if !hasNewTag {
		t.Errorf("New tag not found in: %v", loaded.Tags)
	}
}

// TestTaskRemoveTagsUpdatesCorrectly verifies tag removal works
func TestTaskRemoveTagsUpdatesCorrectly(t *testing.T) {
	store, _, _ := setupTestEnvironment(t)
	defer func() { _ = store.Close() }()

	ctx := context.Background()

	task, _ := store.LoadTask(ctx, 1)

	// Remove a tag
	newTags := []string{}
	for _, tag := range task.Tags {
		if tag != "urgent" {
			newTags = append(newTags, tag)
		}
	}
	task.Tags = newTags
	task.UpdatedAt = time.Now().UTC()

	if err := store.SaveTask(ctx, task); err != nil {
		t.Fatalf("SaveTask failed: %v", err)
	}

	// Load and verify
	loaded, _ := store.LoadTask(ctx, 1)

	hasUrgent := false
	for _, tag := range loaded.Tags {
		if tag == "urgent" {
			hasUrgent = true
			break
		}
	}

	if hasUrgent {
		t.Errorf("Urgent tag should have been removed: %v", loaded.Tags)
	}

	// Should still have "test" tag
	hasTest := false
	for _, tag := range loaded.Tags {
		if tag == "test" {
			hasTest = true
			break
		}
	}

	if !hasTest {
		t.Errorf("Test tag should still exist: %v", loaded.Tags)
	}
}
