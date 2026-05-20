package llm

import (
	"strings"
	"testing"

	"github.com/AdeshDeshmukh/cortex/pkg/types"
)

func TestMockEngine_IsAvailable(t *testing.T) {
	engine := &MockEngine{}

	if !engine.IsAvailable() {
		t.Error("MockEngine should always be available")
	}
}

func TestMockEngine_Name(t *testing.T) {
	engine := &MockEngine{}

	if engine.Name() != "mock" {
		t.Errorf("expected name 'mock', got '%s'", engine.Name())
	}
}

func TestMockEngine_Review_EmptyChange(t *testing.T) {
	engine := &MockEngine{}

	change := types.DiffChange{
		FilePath:   "main.go",
		FileType:   "go",
		AddedLines: []string{},
	}

	suggestions, err := engine.Review(change)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(suggestions) != 0 {
		t.Errorf("expected 0 suggestions for empty change, got %d", len(suggestions))
	}
}

func TestMockEngine_Review_DetectsPanic(t *testing.T) {
	engine := &MockEngine{}

	change := types.DiffChange{
		FilePath:   "main.go",
		FileType:   "go",
		AddedLines: []string{"panic(\"something went wrong\")"},
	}

	suggestions, err := engine.Review(change)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(suggestions) != 1 {
		t.Fatalf("expected 1 suggestion, got %d", len(suggestions))
	}

	if suggestions[0].Severity != "high" {
		t.Errorf("expected high severity, got %s", suggestions[0].Severity)
	}

	if suggestions[0].Source != "llm" {
		t.Errorf("expected source 'llm', got '%s'", suggestions[0].Source)
	}
}

func TestMockEngine_Review_DetectsTimeSleep(t *testing.T) {
	engine := &MockEngine{}

	change := types.DiffChange{
		FilePath:   "handler.go",
		FileType:   "go",
		AddedLines: []string{"time.Sleep(5 * time.Second)"},
	}

	suggestions, err := engine.Review(change)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(suggestions) != 1 {
		t.Fatalf("expected 1 suggestion, got %d", len(suggestions))
	}

	if suggestions[0].Type != "performance" {
		t.Errorf("expected performance type, got %s", suggestions[0].Type)
	}
}

func TestMockEngine_Review_DetectsHTTPWithoutTimeout(t *testing.T) {
	engine := &MockEngine{}

	change := types.DiffChange{
		FilePath:   "client.go",
		FileType:   "go",
		AddedLines: []string{"resp, err := http.Get(url)"},
	}

	suggestions, err := engine.Review(change)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(suggestions) != 1 {
		t.Fatalf("expected 1 suggestion, got %d", len(suggestions))
	}

	if suggestions[0].Severity != "high" {
		t.Errorf("expected high severity, got %s", suggestions[0].Severity)
	}
}

func TestMockEngine_Review_SuggestionFields(t *testing.T) {
	engine := &MockEngine{}

	change := types.DiffChange{
		FilePath:   "service.go",
		FileType:   "go",
		AddedLines: []string{"panic(\"critical error\")"},
	}

	suggestions, err := engine.Review(change)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(suggestions) != 1 {
		t.Fatalf("expected 1 suggestion, got %d", len(suggestions))
	}

	s := suggestions[0]

	if s.ID == "" {
		t.Error("suggestion ID should not be empty")
	}

	if s.FilePath != "service.go" {
		t.Errorf("expected FilePath 'service.go', got '%s'", s.FilePath)
	}

	if s.Confidence != 0.6 {
		t.Errorf("expected Confidence 0.6, got %f", s.Confidence)
	}

	if s.Message == "" {
		t.Error("suggestion Message should not be empty")
	}
}

func TestNewEngine_ReturnsMock(t *testing.T) {
	config := EngineConfig{
		UseMock: true,
	}

	engine := NewEngine(config)

	if engine.Name() != "mock" {
		t.Errorf("expected mock engine, got %s", engine.Name())
	}
}

func TestNewEngine_EmptyPathReturnsMock(t *testing.T) {
	config := EngineConfig{
		UseMock:   false,
		ModelPath: "",
	}

	engine := NewEngine(config)

	if engine.Name() != "mock" {
		t.Errorf("expected mock engine for empty path, got %s", engine.Name())
	}
}

func TestPromptBuilder_Build_ContainsFilePath(t *testing.T) {
	builder := NewPromptBuilder()

	change := types.DiffChange{
		FilePath:   "internal/service.go",
		FileType:   "go",
		AddedLines: []string{"func newService() {}"},
	}

	prompt := builder.Build(change)

	if !strings.Contains(prompt, "internal/service.go") {
		t.Error("prompt should contain file path")
	}
}

func TestPromptBuilder_Build_ContainsLanguageGuide(t *testing.T) {
	builder := NewPromptBuilder()

	change := types.DiffChange{
		FilePath:   "main.go",
		FileType:   "go",
		AddedLines: []string{"func main() {}"},
	}

	prompt := builder.Build(change)

	if !strings.Contains(prompt, "Go-specific") {
		t.Error("prompt should contain Go-specific instructions")
	}
}

func TestPromptBuilder_Build_ContainsAddedLines(t *testing.T) {
	builder := NewPromptBuilder()

	change := types.DiffChange{
		FilePath:   "main.go",
		FileType:   "go",
		AddedLines: []string{"func myFunction() {}"},
	}

	prompt := builder.Build(change)

	if !strings.Contains(prompt, "func myFunction()") {
		t.Error("prompt should contain added lines")
	}
}

func TestPromptBuilder_Build_RespectsMaxLines(t *testing.T) {
	builder := NewPromptBuilder()
	builder.maxLines = 3

	lines := make([]string, 10)
	for i := range lines {
		lines[i] = "some code line"
	}

	change := types.DiffChange{
		FilePath:   "main.go",
		FileType:   "go",
		AddedLines: lines,
	}

	prompt := builder.Build(change)

	if !strings.Contains(prompt, "7 more lines") {
		t.Error("prompt should indicate truncated lines")
	}
}

func TestPromptBuilder_Build_UnknownLanguage(t *testing.T) {
	builder := NewPromptBuilder()

	change := types.DiffChange{
		FilePath:   "script.sh",
		FileType:   "unknown",
		AddedLines: []string{"echo hello"},
	}

	prompt := builder.Build(change)

	if prompt == "" {
		t.Error("prompt should not be empty for unknown language")
	}
}

func TestResponseParser_Parse_NoIssues(t *testing.T) {
	parser := NewResponseParser()

	suggestions := parser.Parse("NO_ISSUES", "main.go")

	if len(suggestions) != 0 {
		t.Errorf("expected 0 suggestions for NO_ISSUES, got %d", len(suggestions))
	}
}

func TestResponseParser_Parse_SingleIssue(t *testing.T) {
	parser := NewResponseParser()

	raw := `ISSUE: Missing error check
SEVERITY: high
TYPE: error-handling
SUGGESTION: Add if err != nil check
---`

	suggestions := parser.Parse(raw, "main.go")

	if len(suggestions) != 1 {
		t.Fatalf("expected 1 suggestion, got %d", len(suggestions))
	}

	s := suggestions[0]

	if s.Severity != "high" {
		t.Errorf("expected high severity, got %s", s.Severity)
	}

	if s.Type != "error-handling" {
		t.Errorf("expected error-handling type, got %s", s.Type)
	}

	if s.FilePath != "main.go" {
		t.Errorf("expected FilePath 'main.go', got '%s'", s.FilePath)
	}
}

func TestResponseParser_Parse_MultipleIssues(t *testing.T) {
	parser := NewResponseParser()

	raw := `ISSUE: Missing error check
SEVERITY: high
TYPE: error-handling
SUGGESTION: Add error handling
---
ISSUE: Hardcoded timeout
SEVERITY: medium
TYPE: reliability
SUGGESTION: Use configuration value
---`

	suggestions := parser.Parse(raw, "service.go")

	if len(suggestions) != 2 {
		t.Fatalf("expected 2 suggestions, got %d", len(suggestions))
	}
}

func TestResponseParser_NormalizeSeverity(t *testing.T) {
	parser := NewResponseParser()

	tests := []struct {
		input    string
		expected string
	}{
		{"HIGH", "high"},
		{"high", "high"},
		{"Critical", "critical"},
		{"MEDIUM", "medium"},
		{"low", "low"},
		{"urgent", "medium"},
		{"important", "medium"},
		{"  high  ", "high"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parser.normalizeSeverity(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeSeverity(%s) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}