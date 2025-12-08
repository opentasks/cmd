package graph

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/opentasks/cmd/internal/model"
)

func TestRenderDOTNilData(t *testing.T) {
	_, err := RenderDOT(nil, "")
	if err == nil {
		t.Errorf("expected error for nil data")
	}
}

func TestRenderDOTBasicGraph(t *testing.T) {
	// Create a simple graph
	graph := NewGraphData(1, 4)

	task := &model.Task{
		ID:     1,
		Title:  "Test Task",
		Type:   "task",
		Status: "todo",
	}
	node := NewNode(task, true)
	graph.AddNode(node)

	// Render the graph
	dot, err := RenderDOT(graph, "")
	if err != nil {
		t.Fatalf("RenderDOT failed: %v", err)
	}

	// Check that it's valid DOT output
	if !strings.Contains(dot, "digraph") {
		t.Errorf("expected 'digraph' in output")
	}

	if !strings.Contains(dot, "task_1") {
		t.Errorf("expected 'task_1' node in output")
	}

	if !strings.Contains(dot, "rankdir=TB") {
		t.Errorf("expected 'rankdir=TB' in output")
	}
}

func TestRenderDOTCenterNode(t *testing.T) {
	graph := NewGraphData(1, 4)

	// Create center node
	centerTask := &model.Task{
		ID:     1,
		Title:  "Center",
		Type:   "plan",
		Status: "todo",
	}
	centerNode := NewNode(centerTask, true)
	graph.AddNode(centerNode)

	dot, err := RenderDOT(graph, "")
	if err != nil {
		t.Fatalf("RenderDOT failed: %v", err)
	}

	// Center node should have special styling
	if !strings.Contains(dot, `fillcolor="#e8f4f8"`) {
		t.Errorf("expected center node color styling in output")
	}

	if !strings.Contains(dot, `style="rounded,filled"`) {
		t.Errorf("expected center node style in output")
	}
}

func TestRenderDOTMissingNode(t *testing.T) {
	graph := NewGraphData(1, 4)

	// Add root node
	rootTask := &model.Task{
		ID:     1,
		Title:  "Root",
		Type:   "task",
		Status: "todo",
	}
	rootNode := NewNode(rootTask, true)
	graph.AddNode(rootNode)

	// Add missing node
	missingNode := NewMissingNode(999)
	graph.AddNode(missingNode)

	// Add edge
	edge := NewEdge(1, 999, "blocks", "blocks")
	graph.AddEdge(edge)

	dot, err := RenderDOT(graph, "")
	if err != nil {
		t.Fatalf("RenderDOT failed: %v", err)
	}

	// Missing node should have special styling
	if !strings.Contains(dot, `fillcolor="#ffcccc"`) {
		t.Errorf("expected missing node color styling in output")
	}

	if !strings.Contains(dot, `style="rounded,filled,dashed"`) {
		t.Errorf("expected missing node dashed styling in output")
	}

	if !strings.Contains(dot, "<missing>") {
		t.Errorf("expected '<missing>' label in output")
	}
}

func TestRenderDOTBlockingEdge(t *testing.T) {
	graph := NewGraphData(1, 4)

	// Add nodes
	task1 := &model.Task{ID: 1, Title: "Task 1", Type: "task", Status: "todo"}
	task2 := &model.Task{ID: 2, Title: "Task 2", Type: "task", Status: "todo"}

	graph.AddNode(NewNode(task1, true))
	graph.AddNode(NewNode(task2, false))

	// Add blocking edge
	edge := NewEdge(1, 2, "blocks", "blocks")
	graph.AddEdge(edge)

	dot, err := RenderDOT(graph, "")
	if err != nil {
		t.Fatalf("RenderDOT failed: %v", err)
	}

	// Blocking edge should have special styling
	if !strings.Contains(dot, `style=dashed color=red`) {
		t.Errorf("expected dashed red style for blocking edge")
	}

	if !strings.Contains(dot, `label="blocks"`) {
		t.Errorf("expected blocks label")
	}
}

func TestRenderDOTRelatedEdge(t *testing.T) {
	graph := NewGraphData(1, 4)

	// Add nodes
	task1 := &model.Task{ID: 1, Title: "Task 1", Type: "task", Status: "todo"}
	task2 := &model.Task{ID: 2, Title: "Task 2", Type: "task", Status: "todo"}

	graph.AddNode(NewNode(task1, true))
	graph.AddNode(NewNode(task2, false))

	// Add related edge
	edge := NewEdge(1, 2, "relates-to", "relates-to")
	graph.AddEdge(edge)

	dot, err := RenderDOT(graph, "")
	if err != nil {
		t.Fatalf("RenderDOT failed: %v", err)
	}

	// Related edge should have dotted style
	if !strings.Contains(dot, `style=dotted`) {
		t.Errorf("expected dotted style for relates-to edge")
	}

	if !strings.Contains(dot, `label="relates-to"`) {
		t.Errorf("expected relates-to label")
	}
}

func TestRenderDOTParentEdge(t *testing.T) {
	graph := NewGraphData(1, 4)

	// Add nodes
	parent := &model.Task{ID: 1, Title: "Parent", Type: "plan", Status: "todo"}
	child := &model.Task{ID: 2, Title: "Child", Type: "task", Status: "todo"}

	graph.AddNode(NewNode(parent, false))
	graph.AddNode(NewNode(child, true))

	// Add parent edge
	edge := NewEdge(2, 1, "parent", "parent")
	graph.AddEdge(edge)

	dot, err := RenderDOT(graph, "")
	if err != nil {
		t.Fatalf("RenderDOT failed: %v", err)
	}

	// Parent edge should have normal style (no special styling)
	if !strings.Contains(dot, `label="parent"`) {
		t.Errorf("expected parent label")
	}

	// Make sure we have the edge
	if !strings.Contains(dot, "task_2 -> task_1") {
		t.Errorf("expected parent edge from task_2 to task_1")
	}
}

func TestRenderDOTComplexGraph(t *testing.T) {
	graph := NewGraphData(1, 4)

	// Create hierarchy: 1 (plan) -> 2 (task) -> 3 (task)
	task1 := &model.Task{ID: 1, Title: "Epic", Type: "plan", Status: "todo"}
	task2 := &model.Task{ID: 2, Title: "Task 2", Type: "task", Status: "todo"}
	task3 := &model.Task{ID: 3, Title: "Task 3", Type: "task", Status: "todo"}

	graph.AddNode(NewNode(task1, true))
	graph.AddNode(NewNode(task2, false))
	graph.AddNode(NewNode(task3, false))

	// Add edges
	graph.AddEdge(NewEdge(1, 2, "blocks", "blocks"))
	graph.AddEdge(NewEdge(2, 3, "relates-to", "relates-to"))

	dot, err := RenderDOT(graph, "")
	if err != nil {
		t.Fatalf("RenderDOT failed: %v", err)
	}

	// Check all nodes are present
	for i := 1; i <= 3; i++ {
		if !strings.Contains(dot, "task_"+string(rune('0'+i))) {
			t.Errorf("expected task_%d in output", i)
		}
	}

	// Check edges
	if !strings.Contains(dot, "task_1 -> task_2") {
		t.Errorf("expected edge from task_1 to task_2")
	}

	if !strings.Contains(dot, "task_2 -> task_3") {
		t.Errorf("expected edge from task_2 to task_3")
	}
}

func TestRenderDOTNodeLabels(t *testing.T) {
	graph := NewGraphData(1, 4)

	task := &model.Task{
		ID:     42,
		Title:  "My Task Title",
		Type:   "plan",
		Status: "in-progress",
	}

	graph.AddNode(NewNode(task, false))

	dot, err := RenderDOT(graph, "")
	if err != nil {
		t.Fatalf("RenderDOT failed: %v", err)
	}

	// Check label contains ID
	if !strings.Contains(dot, "42") {
		t.Errorf("expected task ID in label")
	}

	// Check label contains title
	if !strings.Contains(dot, "My Task Title") {
		t.Errorf("expected task title in label")
	}

	// Check label contains type
	if !strings.Contains(dot, "plan") {
		t.Errorf("expected task type in label")
	}
}

func TestRenderDOTCustomTemplatesNotFound(t *testing.T) {
	graph := NewGraphData(1, 4)

	task := &model.Task{ID: 1, Title: "Test", Type: "task", Status: "todo"}
	graph.AddNode(NewNode(task, true))

	// Try to use non-existent template directory
	_, err := RenderDOT(graph, "/non/existent/path")
	if err == nil {
		t.Errorf("expected error for non-existent template directory")
	}
}

func TestRenderDOTCustomTemplates(t *testing.T) {
	// Create a temporary directory with custom templates
	tmpDir, err := os.MkdirTemp("", "graph-templates-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create custom graph template with define block
	graphTmpl := `{{- define "graph" }}
digraph {
  rankdir=LR;
  {{- range .Nodes }}
  {{ template "node" . }}
  {{- end }}
  {{- range .Edges }}
  {{ template "edge" . }}
  {{- end }}
}
{{- end }}`

	// Create custom node template with define block
	nodeTmpl := `{{- define "node" }}
{{ printf "task_%d" .ID }} [label="{{ .ID }}"];
{{- end }}`

	// Create custom edge template with define block
	edgeTmpl := `{{- define "edge" }}
{{ printf "task_%d" .SourceID }} -> {{ printf "task_%d" .TargetID }};
{{- end }}`

	if err := os.WriteFile(filepath.Join(tmpDir, "graph.dot"), []byte(graphTmpl), 0644); err != nil {
		t.Fatalf("failed to write graph template: %v", err)
	}

	if err := os.WriteFile(filepath.Join(tmpDir, "node.dot"), []byte(nodeTmpl), 0644); err != nil {
		t.Fatalf("failed to write node template: %v", err)
	}

	if err := os.WriteFile(filepath.Join(tmpDir, "edge.dot"), []byte(edgeTmpl), 0644); err != nil {
		t.Fatalf("failed to write edge template: %v", err)
	}

	// Create a simple graph
	graph := NewGraphData(1, 4)
	task := &model.Task{ID: 1, Title: "Test", Type: "task", Status: "todo"}
	graph.AddNode(NewNode(task, true))

	// Render with custom templates
	dot, err := RenderDOT(graph, tmpDir)
	if err != nil {
		t.Fatalf("RenderDOT with custom templates failed: %v", err)
	}

	// Check that it uses the custom direction
	if !strings.Contains(dot, "rankdir=LR") {
		t.Errorf("expected custom rankdir=LR in output")
	}
}

func TestRenderDOTEmptyGraph(t *testing.T) {
	graph := NewGraphData(1, 4)

	dot, err := RenderDOT(graph, "")
	if err != nil {
		t.Fatalf("RenderDOT failed: %v", err)
	}

	// Should still be valid DOT
	if !strings.Contains(dot, "digraph") {
		t.Errorf("expected 'digraph' in output")
	}
}

func TestRenderDOTValidSyntax(t *testing.T) {
	graph := NewGraphData(1, 4)

	// Add multiple nodes and edges
	for i := 1; i <= 5; i++ {
		task := &model.Task{
			ID:     i,
			Title:  "Task " + string(rune('0'+i)),
			Type:   "task",
			Status: "todo",
		}
		graph.AddNode(NewNode(task, i == 1))
	}

	// Add edges
	graph.AddEdge(NewEdge(1, 2, "blocks", "blocks"))
	graph.AddEdge(NewEdge(2, 3, "parent", "parent"))
	graph.AddEdge(NewEdge(3, 4, "relates-to", "relates-to"))

	dot, err := RenderDOT(graph, "")
	if err != nil {
		t.Fatalf("RenderDOT failed: %v", err)
	}

	// Basic syntax checks
	if !strings.Contains(dot, "digraph {") {
		t.Errorf("expected 'digraph {' in output")
	}

	if !strings.HasSuffix(strings.TrimSpace(dot), "}") {
		t.Errorf("expected closing brace at end of output")
	}

	// Check that all nodes are properly quoted
	nodeCount := strings.Count(dot, "task_")
	if nodeCount < 4 {
		t.Errorf("expected at least 4 nodes in output, got %d", nodeCount)
	}

	// Check that all edges are properly formatted
	arrowCount := strings.Count(dot, "->")
	if arrowCount < 3 {
		t.Errorf("expected at least 3 edges in output, got %d", arrowCount)
	}
}
