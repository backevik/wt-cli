package iterm

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// HelperPath resolves the path to the iterm-worktree-open Python script.
// Search order:
//  1. $WT_ITERM_HELPER env var
//  2. ~/.local/bin/iterm-worktree-open
//  3. {executable-dir}/scripts/iterm-worktree-open
func HelperPath() (string, error) {
	if v := os.Getenv("WT_ITERM_HELPER"); v != "" {
		return v, nil
	}

	home, _ := os.UserHomeDir()
	candidate := filepath.Join(home, ".local", "bin", "iterm-worktree-open")
	if _, err := os.Stat(candidate); err == nil {
		return candidate, nil
	}

	exe, err := os.Executable()
	if err == nil {
		candidate = filepath.Join(filepath.Dir(exe), "scripts", "iterm-worktree-open")
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
	}

	return "", fmt.Errorf("iterm-worktree-open not found; run 'make install' or set WT_ITERM_HELPER")
}

// Open launches the iTerm2 helper script for the given worktree path.
// setupScript may be empty.
func Open(wtPath, setupScript string) error {
	helperPath, err := HelperPath()
	if err != nil {
		return err
	}

	args := []string{helperPath, wtPath}
	if setupScript != "" {
		args = append(args, setupScript)
	}

	cmd := exec.Command("python3", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
