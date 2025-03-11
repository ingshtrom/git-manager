package worktree

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// setupTestRepo creates a temporary git repository for testing
// It returns the path to the repository and a cleanup function
func setupTestRepo(t *testing.T) (string, func()) {
	t.Helper()

	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "git-manager-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	// Initialize git repository
	cmd := exec.Command("git", "init", tempDir)
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("Failed to initialize git repository: %v", err)
	}

	// Set git config for tests
	cmd = exec.Command("git", "-C", tempDir, "config", "user.name", "Test User")
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("Failed to set git config user.name: %v", err)
	}

	cmd = exec.Command("git", "-C", tempDir, "config", "user.email", "test@example.com")
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("Failed to set git config user.email: %v", err)
	}

	// Create initial commit
	readmePath := filepath.Join(tempDir, "README.md")
	if err := os.WriteFile(readmePath, []byte("# Test Repository\n"), 0644); err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("Failed to create README.md: %v", err)
	}

	cmd = exec.Command("git", "-C", tempDir, "add", "README.md")
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("Failed to git add: %v", err)
	}

	cmd = exec.Command("git", "-C", tempDir, "commit", "-m", "Initial commit")
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("Failed to git commit: %v", err)
	}

	// Return cleanup function
	cleanup := func() {
		os.RemoveAll(tempDir)
	}

	return tempDir, cleanup
}

// TestGetWorktreeInfo tests the GetWorktreeInfo function
func TestGetWorktreeInfo(t *testing.T) {
	// Set up test repository
	repoPath, cleanup := setupTestRepo(t)
	defer cleanup()

	// Test GetWorktreeInfo
	worktrees, err := GetWorktreeInfo(repoPath)
	if err != nil {
		t.Fatalf("GetWorktreeInfo failed: %v", err)
	}

	// Verify we have at least one worktree (the main one)
	if len(worktrees) < 1 {
		t.Errorf("Expected at least one worktree, got %d", len(worktrees))
	}

	// Verify the main worktree path
	if worktrees[0].Path != repoPath {
		t.Errorf("Expected worktree path %s, got %s", repoPath, worktrees[0].Path)
	}
}

// TestIsGitRepository tests the IsGitRepository function
func TestIsGitRepository(t *testing.T) {
	// Set up test repository
	repoPath, cleanup := setupTestRepo(t)
	defer cleanup()

	// Test with a valid git repository
	if !IsGitRepository(repoPath) {
		t.Errorf("Expected %s to be a git repository", repoPath)
	}

	// Test with a non-git directory
	tempDir, err := os.MkdirTemp("", "not-a-git-repo-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	if IsGitRepository(tempDir) {
		t.Errorf("Expected %s to not be a git repository", tempDir)
	}
}
