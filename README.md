# opentask

A markdown-based task management system written in Go that tracks tasks with metadata, relationships, and flexible status workflows.

## Core Vision

opentask lets you manage tasks as markdown files with YAML frontmatter. Tasks live anywhere—in your project directories, XDG data directories, via environment variables, or command-line arguments. This flexibility makes it ideal for integrating task management into any workflow, whether solo developers or AI agents collaborating with humans.

## Key Features

### Worktree Friendly

- option to not store task data in the same git repository where your agents or users work, avoiding pollution of their project repos.

### Task Storage

Task storage is pluggable, allowing different backends to be used. The default is a simple file-based storage using markdown files.

- **Markdown-based**: Tasks stored as `.md` files with YAML frontmatter metadata
- **Flexible location discovery**: Load projects from (in order of precedence):
  - Explicit paths via:
    - `opentask_PROJECT_PATH` environment variable
    - `--path` CLI argument (cli, mcp, etc)
    - `#path` field in a `.opentask.toml` config file
    - Defaults to creating defined project if not found.
    - if `config.strict`, Errors if not found.

  - implicit paths:
    - `${XDG_DATA_HOME}/opentask/projects/<derived_git_repo_url_project_id>/` 
    - Local directories (`.tasks` convention)
    - if `config.strict`, Errors if not found, otherwise creates at highest precedence location.

### Task Metadata
- **Status**: Customizable per project (e.g., todo → in-progress → done → archived)
- **Type**: Predefined categories - research, spec, plan, story, epic, decision
- **Tags**: Flexible labels for organization and filtering
- **Relationships**: Link tasks with other tasks
  - `blocks` - this task blocks others
  - `relates-to` - related but independent
  - `parent` - hierarchical relationships
  - all tasks have at least one link: to their initiative/epic

### Project Organization
- **Kanban-style workflow**: Status determines column position
- **Type-driven filtering**: Filter by research, specs, plans, etc.
- **simple project rules**: "When creating tasks, read the plans [linked in the epic](opentask://the-epic-task-id) to guide task creation."
- **Templates**: Predefined task templates for common types

## Example Task File

```markdown
---
id: s-1234
title: Implement task linking
type: story
status: in-progress
tags: [feature, core]
parent: setup-data-model
links:
  - relates-to: s-5678
  - blocks: s-91011
---

# Implement task linking

Add support for linking tasks together (blocks, relates-to, parent relationships).

This enables building task dependency graphs and hierarchical task organization.
```

## Project Structure

```sh
project_id/
  config.toml                # project config (undecided what this contains yet)
  1234-some-name/
    1234.epic.md             # epic task
    1234.1.plan.md           # plan
    1234.2.research.md       # research
    1234.3.story.md          # story task
    1234.4.story.md          # story task
    1234.5.decision.md       # decision
    1234.6.story.md          # story task
    1234.7.task.md           # task
    1234.8.task.md           # task
    1234.9.task.md           # task
```

Depending on where the project is resolved this might be in: 

```sh
${XDG_DATA_HOME}/opentask/
  config.toml                 # global user opentask config 
  templates/
    epic.md
    plan.md
    research.md
    story.md
    decision.md
    task.md

  projects/
    <derived_git_repo_url_project_id>/
      config.toml              # project config
      1234-some-name/
        1234.epic.md
        ...
```

or at `opentask_PROJECT_PATH` or `--path` specified location.

```sh
/my/custom/path/
  config.toml                   # project config
  1234-some-name/
    1234.epic.md
    ...
```

## Planned Components

- **Core**: Task model, parsing, relationship resolution
- **Storage**: Load/save from multiple locations
  - BaseStorage interface
  - MarkdownFileStorage implementation
- **Query**: Filter by status, type, tags, relationships
- **CLI**: Interactive and batch operations
  - Use viper/cobra for config and CLI
- **MCP STDIO**: Multi-Client Proxy for AI agent integration

## Near Term Roadmap

- Task delegation: agents can assign tasks to themselves
- Status workflows: customizable per project
- Task templates: predefined structures for common task types


## Non-Goals (For Now)

- Cloud sync or authentication
- Web UI (CLI-first)
- Real-time collaboration (file-based for now)
- Integration with other tools (extensible later)
- **Export**: Various formats (JSON, YAML, etc.) (this is up to the storage backend)

---

**Status**: Planning phase
