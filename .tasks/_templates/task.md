---
id: <<NEXT_ID>>
title: <<TASK_TITLE>>
type: task
status: todo
tags: [<<TAG1>>, <<TAG2>>]
relationships:
  - type: parent
    taskID: <<PARENT_EPIC_ID>>
createdAt: <<NOW_RFC3339>>
updatedAt: <<NOW_RFC3339>>
---

# <<TASK_TITLE>>

## Description

**What needs to be done?**

Clear, specific description of the work.

**Why is this necessary?**

Context or justification for the task.

## Requirements

- [ ] Requirement 1
- [ ] Requirement 2
- [ ] Requirement 3

## Work Items

- [ ] Step 1
- [ ] Step 2
- [ ] Step 3

## Definition of Done

- [ ] Work completed
- [ ] Code committed
- [ ] Tests passing
- [ ] Peer reviewed
- [ ] Documentation updated

## Technical Notes

Implementation details, gotchas, or approach.

## Dependencies

Other tasks that must be completed first:
- Task X (blocking)
- Task Y (helpful context)

## Implementation Log

Track progress:
- Started: <<DATE>>
- Completed: <<DATE>>

## Notes

Additional context or considerations.

---

**Template Variables to Replace:**
- `<<NEXT_ID>>` - Sequential ID
- `<<TASK_TITLE>>` - Task title
- `<<TAG1>>, <<TAG2>>` - Tags like "chore", "refactor", "testing" (optional)
- `<<PARENT_EPIC_ID>>` - ID of parent epic
- `<<NOW_RFC3339>>` - Current timestamp (e.g., 2025-11-02T19:30:00Z)
- `<<DATE>>` - Date when status changed
