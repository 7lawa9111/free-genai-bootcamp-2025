package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"lang-portal/backend_go/internal/services"
	"lang-portal/backend_go/internal/errors"
)

type GroupHandler struct {
	groupService *services.GroupService
}

func NewGroupHandler(groupService *services.GroupService) *GroupHandler {
	return &GroupHandler{
		groupService: groupService,
	}
}

func (h *GroupHandler) GetGroups(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit := 100 // As per spec

	groups, totalCount, err := h.groupService.GetGroups(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	totalPages := (totalCount + limit - 1) / limit

	c.JSON(http.StatusOK, ListResponse{
		Items: groups,
		Pagination: PaginationResponse{
			CurrentPage:  page,
			TotalPages:   totalPages,
			TotalItems:   totalCount,
			ItemsPerPage: limit,
		},
	})
}

func (h *GroupHandler) GetGroupByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid ID"})
		return
	}

	group, err := h.groupService.GetGroupByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, group)
}

func (h *GroupHandler) GetGroupWords(c *gin.Context) {
	groupID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid group ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit := 100

	// First check if group exists
	_, err = h.groupService.GetGroupByID(groupID)
	if err != nil {
		if err == errors.ErrGroupNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "Group not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	words, totalCount, err := h.groupService.GetGroupWords(groupID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: fmt.Sprintf("Error fetching words: %v", err)})
		return
	}

	totalPages := (totalCount + limit - 1) / limit

	c.JSON(http.StatusOK, ListResponse{
		Items: words,
		Pagination: PaginationResponse{
			CurrentPage:  page,
			TotalPages:   totalPages,
			TotalItems:   totalCount,
			ItemsPerPage: limit,
		},
	})
}

func (h *GroupHandler) GetGroupStudySessions(c *gin.Context) {
	groupID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid group ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit := 100

	sessions, totalCount, err := h.groupService.GetGroupStudySessions(groupID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	totalPages := (totalCount + limit - 1) / limit

	c.JSON(http.StatusOK, ListResponse{
		Items: sessions,
		Pagination: PaginationResponse{
			CurrentPage:  page,
			TotalPages:   totalPages,
			TotalItems:   totalCount,
			ItemsPerPage: limit,
		},
	})
} 