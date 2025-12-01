package display

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/opentasks/cmd/internal/model"
)

// Truncate truncates a string to maxLen characters, adding "..." if truncated
func Truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// FormatPath formats a path for display, using ~ for home directory if applicable
func FormatPath(path string) string {
	home, err := os.UserHomeDir()
	if err == nil && strings.HasPrefix(path, home) {
		return "~" + path[len(home):]
	}
	return path
}

// ConfigFileTree creates a list of resolved config files with merge flow indicators
// Files are shown in priority order (highest first)
func ConfigFileTree(files []string) string {
	var result strings.Builder

	// Build list with all items (config files + defaults)
	allItems := make([]string, len(files)+1)
	for i, file := range files {
		// Get path relative to current directory for most readable display
		cwd, err := os.Getwd()
		var displayPath string

		// First, try to show as relative path from cwd
		if err == nil {
			relPath, err := filepath.Rel(cwd, file)
			if err == nil && !strings.HasPrefix(relPath, "..") {
				// Only use relative path if it doesn't go up many directories
				displayPath = relPath
			}
		}

		// If not a good relative path, try to use ~ for home directory
		if displayPath == "" {
			homeDir, err := os.UserHomeDir()
			if err == nil && strings.HasPrefix(file, homeDir) {
				displayPath = "~" + file[len(homeDir):]
			}
		}

		// Fall back to absolute path if nothing else worked
		if displayPath == "" {
			displayPath = file
		}

		allItems[i] = displayPath
	}
	allItems[len(files)] = "(builtin) defaults"

	// Render as vertical list
	for i, item := range allItems {
		if i == len(allItems)-1 {
			// Last item
			result.WriteString(fmt.Sprintf("└── %s\n", item))
		} else {
			// Not last item
			result.WriteString(fmt.Sprintf("├── %s\n", item))
			result.WriteString("│   ↓\n")
		}
	}

	return result.String()
}

// TaskTable formats tasks as a formatted table for display
func TaskTable(tasks []*model.Task) string {
	if len(tasks) == 0 {
		return ""
	}

	var result strings.Builder

	// Print table header
	result.WriteString(fmt.Sprintf("%-5s %-10s %-30s %-15s %-15s\n", "ID", "Type", "Title", "Status", "Created"))
	result.WriteString(strings.Repeat("-", 75) + "\n")

	// Print tasks
	for _, task := range tasks {
		result.WriteString(fmt.Sprintf("%-5d %-10s %-30s %-15s %-15s\n",
			task.ID,
			task.Type,
			Truncate(task.Title, 30),
			task.Status,
			task.CreatedAt.Format("2006-01-02")))
	}

	return result.String()
}

// TaskDetails formats detailed task information for display
func TaskDetails(task *model.Task) string {
	var result strings.Builder

	result.WriteString(fmt.Sprintf("ID: %d\n", task.ID))
	result.WriteString(fmt.Sprintf("Title: %s\n", task.Title))
	result.WriteString(fmt.Sprintf("Type: %s\n", task.Type))
	result.WriteString(fmt.Sprintf("Status: %s\n", task.Status))
	result.WriteString(fmt.Sprintf("Created: %s\n", task.CreatedAt.Format(time.RFC3339)))
	result.WriteString(fmt.Sprintf("Updated: %s\n", task.UpdatedAt.Format(time.RFC3339)))

	if len(task.Tags) > 0 {
		result.WriteString(fmt.Sprintf("Tags: %s\n", strings.Join(task.Tags, ", ")))
	}

	if len(task.Relationships) > 0 {
		result.WriteString("Relationships:\n")
		for _, rel := range task.Relationships {
			result.WriteString(fmt.Sprintf("  - %s: %d\n", rel.Type, rel.TaskID))
		}
	}

	if task.Description != "" {
		result.WriteString("\nDescription:\n")
		result.WriteString(task.Description)
	}

	return result.String()
}

// PrintOnboardingBox prints a helpful onboarding message when no active project is found.
// This guides users through the options for getting started with opentask.
func PrintOnboardingBox(w io.Writer, cwd string) {
	// Define styles
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Padding(0, 1)

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("170")).
		MarginTop(1)

	commandStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("86")).
		Padding(0, 2)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(1, 2).
		Width(60)

	// Format the current working directory
	displayCwd := FormatPath(cwd)
	if len(displayCwd) > 50 {
		displayCwd = Truncate(displayCwd, 47) + "..."
	}

	// Build content
	var content strings.Builder

	content.WriteString(titleStyle.Render("NO OPENTASK PROJECT FOUND"))
	content.WriteString("\n\n")
	content.WriteString(fmt.Sprintf("Current directory: %s\n", displayCwd))
	content.WriteString("\n")
	content.WriteString("No project configuration found for this location.\n")

	content.WriteString(headerStyle.Render("Choose an option:"))
	content.WriteString("\n\n")

	content.WriteString("1. Initialize a new project here\n")
	content.WriteString(commandStyle.Render("$ opentask config init"))
	content.WriteString("\n\n")

	content.WriteString("2. Attach this directory to an existing project\n")
	content.WriteString(commandStyle.Render("$ opentask project attach <project-id>"))
	content.WriteString("\n")
	content.WriteString(commandStyle.Render("$ opentask project list  # see available projects"))
	content.WriteString("\n\n")

	content.WriteString("3. Work in a directory with an existing project\n")
	content.WriteString(commandStyle.Render("$ cd /path/to/existing/project"))
	content.WriteString("\n\n")

	content.WriteString("Learn more: opentask config --help")

	// Render the box
	fmt.Fprintln(w, boxStyle.Render(content.String()))
}
