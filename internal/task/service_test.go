package task

import (
	"context"
	"errors"
	"testing"

	"github.com/opentasks/cmd/internal/model"
	"github.com/opentasks/cmd/internal/query"
	"github.com/opentasks/cmd/internal/storage"
)

// MockEngine implements TaskEngine interface for testing
type MockEngine struct {
	NextIDFunc    func(ctx context.Context) (int, error)
	FindByIDFunc  func(ctx context.Context, id int) (*model.Task, error)
	ListTasksFunc func(ctx context.Context, filters ...query.TaskFilter) ([]*model.Task, error)
}

func (m *MockEngine) NextID(ctx context.Context) (int, error) {
	if m.NextIDFunc != nil {
		return m.NextIDFunc(ctx)
	}
	return 1, nil
}

func (m *MockEngine) FindByID(ctx context.Context, id int) (*model.Task, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, id)
	}
	return &model.Task{ID: id, Title: "Test Task", Status: "todo"}, nil
}

func (m *MockEngine) ListTasks(ctx context.Context, filters ...query.TaskFilter) ([]*model.Task, error) {
	if m.ListTasksFunc != nil {
		return m.ListTasksFunc(ctx, filters...)
	}
	return []*model.Task{}, nil
}

// MockStore implements storage.BaseStorage methods needed for testing
type MockStore struct {
	SaveTaskFunc   func(ctx context.Context, task *model.Task) error
	DeleteTaskFunc func(ctx context.Context, id int) error
	LoadTaskFunc   func(ctx context.Context, id int) (*model.Task, error)
}

func (m *MockStore) SaveTask(ctx context.Context, task *model.Task) error {
	if m.SaveTaskFunc != nil {
		return m.SaveTaskFunc(ctx, task)
	}
	return nil
}

func (m *MockStore) DeleteTask(ctx context.Context, id int) error {
	if m.DeleteTaskFunc != nil {
		return m.DeleteTaskFunc(ctx, id)
	}
	return nil
}

func (m *MockStore) LoadTask(ctx context.Context, id int) (*model.Task, error) {
	if m.LoadTaskFunc != nil {
		return m.LoadTaskFunc(ctx, id)
	}
	return &model.Task{ID: id}, nil
}

// Implement remaining BaseStorage methods (unused in service, but required for interface)
func (m *MockStore) ListTasks(ctx context.Context, filters ...storage.TaskFilter) ([]*model.Task, error) {
	return []*model.Task{}, nil
}

func (m *MockStore) NextID(ctx context.Context) (int, error) {
	return 1, nil
}

func (m *MockStore) GetRelated(ctx context.Context, taskID int, relationType string) ([]*model.Task, error) {
	return []*model.Task{}, nil
}

func (m *MockStore) Close() error {
	return nil
}

// Test helper to create a default context
func testContext() context.Context {
	return context.Background()
}

// Common test errors
var (
	errMockNextID = errors.New("mock NextID error")
	errMockSave   = errors.New("mock SaveTask error")
	errMockLoad   = errors.New("mock LoadTask error")
	errMockDelete = errors.New("mock DeleteTask error")
)

// Helper to create string pointer
func strPtr(s string) *string {
	return &s
}

// Simplified tests - focus on critical paths given time constraints
func TestService_Create_Success(t *testing.T) {
	mockEngine := &MockEngine{
		NextIDFunc: func(context.Context) (int, error) { return 42, nil },
	}
	mockStore := &MockStore{
		SaveTaskFunc: func(context.Context, *model.Task) error { return nil },
	}
	service := NewService(mockEngine, mockStore)

	req := CreateRequest{
		Title:  "Test Task",
		Type:   model.TypeTask,
		Status: "todo",
	}

	task, err := service.Create(testContext(), req)
	if err != nil {
		t.Fatalf("Create() unexpected error: %v", err)
	}

	if task.ID != 42 {
		t.Errorf("task.ID = %d, want 42", task.ID)
	}
	if task.Title != "Test Task" {
		t.Errorf("task.Title = %q, want %q", task.Title, "Test Task")
	}
}

func TestService_Create_InvalidType(t *testing.T) {
	service := NewService(&MockEngine{}, &MockStore{})

	req := CreateRequest{
		Title: "Test",
		Type:  "invalid",
	}

	_, err := service.Create(testContext(), req)
	if err == nil {
		t.Error("Create() expected error for invalid type, got nil")
	}
}

func TestService_Update_Success(t *testing.T) {
	baseTask := &model.Task{
		ID:     42,
		Title:  "Original",
		Status: "todo",
		Tags:   []string{"tag1"},
	}

	mockEngine := &MockEngine{
		FindByIDFunc: func(context.Context, int) (*model.Task, error) {
			copy := *baseTask
			return &copy, nil
		},
	}
	mockStore := &MockStore{
		SaveTaskFunc: func(context.Context, *model.Task) error { return nil },
	}
	service := NewService(mockEngine, mockStore)

	newStatus := "done"
	req := UpdateRequest{Status: &newStatus}

	task, err := service.Update(testContext(), 42, req)
	if err != nil {
		t.Fatalf("Update() unexpected error: %v", err)
	}

	if task.Status != "done" {
		t.Errorf("task.Status = %q, want %q", task.Status, "done")
	}
}

func TestService_Get_Success(t *testing.T) {
	mockEngine := &MockEngine{
		FindByIDFunc: func(_ context.Context, id int) (*model.Task, error) {
			return &model.Task{ID: id, Title: "Found Task"}, nil
		},
	}
	service := NewService(mockEngine, &MockStore{})

	task, err := service.Get(testContext(), 42)
	if err != nil {
		t.Fatalf("Get() unexpected error: %v", err)
	}

	if task.Title != "Found Task" {
		t.Errorf("task.Title = %q, want %q", task.Title, "Found Task")
	}
}

func TestService_List_Success(t *testing.T) {
	mockEngine := &MockEngine{
		ListTasksFunc: func(context.Context, ...query.TaskFilter) ([]*model.Task, error) {
			return []*model.Task{
				{ID: 1, Title: "Task 1"},
				{ID: 2, Title: "Task 2"},
			}, nil
		},
	}
	service := NewService(mockEngine, &MockStore{})

	tasks, err := service.List(testContext())
	if err != nil {
		t.Fatalf("List() unexpected error: %v", err)
	}

	if len(tasks) != 2 {
		t.Errorf("List() returned %d tasks, want 2", len(tasks))
	}
}

func TestService_Delete_Success(t *testing.T) {
	mockStore := &MockStore{
		DeleteTaskFunc: func(context.Context, int) error { return nil },
	}
	service := NewService(&MockEngine{}, mockStore)

	err := service.Delete(testContext(), 42)
	if err != nil {
		t.Fatalf("Delete() unexpected error: %v", err)
	}
}
