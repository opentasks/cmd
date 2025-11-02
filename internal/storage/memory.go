package storage

import (
	"context"
	"sync"

	"github.com/zenobi-us/opentask/internal/model"
)

// MemoryStorage implements BaseStorage using in-memory storage
// Useful for testing and small projects
type MemoryStorage struct {
	mu    sync.RWMutex
	tasks map[int]*model.Task
	maxID int
}

// NewMemoryStorage creates a new in-memory storage
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		tasks: make(map[int]*model.Task),
		maxID: 0,
	}
}

// LoadTask retrieves a task by ID
func (m *MemoryStorage) LoadTask(ctx context.Context, id int) (*model.Task, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	task, ok := m.tasks[id]
	if !ok {
		return nil, ErrTaskNotFound
	}

	// Return a copy to prevent external modifications
	taskCopy := *task
	return &taskCopy, nil
}

// SaveTask stores a task
func (m *MemoryStorage) SaveTask(ctx context.Context, task *model.Task) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Store a copy
	taskCopy := *task
	m.tasks[task.ID] = &taskCopy

	// Update max ID
	if task.ID > m.maxID {
		m.maxID = task.ID
	}

	return nil
}

// DeleteTask removes a task
func (m *MemoryStorage) DeleteTask(ctx context.Context, id int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.tasks[id]; !ok {
		return ErrTaskNotFound
	}

	delete(m.tasks, id)
	return nil
}

// ListTasks returns all tasks matching filters
func (m *MemoryStorage) ListTasks(ctx context.Context, filters ...TaskFilter) ([]*model.Task, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*model.Task

	for _, task := range m.tasks {
		// Apply filters
		match := true
		for _, filter := range filters {
			if !filter(task) {
				match = false
				break
			}
		}

		if match {
			// Add a copy
			taskCopy := *task
			result = append(result, &taskCopy)
		}
	}

	return result, nil
}

// NextID generates the next ID
func (m *MemoryStorage) NextID(ctx context.Context) (int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.maxID + 1, nil
}

// GetRelated returns tasks related by type
func (m *MemoryStorage) GetRelated(ctx context.Context, taskID int, relationType string) ([]*model.Task, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var related []*model.Task

	for _, task := range m.tasks {
		for _, rel := range task.Relationships {
			if rel.Type == relationType && rel.TaskID == taskID {
				taskCopy := *task
				related = append(related, &taskCopy)
			}
		}
	}

	return related, nil
}

// Close is a no-op for memory storage
func (m *MemoryStorage) Close() error {
	return nil
}
