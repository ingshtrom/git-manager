package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/ingshtrom/git-manager/internal/worktree"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all worktrees in the current git-manager workspace",
	Long: `List all worktrees in the current git-manager workspace.
This command will display all worktrees, their paths, and their current branch.`,
	Run: func(cmd *cobra.Command, args []string) {
		listWorktrees()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func listWorktrees() {
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

	// Get worktree information
	worktrees, err := worktree.GetWorktreeInfo(gitDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Print worktree information in a tabular format
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	fmt.Fprintln(w, "PATH\tBRANCH\tCOMMIT")
	for _, wt := range worktrees {
		branchInfo := wt.Branch
		if branchInfo == "" {
			branchInfo = "(detached)"
		}
		if wt.IsBare {
			fmt.Fprintf(w, "%s\t%s\t%s\n", wt.Path, branchInfo, "n/a")
		} else {
			fmt.Fprintf(w, "%s\t%s\t%s\n", wt.Path, branchInfo, wt.Commit[:7])
		}
	}
	w.Flush()
}
