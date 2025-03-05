package handlers

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetGroup(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID format"})
		return
	}

	// For testing, only return data for group ID 1
	if id != 1 {
		c.JSON(404, gin.H{"error": "Group not found"})
		return
	}

	c.JSON(200, gin.H{
		"id":   id,
		"name": "Basic Greetings",
		"stats": gin.H{
			"total_word_count": 5,
		},
	})
}

func GetGroups(c *gin.Context) {
	response := EmptyPaginatedResponse()
	c.JSON(200, response)
}

func GetGroupWords(c *gin.Context) {
	response := EmptyPaginatedResponse()
	c.JSON(200, response)
}

func GetGroupStudySessions(c *gin.Context) {
	response := EmptyPaginatedResponse()
	c.JSON(200, response)
} 