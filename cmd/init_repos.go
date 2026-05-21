package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/backevik/wt-cli/internal/config"
	"github.com/backevik/wt-cli/internal/ui"
	"github.com/spf13/cobra"
)

var initReposWorkDir string

var initReposCmd = &cobra.Command{
	Use:   "init-repos",
	Short: "Scan work directory and save repo aliases",
	Long:  "Scans the work directory for git repositories and saves them to ~/.config/worktree/repos.json.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if initReposWorkDir == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("could not determine home directory: %w", err)
			}
			initReposWorkDir = filepath.Join(home, "projects", "work")
		}

		ui.Info("Scanning %s for git repositories...", initReposWorkDir)

		entries, err := os.ReadDir(initReposWorkDir)
		if err != nil {
			return fmt.Errorf("could not read work directory %s: %w", initReposWorkDir, err)
		}

		var repos []config.Repo
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			dirPath := filepath.Join(initReposWorkDir, entry.Name())
			gitPath := filepath.Join(dirPath, ".git")
			if _, err := os.Stat(gitPath); err != nil {
				continue
			}
			realPath, err := filepath.EvalSymlinks(dirPath)
			if err != nil {
				realPath = dirPath
			}
			repos = append(repos, config.Repo{
				Alias: entry.Name(),
				Path:  realPath,
			})
		}

		worktreesDir := filepath.Join(initReposWorkDir, "worktrees")

		cfg := &config.Config{
			WorkDir:      initReposWorkDir,
			WorktreesDir: worktreesDir,
			Repos:        repos,
		}
		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("could not save config: %w", err)
		}

		ui.Success("Found %d repositories", len(repos))
		for _, r := range repos {
			fmt.Printf("  %s → %s\n", r.Alias, r.Path)
		}
		return nil
	},
}

func init() {
	initReposCmd.Flags().StringVar(&initReposWorkDir, "work-dir", "", "Directory to scan (default: ~/projects/work)")
}
