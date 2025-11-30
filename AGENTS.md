
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

## Searching Memory

- use `grep -r "<search-term>" .memory/` to find relevant notes
- use `grep -r "TODO" TODO.md` to find outstanding tasks

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

**Phase**: 5 - Task Command Refactoring ✅ COMPLETE
**Status**: All implementation complete (3 atomic commits)
**Worker**: OpenCode Agent
**Last Updated**: 2025-11-28
**Deliverables**:
1. Research (480 lines): Single TaskService pattern analysis
2. ID Parser: 17 LOC + 76 test LOC (commit: 9d0a014)
3. TaskService: 154 LOC + 228 test LOC (commit: 372771a)
4. cmd integration: 124 lines refactored (commit: 0b0ff0e)

**Phase**: 6 & 7 - Research ✅ COMPLETE
**Status**: 4 research documents delivered (3,564 total lines)
**Worker**: OpenCode Agent
**Last Updated**: 2025-11-28
**Deliverables**:
1. Deployment research: 774 lines (single binary, goreleaser, Docker, systemd)
2. API integration research: 1100 lines (constructor DI, error handling, validation)
3. Authentication patterns: 1,150 lines (Unix socket, static tokens, OAuth2 analysis)
4. API docs & clients: 540 lines (oapi-codegen, openapi-generator, CI/CD)

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

**Phase 3: Project Subcommand Refactoring** ✅
- ProjectService: 94 LOC, 11 tests (stateful design)
- GlobalConfigSaver: 30 LOC, 4 tests
- cmd/project.go reduced: 332 → 285 LOC (47 line reduction)
- 362 new test LOC (100% coverage)
- Commits: `2f01fe1`, `2c241a4`
- Pattern: Stateful services (instantiated with config)

**Phase 5: Task Command Refactoring** 🔄 (2025-11-28)
- ✅ Research complete: 480-line analysis (`.memory/research-phase5-task-refactoring.md`)
- ✅ ID Parser: 15 LOC + 56 test LOC (commit: `9d0a014`)
- ✅ TaskService: 157 LOC + 210 test LOC (commit: `372771a`)
  - 5 CRUD methods with interface-based DI
  - Pattern: TaskEngine interface for testability
- ⏳ **PENDING**: Refactor cmd/task.go (5 commands, ~120 LOC reduction)

### Next Phase

**Phase 5 Completion** (Next session - 1-2 hours):
1. Refactor cmd/task.go to use TaskService
2. Run full test suite and verify coverage
3. Final commit
4. Update completion metrics

**Phase 6 & 7 Research** (Complete):
1. ✅ **Deployment & Distribution** (Phase 6) - COMPLETE (2025-11-27)
   - Research doc: `.memory/research-phase6-deployment-distribution.md` (774 lines)
   - Recommendation: Single binary with `opentask serve` subcommand
   - Evidence: PocketBase (53k⭐), Ollama (106k⭐), Woodpecker CI (4k⭐)

2. ✅ **API Layer Integration** (Phase 6) - COMPLETE (2025-11-28)
   - Research doc: `.memory/research-phase6-api-layer-integration.md` (1100 lines)
   - Recommendation: Constructor-based DI with Handler Struct
   - Evidence: Grafana (65k⭐), Portainer (30k⭐), Kubernetes, Rancher

3. ✅ **API Documentation & Client Generation** (Phase 7) - COMPLETE (2025-11-28)
   - Research doc: `.memory/research-api-docs-client-generation.md` (540 lines)
   - Recommendation: oapi-codegen (OpenAPI 3.x) + openapi-generator (multi-language clients)
   - Evidence: GUAC, Pomerium, Carbon Aware SDK, Hatchet, Argo CD

**Next Implementation Phases**:
- Phase 6a: API layer implementation (~1-2 weeks, ~600 LOC)
- Phase 7: Multi-language client SDKs (~1 week)
- Phase 6b: Distribution packaging (~1-2 weeks)

See TODO.md for detailed opportunities.