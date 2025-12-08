package graph

import (
	"testing"

	"github.com/opentasks/cmd/internal/model"
)

func TestNewGraphData(t *testing.T) {
	graph := NewGraphData(42, 4)

	if graph.RootTaskID != 42 {
		t.Errorf("expected RootTaskID=42, got %d", graph.RootTaskID)
	}
	if graph.Depth != 4 {
		t.Errorf("expected Depth=4, got %d", graph.Depth)
	}
	if len(graph.Nodes) != 0 {
		t.Errorf("expected 0 nodes, got %d", len(graph.Nodes))
	}
	if len(graph.Edges) != 0 {
		t.Errorf("expected 0 edges, got %d", len(graph.Edges))
	}
}

func TestNewNode(t *testing.T) {
	task := &model.Task{
		ID:     1,
		Title:  "Test Task",
		Type:   "task",
		Status: "todo",
	}

	node := NewNode(task, true)

	if node.ID != 1 {
		t.Errorf("expected ID=1, got %d", node.ID)
	}
	if node.Title != "Test Task" {
		t.Errorf("expected Title='Test Task', got '%s'", node.Title)
	}
	if node.Type != "task" {
		t.Errorf("expected Type='task', got '%s'", node.Type)
	}
	if node.Status != "todo" {
		t.Errorf("expected Status='todo', got '%s'", node.Status)
	}
	if !node.IsCenter {
		t.Errorf("expected IsCenter=true")
	}
	if node.IsMissing {
		t.Errorf("expected IsMissing=false")
	}
}

func TestNewMissingNode(t *testing.T) {
	node := NewMissingNode(99)

	if node.ID != 99 {
		t.Errorf("expected ID=99, got %d", node.ID)
	}
	if node.Title != "" {
		t.Errorf("expected Title='', got '%s'", node.Title)
	}
	if !node.IsMissing {
		t.Errorf("expected IsMissing=true")
	}
	if node.IsCenter {
		t.Errorf("expected IsCenter=false")
	}
}

func TestNewEdge(t *testing.T) {
	edge := NewEdge(1, 2, "parent", "is parent of")

	if edge.SourceID != 1 {
		t.Errorf("expected SourceID=1, got %d", edge.SourceID)
	}
	if edge.TargetID != 2 {
		t.Errorf("expected TargetID=2, got %d", edge.TargetID)
	}
	if edge.RelationType != "parent" {
		t.Errorf("expected RelationType='parent', got '%s'", edge.RelationType)
	}
	if edge.Label != "is parent of" {
		t.Errorf("expected Label='is parent of', got '%s'", edge.Label)
	}
}

func TestNodeByID(t *testing.T) {
	graph := NewGraphData(1, 4)

	node1 := NewNode(&model.Task{ID: 1, Title: "Task 1", Type: "task", Status: "todo"}, true)
	node2 := NewNode(&model.Task{ID: 2, Title: "Task 2", Type: "task", Status: "todo"}, false)

	graph.AddNode(node1)
	graph.AddNode(node2)

	found := graph.NodeByID(2)
	if found == nil {
		t.Errorf("expected to find node with ID=2")
	}
	if found.Title != "Task 2" {
		t.Errorf("expected found node title='Task 2', got '%s'", found.Title)
	}

	notFound := graph.NodeByID(999)
	if notFound != nil {
		t.Errorf("expected nil for missing node")
	}
}

func TestAddNode(t *testing.T) {
	graph := NewGraphData(1, 4)

	node := NewNode(&model.Task{ID: 1, Title: "Task 1", Type: "task", Status: "todo"}, true)

	// First add should succeed
	added := graph.AddNode(node)
	if !added {
		t.Errorf("expected AddNode to return true for new node")
	}
	if len(graph.Nodes) != 1 {
		t.Errorf("expected 1 node, got %d", len(graph.Nodes))
	}

	// Second add of same node should fail
	added = graph.AddNode(node)
	if added {
		t.Errorf("expected AddNode to return false for duplicate node")
	}
	if len(graph.Nodes) != 1 {
		t.Errorf("expected 1 node after duplicate add, got %d", len(graph.Nodes))
	}
}

func TestEdgeExists(t *testing.T) {
	graph := NewGraphData(1, 4)

	edge := NewEdge(1, 2, "parent", "is parent of")
	graph.AddEdge(edge)

	// Edge should exist
	if !graph.EdgeExists(1, 2, "parent") {
		t.Errorf("expected edge (1->2, parent) to exist")
	}

	// Different type should not exist
	if graph.EdgeExists(1, 2, "blocks") {
		t.Errorf("expected edge (1->2, blocks) to not exist")
	}

	// Different nodes should not exist
	if graph.EdgeExists(2, 1, "parent") {
		t.Errorf("expected edge (2->1, parent) to not exist")
	}
}

func TestAddEdge(t *testing.T) {
	graph := NewGraphData(1, 4)

	edge := NewEdge(1, 2, "parent", "is parent of")

	// First add should succeed
	added := graph.AddEdge(edge)
	if !added {
		t.Errorf("expected AddEdge to return true for new edge")
	}
	if len(graph.Edges) != 1 {
		t.Errorf("expected 1 edge, got %d", len(graph.Edges))
	}

	// Second add of same edge should fail
	added = graph.AddEdge(edge)
	if added {
		t.Errorf("expected AddEdge to return false for duplicate edge")
	}
	if len(graph.Edges) != 1 {
		t.Errorf("expected 1 edge after duplicate add, got %d", len(graph.Edges))
	}

	// Different edge should succeed
	edge2 := NewEdge(1, 2, "blocks", "blocks")
	added = graph.AddEdge(edge2)
	if !added {
		t.Errorf("expected AddEdge to return true for different edge type")
	}
	if len(graph.Edges) != 2 {
		t.Errorf("expected 2 edges, got %d", len(graph.Edges))
	}
}

func TestGraphDataIntegration(t *testing.T) {
	graph := NewGraphData(1, 4)

	// Add multiple nodes
	node1 := NewNode(&model.Task{ID: 1, Title: "Root", Type: "plan", Status: "todo"}, true)
	node2 := NewNode(&model.Task{ID: 2, Title: "Child", Type: "task", Status: "todo"}, false)
	node3 := NewMissingNode(3)

	graph.AddNode(node1)
	graph.AddNode(node2)
	graph.AddNode(node3)

	if len(graph.Nodes) != 3 {
		t.Errorf("expected 3 nodes, got %d", len(graph.Nodes))
	}

	// Add edges
	edge1 := NewEdge(1, 2, "parent", "is parent of")
	edge2 := NewEdge(2, 3, "blocks", "blocks")

	graph.AddEdge(edge1)
	graph.AddEdge(edge2)

	if len(graph.Edges) != 2 {
		t.Errorf("expected 2 edges, got %d", len(graph.Edges))
	}

	// Verify node lookups
	found := graph.NodeByID(2)
	if found == nil || found.Title != "Child" {
		t.Errorf("expected to find Child node")
	}

	missing := graph.NodeByID(3)
	if missing == nil || !missing.IsMissing {
		t.Errorf("expected to find missing node 3")
	}

	// Verify edge existence
	if !graph.EdgeExists(1, 2, "parent") {
		t.Errorf("expected parent edge to exist")
	}
	if !graph.EdgeExists(2, 3, "blocks") {
		t.Errorf("expected blocks edge to exist")
	}
}
