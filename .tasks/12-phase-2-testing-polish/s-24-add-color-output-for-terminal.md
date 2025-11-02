---
id: 24
title: Add color output for terminal
type: story
status: done
tags:
    - feature
    - ui
relationships:
    - type: parent
      taskID: 12
createdAt: "2025-11-02T21:00:00Z"
updatedAt: "2025-11-02T12:08:08Z"
---

## Objective
Enhance CLI output with color coding for improved readability and visual hierarchy.

## Color Scheme

### Task Type Colors
- Epic: Blue (`#3498db`)
- Plan: Cyan (`#1abc9c`)
- Research: Yellow (`#f39c12`)
- Story: Green (`#27ae60`)
- Decision: Purple (`#8e44ad`)
- Task: Gray (`#7f8c8d`)

### Status Colors
- todo: Gray
- in-progress: Yellow
- done: Green
- archived: Dim gray

### Output Elements
- Task IDs: Bold white
- Titles: Default (inherit from type color context)
- Status: Colored per status
- Tags: Dim colored (lighter variant of color)
- Errors: Red
- Warnings: Orange/Yellow

## Implementation
Use library like `fatih/color` or `charmbracelet/lipgloss` for:
- Color output with terminal detection
- Graceful fallback for non-TTY (pipes, redirects)
- TERM environment variable support
- `--no-color` flag to disable colors

## Examples

### List Output
```
ID   Type     Status       Title
─────────────────────────────────
1    [Epic]   in-progress  Phase 1: Design
5    [Story]  done         Implement config
10   [Task]   todo         Write tests
```

Colors applied to:
- ID numbers: bright white
- Type badges: per type color
- Status: per status color
- Clickable IDs: potentially underlined

### Task View Output
- Title in type color
- Status in status color
- Metadata (created, updated) in dim gray
- Description: normal text
- Relationships: with colored task types

## Configuration
- `--color`: always, never, auto (default)
- `CLICOLOR` env var support (standard)
- Auto-detect TTY and disable for pipes

## Acceptance Criteria
- [ ] Color output on TTY by default
- [ ] `--color` flag controls output
- [ ] Colors disabled in pipes/non-TTY
- [ ] `--no-color` flag works
- [ ] `CLICOLOR` environment variable respected
- [ ] Color scheme matches type/status semantics
- [ ] Tests verify color code generation
- [ ] Backward compatible (doesn't break parsers)
