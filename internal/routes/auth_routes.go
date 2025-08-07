package routes

import (
	"YufungProject/configs"
	"YufungProject/internal/controller"
	"YufungProject/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupAuthRoutes 设置认证相关路由
func SetupAuthRoutes(router *gin.Engine, authController *controller.AuthController, config *configs.Config) {
	// 公开路由（不需要认证）
	authGroup := router.Group("/api/auth")
	{
		authGroup.POST("/login", authController.Login)
		authGroup.POST("/register", authController.Register)
		authGroup.POST("/refresh", authController.RefreshToken)
	}

	// 需要认证的路由
	authProtectedGroup := router.Group("/api/auth")
	authProtectedGroup.Use(middleware.AuthMiddleware(config))
	{
		authProtectedGroup.POST("/logout", authController.Logout)
		authProtectedGroup.POST("/change-password", authController.ChangePassword)
		authProtectedGroup.GET("/user-info", authController.GetUserInfo)
	}
}
