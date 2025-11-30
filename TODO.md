# OpenTask Project - TODO

## Current Status

✅ **Phase 6 Research: Deployment & Distribution - COMPLETE** (2025-11-28)

Research deliverable: 774 lines analyzing CLI+API hybrid patterns.

✅ **Phase 5: Task Command Refactoring - COMPLETE** (2025-11-28)

All 5 task commands refactored to use TaskService. 16 LOC reduction in cmd layer.

## Phase 6: Deployment & Distribution Research ✅ COMPLETE (2025-11-28)

### Deliverable
**Research document**: `.memory/research-phase6-deployment-distribution.md` (774 lines)

### Key Findings
- **Recommended pattern**: Single binary with `opentask serve` subcommand
- **Evidence**: PocketBase (53k⭐), Ollama (106k⭐), Woodpecker CI (4k⭐)
- **Distribution**: goreleaser + Docker + systemd + Homebrew
- **API versioning**: `/api/v1` URL path pattern (Concourse, Gin examples)
- **Configuration**: Viper (flags > env > config file > defaults)

### Architecture Decision
```
opentask              # CLI mode (existing)
opentask serve        # API server mode (new)
opentask --help       # Shows all commands
```

### Citations Provided
- 10+ production tools analyzed (GitHub stars: 1k-106k)
- goreleaser configs from PocketBase, CLI
- Dockerfile patterns from Ollama, PocketBase
- systemd units from Redis, Valkey
- Kubernetes deployments from Concourse
- Homebrew formula patterns from GitHub CLI

### Next Steps (Phase 6a - Implementation)
1. Add `cmd/opentask/serve.go` (API server command)
2. Extract HTTP handlers to `internal/api/handlers/`
3. Integrate TaskService from Phase 5
4. Basic routes: `/api/v1/tasks` (GET, POST)
5. Viper config for server settings

### Deferred to Phase 6b (Distribution)
- Docker multi-stage build
- systemd service unit
- goreleaser packaging updates
- Deployment documentation

---

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
- [ ] Refactor additional commands following Phase 5 pattern
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
