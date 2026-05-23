package rl

import (
	"testing"

	"github.com/AdeshDeshmukh/cortex/pkg/types"
)

func TestNewRanker(t *testing.T) {
	ranker := NewRanker()

	if ranker.pythonPath != "python3" {
		t.Errorf("expected pythonPath 'python3', got '%s'", ranker.pythonPath)
	}

	if ranker.scriptPath != "python/bandit/linucb.py" {
		t.Errorf("expected scriptPath 'python/bandit/linucb.py', got '%s'", ranker.scriptPath)
	}
}

func TestRankSuggestions_EmptyInput(t *testing.T) {
	ranker := NewRanker()

	suggestions := []types.Suggestion{}
	context := types.ReviewContext{}

	ranked, err := ranker.RankSuggestions(suggestions, context)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(ranked) != 0 {
		t.Errorf("expected empty result, got %d suggestions", len(ranked))
	}
}

func TestRankSuggestions_WithData(t *testing.T) {
	ranker := NewRanker()

	suggestions := []types.Suggestion{
		{ID: "sugg_1", Type: "error-handling", Severity: "high"},
		{ID: "sugg_2", Type: "performance", Severity: "medium"},
		{ID: "sugg_3", Type: "security", Severity: "critical"},
	}

	context := types.ReviewContext{
		Language: "go",
		DiffSize: 50,
	}

	ranked, err := ranker.RankSuggestions(suggestions, context)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(ranked) != 3 {
		t.Fatalf("expected 3 ranked suggestions, got %d", len(ranked))
	}

	for i, s := range ranked {
		t.Logf("Rank %d: %s (type: %s, severity: %s)", i+1, s.ID, s.Type, s.Severity)
	}
}

func TestExtractContext(t *testing.T) {
	ranker := NewRanker()

	context := types.ReviewContext{
		Language: "python",
		DiffSize: 100,
	}

	suggestions := []types.Suggestion{
		{ID: "sugg_1", Type: "error-handling"},
	}

	linucbCtx := ranker.extractContext(context, suggestions)

	if linucbCtx.Language != "python" {
		t.Errorf("expected language 'python', got '%s'", linucbCtx.Language)
	}

	if linucbCtx.DiffSize != 100 {
		t.Errorf("expected diff_size 100, got %d", linucbCtx.DiffSize)
	}

	if linucbCtx.TimeOfDay != 14 {
		t.Errorf("expected time_of_day 14, got %d", linucbCtx.TimeOfDay)
	}
}

func TestExtractContext_Defaults(t *testing.T) {
	ranker := NewRanker()

	context := types.ReviewContext{}

	suggestions := []types.Suggestion{
		{ID: "sugg_1", Type: "error-handling"},
	}

	linucbCtx := ranker.extractContext(context, suggestions)

	if linucbCtx.Language != "go" {
		t.Errorf("expected default language 'go', got '%s'", linucbCtx.Language)
	}

	if linucbCtx.DiffSize != 50 {
		t.Errorf("expected default diff_size 50, got %d", linucbCtx.DiffSize)
	}
}
