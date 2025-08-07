package routes

import (
	"YufungProject/configs"
	"YufungProject/internal/controller"
	"YufungProject/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupUserRoutes 设置用户管理相关路由
func SetupUserRoutes(router *gin.Engine, userController *controller.UserController, config *configs.Config) {
	// 用户管理路由组 - 需要认证
	userGroup := router.Group("/api/v1/users")
	{
		// 添加JWT认证中间件
		userGroup.Use(middleware.JWTAuthMiddleware(config))
		// TODO: 暂时注释管理员权限验证，方便测试功能
		// userGroup.Use(middleware.AdminRequiredMiddleware())

		// 用户CRUD操作
		userGroup.POST("", userController.CreateUser)       // 创建用户
		userGroup.GET("", userController.GetUserList)       // 获取用户列表
		userGroup.GET("/:id", userController.GetUserByID)   // 获取用户详情
		userGroup.PUT("/:id", userController.UpdateUser)    // 更新用户信息
		userGroup.DELETE("/:id", userController.DeleteUser) // 删除用户

		// 用户密码管理
		userGroup.PUT("/:id/reset-password", userController.ResetUserPassword) // 重置用户密码

		// 批量操作
		userGroup.PUT("/batch-status", userController.BatchUpdateUserStatus) // 批量更新用户状态

		// 快捷操作
		userGroup.PUT("/:id/quick-disable", userController.QuickDisableUser) // 快捷停用用户

		// 数据导出
		userGroup.GET("/export", userController.ExportUsers) // 导出用户数据

		// 高级导入导出功能
		userGroup.POST("/export-advanced", userController.ExportUsersAdvanced) // 高级导出
		userGroup.GET("/template", userController.DownloadUserTemplate)        // 下载模板
		userGroup.POST("/import/preview", userController.PreviewUserImport)    // 预览导入
		userGroup.POST("/import", userController.ImportUsers)                  // 导入用户
	}
}
