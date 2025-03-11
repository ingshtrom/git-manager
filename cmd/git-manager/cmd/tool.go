package cmd

import (
	"github.com/spf13/cobra"
)

var toolCmd = &cobra.Command{
	Use:     "tool",
	Aliases: []string{"t"},
	Short:   "Run various tools that help git-manager",
}

func init() {
	rootCmd.AddCommand(toolCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// toolCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// toolCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
