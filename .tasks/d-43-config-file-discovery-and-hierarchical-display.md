---
id: 43
title: Config file discovery and hierarchical display
type: decision
status: done
tags: [config, discovery, cli]
relationships: []
createdAt: "2025-11-03T08:00:00Z"
updatedAt: "2025-11-03T08:00:00Z"
---

# Decision: Config File Discovery and Hierarchical Display

## Problem

How should opentask discover and display configuration files to users? Configuration can come from multiple sources at different directory levels and needs to be merged in the correct priority order.

## Decision

### Discovery Method

**Walk the directory tree upward from the current working directory, stopping only at the filesystem root. Additionally check for a user global config at `~/.config/opentask/config.toml`.**

Config files are discovered in this order (highest to lowest priority):
1. `.opentask.toml` in current directory and parent directories (up to filesystem root)
2. User global config at `~/.config/opentask/config.toml`
3. Built-in defaults (virtual layer)

Files found closer to the current directory have higher priority and override files further up the tree.

### Why This Approach

- **No git root stopping**: Allows configs to be shared across git repository boundaries, enabling monorepo and multi-project setups
- **XDG compliance**: User global config at `~/.config/opentask/config.toml` follows XDG Base Directory specification
- **Flexible scoping**: Teams can put shared defaults at project root, subteams can override with their own configs
- **Backward compatible**: Works with single-file or multi-file configurations seamlessly

### Config File Display

**Show discovered configs as a straight vertical list with tree connectors and merge flow arrows.**

Example output from `opentask config view`:

```
├── ./.opentask.toml
│   ↓
├── ../.opentask.toml
│   ↓
├── ../../../.opentask.toml
│   ↓
├── ~/.config/opentask/config.toml
│   ↓
└── (builtin) defaults
```

### Why This Display Format

- **Visual clarity**: Tree connectors (`├──`, `└──`, `│`) show structure without nesting
- **Flow indication**: Arrow (`↓`) clearly shows merge direction (configs above override below)
- **Scan-friendly**: Straight vertical list is easy to read and understand at a glance
- **Priority obvious**: First item is highest priority, last is lowest (defaults)
- **Professional**: Uses standard ASCII conventions

## Implementation Details

### Discovery Algorithm

1. Start from current working directory
2. Walk upward through parent directories
3. At each level, check for `.opentask.toml` file
4. Continue until filesystem root is reached
5. Additionally check for user global config at `~/.config/opentask/config.toml`
6. Return all found files in order: closest first, then user config

### Merging Order

Configs are merged left-to-right, with later configs overriding earlier ones:
- Start with built-in defaults
- Apply user global config (if exists)
- Apply discovered project configs from furthest to closest
- Result: closest config has final say on all values

### Relative Path Display

Config files are displayed as paths relative to the current working directory:
- `./.opentask.toml` - in current directory
- `../.opentask.toml` - in parent directory  
- `../../parent/.opentask.toml` - further up
- `~/.config/opentask/config.toml` - user home directory (with `~` expansion)

This is more readable than absolute paths and helps users understand the hierarchy relative to where they are working.

## Related Tasks

- Task 1: Hierarchical config loading via `LoadConfigHierarchical()`
- Task 2: Config discovery via `DiscoverConfigFiles()`
- Task 3: Config merging via `MergeConfigs()`
- Task 4: CLI initialization in `initializeStorage()`
- Task 5: Config view display in `config view` command

## Changes Made

- Changed config filename from `opentask.toml` to `.opentask.toml` (hidden file convention)
- Added user global config discovery at `~/.config/opentask/config.toml`
- Implemented hierarchical discovery walking up directory tree
- Removed git repository root as stop condition
- Implemented `DiscoverConfigFiles()` returning files closest-first
- Implemented `MergeConfigs()` to merge multiple configs with proper override priority
- Implemented hierarchical `LoadConfigHierarchical()` for convenient loading + merging
- Added `DiscoverAndAnalyze()` for config analysis and display
- Implemented file tree display with ASCII characters and merge indicators
- Updated documentation to reflect actual behavior

## Testing

All changes covered by comprehensive tests:
- Config discovery with parent directories
- Config merging with override priority
- Hierarchical loading with explicit config paths
- File tree display formatting
- User config discovery
- Relative path resolution

Tests verify:
- Files are discovered in correct order (closest first)
- Merging applies configs in correct order
- Explicit `--config` flag overrides discovery
- User config is included as lowest priority
- Built-in defaults are available as fallback
