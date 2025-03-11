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
	createBranch      bool
	baseBranch        string
	switchAfterCreate bool
)

var createCmd = &cobra.Command{
	Use:   "create [branch-name]",
	Short: "Create a new worktree",
	Long: `Create a new worktree in the current git repository.
This command will create a new worktree with the specified branch name.

When used with shell integration, it can automatically change the directory to the new worktree.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		branchName := args[0]
		createWorktree(branchName, createBranch, baseBranch, switchAfterCreate)
	},
}

func init() {
	worktreeCmd.AddCommand(createCmd)

	// Add flags
	createCmd.Flags().BoolVarP(&createBranch, "create-branch", "b", true, "Create a new branch for the worktree")
	createCmd.Flags().StringVarP(&baseBranch, "base", "", "main", "Base branch to create the new branch from (used with --create-branch)")
	createCmd.Flags().BoolVarP(&switchAfterCreate, "switch", "s", true, "Switch to the new worktree after creation")
}

func createWorktree(branchName string, createBranch bool, baseBranch string, switchAfterCreate bool) {
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
	worktreePath := filepath.Join(parentDir, branchName)

	// Check if the directory already exists
	if _, err := os.Stat(worktreePath); err == nil {
		fmt.Fprintf(os.Stderr, "Error: Directory %s already exists\n", worktreePath)
		os.Exit(1)
	}

	var cmd *exec.Cmd

	if createBranch {
		// Create a new branch and worktree
		fmt.Printf("Creating new branch '%s' based on '%s' and adding worktree...\n", branchName, baseBranch)
		cmd = exec.Command("git", "-C", gitDir, "worktree", "add", "-b", branchName, worktreePath, baseBranch)
	} else {
		// Add worktree for existing branch
		fmt.Printf("Adding worktree for branch '%s'...\n", branchName)
		cmd = exec.Command("git", "-C", gitDir, "worktree", "add", worktreePath, branchName)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating worktree: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nWorktree created successfully at %s\n", worktreePath)

	if switchAfterCreate {
		// Output the special command for shell integration to evaluate
		// This will be captured by the shell wrapper and executed
		fmt.Printf("git-manager-eval:cd %q\n", worktreePath)
	}

	// Print instructions for users without shell integration
	fmt.Println("\nIf you're not using shell integration, run:")
	fmt.Printf("  cd %s\n", worktreePath)

	if !switchAfterCreate {
		fmt.Println("\nTo automatically switch to new worktrees, use the --switch flag:")
		fmt.Printf("  git-manager create --switch %s\n", branchName)
	}

	fmt.Println("\nTo enable shell integration, run:")
	fmt.Println("  git-manager shell [your-shell]")
}
