# Project Selection Mechanism Analysis

## Current State
- Must use `--path` flag explicitly: `opentask --path ~/.local/share/tasks task list`
- Requires knowing exact path to project's `.tasks` directory
- Not ergonomic for frequent use

## Option 1: Global Config Project Matching
**Concept:** User defines projects in `~/.config/opentask/config.toml`, system matches current directory to registered project

**Example Global Config:**
```toml
[[projects]]
id = "work"
name = "Work Tasks"
[projects.storage]
path = "~/projects/work/.tasks"

[[projects]]
id = "personal"
name = "Personal Tasks"
[projects.storage]
path = "~/tasks/.tasks"
```

**Usage:**
```bash
cd ~/projects/work
opentask task list  # Automatically finds "work" project
cd ~/tasks
opentask task list  # Automatically finds "personal" project
```

**Pros:**
- Centralized project definition
- Works across multiple directories
- Feels like workspace management
- Can set `active_project` default

**Cons:**
- Requires global config setup
- Path matching logic needed (exact match? prefix match?)
- What if cwd matches multiple projects?
- Only works if projects are at known paths

---

## Option 2: Environment Variable
**Concept:** User sets `OPENTASK_PATH` or `OPENTASK_PROJECT` env var

**Usage:**
```bash
export OPENTASK_PATH=~/.local/share/tasks
# or
export OPENTASK_PROJECT=work  # looks up in global config

opentask task list
opentask task new "Task"
```

**Pros:**
- Simple to implement
- Works with any shell/IDE integration
- Easy to script
- Can be set per-terminal session

**Cons:**
- Not discoverable
- Users have to remember to set it
- No persistence between sessions
- Could conflict with local `.opentask.toml` in cwd
- Becomes clutter in environment

---

## Option 3: Project Context / Active Project (RECOMMENDED)
**Concept:** Combination approach using global config + active project selection

**Global Config with Active Project:**
```toml
active_project = "personal"

[[projects]]
id = "work"
name = "Work Tasks"
[projects.storage]
path = "~/projects/work/.tasks"

[[projects]]
id = "personal"
name = "Personal Tasks"
[projects.storage]
path = "~/tasks/.tasks"
```

**Commands:**
```bash
# List projects
opentask config projects
# Output:
#   work (Work Tasks)
#   personal (Personal Tasks)
# * personal (active)

# Switch active project
opentask config projects --set-active work

# Use explicit project
opentask --project work task list
opentask -p work task new "Task"

# Use active project (default)
opentask task list
opentask task new "Task"
```

**Pros:**
- Centralized configuration
- Clear, explicit project selection
- Can have sensible default (active_project)
- Works well with shell aliases/functions
- Discoverable via `opentask config projects`
- Can override with `--project` flag per command

**Cons:**
- Requires more config setup initially
- Need to implement `--project` flag on all commands
- Ambiguous if multiple definitions (global + local)

---

## Option 4: Workspace / Directory Scoping
**Concept:** Use `.opentask-workspace` marker file at root of project workspace

**Structure:**
```
~/projects/my-workspace/
  .opentask-workspace   # Marker file defining workspace
  .opentask.toml        # Project config
  .tasks/
  docs/
  src/
```

**Content of .opentask-workspace:**
```
name: my-workspace
path: .tasks
```

**Usage:**
```bash
cd ~/projects/my-workspace/docs
opentask task list  # Walks up, finds .opentask-workspace, uses it
```

**Pros:**
- Discoverable pattern (like .git, package.json)
- Works from any subdirectory
- Can define workspace separately from project config
- Clear intent

**Cons:**
- New file format/concept to explain
- Still requires per-workspace setup

---

## Option 5: Symlink / Alias Pattern
**Concept:** Use well-known symlink location

```bash
# One-time setup:
ln -s /path/to/real/tasks ~/.tasks
# or create standardized location like XDG spec

opentask --path ~/.tasks task list
```

**Pros:**
- Minimal code change
- Works with existing system tools

**Cons:**
- Still requires --path flag
- Doesn't solve the core problem
- More of a workaround

---

## Recommended Approach

**Use Option 3 (Active Project) with smart defaults:**

1. **Phase 1: Active Project Selection**
   - Add `--project` / `-p` flag to all commands
   - Define projects in global config with storage paths
   - Implement `opentask config projects --set-active <id>`
   - If no --project flag: use active_project from global config

2. **Phase 2: Smart Path Matching**
   - When active_project is not set, try to auto-detect from cwd
   - Match current directory to project paths in global config
   - If no match, fall back to discovery from current directory

3. **Phase 3: Enhanced UX**
   - Shell completion for project names
   - Quick-switch aliases: `alias work='opentask --project work'`
   - Show which project is active in CLI output

---

## Decision Matrix

| Feature | Option 1 | Option 2 | Option 3 | Option 4 | Option 5 |
|---------|----------|----------|----------|----------|----------|
| Simplicity | ⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐ | ⭐⭐⭐⭐⭐ |
| Discovery | ⭐⭐ | ⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐ |
| Ergonomics | ⭐⭐⭐ | ⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐ |
| Flexibility | ⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐ | ⭐⭐ |
| Setup Burden | ⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ |

**Recommendation: Option 3 (Active Project) - Best balance of UX and implementation**
