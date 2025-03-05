package handlers

import (
	"log"
	"github.com/gin-gonic/gin"
	"github.com/mohawa/lang-portal/backend_go/internal/database"
)

func GetQuickStats(c *gin.Context) {
	log.Printf("Getting dashboard quick stats...")

	// Get total words
	var totalWords int
	err := database.DB.QueryRow("SELECT COUNT(*) FROM words").Scan(&totalWords)
	if err != nil {
		log.Printf("Error getting total words: %v", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Get words studied (words with reviews)
	var wordsStudied int
	err = database.DB.QueryRow("SELECT COUNT(DISTINCT word_id) FROM word_review_items").Scan(&wordsStudied)
	if err != nil {
		log.Printf("Error getting words studied: %v", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Get study sessions count
	var studySessions int
	err = database.DB.QueryRow("SELECT COUNT(*) FROM study_sessions").Scan(&studySessions)
	if err != nil {
		log.Printf("Error getting study sessions: %v", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Calculate accuracy rate
	var correctCount, totalCount int
	err = database.DB.QueryRow(`
		SELECT 
			COUNT(CASE WHEN correct = true THEN 1 END),
			COUNT(*)
		FROM word_review_items
	`).Scan(&correctCount, &totalCount)
	if err != nil {
		log.Printf("Error calculating accuracy rate: %v", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var accuracyRate float64
	if totalCount > 0 {
		accuracyRate = float64(correctCount) * 100 / float64(totalCount)
	}

	log.Printf("Returning dashboard stats: words=%d, studied=%d, sessions=%d, accuracy=%.2f%%",
		totalWords, wordsStudied, studySessions, accuracyRate)

	c.JSON(200, gin.H{
		"total_words":    totalWords,
		"words_studied":  wordsStudied,
		"study_sessions": studySessions,
		"accuracy_rate":  accuracyRate,
	})
}

func GetStudyProgress(c *gin.Context) {
	log.Printf("Getting study progress...")

	// Get total words studied
	var totalWordsStudied int
	err := database.DB.QueryRow(`
		SELECT COUNT(DISTINCT word_id) 
		FROM word_review_items
	`).Scan(&totalWordsStudied)

	if err != nil {
		log.Printf("Error getting total words studied: %v", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"total_words_studied": totalWordsStudied,
	})
} 