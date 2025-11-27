package config

import (
	"fmt"

	"github.com/opentasks/cmd/internal/project"
)

// ProjectService provides business logic for project CRUD operations on global config.
// This service is stateful - it's instantiated with a global config and mutates it in-place.
type ProjectService struct {
	globalConfig *OpentaskGlobalConfigFile
	pm           *project.Manager
}

// NewProjectService creates a service instance bound to a global config.
// The service will mutate the provided config in-place. Caller is responsible for
// persisting changes via GlobalConfigSaver.
func NewProjectService(cfg *OpentaskGlobalConfigFile) *ProjectService {
	return &ProjectService{
		globalConfig: cfg,
		pm:           project.NewManager(),
	}
}

// AttachContext attaches a working directory to a project.
// The resolvedPath should already be absolute (use Manager.ResolvePath first).
func (s *ProjectService) AttachContext(projectID, resolvedPath string) error {
	proj := s.findProject(projectID)
	if proj == nil {
		return fmt.Errorf("project not found: %s", projectID)
	}

	return proj.AddContextPath(resolvedPath)
}

// DetachContext removes a working directory from a project.
// The resolvedPath should already be absolute (use Manager.ResolvePath first).
func (s *ProjectService) DetachContext(projectID, resolvedPath string) error {
	proj := s.findProject(projectID)
	if proj == nil {
		return fmt.Errorf("project not found: %s", projectID)
	}

	return proj.RemoveContextPath(resolvedPath)
}

// RemoveProject removes a project from the global config.
// If the removed project was the active project, active_project is cleared.
func (s *ProjectService) RemoveProject(projectID string) error {
	projectIndex := s.findProjectIndex(projectID)
	if projectIndex == -1 {
		return fmt.Errorf("project not found: %s", projectID)
	}

	// Remove project from slice
	s.globalConfig.Projects = append(
		s.globalConfig.Projects[:projectIndex],
		s.globalConfig.Projects[projectIndex+1:]...,
	)

	// Clear active project if it was the removed project
	if s.globalConfig.ActiveProject == projectID {
		s.globalConfig.ActiveProject = ""
	}

	return nil
}

// GetProject returns a project by ID, or nil if not found.
// Useful for confirmation prompts in cmd handlers.
func (s *ProjectService) GetProject(projectID string) *GlobalProjectConfig {
	return s.findProject(projectID)
}

// findProject finds a project by ID in the global config.
func (s *ProjectService) findProject(projectID string) *GlobalProjectConfig {
	for i := range s.globalConfig.Projects {
		if s.globalConfig.Projects[i].ID == projectID {
			return &s.globalConfig.Projects[i]
		}
	}
	return nil
}

// findProjectIndex finds the index of a project by ID, or -1 if not found.
func (s *ProjectService) findProjectIndex(projectID string) int {
	for i := range s.globalConfig.Projects {
		if s.globalConfig.Projects[i].ID == projectID {
			return i
		}
	}
	return -1
}
