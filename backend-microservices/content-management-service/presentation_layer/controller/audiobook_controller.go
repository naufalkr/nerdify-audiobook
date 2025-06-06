package controller

import (
	"content-management-service/data_layer/dto"
	"content-management-service/domain_layer/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AudiobookController struct {
	audiobookService *service.AudiobookService
}

func NewAudiobookController(audiobookService *service.AudiobookService) *AudiobookController {
	return &AudiobookController{
		audiobookService: audiobookService,
	}
}

// CreateAudiobook creates a new audiobook
func (ac *AudiobookController) CreateAudiobook(c *gin.Context) {
	var req dto.CreateAudiobookRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	audiobook, err := ac.audiobookService.CreateAudiobook(req)
	if err != nil {
		if err.Error() == "author not found" || err.Error() == "reader not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, audiobook)
}

// GetAudiobookByID retrieves an audiobook by ID
func (ac *AudiobookController) GetAudiobookByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid audiobook ID"})
		return
	}

	audiobook, err := ac.audiobookService.GetAudiobookByID(uint(id))
	if err != nil {
		if err.Error() == "audiobook not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, audiobook)
}

// GetAllAudiobooks retrieves all audiobooks with pagination and optional filtering
func (ac *AudiobookController) GetAllAudiobooks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	authorID, _ := strconv.ParseUint(c.Query("author_id"), 10, 32)
	readerID, _ := strconv.ParseUint(c.Query("reader_id"), 10, 32)
	genreID, _ := strconv.ParseUint(c.Query("genre_id"), 10, 32)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	filter := dto.AudiobookFilter{
		AuthorID: uint(authorID),
		ReaderID: uint(readerID),
		GenreID:  uint(genreID),
	}

	audiobooks, err := ac.audiobookService.GetAudiobooks(filter, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, audiobooks)
}

// UpdateAudiobook updates an existing audiobook
func (ac *AudiobookController) UpdateAudiobook(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid audiobook ID"})
		return
	}

	var req dto.UpdateAudiobookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	audiobook, err := ac.audiobookService.UpdateAudiobook(uint(id), req)
	if err != nil {
		if err.Error() == "audiobook not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "author not found" || err.Error() == "reader not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, audiobook)
}

// DeleteAudiobook deletes an audiobook
func (ac *AudiobookController) DeleteAudiobook(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid audiobook ID"})
		return
	}

	if err := ac.audiobookService.DeleteAudiobook(uint(id)); err != nil {
		if err.Error() == "audiobook not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Audiobook deleted successfully"})
}

// SearchAudiobooks searches audiobooks by title
func (ac *AudiobookController) SearchAudiobooks(c *gin.Context) {
	title := c.Query("title")
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title parameter is required"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	searchReq := dto.SearchRequest{
		Query: title,
		PaginationRequest: dto.PaginationRequest{
			Page:  page,
			Limit: limit,
		},
	}

	audiobooks, err := ac.audiobookService.SearchAudiobooks(searchReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, audiobooks)
}

// AddGenresToAudiobook adds genres to an audiobook
func (ac *AudiobookController) AddGenresToAudiobook(c *gin.Context) {
	audiobookID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid audiobook ID"})
		return
	}

	var req struct {
		GenreIDs []uint `json:"genre_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ac.audiobookService.AddGenresToAudiobook(uint(audiobookID), req.GenreIDs); err != nil {
		if err.Error() == "audiobook not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Genres added successfully"})
}

// RemoveGenresFromAudiobook removes genres from an audiobook
func (ac *AudiobookController) RemoveGenresFromAudiobook(c *gin.Context) {
	audiobookID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid audiobook ID"})
		return
	}

	var req struct {
		GenreIDs []uint `json:"genre_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ac.audiobookService.RemoveGenresFromAudiobook(uint(audiobookID), req.GenreIDs); err != nil {
		if err.Error() == "audiobook not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Genres removed successfully"})
}
