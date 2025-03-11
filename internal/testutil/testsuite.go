package testutil

import (
	"os"
	"path/filepath"
	"testing"
)

// TestSuite represents a test suite for the entire codebase
type TestSuite struct {
	// RootRepo is the main git repository for testing
	RootRepo *GitRepo

	// WorktreeRepos contains additional worktree repositories
	WorktreeRepos map[string]*GitRepo

	// TempDir is the temporary directory containing all test repositories
	TempDir string
}

// NewTestSuite creates a new test suite with a main repository and optional worktrees
func NewTestSuite(t *testing.T) (*TestSuite, func()) {
	t.Helper()

	// Create a temporary directory for the test suite
	tempDir, err := os.MkdirTemp("", "git-manager-testsuite-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory for test suite: %v", err)
	}

	// Create the main repository
	rootRepoPath := filepath.Join(tempDir, "main-repo")
	if err := os.MkdirAll(rootRepoPath, 0755); err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("Failed to create main repo directory: %v", err)
	}

	// Initialize git in the main repository
	cmd := NewCommand("git", "init", rootRepoPath)
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("Failed to initialize git repository: %v", err)
	}

	// Set git config
	cmd = NewCommand("git", "-C", rootRepoPath, "config", "user.name", "Test User")
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("Failed to set git config user.name: %v", err)
	}

	cmd = NewCommand("git", "-C", rootRepoPath, "config", "user.email", "test@example.com")
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("Failed to set git config user.email: %v", err)
	}

	// Create the root repo object
	rootRepo := &GitRepo{
		Path: rootRepoPath,
	}

	// Create initial commit
	rootRepo.CreateFile(t, "README.md", "# Test Repository\n")
	rootRepo.AddAndCommit(t, "Initial commit", "README.md")

	// Create the test suite
	suite := &TestSuite{
		RootRepo:      rootRepo,
		WorktreeRepos: make(map[string]*GitRepo),
		TempDir:       tempDir,
	}

	// Return cleanup function
	cleanup := func() {
		os.RemoveAll(tempDir)
	}

	return suite, cleanup
}

// AddWorktree adds a new worktree to the test suite
func (s *TestSuite) AddWorktree(t *testing.T, name, branch string) *GitRepo {
	t.Helper()

	// Create a new worktree
	worktreePath := s.RootRepo.CreateWorktree(t, name, branch)

	// Create a GitRepo object for the worktree
	worktreeRepo := &GitRepo{
		Path: worktreePath,
	}

	// Add to the map of worktree repos
	s.WorktreeRepos[name] = worktreeRepo

	return worktreeRepo
}

// SetupStandardSuite creates a standard test suite with a main repository and two worktrees
func SetupStandardSuite(t *testing.T) (*TestSuite, func()) {
	t.Helper()

	// Create a new test suite
	suite, cleanup := NewTestSuite(t)

	// Create a feature branch
	suite.RootRepo.CreateBranch(t, "feature")

	// Add a worktree for the feature branch
	featureRepo := suite.AddWorktree(t, "feature-worktree", "feature")

	// Add some files to the feature worktree
	featureRepo.CreateFile(t, "feature.txt", "This is a feature\n")
	featureRepo.AddAndCommit(t, "Add feature", "feature.txt")

	// Create a bugfix branch
	suite.RootRepo.CreateBranch(t, "bugfix")

	// Add a worktree for the bugfix branch
	bugfixRepo := suite.AddWorktree(t, "bugfix-worktree", "bugfix")

	// Add some files to the bugfix worktree
	bugfixRepo.CreateFile(t, "bugfix.txt", "This is a bugfix\n")
	bugfixRepo.AddAndCommit(t, "Add bugfix", "bugfix.txt")

	return suite, cleanup
}
