# Phase 4: Storage Backend Plugins & Extensibility

**Status**: Planned (Q2 2026)  
**Epic ID**: 51  
**Scope**: 8 tasks (research, plan, stories, implementation)

## Overview

Transform opentask from a single-backend system into an extensible plugin platform. Enable organizations to choose their optimal storage backend while maintaining a single CLI and API.

## Architecture

### HashiCorp go-plugin

We'll use `github.com/hashicorp/go-plugin` for:

- **Process Isolation** – Plugins run in subprocess, crashes don't affect core
- **RPC/gRPC** – Plugins communicate via standard protocols
- **Version Management** – Handle protocol evolution gracefully
- **Security** – Sandbox untrusted plugins

```go
type StoragePlugin interface {
    Create(ctx context.Context, task *Task) error
    Read(ctx context.Context, id string) (*Task, error)
    Update(ctx context.Context, task *Task) error
    Delete(ctx context.Context, id string) error
    List(ctx context.Context, filter Filter) ([]*Task, error)
    Close() error
}
```

## Storage Backends

### Tier 1: Core (Phase 4, Month 1)
- **Markdown-FS** (Built-in, always available)
- **SQLite** – Plugin: Single-file, portable, zero-setup at scale

### Tier 2: Team Scale (Phase 4, Month 2)
- **DuckDB** – Plugin: Analytical queries, OLAP-friendly
- **PostgreSQL** – Plugin: Team database, concurrent access, RBAC

### Tier 3: Enterprise (Phase 4, Month 3)
- **S3** – Plugin: AWS artifact storage
- **GCS** – Plugin: Google Cloud storage
- **Azure Blob** – Plugin: Microsoft cloud storage

## Timeline

```
Q2 2026 (3-4 months)
├─ Month 1: Plugin architecture + go-plugin + SQLite
├─ Month 2: DuckDB + PostgreSQL + migration tools
├─ Month 3: S3/GCS/Azure implementations
└─ Month 4: Documentation + examples + optimization
```

## Features

### Plugin System
- [ ] Plugin discovery and loading
- [ ] Plugin lifecycle (init, validate, shutdown)
- [ ] Plugin configuration via TOML
- [ ] Error handling and logging
- [ ] Plugin marketplace/registry (future)

### Storage Backends
- [ ] SQLite backend plugin
- [ ] DuckDB backend plugin
- [ ] PostgreSQL backend plugin
- [ ] S3 backend plugin
- [ ] GCS backend plugin
- [ ] Azure Blob backend plugin

### Tools
- [ ] Export/import utilities (all backends)
- [ ] Zero-downtime migration scripts
- [ ] Data validation and integrity checks
- [ ] Rollback procedures
- [ ] Performance benchmarks

### Documentation
- [ ] Plugin developer guide
- [ ] API reference
- [ ] Example custom plugin
- [ ] Migration guide
- [ ] Best practices

## Use Cases

### Solo Developer
```bash
# Start simple
opentask --path my_project task new "Build feature" --type story

# Scale when needed
opentask config set storage.backend sqlite
opentask migrate markdown-fs sqlite
```

### Small Team
```bash
# Shared PostgreSQL
[storage]
backend = "postgresql"
[storage.postgresql]
host = "team-db.internal"
port = 5432
database = "opentask"
```

### Enterprise
```bash
# Cloud-native with S3
[storage]
backend = "s3"
[storage.s3]
bucket = "company-tasks"
region = "us-west-2"
# Uses AWS IAM/STS for auth
```

## Success Criteria

✓ Plugin system architecture is documented  
✓ HashiCorp go-plugin integration is working  
✓ 4+ backend plugins are functional  
✓ Migration tool is tested across all backends  
✓ Comprehensive plugin documentation  
✓ Example custom plugin provided  
✓ End-to-end tests for plugin lifecycle  
✓ Performance benchmarks published  

## Risks & Mitigation

| Risk | Mitigation |
|------|-----------|
| Plugin API breaks ecosystem | Strict versioning, 2-version support |
| RPC overhead impacts performance | Caching layer, async operations |
| Security vulnerabilities | Sandboxing, optional code signing |
| Data loss during migration | Validation, dry-run mode, backups |

## Dependencies

**Required:**
- `github.com/hashicorp/go-plugin` (Apache 2.0)

**Optional (per plugin):**
- `github.com/mattn/go-sqlite3` (SQLite)
- `github.com/marcboeker/go-duckdb` (DuckDB)
- `github.com/lib/pq` (PostgreSQL)
- `github.com/aws/aws-sdk-go-v2` (S3)
- `cloud.google.com/go/storage` (GCS)
- `github.com/Azure/azure-sdk-for-go` (Azure)

All backend-specific dependencies are plugin responsibilities, keeping core lightweight.

## Related Tasks

From opentask epic 51:

- **Task 52**: Research: Plugin Architecture & HashiCorp go-plugin
- **Task 53**: Plan: Storage Plugin System Design
- **Task 54**: Story: Implement HashiCorp go-plugin Integration
- **Task 55**: Story: SQLite Backend Plugin
- **Task 56**: Story: DuckDB Backend Plugin
- **Task 57**: Story: PostgreSQL Backend Plugin
- **Task 58**: Story: Migration Tool Between Storage Backends
- **Task 59**: Task: Cloud Provider Integrations

## Next Steps (After Phase 3)

1. Research go-plugin ecosystem and design (task 52)
2. Create detailed architecture document (task 53)
3. Implement plugin framework and go-plugin integration (task 54)
4. Build SQLite plugin as reference (task 55)
5. Iterate through remaining backends
6. Write comprehensive documentation
7. Release Phase 4 with full plugin system

---

**Phase**: 4 | **Timeline**: Q2 2026 | **Owner**: TBD | **Priority**: Medium
