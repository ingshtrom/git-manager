package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/ingshtrom/git-manager/internal/worktree"
	"github.com/spf13/cobra"
)

var (
	force        bool
	deleteBranch bool
)

var removeCmd = &cobra.Command{
	Use:   "remove [worktree-name]",
	Short: "Remove a worktree",
	Long: `Remove a worktree from the current git repository.
This command will remove the specified worktree.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		worktreeName := args[0]
		removeWorktree(worktreeName, force, deleteBranch)
	},
}

func init() {
	worktreeCmd.AddCommand(removeCmd)

	// Add flags
	removeCmd.Flags().BoolVarP(&force, "force", "f", false, "Force removal even if the worktree is dirty")
	removeCmd.Flags().BoolVarP(&deleteBranch, "delete-branch", "d", false, "Delete the branch associated with the worktree")
}

func removeWorktree(worktreeName string, force bool, deleteBranch bool) {
	// Get current directory
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
		os.Exit(1)
	}

	// Find the git directory
	gitDir, err := worktree.FindGitDir(currentDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Get the parent directory of the git directory
	parentDir := filepath.Dir(gitDir)

	// Create the worktree path
	worktreePath := filepath.Join(parentDir, worktreeName)

	// Check if the directory exists
	if _, err := os.Stat(worktreePath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: Worktree directory %s does not exist\n", worktreePath)
		os.Exit(1)
	}

	// Remove the worktree
	fmt.Printf("Removing worktree '%s'...\n", worktreeName)

	args := []string{"-C", gitDir, "worktree", "remove"}
	if force {
		args = append(args, "--force")
	}
	args = append(args, worktreePath)

	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error removing worktree: %v\n", err)
		os.Exit(1)
	}

	// Delete the branch if requested
	if deleteBranch {
		fmt.Printf("Deleting branch '%s'...\n", worktreeName)

		deleteCmd := exec.Command("git", "-C", gitDir, "branch", "-D", worktreeName)
		deleteCmd.Stdout = os.Stdout
		deleteCmd.Stderr = os.Stderr

		if err := deleteCmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error deleting branch: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Printf("\nWorktree '%s' removed successfully\n", worktreeName)
}
