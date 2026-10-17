package routes

import (
	adminControllers "gin-skeleton/controllers/admin"
	"gin-skeleton/middleware"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(router *gin.Engine) {
	admin := router.Group("/admin")

	userController := &adminControllers.UserController{}
	roleController := &adminControllers.RoleController{}

	admin.GET("/roles", middleware.AuthMiddleware, middleware.AdminOnly, roleController.GetData)

	admin.GET("/users", middleware.AuthMiddleware, middleware.AdminOnly, userController.GetData)
	admin.POST("/user/action", middleware.AuthMiddleware, middleware.AdminOnly, userController.ActionData)
	admin.GET("/user/:id", middleware.AuthMiddleware, middleware.AdminOnly, userController.GetDetail)
	admin.GET("/user/delete/:id", middleware.AuthMiddleware, middleware.AdminOnly, userController.DeleteData)
}
