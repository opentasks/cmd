# Mise Task Runner Configuration

This project uses [mise](https://mise.jdx.dev/) for development task automation and tool version management.

## Setup

First, ensure you have mise installed:

```bash
curl https://mise.jdx.dev/install.sh | sh
```

Then trust the project configuration:

```bash
cd /path/to/opentask
mise trust
```

## Tools

The project pins the following tools:

- **Go 1.21** - Golang compiler and runtime

When you enter the project directory, mise automatically activates Go 1.21:

```bash
cd opentask
go version  # Should show go version go1.21.x
```

## Available Tasks

Run any task with `mise run <task>`:

### Build Tasks

**`mise run build`** - Build the opentask CLI binary
- Compiles `./cmd/opentask` to `./opentask`
- Includes all packages and dependencies

### Testing Tasks

**`mise run test`** - Run all tests across all packages
- Runs `go test -v ./...`
- Verbose output showing each test

**`mise run test-model`** - Test the model package
- `go test -v ./internal/model/...`

**`mise run test-storage`** - Test the storage package
- `go test -v ./internal/storage/...`

**`mise run test-config`** - Test the config package
- `go test -v ./internal/config/...`

**`mise run test-query`** - Test the query package
- `go test -v ./internal/query/...`

### Code Quality Tasks

**`mise run lint`** - Run Go linter (go vet)
- Checks for common errors and suspicious code
- `go vet ./...`

**`mise run fmt`** - Format code with gofmt
- Automatically formats all Go files
- `gofmt -w .`

### Utility Tasks

**`mise run clean`** - Clean build artifacts
- Removes `opentask` binary
- Removes `test_project/` and `demo_project/` directories
- Cleans up temporary test files

**`mise run demo`** - Run an interactive demo
- Creates `demo_project/` directory
- Creates a sample task hierarchy:
  - Epic: "Design Phase"
  - Plan: "Plan architecture" (parent: epic)
  - Research: "Research patterns" (parent: epic)
  - Story: "Implement core" (parent: epic, tagged: backend)
- Displays all created tasks
- Useful for testing and demonstrations

**`mise run help`** - Show all available tasks
- Lists all tasks with descriptions
- `mise tasks`

## Common Workflows

### Development Iteration

```bash
# Build the CLI
mise run build

# Run demo to test
mise run demo

# Format code
mise run fmt

# Check for issues
mise run lint

# Run all tests
mise run test

# Clean up
mise run clean
```

### Before Committing

```bash
# Format code
mise run fmt

# Run linter
mise run lint

# Run all tests
mise run test

# Build to ensure it compiles
mise run build
```

### Quick Testing of Specific Package

```bash
# Test just the storage package
mise run test-storage

# Test just the model package
mise run test-model

# Test query engine
mise run test-query
```

## Task Definitions

All tasks are defined in `.mise.toml` in the project root. Each task specifies:

- **description**: What the task does
- **run**: The shell command(s) to execute

Example:

```toml
[tasks.build]
description = "Build the opentask CLI"
run = "go build -o opentask ./cmd/opentask"
```

## Adding New Tasks

To add a new task:

1. Edit `.mise.toml`
2. Add a new `[tasks.name]` section
3. Set `description` and `run`
4. Trust the config if needed: `mise trust`
5. Test with `mise run name`

Example:

```toml
[tasks.mynewtask]
description = "Description of what it does"
run = "command to run"
```

## Environment Variables

If you need to set environment variables for tasks, you can add them to `.mise.toml`:

```toml
[env]
MYVAR = "value"
```

## IDE Integration

Most IDEs support mise tasks:

- **VS Code**: Use the Terminal or install a mise extension
- **GoLand/IntelliJ**: Recognized as external tools
- **Vim/Neovim**: Use `:!mise run taskname`

## Troubleshooting

**"mise not found"** - Install mise from https://mise.jdx.dev/

**"Config not trusted"** - Run `mise trust` in the project directory

**"Go version mismatch"** - Ensure you're in the project directory. `mise activate` should be in your shell config.

**Task not found** - Run `mise tasks` to list all available tasks

**Permission denied** - The binary may not be executable. Run `chmod +x opentask`

## More Information

- Official mise docs: https://mise.jdx.dev/
- Task runner guide: https://mise.jdx.dev/tasks.html
