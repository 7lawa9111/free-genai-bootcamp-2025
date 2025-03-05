package main

import (
	"log"
	"os"
	"github.com/gin-gonic/gin"
	"github.com/mohawa/lang-portal/backend_go/internal/database"
	"github.com/mohawa/lang-portal/backend_go/internal/handlers"
)

func main() {
	// Set environment
	if os.Getenv("APP_ENV") == "" {
		os.Setenv("APP_ENV", "development")
	}

	// Initialize database
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// API routes
	api := r.Group("/api")
	{
		// Test routes (only in test environment)
		if os.Getenv("APP_ENV") == "test" {
			api.POST("/test/init_data", handlers.InitTestData)
		}

		// Words routes
		api.GET("/words", handlers.GetWords)
		api.GET("/words/:id", handlers.GetWord)

		// Groups routes
		api.GET("/groups", handlers.GetGroups)
		api.GET("/groups/:id", handlers.GetGroup)
		api.GET("/groups/:id/words", handlers.GetGroupWords)
		api.GET("/groups/:id/study_sessions", handlers.GetGroupStudySessions)

		// Study sessions routes
		api.GET("/study_sessions", handlers.GetStudySessions)
		api.GET("/study_sessions/:id", handlers.GetStudySession)
		api.POST("/study_sessions/:id/words/:word_id/review", handlers.ReviewWord)

		// Study activities routes
		api.GET("/study_activities/:id", handlers.GetStudyActivity)
		api.GET("/study_activities/:id/study_sessions", handlers.GetStudyActivitySessions)
		api.POST("/study_activities", handlers.CreateStudyActivity)

		// Dashboard routes
		api.GET("/dashboard/quick-stats", handlers.GetQuickStats)
		api.GET("/dashboard/study_progress", handlers.GetStudyProgress)

		// Reset routes
		api.POST("/reset_history", handlers.ResetHistory)
		api.POST("/full_reset", handlers.FullReset)
	}

	log.Printf("Server starting on http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
} 