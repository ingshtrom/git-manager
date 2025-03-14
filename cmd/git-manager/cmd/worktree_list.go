package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/ingshtrom/git-manager/internal/worktree"
	"github.com/spf13/cobra"
)

var worktreeListCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list", "l"},
	Short:   "List all worktrees in the current git repository (list, l)",
	Long: `List all worktrees in the current git repository.
This command will display all worktrees, their paths, and their current branch.`,
	Run: func(cmd *cobra.Command, args []string) {
		listWorktrees()
	},
}

func init() {
	worktreeCmd.AddCommand(worktreeListCmd)
}

func listWorktrees() {
	// Get current directory
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
		os.Exit(1)
	}

	// Get worktree information
	worktrees, err := worktree.GetWorktreeInfo(currentDir)
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
