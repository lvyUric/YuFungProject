package controller

import (
	"YufungProject/internal/model"
	"YufungProject/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ActivityLogController struct {
	service *service.ActivityLogService
}

func NewActivityLogController() *ActivityLogController {
	return &ActivityLogController{
		service: service.NewActivityLogService(),
	}
}

// GetActivityLogList 获取活动记录列表
// @Summary 获取活动记录列表
// @Description 获取系统活动记录列表，支持分页和过滤
// @Tags 活动记录
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param company_id query string false "公司ID"
// @Param user_id query string false "用户ID"
// @Param operation_type query string false "操作类型"
// @Param module_name query string false "模块名称"
// @Param start_time query string false "开始时间"
// @Param end_time query string false "结束时间"
// @Success 200 {object} model.Response{data=model.ActivityLogResponse}
// @Router /api/activity-logs [get]
func (c *ActivityLogController) GetActivityLogList(ctx *gin.Context) {
	var query model.ActivityLogQuery

	// 获取查询参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	query.Page = page
	query.PageSize = pageSize
	query.CompanyID = ctx.Query("company_id")
	query.UserID = ctx.Query("user_id")
	query.OperationType = ctx.Query("operation_type")
	query.ModuleName = ctx.Query("module_name")
	query.StartTime = ctx.Query("start_time")
	query.EndTime = ctx.Query("end_time")

	// 权限控制：非平台管理员只能查看本公司的记录
	user := GetCurrentUser(ctx)
	if user != nil && !IsPlatformAdmin(user) && user.CompanyID != "" {
		query.CompanyID = user.CompanyID
	}

	result, err := c.service.GetActivityLogList(ctx, &query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.Response{
			Code:    500,
			Message: "获取活动记录失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Code:    200,
		Message: "获取活动记录成功",
		Data:    result,
	})
}

// GetRecentActivityLogs 获取最近的活动记录（用于仪表盘）
// @Summary 获取最近的活动记录
// @Description 获取最近的活动记录，用于首页仪表盘显示
// @Tags 活动记录
// @Accept json
// @Produce json
// @Param limit query int false "记录数量" default(5)
// @Success 200 {object} model.Response{data=[]model.ActivityLog}
// @Router /api/activity-logs/recent [get]
func (c *ActivityLogController) GetRecentActivityLogs(ctx *gin.Context) {
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "5"))

	// 权限控制：非平台管理员只能查看本公司的记录
	user := GetCurrentUser(ctx)
	var companyID string
	if user != nil && !IsPlatformAdmin(user) {
		companyID = user.CompanyID
	}

	logs, err := c.service.GetRecentActivityLogs(ctx, companyID, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.Response{
			Code:    500,
			Message: "获取最近活动记录失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Code:    200,
		Message: "获取最近活动记录成功",
		Data:    logs,
	})
}

// GetActivityLogByID 根据ID获取活动记录详情
// @Summary 获取活动记录详情
// @Description 根据ID获取活动记录详情
// @Tags 活动记录
// @Accept json
// @Produce json
// @Param id path string true "记录ID"
// @Success 200 {object} model.Response{data=model.ActivityLog}
// @Router /api/activity-logs/{id} [get]
func (c *ActivityLogController) GetActivityLogByID(ctx *gin.Context) {
	id := ctx.Param("id")

	log, err := c.service.GetActivityLogByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, model.Response{
			Code:    404,
			Message: "活动记录不存在: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Code:    200,
		Message: "获取活动记录详情成功",
		Data:    log,
	})
}

// GetActivityLogStatistics 获取活动记录统计
// @Summary 获取活动记录统计
// @Description 获取活动记录统计信息
// @Tags 活动记录
// @Accept json
// @Produce json
// @Param days query int false "统计天数" default(7)
// @Success 200 {object} model.Response{data=object}
// @Router /api/activity-logs/statistics [get]
func (c *ActivityLogController) GetActivityLogStatistics(ctx *gin.Context) {
	days, _ := strconv.Atoi(ctx.DefaultQuery("days", "7"))

	// 权限控制：非平台管理员只能查看本公司的统计
	user := GetCurrentUser(ctx)
	var companyID string
	if user != nil && !IsPlatformAdmin(user) {
		companyID = user.CompanyID
	}

	statistics, err := c.service.GetActivityLogStatistics(ctx, companyID, days)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.Response{
			Code:    500,
			Message: "获取活动记录统计失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Code:    200,
		Message: "获取活动记录统计成功",
		Data:    statistics,
	})
}

// DeleteActivityLogsByCompanyID 删除指定公司的活动记录
// @Summary 删除公司活动记录
// @Description 删除指定公司的所有活动记录（仅平台管理员可用）
// @Tags 活动记录
// @Accept json
// @Produce json
// @Param company_id path string true "公司ID"
// @Success 200 {object} model.Response
// @Router /api/activity-logs/company/{company_id} [delete]
func (c *ActivityLogController) DeleteActivityLogsByCompanyID(ctx *gin.Context) {
	// 权限检查：只有平台管理员可以删除
	user := GetCurrentUser(ctx)
	if user == nil || !IsPlatformAdmin(user) {
		ctx.JSON(http.StatusForbidden, model.Response{
			Code:    403,
			Message: "权限不足，只有平台管理员可以删除活动记录",
		})
		return
	}

	companyID := ctx.Param("company_id")

	err := c.service.DeleteActivityLogsByCompanyID(ctx, companyID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.Response{
			Code:    500,
			Message: "删除活动记录失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Code:    200,
		Message: "删除活动记录成功",
	})
}

// GetCurrentUser 获取当前用户（需要根据你的用户认证逻辑实现）
func GetCurrentUser(ctx *gin.Context) *model.User {
	// 这里需要根据你的用户认证中间件实现
	// 示例实现，你需要根据实际情况修改
	userID, exists := ctx.Get("user_id")
	if !exists {
		return nil
	}

	// 这里应该从数据库或缓存中获取用户信息
	// 暂时返回示例数据
	return &model.User{
		UserID:    userID.(string),
		Username:  "admin",
		CompanyID: "company_001",
		Status:    "active",
	}
}

// IsPlatformAdmin 判断是否为平台管理员
func IsPlatformAdmin(user *model.User) bool {
	// 根据你的业务逻辑判断平台管理员
	// 这里可以根据角色ID、用户名等来判断
	// 示例：用户名为admin或包含特定角色ID的用户为平台管理员
	if user == nil {
		return false
	}

	// 示例判断逻辑，你需要根据实际情况修改
	return user.Username == "admin" || user.CompanyID == ""
}
