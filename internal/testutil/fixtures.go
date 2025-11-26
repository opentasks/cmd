package testutil

import (
	"time"

	"github.com/opentasks/cmd/internal/model"
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
