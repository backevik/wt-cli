package cmd

import (
	"fmt"

	"github.com/backevik/wt-cli/internal/config"
	"github.com/backevik/wt-cli/internal/ui"
	"github.com/spf13/cobra"
)

var newAliasCmd = &cobra.Command{
	Use:   "new-alias <alias> [branch]",
	Short: "Create a worktree in an aliased repository",
	Long: `Creates a git worktree in the repository identified by alias, as stored in
~/.config/worktree/repos.json. Useful when you are not currently in the target repo.`,
	Args: cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		alias := args[0]
		var branchArgs []string
		if len(args) > 1 {
			branchArgs = args[1:]
		}

		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("could not load config: %w", err)
		}
		if len(cfg.Repos) == 0 {
			return fmt.Errorf("no repos configured; run 'wt init-repos' first")
		}

		var repoPath string
		for _, r := range cfg.Repos {
			if r.Alias == alias {
				repoPath = r.Path
				break
			}
		}
		if repoPath == "" {
			ui.Error("unknown alias %q", alias)
			fmt.Println("Available aliases:")
			for _, r := range cfg.Repos {
				fmt.Printf("  %s → %s\n", r.Alias, r.Path)
			}
			return fmt.Errorf("alias not found: %s", alias)
		}

		return createWorktree(repoPath, branchArgs)
	},
}
