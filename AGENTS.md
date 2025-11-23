
# AGENTS.md

> [!NOTE]
> **CRITICAL** Before doing any work, read the entire contents of this file carefully.

<CurrentTasks>
!`cat TODO.md`
</CurrentTasks>


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

**Phase**: Phase 2 - Config & Project Commands Refactoring  
**Status**: ✅ COMPLETE - All Implementation Tasks Finished  
**Worker**: Claude (Central Command)  
**Last Updated**: 2025-11-23

### Phase 2 Completion Summary
- ✅ ConfigInitializer class created and tested (6 tests)
- ✅ ProjectLister class created and tested (10 tests)
- ✅ Template extracted to config.tmpl
- ✅ configInitCmd refactored (100→15 LOC, -85%)
- ✅ configProjectsCmd refactored with ProjectLister
- ✅ projectListCmd refactored with ProjectLister (65→35 LOC, -46%)
- ✅ All 21 new unit tests passing
- ✅ All existing tests passing (100% pass rate)
- ✅ Manual testing completed - output identical to before
- ✅ 5 atomic commits with clear conventional commit messages
- ✅ Code quality maintained, no breaking changes

See `.memory/summary.md` for full project status and detailed metrics.