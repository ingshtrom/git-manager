package internal_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ingshtrom/git-manager/internal/testutil"
	"github.com/ingshtrom/git-manager/internal/worktree"
)

// TestWorktreeIntegration tests the worktree package with a real git repository
func TestWorktreeIntegration(t *testing.T) {
	// Set up a standard test suite
	suite, cleanup := testutil.SetupStandardSuite(t)
	defer cleanup()

	// Test GetWorktreeInfo
	worktrees, err := worktree.GetWorktreeInfo(suite.RootRepo.Path)
	if err != nil {
		t.Fatalf("GetWorktreeInfo failed: %v", err)
	}

	// Verify we have the expected number of worktrees (main + 2 additional)
	if len(worktrees) != 3 {
		t.Errorf("Expected 3 worktrees, got %d", len(worktrees))
	}

	// Verify the main worktree
	foundMain := false
	for _, wt := range worktrees {
		if wt.Path == suite.RootRepo.Path {
			foundMain = true
			if wt.Branch != "main" && wt.Branch != "master" {
				t.Errorf("Expected main worktree branch to be main or master, got %s", wt.Branch)
			}
		}
	}
	if !foundMain {
		t.Errorf("Main worktree not found in worktree list")
	}

	// Verify the feature worktree
	foundFeature := false
	featureWorktreePath := suite.WorktreeRepos["feature-worktree"].Path
	for _, wt := range worktrees {
		if wt.Path == featureWorktreePath {
			foundFeature = true
			if wt.Branch != "feature" {
				t.Errorf("Expected feature worktree branch to be feature, got %s", wt.Branch)
			}
		}
	}
	if !foundFeature {
		t.Errorf("Feature worktree not found in worktree list")
	}

	// Verify the bugfix worktree
	foundBugfix := false
	bugfixWorktreePath := suite.WorktreeRepos["bugfix-worktree"].Path
	for _, wt := range worktrees {
		if wt.Path == bugfixWorktreePath {
			foundBugfix = true
			if wt.Branch != "bugfix" {
				t.Errorf("Expected bugfix worktree branch to be bugfix, got %s", wt.Branch)
			}
		}
	}
	if !foundBugfix {
		t.Errorf("Bugfix worktree not found in worktree list")
	}

	// Test IsGitRepository
	if !worktree.IsGitRepository(suite.RootRepo.Path) {
		t.Errorf("Expected %s to be a git repository", suite.RootRepo.Path)
	}

	// Test with a non-git directory
	tempDir, err := os.MkdirTemp("", "not-a-git-repo-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	if worktree.IsGitRepository(tempDir) {
		t.Errorf("Expected %s to not be a git repository", tempDir)
	}

	// Test FindGitDir
	gitDir, err := worktree.FindGitDir(suite.RootRepo.Path)
	if err != nil {
		t.Fatalf("FindGitDir failed: %v", err)
	}

	expectedGitDir := filepath.Join(suite.RootRepo.Path, ".git")
	if gitDir != expectedGitDir {
		t.Errorf("Expected git dir %s, got %s", expectedGitDir, gitDir)
	}
}
