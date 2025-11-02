package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// projectCmd represents the project command group
var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Manage projects",
	Long:  "Commands for creating and managing projects",
}

// projectNewCmd creates a new project
var projectNewCmd = &cobra.Command{
	Use:   "new [name]",
	Short: "Create a new project",
	Long:  "Initialize a new OpenTasks project",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		fmt.Printf("Creating new project: %s\n", name)
		// TODO: Implement project initialization
		return nil
	},
}

// projectListCmd lists projects
var projectListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List projects",
	Long:    "List all available projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Listing projects...")
		// TODO: Implement project listing
		return nil
	},
}

func init() {
	projectCmd.AddCommand(projectNewCmd)
	projectCmd.AddCommand(projectListCmd)
}
