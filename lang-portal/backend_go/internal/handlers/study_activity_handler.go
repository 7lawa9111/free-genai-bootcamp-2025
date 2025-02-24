package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"lang-portal/backend_go/internal/services"
)

type StudyActivityHandler struct {
	studyActivityService *services.StudyActivityService
}

func NewStudyActivityHandler(studyActivityService *services.StudyActivityService) *StudyActivityHandler {
	return &StudyActivityHandler{
		studyActivityService: studyActivityService,
	}
}

// GetStudyActivityByID godoc
// @Summary Get study activity details
// @Description Get detailed information about a specific study activity
// @Tags study-activities
// @Accept json
// @Produce json
// @Param id path int true "Study Activity ID"
// @Success 200 {object} models.StudyActivityDetails
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /study_activities/{id} [get]
func (h *StudyActivityHandler) GetStudyActivityByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid ID format"})
		return
	}

	activity, err := h.studyActivityService.GetActivityByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, activity)
}

// GetActivitySessions godoc
// @Summary Get study sessions for activity
// @Description Get paginated list of study sessions for a specific activity
// @Tags study-activities
// @Accept json
// @Produce json
// @Param id path int true "Study Activity ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(100)
// @Success 200 {object} ListResponse{items=[]models.StudySession}
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /study_activities/{id}/study_sessions [get]
func (h *StudyActivityHandler) GetActivitySessions(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid ID format"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))

	sessions, total, err := h.studyActivityService.GetActivitySessions(id, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, ListResponse{
		Items: sessions,
		Pagination: PaginationResponse{
			CurrentPage:  page,
			TotalPages:   (total + limit - 1) / limit,
			TotalItems:   total,
			ItemsPerPage: limit,
		},
	})
}

// CreateStudyActivity godoc
// @Summary Create a new study activity
// @Description Create a new study activity for a group
// @Tags study-activities
// @Accept json
// @Produce json
// @Param request body CreateStudyActivityRequest true "Study activity details"
// @Success 201 {object} models.StudyActivity
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /study_activities [post]
func (h *StudyActivityHandler) CreateStudyActivity(c *gin.Context) {
	var request CreateStudyActivityRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}

	activity, err := h.studyActivityService.CreateStudyActivity(request.GroupID, request.Type)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, activity)
}

// GetActivityProgress godoc
// @Summary Get activity progress
// @Description Get completion progress for a study activity
// @Tags study-activities
// @Accept json
// @Produce json
// @Param id path int true "Study Activity ID"
// @Success 200 {object} ProgressResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /study_activities/{id}/progress [get]
func (h *StudyActivityHandler) GetActivityProgress(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid ID format"})
		return
	}

	progress, err := h.studyActivityService.GetActivityProgress(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, ProgressResponse{
		Progress:    progress,
		IsComplete: progress >= 1.0,
	})
}

// Request/Response types
type CreateStudyActivityRequest struct {
	GroupID int64  `json:"group_id" binding:"required"`
	Type    string `json:"type" binding:"required"`
}

type ProgressResponse struct {
	Progress   float64 `json:"progress"`
	IsComplete bool    `json:"is_complete"`
} 