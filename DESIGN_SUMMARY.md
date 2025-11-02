# OpenTasks Design Summary

Quick reference for understanding the system architecture. Read this first, then consult `.tasks/design/` for detailed specifications.

## What is OpenTasks?

A markdown-based task management system written in Go. Tasks are stored as `.md` files with YAML frontmatter metadata. The system supports customizable workflows, task relationships, and pluggable storage backends.

**Key insight**: Tasks are first-class citizens in version control. They live alongside code.

## Core Concepts

### Tasks (`.md` files)

Every task is a markdown file with YAML frontmatter:

```markdown
---
id: s-42              # Semantic ID: type-prefix + sequence number
title: My task
type: story           # epic|plan|research|story|decision|task
status: in-progress   # Customizable per project
tags: [feature, core]
relationships:
  - type: parent      # "parent", "blocks", "relates-to"
    taskID: e-5       # Links to other tasks
createdAt: 2025-11-02T10:00:00Z
updatedAt: 2025-11-02T10:30:00Z
---

# My task

Markdown content here.
```

### Semantic IDs

- **Format**: `{type-prefix}-{sequence}{collision-suffix}`
- **Type prefixes**: e(pic), p(lan), r(esearch), s(tory), d(ecision), t(ask)
- **Sequence**: Per-type counter (first story is s-1, second is s-2)
- **Collision suffix**: Letter appended if ID already exists (s-42a, s-42b)
- **Generation**: Storage backend counts matching files, no persistent state

### Relationships

Tasks can link to other tasks:

```yaml
relationships:
  - type: parent      # Hierarchical (epic contains stories)
    taskID: e-5
  - type: blocks      # This task blocks another
    taskID: s-10
  - type: relates-to  # Related but independent
    taskID: p-3
```

### Project Structure

```
project-root/
├── .tasks/                   # Default local task location
│   ├── config.toml          # Optional (all defaults work without it)
│   ├── 1.epic.md            # File naming: {seq}.{type}.md
│   ├── 1.plan.md
│   ├── 1.story.md
│   └── templates/           # Optional local templates
│       └── story.md
└── ... (project files)
```

Alternative locations:
- `${XDG_DATA_HOME}/opentasks/projects/{derived-id}/` (user-level, git-aware)
- Any path via `--path` flag or `OPENTASKS_PROJECT_PATH` env var

### Configuration (config.toml)

Optional. Every section/field has sensible defaults.

```toml
[project]
name = "My Project"
description = "Optional"
owner = "Optional"

[workflow]
statuses = ["todo", "in-progress", "reviewing", "done", "archived"]
initial = "todo"
[[workflow.transitions]]
from = "todo"
to = ["in-progress", "archived"]
# ... more transitions

[templates]
story = "./templates/story.md"  # Optional paths

[storage]
backend = "markdown-fs"         # Default backend
path = "."                      # Project root
```

## Architecture

### BaseStorage Interface

The abstraction for task persistence. All backends must implement:

```go
type BaseStorage interface {
    LoadTask(id string) (*Task, error)
    SaveTask(task *Task) error
    DeleteTask(id string) error
    ListTasks(filters ...TaskFilter) ([]*Task, error)  // Functional options
    NextID(taskType string) (string, error)             // Returns next ID for type
    GetRelated(taskID, relationType string) ([]*Task, error)
    Close() error
}
```

### MarkdownFileStorage (Default)

Implementation using filesystem:
- Parses `*.md` files in project directory
- ID generation via file globbing (no persistent counter)
- Handles collision detection automatically
- YAML frontmatter parsing/writing

### QueryEngine

Simple wrapper around storage with convenience methods:

```go
engine := NewQueryEngine(storage)
stories := engine.ListTasks(
    WithStatus("in-progress"),
    WithType("story"),
    WithParent("e-5"),
)
```

Filter helpers: `WithStatus()`, `WithType()`, `WithTag()`, `WithParent()`, `WithRelationship()`

### CLI (Viper/Cobra)

Commands follow this structure:

```
opentasks
├── task
│   ├── new <title>        # Create new task
│   ├── list               # List tasks with filters
│   ├── show <id>          # Show task details
│   ├── update <id>        # Update task
│   ├── delete <id>        # Delete task
│   └── link <id>          # Link tasks
├── project
│   ├── new                # Create project
│   └── list               # List projects
└── config
    ├── show               # Show config
    └── set <key> <value>  # Update config
```

## Data Model Types

See `.tasks/design/1.research.md` for complete Go code.

**Key structs**:
- `Task`: ID, Title, Type, Status, Tags, Relationships, Description, Timestamps
- `Relationship`: Type (string), TaskID (string)
- `ProjectConfig`: Workflow, Templates, Storage, Project metadata

## File Organization

### In Git Repo

```
opentasks/
├── README.md
├── DESIGN_SUMMARY.md       # This file
├── go.mod, go.sum
├── main.go
├── cmd/                    # CLI commands
├── internal/
│   ├── model/
│   ├── storage/
│   ├── config/
│   └── query/
├── templates/              # Built-in task templates
└── .tasks/design/          # Our own tasks (dog-food)
```

### In Projects

```
any-project/
├── .tasks/                 # Default location
│   ├── config.toml
│   ├── 1.epic.md
│   ├── 1.story.md
│   └── templates/          # Optional
└── ... project files
```

## Implementation Roadmap

**Phase 1 (Core)**:
- [ ] Data models (Task, Relationship)
- [ ] MarkdownFileStorage
- [ ] BaseStorage interface
- [ ] QueryEngine with filters
- [ ] Config loading (TOML)
- [ ] CLI commands (create, list, show, update, delete)

**Phase 2 (Refinement)**:
- [ ] Task linking/relationships CLI
- [ ] Template system
- [ ] Better error messages
- [ ] Config validation
- [ ] Tests

**Phase 3 (Integration)**:
- [ ] MCP (Multi-Client Proxy) for AI agents
- [ ] Project discovery/management
- [ ] Advanced queries (full-text, dependency graphs)

**Phase 4+ (Futures)**:
- [ ] SQLite backend
- [ ] Web UI (secondary to CLI)
- [ ] Sync/collaboration
- [ ] Cloud storage backends

## Important Design Decisions

| Decision | Choice | Why |
|----------|--------|-----|
| ID system | Semantic sequential per-type | Human-readable, collision-resistant, no state |
| Storage | Pluggable interface | Different backends for different needs |
| Config | Optional with defaults | Works without config, composable |
| Relationships | Slice of structs | Single source of truth, flexible types |
| Query | Functional options | Composable, simple, no query language |
| Frontmatter | YAML | Human-readable, standard, tool-friendly |

## For New Sessions

1. **Start here**: This file (DESIGN_SUMMARY.md)
2. **Data models**: `.tasks/design/1.research.md`
3. **Storage details**: `.tasks/design/2.research.md`
4. **File structure**: `.tasks/design/3.research.md`
5. **Design rationale**: `.tasks/design/*story.md`

All tasks are tracked in `.tasks/design/` using the OpenTasks format itself (dog-fooding).

## Quick Reference

**Create task file manually**:
```markdown
---
id: s-1
title: First story
type: story
status: todo
tags: [example]
relationships: []
createdAt: 2025-11-02T10:00:00Z
updatedAt: 2025-11-02T10:00:00Z
---

# First story

This is a test task.
```

**ID format**:
- Valid: `e-1`, `s-42`, `p-3a`, `d-1b`
- Invalid: `story-1`, `s_1`, `s-01` (no leading zero)

**Common operations** (future CLI):
```bash
opentasks task new "Build API" --type story --parent e-5
opentasks task ls --status in-progress --type story
opentasks task show s-42
opentasks task link s-42 --blocks s-10
```
