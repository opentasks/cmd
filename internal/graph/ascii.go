package graph

import (
	"fmt"
	"sort"
	"strings"
)

// ASCIIRenderer renders a graph as an ASCII tree for terminal display
type ASCIIRenderer struct {
	graph *GraphData
}

// NewASCIIRenderer creates a new ASCII renderer
func NewASCIIRenderer(graph *GraphData) *ASCIIRenderer {
	return &ASCIIRenderer{graph: graph}
}

// Render generates ASCII tree output
func (r *ASCIIRenderer) Render() string {
	if len(r.graph.Nodes) == 0 {
		return "No nodes in graph"
	}

	var output strings.Builder

	// Find root node (the center node)
	root := r.graph.NodeByID(r.graph.RootTaskID)
	if root == nil {
		return fmt.Sprintf("Error: Root task %d not found", r.graph.RootTaskID)
	}

	// Build adjacency list for rendering
	children := r.buildAdjacencyList()

	// Render from root
	output.WriteString(r.formatNode(root, true))
	output.WriteString("\n")
	r.renderSubtree(&output, root.ID, children, "", false)

	return output.String()
}

// buildAdjacencyList creates a map of task ID -> child node IDs
// For parent relationships, the edge goes child->parent, so we need to reverse it
func (r *ASCIIRenderer) buildAdjacencyList() map[int][]int {
	children := make(map[int][]int)

	// Add children based on edges
	for _, edge := range r.graph.Edges {
		if edge.RelationType == "parent" {
			// For parent relationships, reverse the direction:
			// edge goes child->parent, but we want parent->children
			children[edge.TargetID] = append(children[edge.TargetID], edge.SourceID)
		} else {
			// For other relationships, use the edge as-is
			children[edge.SourceID] = append(children[edge.SourceID], edge.TargetID)
		}
	}

	// Sort children for consistent output
	for _, childList := range children {
		sort.Ints(childList)
	}

	return children
}

// renderSubtree recursively renders child nodes
func (r *ASCIIRenderer) renderSubtree(output *strings.Builder, nodeID int, children map[int][]int, prefix string, isLast bool) {
	// Get children for this node
	childList := children[nodeID]
	if len(childList) == 0 {
		return
	}

	for i, childID := range childList {
		child := r.graph.NodeByID(childID)
		if child == nil {
			continue
		}

		isLastChild := (i == len(childList)-1)

		// Build the tree characters
		connector := "├── "
		if isLastChild {
			connector = "└── "
		}

		// Write this node
		output.WriteString(prefix)
		output.WriteString(connector)
		output.WriteString(r.formatNode(child, false))
		output.WriteString("\n")

		// Recursively render children of this child
		var newPrefix string
		if isLastChild {
			newPrefix = prefix + "    "
		} else {
			newPrefix = prefix + "│   "
		}

		r.renderSubtree(output, childID, children, newPrefix, isLastChild)
	}
}

// formatNode formats a single node for display
func (r *ASCIIRenderer) formatNode(node *Node, isRoot bool) string {
	var parts []string

	// Add ID
	parts = append(parts, fmt.Sprintf("[%d]", node.ID))

	// Add title
	if node.Title != "" {
		parts = append(parts, node.Title)
	}

	// Add type in brackets
	if node.Type != "" {
		parts = append(parts, fmt.Sprintf("(%s)", node.Type))
	}

	// Add status in braces
	if node.Status != "" {
		parts = append(parts, fmt.Sprintf("{%s}", node.Status))
	}

	// Add missing indicator
	if node.IsMissing {
		parts = append(parts, "⚠️ MISSING")
	}

	// Add root indicator
	if isRoot {
		parts = append(parts, "★")
	}

	return strings.Join(parts, " ")
}

// RenderASCII is a convenience function that creates a renderer and renders the graph
func RenderASCII(data *GraphData) (string, error) {
	if data == nil {
		return "", fmt.Errorf("data cannot be nil")
	}

	renderer := NewASCIIRenderer(data)
	return renderer.Render(), nil
}
