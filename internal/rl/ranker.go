package rl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"sort"

	"github.com/AdeshDeshmukh/cortex/internal/utils"
	"github.com/AdeshDeshmukh/cortex/pkg/types"
)

type Ranker struct {
	pythonPath string
	scriptPath string
}

type linucbContext struct {
	Language  string `json:"language"`
	FileType  string `json:"file_type"`
	DiffSize  int    `json:"diff_size"`
	TimeOfDay int    `json:"time_of_day"`
}

type predictRequest struct {
	Context     linucbContext          `json:"context"`
	Suggestions []suggestionForRanking `json:"suggestions"`
}

type suggestionForRanking struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type scoreResponse map[string]float64

func NewRanker() *Ranker {
	return &Ranker{
		pythonPath: utils.DetectPython(),
		scriptPath: "python/bandit/linucb.py",
	}
}

func (r *Ranker) RankSuggestions(suggestions []types.Suggestion, context types.ReviewContext) ([]types.Suggestion, error) {
	if len(suggestions) == 0 {
		return suggestions, nil
	}

	ctx := r.extractContext(context, suggestions)
	suggForRanking := make([]suggestionForRanking, len(suggestions))
	for i, s := range suggestions {
		suggForRanking[i] = suggestionForRanking{
			ID:   s.ID,
			Type: s.Type,
		}
	}

	req := predictRequest{
		Context:     ctx,
		Suggestions: suggForRanking,
	}

	scores, err := r.callPython("predict", req)
	if err != nil {
		return suggestions, nil
	}

	ranked := make([]types.Suggestion, len(suggestions))
	copy(ranked, suggestions)

	sort.Slice(ranked, func(i, j int) bool {
		scoreI := scores[ranked[i].ID]
		scoreJ := scores[ranked[j].ID]
		return scoreI > scoreJ
	})

	return ranked, nil
}

func (r *Ranker) callPython(command string, data interface{}) (scoreResponse, error) {
	inputJSON, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("marshal input: %w", err)
	}

	cmd := exec.Command(r.pythonPath, r.scriptPath, command)
	cmd.Stdin = bytes.NewReader(inputJSON)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("python execution failed: %w, stderr: %s", err, stderr.String())
	}

	var scores scoreResponse
	if err := json.Unmarshal(stdout.Bytes(), &scores); err != nil {
		return nil, fmt.Errorf("unmarshal output: %w, stdout: %s", err, stdout.String())
	}

	return scores, nil
}

func (r *Ranker) extractContext(context types.ReviewContext, suggestions []types.Suggestion) linucbContext {
	language := context.Language
	if language == "" && len(suggestions) > 0 {
		language = "go"
	}

	diffSize := context.DiffSize
	if diffSize == 0 {
		diffSize = 50
	}

	return linucbContext{
		Language:  language,
		FileType:  language,
		DiffSize:  diffSize,
		TimeOfDay: 14,
	}
}