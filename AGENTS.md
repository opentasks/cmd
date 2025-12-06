

# AGENTS.md

> [!NOTE]
> **CRITICAL** Before doing any work:
> - Check for active tasks: `opentask task list --status in_progress`
> - Review pending work: `opentask task list --status pending`
> - Check for blocked work: `opentask task list --tag BLOCKED`


## Project Guidelines

- golang project
- ALWAYS use mise tasks for any work done
- never commit binaries
- when adjusting .gitignore, stop. get human help.

## Research & Documentation Guidelines

- [knowledge] Document findings in opentask task descriptions using `--description` flag
- [knowledge] All notes must be in markdown format
- [knowledge] Use opentask tags to categorize work: `research`, `documentation`, `implementation`, etc.
- [knowledge] Keep task descriptions up to date with current status and findings
- [tasks] Break down work into manageable phases using opentask hierarchy (parent/child tasks)
- [tasks] Use opentask task types to organize: `phase` for multi-step work, `task` for atomic units
- [tasks] Track all work through opentask - use `opentask task list` to find remaining items
- [git] When committing changes, follow conventional commit guidelines
- [git] Use clear commit messages referencing relevant task IDs and files (e.g., "fix: improve auth flow (task #42)")

## Task Management with Opentask

**Core opentask commands for workflow:**
- `opentask task new "<title>" --type phase` - Create tracking tasks for work phases
- `opentask task new "<title>" --type task --parent <id>` - Create subtasks for atomic work
- `opentask task list --status pending` - View all pending work items
- `opentask task list --status in_progress` - View active work
- `opentask task update <id> --status in_progress` - Mark work as active
- `opentask task update <id> --status completed` - Mark work as complete
- `opentask task show <id>` - View task details and context

**Task types for agent workflows:**
- `type:phase` - Multi-step implementation phases with clear objectives
- `type:task` - Atomic work units within phases
- `type:research` - Investigation/discovery work
- `type:epic` - Large features spanning multiple phases
- `type:story` - User-facing features or use cases

**Tags for workflow control:**
- `tag:NEEDS-HUMAN` - Work blocked on human decision/intervention
- `tag:BLOCKED` - Work blocked on external dependencies
- `tag:CRITICAL` - High-priority work requiring immediate attention

**Benefits:**
- ✅ Native persistence across sessions
- ✅ Structured task hierarchy and filtering
- ✅ Integrated with project context
- ✅ No manual file management overhead
- ✅ Clear status tracking for agent coordination

**Workflow example:**
```bash
# Create a new phase
opentask task new "Implement authentication feature" --type phase --description "Multi-step implementation"

# Create subtasks for the phase (get phase ID from list)
opentask task new "Research OAuth libraries" --type research --parent <phase-id>
opentask task new "Implement OAuth flow" --type task --parent <phase-id>

# Track progress
opentask task update <task-id> --status in_progress
opentask task update <task-id> --status completed

# Find blocked work
opentask task list --tag BLOCKED
```

## Execution Steps

0. Check task status: `opentask task list --status pending` to understand remaining work.
1. Create/update tasks as needed to track work phases and atomic units.
2. If there are any `--tag BLOCKED` or `--tag NEEDS-HUMAN` tasks, stop and wait for human intervention.
3. Follow the research and documentation guidelines above.
4. When blocked by actions requiring human intervention, create a task with `--tag NEEDS-HUMAN` describing what needs to be done.
5. Mark tasks as `--status completed` when work is finished.
6. Commit changes with clear messages referencing relevant task IDs (e.g., "task #42").

## Human Interaction

- If you need clarification or additional information, please ask a human for assistance.
- When human intervention is needed, create a task with `--tag NEEDS-HUMAN` describing what must be done.
- Print a large ascii box in chat indicating that human intervention is needed, and run `opentask task list --tag NEEDS-HUMAN` to show blocking tasks.
- Wait for human to complete the tasks before proceeding.
