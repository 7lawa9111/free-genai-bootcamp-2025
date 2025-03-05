package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	_ "github.com/mattn/go-sqlite3"
)

var (
	DB *sql.DB
)

func InitDB() (*sql.DB, error) {
	// If in test mode, use test database
	if os.Getenv("APP_ENV") == "test" {
		log.Printf("Using test database")
		return InitTestDB()
	}

	// Get project root directory (backend_go)
	pwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %v", err)
	}
	
	// Find backend_go directory
	projectRoot := pwd
	for filepath.Base(projectRoot) != "backend_go" && projectRoot != "/" {
		projectRoot = filepath.Dir(projectRoot)
	}
	if projectRoot == "/" {
		return nil, fmt.Errorf("could not find backend_go directory")
	}
	log.Printf("Project root directory: %s", projectRoot)

	// Get database path based on environment
	dbPath := filepath.Join(projectRoot, getDatabasePath())
	log.Printf("Opening database at: %s", dbPath)

	// Create database directory if it doesn't exist
	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %v", err)
	}

	// Check if database exists, if not, create it and initialize schema
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		log.Printf("Database does not exist, creating new database")
		db, err := createNewDatabase(dbPath, projectRoot)
		if err != nil {
			return nil, fmt.Errorf("failed to create database: %v", err)
		}
		return db, nil
	}

	// Open existing database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	DB = db
	return db, nil
}

func getDatabasePath() string {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	baseDir := "db/data"
	switch env {
	case "test":
		return filepath.Join(baseDir, "words.test.db")
	case "production":
		return filepath.Join(baseDir, "words.prod.db")
	default:
		return filepath.Join(baseDir, "words.dev.db")
	}
}

func createNewDatabase(dbPath string, projectRoot string) (*sql.DB, error) {
	// Create empty database file
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
	log.Printf("Schema executed successfully")

	// If it's test environment, insert test data
	if os.Getenv("APP_ENV") == "test" {
		log.Printf("Inserting test data...")
		if err := insertTestData(db); err != nil {
			return nil, fmt.Errorf("failed to insert test data: %v", err)
		}
		log.Printf("Test data inserted successfully")
	} else if os.Getenv("APP_ENV") == "production" && os.Getenv("SEED_DB") == "true" {
		log.Printf("Seeding production database...")
		if err := SeedProductionData(db); err != nil {
			return nil, fmt.Errorf("failed to seed production data: %v", err)
		}
		log.Printf("Production database seeded successfully")
	}

	DB = db
	return db, nil
}

func insertTestData(db *sql.DB) error {
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
	}

	for _, query := range testData {
		log.Printf("Executing query: %s", query)
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("error executing test data query: %v", err)
		}
	}

	return nil
} 