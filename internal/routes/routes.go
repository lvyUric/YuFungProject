package routes

import (
	"YufungProject/configs"
	"YufungProject/internal/controller"
	"YufungProject/internal/middleware"
	"YufungProject/internal/repository"
	"YufungProject/internal/service"
	"YufungProject/pkg/logger"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
)

// SetupRoutes 设置所有路由
func SetupRoutes(db *mongo.Database, config *configs.Config) *gin.Engine {
	// 设置Gin运行模式
	gin.SetMode(config.Server.Mode)

	// 创建Gin实例
	router := gin.New()

	// 添加日志中间件
	router.Use(LoggerMiddleware())

	// 添加恢复中间件
	router.Use(gin.Recovery())

	// 添加CORS中间件
	router.Use(middleware.CORSMiddleware())

	// 添加全局OPTIONS处理
	router.OPTIONS("/*path", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Length, Content-Type, Authorization, Accept, X-Requested-With")
		c.Header("Access-Control-Max-Age", "43200")
		c.Status(204)
	})

	// 健康检查接口
	router.GET("/health", func(c *gin.Context) {
		logger.Info("健康检查请求")
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "保险经纪管理系统",
			"version": "1.0.0",
		})
	})

	// Swagger文档路由 - 只在debug模式下启用
	if config.Server.Mode == "debug" {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		logger.Info("Swagger文档已启用: /swagger/index.html")
	}

	// 初始化仓库层
	userRepo := repository.NewUserRepository(db)
	companyRepo := repository.NewCompanyRepository(db, userRepo)
	roleRepo := repository.NewRoleRepository(db)
	menuRepo := repository.NewMenuRepository(db)
	rbacRepo := repository.NewRBACRepository(db)                 // 启用RBAC仓库
	policyRepo := repository.NewPolicyRepository(db)             // 添加保单仓库
	systemConfigRepo := repository.NewSystemConfigRepository(db) // 添加系统配置仓库
	changeRecordRepo := repository.NewChangeRecordRepository(db) // 添加变更记录仓库

	// 初始化服务层
	authService := service.NewAuthService(userRepo, config)
	companyService := service.NewCompanyService(companyRepo, userRepo)
	userService := service.NewUserService(userRepo, companyRepo)
	roleService := service.NewRoleService(roleRepo, companyRepo, rbacRepo)
	menuService := service.NewMenuService(menuRepo)
	changeRecordService := service.NewChangeRecordService(changeRecordRepo, userRepo) // 添加变更记录服务
	policyService := service.NewPolicyService(policyRepo, changeRecordService)        // 修改保单服务注入变更记录服务
	systemConfigService := service.NewSystemConfigService(systemConfigRepo)           // 添加系统配置服务

	// 初始化控制器层
	authController := controller.NewAuthController(authService)
	companyController := controller.NewCompanyController(companyService)
	userController := controller.NewUserController(userService, companyService)
	roleController := controller.NewRoleController(roleService)
	menuController := controller.NewMenuController(menuService)
	policyController := controller.NewPolicyController(policyService)                   // 添加保单控制器
	systemConfigController := controller.NewSystemConfigController(systemConfigService) // 添加系统配置控制器
	changeRecordController := controller.NewChangeRecordController(changeRecordService) // 添加变更记录控制器

	// 设置认证相关路由
	SetupAuthRoutes(router, authController, config)

	// 设置公司管理相关路由
	SetupCompanyRoutes(router, companyController, config)

	// 设置用户管理相关路由
	SetupUserRoutes(router, userController, config)

	// 设置角色管理相关路由
	SetupRoleRoutes(router, roleController, config)

	// 设置菜单管理相关路由
	SetupMenuRoutes(router, menuController, config)

	// 设置保单管理相关路由
	SetupPolicyRoutes(router, policyController, changeRecordController, config)

	// 设置变更记录相关路由
	SetupChangeRecordRoutes(router, changeRecordController, config)

	// 设置系统配置相关路由
	api := router.Group("/api")
	RegisterSystemConfigRoutes(api, systemConfigController, config)

	logger.Info("所有路由设置完成")
	return router
}

// LoggerMiddleware 自定义日志中间件
func LoggerMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// 记录API访问日志
		logger.APILog(
			param.Method,
			param.Path,
			param.ClientIP,
			param.Request.UserAgent(),
			param.StatusCode,
			param.Latency,
		)
		return ""
	})
}
