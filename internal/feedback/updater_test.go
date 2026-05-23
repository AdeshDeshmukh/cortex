package feedback

import (
	"testing"

	"github.com/AdeshDeshmukh/cortex/pkg/types"
)

func TestNewUpdater(t *testing.T) {
	updater := NewUpdater()

	if updater.pythonPath != "/opt/homebrew/bin/python3" {
		t.Errorf("expected pythonPath '/opt/homebrew/bin/python3', got '%s'", updater.pythonPath)
	}

	if updater.scriptPath != "python/bandit/linucb.py" {
		t.Errorf("expected scriptPath 'python/bandit/linucb.py', got '%s'", updater.scriptPath)
	}
}

func TestRecordFeedback_SkipAction(t *testing.T) {
	updater := NewUpdater()

	suggestion := types.Suggestion{
		ID:   "sugg_1",
		Type: "error-handling",
	}

	context := types.ReviewContext{
		Language: "go",
		DiffSize: 50,
	}

	err := updater.RecordFeedback(suggestion, "skip", context)
	if err != nil {
		t.Errorf("skip action should not return error, got: %v", err)
	}
}

func TestRecordFeedback_AcceptAction(t *testing.T) {
	updater := NewUpdater()

	suggestion := types.Suggestion{
		ID:   "sugg_1",
		Type: "error-handling",
	}

	context := types.ReviewContext{
		Language: "go",
		DiffSize: 50,
	}

	err := updater.RecordFeedback(suggestion, "accept", context)
	if err != nil {
		t.Errorf("accept action should not return error, got: %v", err)
	}
}

func TestRecordFeedback_RejectAction(t *testing.T) {
	updater := NewUpdater()

	suggestion := types.Suggestion{
		ID:   "sugg_1",
		Type: "security",
	}

	context := types.ReviewContext{
		Language: "go",
		DiffSize: 30,
	}

	err := updater.RecordFeedback(suggestion, "reject", context)
	if err != nil {
		t.Errorf("reject action should not return error, got: %v", err)
	}
}

func TestBuildUpdatePayload_Defaults(t *testing.T) {
	updater := NewUpdater()

	suggestion := types.Suggestion{
		ID:   "sugg_1",
		Type: "performance",
	}

	context := types.ReviewContext{}

	payload := updater.buildUpdatePayload(suggestion, "accept", context)

	if payload.Context.Language != "go" {
		t.Errorf("expected default language 'go', got '%s'", payload.Context.Language)
	}

	if payload.Context.DiffSize != 50 {
		t.Errorf("expected default diff_size 50, got %d", payload.Context.DiffSize)
	}

	if payload.SuggestionType != "performance" {
		t.Errorf("expected suggestion_type 'performance', got '%s'", payload.SuggestionType)
	}

	if payload.Action != "accept" {
		t.Errorf("expected action 'accept', got '%s'", payload.Action)
	}
}

func TestBuildUpdatePayload_WithContext(t *testing.T) {
	updater := NewUpdater()

	suggestion := types.Suggestion{
		ID:   "sugg_1",
		Type: "security",
	}

	context := types.ReviewContext{
		Language: "python",
		DiffSize: 120,
	}

	payload := updater.buildUpdatePayload(suggestion, "reject", context)

	if payload.Context.Language != "python" {
		t.Errorf("expected language 'python', got '%s'", payload.Context.Language)
	}

	if payload.Context.DiffSize != 120 {
		t.Errorf("expected diff_size 120, got %d", payload.Context.DiffSize)
	}

	if payload.Action != "reject" {
		t.Errorf("expected action 'reject', got '%s'", payload.Action)
	}
}