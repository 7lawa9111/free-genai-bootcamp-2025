package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/mohawa/lang-portal/backend_go/internal/services"
)

type ResetHandler struct {
	resetService *services.ResetService
}

func NewResetHandler() *ResetHandler {
	return &ResetHandler{
		resetService: services.NewResetService(),
	}
}

// ResetHistory handles POST /api/reset_history
func ResetHistory(c *gin.Context) {
	handler := NewResetHandler()

	err := handler.resetService.ResetHistory()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Study history has been reset",
	})
}

// FullReset handles POST /api/full_reset
func FullReset(c *gin.Context) {
	handler := NewResetHandler()

	err := handler.resetService.FullReset()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "System has been fully reset",
	})
} 