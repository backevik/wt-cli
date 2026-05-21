#!/usr/bin/env python3
# worktree-cleanup-monitor.py - Monitor iTerm2 tab and cleanup worktree when tab closes

import iterm2
import sys
import os
import json
import shutil
import subprocess
import asyncio

def _get_worktrees_dir():
    config_path = os.path.expanduser("~/.config/worktree/repos.json")
    try:
        with open(config_path) as f:
            cfg = json.load(f)
        return cfg.get("worktrees_dir", os.path.expanduser("~/projects/work/worktrees"))
    except (FileNotFoundError, json.JSONDecodeError):
        return os.path.expanduser("~/projects/work/worktrees")


async def main(connection):
    if len(sys.argv) < 2:
        print("Usage: worktree-cleanup-monitor.py <worktree-path>")
        sys.exit(1)

    worktree_path = sys.argv[1]

    app = await iterm2.async_get_app(connection)

    # Find the tab that contains a session with our worktree path
    target_tab_id = None
    for window in app.terminal_windows:
        for tab in window.tabs:
            for session in tab.sessions:
                try:
                    stored_path = await session.async_get_variable("user.WORKTREE_PATH")
                    if stored_path == worktree_path:
                        target_tab_id = tab.tab_id
                        break
                except:
                    continue
            if target_tab_id:
                break
        if target_tab_id:
            break

    if not target_tab_id:
        print(f"Warning: Could not find tab for {worktree_path}")
        sys.exit(1)

    # Wait for the entire tab to close by polling
    while True:
        await asyncio.sleep(2)

        tab_still_alive = False
        for window in app.terminal_windows:
            for tab in window.tabs:
                if tab.tab_id == target_tab_id:
                    tab_still_alive = True
                    break
            if tab_still_alive:
                break

        if not tab_still_alive:
            break

    # Verify worktree still exists
    if not os.path.exists(worktree_path):
        sys.exit(0)

    # Find git dir and remove worktree
    try:
        git_dir_result = subprocess.run(
            ['git', 'rev-parse', '--git-common-dir'],
            cwd=worktree_path,
            capture_output=True,
            text=True
        )

        if git_dir_result.returncode == 0:
            git_common_dir = git_dir_result.stdout.strip()
            repo_dir = os.path.dirname(git_common_dir)

            subprocess.run(
                ['git', 'worktree', 'remove', '--force', worktree_path],
                cwd=repo_dir
            )

        # Fallback: remove directory if it still exists.
        # Safety guard: only delete paths under the expected worktrees base dir.
        worktrees_base = _get_worktrees_dir()
        real_path = os.path.realpath(worktree_path)
        real_base = os.path.realpath(worktrees_base)
        if os.path.exists(worktree_path) and real_path.startswith(real_base + os.sep):
            shutil.rmtree(worktree_path)

        if not os.path.exists(worktree_path):
            # Show notification
            subprocess.run([
                'osascript', '-e',
                f'display notification "Worktree deleted: {os.path.basename(worktree_path)}" with title "Worktree Cleanup"'
            ])
    except Exception as e:
        print(f"Error during cleanup: {e}")

if __name__ == "__main__":
    iterm2.run_until_complete(main)
