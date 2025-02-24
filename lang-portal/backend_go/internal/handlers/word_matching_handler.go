package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"lang-portal/backend_go/internal/models"
	"lang-portal/backend_go/internal/services"
)

type WordMatchingHandler struct {
	wordMatchingService *services.WordMatchingService
}

func NewWordMatchingHandler(wordMatchingService *services.WordMatchingService) *WordMatchingHandler {
	return &WordMatchingHandler{
		wordMatchingService: wordMatchingService,
	}
}

// CreateActivity godoc
// @Summary Create a new word matching activity
// @Description Create a new word matching activity for a group
// @Tags word-matching
// @Accept json
// @Produce json
// @Param request body CreateActivityRequest true "Group ID for activity"
// @Success 201 {object} models.WordMatchingActivity
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /word-matching/activities [post]
func (h *WordMatchingHandler) CreateActivity(c *gin.Context) {
	var request CreateActivityRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}

	activity, err := h.wordMatchingService.CreateActivity(request.GroupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, activity)
}

// SaveResult godoc
// @Summary Save word matching activity result
// @Description Save the result of a completed word matching activity
// @Tags word-matching
// @Accept json
// @Produce json
// @Param request body models.WordMatchingResult true "Activity result"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /word-matching/activities/{id}/result [post]
func (h *WordMatchingHandler) SaveResult(c *gin.Context) {
	var result models.WordMatchingResult
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

	if err := h.wordMatchingService.SaveResult(&result); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Message: "Result saved successfully"})
}

// GetActivityStats godoc
// @Summary Get word matching activity statistics
// @Description Get statistics for a completed word matching activity
// @Tags word-matching
// @Accept json
// @Produce json
// @Param id path int true "Activity ID"
// @Success 200 {object} models.WordMatchingResult
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /word-matching/activities/{id}/stats [get]
func (h *WordMatchingHandler) GetActivityStats(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid activity ID"})
		return
	}

	stats, err := h.wordMatchingService.GetActivityStats(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetProgress godoc
// @Summary Get activity progress
// @Description Get completion progress for a word matching activity
// @Tags word-matching
// @Accept json
// @Produce json
// @Param id path int true "Activity ID"
// @Success 200 {object} ProgressResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /word-matching/activities/{id}/progress [get]
func (h *WordMatchingHandler) GetProgress(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid activity ID"})
		return
	}

	progress, err := h.wordMatchingService.GetProgress(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	isComplete, err := h.wordMatchingService.IsActivityComplete(id)
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
type CreateActivityRequest struct {
	GroupID int64 `json:"group_id" binding:"required"`
} 