# Git Manager Test Utilities

This package provides utilities for testing code that interacts with git repositories.

## Overview

The test utilities in this package allow you to:

1. Create temporary git repositories for testing
2. Manipulate git repositories (create files, branches, commits, etc.)
3. Set up complex test scenarios with multiple worktrees
4. Run assertions against git repositories

## Basic Usage

### Creating a Simple Test Repository

```go
import (
    "testing"
    "github.com/ingshtrom/git-manager/internal/testutil"
)

func TestMyFunction(t *testing.T) {
    // Set up a test repository
    repo, cleanup := testutil.SetupGitRepo(t)
    defer cleanup() // Always call cleanup to remove temporary files
    
    // Create a file and commit it
    repo.CreateFile(t, "example.txt", "Example content\n")
    repo.AddAndCommit(t, "Add example file", "example.txt")
    
    // Test your code that interacts with the repository
    // ...
}
```

### Using the Test Suite

For more complex scenarios, use the TestSuite which sets up a main repository with multiple worktrees:

```go
import (
    "testing"
    "github.com/ingshtrom/git-manager/internal/testutil"
)

func TestComplexScenario(t *testing.T) {
    // Set up a standard test suite with a main repository and two worktrees
    suite, cleanup := testutil.SetupStandardSuite(t)
    defer cleanup()
    
    // Access the main repository
    mainRepo := suite.RootRepo
    
    // Access worktree repositories
    featureRepo := suite.WorktreeRepos["feature-worktree"]
    bugfixRepo := suite.WorktreeRepos["bugfix-worktree"]
    
    // Test your code that interacts with these repositories
    // ...
}
```

## Available Utilities

### GitRepo

The `GitRepo` struct represents a git repository and provides methods for manipulating it:

- `CreateFile(t *testing.T, relativePath, content string)`: Creates a file in the repository
- `AddAndCommit(t *testing.T, message string, files ...string)`: Adds and commits files
- `CreateBranch(t *testing.T, branchName string)`: Creates a new branch
- `Checkout(t *testing.T, ref string)`: Checks out a branch or commit
- `CreateWorktree(t *testing.T, path, branch string) string`: Creates a new worktree
- `RunGit(t *testing.T, args ...string) string`: Runs a git command and returns its output
- `AssertBranchExists(t *testing.T, branch string)`: Asserts that a branch exists
- `AssertFileContent(t *testing.T, relativePath, expectedContent string)`: Asserts file content

### TestSuite

The `TestSuite` struct represents a test suite with multiple repositories:

- `RootRepo`: The main git repository
- `WorktreeRepos`: A map of worktree repositories by name
- `TempDir`: The temporary directory containing all repositories
- `AddWorktree(t *testing.T, name, branch string) *GitRepo`: Adds a new worktree to the suite

## Best Practices

1. Always defer the cleanup function to ensure temporary files are removed
2. Use `t.Helper()` in your test helper functions to get better error reporting
3. Create reusable test scenarios for common test cases
4. Use assertions to verify the state of repositories
5. Isolate tests by creating a new repository for each test

## Example: Integration Testing

```go
func TestIntegration(t *testing.T) {
    // Set up a standard test suite
    suite, cleanup := testutil.SetupStandardSuite(t)
    defer cleanup()
    
    // Test your code against the repositories
    result, err := YourFunction(suite.RootRepo.Path)
    if err != nil {
        t.Fatalf("YourFunction failed: %v", err)
    }
    
    // Assert the results
    if result != expectedResult {
        t.Errorf("Expected %v, got %v", expectedResult, result)
    }
    
    // Verify the state of the repositories
    suite.RootRepo.AssertBranchExists(t, "main")
    suite.RootRepo.AssertFileContent(t, "README.md", "# Test Repository\n")
}
``` 
