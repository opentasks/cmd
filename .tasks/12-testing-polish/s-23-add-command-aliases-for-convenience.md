---
id: 23
title: Add command aliases for convenience
type: story
status: done
tags:
    - feature
    - cli
relationships:
    - type: parent
      taskID: 12
createdAt: "2025-11-02T21:00:00Z"
updatedAt: "2025-11-02T21:00:00Z"
---

## Objective
Add convenient command aliases to reduce typing and improve user experience.

## Design Aliases

### Task Command Aliases
- `task create` → `task new`
- `task ls` → `task list`
- `task view` → `task show`
- `task edit` → `task update`
- `task rm` → `task delete`

### Project Command Aliases
- `project create` → `project new`
- `project ls` → `project list`

### Implementation Approach
Using Cobra's alias mechanism - each alias is a separate command that routes to the primary command.

**Example:**
```go
var taskCreateCmd = &cobra.Command{
    Use: "create",
    Aliases: []string{"new"},  // or reverse: main is "create", alias is "new"
    RunE: taskNewCmd.RunE,
}
```

## Testing Requirements
- Each alias invokes the same functionality as primary
- Aliases work with all flags
- Help text references canonical command name
- Aliases appear in help (or don't, if cleaner)

## Acceptance Criteria
- [ ] All designed aliases work
- [ ] Aliases pass all existing tests
- [ ] Help text documents aliases
- [ ] No performance impact from aliases
- [ ] Tests verify alias functionality
