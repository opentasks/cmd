package cmd

import (
	"fmt"

	"github.com/opentasks/cmd/internal/project"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Display agent initialization guidelines",
	Long:  "Render and display the AGENTS.md guidelines as a template-driven primer for bootstrapping AI agents with project context.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runInit(cmd)
	},
}

func runInit(cmd *cobra.Command) error {
	// Build template context data
	data := map[string]string{
		"ProjectName": "opentasks",
	}

	// Render the agents guide template
	output, err := project.RenderAgentsGuide(data)
	if err != nil {
		return fmt.Errorf("failed to render agents guide: %w", err)
	}

	// Write to stdout
	if _, err := fmt.Fprint(cmd.OutOrStdout(), output); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}
	return nil
}

func init() {
	rootCmd.AddCommand(initCmd)
}
