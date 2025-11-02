package storage

import (
	"context"
	"testing"
	"time"

	"github.com/zenobi-us/opentask/internal/model"
)

func newTestTask(id int, title, taskType string) *model.Task {
	now := time.Now().UTC()
	return &model.Task{
		ID:          id,
		Title:       title,
		Type:        taskType,
		Status:      "todo",
		Description: "Test task",
		Tags:        []string{},
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func TestMemoryStorageSaveAndLoad(t *testing.T) {
	store := NewMemoryStorage()
	defer store.Close()

	ctx := context.Background()
	task := newTestTask(1, "Test Task", model.TypeTask)

	// Save task
	if err := store.SaveTask(ctx, task); err != nil {
		t.Fatalf("SaveTask() error = %v, want nil", err)
	}

	// Load task
	loaded, err := store.LoadTask(ctx, 1)
	if err != nil {
		t.Fatalf("LoadTask() error = %v, want nil", err)
	}

	if loaded.ID != task.ID || loaded.Title != task.Title {
		t.Errorf("Loaded task mismatch: got %v, want %v", loaded, task)
	}
}

func TestMemoryStorageLoadNonexistent(t *testing.T) {
	store := NewMemoryStorage()
	defer store.Close()

	ctx := context.Background()
	_, err := store.LoadTask(ctx, 999)

	if err != ErrTaskNotFound {
		t.Errorf("LoadTask() error = %v, want ErrTaskNotFound", err)
	}
}

func TestMemoryStorageDelete(t *testing.T) {
	store := NewMemoryStorage()
	defer store.Close()

	ctx := context.Background()
	task := newTestTask(1, "Test Task", model.TypeTask)

	// Save task
	if err := store.SaveTask(ctx, task); err != nil {
		t.Fatalf("SaveTask() error = %v", err)
	}

	// Delete task
	if err := store.DeleteTask(ctx, 1); err != nil {
		t.Fatalf("DeleteTask() error = %v", err)
	}

	// Verify it's gone
	_, err := store.LoadTask(ctx, 1)
	if err != ErrTaskNotFound {
		t.Errorf("LoadTask() after delete error = %v, want ErrTaskNotFound", err)
	}
}

func TestMemoryStorageDeleteNonexistent(t *testing.T) {
	store := NewMemoryStorage()
	defer store.Close()

	ctx := context.Background()
	err := store.DeleteTask(ctx, 999)

	if err != ErrTaskNotFound {
		t.Errorf("DeleteTask() error = %v, want ErrTaskNotFound", err)
	}
}

func TestMemoryStorageListTasks(t *testing.T) {
	store := NewMemoryStorage()
	defer store.Close()

	ctx := context.Background()

	// Add test tasks
	tasks := []*model.Task{
		newTestTask(1, "Epic Task", model.TypeEpic),
		newTestTask(2, "Story Task", model.TypeStory),
		newTestTask(3, "Another Story", model.TypeStory),
	}

	for _, task := range tasks {
		if err := store.SaveTask(ctx, task); err != nil {
			t.Fatalf("SaveTask() error = %v", err)
		}
	}

	// List all tasks
	all, err := store.ListTasks(ctx)
	if err != nil {
		t.Fatalf("ListTasks() error = %v", err)
	}

	if len(all) != 3 {
		t.Errorf("ListTasks() returned %d tasks, want 3", len(all))
	}
}

func TestMemoryStorageListTasksWithFilter(t *testing.T) {
	store := NewMemoryStorage()
	defer store.Close()

	ctx := context.Background()

	// Add test tasks
	epic := newTestTask(1, "Epic Task", model.TypeEpic)
	story1 := newTestTask(2, "Story Task", model.TypeStory)
	story2 := newTestTask(3, "Another Story", model.TypeStory)

	for _, task := range []*model.Task{epic, story1, story2} {
		if err := store.SaveTask(ctx, task); err != nil {
			t.Fatalf("SaveTask() error = %v", err)
		}
	}

	// Filter by type
	stories, err := store.ListTasks(ctx, func(t *model.Task) bool {
		return t.Type == model.TypeStory
	})
	if err != nil {
		t.Fatalf("ListTasks() error = %v", err)
	}

	if len(stories) != 2 {
		t.Errorf("ListTasks() with filter returned %d tasks, want 2", len(stories))
	}

	for _, task := range stories {
		if task.Type != model.TypeStory {
			t.Errorf("Filter returned task with type %s, want %s", task.Type, model.TypeStory)
		}
	}
}

func TestMemoryStorageNextID(t *testing.T) {
	store := NewMemoryStorage()
	defer store.Close()

	ctx := context.Background()

	// Add tasks with different IDs
	for i := 1; i <= 5; i++ {
		task := newTestTask(i, "Task "+string(rune(i+'0')), model.TypeTask)
		if err := store.SaveTask(ctx, task); err != nil {
			t.Fatalf("SaveTask() error = %v", err)
		}
	}

	// Get next ID
	nextID, err := store.NextID(ctx)
	if err != nil {
		t.Fatalf("NextID() error = %v", err)
	}

	if nextID != 6 {
		t.Errorf("NextID() = %d, want 6", nextID)
	}
}

func TestMemoryStorageUpdate(t *testing.T) {
	store := NewMemoryStorage()
	defer store.Close()

	ctx := context.Background()
	task := newTestTask(1, "Original Title", model.TypeTask)

	// Save task
	if err := store.SaveTask(ctx, task); err != nil {
		t.Fatalf("SaveTask() error = %v", err)
	}

	// Update task
	task.Title = "Updated Title"
	task.Status = "in-progress"
	if err := store.SaveTask(ctx, task); err != nil {
		t.Fatalf("SaveTask() error = %v", err)
	}

	// Load and verify
	loaded, err := store.LoadTask(ctx, 1)
	if err != nil {
		t.Fatalf("LoadTask() error = %v", err)
	}

	if loaded.Title != "Updated Title" {
		t.Errorf("Update failed: Title = %q, want 'Updated Title'", loaded.Title)
	}

	if loaded.Status != "in-progress" {
		t.Errorf("Update failed: Status = %q, want 'in-progress'", loaded.Status)
	}
}

func TestMemoryStorageGetRelated(t *testing.T) {
	store := NewMemoryStorage()
	defer store.Close()

	ctx := context.Background()

	// Create tasks with relationships
	epic := newTestTask(1, "Epic", model.TypeEpic)
	if err := store.SaveTask(ctx, epic); err != nil {
		t.Fatalf("SaveTask() error = %v", err)
	}

	// Create story with parent relationship TO epic (ID 1)
	story := newTestTask(2, "Story", model.TypeStory)
	story.Relationships = []model.Relationship{
		{Type: model.RelParent, TaskID: 1},
	}
	if err := store.SaveTask(ctx, story); err != nil {
		t.Fatalf("SaveTask() error = %v", err)
	}

	// Create another story with parent relationship TO epic (ID 1)
	story2 := newTestTask(3, "Story 2", model.TypeStory)
	story2.Relationships = []model.Relationship{
		{Type: model.RelParent, TaskID: 1},
	}
	if err := store.SaveTask(ctx, story2); err != nil {
		t.Fatalf("SaveTask() error = %v", err)
	}

	// GetRelated(1, RelParent) should return tasks that have RelParent pointing to task 1
	related, err := store.GetRelated(ctx, 1, model.RelParent)
	if err != nil {
		t.Fatalf("GetRelated() error = %v", err)
	}

	if len(related) != 2 {
		t.Errorf("GetRelated() returned %d tasks, want 2", len(related))
	}

	// Should return story and story2
	ids := make(map[int]bool)
	for _, task := range related {
		ids[task.ID] = true
	}

	if !ids[2] || !ids[3] {
		t.Errorf("GetRelated() returned wrong task IDs: %v, want 2 and 3", ids)
	}
}

func TestMemoryStorageEmptyListTasks(t *testing.T) {
	store := NewMemoryStorage()
	defer store.Close()

	ctx := context.Background()

	// List from empty store
	tasks, err := store.ListTasks(ctx)
	if err != nil {
		t.Fatalf("ListTasks() error = %v", err)
	}

	if len(tasks) != 0 {
		t.Errorf("ListTasks() on empty store returned %d tasks, want 0", len(tasks))
	}
}

func TestMemoryStorageClose(t *testing.T) {
	store := NewMemoryStorage()

	// Close should always succeed
	err := store.Close()
	if err != nil {
		t.Errorf("Close() error = %v, want nil", err)
	}
}
