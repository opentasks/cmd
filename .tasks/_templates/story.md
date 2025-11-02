---
id: <<NEXT_ID>>
title: <<STORY_TITLE>>
type: story
status: todo
tags: [<<TAG1>>, <<TAG2>>]
relationships:
  - type: parent
    taskID: <<PARENT_EPIC_ID>>
createdAt: <<NOW_RFC3339>>
updatedAt: <<NOW_RFC3339>>
---

# <<STORY_TITLE>>

## Description

**What is this story?**

User-facing feature or capability being implemented.

**Why is it needed?**

Business value, user need, or improvement.

## Acceptance Criteria

- [ ] Criterion 1
- [ ] Criterion 2
- [ ] Criterion 3

## Definition of Done

- [ ] Code implemented
- [ ] Code reviewed and approved
- [ ] Unit tests written and passing
- [ ] Integration tests passing
- [ ] Documentation updated
- [ ] Ready for review/QA

## Implementation Notes

Technical approach, design decisions, or special considerations.

## Related Tasks

Link to related stories or tasks using relationships:
- `relationships: {type: relates-to, taskID: X}`
- `relationships: {type: blocks, taskID: Y}`

## Implementation Log

Track progress and changes:
- Started: <<DATE>>
- In review: <<DATE>>
- Completed: <<DATE>>

## Notes

Additional context or decisions.

---

**Template Variables to Replace:**
- `<<NEXT_ID>>` - Sequential ID
- `<<STORY_TITLE>>` - Story title
- `<<TAG1>>, <<TAG2>>` - Tags like "feature", "backend", "ui" (optional)
- `<<PARENT_EPIC_ID>>` - ID of parent epic
- `<<NOW_RFC3339>>` - Current timestamp (e.g., 2025-11-02T19:30:00Z)
- `<<DATE>>` - Date when status changed
