package routes

import (
	"microservice/user/domain-layer/middleware"
	"microservice/user/presentation-layer/controller"
	"reflect"

	"github.com/gin-gonic/gin"
)

// getMethodIfExists is a helper to check if a method exists on a controller
func getMethodIfExists(controller interface{}, methodName string) interface{} {
	controllerValue := reflect.ValueOf(controller)
	method := controllerValue.MethodByName(methodName)
	if method.IsValid() {
		return method.Interface()
	}
	return nil
}

// RoleRoutes sets up all role-related routes
func RoleRoutes(router *gin.Engine, controller *controller.RoleController) {
	// Basic role management routes
	roleRoutes := router.Group("/roles")
	roleRoutes.Use(middleware.RequireSuperAdmin()) // Semua endpoint hanya dapat diakses oleh superadmin
	{
		roleRoutes.POST("/", controller.CreateRole)
		roleRoutes.GET("/name/:name", controller.GetRoleByName)
		roleRoutes.GET("/id/:id", controller.GetRoleByID)
		roleRoutes.GET("/", controller.ListAllRoles)
		roleRoutes.POST("/seed", controller.SeedDefaultRoles)
		roleRoutes.PUT("/:id", controller.UpdateRole)
		roleRoutes.GET("/system", controller.GetSystemRoles)
	}

	// Add routes for changing user roles
	changeRoleMethod := getMethodIfExists(controller, "ChangeUserRole")
	if changeRoleMethod != nil {
		// For API path without v1
		adminRoutes := router.Group("/api/admin/users")
		adminRoutes.Use(middleware.RequireSuperAdmin())
		{
			adminRoutes.PUT("/:id/role", changeRoleMethod.(func(*gin.Context)))
		}
	}
}
