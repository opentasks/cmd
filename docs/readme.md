# OpenTask Documentation

Welcome to the OpenTask documentation. This folder contains comprehensive guides for using and understanding the OpenTask task management system.

## Getting Started

**New to OpenTask?** Start here:
- **[quickstart.md](quickstart.md)** - Learn the basics in 5 minutes

## Core Concepts

Understand how OpenTask works:
- **[design-summary.md](design-summary.md)** - Architecture and design philosophy
- **[config.md](config.md)** - Configuration files and options
- **[global-config.md](global-config.md)** - Global configuration reference

## Features

### Task Management
- [quickstart.md](quickstart.md) - Creating and managing tasks

### Project Contexts
Multiple projects? Use project contexts to avoid `--path` on every command:
- **[project-contexts.md](project-contexts.md)** - Complete guide to project contexts
  - Multi-project setup
  - Git worktree support
  - Usage examples
  - See also: Multi-Project Setup section in [quickstart.md](quickstart.md)

## Testing & Verification

Verify OpenTask is working correctly:
- **[testing.md](testing.md)** - Comprehensive testing guide

## Implementation & Development

Understanding the codebase:
- **[implementation-summary.md](implementation-summary.md)** - Implementation details
- **[phase4-roadmap.md](phase4-roadmap.md)** - Future roadmap

## Setup & Installation

Environment setup:
- **[mise.md](mise.md)** - Using Mise for environment management

## Quick Reference

### Common Tasks

**Create a project:**
```bash
mkdir my_project && cd my_project
opentask config init --name "My Project" --storage "./.tasks"
opentask task new "First task" --type story
```

**Use with multiple projects:**
```bash
# One-time setup
opentask project attach --project my-project

# Daily usage
opentask task list      # Auto-finds your project
opentask task new "Bug" # Creates in correct location
```

## Document Guide

| Document | Purpose | Audience |
|----------|---------|----------|
| quickstart.md | Get started quickly | Everyone |
| project-contexts.md | Multi-project setup | Multi-project teams |
| config.md | Configuration details | Advanced users |
| global-config.md | Global config reference | Advanced users |
| design-summary.md | Architecture overview | Developers |
| implementation-summary.md | Code details | Contributors |
| testing.md | Running tests | Developers |
| phase4-roadmap.md | Future plans | Everyone |
| mise.md | Development setup | Contributors |

## Key Features

✅ **Hierarchical task organization** - Epics, plans, stories, and tasks
✅ **Flexible workflows** - Customize task statuses and transitions  
✅ **Git-friendly** - Tasks stored as markdown files
✅ **Multi-project support** - Manage multiple projects seamlessly
✅ **Config flexibility** - Project-local or global configuration
✅ **Git worktree compatible** - Works with git worktrees and branches

## Troubleshooting

**Tasks not showing up?**
- Check config: `opentask config view --path`
- Verify storage location exists
- See [testing.md](testing.md)

**Can't find a task?**
- List all: `opentask task list`
- Search by type: `opentask task list --type story`
- See [quickstart.md](quickstart.md) for filtering options

**Config issues?**
- View resolved config: `opentask config view`
- Check [config.md](config.md) for syntax

## Getting Help

1. Check the relevant documentation above
2. Check test examples in [testing.md](testing.md)

## Contributing

Interested in contributing?
- See [design-summary.md](design-summary.md) for architecture
- See [testing.md](testing.md) for test patterns
- See [phase4-roadmap.md](phase4-roadmap.md) for planned features

---

**Happy task tracking!** 🚀
