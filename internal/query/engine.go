package query

import (
	"context"

	"github.com/opentasks/cmd/internal/model"
	"github.com/opentasks/cmd/internal/storage"
)

// QueryEngine provides a convenient query interface on top of storage
type QueryEngine struct {
	storage storage.BaseStorage
}

// NewQueryEngine creates a new query engine
func NewQueryEngine(storage storage.BaseStorage) *QueryEngine {
	return &QueryEngine{
		storage: storage,
	}
}

// FindByID finds a task by ID
func (q *QueryEngine) FindByID(ctx context.Context, id int) (*model.Task, error) {
	return q.storage.LoadTask(ctx, id)
}

// FindChildren finds all children of a parent task
func (q *QueryEngine) FindChildren(ctx context.Context, parentID int) ([]*model.Task, error) {
	return q.storage.GetRelated(ctx, parentID, model.RelParent)
}

// FindBlocking finds all tasks that are blocked by the given task
func (q *QueryEngine) FindBlocking(ctx context.Context, taskID int) ([]*model.Task, error) {
	return q.storage.GetRelated(ctx, taskID, model.RelBlocks)
}

// FindBlockedBy finds all tasks that block the given task
// This searches for all tasks that have a "blocks" relationship pointing to this task
func (q *QueryEngine) FindBlockedBy(ctx context.Context, taskID int) ([]*model.Task, error) {
	allTasks, err := q.storage.ListTasks(ctx)
	if err != nil {
		return nil, err
	}

	var blockers []*model.Task
	for _, task := range allTasks {
		for _, rel := range task.Relationships {
			if rel.Type == model.RelBlocks && rel.TaskID == taskID {
				blockers = append(blockers, task)
				break
			}
		}
	}

	return blockers, nil
}

// FindRelated finds all tasks related to the given task
func (q *QueryEngine) FindRelated(ctx context.Context, taskID int) ([]*model.Task, error) {
	return q.storage.GetRelated(ctx, taskID, model.RelRelatedTo)
}

// ListTasks returns all tasks matching the given filters
func (q *QueryEngine) ListTasks(ctx context.Context, filters ...TaskFilter) ([]*model.Task, error) {
	return q.storage.ListTasks(ctx, filters...)
}

// GetAllTasks returns all tasks in the project
func (q *QueryEngine) GetAllTasks(ctx context.Context) ([]*model.Task, error) {
	return q.storage.ListTasks(ctx)
}

// GetTasksByType returns all tasks of the given type
func (q *QueryEngine) GetTasksByType(ctx context.Context, taskType string) ([]*model.Task, error) {
	return q.storage.ListTasks(ctx, WithType(taskType))
}

// GetTasksByStatus returns all tasks with the given status
func (q *QueryEngine) GetTasksByStatus(ctx context.Context, status string) ([]*model.Task, error) {
	return q.storage.ListTasks(ctx, WithStatus(status))
}

// GetTasksByTag returns all tasks with the given tag
func (q *QueryEngine) GetTasksByTag(ctx context.Context, tag string) ([]*model.Task, error) {
	return q.storage.ListTasks(ctx, WithTag(tag))
}

// NextID generates the next task ID
func (q *QueryEngine) NextID(ctx context.Context) (int, error) {
	return q.storage.NextID(ctx)
}
