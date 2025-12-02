package testutil

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/opentasks/cmd/internal/model"
	"github.com/opentasks/cmd/internal/storage"
)

// NewTestTask creates a task with sensible defaults for testing
func NewTestTask(id int, title string) *model.Task {
	now := time.Now().UTC()
	return &model.Task{
		ID:            id,
		Title:         title,
		Type:          model.TypeTask,
		Status:        "todo",
		Description:   "Test task description",
		Tags:          []string{},
		Relationships: []model.Relationship{},
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

// NewTestTaskWithType creates a task with a specific type
func NewTestTaskWithType(id int, title, taskType string) *model.Task {
	task := NewTestTask(id, title)
	task.Type = taskType
	return task
}

// NewTestTaskWithStatus creates a task with a specific status
func NewTestTaskWithStatus(id int, title, status string) *model.Task {
	task := NewTestTask(id, title)
	task.Status = status
	return task
}

// NewTestTaskWithTags creates a task with specific tags
func NewTestTaskWithTags(id int, title string, tags []string) *model.Task {
	task := NewTestTask(id, title)
	task.Tags = tags
	return task
}

// NewTestTaskWithRelationships creates a task with relationships
func NewTestTaskWithRelationships(id int, title string, relationships []model.Relationship) *model.Task {
	task := NewTestTask(id, title)
	task.Relationships = relationships
	return task
}

// AddRelationship adds a relationship to a task
func AddRelationship(task *model.Task, relType string, targetID int) {
	task.Relationships = append(task.Relationships, model.Relationship{
		Type:   relType,
		TaskID: targetID,
	})
}

// SampleTasks returns a collection of test tasks for use in tests
func SampleTasks() []*model.Task {
	epic := NewTestTaskWithType(1, "Build Task System", model.TypeEpic)
	epic.Status = "in-progress"

	plan := NewTestTaskWithType(2, "Plan Phase 1", model.TypePlan)
	plan.Status = "done"
	AddRelationship(plan, model.RelParent, 1)

	research := NewTestTaskWithType(3, "Research Storage", model.TypeResearch)
	research.Status = "done"
	AddRelationship(research, model.RelParent, 1)

	story1 := NewTestTaskWithType(4, "Implement Core Model", model.TypeStory)
	story1.Status = "done"
	story1.Tags = []string{"feature", "core"}
	AddRelationship(story1, model.RelParent, 1)

	story2 := NewTestTaskWithType(5, "Implement Storage", model.TypeStory)
	story2.Status = "in-progress"
	story2.Tags = []string{"feature", "core"}
	AddRelationship(story2, model.RelParent, 1)

	decision := NewTestTaskWithType(6, "Choose Database", model.TypeDecision)
	decision.Status = "done"
	decision.Tags = []string{"architecture"}
	AddRelationship(decision, model.RelParent, 1)

	task := NewTestTaskWithType(7, "Write Unit Tests", model.TypeTask)
	task.Status = "todo"
	task.Tags = []string{"testing"}
	AddRelationship(task, model.RelBlocks, 5)

	return []*model.Task{epic, plan, research, story1, story2, decision, task}
}

// TaskOption is a functional option for configuring task builders
type TaskOption func(*model.Task)

// WithID sets the task ID
func WithID(id int) TaskOption {
	return func(t *model.Task) {
		t.ID = id
	}
}

// WithTitle sets the task title
func WithTitle(title string) TaskOption {
	return func(t *model.Task) {
		t.Title = title
	}
}

// WithType sets the task type
func WithType(taskType string) TaskOption {
	return func(t *model.Task) {
		t.Type = taskType
	}
}

// WithStatus sets the task status
func WithStatus(status string) TaskOption {
	return func(t *model.Task) {
		t.Status = status
	}
}

// WithParent adds a parent relationship
func WithParent(parentID int) TaskOption {
	return func(t *model.Task) {
		t.Relationships = append(t.Relationships, model.Relationship{
			Type:   model.RelParent,
			TaskID: parentID,
		})
	}
}

// WithBlocks adds a blocks relationship
func WithBlocks(blockedTaskID int) TaskOption {
	return func(t *model.Task) {
		t.Relationships = append(t.Relationships, model.Relationship{
			Type:   model.RelBlocks,
			TaskID: blockedTaskID,
		})
	}
}

// WithTags sets task tags
func WithTags(tags ...string) TaskOption {
	return func(t *model.Task) {
		t.Tags = tags
	}
}

// WithDescription sets task description
func WithDescription(desc string) TaskOption {
	return func(t *model.Task) {
		t.Description = desc
	}
}

// CreateEpic creates an epic and saves it to storage
// Follows the builder pattern from research findings
func CreateEpic(t *testing.T, store storage.BaseStorage, opts ...TaskOption) *model.Task {
	t.Helper()

	ctx := context.Background()

	// Get next ID from storage
	nextID, err := store.NextID(ctx)
	if err != nil {
		t.Fatalf("Failed to get next ID: %v", err)
	}

	// Create epic with defaults
	now := time.Now().UTC()
	epic := &model.Task{
		ID:            nextID,
		Title:         "Test Epic",
		Type:          model.TypeEpic,
		Status:        "planning",
		Description:   "Test epic description",
		Tags:          []string{},
		Relationships: []model.Relationship{},
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	// Apply options
	for _, opt := range opts {
		opt(epic)
	}

	// Save to storage
	if err := store.SaveTask(ctx, epic); err != nil {
		t.Fatalf("Failed to save epic: %v", err)
	}

	return epic
}

// CreateStory creates a story and saves it to storage
func CreateStory(t *testing.T, store storage.BaseStorage, opts ...TaskOption) *model.Task {
	t.Helper()

	ctx := context.Background()

	// Get next ID from storage
	nextID, err := store.NextID(ctx)
	if err != nil {
		t.Fatalf("Failed to get next ID: %v", err)
	}

	// Create story with defaults
	now := time.Now().UTC()
	story := &model.Task{
		ID:            nextID,
		Title:         "Test Story",
		Type:          model.TypeStory,
		Status:        "todo",
		Description:   "Test story description",
		Tags:          []string{},
		Relationships: []model.Relationship{},
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	// Apply options
	for _, opt := range opts {
		opt(story)
	}

	// Save to storage
	if err := store.SaveTask(ctx, story); err != nil {
		t.Fatalf("Failed to save story: %v", err)
	}

	return story
}

// CreateTask creates a task and saves it to storage
func CreateTask(t *testing.T, store storage.BaseStorage, opts ...TaskOption) *model.Task {
	t.Helper()

	ctx := context.Background()

	// Get next ID from storage
	nextID, err := store.NextID(ctx)
	if err != nil {
		t.Fatalf("Failed to get next ID: %v", err)
	}

	// Create task with defaults
	now := time.Now().UTC()
	task := &model.Task{
		ID:            nextID,
		Title:         "Test Task",
		Type:          model.TypeTask,
		Status:        "todo",
		Description:   "Test task description",
		Tags:          []string{},
		Relationships: []model.Relationship{},
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	// Apply options
	for _, opt := range opts {
		opt(task)
	}

	// Save to storage
	if err := store.SaveTask(ctx, task); err != nil {
		t.Fatalf("Failed to save task: %v", err)
	}

	return task
}

// BulkCreateTasks creates N tasks of the same type with pattern-based titles
// Reduces boilerplate for creating multiple similar tasks
//
// Example:
//
//	tasks := BulkCreateTasks(t, store, 3, model.TypeTask,
//	    WithParent(storyID),
//	    WithStatus("todo"),
//	)
//
// Creates: "Test Task 1", "Test Task 2", "Test Task 3"
func BulkCreateTasks(t *testing.T, store storage.BaseStorage, count int, taskType string, opts ...TaskOption) []*model.Task {
	t.Helper()

	if count <= 0 {
		t.Fatal("BulkCreateTasks: count must be > 0")
	}

	ctx := context.Background()
	tasks := make([]*model.Task, count)

	// Default title pattern based on task type
	var titlePattern string
	switch taskType {
	case model.TypeEpic:
		titlePattern = "Test Epic %d"
	case model.TypeStory:
		titlePattern = "Test Story %d"
	case model.TypeTask:
		titlePattern = "Test Task %d"
	case model.TypePlan:
		titlePattern = "Test Plan %d"
	case model.TypeResearch:
		titlePattern = "Test Research %d"
	case model.TypeDecision:
		titlePattern = "Test Decision %d"
	default:
		titlePattern = "Test Task %d"
	}

	for i := 0; i < count; i++ {
		// Get next ID from storage
		nextID, err := store.NextID(ctx)
		if err != nil {
			t.Fatalf("BulkCreateTasks[%d]: Failed to get next ID: %v", i, err)
		}

		// Create task with defaults
		now := time.Now().UTC()
		task := &model.Task{
			ID:            nextID,
			Title:         "", // Will be set after options
			Type:          taskType,
			Status:        "todo",
			Description:   "",
			Tags:          []string{},
			Relationships: []model.Relationship{},
			CreatedAt:     now,
			UpdatedAt:     now,
		}

		// Apply common options first
		for _, opt := range opts {
			opt(task)
		}

		// Apply title pattern if title wasn't set by options
		if task.Title == "" {
			if count > 1 {
				task.Title = fmt.Sprintf(titlePattern, i+1)
			} else {
				task.Title = titlePattern
			}
		}

		// Save to storage
		if err := store.SaveTask(ctx, task); err != nil {
			t.Fatalf("BulkCreateTasks[%d]: Failed to save task: %v", i, err)
		}

		tasks[i] = task
	}

	return tasks
}

// AssertRelationship verifies a task has a specific relationship
// Provides rich error output showing all current relationships
//
// Example:
//
//	testutil.AssertRelationship(t, story, model.RelParent, epic.ID)
func AssertRelationship(t *testing.T, task *model.Task, relType string, targetID int) {
	t.Helper()

	// Search for the relationship
	found := false
	for _, rel := range task.Relationships {
		if rel.Type == relType && rel.TaskID == targetID {
			found = true
			break
		}
	}

	if !found {
		// Build current relationships list
		relList := ""
		if len(task.Relationships) == 0 {
			relList = "  (none)"
		} else {
			for i, rel := range task.Relationships {
				relList += fmt.Sprintf("  %d. %s → Task #%d\n", i+1, rel.Type, rel.TaskID)
			}
		}

		t.Fatalf(`
Relationship Missing
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Expected: %s → Task #%d
Task: #%d %q (%s)

Current Relationships:
%s
Test cannot continue without this relationship.
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━`,
			relType,
			targetID,
			task.ID,
			task.Title,
			task.Type,
			relList)
	}
}

// TransitionTaskState updates task status and verifies persistence
// Loads task, updates status, saves, then reloads to verify
//
// Example:
//
//	task := testutil.TransitionTaskState(t, env.Store, taskID, "in-progress")
func TransitionTaskState(t *testing.T, store storage.BaseStorage, taskID int, newStatus string) *model.Task {
	t.Helper()

	ctx := context.Background()

	// Load task
	task, err := store.LoadTask(ctx, taskID)
	if err != nil {
		t.Fatalf("TransitionTaskState: Failed to load task #%d: %v", taskID, err)
	}

	oldStatus := task.Status

	// Update status
	task.Status = newStatus
	task.UpdatedAt = time.Now().UTC()

	// Save to storage
	if err := store.SaveTask(ctx, task); err != nil {
		t.Fatalf("TransitionTaskState: Failed to save task #%d: %v", taskID, err)
	}

	// Reload to verify persistence
	reloaded, err := store.LoadTask(ctx, taskID)
	if err != nil {
		t.Fatalf("TransitionTaskState: Failed to reload task #%d after save: %v", taskID, err)
	}

	// Verify status was persisted
	if reloaded.Status != newStatus {
		t.Fatalf(`
Status Transition Failed
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Task: #%d %q
Old Status: %q
New Status: %q
Reloaded:   %q (MISMATCH!)

Status was not persisted correctly.
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━`,
			taskID,
			reloaded.Title,
			oldStatus,
			newStatus,
			reloaded.Status)
	}

	return reloaded
}
