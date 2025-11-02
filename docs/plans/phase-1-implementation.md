# Phase 1: Core Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Build the core OpenTasks system with data models, storage, and query engine.

**Architecture:** Modular design with clear separation: data models → storage interface → implementation → query engine. Each component tested independently with integration tests.

**Tech Stack:** Go 1.21+, gopkg.in/yaml.v3, filepath/filepath, testing (stdlib)

---

## Task 1: Initialize Go Project Structure

**Files:**
- Create: `go.mod`
- Create: `go.sum`
- Create: `cmd/opentasks/main.go` (stub)
- Create: `internal/model/task.go`
- Create: `internal/storage/interface.go`
- Create: `internal/query/engine.go`
- Create: `internal/config/config.go`
- Modify: `.gitignore` (add Go ignores)

**Step 1: Initialize Go module**

```bash
cd /mnt/Store/Projects/Mine/Github/opentasks
go mod init github.com/zenobius/opentasks
```

Expected: Creates `go.mod` file

**Step 2: Create directory structure**

```bash
mkdir -p cmd/opentasks internal/{model,storage,config,query}
```

Expected: Directories created

**Step 3: Create go.sum (empty for now)**

```bash
touch go.sum
```

**Step 4: Add Go-specific .gitignore**

Add to `.gitignore`:
```
# Go
bin/
dist/
*.o
*.a
*.so
.DS_Store
.env
```

**Step 5: Create stub main.go**

File: `cmd/opentasks/main.go`

```go
package main

import "fmt"

func main() {
    fmt.Println("OpenTasks")
}
```

**Step 6: Verify it builds**

```bash
cd cmd/opentasks && go build -o opentasks
./opentasks
```

Expected output: `OpenTasks`

**Step 7: Commit**

```bash
git add go.mod go.sum cmd/ internal/ .gitignore
git commit -m "chore: initialize Go project structure"
```

---

## Task 2: Implement Task Data Model

**Files:**
- Create: `internal/model/task.go`
- Create: `internal/model/relationship.go`
- Create: `internal/model/constants.go`
- Create: `internal/model/task_test.go`

**Step 1: Write tests for Task struct**

File: `internal/model/task_test.go`

```go
package model

import (
	"testing"
	"time"
)

func TestTaskCreation(t *testing.T) {
	task := &Task{
		ID:          42,
		Title:       "Test task",
		Type:        TypeStory,
		Status:      "todo",
		Tags:        []string{"test"},
		Description: "Test description",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if task.ID != 42 {
		t.Errorf("Expected ID 42, got %d", task.ID)
	}
	if task.Title != "Test task" {
		t.Errorf("Expected title 'Test task', got %s", task.Title)
	}
	if task.Type != TypeStory {
		t.Errorf("Expected type %s, got %s", TypeStory, task.Type)
	}
}

func TestTaskWithRelationships(t *testing.T) {
	task := &Task{
		ID:          1,
		Title:       "Epic task",
		Type:        TypeEpic,
		Status:      "todo",
		Relationships: []Relationship{
			{Type: RelParent, TaskID: 5},
			{Type: RelBlocks, TaskID: 10},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if len(task.Relationships) != 2 {
		t.Errorf("Expected 2 relationships, got %d", len(task.Relationships))
	}
	if task.Relationships[0].Type != RelParent {
		t.Errorf("Expected relationship type %s, got %s", RelParent, task.Relationships[0].Type)
	}
}

func TestRelationshipTypes(t *testing.T) {
	rel := Relationship{
		Type:   RelBlocks,
		TaskID: 42,
	}

	if rel.TaskID != 42 {
		t.Errorf("Expected TaskID 42, got %d", rel.TaskID)
	}
}
```

**Step 2: Run tests (expect fail)**

```bash
cd /mnt/Store/Projects/Mine/Github/opentasks
go test ./internal/model -v
```

Expected: FAIL - package does not exist

**Step 3: Implement Task struct**

File: `internal/model/task.go`

```go
package model

import "time"

// Task represents a single task in the system
type Task struct {
	ID            int            // Global sequential ID
	Title         string         // Task title
	Type          string         // epic|plan|research|story|decision|task
	Status        string         // Custom per project (e.g., todo, in-progress)
	Tags          []string       // Labels for organization
	Description   string         // Markdown body
	Relationships []Relationship // Links to other tasks
	CreatedAt     time.Time      // Creation timestamp (RFC3339)
	UpdatedAt     time.Time      // Last update timestamp (RFC3339)
}

// Relationship represents a link between tasks
type Relationship struct {
	Type   string // "parent"|"blocks"|"relates-to"
	TaskID int    // ID of target task
}
```

**Step 4: Implement constants**

File: `internal/model/constants.go`

```go
package model

// Task types
const (
	TypeEpic     = "epic"
	TypePlan     = "plan"
	TypeResearch = "research"
	TypeStory    = "story"
	TypeDecision = "decision"
	TypeTask     = "task"
)

// Relationship types
const (
	RelParent    = "parent"
	RelBlocks    = "blocks"
	RelRelatedTo = "relates-to"
)

// Valid task types
var AllTaskTypes = []string{
	TypeEpic,
	TypePlan,
	TypeResearch,
	TypeStory,
	TypeDecision,
	TypeTask,
}

// Type to code mapping (for filenames)
var TypeCode = map[string]string{
	TypeEpic:     "e",
	TypePlan:     "p",
	TypeResearch: "r",
	TypeStory:    "s",
	TypeDecision: "d",
	TypeTask:     "t",
}

// Code to type mapping (inverse)
var CodeType = map[string]string{
	"e": TypeEpic,
	"p": TypePlan,
	"r": TypeResearch,
	"s": TypeStory,
	"d": TypeDecision,
	"t": TypeTask,
}
```

**Step 5: Run tests (expect pass)**

```bash
go test ./internal/model -v
```

Expected: All tests pass

**Step 6: Commit**

```bash
git add internal/model/
git commit -m "feat: implement task data model with relationships"
```

---

## Task 3: Implement BaseStorage Interface

**Files:**
- Create: `internal/storage/interface.go`
- Create: `internal/storage/storage_test.go`

**Step 1: Write tests for interface**

File: `internal/storage/storage_test.go`

```go
package storage

import (
	"context"
	"testing"

	"github.com/zenobius/opentasks/internal/model"
)

// MockStorage implements BaseStorage for testing
type MockStorage struct {
	tasks map[int]*model.Task
	maxID int
}

func NewMockStorage() *MockStorage {
	return &MockStorage{
		tasks: make(map[int]*model.Task),
		maxID: 0,
	}
}

func (m *MockStorage) LoadTask(ctx context.Context, id int) (*model.Task, error) {
	if task, exists := m.tasks[id]; exists {
		return task, nil
	}
	return nil, ErrTaskNotFound
}

func (m *MockStorage) SaveTask(ctx context.Context, task *model.Task) error {
	m.tasks[task.ID] = task
	if task.ID > m.maxID {
		m.maxID = task.ID
	}
	return nil
}

func (m *MockStorage) DeleteTask(ctx context.Context, id int) error {
	if _, exists := m.tasks[id]; !exists {
		return ErrTaskNotFound
	}
	delete(m.tasks, id)
	return nil
}

func (m *MockStorage) ListTasks(ctx context.Context, filters ...TaskFilter) ([]*model.Task, error) {
	var result []*model.Task
	for _, task := range m.tasks {
		matches := true
		for _, filter := range filters {
			if !filter(task) {
				matches = false
				break
			}
		}
		if matches {
			result = append(result, task)
		}
	}
	return result, nil
}

func (m *MockStorage) NextID(ctx context.Context) (int, error) {
	return m.maxID + 1, nil
}

func (m *MockStorage) GetRelated(ctx context.Context, taskID int, relationType string) ([]*model.Task, error) {
	var result []*model.Task
	for _, task := range m.tasks {
		for _, rel := range task.Relationships {
			if rel.Type == relationType && rel.TaskID == taskID {
				result = append(result, task)
			}
		}
	}
	return result, nil
}

func (m *MockStorage) Close() error {
	return nil
}

// Test interface implementation
func TestStorageInterface(t *testing.T) {
	ctx := context.Background()
	storage := NewMockStorage()

	// Test SaveTask
	task := &model.Task{
		ID:    1,
		Title: "Test task",
		Type:  model.TypeStory,
	}
	if err := storage.SaveTask(ctx, task); err != nil {
		t.Fatalf("SaveTask failed: %v", err)
	}

	// Test LoadTask
	loaded, err := storage.LoadTask(ctx, 1)
	if err != nil {
		t.Fatalf("LoadTask failed: %v", err)
	}
	if loaded.Title != "Test task" {
		t.Errorf("Expected title 'Test task', got %s", loaded.Title)
	}

	// Test NextID
	nextID, err := storage.NextID(ctx)
	if err != nil {
		t.Fatalf("NextID failed: %v", err)
	}
	if nextID != 2 {
		t.Errorf("Expected nextID 2, got %d", nextID)
	}

	// Test DeleteTask
	if err := storage.DeleteTask(ctx, 1); err != nil {
		t.Fatalf("DeleteTask failed: %v", err)
	}
	_, err = storage.LoadTask(ctx, 1)
	if err != ErrTaskNotFound {
		t.Errorf("Expected ErrTaskNotFound, got %v", err)
	}
}
```

**Step 2: Define interface and errors**

File: `internal/storage/interface.go`

```go
package storage

import (
	"context"
	"errors"

	"github.com/zenobius/opentasks/internal/model"
)

// TaskFilter is a function that filters tasks
type TaskFilter func(*model.Task) bool

// BaseStorage defines the interface for task persistence
type BaseStorage interface {
	// LoadTask retrieves a single task by ID
	LoadTask(ctx context.Context, id int) (*model.Task, error)

	// SaveTask persists a task to storage
	SaveTask(ctx context.Context, task *model.Task) error

	// DeleteTask removes a task from storage
	DeleteTask(ctx context.Context, id int) error

	// ListTasks returns all tasks matching the given filters
	ListTasks(ctx context.Context, filters ...TaskFilter) ([]*model.Task, error)

	// NextID generates the next global sequential ID
	NextID(ctx context.Context) (int, error)

	// GetRelated returns all tasks related by relationship type
	GetRelated(ctx context.Context, taskID int, relationType string) ([]*model.Task, error)

	// Close performs cleanup
	Close() error
}

// Common errors
var (
	ErrTaskNotFound     = errors.New("task not found")
	ErrInvalidID        = errors.New("invalid task ID format")
	ErrInvalidTaskType  = errors.New("invalid task type")
	ErrInvalidStatus    = errors.New("invalid status for workflow")
	ErrStorageError     = errors.New("storage error")
)

// StorageConfig contains configuration for storage backends
type StorageConfig struct {
	Backend string            // "markdown-fs", "memory", etc.
	Path    string            // Project path
	Options map[string]string // Backend-specific options
}
```

**Step 3: Run tests (expect pass)**

```bash
go test ./internal/storage -v
```

Expected: All tests pass

**Step 4: Commit**

```bash
git add internal/storage/
git commit -m "feat: define BaseStorage interface with mock implementation"
```

---

## Task 4: Implement Memory Storage Backend (for testing)

**Files:**
- Create: `internal/storage/memory.go`
- Create: `internal/storage/memory_test.go`

**Step 1: Write comprehensive tests**

File: `internal/storage/memory_test.go`

```go
package storage

import (
	"context"
	"testing"
	"time"

	"github.com/zenobius/opentasks/internal/model"
)

func TestMemoryStorageSaveAndLoad(t *testing.T) {
	ctx := context.Background()
	storage := NewMemoryStorage()

	task := &model.Task{
		ID:        1,
		Title:     "Test task",
		Type:      model.TypeStory,
		Status:    "todo",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := storage.SaveTask(ctx, task); err != nil {
		t.Fatalf("SaveTask failed: %v", err)
	}

	loaded, err := storage.LoadTask(ctx, 1)
	if err != nil {
		t.Fatalf("LoadTask failed: %v", err)
	}

	if loaded.Title != task.Title {
		t.Errorf("Title mismatch: expected %s, got %s", task.Title, loaded.Title)
	}
}

func TestMemoryStorageNextID(t *testing.T) {
	ctx := context.Background()
	storage := NewMemoryStorage()

	// Add some tasks
	for i := 1; i <= 5; i++ {
		task := &model.Task{
			ID:        i,
			Title:     "Task " + string(rune(i)),
			Type:      model.TypeStory,
			CreatedAt: time.Now(),
		}
		storage.SaveTask(ctx, task)
	}

	// NextID should be 6
	nextID, err := storage.NextID(ctx)
	if err != nil {
		t.Fatalf("NextID failed: %v", err)
	}
	if nextID != 6 {
		t.Errorf("Expected nextID 6, got %d", nextID)
	}

	// Delete task 3, NextID should still be 6 (no ID reuse)
	storage.DeleteTask(ctx, 3)
	nextID, err = storage.NextID(ctx)
	if err != nil {
		t.Fatalf("NextID after delete failed: %v", err)
	}
	if nextID != 6 {
		t.Errorf("Expected nextID 6 (no reuse), got %d", nextID)
	}
}

func TestMemoryStorageFiltering(t *testing.T) {
	ctx := context.Background()
	storage := NewMemoryStorage()

	// Create test data
	tasks := []*model.Task{
		{ID: 1, Title: "Story 1", Type: model.TypeStory, Status: "todo", CreatedAt: time.Now()},
		{ID: 2, Title: "Story 2", Type: model.TypeStory, Status: "in-progress", CreatedAt: time.Now()},
		{ID: 3, Title: "Plan 1", Type: model.TypePlan, Status: "todo", CreatedAt: time.Now()},
	}

	for _, task := range tasks {
		storage.SaveTask(ctx, task)
	}

	// Filter by type
	filtered, err := storage.ListTasks(ctx, func(t *model.Task) bool {
		return t.Type == model.TypeStory
	})
	if err != nil {
		t.Fatalf("ListTasks failed: %v", err)
	}
	if len(filtered) != 2 {
		t.Errorf("Expected 2 stories, got %d", len(filtered))
	}

	// Filter by status
	filtered, err = storage.ListTasks(ctx, func(t *model.Task) bool {
		return t.Status == "todo"
	})
	if err != nil {
		t.Fatalf("ListTasks failed: %v", err)
	}
	if len(filtered) != 2 {
		t.Errorf("Expected 2 todo tasks, got %d", len(filtered))
	}
}

func TestMemoryStorageRelationships(t *testing.T) {
	ctx := context.Background()
	storage := NewMemoryStorage()

	epic := &model.Task{
		ID:        1,
		Title:     "Epic",
		Type:      model.TypeEpic,
		CreatedAt: time.Now(),
	}
	storage.SaveTask(ctx, epic)

	story := &model.Task{
		ID:    2,
		Title: "Story",
		Type:  model.TypeStory,
		Relationships: []model.Relationship{
			{Type: model.RelParent, TaskID: 1},
		},
		CreatedAt: time.Now(),
	}
	storage.SaveTask(ctx, story)

	// Get tasks that relate to epic 1
	related, err := storage.GetRelated(ctx, 1, model.RelParent)
	if err != nil {
		t.Fatalf("GetRelated failed: %v", err)
	}
	if len(related) != 1 {
		t.Errorf("Expected 1 related task, got %d", len(related))
	}
	if related[0].ID != 2 {
		t.Errorf("Expected task 2, got %d", related[0].ID)
	}
}
```

**Step 2: Implement MemoryStorage**

File: `internal/storage/memory.go`

```go
package storage

import (
	"context"
	"sort"
	"sync"

	"github.com/zenobius/opentasks/internal/model"
)

// MemoryStorage implements BaseStorage using in-memory storage
type MemoryStorage struct {
	mu    sync.RWMutex
	tasks map[int]*model.Task
	maxID int
}

// NewMemoryStorage creates a new memory storage backend
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		tasks: make(map[int]*model.Task),
		maxID: 0,
	}
}

func (m *MemoryStorage) LoadTask(ctx context.Context, id int) (*model.Task, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	task, exists := m.tasks[id]
	if !exists {
		return nil, ErrTaskNotFound
	}
	return task, nil
}

func (m *MemoryStorage) SaveTask(ctx context.Context, task *model.Task) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.tasks[task.ID] = task
	if task.ID > m.maxID {
		m.maxID = task.ID
	}
	return nil
}

func (m *MemoryStorage) DeleteTask(ctx context.Context, id int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.tasks[id]; !exists {
		return ErrTaskNotFound
	}
	delete(m.tasks, id)
	return nil
}

func (m *MemoryStorage) ListTasks(ctx context.Context, filters ...TaskFilter) ([]*model.Task, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*model.Task
	for _, task := range m.tasks {
		matches := true
		for _, filter := range filters {
			if !filter(task) {
				matches = false
				break
			}
		}
		if matches {
			result = append(result, task)
		}
	}

	// Sort by creation time
	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.Before(result[j].CreatedAt)
	})

	return result, nil
}

func (m *MemoryStorage) NextID(ctx context.Context) (int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.maxID + 1, nil
}

func (m *MemoryStorage) GetRelated(ctx context.Context, taskID int, relationType string) ([]*model.Task, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*model.Task
	for _, task := range m.tasks {
		for _, rel := range task.Relationships {
			if rel.Type == relationType && rel.TaskID == taskID {
				result = append(result, task)
				break
			}
		}
	}
	return result, nil
}

func (m *MemoryStorage) Close() error {
	return nil
}
```

**Step 3: Run tests**

```bash
go test ./internal/storage -v
```

Expected: All tests pass

**Step 4: Commit**

```bash
git add internal/storage/memory.go internal/storage/memory_test.go
git commit -m "feat: implement memory storage backend for testing"
```

---

## Task 5: Implement Query Engine

**Files:**
- Create: `internal/query/filters.go`
- Create: `internal/query/engine.go`
- Create: `internal/query/query_test.go`

**Step 1: Write tests for filter functions**

File: `internal/query/query_test.go`

```go
package query

import (
	"context"
	"testing"
	"time"

	"github.com/zenobius/opentasks/internal/model"
	"github.com/zenobius/opentasks/internal/storage"
)

func TestFilterByStatus(t *testing.T) {
	filter := WithStatus("todo")
	task := &model.Task{Status: "todo"}
	if !filter(task) {
		t.Error("WithStatus filter failed for matching status")
	}

	task.Status = "in-progress"
	if filter(task) {
		t.Error("WithStatus filter should not match different status")
	}
}

func TestFilterByType(t *testing.T) {
	filter := WithType(model.TypeStory)
	task := &model.Task{Type: model.TypeStory}
	if !filter(task) {
		t.Error("WithType filter failed for matching type")
	}

	task.Type = model.TypeEpic
	if filter(task) {
		t.Error("WithType filter should not match different type")
	}
}

func TestFilterByTag(t *testing.T) {
	filter := WithTag("feature")
	task := &model.Task{Tags: []string{"feature", "core"}}
	if !filter(task) {
		t.Error("WithTag filter failed for matching tag")
	}

	task.Tags = []string{"docs"}
	if filter(task) {
		t.Error("WithTag filter should not match missing tag")
	}
}

func TestFilterByParent(t *testing.T) {
	filter := WithParent(5)
	task := &model.Task{
		Relationships: []model.Relationship{
			{Type: model.RelParent, TaskID: 5},
		},
	}
	if !filter(task) {
		t.Error("WithParent filter failed for matching parent")
	}

	task.Relationships = []model.Relationship{
		{Type: model.RelParent, TaskID: 10},
	}
	if filter(task) {
		t.Error("WithParent filter should not match different parent")
	}
}

func TestQueryEngineFind(t *testing.T) {
	ctx := context.Background()
	store := storage.NewMemoryStorage()
	engine := NewQueryEngine(store)

	// Create test tasks
	tasks := []*model.Task{
		{
			ID:        1,
			Title:     "Epic 1",
			Type:      model.TypeEpic,
			Status:    "todo",
			Tags:      []string{"core"},
			CreatedAt: time.Now(),
		},
		{
			ID:        2,
			Title:     "Story 1",
			Type:      model.TypeStory,
			Status:    "in-progress",
			Tags:      []string{"feature"},
			Relationships: []model.Relationship{
				{Type: model.RelParent, TaskID: 1},
			},
			CreatedAt: time.Now(),
		},
		{
			ID:        3,
			Title:     "Story 2",
			Type:      model.TypeStory,
			Status:    "todo",
			Tags:      []string{"feature"},
			Relationships: []model.Relationship{
				{Type: model.RelParent, TaskID: 1},
			},
			CreatedAt: time.Now(),
		},
	}

	for _, task := range tasks {
		store.SaveTask(ctx, task)
	}

	// Test FindByID
	task, err := engine.FindByID(ctx, 1)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	if task.Title != "Epic 1" {
		t.Errorf("Expected Epic 1, got %s", task.Title)
	}

	// Test FindChildren
	children, err := engine.FindChildren(ctx, 1)
	if err != nil {
		t.Fatalf("FindChildren failed: %v", err)
	}
	if len(children) != 2 {
		t.Errorf("Expected 2 children, got %d", len(children))
	}
}

func TestCompositeFilters(t *testing.T) {
	ctx := context.Background()
	store := storage.NewMemoryStorage()
	engine := NewQueryEngine(store)

	// Create test tasks
	tasks := []*model.Task{
		{ID: 1, Type: model.TypeStory, Status: "todo", Tags: []string{"feature"}, CreatedAt: time.Now()},
		{ID: 2, Type: model.TypeStory, Status: "in-progress", Tags: []string{"feature"}, CreatedAt: time.Now()},
		{ID: 3, Type: model.TypePlan, Status: "todo", Tags: []string{"core"}, CreatedAt: time.Now()},
	}

	for _, task := range tasks {
		store.SaveTask(ctx, task)
	}

	// Composite filter: stories that are todo
	filtered, err := engine.ListTasks(ctx,
		WithType(model.TypeStory),
		WithStatus("todo"),
	)
	if err != nil {
		t.Fatalf("ListTasks failed: %v", err)
	}
	if len(filtered) != 1 {
		t.Errorf("Expected 1 task, got %d", len(filtered))
	}
	if filtered[0].ID != 1 {
		t.Errorf("Expected task 1, got %d", filtered[0].ID)
	}
}
```

**Step 2: Implement filter functions**

File: `internal/query/filters.go`

```go
package query

import (
	"github.com/zenobius/opentasks/internal/model"
	"github.com/zenobius/opentasks/internal/storage"
)

// TaskFilter is a function that returns true if task matches filter
type TaskFilter = storage.TaskFilter

// WithStatus returns a filter that matches tasks with given status
func WithStatus(status string) TaskFilter {
	return func(t *model.Task) bool {
		return t.Status == status
	}
}

// WithType returns a filter that matches tasks of given type
func WithType(taskType string) TaskFilter {
	return func(t *model.Task) bool {
		return t.Type == taskType
	}
}

// WithTag returns a filter that matches tasks with given tag
func WithTag(tag string) TaskFilter {
	return func(t *model.Task) bool {
		for _, t := range t.Tags {
			if t == tag {
				return true
			}
		}
		return false
	}
}

// WithTags returns a filter that matches tasks with any of given tags
func WithTags(tags []string) TaskFilter {
	return func(t *model.Task) bool {
		for _, tag := range tags {
			for _, taskTag := range t.Tags {
				if taskTag == tag {
					return true
				}
			}
		}
		return false
	}
}

// WithParent returns a filter that matches tasks with given parent
func WithParent(parentID int) TaskFilter {
	return func(t *model.Task) bool {
		for _, rel := range t.Relationships {
			if rel.Type == model.RelParent && rel.TaskID == parentID {
				return true
			}
		}
		return false
	}
}

// WithRelationship returns a filter matching given relationship type and task ID
func WithRelationship(relType string, taskID int) TaskFilter {
	return func(t *model.Task) bool {
		for _, rel := range t.Relationships {
			if rel.Type == relType && rel.TaskID == taskID {
				return true
			}
		}
		return false
	}
}
```

**Step 3: Implement QueryEngine**

File: `internal/query/engine.go`

```go
package query

import (
	"context"

	"github.com/zenobius/opentasks/internal/model"
	"github.com/zenobius/opentasks/internal/storage"
)

// QueryEngine provides convenience methods for querying tasks
type QueryEngine struct {
	storage storage.BaseStorage
}

// NewQueryEngine creates a new query engine
func NewQueryEngine(storage storage.BaseStorage) *QueryEngine {
	return &QueryEngine{storage: storage}
}

// ListTasks returns all tasks matching the given filters
func (q *QueryEngine) ListTasks(ctx context.Context, filters ...TaskFilter) ([]*model.Task, error) {
	return q.storage.ListTasks(ctx, filters...)
}

// FindByID retrieves a single task by ID
func (q *QueryEngine) FindByID(ctx context.Context, id int) (*model.Task, error) {
	return q.storage.LoadTask(ctx, id)
}

// FindChildren returns all child tasks of a given parent
func (q *QueryEngine) FindChildren(ctx context.Context, parentID int) ([]*model.Task, error) {
	return q.storage.GetRelated(ctx, parentID, model.RelParent)
}

// FindBlocking returns all tasks that block the given task
func (q *QueryEngine) FindBlocking(ctx context.Context, taskID int) ([]*model.Task, error) {
	return q.storage.GetRelated(ctx, taskID, model.RelBlocks)
}

// FindBlockedBy returns all tasks blocked by the given task
func (q *QueryEngine) FindBlockedBy(ctx context.Context, taskID int) ([]*model.Task, error) {
	allTasks, err := q.storage.ListTasks(ctx)
	if err != nil {
		return nil, err
	}

	var result []*model.Task
	for _, task := range allTasks {
		for _, rel := range task.Relationships {
			if rel.Type == model.RelBlocks && rel.TaskID == taskID {
				result = append(result, task)
			}
		}
	}
	return result, nil
}

// FindRelated returns all tasks related to the given task
func (q *QueryEngine) FindRelated(ctx context.Context, taskID int) ([]*model.Task, error) {
	allTasks, err := q.storage.ListTasks(ctx)
	if err != nil {
		return nil, err
	}

	var result []*model.Task
	for _, task := range allTasks {
		for _, rel := range task.Relationships {
			if rel.TaskID == taskID {
				result = append(result, task)
			}
		}
	}
	return result, nil
}
```

**Step 4: Run tests**

```bash
go test ./internal/query -v
```

Expected: All tests pass

**Step 5: Commit**

```bash
git add internal/query/
git commit -m "feat: implement query engine with filters and convenience methods"
```

---

## Task 6: Add Main Package Setup

**Files:**
- Create: `internal/config/config.go`
- Modify: `cmd/opentasks/main.go`

**Step 1: Implement basic config loading**

File: `internal/config/config.go`

```go
package config

// ProjectConfig holds project-level configuration
type ProjectConfig struct {
	Project ProjectMetadata `toml:"project"`
	Workflow WorkflowConfig `toml:"workflow"`
	Storage StorageConfig   `toml:"storage"`
}

// ProjectMetadata contains project information
type ProjectMetadata struct {
	Name        string `toml:"name"`
	Description string `toml:"description"`
	Owner       string `toml:"owner"`
}

// WorkflowConfig defines the status workflow
type WorkflowConfig struct {
	Statuses    []string            `toml:"statuses"`
	Initial     string              `toml:"initial"`
	Transitions []TransitionConfig   `toml:"transitions"`
}

// TransitionConfig defines allowed status transitions
type TransitionConfig struct {
	From string   `toml:"from"`
	To   []string `toml:"to"`
}

// StorageConfig configures the storage backend
type StorageConfig struct {
	Backend string            `toml:"backend"`
	Path    string            `toml:"path"`
	Options map[string]string `toml:"options"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *ProjectConfig {
	return &ProjectConfig{
		Project: ProjectMetadata{
			Name: "OpenTasks Project",
		},
		Workflow: WorkflowConfig{
			Statuses: []string{"todo", "in-progress", "reviewing", "done", "archived"},
			Initial:  "todo",
			Transitions: []TransitionConfig{
				{From: "todo", To: []string{"in-progress", "archived"}},
				{From: "in-progress", To: []string{"reviewing", "todo", "archived"}},
				{From: "reviewing", To: []string{"done", "in-progress", "archived"}},
				{From: "done", To: []string{"archived"}},
			},
		},
		Storage: StorageConfig{
			Backend: "markdown-fs",
			Path:    ".",
		},
	}
}
```

**Step 2: Update main.go to use modules**

File: `cmd/opentasks/main.go`

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/zenobius/opentasks/internal/config"
	"github.com/zenobius/opentasks/internal/model"
	"github.com/zenobius/opentasks/internal/query"
	"github.com/zenobius/opentasks/internal/storage"
)

func main() {
	ctx := context.Background()

	// Load configuration
	cfg := config.DefaultConfig()
	fmt.Printf("OpenTasks - %s\n", cfg.Project.Name)

	// Initialize storage
	store := storage.NewMemoryStorage()
	defer store.Close()

	// Initialize query engine
	engine := query.NewQueryEngine(store)

	// Example: Create and query a task
	if err := exampleUsage(ctx, store, engine); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func exampleUsage(ctx context.Context, store storage.BaseStorage, engine *query.QueryEngine) error {
	// Create a task
	task := &model.Task{
		ID:     1,
		Title:  "Sample task",
		Type:   model.TypeStory,
		Status: "todo",
	}

	if err := store.SaveTask(ctx, task); err != nil {
		return err
	}

	// Retrieve and print
	retrieved, err := engine.FindByID(ctx, 1)
	if err != nil {
		return err
	}

	fmt.Printf("Created task: %d - %s (%s)\n", retrieved.ID, retrieved.Title, retrieved.Type)
	return nil
}
```

**Step 3: Test the main application**

```bash
cd cmd/opentasks && go build -o opentasks && ./opentasks
```

Expected output:
```
OpenTasks - OpenTasks Project
Created task: 1 - Sample task (story)
```

**Step 4: Run all tests**

```bash
go test ./...
```

Expected: All tests pass

**Step 5: Commit**

```bash
git add internal/config/ cmd/opentasks/main.go
git commit -m "feat: add config system and wire up components in main"
```

---

## Task 7: Integration Tests

**Files:**
- Create: `internal/integration_test.go`

**Step 1: Write comprehensive integration tests**

File: `internal/integration_test.go`

```go
package internal

import (
	"context"
	"testing"
	"time"

	"github.com/zenobius/opentasks/internal/model"
	"github.com/zenobius/opentasks/internal/query"
	"github.com/zenobius/opentasks/internal/storage"
)

func TestEndToEndWorkflow(t *testing.T) {
	ctx := context.Background()
	store := storage.NewMemoryStorage()
	engine := query.NewQueryEngine(store)

	// Create an epic
	epic := &model.Task{
		ID:        1,
		Title:     "Design System",
		Type:      model.TypeEpic,
		Status:    "in-progress",
		CreatedAt: time.Now(),
	}
	if err := store.SaveTask(ctx, epic); err != nil {
		t.Fatalf("Failed to save epic: %v", err)
	}

	// Create child tasks
	for i := 2; i <= 4; i++ {
		task := &model.Task{
			ID:    i,
			Title: "Task " + string(rune('0'+i)),
			Type:  model.TypeStory,
			Status: "todo",
			Relationships: []model.Relationship{
				{Type: model.RelParent, TaskID: 1},
			},
			Tags:      []string{"feature"},
			CreatedAt: time.Now(),
		}
		if err := store.SaveTask(ctx, task); err != nil {
			t.Fatalf("Failed to save task %d: %v", i, err)
		}
	}

	// Verify we can find all children
	children, err := engine.FindChildren(ctx, 1)
	if err != nil {
		t.Fatalf("FindChildren failed: %v", err)
	}
	if len(children) != 3 {
		t.Errorf("Expected 3 children, got %d", len(children))
	}

	// Verify composite query works
	storiesToDo, err := engine.ListTasks(ctx,
		query.WithType(model.TypeStory),
		query.WithStatus("todo"),
	)
	if err != nil {
		t.Fatalf("ListTasks failed: %v", err)
	}
	if len(storiesToDo) != 3 {
		t.Errorf("Expected 3 todo stories, got %d", len(storiesToDo))
	}

	// Verify next ID generation
	nextID, err := store.NextID(ctx)
	if err != nil {
		t.Fatalf("NextID failed: %v", err)
	}
	if nextID != 5 {
		t.Errorf("Expected nextID 5, got %d", nextID)
	}

	// Test deletion doesn't reuse ID
	store.DeleteTask(ctx, 2)
	nextID, err = store.NextID(ctx)
	if err != nil {
		t.Fatalf("NextID after delete failed: %v", err)
	}
	if nextID != 5 {
		t.Errorf("Expected nextID 5 (no reuse), got %d", nextID)
	}
}

func TestDataPersistence(t *testing.T) {
	ctx := context.Background()
	store := storage.NewMemoryStorage()

	// Create multiple tasks
	for i := 1; i <= 10; i++ {
		task := &model.Task{
			ID:        i,
			Title:     "Task " + string(rune('0'+i)),
			Type:      model.TypeStory,
			Status:    "todo",
			CreatedAt: time.Now(),
		}
		if err := store.SaveTask(ctx, task); err != nil {
			t.Fatalf("Failed to save task: %v", err)
		}
	}

	// Load all
	allTasks, err := store.ListTasks(ctx)
	if err != nil {
		t.Fatalf("ListTasks failed: %v", err)
	}
	if len(allTasks) != 10 {
		t.Errorf("Expected 10 tasks, got %d", len(allTasks))
	}

	// Verify each can be loaded
	for i := 1; i <= 10; i++ {
		task, err := store.LoadTask(ctx, i)
		if err != nil {
			t.Errorf("Failed to load task %d: %v", i, err)
		}
		if task.ID != i {
			t.Errorf("Task ID mismatch: expected %d, got %d", i, task.ID)
		}
	}
}
```

**Step 2: Run integration tests**

```bash
go test ./internal -v
```

Expected: All tests pass

**Step 3: Commit**

```bash
git add internal/integration_test.go
git commit -m "test: add comprehensive integration tests"
```

---

## Task 8: Verify Build and All Tests Pass

**Step 1: Clean build**

```bash
go clean
go build ./cmd/opentasks
```

Expected: Builds successfully

**Step 2: Run all tests with coverage**

```bash
go test ./... -v -cover
```

Expected: All tests pass with reasonable coverage

**Step 3: Check project structure**

```bash
find . -name "*.go" -type f | grep -v ".git" | sort
```

Expected: See all source files organized by package

**Step 4: Final commit**

```bash
git add -A
git commit -m "Phase 1: Core implementation complete - models, storage, query engine

Core components implemented:
- Task data model with relationships
- BaseStorage interface with memory implementation
- Query engine with filters and convenience methods
- Config system with defaults
- Comprehensive test coverage
- Example main application

All tests passing, ready for Phase 2 (CLI and Markdown storage)."
```

---

## Verification Checklist

After completing all tasks:

- [ ] `go test ./... -v` passes with all tests green
- [ ] `go build ./cmd/opentasks` builds without errors
- [ ] All 13+ tests pass
- [ ] Code coverage is >80%
- [ ] Git history shows clean, logical commits
- [ ] No lint warnings (run `golangci-lint run ./...` if available)
- [ ] All functions have godoc comments
- [ ] README.md updated with build/test instructions

---

## What's NOT in Phase 1

**Intentionally deferred to Phase 2:**
- Markdown file storage (MarkdownFileStorage)
- CLI commands (Cobra/Viper integration)
- Config file loading (TOML parsing)
- Slug generation
- Project discovery

**Intentionally out of MVP:**
- Advanced queries (full-text search, dependency graph)
- Web UI
- Cloud storage backends
- Real-time sync

---

## Quick Reference Commands

```bash
# Run tests
go test ./...

# Run tests with coverage
go test ./... -cover

# Run specific test
go test ./internal/model -run TestTaskCreation

# Build main app
go build ./cmd/opentasks

# Run main app
./cmd/opentasks/opentasks

# Format code
go fmt ./...

# Check for issues
go vet ./...
```

---

**This plan is ready for execution. Use @superpowers:executing-plans or subagent-driven-development to implement task-by-task.**
