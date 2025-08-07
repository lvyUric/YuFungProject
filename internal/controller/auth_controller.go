package controller

import (
	"net/http"
	"strings"

	"YufungProject/internal/model"
	"YufungProject/internal/service"
	"YufungProject/pkg/logger"

	"github.com/gin-gonic/gin"
)

// AuthController 认证控制器
type AuthController struct {
	authService service.AuthService
}

// NewAuthController 创建认证控制器实例
func NewAuthController(authService service.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

// Login 用户登录
//
//	@Summary		用户登录
//	@Description	使用用户名和密码进行登录认证，返回JWT令牌
//	@Tags			认证管理
//	@Accept			json
//	@Produce		json
//	@Param			request	body		model.LoginRequest		true	"登录请求参数"
//	@Success		200		{object}	model.Response{data=model.LoginResponse}	"登录成功"
//	@Failure		400		{object}	model.Response{data=string}					"请求参数错误"
//	@Failure		401		{object}	model.Response{data=string}					"认证失败"
//	@Failure		423		{object}	model.Response{data=string}					"账户被锁定"
//	@Failure		500		{object}	model.Response{data=string}					"服务器内部错误"
//	@Router			/auth/login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	// 添加详细的请求日志
	logger.Infof("收到登录请求 - Method: %s, Path: %s, ClientIP: %s", ctx.Request.Method, ctx.Request.URL.Path, ctx.ClientIP())
	logger.Infof("请求头信息 - Content-Type: %s, Origin: %s, User-Agent: %s",
		ctx.GetHeader("Content-Type"),
		ctx.GetHeader("Origin"),
		ctx.GetHeader("User-Agent"))

	var req model.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Warnf("登录请求参数错误: %v", err)
		ctx.JSON(http.StatusBadRequest, model.ErrorResponse(model.CodeInvalidParams, "请求参数错误", err.Error()))
		return
	}

	clientIP := ctx.ClientIP()

	// 记录登录尝试
	logger.AuthLog("login_attempt", req.Username, clientIP, false, "开始登录验证")

	loginResp, err := c.authService.Login(ctx, &req)
	if err != nil {
		// 记录登录失败
		logger.AuthLog("login_failed", req.Username, clientIP, false, err.Error())

		switch err.Error() {
		case "用户不存在", "密码错误":
			ctx.JSON(http.StatusUnauthorized, model.ErrorResponse(model.CodeAuthFailed, "用户名或密码错误", nil))
		case "账户已被锁定":
			ctx.JSON(http.StatusLocked, model.ErrorResponse(model.CodeAccountLocked, "账户已被锁定，请稍后再试", nil))
		case "账户已被禁用":
			ctx.JSON(http.StatusForbidden, model.ErrorResponse(model.CodeAccountDisabled, "账户已被禁用", nil))
		default:
			ctx.JSON(http.StatusInternalServerError, model.ErrorResponse(model.CodeServerError, "登录失败", err.Error()))
		}
		return
	}

	// 记录登录成功
	logger.AuthLog("login_success", req.Username, clientIP, true, "登录成功")
	logger.BusinessLog("认证管理", "用户登录", loginResp.User.UserID, "用户登录成功")

	ctx.JSON(http.StatusOK, model.SuccessResponse("登录成功", loginResp))
}

// Register 用户注册
//
//	@Summary		用户注册
//	@Description	创建新用户账户，需要提供用户名、密码等基本信息
//	@Tags			认证管理
//	@Accept			json
//	@Produce		json
//	@Param			request	body		model.RegisterRequest	true	"注册请求参数"
//	@Success		201		{object}	model.Response{data=model.UserInfo}	"注册成功"
//	@Failure		400		{object}	model.Response{data=string}				"请求参数错误"
//	@Failure		409		{object}	model.Response{data=string}				"用户名或邮箱已存在"
//	@Failure		500		{object}	model.Response{data=string}				"服务器内部错误"
//	@Router			/auth/register [post]
func (c *AuthController) Register(ctx *gin.Context) {
	var req model.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Warnf("注册请求参数错误: %v", err)
		ctx.JSON(http.StatusBadRequest, model.ErrorResponse(model.CodeInvalidParams, "请求参数错误", err.Error()))
		return
	}

	clientIP := ctx.ClientIP()

	// 记录注册尝试
	logger.AuthLog("register_attempt", req.Username, clientIP, false, "开始用户注册")

	user, err := c.authService.Register(ctx, &req)
	if err != nil {
		// 记录注册失败
		logger.AuthLog("register_failed", req.Username, clientIP, false, err.Error())

		switch err.Error() {
		case "用户名已存在":
			ctx.JSON(http.StatusConflict, model.ErrorResponse(model.CodeUserExists, "用户名已存在", nil))
		case "邮箱已存在":
			ctx.JSON(http.StatusConflict, model.ErrorResponse(model.CodeEmailExists, "邮箱已存在", nil))
		default:
			ctx.JSON(http.StatusInternalServerError, model.ErrorResponse(model.CodeServerError, "注册失败", err.Error()))
		}
		return
	}

	// 记录注册成功
	logger.AuthLog("register_success", req.Username, clientIP, true, "用户注册成功")
	logger.BusinessLog("认证管理", "用户注册", user.UserID, "新用户注册成功")

	ctx.JSON(http.StatusCreated, model.SuccessResponse("注册成功", user))
}

// ChangePassword 修改密码
//
//	@Summary		修改密码
//	@Description	用户修改密码，需要提供当前密码和新密码
//	@Tags			认证管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string						true	"Bearer JWT令牌"
//	@Param			request			body		model.ChangePasswordRequest	true	"修改密码请求参数"
//	@Success		200				{object}	model.Response{data=string}		"密码修改成功"
//	@Failure		400				{object}	model.Response{data=string}		"请求参数错误"
//	@Failure		401				{object}	model.Response{data=string}		"未登录或token无效"
//	@Failure		403				{object}	model.Response{data=string}		"当前密码错误"
//	@Failure		500				{object}	model.Response{data=string}		"服务器内部错误"
//	@Security		BearerAuth
//	@Router			/auth/change-password [post]
func (c *AuthController) ChangePassword(ctx *gin.Context) {
	var req model.ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Warnf("修改密码请求参数错误: %v", err)
		ctx.JSON(http.StatusBadRequest, model.ErrorResponse(model.CodeInvalidParams, "请求参数错误", err.Error()))
		return
	}

	// 从中间件获取用户信息
	userID := ctx.GetString("user_id")
	username := ctx.GetString("username")
	clientIP := ctx.ClientIP()

	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, model.ErrorResponse(model.CodeUnauthorized, "未登录", nil))
		return
	}

	// 记录密码修改尝试
	logger.AuthLog("change_password_attempt", username, clientIP, false, "开始修改密码")

	err := c.authService.ChangePassword(ctx, userID, &req)
	if err != nil {
		// 记录密码修改失败
		logger.AuthLog("change_password_failed", username, clientIP, false, err.Error())

		switch err.Error() {
		case "当前密码错误":
			ctx.JSON(http.StatusForbidden, model.ErrorResponse(model.CodeAuthFailed, "当前密码错误", nil))
		default:
			ctx.JSON(http.StatusInternalServerError, model.ErrorResponse(model.CodeServerError, "密码修改失败", err.Error()))
		}
		return
	}

	// 记录密码修改成功
	logger.AuthLog("change_password_success", username, clientIP, true, "密码修改成功")
	logger.BusinessLog("认证管理", "修改密码", userID, "用户修改密码成功")

	ctx.JSON(http.StatusOK, model.SuccessResponse("密码修改成功", nil))
}

// RefreshToken 刷新令牌
//
//	@Summary		刷新访问令牌
//	@Description	使用刷新令牌获取新的访问令牌
//	@Tags			认证管理
//	@Accept			json
//	@Produce		json
//	@Param			request	body		model.RefreshTokenRequest	true	"刷新令牌请求参数"
//	@Success		200		{object}	model.Response{data=model.LoginResponse}	"令牌刷新成功"
//	@Failure		400		{object}	model.Response{data=string}					"请求参数错误"
//	@Failure		401		{object}	model.Response{data=string}					"刷新令牌无效或已过期"
//	@Failure		500		{object}	model.Response{data=string}					"服务器内部错误"
//	@Router			/auth/refresh [post]
func (c *AuthController) RefreshToken(ctx *gin.Context) {
	var req model.RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Warnf("刷新令牌请求参数错误: %v", err)
		ctx.JSON(http.StatusBadRequest, model.ErrorResponse(model.CodeInvalidParams, "请求参数错误", err.Error()))
		return
	}

	clientIP := ctx.ClientIP()

	loginResp, err := c.authService.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		logger.AuthLog("refresh_token_failed", "", clientIP, false, err.Error())
		ctx.JSON(http.StatusUnauthorized, model.ErrorResponse(model.CodeTokenInvalid, "刷新令牌无效或已过期", err.Error()))
		return
	}

	logger.AuthLog("refresh_token_success", loginResp.User.Username, clientIP, true, "令牌刷新成功")
	logger.BusinessLog("认证管理", "刷新令牌", loginResp.User.UserID, "用户刷新令牌成功")

	ctx.JSON(http.StatusOK, model.SuccessResponse("令牌刷新成功", loginResp))
}

// Logout 用户登出
//
//	@Summary		用户登出
//	@Description	用户退出登录，清除相关状态
//	@Tags			认证管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string					true	"Bearer JWT令牌"
//	@Success		200				{object}	model.Response{data=string}	"登出成功"
//	@Failure		401				{object}	model.Response{data=string}	"未登录或token无效"
//	@Failure		500				{object}	model.Response{data=string}	"服务器内部错误"
//	@Security		BearerAuth
//	@Router			/auth/logout [post]
func (c *AuthController) Logout(ctx *gin.Context) {
	// 从中间件获取用户信息
	userID := ctx.GetString("user_id")
	username := ctx.GetString("username")
	clientIP := ctx.ClientIP()

	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, model.ErrorResponse(model.CodeUnauthorized, "未登录", nil))
		return
	}

	// 获取token
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.JSON(http.StatusUnauthorized, model.ErrorResponse(model.CodeUnauthorized, "缺少认证头", nil))
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == authHeader {
		ctx.JSON(http.StatusUnauthorized, model.ErrorResponse(model.CodeUnauthorized, "认证格式错误", nil))
		return
	}

	err := c.authService.Logout(ctx, userID, token)
	if err != nil {
		logger.AuthLog("logout_failed", username, clientIP, false, err.Error())
		ctx.JSON(http.StatusInternalServerError, model.ErrorResponse(model.CodeServerError, "登出失败", err.Error()))
		return
	}

	logger.AuthLog("logout_success", username, clientIP, true, "用户登出成功")
	logger.BusinessLog("认证管理", "用户登出", userID, "用户登出成功")

	ctx.JSON(http.StatusOK, model.SuccessResponse("登出成功", nil))
}

// GetUserInfo 获取用户信息
//
//	@Summary		获取当前用户信息
//	@Description	获取当前登录用户的详细信息
//	@Tags			认证管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string					true	"Bearer JWT令牌"
//	@Success		200				{object}	model.Response{data=model.UserInfo}	"获取成功"
//	@Failure		401				{object}	model.Response{data=string}			"未登录或token无效"
//	@Failure		404				{object}	model.Response{data=string}			"用户不存在"
//	@Failure		500				{object}	model.Response{data=string}			"服务器内部错误"
//	@Security		BearerAuth
//	@Router			/auth/user-info [get]
func (c *AuthController) GetUserInfo(ctx *gin.Context) {
	// 从中间件获取用户信息
	userID := ctx.GetString("user_id")

	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, model.ErrorResponse(model.CodeUnauthorized, "未登录", nil))
		return
	}

	userInfo, err := c.authService.GetUserInfo(ctx, userID)
	if err != nil {
		logger.Errorf("获取用户信息失败: %v, UserID: %s", err, userID)
		if err.Error() == "用户不存在" {
			ctx.JSON(http.StatusNotFound, model.ErrorResponse(model.CodeUserNotExists, "用户不存在", nil))
		} else {
			ctx.JSON(http.StatusInternalServerError, model.ErrorResponse(model.CodeServerError, "获取用户信息失败", err.Error()))
		}
		return
	}

	logger.BusinessLog("认证管理", "获取用户信息", userID, "获取用户信息成功")

	ctx.JSON(http.StatusOK, model.SuccessResponse("获取成功", userInfo))
}
