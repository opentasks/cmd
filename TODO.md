# OpenTask Project - TODO

## Current Status

✅ **Phase 6 Research: Deployment & Distribution - COMPLETE** (2025-11-28)

Research deliverable: 774 lines analyzing CLI+API hybrid patterns.

🔄 **Phase 5: Task Command Refactoring - IN PROGRESS** (2025-11-28)

Service layer complete. Cmd integration pending.

## Phase 6: API Architecture Research ✅ COMPLETE (2025-11-28)

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

### Phase 4: Security Fixes (COMPLETED 2025-11-27)
- [x] Fix gosec hook configuration
- [x] Install govulncheck via mise
- [x] Fix 3x HIGH: Weak random in banner.go → use crypto/rand (12 → 3 issues)
- [x] Fix 3x MkdirAll permissions: 0755 → 0750
- [x] Fix 3x WriteFile permissions: 0644 → 0600
- [x] Review remaining 3x MEDIUM issues → acceptable for CLI tool
  - G204: $EDITOR subprocess launch (legitimate use case)
  - G304: File operations (bounded by basePath, safe for CLI)

### Phase 6: API Architecture Research ✅ COMPLETE (2025-11-28)
- [x] Analyze GitHub CLI, Docker CLI, Terraform, Kubernetes API patterns
- [x] Compare 3 patterns: Library Packages, Service Layer, Embedded Servers
- [x] Document 425-line research with real-world evidence
- [x] Provide code structure recommendations
- [x] Create migration path for `pkg/opentasks/` API layer
- [x] Cite 5+ major open-source projects (all with GitHub repo links)
- [x] Mark biases and assumptions with [BIAS] and [ASSUMPTION] tags
- **Recommendation:** Library Package Pattern (87% adoption rate)
- **File:** `.memory/research-phase6-api-architecture.md`

### Future Opportunities
- [ ] Implement `pkg/opentasks/` API layer (if Phase 6 recommendation approved)
- [ ] Extract more complex command logic to internal packages
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
