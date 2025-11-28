# OpenTask Project - TODO

## Current Status

🔄 **Phase 5: Task Command Refactoring - IN PROGRESS** (2025-11-28)

Service layer complete. Cmd integration pending.

## Phase 5: Task Command Refactoring (Partial)

### Completed Work (2025-11-28)
- [x] Research task refactoring approaches (480-line analysis)
- [x] Extract ID parser utility (internal/task/id.go + 8 tests)
- [x] Create TaskService with 5 CRUD methods (157 LOC)
- [x] Add comprehensive service tests (210 LOC, 6 tests)
- [x] Commit ID parser (commit: 9d0a014)
- [x] Commit TaskService (commit: 372771a)

### Pending Work
- [ ] Refactor cmd/task.go to use TaskService (5 commands)
  - [ ] taskNewCmd → service.Create()
  - [ ] taskUpdateCmd → service.Update()
  - [ ] taskShowCmd → service.Get()
  - [ ] taskListCmd → service.List()
  - [ ] taskDeleteCmd → service.Delete()
- [ ] Expected reduction: 279 → ~160 LOC (43%)
- [ ] Run full test suite
- [ ] Verify 100% coverage maintained
- [ ] Final commit with conventional message
- [ ] Update .memory/summary.md with completion metrics

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
