package middleware

import (
	"net/http"
	"strings"
	"time"

	"YufungProject/configs"
	"YufungProject/internal/model"
	"YufungProject/pkg/logger"
	"YufungProject/pkg/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware JWT认证中间件
func AuthMiddleware(config *configs.Config) gin.HandlerFunc {
	// 解析时间配置
	expiresIn, _ := time.ParseDuration(config.JWT.ExpiresIn)
	refreshExpiresIn, _ := time.ParseDuration(config.JWT.RefreshExpiresIn)

	jwtUtil := utils.NewJWTUtil(
		config.JWT.Secret,
		expiresIn,
		refreshExpiresIn,
	)

	return func(c *gin.Context) {
		// 获取Authorization头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Warnf("认证失败 - 缺少Authorization头: %s", c.ClientIP())
			c.JSON(http.StatusUnauthorized, model.ErrorResponse(model.CodeUnauthorized, "未提供认证令牌", nil))
			c.Abort()
			return
		}

		// 检查Bearer前缀
		tokenParts := strings.SplitN(authHeader, " ", 2)
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			logger.Warnf("认证失败 - 令牌格式错误: %s", c.ClientIP())
			c.JSON(http.StatusUnauthorized, model.ErrorResponse(model.CodeTokenInvalid, "令牌格式错误", nil))
			c.Abort()
			return
		}

		// 解析令牌
		claims, err := jwtUtil.ParseToken(tokenParts[1])
		if err != nil {
			logger.Warnf("认证失败 - 令牌解析错误: %v, IP: %s", err, c.ClientIP())
			c.JSON(http.StatusUnauthorized, model.ErrorResponse(model.CodeTokenInvalid, "令牌无效或已过期", nil))
			c.Abort()
			return
		}

		// 记录认证成功日志
		logger.Debugf("用户认证成功: UserID=%s, Username=%s, IP=%s", claims.UserID, claims.Username, c.ClientIP())

		// 将用户信息存储在上下文中
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("company_id", claims.CompanyID)
		c.Set("role_ids", claims.RoleIDs)

		c.Next()
	}
}

// OptionalAuthMiddleware 可选认证中间件（用于某些不强制登录的接口）
func OptionalAuthMiddleware(config *configs.Config) gin.HandlerFunc {
	// 解析时间配置
	expiresIn, _ := time.ParseDuration(config.JWT.ExpiresIn)
	refreshExpiresIn, _ := time.ParseDuration(config.JWT.RefreshExpiresIn)

	jwtUtil := utils.NewJWTUtil(
		config.JWT.Secret,
		expiresIn,
		refreshExpiresIn,
	)

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			tokenParts := strings.SplitN(authHeader, " ", 2)
			if len(tokenParts) == 2 && tokenParts[0] == "Bearer" {
				claims, err := jwtUtil.ParseToken(tokenParts[1])
				if err == nil {
					c.Set("user_id", claims.UserID)
					c.Set("username", claims.Username)
					c.Set("company_id", claims.CompanyID)
					c.Set("role_ids", claims.RoleIDs)
					logger.Debugf("可选认证成功: UserID=%s, Username=%s", claims.UserID, claims.Username)
				} else {
					logger.Debugf("可选认证失败，但继续处理: %v", err)
				}
			}
		}
		c.Next()
	}
}

// JWTAuthMiddleware JWT认证中间件（别名，兼容性）
func JWTAuthMiddleware(config *configs.Config) gin.HandlerFunc {
	return AuthMiddleware(config)
}

// AdminRequiredMiddleware 管理员权限验证中间件
func AdminRequiredMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户角色信息
		roleIDs, exists := GetRoleIDs(c)
		if !exists {
			logger.Warnf("权限验证失败 - 无法获取用户角色信息: IP=%s", c.ClientIP())
			c.JSON(http.StatusForbidden, model.ErrorResponse(model.CodePermissionDeny, "无法获取用户权限信息", nil))
			c.Abort()
			return
		}

		// 检查是否为平台管理员（超级管理员）
		// 假设超级管理员的角色ID为 "ADMIN" 或 "SUPER_ADMIN"
		isAdmin := false
		for _, roleID := range roleIDs {
			if roleID == "ADMIN" || roleID == "SUPER_ADMIN" || roleID == "platform_admin" {
				isAdmin = true
				break
			}
		}

		if !isAdmin {
			userID, _ := GetUserID(c)
			logger.Warnf("权限验证失败 - 非管理员用户: UserID=%s, RoleIDs=%v, IP=%s", userID, roleIDs, c.ClientIP())
			c.JSON(http.StatusForbidden, model.ErrorResponse(model.CodePermissionDeny, "需要管理员权限", nil))
			c.Abort()
			return
		}

		// 管理员权限验证通过
		userID, _ := GetUserID(c)
		logger.Debugf("管理员权限验证通过: UserID=%s, RoleIDs=%v", userID, roleIDs)
		c.Next()
	}
}

// GetUserID 从上下文获取用户ID
func GetUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return "", false
	}
	return userID.(string), true
}

// GetUsername 从上下文获取用户名
func GetUsername(c *gin.Context) (string, bool) {
	username, exists := c.Get("username")
	if !exists {
		return "", false
	}
	return username.(string), true
}

// GetCompanyID 从上下文获取公司ID
func GetCompanyID(c *gin.Context) (string, bool) {
	companyID, exists := c.Get("company_id")
	if !exists {
		return "", false
	}
	return companyID.(string), true
}

// GetRoleIDs 从上下文获取角色ID列表
func GetRoleIDs(c *gin.Context) ([]string, bool) {
	roleIDs, exists := c.Get("role_ids")
	if !exists {
		return nil, false
	}
	return roleIDs.([]string), true
}
