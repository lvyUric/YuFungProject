package routes

import (
	"YufungProject/configs"
	"YufungProject/internal/controller"
	"YufungProject/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupCompanyRoutes 设置公司管理相关路由
func SetupCompanyRoutes(router *gin.Engine, companyController *controller.CompanyController, config *configs.Config) {
	// 需要认证的路由
	companyGroup := router.Group("/api/company")
	companyGroup.Use(middleware.AuthMiddleware(config))
	{
		// 公司统计（放在参数路由前面，避免被 :id 匹配）
		companyGroup.GET("/stats", companyController.GetCompanyStats)

		// 导入导出功能
		companyGroup.POST("/export", companyController.ExportCompany)         // 导出公司数据
		companyGroup.GET("/template", companyController.DownloadTemplate)     // 下载导入模板
		companyGroup.POST("/import/preview", companyController.PreviewImport) // 预览导入数据
		companyGroup.POST("/import", companyController.ImportCompany)         // 导入公司数据

		// 公司基本操作
		companyGroup.POST("", companyController.CreateCompany)       // 创建公司
		companyGroup.GET("", companyController.GetCompanyList)       // 获取公司列表
		companyGroup.GET("/:id", companyController.GetCompanyByID)   // 获取公司详情
		companyGroup.PUT("/:id", companyController.UpdateCompany)    // 更新公司
		companyGroup.DELETE("/:id", companyController.DeleteCompany) // 删除公司
	}
}
