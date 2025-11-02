# OpenTasks Testing Guide

## Overview

This document describes the testing strategy and patterns used in the OpenTasks project.

## Test Structure

### Test Organization

Tests are organized by package, with test files collocated with the code they test:

```
internal/
  ├── model/
  │   ├── task.go
  │   ├── task_test.go           # Tests for task types and validation
  │   ├── relationship.go
  │   └── relationship_test.go   # Tests for relationship types
  ├── config/
  │   ├── config.go
  │   └── config_test.go         # Tests for configuration loading and defaults
  ├── storage/
  │   ├── memory.go
  │   ├── memory_test.go         # Tests for MemoryStorage backend
  │   ├── interface.go
  │   └── ...
  ├── query/
  │   ├── engine.go
  │   ├── query_test.go          # Tests for QueryEngine
  │   ├── filters.go
  │   ├── filters_test.go        # Tests for filter functions
  │   └── ...
  └── testutil/
      └── fixtures.go            # Shared test fixtures and helpers
```

### Running Tests

Run all tests:
```bash
go test -v ./...
```

Run tests for a specific package:
```bash
go test -v ./internal/model
go test -v ./internal/storage
```

Run tests with coverage:
```bash
go test -cover ./...
```

Generate coverage report:
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Test Patterns

### 1. Fixtures and Helpers (`internal/testutil/fixtures.go`)

Reusable test fixtures reduce boilerplate and ensure consistency:

```go
// Create a basic test task
task := NewTestTask(1, "Test Task")

// Create task with specific type
story := NewTestTaskWithType(2, "Story", model.TypeStory)

// Create task with tags
tagged := NewTestTaskWithTags(3, "Tagged", []string{"feature", "core"})

// Get a full sample dataset
tasks := SampleTasks()
```

**Benefits:**
- Consistent test data across all tests
- Easier to maintain test data
- Clear intent in test code

### 2. Table-Driven Tests

Many validation tests use table-driven patterns for comprehensive coverage:

```go
func TestIsValidType(t *testing.T) {
    tests := []struct {
        name     string
        taskType string
        want     bool
    }{
        {"epic type", TypeEpic, true},
        {"invalid type", "invalid", false},
        // ... more cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := IsValidType(tt.taskType); got != tt.want {
                t.Errorf("got %v, want %v", got, tt.want)
            }
        })
    }
}
```

**Benefits:**
- Easy to add more test cases
- Clear test case organization
- Self-documenting edge cases

### 3. Mock Storage Backend

Tests use `MemoryStorage` as a lightweight mock:

```go
func setupQueryEngine(ctx context.Context, t *testing.T) (*QueryEngine, storage.BaseStorage) {
    store := storage.NewMemoryStorage()
    qe := NewQueryEngine(store)
    return qe, store
}
```

**Benefits:**
- No filesystem required
- Fast test execution
- Deterministic behavior

### 4. Context Passing

All tests that interact with storage pass `context.Background()`:

```go
ctx := context.Background()
task, err := store.LoadTask(ctx, id)
```

This ensures tests are ready for cancellation/timeout handling in the future.

## Coverage Analysis

### Current Coverage

| Package | Coverage | Status |
|---------|----------|--------|
| `model` | 100% | ✅ Complete |
| `query` | 96.4% | ✅ Nearly complete |
| `storage` | 22.1% | ⚠️ MarkdownFileStorage not tested |
| `config` | 18.9% | ⚠️ Config file I/O paths not fully tested |
| `cmd` | 0% | ⚠️ CLI tests not yet implemented |

### Low Coverage Areas

**MarkdownFileStorage** - Requires filesystem interaction:
- Could be tested with `t.TempDir()` for isolated filesystem tests
- Should test file parsing, markdown format, error handling

**Configuration** - Some code paths not covered:
- File parsing edge cases
- Path resolution with environment variables
- Configuration hierarchy merging

**CLI Commands** - No tests yet:
- Integration tests for each command
- Error handling and user feedback
- Output formatting

## Testing Best Practices

### 1. Test Independence
- Each test is completely independent
- No shared state between tests
- Tests can run in any order

### 2. Clear Test Names
Test names describe what is being tested and what's expected:
```go
func TestMemoryStorageSaveAndLoad(t *testing.T) // ✅ Clear
func TestSave(t *testing.T)                     // ❌ Ambiguous
```

### 3. Isolation
Tests create their own fixtures:
```go
func TestFindByID(t *testing.T) {
    ctx := context.Background()
    qe, store := setupQueryEngine(ctx, t)  // Create fresh store
    defer store.Close()
    // ... test code
}
```

### 4. Error Checking
Consistent pattern for error checking:
```go
if err := store.SaveTask(ctx, task); err != nil {
    t.Fatalf("SaveTask() error = %v, want nil", err)
}
```

### 5. Assertions Over Mocks
- Use assertion-based tests rather than mocks
- Mocks add complexity without benefit in this codebase
- Real in-memory storage is simple and fast enough

## Adding New Tests

### Checklist for New Features

1. **Write tests first** (TDD approach)
   - Tests define the contract
   - Guides implementation

2. **Use existing fixtures**
   - Leverage `testutil` helpers
   - Keep tests consistent

3. **Test edge cases**
   - Empty inputs
   - Boundary conditions
   - Error conditions

4. **Update coverage**
   - Aim for >90% coverage
   - Document low-coverage areas

### Example: Testing New Query Filter

```go
func TestNewFilter(t *testing.T) {
    tests := []struct {
        name string
        task *model.Task
        want bool
    }{
        {"matches", createTestTask(...), true},
        {"no match", createTestTask(...), false},
        {"edge case", createTestTask(...), false},
    }

    filter := WithNewFilter("value")
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := filter(tt.task); got != tt.want {
                t.Errorf("got %v, want %v", got, tt.want)
            }
        })
    }
}
```

## Future Improvements

1. **Increase MarkdownFileStorage coverage**
   - Add filesystem-based tests
   - Test file parsing edge cases

2. **Add integration tests**
   - Test component interaction
   - Test CLI commands end-to-end

3. **Add stress tests**
   - Test with many tasks (1000+)
   - Test with deep nesting
   - Memory usage validation

4. **Add benchmarks**
   - Query performance benchmarks
   - Storage operation benchmarks
   - Help identify performance bottlenecks

## Test Execution in CI/CD

All tests should pass in CI/CD:
```bash
go test -v ./...
go test -race ./...  # Detect race conditions
```

Expected output: All tests passing, no race conditions detected.
