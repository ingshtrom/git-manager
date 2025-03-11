# STATUS
> ⚠️ **WARNING**: This project is no longer maintained. The core functionality of automatically changing directories when switching worktrees is not possible due to fundamental limitations in how shells interact with child processes. A child process cannot modify its parent shell's working directory. While shell integration scripts can work around this, they add significant complexity and are not reliable across all environments. For a simpler solution, consider using git worktrees directly or building your own aliases/commands to provide the experience you want.


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

## Shell Integration

Since a command-line tool cannot directly change the parent shell's directory, `git-manager` provides shell integration to make directory switching seamless.

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
- [Task](https://taskfile.dev/#/installation) (for running tasks)

### Getting Started

1. get the dev task running, this gets the binary built in near-real-time. It ouputs to the `./build` directory:
```bash
$ task dev
```
2. [install the shell integration](#using-shell-integration)
3. Modify your shell integration script to point to `<cwd>/build/git-manager` instead of the globally installed git-manager
4. restart your shell so that the newest shell integration is active
5. test with `gm` command

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
