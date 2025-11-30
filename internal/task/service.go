package task

import (
	"context"
	"fmt"
	"time"

	"github.com/opentasks/cmd/internal/model"
	"github.com/opentasks/cmd/internal/query"
	"github.com/opentasks/cmd/internal/storage"
)

// TaskEngine defines the query operations needed by the service
type TaskEngine interface {
	NextID(ctx context.Context) (int, error)
	FindByID(ctx context.Context, id int) (*model.Task, error)
	ListTasks(ctx context.Context, filters ...query.TaskFilter) ([]*model.Task, error)
}

// Service provides high-level task operations
// Coordinates between query engine, storage, and task manager
type Service struct {
	engine  TaskEngine
	store   storage.BaseStorage
	manager *Manager
}

// NewService creates a new task service
func NewService(engine TaskEngine, store storage.BaseStorage) *Service {
	return &Service{
		engine:  engine,
		store:   store,
		manager: NewManager(),
	}
}

// CreateRequest contains parameters for creating a new task
type CreateRequest struct {
	Title       string
	Type        string
	Status      string // Optional - auto-determined if empty
	Description string
	ParentID    int // 0 = no parent
	Tags        []string
}

// Create creates a new task with the given parameters
// Returns the created task or an error
func (s *Service) Create(ctx context.Context, req CreateRequest) (*model.Task, error) {
	// Validate title
	if req.Title == "" {
		return nil, fmt.Errorf("task title cannot be empty")
	}

	// Validate task type
	if !model.IsValidType(req.Type) {
		return nil, fmt.Errorf("invalid task type: %s", req.Type)
	}

	// Generate next ID
	nextID, err := s.engine.NextID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to generate task ID: %w", err)
	}

	// Determine status using task manager
	finalStatus := req.Status
	if finalStatus == "" {
		finalStatus = s.manager.DetermineInitialStatus(req.Type, req.Description)
	}

	// Create task using task manager
	task := s.manager.CreateTask(nextID, req.Title, req.Type, finalStatus, req.Description, req.Tags, req.ParentID)

	// Save task
	if err := s.store.SaveTask(ctx, task); err != nil {
		return nil, fmt.Errorf("failed to save task: %w", err)
	}

	return task, nil
}

// UpdateRequest contains optional updates to apply to a task
// Nil pointer fields indicate "no change"
type UpdateRequest struct {
	Status      *string
	Title       *string
	Description *string
	AddTags     []string
	RemoveTags  []string
}

// Update applies updates to an existing task
// Returns the updated task or an error
func (s *Service) Update(ctx context.Context, taskID int, req UpdateRequest) (*model.Task, error) {
	// Load existing task
	task, err := s.engine.FindByID(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to load task: %w", err)
	}

	// Apply field updates
	if req.Status != nil {
		task.Status = *req.Status
	}

	if req.Title != nil {
		if *req.Title == "" {
			return nil, fmt.Errorf("task title cannot be empty")
		}
		task.Title = *req.Title
	}

	if req.Description != nil {
		task.Description = *req.Description
	}

	// Apply tag operations
	if len(req.AddTags) > 0 {
		task.Tags = s.manager.MergeTags(task.Tags, req.AddTags)
	}

	if len(req.RemoveTags) > 0 {
		task.Tags = s.manager.RemoveTags(task.Tags, req.RemoveTags)
	}

	// Update timestamp
	task.UpdatedAt = time.Now().UTC()

	// Save task
	if err := s.store.SaveTask(ctx, task); err != nil {
		return nil, fmt.Errorf("failed to save task: %w", err)
	}

	return task, nil
}

// Get retrieves a task by ID
func (s *Service) Get(ctx context.Context, taskID int) (*model.Task, error) {
	task, err := s.engine.FindByID(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to load task: %w", err)
	}
	return task, nil
}

// List returns all tasks matching the given filters
func (s *Service) List(ctx context.Context, filters ...query.TaskFilter) ([]*model.Task, error) {
	tasks, err := s.engine.ListTasks(ctx, filters...)
	if err != nil {
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}
	return tasks, nil
}

// Delete removes a task by ID
func (s *Service) Delete(ctx context.Context, taskID int) error {
	if err := s.store.DeleteTask(ctx, taskID); err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}
	return nil
}
