package e2e

import (
	"testing"

	"github.com/opentasks/cmd/internal/model"
	"github.com/opentasks/cmd/internal/testutil"
	"github.com/stretchr/testify/assert"
)

// TestE2E_ResearchPlanningImplementationWorkflow tests a complete 15-operation workflow
// representing a realistic feature development cycle:
// Research → Planning → Implementation with status transitions and relationships
func TestE2E_ResearchPlanningImplementationWorkflow(t *testing.T) {
	env := SetupE2EEnvironment(t)
	defer env.Cleanup()

	// 1. Create epic for the feature
	epic := testutil.CreateEpic(t, env.Store,
		testutil.WithTitle("User Authentication System"),
		testutil.WithStatus("planning"),
		testutil.WithTags("feature", "auth"),
	)

	Assert(t, epic, env).
		HasStatus("planning").
		HasType(model.TypeEpic).
		HasTags("feature", "auth")

	t.Logf("✓ Step 1: Created epic #%d", epic.ID)

	// 2. Create research task under epic
	research := testutil.CreateTask(t, env.Store,
		testutil.WithTitle("OAuth vs JWT analysis"),
		testutil.WithType(model.TypeResearch),
		testutil.WithParent(epic.ID),
		testutil.WithStatus("todo"),
	)

	Assert(t, research, env).
		HasType(model.TypeResearch).
		HasParent(epic.ID).
		HasStatus("todo")

	t.Logf("✓ Step 2: Created research task #%d", research.ID)

	// 3. Transition research: todo → in-progress
	research = testutil.TransitionTaskState(t, env.Store, research.ID, "in-progress")
	Assert(t, research, env).HasStatus("in-progress")
	t.Logf("✓ Step 3: Research transitioned to in-progress")

	// 4. Add tags to research task
	research.Tags = []string{"security", "investigation"}
	err := env.Store.SaveTask(env.Ctx, research)
	assert.NoError(t, err)

	// Reload and verify
	research, err = env.Store.LoadTask(env.Ctx, research.ID)
	assert.NoError(t, err)
	Assert(t, research, env).HasTags("security", "investigation")
	t.Logf("✓ Step 4: Added tags to research")

	// 5. Transition research: in-progress → done
	research = testutil.TransitionTaskState(t, env.Store, research.ID, "done")
	Assert(t, research, env).HasStatus("done")
	t.Logf("✓ Step 5: Research completed")

	// 6. Create plan task under epic
	plan := testutil.CreateTask(t, env.Store,
		testutil.WithTitle("Authentication System Design"),
		testutil.WithType(model.TypePlan),
		testutil.WithParent(epic.ID),
		testutil.WithStatus("todo"),
		testutil.WithDescription("Design OAuth 2.0 + JWT implementation"),
	)

	Assert(t, plan, env).
		HasType(model.TypePlan).
		HasParent(epic.ID).
		HasDescription("Design OAuth 2.0 + JWT implementation")

	t.Logf("✓ Step 6: Created plan task #%d", plan.ID)

	// 7. Complete plan (skip in-progress for this test)
	plan = testutil.TransitionTaskState(t, env.Store, plan.ID, "done")
	Assert(t, plan, env).HasStatus("done")
	t.Logf("✓ Step 7: Plan completed")

	// 8. Create story under epic (implementation begins)
	story := testutil.CreateStory(t, env.Store,
		testutil.WithTitle("Implement OAuth 2.0 Flow"),
		testutil.WithParent(epic.ID),
		testutil.WithStatus("todo"),
	)

	Assert(t, story, env).
		HasType(model.TypeStory).
		HasParent(epic.ID).
		HasStatus("todo")

	t.Logf("✓ Step 8: Created story #%d", story.ID)

	// 9. Transition epic: planning → active (implementation started)
	epic = testutil.TransitionTaskState(t, env.Store, epic.ID, "active")
	Assert(t, epic, env).HasStatus("active")
	t.Logf("✓ Step 9: Epic transitioned to active")

	// 10. Transition story: todo → in-progress
	story = testutil.TransitionTaskState(t, env.Store, story.ID, "in-progress")
	Assert(t, story, env).HasStatus("in-progress")
	t.Logf("✓ Step 10: Story in progress")

	// 11. Create blocking task (database setup blocks OAuth)
	blockingTask := testutil.CreateTask(t, env.Store,
		testutil.WithTitle("Setup user database schema"),
		testutil.WithParent(story.ID),
		testutil.WithStatus("todo"),
	)

	Assert(t, blockingTask, env).
		HasTitle("Setup user database schema").
		HasParent(story.ID)

	t.Logf("✓ Step 11: Created blocking task #%d", blockingTask.ID)

	// 12. Create OAuth implementation task that is blocked
	oauthTask := testutil.CreateTask(t, env.Store,
		testutil.WithTitle("Implement OAuth endpoints"),
		testutil.WithParent(story.ID),
		testutil.WithStatus("blocked"),
	)

	// Add blocks relationship from blockingTask to oauthTask
	blockingTask.Relationships = append(blockingTask.Relationships, model.Relationship{
		Type:   model.RelBlocks,
		TaskID: oauthTask.ID,
	})
	err = env.Store.SaveTask(env.Ctx, blockingTask)
	assert.NoError(t, err)

	// Reload and verify relationship
	blockingTask, err = env.Store.LoadTask(env.Ctx, blockingTask.ID)
	assert.NoError(t, err)
	Assert(t, blockingTask, env).HasBlocksRelationship(oauthTask.ID)
	t.Logf("✓ Step 12: Created OAuth task #%d with blocking relationship", oauthTask.ID)

	// 13. Complete blocking task
	blockingTask = testutil.TransitionTaskState(t, env.Store, blockingTask.ID, "done")
	Assert(t, blockingTask, env).HasStatus("done")
	t.Logf("✓ Step 13: Blocking task completed")

	// 14. Unblock and complete OAuth task
	oauthTask = testutil.TransitionTaskState(t, env.Store, oauthTask.ID, "in-progress")
	oauthTask = testutil.TransitionTaskState(t, env.Store, oauthTask.ID, "done")
	Assert(t, oauthTask, env).HasStatus("done")
	t.Logf("✓ Step 14: OAuth task completed")

	// 15. Complete story and epic
	story = testutil.TransitionTaskState(t, env.Store, story.ID, "done")
	Assert(t, story, env).HasStatus("done")

	epic = testutil.TransitionTaskState(t, env.Store, epic.ID, "done")
	Assert(t, epic, env).HasStatus("done")
	t.Logf("✓ Step 15: Story and epic completed")

	// Final verification: All tasks should be done
	allTasks, err := env.Store.ListTasks(env.Ctx)
	assert.NoError(t, err)

	doneCount := 0
	for _, task := range allTasks {
		if task.Status == "done" {
			doneCount++
		}
	}

	// We created: 1 epic + 1 research + 1 plan + 1 story + 2 tasks = 6 tasks
	// All should be done
	assert.GreaterOrEqual(t, doneCount, 6, "At least 6 tasks should be done")

	t.Logf("✓ Workflow complete: %d tasks done", doneCount)
}

// TestE2E_ParallelStoryWorkflow tests multiple stories being worked on in parallel
func TestE2E_ParallelStoryWorkflow(t *testing.T) {
	env := SetupE2EEnvironment(t)
	defer env.Cleanup()

	// Create epic
	epic := testutil.CreateEpic(t, env.Store,
		testutil.WithTitle("User Management System"),
	)

	// Create 3 stories in parallel
	stories := testutil.BulkCreateTasks(t, env.Store, 3, model.TypeStory,
		testutil.WithParent(epic.ID),
	)

	// Verify all stories have correct parent
	for i, story := range stories {
		Assert(t, story, env).HasParent(epic.ID)
		t.Logf("✓ Story %d created under epic", i+1)
	}

	// Transition all to different states (simulating parallel work)
	stories[0] = testutil.TransitionTaskState(t, env.Store, stories[0].ID, "done")
	stories[1] = testutil.TransitionTaskState(t, env.Store, stories[1].ID, "in-progress")
	stories[2] = testutil.TransitionTaskState(t, env.Store, stories[2].ID, "todo")

	// Verify states
	Assert(t, stories[0], env).HasStatus("done")
	Assert(t, stories[1], env).HasStatus("in-progress")
	Assert(t, stories[2], env).HasStatus("todo")

	// Verify epic still has 3 children
	Assert(t, epic, env).HasChildCount(3)

	t.Logf("✓ Parallel workflow: 3 stories in different states")
}
