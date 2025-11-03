# OpenTask Documentation

Welcome to the OpenTask documentation. This folder contains comprehensive guides for using and understanding the OpenTask task management system.

## Getting Started

**New to OpenTask?** Start here:
- **[QUICKSTART.md](QUICKSTART.md)** - Learn the basics in 5 minutes

## Core Concepts

Understand how OpenTask works:
- **[DESIGN_SUMMARY.md](DESIGN_SUMMARY.md)** - Architecture and design philosophy
- **[Config.md](Config.md)** - Configuration files and options
- **[GlobalConfig.md](GlobalConfig.md)** - Global configuration reference

## Features

### Task Management
- [QUICKSTART.md](QUICKSTART.md) - Creating and managing tasks

### Project Contexts (New!)
Multiple projects? Use project contexts to avoid `--path` on every command:
- **[ProjectContexts.md](ProjectContexts.md)** - Complete guide to project contexts
  - Multi-project setup
  - Git worktree support
  - Usage examples
  - See also: Multi-Project Setup section in [QUICKSTART.md](QUICKSTART.md)

### Recent Fixes & Features
- **[TaskConfigFixes.md](TaskConfigFixes.md)** - Task 21 bug fixes
  - Task creation location fixes
  - Config discovery improvements
  - Verification procedures

## Testing & Verification

Verify OpenTask is working correctly:
- **[VerificationGuide.md](VerificationGuide.md)** - Quick verification steps
- **[TESTING.md](TESTING.md)** - Comprehensive testing guide

## Implementation & Development

Understanding the codebase:
- **[IMPLEMENTATION_SUMMARY.md](IMPLEMENTATION_SUMMARY.md)** - Implementation details
- **[PHASE4_ROADMAP.md](PHASE4_ROADMAP.md)** - Future roadmap

## Setup & Installation

Environment setup:
- **[MISE.md](MISE.md)** - Using Mise for environment management

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

**Verify everything is working:**
```bash
bash scripts/verify_task21.sh
```

## Document Guide

| Document | Purpose | Audience |
|----------|---------|----------|
| QUICKSTART.md | Get started quickly | Everyone |
| ProjectContexts.md | Multi-project setup | Multi-project teams |
| TaskConfigFixes.md | Understand recent fixes | Developers |
| VerificationGuide.md | Verify installation | Users & QA |
| Config.md | Configuration details | Advanced users |
| GlobalConfig.md | Global config reference | Advanced users |
| DESIGN_SUMMARY.md | Architecture overview | Developers |
| IMPLEMENTATION_SUMMARY.md | Code details | Contributors |
| TESTING.md | Running tests | Developers |
| PHASE4_ROADMAP.md | Future plans | Everyone |
| MISE.md | Development setup | Contributors |

## Key Features

✅ **Hierarchical task organization** - Epics, plans, stories, and tasks
✅ **Flexible workflows** - Customize task statuses and transitions  
✅ **Git-friendly** - Tasks stored as markdown files
✅ **Multi-project support** - Manage multiple projects seamlessly
✅ **Config flexibility** - Project-local or global configuration
✅ **Git worktree compatible** - Works with git worktrees and branches

## Recent Improvements (Session Latest)

### Task 21: Task Creation Location Fixes
- ✅ Fixed config init to generate correct TOML format
- ✅ Fixed config discovery to respect git boundaries
- ✅ 53 tests passing (up from 47)

### Project Context Feature
- ✅ Foundation implemented and tested
- ✅ Context matching algorithm ready
- ✅ CLI commands designed and documented
- ⏳ CLI implementation coming next

## Troubleshooting

**Tasks not showing up?**
- Check config: `opentask config view --path`
- Verify storage location exists
- See [VerificationGuide.md](VerificationGuide.md)

**Can't find a task?**
- List all: `opentask task list`
- Search by type: `opentask task list --type story`
- See [QUICKSTART.md](QUICKSTART.md) for filtering options

**Config issues?**
- View resolved config: `opentask config view`
- Check [Config.md](Config.md) for syntax
- See [TaskConfigFixes.md](TaskConfigFixes.md) for common issues

## Getting Help

1. Check the relevant documentation above
2. Run verification: `bash scripts/verify_task21.sh`
3. Review [VerificationGuide.md](VerificationGuide.md)
4. Check test examples in [TESTING.md](TESTING.md)

## Contributing

Interested in contributing? 
- See [DESIGN_SUMMARY.md](DESIGN_SUMMARY.md) for architecture
- See [TESTING.md](TESTING.md) for test patterns
- See [PHASE4_ROADMAP.md](PHASE4_ROADMAP.md) for planned features

## Version Information

- **OpenTask Version:** Phase 4 (Development)
- **Last Updated:** November 3, 2025
- **Tests Passing:** 53/53 ✅

---

**Happy task tracking!** 🚀
