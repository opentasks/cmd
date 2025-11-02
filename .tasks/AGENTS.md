# Task Management Guide for Agents

This file documents conventions and best practices for managing tasks in the `.tasks/` directory.

## Overview

The `.tasks/` directory uses the OpenTasks format itself to track project work. All files in this directory are task files (`.md` format) with YAML frontmatter metadata.

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

### Using OpenTasks CLI

The easiest way to create new tasks:

```bash
# Create a simple task
./opentasks --path .tasks task new "Task title" --type story

# Create with parent epic
./opentasks --path .tasks task new "Task title" --type story --parent 1

# Create with tags
./opentasks --path .tasks task new "Task title" --type story --tag feature --tag urgent
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
./opentasks --path .tasks task list | tail -1  # Shows max current ID

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
./opentasks --path .tasks task update <id> --status in-progress
```

### Keep Files Organized

- Parent tasks (epics) should have brief titles and summaries
- Child tasks go in epic directories
- Related tasks should have `relates-to` relationships
- Use tags for cross-cutting concerns

### Consistent IDs

IDs are globally sequential integers (1, 2, 3, ...). Use the CLI to get the next ID:

```bash
./opentasks --path .tasks task list | wc -l  # Current count
```

Next ID = current count + 1

## Workflow for Adding New Work

### For Agents/Future Sessions

1. **List current tasks**
   ```bash
   ./opentasks --path .tasks task list
   ./opentasks --path .tasks task list --status todo
   ```

2. **Create new epic** (for major features)
   ```bash
   ./opentasks --path .tasks task new "Epic Name" --type epic
   ```

3. **Create subtasks** (under epic)
   ```bash
   ./opentasks --path .tasks task new "Task name" --type story --parent <epic_id>
   ```

4. **Mark in progress** (when starting work)
   ```bash
   ./opentasks --path .tasks task update <id> --status in-progress
   ```

5. **Update content** (add details, notes)
   ```bash
   # Edit the .md file directly to add content
   vim ".tasks/<path-to-file>.md"
   ```

6. **Mark complete** (when done)
   ```bash
   ./opentasks --path .tasks task update <id> --status done
   ```

7. **Archive old tasks** (after long period)
   ```bash
   ./opentasks --path .tasks task update <id> --status archived
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
./opentasks --path .tasks task new "Add user authentication" --type epic

# 2. Check it was created (note the ID)
./opentasks --path .tasks task list --type epic

# 3. Create research task
./opentasks --path .tasks task new "Research OAuth providers" \
  --type research --parent 1 --tag research

# 4. Create implementation stories
./opentasks --path .tasks task new "Implement OAuth login" \
  --type story --parent 1 --tag feature

./opentasks --path .tasks task new "Add password reset flow" \
  --type story --parent 1 --tag feature

# 5. Create testing task
./opentasks --path .tasks task new "Test authentication flows" \
  --type task --parent 1 --tag testing

# 6. View all tasks for this epic
./opentasks --path .tasks task list --parent 1

# 7. Mark one as in progress
./opentasks --path .tasks task update 2 --status in-progress

# 8. When complete, mark done
./opentasks --path .tasks task update 2 --status done
```

## For Future Reference

- **OpenTasks CLI**: `./opentasks --help`
- **List tasks**: `./opentasks --path .tasks task list`
- **Show task details**: `./opentasks --path .tasks task show <id>`
- **List all commands**: `./opentasks --path .tasks task --help`
- **Mise tasks**: `mise run help` (for development tasks)

## Conventions

✅ **DO**:
- Use CLI to create/update tasks when possible
- Keep descriptions clear and concise
- Update timestamps when modifying files
- Use relationships to link related work
- Tag tasks consistently
- Review status before marking done

❌ **DON'T**:
- Delete or remove task files
- Change task IDs
- Create tasks without proper structure
- Skip the YAML frontmatter
- Use non-standard status values
- Leave orphan tasks without parent epic

---

**Last Updated**: 2025-11-02  
**Applies to**: All agents and humans managing `.tasks/` directory
