package query

import (
	"context"
	"testing"
	"time"

	"github.com/zenobi-us/opentasks/internal/model"
	"github.com/zenobi-us/opentasks/internal/storage"
)

func newTestTask(id int, title, taskType, status string) *model.Task {
	now := time.Now().UTC()
	return &model.Task{
		ID:          id,
		Title:       title,
		Type:        taskType,
		Status:      status,
		Description: "Test task",
		Tags:        []string{},
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func setupQueryEngine(ctx context.Context, t *testing.T) (*QueryEngine, storage.BaseStorage) {
	store := storage.NewMemoryStorage()
	qe := NewQueryEngine(store)
	return qe, store
}

func TestNewQueryEngine(t *testing.T) {
	store := storage.NewMemoryStorage()
	defer store.Close()

	qe := NewQueryEngine(store)
	if qe == nil {
		t.Error("NewQueryEngine() returned nil")
	}
}

func TestFindByID(t *testing.T) {
	ctx := context.Background()
	qe, store := setupQueryEngine(ctx, t)
	defer store.Close()

	task := newTestTask(1, "Test Task", model.TypeTask, "todo")
	if err := store.SaveTask(ctx, task); err != nil {
		t.Fatalf("SaveTask() error = %v", err)
	}

	found, err := qe.FindByID(ctx, 1)
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}

	if found.ID != 1 || found.Title != "Test Task" {
		t.Errorf("FindByID() got %v, want task with ID=1 and Title='Test Task'", found)
	}
}

func TestFindByIDNotFound(t *testing.T) {
	ctx := context.Background()
	qe, store := setupQueryEngine(ctx, t)
	defer store.Close()

	_, err := qe.FindByID(ctx, 999)
	if err != storage.ErrTaskNotFound {
		t.Errorf("FindByID() error = %v, want ErrTaskNotFound", err)
	}
}

func TestGetAllTasks(t *testing.T) {
	ctx := context.Background()
	qe, store := setupQueryEngine(ctx, t)
	defer store.Close()

	// Add test tasks
	for i := 1; i <= 5; i++ {
		task := newTestTask(i, "Task "+string(rune(i+'0')), model.TypeTask, "todo")
		if err := store.SaveTask(ctx, task); err != nil {
			t.Fatalf("SaveTask() error = %v", err)
		}
	}

	all, err := qe.GetAllTasks(ctx)
	if err != nil {
		t.Fatalf("GetAllTasks() error = %v", err)
	}

	if len(all) != 5 {
		t.Errorf("GetAllTasks() returned %d tasks, want 5", len(all))
	}
}

func TestGetTasksByType(t *testing.T) {
	ctx := context.Background()
	qe, store := setupQueryEngine(ctx, t)
	defer store.Close()

	// Add tasks of different types
	epic := newTestTask(1, "Epic", model.TypeEpic, "todo")
	story1 := newTestTask(2, "Story 1", model.TypeStory, "todo")
	story2 := newTestTask(3, "Story 2", model.TypeStory, "in-progress")
	task := newTestTask(4, "Task", model.TypeTask, "todo")

	for _, task := range []*model.Task{epic, story1, story2, task} {
		if err := store.SaveTask(ctx, task); err != nil {
			t.Fatalf("SaveTask() error = %v", err)
		}
	}

	// Get stories
	stories, err := qe.GetTasksByType(ctx, model.TypeStory)
	if err != nil {
		t.Fatalf("GetTasksByType() error = %v", err)
	}

	if len(stories) != 2 {
		t.Errorf("GetTasksByType(story) returned %d tasks, want 2", len(stories))
	}

	for _, task := range stories {
		if task.Type != model.TypeStory {
			t.Errorf("GetTasksByType() returned task with type %s, want %s", task.Type, model.TypeStory)
		}
	}
}

func TestGetTasksByStatus(t *testing.T) {
	ctx := context.Background()
	qe, store := setupQueryEngine(ctx, t)
	defer store.Close()

	// Add tasks with different statuses
	task1 := newTestTask(1, "Task 1", model.TypeTask, "todo")
	task2 := newTestTask(2, "Task 2", model.TypeTask, "in-progress")
	task3 := newTestTask(3, "Task 3", model.TypeTask, "in-progress")

	for _, task := range []*model.Task{task1, task2, task3} {
		if err := store.SaveTask(ctx, task); err != nil {
			t.Fatalf("SaveTask() error = %v", err)
		}
	}

	// Get in-progress tasks
	inProgress, err := qe.GetTasksByStatus(ctx, "in-progress")
	if err != nil {
		t.Fatalf("GetTasksByStatus() error = %v", err)
	}

	if len(inProgress) != 2 {
		t.Errorf("GetTasksByStatus('in-progress') returned %d tasks, want 2", len(inProgress))
	}
}

func TestGetTasksByTag(t *testing.T) {
	ctx := context.Background()
	qe, store := setupQueryEngine(ctx, t)
	defer store.Close()

	// Add tasks with different tags
	now := time.Now().UTC()
	task1 := &model.Task{
		ID:        1,
		Title:     "Task 1",
		Type:      model.TypeTask,
		Status:    "todo",
		Tags:      []string{"feature", "core"},
		CreatedAt: now,
		UpdatedAt: now,
	}

	task2 := &model.Task{
		ID:        2,
		Title:     "Task 2",
		Type:      model.TypeTask,
		Status:    "todo",
		Tags:      []string{"feature", "ui"},
		CreatedAt: now,
		UpdatedAt: now,
	}

	task3 := &model.Task{
		ID:        3,
		Title:     "Task 3",
		Type:      model.TypeTask,
		Status:    "todo",
		Tags:      []string{"testing"},
		CreatedAt: now,
		UpdatedAt: now,
	}

	for _, task := range []*model.Task{task1, task2, task3} {
		if err := store.SaveTask(ctx, task); err != nil {
			t.Fatalf("SaveTask() error = %v", err)
		}
	}

	// Get tasks with "feature" tag
	features, err := qe.GetTasksByTag(ctx, "feature")
	if err != nil {
		t.Fatalf("GetTasksByTag() error = %v", err)
	}

	if len(features) != 2 {
		t.Errorf("GetTasksByTag('feature') returned %d tasks, want 2", len(features))
	}
}

func TestListTasksWithMultipleFilters(t *testing.T) {
	ctx := context.Background()
	qe, store := setupQueryEngine(ctx, t)
	defer store.Close()

	// Add tasks
	now := time.Now().UTC()
	task1 := &model.Task{
		ID:        1,
		Title:     "Story 1",
		Type:      model.TypeStory,
		Status:    "in-progress",
		Tags:      []string{"feature"},
		CreatedAt: now,
		UpdatedAt: now,
	}

	task2 := &model.Task{
		ID:        2,
		Title:     "Story 2",
		Type:      model.TypeStory,
		Status:    "todo",
		Tags:      []string{"feature"},
		CreatedAt: now,
		UpdatedAt: now,
	}

	task3 := &model.Task{
		ID:        3,
		Title:     "Task 1",
		Type:      model.TypeTask,
		Status:    "in-progress",
		Tags:      []string{"testing"},
		CreatedAt: now,
		UpdatedAt: now,
	}

	for _, task := range []*model.Task{task1, task2, task3} {
		if err := store.SaveTask(ctx, task); err != nil {
			t.Fatalf("SaveTask() error = %v", err)
		}
	}

	// Query: stories that are in-progress
	results, err := qe.ListTasks(ctx, WithType(model.TypeStory), WithStatus("in-progress"))
	if err != nil {
		t.Fatalf("ListTasks() error = %v", err)
	}

	if len(results) != 1 {
		t.Errorf("ListTasks() with filters returned %d tasks, want 1", len(results))
	}

	if results[0].ID != 1 {
		t.Errorf("ListTasks() returned task %d, want 1", results[0].ID)
	}
}

func TestFindChildren(t *testing.T) {
	ctx := context.Background()
	qe, store := setupQueryEngine(ctx, t)
	defer store.Close()

	// Create parent and children
	epic := newTestTask(1, "Epic", model.TypeEpic, "todo")
	if err := store.SaveTask(ctx, epic); err != nil {
		t.Fatalf("SaveTask() error = %v", err)
	}

	child1 := newTestTask(2, "Child 1", model.TypeStory, "todo")
	child1.Relationships = []model.Relationship{{Type: model.RelParent, TaskID: 1}}

	child2 := newTestTask(3, "Child 2", model.TypeStory, "todo")
	child2.Relationships = []model.Relationship{{Type: model.RelParent, TaskID: 1}}

	for _, child := range []*model.Task{child1, child2} {
		if err := store.SaveTask(ctx, child); err != nil {
			t.Fatalf("SaveTask() error = %v", err)
		}
	}

	// Find children of epic
	children, err := qe.FindChildren(ctx, 1)
	if err != nil {
		t.Fatalf("FindChildren() error = %v", err)
	}

	if len(children) != 2 {
		t.Errorf("FindChildren() returned %d tasks, want 2", len(children))
	}
}

func TestFindBlocking(t *testing.T) {
	ctx := context.Background()
	qe, store := setupQueryEngine(ctx, t)
	defer store.Close()

	// Create epic task
	blocker := newTestTask(1, "Blocker", model.TypeTask, "todo")
	if err := store.SaveTask(ctx, blocker); err != nil {
		t.Fatalf("SaveTask() error = %v", err)
	}

	// Create tasks that are blocked by blocker (ID 1)
	blocked1 := newTestTask(2, "Blocked 1", model.TypeTask, "todo")
	blocked1.Relationships = []model.Relationship{
		{Type: model.RelBlocks, TaskID: 1},
	}

	blocked2 := newTestTask(3, "Blocked 2", model.TypeTask, "todo")
	blocked2.Relationships = []model.Relationship{
		{Type: model.RelBlocks, TaskID: 1},
	}

	if err := store.SaveTask(ctx, blocked1); err != nil {
		t.Fatalf("SaveTask() error = %v", err)
	}

	if err := store.SaveTask(ctx, blocked2); err != nil {
		t.Fatalf("SaveTask() error = %v", err)
	}

	// FindBlocking(1) returns tasks that have RelBlocks pointing to task 1
	// This means tasks that block task 1 (blocked1 and blocked2 block task 1)
	blocking, err := qe.FindBlocking(ctx, 1)
	if err != nil {
		t.Fatalf("FindBlocking() error = %v", err)
	}

	if len(blocking) != 2 {
		t.Errorf("FindBlocking() returned %d tasks, want 2", len(blocking))
	}
}

func TestFindBlockedBy(t *testing.T) {
	ctx := context.Background()
	qe, store := setupQueryEngine(ctx, t)
	defer store.Close()

	// Create blocking relationships
	blocker1 := newTestTask(1, "Blocker 1", model.TypeTask, "todo")
	blocker1.Relationships = []model.Relationship{{Type: model.RelBlocks, TaskID: 3}}

	blocker2 := newTestTask(2, "Blocker 2", model.TypeTask, "todo")
	blocker2.Relationships = []model.Relationship{{Type: model.RelBlocks, TaskID: 3}}

	blocked := newTestTask(3, "Blocked", model.TypeTask, "todo")

	for _, task := range []*model.Task{blocker1, blocker2, blocked} {
		if err := store.SaveTask(ctx, task); err != nil {
			t.Fatalf("SaveTask() error = %v", err)
		}
	}

	// Find blockers
	blockers, err := qe.FindBlockedBy(ctx, 3)
	if err != nil {
		t.Fatalf("FindBlockedBy() error = %v", err)
	}

	if len(blockers) != 2 {
		t.Errorf("FindBlockedBy() returned %d tasks, want 2", len(blockers))
	}
}

func TestFindRelated(t *testing.T) {
	ctx := context.Background()
	qe, store := setupQueryEngine(ctx, t)
	defer store.Close()

	// Create related tasks
	task1 := newTestTask(1, "Task 1", model.TypeTask, "todo")
	if err := store.SaveTask(ctx, task1); err != nil {
		t.Fatalf("SaveTask() error = %v", err)
	}

	task2 := newTestTask(2, "Task 2", model.TypeTask, "todo")
	task2.Relationships = []model.Relationship{{Type: model.RelRelatedTo, TaskID: 1}}

	task3 := newTestTask(3, "Task 3", model.TypeTask, "todo")
	task3.Relationships = []model.Relationship{{Type: model.RelRelatedTo, TaskID: 1}}

	for _, task := range []*model.Task{task2, task3} {
		if err := store.SaveTask(ctx, task); err != nil {
			t.Fatalf("SaveTask() error = %v", err)
		}
	}

	// Find related
	related, err := qe.FindRelated(ctx, 1)
	if err != nil {
		t.Fatalf("FindRelated() error = %v", err)
	}

	if len(related) != 2 {
		t.Errorf("FindRelated() returned %d tasks, want 2", len(related))
	}
}
