package llm

import (
	"fmt"
	"strings"

	"github.com/AdeshDeshmukh/cortex/pkg/types"
)

type Engine interface {
	Review(change types.DiffChange) ([]types.Suggestion, error)
	IsAvailable() bool
	Name() string
}

type MockEngine struct{}

type EngineConfig struct {
	UseMock   bool
	ModelPath string
	Threads   int
}

func NewEngine(config EngineConfig) Engine {
	if config.UseMock || config.ModelPath == "" {
		return &MockEngine{}
	}

	return &MockEngine{}
}

func (m *MockEngine) Review(change types.DiffChange) ([]types.Suggestion, error) {
	var suggestions []types.Suggestion

	for i, line := range change.AddedLines {
		if s := m.analyzeLine(line, change.FilePath, i); s != nil {
			suggestions = append(suggestions, *s)
		}
	}

	return suggestions, nil
}

func (m *MockEngine) IsAvailable() bool {
	return true
}

func (m *MockEngine) Name() string {
	return "mock"
}

func (m *MockEngine) analyzeLine(line, filePath string, lineNum int) *types.Suggestion {
	line = strings.TrimSpace(line)

	checks := []struct {
		contains string
		message  string
		severity string
		kind     string
	}{
		{
			contains: "if err",
			message:  "Error handling detected - ensure all error paths are covered",
			severity: "medium",
			kind:     "error-handling",
		},
		{
			contains: "panic(",
			message:  "Avoid panic in production code - return errors instead",
			severity: "high",
			kind:     "error-handling",
		},
		{
			contains: "interface{}",
			message:  "Consider using a specific type instead of empty interface",
			severity: "low",
			kind:     "type-safety",
		},
		{
			contains: "time.Sleep",
			message:  "time.Sleep in production code may indicate a design issue",
			severity: "medium",
			kind:     "performance",
		},
		{
			contains: "SELECT *",
			message:  "Avoid SELECT * - specify columns explicitly for better performance",
			severity: "medium",
			kind:     "performance",
		},
		{
			contains: "http.Get(",
			message:  "HTTP call without timeout context - consider adding deadline",
			severity: "high",
			kind:     "reliability",
		},
	}

	for _, check := range checks {
		if strings.Contains(line, check.contains) {
			return &types.Suggestion{
				ID:         fmt.Sprintf("llm_mock_%d", lineNum),
				Type:       check.kind,
				Severity:   check.severity,
				Message:    check.message,
				FilePath:   filePath,
				LineNumber: lineNum,
				Confidence: 0.6,
				Source:     "llm",
			}
		}
	}

	return nil
}			