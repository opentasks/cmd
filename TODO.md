# OpenTask Project - TODO

## Current Status

✅ **Phase 5: Task Command Refactoring - COMPLETE** (2025-11-28)

All 5 task commands refactored to use TaskService. 16 LOC reduction in cmd layer.

## Phase 5: Task Command Refactoring (Complete)

### Completed Work (2025-11-28)
- [x] Research task refactoring approaches (480-line analysis)
- [x] Extract ID parser utility (internal/task/id.go + 8 tests)
- [x] Create TaskService with 5 CRUD methods (157 LOC)
- [x] Add comprehensive service tests (210 LOC, 6 tests)
- [x] Commit ID parser (commit: 9d0a014)
- [x] Commit TaskService (commit: 372771a)
- [x] Refactor cmd/task.go to use TaskService (5 commands)
  - [x] taskNewCmd → service.Create()
  - [x] taskUpdateCmd → service.Update()
  - [x] taskShowCmd → service.Get()
  - [x] taskListCmd → service.List()
  - [x] taskDeleteCmd → service.Delete()
- [x] Reduction: 279 → 263 LOC (16 line / 6% reduction)
- [x] All tests passing (83.3% task package coverage)
- [x] Final commit (0b0ff0e)

## Phase 4: Security Fixes (COMPLETED 2025-11-27)
- [x] Fix gosec hook configuration
- [x] Install govulncheck via mise
- [x] Fix 3x HIGH: Weak random in banner.go → use crypto/rand (12 → 3 issues)
- [x] Fix 3x MkdirAll permissions: 0755 → 0750
- [x] Fix 3x WriteFile permissions: 0644 → 0600
- [x] Review remaining 3x MEDIUM issues → acceptable for CLI tool
  - G204: $EDITOR subprocess launch (legitimate use case)
  - G304: File operations (bounded by basePath, safe for CLI)
### Future Opportunities
- [ ] Refactor cmd/task.go (279 LOC) - Extract TaskCreator, TaskUpdater services
- [ ] Extract more complex command logic to internal packages
- [ ] Consider API layer for programmatic access to domain logic
- [ ] Add support for more template types
- [ ] Implement configuration hot-reload capability
- [ ] Add configuration validation and schema checking

## Archive

**Phase 2 Tasks** (all complete): See `.memory/summary.md` for detailed metrics
- Analysis: 7 tasks completed
- Implementation: 6 tasks completed
- Testing & Verification: 3 sections completed
- Build & Commit: 5 commits delivered
- Results: 152 LOC reduction, 21 new tests, 100% pass rate

**Phase 1 Tasks** (all complete): See `.memory/implementation-phase1-complete.md`
