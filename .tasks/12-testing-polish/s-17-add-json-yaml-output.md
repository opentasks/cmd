---
id: 17
title: Add JSON/YAML output
type: story
status: done
tags:
    - feature
relationships:
    - type: parent
      taskID: 12
createdAt: "2025-11-02T08:41:05Z"
updatedAt: "2025-11-02T22:00:00Z"
---

## Objective
Add structured output formats (JSON and YAML) to task listing and viewing commands. This enables integration with other tools, scripting, and machine parsing.

## Current State
- All task commands output human-readable text format
- No structured output options
- No way to programmatically consume task data from CLI

## Output Format Implementation

### Format Support
Add `--format` flag to task commands:
```bash
task list --format json
task list --format yaml
task list --format text    # Default, current output

task view <id> --format json
task view <id> --format yaml
task view <id> --format text  # Default

task new "Title" --format json  # Output created task as JSON
```

### JSON Schema

#### Single Task
```json
{
  "id": 42,
  "title": "Implement feature",
  "type": "story",
  "status": "in-progress",
  "tags": ["feature", "core"],
  "description": "Full markdown description here",
  "relationships": [
    {
      "type": "parent",
      "taskID": 5
    },
    {
      "type": "blocks",
      "taskID": 43
    }
  ],
  "createdAt": "2025-11-02T10:00:00Z",
  "updatedAt": "2025-11-02T12:00:00Z"
}
```

#### Task List
```json
{
  "tasks": [
    { /* task 1 */ },
    { /* task 2 */ }
  ],
  "count": 2,
  "filters": {
    "status": "in-progress",
    "type": "story"
  }
}
```

### YAML Schema
Same structure as JSON, rendered as YAML:

```yaml
id: 42
title: Implement feature
type: story
status: in-progress
tags:
  - feature
  - core
description: Full markdown description here
relationships:
  - type: parent
    taskID: 5
  - type: blocks
    taskID: 43
createdAt: 2025-11-02T10:00:00Z
updatedAt: 2025-11-02T12:00:00Z
```

## Implementation Requirements

### CLI Changes (cmd/task.go, cmd/root.go)
- Add `--format` global flag or per-command flag
- Modify `taskListCmd` to support JSON/YAML output
- Modify `taskViewCmd` to support JSON/YAML output
- Modify `taskNewCmd` to output created task in requested format
- Modify `taskDeleteCmd` to output deleted task info in requested format

### Output Formatter (new)
Create `cmd/format.go` or similar:
- `FormatTask(task *model.Task, format string) (string, error)`
- `FormatTasks(tasks []*model.Task, format string, filters map[string]interface{}) (string, error)`
- Support: "text", "json", "yaml"
- Proper error handling for unknown formats
- Pretty-printed output (indented) for readability

### Model Changes
- No model changes needed
- Use existing Task struct for marshaling

### Integration with Filters
When listing with `--format json`:
- Include applied filters in output metadata
- Show filtered count vs total count
- Example: `task list --type story --status done --format json`

## Example Usage

```bash
# List all stories as JSON
$ task list --type story --format json
{
  "tasks": [...],
  "count": 5,
  "filters": {
    "type": "story"
  }
}

# View single task as YAML
$ task view 42 --format yaml
id: 42
title: Implement feature
...

# Create task and output as JSON
$ task new "Title" --type story --format json
{
  "id": 99,
  "title": "Title",
  ...
}
```

## Acceptance Criteria
- [ ] `--format` flag works with `task list`, `task view`, `task new`
- [ ] JSON output matches schema and is valid JSON
- [ ] YAML output matches schema and is valid YAML
- [ ] Text format (default) shows human-readable output
- [ ] Pretty-printed output with proper indentation
- [ ] Filter metadata included in list JSON/YAML output
- [ ] Invalid format flag shows helpful error
- [ ] All timestamps in RFC3339 format (ISO 8601)
- [ ] Help text documents format options: `task list --help`
- [ ] Unit tests verify JSON/YAML parsing
- [ ] Integration tests verify end-to-end output
