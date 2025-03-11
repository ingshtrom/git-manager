package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:     worktreeAddCmd.Use,
	Aliases: worktreeAddCmd.Aliases,
	Short:   worktreeAddCmd.Short,
	Long:    worktreeAddCmd.Long,
	Run:     worktreeAddCmd.Run,
}

func init() {
	rootCmd.AddCommand(addCmd)

	worktreeAddCmd.Flags().VisitAll(func(flag *pflag.Flag) {
		addCmd.Flags().AddFlag(flag)
	})
	worktreeAddCmd.PersistentFlags().VisitAll(func(flag *pflag.Flag) {
		addCmd.PersistentFlags().AddFlag(flag)
	})

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
