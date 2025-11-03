---
id: 47
title: Task deletion by ID fails due to config discovery assumptions
type: task
status: todo
tags: [bug, task-deletion, config, critical]
relationships: []
createdAt: "2025-11-03T09:30:00Z"
updatedAt: "2025-11-03T09:30:00Z"
---

# Bug: Task Deletion by ID Fails

## Description

When deleting a task by ID (e.g., `opentask task delete 1`), the command fails or deletes from the wrong location. The issue appears to be that task deletion assumes a specific directory structure or config discovery pattern.

## Reproduction Steps

1. Have a project with `.opentask.toml` and `.tasks/` directory
2. From the project root: `opentask task delete 1` (works)
3. From a subdirectory: `opentask task delete 1` (fails or unexpected behavior)
4. From a completely different directory with explicit `--path`: `opentask task delete 1 --path /path/to/project` (may still fail)

## Expected Behavior

- Task deletion should work regardless of current working directory
- Task should be deleted from the resolved storage path
- Command should confirm deletion and show which task was deleted

## Actual Behavior

- Deletion may fail with config resolution error
- May delete wrong task or from wrong location
- No clear error message about what went wrong
- Appears to assume current directory contains the project

## Root Cause Analysis

The bug likely stems from:

1. **Config Discovery Logic** - `initializeStorage()` uses current directory as starting point for config discovery
2. **Storage Path Resolution** - Storage path may not be correctly resolved when called from non-project directory
3. **Task ID Resolution** - Assumes task files are in current directory or uses hardcoded relative paths
4. **Missing `--path` Handling** - The `--path` flag may not properly override config discovery

## Code Areas to Investigate

### `cmd/root.go` - `initializeStorage()`
```go
// Current behavior:
path := projectPath
if path == "" {
  path = "."  // <-- Assumes current directory is project
}

// Should allow:
// - --path flag to override
// - --config flag to specify exact config
// - discovery to work from anywhere
```

### `cmd/task.go` - `deleteCmd`
- Check how task ID is mapped to file path
- Verify storage path is used correctly
- Ensure error handling for missing tasks

### `internal/storage/markdown.go` - `DeleteTask()`
- Verify it uses correct storage path
- Check file naming assumptions
- Ensure ID resolution works

## Impact

**Severity:** HIGH
- Core functionality broken in non-trivial scenarios
- Workaround: Must run commands from project root

**Affected Commands:**
- `opentask task delete <id>`
- Possibly other commands that look up tasks by ID

## Acceptance Criteria

- [ ] `opentask task delete <id>` works from project root
- [ ] `opentask task delete <id>` works from subdirectory of project
- [ ] `opentask task delete <id>` works from arbitrary directory with `--path /project`
- [ ] `opentask task delete <id>` works with explicit `--config /path/to/.opentask.toml`
- [ ] Clear error message if task ID doesn't exist
- [ ] Clear success message showing what was deleted
- [ ] All existing tests pass
- [ ] New tests verify deletion works from different directories

## Investigation Checklist

- [ ] Reproduce the bug and document exact failure mode
- [ ] Trace config discovery flow from non-project directory
- [ ] Check how storage path gets resolved
- [ ] Verify task file lookup uses correct path
- [ ] Test `--path` flag handling
- [ ] Test `--config` flag handling
- [ ] Review error messages for clarity

## Solution Approach

1. **Ensure config resolution works from any directory**
   - Already implemented: `LoadConfigHierarchical()` walks up directories
   - Verify it's being used in `initializeStorage()`

2. **Verify storage initialization uses resolved path**
   - Check that `Store` is initialized with correct path
   - Verify `DeleteTask()` uses storage path, not cwd

3. **Add explicit `--path` and `--config` handling**
   - Make sure flags override discovery
   - Test all combinations

4. **Improve error messages**
   - Show which config file was used
   - Show which storage path is being used
   - Show clear error if task not found

5. **Add integration tests**
   - Test delete from project root
   - Test delete from subdirectory
   - Test delete from external directory with --path
   - Test delete with explicit --config
   - Test delete of non-existent task

## Related Issues

- Task 44: Task editing needs similar fixes for path handling
- Story 46: Config schema redesign will need to address path resolution
- Story 45: CI/CD should catch this with tests

## Files to Modify

- `cmd/root.go` - Verify config resolution
- `cmd/task.go` - Check deleteCmd implementation
- `internal/storage/markdown.go` - Verify DeleteTask uses correct path
- `cmd/task_test.go` - Add integration tests (if doesn't exist)

## Testing Commands

```bash
# Test 1: From project root (should work)
cd ~/Projects/opentask
opentask task delete 1

# Test 2: From subdirectory (currently fails?)
cd ~/Projects/opentask/.tasks
opentask task delete 1

# Test 3: With --path flag
cd /tmp
opentask task delete 1 --path ~/Projects/opentask

# Test 4: With --config flag
cd /tmp
opentask task delete 1 --config ~/Projects/opentask/.opentask.toml

# Test 5: Invalid task ID
opentask task delete 99999

# Test 6: View task before delete
opentask task show 1
opentask task delete 1
opentask task show 1  # Should fail
```

## Notes

- This bug is critical because task deletion is a core feature
- May indicate broader issues with how commands find and access tasks
- Should test ALL commands for similar path resolution issues:
  - `task show`
  - `task update`
  - `task new`
  - `task list`
