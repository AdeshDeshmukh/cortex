package analyzer

import (
	"fmt"
	"regexp"

	"github.com/AdeshDeshmukh/cortex/pkg/types"
)

type Analyzer struct {
	rules []Rule
}

type Rule struct {
	ID        string
	Name      string
	Pattern   *regexp.Regexp
	Severity  string
	Message   string
	Language  string
	CheckType string
}

func NewAnalyzer() *Analyzer {
	return &Analyzer{
		rules: getDefaultRules(),
	}
}

func (a *Analyzer) Analyze(changes []types.DiffChange) []types.Suggestion {
	var suggestions []types.Suggestion

	for _, change := range changes {
		currentLine := change.StartLine

		for _, addedLine := range change.AddedLines {
			lineSuggestions := a.checkLine(addedLine, change.FilePath, change.FileType, currentLine)
			suggestions = append(suggestions, lineSuggestions...)
			currentLine++
		}
	}

	return deduplicateSuggestions(suggestions)
}

func (a *Analyzer) checkLine(line, filePath, fileType string, lineNumber int) []types.Suggestion {
	var suggestions []types.Suggestion

	for _, rule := range a.rules {
		if rule.Language != "" && rule.Language != fileType {
			continue
		}

		if rule.Pattern.MatchString(line) {
			suggestions = append(suggestions, types.Suggestion{
				ID:         generateID(),
				Type:       rule.CheckType,
				Severity:   rule.Severity,
				Message:    rule.Message,
				FilePath:   filePath,
				LineNumber: lineNumber,
				Confidence: 0.8,
				Source:     "static",
			})
		}
	}

	return suggestions
}

func deduplicateSuggestions(suggestions []types.Suggestion) []types.Suggestion {
	seen := make(map[string]bool)
	unique := make([]types.Suggestion, 0, len(suggestions))

	for _, s := range suggestions {
		key := fmt.Sprintf("%s:%d:%s", s.FilePath, s.LineNumber, s.Type)

		if !seen[key] {
			seen[key] = true
			unique = append(unique, s)
		}
	}

	return unique
}

func getDefaultRules() []Rule {
	return []Rule{
		{
			ID:        "go-error-unchecked",
			Name:      "Unchecked error",
			Pattern:   regexp.MustCompile(`\w+\s*:=\s*[\w\.]+\([^)]*\)`),
			Severity:  "high",
			Message:   "Function call may return error but result is not checked",
			Language:  "go",
			CheckType: "error-handling",
		},
		{
			ID:        "todo-comment",
			Name:      "TODO comment",
			Pattern:   regexp.MustCompile(`(//|#)\s*(TODO|FIXME)`),
			Severity:  "medium",
			Message:   "TODO comment found - consider resolving before commit",
			Language:  "",
			CheckType: "documentation",
		},
		{
			ID:        "console-log",
			Name:      "Debug statement",
			Pattern:   regexp.MustCompile(`console\.log\(|fmt\.Println\(|print\(`),
			Severity:  "low",
			Message:   "Debug statement detected - remove before production",
			Language:  "",
			CheckType: "debugging",
		},
		{
			ID:        "hardcoded-password",
			Name:      "Hardcoded credential",
			Pattern:   regexp.MustCompile(`password\s*=\s*["'][^"']{3,}["']|api_key\s*=\s*["']`),
			Severity:  "critical",
			Message:   "Hardcoded credential detected - use environment variables",
			Language:  "",
			CheckType: "security",
		},
	}
}

func generateID() string {
	return "sugg_static_" + randomString(8)
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[i%len(letters)]
	}
	return string(b)
}