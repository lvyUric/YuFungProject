package routes

import (
	"YufungProject/configs"
	"YufungProject/internal/controller"
	"YufungProject/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupActivityLogRoutes 设置活动记录相关路由
func SetupActivityLogRoutes(router *gin.Engine, activityLogController *controller.ActivityLogController, config *configs.Config) {
	// 活动记录路由组 - 需要认证
	activityLogGroup := router.Group("/api/v1/activity-logs")
	{
		// 添加JWT认证中间件
		activityLogGroup.Use(middleware.JWTAuthMiddleware(config))

		// 获取活动记录列表
		activityLogGroup.GET("", activityLogController.GetActivityLogList)

		// 获取最近的活动记录（用于仪表盘）
		activityLogGroup.GET("/recent", activityLogController.GetRecentActivityLogs)

		// 获取活动记录详情
		activityLogGroup.GET("/:id", activityLogController.GetActivityLogByID)

		// 获取活动记录统计
		activityLogGroup.GET("/statistics", activityLogController.GetActivityLogStatistics)

		// 删除指定公司的活动记录（仅平台管理员）
		activityLogGroup.DELETE("/company/:company_id", activityLogController.DeleteActivityLogsByCompanyID)
	}
}
