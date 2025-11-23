package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/zenobi-us/opentask/internal/display"
	"github.com/zenobi-us/opentask/internal/editor"
	"github.com/zenobi-us/opentask/internal/model"
	"github.com/zenobi-us/opentask/internal/query"
	"github.com/zenobi-us/opentask/internal/task"
)

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

		// Determine status using task manager
		taskManager := task.NewManager()
		finalStatus := status
		if finalStatus == "" {
			finalStatus = taskManager.DetermineInitialStatus(taskType, description)
		}

		// Create task using task manager
		newTask := taskManager.CreateTask(nextID, title, taskType, finalStatus, description, tags, parentID)

		// Save task
		if err := Store.SaveTask(ctx, newTask); err != nil {
			return fmt.Errorf("failed to save task: %w", err)
		}

		fmt.Printf("Created task %d: %s\n", newTask.ID, newTask.Title)
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

		// Use display package for formatting
		fmt.Print(display.TaskTable(tasks))

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
		taskItem, err := Engine.FindByID(ctx, taskID)
		if err != nil {
			return fmt.Errorf("failed to load task: %w", err)
		}

		// Use display package for formatting
		fmt.Print(display.TaskDetails(taskItem))

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
		taskItem, err := Engine.FindByID(ctx, taskID)
		if err != nil {
			return fmt.Errorf("failed to load task: %w", err)
		}

		// Apply updates
		if cmd.Flags().Changed("status") {
			status, _ := cmd.Flags().GetString("status")
			taskItem.Status = status
		}

		if cmd.Flags().Changed("title") {
			title, _ := cmd.Flags().GetString("title")
			taskItem.Title = title
		}

		if cmd.Flags().Changed("description") {
			description, _ := cmd.Flags().GetString("description")
			taskItem.Description = description
		}

		// Use task manager for tag operations
		taskManager := task.NewManager()
		if cmd.Flags().Changed("tag") {
			tags, _ := cmd.Flags().GetStringSlice("tag")
			taskItem.Tags = taskManager.MergeTags(taskItem.Tags, tags)
		}

		if cmd.Flags().Changed("remove-tag") {
			removeTags, _ := cmd.Flags().GetStringSlice("remove-tag")
			taskItem.Tags = taskManager.RemoveTags(taskItem.Tags, removeTags)
		}

		// Handle editor flag
		if cmd.Flags().Changed("editor") {
			editorFlag, _ := cmd.Flags().GetBool("editor")
			if editorFlag {
				updatedDescription, err := editor.EditFile(taskItem.Description)
				if err != nil {
					return fmt.Errorf("failed to edit content: %w", err)
				}
				taskItem.Description = updatedDescription
			}
		}

		// Update timestamp
		taskItem.UpdatedAt = time.Now().UTC()

		// Save task
		if err := Store.SaveTask(ctx, taskItem); err != nil {
			return fmt.Errorf("failed to save task: %w", err)
		}

		fmt.Printf("Updated task %d\n", taskItem.ID)
		return nil
	},
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
