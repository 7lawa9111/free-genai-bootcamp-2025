//go:build mage
package main

import (
	"fmt"
	"os"
	"github.com/magefile/mage/sh"
)

// Dbinit initializes the database with schema
func Dbinit() error {
	fmt.Println("Initializing production database...")
	
	// Create db directory if it doesn't exist
	if err := os.MkdirAll("db/data", 0755); err != nil {
		return fmt.Errorf("failed to create db directory: %v", err)
	}

	// Execute schema file
	if err := sh.Run("sqlite3", "db/data/words.prod.db", ".read db/migrations/0001_init.sql"); err != nil {
		return fmt.Errorf("failed to initialize database: %v", err)
	}

	fmt.Println("Production database initialized successfully")
	return nil
}

// Seed adds test data to the database
func Seed() error {
	fmt.Println("Seeding production database...")

	// Test data SQL
	seedSQL := []string{
		// Groups
		`INSERT INTO groups (id, name) VALUES (1, 'Basic Greetings');`,
		`INSERT INTO groups (id, name) VALUES (2, 'Numbers');`,
		
		// Words
		`INSERT INTO words (id, japanese, romaji, english) VALUES (1, 'こんにちは', 'konnichiwa', 'hello');`,
		`INSERT INTO words (id, japanese, romaji, english) VALUES (2, 'さようなら', 'sayounara', 'goodbye');`,
		`INSERT INTO words (id, japanese, romaji, english) VALUES (3, 'おはよう', 'ohayou', 'good morning');`,
		`INSERT INTO words (id, japanese, romaji, english) VALUES (4, '一', 'ichi', 'one');`,
		`INSERT INTO words (id, japanese, romaji, english) VALUES (5, '二', 'ni', 'two');`,

		// Study Activities
		`INSERT INTO study_activities (id, name, thumbnail_url, description) 
		 VALUES (1, 'Flashcards', '/images/flashcards.png', 'Practice words using flashcards');`,

		// Study Sessions
		`INSERT INTO study_sessions (id, group_id, created_at) 
		 VALUES (1, 1, datetime('now', '-1 day'));`,
		`INSERT INTO study_sessions (id, group_id, created_at) 
		 VALUES (2, 2, datetime('now', '-2 hours'));`,

		// Word Review Items
		`INSERT INTO word_review_items (word_id, study_session_id, correct, created_at) 
		 VALUES (1, 1, true, datetime('now', '-23.5 hours'));`,
		`INSERT INTO word_review_items (word_id, study_session_id, correct, created_at) 
		 VALUES (2, 1, false, datetime('now', '-23.4 hours'));`,
		`INSERT INTO word_review_items (word_id, study_session_id, correct, created_at) 
		 VALUES (3, 1, true, datetime('now', '-23.3 hours'));`,
		`INSERT INTO word_review_items (word_id, study_session_id, correct, created_at) 
		 VALUES (4, 2, true, datetime('now', '-1.5 hours'));`,
		`INSERT INTO word_review_items (word_id, study_session_id, correct, created_at) 
		 VALUES (5, 2, false, datetime('now', '-1 hour'));`,

		// Link words to groups
		`INSERT INTO words_groups (word_id, group_id) VALUES (1, 1);`,
		`INSERT INTO words_groups (word_id, group_id) VALUES (2, 1);`,
		`INSERT INTO words_groups (word_id, group_id) VALUES (3, 1);`,
		`INSERT INTO words_groups (word_id, group_id) VALUES (4, 2);`,
		`INSERT INTO words_groups (word_id, group_id) VALUES (5, 2);`,
	}

	// Write seed SQL to temporary file
	tmpFile := "db/seed_temp.sql"
	f, err := os.Create(tmpFile)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile)

	for _, sql := range seedSQL {
		if _, err := f.WriteString(sql + "\n"); err != nil {
			return fmt.Errorf("failed to write seed SQL: %v", err)
		}
	}
	f.Close()

	// Execute seed file
	if err := sh.Run("sqlite3", "db/data/words.prod.db", ".read " + tmpFile); err != nil {
		return fmt.Errorf("failed to seed database: %v", err)
	}

	fmt.Println("Production database seeded successfully")
	return nil
}

// Testdb initializes and seeds test database
func Testdb() error {
	fmt.Println("Setting up test database...")
	
	os.Setenv("APP_ENV", "test")
	
	// Create test db directory
	if err := os.MkdirAll("db/data", 0755); err != nil {
		return fmt.Errorf("failed to create test db directory: %v", err)
	}

	testDB := "db/data/words.test.db"
	
	// Remove existing test db
	os.Remove(testDB)

	// Initialize schema
	if err := sh.Run("sqlite3", testDB, ".read db/migrations/0001_init.sql"); err != nil {
		return fmt.Errorf("failed to initialize test database: %v", err)
	}

	fmt.Println("Test database initialized, seeding data...")

	// Test data SQL
	seedSQL := []string{
		// Groups
		`INSERT INTO groups (id, name) VALUES (1, 'Basic Greetings');`,
		`INSERT INTO groups (id, name) VALUES (2, 'Numbers');`,
		
		// Words
		`INSERT INTO words (id, japanese, romaji, english) VALUES (1, 'こんにちは', 'konnichiwa', 'hello');`,
		`INSERT INTO words (id, japanese, romaji, english) VALUES (2, 'さようなら', 'sayounara', 'goodbye');`,
		`INSERT INTO words (id, japanese, romaji, english) VALUES (3, 'おはよう', 'ohayou', 'good morning');`,
		`INSERT INTO words (id, japanese, romaji, english) VALUES (4, '一', 'ichi', 'one');`,
		`INSERT INTO words (id, japanese, romaji, english) VALUES (5, '二', 'ni', 'two');`,

		// Study Activities
		`INSERT INTO study_activities (id, name, thumbnail_url, description) 
		 VALUES (1, 'Flashcards', '/images/flashcards.png', 'Practice words using flashcards');`,

		// Study Sessions
		`INSERT INTO study_sessions (id, group_id, created_at) 
		 VALUES (1, 1, datetime('now', '-1 day'));`,
		`INSERT INTO study_sessions (id, group_id, created_at) 
		 VALUES (2, 2, datetime('now', '-2 hours'));`,

		// Word Review Items
		`INSERT INTO word_review_items (word_id, study_session_id, correct, created_at) 
		 VALUES (1, 1, true, datetime('now', '-23.5 hours'));`,
		`INSERT INTO word_review_items (word_id, study_session_id, correct, created_at) 
		 VALUES (2, 1, false, datetime('now', '-23.4 hours'));`,
		`INSERT INTO word_review_items (word_id, study_session_id, correct, created_at) 
		 VALUES (3, 1, true, datetime('now', '-23.3 hours'));`,
		`INSERT INTO word_review_items (word_id, study_session_id, correct, created_at) 
		 VALUES (4, 2, true, datetime('now', '-1.5 hours'));`,
		`INSERT INTO word_review_items (word_id, study_session_id, correct, created_at) 
		 VALUES (5, 2, false, datetime('now', '-1 hour'));`,

		// Link words to groups
		`INSERT INTO words_groups (word_id, group_id) VALUES (1, 1);`,
		`INSERT INTO words_groups (word_id, group_id) VALUES (2, 1);`,
		`INSERT INTO words_groups (word_id, group_id) VALUES (3, 1);`,
		`INSERT INTO words_groups (word_id, group_id) VALUES (4, 2);`,
		`INSERT INTO words_groups (word_id, group_id) VALUES (5, 2);`,
	}

	// Write seed SQL to temporary file
	tmpFile := "db/seed_temp.sql"
	f, err := os.Create(tmpFile)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile)

	for _, sql := range seedSQL {
		if _, err := f.WriteString(sql + "\n"); err != nil {
			return fmt.Errorf("failed to write seed SQL: %v", err)
		}
	}
	f.Close()

	// Execute seed file
	if err := sh.Run("sqlite3", testDB, ".read " + tmpFile); err != nil {
		return fmt.Errorf("failed to seed test database: %v", err)
	}

	fmt.Println("Test database setup complete with seed data")
	return nil
} 