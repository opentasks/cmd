---
id: 45
title: Set up CI/CD pipeline with GitHub Actions
type: story
status: todo
tags: [devops, ci-cd, github-actions, testing, release]
relationships: []
createdAt: "2025-11-03T08:35:00Z"
updatedAt: "2025-11-03T09:00:00Z"
---

# Story: Set Up CI/CD Pipeline with GitHub Actions

## Problem

Currently there's no automated testing, linting, or build verification when code is pushed. This means:
- Broken builds can be committed to main
- Test regressions aren't caught automatically
- Code quality isn't enforced
- Release process is manual
- No confidence in pull requests
- No automated artifact builds

## Acceptance Criteria

### Required

- [ ] Tests run automatically on push and PRs
- [ ] Build verification (Go build succeeds)
- [ ] Test coverage reporting
- [ ] Linting checks (golangci-lint or gofmt)
- [ ] Workflow blocks merging on failure
- [ ] Clear success/failure status on commits
- [ ] Automated release PRs via release-please
- [ ] Automated artifact building with go-releaser
- [ ] GitHub releases with binaries for multiple platforms

### Nice to Have

- [ ] Code coverage badge
- [ ] Dependency vulnerability scanning
- [ ] Documentation building and deployment
- [ ] Performance regression detection
- [ ] Docker image builds and pushes

## Implementation Plan

### Phase 1: Basic Test & Build Pipeline

1. Create `.github/workflows/pr.yml`
   - Trigger on: pull_request
   - Run on: ubuntu-latest, use mise action
   - Steps:
	   - Checkout code
	   - setup mise
     - Run `go test ./...`
     - Run `go build ./cmd`
     - Upload test results
   - Run golangci-lint on all Go files
   - Check gofmt compliance
   - Report style violations use problem matchers

3. Configure branch protection
   - Require test status checks
   - Require PR reviews (optional)

### Phase 2: Coverage & Reporting

1. Add test coverage collection
   - Use `go test -coverprofile` 
   - Generate coverage report
   - Upload to Codecov or similar

2. Create coverage badge in README

3. Track trends over time

### Phase 3: Automated Release Pipeline

#### Release Management with release-please

 1. Create `.github/workflows/release.yml`
	- Trigger on: push to main
	- Uses `google-github-actions/release-please-action@v3`
	- Automatically creates release PRs with:
		- Bumped version (semantic versioning)
		- Generated changelog
		- Updated package.json (if needed)
	- On merge of release PR, creates GitHub release
2. Configure release-please
	- Use conventional commits for auto-changelog
	- Semantic versioning (major.minor.patch)
	- Track version in VERSION file or git tags
3. Adjust agents.md and opentask skill so that:
	- any git commits follow conventional commits. 
	- if it merges back to branches, use squash merge.
	- never merge back to mainline, always use PR

#### Binary Artifact Building with go-releaser

1. Update release.yml workflow 
   - Uses `goreleaser/goreleaser-action@v4`
   - for non release merges, upload patch release build
   - for release merges, upload minor release build
   - Builds binaries for multiple platforms:
     - linux/amd64, linux/arm64
     - darwin/amd64, darwin/arm64 (macOS)
     - windows/amd64, windows/arm64
   - Creates checksums for verification
   - Uploads to GitHub release assets

2. Create `.goreleaser.yaml` config
   - Binary name: `opentask`
   - Output directory: `dist/`
   - Include version in binaries
   - Create SBOMs (Software Bill of Materials)
   - Sign artifacts (optional)

3. Benefits:
   - Single command builds all platforms
   - Consistent artifact naming
   - Automatic checksum generation
   - GitHub release integration
   - Clean dist/ folder management

## Workflow Structure

```
.github/
└── workflows/
    ├── test.yml              # Run tests and build
    ├── lint.yml              # Code quality checks
    ├── release-please.yml    # Automated release PRs
    └── goreleaser.yml        # Build release artifacts
    
.goreleaser.yaml             # Go-releaser configuration
```

### test.yml Outline

```yaml
name: Tests & Build
on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main, develop]

jobs:
  test:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        go-version: ['1.21', '1.22']
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v4
        with:
          go-version: ${{ matrix.go-version }}
      - run: go test -v -coverprofile=coverage.out ./...
      - run: go build ./cmd/opentask
      - uses: codecov/codecov-action@v3
        with:
          files: ./coverage.out
```

### lint.yml Outline

```yaml
name: Lint & Format
on: [push, pull_request]

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: golangci/golangci-lint-action@v3
```

### release-please.yml Outline

```yaml
name: Release Please
on:
  push:
    branches: [main]

permissions:
  contents: write
  pull-requests: write

jobs:
  release-please:
    runs-on: ubuntu-latest
    steps:
      - uses: google-github-actions/release-please-action@v3
        with:
          release-type: go
          package-name: opentask
```

### goreleaser.yml Outline

```yaml
name: GoReleaser Build
on:
  release:
    types: [created]

permissions:
  contents: write

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - uses: goreleaser/goreleaser-action@v4
        with:
          version: latest
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

### .goreleaser.yaml Config Outline

```yaml
project_name: opentask
dist: dist

builds:
  - id: opentask
    main: ./cmd/opentask/main.go
    binary: opentask
    goos:
      - linux
      - darwin
      - windows
    goarch:
      - amd64
      - arm64
    ldflags:
      - -s -w -X main.version={{.Version}}

archives:
  - format: tar.gz
    format_overrides:
      - goos: windows
        format: zip

checksum:
  name_template: 'checksums.txt'

changelog:
  use: conventional
```

## Release Workflow

1. **Developer pushes commits** with conventional commit messages
2. **release-please action** detects commits and creates release PR
3. **Release PR includes**:
   - Version bump (semantic versioning)
   - Auto-generated changelog
   - Updated files
4. **On merge of release PR**:
   - GitHub release is created
   - Tag is pushed
5. **goreleaser action** triggered by release
6. **go-releaser builds**:
   - Binaries for all platforms
   - Checksums
   - Archives (tar.gz for Unix, zip for Windows)
7. **Artifacts uploaded** to GitHub release

## Success Metrics

- [ ] All tests pass on every push
- [ ] Build failures caught before merge
- [ ] PRs show CI status
- [ ] Coverage tracked and reported
- [ ] Release PRs created automatically
- [ ] Binaries built for all platforms on release
- [ ] GitHub releases populated with assets
- [ ] No manual release process needed

## Related Tasks

- Task 44: Add task content editing to CLI (will need testing)
- All future features will require CI verification

## Testing Strategy

- Unit tests: run via `go test ./...`
- Integration tests: included in unit test suite
- Build verification: `go build ./cmd/opentask`
- No e2e tests needed yet

## References

- GitHub Actions documentation: https://docs.github.com/actions
- release-please: https://github.com/googleapis/release-please
- go-releaser: https://goreleaser.com/
- golangci-lint: https://golangci-lint.run/
- Go testing: https://golang.org/pkg/testing/
- Codecov: https://codecov.io/
