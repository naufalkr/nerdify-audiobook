package controller

import (
	"microservice/user/domain-layer/service"
	"microservice/user/helpers/dto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RoleController handles HTTP requests related to role management
type RoleController struct {
	roleService service.RoleService
	userService *service.UserService
}

// NewRoleController creates a new RoleController with dependency injection
func NewRoleController(roleService service.RoleService, userService *service.UserService) *RoleController {
	return &RoleController{
		roleService: roleService,
		userService: userService,
	}
}

// CreateRole creates a new role
func (c *RoleController) CreateRole(ctx *gin.Context) {
	var request dto.RoleRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if role with the same name already exists
	exists, err := c.roleService.ExistsByName(ctx, request.Name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check role existence"})
		return
	}
	if exists {
		ctx.JSON(http.StatusConflict, gin.H{"error": "Role with this name already exists"})
		return
	}

	role, err := c.roleService.CreateRole(ctx, request.Name, request.Description, request.IsSystem)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create role"})
		return
	}

	response := dto.RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		IsSystem:    role.IsSystem,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}

	ctx.JSON(http.StatusCreated, response)
}

// GetRoleByName retrieves a role by name
func (c *RoleController) GetRoleByName(ctx *gin.Context) {
	name := ctx.Param("name")

	role, err := c.roleService.FindByName(ctx, name)
	if err != nil {
		if err.Error() == "role not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve role"})
		return
	}

	response := dto.RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		IsSystem:    role.IsSystem,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, response)
}

// GetRoleByID retrieves a role by ID
func (c *RoleController) GetRoleByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID format"})
		return
	}

	role, err := c.roleService.FindByID(ctx, id)
	if err != nil {
		if err.Error() == "role not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve role"})
		return
	}

	response := dto.RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		IsSystem:    role.IsSystem,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, response)
}

// ListAllRoles retrieves all roles with optional pagination and search
func (c *RoleController) ListAllRoles(ctx *gin.Context) {
	// Check if pagination is requested
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("pageSize", "10")
	query := ctx.DefaultQuery("query", "")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	// If pagination or search is requested, use SearchRoles
	if query != "" || ctx.Query("page") != "" || ctx.Query("pageSize") != "" {
		roles, totalCount, err := c.roleService.SearchRoles(ctx, query, page, pageSize)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve roles"})
			return
		}

		// Convert to response
		roleResponses := make([]dto.RoleResponse, len(roles))
		for i, role := range roles {
			roleResponses[i] = dto.RoleResponse{
				ID:          role.ID,
				Name:        role.Name,
				Description: role.Description,
				IsSystem:    role.IsSystem,
				CreatedAt:   role.CreatedAt,
				UpdatedAt:   role.UpdatedAt,
			}
		}

		response := dto.RoleListResponse{
			Roles:      roleResponses,
			TotalCount: totalCount,
			Page:       page,
			PageSize:   pageSize,
		}

		ctx.JSON(http.StatusOK, response)
		return
	}

	// If no pagination or search, use ListAll
	roles, err := c.roleService.ListAll(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve roles"})
		return
	}

	// Convert to response
	roleResponses := make([]dto.RoleResponse, len(roles))
	for i, role := range roles {
		roleResponses[i] = dto.RoleResponse{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
			IsSystem:    role.IsSystem,
			CreatedAt:   role.CreatedAt,
			UpdatedAt:   role.UpdatedAt,
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"roles": roleResponses, "totalCount": len(roles)})
}

// UpdateRole updates an existing role
func (c *RoleController) UpdateRole(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID format"})
		return
	}

	var request dto.RoleRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the role exists
	_, err = c.roleService.FindByID(ctx, id)
	if err != nil {
		if err.Error() == "role not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check role existence"})
		return
	}

	// Update the role
	updatedRole, err := c.roleService.UpdateRole(ctx, id, request.Name, request.Description)
	if err != nil {
		if err.Error() == "cannot change the name of a system role" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update role"})
		return
	}

	response := dto.RoleResponse{
		ID:          updatedRole.ID,
		Name:        updatedRole.Name,
		Description: updatedRole.Description,
		IsSystem:    updatedRole.IsSystem,
		CreatedAt:   updatedRole.CreatedAt,
		UpdatedAt:   updatedRole.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, response)
}

// DeleteRole deletes a role by ID
func (c *RoleController) DeleteRole(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID format"})
		return
	}

	err = c.roleService.DeleteByID(ctx, id)
	if err != nil {
		if err.Error() == "role not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
			return
		} else if err.Error() == "cannot delete a system-defined role" {
			ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete role"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Role successfully deleted"})
}

// SeedDefaultRoles creates default system roles
func (c *RoleController) SeedDefaultRoles(ctx *gin.Context) {
	err := c.roleService.SeedDefaultRoles(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to seed default roles"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Default roles successfully created"})
}

// GetSystemRoles retrieves all system-defined roles
func (c *RoleController) GetSystemRoles(ctx *gin.Context) {
	roles, err := c.roleService.GetSystemRoles(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve system roles"})
		return
	}

	// Convert to response
	roleResponses := make([]dto.RoleResponse, len(roles))
	for i, role := range roles {
		roleResponses[i] = dto.RoleResponse{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
			IsSystem:    role.IsSystem,
			CreatedAt:   role.CreatedAt,
			UpdatedAt:   role.UpdatedAt,
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"roles": roleResponses})
}

// BulkDeleteRoles deletes multiple roles at once
func (c *RoleController) BulkDeleteRoles(ctx *gin.Context) {
	var request dto.BulkDeleteRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the request contains valid IDs
	if len(request.IDs) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No role IDs provided"})
		return
	}

	// Delete the roles
	deletedCount, err := c.roleService.BulkDeleteRoles(ctx, request.IDs)
	if err != nil {
		if err.Error() == "cannot delete system-defined roles" {
			ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete roles"})
		return
	}

	response := dto.BulkDeleteResponse{
		DeletedCount: deletedCount,
	}

	ctx.JSON(http.StatusOK, response)
}
