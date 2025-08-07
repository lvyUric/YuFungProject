package controller

import (
	"fmt"
	"net/http"

	"YufungProject/internal/model"
	"YufungProject/internal/service"
	"YufungProject/pkg/logger"

	"github.com/gin-gonic/gin"
)

// RoleController 角色管理控制器
type RoleController struct {
	roleService service.RoleService
}

// NewRoleController 创建角色控制器实例
func NewRoleController(roleService service.RoleService) *RoleController {
	return &RoleController{
		roleService: roleService,
	}
}

// CreateRole 创建角色
// @Summary 创建角色
// @Description 创建新角色（平台管理员专用）
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param role body model.RoleCreateRequest true "角色信息"
// @Success 200 {object} model.Response{data=model.RoleInfo}
// @Failure 400 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /api/v1/roles [post]
func (rc *RoleController) CreateRole(c *gin.Context) {
	var req model.RoleCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("参数绑定失败", err)
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: "参数格式错误: " + err.Error(),
		})
		return
	}

	// 创建角色
	role, err := rc.roleService.CreateRole(c.Request.Context(), &req)
	if err != nil {
		logger.Error("创建角色失败", err)
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    http.StatusInternalServerError,
			Message: "创建角色失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code:    http.StatusOK,
		Message: "创建角色成功",
		Data:    role,
	})
}

// GetRoleList 获取角色列表
// @Summary 获取角色列表
// @Description 获取角色列表，支持分页和搜索
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页大小" default(10)
// @Param role_name query string false "角色名称搜索"
// @Param role_key query string false "角色标识符搜索"
// @Param company_id query string false "公司ID筛选"
// @Param data_scope query string false "数据权限范围筛选"
// @Param status query string false "状态筛选"
// @Success 200 {object} model.Response{data=model.RoleListResponse}
// @Failure 400 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /api/v1/roles [get]
func (rc *RoleController) GetRoleList(c *gin.Context) {
	var req model.RoleQueryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Error("参数绑定失败", err)
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: "参数格式错误: " + err.Error(),
		})
		return
	}

	// 获取角色列表
	response, err := rc.roleService.GetRoleList(c.Request.Context(), &req)
	if err != nil {
		logger.Error("获取角色列表失败", err)
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    http.StatusInternalServerError,
			Message: "获取角色列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code:    http.StatusOK,
		Message: "获取角色列表成功",
		Data:    response,
	})
}

// GetRoleByID 根据ID获取角色
// @Summary 根据ID获取角色
// @Description 根据角色ID获取角色详细信息
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param id path string true "角色ID"
// @Success 200 {object} model.Response{data=model.RoleInfo}
// @Failure 400 {object} model.Response
// @Failure 404 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /api/v1/roles/{id} [get]
func (rc *RoleController) GetRoleByID(c *gin.Context) {
	roleID := c.Param("id")
	if roleID == "" {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: "角色ID不能为空",
		})
		return
	}

	// 获取角色信息
	role, err := rc.roleService.GetRoleByID(c.Request.Context(), roleID)
	if err != nil {
		logger.Error("获取角色信息失败", err)
		if err.Error() == "角色不存在" {
			c.JSON(http.StatusNotFound, model.Response{
				Code:    http.StatusNotFound,
				Message: "角色不存在",
			})
		} else {
			c.JSON(http.StatusInternalServerError, model.Response{
				Code:    http.StatusInternalServerError,
				Message: "获取角色信息失败: " + err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code:    http.StatusOK,
		Message: "获取角色信息成功",
		Data:    role,
	})
}

// UpdateRole 更新角色
// @Summary 更新角色
// @Description 更新角色信息
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param id path string true "角色ID"
// @Param role body model.RoleUpdateRequest true "角色更新信息"
// @Success 200 {object} model.Response{data=model.RoleInfo}
// @Failure 400 {object} model.Response
// @Failure 404 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /api/v1/roles/{id} [put]
func (rc *RoleController) UpdateRole(c *gin.Context) {
	roleID := c.Param("id")
	if roleID == "" {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: "角色ID不能为空",
		})
		return
	}

	var req model.RoleUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("参数绑定失败", err)
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: "参数格式错误: " + err.Error(),
		})
		return
	}

	// 更新角色
	role, err := rc.roleService.UpdateRole(c.Request.Context(), roleID, &req)
	if err != nil {
		logger.Error("更新角色失败", err)
		if err.Error() == "角色不存在" {
			c.JSON(http.StatusNotFound, model.Response{
				Code:    http.StatusNotFound,
				Message: "角色不存在",
			})
		} else {
			c.JSON(http.StatusInternalServerError, model.Response{
				Code:    http.StatusInternalServerError,
				Message: "更新角色失败: " + err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code:    http.StatusOK,
		Message: "更新角色成功",
		Data:    role,
	})
}

// DeleteRole 删除角色
// @Summary 删除角色
// @Description 删除指定角色
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param id path string true "角色ID"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Failure 404 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /api/v1/roles/{id} [delete]
func (rc *RoleController) DeleteRole(c *gin.Context) {
	roleID := c.Param("id")
	if roleID == "" {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: "角色ID不能为空",
		})
		return
	}

	// 删除角色
	err := rc.roleService.DeleteRole(c.Request.Context(), roleID)
	if err != nil {
		logger.Error("删除角色失败", err)
		if err.Error() == "角色不存在" {
			c.JSON(http.StatusNotFound, model.Response{
				Code:    http.StatusNotFound,
				Message: "角色不存在",
			})
		} else {
			c.JSON(http.StatusInternalServerError, model.Response{
				Code:    http.StatusInternalServerError,
				Message: "删除角色失败: " + err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code:    http.StatusOK,
		Message: "删除角色成功",
	})
}

// BatchUpdateRoleStatus 批量更新角色状态
// @Summary 批量更新角色状态
// @Description 批量启用或禁用角色
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param request body model.BatchUpdateRoleStatusRequest true "批量更新请求"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /api/v1/roles/batch-status [put]
func (rc *RoleController) BatchUpdateRoleStatus(c *gin.Context) {
	var req model.BatchUpdateRoleStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("参数绑定失败", err)
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: "参数格式错误: " + err.Error(),
		})
		return
	}

	// 批量更新角色状态
	err := rc.roleService.BatchUpdateRoleStatus(c.Request.Context(), &req)
	if err != nil {
		logger.Error("批量更新角色状态失败", err)
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    http.StatusInternalServerError,
			Message: "批量更新角色状态失败: " + err.Error(),
		})
		return
	}

	statusDesc := "启用"
	if req.Status == "disable" {
		statusDesc = "禁用"
	}

	c.JSON(http.StatusOK, model.Response{
		Code:    http.StatusOK,
		Message: fmt.Sprintf("批量%s角色成功", statusDesc),
	})
}

// GetRoleStats 获取角色统计信息
// @Summary 获取角色统计信息
// @Description 获取角色统计数据，包括总数、启用数、禁用数等
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param company_id query string false "公司ID"
// @Success 200 {object} model.Response{data=model.RoleStatsResponse}
// @Failure 500 {object} model.Response
// @Router /api/v1/roles/stats [get]
func (rc *RoleController) GetRoleStats(c *gin.Context) {
	companyID := c.Query("company_id")

	// 获取角色统计信息
	stats, err := rc.roleService.GetRoleStats(c.Request.Context(), companyID)
	if err != nil {
		logger.Error("获取角色统计信息失败", err)
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    http.StatusInternalServerError,
			Message: "获取角色统计信息失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code:    http.StatusOK,
		Message: "获取角色统计信息成功",
		Data:    stats,
	})
}

// GetRolesByCompanyID 根据公司ID获取角色列表
// @Summary 根据公司ID获取角色列表
// @Description 获取指定公司的所有可用角色（包括平台角色）
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param company_id path string true "公司ID"
// @Success 200 {object} model.Response{data=[]model.RoleInfo}
// @Failure 400 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /api/v1/roles/company/{company_id} [get]
func (rc *RoleController) GetRolesByCompanyID(c *gin.Context) {
	companyID := c.Param("company_id")
	if companyID == "" {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: "公司ID不能为空",
		})
		return
	}

	// 构建查询参数，获取指定公司和平台级角色
	req := &model.RoleQueryRequest{
		Page:      1,
		PageSize:  1000, // 设置一个较大的值来获取所有角色
		CompanyID: companyID,
		Status:    "enable", // 只获取启用的角色
	}

	// 获取公司角色列表
	response, err := rc.roleService.GetRoleList(c.Request.Context(), req)
	if err != nil {
		logger.Error("获取公司角色列表失败", err)
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    http.StatusInternalServerError,
			Message: "获取公司角色列表失败: " + err.Error(),
		})
		return
	}

	// 同时获取平台级角色（company_id为空）
	platformReq := &model.RoleQueryRequest{
		Page:      1,
		PageSize:  1000,
		CompanyID: "", // 平台级角色
		Status:    "enable",
	}

	platformResponse, err := rc.roleService.GetRoleList(c.Request.Context(), platformReq)
	if err != nil {
		logger.Error("获取平台角色列表失败", err)
		// 平台角色获取失败不影响结果，只记录错误
		platformResponse = &model.RoleListResponse{Roles: []model.RoleInfo{}}
	}

	// 合并角色列表
	allRoles := append(response.Roles, platformResponse.Roles...)

	c.JSON(http.StatusOK, model.Response{
		Code:    http.StatusOK,
		Message: "获取公司角色列表成功",
		Data:    allRoles,
	})
}
