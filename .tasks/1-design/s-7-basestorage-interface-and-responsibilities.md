---
id: 7
title: BaseStorage Interface and Responsibilities
type: story
status: done
tags: [design, storage, interface]
relationships:
  - type: parent
    taskID: 1

createdAt: 2025-11-02T08:00:00Z
updatedAt: 2025-11-02T08:00:00Z
---

# BaseStorage Interface and Responsibilities

## Design Decision

Storage is abstracted behind a pluggable interface. The default implementation uses markdown files on the filesystem.

## BaseStorage Interface

```go
type TaskFilter func(*Task) bool

type BaseStorage interface {
    // Load/Save
    LoadTask(id int) (*Task, error)
    SaveTask(task *Task) error
    DeleteTask(id int) error
    
    // Querying with functional options
    ListTasks(filters ...TaskFilter) ([]*Task, error)
    
    // ID generation (storage backend responsible)
    NextID(ctx context.Context) (int, error)
    
    // Relationship resolution
    GetRelated(taskID int, relationType string) ([]*Task, error)
    
    Close() error
}
```

## Core Responsibilities

1. **Persistence**: Load and save tasks from/to underlying storage
2. **ID Generation**: Generate next global sequential ID (counts all files in project)
3. **Querying**: List tasks with functional option filters
4. **Relationship Resolution**: Resolve task links by loading related tasks
5. **File Organization**: Organize files by epic hierarchy (MarkdownFileStorage detail)

## Initialization

```go
type StorageConfig struct {
    Backend string            // e.g., "markdown-fs"
    Path    string            // Project path (scopes storage to one project)
    // Backend-specific config fields
}

storage, err := NewStorage(config)
```

## MarkdownFileStorage Implementation Details

- **Task files**: Stored as `{epic_id}-{slug}/{type}-{id}-{slug}.md` organized by epic
- **Frontmatter parsing**: Extract YAML metadata from task file header (id is integer)
- **Description**: Markdown body below frontmatter separator (`---`)
- **ID generation**: Count all `.md` files recursively to determine next global ID
- **Listing**: Walk directory tree, parse frontmatter for filtering
- **Epic hierarchy**: Related tasks grouped in epic subdirectories for discoverability

## Future Storage Backends

- SQLite backend (indexed queries, better performance)
- PostgreSQL backend (shared projects, concurrent access)
- Memory backend (testing)
- Cloud storage backend (Google Drive, Notion, etc.)

## Design Considerations

- Each BaseStorage instance is scoped to one project via config.Path
- Multi-project support handled at higher layer (ProjectManager)
- Storage backends are interchangeable—CLI doesn't care which is used
- Config passed to storage backend is immutable after initialization
