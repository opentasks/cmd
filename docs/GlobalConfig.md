# Global Configuration Guide

This guide explains how to set up and use global configuration for managing multiple projects with opentask.

## What is Global Configuration?

Global configuration is a single configuration file stored at `~/.config/opentask/config.toml` that:

- Defines default workflow and template settings used by all projects
- Lists all your known projects with their storage locations
- Tracks which project is currently active
- Provides a central place to manage project-level settings

## Why Use Global Configuration?

Global config is useful when:

- **Multiple projects**: You have 3+ projects and want to manage them from one place
- **Shared workflow**: Multiple projects use the same task statuses and transitions
- **Consistent templates**: You want the same task templates across projects
- **Quick switching**: You frequently switch between different projects
- **Monorepo with subprojects**: You have a large repo with multiple logical projects

## Getting Started

### Step 1: Create the Directory

```bash
mkdir -p ~/.config/opentask
```

### Step 2: Create Global Config

Create `~/.config/opentask/config.toml`:

```toml
# Global opentask configuration
# Defines default settings and projects

[global]
active_project = "work"

# Default workflow for all projects
[workflow]
statuses = ["todo", "in-progress", "reviewing", "done", "archived"]
initial = "todo"

[[workflow.transitions]]
from = "todo"
to = ["in-progress", "archived"]

[[workflow.transitions]]
from = "in-progress"
to = ["reviewing", "todo", "archived"]

[[workflow.transitions]]
from = "reviewing"
to = ["done", "in-progress", "archived"]

[[workflow.transitions]]
from = "done"
to = ["archived"]

# Default templates for all projects
[templates]
epic = "~/.local/share/opentask/templates/epic.md"
plan = "~/.local/share/opentask/templates/plan.md"
task = "~/.local/share/opentask/templates/task.md"

# Define your projects
[[global.projects]]
id = "work"
name = "Work Tasks"
description = "Tasks for my day job"

[global.projects.storage]
backend = "markdown-fs"
path = "~/work/.tasks"

[[global.projects]]
id = "personal"
name = "Personal Projects"
description = "My personal projects and goals"

[global.projects.storage]
backend = "markdown-fs"
path = "~/personal/.tasks"

[[global.projects]]
id = "opensource"
name = "Open Source"
description = "Open source projects I contribute to"

[global.projects.storage]
backend = "markdown-fs"
path = "~/projects/opensource/.tasks"
```

### Step 3: Initialize Project Directories

For each project, create `.opentask.toml` in the project directory:

```bash
cd ~/work
opentask config init
```

This creates a project config with the correct structure.

## Configuration Options

### [global] Section

Required section that defines global settings.

| Field | Type | Description |
|-------|------|-------------|
| `active_project` | string | Currently selected project ID |

### [[global.projects]] Section

Array of project definitions. Each project needs:

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | ✓ | Unique project identifier (slug-like: `work`, `my-project`, etc.) |
| `name` | string | ✗ | Human-readable project name |
| `description` | string | ✗ | Project description |
| `storage` | table | ✓ | Storage configuration |
| `storage.backend` | string | ✓ | Backend type (always `"markdown-fs"`) |
| `storage.path` | string | ✓ | Path to tasks directory (can use `~`) |

### [workflow] Section

Global default workflow settings. Applied to all projects unless overridden.

| Field | Type | Description |
|-------|------|-------------|
| `statuses` | string[] | Available task statuses |
| `initial` | string | Initial status for new tasks |
| `transitions` | table[] | Allowed state transitions |

### [templates] Section

Global default template paths. Applied to all projects unless overridden.

| Field | Type | Description |
|-------|------|-------------|
| `epic` | string | Path to epic template |
| `plan` | string | Path to plan template |
| `story` | string | Path to story template |
| `research` | string | Path to research template |
| `task` | string | Path to task template |
| `decision` | string | Path to decision template |

## Managing Projects

### List All Projects

```bash
opentask config projects
```

Output:
```
Configured projects:

* Work Tasks (work)
  Path: ~/work/.tasks
  Personal Projects (personal)
  Path: ~/personal/.tasks
  Open Source (opensource)
  Path: ~/projects/opensource/.tasks

Active project: work
```

The asterisk (*) indicates the currently active project.

### Switch Active Project

```bash
opentask config projects --active personal
```

This updates the `active_project` field in your global config.

## Project-Specific Overrides

Each project can have its own `.opentask.toml` to override global settings:

```toml
# ~/work/.opentask.toml

[project.project]
name = "Work Tasks"
owner = "John Doe"

# Override workflow for this project
[project.workflow]
statuses = ["backlog", "todo", "in-progress", "review", "done"]
initial = "todo"

[[project.workflow.transitions]]
from = "backlog"
to = ["todo"]

[[project.workflow.transitions]]
from = "todo"
to = ["in-progress"]

[[project.workflow.transitions]]
from = "in-progress"
to = ["review"]

[[project.workflow.transitions]]
from = "review"
to = ["done"]

# Optional: Override storage location
[project.storage]
backend = "markdown-fs"
path = "~/work/.tasks"
```

When you work in the `~/work` directory:
1. The project config is loaded
2. Settings from `[project.*]` override global settings
3. Settings not in project config inherit from global
4. Active project is matched from global projects list

## Examples

### Small Single-User Setup

```toml
# Simple global config for personal projects
[global]
active_project = "personal"

[workflow]
statuses = ["todo", "in-progress", "done"]
initial = "todo"

[[workflow.transitions]]
from = "todo"
to = ["in-progress"]

[[workflow.transitions]]
from = "in-progress"
to = ["done"]

[[global.projects]]
id = "personal"
name = "Personal"

[global.projects.storage]
path = "~/tasks/.tasks"
```

### Team Project Setup

```toml
# Global config for team projects
[global]
active_project = "main"

[workflow]
statuses = ["backlog", "todo", "in-progress", "review", "testing", "done"]
initial = "backlog"

# ... transitions ...

[templates]
epic = "~/.local/share/opentask/templates/epic.md"
task = "~/.local/share/opentask/templates/task.md"

[[global.projects]]
id = "main"
name = "Main Project"
[global.projects.storage]
path = "~/projects/main/.tasks"

[[global.projects]]
id = "infrastructure"
name = "Infrastructure"
[global.projects.storage]
path = "~/projects/infrastructure/.tasks"

[[global.projects]]
id = "documentation"
name = "Documentation"
[global.projects.storage]
path = "~/projects/documentation/.tasks"
```

### Monorepo Structure

```toml
# Global config for monorepo with multiple packages
[global]
active_project = "packages"

[workflow]
statuses = ["todo", "in-progress", "done"]
initial = "todo"

[[global.projects]]
id = "packages"
name = "Monorepo Packages"
[global.projects.storage]
path = "~/monorepo/.tasks"

[[global.projects]]
id = "backend"
name = "Backend Services"
[global.projects.storage]
path = "~/monorepo/services/.tasks"

[[global.projects]]
id = "frontend"
name = "Frontend Apps"
[global.projects.storage]
path = "~/monorepo/apps/.tasks"
```

Each directory also has its own `.opentask.toml`:
- `~/monorepo/.opentask.toml` - Shared settings
- `~/monorepo/services/.opentask.toml` - Backend-specific overrides
- `~/monorepo/apps/.opentask.toml` - Frontend-specific overrides

## Tips

- **Use descriptive project IDs**: Use slugs like `my-project`, `team-work`, not `proj1`, `proj2`
- **Path expansion**: Use `~` in paths - it's automatically expanded to your home directory
- **Relative project paths**: Each project config's storage path is relative to that config file
- **View active project**: Run `opentask config projects` to see which project is active
- **Test resolution**: Run `opentask config view` from any project directory to see what settings will be used
- **Consistency**: If multiple projects share settings, define them in global config rather than duplicating

## Troubleshooting

### "No global configuration found"

**Solution**: Create `~/.config/opentask/config.toml` with at least:
```toml
[global]
active_project = "myproject"

[[global.projects]]
id = "myproject"
name = "My Project"
[global.projects.storage]
path = "~/myproject/.tasks"
```

### "Project not found in global config"

**Solution**: Check that:
1. The project ID in `.opentask.toml` matches a `[[global.projects]] id`
2. The ID is spelled correctly (case-sensitive)
3. Run `opentask config projects` to list available project IDs

### Settings not being applied

**Solution**: Run `opentask config view` to see the merged configuration and which files were discovered. This helps you understand the merging order.

### Can't switch projects

**Solution**: The `--active` flag currently only shows the option. To permanently change the active project, edit `~/.config/opentask/config.toml` and update the `active_project` field manually.
