---
id: <<NEXT_ID>>
title: <<DECISION_TITLE>>
type: decision
status: todo
tags: [<<TAG1>>, <<TAG2>>]
relationships:
  - type: parent
    taskID: <<PARENT_EPIC_ID>>
createdAt: <<NOW_RFC3339>>
updatedAt: <<NOW_RFC3339>>
---

# <<DECISION_TITLE>>

## Decision Statement

We will **<<DECISION>>** because **<<RATIONALE>>**.

## Context

**What problem are we trying to solve?**

Background and business context that led to this decision.

**Why does this matter?**

Impact on the project, team, or architecture.

## Options Considered

### Option 1: <<OPTION_NAME>>

**Approach:**
Description of this approach.

**Pros:**
- Pro 1
- Pro 2

**Cons:**
- Con 1
- Con 2

**Effort:** Medium

### Option 2: <<OPTION_NAME>>

**Approach:**
Description of this approach.

**Pros:**
- Pro 1
- Pro 2

**Cons:**
- Con 1
- Con 2

**Effort:** High

### Option 3: <<OPTION_NAME>> (Rejected)

**Reason for Rejection:**
Why this wasn't chosen.

## Decision

**Selected:** Option <<NUM>> - <<OPTION_NAME>>

**Decided by:** <<PERSON/TEAM>>

**Decision date:** <<DATE>>

**Status:** <<APPROVED/PENDING/IMPLEMENTED>>

## Rationale

Detailed explanation of why this decision was made:

1. Primary reasons
2. Trade-offs accepted
3. Long-term implications

## Consequences

**Positive Impacts:**
- Impact 1
- Impact 2

**Negative Impacts:**
- Impact 1
- Impact 2

**Mitigation:**
How we'll address negative impacts.

## Alternatives for Future Reconsideration

Conditions under which we might revisit this decision:
- Condition 1
- Condition 2

## Implementation Notes

How this decision will be implemented:
- Step 1
- Step 2

## Review and Approval

- Proposed by: <<PERSON>>
- Reviewed by: <<PERSON>>
- Approved by: <<PERSON/TEAM>>
- Implementation owner: <<PERSON>>

## Decision Log

- Proposed: <<DATE>>
- Reviewed: <<DATE>>
- Approved: <<DATE>>
- Implemented: <<DATE>>

## References

Related decisions or documentation:
- Decision X
- Design doc
- RFC or proposal

## Notes

Additional context or follow-up items.

---

**Template Variables to Replace:**
- `<<NEXT_ID>>` - Sequential ID
- `<<DECISION_TITLE>>` - Decision title
- `<<DECISION>>` - Clear statement of decision
- `<<RATIONALE>>` - Brief reason
- `<<TAG1>>, <<TAG2>>` - Tags like "architecture", "technology" (optional)
- `<<PARENT_EPIC_ID>>` - ID of parent epic
- `<<NOW_RFC3339>>` - Current timestamp (e.g., 2025-11-02T19:30:00Z)
- `<<OPTION_NAME>>` - Name of option
- `<<NUM>>` - Option number
- `<<PERSON/TEAM>>` - Name or team
- `<<DATE>>` - Specific date
- `<<APPROVED/PENDING/IMPLEMENTED>>` - Status
