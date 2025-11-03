package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/zenobi-us/opentask/internal/model"
	"github.com/zenobi-us/opentask/internal/query"
)

// editInEditor opens the content in the user's editor
func editInEditor(initialContent string) (string, error) {
	// Get the editor from environment variable
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi" // Fallback to vi if EDITOR is not set
	}

	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "opentask-*.md")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpFilePath := tmpFile.Name()
	defer os.Remove(tmpFilePath)

	// Write initial content to the temp file
	if _, err := tmpFile.WriteString(initialContent); err != nil {
		tmpFile.Close()
		return "", fmt.Errorf("failed to write temp file: %w", err)
	}
	tmpFile.Close()

	// Launch the editor
	cmd := exec.Command(editor, tmpFilePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("editor command failed: %w", err)
	}

	// Read the updated content
	updatedContent, err := os.ReadFile(tmpFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read updated content: %w", err)
	}

	return string(updatedContent), nil
}

// taskCmd represents the task command group
var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "Manage tasks",
	Long:  "Commands for creating, listing, and managing tasks",
}

// taskNewCmd creates a new task
var taskNewCmd = &cobra.Command{
	Use:   "new [title]",
	Short: "Create a new task",
	Long:  "Create a new task with the given title",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		title := args[0]
		taskType, _ := cmd.Flags().GetString("type")
		status, _ := cmd.Flags().GetString("status")
		description, _ := cmd.Flags().GetString("description")
		parentID, _ := cmd.Flags().GetInt("parent")
		tags, _ := cmd.Flags().GetStringSlice("tag")

		ctx := GetContext()

		// Validate task type
		if !model.IsValidType(taskType) {
			return fmt.Errorf("invalid task type: %s", taskType)
		}

		// Get next ID
		nextID, err := Engine.NextID(ctx)
		if err != nil {
			return fmt.Errorf("failed to generate task ID: %w", err)
		}

		// Determine status: use backlog for high-level tasks without description (unless explicitly set)
		finalStatus := status
		if finalStatus == "" { // Only auto-assign if not explicitly set
			highLevelTypes := []string{"epic", "plan", "research", "story", "decision"}
			isHighLevel := false
			for _, t := range highLevelTypes {
				if t == taskType {
					isHighLevel = true
					break
				}
			}

			// High-level tasks without description start in backlog
			if isHighLevel && description == "" {
				finalStatus = "backlog"
			} else {
				finalStatus = "todo"
			}
		}

		// Create task
		task := &model.Task{
			ID:          nextID,
			Title:       title,
			Type:        taskType,
			Status:      finalStatus,
			Tags:        tags,
			Description: description,
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		}

		// Add parent relationship if specified
		if parentID > 0 {
			task.Relationships = []model.Relationship{
				{
					Type:   model.RelParent,
					TaskID: parentID,
				},
			}
		}

		// Save task
		if err := Store.SaveTask(ctx, task); err != nil {
			return fmt.Errorf("failed to save task: %w", err)
		}

		fmt.Printf("Created task %d: %s\n", task.ID, task.Title)
		return nil
	},
}

// taskListCmd lists tasks
var taskListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List tasks",
	Long:    "List tasks matching the given criteria",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := GetContext()

		// Get filter flags
		status, _ := cmd.Flags().GetString("status")
		taskType, _ := cmd.Flags().GetString("type")
		parentID, _ := cmd.Flags().GetInt("parent")
		tags, _ := cmd.Flags().GetStringSlice("tag")

		// Build filters
		var filters []query.TaskFilter
		if status != "" {
			filters = append(filters, query.WithStatus(status))
		}
		if taskType != "" {
			filters = append(filters, query.WithType(taskType))
		}
		if parentID > 0 {
			filters = append(filters, query.WithParent(parentID))
		}
		if len(tags) > 0 {
			filters = append(filters, query.WithTags(tags))
		}

		// List tasks
		tasks, err := Engine.ListTasks(ctx, filters...)
		if err != nil {
			return fmt.Errorf("failed to list tasks: %w", err)
		}

		// Display results
		if len(tasks) == 0 {
			fmt.Println("No tasks found")
			return nil
		}

		// Print table header
		fmt.Printf("%-5s %-10s %-30s %-15s %-15s\n", "ID", "Type", "Title", "Status", "Created")
		fmt.Println(strings.Repeat("-", 75))

		// Print tasks
		for _, task := range tasks {
			fmt.Printf("%-5d %-10s %-30s %-15s %-15s\n",
				task.ID,
				task.Type,
				truncate(task.Title, 30),
				task.Status,
				task.CreatedAt.Format("2006-01-02"))
		}

		return nil
	},
}

// taskShowCmd shows task details
var taskShowCmd = &cobra.Command{
	Use:     "show [id]",
	Aliases: []string{"view"},
	Short:   "Show task details",
	Long:    "Display detailed information about a task",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := GetContext()

		// Parse task ID
		var taskID int
		_, err := fmt.Sscanf(args[0], "%d", &taskID)
		if err != nil {
			return fmt.Errorf("invalid task ID: %s", args[0])
		}

		// Load task
		task, err := Engine.FindByID(ctx, taskID)
		if err != nil {
			return fmt.Errorf("failed to load task: %w", err)
		}

		// Display task
		fmt.Printf("ID: %d\n", task.ID)
		fmt.Printf("Title: %s\n", task.Title)
		fmt.Printf("Type: %s\n", task.Type)
		fmt.Printf("Status: %s\n", task.Status)
		fmt.Printf("Created: %s\n", task.CreatedAt.Format(time.RFC3339))
		fmt.Printf("Updated: %s\n", task.UpdatedAt.Format(time.RFC3339))

		if len(task.Tags) > 0 {
			fmt.Printf("Tags: %s\n", strings.Join(task.Tags, ", "))
		}

		if len(task.Relationships) > 0 {
			fmt.Println("Relationships:")
			for _, rel := range task.Relationships {
				fmt.Printf("  - %s: %d\n", rel.Type, rel.TaskID)
			}
		}

		if task.Description != "" {
			fmt.Println("\nDescription:")
			fmt.Println(task.Description)
		}

		return nil
	},
}

// taskDeleteCmd deletes a task
var taskDeleteCmd = &cobra.Command{
	Use:     "delete [id]",
	Aliases: []string{"rm"},
	Short:   "Delete a task",
	Long:    "Delete a task by ID",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := GetContext()

		// Parse task ID
		var taskID int
		_, err := fmt.Sscanf(args[0], "%d", &taskID)
		if err != nil {
			return fmt.Errorf("invalid task ID: %s", args[0])
		}

		// Delete task
		if err := Store.DeleteTask(ctx, taskID); err != nil {
			return fmt.Errorf("failed to delete task: %w", err)
		}

		fmt.Printf("Deleted task %d\n", taskID)
		return nil
	},
}

// taskUpdateCmd updates a task's status
var taskUpdateCmd = &cobra.Command{
	Use:     "update [id]",
	Aliases: []string{"edit"},
	Short:   "Update a task",
	Long:    "Update a task's status, tags, or other properties",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := GetContext()

		// Parse task ID
		var taskID int
		_, err := fmt.Sscanf(args[0], "%d", &taskID)
		if err != nil {
			return fmt.Errorf("invalid task ID: %s", args[0])
		}

		// Load task
		task, err := Engine.FindByID(ctx, taskID)
		if err != nil {
			return fmt.Errorf("failed to load task: %w", err)
		}

		// Apply updates
		if cmd.Flags().Changed("status") {
			status, _ := cmd.Flags().GetString("status")
			task.Status = status
		}

		if cmd.Flags().Changed("title") {
			title, _ := cmd.Flags().GetString("title")
			task.Title = title
		}

		if cmd.Flags().Changed("description") {
			description, _ := cmd.Flags().GetString("description")
			task.Description = description
		}

		if cmd.Flags().Changed("tag") {
			tags, _ := cmd.Flags().GetStringSlice("tag")
			// Add tags to existing tags (union)
			tagMap := make(map[string]bool)
			for _, t := range task.Tags {
				tagMap[t] = true
			}
			for _, t := range tags {
				tagMap[t] = true
			}
			// Rebuild tags list
			task.Tags = []string{}
			for t := range tagMap {
				task.Tags = append(task.Tags, t)
			}
		}

		if cmd.Flags().Changed("remove-tag") {
			removeTags, _ := cmd.Flags().GetStringSlice("remove-tag")
			// Remove tags from existing tags
			removeTagMap := make(map[string]bool)
			for _, t := range removeTags {
				removeTagMap[t] = true
			}
			// Rebuild tags list without removed tags
			newTags := []string{}
			for _, t := range task.Tags {
				if !removeTagMap[t] {
					newTags = append(newTags, t)
				}
			}
			task.Tags = newTags
		}

		// Handle editor flag
		if cmd.Flags().Changed("editor") {
			editor, _ := cmd.Flags().GetBool("editor")
			if editor {
				updatedDescription, err := editInEditor(task.Description)
				if err != nil {
					return fmt.Errorf("failed to edit content: %w", err)
				}
				task.Description = updatedDescription
			}
		}

		task.UpdatedAt = time.Now().UTC()

		// Save task
		if err := Store.SaveTask(ctx, task); err != nil {
			return fmt.Errorf("failed to save task: %w", err)
		}

		fmt.Printf("Updated task %d\n", task.ID)
		return nil
	},
}

// Helper function to truncate strings
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func init() {
	// task new flags
	taskNewCmd.Flags().StringP("type", "t", "task", "Task type (epic, plan, research, story, decision, task)")
	taskNewCmd.Flags().StringP("status", "s", "", "Task status (auto-assigned if not specified)")
	taskNewCmd.Flags().StringP("description", "d", "", "Task description")
	taskNewCmd.Flags().IntP("parent", "p", 0, "Parent task ID (for creating subtasks)")
	taskNewCmd.Flags().StringSliceP("tag", "g", []string{}, "Tags to add to the task")

	// task list flags
	taskListCmd.Flags().StringP("status", "s", "", "Filter by status")
	taskListCmd.Flags().StringP("type", "t", "", "Filter by type")
	taskListCmd.Flags().IntP("parent", "p", 0, "Filter by parent task ID")
	taskListCmd.Flags().StringSliceP("tag", "g", []string{}, "Filter by tags")

	// task update flags
	taskUpdateCmd.Flags().StringP("status", "s", "", "New status")
	taskUpdateCmd.Flags().StringP("title", "t", "", "New title")
	taskUpdateCmd.Flags().StringP("description", "d", "", "New description")
	taskUpdateCmd.Flags().StringSliceP("tag", "g", []string{}, "Tags to add to the task")
	taskUpdateCmd.Flags().StringSliceP("remove-tag", "r", []string{}, "Tags to remove from the task")
	taskUpdateCmd.Flags().BoolP("editor", "e", false, "Open task content in $EDITOR for editing")

	// Add subcommands to task command
	taskCmd.AddCommand(taskNewCmd)
	taskCmd.AddCommand(taskListCmd)
	taskCmd.AddCommand(taskShowCmd)
	taskCmd.AddCommand(taskDeleteCmd)
	taskCmd.AddCommand(taskUpdateCmd)
}
