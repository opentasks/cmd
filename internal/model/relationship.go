package model

// Relationship represents a link between tasks
type Relationship struct {
	Type   string // "parent"|"blocks"|"relates-to"
	TaskID int    // Target task ID (e.g., 42, 5)
}

// Relationship type constants
const (
	RelParent    = "parent"     // Hierarchical parent
	RelBlocks    = "blocks"     // This task blocks another
	RelRelatedTo = "relates-to" // Related but independent
)

// AllRelationshipTypes contains all valid relationship types
var AllRelationshipTypes = []string{
	RelParent,
	RelBlocks,
	RelRelatedTo,
}

// IsValidRelationshipType returns true if the given type is a valid relationship type
func IsValidRelationshipType(t string) bool {
	for _, validType := range AllRelationshipTypes {
		if t == validType {
			return true
		}
	}
	return false
}
