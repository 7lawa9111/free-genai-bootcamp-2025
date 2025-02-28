package handlers

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/mohawa/lang-portal/backend_go/internal/services"
)

type WordHandler struct {
	wordService *services.WordService
}

func NewWordHandler() *WordHandler {
	return &WordHandler{
		wordService: services.NewWordService(),
	}
}

func GetWords(c *gin.Context) {
	handler := NewWordHandler()
	
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "100"))

	response, err := handler.wordService.GetWords(page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

func GetWord(c *gin.Context) {
	handler := NewWordHandler()
	
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid word ID"})
		return
	}

	word, err := handler.wordService.GetWord(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, word)
} 