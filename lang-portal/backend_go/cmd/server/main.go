package main

import (
	"log"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"lang-portal/backend_go/internal/database"
	"lang-portal/backend_go/internal/handlers"
	"lang-portal/backend_go/internal/middleware"
	"lang-portal/backend_go/internal/repository"
	"lang-portal/backend_go/internal/services"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "lang-portal/backend_go/docs" // Import swagger docs
)

func main() {
	dbPath := filepath.Join("db", "words.db")
	if err := database.InitDB(dbPath); err != nil {
		log.Fatal("Failed to initialize database: ", err)
	}

	// Initialize repositories
	wordRepo := repository.NewWordRepository(database.DB)
	groupRepo := repository.NewGroupRepository(database.DB)
	studySessionRepo := repository.NewStudySessionRepository(database.DB)
	wordMatchingRepo := repository.NewWordMatchingRepository(database.DB)
	vocabularyQuizRepo := repository.NewVocabularyQuizRepository(database.DB)
	flashcardRepo := repository.NewFlashcardRepository(database.DB)
	writingRepo := repository.NewWritingPracticeRepository(database.DB)
	sentenceRepo := repository.NewSentenceConstructionRepository(database.DB)
	statsRepo := repository.NewStatisticsRepository(database.DB)

	// Initialize services
	wordService := services.NewWordService(wordRepo)
	groupService := services.NewGroupService(
		groupRepo,
		wordRepo,
		studySessionRepo,
	)
	studySessionService := services.NewStudySessionService(studySessionRepo)
	dashboardService := services.NewDashboardService(studySessionRepo, wordRepo, groupRepo)
	wordMatchingService := services.NewWordMatchingService(wordMatchingRepo)
	vocabularyQuizService := services.NewVocabularyQuizService(vocabularyQuizRepo)
	flashcardService := services.NewFlashcardService(flashcardRepo)
	writingService := services.NewWritingPracticeService(writingRepo)
	sentenceService := services.NewSentenceConstructionService(sentenceRepo)
	statsService := services.NewStatisticsService(statsRepo)

	// Initialize handlers
	wordHandler := handlers.NewWordHandler(wordService)
	groupHandler := handlers.NewGroupHandler(groupService)
	studySessionHandler := handlers.NewStudySessionHandler(studySessionService)
	dashboardHandler := handlers.NewDashboardHandler(dashboardService)
	wordMatchingHandler := handlers.NewWordMatchingHandler(wordMatchingService)
	vocabularyQuizHandler := handlers.NewVocabularyQuizHandler(vocabularyQuizService)
	flashcardHandler := handlers.NewFlashcardHandler(flashcardService)
	writingHandler := handlers.NewWritingPracticeHandler(writingService)
	sentenceHandler := handlers.NewSentenceConstructionHandler(sentenceService)
	statsHandler := handlers.NewStatisticsHandler(statsService)

	// Initialize router with middleware
	r := gin.Default()
	
	// Add middleware with custom CORS config if needed
	corsConfig := middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"}, // Allow only our frontend
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
		MaxAge:      3600, // 1 hour
	}
	r.Use(middleware.CORSWithConfig(corsConfig))
	r.Use(middleware.ErrorHandler())

	// Add Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes will be grouped under /api
	api := r.Group("/api")

	// Dashboard routes
	api.GET("/dashboard/last_study_session", dashboardHandler.GetLastStudySession)
	api.GET("/dashboard/study_progress", dashboardHandler.GetStudyProgress)
	api.GET("/dashboard/quick-stats", dashboardHandler.GetQuickStats)

	// Statistics routes
	statistics := api.Group("/statistics")
	{
		statistics.GET("/overview", statsHandler.GetOverview)
		statistics.GET("/activities", statsHandler.GetActivityStats)
		statistics.GET("/progress", statsHandler.GetStudyProgress)
	}

	// Groups routes
	groups := api.Group("/groups")
	{
		groups.GET("", groupHandler.GetGroups)
		groups.GET("/:id", groupHandler.GetGroupByID)
		groups.GET("/:id/words", groupHandler.GetGroupWords)
		groups.GET("/:id/study-sessions", groupHandler.GetGroupStudySessions)
	}

	// Words routes
	words := api.Group("/words")
	{
		words.GET("", wordHandler.GetWords)
		words.GET("/:id", wordHandler.GetWordByID)
	}

	// Study sessions routes
	studySessions := api.Group("/study-sessions")
	{
		studySessions.GET("", studySessionHandler.GetStudySessions)
		studySessions.GET("/:id", studySessionHandler.GetStudySessionByID)
		studySessions.GET("/:id/words", studySessionHandler.GetStudySessionWords)
	}

	// Study activities routes
	study := api.Group("/study")
	{
		// Vocabulary Quiz
		vocabularyQuiz := study.Group("/vocabulary-quiz")
		{
			vocabularyQuiz.POST("", vocabularyQuizHandler.CreateQuiz)
			vocabularyQuiz.POST("/:id/result", vocabularyQuizHandler.SaveResult)
			vocabularyQuiz.GET("/:id/stats", vocabularyQuizHandler.GetQuizStats)
			vocabularyQuiz.GET("/:id/progress", vocabularyQuizHandler.GetProgress)
			vocabularyQuiz.GET("/:id/debug", vocabularyQuizHandler.GetQuizDebug)
		}

		// Word Matching
		wordMatching := study.Group("/word-matching")
		{
			wordMatching.POST("", wordMatchingHandler.CreateActivity)
			wordMatching.POST("/:id/result", wordMatchingHandler.SaveResult)
			wordMatching.GET("/:id/stats", wordMatchingHandler.GetActivityStats)
			wordMatching.GET("/:id/progress", wordMatchingHandler.GetProgress)
		}

		// Writing Practice
		writingPractice := study.Group("/writing-practice")
		{
			writingPractice.POST("", writingHandler.CreateActivity)
			writingPractice.POST("/:id/result", writingHandler.SaveResult)
			writingPractice.GET("/:id/stats", writingHandler.GetActivityStats)
			writingPractice.GET("/:id/progress", writingHandler.GetProgress)
		}

		// Flashcards
		flashcards := study.Group("/flashcards")
		{
			flashcards.POST("", flashcardHandler.CreateActivity)
			flashcards.POST("/:id/result", flashcardHandler.SaveResult)
			flashcards.GET("/:id/stats", flashcardHandler.GetActivityStats)
			flashcards.GET("/:id/progress", flashcardHandler.GetProgress)
		}

		// Sentence Construction
		sentenceConstruction := study.Group("/sentence-construction")
		{
			sentenceConstruction.POST("", sentenceHandler.CreateActivity)
			sentenceConstruction.POST("/:id/result", sentenceHandler.SaveResult)
			sentenceConstruction.GET("/:id/stats", sentenceHandler.GetActivityStats)
			sentenceConstruction.GET("/:id/progress", sentenceHandler.GetProgress)
		}
	}

	// System routes
	api.POST("/reset_history", studySessionHandler.ResetHistory)
	api.POST("/full_reset", studySessionHandler.FullReset)

	// Start server
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
} 