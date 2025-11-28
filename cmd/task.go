package cmd

import (
	"fmt"

	"github.com/opentasks/cmd/internal/display"
	"github.com/opentasks/cmd/internal/editor"
	"github.com/opentasks/cmd/internal/query"
	"github.com/opentasks/cmd/internal/task"
	"github.com/spf13/cobra"
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

		// Create service
		service := task.NewService(Engine, Store)

		// Create task
		newTask, err := service.Create(ctx, task.CreateRequest{
			Title:       title,
			Type:        taskType,
			Status:      status,
			Description: description,
			ParentID:    parentID,
			Tags:        tags,
		})
		if err != nil {
			return err
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

		// Create service and list tasks
		service := task.NewService(Engine, Store)
		tasks, err := service.List(ctx, filters...)
		if err != nil {
			return err
		}

		// Display results
		if len(tasks) == 0 {
			fmt.Println("No tasks found")
			return nil
		}

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
		taskID, err := task.ParseID(args[0])
		if err != nil {
			return err
		}

		// Create service and get task
		service := task.NewService(Engine, Store)
		taskItem, err := service.Get(ctx, taskID)
		if err != nil {
			return err
		}

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
		taskID, err := task.ParseID(args[0])
		if err != nil {
			return err
		}

		// Create service and delete task
		service := task.NewService(Engine, Store)
		if err := service.Delete(ctx, taskID); err != nil {
			return err
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
		taskID, err := task.ParseID(args[0])
		if err != nil {
			return err
		}

		// Build update request
		req := task.UpdateRequest{}

		if cmd.Flags().Changed("status") {
			status, _ := cmd.Flags().GetString("status")
			req.Status = &status
		}

		if cmd.Flags().Changed("title") {
			title, _ := cmd.Flags().GetString("title")
			req.Title = &title
		}

		if cmd.Flags().Changed("description") {
			description, _ := cmd.Flags().GetString("description")
			req.Description = &description
		}

		if cmd.Flags().Changed("tag") {
			tags, _ := cmd.Flags().GetStringSlice("tag")
			req.AddTags = tags
		}

		if cmd.Flags().Changed("remove-tag") {
			removeTags, _ := cmd.Flags().GetStringSlice("remove-tag")
			req.RemoveTags = removeTags
		}

		// Create service once for efficiency
		service := task.NewService(Engine, Store)

		// Handle editor flag separately (before service call)
		if cmd.Flags().Changed("editor") {
			editorFlag, _ := cmd.Flags().GetBool("editor")
			if editorFlag {
				// Get current task to edit description
				currentTask, err := service.Get(ctx, taskID)
				if err != nil {
					return err
				}

				updatedDescription, err := editor.EditFile(currentTask.Description)
				if err != nil {
					return fmt.Errorf("failed to edit content: %w", err)
				}
				req.Description = &updatedDescription
			}
		}

		// Update task using same service instance
		taskItem, err := service.Update(ctx, taskID, req)
		if err != nil {
			return err
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
