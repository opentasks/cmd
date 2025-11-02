---
id: 10
title: CLI Architecture with Viper/Cobra
type: story
status: todo
tags: [design, cli]
relationships:
  - type: parent
    taskID: 1

createdAt: 2025-11-02T08:00:00Z
updatedAt: 2025-11-02T08:00:00Z
---

# CLI Architecture with Viper/Cobra

## Design Decision

CLI uses Cobra for command structure and Viper for configuration management. This provides a clean, extensible interface for task operations.

## Command Structure

### Current Implementation
```
opentasks
├── task
│   ├── new
│   ├── list
│   ├── show
│   ├── update
│   ├── delete
│   └── link (pending: s-15)
├── project
│   ├── new
│   └── list
└── config
    ├── show
    ├── set
    └── get
```

### Future Enhancements
- task link command (s-15: Add task linking command)
- config view (b-21: Show resolved config path)
- config init (b-21: Initialize local config)
- status transitions (currently not designed in detail)
- project config/path subcommands (currently not designed in detail)

### Design Notes on Aliases
- Design originally specified aliases (new|create, list|ls, show|view, delete|rm)
- Current implementation uses primary commands only
- Can add aliases later if needed for convenience

## Viper Integration

```go
// Global config binding
viper.SetConfigName("config")
viper.AddConfigPath("$HOME/.config/opentasks")
viper.AddConfigPath(".")

// Environment variable binding
viper.BindEnv("project.path", "OPENTASKS_PROJECT_PATH")
viper.BindEnv("storage.backend", "OPENTASKS_STORAGE_BACKEND")

// Flag binding (per command)
cmd.Flags().StringVar(&path, "path", "", "Project path")
viper.BindPFlag("project.path", cmd.Flags().Lookup("path"))
```

## Precedence (Viper)

Config loaded in order:
1. Built-in defaults
2. Global config (`~/.config/opentasks/config.toml`)
3. Project configs via hierarchical walk:
   - Walk up from current directory looking for `opentasks.toml`
   - Stop at filesystem root (`/`) or git repository root (`.git`)
   - Merge configs with closest (current dir) having highest priority
4. Environment variables (`OPENTASKS_*`)
5. CLI flags (highest priority)

**Note:** Hierarchical config discovery and merging is implemented but not transparent. See b-21 for `config view` and `config init` commands to help users understand resolved config.

## Common Commands

### Create Task
```bash
opentasks task new "My story title" --type story --parent 5 --tag feature
```

### List Tasks
```bash
opentasks task ls --status in-progress --type story
opentasks task ls --parent 5     # Tasks under epic 5
opentasks task ls --tag feature
```

### Show Task Details
```bash
opentasks task show 42
```

### Link Tasks
```bash
opentasks task link 42 --blocks 10
opentasks task link 42 --relates-to 3
```

## Implementation Status

### ✅ Implemented
- Cobra root command initializes storage in PersistentPreRunE hook
- All subcommands have access to global Engine and Store
- Global flags: `--path`, `--config`, `--verbose`
- Flag-to-Viper binding for configuration
- Basic error handling in CLI commands

### 🔄 In Progress / Future Work
- Output formats (JSON/YAML): Being enhanced with `--format` flag (s-16)
- Task linking: Implemented as separate task (s-15)
- Config transparency: `config view` and `config init` commands (b-21)
- Color-coded terminal output (currently plain text)
- Enhanced help text and command examples

## Design Validation

This CLI architecture successfully achieves:
1. **Clean separation of concerns** - Cobra handles CLI, Viper handles config, storage layer separate
2. **Extensibility** - New commands easily added as subcommands
3. **Configuration flexibility** - Supports flags, env vars, config files with proper precedence
4. **Storage abstraction** - Pluggable storage backends via interface

**Known limitation:** Config discovery and merging is implicit rather than transparent. The b-21 bug task addresses this with `config view` and `config init` commands.
