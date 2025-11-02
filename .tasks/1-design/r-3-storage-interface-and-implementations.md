---
id: 3
title: Storage Interface and Implementations
type: research
status: done
tags: [design, storage, interface, reference]
relationships:
  - type: parent
    taskID: 1
createdAt: 2025-11-02T08:00:00Z
updatedAt: 2025-11-02T08:00:00Z
---

# Storage Interface and Implementations

Complete specification of the BaseStorage interface and MarkdownFileStorage implementation details.

## BaseStorage Interface

The abstraction layer for task persistence. All storage backends must implement this interface.

```go
// NOTE: This is a design specification, not production code.
// These interfaces and types are guidelines for implementation.
// Adjust method signatures, error handling, and organization based on actual needs.
// Verify against requirements before implementing.

// TaskFilter is a functional option for filtering tasks
type TaskFilter func(*model.Task) bool

// BaseStorage defines the interface all storage backends must implement
type BaseStorage interface {
    // LoadTask retrieves a single task by ID
    // Returns ErrTaskNotFound if task doesn't exist
    LoadTask(ctx context.Context, id int) (*model.Task, error)
    
    // SaveTask persists a task to storage
    // Creates if task doesn't exist, updates if it does
    SaveTask(ctx context.Context, task *model.Task) error
    
    // DeleteTask removes a task from storage
    // Returns ErrTaskNotFound if task doesn't exist
    DeleteTask(ctx context.Context, id int) error
    
    // ListTasks returns all tasks matching the given filters
    // If no filters provided, returns all tasks
    // Filters are applied in AND fashion (all must match)
    ListTasks(ctx context.Context, filters ...TaskFilter) ([]*model.Task, error)
    
    // NextID generates the next global sequential ID
    // Counts all task files in project, returns count + 1
    // Guaranteed globally unique (no collisions)
    NextID(ctx context.Context) (int, error)
    
    // GetRelated returns all tasks related to the given task by relationship type
    // relationType must be one of: "parent", "blocks", "relates-to"
    GetRelated(ctx context.Context, taskID, relationType string) ([]*model.Task, error)
    
    // Close performs cleanup (if needed for this backend)
    Close() error
}

// Common errors
var (
    ErrTaskNotFound = errors.New("task not found")
    ErrInvalidID = errors.New("invalid task ID format")
    ErrInvalidTaskType = errors.New("invalid task type")
    ErrInvalidStatus = errors.New("invalid status for workflow")
    ErrCircularRelationship = errors.New("circular relationship detected")
)

// StorageConfig contains backend-agnostic configuration
type StorageConfig struct {
    Backend string            // "markdown-fs", "sqlite", etc.
    Path    string            // Project path or database location
    Options map[string]string // Backend-specific options
}

// NewStorage factory function creates storage backend based on config
func NewStorage(config StorageConfig) (BaseStorage, error) {
    switch config.Backend {
    case "markdown-fs":
        return NewMarkdownFileStorage(config.Path)
    case "memory":
        return NewMemoryStorage(), nil
    default:
        return nil, errors.New("unknown storage backend: " + config.Backend)
    }
}
```

## MarkdownFileStorage Implementation

Reference implementation using markdown files with YAML frontmatter on the filesystem.

```go
// NOTE: This is a design specification and reference implementation.
// Use this as a guide, but implement based on actual requirements.
// Error handling, performance considerations, and edge cases may need adjustment.
// Test thoroughly and verify against design before considering complete.

package storage

import (
    "context"
    "fmt"
    "os"
    "path/filepath"
    "regexp"
    "sort"
    "strconv"
    "strings"
    "time"
    
    "github.com/yourusername/opentask/model"
    "gopkg.in/yaml.v3"
)

type MarkdownFileStorage struct {
    basePath string // Project root path
}

// NewMarkdownFileStorage creates a new markdown file storage backend
func NewMarkdownFileStorage(basePath string) (*MarkdownFileStorage, error) {
    // Ensure path exists
    if err := os.MkdirAll(basePath, 0755); err != nil {
        return nil, err
    }
    
    return &MarkdownFileStorage{basePath: basePath}, nil
}

// File naming: <epic_id>-<epic_slug>/<typecode>-<id>-<slug>.md
// Examples:
//   1-design-opentask/s-5-task-data-model.md
//   1-design-opentask/s-6-semantic-id-system.md
//   e-1-design-opentask.md (epic at root)
// NOTE: Implementation guidance - adjust as needed for actual requirements
func (s *MarkdownFileStorage) taskToPath(task *model.Task) (string, error) {
    // Determine parent epic directory
    var epicDir string
    var parentEpicID int
    
    for _, rel := range task.Relationships {
        if rel.Type == model.RelParent {
            parentEpicID = rel.TaskID
            break
        }
    }
    
    if parentEpicID > 0 {
        // Load parent epic to get its title
        parentEpic, err := s.LoadTask(context.Background(), parentEpicID)
        if err != nil {
            return "", fmt.Errorf("parent epic not found: %w", err)
        }
        
        epicSlug := slugify(parentEpic.Title)
        epicDir = fmt.Sprintf("%d-%s", parentEpicID, epicSlug)
    }
    
    // Build filename: typecode-id-slug.md
    typeCode := model.TypeCode[task.Type]
    slug := slugify(task.Title)
    filename := fmt.Sprintf("%s-%d-%s.md", typeCode, task.ID, slug)
    
    // Combine directory and filename
    if epicDir != "" {
        return filepath.Join(s.basePath, epicDir, filename), nil
    }
    return filepath.Join(s.basePath, filename), nil
}

// pathToTaskID extracts task ID from file path
// Parses filename like "s-42-task-title.md" and extracts ID 42
func (s *MarkdownFileStorage) pathToTaskID(filename string) (int, error) {
    // Remove .md extension
    if !strings.HasSuffix(filename, ".md") {
        return 0, ErrInvalidID
    }
    
    base := strings.TrimSuffix(filename, ".md")
    parts := strings.SplitN(base, "-", 3)
    
    if len(parts) < 2 {
        return 0, ErrInvalidID
    }
    
    // parts[0] is typecode (s, e, p, etc.)
    // parts[1] is the ID number
    
    id, err := strconv.Atoi(parts[1])
    if err != nil {
        return 0, ErrInvalidID
    }
    
    return id, nil
}

// LoadTask loads task frontmatter and content from markdown file
// Task ID is extracted from frontmatter, not from filename
func (s *MarkdownFileStorage) LoadTask(ctx context.Context, id int) (*model.Task, error) {
    // Walk project to find file with matching ID in frontmatter
    var foundPath string
    
    err := filepath.Walk(s.basePath, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        
        if !strings.HasSuffix(path, ".md") {
            return nil
        }
        
        task, err := s.parseTaskFile(path)
        if err != nil {
            return nil // Skip files that fail to parse
        }
        
        if task.ID == id {
            foundPath = path
            return filepath.SkipDir // Found it, stop walking
        }
        
        return nil
    })
    
    if err != nil {
        return nil, err
    }
    
    if foundPath == "" {
        return nil, ErrTaskNotFound
    }
    
    content, err := os.ReadFile(foundPath)
    if err != nil {
        return nil, err
    }
    
    return s.parseTaskFile(string(content))
}

// parseTaskFile parses markdown file into Task
// Format: frontmatter (YAML) separator (---) markdown body
func (s *MarkdownFileStorage) parseTaskFile(filePath string) (*model.Task, error) {
    content, err := os.ReadFile(filePath)
    if err != nil {
        return nil, err
    }
    
    // Split by --- separator
    parts := strings.SplitN(string(content), "---", 3)
    if len(parts) < 3 {
        return nil, errors.New("invalid task file format: missing frontmatter")
    }
    
    frontmatterYAML := strings.TrimSpace(parts[1])
    description := strings.TrimSpace(parts[2])
    
    // Parse YAML frontmatter
    var frontmatter struct {
        ID            int `yaml:"id"`
        Title         string `yaml:"title"`
        Type          string `yaml:"type"`
        Status        string `yaml:"status"`
        Tags          []string `yaml:"tags"`
        Relationships []struct {
            Type   string `yaml:"type"`
            TaskID int `yaml:"taskID"`
        } `yaml:"relationships"`
        CreatedAt string `yaml:"createdAt"`
        UpdatedAt string `yaml:"updatedAt"`
    }
    
    if err := yaml.Unmarshal([]byte(frontmatterYAML), &frontmatter); err != nil {
        return nil, fmt.Errorf("failed to parse frontmatter: %w", err)
    }
    
    // Parse timestamps
    createdAt, err := time.Parse(time.RFC3339, frontmatter.CreatedAt)
    if err != nil {
        return nil, fmt.Errorf("invalid createdAt: %w", err)
    }
    
    updatedAt, err := time.Parse(time.RFC3339, frontmatter.UpdatedAt)
    if err != nil {
        return nil, fmt.Errorf("invalid updatedAt: %w", err)
    }
    
    // Convert frontmatter relationships
    relationships := make([]model.Relationship, len(frontmatter.Relationships))
    for i, rel := range frontmatter.Relationships {
        relationships[i] = model.Relationship{
            Type:   rel.Type,
            TaskID: rel.TaskID,
        }
    }
    
    task := &model.Task{
        ID:            frontmatter.ID,
        Title:         frontmatter.Title,
        Type:          frontmatter.Type,
        Status:        frontmatter.Status,
        Tags:          frontmatter.Tags,
        Relationships: relationships,
        Description:   description,
        CreatedAt:     createdAt,
        UpdatedAt:     updatedAt,
    }
    
    return task, nil
}

// SaveTask writes task to markdown file
func (s *MarkdownFileStorage) SaveTask(ctx context.Context, task *model.Task) error {
    filename, err := s.idToFilename(task.ID)
    if err != nil {
        return err
    }
    
    filepath := filepath.Join(s.basePath, filename)
    
    // Build frontmatter
    frontmatter := struct {
        ID            string `yaml:"id"`
        Title         string `yaml:"title"`
        Type          string `yaml:"type"`
        Status        string `yaml:"status"`
        Tags          []string `yaml:"tags"`
        Relationships []struct {
            Type   string `yaml:"type"`
            TaskID string `yaml:"taskID"`
        } `yaml:"relationships"`
        CreatedAt string `yaml:"createdAt"`
        UpdatedAt string `yaml:"updatedAt"`
    }{
        ID:        task.ID,
        Title:     task.Title,
        Type:      task.Type,
        Status:    task.Status,
        Tags:      task.Tags,
        CreatedAt: task.CreatedAt.Format(time.RFC3339),
        UpdatedAt: task.UpdatedAt.Format(time.RFC3339),
    }
    
    // Convert relationships
    frontmatter.Relationships = make([]struct {
        Type   string `yaml:"type"`
        TaskID string `yaml:"taskID"`
    }, len(task.Relationships))
    for i, rel := range task.Relationships {
        frontmatter.Relationships[i].Type = rel.Type
        frontmatter.Relationships[i].TaskID = rel.TaskID
    }
    
    // Marshal YAML
    yamlBytes, err := yaml.Marshal(frontmatter)
    if err != nil {
        return err
    }
    
    // Build file content
    content := fmt.Sprintf("---\n%s---\n\n%s\n", string(yamlBytes), task.Description)
    
    // Write file
    return os.WriteFile(filepath, []byte(content), 0644)
}

// DeleteTask removes the task file
func (s *MarkdownFileStorage) DeleteTask(ctx context.Context, id int) error {
    task, err := s.LoadTask(ctx, id)
    if err != nil {
        return err
    }
    
    path, err := s.taskToPath(task)
    if err != nil {
        return err
    }
    
    err = os.Remove(path)
    if err != nil && os.IsNotExist(err) {
        return ErrTaskNotFound
    }
    return err
}

// ListTasks returns all tasks matching filters
func (s *MarkdownFileStorage) ListTasks(ctx context.Context, filters ...TaskFilter) ([]*model.Task, error) {
    var tasks []*model.Task
    
    // Walk all directories and files
    err := filepath.Walk(s.basePath, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        
        if !strings.HasSuffix(path, ".md") {
            return nil
        }
        
        task, err := s.parseTaskFile(path)
        if err != nil {
            return nil // Skip files that fail to parse
        }
        
        // Apply filters
        match := true
        for _, filter := range filters {
            if !filter(task) {
                match = false
                break
            }
        }
        
        if match {
            tasks = append(tasks, task)
        }
        
        return nil
    })
    
    if err != nil {
        return nil, err
    }
    
    // Sort by creation time for consistent ordering
    sort.Slice(tasks, func(i, j int) bool {
        return tasks[i].CreatedAt.Before(tasks[j].CreatedAt)
    })
    
    return tasks, nil
}

// NextID generates the next global sequential ID
// Finds the maximum ID from all task frontmatter and returns max + 1
// This ensures IDs are never reused even if tasks are deleted
func (s *MarkdownFileStorage) NextID(ctx context.Context) (int, error) {
    var maxID int
    
    err := filepath.Walk(s.basePath, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return nil  // Continue walking on errors
        }
        
        if !strings.HasSuffix(path, ".md") {
            return nil
        }
        
        task, err := s.parseTaskFile(path)
        if err != nil {
            return nil  // Skip unparseable files
        }
        
        if task.ID > maxID {
            maxID = task.ID
        }
        
        return nil
    })
    
    if err != nil {
        return 0, err
    }
    
    return maxID + 1, nil
}

// GetRelated returns all tasks related by relationship type
func (s *MarkdownFileStorage) GetRelated(ctx context.Context, taskID int, relationType string) ([]*model.Task, error) {
    allTasks, err := s.ListTasks(ctx)
    if err != nil {
        return nil, err
    }
    
    var related []*model.Task
    for _, task := range allTasks {
        for _, rel := range task.Relationships {
            if rel.Type == relationType && rel.TaskID == taskID {
                related = append(related, task)
            }
        }
    }
    
    return related, nil
}

// Close is a no-op for file storage
func (s *MarkdownFileStorage) Close() error {
    return nil
}
```

## Storage Backend Considerations

### Path Resolution

- All paths in config are relative to the config file location
- If no config provided, uses project root (.)
- Environment variable `opentask_PROJECT_PATH` overrides config

### Concurrency

- Current markdown-fs implementation is not thread-safe
- Each operation acquires implicit file-level locks via OS
- For concurrent access, use SQLite or similar backend

### Performance

- Markdown-fs: O(n) for list operations (scans all files)
- Acceptable for projects up to ~1000 tasks
- For larger projects, migrate to indexed backend (SQLite/DB)

### Validation

- Storage doesn't validate business rules (e.g., circular relationships)
- Higher layers (QueryEngine) responsible for validation
- Storage focuses on persistence and retrieval
