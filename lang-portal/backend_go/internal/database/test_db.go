package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func InitTestDB() (*sql.DB, error) {
	log.Printf("Initializing test database...")
	
	// Get project root directory (backend_go)
	pwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %v", err)
	}
	
	projectRoot := pwd
	for filepath.Base(projectRoot) != "backend_go" && projectRoot != "/" {
		projectRoot = filepath.Dir(projectRoot)
	}
	if projectRoot == "/" {
		return nil, fmt.Errorf("could not find backend_go directory")
	}

	dbPath := filepath.Join(projectRoot, "db", "data", "words.test.db")
	log.Printf("Opening test database at: %s", dbPath)

	// Create database directory if it doesn't exist
	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %v", err)
	}

	// Remove existing test database
	os.Remove(dbPath)

	// Create new database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Read and execute schema
	schemaPath := filepath.Join(projectRoot, "db", "migrations", "0001_init.sql")
	log.Printf("Reading schema from: %s", schemaPath)
	schemaSQL, err := os.ReadFile(schemaPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read schema file: %v", err)
	}

	_, err = db.Exec(string(schemaSQL))
	if err != nil {
		return nil, fmt.Errorf("failed to execute schema: %v", err)
	}

	// Insert test data
	if err := InsertTestData(db); err != nil {
		return nil, fmt.Errorf("failed to insert test data: %v", err)
	}

	log.Printf("Test database initialized successfully")
	DB = db
	return db, nil
}

// Make InsertTestData public so it can be called after reset
func InsertTestData(db *sql.DB) error {
	testData := []string{
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
	}

	for _, query := range testData {
		log.Printf("Executing: %s", query)
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("error executing test data query: %v", err)
		}
	}

	return nil
} 