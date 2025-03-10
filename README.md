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
```

## License
See the [LICENSE](LICENSE) file for details.
