package config

import (
	"fmt"
	"strings"

	"github.com/zenobi-us/opentask/internal/project"
)

// ProjectLister handles formatting and display of project listings
type ProjectLister struct {
	globalConfig *OpentaskGlobalConfigFile
	pm           *project.Manager
}

// NewProjectLister creates a new ProjectLister for the given global configuration
func NewProjectLister(cfg *OpentaskGlobalConfigFile) *ProjectLister {
	return &ProjectLister{
		globalConfig: cfg,
		pm:           project.NewManager(),
	}
}

// List returns a formatted string representation of all configured projects
func (pl *ProjectLister) List() string {
	if pl.globalConfig == nil || len(pl.globalConfig.Projects) == 0 {
		return ""
	}

	var entries []string
	for _, proj := range pl.globalConfig.Projects {
		entries = append(entries, pl.formatProjectEntry(proj))
	}

	return strings.Join(entries, "\n")
}

// GetActive returns the name of the currently active project, or an empty string if none
func (pl *ProjectLister) GetActive() string {
	if pl.globalConfig == nil {
		return ""
	}
	return pl.globalConfig.ActiveProject
}

// formatProjectEntry formats a single project entry for display
func (pl *ProjectLister) formatProjectEntry(proj GlobalProjectConfig) string {
	// Mark active project
	activeMarker := ""
	if pl.globalConfig != nil && pl.globalConfig.ActiveProject == proj.ID {
		activeMarker = " *"
	}

	// Use project name if available, otherwise use ID
	projectName := proj.Name
	if projectName == "" {
		projectName = proj.ID
	}

	var sb strings.Builder

	// Project header line
	sb.WriteString(fmt.Sprintf("%s (%s)%s\n", proj.ID, projectName, activeMarker))

	// Storage path
	if proj.Storage != nil && proj.Storage.Path != "" {
		displayPath := pl.pm.FormatPathForDisplay(proj.Storage.Path)
		sb.WriteString(fmt.Sprintf("  Storage: %s\n", displayPath))
	}

	// Contexts
	if len(proj.Context) > 0 {
		sb.WriteString("  Contexts:\n")
		for _, ctx := range proj.Context {
			displayPath := pl.pm.FormatPathForDisplay(ctx.Path)
			sb.WriteString(fmt.Sprintf("    - %s\n", displayPath))
		}
	} else {
		sb.WriteString("  Contexts: (none)\n")
	}

	return strings.TrimRight(sb.String(), "\n")
}
