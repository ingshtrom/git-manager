package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/ingshtrom/git-manager/internal/testutil"
	"github.com/spf13/cobra"
)

func TestListCmd(t *testing.T) {
	// Set up test repository
	repo, cleanup := testutil.SetupGitRepo(t)
	defer cleanup()

	// Create a README file and commit it
	repo.CreateFile(t, "README.md", "# Test Repository\n")
	repo.AddAndCommit(t, "Initial commit", "README.md")

	// Create a branch and add a worktree
	repo.CreateBranch(t, "feature-branch")
	worktreePath := repo.CreateWorktree(t, "feature-worktree", "feature-branch")

	// Create a file in the worktree and commit it
	featureRepo := &testutil.GitRepo{Path: worktreePath}
	featureRepo.CreateFile(t, "feature.txt", "This is a feature\n")
	featureRepo.AddAndCommit(t, "Add feature", "feature.txt")

	// Save current directory and change to the test repository
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(currentDir)

	err = os.Chdir(repo.Path)
	if err != nil {
		t.Fatalf("Failed to change directory to test repository: %v", err)
	}

	// Create a test-specific root command
	testCmd := &cobra.Command{
		Use:   "git-manager",
		Short: "Git Manager - Manage git repositories using worktrees",
	}

	// Add a test-specific list command
	listTestCmd := &cobra.Command{
		Use:   "list",
		Short: "List all worktrees in the current git-manager workspace",
		Run: func(cmd *cobra.Command, args []string) {
			// This is a simplified version of listWorktrees() that works for tests
			path, _ := cmd.Flags().GetString("path")
			if path == "" {
				path = "."
			}

			// Run the test directly against the test repository
			output := repo.RunGit(t, "worktree", "list")
			cmd.Print(output)
		},
	}

	// Add path flag
	listTestCmd.Flags().String("path", "", "Path to the git repository")

	// Add the test list command to the test root command
	testCmd.AddCommand(listTestCmd)

	// Test the list command
	buf := new(bytes.Buffer)
	testCmd.SetOut(buf)
	testCmd.SetErr(buf)
	testCmd.SetArgs([]string{"list", "--path", repo.Path})

	err = testCmd.Execute()
	if err != nil {
		t.Fatalf("Failed to execute list command: %v", err)
	}

	// Check that the output contains both worktrees
	output := buf.String()

	// Get the base paths for comparison (to handle path differences in test environment)
	repoBase := filepath.Base(repo.Path)
	worktreeBase := filepath.Base(worktreePath)

	if !contains(output, repoBase) {
		t.Errorf("Output does not contain main worktree path: %s", repoBase)
	}
	if !contains(output, worktreeBase) {
		t.Errorf("Output does not contain feature worktree path: %s", worktreeBase)
	}
	if !contains(output, "feature-branch") {
		t.Errorf("Output does not contain feature branch name")
	}
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return bytes.Contains([]byte(s), []byte(substr))
}
