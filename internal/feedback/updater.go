package feedback

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/AdeshDeshmukh/cortex/pkg/types"
)

type Updater struct {
	pythonPath string
	scriptPath string
}

type updatePayload struct {
	Context        updateContext `json:"context"`
	SuggestionType string       `json:"suggestion_type"`
	Action         string       `json:"action"`
}

type updateContext struct {
	Language  string `json:"language"`
	FileType  string `json:"file_type"`
	DiffSize  int    `json:"diff_size"`
	TimeOfDay int    `json:"time_of_day"`
}

func NewUpdater() *Updater {
	return &Updater{
		pythonPath: "/opt/homebrew/bin/python3",
		scriptPath: "python/bandit/linucb.py",
	}
}

func (u *Updater) RecordFeedback(suggestion types.Suggestion, action string, reviewContext types.ReviewContext) error {
	if action != "accept" && action != "reject" {
		return nil
	}

	payload := u.buildUpdatePayload(suggestion, action, reviewContext)

	if err := u.callPythonUpdate(payload); err != nil {
		return fmt.Errorf("update LinUCB model: %w", err)
	}

	return nil
}

func (u *Updater) buildUpdatePayload(suggestion types.Suggestion, action string, reviewContext types.ReviewContext) updatePayload {
	language := reviewContext.Language
	if language == "" {
		language = "go"
	}

	diffSize := reviewContext.DiffSize
	if diffSize == 0 {
		diffSize = 50
	}

	return updatePayload{
		Context: updateContext{
			Language:  language,
			FileType:  language,
			DiffSize:  diffSize,
			TimeOfDay: 14,
		},
		SuggestionType: suggestion.Type,
		Action:         action,
	}
}

func (u *Updater) callPythonUpdate(payload updatePayload) error {
	inputJSON, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	// Find the project root by looking for go.mod
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get working directory: %w", err)
	}

	scriptPath := u.scriptPath
	// If relative path doesn't exist, try to find it from project root
	if _, err := os.Stat(scriptPath); err != nil {
		projectRoot := cwd
		for projectRoot != "/" {
			if _, err := os.Stat(filepath.Join(projectRoot, "go.mod")); err == nil {
				scriptPath = filepath.Join(projectRoot, u.scriptPath)
				break
			}
			projectRoot = filepath.Dir(projectRoot)
		}
	}

	cmd := exec.Command(u.pythonPath, scriptPath, "update")

	cmd.Stdin = bytes.NewReader(inputJSON)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("python update failed: %w, stderr: %s", err, stderr.String())
	}

	var result map[string]interface{}
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		return fmt.Errorf("parse python response: %w", err)
	}

	if status, ok := result["status"].(string); !ok || status != "ok" {
		return fmt.Errorf("unexpected response: %s", stdout.String())
	}

	return nil
}