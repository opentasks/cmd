package task

import (
	"time"

	"github.com/zenobi-us/opentask/internal/model"
)

// Manager provides pure business logic for task operations
type Manager struct{}

// NewManager creates a new task manager
func NewManager() *Manager {
	return &Manager{}
}

// DetermineInitialStatus determines the initial status for a new task based on type and description
// High-level tasks (epic, plan, research, story, decision) without description start in backlog
// All other tasks start in todo unless explicitly set
func (m *Manager) DetermineInitialStatus(taskType, description string) string {
	highLevelTypes := map[string]bool{
		model.TypeEpic:     true,
		model.TypePlan:     true,
		model.TypeResearch: true,
		model.TypeStory:    true,
		model.TypeDecision: true,
	}

	// High-level tasks without description start in backlog
	if highLevelTypes[taskType] && description == "" {
		return "backlog"
	}

	return "todo"
}

// CreateTask creates a new task with the given parameters
func (m *Manager) CreateTask(id int, title, taskType, status, description string, tags []string, parentID int) *model.Task {
	now := time.Now().UTC()
	task := &model.Task{
		ID:          id,
		Title:       title,
		Type:        taskType,
		Status:      status,
		Tags:        tags,
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
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

	return task
}

// MergeTags merges existing tags with new tags (union operation)
func (m *Manager) MergeTags(existing, newTags []string) []string {
	tagMap := make(map[string]bool)
	for _, t := range existing {
		tagMap[t] = true
	}
	for _, t := range newTags {
		tagMap[t] = true
	}

	// Rebuild tags list
	result := []string{}
	for t := range tagMap {
		result = append(result, t)
	}
	return result
}

// RemoveTags removes specified tags from the existing tag list
func (m *Manager) RemoveTags(existing, tagsToRemove []string) []string {
	removeMap := make(map[string]bool)
	for _, t := range tagsToRemove {
		removeMap[t] = true
	}

	// Rebuild tags list without removed tags
	result := []string{}
	for _, t := range existing {
		if !removeMap[t] {
			result = append(result, t)
		}
	}
	return result
}
