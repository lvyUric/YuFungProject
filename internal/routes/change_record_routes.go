package routes

import (
	"YufungProject/configs"
	"YufungProject/internal/controller"
	"YufungProject/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupChangeRecordRoutes 设置变更记录路由
func SetupChangeRecordRoutes(r *gin.Engine, changeRecordController *controller.ChangeRecordController, config *configs.Config) {
	// JWT认证中间件
	authMiddleware := middleware.JWTAuthMiddleware(config)

	// 变更记录路由组，需要认证
	changeRecordGroup := r.Group("/api/change-records")
	changeRecordGroup.Use(authMiddleware)

	// 通用变更记录查询
	changeRecordGroup.GET("", changeRecordController.GetChangeRecordsList) // 获取变更记录列表
}
