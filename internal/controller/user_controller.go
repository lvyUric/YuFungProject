package controller

import (
	"net/http"
	"strconv"

	"YufungProject/internal/model"
	"YufungProject/internal/service"
	"YufungProject/pkg/logger"

	"github.com/gin-gonic/gin"
)

// UserController 用户管理控制器
type UserController struct {
	userService    service.UserService
	companyService service.CompanyService
}

// NewUserController 创建用户控制器实例
func NewUserController(userService service.UserService, companyService service.CompanyService) *UserController {
	return &UserController{
		userService:    userService,
		companyService: companyService,
	}
}

// CreateUser 创建用户
// @Summary 创建用户
// @Description 创建新用户（平台管理员专用）
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param user body model.UserCreateRequest true "用户信息"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /api/v1/users [post]
func (uc *UserController) CreateUser(c *gin.Context) {
	var req model.UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("参数绑定失败", err)
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: "参数格式错误: " + err.Error(),
		})
		return
	}

	// 验证所属公司是否存在
	if _, err := uc.companyService.GetCompanyByID(c.Request.Context(), req.CompanyID); err != nil {
		logger.Error("公司不存在", err)
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: "所属公司不存在",
		})
		return
	}

	// 创建用户
	user, err := uc.userService.CreateUser(c.Request.Context(), &req)
	if err != nil {
		logger.Error("创建用户失败", err)
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    http.StatusInternalServerError,
			Message: "创建用户失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code:    http.StatusOK,
		Message: "创建用户成功",
		Data:    user,
	})
}

// GetUserList 获取用户列表
// @Summary 获取用户列表
// @Description 获取用户列表，支持分页和搜索
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页大小" default(10)
// @Param username query string false "用户名搜索"
// @Param display_name query string false "显示名称搜索"
// @Param company_id query string false "公司ID筛选"
// @Param status query string false "状态筛选"
// @Success 200 {object} model.Response{data=model.UserListResponse}
// @Failure 400 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /api/v1/users [get]
func (uc *UserController) GetUserList(c *gin.Context) {
	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// 构建筛选条件
	filter := make(map[string]interface{})

	if username := c.Query("username"); username != "" {
		filter["username"] = username
	}

	if displayName := c.Query("display_name"); displayName != "" {
		filter["display_name"] = displayName
	}

	if companyID := c.Query("company_id"); companyID != "" {
		filter["company_id"] = companyID
	}

	if status := c.Query("status"); status != "" {
		filter["status"] = status
	}

	// 获取用户列表
	response, err := uc.userService.GetUserList(c.Request.Context(), filter, page, pageSize)
	if err != nil {
		logger.Error("获取用户列表失败", err)
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    http.StatusInternalServerError,
			Message: "获取用户列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code:    http.StatusOK,
		Message: "获取用户列表成功",
		Data:    response,
	})
}

// GetUserByID 根据ID获取用户详情
// @Summary 获取用户详情
// @Description 根据用户ID获取用户详细信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path string true "用户ID"
// @Success 200 {object} model.Response{data=model.UserInfo}
// @Failure 400 {object} model.Response
// @Failure 404 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /api/v1/users/{id} [get]
func (uc *UserController) GetUserByID(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: "用户ID不能为空",
		})
		return
	}

	user, err := uc.userService.GetByUserID(c.Request.Context(), userID)
	if err != nil {
		logger.Error("获取用户详情失败", err)
		c.JSON(http.StatusNotFound, model.Response{
			Code:    http.StatusNotFound,
			Message: "用户不存在",
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code:    http.StatusOK,
		Message: "获取用户详情成功",
		Data:    user,
	})
}

// UpdateUser 更新用户信息
// @Summary 更新用户信息
// @Description 更新用户基本信息（不包括密码）
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path string true "用户ID"
// @Param user body model.UserUpdateRequest true "用户更新信息"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Failure 404 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /api/v1/users/{id} [put]
func (uc *UserController) UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: "用户ID不能为空",
		})
		return
	}

	var req model.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("参数绑定失败", err)
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: "参数格式错误: " + err.Error(),
		})
		return
	}

	err := uc.userService.UpdateUser(c.Request.Context(), userID, &req)
	if err != nil {
		logger.Error("更新用户失败", err)
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    http.StatusInternalServerError,
			Message: "更新用户失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code:    http.StatusOK,
		Message: "更新用户成功",
	})
}

// DeleteUser 删除用户
// @Summary 删除用户
// @Description 删除指定用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path string true "用户ID"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Failure 404 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /api/v1/users/{id} [delete]
func (uc *UserController) DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: "用户ID不能为空",
		})
		return
	}

	err := uc.userService.DeleteUser(c.Request.Context(), userID)
	if err != nil {
		logger.Error("删除用户失败", err)
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    http.StatusInternalServerError,
			Message: "删除用户失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code:    http.StatusOK,
		Message: "删除用户成功",
	})
}

// BatchUpdateUserStatus 批量更新用户状态
// @Summary 批量更新用户状态
// @Description 批量启用/停用用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body model.BatchUpdateUserStatusRequest true "批量更新请求"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /api/v1/users/batch-status [put]
func (uc *UserController) BatchUpdateUserStatus(c *gin.Context) {
	var req model.BatchUpdateUserStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("参数绑定失败", err)
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: "参数格式错误: " + err.Error(),
		})
		return
	}

	if len(req.UserIDs) == 0 {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: "用户ID列表不能为空",
		})
		return
	}

	err := uc.userService.BatchUpdateUserStatus(c.Request.Context(), req.UserIDs, req.Status)
	if err != nil {
		logger.Error("批量更新用户状态失败", err)
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    http.StatusInternalServerError,
			Message: "批量更新用户状态失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code:    http.StatusOK,
		Message: "批量更新用户状态成功",
	})
}

// ResetUserPassword 重置用户密码
// @Summary 重置用户密码
// @Description 管理员重置用户密码
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path string true "用户ID"
// @Param request body model.ResetPasswordRequest true "重置密码请求"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /api/v1/users/{id}/reset-password [put]
func (uc *UserController) ResetUserPassword(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: "用户ID不能为空",
		})
		return
	}

	var req model.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("参数绑定失败", err)
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: "参数格式错误: " + err.Error(),
		})
		return
	}

	// 确保路径参数和请求体中的用户ID一致
	req.UserID = userID

	err := uc.userService.ResetPassword(c.Request.Context(), &req)
	if err != nil {
		logger.Error("重置密码失败", err)
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    http.StatusInternalServerError,
			Message: "重置密码失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code:    http.StatusOK,
		Message: "重置密码成功",
	})
}

// ExportUsers 导出用户数据
// @Summary 导出用户数据
// @Description 导出用户数据为Excel格式
// @Tags 用户管理
// @Accept json
// @Produce application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Param company_id query string false "公司ID筛选"
// @Param status query string false "状态筛选"
// @Success 200 {file} file "Excel文件"
// @Failure 500 {object} model.Response
// @Router /api/v1/users/export [get]
func (uc *UserController) ExportUsers(c *gin.Context) {
	// 构建筛选条件
	filter := make(map[string]interface{})

	if companyID := c.Query("company_id"); companyID != "" {
		filter["company_id"] = companyID
	}

	if status := c.Query("status"); status != "" {
		filter["status"] = status
	}

	// 导出用户数据
	fileData, filename, err := uc.userService.ExportUsers(c.Request.Context(), filter)
	if err != nil {
		logger.Error("导出用户数据失败", err)
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    http.StatusInternalServerError,
			Message: "导出用户数据失败: " + err.Error(),
		})
		return
	}

	// 设置响应头
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Length", strconv.Itoa(len(fileData)))

	// 返回文件数据
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", fileData)
}

// QuickDisableUser 快捷停用用户
// @Summary 快捷停用用户
// @Description 一键停用用户账户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path string true "用户ID"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /api/v1/users/{id}/quick-disable [put]
func (uc *UserController) QuickDisableUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: "用户ID不能为空",
		})
		return
	}

	err := uc.userService.QuickDisableUser(c.Request.Context(), userID)
	if err != nil {
		logger.Error("快捷停用用户失败", err)
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    http.StatusInternalServerError,
			Message: "快捷停用用户失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code:    http.StatusOK,
		Message: "快捷停用用户成功",
	})
}

// ExportUsersAdvanced 高级导出用户数据
// @Summary 高级导出用户数据
// @Description 根据条件导出用户数据到Excel或CSV文件
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body model.UserExportRequest true "导出请求"
// @Success 200 {object} model.Response{data=model.UserExportResponse}
// @Failure 400 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /api/v1/users/export-advanced [post]
func (uc *UserController) ExportUsersAdvanced(c *gin.Context) {
	var req model.UserExportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("参数绑定失败", err)
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: "参数格式错误: " + err.Error(),
		})
		return
	}

	response, err := uc.userService.ExportUsersAdvanced(c.Request.Context(), &req)
	if err != nil {
		logger.Error("高级导出用户数据失败", err)
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    http.StatusInternalServerError,
			Message: "导出失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code:    http.StatusOK,
		Message: "导出成功",
		Data:    response,
	})
}

// DownloadUserTemplate 下载用户导入模板
// @Summary 下载用户导入模板
// @Description 下载用户数据导入模板文件
// @Tags 用户管理
// @Accept json
// @Produce application/octet-stream
// @Param format query string false "文件格式：xlsx, csv" default(xlsx)
// @Success 200 {file} binary "模板文件"
// @Failure 400 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /api/v1/users/template [get]
func (uc *UserController) DownloadUserTemplate(c *gin.Context) {
	format := c.DefaultQuery("format", "xlsx")
	if format != "xlsx" && format != "csv" {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: "不支持的文件格式",
		})
		return
	}

	fileData, fileName, err := uc.userService.GenerateUserTemplate(c.Request.Context(), format)
	if err != nil {
		logger.Error("生成用户模板失败", err)
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    http.StatusInternalServerError,
			Message: "生成模板失败: " + err.Error(),
		})
		return
	}

	// 设置响应头
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Length", strconv.Itoa(len(fileData)))

	c.Data(http.StatusOK, "application/octet-stream", fileData)
}

// PreviewUserImport 预览用户导入数据
// @Summary 预览用户导入数据
// @Description 预览导入的用户数据，检查格式和错误
// @Tags 用户管理
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "导入文件"
// @Param skip_header formData bool false "是否跳过表头行"
// @Param update_existing formData bool false "是否更新已存在的用户"
// @Success 200 {object} model.Response{data=model.UserImportResponse}
// @Failure 400 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /api/v1/users/import/preview [post]
func (uc *UserController) PreviewUserImport(c *gin.Context) {
	// 获取上传文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		logger.Error("获取上传文件失败", err)
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: "请选择要上传的文件: " + err.Error(),
		})
		return
	}
	defer file.Close()

	// 获取其他参数
	var req model.UserImportRequest
	if err := c.ShouldBind(&req); err != nil {
		logger.Error("导入预览请求参数错误", err)
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	response, err := uc.userService.PreviewUserImport(c.Request.Context(), file, header, &req)
	if err != nil {
		logger.Error("预览用户导入失败", err)
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    http.StatusInternalServerError,
			Message: "预览失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code:    http.StatusOK,
		Message: "预览成功",
		Data:    response,
	})
}

// ImportUsers 导入用户数据
// @Summary 导入用户数据
// @Description 导入用户数据到系统中
// @Tags 用户管理
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "导入文件"
// @Param skip_header formData bool false "是否跳过表头行"
// @Param update_existing formData bool false "是否更新已存在的用户"
// @Success 200 {object} model.Response{data=model.UserImportResponse}
// @Failure 400 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /api/v1/users/import [post]
func (uc *UserController) ImportUsers(c *gin.Context) {
	// 获取上传文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		logger.Error("获取上传文件失败", err)
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: "请选择要上传的文件: " + err.Error(),
		})
		return
	}
	defer file.Close()

	// 获取其他参数
	var req model.UserImportRequest
	if err := c.ShouldBind(&req); err != nil {
		logger.Error("导入请求参数错误", err)
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	response, err := uc.userService.ImportUsers(c.Request.Context(), file, header, &req)
	if err != nil {
		logger.Error("导入用户数据失败", err)
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    http.StatusInternalServerError,
			Message: "导入失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code:    http.StatusOK,
		Message: "导入完成",
		Data:    response,
	})
}
