package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func readMigrationFile(filename string) string {
	path := filepath.Join("db", "migrations", filename)
	content, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("Error reading migration file %s: %v", filename, err))
	}
	return string(content)
}

func loadSeedFile(db *sql.DB, filename string) error {
	path := filepath.Join("db", "seeds", filename)
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error reading seed file %s: %v", filename, err)
	}

	var seedData struct {
		GroupName string `json:"group_name"`
		Words []struct {
			Japanese string `json:"japanese"`
			Romaji string `json:"romaji"`
			English string `json:"english"`
			Parts interface{} `json:"parts"`
		} `json:"words"`
	}

	if err := json.Unmarshal(content, &seedData); err != nil {
		return fmt.Errorf("error parsing seed file %s: %v", filename, err)
	}

	// Insert data in a transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert group
	result, err := tx.Exec(`INSERT INTO groups (name) VALUES (?)`, seedData.GroupName)
	if err != nil {
		return err
	}

	groupID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	// Insert words and word-group relationships
	for _, word := range seedData.Words {
		result, err := tx.Exec(`
			INSERT INTO words (japanese, romaji, english, parts)
			VALUES (?, ?, ?, ?)`,
			word.Japanese, word.Romaji, word.English, word.Parts)
		if err != nil {
			return err
		}

		wordID, err := result.LastInsertId()
		if err != nil {
			return err
		}

		_, err = tx.Exec(`
			INSERT INTO words_groups (word_id, group_id)
			VALUES (?, ?)`,
			wordID, groupID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
} 