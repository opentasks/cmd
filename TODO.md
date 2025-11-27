# OpenTask Project - TODO

## Current Status

✅ **Phase 2: Config & Project Commands Refactoring - COMPLETE**

All tasks implemented, tested, and committed. See `.memory/summary.md` for full completion details.

## Phase 4: Security Fixes & Task Refactoring

### Completed
- [x] Fix gosec hook configuration (2025-11-27)
- [x] Install govulncheck via mise (2025-11-27)

### Security Issues (gosec findings)
- [ ] Address 3x HIGH severity: Weak random in banner.go (lines 94, 100, 107) - use crypto/rand
- [ ] Address 9x MEDIUM severity: File permissions too open
  - [ ] MkdirAll should use 0750 not 0755 (3 locations)
  - [ ] WriteFile should use 0600 not 0644 (3 locations)
  - [ ] File inclusion via variable (3 locations - may be false positive)

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
