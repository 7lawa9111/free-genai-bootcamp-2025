package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"lang-portal/backend_go/internal/models"
	"lang-portal/backend_go/internal/services"
)

type FlashcardHandler struct {
	flashcardService *services.FlashcardService
}

func NewFlashcardHandler(flashcardService *services.FlashcardService) *FlashcardHandler {
	return &FlashcardHandler{
		flashcardService: flashcardService,
	}
}

// CreateActivity godoc
// @Summary Create a new flashcard activity
// @Description Create a new flashcard activity for a group
// @Tags flashcards
// @Accept json
// @Produce json
// @Param request body CreateFlashcardRequest true "Activity details"
// @Success 201 {object} models.FlashcardActivity
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /flashcards/activities [post]
func (h *FlashcardHandler) CreateActivity(c *gin.Context) {
	var request CreateFlashcardRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}

	activity, err := h.flashcardService.CreateActivity(request.GroupID, request.Direction)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, activity)
}

// SaveResult godoc
// @Summary Save flashcard activity result
// @Description Save the result of a completed flashcard activity
// @Tags flashcards
// @Accept json
// @Produce json
// @Param id path int true "Activity ID"
// @Param request body models.FlashcardResult true "Activity result"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /flashcards/activities/{id}/result [post]
func (h *FlashcardHandler) SaveResult(c *gin.Context) {
	var result models.FlashcardResult
	if err := c.ShouldBindJSON(&result); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}

	// Validate activity ID from path matches body
	activityID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid activity ID"})
		return
	}
	if result.ActivityID != activityID {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "activity ID mismatch"})
		return
	}

	if err := h.flashcardService.SaveResult(&result); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Message: "Result saved successfully"})
}

// GetActivityStats godoc
// @Summary Get flashcard activity statistics
// @Description Get statistics for a completed flashcard activity
// @Tags flashcards
// @Accept json
// @Produce json
// @Param id path int true "Activity ID"
// @Success 200 {object} models.FlashcardResult
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /flashcards/activities/{id}/stats [get]
func (h *FlashcardHandler) GetActivityStats(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid activity ID"})
		return
	}

	stats, err := h.flashcardService.GetActivityStats(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetProgress godoc
// @Summary Get activity progress
// @Description Get completion progress for a flashcard activity
// @Tags flashcards
// @Accept json
// @Produce json
// @Param id path int true "Activity ID"
// @Success 200 {object} ProgressResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /flashcards/activities/{id}/progress [get]
func (h *FlashcardHandler) GetProgress(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid activity ID"})
		return
	}

	progress, err := h.flashcardService.GetProgress(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	isComplete, err := h.flashcardService.IsActivityComplete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, ProgressResponse{
		Progress:   progress,
		IsComplete: isComplete,
	})
}

// Request/Response types
type CreateFlashcardRequest struct {
	GroupID   int64  `json:"group_id" binding:"required"`
	Direction string `json:"direction" binding:"required"`
} 