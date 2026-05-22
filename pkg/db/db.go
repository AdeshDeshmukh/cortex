package db

import (
	"database/sql"
	"embed"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed schema.sql
var schemaFS embed.FS

type DB struct {
	conn *sql.DB
	path string
}

func New(repoPath string) (*DB, error) {
	cortexDir := filepath.Join(repoPath, ".cortex")
	if err := os.MkdirAll(cortexDir, 0755); err != nil {
		return nil, fmt.Errorf("create .cortex directory: %w", err)
	}

	dbPath := filepath.Join(cortexDir, "feedback.db")

	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	db := &DB{
		conn: conn,
		path: dbPath,
	}

	if err := db.migrate(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("migrate database: %w", err)
	}

	return db, nil
}

func (db *DB) migrate() error {
	schema, err := schemaFS.ReadFile("schema.sql")
	if err != nil {
		return fmt.Errorf("read schema: %w", err)
	}

	if _, err := db.conn.Exec(string(schema)); err != nil {
		return fmt.Errorf("execute schema: %w", err)
	}

	return nil
}

func (db *DB) Close() error {
	if db.conn != nil {
		return db.conn.Close()
	}
	return nil
}

func (db *DB) SaveReview(commitHash, repository, language string, diffSize, durationMs int) (int64, error) {
	result, err := db.conn.Exec(
		`INSERT INTO reviews (commit_hash, repository, timestamp, diff_size, language, duration_ms)
		 VALUES (?, ?, strftime('%s', 'now'), ?, ?, ?)`,
		commitHash, repository, diffSize, language, durationMs,
	)
	if err != nil {
		return 0, fmt.Errorf("insert review: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get review id: %w", err)
	}

	return id, nil
}

func (db *DB) SaveSuggestion(reviewID int64, suggestionType, severity, message, filePath string, lineNumber int, suggestion string, confidence float64, source string) (int64, error) {
	result, err := db.conn.Exec(
		`INSERT INTO suggestions (review_id, type, severity, message, file_path, line_number, suggestion, confidence, source)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		reviewID, suggestionType, severity, message, filePath, lineNumber, suggestion, confidence, source,
	)
	if err != nil {
		return 0, fmt.Errorf("insert suggestion: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get suggestion id: %w", err)
	}

	return id, nil
}

func (db *DB) SaveFeedback(suggestionID int64, action, userReason string) error {
	_, err := db.conn.Exec(
		`INSERT INTO feedback (suggestion_id, action, user_reason, timestamp)
		 VALUES (?, ?, ?, strftime('%s', 'now'))`,
		suggestionID, action, userReason,
	)
	if err != nil {
		return fmt.Errorf("insert feedback: %w", err)
	}

	return nil
}

func (db *DB) GetAcceptanceRate(suggestionType string) (float64, error) {
	var accepted, total int

	err := db.conn.QueryRow(
		`SELECT 
			COUNT(CASE WHEN f.action = 'accept' THEN 1 END) as accepted,
			COUNT(*) as total
		 FROM feedback f
		 JOIN suggestions s ON f.suggestion_id = s.id
		 WHERE s.type = ?`,
		suggestionType,
	).Scan(&accepted, &total)

	if err != nil {
		return 0, fmt.Errorf("query acceptance rate: %w", err)
	}

	if total == 0 {
		return 0, nil
	}

	return float64(accepted) / float64(total), nil
}

func (db *DB) GetTotalReviews() (int, error) {
	var count int
	err := db.conn.QueryRow(`SELECT COUNT(*) FROM reviews`).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("query total reviews: %w", err)
	}
	return count, nil
}

func (db *DB) GetTotalSuggestions() (int, error) {
	var count int
	err := db.conn.QueryRow(`SELECT COUNT(*) FROM suggestions`).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("query total suggestions: %w", err)
	}
	return count, nil
}

func (db *DB) GetTotalFeedback() (int, error) {
	var count int
	err := db.conn.QueryRow(`SELECT COUNT(*) FROM feedback`).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("query total feedback: %w", err)
	}
	return count, nil
}