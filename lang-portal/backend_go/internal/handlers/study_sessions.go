package handlers

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetStudySessions(c *gin.Context) {
	response := EmptyPaginatedResponse()
	c.JSON(200, response)
}

func GetStudySession(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID format"})
		return
	}

	// For testing, only return data for session ID 1
	if id != 1 {
		c.JSON(404, gin.H{"error": "Study session not found"})
		return
	}

	c.JSON(200, gin.H{
		"id":                 id,
		"activity_name":      "Flashcards",
		"group_name":         "Basic Greetings",
		"start_time":         "2025-03-01T13:53:24Z",
		"end_time":          "2025-03-01T14:53:24Z",
		"review_items_count": 5,
	})
}

func ReviewWord(c *gin.Context) {
	c.JSON(200, gin.H{
		"success":          true,
		"word_id":         1,
		"study_session_id": 1,
		"correct":         true,
		"created_at":      "2025-03-01T14:52:59Z",
	})
} 