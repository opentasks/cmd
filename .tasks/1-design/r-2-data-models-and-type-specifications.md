---
id: 2
title: Data Models and Type Specifications
type: research
status: done
tags: [design, data-model, reference]
relationships:
  - type: parent
    taskID: 1
createdAt: 2025-11-02T08:00:00Z
updatedAt: 2025-11-02T08:00:00Z
---

# Data Models and Type Specifications

Complete type specifications for the OpenTasks system. This is a reference document for implementation.

## Core Types

### Task

The fundamental unit in OpenTasks. Represents any trackable work item.

```go
// NOTE: This is a design specification, not production code.
// Implement based on this design but verify against actual requirements.
// Type names, field names, and organization are guidelines - adjust as needed.

package model

import "time"

type Task struct {
    // Identity and Metadata
    ID          int            // Global sequential ID (e.g., 42, 5)
    Title       string         // Short description
    Type        string         // epic|plan|research|story|decision|task
    Status      string         // Customizable per project (e.g., "todo", "in-progress")
    
    // Content
    Description string         // Markdown body (everything after frontmatter)
    
    // Organization
    Tags          []string     // Labels (e.g., ["feature", "core", "urgent"])
    Relationships []Relationship // Links to other tasks
    
    // Timestamps
    CreatedAt   time.Time      // RFC3339 UTC
    UpdatedAt   time.Time      // RFC3339 UTC
}

// Task type constants
const (
    TypeEpic     = "epic"
    TypePlan     = "plan"
    TypeResearch = "research"
    TypeStory    = "story"
    TypeDecision = "decision"
    TypeTask     = "task"
)

// All valid task types
var AllTaskTypes = []string{
    TypeEpic,
    TypePlan,
    TypeResearch,
    TypeStory,
    TypeDecision,
    TypeTask,
}

// Type to code mapping (for file naming)
var TypeCode = map[string]string{
    TypeEpic:     "e",
    TypePlan:     "p",
    TypeResearch: "r",
    TypeStory:    "s",
    TypeDecision: "d",
    TypeTask:     "t",
}

// Code to type mapping (inverse, for file parsing)
var CodeType = map[string]string{
    "e": TypeEpic,
    "p": TypePlan,
    "r": TypeResearch,
    "s": TypeStory,
    "d": TypeDecision,
    "t": TypeTask,
}
```

### Relationship

Links between tasks. Represents dependencies, hierarchies, and connections.

```go
type Relationship struct {
    Type   string // "parent"|"blocks"|"relates-to"
    TaskID int    // Target task ID (e.g., 42, 5)
}

// Relationship type constants
const (
    RelParent    = "parent"      // Hierarchical parent
    RelBlocks    = "blocks"      // This task blocks another
    RelRelatedTo = "relates-to"  // Related but independent
)
```

### Semantic ID

Task IDs are simple global sequential integers. Type is stored in frontmatter metadata, not in the ID.

```go
// Examples:
// 1       (first task - likely an epic)
// 42      (42nd task - could be any type)
// 1000    (1000th task)

// ID validation regex pattern
const IDPattern = `^\d+$`

// ID format components:
// Just a number: \d+ (one or more digits, no leading zeros)
// Type determined by: type field in frontmatter (epic|plan|research|story|decision|task)
// Uniqueness guaranteed by: global sequential counter across entire project
```

## Configuration Types

### ProjectConfig

Project-level configuration loaded from `config.toml`.

```go
package config

import "github.com/BurntSushi/toml"

type ProjectConfig struct {
    Project struct {
        Name        string `toml:"name"`
        Description string `toml:"description"`
        Owner       string `toml:"owner"`
    } `toml:"project"`
    
    Workflow WorkflowConfig `toml:"workflow"`
    Templates TemplateConfig `toml:"templates"`
    Storage StorageConfig `toml:"storage"`
}

type WorkflowConfig struct {
    Statuses []string `toml:"statuses"`
    Initial  string   `toml:"initial"`
    Transitions []TransitionConfig `toml:"transitions"`
}

type TransitionConfig struct {
    From string   `toml:"from"`
    To   []string `toml:"to"`
}

type TemplateConfig struct {
    Epic     string `toml:"epic"`
    Plan     string `toml:"plan"`
    Research string `toml:"research"`
    Story    string `toml:"story"`
    Decision string `toml:"decision"`
    Task     string `toml:"task"`
}

type StorageConfig struct {
    Backend string `toml:"backend"`
    Path    string `toml:"path"`
}
```

### Default Workflow

When no config is provided:

```go
var DefaultWorkflow = WorkflowConfig{
    Statuses: []string{"todo", "in-progress", "reviewing", "done", "archived"},
    Initial: "todo",
    Transitions: []TransitionConfig{
        {From: "todo", To: []string{"in-progress", "archived"}},
        {From: "in-progress", To: []string{"reviewing", "todo", "archived"}},
        {From: "reviewing", To: []string{"done", "in-progress", "archived"}},
        {From: "done", To: []string{"archived"}},
    },
}
```

### Default Storage Config

```go
var DefaultStorageConfig = StorageConfig{
    Backend: "markdown-fs",
    Path: ".",  // Project root
}
```

## Task Frontmatter Schema

Tasks are stored as markdown files with YAML frontmatter. Below is the canonical schema:

```yaml
# NOTE: This is a design specification. Adjust field names and structure as needed
# during implementation. These are guidelines, not absolute requirements.
id: 42
title: Implement task linking
type: story
status: in-progress
tags: [feature, core, urgent]
relationships:
  - type: parent
    taskID: 5
  - type: blocks
    taskID: 10
  - type: relates-to
    taskID: 3
createdAt: 2025-11-02T10:00:00Z
updatedAt: 2025-11-02T10:30:00Z
```

All fields are required in the frontmatter. Default values:
- `status`: First status in workflow (usually "todo")
- `tags`: Empty array
- `relationships`: Empty array
- `createdAt`/`updatedAt`: Current time when created

## File Organization (MarkdownFileStorage)

Files are organized hierarchically by epic with human-readable names:

```
<epic_id>-<epic_slug>/<typecode>-<id>-<slug>.md
```

Example structure:

```
1-design-opentasks/
  p-3-design-roadmap.md
  r-4-evaluate-storage.md
  s-5-task-data-model.md
2-implement-core/
  s-11-build-storage.md
  t-13-write-tests.md
e-1-design-opentasks.md
e-2-implement-core.md
```

The filename is a storage backend concern. The ID in the frontmatter is the source of truth.

## Error Types

```go
package errors

type ValidationError struct {
    Field   string
    Message string
}

type TaskError struct {
    TaskID    string
    Operation string
    Err       error
}

// Common error scenarios:
// - Task not found
// - Invalid ID format
// - Invalid status transition
// - Circular relationships
// - Template not found
// - Config parsing failed
```

## Common Patterns

### Filtering by Parent

To find all children of epic e-5:

```go
children := make([]*Task, 0)
for _, task := range allTasks {
    for _, rel := range task.Relationships {
        if rel.Type == RelParent && rel.TaskID == "e-5" {
            children = append(children, task)
        }
    }
}
```

### Finding Blocking Tasks

To find all tasks that block task s-42:

```go
blockers := make([]*Task, 0)
for _, task := range allTasks {
    for _, rel := range task.Relationships {
        if rel.Type == RelBlocks && rel.TaskID == "s-42" {
            blockers = append(blockers, task)
        }
    }
}
```

## TOML Parsing Considerations

- Use `BurntSushi/toml` for parsing
- All config sections are optional
- Unset fields use Go zero values or custom defaults
- Path fields are resolved relative to config file location
