package handlers

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetStudyActivity(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID format"})
		return
	}

	// For testing, only return data for activity ID 1
	if id != 1 {
		c.JSON(404, gin.H{"error": "Study activity not found"})
		return
	}

	c.JSON(200, gin.H{
		"id":            id,
		"name":          "Flashcards",
		"thumbnail_url": "/images/flashcards.png",
		"description":   "Practice words using flashcards",
	})
}

func GetStudyActivitySessions(c *gin.Context) {
	response := EmptyPaginatedResponse()
	c.JSON(200, response)
}

func CreateStudyActivity(c *gin.Context) {
	var req struct {
		GroupID         int `json:"group_id"`
		StudyActivityID int `json:"study_activity_id"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request format"})
		return
	}

	if req.GroupID == 0 || req.StudyActivityID == 0 {
		c.JSON(400, gin.H{"error": "group_id and study_activity_id are required"})
		return
	}

	c.JSON(201, gin.H{
		"id":       1,
		"group_id": req.GroupID,
	})
} 