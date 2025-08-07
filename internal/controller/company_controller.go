package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"YufungProject/internal/model"
	"YufungProject/internal/service"
	"YufungProject/pkg/logger"

	"github.com/gin-gonic/gin"
)

// CompanyController 公司控制器
type CompanyController struct {
	companyService service.CompanyService
}

// NewCompanyController 创建公司控制器实例
func NewCompanyController(companyService service.CompanyService) *CompanyController {
	return &CompanyController{
		companyService: companyService,
	}
}

// CreateCompany 创建公司
//
//	@Summary		创建公司
//	@Description	创建新的保险经纪公司
//	@Tags			公司管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string							true	"Bearer JWT令牌"
//	@Param			request			body		model.CreateCompanyRequest		true	"创建公司请求参数"
//	@Success		201				{object}	model.Response{data=model.CompanyInfo}	"创建成功"
//	@Failure		400				{object}	model.Response{data=string}				"请求参数错误"
//	@Failure		409				{object}	model.Response{data=string}				"公司名称已存在"
//	@Failure		500				{object}	model.Response{data=string}				"服务器内部错误"
//	@Router			/company [post]
func (c *CompanyController) CreateCompany(ctx *gin.Context) {
	var req model.CreateCompanyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Warnf("创建公司请求参数错误: %v", err)
		ctx.JSON(http.StatusBadRequest, model.ErrorResponse(model.CodeInvalidParams, "请求参数错误", err.Error()))
		return
	}

	// 记录操作日志
	userID, _ := ctx.Get("user_id")
	logger.BusinessLog("公司管理", "创建公司", userID.(string), "开始创建公司: "+req.CompanyName)

	company, err := c.companyService.CreateCompany(ctx, &req)
	if err != nil {
		logger.Errorf("创建公司失败: %v", err)
		switch err.Error() {
		case "公司名称已存在":
			ctx.JSON(http.StatusConflict, model.ErrorResponse(model.CodeCompanyExists, "公司名称已存在", nil))
		default:
			ctx.JSON(http.StatusInternalServerError, model.ErrorResponse(model.CodeServerError, "创建公司失败", err.Error()))
		}
		return
	}

	ctx.JSON(http.StatusCreated, model.SuccessResponse("公司创建成功", company))
}

// GetCompanyByID 获取公司详情
//
//	@Summary		获取公司详情
//	@Description	根据公司ID获取公司详细信息
//	@Tags			公司管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string						true	"Bearer JWT令牌"
//	@Param			id				path		string						true	"公司ID"
//	@Success		200				{object}	model.Response{data=model.CompanyInfo}	"查询成功"
//	@Failure		404				{object}	model.Response{data=string}				"公司不存在"
//	@Failure		500				{object}	model.Response{data=string}				"服务器内部错误"
//	@Router			/company/{id} [get]
func (c *CompanyController) GetCompanyByID(ctx *gin.Context) {
	companyID := ctx.Param("id")
	if companyID == "" {
		ctx.JSON(http.StatusBadRequest, model.ErrorResponse(model.CodeInvalidParams, "公司ID不能为空", nil))
		return
	}

	company, err := c.companyService.GetCompanyByID(ctx, companyID)
	if err != nil {
		logger.Errorf("查询公司详情失败: %v", err)
		switch err.Error() {
		case "公司不存在":
			ctx.JSON(http.StatusNotFound, model.ErrorResponse(model.CodeCompanyNotExists, "公司不存在", nil))
		default:
			ctx.JSON(http.StatusInternalServerError, model.ErrorResponse(model.CodeServerError, "查询公司失败", err.Error()))
		}
		return
	}

	ctx.JSON(http.StatusOK, model.SuccessResponse("查询成功", company))
}

// GetCompanyList 获取公司列表
//
//	@Summary		获取公司列表
//	@Description	分页查询公司列表，支持状态筛选和关键词搜索
//	@Tags			公司管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string							true	"Bearer JWT令牌"
//	@Param			page			query		int								false	"页码"
//	@Param			page_size		query		int								false	"每页大小"
//	@Param			status			query		string							false	"状态筛选"
//	@Param			keyword			query		string							false	"关键词搜索"
//	@Success		200				{object}	model.Response{data=model.CompanyListResponse}	"查询成功"
//	@Failure		400				{object}	model.Response{data=string}						"请求参数错误"
//	@Failure		500				{object}	model.Response{data=string}						"服务器内部错误"
//	@Router			/company [get]
func (c *CompanyController) GetCompanyList(ctx *gin.Context) {
	var req model.CompanyQueryRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		logger.Warnf("公司列表查询参数错误: %v", err)
		ctx.JSON(http.StatusBadRequest, model.ErrorResponse(model.CodeInvalidParams, "请求参数错误", err.Error()))
		return
	}

	response, err := c.companyService.GetCompanyList(ctx, &req)
	if err != nil {
		logger.Errorf("查询公司列表失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, model.ErrorResponse(model.CodeServerError, "查询公司列表失败", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, model.SuccessResponse("查询成功", response))
}

// UpdateCompany 更新公司
//
//	@Summary		更新公司
//	@Description	更新公司信息
//	@Tags			公司管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string							true	"Bearer JWT令牌"
//	@Param			id				path		string							true	"公司ID"
//	@Param			request			body		model.UpdateCompanyRequest		true	"更新公司请求参数"
//	@Success		200				{object}	model.Response{data=model.CompanyInfo}	"更新成功"
//	@Failure		400				{object}	model.Response{data=string}				"请求参数错误"
//	@Failure		404				{object}	model.Response{data=string}				"公司不存在"
//	@Failure		409				{object}	model.Response{data=string}				"公司名称已存在"
//	@Failure		500				{object}	model.Response{data=string}				"服务器内部错误"
//	@Router			/company/{id} [put]
func (c *CompanyController) UpdateCompany(ctx *gin.Context) {
	companyID := ctx.Param("id")
	if companyID == "" {
		ctx.JSON(http.StatusBadRequest, model.ErrorResponse(model.CodeInvalidParams, "公司ID不能为空", nil))
		return
	}

	var req model.UpdateCompanyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Warnf("更新公司请求参数错误: %v", err)
		ctx.JSON(http.StatusBadRequest, model.ErrorResponse(model.CodeInvalidParams, "请求参数错误", err.Error()))
		return
	}

	// 记录操作日志
	userID, _ := ctx.Get("user_id")
	logger.BusinessLog("公司管理", "更新公司", userID.(string), "更新公司: "+companyID)

	company, err := c.companyService.UpdateCompany(ctx, companyID, &req)
	if err != nil {
		logger.Errorf("更新公司失败: %v", err)
		switch err.Error() {
		case "公司不存在":
			ctx.JSON(http.StatusNotFound, model.ErrorResponse(model.CodeCompanyNotExists, "公司不存在", nil))
		case "公司名称已存在":
			ctx.JSON(http.StatusConflict, model.ErrorResponse(model.CodeCompanyExists, "公司名称已存在", nil))
		default:
			ctx.JSON(http.StatusInternalServerError, model.ErrorResponse(model.CodeServerError, "更新公司失败", err.Error()))
		}
		return
	}

	ctx.JSON(http.StatusOK, model.SuccessResponse("更新成功", company))
}

// DeleteCompany 删除公司
//
//	@Summary		删除公司
//	@Description	删除指定公司（需要确保公司下没有用户）
//	@Tags			公司管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string					true	"Bearer JWT令牌"
//	@Param			id				path		string					true	"公司ID"
//	@Success		200				{object}	model.Response{data=string}	"删除成功"
//	@Failure		400				{object}	model.Response{data=string}	"请求参数错误"
//	@Failure		404				{object}	model.Response{data=string}	"公司不存在"
//	@Failure		409				{object}	model.Response{data=string}	"公司下还有用户，无法删除"
//	@Failure		500				{object}	model.Response{data=string}	"服务器内部错误"
//	@Router			/company/{id} [delete]
func (c *CompanyController) DeleteCompany(ctx *gin.Context) {
	companyID := ctx.Param("id")
	if companyID == "" {
		ctx.JSON(http.StatusBadRequest, model.ErrorResponse(model.CodeInvalidParams, "公司ID不能为空", nil))
		return
	}

	// 记录操作日志
	userID, _ := ctx.Get("user_id")
	logger.BusinessLog("公司管理", "删除公司", userID.(string), "删除公司: "+companyID)

	err := c.companyService.DeleteCompany(ctx, companyID)
	if err != nil {
		logger.Errorf("删除公司失败: %v", err)
		switch {
		case err.Error() == "公司不存在":
			ctx.JSON(http.StatusNotFound, model.ErrorResponse(model.CodeCompanyNotExists, "公司不存在", nil))
		case err.Error() == "公司下还有用户，无法删除" ||
			err.Error() == "公司下还有 1 个用户，无法删除":
			ctx.JSON(http.StatusConflict, model.ErrorResponse(model.CodeConflict, err.Error(), nil))
		default:
			ctx.JSON(http.StatusInternalServerError, model.ErrorResponse(model.CodeServerError, "删除公司失败", err.Error()))
		}
		return
	}

	ctx.JSON(http.StatusOK, model.SuccessResponse("删除成功", nil))
}

// GetCompanyStats 获取公司统计
//
//	@Summary		获取公司统计
//	@Description	获取公司相关统计数据
//	@Tags			公司管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string						true	"Bearer JWT令牌"
//	@Success		200				{object}	model.Response{data=model.CompanyStatsResponse}	"查询成功"
//	@Failure		500				{object}	model.Response{data=string}						"服务器内部错误"
//	@Router			/company/stats [get]
func (c *CompanyController) GetCompanyStats(ctx *gin.Context) {
	stats, err := c.companyService.GetCompanyStats(ctx)
	if err != nil {
		logger.Errorf("查询公司统计失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, model.ErrorResponse(model.CodeServerError, "查询统计失败", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, model.SuccessResponse("查询成功", stats))
}

// ExportCompany 导出公司数据
//
//	@Summary		导出公司数据
//	@Description	导出公司数据到Excel或CSV文件
//	@Tags			公司管理
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string						true	"Bearer JWT令牌"
//	@Param			request			body		model.CompanyExportRequest	true	"导出请求参数"
//	@Success		200				{object}	model.Response{data=model.CompanyExportResponse}	"导出成功"
//	@Failure		400				{object}	model.Response{data=string}							"请求参数错误"
//	@Failure		500				{object}	model.Response{data=string}							"服务器内部错误"
//	@Router			/company/export [post]
func (c *CompanyController) ExportCompany(ctx *gin.Context) {
	var req model.CompanyExportRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Warnf("导出公司请求参数错误: %v", err)
		ctx.JSON(http.StatusBadRequest, model.ErrorResponse(model.CodeInvalidParams, "请求参数错误", err.Error()))
		return
	}

	// 记录操作日志
	userID, _ := ctx.Get("user_id")
	logger.BusinessLog("公司管理", "导出公司", userID.(string), fmt.Sprintf("导出类型: %s, 格式: %s", req.ExportType, req.Format))

	response, err := c.companyService.ExportCompany(ctx, &req)
	if err != nil {
		logger.Errorf("导出公司数据失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, model.ErrorResponse(model.CodeServerError, "导出失败", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, model.SuccessResponse("导出成功", response))
}

// DownloadTemplate 下载导入模板
//
//	@Summary		下载导入模板
//	@Description	下载公司数据导入模板文件
//	@Tags			公司管理
//	@Accept			json
//	@Produce		application/octet-stream
//	@Param			Authorization	header		string	true	"Bearer JWT令牌"
//	@Param			format			query		string	false	"文件格式：xlsx, csv"	default(xlsx)
//	@Success		200				{file}		binary	"模板文件"
//	@Failure		400				{object}	model.Response{data=string}	"请求参数错误"
//	@Failure		500				{object}	model.Response{data=string}	"服务器内部错误"
//	@Router			/company/template [get]
func (c *CompanyController) DownloadTemplate(ctx *gin.Context) {
	format := ctx.DefaultQuery("format", "xlsx")
	if format != "xlsx" && format != "csv" {
		ctx.JSON(http.StatusBadRequest, model.ErrorResponse(model.CodeInvalidParams, "不支持的文件格式", nil))
		return
	}

	// 记录操作日志
	userID, _ := ctx.Get("user_id")
	logger.BusinessLog("公司管理", "下载模板", userID.(string), "下载格式: "+format)

	fileData, fileName, err := c.companyService.GenerateTemplate(ctx, format)
	if err != nil {
		logger.Errorf("生成模板失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, model.ErrorResponse(model.CodeServerError, "生成模板失败", err.Error()))
		return
	}

	// 设置响应头
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	ctx.Header("Content-Length", strconv.Itoa(len(fileData)))

	ctx.Data(http.StatusOK, "application/octet-stream", fileData)
}

// PreviewImport 预览导入数据
//
//	@Summary		预览导入数据
//	@Description	预览导入的公司数据，检查格式和错误
//	@Tags			公司管理
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			Authorization	header		string	true	"Bearer JWT令牌"
//	@Param			file			formData	file	true	"导入文件"
//	@Param			skip_header		formData	bool	false	"是否跳过表头行"
//	@Param			update_existing	formData	bool	false	"是否更新已存在的公司"
//	@Success		200				{object}	model.Response{data=model.CompanyImportResponse}	"预览成功"
//	@Failure		400				{object}	model.Response{data=string}							"请求参数错误"
//	@Failure		500				{object}	model.Response{data=string}							"服务器内部错误"
//	@Router			/company/import/preview [post]
func (c *CompanyController) PreviewImport(ctx *gin.Context) {
	// 获取上传文件
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		logger.Warnf("获取上传文件失败: %v", err)
		ctx.JSON(http.StatusBadRequest, model.ErrorResponse(model.CodeInvalidParams, "请选择要上传的文件", err.Error()))
		return
	}
	defer file.Close()

	// 获取其他参数
	var req model.CompanyImportRequest
	if err := ctx.ShouldBind(&req); err != nil {
		logger.Warnf("导入预览请求参数错误: %v", err)
		ctx.JSON(http.StatusBadRequest, model.ErrorResponse(model.CodeInvalidParams, "请求参数错误", err.Error()))
		return
	}

	// 记录操作日志
	userID, _ := ctx.Get("user_id")
	logger.BusinessLog("公司管理", "预览导入", userID.(string), "文件名: "+header.Filename)

	response, err := c.companyService.PreviewImport(ctx, file, header, &req)
	if err != nil {
		logger.Errorf("预览导入失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, model.ErrorResponse(model.CodeServerError, "预览失败", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, model.SuccessResponse("预览成功", response))
}

// ImportCompany 导入公司数据
//
//	@Summary		导入公司数据
//	@Description	导入公司数据到系统中
//	@Tags			公司管理
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			Authorization	header		string	true	"Bearer JWT令牌"
//	@Param			file			formData	file	true	"导入文件"
//	@Param			skip_header		formData	bool	false	"是否跳过表头行"
//	@Param			update_existing	formData	bool	false	"是否更新已存在的公司"
//	@Success		200				{object}	model.Response{data=model.CompanyImportResponse}	"导入成功"
//	@Failure		400				{object}	model.Response{data=string}							"请求参数错误"
//	@Failure		500				{object}	model.Response{data=string}							"服务器内部错误"
//	@Router			/company/import [post]
func (c *CompanyController) ImportCompany(ctx *gin.Context) {
	// 获取上传文件
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		logger.Warnf("获取上传文件失败: %v", err)
		ctx.JSON(http.StatusBadRequest, model.ErrorResponse(model.CodeInvalidParams, "请选择要上传的文件", err.Error()))
		return
	}
	defer file.Close()

	// 获取其他参数
	var req model.CompanyImportRequest
	if err := ctx.ShouldBind(&req); err != nil {
		logger.Warnf("导入请求参数错误: %v", err)
		ctx.JSON(http.StatusBadRequest, model.ErrorResponse(model.CodeInvalidParams, "请求参数错误", err.Error()))
		return
	}

	// 记录操作日志
	userID, _ := ctx.Get("user_id")
	logger.BusinessLog("公司管理", "导入公司", userID.(string), "文件名: "+header.Filename)

	response, err := c.companyService.ImportCompany(ctx, file, header, &req)
	if err != nil {
		logger.Errorf("导入公司数据失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, model.ErrorResponse(model.CodeServerError, "导入失败", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, model.SuccessResponse("导入完成", response))
}
