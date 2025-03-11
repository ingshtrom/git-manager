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
	Long: `Switch to a worktree in the current git repository.
This command will print the path to the specified worktree and instructions on how to switch to it.

When used with shell integration, it will automatically change the directory to the worktree.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		worktreeName := args[0]
		switchToWorktree(worktreeName)
	},
}

func init() {
	worktreeCmd.AddCommand(switchCmd)
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

	// Print information about the worktree
	fmt.Printf("Worktree path: %s\n", worktreePath)

	// Output the special command for shell integration to evaluate
	// This will be captured by the shell wrapper and executed
	fmt.Printf("git-manager-eval:cd %q\n", worktreePath)

	// Print instructions for users without shell integration
	fmt.Println("\nIf you're not using shell integration, run:")
	fmt.Printf("  cd %s\n", worktreePath)
	fmt.Println("\nTo enable shell integration, run:")
	fmt.Println("  git-manager shell [your-shell]")
}
