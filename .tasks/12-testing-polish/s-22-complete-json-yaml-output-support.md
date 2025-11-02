---
id: 22
title: Complete JSON/YAML output support
type: story
status: done
tags:
    - feature
    - output
relationships:
    - type: parent
      taskID: 12
    - type: relates-to
      taskID: 16
createdAt: "2025-11-02T21:00:00Z"
updatedAt: "2025-11-02T21:00:00Z"
---

## Objective
Ensure all CLI commands that output task data support JSON and YAML formats consistently.

## Current State
- s-16 (Add JSON/YAML output) designs the feature
- Implementation may be incomplete or inconsistent across commands
- Need to verify all commands that output tasks support `--format` flag

## Scope: Output Format Support

### Commands that need `--format` support
- `task new` - output created task
- `task list` - output task list
- `task show` - output single task
- `task update` - output updated task
- `task delete` - output deleted task info

### Verification Requirements
- Each command supports: `--format text|json|yaml`
- Text format is default
- JSON output is valid and parseable
- YAML output is valid and parseable
- Timestamps are RFC3339 format
- Task relationships included in output
- Tags array serialized correctly

### Error Handling
- Invalid format flag shows helpful error
- Format conversion errors caught and reported
- Empty results handled properly (empty array for JSON/YAML)

## Acceptance Criteria
- [ ] All output commands support `--format` flag
- [ ] JSON output is valid per spec
- [ ] YAML output is valid per spec
- [ ] Text output unchanged (backward compatible)
- [ ] Tests verify each format for each command
- [ ] Help text documents format options
- [ ] No breaking changes to existing behavior
