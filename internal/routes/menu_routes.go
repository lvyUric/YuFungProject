package routes

import (
	"YufungProject/configs"
	"YufungProject/internal/controller"
	"YufungProject/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupMenuRoutes 设置菜单管理相关路由
func SetupMenuRoutes(router *gin.Engine, menuController *controller.MenuController, config *configs.Config) {
	// 创建认证中间件
	authMiddleware := middleware.JWTAuthMiddleware(config)

	// API v1 路由组
	v1 := router.Group("/api/v1")
	v1.Use(authMiddleware) // 所有菜单相关接口都需要认证

	// 菜单管理路由
	menus := v1.Group("/menus")
	{
		// 基本CRUD操作
		menus.POST("", menuController.CreateMenu)       // 创建菜单
		menus.GET("", menuController.GetMenuList)       // 获取菜单列表
		menus.GET("/tree", menuController.GetMenuTree)  // 获取菜单树
		menus.GET("/user", menuController.GetUserMenus) // 获取用户菜单
		menus.GET("/:id", menuController.GetMenuByID)   // 根据ID获取菜单
		menus.PUT("/:id", menuController.UpdateMenu)    // 更新菜单
		menus.DELETE("/:id", menuController.DeleteMenu) // 删除菜单

		// 批量操作
		menus.PUT("/batch-status", menuController.BatchUpdateMenuStatus) // 批量更新菜单状态

		// 统计信息
		menus.GET("/stats", menuController.GetMenuStats) // 获取菜单统计信息
	}

	// 单独的菜单树路由（适配前端调用）
	v1.GET("/menu/tree", menuController.GetMenuTree)
}
