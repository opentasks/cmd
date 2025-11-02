# Context for Next Session

**Status**: Design phase complete. Ready for implementation.

## Quick Start for New Session

1. **Read first**: `DESIGN_SUMMARY.md` (5 min)
2. **Reference**: `.tasks/design/` folder for detailed specs
3. **Implementation**: See implementation plan below

## Project State

- **Repository**: `/mnt/Store/Projects/Mine/Github/opentask`
- **Current branch**: `master`
- **Last commit**: Design summary documentation

## What's Complete

✅ **Architecture designed**:
- Data models (Task, Relationship, Config structs)
- BaseStorage interface with MD5FileStorage implementation
- Query engine with functional options pattern
- Configuration system (config.toml schema)
- Semantic ID system (per-type sequential with collision detection)
- CLI structure (Viper/Cobra)

✅ **Documentation complete**:
- `README.md`: Project vision and overview
- `DESIGN_SUMMARY.md`: Quick reference guide
- `.tasks/design/`: Full design docs with rationale
- `.tasks/design/r-*.research.md`: Implementation specification

## What's Next (Implementation)

### Phase 1: Core System (MVP)

**1. Project setup**
- [ ] Initialize Go module
- [ ] Add dependencies (viper, cobra, yaml, etc.)
- [ ] Create package structure

**2. Data models**
- [ ] `internal/model/task.go`: Task struct
- [ ] `internal/model/relationship.go`: Relationship struct
- [ ] Constants and helper functions

**3. Storage layer**
- [ ] `internal/storage/interface.go`: BaseStorage interface
- [ ] `internal/storage/markdown.go`: MarkdownFileStorage implementation
- [ ] File parsing/writing with YAML frontmatter
- [ ] ID generation via globbing
- [ ] Collision detection

**4. Configuration**
- [ ] `internal/config/config.go`: Config loading
- [ ] TOML parsing with defaults
- [ ] Template resolution hierarchy
- [ ] Workflow validation

**5. Query engine**
- [ ] `internal/query/filters.go`: TaskFilter implementations
- [ ] `internal/query/engine.go`: QueryEngine struct and methods
- [ ] Convenience methods (FindByID, FindChildren, etc.)

**6. CLI commands**
- [ ] `cmd/root.go`: Root command
- [ ] `cmd/task.go`: Task subcommands (new, list, show, update, delete)
- [ ] `cmd/project.go`: Project management
- [ ] `cmd/config.go`: Config commands
- [ ] Viper integration for config loading

**7. Testing**
- [ ] Unit tests for storage
- [ ] Unit tests for query engine
- [ ] Unit tests for config loading
- [ ] Integration tests for CLI

### Phase 2: Enhancements

- [ ] Task linking (create relationships via CLI)
- [ ] Template system implementation
- [ ] Better error handling
- [ ] Output formatting (table, JSON, YAML)

### Phase 3: Integration

- [ ] MCP (Multi-Client Proxy) server
- [ ] Project discovery
- [ ] Advanced queries

## Key Files to Know

**Documentation**:
- `README.md` - Project vision
- `DESIGN_SUMMARY.md` - Quick reference
- `CONTEXT_FOR_NEXT_SESSION.md` - This file
- `.tasks/design/1.epic.md` - Design epic overview

**Specifications** (detailed, with code):
- `.tasks/design/1.research.md` - Data models
- `.tasks/design/2.research.md` - Storage interface & implementation
- `.tasks/design/3.research.md` - Project structure & config

**Design Rationale**:
- `.tasks/design/1.story.md` - Task model
- `.tasks/design/2.story.md` - ID system
- `.tasks/design/3.story.md` - Storage design
- `.tasks/design/4.story.md` - Config system
- `.tasks/design/5.story.md` - Query engine
- `.tasks/design/6.story.md` - CLI architecture

## Git Workflow

```bash
# Create isolated worktree for implementation
git worktree add ../opentask-impl

# In isolated worktree
cd ../opentask-impl
git checkout -b feat/core-implementation

# Make changes, commit, push
# Create PR when ready

# Back in main worktree to review/merge
cd ../opentask
```

## Dependencies

From design docs, will need:
- `github.com/spf13/cobra` - CLI framework
- `github.com/spf13/viper` - Config management
- `github.com/BurntSushi/toml` - TOML parsing
- `gopkg.in/yaml.v3` - YAML parsing
- Standard library: `os`, `io`, `time`, `context`, `regexp`, `path/filepath`

## Known Decisions

| What | Decision | Reason |
|------|----------|--------|
| Task IDs | Semantic sequential per-type (s-1, s-2, ...) | Human-readable, no state needed |
| Storage | Pluggable interface with FS as default | Flexibility, testability |
| Config | Optional with sensible defaults | Works without setup |
| Relationships | Slice of Relationship structs | Single source of truth |
| Query filters | Functional options pattern | Composable, extensible |

## Implementation Tips

1. **Start with tests**: Write test files first (table-driven tests work well)
2. **Use errors as values**: Define specific error types in each package
3. **Keep storage simple**: Focus on correctness, performance optimizations later
4. **Dog-food the system**: Use `.tasks/` for tracking implementation work
5. **Type safety**: Use constants for task types, relationship types, etc.

## Common Questions

**Q: Where should I start?**
A: Start with `internal/model/` package. Define Task and Relationship types. Then build storage around them.

**Q: How do I handle project paths?**
A: See `.tasks/design/3.research.md` for complete path resolution logic. Key: StorageConfig.Path scopes each storage instance to one project.

**Q: What about validation?**
A: Keep storage focused on persistence. Higher layers (QueryEngine, CLI) handle business logic validation.

**Q: Should I implement all backends?**
A: No. Markdown-FS is enough for MVP. Design allows adding SQLite/others later.

## Testing Strategy

- **Unit tests**: Test each package independently (mocks for dependencies)
- **Integration tests**: Test CLI commands end-to-end
- **Test fixtures**: Create `.tasks/test/` with sample task files
- **Acceptance tests**: Could use actual `.tasks/design/` tasks

## Success Criteria for MVP

- [ ] All core types implemented and tested
- [ ] MarkdownFileStorage fully working
- [ ] Can create/read/update/delete tasks via API
- [ ] CLI can list and create tasks
- [ ] Config system loads and applies defaults
- [ ] Relationship links stored and retrievable
- [ ] Comprehensive README for using the system

## Contact/Questions

All design decisions documented in `.tasks/design/`. If something is unclear:
1. Check the relevant research doc (r-*)
2. Check the design story (s-*) for rationale
3. Review git commit messages for context
