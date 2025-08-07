package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"YufungProject/configs"
	// _ "YufungProject/docs" // 导入swagger文档 - 临时注释
	"YufungProject/internal/routes"
	"YufungProject/pkg/database"
	"YufungProject/pkg/logger"
	customValidator "YufungProject/pkg/validator"

	"github.com/gin-gonic/gin"
)

// @title 保险经纪管理系统API
// @version 1.0
// @description 基于Gin + MongoDB构建的保险经纪公司管理平台
func main() {
	// 1. 加载配置
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 2. 初始化日志系统
	if err := logger.InitLogger(logger.LogConfig{
		Level:      config.Log.Level,
		Format:     config.Log.Format,
		Output:     config.Log.Output,
		FilePath:   config.Log.FilePath,
		MaxSize:    config.Log.MaxSize,
		MaxAge:     config.Log.MaxAge,
		MaxBackups: config.Log.MaxBackups,
		Compress:   config.Log.Compress,
	}); err != nil {
		log.Fatalf("初始化日志系统失败: %v", err)
	}

	logger.Info("🏥 保险经纪管理系统启动中...")
	logger.Infof("配置信息: 服务端口=%d, 运行模式=%s, 日志级别=%s",
		config.Server.Port, config.Server.Mode, config.Log.Level)

	// 3. 初始化数据库
	logger.Info("初始化MongoDB数据库连接...")
	db, err := database.InitMongoDB(config.Database.MongoDB)
	if err != nil {
		logger.Fatalf("数据库连接失败: %v", err)
	}
	logger.Info("✅ MongoDB数据库连接成功")

	// 4. 初始化自定义验证器
	logger.Info("初始化自定义验证器...")
	if err := customValidator.InitCustomValidators(); err != nil {
		logger.Fatalf("初始化自定义验证器失败: %v", err)
	}
	logger.Info("✅ 自定义验证器初始化成功")

	// 设置Gin运行模式
	gin.SetMode(config.Server.Mode)

	// 5. 初始化路由
	logger.Info("初始化应用路由...")
	router := routes.SetupRoutes(db, config)
	logger.Info("✅ 路由初始化完成")

	// 6. 创建HTTP服务器
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Server.Port),
		Handler: router,
	}

	// 7. 启动服务器
	go func() {
		logger.Infof("🚀 服务器启动成功")
		logger.Infof("📍 服务地址: http://localhost:%d", config.Server.Port)
		logger.Infof("📖 API文档: http://localhost:%d/swagger/index.html", config.Server.Port)
		logger.Infof("🔍 健康检查: http://localhost:%d/health", config.Server.Port)
		logger.Info("按 Ctrl+C 停止服务")

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("服务器启动失败: %v", err)
		}
	}()

	// 8. 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("🛑 正在关闭服务器...")

	// 设置关闭超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 关闭HTTP服务器
	if err := srv.Shutdown(ctx); err != nil {
		logger.Errorf("服务器强制关闭: %v", err)
	}

	// 关闭数据库连接
	if err := database.DisconnectMongoDB(); err != nil {
		logger.Errorf("关闭数据库连接失败: %v", err)
	}

	logger.Info("✅ 服务器已安全关闭")
}
