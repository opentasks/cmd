# Project Context Design

## Concept

Projects without their own `.opentask.toml` are matched to task storage via explicit **context paths** in the global config.

A "context" is simply a working directory that should use a specific project's tasks.

## Global Config Structure

```toml
# ~/.config/opentask/config.toml
active_project = "personal-notes"

[[projects]]
id = "personal-notes"
name = "Personal Notes"

[projects.storage]
backend = "markdown-fs"
path = "/home/zenobius/Notes/.tasks"

# This project has no contexts - uses only with --project flag
# or when it's the active_project


[[projects]]
id = "opentask"
name = "OpenTask Project"

[projects.storage]
backend = "markdown-fs"
path = "/home/zenobius/Notes/Projects/OpenTask/.tasks"

# Context paths: where this project's tasks apply
[[projects.context]]
path = "/mnt/Store/Projects/Mine/Github/opentasks"

[[projects.context]]
path = "/mnt/Store/Projects/Mine/Github/opentasks.worktrees/feature-x"

[[projects.context]]
path = "/mnt/Store/Projects/Mine/Github/opentasks.worktrees/bugfix-y"


[[projects]]
id = "work-tasks"
name = "Work Project"

[projects.storage]
backend = "markdown-fs"
path = "/home/zenobius/Notes/Work/.tasks"

[[projects.context]]
path = "/mnt/projects/work-repo"

[[projects.context]]
path = "/mnt/projects/work-repo.worktrees/main"

[[projects.context]]
path = "/mnt/projects/work-repo.worktrees/develop"
```

## Matching Algorithm

When user runs `opentask task list`:

1. **Check if current dir has .opentask.toml**
   - If yes, use that project's storage
   
2. **Check context paths in active_project**
   - If cwd matches any context path, use active_project
   
3. **Check context paths in all projects**
   - Find longest matching context path (most specific wins)
   - Use that project's storage
   
4. **Fallback**
   - Use active_project
   - If no active_project, error: "No project found for this context"

## Usage Examples

### Setup (One-time)

**1. Create global config with projects:**
```bash
mkdir -p ~/.config/opentask
cat > ~/.config/opentask/config.toml << 'TOML'
active_project = "personal-notes"

[[projects]]
id = "personal-notes"
name = "Personal Notes"
[projects.storage]
backend = "markdown-fs"
path = "/home/zenobius/Notes/.tasks"

[[projects]]
id = "opentask"
name = "OpenTask Project"
[projects.storage]
backend = "markdown-fs"
path = "/home/zenobius/Notes/Projects/OpenTask/.tasks"

[[projects.context]]
path = "/mnt/Store/Projects/Mine/Github/opentasks"
TOML
```

**2. Attach current directory to a project:**
```bash
cd /mnt/Store/Projects/Mine/Github/opentasks
opentask project attach
# Adds current directory to opentask project's contexts
# (must determine which project - maybe require --project flag)

opentask project attach --project opentask
# Explicitly specify project
```

**3. Attach a different directory:**
```bash
cd /mnt/Store/Projects/Mine/Github/opentasks/src
opentask project attach /mnt/Store/Projects/Mine/Github/opentasks.worktrees/feature-x --project opentask
# Adds feature-x worktree to opentask project's contexts
```

### Daily Usage

**In opentask repo (context matches "opentask" project):**
```bash
cd /mnt/Store/Projects/Mine/Github/opentasks/src
opentask task list       # Uses opentask project
opentask task new "Bug"  # Creates in /home/zenobius/Notes/Projects/OpenTask/.tasks
```

**In personal notes (uses active_project):**
```bash
cd /home/zenobius/Notes
opentask task list       # Uses personal-notes (active)
opentask task new "Note" # Creates in /home/zenobius/Notes/.tasks
```

**Override with --project flag:**
```bash
cd /mnt/Store/Projects/Mine/Github/opentasks
opentask --project personal-notes task list
# Temporarily uses personal-notes project instead
```

**Switch active project:**
```bash
opentask config projects --set-active work-tasks
opentask task list  # Now uses work-tasks by default
```

## Commands

### `opentask project attach [PATH]`

Attach a working directory to a project.

**Syntax:**
```bash
opentask project attach [PATH] [--project PROJECT_ID]
```

**Arguments:**
- `PATH` (optional): Directory to attach. Defaults to current directory.
  - Can be relative: `../`, `../../other-project`
  - Can be absolute: `/mnt/projects/repo`

**Options:**
- `--project PROJECT_ID` (required if ambiguous): Which project to attach to

**Examples:**
```bash
# Attach current directory to opentask project
cd /mnt/Store/Projects/Mine/Github/opentasks
opentask project attach --project opentask

# Attach specific directory
opentask project attach /mnt/Store/Projects/Mine/Github/opentasks.worktrees/feature-x --project opentask

# Attach relative path
cd /mnt/Store/Projects/Mine/Github/opentasks/src
opentask project attach ../../ --project opentask

# If only one project, --project can be omitted
opentask project attach  # Error: multiple projects, specify --project
```

### `opentask project detach [PATH]`

Remove a working directory from a project's contexts.

**Syntax:**
```bash
opentask project detach [PATH] [--project PROJECT_ID]
```

**Examples:**
```bash
# Remove current directory
opentask project detach --project opentask

# Remove specific directory
opentask project detach /mnt/Store/Projects/Mine/Github/opentasks.worktrees/old-branch --project opentask
```

### `opentask project list`

List all projects and their contexts.

**Output:**
```
Projects:

work-tasks (Work Project)
  Storage: /home/zenobius/Notes/Work/.tasks
  Contexts:
    - /mnt/projects/work-repo
    - /mnt/projects/work-repo.worktrees/main

opentask (OpenTask Project)
  Storage: /home/zenobius/Notes/Projects/OpenTask/.tasks
  Contexts:
    - /mnt/Store/Projects/Mine/Github/opentasks
    - /mnt/Store/Projects/Mine/Github/opentasks.worktrees/feature-x

personal-notes (Personal Notes) *
  Storage: /home/zenobius/Notes/.tasks
  Contexts: (none)

* = active_project
```

### `opentask config projects --set-active PROJECT_ID`

Set the default/active project.

```bash
opentask config projects --set-active work-tasks
opentask task list  # Now uses work-tasks
```

## Advantages

✅ **Explicit mapping** - No magic, users control what goes where
✅ **Multiple contexts per project** - Handle git worktrees, branches, etc.
✅ **Easy to attach** - One command: `opentask project attach`
✅ **Longest match wins** - More specific paths take precedence
✅ **Works from anywhere** - Any directory under a context path works
✅ **Simple to understand** - Just a list of paths in config
✅ **No auto-detection** - Predictable behavior

## Implementation Details

### Config Validation
- Validate that all context paths are absolute and resolvable
- Warn if path doesn't exist (but allow - might be created later)
- Prevent duplicate context paths

### Path Matching Algorithm
```
Given: current working directory (cwd)

1. If cwd or parent has .opentask.toml:
   → Use that project
   
2. Load global config projects
   
3. For each project:
   - Check if cwd matches any context path
   - Record longest matching path and project
   
4. If match found:
   → Use that project
   
5. Otherwise:
   → Use active_project
   → If no active_project: error
```

### `opentask project attach` Logic
```bash
opentask project attach [PATH] --project PROJECT_ID

1. Resolve PATH to absolute path
2. Validate path exists
3. Load global config
4. Find project by PROJECT_ID
5. Check if path already in project.context
6. Add path to project.context
7. Save global config
8. Show confirmation
```

## Migration Path

**Current workflow:**
```bash
opentask --path /home/zenobius/Notes/Projects/OpenTask/.tasks task list
```

**After implementation:**
```bash
# One-time setup:
opentask project attach --project opentask

# Then:
opentask task list
```

## Example Workflow: Multi-Worktree Project

**Initial setup:**
```bash
# Create project in global config
cat > ~/.config/opentask/config.toml << 'TOML'
[[projects]]
id = "myproject"
name = "My Project"
[projects.storage]
path = "/home/user/Notes/Projects/MyProject/.tasks"
TOML

# Attach main worktree
cd /mnt/repos/myproject
opentask project attach --project myproject

# Attach feature branch worktree
cd /mnt/repos/myproject.worktrees/feature-x
opentask project attach --project myproject

# Attach another worktree
cd /mnt/repos/myproject.worktrees/bugfix-y
opentask project attach --project myproject
```

**Daily usage:**
```bash
# In any worktree
cd /mnt/repos/myproject.worktrees/feature-x/src
opentask task list          # ✓ Uses myproject
opentask task new "Feature" # ✓ Creates in shared storage
opentask task done 42       # ✓ Marks task as done

# All worktrees share the same tasks
cd /mnt/repos/myproject
opentask task list  # ✓ Shows same tasks
# Task 42 is already marked done
```

## FAQ

**Q: Why not just have --path be persistent (with --set-path)?**
A: Context list is more flexible. Allows multiple directories for one project (worktrees, branches).

**Q: What if I move a project directory?**
A: Context paths become stale. Need to run `opentask project detach` + `opentask project attach` at new location.

**Q: Can I edit contexts manually in config?**
A: Yes, but it's safer to use `opentask project attach/detach`.

**Q: How are context paths matched - exact or prefix?**
A: Prefix match, with longest path winning (most specific wins).

**Q: What about symlinks?**
A: Paths are resolved to absolute canonical paths.

**Q: Can I have overlapping contexts?**
A: Yes, but longest match wins. Useful for subdirectories.

**Q: What if project has no contexts and no active_project?**
A: Error: "No project found for this context. Use --project flag or set active_project."

