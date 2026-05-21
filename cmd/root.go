package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "wt",
	Short: "Git worktree manager for iTerm2",
	Long:  "wt manages git worktrees and opens them in a 2x2 iTerm2 grid.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(newAliasCmd)
	rootCmd.AddCommand(rmCmd)
	rootCmd.AddCommand(initReposCmd)
	rootCmd.AddCommand(reviewCmd)
	rootCmd.AddCommand(versionCmd)
}
