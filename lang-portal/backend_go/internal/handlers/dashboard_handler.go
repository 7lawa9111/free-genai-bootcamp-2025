package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"lang-portal/backend_go/internal/services"
)

type DashboardHandler struct {
	dashboardService *services.DashboardService
}

func NewDashboardHandler(dashboardService *services.DashboardService) *DashboardHandler {
	return &DashboardHandler{
		dashboardService: dashboardService,
	}
}

// GetLastStudySession godoc
// @Summary Get the most recent study session
// @Description Returns information about the most recent study session
// @Tags dashboard
// @Accept json
// @Produce json
// @Success 200 {object} models.StudySession
// @Failure 500 {object} ErrorResponse
// @Router /dashboard/last_study_session [get]
func (h *DashboardHandler) GetLastStudySession(c *gin.Context) {
	session, err := h.dashboardService.GetLastStudySession()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, session)
}

// GetStudyProgress godoc
// @Summary Get study progress statistics
// @Description Returns study progress including total words studied and available
// @Tags dashboard
// @Accept json
// @Produce json
// @Success 200 {object} map[string]int
// @Failure 500 {object} ErrorResponse
// @Router /dashboard/study_progress [get]
func (h *DashboardHandler) GetStudyProgress(c *gin.Context) {
	progress, err := h.dashboardService.GetStudyProgress()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, progress)
}

// GetQuickStats godoc
// @Summary Get quick overview statistics
// @Description Returns various statistics about study progress
// @Tags dashboard
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} ErrorResponse
// @Router /dashboard/quick-stats [get]
func (h *DashboardHandler) GetQuickStats(c *gin.Context) {
	stats, err := h.dashboardService.GetQuickStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
} 