---
id: 15
title: Add task linking command
type: story
status: done
tags:
    - feature
relationships:
    - type: parent
      taskID: 12
createdAt: "2025-11-02T08:41:05Z"
updatedAt: "2025-11-02T22:00:00Z"
---

## Objective
Implement a CLI command to create, view, and remove relationships between tasks. This enables users to build task hierarchies, mark dependencies, and associate related work.

## Current State
- Task model supports three relationship types: `parent`, `blocks`, `relates-to`
- Relationships can be created during task creation with `--parent` flag
- No command exists to manage relationships independently
- No way to remove relationships once created

## New Command: `task link`

### Subcommands to Implement

#### 1. `task link add <from-id> <to-id>`
Create a relationship from one task to another.

**Flags:**
- `--type` (required): Relationship type - `parent`, `blocks`, or `relates-to`
  - `parent`: Task A is parent/epic of Task B (hierarchical)
  - `blocks`: Task A blocks/prevents Task B (dependency)
  - `relates-to`: Task A is related to Task B (loose association)

**Behavior:**
- Load both tasks from storage
- Validate relationship doesn't already exist
- Prevent circular relationships (A→B and B→A parent links)
- Add relationship to "from" task's relationships list
- Save updated task
- Output: "Created relationship: task <from-id> [parent/blocks/relates-to] task <to-id>"

**Error Cases:**
- Either task ID doesn't exist → "Task not found: <id>"
- Invalid relationship type → "Invalid relationship type: <type>"
- Circular parent relationship → "Circular parent relationship detected"
- Duplicate relationship → "Relationship already exists"

#### 2. `task link remove <from-id> <to-id>`
Remove a relationship between tasks.

**Flags:**
- `--type` (optional): If provided, only remove relationship of this type. If omitted, remove any relationship.

**Behavior:**
- Load "from" task
- Find and remove matching relationship(s)
- Save updated task
- Output: "Removed relationship: task <from-id> [X] task <to-id>" or "No relationship found"

**Error Cases:**
- From task doesn't exist → "Task not found: <id>"
- Relationship doesn't exist → "Relationship not found"

#### 3. `task link view <task-id>`
Display all relationships for a task (incoming and outgoing).

**Output format:**
```
Task <id>: Title

Outgoing relationships:
  ✓ parent → Task <id>: Title
  → blocks Task <id>: Title
  → relates-to Task <id>: Title
  (none if no outgoing)

Incoming relationships:
  ← parent Task <id>: Title (child of this task)
  ← blocks Task <id>: Title (blocked by)
  ← relates-to Task <id>: Title (related from)
  (none if none found)
```

**Error Cases:**
- Task doesn't exist → "Task not found: <id>"

## Implementation Requirements

### Storage/Model Changes
- No changes needed - relationships already supported in model

### CLI Changes (cmd/task.go)
- Add `taskLinkCmd` command group
- Add `taskLinkAddCmd` subcommand
- Add `taskLinkRemoveCmd` subcommand
- Add `taskLinkViewCmd` subcommand
- Wire into root task command

### Tests to Write
- Unit tests for CLI command handlers
- Integration tests for complete link workflows
- Test circular relationship detection
- Test relationship validation

## Acceptance Criteria
- [ ] `task link add <from-id> <to-id> --type <type>` works correctly
- [ ] `task link remove <from-id> <to-id>` works correctly
- [ ] `task link view <task-id>` shows both incoming and outgoing relationships
- [ ] Circular parent relationships are detected and rejected
- [ ] Duplicate relationships are prevented
- [ ] All error cases produce helpful error messages
- [ ] Help text and examples provided: `task link --help`
- [ ] Unit and integration tests pass
