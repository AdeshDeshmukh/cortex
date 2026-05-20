package analyzer

import (
	"testing"

	"github.com/AdeshDeshmukh/cortex/pkg/types"
)

func TestAnalyzer_GoErrorDetection(t *testing.T) {
	analyzer := NewAnalyzer()
	
	changes := []types.DiffChange{
		{
			FilePath:   "main.go",
			FileType:   "go",
			AddedLines: []string{"x := doSomething()"},
		},
	}

	suggestions := analyzer.Analyze(changes)

	if len(suggestions) != 1 {
		t.Errorf("expected 1 suggestion, got %d", len(suggestions))
		return
	}

	suggestion := suggestions[0]

	if suggestion.Severity != "high" {
		t.Errorf("expected high severity, got %s", suggestion.Severity)
	}

	if suggestion.Type != "error-handling" {
		t.Errorf("expected error-handling type, got %s", suggestion.Type)
	}

	if suggestion.Source != "static" {
		t.Errorf("expected static source, got %s", suggestion.Source)
	}
}

func TestAnalyzer_TODODetection(t *testing.T) {
	analyzer := NewAnalyzer()
	
	changes := []types.DiffChange{
		{
			FilePath:   "main.go",
			FileType:   "go",
			AddedLines: []string{"// TODO: fix this later"},
		},
	}

	suggestions := analyzer.Analyze(changes)

	if len(suggestions) != 1 {
		t.Errorf("expected 1 suggestion, got %d", len(suggestions))
		return
	}

	suggestion := suggestions[0]

	if suggestion.Severity != "medium" {
		t.Errorf("expected medium severity, got %s", suggestion.Severity)
	}

	if suggestion.Type != "documentation" {
		t.Errorf("expected documentation type, got %s", suggestion.Type)
	}
}

func TestAnalyzer_DebugStatementDetection(t *testing.T) {
	analyzer := NewAnalyzer()
	
	changes := []types.DiffChange{
		{
			FilePath:   "main.js",
			FileType:   "javascript",
			AddedLines: []string{"console.log('debug info')"},
		},
	}

	suggestions := analyzer.Analyze(changes)

	if len(suggestions) != 1 {
		t.Errorf("expected 1 suggestion, got %d", len(suggestions))
		return
	}

	suggestion := suggestions[0]

	if suggestion.Severity != "low" {
		t.Errorf("expected low severity, got %s", suggestion.Severity)
	}

	if suggestion.Type != "debugging" {
		t.Errorf("expected debugging type, got %s", suggestion.Type)
	}
}

func TestAnalyzer_HardcodedPasswordDetection(t *testing.T) {
	analyzer := NewAnalyzer()
	
	changes := []types.DiffChange{
		{
			FilePath:   "config.py",
			FileType:   "python",
			AddedLines: []string{"password = 'admin123'"},
		},
	}

	suggestions := analyzer.Analyze(changes)

	if len(suggestions) != 1 {
		t.Errorf("expected 1 suggestion, got %d", len(suggestions))
		return
	}

	suggestion := suggestions[0]

	if suggestion.Severity != "critical" {
		t.Errorf("expected critical severity, got %s", suggestion.Severity)
	}

	if suggestion.Type != "security" {
		t.Errorf("expected security type, got %s", suggestion.Type)
	}
}

func TestAnalyzer_LanguageFiltering(t *testing.T) {
	analyzer := NewAnalyzer()
	
	changes := []types.DiffChange{
		{
			FilePath:   "script.py",
			FileType:   "python",
			AddedLines: []string{"x := doSomething()"},
		},
	}

	suggestions := analyzer.Analyze(changes)

	if len(suggestions) != 0 {
		t.Errorf("expected 0 suggestions for Go syntax in Python file, got %d", len(suggestions))
	}
}

func TestAnalyzer_MultipleRules(t *testing.T) {
	analyzer := NewAnalyzer()
	
	changes := []types.DiffChange{
		{
			FilePath:   "debug.go",
			FileType:   "go",
			AddedLines: []string{
				"x := db.Query()",
				"fmt.Println(x)",
				"// TODO: handle errors",
			},
		},
	}

	suggestions := analyzer.Analyze(changes)

	if len(suggestions) != 3 {
		t.Errorf("expected 3 suggestions, got %d", len(suggestions))
		return
	}

	severities := make(map[string]bool)
	for _, s := range suggestions {
		severities[s.Severity] = true
	}

	if !severities["high"] {
		t.Error("expected high severity suggestion for error handling")
	}

	if !severities["medium"] {
		t.Error("expected medium severity suggestion for TODO")
	}

	if !severities["low"] {
		t.Error("expected low severity suggestion for debug statement")
	}
}

func TestAnalyzer_EmptyInput(t *testing.T) {
	analyzer := NewAnalyzer()
	
	suggestions := analyzer.Analyze([]types.DiffChange{})

	if len(suggestions) != 0 {
		t.Errorf("expected 0 suggestions for empty input, got %d", len(suggestions))
	}
}

func TestAnalyzer_NoAddedLines(t *testing.T) {
	analyzer := NewAnalyzer()
	
	changes := []types.DiffChange{
		{
			FilePath:     "main.go",
			FileType:     "go",
			AddedLines:   []string{},
			RemovedLines: []string{"x := something()"},
		},
	}

	suggestions := analyzer.Analyze(changes)

	if len(suggestions) != 0 {
		t.Errorf("expected 0 suggestions when no lines added, got %d", len(suggestions))
	}
}

func TestAnalyzer_SuggestionFields(t *testing.T) {
	analyzer := NewAnalyzer()
	
	changes := []types.DiffChange{
		{
			FilePath:   "test.go",
			FileType:   "go",
			AddedLines: []string{"x := doSomething()"},
		},
	}

	suggestions := analyzer.Analyze(changes)

	if len(suggestions) != 1 {
		t.Fatalf("expected 1 suggestion, got %d", len(suggestions))
	}

	s := suggestions[0]

	if s.ID == "" {
		t.Error("suggestion ID should not be empty")
	}

	if s.FilePath != "test.go" {
		t.Errorf("expected FilePath 'test.go', got '%s'", s.FilePath)
	}

	if s.Confidence != 0.8 {
		t.Errorf("expected Confidence 0.8, got %f", s.Confidence)
	}

	if s.Message == "" {
		t.Error("suggestion Message should not be empty")
	}
}