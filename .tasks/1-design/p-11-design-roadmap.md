---
id: 11
title: Design Roadmap
type: plan
status: done
tags: [design, planning]
relationships:
  - type: parent
    taskID: 1

createdAt: 2025-11-02T08:00:00Z
updatedAt: 2025-11-02T08:00:00Z
---

# Design Roadmap

Overview of the design phase tasks and their dependencies.

## Design Phase Breakdown

1. **Task Data Model** (5): Define Task struct, relationships, metadata
2. **Semantic ID System** (6): Global sequential IDs with simple generation
3. **Storage Interface** (7): BaseStorage interface and responsibilities
4. **Config System** (8): config.toml structure and defaults
5. **Query Engine** (9): Simple filtering with functional options
6. **CLI Architecture** (10): Viper/Cobra integration

Each story documents the design decision, rationale, and implementation considerations.

---

## Phase 2 Roadmap (Implementation Priorities)

After core system completion, focus on:

### High Priority
1. **Unit Testing** - Test each package independently
2. **Integration Testing** - End-to-end CLI testing
3. **Task Linking CLI** - `task link` command for creating relationships interactively

### Medium Priority
4. **Template System** - Built-in templates for each task type, custom template support
5. **Output Formats** - JSON, YAML, CSV export options
6. **Error Handling** - Better error messages and recovery

### Low Priority
7. **Advanced Features** - Full-text search, dependency graphs, critical path analysis
8. **MCP Server** - Multi-Client Proxy integration for tool-based access
9. **Web UI** - Browser-based interface (future milestone)

### Completed Phases
- ✅ **Phase 0**: Research and Design (completed Nov 2, 2025)
- ✅ **Phase 1**: Core Implementation (completed Nov 2, 2025)
  - Data models and storage layer
  - Configuration system
  - Query engine
  - CLI with Cobra/Viper
  - File-based persistence with YAML frontmatter
  - All CRUD operations working

Current implementation tracked in `.tasks/1-phase-2-testing-polish/` epic.
