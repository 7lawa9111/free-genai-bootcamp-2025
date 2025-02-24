package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"lang-portal/backend_go/internal/services"
)

type WordHandler struct {
	wordService *services.WordService
}

func NewWordHandler(wordService *services.WordService) *WordHandler {
	return &WordHandler{
		wordService: wordService,
	}
}

// GetWords godoc
// @Summary List vocabulary words
// @Description Get paginated list of vocabulary words
// @Tags words
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(100)
// @Success 200 {object} ListResponse{items=[]models.Word}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /words [get]
func (h *WordHandler) GetWords(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit := 100 // As per spec

	words, totalCount, err := h.wordService.GetWords(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
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

// GetWordByID godoc
// @Summary Get a specific word
// @Description Get detailed information about a specific word
// @Tags words
// @Accept json
// @Produce json
// @Param id path int true "Word ID"
// @Success 200 {object} models.Word
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /words/{id} [get]
func (h *WordHandler) GetWordByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid ID"})
		return
	}

	word, err := h.wordService.GetWordByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, word)
} 