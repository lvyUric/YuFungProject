package routes

import (
	"github.com/gin-gonic/gin"

	"YufungProject/configs"
	"YufungProject/internal/controller"
	"YufungProject/internal/middleware"
)

// RegisterSystemConfigRoutes 注册系统配置相关路由
func RegisterSystemConfigRoutes(r *gin.RouterGroup, systemConfigController *controller.SystemConfigController, config *configs.Config) {
	// 系统配置管理路由组
	systemConfigGroup := r.Group("/system-configs")
	systemConfigGroup.Use(middleware.AuthMiddleware(config)) // 需要认证

	{
		systemConfigGroup.GET("", systemConfigController.ListSystemConfigs)         // 获取系统配置列表
		systemConfigGroup.POST("", systemConfigController.CreateSystemConfig)       // 创建系统配置
		systemConfigGroup.GET("/:id", systemConfigController.GetSystemConfig)       // 获取系统配置详情
		systemConfigGroup.PUT("/:id", systemConfigController.UpdateSystemConfig)    // 更新系统配置
		systemConfigGroup.DELETE("/:id", systemConfigController.DeleteSystemConfig) // 删除系统配置

		// 获取配置选项
		systemConfigGroup.GET("/options/:type", systemConfigController.GetConfigOptions) // 根据类型获取配置选项
	}
}
