# git-manager
Manage git repositories using Git worktrees.

## Overview
Git worktrees are awesome, but to use them effectively, some prior organization helps. `git-manager` is a command-line tool that helps you organize and manage multiple git repositories using git worktrees. It puts them in a special directory to ensure that it does not conflict with current git repositories that might be managed in different ways.

## Features
- Create and manage git repositories
- Create and manage git worktrees for a given repository
- Navigate between different worktrees and repositories
- Organize repositories efficiently
- Streamline git workflow
- Automatic directory switching with shell integration

## Installation
```bash
go install github.com/ingshtrom/git-manager@latest
```

## Usage
```bash
# Initialize a new git-manager workspace
git-manager init [repository-url]

# List all worktrees
git-manager list

# Create a new worktree
git-manager create [branch-name]

# Switch to a worktree
git-manager switch [worktree-name]

# Remove a worktree
git-manager remove [worktree-name]

# Generate shell integration scripts
git-manager shell [shell-type]
```

## Shell Integration

Since a command-line tool cannot directly change the parent shell's directory, `git-manager` provides shell integration to make directory switching seamless.

### Setting Up Shell Integration

Generate and install shell integration for your preferred shell:

```bash
# For Bash
git-manager shell bash >> ~/.bashrc

# For Zsh
git-manager shell zsh >> ~/.zshrc

# For Fish
git-manager shell fish > ~/.config/fish/functions/git-manager.fish

# For Nushell
git-manager shell nushell >> ~/.config/nushell/config.nu
```

### How It Works

The shell integration works by:

1. Wrapping the `git-manager` command with a shell function
2. Capturing the output of the command
3. Looking for special `git-manager-eval:` prefixed lines
4. Evaluating any commands with this prefix in the shell
5. Displaying all other output normally

This allows commands like `switch` to change your current directory automatically when shell integration is enabled.

### Using Shell Integration

Once shell integration is set up, you can use `git-manager` (or the shorter alias `gm`) as usual:

```bash
# Switch to a worktree and automatically change directory
git-manager switch feature-branch

# Or use the shorter alias
gm switch feature-branch
```

## Development

### Prerequisites
- [Go](https://golang.org/doc/install) (version 1.24 or later)
- [Task](https://taskfile.dev/#/installation) (optional, for running tasks)

### Using Taskfile
This project includes a Taskfile.yml for common development tasks. If you have [Task](https://taskfile.dev) installed, you can use the following commands:

```bash
# List all available tasks
task

# Build the application
task build

# Run the application
task run

# Run tests
task test

# Run tests with coverage
task test:cover

# Format code
task fmt

# Lint code
task lint

# Install git-manager globally
task install

# Clean build artifacts
task clean

# Run in development mode with hot reload (requires Air)
task dev
```

### Testing

You can run tests locally:

```bash
# Run all tests
task test

# Run tests with coverage
task test:cover
```

#### Docker-Based Testing

For consistent testing across environments, you can run tests in Docker:

```bash
# Run all tests in Docker
task test:docker

# Run specific tests in Docker
task test:docker:specific -- ./internal/worktree

# Run tests with coverage in Docker
task test:docker:cover
```

For more information about Docker-based testing, see [Docker Testing Documentation](docs/docker-testing.md).

## License
See the [LICENSE](LICENSE) file for details.
