package middleware

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"lang-portal/backend_go/internal/validation"
)

// ValidatePagination validates page and limit query parameters
func ValidatePagination() gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", strconv.Itoa(validation.DefaultPageSize)))

		validatedPage, validatedLimit, err := validation.ValidatePagination(page, limit)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// Store validated values in context
		c.Set("page", validatedPage)
		c.Set("limit", validatedLimit)
		c.Next()
	}
}

// ValidateID validates ID parameter in URL
func ValidateID(paramName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param(paramName), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
			c.Abort()
			return
		}

		if err := validation.ValidateID(id); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// Store validated ID in context
		c.Set(paramName, id)
		c.Next()
	}
}

// ValidateJSONRequest validates JSON request body
func ValidateJSONRequest(model interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := c.ShouldBindJSON(model); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			c.Abort()
			return
		}
		c.Next()
	}
} 