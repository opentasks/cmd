---
id: 18
title: Improve error messages
type: task
status: todo
tags:
    - polish
relationships:
    - type: parent
      taskID: 12
createdAt: "2025-11-02T08:41:05Z"
updatedAt: "2025-11-02T20:35:00Z"
---

## Objective
Enhance all error messages to be helpful, consistent, and actionable. Create an error glossary that maps error codes to detailed explanations and recovery steps.

## Current Error Patterns

Errors currently follow inconsistent patterns:
- Some are generic: "failed to generate task ID: <wrapped error>"
- Some lack context: "task not found"
- Some don't suggest solutions: "invalid task ID: foo"

## Error Message Philosophy

**Each error message should answer:**
1. **What** happened? (Be specific)
2. **Why** did it fail? (Provide context)
3. **How** to fix it? (Suggest action)

**Good error:** "Task 42 not found in project at ~/.tasks/. Create it with: `task new 'Title'`"
**Bad error:** "Task not found"

## Error Glossary

Create `docs/ERROR_GLOSSARY.md` documenting all error codes:

### Format
```markdown
## ERR-001: TASK_NOT_FOUND

### Message
"Task <id> not found in <storage_path>"

### Cause
The task ID was not found in the task storage.

### Solutions
1. Verify the task ID is correct: `task list | grep <id>`
2. Create the task: `task new "Title" --id <id>` (if ID auto-assignment)
3. Check the --path flag points to correct project

### Example
$ task view 999
Error: Task 999 not found in ~/.tasks/
Try: task list  # to see all available tasks
```

## Error Code Registry

Categorized by severity and context:

### Configuration Errors (ERR-100 to ERR-199)
- **ERR-101: CONFIG_NOT_FOUND** - Config file missing (use defaults)
- **ERR-102: CONFIG_PARSE_ERROR** - TOML syntax invalid
- **ERR-103: INVALID_WORKFLOW** - Workflow config invalid
- **ERR-104: TEMPLATE_NOT_FOUND** - Task template missing
- **ERR-105: INVALID_PATH** - Path resolution failed

### Task Errors (ERR-200 to ERR-299)
- **ERR-201: TASK_NOT_FOUND** - Task ID doesn't exist
- **ERR-202: INVALID_TASK_ID** - ID format invalid (not a number)
- **ERR-203: INVALID_TASK_TYPE** - Type not in (epic, plan, research, story, decision, task)
- **ERR-204: INVALID_STATUS** - Status not allowed by workflow
- **ERR-205: TASK_ALREADY_EXISTS** - Task ID already in use
- **ERR-206: DUPLICATE_RELATIONSHIP** - Relationship already exists

### Relationship Errors (ERR-300 to ERR-399)
- **ERR-301: CIRCULAR_RELATIONSHIP** - Parent relationship would create cycle
- **ERR-302: INVALID_RELATIONSHIP_TYPE** - Type not in (parent, blocks, relates-to)
- **ERR-303: PARENT_NOT_FOUND** - Parent task referenced doesn't exist
- **ERR-304: RELATIONSHIP_NOT_FOUND** - Relationship doesn't exist

### Storage Errors (ERR-400 to ERR-499)
- **ERR-401: STORAGE_NOT_INITIALIZED** - Storage backend failed to start
- **ERR-402: FILE_WRITE_ERROR** - Can't write task file
- **ERR-403: FILE_READ_ERROR** - Can't read task file
- **ERR-404: FILE_PARSE_ERROR** - Markdown/YAML format invalid
- **ERR-405: DIRECTORY_CREATE_ERROR** - Can't create storage directory

### Query Errors (ERR-500 to ERR-599)
- **ERR-501: INVALID_FILTER** - Filter syntax invalid
- **ERR-502: INVALID_QUERY** - Query malformed

## Implementation Plan

### 1. Define Error Types
Create `internal/errors/errors.go`:
```go
type ErrorCode string

const (
    TaskNotFound      ErrorCode = "ERR-201"
    InvalidTaskID     ErrorCode = "ERR-202"
    InvalidTaskType   ErrorCode = "ERR-203"
    CircularRelationship ErrorCode = "ERR-301"
    // ... more codes
)

type AppError struct {
    Code    ErrorCode
    Message string
    Context map[string]string  // Additional context
    Cause   error              // Wrapped error
}

func (e *AppError) Error() string { /* formatted output */ }
func (e *AppError) Suggest() string { /* recovery suggestions */ }
```

### 2. Update Error Creation
Replace generic `fmt.Errorf()` calls with structured errors:

```go
// Old
return fmt.Errorf("invalid task type: %s", taskType)

// New
return &AppError{
    Code: InvalidTaskType,
    Message: fmt.Sprintf("Task type '%s' is invalid", taskType),
    Context: map[string]string{
        "provided": taskType,
        "valid": strings.Join(model.AllTaskTypes, ", "),
    },
}
```

### 3. CLI Error Display
Update root command error handler to show helpful output:

```
Error [ERR-203]: Invalid task type 'epic2'
Valid types: epic, plan, research, story, decision, task

See: task new --help
For more info: docs/ERROR_GLOSSARY.md#ERR-203
```

## Acceptance Criteria
- [ ] Error glossary document created with all error codes
- [ ] All errors have unique error codes (ERR-XXX format)
- [ ] Each error message includes: what, why, how to fix
- [ ] Structured error type with code, message, suggestions
- [ ] All CLI commands use structured errors
- [ ] CLI displays error code and suggestions
- [ ] Tests verify error messages are helpful and consistent
- [ ] Help text links to error glossary
- [ ] Documentation complete and searchable by error code
