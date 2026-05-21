package feedback

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/AdeshDeshmukh/cortex/pkg/types"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(repoPath string) (*Storage, error) {
	cortexDir := filepath.Join(repoPath, ".cortex")

	if err := os.MkdirAll(cortexDir, 0755); err != nil {
		return nil, fmt.Errorf("create cortex directory: %w", err)
	}

	dbPath := filepath.Join(cortexDir, "feedback.db")

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	storage := &Storage{db: db}

	if err := storage.migrate(); err != nil {
		return nil, fmt.Errorf("migrate database: %w", err)
	}

	return storage, nil
}

func (s *Storage) migrate() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS reviews (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			commit_hash TEXT,
			branch TEXT,
			timestamp INTEGER NOT NULL,
			files_changed INTEGER,
			total_lines INTEGER
		)`,

		`CREATE TABLE IF NOT EXISTS suggestions (
			id TEXT PRIMARY KEY,
			review_id INTEGER,
			type TEXT,
			severity TEXT,
			message TEXT,
			file_path TEXT,
			line_number INTEGER,
			confidence REAL,
			source TEXT,
			FOREIGN KEY(review_id) REFERENCES reviews(id)
		)`,

		`CREATE TABLE IF NOT EXISTS feedback (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			suggestion_id TEXT,
			action TEXT NOT NULL,
			reason TEXT,
			timestamp INTEGER NOT NULL,
			FOREIGN KEY(suggestion_id) REFERENCES suggestions(id)
		)`,

		`CREATE INDEX IF NOT EXISTS idx_feedback_action ON feedback(action)`,
		`CREATE INDEX IF NOT EXISTS idx_suggestions_type ON suggestions(type)`,
		`CREATE INDEX IF NOT EXISTS idx_suggestions_source ON suggestions(source)`,
	}

	for _, query := range queries {
		if _, err := s.db.Exec(query); err != nil {
			return fmt.Errorf("execute migration: %w", err)
		}
	}

	return nil
}

func (s *Storage) SaveReview(review types.CodeReview) (int64, error) {
	result, err := s.db.Exec(
		`INSERT INTO reviews (commit_hash, branch, timestamp, files_changed, total_lines)
		 VALUES (?, ?, ?, ?, ?)`,
		review.CommitHash,
		"",
		review.Timestamp.Unix(),
		review.FilesChanged,
		review.TotalLines,
	)
	if err != nil {
		return 0, fmt.Errorf("save review: %w", err)
	}

	return result.LastInsertId()
}

func (s *Storage) SaveSuggestion(suggestion types.Suggestion, reviewID int64) error {
	_, err := s.db.Exec(
		`INSERT OR IGNORE INTO suggestions
		 (id, review_id, type, severity, message, file_path, line_number, confidence, source)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		suggestion.ID,
		reviewID,
		suggestion.Type,
		suggestion.Severity,
		suggestion.Message,
		suggestion.FilePath,
		suggestion.LineNumber,
		suggestion.Confidence,
		suggestion.Source,
	)
	if err != nil {
		return fmt.Errorf("save suggestion: %w", err)
	}

	return nil
}

func (s *Storage) SaveFeedback(feedback types.Feedback) error {
	_, err := s.db.Exec(
		`INSERT INTO feedback (suggestion_id, action, reason, timestamp)
		 VALUES (?, ?, ?, ?)`,
		feedback.SuggestionID,
		feedback.Action,
		feedback.Reason,
		feedback.Timestamp.Unix(),
	)
	if err != nil {
		return fmt.Errorf("save feedback: %w", err)
	}

	return nil
}

func (s *Storage) GetAcceptanceRate(suggestionType string) (float64, error) {
	var accepted, total int

	err := s.db.QueryRow(
		`SELECT
			COUNT(CASE WHEN f.action = 'accept' THEN 1 END),
			COUNT(*)
		FROM feedback f
		JOIN suggestions s ON f.suggestion_id = s.id
		WHERE s.type = ?`,
		suggestionType,
	).Scan(&accepted, &total)

	if err != nil {
		return 0, fmt.Errorf("get acceptance rate: %w", err)
	}

	if total == 0 {
		return 0, nil
	}

	return float64(accepted) / float64(total), nil
}

func (s *Storage) GetFeedbackStats() (map[string]int, error) {
	rows, err := s.db.Query(
		`SELECT action, COUNT(*) FROM feedback GROUP BY action`,
	)
	if err != nil {
		return nil, fmt.Errorf("get feedback stats: %w", err)
	}
	defer rows.Close()

	stats := make(map[string]int)

	for rows.Next() {
		var action string
		var count int
		if err := rows.Scan(&action, &count); err != nil {
			return nil, fmt.Errorf("scan feedback stats: %w", err)
		}
		stats[action] = count
	}

	return stats, nil
}

func (s *Storage) CreateReview(changes []types.DiffChange, suggestions []types.Suggestion) (types.CodeReview, error) {
	totalLines := 0
	for _, change := range changes {
		totalLines += len(change.AddedLines) + len(change.RemovedLines)
	}

	review := types.CodeReview{
		Timestamp:    time.Now(),
		Changes:      changes,
		Suggestions:  suggestions,
		FilesChanged: len(changes),
		TotalLines:   totalLines,
	}

	return review, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}