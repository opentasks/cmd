package model

import (
	"testing"
)

func TestIsValidType(t *testing.T) {
	tests := []struct {
		name     string
		taskType string
		want     bool
	}{
		{"epic type", TypeEpic, true},
		{"plan type", TypePlan, true},
		{"research type", TypeResearch, true},
		{"story type", TypeStory, true},
		{"decision type", TypeDecision, true},
		{"task type", TypeTask, true},
		{"invalid type", "invalid", false},
		{"empty type", "", false},
		{"uppercase type", "EPIC", false},
		{"typo type", "epics", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidType(tt.taskType); got != tt.want {
				t.Errorf("IsValidType(%q) = %v, want %v", tt.taskType, got, tt.want)
			}
		})
	}
}

func TestAllTaskTypes(t *testing.T) {
	// Verify AllTaskTypes contains all valid types
	expectedCount := 6
	if len(AllTaskTypes) != expectedCount {
		t.Errorf("AllTaskTypes has %d types, want %d", len(AllTaskTypes), expectedCount)
	}

	// Verify each type in AllTaskTypes is recognized as valid
	for _, taskType := range AllTaskTypes {
		if !IsValidType(taskType) {
			t.Errorf("Type %q in AllTaskTypes is not recognized as valid", taskType)
		}
	}
}

func TestTypeCode(t *testing.T) {
	tests := []struct {
		taskType string
		code     string
	}{
		{TypeEpic, "e"},
		{TypePlan, "p"},
		{TypeResearch, "r"},
		{TypeStory, "s"},
		{TypeDecision, "d"},
		{TypeTask, "t"},
	}

	for _, tt := range tests {
		t.Run(tt.taskType, func(t *testing.T) {
			if got := TypeCode[tt.taskType]; got != tt.code {
				t.Errorf("TypeCode[%q] = %q, want %q", tt.taskType, got, tt.code)
			}
		})
	}

	// Verify all AllTaskTypes have corresponding codes
	for _, taskType := range AllTaskTypes {
		if _, exists := TypeCode[taskType]; !exists {
			t.Errorf("TypeCode missing entry for %q", taskType)
		}
	}
}

func TestCodeType(t *testing.T) {
	tests := []struct {
		code     string
		taskType string
	}{
		{"e", TypeEpic},
		{"p", TypePlan},
		{"r", TypeResearch},
		{"s", TypeStory},
		{"d", TypeDecision},
		{"t", TypeTask},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			if got := CodeType[tt.code]; got != tt.taskType {
				t.Errorf("CodeType[%q] = %q, want %q", tt.code, got, tt.taskType)
			}
		})
	}
}

func TestTypeCodeRoundTrip(t *testing.T) {
	// Verify that TypeCode -> CodeType is bidirectional
	for _, taskType := range AllTaskTypes {
		code := TypeCode[taskType]
		roundTripped := CodeType[code]
		if roundTripped != taskType {
			t.Errorf("Round trip failed: %s -> %s -> %s", taskType, code, roundTripped)
		}
	}
}
