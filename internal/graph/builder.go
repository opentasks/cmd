package graph

import (
	"context"
	"fmt"

	"github.com/opentasks/cmd/internal/model"
)

// BFSItem represents a task to process in the BFS queue
type bfsItem struct {
	taskID       int
	currentDepth int
	relationType string // The relationship type that led to this task
}

// BuildGraph constructs a graph of task relationships starting from a root task.
// It uses BFS traversal to include relationships up to the specified depth.
//
// Parameters:
// - ctx: Context for cancellation and deadline support
// - taskID: The starting task ID
// - depth: Maximum relationship depth to traverse (0 means use default)
// - taskLoader: Interface to load tasks
// - excludeTypes: Set of relationship types to skip (e.g., {"blocks": true})
//
// Returns:
// - *GraphData: The constructed graph
// - error: Any error encountered during traversal
func BuildGraph(
	ctx context.Context,
	taskID int,
	depth int,
	taskLoader TaskLoader,
	excludeTypes map[string]bool,
) (*GraphData, error) {
	if depth <= 0 {
		return nil, fmt.Errorf("depth must be greater than 0")
	}

	if taskLoader == nil {
		return nil, fmt.Errorf("taskLoader cannot be nil")
	}

	if excludeTypes == nil {
		excludeTypes = make(map[string]bool)
	}

	// Create the graph structure
	graph := NewGraphData(taskID, depth)

	// Load the root task
	rootTask, err := taskLoader.GetTask(ctx, taskID)
	if err != nil {
		// If we can't load the root task, return error
		return nil, fmt.Errorf("failed to load root task %d: %w", taskID, err)
	}

	// Add root node to the graph
	rootNode := NewNode(rootTask, true)
	graph.AddNode(rootNode)

	// BFS queue: items to process
	// Track visited nodes to avoid cycles and duplicate processing
	queue := []bfsItem{{taskID: taskID, currentDepth: 0}}
	visited := make(map[int]bool)
	visited[taskID] = true

	// BFS traversal
	for len(queue) > 0 {
		// Check context cancellation
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("context canceled: %w", ctx.Err())
		default:
		}

		// Dequeue item
		item := queue[0]
		queue = queue[1:]

		// Load the task (it may be different from root if we're processing via relationships)
		task, err := taskLoader.GetTask(ctx, item.taskID)
		if err != nil {
			// Task not found - create a missing node and continue
			missingNode := NewMissingNode(item.taskID)
			graph.AddNode(missingNode)
			continue
		}

		// Ensure node exists in graph
		if graph.NodeByID(task.ID) == nil {
			node := NewNode(task, task.ID == taskID)
			graph.AddNode(node)
		}

		// If we haven't reached max depth, traverse relationships
		if item.currentDepth < depth {
			// Process parent relationship
			if !excludeTypes["parent"] {
				// Any task type can have a parent
				for _, rel := range task.Relationships {
					if rel.Type == model.RelParent {
						// Process parent relationship
						if !visited[rel.TaskID] {
							queue = append(queue, bfsItem{
								taskID:       rel.TaskID,
								currentDepth: item.currentDepth + 1,
								relationType: model.RelParent,
							})
							visited[rel.TaskID] = true
						}

						// Add edge: this task -> parent
						edge := NewEdge(item.taskID, rel.TaskID, "parent", "parent")
						graph.AddEdge(edge)
					}
				}
			}

			// Process blocking relationships
			if !excludeTypes["blocks"] {
				for _, rel := range task.Relationships {
					if rel.Type == model.RelBlocks {
						// This task blocks another task
						if !visited[rel.TaskID] {
							queue = append(queue, bfsItem{
								taskID:       rel.TaskID,
								currentDepth: item.currentDepth + 1,
								relationType: model.RelBlocks,
							})
							visited[rel.TaskID] = true
						}

						// Add edge: this task -> blocked task
						edge := NewEdge(item.taskID, rel.TaskID, "blocks", "blocks")
						graph.AddEdge(edge)
					}
				}
			}

			// Process related-to relationships
			if !excludeTypes["relates-to"] {
				for _, rel := range task.Relationships {
					if rel.Type == model.RelRelatedTo {
						if !visited[rel.TaskID] {
							queue = append(queue, bfsItem{
								taskID:       rel.TaskID,
								currentDepth: item.currentDepth + 1,
								relationType: model.RelRelatedTo,
							})
							visited[rel.TaskID] = true
						}

						// Add edge: this task -> related task
						edge := NewEdge(item.taskID, rel.TaskID, "relates-to", "relates-to")
						graph.AddEdge(edge)
					}
				}
			}
		}
	}

	return graph, nil
}
