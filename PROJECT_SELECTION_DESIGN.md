# Project Selection - Design Document

## The Core Issue

**Global config storage path** and **current working directory** are usually unrelated:

```
Global config:
  project: my-notes
  path: /home/zenobius/Notes/Something/

Your cwd:
  /mnt/Store/Projects/Mine/Github/SomeReponame/src/somepath/
  
❌ No reliable correlation possible
```

## Recommended Solution: Three-Layer Approach

### Priority Order:
1. **`--project` flag** (explicit per-command override)
2. **`.opentask-project` marker file** (walked up from cwd)
3. **`active_project` from global config** (sensible default)

---

## Layer 1: Explicit Override (--project flag)

```bash
opentask --project my-notes task list
opentask -p my-notes task new "Quick task"
```

**Use when:** You need to work with a specific project regardless of cwd

---

## Layer 2: Project Marker File (.opentask-project)

Place `.opentask-project` in the root of your project:

```
/mnt/Store/Projects/Mine/Github/SomeReponame/
  .opentask-project    ← Create this
  .git/
  src/
  ...
```

**Content of .opentask-project:**
```toml
id = "my-github-project"
name = "GitHub Project Tasks"

# Optional: override storage path from global config
# [storage]
# path = "/home/zenobius/Notes/Projects/SomeProject/.tasks"
```

**Usage:**
```bash
cd /mnt/Store/Projects/Mine/Github/SomeReponame/src/some/deep/path
opentask task list
# ✓ Walks up directory tree
# ✓ Finds .opentask-project
# ✓ Reads id = "my-github-project"
# ✓ Looks up in global config
# ✓ Uses correct task storage
```

**Advantages:**
- Works from any subdirectory of the project
- Explicit and unambiguous
- Same pattern as `.git`, `package.json`, etc.
- Can override storage path if needed

---

## Layer 3: Active Project Default

Set a default project in global config:

**~/.config/opentask/config.toml:**
```toml
active_project = "my-notes"

[[projects]]
id = "my-notes"
name = "My Notes"
[projects.storage]
path = "/home/zenobius/Notes/Something/.tasks"

[[projects]]
id = "my-github-project"
name = "GitHub Project Tasks"
[projects.storage]
path = "/home/zenobius/Notes/Projects/SomeProject/.tasks"

[[projects]]
id = "work-tasks"
name = "Work Tasks"
[projects.storage]
path = "/home/zenobius/Notes/Work/.tasks"
```

**Usage:**
```bash
# Uses active_project (my-notes)
opentask task list

# Switch to different project
opentask config projects --set-active work-tasks
opentask task list  # Now shows work-tasks

# Switch back
opentask config projects --set-active my-notes
```

**View current settings:**
```bash
opentask config projects
# Output:
#   my-notes (My Notes) *
#   my-github-project (GitHub Project Tasks)
#   work-tasks (Work Tasks)
```

---

## Complete Usage Example

### Setup (One-time)

**1. Create global config:**
```bash
mkdir -p ~/.config/opentask
cat > ~/.config/opentask/config.toml << 'TOML'
active_project = "personal"

[[projects]]
id = "personal"
name = "Personal Tasks"
[projects.storage]
path = "/home/zenobius/Notes/.tasks"

[[projects]]
id = "opentask"
name = "OpenTask Project"
[projects.storage]
path = "/home/zenobius/Notes/Projects/OpenTask/.tasks"
TOML
```

**2. In project repos, add .opentask-project:**
```bash
# In /mnt/Store/Projects/Mine/Github/opentasks/
echo 'id = "opentask"' > .opentask-project
```

### Daily Usage

**Working with default project:**
```bash
opentask task list
opentask task new "My task"
```

**Working with different project:**
```bash
opentask --project opentask task list
opentask --project opentask task new "Feature request"
```

**In project repo (automatically finds tasks):**
```bash
cd /mnt/Store/Projects/Mine/Github/opentasks/src
opentask task list  # Automatically uses opentask project
opentask task new "Bug"
```

**Switch default for session:**
```bash
opentask config projects --set-active opentask
opentask task list
opentask task new "Task"
opentask config projects --set-active personal
```

---

## Implementation Roadmap

### Phase 1: Active Project (MVP)
- [x] Projects defined in global config
- [ ] `active_project` field in global config
- [ ] `--project` flag on root command
- [ ] `opentask config projects --set-active <id>`

### Phase 2: Project Markers
- [ ] Search for `.opentask-project` walking up from cwd
- [ ] Parse TOML format
- [ ] Use project ID from marker
- [ ] Optional storage path override in marker

### Phase 3: Enhanced UX
- [ ] Show current project in output
- [ ] Shell completion for project names
- [ ] `opentask config projects list` command
- [ ] Better error messages when project not found

---

## Migration Path

**Current workflow with --path flag:**
```bash
opentask --path /home/zenobius/Notes/.tasks task list
```

**After Phase 1 (Active Project):**
```bash
# Set in global config, then:
opentask task list
# Or override:
opentask --project opentask task list
```

**After Phase 2 (Project Markers):**
```bash
# In project repos with .opentask-project marker:
cd /mnt/Store/Projects/Mine/Github/opentasks
opentask task list  # Auto-detects project
```

---

## FAQ

**Q: Why not auto-detect from cwd?**
A: Storage path and working directory are disconnected. A reliable heuristic is impossible.

**Q: Can I use environment variables?**
A: Not in this design, but could be added. Would you like that?

**Q: What if I'm in a repo without .opentask-project?**
A: Falls back to active_project default, then current directory discovery.

**Q: Do I need .opentask.toml and .opentask-project?**
A: No. They serve different purposes:
- `.opentask.toml`: Local project config (in repo root)
- `.opentask-project`: Links repo to global task storage (in repo root)

**Q: Can I have multiple projects?**
A: Yes, define as many as you want in global config, switch with `--project` flag.

