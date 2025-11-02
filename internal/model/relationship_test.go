package model

import (
	"testing"
)

func TestIsValidRelationshipType(t *testing.T) {
	tests := []struct {
		name    string
		relType string
		want    bool
	}{
		{"parent relation", RelParent, true},
		{"blocks relation", RelBlocks, true},
		{"relates-to relation", RelRelatedTo, true},
		{"invalid relation", "invalid", false},
		{"empty relation", "", false},
		{"uppercase relation", "PARENT", false},
		{"typo relation", "parents", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidRelationshipType(tt.relType); got != tt.want {
				t.Errorf("IsValidRelationshipType(%q) = %v, want %v", tt.relType, got, tt.want)
			}
		})
	}
}

func TestAllRelationshipTypes(t *testing.T) {
	// Verify AllRelationshipTypes contains all valid types
	expectedCount := 3
	if len(AllRelationshipTypes) != expectedCount {
		t.Errorf("AllRelationshipTypes has %d types, want %d", len(AllRelationshipTypes), expectedCount)
	}

	// Verify each type in AllRelationshipTypes is recognized as valid
	for _, relType := range AllRelationshipTypes {
		if !IsValidRelationshipType(relType) {
			t.Errorf("Type %q in AllRelationshipTypes is not recognized as valid", relType)
		}
	}
}

func TestRelationshipCreation(t *testing.T) {
	rel := Relationship{
		Type:   RelParent,
		TaskID: 42,
	}

	if rel.Type != RelParent {
		t.Errorf("Relationship.Type = %q, want %q", rel.Type, RelParent)
	}

	if rel.TaskID != 42 {
		t.Errorf("Relationship.TaskID = %d, want %d", rel.TaskID, 42)
	}
}

func TestRelationshipTypeConstants(t *testing.T) {
	if RelParent != "parent" {
		t.Errorf("RelParent = %q, want 'parent'", RelParent)
	}

	if RelBlocks != "blocks" {
		t.Errorf("RelBlocks = %q, want 'blocks'", RelBlocks)
	}

	if RelRelatedTo != "relates-to" {
		t.Errorf("RelRelatedTo = %q, want 'relates-to'", RelRelatedTo)
	}
}
