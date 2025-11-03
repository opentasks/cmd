---
id: <<NEXT_ID>>
title: <<PLAN_TITLE>>
type: plan
status: todo
tags: [<<TAG1>>, <<TAG2>>]
relationships:
  - type: parent
    taskID: <<PARENT_EPIC_ID>>
createdAt: <<NOW_RFC3339>>
updatedAt: <<NOW_RFC3339>>
---

# <<PLAN_TITLE>>

## Overview

**What is being planned?**

High-level overview of the plan.

**Why is this necessary?**

Strategic importance or business value.

**Timeline**

Expected duration or target completion date.

## Goals

- [ ] Goal 1
- [ ] Goal 2
- [ ] Goal 3

## Phases

### Phase 1: <<PHASE_NAME>>

**Duration:** <<DURATION>>

**Deliverables:**
- Deliverable 1
- Deliverable 2

**Success Criteria:**
- Criterion 1
- Criterion 2

### Phase 2: <<PHASE_NAME>>

**Duration:** <<DURATION>>

**Deliverables:**
- Deliverable 1
- Deliverable 2

**Success Criteria:**
- Criterion 1
- Criterion 2

## Milestones

| Milestone | Target Date | Status |
|-----------|------------|--------|
| Milestone 1 | <<DATE>> | - |
| Milestone 2 | <<DATE>> | - |
| Milestone 3 | <<DATE>> | - |

## Resources Required

- Team members: <<COUNT>>
- Budget: <<AMOUNT>>
- Tools/services: <<LIST>>

## Risks and Mitigation

| Risk | Probability | Impact | Mitigation |
|------|------------|--------|-----------|
| Risk 1 | High | High | Mitigation strategy |
| Risk 2 | Medium | Medium | Mitigation strategy |

## Dependencies

- Task X (must complete first)
- Resource Y (needed for implementation)

## Communication Plan

How and when updates will be shared:
- Weekly updates via <<CHANNEL>>
- Milestone reviews with stakeholders
- Risk escalation if needed

## Planning Log

Track planning progress:
- Plan created: <<DATE>>
- Approved: <<DATE>>
- Execution started: <<DATE>>

## Notes

Additional context, assumptions, or decisions.

---

**Template Variables to Replace:**
- `<<NEXT_ID>>` - Sequential ID
- `<<PLAN_TITLE>>` - Plan name
- `<<TAG1>>, <<TAG2>>` - Tags like "planning", "roadmap" (optional)
- `<<PARENT_EPIC_ID>>` - ID of parent epic
- `<<NOW_RFC3339>>` - Current timestamp (e.g., 2025-11-02T19:30:00Z)
- `<<PHASE_NAME>>` - Name of phase
- `<<DURATION>>` - Duration (e.g., "1 week", "2 sprints")
- `<<DATE>>` - Specific date
- `<<COUNT>>` - Number of resources
- `<<AMOUNT>>` - Budget amount
- `<<LIST>>` - List of items
- `<<CHANNEL>>` - Communication channel
