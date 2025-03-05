package handlers

import (
	"log"
	"github.com/gin-gonic/gin"
	"github.com/mohawa/lang-portal/backend_go/internal/database"
)

func InitTestData(c *gin.Context) {
	log.Printf("Initializing test data...")

	if err := database.InsertTestData(database.DB); err != nil {
		log.Printf("Error initializing test data: %v", err)
		c.JSON(500, gin.H{"error": "Failed to initialize test data"})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "Test data initialized",
	})
} 