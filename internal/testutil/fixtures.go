package testutil

import (
	"context"
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
