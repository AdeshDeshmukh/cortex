package commands

import (
	"fmt"

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
	fmt.Println("🧠 Cortex Review")
	fmt.Println()
	fmt.Println("Status: Placeholder implementation")
	fmt.Println()
	fmt.Println("Upcoming features:")
	fmt.Println("  → Git diff parsing")
	fmt.Println("  → Static analysis")
	fmt.Println("  → LLM integration")
	fmt.Println("  → Interactive feedback")
	fmt.Println("  → Reinforcement learning")

	return nil
}
