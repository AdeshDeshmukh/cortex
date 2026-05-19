package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const preCommitHookTemplate = `#!/bin/bash
set -e

REPO_ROOT="$(git rev-parse --show-toplevel)"
CORTEX_BIN="$REPO_ROOT/cortex"

if [ ! -f "$CORTEX_BIN" ]; then
    echo "⚠️  Cortex binary not found at $CORTEX_BIN"
    echo "   Run 'make build' to create it"
    exit 1
fi

"$CORTEX_BIN" review --hook
EXIT_CODE=$?

case $EXIT_CODE in
    0) exit 0 ;;
    1) 
        echo "❌ Critical issues found. Commit aborted."
        echo "   Fix issues or use: git commit --no-verify"
        exit 1
        ;;
    2)
        echo "⚠️  Commit aborted by user"
        exit 1
        ;;
    *)
        exit $EXIT_CODE
        ;;
esac
`

type HookManager struct {
	repoPath string
}

func NewHookManager(repoPath string) (*HookManager, error) {
	if !isGitRepo(repoPath) {
		return nil, fmt.Errorf("not a git repository: %s", repoPath)
	}

	return &HookManager{repoPath: repoPath}, nil
}

func (h *HookManager) Install(force bool) error {
	hookPath := h.hookPath()

	if err := os.MkdirAll(filepath.Dir(hookPath), 0755); err != nil {
		return fmt.Errorf("create hooks directory: %w", err)
	}

	if exists, isOurs := h.hookExists(); exists && !force {
		if isOurs {
			return fmt.Errorf("hook already installed")
		}
		return fmt.Errorf("pre-commit hook exists (use --force to overwrite)")
	}

	if err := os.WriteFile(hookPath, []byte(preCommitHookTemplate), 0755); err != nil {
		return fmt.Errorf("write hook: %w", err)
	}

	return nil
}

func (h *HookManager) Uninstall() error {
	hookPath := h.hookPath()

	exists, isOurs := h.hookExists()
	if !exists {
		return fmt.Errorf("no hook installed")
	}

	if !isOurs {
		return fmt.Errorf("hook is not managed by cortex")
	}

	if err := os.Remove(hookPath); err != nil {
		return fmt.Errorf("remove hook: %w", err)
	}

	return nil
}

func (h *HookManager) IsInstalled() bool {
	_, isOurs := h.hookExists()
	return isOurs
}

func (h *HookManager) hookPath() string {
	return filepath.Join(h.repoPath, ".git", "hooks", "pre-commit")
}

func (h *HookManager) hookExists() (exists bool, isOurs bool) {
	content, err := os.ReadFile(h.hookPath())
	if err != nil {
		return false, false
	}

	return true, strings.TrimSpace(string(content)) == strings.TrimSpace(preCommitHookTemplate)
}

func isGitRepo(path string) bool {
	gitDir := filepath.Join(path, ".git")
	info, err := os.Stat(gitDir)
	return err == nil && info.IsDir()
}

func GetRepoRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("not a git repository")
	}

	return strings.TrimSpace(string(output)), nil
}

func GetCurrentBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("get branch: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

func GetStagedDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--cached")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("get diff: %w", err)
	}

	return string(output), nil
}