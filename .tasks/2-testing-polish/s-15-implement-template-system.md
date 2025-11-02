---
id: 15
title: Implement template system
type: story
status: todo
tags:
    - feature
relationships:
    - type: parent
      taskID: 1
createdAt: "2025-11-02T08:41:05Z"
updatedAt: "2025-11-02T20:35:00Z"
---

## Objective
Implement a template system that allows users to define task templates for different task types. This enables consistent task structure and reduces repetitive work when creating similar tasks.

## Current State
- Config system supports optional template paths for each task type
- No template loading or application logic exists
- No way to create tasks from templates
- Default templates are embedded in code (hardcoded or none)

## Template System Design

### Template File Format
Templates are markdown files with YAML frontmatter (same format as tasks):

```markdown
---
# These fields are used as defaults when creating tasks from this template
title: "Template: {{ .type | title }} Task"  # Optional placeholder for task type
type: story
tags:
  - feature
  - needs-review
status: null  # null means use workflow.initial status
relationships: []
---

## Summary
Brief description of work

## Requirements
- Requirement 1
- Requirement 2

## Acceptance Criteria
- [ ] Criterion 1
- [ ] Criterion 2

## Notes
Optional notes for task creator
```

### Template Placeholders (Optional Enhancement)
Simple template variable support (can be added later if needed):
- `{{ .type }}` - task type being created
- `{{ .title }}` - task title from CLI
- `{{ .date }}` - current date
- `{{ .user }}` - current user

For MVP, templates can be static with no placeholders.

## Template Resolution

### Template Locations (in order of precedence)
1. Explicit path in config file: `templates.story = "./custom/story.md"`
2. Project local: `.tasks/templates/story.md`
3. User XDG data: `${XDG_DATA_HOME}/opentasks/templates/story.md`
4. System default: Built-in templates in code

### Task Type to Template Mapping
Maps task types to template filenames:
```
epic     → epic.md
plan     → plan.md
research → research.md
story    → story.md
decision → decision.md
task     → task.md
```

## CLI Integration

### New Flag: `task new --template`
Create task from template:

```bash
task new "Implement feature X" --type story
# Uses: .tasks/templates/story.md or default

task new "Implement feature X" --type story --template ./templates/custom-story.md
# Uses custom template
```

### Behavior
1. Parse template from configured/provided path
2. Extract frontmatter fields (tags, status, relationships)
3. Use template markdown body as initial description
4. CLI flags override template fields
5. Create task with merged content

**Example:**
Template has: `tags: [feature, needs-review]`
CLI: `task new "Title" --tag urgent`
Result: `tags: [feature, needs-review, urgent]`

## Implementation Requirements

### Config System Changes
- Already supports `templates.{type}` in config.toml
- Verify path resolution works (relative, absolute, XDG)
- Add validation that template file exists

### Storage Changes
- No storage changes needed
- Template loading happens in CLI layer

### CLI Changes (cmd/task.go)
- Modify `taskNewCmd` to accept `--template` flag
- Add template loading logic
- Merge template content with CLI flags
- Handle template not found errors

### Model/Template Loader (new)
Create `internal/config/templates.go`:
- `LoadTemplate(ctx, taskType, configPath) (*TemplateContent, error)`
- `TemplateContent` struct with frontmatter fields and body
- Template resolution logic (check all locations)
- YAML parsing with validation

### Tests to Write
- Unit tests for template loading from different locations
- Template content merging (CLI flags override template fields)
- Task creation from template integration test
- Error handling: missing template, invalid YAML, etc.
- Path resolution with environment variables and relative paths

## Acceptance Criteria
- [ ] Template files can be loaded from all configured locations
- [ ] `task new` accepts `--template` flag
- [ ] Task created from template includes template content
- [ ] CLI flags override template frontmatter fields
- [ ] Template tags/status merge correctly with CLI arguments
- [ ] Default templates provided for all task types
- [ ] Config supports custom template paths per type
- [ ] Helpful error message if template not found
- [ ] Help text shows template usage: `task new --help`
- [ ] Unit and integration tests pass
