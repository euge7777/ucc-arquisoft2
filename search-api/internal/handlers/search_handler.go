package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/gym-management/search-api/internal/models"
	"github.com/yourusername/gym-management/search-api/internal/services"
)

type SearchHandler struct {
	searchService *services.SearchService
	cacheService  *services.CacheService
}

func NewSearchHandler(searchService *services.SearchService, cacheService *services.CacheService) *SearchHandler {
	return &SearchHandler{
		searchService: searchService,
		cacheService:  cacheService,
	}
}

func (h *SearchHandler) Search(c *gin.Context) {
	var req models.SearchRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generar clave de caché
	cacheKey := services.GenerateCacheKey("search", req)

	// Intentar obtener del caché
	if cachedData, found := h.cacheService.Get(cacheKey); found {
		var response models.SearchResponse
		if err := json.Unmarshal(cachedData, &response); err == nil {
			c.Header("X-Cache", "HIT")
			c.JSON(http.StatusOK, response)
			return
		}
	}

	// Si no está en caché, realizar búsqueda
	response, err := h.searchService.Search(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Guardar en caché
	if responseData, err := json.Marshal(response); err == nil {
		h.cacheService.Set(cacheKey, responseData)
	}

	c.Header("X-Cache", "MISS")
	c.JSON(http.StatusOK, response)
}

func (h *SearchHandler) QuickSearch(c *gin.Context) {
	query := c.Query("q")
	typeFilter := c.Query("type")

	req := models.SearchRequest{
		Query:    query,
		Type:     typeFilter,
		Page:     1,
		PageSize: 20,
	}

	response, err := h.searchService.Search(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *SearchHandler) GetDocument(c *gin.Context) {
	docID := c.Param("id")

	doc, err := h.searchService.GetDocumentByID(docID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, doc)
}

func (h *SearchHandler) IndexDocument(c *gin.Context) {
	var doc models.SearchDocument

	if err := c.ShouldBindJSON(&doc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.searchService.IndexDocument(doc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Invalidar caché
	h.cacheService.InvalidatePattern(doc.Type)

	c.JSON(http.StatusCreated, gin.H{"message": "Documento indexado correctamente"})
}

func (h *SearchHandler) DeleteDocument(c *gin.Context) {
	docID := c.Param("id")

	err := h.searchService.DeleteDocument(docID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Documento eliminado correctamente"})
}

func (h *SearchHandler) GetStats(c *gin.Context) {
	stats := h.searchService.GetStats()
	c.JSON(http.StatusOK, stats)
}

func (h *SearchHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "search-api",
	})
}
