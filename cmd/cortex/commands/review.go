package commands

import (
	"fmt"

	"github.com/AdeshDeshmukh/cortex/internal/git"
	"github.com/AdeshDeshmukh/cortex/pkg/types"
	"github.com/spf13/cobra"
)

var reviewCmd = &cobra.Command{
	Use:   "review",
	Short: "Review staged changes",
	Long: `Analyze staged git changes and provide AI-powered suggestions.

The review process:
  1. Parse git diff
  2. Run static analysis
  3. Generate AI suggestions
  4. Collect feedback for learning

Examples:
  cortex review
  cortex review --hook`,
	RunE: runReview,
}

var hookMode bool

func init() {
	rootCmd.AddCommand(reviewCmd)

	reviewCmd.Flags().BoolVar(&hookMode, "hook", false, "run in hook mode")
}

func runReview(cmd *cobra.Command, args []string) error {
	diff, err := git.GetStagedDiff()
	if err != nil {
		return fmt.Errorf("get staged changes: %w", err)
	}

	if diff == "" {
		fmt.Println("ℹ️  No staged changes to review")
		fmt.Println()
		fmt.Println("Hint: Stage files with 'git add <file>'")
		return nil
	}

	parser := git.NewDiffParser(diff)
	changes, err := parser.Parse()
	if err != nil {
		return fmt.Errorf("parse diff: %w", err)
	}

	if len(changes) == 0 {
		fmt.Println("ℹ️  No code changes detected")
		return nil
	}

	displayReviewSummary(changes)
	displayFileBreakdown(changes)

	fmt.Println()
	fmt.Println("⏳ Analysis pipeline:")
	fmt.Println("   ✅ Diff parsing complete")
	fmt.Println("   ⏹️  Static analysis (coming soon)")
	fmt.Println("   ⏹️  LLM suggestions (coming soon)")
	fmt.Println("   ⏹️  Interactive feedback (coming soon)")

	return nil
}

func displayReviewSummary(changes []types.DiffChange) {
	fmt.Println("🧠 Cortex Review")
	fmt.Println()

	totalAdded := 0
	totalRemoved := 0
	languages := make(map[string]int)

	for _, change := range changes {
		totalAdded += len(change.AddedLines)
		totalRemoved += len(change.RemovedLines)
		if change.FileType != "unknown" {
			languages[change.FileType]++
		}
	}

	fmt.Println("📊 Summary:")
	fmt.Printf("   Files changed:  %d\n", len(changes))
	fmt.Printf("   Lines added:    +%d\n", totalAdded)
	fmt.Printf("   Lines removed:  -%d\n", totalRemoved)

	if len(languages) > 0 {
		fmt.Print("   Languages:      ")
		first := true
		for lang, count := range languages {
			if !first {
				fmt.Print(", ")
			}
			fmt.Printf("%s (%d)", lang, count)
			first = false
		}
		fmt.Println()
	}
	fmt.Println()
}

func displayFileBreakdown(changes []types.DiffChange) {
	fmt.Println("📁 Files:")
	for i, change := range changes {
		fmt.Printf("   %d. %s (%s)\n", i+1, change.FilePath, change.FileType)
		fmt.Printf("      +%d -%d lines", len(change.AddedLines), len(change.RemovedLines))
		if change.StartLine > 0 {
			fmt.Printf(" starting at line %d", change.StartLine)
		}
		fmt.Println()
	}
}