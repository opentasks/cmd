package graph

import (
	"context"
	"fmt"
	"testing"

	"github.com/opentasks/cmd/internal/model"
)

// MockTaskLoader is a test implementation of TaskLoader
type MockTaskLoader struct {
	tasks map[int]*model.Task
	fail  map[int]bool // IDs that should fail to load
}

// NewMockTaskLoader creates a new mock task loader
func NewMockTaskLoader() *MockTaskLoader {
	return &MockTaskLoader{
		tasks: make(map[int]*model.Task),
		fail:  make(map[int]bool),
	}
}

// AddTask adds a task to the mock loader
func (m *MockTaskLoader) AddTask(task *model.Task) {
	m.tasks[task.ID] = task
}

// FailToLoad marks an ID as unable to load
func (m *MockTaskLoader) FailToLoad(id int) {
	m.fail[id] = true
}

// GetTask implements TaskLoader interface
func (m *MockTaskLoader) GetTask(ctx context.Context, id int) (*model.Task, error) {
	if m.fail[id] {
		return nil, fmt.Errorf("task not found")
	}
	task, ok := m.tasks[id]
	if !ok {
		return nil, fmt.Errorf("task not found")
	}
	return task, nil
}

func TestBuildGraphBasic(t *testing.T) {
	loader := NewMockTaskLoader()

	// Create root task
	rootTask := &model.Task{
		ID:     1,
		Title:  "Root Task",
		Type:   "plan",
		Status: "todo",
	}
	loader.AddTask(rootTask)

	// Build graph
	graph, err := BuildGraph(context.Background(), 1, 4, loader, nil)

	if err != nil {
		t.Fatalf("BuildGraph failed: %v", err)
	}

	if graph.RootTaskID != 1 {
		t.Errorf("expected RootTaskID=1, got %d", graph.RootTaskID)
	}

	if len(graph.Nodes) != 1 {
		t.Errorf("expected 1 node, got %d", len(graph.Nodes))
	}

	// Check root node is marked as center
	root := graph.NodeByID(1)
	if root == nil {
		t.Errorf("expected to find root node")
	}
	if !root.IsCenter {
		t.Errorf("expected root node to have IsCenter=true")
	}
}

func TestBuildGraphWithParent(t *testing.T) {
	loader := NewMockTaskLoader()

	// Create parent task
	parentTask := &model.Task{
		ID:     1,
		Title:  "Parent Task",
		Type:   "plan",
		Status: "todo",
	}

	// Create child task with parent relationship
	childTask := &model.Task{
		ID:     2,
		Title:  "Child Task",
		Type:   "task",
		Status: "todo",
		Relationships: []model.Relationship{
			{Type: model.RelParent, TaskID: 1},
		},
	}

	loader.AddTask(parentTask)
	loader.AddTask(childTask)

	// Build graph starting from child
	graph, err := BuildGraph(context.Background(), 2, 4, loader, nil)

	if err != nil {
		t.Fatalf("BuildGraph failed: %v", err)
	}

	// Should have both parent and child
	if len(graph.Nodes) != 2 {
		t.Errorf("expected 2 nodes, got %d", len(graph.Nodes))
	}

	// Should have parent edge
	if !graph.EdgeExists(2, 1, "parent") {
		t.Errorf("expected parent edge from 2 to 1")
	}
}

func TestBuildGraphWithBlocking(t *testing.T) {
	loader := NewMockTaskLoader()

	// Create blocking task
	blockingTask := &model.Task{
		ID:     1,
		Title:  "Blocking Task",
		Type:   "task",
		Status: "todo",
		Relationships: []model.Relationship{
			{Type: model.RelBlocks, TaskID: 2},
		},
	}

	// Create blocked task
	blockedTask := &model.Task{
		ID:     2,
		Title:  "Blocked Task",
		Type:   "task",
		Status: "todo",
	}

	loader.AddTask(blockingTask)
	loader.AddTask(blockedTask)

	// Build graph starting from blocking task
	graph, err := BuildGraph(context.Background(), 1, 4, loader, nil)

	if err != nil {
		t.Fatalf("BuildGraph failed: %v", err)
	}

	// Should have both tasks
	if len(graph.Nodes) != 2 {
		t.Errorf("expected 2 nodes, got %d", len(graph.Nodes))
	}

	// Should have blocks edge
	if !graph.EdgeExists(1, 2, "blocks") {
		t.Errorf("expected blocks edge from 1 to 2")
	}
}

func TestBuildGraphWithMissingTask(t *testing.T) {
	loader := NewMockTaskLoader()

	// Create task that references missing task
	rootTask := &model.Task{
		ID:     1,
		Title:  "Root Task",
		Type:   "plan",
		Status: "todo",
		Relationships: []model.Relationship{
			{Type: model.RelBlocks, TaskID: 999}, // Non-existent task
		},
	}
	loader.AddTask(rootTask)

	// Build graph
	graph, err := BuildGraph(context.Background(), 1, 4, loader, nil)

	if err != nil {
		t.Fatalf("BuildGraph failed: %v", err)
	}

	// Should have 2 nodes: root and missing node
	if len(graph.Nodes) != 2 {
		t.Errorf("expected 2 nodes (root + missing), got %d", len(graph.Nodes))
	}

	// Missing node should exist
	missing := graph.NodeByID(999)
	if missing == nil {
		t.Errorf("expected to find missing node 999")
	}
	if !missing.IsMissing {
		t.Errorf("expected missing node to have IsMissing=true")
	}

	// Should still have edge
	if !graph.EdgeExists(1, 999, "blocks") {
		t.Errorf("expected blocks edge even to missing task")
	}
}

func TestBuildGraphDepthLimit(t *testing.T) {
	loader := NewMockTaskLoader()

	// Create a chain of tasks: 1 -> 2 -> 3 -> 4 -> 5
	task1 := &model.Task{
		ID:    1,
		Title: "Task 1",
		Type:  "plan",
		Relationships: []model.Relationship{
			{Type: model.RelBlocks, TaskID: 2},
		},
	}
	task2 := &model.Task{
		ID:    2,
		Title: "Task 2",
		Type:  "task",
		Relationships: []model.Relationship{
			{Type: model.RelBlocks, TaskID: 3},
		},
	}
	task3 := &model.Task{
		ID:    3,
		Title: "Task 3",
		Type:  "task",
		Relationships: []model.Relationship{
			{Type: model.RelBlocks, TaskID: 4},
		},
	}
	task4 := &model.Task{
		ID:    4,
		Title: "Task 4",
		Type:  "task",
		Relationships: []model.Relationship{
			{Type: model.RelBlocks, TaskID: 5},
		},
	}
	task5 := &model.Task{
		ID:    5,
		Title: "Task 5",
		Type:  "task",
	}

	loader.AddTask(task1)
	loader.AddTask(task2)
	loader.AddTask(task3)
	loader.AddTask(task4)
	loader.AddTask(task5)

	// Test with depth 2: should get tasks 1, 2, 3
	graph, err := BuildGraph(context.Background(), 1, 2, loader, nil)
	if err != nil {
		t.Fatalf("BuildGraph failed: %v", err)
	}

	if len(graph.Nodes) != 3 {
		t.Errorf("depth=2: expected 3 nodes, got %d", len(graph.Nodes))
	}

	if graph.NodeByID(3) == nil {
		t.Errorf("depth=2: expected task 3 in graph")
	}

	if graph.NodeByID(4) != nil {
		t.Errorf("depth=2: task 4 should not be in graph at depth 2")
	}
}

func TestBuildGraphExcludeTypes(t *testing.T) {
	loader := NewMockTaskLoader()

	// Create task with multiple relationship types
	rootTask := &model.Task{
		ID:    1,
		Title: "Root",
		Type:  "plan",
		Relationships: []model.Relationship{
			{Type: model.RelBlocks, TaskID: 2},
			{Type: model.RelRelatedTo, TaskID: 3},
		},
	}
	task2 := &model.Task{ID: 2, Title: "Task 2", Type: "task"}
	task3 := &model.Task{ID: 3, Title: "Task 3", Type: "task"}

	loader.AddTask(rootTask)
	loader.AddTask(task2)
	loader.AddTask(task3)

	// Build graph excluding "blocks"
	excludeTypes := map[string]bool{"blocks": true}
	graph, err := BuildGraph(context.Background(), 1, 4, loader, excludeTypes)

	if err != nil {
		t.Fatalf("BuildGraph failed: %v", err)
	}

	// Should have root and only relates-to task
	if len(graph.Nodes) != 2 {
		t.Errorf("expected 2 nodes (excluding blocked), got %d", len(graph.Nodes))
	}

	// Should NOT have blocks edge
	if graph.EdgeExists(1, 2, "blocks") {
		t.Errorf("blocks edge should be excluded")
	}

	// Should have relates-to edge
	if !graph.EdgeExists(1, 3, "relates-to") {
		t.Errorf("relates-to edge should exist")
	}
}

func TestBuildGraphCyclePrevention(t *testing.T) {
	loader := NewMockTaskLoader()

	// Create cycle: 1 -> 2 -> 3 -> 1
	task1 := &model.Task{
		ID:    1,
		Title: "Task 1",
		Type:  "task",
		Relationships: []model.Relationship{
			{Type: model.RelBlocks, TaskID: 2},
		},
	}
	task2 := &model.Task{
		ID:    2,
		Title: "Task 2",
		Type:  "task",
		Relationships: []model.Relationship{
			{Type: model.RelBlocks, TaskID: 3},
		},
	}
	task3 := &model.Task{
		ID:    3,
		Title: "Task 3",
		Type:  "task",
		Relationships: []model.Relationship{
			{Type: model.RelBlocks, TaskID: 1},
		},
	}

	loader.AddTask(task1)
	loader.AddTask(task2)
	loader.AddTask(task3)

	// Build graph - should not infinite loop
	graph, err := BuildGraph(context.Background(), 1, 10, loader, nil)

	if err != nil {
		t.Fatalf("BuildGraph failed: %v", err)
	}

	// Should have exactly 3 nodes (cycle handled by visited tracking)
	if len(graph.Nodes) != 3 {
		t.Errorf("expected 3 nodes, got %d", len(graph.Nodes))
	}
}

func TestBuildGraphInvalidInput(t *testing.T) {
	loader := NewMockTaskLoader()

	rootTask := &model.Task{
		ID:    1,
		Title: "Root",
		Type:  "plan",
	}
	loader.AddTask(rootTask)

	// Test depth <= 0
	_, err := BuildGraph(context.Background(), 1, 0, loader, nil)
	if err == nil {
		t.Errorf("expected error for depth <= 0")
	}

	// Test nil loader
	_, err = BuildGraph(context.Background(), 1, 4, nil, nil)
	if err == nil {
		t.Errorf("expected error for nil loader")
	}

	// Test missing root task
	_, err = BuildGraph(context.Background(), 999, 4, loader, nil)
	if err == nil {
		t.Errorf("expected error for missing root task")
	}
}

func TestBuildGraphContextCancellation(t *testing.T) {
	loader := NewMockTaskLoader()

	rootTask := &model.Task{
		ID:    1,
		Title: "Root",
		Type:  "plan",
	}
	loader.AddTask(rootTask)

	// Create canceled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Should fail due to context cancellation
	_, err := BuildGraph(ctx, 1, 4, loader, nil)
	if err == nil {
		t.Errorf("expected error for canceled context")
	}
}

func TestBuildGraphRelatedTo(t *testing.T) {
	loader := NewMockTaskLoader()

	// Create tasks with relates-to relationship
	rootTask := &model.Task{
		ID:    1,
		Title: "Root",
		Type:  "plan",
		Relationships: []model.Relationship{
			{Type: model.RelRelatedTo, TaskID: 2},
		},
	}
	relatedTask := &model.Task{
		ID:    2,
		Title: "Related",
		Type:  "task",
	}

	loader.AddTask(rootTask)
	loader.AddTask(relatedTask)

	graph, err := BuildGraph(context.Background(), 1, 4, loader, nil)
	if err != nil {
		t.Fatalf("BuildGraph failed: %v", err)
	}

	if len(graph.Nodes) != 2 {
		t.Errorf("expected 2 nodes, got %d", len(graph.Nodes))
	}

	if !graph.EdgeExists(1, 2, "relates-to") {
		t.Errorf("expected relates-to edge")
	}
}

func TestBuildGraphComplexHierarchy(t *testing.T) {
	loader := NewMockTaskLoader()

	// Create a hierarchy starting from Plan (2):
	//   Plan (2) [parent: 1]
	//   ├── Parent Epic (1)
	//   └── Task (5) [blocks: 6]
	//       └── Missing (6)

	epic := &model.Task{
		ID:    1,
		Title: "Epic",
		Type:  "epic",
	}

	plan := &model.Task{
		ID:    2,
		Title: "Plan",
		Type:  "plan",
		Relationships: []model.Relationship{
			{Type: model.RelParent, TaskID: 1},
			{Type: model.RelBlocks, TaskID: 5}, // Plan blocks Task 5
		},
	}

	task5 := &model.Task{
		ID:    5,
		Title: "Task 5",
		Type:  "task",
		Relationships: []model.Relationship{
			{Type: model.RelBlocks, TaskID: 6},
		},
	}

	loader.AddTask(epic)
	loader.AddTask(plan)
	loader.AddTask(task5)
	// Task 6 is missing

	// Start from Plan (2)
	graph, err := BuildGraph(context.Background(), 2, 4, loader, nil)
	if err != nil {
		t.Fatalf("BuildGraph failed: %v", err)
	}

	// Should have 4 nodes: Plan (2), Epic (1), Task (5), Missing (6)
	if len(graph.Nodes) != 4 {
		t.Errorf("expected 4 nodes, got %d", len(graph.Nodes))
	}

	// Verify structure
	edges := map[string]bool{
		"2->1:parent": graph.EdgeExists(2, 1, "parent"),
		"2->5:blocks": graph.EdgeExists(2, 5, "blocks"),
		"5->6:blocks": graph.EdgeExists(5, 6, "blocks"),
	}

	for edgeStr, exists := range edges {
		if !exists {
			t.Errorf("expected edge %s to exist", edgeStr)
		}
	}

	// Verify missing node
	missing := graph.NodeByID(6)
	if missing == nil || !missing.IsMissing {
		t.Errorf("expected missing node 6")
	}
}
