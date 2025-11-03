---
name: tmp-opentask
description: Use when needing to create temporary opentask documentation for agents.
license: MIT
---
# Task Management Guide for Agents

> **NOTE**: This guide uses the opentask CLI. Use `go run cmd/opentask ...` to run commands instead of `./opentask`.

## ⚡ Quick Start

```bash
# START a task (mark in progress FIRST)
go run cmd/opentask --path .tasks task update <id> --status in-progress

# FINISH a task (mark done LAST)
go run cmd/opentask --path .tasks task update <id> --status done

# See what's being worked on
go run cmd/opentask --path .tasks task list --status in-progress

# See what's available
go run cmd/opentask --path .tasks task list --status todo

# Create a new task
go run cmd/opentask --path .tasks task new "Task title" --type story

# List all tasks
go run cmd/opentask --path .tasks task list
```

**Key Rule**: Mark `in-progress` BEFORE work, mark `done` AFTER work. Track progress in todowrite within the session.

## Task Basics

**Use `go run cmd/opentask ...` for all task operations.**

If any needed functionality is missing from the CLI, **ask the user what to do** instead of working around it.

### Common Operations

```bash
# View a specific task
go run cmd/opentask --path .tasks task show <id>

# Create a task with parent epic
go run cmd/opentask --path .tasks task new "Title" --type story --parent <epic_id>

# Create a task with tags
go run cmd/opentask --path .tasks task new "Title" --type story --tag feature --tag urgent

# Filter tasks by status
go run cmd/opentask --path .tasks task list --status todo
go run cmd/opentask --path .tasks task list --status done

# Filter by parent epic
go run cmd/opentask --path .tasks task list --parent <epic_id>

# Filter by type
go run cmd/opentask --path .tasks task list --type story
```

### If CLI lacks needed functionality

**Example scenarios**:
- Need to create relationships between tasks? Ask user how to proceed
- Need to bulk update tasks? Ask user if CLI supports it
- Need to export task data? Ask user what format needed

**Pattern**:
```
1. Try the CLI command
2. If it fails or unsupported: "The opentask CLI doesn't support [feature]. What would you like me to do?"
3. Don't work around missing functionality - get user input
```

## Task Types

- **epic**: Large body of work  
- **plan**: Planning/roadmap items  
- **research**: Investigation tasks  
- **story**: User-facing features (typical work unit)  
- **decision**: Architecture/design decisions  
- **task**: Implementation tasks or chores

## Task Status

Workflow: `todo` → `in-progress` → `done` → `archived`

Use CLI to update:
```bash
go run cmd/opentask --path .tasks task update <id> --status done
```

## Session Workflow

1. **Check current status**:
   ```bash
   go run cmd/opentask --path .tasks task list --status in-progress
   ```

2. **Start a task** (mark BEFORE working):
   ```bash
   go run cmd/opentask --path .tasks task update <id> --status in-progress
   ```

3. **Track progress** with `todowrite` tool (within session)

4. **Finish task** (mark AFTER completing):
   ```bash
   go run cmd/opentask --path .tasks task update <id> --status done
   ```

5. **Session handoff**: If task incomplete, leave as `in-progress` and document next steps in SESSION_SUMMARY.md

## Principles

- **Never delete tasks**: Mark as `archived` instead. Files are audit trail records.
- **CLI updates timestamps**: When you use CLI to update, timestamps are handled automatically.
- **IDs are auto-managed**: When creating tasks, CLI generates next ID.

## Common Tags

Use consistently: `feature`, `bug`, `testing`, `documentation`, `backend`, `frontend`, `chore`, `urgent`

## Get Help

Run:
```bash
go run cmd/opentask --help
go run cmd/opentask task --help
go run cmd/opentask task list --help
```

---

**Applies to**: All agents managing `.tasks/` directory
