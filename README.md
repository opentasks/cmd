# opentask

**Git-agnostic task management for developers and AI agents.**

Markdown-based tasks that live anywhere—in your projects, user directories, or CI pipelines. No database, no setup, no pollution of your Git repos. Perfect for solo developers, teams tracking work in code, and organizations planning AI agent integration.

## Why opentask?

- 📁 **Git-clean** – Store tasks outside project repos, keep them in sync anywhere
- 🔗 **Flexible paths** – XDG directories, environment variables, explicit paths—your choice
- 🤖 **AI-ready** ⏳ – Built for agent collaboration via MCP (Phase 3 roadmap)
- 📝 **Markdown native** – Edit in VS Code, Vim, or any editor
- ⚡ **Zero setup** – Works with sensible defaults, optional TOML configuration
- 🔌 **Pluggable storage** ⏳ – Swap backends (file → SQLite → cloud) without changing code

## Quick Start

```bash
# Build the CLI
go build -o opentask ./cmd/opentask

# Create your first task
./opentask --path my_project task new "Build login page" --type story

# List all tasks
./opentask --path my_project task list

# Filter by status or type
./opentask --path my_project task list --status todo --type story

# Show task details
./opentask --path my_project task show 1

# Update task status
./opentask --path my_project task update 1 --status in-progress
```

For detailed usage, see [QUICKSTART.md](docs/QUICKSTART.md).

## Core Features

### ✅ Ready Now (MVP)

**Task Management**
- Full CRUD operations (create, read, update, delete)
- Markdown files with YAML frontmatter
- Customizable status workflows (todo → in-progress → done → archived)
- Type-driven organization (research, spec, plan, story, epic, decision)
- Flexible tagging and filtering

**Relationships**
- Link tasks together: `blocks`, `relates-to`, `parent`
- Hierarchical task organization (epics contain stories and tasks)
- Query and traverse task dependencies

**Flexible Storage**
- Default: File-based markdown storage (no database needed)
- Project discovery from multiple locations (in order of precedence):
  - Explicit paths: `--path` CLI flag or `opentask_PROJECT_PATH` env var
  - Config file: `.opentask.toml` with `#path` field
  - Implicit paths: `${XDG_DATA_HOME}/opentask/projects/`, `.tasks/` directories
  - Sensible defaults: Creates projects automatically or errors if `config.strict` is set

**Configuration**
- Optional TOML-based config for project settings
- Per-project workflow customization
- Global defaults in `${XDG_DATA_HOME}/opentask/config.toml`

### ⏳ In Development (Phase 2)

- Task templating system with built-in templates
- Enhanced CLI UX for epic/story workflows
- Improved error messages and guidance
- End-to-end testing suite
- Interactive task editing within the CLI
- Status transition validation
- Better task list printing and formatting

### 🔮 Planned (Phase 3)

- **MCP (Model Context Protocol)** – AI agents discover and manage tasks
- Advanced query capabilities and saved filters
- Additional storage backends (SQLite, DuckDB, cloud)
- Interactive terminal UI (alongside CLI)
- Task delegation and assignment tracking

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
./opentask --path .tasks task new "Add dark mode" --type story
```

### Teams Coordinating in Code
Track work that's tied to specific codebases without polluting Git with task files.

```bash
# Separate task repo
export opentask_PROJECT_PATH=/mnt/team-tasks
opentask task list --status in-progress
```

### AI Agent Collaboration ⏳ (Phase 3)
Agents will discover and autonomously manage task workflows via MCP.

```python
# Coming in Phase 3
async with MCPClient('opentask') as client:
    tasks = await client.get_tasks(status='todo')
    await client.update_task(tasks[0].id, status='in-progress')
```

## Project Structure

Tasks organize hierarchically. By default, they live in `.tasks/`:

```
.tasks/
├── .opentask.toml              # Project configuration
└── 1-my-epic/
    ├── 1.epic.md               # Epic task
    ├── 1.1.plan.md             # Plans tied to epic
    ├── 1.2.research.md         # Research
    ├── 1.3.story.md            # Implementation stories
    ├── 1.4.story.md
    └── 1.5.task.md             # Specific tasks
```

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

## AI-Ready Workflow (Coming Soon)

opentask is designed with AI agents in mind. The core task management is **ready today**, but the full AI experience (MCP support, autonomous task management) ships in **Phase 3**.

**Today**: Humans organize tasks that agents will eventually manage  
**Phase 3**: Agents can discover projects via MCP and autonomously coordinate work

See [AGENTS.md](AGENTS.md) for current collaboration patterns and the roadmap.

## Roadmap

### ✅ Phase 1: MVP (Complete)
- [x] Task CRUD operations
- [x] Markdown-based storage
- [x] Flexible project discovery
- [x] Task relationships and filtering
- [x] CLI with Cobra/Viper
- [x] Configuration system

### ⏳ Phase 2: Polish & Features (In Progress)
- [ ] Task templating system
- [ ] Enhanced CLI UX for epics/stories
- [ ] Improved error messages
- [ ] End-to-end testing
- [ ] Interactive task editing
- [ ] Status transition validation

### 🔮 Phase 3: AI Integration
- [ ] MCP (Model Context Protocol) server
- [ ] Agent task discovery and delegation
- [ ] Advanced query engine
- [ ] Multiple storage backends (SQLite, DuckDB, etc.)
- [ ] Collaboration features

For detailed progress, check [SESSION_SUMMARY.md](SESSION_SUMMARY.md).

## Documentation

- **[QUICKSTART.md](docs/QUICKSTART.md)** – Get started in 5 minutes
- **[DESIGN_SUMMARY.md](docs/DESIGN_SUMMARY.md)** – Architecture overview
- **[IMPLEMENTATION_SUMMARY.md](docs/IMPLEMENTATION_SUMMARY.md)** – What's built
- **[MISE.md](docs/MISE.md)** – Development task runner guide
- **[AGENTS.md](AGENTS.md)** – AI agent integration (roadmap)

## Build & Run

**Requirements**: Go 1.21+ (or use `mise` for task automation)

```bash
# Build
go build -o opentask ./cmd/opentask

# Run tests
go test ./...

# Use it
./opentask --help
```

For development: see [MISE.md](docs/MISE.md) for task automation.

---

**Status**: MVP Complete (Phase 1 ✅)  
**Next**: Phase 2 (Polish & Features) – Estimated Q4 2025  
**Future**: Phase 3 (AI Integration) – Estimated Q1 2026
