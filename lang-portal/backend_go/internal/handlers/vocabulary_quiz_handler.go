package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"lang-portal/backend_go/internal/models"
	"lang-portal/backend_go/internal/services"
	"lang-portal/backend_go/internal/validation"
)

type VocabularyQuizHandler struct {
	quizService *services.VocabularyQuizService
}

func NewVocabularyQuizHandler(quizService *services.VocabularyQuizService) *VocabularyQuizHandler {
	return &VocabularyQuizHandler{
		quizService: quizService,
	}
}

// CreateQuiz godoc
// @Summary Create a new vocabulary quiz
// @Description Create a new vocabulary quiz for a group
// @Tags vocabulary-quiz
// @Accept json
// @Produce json
// @Param request body CreateQuizRequest true "Group ID for quiz"
// @Success 201 {object} models.VocabularyQuiz
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /vocabulary-quiz/quizzes [post]
func (h *VocabularyQuizHandler) CreateQuiz(c *gin.Context) {
	var req CreateQuizRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}

	if err := validation.ValidateID(req.GroupID); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid group ID"})
		return
	}

	quiz, err := h.quizService.CreateQuiz(req.GroupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, quiz)
}

// SaveResult godoc
// @Summary Save quiz result
// @Description Save the result of a completed vocabulary quiz
// @Tags vocabulary-quiz
// @Accept json
// @Produce json
// @Param id path int true "Quiz ID"
// @Param request body models.QuizResult true "Quiz result"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /vocabulary-quiz/quizzes/{id}/result [post]
func (h *VocabularyQuizHandler) SaveResult(c *gin.Context) {
	var result models.QuizResult
	if err := c.ShouldBindJSON(&result); err != nil {
		fmt.Printf("Debug - JSON binding error: %v\n", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}

	fmt.Printf("Debug - Received result: %+v\n", result)

	// Validate quiz ID from path matches body
	quizID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid quiz ID"})
		return
	}
	if result.ActivityID != quizID {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "quiz ID mismatch"})
		return
	}

	if err := h.quizService.SaveResult(&result); err != nil {
		fmt.Printf("Debug - Save error: %v\n", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Message: "Result saved successfully"})
}

// GetQuizStats godoc
// @Summary Get quiz statistics
// @Description Get statistics for a completed vocabulary quiz
// @Tags vocabulary-quiz
// @Accept json
// @Produce json
// @Param id path int true "Quiz ID"
// @Success 200 {object} models.QuizResult
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /vocabulary-quiz/quizzes/{id}/stats [get]
func (h *VocabularyQuizHandler) GetQuizStats(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid quiz ID"})
		return
	}

	stats, err := h.quizService.GetQuizStats(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetProgress godoc
// @Summary Get quiz progress
// @Description Get completion progress for a vocabulary quiz
// @Tags vocabulary-quiz
// @Accept json
// @Produce json
// @Param id path int true "Quiz ID"
// @Success 200 {object} ProgressResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /vocabulary-quiz/quizzes/{id}/progress [get]
func (h *VocabularyQuizHandler) GetProgress(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid quiz ID"})
		return
	}

	progress, err := h.quizService.GetProgress(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	isComplete, err := h.quizService.IsQuizComplete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, ProgressResponse{
		Progress:   progress,
		IsComplete: isComplete,
	})
}

// Debug endpoint
func (h *VocabularyQuizHandler) GetQuizDebug(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid quiz ID"})
		return
	}

	debug, err := h.quizService.GetQuizDebug(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, debug)
}

// Request/Response types
type CreateQuizRequest struct {
	GroupID int64 `json:"group_id" binding:"required"`
} 