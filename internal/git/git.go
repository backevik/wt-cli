package git

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

// run executes a git command in the given directory and returns trimmed stdout.
func run(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("git %s: %s", strings.Join(args, " "), strings.TrimSpace(string(ee.Stderr)))
		}
		return "", fmt.Errorf("git %s: %w", strings.Join(args, " "), err)
	}
	return strings.TrimSpace(string(out)), nil
}

// IsRepo returns true if path is inside a git repository.
func IsRepo(path string) bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	cmd.Dir = path
	return cmd.Run() == nil
}

// IsWorktree returns true if path is a linked worktree (not the main repo).
func IsWorktree(path string) (bool, error) {
	gitDir, err := run(path, "rev-parse", "--git-dir")
	if err != nil {
		return false, err
	}
	commonDir, err := run(path, "rev-parse", "--git-common-dir")
	if err != nil {
		return false, err
	}
	// Make absolute for reliable comparison
	if !filepath.IsAbs(gitDir) {
		gitDir = filepath.Join(path, gitDir)
	}
	if !filepath.IsAbs(commonDir) {
		commonDir = filepath.Join(path, commonDir)
	}
	return filepath.Clean(gitDir) != filepath.Clean(commonDir), nil
}

// RepoRoot returns the top-level directory of the repository at path.
func RepoRoot(path string) (string, error) {
	return run(path, "rev-parse", "--show-toplevel")
}

// CommonDir returns the common git directory (main .git dir even from a worktree).
// The returned path is always absolute.
func CommonDir(path string) (string, error) {
	dir, err := run(path, "rev-parse", "--git-common-dir")
	if err != nil {
		return "", err
	}
	if !filepath.IsAbs(dir) {
		dir = filepath.Join(path, dir)
	}
	return filepath.Clean(dir), nil
}

// MainBranch returns the main branch name by following origin/HEAD.
// Falls back to "main" if it can't be determined.
func MainBranch(repoPath string) (string, error) {
	out, err := run(repoPath, "symbolic-ref", "refs/remotes/origin/HEAD")
	if err != nil {
		return "main", nil
	}
	// out is like "refs/remotes/origin/main"
	parts := strings.Split(out, "/")
	if len(parts) >= 1 {
		return parts[len(parts)-1], nil
	}
	return "main", nil
}

// Fetch fetches the given branch from origin.
func Fetch(repoPath, branch string) error {
	cmd := exec.Command("git", "fetch", "origin", branch)
	cmd.Dir = repoPath
	cmd.Stdout = nil
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git fetch origin %s: %s", branch, strings.TrimSpace(string(out)))
	}
	return nil
}

// WorktreeAdd creates a new worktree at wtPath with branchName starting from startPoint.
func WorktreeAdd(repoPath, wtPath, branchName, startPoint string) error {
	cmd := exec.Command("git", "worktree", "add", wtPath, "-b", branchName, startPoint)
	cmd.Dir = repoPath
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git worktree add: %s", strings.TrimSpace(string(out)))
	}
	return nil
}

// WorktreeRemove removes the worktree at wtPath from repoPath.
func WorktreeRemove(repoPath, wtPath string, force bool) error {
	args := []string{"worktree", "remove"}
	if force {
		args = append(args, "--force")
	}
	args = append(args, wtPath)
	cmd := exec.Command("git", args...)
	cmd.Dir = repoPath
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git worktree remove: %s", strings.TrimSpace(string(out)))
	}
	return nil
}

// SetupScript reads the worktree.setup git config value, returning empty string if unset.
func SetupScript(repoPath string) (string, error) {
	out, err := run(repoPath, "config", "worktree.setup")
	if err != nil {
		// config not set is not an error
		return "", nil
	}
	return out, nil
}

var nonAlphanumHyphen = regexp.MustCompile(`[^a-z0-9-]`)

// GitHubUsername returns the github.user config value, falling back to a sanitized user.name.
func GitHubUsername(repoPath string) (string, error) {
	username, err := run(repoPath, "config", "github.user")
	if err == nil && username != "" {
		return username, nil
	}
	name, err := run(repoPath, "config", "user.name")
	if err != nil || name == "" {
		return "dev", nil
	}
	sanitized := strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' {
			return unicode.ToLower(r)
		}
		return -1
	}, name)
	sanitized = nonAlphanumHyphen.ReplaceAllString(strings.ToLower(name), "")
	if sanitized == "" {
		return "dev", nil
	}
	return sanitized, nil
}
