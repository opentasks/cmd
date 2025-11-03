---
id: 25
title: Implement status transitions command
type: story
status: done
tags:
    - feature
    - workflow
relationships:
    - type: parent
      taskID: 12
createdAt: "2025-11-02T21:00:00Z"
updatedAt: "2025-11-02T12:08:08Z"
---

## Objective
Implement `status transitions` command to help users understand and validate workflow transitions.

## Current State
- Config system supports workflow with transitions (from → to mapping)
- No command exists to view or validate transitions
- Users have no visibility into what status changes are allowed

## New Command: `status transitions`

### Subcommands

#### 1. `status transitions list` (or `status transitions show`)
Display all configured status transitions for the project.

**Output:**
```
Project Workflow Transitions:

todo
  → in-progress
  → archived

in-progress
  → done
  → in-progress (self)
  → todo

done
  → archived

archived
  (no transitions)

Initial status: todo
```

**Format support:** `--format text|json|yaml`

#### 2. `status transitions validate <from> <to>`
Check if a specific transition is allowed.

**Examples:**
```bash
status transitions validate todo in-progress
# Output: ✓ Valid transition

status transitions validate done todo
# Output: ✗ Invalid transition
# Did you mean: archived
```

**Exit codes:**
- 0: transition is valid
- 1: transition is invalid
- 2: status not found

#### 3. `status transitions describe <status>`
Show what a specific status means and what transitions are available.

**Output:**
```
Status: in-progress

Description: Work currently being done

From these statuses:
  ← todo
  ← in-progress (self)

To these statuses:
  → done
  → in-progress (self)
  → todo

Total transitions involving this status: 5
```

## Implementation Details

### Data Structure
Use config's workflow transitions to build transition graph.

### Validation Logic
- Check if (from, to) pair exists in workflow config
- Detect and allow self-transitions if configured
- Suggest similar statuses on invalid input

## Testing Requirements
- Unit tests for transition validation
- Tests with various workflow configurations
- Tests with no transitions defined
- Tests with self-transitions
- E2E test that shows workflow accurately

## Acceptance Criteria
- [ ] `status transitions list` shows all transitions
- [ ] `status transitions validate` checks specific transitions
- [ ] `status transitions describe` shows status details
- [ ] All commands support `--format` flag
- [ ] Helpful error messages for invalid status
- [ ] Suggestions for similar status names
- [ ] Help text provided: `status transitions --help`
- [ ] Exit codes correct for scripting
- [ ] Tests verify all workflows display correctly
