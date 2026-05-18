package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func setupTestRepo(t *testing.T) string {
	t.Helper()

	tmpDir := t.TempDir()

	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("git init: %v", err)
	}

	cmd = exec.Command("git", "config", "user.email", "test@cortex.dev")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("git config email: %v", err)
	}

	cmd = exec.Command("git", "config", "user.name", "Cortex Test")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("git config name: %v", err)
	}

	return tmpDir
}

func TestHookManager_Install(t *testing.T) {
	repoPath := setupTestRepo(t)

	hm, err := NewHookManager(repoPath)
	if err != nil {
		t.Fatalf("NewHookManager: %v", err)
	}

	if err := hm.Install(false); err != nil {
		t.Fatalf("Install: %v", err)
	}

	hookPath := filepath.Join(repoPath, ".git", "hooks", "pre-commit")
	info, err := os.Stat(hookPath)
	if err != nil {
		t.Fatalf("hook file not created: %v", err)
	}

	if info.Mode().Perm()&0111 == 0 {
		t.Error("hook is not executable")
	}

	if !hm.IsInstalled() {
		t.Error("IsInstalled returned false after installation")
	}
}

func TestHookManager_InstallTwice(t *testing.T) {
	repoPath := setupTestRepo(t)

	hm, err := NewHookManager(repoPath)
	if err != nil {
		t.Fatalf("NewHookManager: %v", err)
	}

	if err := hm.Install(false); err != nil {
		t.Fatalf("first install: %v", err)
	}

	if err := hm.Install(false); err == nil {
		t.Error("second install should fail without --force")
	}

	if err := hm.Install(true); err != nil {
		t.Errorf("install with force: %v", err)
	}
}

func TestHookManager_Uninstall(t *testing.T) {
	repoPath := setupTestRepo(t)

	hm, err := NewHookManager(repoPath)
	if err != nil {
		t.Fatalf("NewHookManager: %v", err)
	}

	if err := hm.Install(false); err != nil {
		t.Fatalf("Install: %v", err)
	}

	if err := hm.Uninstall(); err != nil {
		t.Fatalf("Uninstall: %v", err)
	}

	if hm.IsInstalled() {
		t.Error("hook still installed after uninstall")
	}
}

func TestHookManager_NotGitRepo(t *testing.T) {
	tmpDir := t.TempDir()

	_, err := NewHookManager(tmpDir)
	if err == nil {
		t.Error("expected error for non-git directory")
	}
}