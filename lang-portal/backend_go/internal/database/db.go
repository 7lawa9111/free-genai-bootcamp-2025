package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// ConfigDB holds database configuration
type ConfigDB struct {
	*sql.DB
}

func InitDB(dbPath string) error {
	var err error
	// Ensure the database file exists
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}

	// Test the connection
	if err := DB.Ping(); err != nil {
		return fmt.Errorf("error connecting to database: %v", err)
	}

	// Configure connection pool
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(10)
	DB.SetConnMaxLifetime(5 * time.Minute)
	DB.SetConnMaxIdleTime(5 * time.Minute)

	// Run migrations
	if err := runMigrations(DB); err != nil {
		return fmt.Errorf("error running migrations: %v", err)
	}

	// Seed initial data if needed
	if err := seedData(DB); err != nil {
		return fmt.Errorf("error seeding data: %v", err)
	}

	return nil
}

func runMigrations(db *sql.DB) error {
	// Enable SQLite foreign key support
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return fmt.Errorf("error enabling foreign keys: %v", err)
	}

	// First migration
	if _, err := db.Exec(readMigrationFile("0001_init.sql")); err != nil {
		return fmt.Errorf("error in migration 1: %v", err)
	}
	fmt.Println("Migration 1 completed successfully")

	// Second migration - run each statement separately
	fmt.Println("Running migration 2...")

	// Drop tables
	if _, err := db.Exec("DROP TABLE IF EXISTS vocabulary_quiz_answers"); err != nil {
		return fmt.Errorf("error dropping vocabulary_quiz_answers: %v", err)
	}
	if _, err := db.Exec("DROP TABLE IF EXISTS study_sessions"); err != nil {
		return fmt.Errorf("error dropping study_sessions: %v", err)
	}
	if _, err := db.Exec("DROP TABLE IF EXISTS study_activities"); err != nil {
		return fmt.Errorf("error dropping study_activities: %v", err)
	}

	// Create study_activities table
	if _, err := db.Exec(`
        CREATE TABLE study_activities (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            group_id INTEGER NOT NULL,
            activity_type TEXT NOT NULL,
            created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
            completed_at DATETIME,
            score REAL,
            settings TEXT,
            confidence_score FLOAT,
            FOREIGN KEY (group_id) REFERENCES groups(id)
        )`); err != nil {
		return fmt.Errorf("error creating study_activities: %v", err)
	}

	// Create study_sessions table
	if _, err := db.Exec(`
        CREATE TABLE study_sessions (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            group_id INTEGER NOT NULL,
            study_activity_id INTEGER NOT NULL,
            created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
            completed_at DATETIME,
            FOREIGN KEY (group_id) REFERENCES groups(id),
            FOREIGN KEY (study_activity_id) REFERENCES study_activities(id)
        )`); err != nil {
		return fmt.Errorf("error creating study_sessions: %v", err)
	}

	// Create vocabulary_quiz_answers table
	if _, err := db.Exec(`
        CREATE TABLE vocabulary_quiz_answers (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            study_session_id INTEGER NOT NULL,
            word_id INTEGER NOT NULL,
            answer TEXT NOT NULL,
            correct BOOLEAN NOT NULL,
            created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (study_session_id) REFERENCES study_sessions(id),
            FOREIGN KEY (word_id) REFERENCES words(id)
        )`); err != nil {
		return fmt.Errorf("error creating vocabulary_quiz_answers: %v", err)
	}

	// Create indices
	if _, err := db.Exec("CREATE INDEX idx_study_activities_type ON study_activities(activity_type)"); err != nil {
		return fmt.Errorf("error creating activity_type index: %v", err)
	}
	if _, err := db.Exec("CREATE INDEX idx_study_activities_completed ON study_activities(completed_at)"); err != nil {
		return fmt.Errorf("error creating completed_at index: %v", err)
	}
	if _, err := db.Exec("CREATE INDEX idx_study_activities_group_id ON study_activities(group_id)"); err != nil {
		return fmt.Errorf("error creating group_id index: %v", err)
	}
	if _, err := db.Exec("CREATE INDEX idx_study_sessions_activity_id ON study_sessions(study_activity_id)"); err != nil {
		return fmt.Errorf("error creating study_activity_id index: %v", err)
	}
	if _, err := db.Exec("CREATE INDEX idx_quiz_answers_session_id ON vocabulary_quiz_answers(study_session_id)"); err != nil {
		return fmt.Errorf("error creating study_session_id index: %v", err)
	}

	fmt.Println("Migration 2 completed successfully")
	return nil
}

func seedData(db *sql.DB) error {
	// First check if we already have data
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM groups").Scan(&count)
	if err != nil {
		return fmt.Errorf("error checking existing data: %v", err)
	}

	// Skip seeding if we already have data
	if count > 0 {
		return nil
	}

	seeds := []string{
		"basic_greetings.json",
	}

	for _, seed := range seeds {
		if err := loadSeedFile(db, seed); err != nil {
			return fmt.Errorf("error loading seed file %s: %v", seed, err)
		}
	}
	return nil
} 