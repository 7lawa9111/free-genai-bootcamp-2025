package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"lang-portal/backend_go/internal/services"
)

type StatisticsHandler struct {
	statsService *services.StatisticsService
}

func NewStatisticsHandler(statsService *services.StatisticsService) *StatisticsHandler {
	return &StatisticsHandler{
		statsService: statsService,
	}
}

// GetUserStats godoc
// @Summary Get user statistics
// @Description Get overall user study statistics
// @Tags statistics
// @Accept json
// @Produce json
// @Success 200 {object} models.UserStats
// @Failure 500 {object} ErrorResponse
// @Router /statistics/user [get]
func (h *StatisticsHandler) GetUserStats(c *gin.Context) {
	stats, err := h.statsService.GetUserStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

// GetActivityStats godoc
// @Summary Get activity statistics
// @Description Get statistics broken down by activity type
// @Tags statistics
// @Accept json
// @Produce json
// @Success 200 {object} map[string]models.ActivityStats
// @Failure 500 {object} ErrorResponse
// @Router /statistics/activities [get]
func (h *StatisticsHandler) GetActivityStats(c *gin.Context) {
	stats, err := h.statsService.GetActivityStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

// GetStudyProgress godoc
// @Summary Get study progress
// @Description Get daily and weekly study progress statistics
// @Tags statistics
// @Accept json
// @Produce json
// @Success 200 {object} models.StudyProgress
// @Failure 500 {object} ErrorResponse
// @Router /statistics/progress [get]
func (h *StatisticsHandler) GetStudyProgress(c *gin.Context) {
	progress, err := h.statsService.GetStudyProgress()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, progress)
}

// GetStudyMetrics godoc
// @Summary Get study metrics
// @Description Get additional study metrics like average study time and completion rate
// @Tags statistics
// @Accept json
// @Produce json
// @Success 200 {object} StudyMetricsResponse
// @Failure 500 {object} ErrorResponse
// @Router /statistics/metrics [get]
func (h *StatisticsHandler) GetStudyMetrics(c *gin.Context) {
	avgTime, err := h.statsService.GetAverageStudyTime()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	completionRate, err := h.statsService.GetCompletionRate()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, StudyMetricsResponse{
		AverageStudyTimeMinutes: avgTime,
		CompletionRatePercent:   completionRate,
	})
}

// Response types
type StudyMetricsResponse struct {
	AverageStudyTimeMinutes float64 `json:"average_study_time_minutes"`
	CompletionRatePercent   float64 `json:"completion_rate_percent"`
}

// GetOverview is an alias for GetUserStats for backward compatibility
func (h *StatisticsHandler) GetOverview(c *gin.Context) {
	h.GetUserStats(c)
} 