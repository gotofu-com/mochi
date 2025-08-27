# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

Mochi is a Go-based CLI tool for managing release notes in a monorepo. It provides commands to create change entries and manage releases with support for multiple targets and configurable change types.

## Common Commands

### Build and Run
```bash
# Build the project
go build

# Run the CLI
go run main.go

# Install locally
go install

# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test -run TestName ./...
```

### Development Workflow
```bash
# Create a new release note entry
mochi new [type] [target] [message]

# Start a release
mochi release start [target]

# Preview release notes
mochi release preview

# Finish a release
mochi release finish
```

## Architecture

### Core Components

1. **Command Layer** (`/cmd/`): Cobra-based CLI commands
   - `root.go`: Main command setup with debug flag support
   - `new.go`: Creates new change entries with interactive prompts
   - `release.go`: Manages release lifecycle (start/preview/finish)

2. **Domain Models** (`/domain/`): Core business entities
   - `change.go`: Individual change entry with type, target, message, and optional ticket ID
   - `change_type.go`: Types of changes (feature, bugfix, doc, removal, misc)
   - `target.go`: Release targets with tag and ticket prefix support
   - `version.go`: Semantic versioning logic with date-based format support
   - `tag.go`: Git tag representation combining target and version
   - `release.go`: Release notes aggregation and rendering

3. **Services** (`/change/`, `/release/`, `/tag/`, `/version/`): Business logic
   - Each service handles operations for its respective domain entity
   - Services interact with git and file system for persistence

4. **Configuration** (`/config/`): 
   - Looks for `.mochi/config.yaml` recursively up the directory tree
   - Supports environment variables (e.g., `GITHUB_TOKEN`)
   - Defaults provided for change types and base branch

5. **Utilities** (`/utils/git/`): Git operations wrapper
   - Branch management, tag operations, and repository state checks

### Key Design Patterns

- **Configuration Discovery**: Recursively searches for `.mochi/config.yaml` from current directory upward
- **Change Storage**: Changes are stored as markdown files with YAML frontmatter in `.mochi/` directories
- **Version Format**: Supports both semantic (x.y.z) and date-based (YYYY.MM.PATCH) versioning
- **Interactive CLI**: Uses promptui for user-friendly prompts when arguments not provided
- **Branch Naming**: Release branches follow `release/<target>/<version>` pattern
- **Ticket ID Extraction**: Automatically extracts ticket IDs from branch names using configurable prefixes

### Release Workflow

1. Developer creates change entries with `mochi new` during development
2. Release manager starts release with `mochi release start [target]`
3. Preview notes with `mochi release preview`
4. Finalize with `mochi release finish` which:
   - Gathers all change files for the target
   - Creates a commit removing change files
   - Tags the commit
   - Merges to base branch

## Configuration

The tool uses Viper for configuration management with the following hierarchy:
1. `.mochi/config.yaml` (searched recursively)
2. Environment variables
3. Default values

Example config structure:
```yaml
baseBranch: main
baseTicketUrl: https://jira.example.com/browse/
types:
  - id: feature
    name: Feature
    title: Features
targets:
  - id: backend
    name: Backend Service
    tagPrefix: backend
    ticketPrefix: BE
```

## Testing Approach

The project follows Go testing conventions with test files alongside implementation files. Run tests with standard Go test commands.

## Git Integration

The tool integrates deeply with git for:
- Branch management (creating release branches)
- Tag creation and parsing
- Extracting ticket IDs from branch names
- Ensuring clean working directory before operations