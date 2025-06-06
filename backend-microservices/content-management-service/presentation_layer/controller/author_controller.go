package controller

import (
	"content-management-service/data_layer/dto"
	"content-management-service/domain_layer/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AuthorController struct {
	authorService *service.AuthorService
}

func NewAuthorController(authorService *service.AuthorService) *AuthorController {
	return &AuthorController{
		authorService: authorService,
	}
}

// CreateAuthor creates a new author
// @Summary Create a new author
// @Description Create a new author with the provided data
// @Tags authors
// @Accept json
// @Produce json
// @Param request body dto.CreateAuthorRequest true "Author data"
// @Success 201 {object} dto.APIResponse
// @Failure 400 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /authors [post]
func (c *AuthorController) CreateAuthor(ctx *gin.Context) {
	var req dto.CreateAuthorRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid request data",
			Error:   err.Error(),
		})
		return
	}

	author, err := c.authorService.CreateAuthor(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: "Failed to create author",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, dto.APIResponse{
		Success: true,
		Message: "Author created successfully",
		Data:    author,
	})
}

// GetAuthor retrieves an author by ID
// @Summary Get author by ID
// @Description Get author details by ID
// @Tags authors
// @Produce json
// @Param id path int true "Author ID"
// @Success 200 {object} dto.APIResponse
// @Failure 400 {object} dto.APIResponse
// @Failure 404 {object} dto.APIResponse
// @Router /authors/{id} [get]
func (c *AuthorController) GetAuthor(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid author ID",
			Error:   err.Error(),
		})
		return
	}

	author, err := c.authorService.GetAuthorByID(uint(id))
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "author not found" {
			statusCode = http.StatusNotFound
		}
		ctx.JSON(statusCode, dto.APIResponse{
			Success: false,
			Message: "Failed to get author",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Author retrieved successfully",
		Data:    author,
	})
}

// GetAllAuthors retrieves all authors with pagination
// @Summary Get all authors
// @Description Get all authors with pagination
// @Tags authors
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} dto.APIResponse
// @Failure 400 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /authors [get]
func (c *AuthorController) GetAllAuthors(ctx *gin.Context) {
	var req dto.PaginationRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid query parameters",
			Error:   err.Error(),
		})
		return
	}

	result, err := c.authorService.GetAllAuthors(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: "Failed to get authors",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Authors retrieved successfully",
		Data:    result,
	})
}

// UpdateAuthor updates an existing author
// @Summary Update author
// @Description Update author data by ID
// @Tags authors
// @Accept json
// @Produce json
// @Param id path int true "Author ID"
// @Param request body dto.UpdateAuthorRequest true "Author data"
// @Success 200 {object} dto.APIResponse
// @Failure 400 {object} dto.APIResponse
// @Failure 404 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /authors/{id} [put]
func (c *AuthorController) UpdateAuthor(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid author ID",
			Error:   err.Error(),
		})
		return
	}

	var req dto.UpdateAuthorRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid request data",
			Error:   err.Error(),
		})
		return
	}

	author, err := c.authorService.UpdateAuthor(uint(id), req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "author not found" {
			statusCode = http.StatusNotFound
		}
		ctx.JSON(statusCode, dto.APIResponse{
			Success: false,
			Message: "Failed to update author",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Author updated successfully",
		Data:    author,
	})
}

// DeleteAuthor deletes an author
// @Summary Delete author
// @Description Delete author by ID
// @Tags authors
// @Produce json
// @Param id path int true "Author ID"
// @Success 200 {object} dto.APIResponse
// @Failure 400 {object} dto.APIResponse
// @Failure 404 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /authors/{id} [delete]
func (c *AuthorController) DeleteAuthor(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid author ID",
			Error:   err.Error(),
		})
		return
	}

	err = c.authorService.DeleteAuthor(uint(id))
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "author not found" {
			statusCode = http.StatusNotFound
		}
		ctx.JSON(statusCode, dto.APIResponse{
			Success: false,
			Message: "Failed to delete author",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Author deleted successfully",
	})
}

// SearchAuthors searches authors by name
// @Summary Search authors
// @Description Search authors by name
// @Tags authors
// @Produce json
// @Param q query string false "Search query"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} dto.APIResponse
// @Failure 400 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /authors/search [get]
func (c *AuthorController) SearchAuthors(ctx *gin.Context) {
	var req dto.SearchRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid query parameters",
			Error:   err.Error(),
		})
		return
	}

	result, err := c.authorService.SearchAuthors(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: "Failed to search authors",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Authors search completed successfully",
		Data:    result,
	})
}
