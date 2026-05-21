package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/backevik/wt-cli/internal/git"
	"github.com/backevik/wt-cli/internal/ui"
	"github.com/spf13/cobra"
)

var reviewCmd = &cobra.Command{
	Use:   "review [base-branch]",
	Short: "Interactive diff review with fzf and delta",
	Long: `Opens an interactive fuzzy-finder to review all changes in the current
worktree against the base branch. Requires fzf and delta to be in PATH.

Keybindings:
  enter    open file in $EDITOR
  ctrl-d   preview diff for selected file
  ctrl-a   preview all changes combined
  ctrl-s   stage selected file
  ctrl-u   unstage selected file`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("could not get working directory: %w", err)
		}

		if !git.IsRepo(cwd) {
			return fmt.Errorf("not a git repository")
		}

		// Check dependencies
		for _, dep := range []string{"fzf", "delta"} {
			if _, err := exec.LookPath(dep); err != nil {
				return fmt.Errorf("%s not found in PATH; please install it first", dep)
			}
		}

		var base string
		if len(args) > 0 {
			base = args[0]
		} else {
			base, err = git.MainBranch(cwd)
			if err != nil {
				base = "main"
			}
		}

		// Compute merge base
		mergeBaseCmd := exec.Command("git", "merge-base", "origin/"+base, "HEAD")
		mergeBaseCmd.Dir = cwd
		mergeBaseOut, err := mergeBaseCmd.Output()
		if err != nil {
			return fmt.Errorf("could not find merge base with origin/%s; make sure you have fetched", base)
		}
		mergeBase := string(mergeBaseOut[:len(mergeBaseOut)-1]) // trim newline

		// Check there are any changes
		changedCmd := exec.Command("sh", "-c",
			fmt.Sprintf(`{ git diff --name-only %s HEAD; git diff --name-only HEAD; git diff --name-only --cached HEAD; git ls-files --others --exclude-standard; } | sort -u`, mergeBase))
		changedCmd.Dir = cwd
		changedOut, _ := changedCmd.Output()
		if len(changedOut) == 0 {
			ui.Info("No changes against origin/%s", base)
			return nil
		}

		// Build fzf pipeline (mirrors wt-review.zsh exactly)
		reloadCmd := fmt.Sprintf(
			`{ git diff --name-only %s HEAD; git diff --name-only HEAD; git diff --name-only --cached HEAD; git ls-files --others --exclude-standard; } | sort -u`,
			mergeBase,
		)
		diffFileCmd := fmt.Sprintf(
			`{ git ls-files --others --exclude-standard | grep -qx {} && cat {} || { git diff %s HEAD -- {}; git diff HEAD -- {}; git diff --cached HEAD -- {}; }; } | delta --width $FZF_PREVIEW_COLUMNS`,
			mergeBase,
		)
		diffAllCmd := fmt.Sprintf(
			`{ git diff %s HEAD; git diff HEAD; git diff --cached HEAD; git ls-files --others --exclude-standard -z | xargs -0 -I%% sh -c 'printf "\n=== %% (untracked) ===\n"; cat %%'; } | delta --width $FZF_PREVIEW_COLUMNS`,
			mergeBase,
		)

		fileCount := len(changedOut)
		_ = fileCount

		pipeline := fmt.Sprintf(
			`{ git diff --name-only %s HEAD; git diff --name-only HEAD; git diff --name-only --cached HEAD; git ls-files --others --exclude-standard; } | sort -u | fzf `+
				`--ansi `+
				`--multi `+
				`--header "Files changed against origin/%s | enter:edit  ^d:diff  ^a:all  ^s:stage  ^u:unstage" `+
				`--preview '%s' `+
				`--preview-window "right:65%%:wrap" `+
				`--bind "enter:execute(%s < /dev/tty > /dev/tty)" `+
				`--bind "ctrl-d:preview(%s)" `+
				`--bind "ctrl-a:preview(%s)" `+
				`--bind "ctrl-s:execute-silent(git add {})+reload(%s)" `+
				`--bind "ctrl-u:execute-silent(git restore --staged {})+reload(%s)"`,
			mergeBase,
			base,
			diffFileCmd,
			`${EDITOR:-vim} {}`,
			diffFileCmd,
			diffAllCmd,
			reloadCmd,
			reloadCmd,
		)

		fzfCmd := exec.Command("zsh", "-c", pipeline)
		fzfCmd.Dir = cwd
		fzfCmd.Stdin = os.Stdin
		fzfCmd.Stdout = os.Stdout
		fzfCmd.Stderr = os.Stderr

		if err := fzfCmd.Run(); err != nil {
			if ee, ok := err.(*exec.ExitError); ok {
				// fzf exits 130 on ctrl-c — treat as normal exit
				if ee.ExitCode() == 130 {
					return nil
				}
			}
			return err
		}
		return nil
	},
}
