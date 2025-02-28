package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mohawa/lang-portal/backend_go/internal/database"
	"github.com/mohawa/lang-portal/backend_go/internal/handlers"
)

func main() {
	// Initialize DB
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize Gin
	r := gin.Default()

	// API routes
	api := r.Group("/api")
	{
		// Dashboard routes
		api.GET("/dashboard/last_study_session", handlers.GetLastStudySession)
		api.GET("/dashboard/study_progress", handlers.GetStudyProgress)
		api.GET("/dashboard/quick-stats", handlers.GetQuickStats)

		// Words routes
		api.GET("/words", handlers.GetWords)
		api.GET("/words/:id", handlers.GetWord)

		// Groups routes
		api.GET("/groups", handlers.GetGroups)
		api.GET("/groups/:id", handlers.GetGroup)
		api.GET("/groups/:id/words", handlers.GetGroupWords)

		// Study sessions routes
		api.GET("/study_sessions", handlers.GetStudySessions)
		api.GET("/study_sessions/:id", handlers.GetStudySession)
		api.GET("/study_sessions/:id/words", handlers.GetStudySessionWords)

		// Study activities routes
		api.GET("/study_activities/:id", handlers.GetStudyActivity)
		api.POST("/study_activities", handlers.CreateStudyActivity)

		// Reset routes
		api.POST("/reset_history", handlers.ResetHistory)
		api.POST("/full_reset", handlers.FullReset)
	}

	log.Println("Server starting on http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
} 