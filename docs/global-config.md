# Global Configuration Guide

This guide explains how to set up and use global configuration for managing multiple projects with opentask.

## What is Global Configuration?

Global configuration is a single configuration file stored at `${XDG_CONFIG_HOME}/opentask/config.toml` (typically `~/.config/opentask/config.toml`) that:

- Defines default workflow and template settings used by all projects
- Lists all your known projects with their storage locations
- Defines working directory contexts for automatic project detection
- Provides a central place to manage project-level settings

## Why Use Global Configuration?

Global config is useful when:

- **Multiple projects**: You have 3+ projects and want to manage them from one place
- **Shared workflow**: Multiple projects use the same task statuses and transitions
- **Consistent templates**: You want the same task templates across projects
- **Directory-based switching**: You want opentask to automatically detect which project based on your current working directory
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
[[projects]]
id = "work"
name = "Work Tasks"
description = "Tasks for my day job"

[projects.storage]
backend = "markdown-fs"
path = "~/work/.tasks"

[[projects.context]]
path = "/home/user/work"

[[projects]]
id = "personal"
name = "Personal Projects"
description = "My personal projects and goals"

[projects.storage]
backend = "markdown-fs"
path = "~/personal/.tasks"

[[projects.context]]
path = "/home/user/personal"

[[projects]]
id = "opensource"
name = "Open Source"
description = "Open source projects I contribute to"

[projects.storage]
backend = "markdown-fs"
path = "~/projects/opensource/.tasks"

[[projects.context]]
path = "/home/user/projects/opensource"
```

### Step 3: Initialize Project Directories (Optional)

For each project, you can optionally create `.opentask.toml` to override global settings:

```bash
cd ~/work
cat > .opentask.toml << 'EOF'
[project]
name = "Work Tasks"
owner = "Your Name"

# Optional: Override workflow for this project only
[workflow]
statuses = ["backlog", "todo", "in-progress", "review", "done"]
initial = "todo"
EOF
```

## Configuration Reference

### [[projects]] Section

Array of project definitions. Each project needs:

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | ✓ | Unique project identifier (slug-like: `work`, `my-project`, etc.) |
| `name` | string | ✗ | Human-readable project name |
| `description` | string | ✗ | Project description |
| `storage` | table | ✓ | Storage configuration |
| `storage.backend` | string | ✓ | Backend type (currently `"markdown-fs"`) |
| `storage.path` | string | ✓ | Path to tasks directory (can use `~`) |

### [[projects.context]] Section

Array of working directory contexts for automatic project detection. When you run opentask from a directory, it matches the `cwd` against context paths to find the active project.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `path` | string | ✓ | Working directory path (can use `~`, expands to full path) |

**How context matching works:**
1. opentask checks your current working directory (cwd)
2. For each project, it checks if cwd matches any `[[projects.context]]` path
3. If a match is found, that project becomes active
4. If no match, opentask looks for `.opentask.toml` in cwd or parent directories
5. If still no match, returns an error (no active project found)

### [workflow] Section

Global default workflow settings. Applied to all projects unless overridden by project-level `.opentask.toml`.

| Field | Type | Description |
|-------|------|-------------|
| `statuses` | string[] | Available task statuses |
| `initial` | string | Initial status for new tasks |
| `transitions` | table[] | Allowed state transitions |

### [templates] Section

Global default template paths. Applied to all projects unless overridden by project-level config.

| Field | Type | Description |
|-------|------|-------------|
| `epic` | string | Path to epic template |
| `plan` | string | Path to plan template |
| `story` | string | Path to story template |
| `research` | string | Path to research template |
| `task` | string | Path to task template |
| `decision` | string | Path to decision template |

## Project Resolution Order

When you run opentask, it resolves the active project in this order:

1. **Explicit flag**: `--project work` or `opentask_PROJECT` env var
2. **Context match**: Current working directory matches a `[[projects.context]]` path
3. **Local config**: `.opentask.toml` found in cwd or parent directories
4. **Error**: No project found - you must specify one or define a matching context

## Examples

### Small Single-User Setup

```toml
# Simple global config for personal projects
[workflow]
statuses = ["todo", "in-progress", "done"]
initial = "todo"

[[workflow.transitions]]
from = "todo"
to = ["in-progress"]

[[workflow.transitions]]
from = "in-progress"
to = ["done"]

[[projects]]
id = "personal"
name = "Personal"

[projects.storage]
backend = "markdown-fs"
path = "~/.local/share/opentask/projects/personal/.tasks"

[[projects.context]]
path = "~/projects"
path = "~/documents"
```

### Team Project Setup with Directory Contexts

```toml
[workflow]
statuses = ["backlog", "todo", "in-progress", "review", "testing", "done"]
initial = "backlog"

[[workflow.transitions]]
from = "backlog"
to = ["todo"]

[[workflow.transitions]]
from = "todo"
to = ["in-progress"]

[[workflow.transitions]]
from = "in-progress"
to = ["review"]

[[workflow.transitions]]
from = "review"
to = ["testing", "in-progress"]

[[workflow.transitions]]
from = "testing"
to = ["done"]

[templates]
epic = "~/.local/share/opentask/templates/epic.md"
task = "~/.local/share/opentask/templates/task.md"

[[projects]]
id = "main"
name = "Main Project"

[projects.storage]
backend = "markdown-fs"
path = "~/projects/main/.tasks"

[[projects.context]]
path = "/home/user/projects/main"
path = "/home/user/work/main"

[[projects]]
id = "infrastructure"
name = "Infrastructure"

[projects.storage]
backend = "markdown-fs"
path = "~/projects/infrastructure/.tasks"

[[projects.context]]
path = "/home/user/projects/infrastructure"

[[projects]]
id = "documentation"
name = "Documentation"

[projects.storage]
backend = "markdown-fs"
path = "~/projects/documentation/.tasks"

[[projects.context]]
path = "/home/user/projects/documentation"
```

### Monorepo Structure with Subproject Contexts

```toml
[workflow]
statuses = ["todo", "in-progress", "done"]
initial = "todo"

[[workflow.transitions]]
from = "todo"
to = ["in-progress"]

[[workflow.transitions]]
from = "in-progress"
to = ["done"]

# Root monorepo project
[[projects]]
id = "monorepo"
name = "Monorepo Root"

[projects.storage]
backend = "markdown-fs"
path = "~/monorepo/.tasks"

[[projects.context]]
path = "/home/user/monorepo"

# Backend services
[[projects]]
id = "backend"
name = "Backend Services"

[projects.storage]
backend = "markdown-fs"
path = "~/monorepo/services/.tasks"

[[projects.context]]
path = "/home/user/monorepo/services"

# Frontend apps
[[projects]]
id = "frontend"
name = "Frontend Apps"

[projects.storage]
backend = "markdown-fs"
path = "~/monorepo/apps/.tasks"

[[projects.context]]
path = "/home/user/monorepo/apps"
```

With this setup:
- `cd ~/monorepo` → activates `monorepo` project
- `cd ~/monorepo/services` → activates `backend` project (context matches first)
- `cd ~/monorepo/apps` → activates `frontend` project (context matches first)

## Project-Specific Overrides

Each project can have its own `.opentask.toml` to override global settings:

```toml
# ~/work/.opentask.toml

[project]
name = "Work Tasks"
owner = "Your Name"

# Override workflow for this project
[workflow]
statuses = ["backlog", "todo", "in-progress", "review", "done"]
initial = "todo"

[[workflow.transitions]]
from = "backlog"
to = ["todo"]

[[workflow.transitions]]
from = "todo"
to = ["in-progress"]

[[workflow.transitions]]
from = "in-progress"
to = ["review"]

[[workflow.transitions]]
from = "review"
to = ["done"]

# Optional: Override storage location
[storage]
backend = "markdown-fs"
path = "~/work/.tasks"
```

When you work in the `~/work` directory:
1. Global config loads first (workflows, templates, projects list)
2. Local `.opentask.toml` overrides global settings
3. Settings not in local config inherit from global defaults
4. Active project is determined by context matching + local config

## Tips

- **Use descriptive project IDs**: Use slugs like `my-project`, `team-work`, not `proj1`, `proj2`
- **Path expansion**: Use `~` in paths - it's automatically expanded to your home directory
- **Multiple contexts per project**: You can specify multiple `[[projects.context]]` paths for the same project (useful for monorepos or shared workspaces)
- **Test resolution**: Run `opentask config view` from any project directory to see what settings will be used
- **Consistency**: If multiple projects share settings, define them in global config rather than duplicating in each `.opentask.toml`

## Troubleshooting

### "No global configuration found"

**Cause**: Global config file doesn't exist at `~/.config/opentask/config.toml`

**Solution**: Create the file with minimal configuration:
```toml
[workflow]
statuses = ["todo", "in-progress", "done"]
initial = "todo"

[[projects]]
id = "myproject"
name = "My Project"

[projects.storage]
backend = "markdown-fs"
path = "~/myproject/.tasks"

[[projects.context]]
path = "~/myproject"
```

### "Project not found"

**Cause**: Current working directory doesn't match any `[[projects.context]]` path, and no `.opentask.toml` found

**Solution**:
1. Check that your cwd matches a context path in global config
2. Or create `.opentask.toml` in your project directory
3. Or use `--project <id>` flag to explicitly specify the project

### Settings not being applied

**Cause**: Project-level `.opentask.toml` may be overriding global settings

**Solution**: Run `opentask config view` to see the merged configuration and understand which settings came from which files.

### Settings inheritance

**Order of precedence** (highest to lowest):
1. Command-line flags (`--project`, `--status`, etc.)
2. Project-level `.opentask.toml` settings
3. Global config `[workflow]`, `[templates]`
4. Built-in defaults
