# Opentasks Configuration

Configuration files are **optional**. If none are found, sensible defaults are used. Create an `opentask.toml` file to customize behavior for your project.

## Config Discovery

The CLI searches for `opentask.toml` files starting from your current working directory and walking up the directory tree. This means:

1. **Hierarchical resolution**: Configs are discovered from closest to furthest
2. **Merging**: All discovered configs are merged (closer/later configs override earlier ones)
3. **Stops at**: The filesystem root only
4. **Optional**: If no config is found, all defaults are applied

### Discovery Examples

```
/                                        # filesystem root
├── home/user/
│   ├── .config/opentask/config.toml     # global user config
│   ├── .local/share/opentask/
│   │   └── templates/
│   └── projects/
│       ├── .opentask.toml                # global projects config
│       └── myproject/
│           ├── .opentask.toml            # project-specific config
│           └── subproject/
│               ├── .opentask.toml        # subproject-specific config
│               └── workingdir/          # current working directory
```

When running from workingdir/:
- Searches: workingdir/ → subproject/ → myproject/ → projects/ → home/user/ → /
- Finds: .opentask.toml (subproject), .opentask.toml (myproject), .config/opentask/config.toml (user)
- Uses: merged config (subproject overrides myproject overrides user)                  (current working directory)
```

When running from special-tasks/:
- Searches: special-tasks/ → subproject/ → root/
- Finds: .opentask.toml (root)
- Uses: root config settings
```

## Configuration Structure

### Minimal Config (Storage Only)

```toml
[storage]
backend = "markdown-fs"  # Storage backend type
path = ".tasks"          # Path to task storage (relative to config file)
```

### Complete Config Example

```toml
[project]
name = "My Project"
description = "Project description"
owner = "your-name"

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

[templates]
epic = "templates/epic.md"
plan = "templates/plan.md"
story = "templates/story.md"
research = "templates/research.md"
task = "templates/task.md"
decision = "templates/decision.md"

[storage]
backend = "markdown-fs"
path = ".tasks"

[storage.options]
# Backend-specific options (empty for markdown-fs)
```

## Configuration Reference

### [project]

Project metadata. All fields are optional.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `name` | string | "" | Project name |
| `description` | string | "" | Project description |
| `owner` | string | "" | Project owner/maintainer |

### [workflow]

Defines task statuses and allowed transitions. Partially optional - can override defaults.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `statuses` | string[] | `["todo", "in-progress", "reviewing", "done", "archived"]` | Available task statuses |
| `initial` | string | `"todo"` | Initial status for new tasks |
| `transitions` | table[] | Default transitions | Allowed status transitions |

**Transition rules:**
- Each transition has a `from` status and array of `to` statuses
- If you override `transitions`, you must define all paths
- If you override `statuses`, they must all be valid in transitions

### [templates]

Paths to custom task type templates. All fields are optional.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `epic` | string | Built-in | Path to epic template |
| `plan` | string | Built-in | Path to plan template |
| `story` | string | Built-in | Path to story template |
| `research` | string | Built-in | Path to research template |
| `task` | string | Built-in | Path to task template |
| `decision` | string | Built-in | Path to decision template |

Paths are resolved relative to the config file location.

### [storage]

Storage backend configuration.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `backend` | string | `"markdown-fs"` | Storage backend type (currently only `"markdown-fs"` supported) |
| `path` | string | Current directory | Path where tasks are stored (relative to config file) |
| `options` | map | Empty | Backend-specific options |

## Tips

- **Relative paths**: All paths (`storage.path`, `templates.*`) are resolved relative to the config file's directory
- **Multi-level projects**: Put a config at the project root and subdirectories will automatically discover it
- **Override specific settings**: You only need to define sections/fields you want to override; others use defaults
- **Validate your config**: Run `opentask config show` to see the resolved configuration
