---
id: 48
title: Update default workflow to include backlog status
type: task
status: todo
tags: [workflow, task-management, ux]
relationships: []
createdAt: "2025-11-03T09:45:00Z"
updatedAt: "2025-11-03T09:45:00Z"
---

# Task: Update Default Workflow to Include Backlog Status

## Problem

The current default workflow is:
```
todo → in-progress → reviewing → done → archived
```

But when creating epics, stories, or other higher-level tasks without a description, they should start in `backlog` status rather than `todo`. This reflects the real-world workflow where:

- **backlog**: Items that are identified but not yet ready for work (missing details, awaiting prioritization)
- **todo**: Items that are ready to be worked on (have clear scope and requirements)
- **in-progress**: Currently being worked on
- **reviewing**: Code/work review in progress
- **done**: Completed work

## Requirements

### Update Default Workflow

Change default workflow from:
```
todo → in-progress → reviewing → done → archived
```

To:
```
backlog → todo → in-progress → reviewing → done → archived
```

With proper transitions:
- `backlog` → `todo` (move to active work)
- `backlog` → `archived` (discard undeveloped ideas)
- `todo` → `in-progress` (start work)
- `todo` → `archived` (cancel planned work)
- `in-progress` → `reviewing` (ready for review)
- `in-progress` → `todo` (pause work, return to backlog)
- `in-progress` → `archived` (abandon work)
- `reviewing` → `done` (approved)
- `reviewing` → `in-progress` (request changes, return to work)
- `reviewing` → `archived` (reject)
- `done` → `archived` (archive completed)

### Update Default Initial Status

Change default initial status from `todo` to `backlog` ONLY for:
- Epics
- Plans
- Research tasks
- Stories
- Decisions
- Other high-level tasks without description

Regular tasks without description should still default to `todo`.

### Logic for Status Assignment

When creating a task:
```
if task_type in ["epic", "plan", "research", "story", "decision"] and no_description:
  initial_status = "backlog"
else:
  initial_status = "todo"
```

## Implementation Plan

### Phase 1: Update Config Defaults

1. **Update `internal/config/config.go`** - `DefaultWorkflow()`:
   ```go
   return WorkflowConfig{
       Statuses: []string{"backlog", "todo", "in-progress", "reviewing", "done", "archived"},
       Initial:  "backlog",  // Changed from "todo"
       Transitions: []TransitionConfig{
           {From: "backlog", To: []string{"todo", "archived"}},
           {From: "todo", To: []string{"in-progress", "archived"}},
           {From: "in-progress", To: []string{"reviewing", "todo", "archived"}},
           {From: "reviewing", To: []string{"done", "in-progress", "archived"}},
           {From: "done", To: []string{"archived"}},
       },
   }
   ```

2. **Update test files** to reflect new workflow

### Phase 2: Update Task Creation Logic

1. **Update `cmd/task.go`** - `newCmd`:
   - Add logic to detect task type
   - Check if description is provided
   - Set `initial_status` based on type and description
   - Ensure proper status is set before saving

2. **Code logic**:
   ```go
   status := "todo"  // default
   
   // High-level tasks without description start in backlog
   highLevelTypes := []string{"epic", "plan", "research", "story", "decision"}
   if contains(highLevelTypes, taskType) && description == "" {
       status = "backlog"
   }
   
   task.Status = status
   ```

### Phase 3: Update Templates

1. **Update default templates** in `./_templates/`:
   - Epic template should mention backlog status
   - Plan template should mention backlog status
   - Research template should mention backlog status
   - Story template should mention backlog status
   - Decision template should mention backlog status

2. **Template note**:
   ```markdown
   > Note: Tasks without a description start in "backlog" status.
   > Complete the description and move to "todo" when ready for work.
   ```

### Phase 4: Update Documentation

1. **Update `docs/Config.md`**:
   - Update default workflow example
   - Explain backlog status and when to use it
   - Show transition diagram

2. **Add workflow diagram** showing:
   ```
   backlog ──→ todo ──→ in-progress ──→ reviewing ──→ done
      ↓          ↓            ↓              ↓          ↓
    archived   archived    archived       archived   archived
   ```

3. **Update README** if it mentions workflow

### Phase 5: Update Tests

1. **Update existing tests** that check default workflow:
   - `internal/config/config_test.go` - Default workflow tests
   - Update expected statuses and transitions

2. **Add new tests** for task creation:
   - Test epic without description starts in backlog
   - Test epic with description starts in todo (or backlog?)
   - Test regular task without description starts in todo
   - Test explicit status overrides automatic assignment

3. **Test transition validation**:
   - Verify all new transitions are valid
   - Verify old transitions still work or are updated
   - Test attempting invalid transitions

### Phase 6: Handle Migration

1. **For existing projects**:
   - If user has custom workflow, leave it unchanged
   - If user has default workflow, update to new default
   - Add migration note in docs

2. **Migration path**:
   - Existing tasks keep their status
   - New workflow applies only to newly created tasks
   - Users can update existing tasks with `task update` if needed

## Files to Modify

### Core Changes
- `internal/config/config.go` - `DefaultWorkflow()` function
- `cmd/task.go` - `newCmd` function (status assignment logic)

### Tests
- `internal/config/config_test.go` - Update workflow tests
- `cmd/task_test.go` - Add status assignment tests (if exists)

### Documentation
- `docs/Config.md` - Update workflow section
- `README.md` - Update workflow reference
- `./_templates/epic.md` - Add backlog note
- `./_templates/plan.md` - Add backlog note
- `./_templates/research.md` - Add backlog note
- `./_templates/story.md` - Add backlog note
- `./_templates/decision.md` - Add backlog note

## Acceptance Criteria

- [ ] Default workflow includes backlog status
- [ ] Default workflow has correct transitions
- [ ] Epic/plan/story/decision without description starts in backlog
- [ ] Regular task without description starts in todo
- [ ] All status transitions work correctly
- [ ] Existing custom workflows unchanged
- [ ] Templates mention backlog status
- [ ] Documentation updated with new workflow
- [ ] All tests pass
- [ ] No migration issues

## Testing Checklist

```bash
# Test 1: Create epic without description
opentask task new "My Epic Idea" --type epic
opentask task show 1  # Should show status: backlog

# Test 2: Create epic with description
opentask task new "My Epic" --type epic --description "Full epic description"
opentask task show 2  # Status: backlog or todo?

# Test 3: Create regular task without description
opentask task new "Quick task" --type task
opentask task show 3  # Should show status: todo

# Test 4: Create task with explicit status
opentask task new "Task" --type task --status in-progress
opentask task show 4  # Should show status: in-progress

# Test 5: Verify transitions
opentask task update 1 --status todo  # backlog → todo (should work)
opentask task update 1 --status reviewing  # todo → reviewing (should fail)

# Test 6: Check workflow
opentask config view  # Should show backlog in statuses
```

## Backlog vs Todo Distinction

**Backlog (Not Ready):**
- Idea hasn't been fully fleshed out
- Missing important details or requirements
- Awaiting prioritization decision
- No clear scope or acceptance criteria

**Todo (Ready):**
- Has clear description and requirements
- Ready to be picked up for work
- Acceptance criteria defined
- All blocking information available

## Related Tasks

- Task 44: Task editing (need status update command)
- Task 47: Task deletion (ensure works with backlog)
- Story 46: Config schema (workflow may be customizable per project)

## Notes

- Some users may prefer different workflow - config allows overriding
- This is a sensible default but not mandatory
- Consider adding `task move-to-todo` convenience command later
- May want different default initial status per task type (optional enhancement)
