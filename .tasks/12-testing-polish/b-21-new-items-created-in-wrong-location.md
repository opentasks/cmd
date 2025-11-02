---
id: 21
title: New items created in wrong location
type: bug
status: todo
tags:
    - config
    - storage
relationships:
    - type: parent
      taskID: 12
createdAt: "2025-11-02T20:50:00Z"
updatedAt: "2025-11-02T20:50:00Z"
---

## Problem
When creating new tasks via CLI, they are not being saved in the correct project location. This is likely due to insufficient configuration validation and unclear config resolution logic.

The system should support hierarchical config discovery (local → parent dirs → home), but it's unclear if:
- The resolved config path is correct
- Multiple config files are being properly merged
- Task storage location is being determined from the right config source

## Root Cause Analysis
Suspected insufficient validation in config loading:
1. Config file discovery may not be walking up parent directories correctly
2. Multiple config sources (./opentasks.toml, parent dir opentasks.toml, ~/.config/opentasks/config.toml) may not be merged properly
3. No visibility into which config is actually being used
4. No way to verify resolved storage path before creating tasks

## Solution: Two New Commands Needed

### 1. `config view` (or `config show`)
Display the resolved configuration after all merging and defaults.

**Purpose:** Help users understand which config is being used and what the final settings are.

**Output:**
```
Config Resolution Path:
  1. Current directory: ./opentasks.toml (found: yes/no)
  2. Parent directories: ../, ../../ (found: yes/no)
  3. Home directory: ~/.config/opentasks/config.toml (found: yes/no)

Resolved Configuration:
[project]
name = "My Project"
...

[workflow]
statuses = ["todo", "in-progress", "done"]
...

[storage]
type = "markdown"
path = "/Users/me/.local/share/opentasks/projects/abc123/"  # RESOLVED PATH

Active config sources (in order of precedence):
  - ./opentasks.toml (local, highest priority)
  - ~/.config/opentasks/config.toml (user-level)
  - Built-in defaults (lowest priority)
```

**Flags:**
- `--json` - Output resolved config as JSON
- `--path` - Show only the resolved storage path
- `--verbose` - Show config resolution details and merging order

### 2. `config init` (or `opentasks init`)
Initialize a new OpenTasks project in the current directory.

**Purpose:** Create local config file with sensible defaults, making project-specific configuration explicit.

**Behavior:**
- Create `./opentasks.toml` in current directory with:
  - Project name (derived from directory name or prompted)
  - Sensible defaults for workflow statuses
  - Local storage path (default: `./.tasks/`)
  - Optional: custom templates path

**Output:**
```
Initialized OpenTasks project in /path/to/current/dir
Created: ./opentasks.toml
Storage: ./.tasks/ (local directory)

Next steps:
  1. Create a task: opentasks task new "Title"
  2. View tasks: opentasks task list
  3. Edit config: ./opentasks.toml
```

**Flags:**
- `--name` - Project name (skip prompt)
- `--storage` - Storage path (default: ./.tasks/)
- `--force` - Overwrite existing config.toml

## Testing Requirements

### Unit Tests
- Config view correctly merges multiple config sources
- Resolved paths are absolute and correct
- JSON output is valid JSON
- Init creates valid TOML config

### Integration Tests
1. **Config Discovery:**
   - Config in current dir is found
   - Config in parent dir is found
   - Multiple levels of parent dirs are checked
   - Precedence order is correct (local > parent > home)

2. **Task Creation After Init:**
   - `config init` creates valid config
   - `task new` creates task in correct location specified by config
   - Tasks appear when listed
   - Storage path from config is used

## Acceptance Criteria
- [ ] `config view` displays resolved configuration
- [ ] `config view --json` outputs valid JSON
- [ ] `config view --path` shows only storage path
- [ ] `config init` creates ./opentasks.toml in current dir
- [ ] Config is properly validated after init
- [ ] Tasks created after init appear in correct location
- [ ] Multiple config sources are merged in correct order (local → parent → home)
- [ ] Helpful error messages if config is invalid
- [ ] `--help` text provided for both commands
- [ ] Integration tests verify full workflow
