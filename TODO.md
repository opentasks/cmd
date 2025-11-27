# OpenTask Project - TODO

## Current Status

✅ **Phase 2: Config & Project Commands Refactoring - COMPLETE**

All tasks implemented, tested, and committed. See `.memory/summary.md` for full completion details.

## Phase 4: Security Fixes & Task Refactoring

### Completed Security Fixes (2025-11-27)
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
