
# AGENTS.md

> [!NOTE]
> **CRITICAL** Before doing any work, read the entire contents of this file carefully.

<CurrentTasks>
!`cat TODO.md`
</CurrentTasks>


## Project Guidelines

- golang project
- ALWAYS use mise tasks for any work done
- never commit binaries
- when adjusting .gitignore, stop. get human help.

## Research Guidelines

- [knowledge] store findings in `.memory/` directory
- [knowledge] all notes in `.memory/` must be in markdown format
- [knowledge] except for `.memory/summary.md`, all notes in `.memory/` must follow the filename convention of `.memory/<type>-<id>-<title>.md`
- [knowledge] where `<type>` is one of: `research`, `phase`, `guide`, `notes`, `implementation`
- [knowledge] Always keep `.memory/summary.md` up to date with current status, prune incorrect or outdated information.
- [tasks] when finishing a phase, compact relevant successful outcomes from implementation, research and phase into the `.memory/summary.md` and delete the other files. empty `TODO.md` of completed tasks.
- [tasks] break down tasks into manageable phases, each with clear objectives and deliverables.
- [tasks] use `TODO.md` to track remaining tasks. [CRITICAL] keep `TODO.md` up to date at every step.
- [git] when committing changes, follow conventional commit guidelines.
- [git] Use clear commit messages referencing relevant files for changes.

## Execution Steps

1. always read `.memory/summary.md` first to understand successful outcomes so far.
2. update `AGENTS.md` to indicate which phase is being worked on and by whom.
3. If there are any `[NEEDS-HUMAN]` tasks in `TODO.md`, stop and wait for human intervention.
4. follow the research guidelines above.
5. when you are blocked by actions that require human intervention, create a `TODO.md` file listing the tasks that need to be done by a human. tag it with `[NEEDS-HUMAN]` on the task line.
6. after completing a phase, update `.memory/summary.md` and prune other files as necessary.
7. commit changes with clear messages referencing relevant files.

## Human Interaction

- If you need clarification or additional information, please ask a human for assistance.
- print a large ascii box in chat indicating that human intervention is needed, and list the tasks from `TODO.md` inside the box.
- wait for human to complete the tasks before proceeding.

## Current Work Status

**Phase**: Ready for Phase 3 Planning  
**Status**: ✅ Phase 2 COMPLETE & ARCHIVED  
**Worker**: (awaiting Phase 3 assignment)  
**Last Updated**: 2025-11-23

### Completed Phases

**Phase 1: Domain Logic Extraction** ✅
- See `.memory/implementation-phase1-complete.md`

**Phase 2: Config & Project Commands Refactoring** ✅
- ConfigInitializer class: 99 LOC, 6 tests
- ProjectLister class: 83 LOC, 10 tests
- config.tmpl template extracted
- 3 command handlers refactored
- 152 LOC reduction in cmd/ layer
- 21 new unit tests (100% coverage)
- 5 atomic commits with conventional messages
- See `.memory/summary.md` for complete details

### Next Phase

Phase 3 planning available in TODO.md with future opportunities.