package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"lang-portal/backend_go/internal/errors"
)

// ErrorHandler handles common errors and converts them to appropriate HTTP responses
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			switch err {
			case errors.ErrNotFound, errors.ErrGroupNotFound, errors.ErrWordNotFound,
				errors.ErrStudySessionNotFound, errors.ErrActivityNotFound:
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			case errors.ErrInvalidInput, errors.ErrInvalidPageSize:
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			case errors.ErrDatabaseError:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			}
			return
		}
	}
} 