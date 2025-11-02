# opentask Quick Start

## Building

```bash
go build -o opentask ./cmd/opentask
```

The binary will be available as `./opentask`

## Basic Usage

### Create a Project

```bash
mkdir my_project
cd my_project
```

### Create Your First Task

```bash
opentask task new "My Epic" --type epic
```

This creates task ID 1. Tasks are stored as markdown files with YAML frontmatter.

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

Create a `config.toml` in your project root to customize:

```toml
[project]
name = "My Project"
description = "Project description"
owner = "my-team"

[workflow]
statuses = ["todo", "in-progress", "review", "done"]
initial = "todo"

[[workflow.transitions]]
from = "todo"
to = ["in-progress", "archived"]

[[workflow.transitions]]
from = "in-progress"
to = ["review", "done", "archived"]

[storage]
backend = "markdown-fs"
path = "."
```

## Command Line Options

```bash
# Specify project path
opentask --path /path/to/project task list

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

- Read `DESIGN_SUMMARY.md` for architecture details
- Check `.tasks/design/` for detailed specifications
- Run tests with `go test ./...`
- Customize `config.toml` for your workflow
- Create project templates in `templates/`

## Troubleshooting

**"task not found"** - Task ID doesn't exist. Use `task list` to see all tasks.

**"invalid task type"** - Use one of: epic, plan, research, story, decision, task

**Files not being found** - Ensure you're in the right project directory or use `--path`

**Config not loading** - Check that `config.toml` exists and is valid TOML syntax
