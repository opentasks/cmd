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
