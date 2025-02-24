package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"lang-portal/backend_go/internal/models"
	"lang-portal/backend_go/internal/services"
)

type WritingPracticeHandler struct {
	writingService *services.WritingPracticeService
}

func NewWritingPracticeHandler(writingService *services.WritingPracticeService) *WritingPracticeHandler {
	return &WritingPracticeHandler{
		writingService: writingService,
	}
}

// CreateActivity godoc
// @Summary Create a new writing practice activity
// @Description Create a new writing practice activity for a group
// @Tags writing-practice
// @Accept json
// @Produce json
// @Param request body CreateWritingRequest true "Group ID for practice"
// @Success 201 {object} models.WritingPractice
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /writing-practice/activities [post]
func (h *WritingPracticeHandler) CreateActivity(c *gin.Context) {
	var request CreateWritingRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}

	activity, err := h.writingService.CreateActivity(request.GroupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, activity)
}

// SaveResult godoc
// @Summary Save writing practice result
// @Description Save the result of a completed writing practice activity
// @Tags writing-practice
// @Accept json
// @Produce json
// @Param id path int true "Activity ID"
// @Param request body models.WritingResult true "Activity result"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /writing-practice/activities/{id}/result [post]
func (h *WritingPracticeHandler) SaveResult(c *gin.Context) {
	var result models.WritingResult
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

	if err := h.writingService.SaveResult(&result); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Message: "Result saved successfully"})
}

// GetActivityStats godoc
// @Summary Get writing practice statistics
// @Description Get statistics for a completed writing practice activity
// @Tags writing-practice
// @Accept json
// @Produce json
// @Param id path int true "Activity ID"
// @Success 200 {object} models.WritingResult
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /writing-practice/activities/{id}/stats [get]
func (h *WritingPracticeHandler) GetActivityStats(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid activity ID"})
		return
	}

	stats, err := h.writingService.GetActivityStats(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetProgress godoc
// @Summary Get activity progress
// @Description Get completion progress for a writing practice activity
// @Tags writing-practice
// @Accept json
// @Produce json
// @Param id path int true "Activity ID"
// @Success 200 {object} ProgressResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /writing-practice/activities/{id}/progress [get]
func (h *WritingPracticeHandler) GetProgress(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid activity ID"})
		return
	}

	progress, err := h.writingService.GetProgress(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	isComplete, err := h.writingService.IsActivityComplete(id)
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
type CreateWritingRequest struct {
	GroupID int64 `json:"group_id" binding:"required"`
} 