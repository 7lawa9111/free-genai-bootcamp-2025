package handlers

import (
	"log"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/mohawa/lang-portal/backend_go/internal/database"
)

func GetWords(c *gin.Context) {
	log.Printf("Getting words from database...")
	
	// First, verify database connection
	if database.DB == nil {
		log.Printf("Error: database.DB is nil")
		c.JSON(500, gin.H{"error": "Database connection not initialized"})
		return
	}

	// Query all words
	rows, err := database.DB.Query("SELECT id, japanese, romaji, english FROM words")
	if err != nil {
		log.Printf("Error querying words: %v", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	// Initialize empty slice for words
	words := make([]gin.H, 0)

	// Scan rows into words slice
	for rows.Next() {
		var id int
		var japanese, romaji, english string
		if err := rows.Scan(&id, &japanese, &romaji, &english); err != nil {
			log.Printf("Error scanning word: %v", err)
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		log.Printf("Found word: %d, %s (%s) - %s", id, japanese, romaji, english)
		words = append(words, gin.H{
			"id":       id,
			"japanese": japanese,
			"romaji":   romaji,
			"english":  english,
		})
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		log.Printf("Error iterating over words: %v", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Returning %d words", len(words))
	c.JSON(200, gin.H{
		"items": words,
		"pagination": gin.H{
			"current_page":   1,
			"total_pages":    1,
			"total_items":    len(words),
			"items_per_page": 100,
		},
	})
}

func GetWord(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("Invalid ID format: %v", err)
		c.JSON(400, gin.H{"error": "Invalid ID format"})
		return
	}

	log.Printf("Getting word with ID: %d", id)

	var word struct {
		Japanese string
		Romaji   string
		English  string
	}

	err = database.DB.QueryRow(
		"SELECT japanese, romaji, english FROM words WHERE id = ?", 
		id,
	).Scan(&word.Japanese, &word.Romaji, &word.English)

	if err != nil {
		log.Printf("Word not found with ID: %d", id)
		c.JSON(404, gin.H{"error": "Word not found"})
		return
	}

	log.Printf("Found word %d: %s (%s) - %s", id, word.Japanese, word.Romaji, word.English)
	c.JSON(200, gin.H{
		"id":       id,
		"japanese": word.Japanese,
		"romaji":   word.Romaji,
		"english":  word.English,
		"stats": gin.H{
			"correct_count": 0,
			"wrong_count":   0,
		},
		"groups": []interface{}{},
	})
} 