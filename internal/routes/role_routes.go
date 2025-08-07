package routes

import (
	"YufungProject/configs"
	"YufungProject/internal/controller"
	"YufungProject/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoleRoutes 设置角色管理相关路由
func SetupRoleRoutes(router *gin.Engine, roleController *controller.RoleController, config *configs.Config) {
	// 创建认证中间件
	authMiddleware := middleware.JWTAuthMiddleware(config)

	// API v1 路由组
	v1 := router.Group("/api/v1")
	v1.Use(authMiddleware) // 所有角色相关接口都需要认证

	// 角色管理路由
	roles := v1.Group("/roles")
	{
		// 基本CRUD操作
		roles.POST("", roleController.CreateRole)       // 创建角色
		roles.GET("", roleController.GetRoleList)       // 获取角色列表
		roles.GET("/:id", roleController.GetRoleByID)   // 根据ID获取角色
		roles.PUT("/:id", roleController.UpdateRole)    // 更新角色
		roles.DELETE("/:id", roleController.DeleteRole) // 删除角色

		// 批量操作
		roles.PUT("/batch-status", roleController.BatchUpdateRoleStatus) // 批量更新角色状态

		// 统计信息
		roles.GET("/stats", roleController.GetRoleStats) // 获取角色统计信息

		// 公司角色
		roles.GET("/company/:company_id", roleController.GetRolesByCompanyID) // 根据公司ID获取角色列表
	}
}
