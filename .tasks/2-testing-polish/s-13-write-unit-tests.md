---
id: 13
title: Write unit tests
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
Ensure all packages have comprehensive unit tests that verify correctness of business logic and validate against requirements. Tests should be isolated, deterministic, and follow established patterns from `docs/TESTING.md`.

## Scope: Package Coverage Requirements

### 1. `internal/model` - **Status: ✅ Complete (100%)**
Already has excellent coverage:
- Task type validation (epic, plan, research, story, decision, task)
- Type code mapping (e↔epic, p↔plan, etc.)
- Relationship type validation (parent, blocks, relates-to)
- Tag and metadata handling

### 2. `internal/query` - **Status: ⚠️ 96.4% (needs final edge cases)**
Existing tests cover filter functions and query engine. Gaps to address:
- Complex multi-filter combinations (AND/OR logic)
- Edge cases: empty task lists, circular relationships in filters
- Query performance with large datasets (1000+ tasks)
- Boundary conditions in date/time filtering

### 3. `internal/storage` - **Status: ⚠️ 22.1% (MemoryStorage complete, MarkdownFileStorage needs tests)**

#### MemoryStorage - **Status: ✅ Complete**
- SaveTask, LoadTask, DeleteTask operations
- Filters applied during list operations
- Relationship handling

#### MarkdownFileStorage - **Status: ❌ No tests**
Tests must verify:
- File path generation for tasks with/without epic parents
- YAML frontmatter parsing (valid and invalid formats)
- Markdown body content preservation
- NextID generation with sequential numbering
- Directory structure creation and cleanup
- File encoding (UTF-8) handling
- Filename slugification (special chars, unicode, spaces)

Use `t.TempDir()` for isolated filesystem tests. Test with realistic markdown files including YAML syntax edge cases.

### 4. `internal/config` - **Status: ⚠️ 18.9% (partial coverage)**
Existing tests cover default loading. Gaps:
- TOML parsing with various syntaxes and nesting
- Path resolution with environment variables (OPENTASKS_PROJECT_PATH)
- Config file location discovery (local .tasks/, XDG paths, explicit paths)
- Configuration hierarchy merging (defaults → project config → env overrides)
- Invalid config file handling (malformed TOML, missing required fields)
- Workflow transition validation against config
- Template path resolution (relative, absolute, XDG paths)

### 5. `cmd` - **Status: ❌ 0% (no tests)**
Unit tests for CLI command handlers. Test each command's core logic:
- **task new**: ID generation, type validation, parent linking, tag handling
- **task list**: Filter application, output formatting
- **task view**: Task loading, relationship display
- **task delete**: Task removal, relationship cleanup
- **project init**: Config file creation, directory setup
- **config view**: Config loading and display

Tests should use in-memory storage and mock filesystem as needed.

## Test Pattern Requirements
- Use table-driven tests for validation (see `TESTING.md`)
- Use fixture helpers from `internal/testutil` for consistency
- Use `context.Background()` for all storage operations
- Test both success and failure paths
- Include edge cases and boundary conditions
- Write tests before implementation (TDD)

## Acceptance Criteria
- [ ] MarkdownFileStorage has >90% test coverage
- [ ] All config code paths tested (env vars, file I/O, merging)
- [ ] CLI command handlers have unit tests for core logic
- [ ] All tests pass: `go test ./...`
- [ ] No race conditions: `go test -race ./...`
- [ ] Generate coverage report: `go test -coverprofile=coverage.out ./...`
- [ ] Document any low-coverage areas and reasons
