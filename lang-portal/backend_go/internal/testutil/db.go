package testutil

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

// SetupTestDB creates a temporary SQLite database for testing
func SetupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	// Create temporary database file
	dbFile, err := os.CreateTemp("", "test-*.db")
	if err != nil {
		t.Fatalf("failed to create temp db file: %v", err)
	}
	dbFile.Close()

	// Open database connection
	db, err := sql.Open("sqlite3", dbFile.Name())
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	// Run migrations
	if err := runMigrations(db); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	// Clean up database file when test completes
	t.Cleanup(func() {
		db.Close()
		os.Remove(dbFile.Name())
	})

	return db
}

// runMigrations runs the schema migrations on the test database
func runMigrations(db *sql.DB) error {
	// Create tables
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS words (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			japanese TEXT NOT NULL,
			romaji TEXT NOT NULL,
			english TEXT NOT NULL,
			parts TEXT
		);

		CREATE TABLE IF NOT EXISTS groups (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL
		);

		CREATE TABLE IF NOT EXISTS words_groups (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			word_id INTEGER,
			group_id INTEGER,
			FOREIGN KEY (word_id) REFERENCES words(id),
			FOREIGN KEY (group_id) REFERENCES groups(id)
		);

		CREATE TABLE IF NOT EXISTS study_sessions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			group_id INTEGER,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			study_activity_id INTEGER,
			FOREIGN KEY (group_id) REFERENCES groups(id)
		);

		CREATE TABLE IF NOT EXISTS study_activities (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			study_session_id INTEGER,
			group_id INTEGER,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (study_session_id) REFERENCES study_sessions(id),
			FOREIGN KEY (group_id) REFERENCES groups(id)
		);

		CREATE TABLE IF NOT EXISTS word_review_items (
			word_id INTEGER,
			study_session_id INTEGER,
			correct BOOLEAN,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (word_id) REFERENCES words(id),
			FOREIGN KEY (study_session_id) REFERENCES study_sessions(id)
		);
	`)
	return err
} 