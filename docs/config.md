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

**Active project resolution:**

The active project is **always derived at runtime** - never persisted to disk. Resolution follows a clear 3-tier priority:

1. **Priority 1: Explicit ID** - If local `.opentask.toml` has `[project] id = "foo"`, that wins
2. **Priority 2: Context Match** - Matches current directory against global config `[[projects.context]]` paths (longest match wins)
3. **Priority 3: Directory Name** - If `.opentask.toml` exists but has no explicit ID, uses directory name (e.g., "module-a")

If none of these resolve, you'll see an onboarding guide to initialize or attach a project.

See "Active Project Resolution" section below for details.

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

# Define projects with context paths
[[projects]]
id = "work"
name = "Work Tasks"

[[projects.context]]
path = "~/work/company-repo"

[[projects.context]]
path = "~/work/client-projects"

[projects.storage]
backend = "markdown-fs"
path = "~/work/.tasks"

[[projects]]
id = "personal"
name = "Personal Tasks"

[[projects.context]]
path = "~/personal/blog"

[[projects.context]]
path = "~/personal/side-projects"

[projects.storage]
backend = "markdown-fs"
path = "~/personal/.tasks"
```

## Active Project Resolution

Opentask determines which project you're working on by analyzing your current directory and configuration files. The active project is **always derived at runtime** - never stored or persisted.

### Resolution Algorithm

The system follows a clear 3-tier priority:

```
ActiveProject = f(cwd, local_config, global_config)
```

#### Priority 1: Explicit Project ID (Highest)

If your local `.opentask.toml` specifies an explicit `[project] id`, that always wins:

```toml
# /home/user/myapp/.opentask.toml
[project]
id = "myapp-v2"  # ← Explicit, overrides everything
name = "My Application"
```

**Use case**: When you want a specific project ID that differs from the directory name.

#### Priority 2: Global Config Context Match

Matches your current directory against `[[projects.context]]` paths in global config. Longest matching path wins:

```toml
# ~/.config/opentask/config.toml
[[projects]]
id = "client-work"
name = "Client Project"

[[projects.context]]
path = "/home/user/work/client-repo"

[[projects.context]]
path = "/home/user/work/client-frontend"
```

```bash
cd /home/user/work/client-repo
opentask task list
# → ActiveProject = "client-work" (matched by context)
```

**Use case**: Multiple directories share the same project (monorepo, multiple clones).

#### Priority 3: Directory Name Fallback

If `.opentask.toml` exists but has no explicit `[project] id`, uses the directory name:

```bash
cd /home/user/projects/awesome-app
ls .opentask.toml  # exists, no explicit id
opentask task list
# → ActiveProject = "awesome-app" (from directory)
```

**Use case**: Simple single-directory projects.

### Onboarding Flow

If no project can be resolved (no `.opentask.toml`, no context match), you'll see an onboarding guide:

```
╭────────────────────────────────────────────────────────╮
│  NO OPENTASK PROJECT FOUND                            │
│                                                       │
│  Current directory: /home/user/random-dir            │
│                                                       │
│  Choose an option:                                    │
│                                                       │
│  1. Initialize a new project here                     │
│     $ opentask config init                           │
│                                                       │
│  2. Attach this directory to an existing project      │
│     $ opentask project attach <project-id>           │
│                                                       │
│  3. Work in a directory with an existing project      │
│     $ cd /path/to/existing/project                   │
╰────────────────────────────────────────────────────────╯
```

**Commands that bypass onboarding:**
- `opentask config init` - Creates local config
- `opentask project list` - Shows available projects
- `opentask project attach` - Adds context to existing project
- `opentask --help` - Help/docs

**Commands that require a resolved project:**
- `opentask task *` - All task operations
- `opentask config view` - Needs to know which config to show

### Examples

#### Example 1: Explicit ID Override

```toml
# /home/user/my-project/.opentask.toml
[project]
id = "prod-app"  # Overrides directory name "my-project"
name = "Production Application"
```

Result: `ActiveProject = "prod-app"`

#### Example 2: Multiple Clones via Context

```toml
# ~/.config/opentask/config.toml
[[projects]]
id = "webapp"
name = "Web Application"

[[projects.context]]
path = "/home/user/repos/webapp-main"

[[projects.context]]
path = "/home/user/repos/webapp-feature-x"

[[projects.context]]
path = "/home/user/repos/webapp-hotfix"
```

All three directories resolve to `ActiveProject = "webapp"` - they share tasks!

#### Example 3: Simple Directory Name

```bash
cd ~/projects/blog-site
ls .opentask.toml  # exists, no [project] id
opentask task list
# → ActiveProject = "blog-site"
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

#### `id` - Project Identifier

Optional string in `[project]` section. Explicitly sets the project ID. If not specified, defaults to directory name.

**Example:**
```toml
[project]
id = "my-app-v2"  # Explicit ID (Priority 1)
name = "My Application"
```

Without explicit ID, the project ID is derived from the directory name.

### Global Config Structure

Global configs use the `[global.*]` and `[workflow]`, `[templates]` namespaces:

#### `[global]` - Global Settings

**Note**: The `active_project` field is **deprecated** and no longer used. Active project is now derived at runtime (see "Active Project Resolution" below).

#### `[[projects]]` - Project Definitions

Array of project definitions available globally.

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | **Required.** Unique project identifier |
| `name` | string | Project display name |
| `context` | array | Array of `[[projects.context]]` tables with `path` fields |
| `storage` | table | Storage configuration for this project |
| `workflow` | table | Optional project-specific workflow override |
| `templates` | table | Optional project-specific templates override |

**Context paths** allow multiple directories to resolve to the same project:

```toml
[[projects]]
id = "webapp"

[[projects.context]]
path = "/home/user/webapp-main"

[[projects.context]]
path = "/home/user/webapp-feature"
```

When you `cd` into either directory, `ActiveProject = "webapp"`.

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

- `opentask config init` - Initialize a new project config in current directory
- `opentask config view` - Show resolved configuration with discovery details
- `opentask project list` - List all projects in global config
- `opentask project attach <project-id>` - Attach current directory to an existing project (adds to context paths)
