package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"YufungProject/internal/model"
	"YufungProject/internal/service"
	"YufungProject/pkg/logger"
	"errors"
)

type PolicyController struct {
	policyService *service.PolicyService
}

func NewPolicyController(policyService *service.PolicyService) *PolicyController {
	return &PolicyController{
		policyService: policyService,
	}
}

// CreatePolicy 创建保单
// @Summary 创建保单
// @Description 创建新的保单记录
// @Tags 保单管理
// @Accept json
// @Produce json
// @Param request body model.PolicyCreateRequest true "创建保单请求"
// @Success 200 {object} model.Response{data=model.PolicyResponse} "成功"
// @Failure 400 {object} model.Response "请求参数错误"
// @Failure 401 {object} model.Response "未授权"
// @Failure 500 {object} model.Response "服务器错误"
// @Router /api/policies [post]
func (c *PolicyController) CreatePolicy(ctx *gin.Context) {
	var req model.PolicyCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, model.ValidationError(err))
		return
	}

	// 获取用户信息
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, model.UnauthorizedError("用户未登录"))
		return
	}

	companyID, exists := ctx.Get("company_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, model.UnauthorizedError("公司信息缺失"))
		return
	}

	policy, err := c.policyService.CreatePolicy(ctx.Request.Context(), &req, userID.(string), companyID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.ServerError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, model.Success(policy))
}

// GetPolicy 获取保单详情
// @Summary 获取保单详情
// @Description 根据保单ID获取保单详细信息
// @Tags 保单管理
// @Accept json
// @Produce json
// @Param id path string true "保单ID"
// @Success 200 {object} model.Response{data=model.PolicyResponse} "成功"
// @Failure 400 {object} model.Response "请求参数错误"
// @Failure 401 {object} model.Response "未授权"
// @Failure 404 {object} model.Response "保单不存在"
// @Failure 500 {object} model.Response "服务器错误"
// @Router /api/policies/{id} [get]
func (c *PolicyController) GetPolicy(ctx *gin.Context) {
	policyID := ctx.Param("id")
	if policyID == "" {
		ctx.JSON(http.StatusBadRequest, model.Error(model.CodeInvalidParams, "保单ID不能为空"))
		return
	}

	companyID, exists := ctx.Get("company_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, model.UnauthorizedError("公司信息缺失"))
		return
	}

	policy, err := c.policyService.GetPolicyByID(ctx.Request.Context(), policyID, companyID.(string))
	if err != nil {
		if err.Error() == "保单不存在" || err.Error() == "无权访问该保单" {
			ctx.JSON(http.StatusNotFound, model.NotFoundError(err.Error()))
			return
		}
		ctx.JSON(http.StatusInternalServerError, model.ServerError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, model.Success(policy))
}

// UpdatePolicy 更新保单
// @Summary 更新保单
// @Description 更新指定ID的保单信息
// @Tags 保单管理
// @Accept json
// @Produce json
// @Param id path string true "保单ID"
// @Param request body model.PolicyUpdateRequest true "更新保单请求"
// @Success 200 {object} model.Response{data=model.PolicyResponse} "成功"
// @Failure 400 {object} model.Response "请求参数错误"
// @Failure 401 {object} model.Response "未授权"
// @Failure 404 {object} model.Response "保单不存在"
// @Failure 500 {object} model.Response "服务器错误"
// @Router /api/policies/{id} [put]
func (c *PolicyController) UpdatePolicy(ctx *gin.Context) {
	policyID := ctx.Param("id")
	if policyID == "" {
		ctx.JSON(http.StatusBadRequest, model.ValidationError(errors.New("保单ID不能为空")))
		return
	}

	var req model.PolicyUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, model.ValidationError(err))
		return
	}

	// 获取用户信息
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, model.UnauthorizedError("用户未登录"))
		return
	}

	companyID, exists := ctx.Get("company_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, model.UnauthorizedError("公司信息缺失"))
		return
	}

	// 获取客户端IP地址
	ipAddress := ctx.ClientIP()

	// 获取User-Agent
	userAgent := ctx.GetHeader("User-Agent")

	policy, err := c.policyService.UpdatePolicy(
		ctx.Request.Context(),
		policyID,
		&req,
		userID.(string),
		companyID.(string),
		ipAddress,
		userAgent,
	)
	if err != nil {
		if err.Error() == "保单不存在" {
			ctx.JSON(http.StatusNotFound, model.NotFoundError("保单不存在"))
			return
		}
		ctx.JSON(http.StatusInternalServerError, model.ServerError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, model.Success(policy))
}

// DeletePolicy 删除保单
// @Summary 删除保单
// @Description 删除指定的保单
// @Tags 保单管理
// @Accept json
// @Produce json
// @Param id path string true "保单ID"
// @Success 200 {object} model.Response "成功"
// @Failure 400 {object} model.Response "请求参数错误"
// @Failure 401 {object} model.Response "未授权"
// @Failure 404 {object} model.Response "保单不存在"
// @Failure 500 {object} model.Response "服务器错误"
// @Router /api/policies/{id} [delete]
func (c *PolicyController) DeletePolicy(ctx *gin.Context) {
	policyID := ctx.Param("id")
	if policyID == "" {
		ctx.JSON(http.StatusBadRequest, model.Error(model.CodeInvalidParams, "保单ID不能为空"))
		return
	}

	companyID, exists := ctx.Get("company_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, model.UnauthorizedError("公司信息缺失"))
		return
	}

	err := c.policyService.DeletePolicy(ctx.Request.Context(), policyID, companyID.(string))
	if err != nil {
		if err.Error() == "保单不存在" || err.Error() == "无权删除该保单" {
			ctx.JSON(http.StatusNotFound, model.NotFoundError(err.Error()))
			return
		}
		ctx.JSON(http.StatusInternalServerError, model.ServerError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, model.SuccessResponse("删除成功", nil))
}

// ListPolicies 获取保单列表
// @Summary 获取保单列表
// @Description 分页查询保单列表，支持多条件筛选
// @Tags 保单管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param account_number query string false "账户号"
// @Param customer_number query string false "客户号"
// @Param customer_name_cn query string false "客户中文名"
// @Param customer_name_en query string false "客户英文名"
// @Param proposal_number query string false "投保单号"
// @Param policy_currency query string false "保单币种" Enums(USD, HKD, CNY)
// @Param partner query string false "合作伙伴"
// @Param referral_code query string false "转介编号"
// @Param hk_manager query string false "港分客户经理"
// @Param referral_pm query string false "转介理财经理"
// @Param referral_branch query string false "转介分行"
// @Param referral_sub_branch query string false "转介支行"
// @Param payment_method query string false "缴费方式" Enums(期缴, 趸缴, 预缴)
// @Param insurance_company query string false "承保公司"
// @Param product_name query string false "保险产品名称"
// @Param product_type query string false "产品类型"
// @Param is_surrendered query bool false "是否退保"
// @Param past_cooling_period query bool false "是否已过冷静期"
// @Param is_paid_commission query bool false "是否支付佣金"
// @Param is_employee query bool false "是否员工"
// @Param referral_date_start query string false "转介日期开始" format(date)
// @Param referral_date_end query string false "转介日期结束" format(date)
// @Param payment_date_start query string false "缴费日期开始" format(date)
// @Param payment_date_end query string false "缴费日期结束" format(date)
// @Param effective_date_start query string false "生效日期开始" format(date)
// @Param effective_date_end query string false "生效日期结束" format(date)
// @Param sort_by query string false "排序字段"
// @Param sort_order query string false "排序方向" Enums(asc, desc)
// @Success 200 {object} model.Response{data=model.PolicyListResponse} "成功"
// @Failure 400 {object} model.Response "请求参数错误"
// @Failure 401 {object} model.Response "未授权"
// @Failure 500 {object} model.Response "服务器错误"
// @Router /api/policies [get]
func (c *PolicyController) ListPolicies(ctx *gin.Context) {
	var req model.PolicyQueryRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, model.ValidationError(err))
		return
	}

	companyID, exists := ctx.Get("company_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, model.UnauthorizedError("公司信息缺失"))
		return
	}

	policies, err := c.policyService.ListPolicies(ctx.Request.Context(), &req, companyID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.ServerError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, model.Success(policies))
}

// GetPolicyStatistics 获取保单统计
// @Summary 获取保单统计
// @Description 获取当前公司的保单统计信息
// @Tags 保单管理
// @Accept json
// @Produce json
// @Success 200 {object} model.Response{data=model.PolicyStatistics} "成功"
// @Failure 401 {object} model.Response "未授权"
// @Failure 500 {object} model.Response "服务器错误"
// @Router /api/policies/statistics [get]
func (c *PolicyController) GetPolicyStatistics(ctx *gin.Context) {
	companyID, exists := ctx.Get("company_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, model.UnauthorizedError("公司信息缺失"))
		return
	}

	stats, err := c.policyService.GetPolicyStatistics(ctx.Request.Context(), companyID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.ServerError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, model.Success(stats))
}

// BatchUpdatePolicyStatus 批量更新保单状态
// @Summary 批量更新保单状态
// @Description 批量更新多个保单的状态字段
// @Tags 保单管理
// @Accept json
// @Produce json
// @Param request body model.BatchUpdatePolicyStatusRequest true "批量更新请求"
// @Success 200 {object} model.Response "成功"
// @Failure 400 {object} model.Response "请求参数错误"
// @Failure 401 {object} model.Response "未授权"
// @Failure 500 {object} model.Response "服务器错误"
// @Router /api/policies/batch-update [post]
func (c *PolicyController) BatchUpdatePolicyStatus(ctx *gin.Context) {
	var req model.BatchUpdatePolicyStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, model.ValidationError(err))
		return
	}

	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, model.UnauthorizedError("用户未登录"))
		return
	}

	companyID, exists := ctx.Get("company_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, model.UnauthorizedError("公司信息缺失"))
		return
	}

	err := c.policyService.BatchUpdatePolicyStatus(ctx.Request.Context(), &req, userID.(string), companyID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.ServerError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, model.SuccessResponse("批量更新成功", nil))
}

// ExportPolicies 导出保单
// @Summary 导出保单
// @Description 导出保单数据为Excel或CSV格式
// @Tags 保单管理
// @Accept json
// @Produce json
// @Param request body model.PolicyExportRequest true "导出请求"
// @Success 200 {object} model.Response{data=[]model.Policy} "成功"
// @Failure 400 {object} model.Response "请求参数错误"
// @Failure 401 {object} model.Response "未授权"
// @Failure 500 {object} model.Response "服务器错误"
// @Router /api/policies/export [post]
func (c *PolicyController) ExportPolicies(ctx *gin.Context) {
	var req model.PolicyExportRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, model.ValidationError(err))
		return
	}

	companyID, exists := ctx.Get("company_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, model.UnauthorizedError("公司信息缺失"))
		return
	}

	// 如果请求包含format参数，则导出文件
	format := ctx.DefaultQuery("format", "")
	if format != "" {
		fileData, fileName, err := c.policyService.ExportPoliciesToFile(ctx.Request.Context(), &req, companyID.(string), format)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, model.ServerError(err.Error()))
			return
		}

		// 设置响应头
		if format == "csv" {
			ctx.Header("Content-Type", "text/csv")
			ctx.Header("Content-Disposition", "attachment; filename="+fileName)
		} else {
			ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
			ctx.Header("Content-Disposition", "attachment; filename="+fileName)
		}

		ctx.Data(http.StatusOK, "application/octet-stream", fileData)
		return
	}

	// 否则返回JSON数据
	policies, err := c.policyService.ExportPolicies(ctx.Request.Context(), &req, companyID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.ServerError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, model.Success(policies))
}

// DownloadPolicyTemplate 下载保单导入模板
// @Summary 下载保单导入模板
// @Description 下载用于批量导入的保单模板文件
// @Tags 保单管理
// @Accept json
// @Produce application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Param type query string false "模板类型" Enums(xlsx, csv) default(xlsx)
// @Success 200 {file} file "模板文件"
// @Failure 500 {object} model.Response "服务器错误"
// @Router /api/policies/template [get]
func (c *PolicyController) DownloadPolicyTemplate(ctx *gin.Context) {
	templateType := ctx.DefaultQuery("type", "xlsx")

	fileData, fileName, err := c.policyService.GeneratePolicyTemplate(ctx.Request.Context(), templateType)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.ServerError(err.Error()))
		return
	}

	// 设置响应头
	if templateType == "csv" {
		ctx.Header("Content-Type", "text/csv")
		ctx.Header("Content-Disposition", "attachment; filename="+fileName)
	} else {
		ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		ctx.Header("Content-Disposition", "attachment; filename="+fileName)
	}

	ctx.Data(http.StatusOK, "application/octet-stream", fileData)
}

// PreviewPolicyImport 预览保单导入数据
// @Summary 预览保单导入数据
// @Description 预览导入的保单数据，检查格式和错误
// @Tags 保单管理
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "Bearer JWT令牌"
// @Param file formData file true "导入文件"
// @Param skip_header formData bool false "是否跳过表头行"
// @Param update_existing formData bool false "是否更新已存在的保单"
// @Success 200 {object} model.Response{data=model.PolicyImportResponse} "预览成功"
// @Failure 400 {object} model.Response{data=string} "请求参数错误"
// @Failure 500 {object} model.Response{data=string} "服务器内部错误"
// @Router /api/policies/import/preview [post]
func (c *PolicyController) PreviewPolicyImport(ctx *gin.Context) {
	// 获取上传文件
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		logger.Warnf("获取上传文件失败: %v", err)
		ctx.JSON(http.StatusBadRequest, model.ErrorResponse(model.CodeInvalidParams, "请选择要上传的文件", err.Error()))
		return
	}
	defer file.Close()

	// 获取其他参数
	var req model.PolicyImportFileRequest
	if err := ctx.ShouldBind(&req); err != nil {
		logger.Warnf("导入预览请求参数错误: %v", err)
		ctx.JSON(http.StatusBadRequest, model.ErrorResponse(model.CodeInvalidParams, "请求参数错误", err.Error()))
		return
	}

	// 记录操作日志
	userID, _ := ctx.Get("user_id")
	logger.BusinessLog("保单管理", "预览导入", userID.(string), "文件名: "+header.Filename)

	response, err := c.policyService.PreviewPolicyImport(ctx.Request.Context(), file, header, &req)
	if err != nil {
		logger.Errorf("预览导入失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, model.ErrorResponse(model.CodeServerError, "预览失败", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, model.SuccessResponse("预览成功", response))
}

// ImportPoliciesFromFile 从文件导入保单
// @Summary 从文件导入保单
// @Description 从Excel或CSV文件批量导入保单数据
// @Tags 保单管理
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "Bearer JWT令牌"
// @Param file formData file true "导入文件"
// @Param skip_header formData bool false "是否跳过表头行"
// @Param update_existing formData bool false "是否更新已存在的保单"
// @Success 200 {object} model.Response{data=model.PolicyImportResponse} "导入成功"
// @Failure 400 {object} model.Response{data=string} "请求参数错误"
// @Failure 500 {object} model.Response{data=string} "服务器内部错误"
// @Router /api/policies/import [post]
func (c *PolicyController) ImportPoliciesFromFile(ctx *gin.Context) {
	// 获取上传文件
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		logger.Warnf("获取上传文件失败: %v", err)
		ctx.JSON(http.StatusBadRequest, model.ErrorResponse(model.CodeInvalidParams, "请选择要上传的文件", err.Error()))
		return
	}
	defer file.Close()

	// 获取其他参数
	var req model.PolicyImportFileRequest
	if err := ctx.ShouldBind(&req); err != nil {
		logger.Warnf("导入请求参数错误: %v", err)
		ctx.JSON(http.StatusBadRequest, model.ErrorResponse(model.CodeInvalidParams, "请求参数错误", err.Error()))
		return
	}

	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, model.ErrorResponse(model.CodeUnauthorized, "用户未登录", ""))
		return
	}

	companyID, exists := ctx.Get("company_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, model.ErrorResponse(model.CodeUnauthorized, "公司信息缺失", ""))
		return
	}

	// 记录操作日志
	logger.BusinessLog("保单管理", "导入保单", userID.(string), "文件名: "+header.Filename)

	response, err := c.policyService.ImportPoliciesFromFile(ctx.Request.Context(), file, header, &req, userID.(string), companyID.(string))
	if err != nil {
		logger.Errorf("导入保单数据失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, model.ErrorResponse(model.CodeServerError, "导入失败", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, model.SuccessResponse("导入完成", response))
}

// GetPolicyValidationRules 获取保单字段验证规则
// @Summary 获取保单字段验证规则
// @Description 获取保单各字段的验证规则，用于前端表单验证
// @Tags 保单管理
// @Accept json
// @Produce json
// @Success 200 {object} model.Response{data=map[string]interface{}} "成功"
// @Router /api/policies/validation-rules [get]
func (c *PolicyController) GetPolicyValidationRules(ctx *gin.Context) {
	rules := map[string]interface{}{
		"account_number": map[string]interface{}{
			"required":  false,
			"type":      "string",
			"maxLength": 50,
		},
		"customer_number": map[string]interface{}{
			"required":  true,
			"type":      "string",
			"maxLength": 50,
		},
		"customer_name_cn": map[string]interface{}{
			"required":  true,
			"type":      "string",
			"maxLength": 100,
		},
		"customer_name_en": map[string]interface{}{
			"required":  false,
			"type":      "string",
			"maxLength": 100,
		},
		"proposal_number": map[string]interface{}{
			"required":  true,
			"type":      "string",
			"maxLength": 50,
		},
		"policy_currency": map[string]interface{}{
			"required": true,
			"type":     "string",
			"enum":     []string{"USD", "HKD", "CNY"},
		},
		"payment_method": map[string]interface{}{
			"required": false,
			"type":     "string",
			"enum":     []string{"期缴", "趸缴", "预缴"},
		},
		"actual_premium": map[string]interface{}{
			"required": false,
			"type":     "number",
			"min":      0,
		},
		"aum": map[string]interface{}{
			"required": false,
			"type":     "number",
			"min":      0,
		},
		"referral_rate": map[string]interface{}{
			"required": false,
			"type":     "number",
			"min":      0,
			"max":      100,
		},
		"exchange_rate": map[string]interface{}{
			"required": false,
			"type":     "number",
			"min":      0,
		},
		"expected_fee": map[string]interface{}{
			"required": false,
			"type":     "number",
			"min":      0,
		},
		"insurance_company": map[string]interface{}{
			"required":  true,
			"type":      "string",
			"maxLength": 100,
		},
		"product_name": map[string]interface{}{
			"required":  true,
			"type":      "string",
			"maxLength": 200,
		},
		"product_type": map[string]interface{}{
			"required":  true,
			"type":      "string",
			"maxLength": 100,
		},
	}

	ctx.JSON(http.StatusOK, model.Success(rules))
}
