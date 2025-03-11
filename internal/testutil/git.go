package testutil

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// GitRepo represents a test git repository
type GitRepo struct {
	Path string
}

// SetupGitRepo creates a new git repository for testing
func SetupGitRepo(t *testing.T) (*GitRepo, func()) {
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

	repo := &GitRepo{
		Path: tempDir,
	}

	// Return cleanup function
	cleanup := func() {
		os.RemoveAll(tempDir)
	}

	return repo, cleanup
}

// CreateFile creates a file in the repository with the given content
func (r *GitRepo) CreateFile(t *testing.T, relativePath, content string) {
	t.Helper()

	fullPath := filepath.Join(r.Path, relativePath)

	// Ensure directory exists
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("Failed to create directory %s: %v", dir, err)
	}

	// Write file
	if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write file %s: %v", fullPath, err)
	}
}

// AddAndCommit adds and commits the specified files
func (r *GitRepo) AddAndCommit(t *testing.T, message string, files ...string) {
	t.Helper()

	// Add files
	args := append([]string{"-C", r.Path, "add"}, files...)
	cmd := exec.Command("git", args...)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to git add: %v", err)
	}

	// Commit
	cmd = exec.Command("git", "-C", r.Path, "commit", "-m", message)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to git commit: %v", err)
	}
}

// CreateBranch creates a new branch
func (r *GitRepo) CreateBranch(t *testing.T, branchName string) {
	t.Helper()

	cmd := exec.Command("git", "-C", r.Path, "branch", branchName)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to create branch %s: %v", branchName, err)
	}
}

// Checkout checks out the specified branch or commit
func (r *GitRepo) Checkout(t *testing.T, ref string) {
	t.Helper()

	cmd := exec.Command("git", "-C", r.Path, "checkout", ref)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to checkout %s: %v", ref, err)
	}
}

// CreateWorktree creates a new worktree
func (r *GitRepo) CreateWorktree(t *testing.T, path, branch string) string {
	t.Helper()

	worktreePath := filepath.Join(r.Path, "..", filepath.Base(path))

	var cmd *exec.Cmd
	if branch == "" {
		cmd = exec.Command("git", "-C", r.Path, "worktree", "add", worktreePath)
	} else {
		cmd = exec.Command("git", "-C", r.Path, "worktree", "add", "-b", branch, worktreePath)
	}

	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to create worktree at %s: %v", worktreePath, err)
	}

	return worktreePath
}

// RunGit runs a git command in the repository and returns its output
func (r *GitRepo) RunGit(t *testing.T, args ...string) string {
	t.Helper()

	cmdArgs := append([]string{"-C", r.Path}, args...)
	cmd := exec.Command("git", cmdArgs...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run git %v: %v\nOutput: %s", args, err, output)
	}

	return string(output)
}

// AssertBranchExists checks if a branch exists
func (r *GitRepo) AssertBranchExists(t *testing.T, branch string) {
	t.Helper()

	cmd := exec.Command("git", "-C", r.Path, "rev-parse", "--verify", branch)
	if err := cmd.Run(); err != nil {
		t.Errorf("Branch %s does not exist: %v", branch, err)
	}
}

// AssertFileContent checks if a file has the expected content
func (r *GitRepo) AssertFileContent(t *testing.T, relativePath, expectedContent string) {
	t.Helper()

	fullPath := filepath.Join(r.Path, relativePath)
	content, err := os.ReadFile(fullPath)
	if err != nil {
		t.Fatalf("Failed to read file %s: %v", fullPath, err)
	}

	if string(content) != expectedContent {
		t.Errorf("File %s content mismatch\nExpected: %s\nActual: %s",
			relativePath, expectedContent, string(content))
	}
}
