package controller

import (
	"catalog-service/data_layer/dto"
	"catalog-service/domain_layer/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReaderController struct {
	readerService *service.ReaderService
}

func NewReaderController(readerService *service.ReaderService) *ReaderController {
	return &ReaderController{
		readerService: readerService,
	}
}

// CreateReader creates a new reader
func (rc *ReaderController) CreateReader(c *gin.Context) {
	var req dto.CreateReaderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reader, err := rc.readerService.CreateReader(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, reader)
}

// GetReaderByID retrieves a reader by ID
func (rc *ReaderController) GetReaderByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reader ID"})
		return
	}

	reader, err := rc.readerService.GetReaderByID(uint(id))
	if err != nil {
		if err.Error() == "reader not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reader)
}

// GetAllReaders retrieves all readers with pagination
func (rc *ReaderController) GetAllReaders(c *gin.Context) {
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

	readers, err := rc.readerService.GetAllReaders(paginationReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, readers)
}

// UpdateReader updates an existing reader
func (rc *ReaderController) UpdateReader(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reader ID"})
		return
	}

	var req dto.UpdateReaderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reader, err := rc.readerService.UpdateReader(uint(id), req)
	if err != nil {
		if err.Error() == "reader not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reader)
}

// DeleteReader deletes a reader
func (rc *ReaderController) DeleteReader(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reader ID"})
		return
	}

	if err := rc.readerService.DeleteReader(uint(id)); err != nil {
		if err.Error() == "reader not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reader deleted successfully"})
}

// SearchReaders searches readers by name
func (rc *ReaderController) SearchReaders(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name parameter is required"})
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

	readers, err := rc.readerService.SearchReadersByName(name, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, readers)
}
