package model

import "time"

// Task represents a trackable work item in opentask
type Task struct {
	// Identity and Metadata
	ID     int    // Global sequential ID (e.g., 42, 5)
	Title  string // Short description
	Type   string // epic|plan|research|story|decision|task
	Status string // Customizable per project (e.g., "todo", "in-progress")

	// Content
	Description string // Markdown body (everything after frontmatter)

	// Organization
	Tags          []string       // Labels (e.g., ["feature", "core", "urgent"])
	Relationships []Relationship // Links to other tasks

	// Timestamps
	CreatedAt time.Time // RFC3339 UTC
	UpdatedAt time.Time // RFC3339 UTC
}

// Task type constants
const (
	TypeEpic     = "epic"
	TypePlan     = "plan"
	TypeResearch = "research"
	TypeStory    = "story"
	TypeDecision = "decision"
	TypeTask     = "task"
)

// AllTaskTypes contains all valid task types
var AllTaskTypes = []string{
	TypeEpic,
	TypePlan,
	TypeResearch,
	TypeStory,
	TypeDecision,
	TypeTask,
}

// TypeCode maps task types to single-character codes for file naming
var TypeCode = map[string]string{
	TypeEpic:     "e",
	TypePlan:     "p",
	TypeResearch: "r",
	TypeStory:    "s",
	TypeDecision: "d",
	TypeTask:     "t",
}

// CodeType maps single-character codes back to task types
var CodeType = map[string]string{
	"e": TypeEpic,
	"p": TypePlan,
	"r": TypeResearch,
	"s": TypeStory,
	"d": TypeDecision,
	"t": TypeTask,
}

// IsValidType returns true if the given type is a valid task type
func IsValidType(t string) bool {
	for _, validType := range AllTaskTypes {
		if t == validType {
			return true
		}
	}
	return false
}
