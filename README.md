<div align="center">
  <img src="docs/banner.png" alt="OpenTask Banner" style="border-radius: 12px; max-width: 600px;" />
</div>

# opentask

**Git-agnostic task management for developers and AI agents.**

Markdown-based tasks that live anywhere—in your projects, user directories, or CI pipelines. No database, no setup, no pollution of your Git repos. Perfect for solo developers, teams tracking work in code, and organizations planning AI agent integration.

## Why opentask?

- 📁 **Git-clean** – Store tasks outside project repos, keep them in sync anywhere
- 🔗 **Flexible paths** – XDG directories, environment variables, explicit paths—your choice
- 🤖 **AI-ready** ⏳ – Built for agent collaboration via cli.
- 📝 **Markdown native** – Edit in VS Code, Vim, or any editor
- ⚡ **Zero setup** – Works with sensible defaults, optional TOML configuration
- 🔌 **Pluggable storage** ⏳ – Swap backends (Planned: SQLite, DuckDB, PostgreSQL, cloud)

## Quick Start

```bash
# Build the CLI
mise build # or download a binary from the release tab

# Create an .opentask.toml in ~/Notes/MyProject using the markdown backend
opentask config init ~/Notes/MyProject

# Create your first task
./opentask task new "Build login page" "type:story"

# List all tasks
./opentask task list

# Filter by status or type
./opentask task list "status:todo AND type:story"

# Show task details
./opentask task show 1

# Update task status
./opentask task update 1 --status in-progress
```

## Core Features

### ✅ Ready Now (MVP)

**Flexible Storage**
- Default: File-based markdown storage (no database needed)
- Project discovery from multiple locations (in order of precedence):
  - Explicit paths: `--path` CLI flag or `opentask_PROJECT_PATH` env var
  - Config file: `.opentask.toml` with `[storage]` section (see [docs/config.md](docs/config.md))
  - Implicit paths: `${XDG_CONFIG_HOME}/opentask/projects/` (typically `~/.config/opentask/projects/`), `.tasks/` directories
  - Sensible defaults: Creates projects automatically or errors if `config.strict` is set

**Configuration**
- Optional TOML-based config for project settings
- Per-project workflow customization
- Global defaults in `${XDG_CONFIG_HOME}/opentask/config.toml` (typically `~/.config/opentask/config.toml`)

opentask uses two levels of configuration:
- **Project-level** (`.opentask.toml` in each project directory): Workflow, templates, storage location
- **Global** (`${XDG_CONFIG_HOME}/opentask/config.toml`): Manage multiple projects, set active project, share defaults across projects

See [docs/config.md](docs/config.md) for project-level configuration and [docs/global-config.md](docs/global-config.md) for multi-project setup with directory contexts.

### ⏳In Development (Phase 2)

- Sync storage backend to normalised SQL memory db
- SQL query search system
- Different pluggable display modes: List, Table, DAG, Obsidian Canvas

### 🔮Planned (Phase 3)

- Task templating system with built-in templates
- Enhanced CLI UX for epic/story workflows
- Improved error messages and guidance
- Status transition validation


## Example Task File

Tasks are simple markdown files with YAML frontmatter:

```markdown
---
id: s-1234
title: Implement task linking
type: story
status: in-progress
tags: [feature, core]
parent: e-1
links:
  - relates-to: s-5678
  - blocks: s-9101
---

# Implement task linking

Add support for linking tasks together (blocks, relates-to, parent relationships).

This enables building task dependency graphs and hierarchical task organization.

## Acceptance Criteria
- [ ] Tasks can reference other tasks by ID
- [ ] CLI shows task relationships
- [ ] Parent-child relationships organize epics
```

Edit directly in your editor, or use the CLI:
```bash
./opentask --path my_project task update s-1234 --status done
```

## Use Cases

### Solo Developers
Keep tasks in version control, alongside your code. No separate tool to update.

```bash
# In your project repo
opentask --path .tasks task new "Add dark mode" --type story
```

### Teams Coordinating in Code
Track work that's tied to specific codebases without polluting Git with task files.

```bash
# Separate task repo
export opentask_PROJECT_PATH=/mnt/team-tasks
opentask task list "status:in-progress"
```

## Project Structure

Tasks organize hierarchically. By default, they live in `.tasks/`:

```
.tasks/
├── .opentask.toml              # Project configuration
└── 1-my-epic/
    ├── e-1-my-epic.epic.md     # Epic task
    ├── p-2-planning.plan.md    # Plans tied to epic
    ├── r-3-research.research.md # Research
    ├── s-4-feature.story.md    # Implementation stories
    ├── s-5-another.story.md
    └── t-6-specific.task.md    # Specific tasks
```

*File naming convention: `{type_prefix}-{id}-{slug}.{type}.md`. See [docs/design-summary.md](docs/design-summary.md) for full naming details and configuration.*

Or store globally, organized by repo:

```
${XDG_DATA_HOME}/opentask/
├── config.toml                 # Global opentask config
├── templates/                  # Task templates
└── projects/
    ├── github.com-org-project/
    │   ├── 1-my-epic/
    │   │   └── ...
    └── github.com-another-org-project/
        └── ...
```

Task file pattern is configurable, so you can have a johnny decimal system or a hierarchically organised system:

```tmpl
{{ epic.id }}-{{ epic.slug }}/{{ task.id }}-{{ task.slug }}.md
```

or

```tmpl
{{ task.type_id }}-{{ task.type }}/{{ task.type_id }}_{{ task.id }}-{{ task.slug }}.md
```

## AI-Ready Workflow (Coming Soon)

opentask is designed with AI agents in mind. The core task management is **ready today**, but the full AI experience (MCP support, autonomous task management) ships in **Phase 3**.

**Today**: Humans organize tasks that agents will eventually manage
**Phase 3**: Agents can discover projects via MCP and autonomously coordinate work

See [AGENTS.md](AGENTS.md) for current collaboration patterns and the roadmap.


## Build & Run

- have mise installed

```bash
# Build
mise build

# Run tests
mise test

# Use it
./bin/opentask --help
```

