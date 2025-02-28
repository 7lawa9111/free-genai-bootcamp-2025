//go:build mage
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type Word struct {
	Japanese string          `json:"japanese"`
	Romaji   string          `json:"romaji"`
	English  string          `json:"english"`
	Parts    json.RawMessage `json:"parts"`
}

// InitDB initializes the SQLite database
func InitDB() error {
	fmt.Println("Initializing database...")
	db, err := sql.Open("sqlite3", "./words.db")
	if err != nil {
		return err
	}
	defer db.Close()
	return nil
}

// Migrate runs all pending migrations
func Migrate() error {
	fmt.Println("Running migrations...")
	db, err := sql.Open("sqlite3", "./words.db")
	if err != nil {
		return err
	}
	defer db.Close()

	files, err := filepath.Glob("db/migrations/*.sql")
	if err != nil {
		return err
	}

	sort.Strings(files)

	for _, file := range files {
		fmt.Printf("Applying migration: %s\n", file)
		content, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		statements := strings.Split(string(content), ";")
		for _, statement := range statements {
			if strings.TrimSpace(statement) != "" {
				_, err = db.Exec(statement)
				if err != nil {
					return fmt.Errorf("error executing %s: %v", file, err)
				}
			}
		}
	}

	return nil
}

// Seed imports data from JSON files
func Seed() error {
	fmt.Println("Seeding database...")
	db, err := sql.Open("sqlite3", "./words.db")
	if err != nil {
		return err
	}
	defer db.Close()

	files, err := filepath.Glob("db/seeds/*.json")
	if err != nil {
		return err
	}

	for _, file := range files {
		groupName := strings.TrimSuffix(filepath.Base(file), ".json")
		fmt.Printf("Seeding %s...\n", groupName)

		// Create group
		result, err := db.Exec("INSERT INTO groups (name) VALUES (?)", groupName)
		if err != nil {
			return err
		}

		groupID, err := result.LastInsertId()
		if err != nil {
			return err
		}

		// Read and parse JSON file
		content, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		var words []Word
		if err := json.Unmarshal(content, &words); err != nil {
			return err
		}

		// Insert words and create word-group associations
		for _, word := range words {
			result, err := db.Exec(
				"INSERT INTO words (japanese, romaji, english, parts) VALUES (?, ?, ?, ?)",
				word.Japanese, word.Romaji, word.English, word.Parts,
			)
			if err != nil {
				return err
			}

			wordID, err := result.LastInsertId()
			if err != nil {
				return err
			}

			_, err = db.Exec(
				"INSERT INTO words_groups (word_id, group_id) VALUES (?, ?)",
				wordID, groupID,
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
} 