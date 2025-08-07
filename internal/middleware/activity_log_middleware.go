package middleware

import (
	"YufungProject/internal/model"
	"YufungProject/internal/service"
	"YufungProject/pkg/logger"
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

// ActivityLogMiddleware 活动记录中间件
func ActivityLogMiddleware(activityLogService *service.ActivityLogService) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// 获取请求信息
		requestURL := c.Request.URL.String()
		requestMethod := c.Request.Method
		ipAddress := c.ClientIP()
		userAgent := c.Request.UserAgent()

		// 获取用户信息（从JWT中间件中获取）
		userID, _ := c.Get("user_id")
		username, _ := c.Get("username")
		companyID, _ := c.Get("company_id")
		companyName, _ := c.Get("company_name")

		// 如果用户信息缺失，跳过记录
		if userID == nil || username == nil || companyID == nil {
			c.Next()
			return
		}

		// 如果company_name为空，使用默认值
		if companyName == nil || companyName.(string) == "" {
			companyName = "未知公司"
		}

		// 跳过不需要记录的操作
		if shouldSkipLogging(requestURL, requestMethod) {
			c.Next()
			return
		}

		// 处理请求
		c.Next()

		// 计算执行时间
		executionTime := time.Since(startTime).Milliseconds()

		// 获取操作结果
		resultStatus := "success"
		if c.Writer.Status() >= 400 {
			resultStatus = "failure"
		}

		// 确定操作类型和模块
		operationType, moduleName, operationDesc := parseOperationInfo(requestURL, requestMethod)

		// 异步记录活动日志
		go func() {
			ctx := context.Background()
			log := &model.ActivityLog{
				UserID:        userID.(string),
				Username:      username.(string),
				CompanyID:     companyID.(string),
				CompanyName:   companyName.(string),
				OperationType: operationType,
				ModuleName:    moduleName,
				OperationDesc: operationDesc,
				RequestURL:    requestURL,
				RequestMethod: requestMethod,
				IPAddress:     ipAddress,
				UserAgent:     userAgent,
				OperationTime: time.Now(),
				ExecutionTime: executionTime,
				ResultStatus:  resultStatus,
			}

			if err := activityLogService.CreateActivityLog(ctx, log); err != nil {
				// 记录错误但不影响主流程
				logger.Errorf("记录活动日志失败: %v", err)
			}
		}()
	}
}

// shouldSkipLogging 判断是否需要跳过日志记录
func shouldSkipLogging(url, method string) bool {
	// 跳过健康检查、静态资源等
	skipPaths := []string{
		"/health",
		"/swagger",
		"/favicon.ico",
		"/static",
		"/assets",
	}

	for _, path := range skipPaths {
		if contains(url, path) {
			return true
		}
	}

	// 跳过OPTIONS请求
	if method == "OPTIONS" {
		return true
	}

	return false
}

// parseOperationInfo 解析操作信息
func parseOperationInfo(url, method string) (operationType, moduleName, operationDesc string) {
	// 根据URL和HTTP方法判断操作类型
	switch method {
	case "GET":
		if contains(url, "/list") || contains(url, "/page") {
			operationType = model.OperationTypeView
			operationDesc = "查看列表"
		} else if contains(url, "/detail") || contains(url, "/info") {
			operationType = model.OperationTypeView
			operationDesc = "查看详情"
		} else if contains(url, "/export") {
			operationType = model.OperationTypeExport
			operationDesc = "导出数据"
		} else {
			operationType = model.OperationTypeView
			operationDesc = "查看"
		}
	case "POST":
		if contains(url, "/create") || contains(url, "/add") {
			operationType = model.OperationTypeCreate
			operationDesc = "新增"
		} else if contains(url, "/import") {
			operationType = model.OperationTypeImport
			operationDesc = "导入数据"
		} else if contains(url, "/login") {
			operationType = model.OperationTypeLogin
			operationDesc = "用户登录"
		} else {
			operationType = model.OperationTypeCreate
			operationDesc = "新增"
		}
	case "PUT", "PATCH":
		operationType = model.OperationTypeUpdate
		operationDesc = "更新"
	case "DELETE":
		operationType = model.OperationTypeDelete
		operationDesc = "删除"
	default:
		operationType = model.OperationTypeView
		operationDesc = "查看"
	}

	// 根据URL判断模块名称
	switch {
	case contains(url, "/users"):
		moduleName = model.ModuleUser
	case contains(url, "/roles"):
		moduleName = model.ModuleRole
	case contains(url, "/menus"):
		moduleName = model.ModuleMenu
	case contains(url, "/companies"):
		moduleName = model.ModuleCompany
	case contains(url, "/policies"):
		moduleName = model.ModulePolicy
	case contains(url, "/customers"):
		moduleName = model.ModuleCustomer
	case contains(url, "/system"):
		moduleName = model.ModuleSystem
	case contains(url, "/auth") || contains(url, "/login") || contains(url, "/logout"):
		moduleName = model.ModuleAuth
	default:
		moduleName = "其他"
	}

	return
}

// contains 检查字符串是否包含子字符串
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) && (s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr ||
			containsSubstring(s, substr))))
}

// containsSubstring 检查字符串中间是否包含子字符串
func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
