package routes

import (
	"gin-skeleton/controllers"
	userControllers "gin-skeleton/controllers/user"
	"gin-skeleton/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	user := router.Group("/")

	authController := &controllers.AuthController{}
	profileController := &userControllers.ProfileController{}

	user.POST("/login", authController.Login)
	user.POST("/register", authController.Register)
	user.POST("/verify", authController.Verify)
	user.POST("/retry-code", authController.RetryCode)
	user.POST("/reset-password", authController.ResetPassword)

	user.GET("/check-auth", middleware.AuthMiddleware, profileController.GetProfile)
	user.POST("/profile/action", middleware.AuthMiddleware, profileController.Action)
}
