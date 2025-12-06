package e2e

import (
	"testing"

	"github.com/opentasks/cmd/internal/model"
	"github.com/opentasks/cmd/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestE2E_FourLevelHierarchy tests a complete 4-level task hierarchy
// Epic → Story → Task → Subtask with relationship and query validation
func TestE2E_FourLevelHierarchy(t *testing.T) {
	env := SetupE2EEnvironment(t)
	defer env.Cleanup()

	t.Run("create four-level hierarchy", func(t *testing.T) {
		// Create epic
		epic := testutil.CreateEpic(t, env.Store,
			testutil.WithTitle("Authentication System"),
			testutil.WithStatus("active"),
		)

		Assert(t, epic, env).
			HasTitle("Authentication System").
			HasType(model.TypeEpic).
			HasStatus("active")

		// Create story under epic
		story := testutil.CreateStory(t, env.Store,
			testutil.WithTitle("OAuth Integration"),
			testutil.WithParent(epic.ID),
		)

		Assert(t, story, env).
			HasTitle("OAuth Integration").
			HasType(model.TypeStory).
			HasParent(epic.ID)

		// Verify relationship helper
		testutil.AssertRelationship(t, story, model.RelParent, epic.ID)

		// Create task under story
		task := testutil.CreateTask(t, env.Store,
			testutil.WithTitle("Implement OAuth Flow"),
			testutil.WithParent(story.ID),
		)

		Assert(t, task, env).
			HasTitle("Implement OAuth Flow").
			HasType(model.TypeTask).
			HasParent(story.ID)

		testutil.AssertRelationship(t, task, model.RelParent, story.ID)

		// Create subtask under task
		subtask := testutil.CreateTask(t, env.Store,
			testutil.WithTitle("Add token refresh"),
			testutil.WithParent(task.ID),
		)

		Assert(t, subtask, env).
			HasTitle("Add token refresh").
			HasType(model.TypeTask).
			HasParent(task.ID)

		testutil.AssertRelationship(t, subtask, model.RelParent, task.ID)

		// Verify queries at each level
		Assert(t, epic, env).HasChildCount(1)
		Assert(t, story, env).HasChildCount(1)
		Assert(t, task, env).HasChildCount(1)
		Assert(t, subtask, env).HasChildCount(0)

		// Verify file system structure (only for file storage)
		if env.TmpDir != "" {
			Assert(t, epic, env).FileExists()
			Assert(t, story, env).FileExists()
			Assert(t, task, env).FileExists()
			Assert(t, subtask, env).FileExists()
		}
	})

	t.Run("query children returns correct results", func(t *testing.T) {
		// Create test hierarchy
		epic := testutil.CreateEpic(t, env.Store,
			testutil.WithTitle("Test Epic"),
		)

		// Create 3 stories under epic
		stories := testutil.BulkCreateTasks(t, env.Store, 3, model.TypeStory,
			testutil.WithParent(epic.ID),
		)

		// Verify we can query all children
		children, err := env.Engine.FindChildren(env.Ctx, epic.ID)
		require.NoError(t, err)
		assert.Len(t, children, 3, "Epic should have 3 story children")

		// Verify each story is in the children list
		for i, story := range stories {
			found := false
			for _, child := range children {
				if child.ID == story.ID {
					found = true
					break
				}
			}
			assert.True(t, found, "Story %d should be in children list", i+1)
		}
	})

	t.Run("deep hierarchy traversal", func(t *testing.T) {
		// Create: Epic → 2 Stories → 2 Tasks each = 5 total tasks
		epic := testutil.CreateEpic(t, env.Store,
			testutil.WithTitle("Feature X"),
		)

		story1 := testutil.CreateStory(t, env.Store,
			testutil.WithTitle("Story 1"),
			testutil.WithParent(epic.ID),
		)

		story2 := testutil.CreateStory(t, env.Store,
			testutil.WithTitle("Story 2"),
			testutil.WithParent(epic.ID),
		)

		// 2 tasks under story1
		testutil.BulkCreateTasks(t, env.Store, 2, model.TypeTask,
			testutil.WithParent(story1.ID),
		)

		// 2 tasks under story2
		testutil.BulkCreateTasks(t, env.Store, 2, model.TypeTask,
			testutil.WithParent(story2.ID),
		)

		// Verify epic has 2 children (stories)
		Assert(t, epic, env).HasChildCount(2)

		// Verify each story has 2 children (tasks)
		Assert(t, story1, env).HasChildCount(2)
		Assert(t, story2, env).HasChildCount(2)

		// Each subtest shares the same environment, so we just verify the relationships
		// rather than counting all tasks in storage
	})
}

// TestE2E_EpicDeletionOrphaning tests current behavior when deleting parent tasks
// Documents the orphaning issue for future cascade delete implementation
func TestE2E_EpicDeletionOrphaning(t *testing.T) {
	env := SetupE2EEnvironment(t)
	defer env.Cleanup()

	t.Run("epic deletion leaves orphaned children", func(t *testing.T) {
		// Create epic with 3 children
		epic := testutil.CreateEpic(t, env.Store,
			testutil.WithTitle("Feature X"),
		)

		story1 := testutil.CreateStory(t, env.Store,
			testutil.WithTitle("Story 1"),
			testutil.WithParent(epic.ID),
		)
		story2 := testutil.CreateStory(t, env.Store,
			testutil.WithTitle("Story 2"),
			testutil.WithParent(epic.ID),
		)
		task := testutil.CreateTask(t, env.Store,
			testutil.WithTitle("Task 1"),
			testutil.WithParent(epic.ID),
		)

		// Verify children exist
		Assert(t, epic, env).HasChildCount(3)

		// Delete epic
		err := env.Store.DeleteTask(env.Ctx, epic.ID)
		require.NoError(t, err)

		// Verify epic deleted
		_, err = env.Store.LoadTask(env.Ctx, epic.ID)
		assert.Error(t, err, "Epic should be deleted")

		// CURRENT BEHAVIOR: Children still exist
		loaded1, err := env.Store.LoadTask(env.Ctx, story1.ID)
		assert.NoError(t, err, "Story 1 should still exist (orphaned)")

		loaded2, err := env.Store.LoadTask(env.Ctx, story2.ID)
		assert.NoError(t, err, "Story 2 should still exist (orphaned)")

		loadedTask, err := env.Store.LoadTask(env.Ctx, task.ID)
		assert.NoError(t, err, "Task should still exist (orphaned)")

		// DOCUMENT ORPHAN ISSUE
		t.Logf("⚠️  ORPHAN DETECTED: Story %d has parent=%d but parent doesn't exist",
			loaded1.ID, epic.ID)
		t.Logf("⚠️  ORPHAN DETECTED: Story %d has parent=%d but parent doesn't exist",
			loaded2.ID, epic.ID)
		t.Logf("⚠️  ORPHAN DETECTED: Task %d has parent=%d but parent doesn't exist",
			loadedTask.ID, epic.ID)

		// CURRENT BEHAVIOR: FindChildren still returns orphaned children!
		// The query engine doesn't verify the parent exists
		orphans, err := env.Engine.FindChildren(env.Ctx, epic.ID)
		assert.NoError(t, err)
		assert.Len(t, orphans, 3, "FindChildren still returns orphans (BUG)")

		t.Logf("⚠️  BUG: FindChildren returned %d orphans for deleted parent %d", len(orphans), epic.ID)

		// But relationships still exist in tasks
		testutil.AssertRelationship(t, loaded1, model.RelParent, epic.ID)
		testutil.AssertRelationship(t, loaded2, model.RelParent, epic.ID)
		testutil.AssertRelationship(t, loadedTask, model.RelParent, epic.ID)
	})

	t.Run("blocks relationship survives task deletion", func(t *testing.T) {
		// Create task A that blocks task B
		taskA := testutil.CreateTask(t, env.Store,
			testutil.WithTitle("Blocking Task"),
		)

		taskB := testutil.CreateTask(t, env.Store,
			testutil.WithTitle("Blocked Task"),
		)

		// Update taskA to block taskB
		taskA.Relationships = append(taskA.Relationships, model.Relationship{
			Type:   model.RelBlocks,
			TaskID: taskB.ID,
		})
		err := env.Store.SaveTask(env.Ctx, taskA)
		require.NoError(t, err)

		// Verify blocks relationship
		reloadedA, err := env.Store.LoadTask(env.Ctx, taskA.ID)
		require.NoError(t, err)
		testutil.AssertRelationship(t, reloadedA, model.RelBlocks, taskB.ID)

		// Delete taskB (the blocked task)
		err = env.Store.DeleteTask(env.Ctx, taskB.ID)
		require.NoError(t, err)

		// Reload taskA - blocks relationship still references deleted task
		reloadedA, err = env.Store.LoadTask(env.Ctx, taskA.ID)
		require.NoError(t, err)

		// CURRENT BEHAVIOR: Stale reference remains
		testutil.AssertRelationship(t, reloadedA, model.RelBlocks, taskB.ID)

		t.Logf("⚠️  STALE REFERENCE: Task %d blocks Task %d but target doesn't exist",
			taskA.ID, taskB.ID)
	})

	t.Run("TODO: add cascade delete or orphan prevention", func(t *testing.T) {
		t.Skip("Feature not implemented - orphan prevention/cascade delete needed")

		// Future implementation ideas:
		// 1. Prevent deletion if children exist
		// 2. Cascade delete all descendants
		// 3. Move orphans to a default parent
		// 4. Clear parent relationships on delete
	})
}
