package controller

import (
	"content-management-service/data_layer/dto"
	"content-management-service/domain_layer/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GenreController struct {
	genreService *service.GenreService
}

func NewGenreController(genreService *service.GenreService) *GenreController {
	return &GenreController{
		genreService: genreService,
	}
}

// CreateGenre creates a new genre
func (gc *GenreController) CreateGenre(c *gin.Context) {
	var req dto.CreateGenreRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	genre, err := gc.genreService.CreateGenre(req)
	if err != nil {
		if err.Error() == "genre already exists" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, genre)
}

// GetGenreByID retrieves a genre by ID
func (gc *GenreController) GetGenreByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid genre ID"})
		return
	}

	genre, err := gc.genreService.GetGenreByID(uint(id))
	if err != nil {
		if err.Error() == "genre not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, genre)
}

// GetAllGenres retrieves all genres with pagination
func (gc *GenreController) GetAllGenres(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	paginationReq := dto.PaginationRequest{
		Page:  page,
		Limit: limit,
	}

	genres, err := gc.genreService.GetAllGenres(paginationReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, genres)
}

// UpdateGenre updates an existing genre
func (gc *GenreController) UpdateGenre(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid genre ID"})
		return
	}

	var req dto.UpdateGenreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	genre, err := gc.genreService.UpdateGenre(uint(id), req)
	if err != nil {
		if err.Error() == "genre not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "genre already exists" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, genre)
}

// DeleteGenre deletes a genre
func (gc *GenreController) DeleteGenre(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid genre ID"})
		return
	}

	if err := gc.genreService.DeleteGenre(uint(id)); err != nil {
		if err.Error() == "genre not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "cannot delete genre with associated audiobooks" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Genre deleted successfully"})
}

// GetGenresByIDs retrieves multiple genres by their IDs
func (gc *GenreController) GetGenresByIDs(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	genres, err := gc.genreService.GetGenresByIDs(req.IDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, genres)
}
