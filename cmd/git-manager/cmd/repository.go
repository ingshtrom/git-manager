package cmd

import (
	"github.com/spf13/cobra"
)

// repositoryCmd represents the repository command
var repositoryCmd = &cobra.Command{
	Use:     "repository",
	Aliases: []string{"repo", "r"},
	Short:   "Manage git repositories (repo, r)",
	Long: `Manage git repositories.

This command allows you to manage your git repositories.
You can initialize a new git repository, switch between worktrees, and more.`,
}

func init() {
	rootCmd.AddCommand(repositoryCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// repositoryCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// repositoryCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
