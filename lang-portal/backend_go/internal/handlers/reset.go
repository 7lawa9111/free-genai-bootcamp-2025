package handlers

import (
	"log"
	"github.com/gin-gonic/gin"
	"github.com/mohawa/lang-portal/backend_go/internal/database"
)

func ResetHistory(c *gin.Context) {
	log.Printf("ResetHistory called with method: %s", c.Request.Method)

	// First verify database connection
	if database.DB == nil {
		log.Printf("Error: database.DB is nil")
		c.JSON(500, gin.H{"error": "Database connection not initialized"})
		return
	}

	// Begin transaction
	tx, err := database.DB.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		c.JSON(500, gin.H{"error": "Failed to start reset operation"})
		return
	}

	// Delete all word review items
	result, err := tx.Exec("DELETE FROM word_review_items")
	if err != nil {
		tx.Rollback()
		log.Printf("Error deleting word review items: %v", err)
		c.JSON(500, gin.H{"error": "Failed to reset study history"})
		return
	}
	reviewsDeleted, _ := result.RowsAffected()

	// Delete all study sessions
	result, err = tx.Exec("DELETE FROM study_sessions")
	if err != nil {
		tx.Rollback()
		log.Printf("Error deleting study sessions: %v", err)
		c.JSON(500, gin.H{"error": "Failed to reset study history"})
		return
	}
	sessionsDeleted, _ := result.RowsAffected()

	// Commit transaction
	if err = tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		c.JSON(500, gin.H{"error": "Failed to complete reset"})
		return
	}

	log.Printf("Study history reset successful. Deleted %d review items and %d study sessions", 
		reviewsDeleted, sessionsDeleted)
	
	c.JSON(200, gin.H{
		"success": true,
		"message": "Study history has been reset",
	})
}

func FullReset(c *gin.Context) {
	// ... FullReset remains the same ...
} 