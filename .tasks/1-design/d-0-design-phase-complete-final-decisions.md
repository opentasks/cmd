---
id: 0
title: Design Phase Complete - Final Decisions
type: decision
status: done
tags: [decision, design, complete, architecture]
relationships: []
createdAt: 2025-11-02T08:00:00Z
updatedAt: 2025-11-02T19:35:00Z
---

# OpenTasks Design Decisions

Final decisions on architecture and implementation strategy. All rationale documented in `.tasks/design/` story files.

## ID System

**Format**: Global sequential integers (1, 2, 3, ..., 42, ...)
- Type stored in frontmatter (`type: story`), not in ID
- Simple generation: find max ID, return max + 1
- IDs never reused, even if tasks deleted

**Why**: Simpler than per-type counters, decoupled from storage format, prevents ID collision after deletion

See: `.tasks/design/2.story.md` (Story 6: Semantic ID System)

## File Organization

**Pattern**: `<epic_id>-<epic_slug>/<typecode>-<id>-<slug>.md`

**Root-level constraint**: Only epics at project root
- Root files: `e-1-design.md`, `e-2-implement.md`
- Child tasks: `1-design/s-5-task-name.md`, `1-design/p-3-plan-name.md`
- All non-epic tasks must have parent epic

**Why**: Clear hierarchy, consistent organization, groups related work

### Epic Directory Naming

Pattern: `{epic_id}-{epic_slug}/`
- Include both ID and slug for discoverability
- Example: `1-design-opentasks/` instead of just `1/`
- Slug is informational (stale slugs acceptable if title changes)

**Why**: Self-documenting directories, browsable without opening files

## Slug Generation

**Length**: 3-5 key words, ~50-60 characters max

**Strategy**:
1. Convert to lowercase
2. Remove articles (a, an, the) and conjunctions (and, or, with, for)
3. Keep first N words until approaching 60 chars
4. Join with hyphens
5. Remove special characters

**Examples**:
- "Task Data Model and Relationships" → `task-data-model`
- "Semantic ID System with Collision Detection" → `semantic-id-system`
- "BaseStorage Interface and Responsibilities" → `base-storage-interface`
- "Design OpenTasks Core System" → `design-opentasks-core`

**Why**: Balance between clarity and path length, human-readable filenames

## NextID Implementation

**Algorithm**: Find maximum ID from all task frontmatter, return `max + 1`

```pseudocode
Walk all .md files in project
  Parse frontmatter, extract ID
  Track maximum ID seen
Return maximum ID + 1
```

**Characteristics**:
- IDs never reused (even if tasks deleted)
- Robust to manual file creation
- Self-healing (works even if filenames wrong)
- Slightly slower than counting (parses all files), but acceptable for MVP

**Why**: Prevents accidental ID reuse, robust to human mistakes, works with any file organization

## Relationship Hierarchy

**Parent relationship**: Every non-epic task must have a parent epic
- Stored in frontmatter as `relationships: [{type: "parent", taskID: 1}]`
- Determines directory location (tasks live in parent epic's directory)
- Enforced by storage backend

**Other relationships**: `blocks`, `relates-to`
- Can link any tasks to any tasks
- Optional (not required)
- Don't affect file organization

**Why**: Clear hierarchy, organized grouping, single parent-child relationship per task

## Configuration

**Location**: Optional `config.toml` in project root

**What it can override**:
- Workflow statuses and transitions
- Template locations
- Storage backend settings
- Project metadata

**Defaults**: Sensible built-ins (5-status workflow, embedded templates, markdown-fs)

**Why**: Works without config, flexible for custom workflows, composable at multiple levels

## Storage Backend

**Interface**: `BaseStorage` (pluggable)

**Default**: `MarkdownFileStorage`
- Stores tasks as `.md` files
- Organizes hierarchically by epic
- Parses YAML frontmatter
- Implements all required methods

**Future backends**: SQLite, PostgreSQL, Cloud storage, etc.

**Why**: Flexible, can adapt storage without changing core logic

## CLI Architecture

**Framework**: Cobra (commands) + Viper (config)

**Commands**:
- `opentasks task new` - Create task
- `opentasks task ls` - List tasks with filters
- `opentasks task show {id}` - Show task details
- `opentasks task update {id}` - Update task
- `opentasks task delete {id}` - Delete task
- `opentasks project` - Project management

**Why**: Standard Go CLI framework, handles config hierarchy automatically

## Query Engine

**Pattern**: Functional options for composable filtering

```go
storage.ListTasks(
    WithStatus("in-progress"),
    WithType("story"),
    WithParent(5),
    WithTag("feature"),
)
```

**Convenience methods**:
- `FindByID(id int)`
- `FindChildren(parentID int)`
- `FindBlocking(taskID int)`
- `FindBlockedBy(taskID int)`
- `FindRelated(taskID int)`

**Why**: Composable, extensible, simple to use, can add filters without API changes

## Data Model

**Task struct**:
```go
type Task struct {
    ID            int            // Global sequential ID
    Title         string
    Type          string         // epic|plan|research|story|decision|task
    Status        string         // Custom per project
    Tags          []string
    Description   string         // Markdown body
    Relationships []Relationship
    CreatedAt     time.Time      // RFC3339 UTC
    UpdatedAt     time.Time      // RFC3339 UTC
}

type Relationship struct {
    Type   string // "parent"|"blocks"|"relates-to"
    TaskID int    // Reference to another task
}
```

**Why**: Simple, flexible, maps naturally to YAML frontmatter, supports all relationship types

---

## Decision Matrix

| Decision | Choice | Reason |
|----------|--------|--------|
| ID format | Global sequential integers | Simple, decoupled from type |
| Root-level tasks | Epic-only | Clear hierarchy, consistent |
| Epic dir naming | ID + slug | Self-documenting, discoverable |
| Slug length | 3-5 words | Balance clarity and path length |
| NextID logic | Max ID + 1 | No accidental ID reuse |
| Relationships | Slice of structs | Single source of truth |
| Query pattern | Functional options | Composable, extensible |
| Config | Optional TOML | Works without config |
| Storage | Pluggable interface | Different backends possible |

---

**All decisions made with consideration for**:
✅ Developer ergonomics
✅ File system discoverability
✅ Robustness to human mistakes
✅ Future extensibility
✅ MVP simplicity

See `.tasks/1-design/` for complete rationale and design documents.

---

## Implementation Status

**Completion Date**: November 2, 2025  
**Status**: ✅ COMPLETE

All design decisions have been successfully implemented in Phase 1 (Core System). See IMPLEMENTATION_COMPLETE.md for full details.

### Packages Implemented
- ✅ `internal/model/` - Task and Relationship structs (300 lines)
- ✅ `internal/storage/` - BaseStorage interface with MarkdownFileStorage and MemoryStorage (700 lines)
- ✅ `internal/config/` - ProjectConfig with TOML loading and defaults (250 lines)
- ✅ `internal/query/` - QueryEngine with 9 functional option filters (150 lines)
- ✅ `cmd/` - Full CLI with Cobra/Viper, all CRUD operations (400 lines)

### Quality Metrics
- Zero external dependencies for core logic (only Go stdlib + design dependencies)
- Clean separation of concerns
- Interface-based for extensibility
- Functional options pattern for filters
- Proper error handling throughout
- Consistent code style
- All functionality manually tested and verified

### What Works
- ✅ Task creation, listing, filtering, updating, deletion
- ✅ Markdown file storage with YAML frontmatter
- ✅ Hierarchical organization by epic
- ✅ Automatic sequential ID generation
- ✅ Relationship tracking (parent, blocks, relates-to)
- ✅ TOML configuration with defaults
- ✅ Environment variable overrides
- ✅ CLI flags override all
- ✅ Table formatting for output
- ✅ Status transitions
- ✅ Tag-based filtering

### Known Limitations (Intentional for MVP)
- Not thread-safe for concurrent writes (acceptable for typical use)
- O(n) list operations (fine for ~1000 tasks)
- No full-text search (Phase 2)
- No CLI relationship creation (Phase 2)
- No web UI (future)
- No MCP server (future)

### Module Info
- **Module**: github.com/zenobi-us/opentasks
- **Go Version**: 1.21+
- **Branch**: design
- **Dependencies**: cobra v1.7.0, viper v1.16.0, toml v1.5.0, yaml.v3 v3.0.1
