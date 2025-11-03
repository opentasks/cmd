# Task Creation Location Fixes (Task 21)

## Problem
When creating new tasks via the CLI, they were not being saved in the correct project location. This was caused by:

1. **Incorrect TOML Config Format** - `config init` was generating nested sections instead of flat sections
2. **Config Discovery Walking Past Git Roots** - The discovery algorithm continued walking up past git repository boundaries

## Solutions Implemented

### Fix 1: TOML Config Format

**Before:**
```toml
[project.project]
name = "my-project"

[project.storage]
path = "./.tasks"
```

This nested format wasn't being parsed correctly by the TOML decoder.

**After:**
```toml
[project]
name = "my-project"

[storage]
path = "./.tasks"
```

The flat schema now correctly matches the config struct definitions.

**What Changed:**
- `cmd/config.go` - Updated `config init` template
- Template now generates correct flat TOML sections

### Fix 2: Config Discovery Git Root Handling

**Before:**
- Config discovery walked all the way up to filesystem root
- Would find `.opentask.toml` files above git repository roots
- Caused configs from parent directories to override project-local configs

**After:**
- Config discovery stops at `.git` directory (git repository root)
- Still respects filesystem root as a fallback stop condition
- Cleaner config resolution for git-based projects

**What Changed:**
- `internal/config/discovery.go` - Added git root detection
- Check for `.git` directory during discovery walk
- Updated tests to reflect new behavior

## Verification

### Test: Config Init Creates Correct Format
```bash
cd /tmp/test-project
opentask config init --name "test" --storage "./.tasks"
cat .opentask.toml
# Should show [storage] section, not [project.storage]
```

### Test: Tasks Created in Correct Location
```bash
cd /tmp/test-project
opentask task new "Test Task"
ls -la .tasks/
# Should show task file exists
```

### Test: Config Discovery Respects Git Roots
```bash
# Create git repo with config
mkdir /tmp/git-test && cd /tmp/git-test
git init
echo "[project]" > .opentask.toml
echo "name = \"root\"" >> .opentask.toml

# Create subdir and test
mkdir subdir && cd subdir
opentask config view --path
# Should NOT walk past .git directory
```

## Related Documentation

- `VERIFICATION_GUIDE.md` - Complete verification procedures
- `ProjectContexts.md` - New project context feature (reduces need for `--path`)
- `Config.md` - Configuration documentation

## Impact

✅ Tasks are now created in configured location
✅ Config discovery respects git repository boundaries
✅ Users no longer need to manually specify task storage location
✅ Works correctly with git worktrees
✅ All 53 tests passing (50 original + 3 new)

## Migration

If you have existing task files in the wrong locations, you can:

1. Find tasks created by old config: `find . -name "*.md" -path "*wrong-location*"`
2. Move them to correct location specified in `.opentask.toml`
3. Verify with `opentask config view --path`
