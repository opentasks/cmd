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

## Config Resolution Strategy

The system implements hierarchical config discovery with merging:

1. **Config File Discovery:** Search for `opentask.toml` starting from current directory and walking up parent directories
2. **Stop Conditions:** Stop searching when reaching either:
   - Filesystem root (`/`)
   - Git repository root (`.git` directory found)
3. **Merging Strategy:** Later configs (closer to current directory) are merged on top of earlier configs (parent directories), so:
   - Closest `opentask.toml` has highest priority
   - Parent directory `opentask.toml` values fill in gaps
   - Further parent configs provide additional values
   - System defaults are lowest priority

**Example:** For `/path/to/project/subfolder/task`:
```
Start: /path/to/project/subfolder/
Found: opentask.toml ← highest priority (merge this first)

Check parent: /path/to/project/
Found: opentask.toml ← merge on top

Check parent: /path/to/
Not found

Stop: /path/to (or if .git found here)
```

## Root Cause Analysis
Suspected insufficient validation in config loading:
1. Config file discovery may not be stopping correctly at git repo root or filesystem root
2. Multiple config files may not be merged in correct order (closest to current dir on top)
3. No visibility into which config files were found or how they were merged
4. No way to verify resolved storage path before creating tasks

## Solution: Two New Commands Needed

### 1. `config view` (or `config show`)
Display the resolved configuration after all walking and merging.

**Purpose:** Help users understand which config files were found, how they were merged, and what the final settings are.

**Output:**
```
Config Resolution Search (starting from current directory, walking up):

Found config files:
  1. /path/to/project/subfolder/opentask.toml (HIGHEST PRIORITY - merged last)
  2. /path/to/project/opentask.toml
  3. /path/to/opentask.toml (stopped here - git repo root found)

Merging order (lowest → highest priority):
  /path/to/opentask.toml
  + /path/to/project/opentask.toml
  + /path/to/project/subfolder/opentask.toml
  = Resolved configuration

Resolved Configuration:
[project]
name = "My Project"

[workflow]
statuses = ["todo", "in-progress", "done"]

[storage]
type = "markdown"
path = "/path/to/project/subfolder/.tasks/"  # RESOLVED ABSOLUTE PATH

Search stopped at: /path/to/ (git repository root)
```

**Flags:**
- `--json` - Output resolved config as JSON
- `--path` - Show only the resolved storage path  
- `--verbose` - Show each config file contents during merging

### 2. `config init` (or `opentask init`)
Initialize a new opentask project in the current directory.

**Purpose:** Create local config file with sensible defaults, making project-specific configuration explicit.

**Behavior:**
- Create `./opentask.toml` in current directory with:
  - Project name (derived from directory name or prompted)
  - Sensible defaults for workflow statuses
  - Local storage path (default: `./.tasks/`)
  - Optional: custom templates path

**Output:**
```
Initialized opentask project in /path/to/current/dir
Created: ./opentask.toml
Storage: ./.tasks/ (local directory)

Next steps:
  1. Create a task: opentask task new "Title"
  2. View tasks: opentask task list
  3. Edit config: ./opentask.toml
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
1. **Config Discovery and Merging:**
   - Config in current dir is found
   - Configs in multiple parent dirs are found and walked
   - Search stops at filesystem root
   - Search stops at git repository root (.git directory)
   - Configs are merged in correct order (furthest parent first, current dir last)
   - Closest config values override parent values
   - Merging respects TOML structure (sections merge correctly)

2. **Config View Output:**
   - Lists all found config files in discovery order
   - Shows merging order and precedence
   - Displays final resolved storage path (absolute)
   - `--path` flag shows only storage path
   - `--json` outputs valid JSON with all resolution details

3. **Task Creation After Init:**
   - `config init` creates valid config
   - `task new` creates task in location specified by resolved config
   - Tasks appear in correct directory when listed
   - Storage path from merged config is respected
   - Multiple config levels work together correctly

## Acceptance Criteria
- [ ] Config discovery walks up from current dir and stops at filesystem root or git repo root
- [ ] Multiple config files found during walk are all listed by `config view`
- [ ] Configs are merged in correct order (furthest parent first, current dir last/highest priority)
- [ ] `config view` displays all found config files and merging order
- [ ] `config view --path` shows only the resolved storage path (absolute)
- [ ] `config view --json` outputs valid JSON with resolution details
- [ ] `config init` creates ./opentask.toml in current dir
- [ ] Config is properly validated and merged after init
- [ ] Tasks created after init appear in correct location from resolved config
- [ ] Helpful error messages if config files are invalid or unreadable
- [ ] `--help` text provided for both commands
- [ ] Integration tests verify config merging with multiple levels
- [ ] Integration tests verify task creation uses merged config path
