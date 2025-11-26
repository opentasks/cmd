package query

import (
	"testing"
	"time"

	"github.com/opentasks/cmd/internal/model"
)

func createTestTask(id int, title, taskType, status string, tags []string) *model.Task {
	now := time.Now().UTC()
	return &model.Task{
		ID:          id,
		Title:       title,
		Type:        taskType,
		Status:      status,
		Description: "Test",
		Tags:        tags,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func TestWithStatus(t *testing.T) {
	filter := WithStatus("todo")

	task1 := createTestTask(1, "Task 1", model.TypeTask, "todo", []string{})
	task2 := createTestTask(2, "Task 2", model.TypeTask, "in-progress", []string{})

	if !filter(task1) {
		t.Error("WithStatus('todo') should match task with status 'todo'")
	}

	if filter(task2) {
		t.Error("WithStatus('todo') should not match task with status 'in-progress'")
	}
}

func TestWithType(t *testing.T) {
	filter := WithType(model.TypeStory)

	story := createTestTask(1, "Story", model.TypeStory, "todo", []string{})
	task := createTestTask(2, "Task", model.TypeTask, "todo", []string{})

	if !filter(story) {
		t.Error("WithType(story) should match story task")
	}

	if filter(task) {
		t.Error("WithType(story) should not match task type")
	}
}

func TestWithTag(t *testing.T) {
	filter := WithTag("feature")

	task1 := createTestTask(1, "Task 1", model.TypeTask, "todo", []string{"feature", "core"})
	task2 := createTestTask(2, "Task 2", model.TypeTask, "todo", []string{"testing"})
	task3 := createTestTask(3, "Task 3", model.TypeTask, "todo", []string{})

	if !filter(task1) {
		t.Error("WithTag('feature') should match task with 'feature' tag")
	}

	if filter(task2) {
		t.Error("WithTag('feature') should not match task without 'feature' tag")
	}

	if filter(task3) {
		t.Error("WithTag('feature') should not match task with no tags")
	}
}

func TestWithTags(t *testing.T) {
	filter := WithTags([]string{"feature", "core"})

	task1 := createTestTask(1, "Task 1", model.TypeTask, "todo", []string{"feature"})
	task2 := createTestTask(2, "Task 2", model.TypeTask, "todo", []string{"core"})
	task3 := createTestTask(3, "Task 3", model.TypeTask, "todo", []string{"testing"})
	task4 := createTestTask(4, "Task 4", model.TypeTask, "todo", []string{})

	if !filter(task1) {
		t.Error("WithTags should match task with 'feature' tag")
	}

	if !filter(task2) {
		t.Error("WithTags should match task with 'core' tag")
	}

	if filter(task3) {
		t.Error("WithTags should not match task without matching tags")
	}

	if filter(task4) {
		t.Error("WithTags should not match task with no tags")
	}
}

func TestWithStatuses(t *testing.T) {
	filter := WithStatuses([]string{"todo", "in-progress"})

	task1 := createTestTask(1, "Task 1", model.TypeTask, "todo", []string{})
	task2 := createTestTask(2, "Task 2", model.TypeTask, "in-progress", []string{})
	task3 := createTestTask(3, "Task 3", model.TypeTask, "done", []string{})

	if !filter(task1) {
		t.Error("WithStatuses should match task with 'todo' status")
	}

	if !filter(task2) {
		t.Error("WithStatuses should match task with 'in-progress' status")
	}

	if filter(task3) {
		t.Error("WithStatuses should not match task with 'done' status")
	}
}

func TestWithParent(t *testing.T) {
	filter := WithParent(1)

	task1 := &model.Task{
		ID:     1,
		Title:  "Task 1",
		Type:   model.TypeStory,
		Status: "todo",
		Relationships: []model.Relationship{
			{Type: model.RelParent, TaskID: 1},
		},
	}

	task2 := &model.Task{
		ID:     2,
		Title:  "Task 2",
		Type:   model.TypeStory,
		Status: "todo",
		Relationships: []model.Relationship{
			{Type: model.RelParent, TaskID: 2},
		},
	}

	task3 := createTestTask(3, "Task 3", model.TypeStory, "todo", []string{})

	if !filter(task1) {
		t.Error("WithParent(1) should match task with parent 1")
	}

	if filter(task2) {
		t.Error("WithParent(1) should not match task with parent 2")
	}

	if filter(task3) {
		t.Error("WithParent(1) should not match task without parent relationship")
	}
}

func TestWithParentMultipleRelationships(t *testing.T) {
	filter := WithParent(1)

	task := &model.Task{
		ID:     1,
		Title:  "Task with multiple rels",
		Type:   model.TypeStory,
		Status: "todo",
		Relationships: []model.Relationship{
			{Type: model.RelRelatedTo, TaskID: 5},
			{Type: model.RelParent, TaskID: 1},
			{Type: model.RelBlocks, TaskID: 3},
		},
	}

	if !filter(task) {
		t.Error("WithParent(1) should find parent relationship among multiple relationships")
	}
}

func TestWithRelationship(t *testing.T) {
	filter := WithRelationship(model.RelBlocks, 5)

	task1 := &model.Task{
		ID:     1,
		Title:  "Task 1",
		Type:   model.TypeTask,
		Status: "todo",
		Relationships: []model.Relationship{
			{Type: model.RelBlocks, TaskID: 5},
		},
	}

	task2 := &model.Task{
		ID:     2,
		Title:  "Task 2",
		Type:   model.TypeTask,
		Status: "todo",
		Relationships: []model.Relationship{
			{Type: model.RelBlocks, TaskID: 3},
		},
	}

	task3 := &model.Task{
		ID:     3,
		Title:  "Task 3",
		Type:   model.TypeTask,
		Status: "todo",
		Relationships: []model.Relationship{
			{Type: model.RelRelatedTo, TaskID: 5},
		},
	}

	if !filter(task1) {
		t.Error("WithRelationship should match task with matching relationship")
	}

	if filter(task2) {
		t.Error("WithRelationship should not match task with different target ID")
	}

	if filter(task3) {
		t.Error("WithRelationship should not match task with different relationship type")
	}
}

func TestWithID(t *testing.T) {
	filter := WithID(42)

	task1 := createTestTask(42, "Task 42", model.TypeTask, "todo", []string{})
	task2 := createTestTask(41, "Task 41", model.TypeTask, "todo", []string{})

	if !filter(task1) {
		t.Error("WithID(42) should match task with ID 42")
	}

	if filter(task2) {
		t.Error("WithID(42) should not match task with ID 41")
	}
}

func TestWithTitle(t *testing.T) {
	filter := WithTitle("Exact Title")

	task1 := createTestTask(1, "Exact Title", model.TypeTask, "todo", []string{})
	task2 := createTestTask(2, "Different Title", model.TypeTask, "todo", []string{})
	task3 := createTestTask(3, "Exact Title But More", model.TypeTask, "todo", []string{})

	if !filter(task1) {
		t.Error("WithTitle should match task with exact title")
	}

	if filter(task2) {
		t.Error("WithTitle should not match task with different title")
	}

	if filter(task3) {
		t.Error("WithTitle should only match exact title, not partial")
	}
}

func TestMultipleFiltersAND(t *testing.T) {
	// Filters should be ANDed together
	filter1 := WithType(model.TypeStory)
	filter2 := WithStatus("in-progress")
	filter3 := WithTag("feature")

	task := &model.Task{
		ID:          1,
		Title:       "Story Task",
		Type:        model.TypeStory,
		Status:      "in-progress",
		Tags:        []string{"feature", "core"},
		Description: "Test",
	}

	// All filters should match
	if !filter1(task) || !filter2(task) || !filter3(task) {
		t.Error("All filters should match the task")
	}

	// If we change status, some filters won't match
	task.Status = "todo"
	if filter2(task) {
		t.Error("WithStatus filter should not match changed status")
	}
}

func TestEmptyTagsList(t *testing.T) {
	filter := WithTag("feature")
	task := createTestTask(1, "Task", model.TypeTask, "todo", []string{})

	if filter(task) {
		t.Error("WithTag should not match task with empty tags list")
	}
}

func TestEmptyRelationshipsList(t *testing.T) {
	filter := WithParent(1)
	task := createTestTask(1, "Task", model.TypeTask, "todo", []string{})

	if filter(task) {
		t.Error("WithParent should not match task with empty relationships list")
	}
}
