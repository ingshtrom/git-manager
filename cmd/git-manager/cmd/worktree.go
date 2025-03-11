package cmd

import (
	"fmt"
	"os"

	"github.com/ingshtrom/git-manager/internal/worktree"
	"github.com/spf13/cobra"
)

// worktreeCmd represents the worktree command
var worktreeCmd = &cobra.Command{
	Use:     "worktree",
	Aliases: []string{"wt", "w"},
	Short:   "Manage git worktrees (wt, w)",
	Long: `Manage git worktrees.

This command allows you to manage your git worktrees.
You can create, list, switch, and remove worktrees.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		dir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to determine your current directory. Please ensure you have permissions to access this directory and try again: %v", err)
		}

		if !worktree.IsGitRepository(dir) {
			return fmt.Errorf("this command must be run from within a git repository. Please navigate to a git repository and try again")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(worktreeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// worktreeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// worktreeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
