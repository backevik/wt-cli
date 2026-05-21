package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/backevik/wt-cli/internal/config"
	"github.com/backevik/wt-cli/internal/git"
	"github.com/backevik/wt-cli/internal/iterm"
	"github.com/backevik/wt-cli/internal/namegen"
	"github.com/backevik/wt-cli/internal/ui"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new [branch]",
	Short: "Create a worktree in the current repository",
	Long: `Creates a git worktree from the latest origin main branch and opens it
in a 2x2 iTerm2 grid. If no branch name is given, a random one is generated.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("could not get working directory: %w", err)
		}
		return createWorktree(cwd, args)
	},
}

// createWorktree contains the shared logic for `new` and `new-alias`.
func createWorktree(repoPath string, args []string) error {
	if !git.IsRepo(repoPath) {
		return fmt.Errorf("not a git repository: %s", repoPath)
	}

	username, err := git.GitHubUsername(repoPath)
	if err != nil {
		username = "dev"
	}

	userProvided := len(args) > 0
	branchName := ""
	if userProvided {
		branchName = args[0]
	}

	mainBranch, err := git.MainBranch(repoPath)
	if err != nil {
		return fmt.Errorf("could not determine main branch: %w", err)
	}

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("could not load config: %w", err)
	}
	if cfg.WorktreesDir == "" {
		return fmt.Errorf("worktrees_dir not configured — run 'wt init-repos' first")
	}
	worktreesBase := cfg.WorktreesDir

	// Attempt to create the worktree, retrying on path collision (auto-generated names only).
	const maxAttempts = 3
	var wtPath string
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if !userProvided || branchName == "" {
			branchName = namegen.RandomBranchName(username)
		}

		dirName := branchName
		if idx := strings.Index(branchName, "/"); idx >= 0 {
			dirName = branchName[idx+1:]
		}
		wtPath = filepath.Join(worktreesBase, dirName)

		if _, err := os.Stat(wtPath); os.IsNotExist(err) {
			break // path is free
		}
		if userProvided {
			return fmt.Errorf("worktree path already exists: %s", wtPath)
		}
		if attempt == maxAttempts {
			return fmt.Errorf("could not find a free worktree path after %d attempts", maxAttempts)
		}
		ui.Info("Path %s already exists, trying a different name...", wtPath)
	}

	ui.Info("Fetching latest from origin/%s...", mainBranch)
	if err := git.Fetch(repoPath, mainBranch); err != nil {
		return fmt.Errorf("fetch failed: %w", err)
	}

	ui.Info("Creating worktree: %s → %s", branchName, wtPath)
	if err := git.WorktreeAdd(repoPath, wtPath, branchName, "origin/"+mainBranch); err != nil {
		return fmt.Errorf("worktree add failed: %w", err)
	}

	setupScript, err := git.SetupScript(repoPath)
	if err != nil {
		return fmt.Errorf("could not read setup script: %w", err)
	}

	if err := iterm.Open(wtPath, setupScript); err != nil {
		return fmt.Errorf("could not open iTerm2: %w", err)
	}

	return nil
}

