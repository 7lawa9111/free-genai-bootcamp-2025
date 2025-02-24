package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"lang-portal/backend_go/internal/services"
)

type StudySessionHandler struct {
	studySessionService *services.StudySessionService
}

func NewStudySessionHandler(studySessionService *services.StudySessionService) *StudySessionHandler {
	return &StudySessionHandler{
		studySessionService: studySessionService,
	}
}

// GetStudySessions godoc
// @Summary List study sessions
// @Description Get paginated list of study sessions
// @Tags study-sessions
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(100)
// @Success 200 {object} ListResponse{items=[]models.StudySession}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /study_sessions [get]
func (h *StudySessionHandler) GetStudySessions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))

	sessions, total, err := h.studySessionService.GetStudySessions(page, limit)
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

// GetStudySessionByID godoc
// @Summary Get study session details
// @Description Get detailed information about a specific study session
// @Tags study-sessions
// @Accept json
// @Produce json
// @Param id path int true "Study Session ID"
// @Success 200 {object} services.StudySessionDetails
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /study_sessions/{id} [get]
func (h *StudySessionHandler) GetStudySessionByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid ID format"})
		return
	}

	session, err := h.studySessionService.GetStudySessionByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, session)
}

// GetStudySessionWords godoc
// @Summary Get words from study session
// @Description Get paginated list of words reviewed in a study session
// @Tags study-sessions
// @Accept json
// @Produce json
// @Param id path int true "Study Session ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(100)
// @Success 200 {object} ListResponse{items=[]models.Word}
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /study_sessions/{id}/words [get]
func (h *StudySessionHandler) GetStudySessionWords(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid ID format"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))

	words, total, err := h.studySessionService.GetStudySessionWords(id, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, ListResponse{
		Items: words,
		Pagination: PaginationResponse{
			CurrentPage:  page,
			TotalPages:   (total + limit - 1) / limit,
			TotalItems:   total,
			ItemsPerPage: limit,
		},
	})
}

// ReviewWord godoc
// @Summary Record word review result
// @Description Record whether a word was correctly reviewed in a study session
// @Tags study-sessions
// @Accept json
// @Produce json
// @Param id path int true "Study Session ID"
// @Param word_id path int true "Word ID"
// @Param correct query bool true "Whether the review was correct"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /study_sessions/{id}/words/{word_id}/review [post]
func (h *StudySessionHandler) ReviewWord(c *gin.Context) {
	sessionID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid session ID format"})
		return
	}

	wordID, err := strconv.ParseInt(c.Param("word_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid word ID format"})
		return
	}

	correct, _ := strconv.ParseBool(c.Query("correct"))

	if err := h.studySessionService.ReviewWord(sessionID, wordID, correct); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Message: "Review recorded successfully"})
}

// ResetHistory godoc
// @Summary Reset study history
// @Description Delete all study sessions and word reviews
// @Tags study-sessions
// @Accept json
// @Produce json
// @Success 200 {object} SuccessResponse
// @Failure 500 {object} ErrorResponse
// @Router /reset_history [post]
func (h *StudySessionHandler) ResetHistory(c *gin.Context) {
	if err := h.studySessionService.ResetHistory(); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Message: "Study history reset successfully"})
}

// CreateStudySession godoc
// @Summary Create a new study session
// @Description Create a new study session for a group and activity
// @Tags study-sessions
// @Accept json
// @Produce json
// @Param request body CreateStudySessionRequest true "Study session details"
// @Success 201 {object} models.StudySession
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /study_sessions [post]
func (h *StudySessionHandler) CreateStudySession(c *gin.Context) {
	var request CreateStudySessionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}

	session, err := h.studySessionService.CreateStudySession(request.GroupID, request.StudyActivityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, session)
}

// GetStudyStats godoc
// @Summary Get study statistics
// @Description Get statistics about study sessions and word reviews
// @Tags study-sessions
// @Accept json
// @Produce json
// @Success 200 {object} StudyStatsResponse
// @Failure 500 {object} ErrorResponse
// @Router /study_sessions/stats [get]
func (h *StudySessionHandler) GetStudyStats(c *gin.Context) {
	stats, err := h.studySessionService.GetStudyStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// FullReset godoc
// @Summary Full system reset
// @Description Reset all study sessions, word reviews, and related data
// @Tags system
// @Accept json
// @Produce json
// @Success 200 {object} SuccessResponse
// @Failure 500 {object} ErrorResponse
// @Router /full_reset [post]
func (h *StudySessionHandler) FullReset(c *gin.Context) {
	if err := h.studySessionService.FullReset(); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Message: "System has been fully reset"})
}

// Request/Response types
type CreateStudySessionRequest struct {
	GroupID         int64 `json:"group_id" binding:"required"`
	StudyActivityID int64 `json:"study_activity_id" binding:"required"`
}

type StudyStatsResponse struct {
	TotalSessions     int     `json:"total_sessions"`
	TotalWordsReviewed int     `json:"total_words_reviewed"`
	CorrectRate       float64 `json:"correct_rate"`
	AverageSessionTime int     `json:"average_session_time_seconds"`
	LastSessionDate    string  `json:"last_session_date"`
} 