package middleware

import (
	"github.com/gin-gonic/gin"
)

type CORSConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge          int
}

var DefaultCORSConfig = CORSConfig{
	AllowOrigins:     []string{"*"},
	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	AllowHeaders:     []string{"Content-Type", "Content-Length", "Accept-Encoding", "Authorization"},
	ExposeHeaders:    []string{},
	AllowCredentials: true,
	MaxAge:          86400, // 24 hours
}

// CORSWithConfig creates a CORS middleware with custom configuration
func CORSWithConfig(config CORSConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		
		// Check if origin is allowed
		allowOrigin := "*"
		for _, o := range config.AllowOrigins {
			if o == origin || o == "*" {
				allowOrigin = o
				break
			}
		}

		// Set CORS headers
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Methods", joinStrings(config.AllowMethods))
		c.Writer.Header().Set("Access-Control-Allow-Headers", joinStrings(config.AllowHeaders))
		
		if len(config.ExposeHeaders) > 0 {
			c.Writer.Header().Set("Access-Control-Expose-Headers", joinStrings(config.ExposeHeaders))
		}
		
		if config.AllowCredentials {
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		}
		
		if config.MaxAge > 0 {
			c.Writer.Header().Set("Access-Control-Max-Age", string(config.MaxAge))
		}

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// Helper function to join strings with comma
func joinStrings(strings []string) string {
	if len(strings) == 0 {
		return ""
	}
	result := strings[0]
	for _, s := range strings[1:] {
		result += ", " + s
	}
	return result
}

// CORS returns the default CORS middleware
func CORS() gin.HandlerFunc {
	return CORSWithConfig(DefaultCORSConfig)
} 