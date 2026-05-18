package commands

import (
	"fmt"

	"github.com/AdeshDeshmukh/cortex/internal/git"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install git hooks in current repository",
	Long: `Install the Cortex pre-commit hook to enable code review.

The hook will analyze code changes before each commit and provide
AI-powered suggestions that improve over time through learning.

Examples:
  cortex install
  cortex install --force
  cortex install --uninstall`,
	RunE: runInstall,
}

var (
	forceInstall bool
	uninstall    bool
)

func init() {
	rootCmd.AddCommand(installCmd)

	installCmd.Flags().BoolVarP(&forceInstall, "force", "f", false, "overwrite existing hook")
	installCmd.Flags().BoolVarP(&uninstall, "uninstall", "u", false, "remove hook")
}

func runInstall(cmd *cobra.Command, args []string) error {
	repoRoot, err := git.GetRepoRoot()
	if err != nil {
		return fmt.Errorf("must be in a git repository: %w", err)
	}

	hm, err := git.NewHookManager(repoRoot)
	if err != nil {
		return fmt.Errorf("hook manager: %w", err)
	}

	if uninstall {
		if err := hm.Uninstall(); err != nil {
			return fmt.Errorf("uninstall: %w", err)
		}
		fmt.Println("✅ Cortex hooks removed")
		return nil
	}

	if err := hm.Install(forceInstall); err != nil {
		return fmt.Errorf("install: %w", err)
	}

	branch, _ := git.GetCurrentBranch()

	fmt.Println("✅ Cortex installed successfully")
	fmt.Println()
	fmt.Printf("📍 Repository: %s\n", repoRoot)
	fmt.Printf("🌿 Branch:     %s\n", branch)
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  1. Make code changes")
	fmt.Println("  2. Run 'git commit'")
	fmt.Println("  3. Cortex will review automatically")
	fmt.Println()
	fmt.Println("Or test now: cortex review")

	return nil
}
