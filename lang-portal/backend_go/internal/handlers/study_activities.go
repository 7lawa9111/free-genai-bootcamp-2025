package handlers

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/mohawa/lang-portal/backend_go/internal/services"
	"github.com/mohawa/lang-portal/backend_go/internal/models"
)

type StudyActivityHandler struct {
	studyService *services.StudyService
}

func NewStudyActivityHandler() *StudyActivityHandler {
	return &StudyActivityHandler{
		studyService: services.NewStudyService(),
	}
}

func GetStudyActivity(c *gin.Context) {
	handler := NewStudyActivityHandler()
	
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity ID"})
		return
	}

	activity, err := handler.studyService.GetStudyActivity(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, activity)
}

func CreateStudyActivity(c *gin.Context) {
	handler := NewStudyActivityHandler()

	var request struct {
		GroupID         int `json:"group_id"`
		StudyActivityID int `json:"study_activity_id"`
	}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	session, err := handler.studyService.CreateStudyActivity(request.GroupID, request.StudyActivityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":       session.ID,
		"group_id": session.GroupID,
	})
}

func GetStudyActivitySessions(c *gin.Context) {
	handler := NewStudyActivityHandler()
	
	activityID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "100"))

	response, err := handler.studyService.GetStudySessions(page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Filter sessions by activity ID
	filteredResponse := filterSessionsByActivityID(response, activityID)
	c.JSON(http.StatusOK, filteredResponse)
}

func filterSessionsByActivityID(response *models.PaginatedResponse, activityID int) *models.PaginatedResponse {
	if sessions, ok := response.Items.([]models.StudySession); ok {
		var filtered []models.StudySession
		for _, session := range sessions {
			if session.StudyActivityID == activityID {
				filtered = append(filtered, session)
			}
		}
		response.Items = filtered
		response.Pagination.TotalItems = len(filtered)
		response.Pagination.TotalPages = (len(filtered) + response.Pagination.ItemsPerPage - 1) / response.Pagination.ItemsPerPage
	}
	return response
} 