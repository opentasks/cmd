# OpenTask Project - TODO

## Current Status

✅ **Phase 5: Task Command Refactoring - COMPLETE** (2025-11-28)

All implementation complete. 3 atomic commits delivered (id parser, service, cmd integration).

✅ **Phase 6 & 7 Research - COMPLETE** (2025-11-28)

4 comprehensive research documents delivered (3,564 total lines).

## Phase 6 & 7: API Research ✅ COMPLETE (2025-11-28)

### Deliverables
1. **Deployment & Distribution**: `.memory/research-phase6-deployment-distribution.md` (774 lines)
2. **API Layer Integration**: `.memory/research-phase6-api-layer-integration.md` (1,100 lines)
3. **Authentication Patterns**: `.memory/research-phase6-auth-patterns.md` (1,150 lines)
4. **API Documentation & Clients**: `.memory/research-api-docs-client-generation.md` (540 lines)

**Total Research**: 3,564 lines across 4 documents
### Key Findings (Deployment)
- **Recommended pattern**: Single binary with `opentask serve` subcommand
- **Evidence**: PocketBase (53k⭐), Ollama (106k⭐), Woodpecker CI (4k⭐)
- **Distribution**: goreleaser + Docker + systemd + Homebrew
- **API versioning**: `/api/v1` URL path pattern (Concourse, Gin examples)
- **Configuration**: Viper (flags > env > config file > defaults)

### Key Findings (API Integration)
- **Recommended pattern**: Constructor-based DI with Handler Struct (80% adoption)
- **Evidence**: Grafana (65k⭐), Portainer (30k⭐), Kubernetes, Rancher
- **Error handling**: Sentinel errors + `errors.Is()` → HTTP status codes
- **Validation**: go-playground/validator (15+ production repos)
- **Middleware**: stdlib patterns (logging, recovery, CORS)
- **gRPC**: `grpc.UnaryInterceptor` for middleware chaining
- **Testing**: Mock services with interfaces, httptest.NewRecorder

### Key Findings (Authentication)
- **Phase 6 (Local API)**: Unix socket + file permissions (Docker, Podman, LXD pattern)
  - OS-enforced auth (no tokens/keys needed)
  - Zero configuration complexity
  - Industry standard for local CLIs
- **Phase 7+ (Remote API)**: Static bearer tokens (if network access needed)
  - Simple 50 LOC middleware
  - Requires TLS + rate limiting
- **OAuth2/JWT/mTLS**: Overkill for single-user CLI
- **Evidence**: 50+ Go codebases, OWASP API Security Top 10 2023

### Key Findings (API Documentation & Clients)
- **Recommended toolchain**: oapi-codegen (OpenAPI 3.x server) + openapi-generator (multi-language clients)
- **Evidence**: GUAC, Pomerium (oapi-codegen), Carbon Aware SDK, Hatchet (openapi-generator)
- **Interactive docs**: Swagger UI (standard for OpenAPI)
- **CI/CD**: Auto-generate on spec changes, validate with pre-commit hooks
- **Alternative**: grpc-gateway for gRPC + REST dual protocol (Argo CD, CockroachDB)

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

### Implementation Roadmap (Total: ~1500 LOC, 3-4 weeks)

#### Phase 6a: HTTP API Foundation (Week 1, ~600 LOC)
**Files to Create**:
1. `api/http/server.go` - Server struct with constructor DI
2. `api/http/dto.go` - Request/response DTOs with go-playground/validator
3. `api/http/errors.go` - Error mapping (service errors → HTTP status)
4. `api/http/middleware.go` - Logging, recovery, CORS
5. `api/http/handlers_task.go` - Task CRUD (5 endpoints)
6. `api/http/handlers_project.go` - Project handlers
7. `cmd/serve.go` - CLI command to start API server

**Testing**:
- Unit tests: Mock TaskService/ProjectService
- Integration tests: Real services with in-memory storage
- API tests: httptest.NewRecorder

#### Phase 6b: Advanced API Features (Week 2, ~400 LOC)
- Pagination: `?page=1&limit=20`
- Filtering: `?status=todo&type=bug`
- Sorting: `?sort=created_at&order=desc`
- API versioning: `/api/v1/tasks`
- OpenAPI/Swagger spec generation

#### Phase 6c: gRPC API (Week 3, ~500 LOC, optional)
- `api/grpc/proto/task.proto` - Protocol buffer definitions
- `api/grpc/server.go` - gRPC server implementation
- `api/grpc/interceptors.go` - Logging, auth, recovery

#### Phase 6d: Distribution (Week 4)
- Docker multi-stage build
- systemd service unit
- goreleaser packaging updates
- Kubernetes manifests
- Deployment documentation

---

## Phase 5: Task Command Refactoring ✅ COMPLETE (2025-11-28)

### All Work Complete
- [x] Research task refactoring approaches (480-line analysis)
- [x] Extract ID parser utility (internal/task/id.go + 76 test LOC)
- [x] Create TaskService with 5 CRUD methods (154 LOC + 228 test LOC)
- [x] Commit ID parser (commit: 9d0a014)
- [x] Commit TaskService (commit: 372771a)
- [x] Refactor cmd/task.go to use TaskService (5 commands)
  - [x] taskNewCmd → service.Create()
  - [x] taskUpdateCmd → service.Update()
  - [x] taskShowCmd → service.Get()
  - [x] taskListCmd → service.List()
  - [x] taskDeleteCmd → service.Delete()
- [x] Actual reduction: 279 → 263 LOC (6% via 124-line refactor)
- [x] All tests pass (83.3% task package coverage)
- [x] Final commit (commit: 0b0ff0e)
- [x] Documentation updated with completion metrics

### Final Metrics
- **Code added**: 171 LOC (service 154 + id 17)
- **Tests added**: 304 LOC (service_test 228 + id_test 76)
- **cmd refactored**: 124 lines changed (-70, +54 = net -16 LOC)
- **Commits**: 3 atomic commits
- **Pattern**: Single TaskService (ecosystem-proven, Phase 6 ready)

### Phase 4: Security Fixes (COMPLETED 2025-11-27)
- [x] Fix gosec hook configuration
- [x] Install govulncheck via mise
- [x] Fix 3x HIGH: Weak random in banner.go → use crypto/rand (12 → 3 issues)
- [x] Fix 3x MkdirAll permissions: 0755 → 0750
- [x] Fix 3x WriteFile permissions: 0644 → 0600
- [x] Review remaining 3x MEDIUM issues → acceptable for CLI tool
  - G204: $EDITOR subprocess launch (legitimate use case)
  - G304: File operations (bounded by basePath, safe for CLI)

### Phase 6: Deployment & Distribution Research ✅ COMPLETE (2025-11-27)
- [x] Analyze CLI+API hybrid patterns (PocketBase, Ollama, Woodpecker CI)
- [x] Document 774-line research with goreleaser, Docker, systemd patterns
- [x] Cite 10+ production tools (GitHub stars: 1k-106k)
- **Recommendation:** Single binary with `opentask serve` subcommand
- **File:** `.memory/research-phase6-deployment-distribution.md`

### Phase 6 (continued): API Layer Integration Research ✅ COMPLETE (2025-11-28)
- [x] Compare 3 patterns: Library Packages, Handler Structs, Embedded Servers
- [x] Document 1100-line research with constructor DI, error handling, validation
- [x] Cite 20+ production examples (Grafana, Portainer, Kubernetes, Rancher)
- **Recommendation:** Constructor-based DI with Handler Struct (80% adoption)
- **File:** `.memory/research-phase6-api-layer-integration.md`

### Phase 7: API Documentation & Client Generation Research ✅ COMPLETE (2025-11-28)
- [x] Compare 5 major tools (swaggo, oapi-codegen, grpc-gateway, gqlgen, openapi-generator)
- [x] Document 540-line research with production examples
- [x] Evaluate OpenAPI 3.x vs Swagger 2.0 vs gRPC vs GraphQL
- [x] Provide CI/CD integration patterns (GitHub Actions, pre-commit hooks)
- [x] Multi-language client generation strategy (Python, TypeScript, Go)
- [x] Cite 15+ production examples (GUAC, Pomerium, Argo CD, Hatchet)
- **Recommendation:** oapi-codegen (server) + openapi-generator (clients)
- **File:** `.memory/research-api-docs-client-generation.md`

### Future Opportunities (Phase 6a - API Implementation)
- [ ] Design OpenAPI 3.x spec (`api/openapi.yaml`)
- [ ] Generate Go server with oapi-codegen (strict-server mode)
- [ ] Implement HTTP handlers using TaskService
- [ ] Add Swagger UI at `/swagger` endpoint
- [ ] CI/CD: Validate OpenAPI spec on pre-commit

### Future Opportunities (Phase 7 - Client SDKs)
- [ ] Generate Python client with openapi-generator
- [ ] Generate TypeScript client with openapi-generator
- [ ] Add client integration tests
- [ ] Document SDK usage examples

### Other Opportunities
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
