package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ingshtrom/git-manager/internal/worktree"
	"github.com/spf13/cobra"
)

var switchCmd = &cobra.Command{
	Use:   "switch [worktree-name]",
	Short: "Switch to a worktree",
	Long: `Switch to a worktree in the current git-manager workspace.
This command will print the path to the specified worktree and instructions on how to switch to it.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		worktreeName := args[0]
		switchToWorktree(worktreeName)
	},
}

func init() {
	rootCmd.AddCommand(switchCmd)
}

func switchToWorktree(worktreeName string) {
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

	// Get worktree information to verify it's a valid worktree
	worktrees, err := worktree.GetWorktreeInfo(gitDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing worktrees: %v\n", err)
		os.Exit(1)
	}

	// Check if the worktree exists in the list
	found := false
	for _, wt := range worktrees {
		if wt.Path == worktreePath {
			found = true
			break
		}
	}

	if !found {
		fmt.Fprintf(os.Stderr, "Error: %s is not a valid worktree\n", worktreePath)
		os.Exit(1)
	}

	// Print instructions for switching to the worktree
	fmt.Printf("Worktree path: %s\n", worktreePath)
	fmt.Println("\nTo switch to this worktree, run:")
	fmt.Printf("  cd %s\n", worktreePath)

	// Since we can't actually change the user's current directory from within the program,
	// we just provide instructions on how to do it.
}
