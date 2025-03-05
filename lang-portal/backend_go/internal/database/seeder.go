package database

import (
	"database/sql"
	"fmt"
	"log"
)

func SeedProductionData(db *sql.DB) error {
	log.Printf("Seeding production database...")

	productionData := []string{
		// Groups
		`INSERT INTO groups (id, name) VALUES (1, 'Basic Greetings');`,
		`INSERT INTO groups (id, name) VALUES (2, 'Numbers');`,
		
		// Words
		`INSERT INTO words (id, japanese, romaji, english) VALUES (1, 'こんにちは', 'konnichiwa', 'hello');`,
		`INSERT INTO words (id, japanese, romaji, english) VALUES (2, 'さようなら', 'sayounara', 'goodbye');`,
		`INSERT INTO words (id, japanese, romaji, english) VALUES (3, 'おはよう', 'ohayou', 'good morning');`,
		
		// Study Activities
		`INSERT INTO study_activities (id, name, thumbnail_url, description) 
		 VALUES (1, 'Flashcards', '/images/flashcards.png', 'Practice words using flashcards');`,
	}

	for _, query := range productionData {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("error seeding production data: %v", err)
		}
	}

	log.Printf("Production database seeded successfully")
	return nil
} 