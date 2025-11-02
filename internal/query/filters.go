package query

import (
	"github.com/zenobi-us/opentask/internal/model"
	"github.com/zenobi-us/opentask/internal/storage"
)

// TaskFilter is a functional option for filtering tasks
type TaskFilter = storage.TaskFilter

// WithStatus returns a filter that matches tasks with the given status
func WithStatus(status string) TaskFilter {
	return func(t *model.Task) bool {
		return t.Status == status
	}
}

// WithType returns a filter that matches tasks of the given type
func WithType(taskType string) TaskFilter {
	return func(t *model.Task) bool {
		return t.Type == taskType
	}
}

// WithTag returns a filter that matches tasks with the given tag
func WithTag(tag string) TaskFilter {
	return func(t *model.Task) bool {
		for _, t := range t.Tags {
			if t == tag {
				return true
			}
		}
		return false
	}
}

// WithTags returns a filter that matches tasks with any of the given tags
func WithTags(tags []string) TaskFilter {
	return func(t *model.Task) bool {
		for _, taskTag := range t.Tags {
			for _, filterTag := range tags {
				if taskTag == filterTag {
					return true
				}
			}
		}
		return false
	}
}

// WithStatuses returns a filter that matches tasks with any of the given statuses
func WithStatuses(statuses []string) TaskFilter {
	return func(t *model.Task) bool {
		for _, status := range statuses {
			if t.Status == status {
				return true
			}
		}
		return false
	}
}

// WithParent returns a filter that matches tasks with the given parent epic ID
func WithParent(parentID int) TaskFilter {
	return func(t *model.Task) bool {
		for _, rel := range t.Relationships {
			if rel.Type == model.RelParent && rel.TaskID == parentID {
				return true
			}
		}
		return false
	}
}

// WithRelationship returns a filter that matches tasks with a specific relationship
func WithRelationship(relType string, taskID int) TaskFilter {
	return func(t *model.Task) bool {
		for _, rel := range t.Relationships {
			if rel.Type == relType && rel.TaskID == taskID {
				return true
			}
		}
		return false
	}
}

// WithID returns a filter that matches a specific task ID
func WithID(id int) TaskFilter {
	return func(t *model.Task) bool {
		return t.ID == id
	}
}

// WithTitle returns a filter that matches tasks by title (exact match)
func WithTitle(title string) TaskFilter {
	return func(t *model.Task) bool {
		return t.Title == title
	}
}
