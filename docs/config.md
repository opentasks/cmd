# Opentasks Configuration

Configuration files are **optional**. If none are found, sensible defaults are used. Opentask supports both **global** and **project-specific** configurations:

- **Global config** (`~/.config/opentask/config.toml`): User-level defaults and project definitions
- **Project config** (`.opentask.toml`): Project-specific settings and overrides

This separation allows you to:
- Share workflow and template defaults across multiple projects
- Define different projects at the global level
- Override settings at the project level
- Support multi-project development with hierarchical configs

## Config Discovery and Resolution

The resolution process searches for `.opentask.toml` files starting from your current working directory and walking up the directory tree. It also checks for a global config at `~/.config/opentask/config.toml`.

**Resolution order (highest to lowest priority):**
1. `.opentask.toml` in current directory (project schema)
2. `.opentask.toml` in parent directories (project schema, walking up)
3. Global config `[[projects]]` matching the active project (global schema)
4. Global config `[global]` section (shared core settings)
5. Built-in defaults

This means:
- Project configs closest to you override those further away
- Project settings override global settings
- Global defaults apply when nothing else is specified

### Discovery and Merging Example

**Directory structure:**
```
~/.config/opentask/
└── config.toml                # Global config (defines projects & defaults)

~/projects/
├── work/
│   ├── .opentask.toml        # Project root config
│   ├── module-a/
│   │   ├── .opentask.toml    # Sub-project config
│   │   └── src/              # Current working directory
│   └── module-b/
│       └── .opentask.toml
└── personal/
    └── .opentask.toml
```

**When running from `~/projects/work/module-a/src/`:**

1. **Discovery finds** (in order of precedence):
   - `~/projects/work/module-a/.opentask.toml` (1st - project schema)
   - `~/projects/work/.opentask.toml` (2nd - project schema)
   - Global config's `[[projects]]` matching active project (3rd - global schema)
   - Global config's `[global]` section (4th - core schema)
   - Built-in defaults (5th)

2. **Merging result**: 
   - Module-A's config overrides Work's config
   - Work's config overrides global project definition
   - Global settings apply where not overridden
   - Final result: One merged `OpentaskResolvedConfig` with all settings

**Active project selection:**
- If `.opentask.toml` specifies `active_project`, that's used
- Otherwise, derived from directory name ("module-a") or global projects list
- Can override with `--project` flag

## Configuration Structure

There are two types of config files:

### Project Config (`.opentask.toml`)

Store project-specific settings that override global defaults.

**Minimal example:**
```toml
[project.project]
name = "My Project"

[project.storage]
backend = "markdown-fs"
path = ".tasks"
```

**Complete example with overrides:**
```toml
# Project metadata
[project.project]
name = "My Project"
description = "Project description"
owner = "your-name"

# Project-specific storage location
[project.storage]
backend = "markdown-fs"
path = ".tasks"

# Project-specific workflow (optional - overrides global)
[project.workflow]
statuses = ["todo", "in-progress", "reviewing", "done", "archived"]
initial = "todo"

[[project.workflow.transitions]]
from = "todo"
to = ["in-progress", "archived"]

[[project.workflow.transitions]]
from = "in-progress"
to = ["reviewing", "todo", "archived"]

[[project.workflow.transitions]]
from = "reviewing"
to = ["done", "in-progress", "archived"]

[[project.workflow.transitions]]
from = "done"
to = ["archived"]

# Project-specific templates (optional)
[project.templates]
epic = "templates/epic.md"
task = "templates/custom-task.md"
```

### Global Config (`~/.config/opentask/config.toml`)

Define user-level defaults and multiple projects.

**Example:**
```toml
# Global settings applied to all projects
[global]
active_project = "work"

# Shared workflow definition
[workflow]
statuses = ["todo", "in-progress", "reviewing", "done", "archived"]
initial = "todo"

[[workflow.transitions]]
from = "todo"
to = ["in-progress", "archived"]

# ... more transitions ...

# Shared templates
[templates]
epic = "~/.local/share/opentask/templates/epic.md"
plan = "~/.local/share/opentask/templates/plan.md"

# Define projects
[[global.projects]]
id = "work"
name = "Work Tasks"

[global.projects.storage]
backend = "markdown-fs"
path = "~/work/.tasks"

[[global.projects]]
id = "personal"
name = "Personal Tasks"

[global.projects.storage]
backend = "markdown-fs"
path = "~/personal/.tasks"
```

## Configuration Reference

### Project Config Structure

Project configs use the `[project.*]` namespace:

#### `[project.project]` - Project Metadata

All fields optional.

| Field | Type | Description |
|-------|------|-------------|
| `name` | string | Project name |
| `description` | string | Project description |
| `owner` | string | Project owner/maintainer |

#### `[project.storage]` - Storage Configuration

Location where tasks are stored.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `backend` | string | `"markdown-fs"` | Storage backend type (currently only `"markdown-fs"` supported) |
| `path` | string | `./.tasks` | Path where tasks are stored (relative to config file) |
| `options` | map | Empty | Backend-specific options |

#### `[project.workflow]` - Workflow Definition

Defines task statuses and transitions. Optional - uses global defaults if not specified.

| Field | Type | Description |
|-------|------|-------------|
| `statuses` | string[] | Available task statuses |
| `initial` | string | Initial status for new tasks |
| `transitions` | table[] | Allowed status transitions |

**Transition format:**
```toml
[[project.workflow.transitions]]
from = "todo"
to = ["in-progress", "archived"]
```

#### `[project.templates]` - Template Paths

Paths to custom task type templates. All optional.

| Field | Type | Description |
|-------|------|-------------|
| `epic` | string | Path to epic template |
| `plan` | string | Path to plan template |
| `story` | string | Path to story template |
| `research` | string | Path to research template |
| `task` | string | Path to task template |
| `decision` | string | Path to decision template |

Paths are resolved relative to the config file location.

#### `active_project` - Project Identifier

Optional string. If not specified, automatically derived from directory name or matched against global projects.

### Global Config Structure

Global configs use the `[global.*]` and `[workflow]`, `[templates]` namespaces:

#### `[global]` - Global Settings

| Field | Type | Description |
|-------|------|-------------|
| `active_project` | string | Currently selected project ID |

#### `[[global.projects]]` - Project Definitions

Array of project definitions available globally.

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | **Required.** Unique project identifier |
| `name` | string | Project display name |
| `storage` | table | Storage configuration for this project |
| `workflow` | table | Optional project-specific workflow override |
| `templates` | table | Optional project-specific templates override |

#### `[workflow]` - Global Workflow Defaults

Workflow configuration for `[project.workflow]` to inherit from. Same structure as project workflow.

#### `[templates]` - Global Template Defaults

Template paths for `[project.templates]` to inherit from. Same structure as project templates.

## Common Use Cases

### Single Project (No Global Config)

For simple single-project setups:
1. Create `.opentask.toml` in your project root
2. Define project metadata and storage path
3. Use global defaults for workflow

### Multi-Project Setup

For managing multiple projects:
1. Create `~/.config/opentask/config.toml` with global defaults
2. Define all projects in `[[global.projects]]`
3. Use `opentask config projects` to list and switch projects
4. Each project can have `.opentask.toml` for local overrides

### Monorepo with Shared Workflow

For monorepos with multiple packages/modules:
1. Create `.opentask.toml` at repo root with shared workflow
2. Create `.opentask.toml` in subdirectories to override project metadata
3. Storage paths can be relative to each config location
4. Subdirectories inherit workflow from root

### Project-Specific Workflow

For projects with custom task statuses:
1. Define global defaults in `~/.config/opentask/config.toml`
2. Override `[project.workflow]` in specific project's `.opentask.toml`
3. Project workflow takes precedence over global

## Tips

- **Relative paths**: All paths (`storage.path`, `templates.*`) are resolved relative to the config file's directory
- **Path expansion**: Tilde (`~`) in paths is expanded to your home directory
- **Auto-detection**: `active_project` is automatically derived if not specified - matches directory name or global projects
- **Hierarchical projects**: Subdirectories automatically discover parent configs and inherit settings
- **Override specific settings**: You only need to define sections/fields you want to override; others use defaults or inherit
- **Validate your config**: Run `opentask config view` to see the resolved configuration including all discovered files

## Commands

- `opentask config init` - Initialize a new project config
- `opentask config view` - Show resolved configuration with discovery details
- `opentask config projects` - List all projects in global config
- `opentask config projects --active [id]` - Switch active project (global config)
