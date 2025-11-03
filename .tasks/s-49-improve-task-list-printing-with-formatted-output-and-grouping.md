---
id: 49
title: Improve task list printing with formatted output and grouping
type: story
status: todo
tags: [cli, ux, formatting, output]
relationships: []
createdAt: "2025-11-03T10:00:00Z"
updatedAt: "2025-11-03T10:00:00Z"
---

# Story: Improve Task List Printing with Formatted Output and Grouping

## Problem

Currently, task list output is basic fixed-width plain text:

```
ID    Type       Title                          Status          Created    
---------------------------------------------------------------------------
1     epic       My Epic Idea                   backlog         2025-11-03
2     story      Implement feature              todo            2025-11-02
3     task       Quick fix                      in-progress     2025-11-01
```

Issues:
- No formatting options (JSON, pretty-printed, etc.)
- No grouping capability (by status, type, parent, etc.)
- No color or styling
- Limited readability for large task lists
- Can't use output in other tools (JSON)
- No way to customize columns or sorting

## Vision

### Output Formats

Support multiple output formats:

1. **Pretty (default)** - Beautiful markdown table rendered with syntax highlighting
2. **JSON** - Structured output for tooling and automation

### Pretty Printing Features

1. **Markdown Table** - Generated as markdown, rendered with glamour
2. **Grouping** - Group results by any field or relationship
3. **Colors** - Status colors, type icons, visual hierarchy
4. **Better formatting** - Smart column widths, word wrapping

### Examples

#### Default Pretty Output (grouped by status)

```
# Tasks

## Backlog (2)

| ID | Type | Title | Created |
|---|---|---|---|
| 1 | epic | My Epic Idea | 2025-11-03 |
| 5 | story | Review proposal | 2025-11-02 |

## Todo (3)

| ID | Type | Title | Created |
|---|---|---|---|
| 2 | story | Implement feature | 2025-11-02 |
| 3 | task | Quick fix | 2025-11-01 |
| 7 | task | Update docs | 2025-11-01 |

## In Progress (1)

| ID | Type | Title | Created |
|---|---|---|---|
| 4 | epic | Large initiative | 2025-11-03 |

## Done (2)

| ID | Type | Title | Created |
|---|---|---|---|
| 6 | task | Deploy v1.0 | 2025-10-30 |
| 8 | decision | Use Go for CLI | 2025-10-28 |
```

#### JSON Output

```json
{
  "total": 8,
  "tasks": [
    {
      "id": 1,
      "title": "My Epic Idea",
      "type": "epic",
      "status": "backlog",
      "created": "2025-11-03T00:00:00Z",
      "updated": "2025-11-03T00:00:00Z",
      "tags": ["high-priority"],
      "relationships": []
    },
    ...
  ]
}
```

#### Grouped by Type

```
# Tasks

## Epic (2)

| ID | Title | Status | Created |
|---|---|---|---|
| 1 | My Epic Idea | backlog | 2025-11-03 |
| 4 | Large initiative | in-progress | 2025-11-03 |

## Story (2)

| ID | Title | Status | Created |
|---|---|---|---|
| 2 | Implement feature | todo | 2025-11-02 |
| 5 | Review proposal | backlog | 2025-11-02 |

## Task (3)

| ID | Title | Status | Created |
|---|---|---|---|
| 3 | Quick fix | in-progress | 2025-11-01 |
| 6 | Deploy v1.0 | done | 2025-10-30 |
| 7 | Update docs | todo | 2025-11-01 |

## Decision (1)

| ID | Title | Status | Created |
|---|---|---|---|
| 8 | Use Go for CLI | done | 2025-10-28 |
```

#### Grouped by Parent (shows hierarchy)

```
# Tasks

## Epic: My Epic Idea (1)

| ID | Type | Title | Status |
|---|---|---|---|
| 2 | story | Implement feature | todo |

## Epic: Large Initiative (2)

| ID | Type | Title | Status |
|---|---|---|---|
| 3 | task | Quick fix | in-progress |
| 7 | task | Update docs | todo |

## Ungrouped (4)

| ID | Type | Title | Status |
|---|---|---|---|
| 1 | epic | My Epic Idea | backlog |
| 4 | epic | Large initiative | in-progress |
| 5 | story | Review proposal | backlog |
| 8 | decision | Use Go for CLI | done |
```

## Implementation Plan

### Phase 1: Output Format Support

1. **Add `--format` flag** to `taskListCmd`:
   ```go
   taskListCmd.Flags().StringP("format", "f", "pretty", "Output format (pretty, json)")
   ```

2. **Add `--group-by` flag** to `taskListCmd`:
   ```go
   taskListCmd.Flags().StringP("group-by", "g", "status", "Group by field (status, type, parent, none)")
   ```

3. **Update list command logic**:
   - Check `--format` flag
   - Route to appropriate formatter (JSON or pretty)
   - Pass grouping preference to formatter

### Phase 2: JSON Output

1. **Create formatter package** `internal/format/` (if doesn't exist):
   - `formatter.go` - Base interfaces
   - `json_formatter.go` - JSON output

2. **Implement JSON formatter**:
   ```go
   type JSONFormatter struct{}
   
   func (f *JSONFormatter) Format(tasks []*model.Task, groupBy string) (string, error) {
       output := map[string]interface{}{
           "total": len(tasks),
           "tasks": tasks,
       }
       return json.MarshalIndent(output, "", "  ")
   }
   ```

3. **Test JSON output**:
   - Valid JSON structure
   - Includes all task fields
   - Pretty-printed with indentation

### Phase 3: Pretty Printing Template

1. **Create template file** `internal/cmd/list.go.tmpl`:
   ```go
   // Template for rendering markdown table
   # {{ .Title }}
   
   {{ range $groupName, $tasks := .Groups }}
   ## {{ $groupName }} ({{ len $tasks }})
   
   | ID | Type | Title | Status | Created |
   |---|---|---|---|---|
   {{ range $tasks }}
   | {{ .ID }} | {{ .Type }} | {{ .Title }} | {{ .Status }} | {{ .CreatedAt.Format "2006-01-02" }} |
   {{ end }}
   
   {{ end }}
   ```

2. **Create pretty formatter** `internal/format/pretty_formatter.go`:
   ```go
   type PrettyFormatter struct {
       template *template.Template
   }
   
   func (f *PrettyFormatter) Format(tasks []*model.Task, groupBy string) (string, error) {
       grouped := f.groupTasks(tasks, groupBy)
       data := map[string]interface{}{
           "Title":  "Tasks",
           "Groups": grouped,
       }
       var buf bytes.Buffer
       if err := f.template.Execute(&buf, data); err != nil {
           return "", err
       }
       return buf.String(), nil
   }
   
   func (f *PrettyFormatter) groupTasks(tasks []*model.Task, groupBy string) map[string][]*model.Task {
       // Group by status, type, parent, or none
   }
   ```

3. **Render with glamour**:
   - Take markdown output
   - Pass through glamour renderer (already used in config view)
   - Output beautifully formatted terminal markdown

### Phase 4: Grouping Logic

1. **Implement grouping strategies**:
   ```go
   func GroupByStatus(tasks []*model.Task) map[string][]*model.Task {
       groups := make(map[string][]*model.Task)
       for _, task := range tasks {
           groups[task.Status] = append(groups[task.Status], task)
       }
       return groups
   }
   
   func GroupByType(tasks []*model.Task) map[string][]*model.Task {
       // Similar implementation
   }
   
   func GroupByParent(tasks []*model.Task) map[string][]*model.Task {
       // Group by parent ID, showing hierarchy
   }
   
   func GroupByNone(tasks []*model.Task) map[string][]*model.Task {
       // Single group with all tasks
   }
   ```

2. **Order groups intelligently**:
   - Status: backlog → todo → in-progress → reviewing → done → archived
   - Type: epic → plan → story → task → decision → research
   - Parent: by parent ID, ungrouped last
   - None: original order

3. **Show group counts** in headers

### Phase 5: Update Task List Command

1. **Refactor `taskListCmd`** in `cmd/task.go`:
   - Remove inline table printing
   - Create `ListOutputter` interface
   - Support pretty and JSON formatters
   - Apply grouping

2. **Example flow**:
   ```go
   // Get tasks
   tasks, err := Engine.ListTasks(ctx, filters...)
   
   // Choose formatter based on flag
   var outputter format.Outputter
   if formatFlag == "json" {
       outputter = format.NewJSONOutputter()
   } else {
       outputter = format.NewPrettyOutputter()
   }
   
   // Format and output
   output, err := outputter.Format(tasks, groupByFlag)
   fmt.Println(output)
   ```

### Phase 6: Tests

1. **Unit tests** for formatters:
   - JSON formatter produces valid JSON
   - Pretty formatter produces valid markdown
   - Grouping works correctly for each strategy

2. **Integration tests** for command:
   - `task list --format json` outputs JSON
   - `task list --format pretty` outputs markdown
   - `task list --group-by status` groups correctly
   - `task list --group-by type` groups correctly
   - `task list --group-by parent` groups correctly
   - Combination flags work together

3. **Output validation**:
   - JSON is valid and parseable
   - Markdown renders without errors
   - All tasks appear in output
   - Grouping counts are correct

## File Structure

```
internal/
└── format/
    ├── formatter.go           # Base interfaces
    ├── json_formatter.go      # JSON output
    ├── pretty_formatter.go    # Markdown/glamour output
    ├── grouper.go             # Grouping logic
    └── formatter_test.go      # Tests

internal/cmd/                  # Or move to internal/format/
└── list.go.tmpl             # Markdown template for list output

cmd/
└── task.go                   # Update taskListCmd
```

## CLI Usage Examples

```bash
# Default: pretty output grouped by status
opentask task list

# Pretty output grouped by type
opentask task list --group-by type

# Pretty output, no grouping
opentask task list --group-by none

# JSON output (ignores grouping, for tooling)
opentask task list --format json

# JSON with filtering
opentask task list --status backlog --format json | jq '.tasks[] | .title'

# Pretty with filters
opentask task list --status in-progress --group-by type

# Show only todos grouped by type
opentask task list --status todo --group-by type
```

## Flags Reference

- `--format <pretty|json>` - Output format (default: pretty)
- `--group-by <status|type|parent|none>` - Group results (default: status)
- `--status <status>` - Filter by status (existing)
- `--type <type>` - Filter by type (existing)
- `--tag <tag>` - Filter by tags (existing)

## Success Criteria

- [ ] `--format` flag supports json and pretty
- [ ] `--group-by` flag supports status, type, parent, none
- [ ] JSON output is valid and complete
- [ ] Pretty output uses markdown tables
- [ ] Pretty output rendered with glamour
- [ ] Grouping works for all strategies
- [ ] Group counts displayed in headers
- [ ] Groups ordered intelligently
- [ ] All filters work with formatting
- [ ] All tests pass
- [ ] Documentation updated

## Future Enhancements

- [ ] CSV output format
- [ ] Custom column selection (`--columns`)
- [ ] Sorting options (`--sort`)
- [ ] Custom grouping order
- [ ] Export to HTML
- [ ] Tree view for parent relationships

## Related Tasks

- Story 46: Config schema (may include output format preferences)
- Task 44: CLI editing (consistency in output)

## References

- Current implementation: `cmd/task.go` lines 78-136
- Glamour usage: `cmd/config.go` (already using for markdown rendering)
- Template example: `internal/config/view.go.tmpl`
