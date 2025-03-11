package cmd

import (
	"github.com/spf13/cobra"
)

// hoist the worktreeListCmd to the rootCmd for convenience
var listCmd = &cobra.Command{
	Use:               worktreeListCmd.Use,
	Aliases:           worktreeListCmd.Aliases,
	Short:             worktreeListCmd.Short,
	Long:              worktreeListCmd.Long,
	PersistentPreRunE: worktreeListCmd.PersistentPreRunE,
	Run:               worktreeListCmd.Run,
}

func init() {
	rootCmd.AddCommand(listCmd)
}
