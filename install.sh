#!/usr/bin/env bash
# Convenience wrapper — runs 'mise run setup' for first-time installs
# without requiring mise to already be active in the shell.
set -euo pipefail

if command -v mise &>/dev/null; then
  mise run setup
else
  echo "mise not found — install it from https://mise.jdx.dev and re-run this script."
  exit 1
fi

# Ensure ~/.local/bin is in PATH
if [[ ":$PATH:" != *":$HOME/.local/bin:"* ]]; then
  echo ""
  echo "Add the following to your ~/.zshrc:"
  echo "  export PATH=\"\$HOME/.local/bin:\$PATH\""
fi

echo ""
echo "Migration steps (if replacing the old Python/Zsh system):"
echo "  1. Remove these lines from ~/.zshrc:"
echo "       source \".../worktree-automation/shell/worktree-functions.zsh\""
echo "       source \".../worktree-automation/shell/wt-review.zsh\""
echo ""
echo "  2. Update Raycast new-worktree.sh:"
echo "       Replace: wt-new-alias \"\$1\" \"\$2\""
echo "       With:    ~/.local/bin/wt new-alias \"\$1\" \"\$2\""
echo ""
echo "  3. Optional aliases for muscle memory (add to ~/.zshrc):"
echo "       alias wt-new='wt new'"
echo "       alias wt-rm='wt rm'"
echo "       alias wt-init-repos='wt init-repos'"
echo "       alias wt-new-alias='wt new-alias'"
echo "       alias wt-review='wt review'"
