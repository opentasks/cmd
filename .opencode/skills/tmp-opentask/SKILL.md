---
name: tmp-opentask
description: Use when needing to create temporary opentask documentation for agents.
license: MIT
---
# Task Management Guide for Agents

> **NOTE**: This guide describes manual task management using the opentask CLI. This is temporary until the opentask tool becomes feature-complete and fully usable through the MCP interface. Once complete, tasks should be managed through the proper tool integration rather than direct CLI commands.

This file documents conventions and best practices for managing tasks in the `.tasks/` directory.

## ⚡ Quick Start (TL;DR)

```bash
# START a task (mark in progress FIRST)
./opentask --path .tasks task update <id> --status in-progress

# FINISH a task (mark done LAST)
./opentask --path .tasks task update <id> --status done

# See what's being worked on
./opentask --path .tasks task list --status in-progress

# See what's available
./opentask --path .tasks task list --status todo
```

**Key Rule**: Mark `in-progress` BEFORE work, mark `done` AFTER work. This is your project's work tracking system.

## Overview

The `.tasks/` directory uses the opentask format itself to track project work. All files in this directory are task files (`.md` format) with YAML frontmatter metadata.

**Core Principle**: Tasks are immutable records. Files are **never removed** - only marked as complete or archived.

## Task File Structure

Every task file has this structure:

```markdown
---
id: <number>
title: <task title>
type: <epic|plan|research|story|decision|task>
status: <todo|in-progress|reviewing|done|archived>
tags: [tag1, tag2, ...]
relationships:
  - type: parent
    taskID: <id>
  - type: blocks
    taskID: <id>
createdAt: <RFC3339 timestamp>
updatedAt: <RFC3339 timestamp>
---

# Task Title

Markdown content describing the task, requirements, notes, etc.
```

## Creating New Tasks

### Using opentask CLI

The easiest way to create new tasks:

```bash
# Create a simple task
./opentask --path .tasks task new "Task title" --type story

# Create with parent epic
./opentask --path .tasks task new "Task title" --type story --parent 1

# Create with tags
./opentask --path .tasks task new "Task title" --type story --tag feature --tag urgent
```

### Using Templates

Templates are stored in `.tasks/templates/` (not yet fully implemented). When creating tasks manually:

1. **Copy appropriate template** from `.tasks/templates/`
2. **Update frontmatter** (id, title, type, status, tags, relationships)
3. **Add content** in markdown section
4. **Save to correct location**:
   - Epics: `.tasks/e-<id>-<slug>.md`
   - Child tasks: `.tasks/<parent_id>-<parent_slug>/<typecode>-<id>-<slug>.md`

### Manual Creation

If creating files manually:

```bash
# Get next ID
./opentask --path .tasks task list | tail -1  # Shows max current ID

# Create file with proper naming
touch ".tasks/1-my-epic/s-<next_id>-my-story.md"

# Edit and add frontmatter and content
vim ".tasks/1-my-epic/s-<next_id>-my-story.md"
```

## Naming Conventions

### Directory Structure

```
.tasks/
├── design/                    # Design phase (epic 0)
├── templates/                 # Task templates
├── <epic_id>-<slug>/         # Child tasks of epic
│   ├── <type>-<id>-<slug>.md
│   └── <type>-<id>-<slug>.md
└── <type>-<id>-<slug>.md     # Epic or root-level tasks
```

### Slug Generation

Slugs are the human-readable part of filenames (e.g., `write-unit-tests`).

**Rules**:
1. Convert to lowercase
2. Replace spaces and special characters with hyphens
3. Keep 3-5 key words
4. Maximum ~50 characters
5. Remove articles (a, the) and conjunctions (and, or)

**Examples**:
- "Write Unit Tests" → `write-unit-tests`
- "Implement Task Linking" → `implement-task-linking`
- "Design Database Schema" → `design-database-schema`

## Task Types and When to Use

| Type | Purpose | Example |
|------|---------|---------|
| **epic** | Large body of work (100+ effort points) | "Phase 2: Testing & Polish" |
| **plan** | Planning and roadmap items | "Plan authentication flow" |
| **research** | Research and exploration tasks | "Research OAuth providers" |
| **story** | User-facing feature (typical dev unit) | "Implement user registration" |
| **decision** | Architecture/design decision | "Decide on database backend" |
| **task** | Implementation task or chore | "Add error logging" |

## Task Status Workflow

Default statuses (can be customized in `config.toml`):

| Status | Meaning |
|--------|---------|
| **todo** | Not started (initial status) |
| **in-progress** | Currently being worked on |
| **reviewing** | Complete, pending review/approval |
| **done** | Finished and verified |
| **archived** | Historical record (completed long ago) |

Transition from status:
- `todo` → `in-progress`, `archived`
- `in-progress` → `reviewing`, `todo`, `archived`
- `reviewing` → `done`, `in-progress`, `archived`
- `done` → `archived`

## Relationships

Tasks can link to other tasks. Store relationships in frontmatter:

```yaml
relationships:
  - type: parent
    taskID: 1
  - type: blocks
    taskID: 5
  - type: relates-to
    taskID: 3
```

### Relationship Types

| Type | Meaning |
|------|---------|
| **parent** | Parent epic (hierarchical) |
| **blocks** | This task blocks another |
| **relates-to** | Related but independent task |

## Task Management Best Practices

### Integration with Session Work

When working on a task within a single session:

1. **Create a todowrite todo for the task**
   ```bash
   # At start of work
   todowrite([{
     id: "my-task-id",
     content: "Description of work",
     status: "in_progress",
     priority: "high"
   }])
   ```

2. **Update todowrite todo status as work progresses**
   ```bash
   # Mid-work: break into subtasks
   todowrite([{
     id: "my-task-1",
     content: "Part 1 of task",
     status: "in_progress",
     priority: "high"
   }, {
     id: "my-task-2", 
     content: "Part 2 of task",
     status: "pending",
     priority: "high"
   }])
   
   # After part 1 done
   todowrite([{
     id: "my-task-1",
     content: "Part 1 of task",
     status: "completed",
     priority: "high"
   }, {
     id: "my-task-2",
     content: "Part 2 of task",
     status: "in_progress",
     priority: "high"
   }])
   ```

3. **When session ends OR task completes** → **Update opentask immediately**
   ```bash
   # Task complete: mark done in opentask
   ./opentask --path .tasks task update <id> --status done
   
   # Task incomplete: leave as in-progress in opentask
   # (future sessions will see it's active)
   ```

### Session Handoff Pattern

When a task is in-progress but will be resumed in a future session:

1. **Leave task as `in-progress` in opentask** - Don't mark done
2. **Add session notes to task file** describing next steps
3. **Create SESSION_SUMMARY.md** in project root noting which task to resume
4. **Next session**: Read SESSION_SUMMARY.md and resume the in-progress task

Example SESSION_SUMMARY.md:
```markdown
# Current Work Status

## In Progress
- Task 42: Implementing feature X
  - What's done: Core logic implemented
  - Next steps: Add tests and error handling
  - Time estimate: 2-3 hours

## Ready to Start
- Task 43: Polish error messages
- Task 44: Add documentation
```

## Important Principles

### Never Delete Files

Files are never removed from `.tasks/`. Instead:
- Mark as `status: archived` when complete
- Keep all historical records
- This creates an audit trail

### Update Timestamps

When modifying a task, update the `updatedAt` field:

```yaml
updatedAt: 2025-11-02T12:34:56Z  # RFC3339 format
```

Or use CLI to update automatically:
```bash
./opentask --path .tasks task update <id> --status in-progress
```

### Keep Files Organized

- Parent tasks (epics) should have brief titles and summaries
- Child tasks go in epic directories
- Related tasks should have `relates-to` relationships
- Use tags for cross-cutting concerns

### Consistent IDs

IDs are globally sequential integers (1, 2, 3, ...). Use the CLI to get the next ID:

```bash
./opentask --path .tasks task list | wc -l  # Current count
```

Next ID = current count + 1

## Workflow for Task Management

### ⚡ CRITICAL: Task Status Tracking Pattern

**This is the most important workflow - follow it strictly for every task:**

1. **When you START working on a task** → **IMMEDIATELY mark it `in-progress`**
   ```bash
   ./opentask --path .tasks task update <id> --status in-progress
   ```
   - Do this BEFORE writing any code or making changes
   - Update `updatedAt` timestamp (CLI does this automatically)
   - Prevents duplicate work and shows current effort

2. **While working** → **Track progress via todo list or session notes**
   - Use `todowrite` tool to manage subtasks
   - Update the task file with progress notes if it's a multi-day task
   - Mark the todo as `in_progress` ONLY while actively working

3. **When you FINISH the task** → **IMMEDIATELY mark it `done`**
   ```bash
   ./opentask --path .tasks task update <id> --status done
   ```
   - Do this BEFORE moving to next task
   - Verify all acceptance criteria met
   - CLI automatically updates `updatedAt` timestamp

4. **Moving between tasks**
   ```bash
   # ALWAYS do this sequence:
   # 1. Mark current task done (if finished)
   ./opentask --path .tasks task update <current_id> --status done
   
   # 2. Mark next task in progress (if starting)
   ./opentask --path .tasks task update <next_id> --status in-progress
   ```
   - Never skip this - it's the project audit trail
   - Helps future sessions understand where work stopped/started

### For Agents/Future Sessions - Complete Workflow

1. **List current tasks and status**
   ```bash
   ./opentask --path .tasks task list
   ./opentask --path .tasks task list --status in-progress  # Shows active work
   ./opentask --path .tasks task list --status todo          # Shows available work
   ./opentask --path .tasks task list --status done          # Shows completed work
   ```

2. **Create new epic** (for major features)
   ```bash
   ./opentask --path .tasks task new "Epic Name" --type epic
   ```

3. **Create subtasks** (under epic)
   ```bash
   ./opentask --path .tasks task new "Task name" --type story --parent <epic_id>
   ```

4. **START TASK** - Mark in progress BEFORE doing work
   ```bash
   ./opentask --path .tasks task update <id> --status in-progress
   # ← DO THIS FIRST, before writing any code
   ```

5. **Do the work**
   ```bash
   # Edit files, write code, implement feature
   # Track progress with todowrite tool
   # Update task file if documenting progress
   ```

6. **FINISH TASK** - Mark done AFTER completing work
   ```bash
   ./opentask --path .tasks task update <id> --status done
   # ← DO THIS IMMEDIATELY after finishing
   ```

7. **Archive old tasks** (after long period)
   ```bash
   ./opentask --path .tasks task update <id> --status archived
   ```

## File Locations and Purposes

| Location | Purpose |
|----------|---------|
| `.tasks/design/` | Design phase work items (completed, reference only) |
| `.tasks/templates/` | Task templates for creating new tasks |
| `.tasks/1-phase-2-testing-polish/` | Phase 2 development tasks |
| `.tasks/<epic_id>-<slug>/` | Any new epic's subtasks |

## Tags for Organization

Use tags to categorize work across epics:

| Tag | Meaning |
|-----|---------|
| **testing** | Test-related work |
| **feature** | New feature |
| **bug** | Bug fix |
| **polish** | Code quality, UX improvements |
| **documentation** | Docs, comments, guides |
| **backend** | Backend/server work |
| **frontend** | UI/frontend work |
| **chore** | Maintenance, refactoring |
| **urgent** | High priority |

Use tags consistently across tasks.

## Example: Adding a New Feature

```bash
# 1. Create epic for feature
./opentask --path .tasks task new "Add user authentication" --type epic

# 2. Check it was created (note the ID)
./opentask --path .tasks task list --type epic

# 3. Create research task
./opentask --path .tasks task new "Research OAuth providers" \
  --type research --parent 1 --tag research

# 4. Create implementation stories
./opentask --path .tasks task new "Implement OAuth login" \
  --type story --parent 1 --tag feature

./opentask --path .tasks task new "Add password reset flow" \
  --type story --parent 1 --tag feature

# 5. Create testing task
./opentask --path .tasks task new "Test authentication flows" \
  --type task --parent 1 --tag testing

# 6. View all tasks for this epic
./opentask --path .tasks task list --parent 1

# 7. Mark one as in progress
./opentask --path .tasks task update 2 --status in-progress

# 8. When complete, mark done
./opentask --path .tasks task update 2 --status done
```

## For Future Reference

- **opentask CLI**: `./opentask --help`
- **List tasks**: `./opentask --path .tasks task list`
- **Show task details**: `./opentask --path .tasks task show <id>`
- **List all commands**: `./opentask --path .tasks task --help`
- **Mise tasks**: `mise run help` (for development tasks)

## Conventions

✅ **DO**:
- **Mark `in-progress` BEFORE starting work on a task** ← CRITICAL
- **Mark `done` IMMEDIATELY after finishing a task** ← CRITICAL
- Use CLI to create/update tasks when possible
- Keep descriptions clear and concise
- Update timestamps when modifying files
- Use relationships to link related work
- Tag tasks consistently
- Review status before marking done
- Move your todo items through statuses as you work
- Complete all subtasks before moving to next epic task

❌ **DON'T**:
- Delete or remove task files
- Change task IDs
- Create tasks without proper structure
- Skip the YAML frontmatter
- Use non-standard status values
- Leave orphan tasks without parent epic
- **Leave a task marked `in-progress` at end of session without notes** ← CRITICAL
- **Forget to mark task `done` before moving to next task** ← CRITICAL
- Have multiple tasks marked `in-progress` without clear reason

## Status Transition Rules (STRICT)

**Always follow this pattern when switching tasks:**

```
Current Task          Next Task
    ↓                   ↓
in-progress    →    done (mark current)
    ↓                   ↓
(cleanup)      →    in-progress (mark next)
```

Never break the chain. Every task should have clear start and end points in the work log.

---

**Last Updated**: 2025-11-02  
**Applies to**: All agents and humans managing `.tasks/` directory

### Pattern Override: Session-Specific Work
If using `todowrite` for session-specific work items:
- `todowrite` tracks sub-tasks within a single task
- `opentask task update` tracks the overall task status
- Both should be kept in sync
- Mark opentask `done` when ALL subtasks are complete
- Mark opentask `in-progress` when starting ANY subtask in that task
