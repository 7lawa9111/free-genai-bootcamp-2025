package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"lang-portal/backend_go/internal/models"
	"lang-portal/backend_go/internal/services"
)

type SentenceConstructionHandler struct {
	sentenceService *services.SentenceConstructionService
}

func NewSentenceConstructionHandler(sentenceService *services.SentenceConstructionService) *SentenceConstructionHandler {
	return &SentenceConstructionHandler{
		sentenceService: sentenceService,
	}
}

// CreateActivity godoc
// @Summary Create a new sentence construction activity
// @Description Create a new sentence construction activity for a group
// @Tags sentence-construction
// @Accept json
// @Produce json
// @Param request body CreateSentenceRequest true "Group ID for practice"
// @Success 201 {object} models.SentenceConstruction
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /sentence-construction/activities [post]
func (h *SentenceConstructionHandler) CreateActivity(c *gin.Context) {
	var request CreateSentenceRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}

	activity, err := h.sentenceService.CreateActivity(request.GroupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, activity)
}

// SaveResult godoc
// @Summary Save sentence construction result
// @Description Save the result of a completed sentence construction activity
// @Tags sentence-construction
// @Accept json
// @Produce json
// @Param id path int true "Activity ID"
// @Param request body models.SentenceResult true "Activity result"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /sentence-construction/activities/{id}/result [post]
func (h *SentenceConstructionHandler) SaveResult(c *gin.Context) {
	var result models.SentenceResult
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

	if err := h.sentenceService.SaveResult(&result); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Message: "Result saved successfully"})
}

// GetActivityStats godoc
// @Summary Get sentence construction statistics
// @Description Get statistics for a completed sentence construction activity
// @Tags sentence-construction
// @Accept json
// @Produce json
// @Param id path int true "Activity ID"
// @Success 200 {object} models.SentenceResult
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /sentence-construction/activities/{id}/stats [get]
func (h *SentenceConstructionHandler) GetActivityStats(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid activity ID"})
		return
	}

	stats, err := h.sentenceService.GetActivityStats(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetProgress godoc
// @Summary Get activity progress
// @Description Get completion progress for a sentence construction activity
// @Tags sentence-construction
// @Accept json
// @Produce json
// @Param id path int true "Activity ID"
// @Success 200 {object} ProgressResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /sentence-construction/activities/{id}/progress [get]
func (h *SentenceConstructionHandler) GetProgress(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid activity ID"})
		return
	}

	progress, err := h.sentenceService.GetProgress(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	isComplete, err := h.sentenceService.IsActivityComplete(id)
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
type CreateSentenceRequest struct {
	GroupID int64 `json:"group_id" binding:"required"`
} 