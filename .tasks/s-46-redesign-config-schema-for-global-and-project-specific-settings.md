---
id: 46
title: Redesign config schema for global and project-specific settings
type: story
status: todo
tags: [config, schema, design, breaking-change]
relationships: []
createdAt: "2025-11-03T09:15:00Z"
updatedAt: "2025-11-03T09:15:00Z"
---

# Story: Redesign Config Schema for Global and Project-Specific Settings

## Problem

The current config schema doesn't distinguish between global user defaults and project-specific settings. This creates issues:

1. **Global settings mixed with project config** - Can't share defaults across projects
2. **Project-specific storage path in global config** - Global config must be copied per project
3. **No way to specify active project** - Can't correlate parent configs with actual project
4. **Can't have project-specific templates or workflows** - All settings are global

## Vision

### Schema Composition

Create a modular config schema built from three composable parts:

```
OpentaskConfigSchema = OpentaskConfigCoreSchema + OpentaskConfigGlobalSchema + OpentaskConfigProjectSchema
```

#### OpentaskConfigCoreSchema
Shared fields used in both global and project contexts:
- `[workflow]` - Task status definitions and transitions
- `[templates]` - Template paths for task types

#### OpentaskConfigGlobalSchema
Global user defaults:
- `active_project` - Currently selected project ID
- List of `[[projects]]` configurations

#### OpentaskConfigProjectSchema
Project-specific overrides:
- `[project]` - Project metadata (name, description, owner)
- `[storage]` - Storage backend and path (project-specific!)
- Project-specific `[workflow]` and `[templates]` overrides

### File Structures

#### Global Config: `~/.config/opentask/config.toml`

```toml
# Global settings applied to all projects
[global]
active_project = "opentask"  # Currently selected project

# Global workflow and templates
[workflow]
statuses = ["todo", "in-progress", "done"]
initial = "todo"

[templates]
epic = "~/.local/share/opentask/templates/epic.md"
task = "~/.local/share/opentask/templates/task.md"

# Define projects at global level
[[projects]]
id = "opentask"
name = "Opentask"
storage.backend = "markdown-fs"
storage.path = "~/Projects/opentask/.tasks"

[[projects]]
id = "my-notes"
name = "Personal Notes"
storage.backend = "markdown-fs"
storage.path = "~/Notes/.tasks"

[[projects]]
id = "work-tasks"
name = "Work Tasks"
storage.backend = "markdown-fs"
storage.path = "~/Work/opentask/.tasks"
```

#### Project Config: `.opentask.toml` in project root

```toml
# Project-specific metadata
[project]
name = "Opentask"
description = "Task management system"
owner = "zenobius"

# Project-specific workflow (optional - overrides global)
[workflow]
statuses = ["todo", "in-progress", "reviewing", "done", "archived"]
initial = "todo"

# Project-specific templates (optional - overrides global)
[templates]
epic = "./_templates/epic.md"
decision = "./_templates/decision.md"

# Project-specific storage
[storage]
backend = "markdown-fs"
path = "./.tasks"

# Active project tracking (correlates with global config)
active_project = "opentask"
```

### Resolution Algorithm

When resolving config for a project:

1. **Discover** `.opentask.toml` files walking up from current directory
2. **Collect** all found config files (closest to furthest)
3. **Determine** which project from:
   - Explicit `--project` flag
   - `active_project` field in closest `.opentask.toml`
   - Fall back to default/first project
4. **Merge** in priority order (highest to lowest):
   - Current directory `.opentask.toml` (project schema fields)
   - Parent `.opentask.toml` files (project schema fields)
   - Global config `[[projects]]` matching project ID (project schema fields)
   - Global config `[global]` section (core schema fields)
   - Built-in defaults

5. **Result**: Single merged `OpentaskProjectConfig` with only project-relevant fields

### Examples

#### Example 1: Simple Global + Project Override

**Global config** (`~/.config/opentask/config.toml`):
```toml
[global]
active_project = "personal"

[workflow]
statuses = ["todo", "done"]
initial = "todo"

[[projects]]
id = "personal"
name = "Personal Projects"
storage.backend = "markdown-fs"
storage.path = "~/personal/.tasks"
```

**Project config** (`~/personal/.opentask.toml`):
```toml
[project]
name = "Personal"
owner = "me"

active_project = "personal"
```

**Result**: Merged config uses:
- Global workflow (statuses, initial)
- Project name "Personal"
- Storage path ~/personal/.tasks
- Global templates

#### Example 2: Nested Project with Specific Workflow

**Global config** (`~/.config/opentask/config.toml`):
```toml
[global]
active_project = "work"

[workflow]
statuses = ["todo", "done"]

[[projects]]
id = "work"
name = "Work"
storage.backend = "markdown-fs"
storage.path = "~/work/.tasks"
```

**Root project config** (`~/work/.opentask.toml`):
```toml
[project]
name = "Work Projects"
active_project = "work"

[workflow]
statuses = ["backlog", "todo", "in-progress", "review", "done"]
initial = "todo"
```

**Sub-project** (`~/work/project1/.opentask.toml`):
```toml
[project]
name = "Project 1"
owner = "team-a"

active_project = "work"

[templates]
task = "./templates/task.md"
```

**Result when in ~/work/project1**:
- Project name: "Project 1"
- Owner: "team-a"
- Workflow: From root (5 statuses)
- Templates: task.md from ./templates, others from global
- Storage: ~/work/.tasks (from global, project doesn't override)

## Implementation Plan

### Phase 1: Schema Definition

1. **Define new schema types** in `internal/config/config.go`:
   - `OpentaskConfigCoreSchema` struct (workflow, templates)
   - `OpentaskConfigGlobalSchema` struct (global, projects array)
   - `OpentaskConfigProjectSchema` struct (project, storage, active_project)
   - `OpentaskGlobalConfig` type (full global file schema)
   - `OpentaskProjectConfig` type (full project file schema)
   - `OpentaskResolvedConfig` type (final merged result, project schema only)

2. **Add type comments** explaining each schema's purpose

3. **Add factory functions**:
   - `NewGlobalConfig()` - Create with defaults
   - `NewProjectConfig()` - Create with defaults
   - `NewResolvedConfig()` - Create resolved config

### Phase 2: Loading and Merging

1. **Update `internal/config/discovery.go`**:
   - Return both `.opentask.toml` and `~/.config/opentask/config.toml` paths
   - Add metadata about each file (type: project or global)

2. **Create `internal/config/merge.go`** with functions:
   - `LoadGlobalConfig(path string)` - Load global config
   - `LoadProjectConfig(path string)` - Load project config
   - `ResolveProjectConfig(cwd string) (*OpentaskResolvedConfig, error)` - Full resolution

3. **Resolution logic**:
   - Accept optional `--project` flag to override
   - Read `active_project` from discovered files
   - Load matching project from global config
   - Merge in correct order
   - Return `OpentaskResolvedConfig` with project-only fields

4. **Validation**:
   - Validate project ID exists when referenced
   - Validate workflow transitions
   - Ensure storage path is set (from project or global)

### Phase 3: CLI Updates

1. **Update `cmd/root.go`**:
   - Pass `--project` flag through to resolution
   - Use new `ResolveProjectConfig()` function

2. **Add `opentask config projects` command**:
   - List all projects from global config
   - Show which is active
   - Allow switching active project

3. **Update `opentask config view`**:
   - Show which project is active
   - Show project-specific overrides
   - Show what's inherited from global
   - Clear separation in output

4. **Update `opentask config init`**:
   - Ask for project ID
   - Set `active_project` field
   - Create starter templates (optional)

### Phase 4: Tests

1. **Unit tests** for schema validation:
   - Load global config files
   - Load project config files
   - Schema validation

2. **Integration tests** for resolution:
   - Global only
   - Project only
   - Global + single project level
   - Global + multiple project levels
   - Project overrides global
   - Explicit --project flag overrides

3. **Edge cases**:
   - Missing project in global config
   - Project ID mismatch
   - Missing storage path
   - Invalid workflow transitions

### Phase 5: Documentation

1. **Update `docs/Config.md`**:
   - Explain global vs project config
   - Show both file structures
   - Explain resolution order
   - Examples of common scenarios

2. **Add `docs/GlobalConfig.md`**:
   - How to set up global defaults
   - Managing multiple projects
   - Switching active project

3. **Update README**:
   - Mention multi-project support
   - Link to config docs

## File Changes Summary

### New Files
- `internal/config/merge.go` - Resolution and merging logic

### Modified Files
- `internal/config/config.go` - New schema types
- `internal/config/discovery.go` - Distinguish global vs project files
- `cmd/root.go` - Use new resolution logic
- `cmd/config.go` - Add projects command, update commands
- `docs/Config.md` - Updated schema docs
- (NEW) `docs/GlobalConfig.md` - Global config guide

### Tests to Add
- `internal/config/merge_test.go` - Resolution tests
- Updates to existing test files

## Success Criteria

- [ ] Global and project configs load independently
- [ ] Project config resolution follows merge order correctly
- [ ] Project ID routing works via active_project field
- [ ] Can have multiple projects in global config
- [ ] Project-specific workflows/templates override global
- [ ] Storage path is always resolved correctly
- [ ] `opentask config projects` shows all projects
- [ ] `opentask config view` shows resolution details
- [ ] All tests pass
- [ ] Documentation is clear and complete

## Breaking Changes

⚠️ **This is a breaking change:**
- Existing `.opentask.toml` files need `active_project` field (or use explicit `--project`)
- Global config structure completely changes
- Users with global config must migrate

## Migration Path

For existing users:
1. Backup `~/.config/opentask/config.toml` if it exists
2. Run `opentask config init --global` to create new global structure
3. Manually migrate any project definitions to global config
4. Add `active_project` to project `.opentask.toml` files

For new users:
- `opentask config init --global` creates global config
- `opentask config init` in project creates project config with project ID

## Related Decisions

- Decision 43: Config discovery patterns (will be updated)
- Task 44: CLI editing features (useful for managing projects)

## References

- Current config structure: `internal/config/config.go`
- Discovery logic: `internal/config/discovery.go`
- Config docs: `docs/Config.md`
