package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"YufungProject/internal/model"
	"YufungProject/internal/service"
	"YufungProject/pkg/logger"
)

type ChangeRecordController struct {
	changeRecordService *service.ChangeRecordService
}

func NewChangeRecordController(changeRecordService *service.ChangeRecordService) *ChangeRecordController {
	return &ChangeRecordController{
		changeRecordService: changeRecordService,
	}
}

// GetPolicyChangeRecords 获取保单变更记录
// @Summary 获取保单变更记录
// @Description 获取指定保单的变更记录，支持分页和时间范围筛选
// @Tags 变更记录
// @Accept json
// @Produce json
// @Param id path string true "保单ID"
// @Param days query int false "查询天数" default(10)
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} model.Response{data=[]model.ChangeRecordResponse} "成功"
// @Failure 400 {object} model.Response "请求参数错误"
// @Failure 401 {object} model.Response "未授权"
// @Failure 500 {object} model.Response "服务器错误"
// @Router /api/policies/{id}/change-records [get]
func (c *ChangeRecordController) GetPolicyChangeRecords(ctx *gin.Context) {
	policyID := ctx.Param("id")
	if policyID == "" {
		ctx.JSON(http.StatusBadRequest, model.ValidationError(errors.New("保单ID不能为空")))
		return
	}

	// 解析查询参数
	days := 10 // 默认查询10天
	if daysStr := ctx.Query("days"); daysStr != "" {
		if parsedDays, err := strconv.Atoi(daysStr); err == nil && parsedDays > 0 {
			days = parsedDays
		}
	}

	page := 1
	if pageStr := ctx.Query("page"); pageStr != "" {
		if parsedPage, err := strconv.Atoi(pageStr); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	pageSize := 10
	if pageSizeStr := ctx.Query("page_size"); pageSizeStr != "" {
		if parsedPageSize, err := strconv.Atoi(pageSizeStr); err == nil && parsedPageSize > 0 && parsedPageSize <= 100 {
			pageSize = parsedPageSize
		}
	}

	records, total, err := c.changeRecordService.GetChangeRecordsByPolicy(ctx.Request.Context(), policyID, days, page, pageSize)
	if err != nil {
		logger.Error("Failed to get policy change records", "policy_id", policyID, "error", err)
		ctx.JSON(http.StatusInternalServerError, model.ServerError("获取变更记录失败"))
		return
	}

	// 构建响应数据
	response := map[string]interface{}{
		"records":   records,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"has_more":  int64(page*pageSize) < total,
	}

	ctx.JSON(http.StatusOK, model.Success(response))
}

// GetChangeRecordsList 获取变更记录列表
// @Summary 获取变更记录列表
// @Description 获取变更记录列表，支持多种筛选条件
// @Tags 变更记录
// @Accept json
// @Produce json
// @Param table_name query string false "表名"
// @Param record_id query string false "记录ID"
// @Param user_id query string false "用户ID"
// @Param change_type query string false "变更类型"
// @Param start_time query string false "开始时间(YYYY-MM-DD)"
// @Param end_time query string false "结束时间(YYYY-MM-DD)"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} model.Response{data=[]model.ChangeRecordResponse} "成功"
// @Failure 400 {object} model.Response "请求参数错误"
// @Failure 401 {object} model.Response "未授权"
// @Failure 500 {object} model.Response "服务器错误"
// @Router /api/change-records [get]
func (c *ChangeRecordController) GetChangeRecordsList(ctx *gin.Context) {
	// 获取用户公司ID（权限控制）
	companyID, exists := ctx.Get("company_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, model.UnauthorizedError("公司信息缺失"))
		return
	}

	// 构建查询参数
	params := &model.ChangeRecordListParams{
		TableName:  ctx.Query("table_name"),
		RecordID:   ctx.Query("record_id"),
		UserID:     ctx.Query("user_id"),
		CompanyID:  companyID.(string), // 强制按公司过滤
		ChangeType: ctx.Query("change_type"),
		StartTime:  ctx.Query("start_time"),
		EndTime:    ctx.Query("end_time"),
		Page:       1,
		PageSize:   10,
	}

	// 解析分页参数
	if pageStr := ctx.Query("page"); pageStr != "" {
		if parsedPage, err := strconv.Atoi(pageStr); err == nil && parsedPage > 0 {
			params.Page = parsedPage
		}
	}

	if pageSizeStr := ctx.Query("page_size"); pageSizeStr != "" {
		if parsedPageSize, err := strconv.Atoi(pageSizeStr); err == nil && parsedPageSize > 0 && parsedPageSize <= 100 {
			params.PageSize = parsedPageSize
		}
	}

	records, total, err := c.changeRecordService.GetChangeRecordsList(ctx.Request.Context(), params)
	if err != nil {
		logger.Error("Failed to get change records list", "params", params, "error", err)
		ctx.JSON(http.StatusInternalServerError, model.ServerError("获取变更记录列表失败"))
		return
	}

	// 构建响应数据
	response := map[string]interface{}{
		"records":   records,
		"total":     total,
		"page":      params.Page,
		"page_size": params.PageSize,
		"has_more":  int64(params.Page*params.PageSize) < total,
	}

	ctx.JSON(http.StatusOK, model.Success(response))
}
