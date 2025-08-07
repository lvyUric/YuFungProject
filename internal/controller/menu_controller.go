package controller

import (
	"fmt"
	"net/http"

	"YufungProject/internal/model"
	"YufungProject/internal/service"
	"YufungProject/pkg/logger"

	"github.com/gin-gonic/gin"
)

// MenuController 菜单管理控制器
type MenuController struct {
	menuService service.MenuService
}

// NewMenuController 创建菜单管理控制器实例
func NewMenuController(menuService service.MenuService) *MenuController {
	return &MenuController{
		menuService: menuService,
	}
}

// CreateMenu 创建菜单
// @Summary 创建菜单
// @Description 创建新菜单（平台管理员专用）
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Param menu body model.MenuCreateRequest true "菜单信息"
// @Success 200 {object} model.Response{data=model.MenuInfo}
// @Failure 400 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /api/v1/menus [post]
func (mc *MenuController) CreateMenu(c *gin.Context) {
	var req model.MenuCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("参数绑定失败", err)
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: "参数格式错误: " + err.Error(),
		})
		return
	}

	menu, err := mc.menuService.CreateMenu(c.Request.Context(), &req)
	if err != nil {
		logger.Error("创建菜单失败", err)
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    http.StatusInternalServerError,
			Message: "创建菜单失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code:    http.StatusOK,
		Message: "创建菜单成功",
		Data:    menu,
	})
}

// GetMenuByID 根据ID获取菜单
// @Summary 根据ID获取菜单
// @Description 根据菜单ID获取菜单详细信息
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Param id path string true "菜单ID"
// @Success 200 {object} model.Response{data=model.MenuInfo}
// @Failure 400 {object} model.Response
// @Failure 404 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /api/v1/menus/{id} [get]
func (mc *MenuController) GetMenuByID(c *gin.Context) {
	menuID := c.Param("id")
	if menuID == "" {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: "菜单ID不能为空",
		})
		return
	}

	menu, err := mc.menuService.GetMenuByID(c.Request.Context(), menuID)
	if err != nil {
		logger.Error("获取菜单失败", err)
		statusCode := http.StatusInternalServerError
		if err.Error() == "菜单不存在" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, model.Response{
			Code:    statusCode,
			Message: "获取菜单失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code:    http.StatusOK,
		Message: "获取菜单成功",
		Data:    menu,
	})
}

// UpdateMenu 更新菜单
// @Summary 更新菜单
// @Description 更新菜单信息
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Param id path string true "菜单ID"
// @Param menu body model.MenuUpdateRequest true "菜单信息"
// @Success 200 {object} model.Response{data=model.MenuInfo}
// @Failure 400 {object} model.Response
// @Failure 404 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /api/v1/menus/{id} [put]
func (mc *MenuController) UpdateMenu(c *gin.Context) {
	menuID := c.Param("id")
	if menuID == "" {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: "菜单ID不能为空",
		})
		return
	}

	var req model.MenuUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("参数绑定失败", err)
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: "参数格式错误: " + err.Error(),
		})
		return
	}

	menu, err := mc.menuService.UpdateMenu(c.Request.Context(), menuID, &req)
	if err != nil {
		logger.Error("更新菜单失败", err)
		statusCode := http.StatusInternalServerError
		if err.Error() == "菜单不存在" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, model.Response{
			Code:    statusCode,
			Message: "更新菜单失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code:    http.StatusOK,
		Message: "更新菜单成功",
		Data:    menu,
	})
}

// DeleteMenu 删除菜单
// @Summary 删除菜单
// @Description 删除指定菜单
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Param id path string true "菜单ID"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Failure 404 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /api/v1/menus/{id} [delete]
func (mc *MenuController) DeleteMenu(c *gin.Context) {
	menuID := c.Param("id")
	if menuID == "" {
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: "菜单ID不能为空",
		})
		return
	}

	err := mc.menuService.DeleteMenu(c.Request.Context(), menuID)
	if err != nil {
		logger.Error("删除菜单失败", err)
		statusCode := http.StatusInternalServerError
		if err.Error() == "菜单不存在" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, model.Response{
			Code:    statusCode,
			Message: "删除菜单失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code:    http.StatusOK,
		Message: "删除菜单成功",
	})
}

// GetMenuList 获取菜单列表
// @Summary 获取菜单列表
// @Description 获取菜单列表（支持条件查询和树形结构）
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Param menu_name query string false "菜单名称"
// @Param menu_type query string false "菜单类型"
// @Param status query string false "状态"
// @Param visible query bool false "是否可见"
// @Param permission_code query string false "权限标识符"
// @Success 200 {object} model.Response{data=model.MenuListResponse}
// @Failure 500 {object} model.Response
// @Router /api/v1/menus [get]
func (mc *MenuController) GetMenuList(c *gin.Context) {
	var req model.MenuQueryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Error("参数绑定失败", err)
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: "参数格式错误: " + err.Error(),
		})
		return
	}

	menuList, err := mc.menuService.GetMenuList(c.Request.Context(), &req)
	if err != nil {
		logger.Error("获取菜单列表失败", err)
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    http.StatusInternalServerError,
			Message: "获取菜单列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code:    http.StatusOK,
		Message: "获取菜单列表成功",
		Data:    menuList,
	})
}

// GetMenuTree 获取菜单树
// @Summary 获取菜单树
// @Description 获取菜单的层级树形结构
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Param menu_name query string false "菜单名称"
// @Param menu_type query string false "菜单类型"
// @Param status query string false "状态"
// @Success 200 {object} model.Response{data=[]model.MenuInfo} "成功"
// @Failure 400 {object} model.Response "请求参数错误"
// @Failure 500 {object} model.Response "内部服务器错误"
// @Router /api/v1/menu/tree [get]
func (c *MenuController) GetMenuTree(ctx *gin.Context) {
	var req model.MenuQueryRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	// 获取菜单树
	menuTree, err := c.menuService.GetMenuTree(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.Response{
			Code:    http.StatusInternalServerError,
			Message: "获取菜单树失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Code:    http.StatusOK,
		Message: "获取菜单树成功",
		Data:    menuTree,
	})
}

// GetUserMenus 获取用户菜单
// @Summary 获取用户菜单
// @Description 根据用户角色权限获取菜单树（用于前端菜单渲染）
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Success 200 {object} model.Response{data=[]model.UserMenuResponse}
// @Failure 500 {object} model.Response
// @Router /api/v1/menus/user [get]
func (mc *MenuController) GetUserMenus(c *gin.Context) {
	// 从JWT token中获取用户角色ID列表
	roleIDs, exists := c.Get("role_ids")
	if !exists {
		logger.Error("无法获取用户角色信息", nil)
		c.JSON(http.StatusUnauthorized, model.Response{
			Code:    http.StatusUnauthorized,
			Message: "无法获取用户角色信息",
		})
		return
	}

	logger.Debugf("从context获取的roleIDs: %+v, type: %T", roleIDs, roleIDs)

	userRoleIDs, ok := roleIDs.([]string)
	if !ok {
		logger.Errorf("角色ID格式错误: %+v, type: %T", roleIDs, roleIDs)
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    http.StatusInternalServerError,
			Message: "角色ID格式错误",
		})
		return
	}

	logger.Debugf("转换后的userRoleIDs: %+v, length: %d", userRoleIDs, len(userRoleIDs))

	// 如果用户没有角色，返回空菜单
	if len(userRoleIDs) == 0 {
		logger.Warn("用户没有分配任何角色")
		c.JSON(http.StatusOK, model.Response{
			Code:    http.StatusOK,
			Message: "获取用户菜单成功",
			Data:    []model.UserMenuResponse{},
		})
		return
	}

	userMenus, err := mc.menuService.GetUserMenusByRoles(c.Request.Context(), userRoleIDs)
	if err != nil {
		logger.Error("获取用户菜单失败", err)
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    http.StatusInternalServerError,
			Message: "获取用户菜单失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code:    http.StatusOK,
		Message: "获取用户菜单成功",
		Data:    userMenus,
	})
}

// BatchUpdateMenuStatus 批量更新菜单状态
// @Summary 批量更新菜单状态
// @Description 批量启用或禁用菜单
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Param request body model.BatchUpdateMenuStatusRequest true "批量更新请求"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /api/v1/menus/batch-status [put]
func (mc *MenuController) BatchUpdateMenuStatus(c *gin.Context) {
	var req model.BatchUpdateMenuStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("参数绑定失败", err)
		c.JSON(http.StatusBadRequest, model.Response{
			Code:    http.StatusBadRequest,
			Message: "参数格式错误: " + err.Error(),
		})
		return
	}

	err := mc.menuService.BatchUpdateMenuStatus(c.Request.Context(), &req)
	if err != nil {
		logger.Error("批量更新菜单状态失败", err)
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    http.StatusInternalServerError,
			Message: "批量更新菜单状态失败: " + err.Error(),
		})
		return
	}

	statusDesc := "启用"
	if req.Status == "disable" {
		statusDesc = "禁用"
	}

	c.JSON(http.StatusOK, model.Response{
		Code:    http.StatusOK,
		Message: fmt.Sprintf("批量%s菜单成功", statusDesc),
	})
}

// GetMenuStats 获取菜单统计信息
// @Summary 获取菜单统计信息
// @Description 获取菜单相关的统计数据
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Success 200 {object} model.Response{data=model.MenuStatsResponse}
// @Failure 500 {object} model.Response
// @Router /api/v1/menus/stats [get]
func (mc *MenuController) GetMenuStats(c *gin.Context) {
	stats, err := mc.menuService.GetMenuStats(c.Request.Context())
	if err != nil {
		logger.Error("获取菜单统计信息失败", err)
		c.JSON(http.StatusInternalServerError, model.Response{
			Code:    http.StatusInternalServerError,
			Message: "获取菜单统计信息失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Code:    http.StatusOK,
		Message: "获取菜单统计信息成功",
		Data:    stats,
	})
}
