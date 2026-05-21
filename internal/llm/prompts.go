package llm

import (
	"fmt"
	"strings"

	"github.com/AdeshDeshmukh/cortex/pkg/types"
)

type PromptBuilder struct {
	maxLines int
}

type ParsedResponse struct {
	Suggestions []types.Suggestion
	Raw         string
}

func NewPromptBuilder() *PromptBuilder {
	return &PromptBuilder{
		maxLines: 50,
	}
}

func (p *PromptBuilder) Build(change types.DiffChange) string {
	var sb strings.Builder

	sb.WriteString(p.systemInstructions(change.FileType))
	sb.WriteString("\n\n")
	sb.WriteString(p.codeContext(change))
	sb.WriteString("\n\n")
	sb.WriteString(p.reviewInstructions())

	return sb.String()
}

func (p *PromptBuilder) systemInstructions(fileType string) string {
	base := `You are an expert code reviewer with deep knowledge of software engineering best practices.
Your goal is to identify real issues, not nitpick style preferences.
Be specific, actionable, and concise.`

	languageGuide := p.languageSpecificGuide(fileType)
	if languageGuide != "" {
		return base + "\n\n" + languageGuide
	}

	return base
}

func (p *PromptBuilder) languageSpecificGuide(fileType string) string {
	guides := map[string]string{
		"go": `Go-specific focus areas:
- Error handling: Every error must be checked
- Goroutine leaks: Goroutines must have clear exit conditions
- Context propagation: Context should flow through call chain
- Interface design: Keep interfaces small and focused`,

		"python": `Python-specific focus areas:
- Exception handling: Use specific exception types
- Type hints: Ensure type annotations are present
- Resource management: Use context managers for resources
- Mutability: Be careful with mutable default arguments`,

		"javascript": `JavaScript-specific focus areas:
- Async/await: Ensure promises are properly handled
- Null checks: Guard against null/undefined access
- Memory leaks: Event listeners must be cleaned up
- Type coercion: Avoid implicit type conversions`,

		"typescript": `TypeScript-specific focus areas:
- Type safety: Avoid 'any' type where possible
- Null safety: Handle undefined and null explicitly
- Generic constraints: Use proper type constraints
- Interface vs Type: Use interfaces for object shapes`,
	}

	if guide, ok := guides[fileType]; ok {
		return guide
	}

	return ""
}

func (p *PromptBuilder) codeContext(change types.DiffChange) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("File: %s (Language: %s)\n", change.FilePath, change.FileType))
	sb.WriteString(fmt.Sprintf("Change type: %s\n", change.ChangeType))
	sb.WriteString("\n")

	if change.Context != "" {
		sb.WriteString("Surrounding context:\n")
		sb.WriteString("```\n")
		sb.WriteString(change.Context)
		sb.WriteString("```\n\n")
	}

	if len(change.RemovedLines) > 0 {
		sb.WriteString("Removed lines:\n")
		sb.WriteString("```\n")
		for _, line := range change.RemovedLines {
			sb.WriteString(fmt.Sprintf("- %s\n", line))
		}
		sb.WriteString("```\n\n")
	}

	if len(change.AddedLines) > 0 {
		sb.WriteString("Added lines:\n")
		sb.WriteString("```\n")

		limit := len(change.AddedLines)
		if limit > p.maxLines {
			limit = p.maxLines
		}

		for _, line := range change.AddedLines[:limit] {
			sb.WriteString(fmt.Sprintf("+ %s\n", line))
		}

		if len(change.AddedLines) > p.maxLines {
			sb.WriteString(fmt.Sprintf("... (%d more lines)\n", len(change.AddedLines)-p.maxLines))
		}

		sb.WriteString("```\n")
	}

	return sb.String()
}

func (p *PromptBuilder) reviewInstructions() string {
	return `Review the code changes above and identify issues.

For each issue found, respond in this exact format:
ISSUE: <one line description of the problem>
SEVERITY: <critical|high|medium|low>
TYPE: <security|error-handling|performance|reliability|type-safety|documentation>
SUGGESTION: <specific actionable fix>
---

Rules:
- Only flag real issues, not style preferences
- Maximum 5 issues per file
- If no issues found, respond with: NO_ISSUES
- Be specific about which line or pattern is problematic`
}

type ResponseParser struct{}

func NewResponseParser() *ResponseParser {
	return &ResponseParser{}
}

func (r *ResponseParser) Parse(raw string, filePath string) []types.Suggestion {
	if strings.Contains(raw, "NO_ISSUES") {
		return nil
	}

	var suggestions []types.Suggestion

	blocks := strings.Split(raw, "---")

	for i, block := range blocks {
		block = strings.TrimSpace(block)
		if block == "" {
			continue
		}

		suggestion := r.parseBlock(block, filePath, i)
		if suggestion != nil {
			suggestions = append(suggestions, *suggestion)
		}
	}

	return suggestions
}

func (r *ResponseParser) parseBlock(block, filePath string, index int) *types.Suggestion {
	lines := strings.Split(block, "\n")

	suggestion := &types.Suggestion{
		ID:         fmt.Sprintf("llm_%d", index),
		FilePath:   filePath,
		Confidence: 0.7,
		Source:     "llm",
	}

	for _, line := range lines {
		line = strings.TrimSpace(line)

		switch {
		case strings.HasPrefix(line, "ISSUE:"):
			suggestion.Message = strings.TrimSpace(strings.TrimPrefix(line, "ISSUE:"))

		case strings.HasPrefix(line, "SEVERITY:"):
			severity := strings.TrimSpace(strings.TrimPrefix(line, "SEVERITY:"))
			suggestion.Severity = r.normalizeSeverity(severity)

		case strings.HasPrefix(line, "TYPE:"):
			suggestion.Type = strings.TrimSpace(strings.TrimPrefix(line, "TYPE:"))

		case strings.HasPrefix(line, "SUGGESTION:"):
			if suggestion.Message != "" {
				suggestion.Message += " - " + strings.TrimSpace(strings.TrimPrefix(line, "SUGGESTION:"))
			}
		}
	}

	if suggestion.Message == "" || suggestion.Severity == "" {
		return nil
	}

	return suggestion
}

func (r *ResponseParser) normalizeSeverity(severity string) string {
	severity = strings.ToLower(strings.TrimSpace(severity))

	validSeverities := map[string]bool{
		"critical": true,
		"high":     true,
		"medium":   true,
		"low":      true,
	}

	if validSeverities[severity] {
		return severity
	}

	return "medium"
}