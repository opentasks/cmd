package graph

import (
	"context"

	"github.com/opentasks/cmd/internal/model"
)

// TaskLoader defines the minimal interface needed to load tasks
// for graph traversal. This allows graph package to work with
// different implementations (service, query engine, mock, etc.)
type TaskLoader interface {
	// GetTask loads a task by ID
	// Returns the task or an error if not found or context canceled
	GetTask(ctx context.Context, id int) (*model.Task, error)
}

// GraphData represents the complete graph structure for a task and its relationships
type GraphData struct {
	// RootTaskID is the original task ID that was queried
	RootTaskID int

	// Nodes represents all tasks in the graph
	Nodes []*Node

	// Edges represents all relationships in the graph
	Edges []*Edge

	// Depth is the maximum relationship depth traversed
	Depth int
}

// Node represents a task as a vertex in the graph
type Node struct {
	// ID is the unique task identifier
	ID int

	// Title is the task's short description
	Title string

	// Type is the task type (e.g., "plan", "task", "epic")
	Type string

	// Status is the task's current status
	Status string

	// IsCenter marks if this is the root/queried task
	// Used for visual distinction (highlight in graph)
	IsCenter bool

	// IsMissing indicates the task couldn't be loaded
	// Task is rendered as a leaf node with visual warning
	IsMissing bool
}

// Edge represents a relationship between two tasks
type Edge struct {
	// SourceID is the task that has the relationship
	SourceID int

	// TargetID is the task being related to
	TargetID int

	// RelationType is the type of relationship
	// Values: "parent", "blocks", "blocked-by", "relates-to"
	RelationType string

	// Label is the human-readable label for this edge
	Label string
}

// NewGraphData creates a new GraphData structure
func NewGraphData(rootTaskID int, depth int) *GraphData {
	return &GraphData{
		RootTaskID: rootTaskID,
		Nodes:      make([]*Node, 0),
		Edges:      make([]*Edge, 0),
		Depth:      depth,
	}
}

// NewNode creates a new node from a task
func NewNode(task *model.Task, isCenter bool) *Node {
	return &Node{
		ID:        task.ID,
		Title:     task.Title,
		Type:      task.Type,
		Status:    task.Status,
		IsCenter:  isCenter,
		IsMissing: false,
	}
}

// NewMissingNode creates a node for a task that couldn't be loaded
func NewMissingNode(taskID int) *Node {
	return &Node{
		ID:        taskID,
		Title:     "",
		Type:      "",
		Status:    "",
		IsCenter:  false,
		IsMissing: true,
	}
}

// NewEdge creates a new edge between two tasks
func NewEdge(sourceID, targetID int, relationType, label string) *Edge {
	return &Edge{
		SourceID:     sourceID,
		TargetID:     targetID,
		RelationType: relationType,
		Label:        label,
	}
}

// NodeByID finds a node in the graph by its ID
// Returns nil if not found
func (g *GraphData) NodeByID(id int) *Node {
	for _, n := range g.Nodes {
		if n.ID == id {
			return n
		}
	}
	return nil
}

// AddNode adds a node to the graph if not already present
// Returns true if added, false if already existed
func (g *GraphData) AddNode(node *Node) bool {
	if g.NodeByID(node.ID) != nil {
		return false
	}
	g.Nodes = append(g.Nodes, node)
	return true
}

// EdgeExists checks if an edge already exists
func (g *GraphData) EdgeExists(sourceID, targetID int, relationType string) bool {
	for _, e := range g.Edges {
		if e.SourceID == sourceID && e.TargetID == targetID && e.RelationType == relationType {
			return true
		}
	}
	return false
}

// AddEdge adds an edge to the graph if not already present
// Returns true if added, false if already existed
func (g *GraphData) AddEdge(edge *Edge) bool {
	if g.EdgeExists(edge.SourceID, edge.TargetID, edge.RelationType) {
		return false
	}
	g.Edges = append(g.Edges, edge)
	return true
}
