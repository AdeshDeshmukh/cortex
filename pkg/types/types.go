package types

import "time"

type DiffChange struct {
	FilePath     string
	FileType     string
	ChangeType   string
	AddedLines   []string
	RemovedLines []string
	Context      string
	StartLine    int
	EndLine      int
}

type CodeReview struct {
	CommitHash   string
	Timestamp    time.Time
	Changes      []DiffChange
	Suggestions  []Suggestion
	TotalLines   int
	FilesChanged int
}

type Suggestion struct {
	ID         string
	Type       string
	Severity   string
	Message    string
	FilePath   string
	LineNumber int
	Suggestion string
	Confidence float64
	CreatedAt  time.Time
	Source     string
}

type Feedback struct {
	SuggestionID string
	Action       string
	Reason       string
	Timestamp    time.Time
}

type Config struct {
	ModelPath      string
	MaxThreads     int
	Verbose        bool
	AutoFix        bool
	ThresholdScore float64
}

type ReviewContext struct {
	RepoPath string
	Branch   string
	Author   string
	IsMerge  bool
	DiffSize int
	Language string
}