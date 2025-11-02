---
id: 19
title: Refine CLI UX for epic item display
type: story
status: todo
tags:
    - cli
    - ux
relationships:
    - type: parent
      taskID: 1
createdAt: "2025-11-02T20:25:00Z"
updatedAt: "2025-11-02T20:25:00Z"
---

## Problem
Currently, when listing tasks, it's not clear which epic (parent) items belong to. Users need better visibility into the epic hierarchy when viewing tasks in the CLI.

## Considerations
Multiple display options to explore:
- **Grouping/collation**: Show tasks grouped by their parent epic
- **Tree view**: Display hierarchical structure showing parent-child relationships
- **Epic column**: Always include an "Epic" column in task listings showing the parent epic
- **Nested display**: Indent child tasks under their parent epic

## Questions to answer
1. Which approach best fits the CLI workflow?
2. Should this be a default behavior or configurable via flags?
3. How should this interact with filtering and sorting?
4. What about tasks without epic parents (root-level tasks)?

## Acceptance Criteria
- [ ] Explore and document at least 2-3 different UX approaches
- [ ] Design the preferred solution
- [ ] Implement the chosen approach
- [ ] Update tests to cover new CLI output formats
- [ ] Document the feature in CLI help/docs
