//go:build mage
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

const dbPath = "words.db"

// SeedData represents the structure of our seed JSON files
type SeedData struct {
	GroupName string      `json:"group_name"`
	Words     []SeedWord `json:"words"`
}

type SeedWord struct {
	Japanese string `json:"japanese"`
	Romaji   string `json:"romaji"`
	English  string `json:"english"`
	Parts    string `json:"parts"`
}

// InitDB initializes the SQLite database
func InitDB() error {
	if _, err := os.Stat(dbPath); err == nil {
		fmt.Printf("Database %s already exists\n", dbPath)
		return nil
	}

	file, err := os.Create(dbPath)
	if err != nil {
		return fmt.Errorf("error creating database file: %v", err)
	}
	file.Close()

	fmt.Printf("Created database file: %s\n", dbPath)
	return nil
}

// Migrate runs all pending migrations
func Migrate() error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}
	defer db.Close()

	// Create migrations table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("error creating migrations table: %v", err)
	}

	// Get list of applied migrations
	rows, err := db.Query("SELECT name FROM migrations")
	if err != nil {
		return fmt.Errorf("error querying migrations: %v", err)
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return fmt.Errorf("error scanning migration row: %v", err)
		}
		applied[name] = true
	}

	// Get list of migration files
	files, err := filepath.Glob("db/migrations/*.sql")
	if err != nil {
		return fmt.Errorf("error listing migration files: %v", err)
	}
	sort.Strings(files)

	// Apply pending migrations
	for _, file := range files {
		name := filepath.Base(file)
		if applied[name] {
			continue
		}

		fmt.Printf("Applying migration: %s\n", name)

		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("error reading migration file: %v", err)
		}

		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("error starting transaction: %v", err)
		}

		// Split the file into individual statements
		statements := strings.Split(string(content), ";")
		for _, stmt := range statements {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}

			if _, err := tx.Exec(stmt); err != nil {
				tx.Rollback()
				return fmt.Errorf("error executing migration: %v", err)
			}
		}

		// Record the migration
		if _, err := tx.Exec("INSERT INTO migrations (name) VALUES (?)", name); err != nil {
			tx.Rollback()
			return fmt.Errorf("error recording migration: %v", err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("error committing migration: %v", err)
		}

		fmt.Printf("Successfully applied migration: %s\n", name)
	}

	return nil
}

// Seed runs all seeding operations
func Seed() error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}
	defer db.Close()

	// Get list of seed files
	files, err := filepath.Glob("db/seeds/*.json")
	if err != nil {
		return fmt.Errorf("error listing seed files: %v", err)
	}

	for _, file := range files {
		fmt.Printf("Processing seed file: %s\n", filepath.Base(file))

		if err := processSeedFile(db, file); err != nil {
			return fmt.Errorf("error processing seed file %s: %v", file, err)
		}
	}

	return nil
}

func processSeedFile(db *sql.DB, filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading seed file: %v", err)
	}

	var seedData SeedData
	if err := json.Unmarshal(content, &seedData); err != nil {
		return fmt.Errorf("error parsing seed file: %v", err)
	}

	// Start transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback()

	// Create group
	var groupID int64
	result := tx.QueryRow("SELECT id FROM groups WHERE name = ?", seedData.GroupName)
	if err := result.Scan(&groupID); err == sql.ErrNoRows {
		res, err := tx.Exec("INSERT INTO groups (name) VALUES (?)", seedData.GroupName)
		if err != nil {
			return fmt.Errorf("error creating group: %v", err)
		}
		groupID, err = res.LastInsertId()
		if err != nil {
			return fmt.Errorf("error getting group ID: %v", err)
		}
	} else if err != nil {
		return fmt.Errorf("error checking group existence: %v", err)
	}

	// Insert words and create word-group associations
	for _, word := range seedData.Words {
		// Check if word already exists
		var wordID int64
		result := tx.QueryRow("SELECT id FROM words WHERE japanese = ? AND romaji = ?", 
			word.Japanese, word.Romaji)
		if err := result.Scan(&wordID); err == sql.ErrNoRows {
			// Create new word
			res, err := tx.Exec(`
				INSERT INTO words (japanese, romaji, english, parts) 
				VALUES (?, ?, ?, ?)`,
				word.Japanese, word.Romaji, word.English, word.Parts)
			if err != nil {
				return fmt.Errorf("error creating word: %v", err)
			}
			wordID, err = res.LastInsertId()
			if err != nil {
				return fmt.Errorf("error getting word ID: %v", err)
			}
		} else if err != nil {
			return fmt.Errorf("error checking word existence: %v", err)
		}

		// Create word-group association if it doesn't exist
		_, err = tx.Exec(`
			INSERT INTO words_groups (word_id, group_id)
			SELECT ?, ?
			WHERE NOT EXISTS (
				SELECT 1 FROM words_groups 
				WHERE word_id = ? AND group_id = ?
			)`,
			wordID, groupID, wordID, groupID)
		if err != nil {
			return fmt.Errorf("error creating word-group association: %v", err)
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	fmt.Printf("Successfully seeded group '%s' with %d words\n", 
		seedData.GroupName, len(seedData.Words))
	return nil
}

// ResetDB drops and recreates all tables, then runs migrations and seeds
func ResetDB() error {
	// Remove existing database
	if err := os.Remove(dbPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("error removing database: %v", err)
	}

	// Initialize new database
	if err := InitDB(); err != nil {
		return err
	}

	// Run migrations
	if err := Migrate(); err != nil {
		return err
	}

	// Run seeds
	if err := Seed(); err != nil {
		return err
	}

	return nil
} 