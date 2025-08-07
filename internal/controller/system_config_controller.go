package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"YufungProject/internal/middleware"
	"YufungProject/internal/model"
	"YufungProject/internal/service"
)

// SystemConfigController 系统配置控制器
type SystemConfigController struct {
	systemConfigService service.SystemConfigService
}

// NewSystemConfigController 创建系统配置控制器实例
func NewSystemConfigController(systemConfigService service.SystemConfigService) *SystemConfigController {
	return &SystemConfigController{
		systemConfigService: systemConfigService,
	}
}

// CreateSystemConfig 创建系统配置
// @Summary 创建系统配置
// @Tags 系统配置管理
// @Accept json
// @Produce json
// @Param request body model.SystemConfigCreateRequest true "创建系统配置请求"
// @Success 200 {object} model.Response{data=model.SystemConfigResponse}
// @Router /api/system-configs [post]
func (ctrl *SystemConfigController) CreateSystemConfig(c *gin.Context) {
	var req model.SystemConfigCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	userID, _ := middleware.GetUserID(c)

	config, err := ctrl.systemConfigService.CreateSystemConfig(c.Request.Context(), &req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code:    http.StatusOK,
		Message: "创建成功",
		Data:    config,
	})
}

// GetSystemConfig 获取系统配置详情
// @Summary 获取系统配置详情
// @Tags 系统配置管理
// @Accept json
// @Produce json
// @Param id path string true "系统配置ID"
// @Success 200 {object} model.Response{data=model.SystemConfigResponse}
// @Router /api/system-configs/{id} [get]
func (ctrl *SystemConfigController) GetSystemConfig(c *gin.Context) {
	configID := c.Param("id")
	if configID == "" {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: "系统配置ID不能为空",
		})
		return
	}

	config, err := ctrl.systemConfigService.GetSystemConfigByID(c.Request.Context(), configID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code:    http.StatusOK,
		Message: "获取成功",
		Data:    config,
	})
}

// UpdateSystemConfig 更新系统配置
// @Summary 更新系统配置
// @Tags 系统配置管理
// @Accept json
// @Produce json
// @Param id path string true "系统配置ID"
// @Param request body model.SystemConfigUpdateRequest true "更新系统配置请求"
// @Success 200 {object} model.Response{data=model.SystemConfigResponse}
// @Router /api/system-configs/{id} [put]
func (ctrl *SystemConfigController) UpdateSystemConfig(c *gin.Context) {
	configID := c.Param("id")
	if configID == "" {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: "系统配置ID不能为空",
		})
		return
	}

	var req model.SystemConfigUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	userID, _ := middleware.GetUserID(c)

	config, err := ctrl.systemConfigService.UpdateSystemConfig(c.Request.Context(), configID, &req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code:    http.StatusOK,
		Message: "更新成功",
		Data:    config,
	})
}

// DeleteSystemConfig 删除系统配置
// @Summary 删除系统配置
// @Tags 系统配置管理
// @Accept json
// @Produce json
// @Param id path string true "系统配置ID"
// @Success 200 {object} model.Response
// @Router /api/system-configs/{id} [delete]
func (ctrl *SystemConfigController) DeleteSystemConfig(c *gin.Context) {
	configID := c.Param("id")
	if configID == "" {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: "系统配置ID不能为空",
		})
		return
	}

	err := ctrl.systemConfigService.DeleteSystemConfig(c.Request.Context(), configID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code:    http.StatusOK,
		Message: "删除成功",
	})
}

// ListSystemConfigs 获取系统配置列表
// @Summary 获取系统配置列表
// @Tags 系统配置管理
// @Accept json
// @Produce json
// @Param config_type query string false "配置类型"
// @Param status query string false "状态"
// @Param keyword query string false "关键词"
// @Param page query int false "页码"
// @Param page_size query int false "页大小"
// @Success 200 {object} model.Response{data=model.SystemConfigListResponse}
// @Router /api/system-configs [get]
func (ctrl *SystemConfigController) ListSystemConfigs(c *gin.Context) {
	var req model.SystemConfigQueryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	// 系统配置是全局的，不需要公司ID限制
	result, err := ctrl.systemConfigService.ListSystemConfigs(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code:    http.StatusOK,
		Message: "获取成功",
		Data:    result,
	})
}

// GetConfigOptions 获取配置选项
// @Summary 根据配置类型获取配置选项
// @Tags 系统配置管理
// @Accept json
// @Produce json
// @Param type path string true "配置类型"
// @Success 200 {object} model.Response{data=[]model.SystemConfigResponse}
// @Router /api/system-configs/options/{type} [get]
func (ctrl *SystemConfigController) GetConfigOptions(c *gin.Context) {
	configType := c.Param("type")
	if configType == "" {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: "配置类型不能为空",
		})
		return
	}

	configs, err := ctrl.systemConfigService.GetConfigsByType(c.Request.Context(), configType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code:    http.StatusOK,
		Message: "获取成功",
		Data:    configs,
	})
}
