package query

import (
	"context"
	"testing"

	"github.com/opentasks/cmd/internal/model"
	"github.com/opentasks/cmd/internal/storage"
	"github.com/opentasks/cmd/internal/testutil"
)

func setupSQLiteService(t *testing.T) (*SQLiteService, storage.BaseStorage) {
	t.Helper()

	ctx := context.Background()
	sourceStore := storage.NewMemoryStorage()

	svc, err := NewSQLiteService(ctx, sourceStore)
	if err != nil {
		sourceStore.Close()
		t.Fatalf("NewSQLiteService() error = %v", err)
	}

	return svc, sourceStore
}

func TestNewSQLiteService(t *testing.T) {
	svc, store := setupSQLiteService(t)
	defer func() {
		_ = svc.Close()
		_ = store.Close()
	}()

	if svc == nil {
		t.Error("NewSQLiteService() returned nil")
	}

	if svc.db == nil {
		t.Error("NewSQLiteService() db is nil")
	}
}

func TestSQLiteSaveAndLoad(t *testing.T) {
	svc, sourceStore := setupSQLiteService(t)
	defer func() {
		_ = svc.Close()
		_ = sourceStore.Close()
	}()

	ctx := context.Background()
	task := testutil.NewTestTask(1, "Test Task")

	// Save task to SQLite service
	if err := svc.SaveTask(ctx, task); err != nil {
		t.Fatalf("SaveTask() error = %v", err)
	}

	// Load and verify
	loaded, err := svc.LoadTask(ctx, 1)
	if err != nil {
		t.Fatalf("LoadTask() error = %v", err)
	}

	if loaded.ID != task.ID || loaded.Title != task.Title {
		t.Errorf("Loaded task mismatch: got %v, want %v", loaded, task)
	}
}

func TestSQLiteLoadNonexistent(t *testing.T) {
	svc, store := setupSQLiteService(t)
	defer func() {
		_ = svc.Close()
		_ = store.Close()
	}()

	ctx := context.Background()
	_, err := svc.LoadTask(ctx, 999)

	if err != storage.ErrTaskNotFound {
		t.Errorf("LoadTask() error = %v, want ErrTaskNotFound", err)
	}
}

func TestSQLiteDelete(t *testing.T) {
	svc, store := setupSQLiteService(t)
	defer func() {
		_ = svc.Close()
		_ = store.Close()
	}()

	ctx := context.Background()
	task := testutil.NewTestTask(1, "Test Task")

	// Save task
	if err := svc.SaveTask(ctx, task); err != nil {
		t.Fatalf("SaveTask() error = %v", err)
	}

	// Delete task
	if err := svc.DeleteTask(ctx, 1); err != nil {
		t.Fatalf("DeleteTask() error = %v", err)
	}

	// Verify it's gone
	_, err := svc.LoadTask(ctx, 1)
	if err != storage.ErrTaskNotFound {
		t.Errorf("LoadTask() after delete error = %v, want ErrTaskNotFound", err)
	}
}

func TestSQLiteTagNormalization(t *testing.T) {
	svc, store := setupSQLiteService(t)
	defer func() {
		_ = svc.Close()
		_ = store.Close()
	}()

	ctx := context.Background()
	task := testutil.NewTestTaskWithTags(1, "Tagged Task", []string{"feature", "core", "urgent"})

	// Save task
	if err := svc.SaveTask(ctx, task); err != nil {
		t.Fatalf("SaveTask() error = %v", err)
	}

	// Load and verify tags are preserved
	loaded, err := svc.LoadTask(ctx, 1)
	if err != nil {
		t.Fatalf("LoadTask() error = %v", err)
	}

	if len(loaded.Tags) != 3 {
		t.Errorf("Loaded task has %d tags, want 3", len(loaded.Tags))
	}

	// Verify all tags are present
	tagMap := make(map[string]bool)
	for _, tag := range loaded.Tags {
		tagMap[tag] = true
	}

	for _, expectedTag := range []string{"feature", "core", "urgent"} {
		if !tagMap[expectedTag] {
			t.Errorf("Tag %q not found in loaded task tags", expectedTag)
		}
	}
}

func TestSQLiteRelationshipNormalization(t *testing.T) {
	svc, store := setupSQLiteService(t)
	defer func() {
		_ = svc.Close()
		_ = store.Close()
	}()

	ctx := context.Background()

	// Create parent task
	parent := testutil.NewTestTaskWithType(1, "Parent Epic", model.TypeEpic)
	if err := svc.SaveTask(ctx, parent); err != nil {
		t.Fatalf("SaveTask() error = %v", err)
	}

	// Create child task with parent relationship
	child := testutil.NewTestTaskWithType(2, "Child Story", model.TypeStory)
	child.Relationships = []model.Relationship{
		{Type: model.RelParent, TaskID: 1},
	}
	if err := svc.SaveTask(ctx, child); err != nil {
		t.Fatalf("SaveTask() error = %v", err)
	}

	// Load child and verify relationship is preserved
	loaded, err := svc.LoadTask(ctx, 2)
	if err != nil {
		t.Fatalf("LoadTask() error = %v", err)
	}

	if len(loaded.Relationships) != 1 {
		t.Errorf("Loaded task has %d relationships, want 1", len(loaded.Relationships))
	}

	if loaded.Relationships[0].Type != model.RelParent || loaded.Relationships[0].TaskID != 1 {
		t.Errorf("Relationship mismatch: got %v, want RelParent→1", loaded.Relationships[0])
	}
}

func TestSQLiteListTasks(t *testing.T) {
	svc, store := setupSQLiteService(t)
	defer func() {
		_ = svc.Close()
		_ = store.Close()
	}()

	ctx := context.Background()

	// Add test tasks
	tasks := []*model.Task{
		testutil.NewTestTaskWithType(1, "Epic Task", model.TypeEpic),
		testutil.NewTestTaskWithType(2, "Story Task", model.TypeStory),
		testutil.NewTestTaskWithType(3, "Another Story", model.TypeStory),
	}

	for _, task := range tasks {
		if err := svc.SaveTask(ctx, task); err != nil {
			t.Fatalf("SaveTask() error = %v", err)
		}
	}

	// List all tasks
	all, err := svc.ListTasks(ctx)
	if err != nil {
		t.Fatalf("ListTasks() error = %v", err)
	}

	if len(all) != 3 {
		t.Errorf("ListTasks() returned %d tasks, want 3", len(all))
	}
}

func TestSQLiteListTasksWithFilter(t *testing.T) {
	svc, store := setupSQLiteService(t)
	defer func() {
		_ = svc.Close()
		_ = store.Close()
	}()

	ctx := context.Background()

	// Add test tasks
	epic := testutil.NewTestTaskWithType(1, "Epic Task", model.TypeEpic)
	story1 := testutil.NewTestTaskWithType(2, "Story Task", model.TypeStory)
	story2 := testutil.NewTestTaskWithType(3, "Another Story", model.TypeStory)

	for _, task := range []*model.Task{epic, story1, story2} {
		if err := svc.SaveTask(ctx, task); err != nil {
			t.Fatalf("SaveTask() error = %v", err)
		}
	}

	// Filter by type - verify backward compatibility
	stories, err := svc.ListTasks(ctx, func(t *model.Task) bool {
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

func TestSQLiteNextID(t *testing.T) {
	svc, store := setupSQLiteService(t)
	defer func() {
		_ = svc.Close()
		_ = store.Close()
	}()

	ctx := context.Background()

	// Add tasks with different IDs
	for i := 1; i <= 5; i++ {
		task := testutil.NewTestTask(i, "Task "+string(rune(i+'0')))
		if err := svc.SaveTask(ctx, task); err != nil {
			t.Fatalf("SaveTask() error = %v", err)
		}
	}

	// Get next ID
	nextID, err := svc.NextID(ctx)
	if err != nil {
		t.Fatalf("NextID() error = %v", err)
	}

	if nextID != 6 {
		t.Errorf("NextID() = %d, want 6", nextID)
	}
}

func TestSQLiteGetRelated(t *testing.T) {
	svc, store := setupSQLiteService(t)
	defer func() {
		_ = svc.Close()
		_ = store.Close()
	}()

	ctx := context.Background()

	// Create epic
	epic := testutil.NewTestTaskWithType(1, "Epic", model.TypeEpic)
	if err := svc.SaveTask(ctx, epic); err != nil {
		t.Fatalf("SaveTask() error = %v", err)
	}

	// Create stories with parent relationship TO epic (ID 1)
	story1 := testutil.NewTestTaskWithType(2, "Story 1", model.TypeStory)
	story1.Relationships = []model.Relationship{{Type: model.RelParent, TaskID: 1}}
	if err := svc.SaveTask(ctx, story1); err != nil {
		t.Fatalf("SaveTask() error = %v", err)
	}

	story2 := testutil.NewTestTaskWithType(3, "Story 2", model.TypeStory)
	story2.Relationships = []model.Relationship{{Type: model.RelParent, TaskID: 1}}
	if err := svc.SaveTask(ctx, story2); err != nil {
		t.Fatalf("SaveTask() error = %v", err)
	}

	// GetRelated(1, RelParent) should return tasks that have RelParent pointing to task 1
	related, err := svc.GetRelated(ctx, 1, model.RelParent)
	if err != nil {
		t.Fatalf("GetRelated() error = %v", err)
	}

	if len(related) != 2 {
		t.Errorf("GetRelated() returned %d tasks, want 2", len(related))
	}

	// Should return story1 and story2
	ids := make(map[int]bool)
	for _, task := range related {
		ids[task.ID] = true
	}

	if !ids[2] || !ids[3] {
		t.Errorf("GetRelated() returned wrong task IDs: %v, want 2 and 3", ids)
	}
}

func TestSQLiteQueryDAG(t *testing.T) {
	svc, store := setupSQLiteService(t)
	defer func() {
		_ = svc.Close()
		_ = store.Close()
	}()

	ctx := context.Background()

	// Create hierarchy:
	// Epic (1)
	//   ├─ Story (2)
	//   │   ├─ Task (4)
	//   │   └─ Task (5)
	//   └─ Story (3)
	//       └─ Task (6)

	epic := testutil.NewTestTaskWithType(1, "Epic", model.TypeEpic)
	if err := svc.SaveTask(ctx, epic); err != nil {
		t.Fatalf("SaveTask() epic error = %v", err)
	}

	story1 := testutil.NewTestTaskWithType(2, "Story 1", model.TypeStory)
	story1.Relationships = []model.Relationship{{Type: model.RelParent, TaskID: 1}}
	if err := svc.SaveTask(ctx, story1); err != nil {
		t.Fatalf("SaveTask() story1 error = %v", err)
	}

	story2 := testutil.NewTestTaskWithType(3, "Story 2", model.TypeStory)
	story2.Relationships = []model.Relationship{{Type: model.RelParent, TaskID: 1}}
	if err := svc.SaveTask(ctx, story2); err != nil {
		t.Fatalf("SaveTask() story2 error = %v", err)
	}

	task4 := testutil.NewTestTaskWithType(4, "Task 4", model.TypeTask)
	task4.Relationships = []model.Relationship{{Type: model.RelParent, TaskID: 2}}
	if err := svc.SaveTask(ctx, task4); err != nil {
		t.Fatalf("SaveTask() task4 error = %v", err)
	}

	task5 := testutil.NewTestTaskWithType(5, "Task 5", model.TypeTask)
	task5.Relationships = []model.Relationship{{Type: model.RelParent, TaskID: 2}}
	if err := svc.SaveTask(ctx, task5); err != nil {
		t.Fatalf("SaveTask() task5 error = %v", err)
	}

	task6 := testutil.NewTestTaskWithType(6, "Task 6", model.TypeTask)
	task6.Relationships = []model.Relationship{{Type: model.RelParent, TaskID: 3}}
	if err := svc.SaveTask(ctx, task6); err != nil {
		t.Fatalf("SaveTask() task6 error = %v", err)
	}

	// QueryDAG(1) should return stories and their children: 2, 3, 4, 5, 6
	descendants, err := svc.QueryDAG(ctx, 1)
	if err != nil {
		t.Fatalf("QueryDAG() error = %v", err)
	}

	if len(descendants) != 5 {
		t.Errorf("QueryDAG() returned %d tasks, want 5", len(descendants))
	}

	// Verify all expected IDs are present
	ids := make(map[int]bool)
	for _, task := range descendants {
		ids[task.ID] = true
	}

	for _, expectedID := range []int{2, 3, 4, 5, 6} {
		if !ids[expectedID] {
			t.Errorf("QueryDAG() missing task ID %d", expectedID)
		}
	}
}

func TestSQLiteSyncWithSourceStorage(t *testing.T) {
	ctx := context.Background()
	sourceStore := storage.NewMemoryStorage()

	// Pre-populate source storage with tasks
	task1 := testutil.NewTestTask(1, "Task 1")
	task2 := testutil.NewTestTaskWithTags(2, "Task 2", []string{"tag1", "tag2"})
	task3 := testutil.NewTestTaskWithType(3, "Task 3", model.TypeStory)
	task3.Relationships = []model.Relationship{{Type: model.RelParent, TaskID: 1}}

	for _, task := range []*model.Task{task1, task2, task3} {
		if err := sourceStore.SaveTask(ctx, task); err != nil {
			t.Fatalf("SaveTask() error = %v", err)
		}
	}

	// Create SQLite service - should load tasks from source
	svc, err := NewSQLiteService(ctx, sourceStore)
	if err != nil {
		sourceStore.Close()
		t.Fatalf("NewSQLiteService() error = %v", err)
	}
	defer func() {
		_ = svc.Close()
		_ = sourceStore.Close()
	}()

	// Verify all tasks were loaded
	all, err := svc.ListTasks(ctx)
	if err != nil {
		t.Fatalf("ListTasks() error = %v", err)
	}

	if len(all) != 3 {
		t.Errorf("ListTasks() returned %d tasks, want 3", len(all))
	}

	// Verify task2 tags were normalized
	task2Loaded, err := svc.LoadTask(ctx, 2)
	if err != nil {
		t.Fatalf("LoadTask(2) error = %v", err)
	}

	if len(task2Loaded.Tags) != 2 {
		t.Errorf("Task 2 has %d tags, want 2", len(task2Loaded.Tags))
	}

	// Verify task3 relationships were normalized
	task3Loaded, err := svc.LoadTask(ctx, 3)
	if err != nil {
		t.Fatalf("LoadTask(3) error = %v", err)
	}

	if len(task3Loaded.Relationships) != 1 {
		t.Errorf("Task 3 has %d relationships, want 1", len(task3Loaded.Relationships))
	}
}

func TestSQLiteUpdate(t *testing.T) {
	svc, store := setupSQLiteService(t)
	defer func() {
		_ = svc.Close()
		_ = store.Close()
	}()

	ctx := context.Background()
	task := testutil.NewTestTask(1, "Original Title")

	// Save task
	if err := svc.SaveTask(ctx, task); err != nil {
		t.Fatalf("SaveTask() error = %v", err)
	}

	// Update task
	task.Title = "Updated Title"
	task.Status = "in-progress"
	if err := svc.SaveTask(ctx, task); err != nil {
		t.Fatalf("SaveTask() error = %v", err)
	}

	// Load and verify
	loaded, err := svc.LoadTask(ctx, 1)
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

func TestSQLiteMultipleRelationshipTypes(t *testing.T) {
	svc, store := setupSQLiteService(t)
	defer func() {
		_ = svc.Close()
		_ = store.Close()
	}()

	ctx := context.Background()

	// Create tasks with multiple relationship types
	task1 := testutil.NewTestTask(1, "Task 1")
	if err := svc.SaveTask(ctx, task1); err != nil {
		t.Fatalf("SaveTask(1) error = %v", err)
	}

	task2 := testutil.NewTestTask(2, "Task 2")
	if err := svc.SaveTask(ctx, task2); err != nil {
		t.Fatalf("SaveTask(2) error = %v", err)
	}

	// Task 3 has multiple relationships to different tasks
	task3 := testutil.NewTestTask(3, "Task 3")
	task3.Relationships = []model.Relationship{
		{Type: model.RelParent, TaskID: 1},
		{Type: model.RelBlocks, TaskID: 2},
		{Type: model.RelRelatedTo, TaskID: 1},
	}
	if err := svc.SaveTask(ctx, task3); err != nil {
		t.Fatalf("SaveTask(3) error = %v", err)
	}

	// Load and verify all relationships are preserved
	loaded, err := svc.LoadTask(ctx, 3)
	if err != nil {
		t.Fatalf("LoadTask() error = %v", err)
	}

	if len(loaded.Relationships) != 3 {
		t.Errorf("Task has %d relationships, want 3", len(loaded.Relationships))
	}

	// Verify relationship types
	relMap := make(map[string]map[int]bool)
	for _, rel := range loaded.Relationships {
		if relMap[rel.Type] == nil {
			relMap[rel.Type] = make(map[int]bool)
		}
		relMap[rel.Type][rel.TaskID] = true
	}

	if !relMap[model.RelParent][1] {
		t.Error("Missing RelParent → 1")
	}
	if !relMap[model.RelBlocks][2] {
		t.Error("Missing RelBlocks → 2")
	}
	if !relMap[model.RelRelatedTo][1] {
		t.Error("Missing RelRelatedTo → 1")
	}
}

func TestSQLiteClose(t *testing.T) {
	svc, store := setupSQLiteService(t)
	_ = store.Close()

	// Close should always succeed
	err := svc.Close()
	if err != nil {
		t.Errorf("Close() error = %v, want nil", err)
	}
}

func TestSQLiteEmptyListTasks(t *testing.T) {
	svc, store := setupSQLiteService(t)
	defer func() {
		_ = svc.Close()
		_ = store.Close()
	}()

	ctx := context.Background()

	// List from empty service
	tasks, err := svc.ListTasks(ctx)
	if err != nil {
		t.Fatalf("ListTasks() error = %v", err)
	}

	if len(tasks) != 0 {
		t.Errorf("ListTasks() on empty service returned %d tasks, want 0", len(tasks))
	}
}
