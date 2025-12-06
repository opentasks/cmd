# opentask Quick Start

## Building

```bash
go build -o opentask ./cmd/opentask
```

The binary will be available as `./opentask`

## Basic Usage

### Initialize a New Project

```bash
mkdir my_project
cd my_project
opentask config init --name "My Project"
```

This creates `.opentask.toml` in the current directory. The project ID is automatically derived from the directory name ("my_project").

**What happens:**
- Creates `.opentask.toml` configuration
- Sets up `.tasks/` directory for task storage
- Project is immediately ready to use

### Create Your First Task

```bash
opentask task new "My Epic" --type epic
```

This creates task ID 1. Tasks are stored as markdown files with YAML frontmatter in `.tasks/`.

### Create Subtasks

```bash
opentask task new "Plan the Work" --type plan --parent 1
opentask task new "Research Design" --type research --parent 1
opentask task new "Write Code" --type story --parent 1 --tag feature
```

### List Tasks

```bash
# List all tasks
opentask task list

# List by type
opentask task list --type story

# List by status
opentask task list --status todo

# List by parent epic
opentask task list --parent 1

# List by tag
opentask task list --tag feature

# Combine filters
opentask task list --type story --status in-progress --parent 1
```

### Show Task Details

```bash
opentask task show 3
```

### Update Task Status

```bash
opentask task update 3 --status in-progress
opentask task update 3 --status done
```

### Delete Tasks

```bash
opentask task delete 3
```

## Task Types

All tasks require a type. Valid types are:

- `epic` - Large bodies of work (usually 100+ points)
- `plan` - Planning and roadmap items
- `research` - Research and spike tasks
- `story` - User stories (typical development unit)
- `decision` - Architecture/design decisions
- `task` - Small implementation tasks

## Task Structure

A task file looks like this:

```markdown
---
id: 42
title: My Story
type: story
status: in-progress
tags: [feature, backend, urgent]
relationships:
  - type: parent
    taskID: 5
createdAt: 2025-11-02T10:00:00Z
updatedAt: 2025-11-02T10:30:00Z
---

# My Story

Markdown content here...

You can write detailed descriptions, acceptance criteria,
implementation notes, etc.
```

## Statuses

Default statuses are:
- `todo` (initial status)
- `in-progress`
- `reviewing`
- `done`
- `archived`

You can customize these in `config.toml`.

## Configuration

The `.opentask.toml` file was created by `config init`. You can edit it to customize:

```toml
[project]
id = "my-project"          # Optional: explicit project ID (defaults to directory name)
name = "My Project"
description = "Project description"
owner = "my-team"

[storage]
backend = "markdown-fs"
path = ".tasks"

[workflow]
statuses = ["todo", "in-progress", "review", "done"]
initial = "todo"

[[workflow.transitions]]
from = "todo"
to = ["in-progress", "archived"]

[[workflow.transitions]]
from = "in-progress"
to = ["review", "done", "archived"]
```

**Pro tip**: If you don't specify `[project] id`, opentask uses the directory name. This makes projects portable!

## Multi-Project Setup

For managing multiple projects, use global configuration with context paths. This lets you work on different projects without switching manually.

### Setup (One-time)

```bash
# Create global config
mkdir -p ~/.config/opentask
cat > ~/.config/opentask/config.toml << 'TOML'
[[projects]]
id = "personal"
name = "Personal Projects"

[[projects.context]]
path = "/home/user/personal/blog"

[[projects.context]]
path = "/home/user/personal/side-projects"

[projects.storage]
path = "/home/user/Notes/.tasks"

[[projects]]
id = "work"
name = "Work Tasks"

[[projects.context]]
path = "/home/user/work/company-repo"

[[projects.context]]
path = "/home/user/work/client-projects"

[projects.storage]
path = "/home/user/work/.tasks"
TOML
```

**What this does:**
- Defines two projects: "personal" and "work"
- Maps multiple directories to each project via context paths
- When you `cd` into any context directory, opentask automatically uses that project!

### Daily Usage

```bash
# Navigate to work directory
cd /home/user/work/company-repo/src/api
opentask task list      # Automatically uses "work" project (matched by context)

# Navigate to personal project
cd /home/user/personal/blog
opentask task list      # Automatically uses "personal" project

# View all projects and their contexts
opentask project list
```

**No manual switching needed!** The active project is derived from your current directory.

### Adding More Context Paths

If you clone a new repository or create a new directory for an existing project:

```bash
cd /home/user/work/new-microservice
opentask project attach work  # Adds this directory to "work" project contexts
```

Now this directory resolves to the "work" project automatically!

See [project-contexts.md](project-contexts.md) for detailed documentation.

## Command Line Options

```bash
# Specify project path (legacy, use project contexts instead)
opentask --path /path/to/project task list

# Specify project by ID (when using project contexts)
opentask --project work task list

# Specify config file
opentask --config /path/to/config.toml task list

# Enable verbose output
opentask --verbose task list
```

## File Organization

Tasks are organized hierarchically:

```
project-root/
├── config.toml                 # Optional config
├── e-1-my-epic.md             # Epic at root
├── 1-my-epic/                  # Subtasks in epic directory
│   ├── p-2-planning.md
│   ├── r-3-research.md
│   └── s-4-implementation.md
├── e-5-second-epic.md
├── 5-second-epic/
│   └── s-6-task.md
└── templates/                  # Optional local templates
    ├── epic.md
    └── story.md
```

## Global Flags

All commands support these flags:

- `--path PATH` - Project path (default: current directory)
- `--config PATH` - Config file path
- `--verbose` - Enable verbose output
- `--help` - Show help

## Tips

1. **Dog-food the system**: Use `.tasks/` in your project to track opentask development itself
2. **Use tags**: Tags help organize related work across epics
3. **Keep descriptions concise**: YAML frontmatter is metadata, markdown content is for details
4. **Commit regularly**: Tasks are files - commit them with your code
5. **Use statuses meaningfully**: They flow through your workflow, don't arbitrary jump states

## Examples

### Create a user story with acceptance criteria

```bash
opentask task new "Display user profile" --type story --parent 5 --tag ui

# Then edit the file to add details:
# ---
# id: 42
# ...
# ---
# # Display User Profile
#
# As a user, I want to see my profile page.
#
# ## Acceptance Criteria
# - [ ] Name and email are displayed
# - [ ] Avatar is shown
# - [ ] Profile is editable
```

### Track research and implementation together

```bash
opentask task new "User Authentication" --type epic
opentask task new "Research auth strategies" --type research --parent 1
opentask task new "Plan auth system" --type plan --parent 1
opentask task new "Implement OAuth" --type story --parent 1
opentask task new "Add password reset" --type story --parent 1
```

### Filter work by status

```bash
# Show all in-progress work
opentask task list --status in-progress

# Show all done work this epic
opentask task list --parent 5 --status done

# Show all pending stories
opentask task list --type story --status todo
```

## Next Steps

- Read [design-summary.md](design-summary.md) for architecture details
- Run tests with `go test ./...`
- Customize `config.toml` for your workflow
- Create project templates in `templates/`

## Onboarding Flow

If you try to run a task command in a directory without a project, you'll see an onboarding guide:

```
╭────────────────────────────────────────────────────────╮
│  NO OPENTASK PROJECT FOUND                            │
│                                                       │
│  Current directory: /home/user/random-dir            │
│                                                       │
│  No project configuration found for this location.    │
│                                                       │
│  Choose an option:                                    │
│                                                       │
│  1. Initialize a new project here                     │
│     $ opentask config init                           │
│                                                       │
│  2. Attach this directory to an existing project      │
│     $ opentask project attach <project-id>           │
│     $ opentask project list  # see available projects│
│                                                       │
│  3. Work in a directory with an existing project      │
│     $ cd /path/to/existing/project                   │
╰────────────────────────────────────────────────────────╯
```

**Option 1**: Create a new project in this directory
```bash
opentask config init --name "My Project"
```

**Option 2**: Link this directory to an existing project
```bash
opentask project list            # See available projects
opentask project attach work     # Attach to "work" project
```

**Option 3**: Navigate to a directory that already has a project
```bash
cd /path/to/existing/project
opentask task list
```

## How Projects Are Resolved

Opentask determines your active project using a clear 3-tier priority:

1. **Explicit ID** - If `.opentask.toml` has `[project] id = "foo"`, that wins
2. **Context Match** - Matches current directory against global config context paths
3. **Directory Name** - If `.opentask.toml` exists but no explicit ID, uses directory name

See [config.md](config.md#active-project-resolution) for detailed explanation.

## Troubleshooting

### "No active project found"

This means opentask couldn't determine which project you're working on.

**Solutions:**
1. Run `opentask config init` to create a new project
2. Run `opentask project attach <id>` to link this directory to an existing project
3. Navigate to a directory with `.opentask.toml`

### "task not found"

Task ID doesn't exist in the current project. Use `opentask task list` to see all tasks.

### "invalid task type"

Use one of the valid types: `epic`, `plan`, `research`, `story`, `decision`, `task`

### Config not loading

Check that `.opentask.toml` exists and is valid TOML syntax. Run `opentask config view` to see what was loaded.
