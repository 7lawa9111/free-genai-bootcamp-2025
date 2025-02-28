package handlers

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/mohawa/lang-portal/backend_go/internal/services"
)

type GroupHandler struct {
	groupService *services.GroupService
}

func NewGroupHandler() *GroupHandler {
	return &GroupHandler{
		groupService: services.NewGroupService(),
	}
}

func GetGroups(c *gin.Context) {
	handler := NewGroupHandler()
	
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "100"))

	response, err := handler.groupService.GetGroups(page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

func GetGroup(c *gin.Context) {
	handler := NewGroupHandler()
	
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	group, err := handler.groupService.GetGroup(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, group)
}

func GetGroupWords(c *gin.Context) {
	handler := NewGroupHandler()
	
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "100"))

	response, err := handler.groupService.GetGroupWords(id, page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
} 