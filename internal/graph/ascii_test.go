package graph

import (
	"strings"
	"testing"

	"github.com/opentasks/cmd/internal/model"
)

func TestASCIIRendererBasic(t *testing.T) {
	graph := NewGraphData(1, 4)

	task := &model.Task{
		ID:     1,
		Title:  "Root Task",
		Type:   "plan",
		Status: "todo",
	}

	graph.AddNode(NewNode(task, true))

	renderer := NewASCIIRenderer(graph)
	output := renderer.Render()

	if !strings.Contains(output, "Root Task") {
		t.Errorf("Expected 'Root Task' in output")
	}

	if !strings.Contains(output, "[1]") {
		t.Errorf("Expected task ID [1] in output")
	}

	if !strings.Contains(output, "(plan)") {
		t.Errorf("Expected type (plan) in output")
	}

	if !strings.Contains(output, "{todo}") {
		t.Errorf("Expected status {todo} in output")
	}

	if !strings.Contains(output, "★") {
		t.Errorf("Expected root indicator ★ in output")
	}
}

func TestASCIIRendererWithChildren(t *testing.T) {
	graph := NewGraphData(1, 4)

	root := &model.Task{
		ID:     1,
		Title:  "Root",
		Type:   "plan",
		Status: "todo",
	}

	child1 := &model.Task{
		ID:     2,
		Title:  "Child 1",
		Type:   "task",
		Status: "todo",
	}

	child2 := &model.Task{
		ID:     3,
		Title:  "Child 2",
		Type:   "task",
		Status: "in-progress",
	}

	graph.AddNode(NewNode(root, true))
	graph.AddNode(NewNode(child1, false))
	graph.AddNode(NewNode(child2, false))

	graph.AddEdge(NewEdge(1, 2, "parent", "parent"))
	graph.AddEdge(NewEdge(1, 3, "parent", "parent"))

	renderer := NewASCIIRenderer(graph)
	output := renderer.Render()

	if !strings.Contains(output, "Root") {
		t.Errorf("Expected 'Root' in output")
	}

	if !strings.Contains(output, "Child 1") {
		t.Errorf("Expected 'Child 1' in output")
	}

	if !strings.Contains(output, "Child 2") {
		t.Errorf("Expected 'Child 2' in output")
	}

	// Check tree structure
	if !strings.Contains(output, "├──") || !strings.Contains(output, "└──") {
		t.Errorf("Expected tree branch characters in output")
	}
}

func TestASCIIRendererWithMissingNode(t *testing.T) {
	graph := NewGraphData(1, 4)

	root := &model.Task{
		ID:     1,
		Title:  "Root",
		Type:   "task",
		Status: "todo",
	}

	graph.AddNode(NewNode(root, true))
	graph.AddNode(NewMissingNode(99))
	graph.AddEdge(NewEdge(1, 99, "blocks", "blocks"))

	renderer := NewASCIIRenderer(graph)
	output := renderer.Render()

	if !strings.Contains(output, "MISSING") {
		t.Errorf("Expected 'MISSING' indicator in output")
	}

	if !strings.Contains(output, "[99]") {
		t.Errorf("Expected missing task ID [99] in output")
	}
}

func TestASCIIRendererEmptyGraph(t *testing.T) {
	graph := NewGraphData(1, 4)

	renderer := NewASCIIRenderer(graph)
	output := renderer.Render()

	if !strings.Contains(output, "No nodes") {
		t.Errorf("Expected 'No nodes' message for empty graph")
	}
}

func TestASCIIRendererNestedChildren(t *testing.T) {
	graph := NewGraphData(1, 4)

	// Create a 3-level hierarchy
	level1 := &model.Task{ID: 1, Title: "Level 1", Type: "plan", Status: "todo"}
	level2 := &model.Task{ID: 2, Title: "Level 2", Type: "task", Status: "todo"}
	level3 := &model.Task{ID: 3, Title: "Level 3", Type: "task", Status: "todo"}

	graph.AddNode(NewNode(level1, true))
	graph.AddNode(NewNode(level2, false))
	graph.AddNode(NewNode(level3, false))

	graph.AddEdge(NewEdge(1, 2, "parent", "parent"))
	graph.AddEdge(NewEdge(2, 3, "parent", "parent"))

	renderer := NewASCIIRenderer(graph)
	output := renderer.Render()

	// Check hierarchy is represented
	if !strings.Contains(output, "Level 1") {
		t.Errorf("Expected 'Level 1' in output")
	}

	if !strings.Contains(output, "Level 2") {
		t.Errorf("Expected 'Level 2' in output")
	}

	if !strings.Contains(output, "Level 3") {
		t.Errorf("Expected 'Level 3' in output")
	}

	// Verify indentation indicates hierarchy
	lines := strings.Split(output, "\n")
	if len(lines) < 3 {
		t.Errorf("Expected at least 3 lines for 3-level hierarchy")
	}
}

func TestRenderASCIIConvenienceFunction(t *testing.T) {
	graph := NewGraphData(1, 4)
	task := &model.Task{ID: 1, Title: "Test", Type: "task", Status: "todo"}
	graph.AddNode(NewNode(task, true))

	output, err := RenderASCII(graph)
	if err != nil {
		t.Fatalf("RenderASCII failed: %v", err)
	}

	if !strings.Contains(output, "Test") {
		t.Errorf("Expected 'Test' in output")
	}
}

func TestRenderASCIINilGraph(t *testing.T) {
	_, err := RenderASCII(nil)
	if err == nil {
		t.Errorf("Expected error for nil graph")
	}
}

func TestASCIIRendererFormatsAllNodeAttributes(t *testing.T) {
	graph := NewGraphData(42, 4)

	task := &model.Task{
		ID:     42,
		Title:  "Complex Task",
		Type:   "research",
		Status: "blocked",
	}

	graph.AddNode(NewNode(task, true))

	renderer := NewASCIIRenderer(graph)
	output := renderer.Render()

	// All attributes should be in output
	if !strings.Contains(output, "[42]") {
		t.Errorf("Expected ID [42]")
	}
	if !strings.Contains(output, "Complex Task") {
		t.Errorf("Expected title")
	}
	if !strings.Contains(output, "(research)") {
		t.Errorf("Expected type")
	}
	if !strings.Contains(output, "{blocked}") {
		t.Errorf("Expected status")
	}
}
