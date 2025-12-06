# End-to-End Tests

## Quick Summary

The E2E test suite validates the complete opentasks system including:
- **Task Management**: Create, read, update, delete tasks via API and CLI
- **Task Hierarchies**: Multi-level parent-child relationships (Epic → Story → Task → Subtask)
- **Task Relationships**: Blocks, relates-to, and parent relationships
- **CLI Interface**: Full command-line interface for task operations
- **Storage**: Both file-based and in-memory storage backends
- **Query Engine**: Hierarchical queries and relationship traversal

## Testing Steps for Humans

### Prerequisites
```bash
# Build the opentask binary
mise run build

# Ensure binary is in PATH or in ./bin/opentask
which opentask  # or ls bin/opentask
```

### Run All E2E Tests
```bash
# Run all e2e tests
mise run test-e2e

# Run specific test file
go test -v ./test/e2e -run TestE2E_CLITaskCreation

# Run with verbose output and fail-fast
go test -v -failfast ./test/e2e
```

### Test Categories

#### 1. Smoke Tests (Basic Infrastructure)
```bash
go test -v ./test/e2e -run TestE2E_SmokeTest
```
Validates that the test infrastructure itself works correctly.

#### 2. CLI Tests (Command Line Interface)
```bash
go test -v ./test/e2e -run TestE2E_CLITaskCreation
go test -v ./test/e2e -run TestE2E_CLIErrorHandling
```
Tests task operations via the command-line interface.

#### 3. Hierarchy Tests (Parent-Child Relationships)
```bash
go test -v ./test/e2e -run TestE2E_FourLevelHierarchy
go test -v ./test/e2e -run TestE2E_EpicDeletionOrphaning
```
Tests multi-level task hierarchies and relationship integrity.

#### 4. Workflow Tests (Complete Feature Cycles)
```bash
go test -v ./test/e2e -run TestE2E_ResearchPlanningImplementationWorkflow
go test -v ./test/e2e -run TestE2E_ParallelStoryWorkflow
```
Tests realistic development workflows with multiple task types.

---

## Detailed Test Steps by Section

### 1. Smoke Tests (`smoke_test.go`)

**Purpose**: Validate the E2E test infrastructure and basic operations work correctly.

#### Test: `TestE2E_SmokeTest`
Validates core functionality:

| Step | Operation | Verification |
|------|-----------|--------------|
| 1 | Create temp directory via `SetupE2EEnvironment()` | Directory exists and is writable |
| 2 | Create epic via `testutil.CreateEpic()` | Epic has ID, title, and status |
| 3 | Load epic from storage | Loaded task matches created task |
| 4 | Verify file created on disk | At least one markdown file exists in temp dir |
| 5 | Create story with parent relationship | Story has parent relationship to epic |
| 6 | Query children via `engine.FindChildren()` | Query returns 1 child (the story) |
| 7 | Create task in memory storage | Task exists in memory without file I/O |

**Key Assertions**:
```go
epic.ID != 0                           // Task assigned unique ID
epic.Title == "Authentication System"  // Title preserved
loaded.Title == epic.Title             // Storage round-trip works
len(files) > 0                         // File created on disk
len(children) == 1                     // Query returns correct count
```

#### Test: `TestE2E_FixtureBuilders`
Validates all fixture builder helper functions:

| Step | Operation | Verification |
|------|-----------|--------------|
| 1 | Create epic with all options | Title, status, description, tags all set |
| 2 | Create story with parent and blocks | 2 relationships present (parent + blocks) |
| 3 | Verify parent relationship exists | `rel.Type == "parent"` and `rel.TaskID == epic.ID` |
| 4 | Verify blocks relationship exists | `rel.Type == "blocks"` relationship found |
| 5 | Create task with minimal options | Default status is "todo" |

**Key Assertions**:
```go
epic.Title == "Complex Epic"
epic.Status == "active"
len(epic.Tags) == 2
len(story.Relationships) == 2
task.Status == "todo"  // Default
```

---

### 2. CLI Tests (`cli_test.go`)

**Purpose**: Test the command-line interface for task operations.

#### Test: `TestE2E_CLITaskCreation` → Subtest: "create epic via CLI"
Tests creating tasks through the CLI:

| Step | Operation | CLI Command | Verification |
|------|-----------|-------------|--------------|
| 1 | Setup project config | Write `.opentask.toml` | Config file created |
| 2 | Execute CLI create command | `opentask task new "My Epic" --type epic --status planning` | Exit code = 0 |
| 3 | Check stdout/stderr | Output contains "created" | No errors reported |
| 4 | Load task from storage | `store.ListTasks()` | Task exists in storage |
| 5 | Verify task properties | Title, type, status match | All fields correct |
| 6 | Verify file on disk | `Assert(epic, env).FileExists()` | Markdown file created |

**Key Assertions**:
```go
exitCode == 0
strings.Contains(stdout+stderr, "created")
epic.Title == "My Epic"
epic.Type == model.TypeEpic
epic.Status == "planning"
```

#### Test: `TestE2E_CLITaskCreation` → Subtest: "list tasks via CLI"
Tests listing tasks via CLI:

| Step | Operation | CLI Command | Verification |
|------|-----------|-------------|--------------|
| 1 | Create 3 test tasks | Use testutil helpers | Tasks in storage |
| 2 | Execute list command | `opentask task list` | Exit code = 0 |
| 3 | Parse output | Check stdout for task titles | All titles present |

**Key Assertions**:
```go
exitCode == 0
output.Contains("Epic One")
output.Contains("Story Two")
output.Contains("Task Three")
```

#### Test: `TestE2E_CLITaskCreation` → Subtest: "show task via CLI"
Tests retrieving task details via CLI:

| Step | Operation | CLI Command | Verification |
|------|-----------|-------------|--------------|
| 1 | Create task with description | `testutil.CreateTask()` | Task in storage |
| 2 | Execute show command | `opentask task show <id>` | Exit code = 0 |
| 3 | Verify output | Check stdout contains title and description | Both fields displayed |

**Key Assertions**:
```go
exitCode == 0
output.Contains("Show Me Task")
output.Contains("This is a test description")
```

#### Test: `TestE2E_CLITaskCreation` → Subtest: "update task via CLI"
Tests updating task properties via CLI:

| Step | Operation | CLI Command | Verification |
|------|-----------|-------------|--------------|
| 1 | Create task with status "todo" | `testutil.CreateTask()` | Task created |
| 2 | Execute update command | `opentask task update <id> --status done` | Exit code = 0 |
| 3 | Reload from storage | `store.LoadTask()` | New status persisted |

**Key Assertions**:
```go
original.Status == "todo"
updated.Status == "done"  // After reload
```

#### Test: `TestE2E_CLIErrorHandling` → Subtest: "no error when no project initialized"
Tests graceful error handling when project is missing:

| Step | Operation | Expected Behavior |
|------|-----------|------------------|
| 1 | Create empty temp directory | No config file |
| 2 | Execute CLI list | Exit code != 0 |
| 3 | Check error message | Output mentions "no active project" or "not found" |

**Key Assertions**:
```go
exitCode != 0
output.ToLower().Contains("no active project") ||
output.ToLower().Contains("not found") ||
output.ToLower().Contains("project")
```

#### Test: `TestE2E_CLIErrorHandling` → Subtest: "error on invalid task ID"
Tests error when showing non-existent task:

| Step | Operation | Expected Behavior |
|------|-----------|------------------|
| 1 | Execute show with ID 999 | Task doesn't exist |
| 2 | Check exit code | Exit code != 0 |
| 3 | Check error message | Output mentions "not found" or "error" |

**Key Assertions**:
```go
exitCode != 0
output.ToLower().Contains("not found") ||
output.ToLower().Contains("error")
```

#### Test: `TestE2E_CLIErrorHandling` → Subtest: "error on missing required flags"
Tests error when required arguments are missing:

| Step | Operation | Expected Behavior |
|------|-----------|------------------|
| 1 | Execute create without title | Missing required positional arg |
| 2 | Check exit code or output | Either exitCode != 0 OR mentions "required" |
| 3 | Verify helpful message | Output shows usage or requirement info |

**Key Assertions**:
```go
exitCode != 0 ||
output.ToLower().Contains("required") ||
output.ToLower().Contains("usage")
```

#### Test: `TestE2E_CLIErrorHandling` → Subtest: "error on invalid task type"
Tests validation of task type field:

| Step | Operation | Expected Behavior |
|------|-----------|------------------|
| 1 | Create task with type "invalid-type" | Invalid enum value |
| 2 | Check exit code or output | Either exitCode != 0 OR mentions "invalid" |
| 3 | Verify error message | Output indicates type validation failure |

**Key Assertions**:
```go
exitCode != 0 ||
output.ToLower().Contains("invalid") ||
output.ToLower().Contains("error")
```

---

### 3. Hierarchy Tests (`hierarchy_test.go`)

**Purpose**: Test parent-child relationships and task hierarchies.

#### Test: `TestE2E_FourLevelHierarchy` → Subtest: "create four-level hierarchy"
Tests creating a complete task hierarchy:

| Level | Task Type | Operation | Verification |
|-------|-----------|-----------|--------------|
| 1 (Epic) | Epic | `CreateEpic("Authentication System")` | Type = Epic, Status = active |
| 2 (Story) | Story | `CreateStory(parent=epic.ID)` | HasParent(epic.ID) |
| 3 (Task) | Task | `CreateTask(parent=story.ID)` | HasParent(story.ID) |
| 4 (Subtask) | Task | `CreateTask(parent=task.ID)` | HasParent(task.ID) |
| Verify | - | `engine.FindChildren()` at each level | Returns correct count at each level |
| Files | - | Check file existence | 4 markdown files created |

**Key Assertions**:
```go
epic.Type == model.TypeEpic
story.HasParent(epic.ID) == true
task.HasParent(story.ID) == true
subtask.HasParent(task.ID) == true
epic.ChildCount == 1
story.ChildCount == 1
task.ChildCount == 1
```

#### Test: `TestE2E_FourLevelHierarchy` → Subtest: "query children returns correct results"
Tests the query engine's ability to find all children:

| Step | Operation | Verification |
|------|-----------|--------------|
| 1 | Create epic | Epic created |
| 2 | Create 3 stories with epic as parent | All stories linked |
| 3 | Query children: `engine.FindChildren(epic.ID)` | Returns list of 3 stories |
| 4 | Verify each story in results | All story IDs found in returned list |

**Key Assertions**:
```go
len(children) == 3
children[0].ID == stories[0].ID
children[1].ID == stories[1].ID
children[2].ID == stories[2].ID
```

#### Test: `TestE2E_FourLevelHierarchy` → Subtest: "deep hierarchy traversal"
Tests complex multi-level hierarchies:

| Step | Operation | Verification |
|------|-----------|--------------|
| 1 | Create epic | Root level |
| 2 | Create 2 stories under epic | Level 2 |
| 3 | Create 2 tasks under story1 | Level 3 (4 total under story1) |
| 4 | Create 2 tasks under story2 | Level 3 (4 total under story2) |
| 5 | Verify epic.ChildCount == 2 | Stories are direct children |
| 6 | Verify story1.ChildCount == 2 | Tasks are story children |
| 7 | Verify story2.ChildCount == 2 | Tasks are story children |

**Key Assertions**:
```go
epic.ChildCount == 2      // 2 stories
story1.ChildCount == 2    // 2 tasks each
story2.ChildCount == 2
```

#### Test: `TestE2E_EpicDeletionOrphaning` → Subtest: "epic deletion leaves orphaned children"
Documents current behavior of task deletion (orphaning issue):

| Step | Operation | Expected Behavior | Status |
|------|-----------|------------------|--------|
| 1 | Create epic with 3 children | Hierarchy created | ✓ |
| 2 | Delete epic | Epic removed from storage | ✓ |
| 3 | Load children | Children still exist in storage | ⚠️ ORPHANED |
| 4 | Check parent relationships | Children still reference deleted parent | ⚠️ STALE REFS |
| 5 | Query children of deleted epic | Query still returns orphaned children | ⚠️ BUG |

**Key Assertions**:
```go
// Before deletion
epic.ChildCount == 3

// After deletion
store.LoadTask(epic.ID) → ERROR (epic deleted)
store.LoadTask(story1.ID) → OK (orphaned!)
story1.relationships[0].TaskID == epic.ID  // Stale reference!
engine.FindChildren(epic.ID) → [3 orphans] // BUG!
```

**Known Issues**:
- Orphaned children remain in storage
- Stale parent references persist
- Query engine doesn't verify parent exists
- **TODO**: Implement cascade delete or orphan prevention

#### Test: `TestE2E_EpicDeletionOrphaning` → Subtest: "blocks relationship survives task deletion"
Tests relationship cleanup on task deletion:

| Step | Operation | Expected Behavior | Status |
|------|-----------|------------------|--------|
| 1 | Create taskA and taskB | Both tasks exist | ✓ |
| 2 | Add blocks relationship: A → B | Relationship stored | ✓ |
| 3 | Delete taskB | taskB removed | ✓ |
| 4 | Reload taskA | taskA still exists | ✓ |
| 5 | Check blocks relationship | Still references deleted taskB | ⚠️ STALE REF |

**Key Assertions**:
```go
taskA.Relationships[0].Type == model.RelBlocks
taskA.Relationships[0].TaskID == taskB.ID  // After taskB deleted!
// Expected: relationship should be cleaned up or marked stale
```

**Known Issues**:
- Stale relationship references not cleaned up on deletion
- **TODO**: Implement relationship cleanup on task deletion

---

### 4. Workflow Tests (`workflow_test.go`)

**Purpose**: Test realistic feature development workflows.

#### Test: `TestE2E_ResearchPlanningImplementationWorkflow`
Tests a complete 15-step feature development cycle:

| Step | Operation | State | Details |
|------|-----------|-------|---------|
| 1 | Create epic | planning | Title: "User Authentication System", Tags: feature, auth |
| 2 | Create research task | todo | Type: research, Parent: epic |
| 3 | Transition research | in-progress | Status: todo → in-progress |
| 4 | Add tags to research | - | Add "security", "investigation" tags |
| 5 | Complete research | done | Status: in-progress → done |
| 6 | Create plan task | todo | Type: plan, Parent: epic |
| 7 | Complete plan | done | Status: todo → done (no in-progress) |
| 8 | Create story | todo | Parent: epic, ready for implementation |
| 9 | Activate epic | active | Status: planning → active (implementation starts) |
| 10 | Start story | in-progress | Status: todo → in-progress |
| 11 | Create blocking task | todo | Title: "Setup user database schema", Parent: story |
| 12 | Create blocked task | blocked | Title: "Implement OAuth endpoints", blocked by task #11 |
| 13 | Complete blocking task | done | Database schema ready |
| 14 | Unblock and complete OAuth | done | Status: blocked → in-progress → done |
| 15 | Complete story and epic | done | All tasks finished, workflow complete |

**Verification at Each Step**:
```go
Step 1:  Assert(epic, env).HasStatus("planning").HasTags("feature", "auth")
Step 3:  Assert(research, env).HasStatus("in-progress")
Step 5:  Assert(research, env).HasStatus("done")
Step 7:  Assert(plan, env).HasStatus("done")
Step 9:  Assert(epic, env).HasStatus("active")
Step 12: Assert(blockingTask, env).HasBlocksRelationship(oauthTask.ID)
Step 15: assert.Equal(t, 6, doneCount)  // All tasks done
```

**Final Verification**:
- Total tasks created: 6 (1 epic + 1 research + 1 plan + 1 story + 2 tasks)
- All tasks status: done
- No orphaned or incomplete tasks

#### Test: `TestE2E_ParallelStoryWorkflow`
Tests multiple stories being worked on simultaneously:

| Step | Operation | Details |
|------|-----------|---------|
| 1 | Create epic | Root for 3 stories |
| 2 | Bulk create 3 stories | All under same epic |
| 3 | Verify all have parent | All stories linked to epic |
| 4 | Transition story 1 | todo → done |
| 5 | Transition story 2 | todo → in-progress |
| 6 | Transition story 3 | stays todo |
| 7 | Verify states | Each story in expected state |
| 8 | Verify epic still has 3 children | Parent count unchanged despite different states |

**Key Assertions**:
```go
// After setup
Assert(stories[0], env).HasParent(epic.ID)
Assert(stories[1], env).HasParent(epic.ID)
Assert(stories[2], env).HasParent(epic.ID)

// After transitions
Assert(stories[0], env).HasStatus("done")
Assert(stories[1], env).HasStatus("in-progress")
Assert(stories[2], env).HasStatus("todo")

// Parent unchanged
Assert(epic, env).HasChildCount(3)
```

---

## Helper Functions and Utilities

### Test Setup (`helpers.go`)

**`SetupE2EEnvironment(t)`**
- Creates temp directory via `t.TempDir()`
- Initializes markdown-fs storage backend
- Creates query engine
- Returns `E2ETestContext` with all components
- Auto-cleanup via defer

**`SetupMemoryEnvironment(t)`**
- Same as above but uses in-memory storage
- No temp directory (faster for unit-style tests)
- Useful when filesystem validation not needed

### Assertions (`assertions.go`)

**`Assert(t, task, env) *AssertTask`**
- Rich assertion failures with formatted output
- Chainable API for multiple checks
- Returns full task JSON on failure

**Available Assertions**:
- `.HasStatus(expected string)` - Check task status
- `.HasTitle(expected string)` - Check task title
- `.HasType(expected string)` - Check task type
- `.HasDescription(expected string)` - Check description
- `.HasParent(parentID int)` - Verify parent relationship
- `.HasChildCount(expected int)` - Query engine verification
- `.HasTag(tag string)` - Single tag check
- `.HasTags(tags ...string)` - Multiple tags check
- `.HasBlocksRelationship(taskID int)` - Verify blocks relationship
- `.HasRelatedToRelationship(taskID int)` - Verify relates-to relationship
- `.FileExists()` - Verify markdown file on disk

### Test Fixtures (`testutil` package)

**Task Creation**:
- `CreateEpic(t, store, options...)` - Create epic task
- `CreateStory(t, store, options...)` - Create story task
- `CreateTask(t, store, options...)` - Create generic task
- `BulkCreateTasks(t, store, count, type, options...)` - Create multiple tasks

**Options** (builder pattern):
- `WithTitle(string)` - Set task title
- `WithStatus(string)` - Set task status
- `WithType(string)` - Set task type
- `WithDescription(string)` - Set description
- `WithTags(tags...)` - Add tags
- `WithParent(int)` - Link to parent
- `WithBlocks(int)` - Add blocks relationship
- `WithRelatedTo(int)` - Add relates-to relationship

**State Transitions**:
- `TransitionTaskState(t, store, taskID, newStatus)` - Change task status

---

## Test Environment Architecture

```
SetupE2EEnvironment
├── Create temp directory (t.TempDir())
├── Initialize Storage
│   └── markdown-fs backend pointing to temp dir
├── Initialize Query Engine
│   └── Queries against storage
├── Create Context
│   └── context.Background()
└── Return E2ETestContext
    ├── Store (BaseStorage interface)
    ├── Engine (QueryEngine)
    ├── TmpDir (temp filesystem path)
    ├── Ctx (context)
    └── Cleanup (defer function)
```

---

## Running Tests with Coverage

```bash
# Generate coverage report
go test -v -coverprofile=coverage.out ./test/e2e
go tool cover -html=coverage.out

# Show coverage per function
go test -v -coverprofile=coverage.out ./test/e2e
go tool cover -func=coverage.out
```

---

## Troubleshooting

### Binary Not Found
```bash
# Ensure binary is built
mise run build

# Check binary location
ls -la bin/opentask

# Add to PATH or specify full path in tests
export PATH=$PWD/bin:$PATH
```

### Tests Skip CLI Tests
Tests automatically skip if `opentask` binary not found. Check:
```bash
which opentask
echo $PATH
ls bin/opentask
```

### Temp Directory Issues
- Tests use `t.TempDir()` - auto-cleanup on test completion
- For debugging, `t.Logf("tmpDir: %s", env.TmpDir)` to keep directory
- Check permissions: `chmod +rw` on temp directories

### Storage Not Persisting
- Verify `storage.SaveTask()` returns no error
- Check temp directory exists: `os.Stat(env.TmpDir)`
- Verify markdown-fs backend initialized: check `storage.go`
- Call `env.Cleanup()` to flush writes

