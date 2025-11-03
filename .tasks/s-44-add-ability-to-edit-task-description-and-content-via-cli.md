---
id: 44
title: Add ability to edit task description and content via CLI
type: story
status: todo
tags: [cli, tasks, ux]
relationships: []
createdAt: "2025-11-03T08:30:00Z"
updatedAt: "2025-11-03T08:30:00Z"
---

# Story: Add Task Content Editing to CLI

## Problem

Currently, the `opentask task update` command only supports updating status. There's no way to edit task descriptions, body content, or metadata (title, tags, etc.) from the command line. Users must edit task files directly.

This is a significant limitation when using opentask to manage design decisions and implementation plans that need detailed content.

## Acceptance Criteria

### Required Functionality

- [ ] `opentask task update <id> --title "new title"` - Update task title
- [ ] `opentask task update <id> --description "description text"` - Set or update description
- [ ] `opentask task update <id> --editor` - Open task in $EDITOR for full content editing
- [ ] `opentask task update <id> --tag "new-tag"` - Add tags to tasks
- [ ] `opentask task update <id> --remove-tag "old-tag"` - Remove tags from tasks
- [ ] Proper error handling for invalid operations
- [ ] All updates reflected immediately in task files

### Nice to Have

- [ ] `opentask task new --description <text>` - Add description when creating tasks
- [ ] Multi-line description support (heredoc or temp file)
- [ ] Batch tag operations
- [ ] Content preview before saving with `--editor`

## Implementation Plan

### Phase 1: Metadata Updates (--title, --tag flags)

1. Add flag definitions to `cmd/task.go` update command
2. Add `UpdateTitle()` method to task model
3. Add `AddTag()` and `RemoveTag()` methods to task model
4. Update `cmd/task.go` to apply these changes
5. Write tests for metadata updates

### Phase 2: Editor Support (--editor flag)

1. Detect EDITOR environment variable
2. Create temp file with task content in markdown
3. Launch editor (vim, nano, emacs, vscode, etc.)
4. Parse updated content back into task
5. Save to storage
6. Write tests with mock editor

### Phase 3: Description Flag (--description)

1. Add `--description` flag to update command
2. Support both single-line and file input
3. Merge with existing description or replace
4. Write tests

## Design Considerations

### Editor Integration

- Use `$EDITOR` environment variable (vim, nano, etc.)
- Fallback options if EDITOR not set
- Show task content in editor with clear formatting
- Preserve YAML frontmatter structure
- Handle user cancellation gracefully

### Multi-line Content

- Support piping content: `echo "description" | opentask task update 1 --description -`
- Support file input: `opentask task update 1 --description @file.txt`
- Support heredoc: `opentask task update 1 --description << 'EOF'` ... `EOF`

### Validation

- Validate tags are valid identifiers
- Prevent invalid status updates
- Preserve relationships during updates
- Maintain ID and timestamps

## Related Tasks

- Decision 43: Config discovery and display patterns
- This enhancement makes task management more practical

## References

- Current `cmd/task.go` update command
- Internal task model in `internal/model/task.go`
- Storage interface in `internal/storage/interface.go`
