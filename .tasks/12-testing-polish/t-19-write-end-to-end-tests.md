---
id: 19
title: Write end-to-end tests
type: task
status: todo
tags:
    - testing
relationships:
    - type: parent
      taskID: 12
createdAt: "2025-11-02T08:41:05Z"
updatedAt: "2025-11-02T20:35:00Z"
---

## Objective
Test complete user workflows from CLI invocation through result verification. E2E tests validate that all components work together correctly from the user's perspective.

## Scope: Real User Workflows

### Testing Approach
- Use temporary directories (`t.TempDir()`) for filesystem isolation
- Invoke CLI via `os/exec` or internal CLI entry point
- Verify output format, exit codes, and file state
- Test with realistic task configurations

### 1. **Project Initialization Workflow**
Verify a user can initialize a project and start working:

```
1. Initialize project: opentask init --path /tmp/myproject
   - Verify: .tasks/ directory created
   - Verify: config.toml created with defaults
   - Verify: No tasks exist yet
   
2. Create tasks:
   - task new "Feature 1" --type story
   - task new "Feature 2" --type story
   - Verify: Files created with IDs 1, 2
   
3. List tasks
   - Verify: Both tasks appear in output
   - Verify: Default format is human-readable
```

### 2. **Task Creation and Management Workflow**

```
1. Create epic: task new "Phase 1" --type epic
2. Create stories under epic:
   - task new "Story A" --type story --parent 1
   - task new "Story B" --type story --parent 1
3. List all tasks (should see epic + 2 stories)
4. View task details: task view 1
   - Verify: Shows epic title and description
5. Update task: task update 1 --status in-progress
6. Delete task: task delete 3
   - Verify: Task removed
   - Verify: Relationships cleaned up
```

### 3. **Task Relationships Workflow**

```
1. Create tasks:
   - task new "Task A" (ID 1)
   - task new "Task B" (ID 2)
   - task new "Task C" (ID 3)
   
2. Link tasks:
   - task link add 1 2 --type blocks  (A blocks B)
   - task view 1  (verify shows B is blocked)
   - task view 2  (verify shows blocked by A)
   
3. Remove link:
   - task link remove 1 2
   - task view 1  (verify B no longer linked)
```

### 4. **Filtering and Querying Workflow**

```
1. Create diverse tasks:
   - task new "Epic 1" --type epic
   - task new "Research topic" --type research --tag analysis
   - task new "Decision point" --type decision --tag urgent
   - task new "Build feature" --type story --parent 1 --tag feature
   
2. Filter by type: task list --type story
   - Verify: Only stories appear
   
3. Filter by status: task list --status todo
   - Verify: Only todo tasks appear
   
4. Filter by tag: task list --tag urgent
   - Verify: Only urgent-tagged tasks appear
   
5. Combined filters: task list --type story --status in-progress
   - Verify: Only stories that are in-progress
```

### 5. **Template Usage Workflow**

```
1. Create custom template: .tasks/templates/feature-story.md
   - Content: pre-filled requirements, acceptance criteria
   
2. Create task from template:
   - task new "Build login page" --type story --template feature-story.md
   - Verify: Task includes template content
   - Verify: Template tags/relationships applied
   
3. Verify task content:
   - task view <id> (shows template sections)
```

### 6. **JSON/YAML Output Workflow**

```
1. Create task: task new "Test" --type story
   
2. View as JSON: task view 1 --format json
   - Verify: Valid JSON with all fields
   
3. View as YAML: task view 1 --format yaml
   - Verify: Valid YAML with all fields
   
4. List as JSON: task list --format json
   - Verify: Array of tasks with metadata
   
5. Parse and verify structure
   - All fields present and correct types
   - Timestamps in RFC3339 format
```

### 7. **Error Handling Workflow**

```
1. Create with invalid type:
   - task new "Title" --type invalid
   - Verify: Error message includes valid types
   
2. Reference non-existent parent:
   - task new "Title" --parent 999
   - Verify: Clear error that parent doesn't exist
   
3. Invalid task ID:
   - task view abc
   - Verify: Error explains ID must be numeric
   
4. Delete non-existent:
   - task delete 999
   - Verify: Helpful error message
```

### 8. **Configuration Workflow**

```
1. Create project with custom config:
   - Write config.toml with custom statuses
   - Create task with custom status
   - Verify: Works with custom workflow
   
2. Test status transitions:
   - Verify workflow limits transitions correctly
   
3. Test path resolution:
   - Use --path flag
   - Use opentask_PROJECT_PATH env var
   - Verify: Same project accessed
```

### 9. **Large Dataset Workflow**

```
1. Create 100+ tasks with relationships
2. Verify queries work correctly with large dataset
3. Verify performance is acceptable
4. Verify file organization with nested epics
```

### 10. **Recovery and Edge Cases**

```
1. Corrupted file handling:
   - Create task, corrupt YAML, attempt load
   - Verify: Graceful error
   
2. Missing parent:
   - Create task with parent
   - Delete parent task
   - Verify: Child task still exists (orphaned)
   
3. Circular dependencies (if applicable):
   - Attempt to create circular parent relationships
   - Verify: Prevented with clear error
```

## Test Implementation Strategy

### Test Structure
Create `internal/tests/e2e/` with tests for each workflow:
- `TestProjectInitialization`
- `TestTaskCreationWorkflow`
- `TestTaskRelationshipsWorkflow`
- `TestFilteringWorkflow`
- `TestTemplateWorkflow`
- `TestJSONYAMLOutput`
- `TestErrorHandling`
- `TestConfigurationWorkflow`
- `TestLargeDatasetWorkflow`
- `TestEdgeCases`

### CLI Invocation Approach
Two options:
1. **Direct CLI call**: Import cmd package, call handlers directly
2. **Exec subprocess**: Run `opentask` binary directly

Use direct calls for speed, subprocess for true integration testing.

### Assertion Patterns
```go
// Verify command succeeds
err := runCommand("task", "new", "Title", "--type", "story")
require.NoError(t, err)

// Verify task exists
tasks, err := store.ListTasks(ctx)
require.Len(t, tasks, 1)
require.Equal(t, "Title", tasks[0].Title)

// Verify output format
output := runCommandOutput("task", "view", "1", "--format", "json")
var task Task
require.NoError(t, json.Unmarshal([]byte(output), &task))
```

## Acceptance Criteria
- [ ] All 10 major workflows have E2E tests
- [ ] Tests use isolated temporary directories
- [ ] Tests are independent and can run in parallel
- [ ] Tests verify both success and error paths
- [ ] Output formatting verified (text, JSON, YAML)
- [ ] File state verified after operations
- [ ] All tests pass: `go test ./...`
- [ ] No race conditions: `go test -race ./...`
- [ ] Tests execute in reasonable time (<5 seconds total)
