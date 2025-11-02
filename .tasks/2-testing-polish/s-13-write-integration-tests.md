---
id: 13
title: Write integration tests
type: story
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
Test interactions between multiple components working together to ensure the system behaves correctly end-to-end. Integration tests verify that queries, storage, config, and models work together as intended.

## Scope: Multi-Component Workflows

### 1. **Storage + Query Integration**
Test that QueryEngine correctly filters tasks from different storage backends:

- Load tasks via MarkdownFileStorage
- Apply complex filters (type, status, tags, relationships)
- Verify filtered results match expected tasks
- Test with tasks in epic parent directories
- Test relationship-based queries (find all children of an epic)
- Test with large task sets (100+ tasks)

### 2. **Config + Storage Integration**
Test that configuration properly directs storage behavior:

- Load config from different locations (local .tasks/, XDG, explicit path)
- Verify storage uses correct base path from config
- Test custom workflow statuses affect task validation
- Verify template paths resolve correctly from config

### 3. **Task Lifecycle Workflows**
Test complete task workflows from creation through deletion:

**Create → Read → Update → Delete:**
- Create a task with tags, relationships, type
- Load it back and verify all fields preserved
- Update status, title, tags
- Load again and verify updates persisted
- Delete task
- Verify it can't be loaded again
- Verify relationships to deleted task are handled gracefully

**Epic Hierarchy Workflow:**
- Create an epic (parent task)
- Create multiple stories (child tasks with parent relationships)
- List tasks and verify children appear under parent
- Update epic title and verify children still link correctly
- Delete epic and verify behavior (orphaned tasks remain, relationships broken)

### 4. **Query + Relationship Integration**
Test complex relationship queries:

- Find all children of an epic using query filters
- Find all tasks that block a given task
- Find related tasks across multiple relationship types
- Verify circular relationships are detected (if applicable)
- Test relationship queries with tasks spread across multiple epics

### 5. **File Format + Parsing Integration**
Test complete markdown file round-trips:

- Create task in memory with all fields
- Save to markdown file
- Parse file and verify format (YAML frontmatter, markdown body)
- Load task back and verify all fields match original
- Modify file (YAML and markdown), load again, verify changes persist
- Test with special characters, unicode, multiline descriptions

### 6. **Error Handling + Recovery**
Test system behavior when things go wrong:

- Load task with corrupted YAML frontmatter (fallback behavior)
- Save task when directory doesn't exist (auto-create)
- Handle missing config file (use defaults)
- Delete non-existent task (graceful error)
- Load task with invalid type code in filename (still load from frontmatter)

## Test Pattern Requirements
- Create realistic workflows that users might perform
- Use `t.TempDir()` for filesystem-based tests
- Test with multiple storage backends (MemoryStorage, MarkdownFileStorage)
- Test with multiple config variations
- Verify state is consistent across operations

## Acceptance Criteria
- [ ] All major workflows have integration tests
- [ ] Storage + Query integration verified
- [ ] Config + Storage integration verified
- [ ] File format round-trips work correctly
- [ ] Error cases handled gracefully
- [ ] All tests pass: `go test ./...`
- [ ] No race conditions: `go test -race ./...`
- [ ] Tests are independent and can run in any order
