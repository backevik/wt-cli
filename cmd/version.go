package cmd

import (
	"github.com/backevik/wt-cli/internal/ui"
	"github.com/spf13/cobra"
)

const version = "0.1.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version",
	Run: func(cmd *cobra.Command, args []string) {
		ui.Banner("wt " + version)
	},
}
