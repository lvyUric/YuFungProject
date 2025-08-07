package routes

import (
	"YufungProject/configs"
	"YufungProject/internal/controller"
	"YufungProject/internal/middleware"
	"YufungProject/internal/service"

	"github.com/gin-gonic/gin"
)

// SetupPolicyRoutes 设置保单管理相关路由
func SetupPolicyRoutes(router *gin.Engine, policyController *controller.PolicyController, changeRecordController *controller.ChangeRecordController, config *configs.Config) {
	// 初始化活动记录服务
	activityLogService := service.NewActivityLogService()

	// 需要认证的路由
	policyGroup := router.Group("/api/policies")
	policyGroup.Use(middleware.AuthMiddleware(config))
	policyGroup.Use(middleware.ActivityLogMiddleware(activityLogService)) // 添加活动记录中间件
	{
		// 保单统计（放在参数路由前面，避免被 :id 匹配）
		policyGroup.GET("/statistics", policyController.GetPolicyStatistics)

		// 导入导出功能
		policyGroup.POST("/export", policyController.ExportPolicies)              // 导出保单数据
		policyGroup.GET("/template", policyController.DownloadPolicyTemplate)     // 下载导入模板
		policyGroup.POST("/import/preview", policyController.PreviewPolicyImport) // 预览导入数据
		policyGroup.POST("/import", policyController.ImportPoliciesFromFile)      // 导入保单数据

		// 获取字段验证规则
		policyGroup.GET("/validation-rules", policyController.GetPolicyValidationRules)

		// 保单基本操作
		policyGroup.POST("", policyController.CreatePolicy)       // 创建保单
		policyGroup.GET("", policyController.ListPolicies)        // 获取保单列表
		policyGroup.GET("/:id", policyController.GetPolicy)       // 获取保单详情
		policyGroup.PUT("/:id", policyController.UpdatePolicy)    // 更新保单
		policyGroup.DELETE("/:id", policyController.DeletePolicy) // 删除保单

		// 保单变更记录（使用不同的路径避免冲突）
		policyGroup.GET("/:id/change-records", changeRecordController.GetPolicyChangeRecords) // 获取保单变更记录

		// 批量操作
		policyGroup.POST("/batch-update", policyController.BatchUpdatePolicyStatus) // 批量更新状态
	}
}
