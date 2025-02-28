package handlers

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/mohawa/lang-portal/backend_go/internal/services"
	"github.com/mohawa/lang-portal/backend_go/internal/models"
)

type StudySessionHandler struct {
	studyService *services.StudyService
}

func NewStudySessionHandler() *StudySessionHandler {
	return &StudySessionHandler{
		studyService: services.NewStudyService(),
	}
}

// GetStudySessions handles GET /api/study_sessions
func GetStudySessions(c *gin.Context) {
	handler := NewStudySessionHandler()
	
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "100"))

	response, err := handler.studyService.GetStudySessions(page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

// GetStudySession handles GET /api/study_sessions/:id
func GetStudySession(c *gin.Context) {
	handler := NewStudySessionHandler()
	
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	session, err := handler.studyService.GetStudyActivity(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, session)
}

// GetStudySessionWords handles GET /api/study_sessions/:id/words
func GetStudySessionWords(c *gin.Context) {
	handler := NewStudySessionHandler()
	
	sessionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "100"))

	response, err := handler.studyService.GetStudySessions(page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	filteredResponse := filterWordsBySessionID(response, sessionID)
	c.JSON(http.StatusOK, filteredResponse)
}

// ReviewWord handles POST /api/study_sessions/:id/words/:word_id/review
func ReviewWord(c *gin.Context) {
	handler := NewStudySessionHandler()
	
	sessionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	wordID, err := strconv.Atoi(c.Param("word_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid word ID"})
		return
	}

	var request struct {
		Correct bool `json:"correct"`
	}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err = handler.studyService.ReviewWord(sessionID, wordID, request.Correct)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"word_id": wordID,
		"study_session_id": sessionID,
		"correct": request.Correct,
	})
}

func filterWordsBySessionID(response *models.PaginatedResponse, sessionID int) *models.PaginatedResponse {
	if sessions, ok := response.Items.([]models.StudySession); ok {
		for _, session := range sessions {
			if session.ID == sessionID {
				return response
			}
		}
	}
	response.Items = []models.StudySession{}
	response.Pagination.TotalItems = 0
	response.Pagination.TotalPages = 0
	return response
} 