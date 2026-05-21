package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/backevik/wt-cli/internal/git"
	"github.com/backevik/wt-cli/internal/ui"
	"github.com/spf13/cobra"
)

var rmForce bool

var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Delete the current worktree",
	Long: `Removes the git worktree for the current directory. Must be run from inside
a linked worktree (not the main repository). Prompts for confirmation unless --force is given.

Note: after removal your shell's working directory will no longer exist.
Run 'cd ~' or open a new tab.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("could not get working directory: %w", err)
		}

		if !git.IsRepo(cwd) {
			return fmt.Errorf("not a git repository")
		}

		isWt, err := git.IsWorktree(cwd)
		if err != nil {
			return fmt.Errorf("could not determine worktree status: %w", err)
		}
		if !isWt {
			return fmt.Errorf("cannot remove: this is the main repository, not a linked worktree")
		}

		wtPath, err := git.RepoRoot(cwd)
		if err != nil {
			return fmt.Errorf("could not determine worktree root: %w", err)
		}

		commonDir, err := git.CommonDir(cwd)
		if err != nil {
			return fmt.Errorf("could not determine git common dir: %w", err)
		}
		mainRepoPath := filepath.Dir(commonDir)

		if !rmForce {
			ok, err := ui.Confirm(fmt.Sprintf("Delete worktree %s?", wtPath))
			if err != nil {
				return err
			}
			if !ok {
				ui.Step("aborted")
				return nil
			}
		}

		if err := git.WorktreeRemove(mainRepoPath, wtPath, true); err != nil {
			return fmt.Errorf("failed to remove worktree: %w", err)
		}

		ui.Success("Worktree nuked: %s", wtPath)
		ui.Step("cwd is gone — run 'cd ~' or open a new tab")
		return nil
	},
}

func init() {
	rmCmd.Flags().BoolVar(&rmForce, "force", false, "Skip confirmation prompt")
}
