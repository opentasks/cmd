

# AGENTS.md

> [!NOTE]
> **CRITICAL** Before doing any work,
> - read `.memory/todo.md`, `.memory/summary.md` and `.memory/team.md`
> - if `.memory/` is missing these files, then create those three.
> - Use relevant in-memoria tools before starting a task or thinking about an answer.


## Project Guidelines

- never commit binaries
- when adjusting .gitignore, stop. get human help.

## Research Guidelines

- [knowledge] store findings in `.memory/` directory
- [knowledge] all notes in `.memory/` must be in markdown format
- [knowledge] except for `.memory/summary.md`, all notes in `.memory/` must follow the filename convention of `.memory/<type>-<id>-<title>.md`
- [knowledge] where `<type>` is one of: `research`, `phase`, `guide`, `notes`, `implementation`, `task`
- [knowledge] Always keep `.memory/summary.md` up to date with current status, prune incorrect or outdated information.
- [tasks] when finishing a phase, compact relevant successful outcomes from implementation, research and phase into the `.memory/summary.md` and delete the other files. empty `.memory/todo.md` of completed tasks.
- [tasks] break down tasks into manageable phases, each with clear objectives and deliverables.
- [tasks] use `.memory/todo.md` to track remaining tasks. This file only contains links to `.memory/task-<id>-<title>.md` files. [CRITICAL] keep `.memory/todo.md` up to date at every step.
- [git] when committing changes, follow conventional commit guidelines.
- [git] Use clear commit messages referencing relevant files for changes.

## Searching Memory

- use `grep -r "<search-term>" .memory/` to find relevant notes
- use `grep -r "TODO" .memory/todo.md` to find outstanding tasks

## In-Memoria Intelligence (Codebase Navigation)

**When to use In-Memoria** (prefer over manual grep/glob):
- [ALWAYS] When asked "where should I..." or "what files..." → Use `memoria_predict_coding_approach`
- [ALWAYS] When searching for patterns across codebase → Use `memoria_search_codebase`
- [START OF SESSION] Get project overview → Use `memoria_get_project_blueprint`
- [BEFORE REFACTORING] Check existing patterns → Use `memoria_get_pattern_recommendations`

**Quick commands:**
```bash
# Project overview (run once per session)
memoria_get_project_blueprint(path="/mnt/Store/Projects/Mine/Github/opentasks")

# Find files for a task
memoria_predict_coding_approach(problemDescription="add task filtering by assignee")

# Search for pattern usage
memoria_search_codebase(query="ProjectService", type="text")

# Get recommended patterns for new code
memoria_get_pattern_recommendations(problemDescription="create new service class")
```

**Benefits:**
- ✅ 10-100x faster than manual grep/find
- ✅ Learns project patterns (Factory, Builder, Service patterns)
- ✅ No more "exploring to understand" - instant routing
- ✅ Confidence scores guide decisions

**When NOT to use:**
- Reading specific known files → Use `read` tool directly
- Exact file path operations → Use bash/filesystem tools

## Execution Steps

0. always read `.memory/summary.md` first to understand successful outcomes so far.
1. update `.memory/team.md` to indicate which phase is being worked on and by whom (use the session id to indicate this, not the agent name).
2. If there are any `[NEEDS-HUMAN]` tasks in `.memory/todo.md`, stop and wait for human intervention.
3. follow the research guidelines above.
4. when you are blocked by actions that require human intervention, create a `.memory/todo.md` file listing the tasks that need to be done by a human. tag it with `[NEEDS-HUMAN]` on the task line.
5. after completing a phase, update `.memory/summary.md` and prune other files as necessary.
6. commit changes with clear messages referencing relevant files.

## Human Interaction

- If you need clarification or additional information, please ask a human for assistance.
- print a large ascii box in chat indicating that human intervention is needed, and list the tasks from `.memory/todo.md` inside the box.
- wait for human to complete the tasks before proceeding.
