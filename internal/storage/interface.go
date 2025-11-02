package storage

import (
	"context"
	"errors"

	"github.com/zenobi-us/opentask/internal/model"
)

// TaskFilter is a functional option for filtering tasks
type TaskFilter func(*model.Task) bool

// BaseStorage defines the interface all storage backends must implement
type BaseStorage interface {
	// LoadTask retrieves a single task by ID
	// Returns ErrTaskNotFound if task doesn't exist
	LoadTask(ctx context.Context, id int) (*model.Task, error)

	// SaveTask persists a task to storage
	// Creates if task doesn't exist, updates if it does
	SaveTask(ctx context.Context, task *model.Task) error

	// DeleteTask removes a task from storage
	// Returns ErrTaskNotFound if task doesn't exist
	DeleteTask(ctx context.Context, id int) error

	// ListTasks returns all tasks matching the given filters
	// If no filters provided, returns all tasks
	// Filters are applied in AND fashion (all must match)
	ListTasks(ctx context.Context, filters ...TaskFilter) ([]*model.Task, error)

	// NextID generates the next global sequential ID
	// Counts all task files in project, returns count + 1
	// Guaranteed globally unique (no collisions)
	NextID(ctx context.Context) (int, error)

	// GetRelated returns all tasks related to the given task by relationship type
	// relationType must be one of: "parent", "blocks", "relates-to"
	GetRelated(ctx context.Context, taskID int, relationType string) ([]*model.Task, error)

	// Close performs cleanup (if needed for this backend)
	Close() error
}

// StorageConfig contains backend-agnostic configuration
type StorageConfig struct {
	Backend string            // "markdown-fs", "sqlite", etc.
	Path    string            // Project path or database location
	Options map[string]string // Backend-specific options
}

// Common errors
var (
	ErrTaskNotFound         = errors.New("task not found")
	ErrInvalidID            = errors.New("invalid task ID format")
	ErrInvalidTaskType      = errors.New("invalid task type")
	ErrInvalidStatus        = errors.New("invalid status for workflow")
	ErrCircularRelationship = errors.New("circular relationship detected")
)

// NewStorage factory function creates storage backend based on config
func NewStorage(config StorageConfig) (BaseStorage, error) {
	switch config.Backend {
	case "markdown-fs":
		return NewMarkdownFileStorage(config.Path)
	case "memory":
		return NewMemoryStorage(), nil
	default:
		return nil, errors.New("unknown storage backend: " + config.Backend)
	}
}
