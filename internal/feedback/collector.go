package feedback

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/AdeshDeshmukh/cortex/pkg/db"
	"github.com/AdeshDeshmukh/cortex/pkg/types"
)

type Collector struct {
	db     *db.DB
	reader *bufio.Reader
}

func NewCollector(database *db.DB) *Collector {
	return &Collector{
		db:     database,
		reader: bufio.NewReader(os.Stdin),
	}
}

type Stats struct {
	Accepted int
	Rejected int
	Skipped  int
	Total    int
}

func (c *Collector) CollectFeedback(reviewID int64, suggestions []types.Suggestion) (*Stats, error) {
	stats := &Stats{Total: len(suggestions)}

	if len(suggestions) == 0 {
		return stats, nil
	}

	fmt.Printf("\n🤖 Found %d suggestion(s)\n\n", len(suggestions))

	for i, sugg := range suggestions {
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		fmt.Printf("[%d/%d] %s (%s)\n", i+1, len(suggestions), sugg.Type, sugg.Severity)
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")

		fmt.Printf("📍 %s:%d\n", sugg.FilePath, sugg.LineNumber)
		fmt.Printf("💬 %s\n", sugg.Message)
		fmt.Printf("💡 %s\n", sugg.Suggestion)
		fmt.Printf("📊 Confidence: %.0f%%\n\n", sugg.Confidence*100)

		suggID, err := c.db.SaveSuggestion(
			reviewID,
			sugg.Type,
			sugg.Severity,
			sugg.Message,
			sugg.FilePath,
			sugg.LineNumber,
			sugg.Suggestion,
			sugg.Confidence,
			sugg.Source,
		)
		if err != nil {
			return stats, fmt.Errorf("save suggestion: %w", err)
		}

		action, reason := c.promptUser()

		if err := c.db.SaveFeedback(suggID, action, reason); err != nil {
			return stats, fmt.Errorf("save feedback: %w", err)
		}

		switch action {
		case "accept":
			stats.Accepted++
			fmt.Println("✅ Accepted")
		case "reject":
			stats.Rejected++
			fmt.Println("❌ Rejected")
		case "skip":
			stats.Skipped++
			fmt.Println("⏭️  Skipped")
		}
	}

	c.printSummary(stats)

	return stats, nil
}

func (c *Collector) promptUser() (action, reason string) {
	fmt.Print("✅ Accept  ❌ Reject  ⏭️  Skip  💬 Feedback\n")
	fmt.Print("> ")

	input, _ := c.reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))

	switch input {
	case "a", "accept":
		return "accept", ""
	case "r", "reject":
		fmt.Print("❌ Why reject? (optional, press Enter to skip): ")
		reasonInput, _ := c.reader.ReadString('\n')
		return "reject", strings.TrimSpace(reasonInput)
	case "s", "skip", "":
		return "skip", ""
	default:
		fmt.Println("⚠️  Invalid input. Skipping...")
		return "skip", ""
	}
}

func (c *Collector) printSummary(stats *Stats) {
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("📊 Feedback Summary\n")
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")
	fmt.Printf("✅ Accepted: %d\n", stats.Accepted)
	fmt.Printf("❌ Rejected: %d\n", stats.Rejected)
	fmt.Printf("⏭️  Skipped:  %d\n", stats.Skipped)
	fmt.Printf("📈 Total:    %d\n\n", stats.Total)

	if stats.Total > 0 {
		acceptRate := float64(stats.Accepted) / float64(stats.Total) * 100
		fmt.Printf("💯 Acceptance Rate: %.1f%%\n\n", acceptRate)
	}
}