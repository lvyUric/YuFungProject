package service

import (
	"context"
	"fmt"
	"sort"

	"YufungProject/internal/model"
	"YufungProject/internal/repository"
	"YufungProject/pkg/logger"

	"go.mongodb.org/mongo-driver/bson"
)

// MenuService 菜单服务接口
type MenuService interface {
	// 基本CRUD操作
	CreateMenu(ctx context.Context, req *model.MenuCreateRequest) (*model.MenuInfo, error)
	GetMenuByID(ctx context.Context, menuID string) (*model.MenuInfo, error)
	UpdateMenu(ctx context.Context, menuID string, req *model.MenuUpdateRequest) (*model.MenuInfo, error)
	DeleteMenu(ctx context.Context, menuID string) error

	// 查询操作
	GetMenuList(ctx context.Context, req *model.MenuQueryRequest) (*model.MenuListResponse, error)
	GetMenuTree(ctx context.Context, req *model.MenuQueryRequest) ([]model.MenuInfo, error)
	GetUserMenus(ctx context.Context, menuIDs []string) ([]model.UserMenuResponse, error)
	GetUserMenusByRoles(ctx context.Context, roleIDs []string) ([]model.UserMenuResponse, error)

	// 批量操作
	BatchUpdateMenuStatus(ctx context.Context, req *model.BatchUpdateMenuStatusRequest) error

	// 统计操作
	GetMenuStats(ctx context.Context) (*model.MenuStatsResponse, error)
}

// menuService 菜单服务实现
type menuService struct {
	menuRepo repository.MenuRepository
}

// NewMenuService 创建菜单服务实例
func NewMenuService(menuRepo repository.MenuRepository) MenuService {
	return &menuService{
		menuRepo: menuRepo,
	}
}

// CreateMenu 创建菜单
func (s *menuService) CreateMenu(ctx context.Context, req *model.MenuCreateRequest) (*model.MenuInfo, error) {
	// 验证父菜单是否存在
	if req.ParentID != "" {
		_, err := s.menuRepo.GetMenuByID(ctx, req.ParentID)
		if err != nil {
			logger.Error("父菜单不存在", err)
			return nil, fmt.Errorf("父菜单不存在")
		}
	}

	// 验证菜单名称在同一层级下是否唯一
	exists, err := s.menuRepo.CheckMenuNameExists(ctx, req.MenuName, req.ParentID, "")
	if err != nil {
		logger.Error("检查菜单名称是否存在失败", err)
		return nil, fmt.Errorf("检查菜单名称失败: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("同一层级下菜单名称已存在")
	}

	// 如果设置了权限标识符，验证其唯一性
	if req.PermissionCode != "" {
		exists, err := s.menuRepo.CheckPermissionCodeExists(ctx, req.PermissionCode, "")
		if err != nil {
			logger.Error("检查权限标识符是否存在失败", err)
			return nil, fmt.Errorf("检查权限标识符失败: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("权限标识符已存在")
		}
	}

	// 创建菜单实体
	menu := &model.Menu{
		ParentID:       req.ParentID,
		MenuName:       req.MenuName,
		MenuType:       req.MenuType,
		RoutePath:      req.RoutePath,
		Component:      req.Component,
		PermissionCode: req.PermissionCode,
		Icon:           req.Icon,
		SortOrder:      req.SortOrder,
		Visible:        req.Visible,
		Status:         "enable", // 默认启用
	}

	// 如果请求中指定了状态，使用请求的状态
	if req.Status != "" {
		menu.Status = req.Status
	}

	// 创建菜单
	err = s.menuRepo.CreateMenu(ctx, menu)
	if err != nil {
		logger.Error("创建菜单失败", err)
		return nil, fmt.Errorf("创建菜单失败: %w", err)
	}

	// 记录业务日志
	logger.BusinessLog("菜单管理", "创建菜单", "", fmt.Sprintf("创建菜单成功: %s", menu.MenuName))

	return s.convertToMenuInfo(menu), nil
}

// GetMenuByID 根据ID获取菜单
func (s *menuService) GetMenuByID(ctx context.Context, menuID string) (*model.MenuInfo, error) {
	menu, err := s.menuRepo.GetMenuByID(ctx, menuID)
	if err != nil {
		return nil, err
	}

	return s.convertToMenuInfo(menu), nil
}

// UpdateMenu 更新菜单
func (s *menuService) UpdateMenu(ctx context.Context, menuID string, req *model.MenuUpdateRequest) (*model.MenuInfo, error) {
	// 验证菜单是否存在
	existingMenu, err := s.menuRepo.GetMenuByID(ctx, menuID)
	if err != nil {
		return nil, err
	}

	updates := bson.M{}

	// 验证父菜单（如果更改了父菜单）
	if req.ParentID != existingMenu.ParentID {
		// 防止循环引用
		if req.ParentID == menuID {
			return nil, fmt.Errorf("不能将菜单设置为自己的父菜单")
		}

		// 防止将菜单设置为其子菜单的父菜单
		if req.ParentID != "" {
			isChild, err := s.isChildMenu(ctx, menuID, req.ParentID)
			if err != nil {
				return nil, fmt.Errorf("检查菜单层级关系失败: %w", err)
			}
			if isChild {
				return nil, fmt.Errorf("不能将菜单设置为其子菜单的父菜单")
			}
		}

		// 验证父菜单是否存在
		if req.ParentID != "" {
			_, err := s.menuRepo.GetMenuByID(ctx, req.ParentID)
			if err != nil {
				logger.Error("父菜单不存在", err)
				return nil, fmt.Errorf("父菜单不存在")
			}
		}

		updates["parent_id"] = req.ParentID
	}

	// 验证菜单名称在同一层级下是否唯一（排除当前菜单）
	if req.MenuName != "" && req.MenuName != existingMenu.MenuName {
		parentID := req.ParentID
		if parentID == "" {
			parentID = existingMenu.ParentID
		}

		exists, err := s.menuRepo.CheckMenuNameExists(ctx, req.MenuName, parentID, menuID)
		if err != nil {
			logger.Error("检查菜单名称是否存在失败", err)
			return nil, fmt.Errorf("检查菜单名称失败: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("同一层级下菜单名称已存在")
		}
		updates["menu_name"] = req.MenuName
	}

	// 验证权限标识符唯一性（排除当前菜单）
	if req.PermissionCode != existingMenu.PermissionCode {
		if req.PermissionCode != "" {
			exists, err := s.menuRepo.CheckPermissionCodeExists(ctx, req.PermissionCode, menuID)
			if err != nil {
				logger.Error("检查权限标识符是否存在失败", err)
				return nil, fmt.Errorf("检查权限标识符失败: %w", err)
			}
			if exists {
				return nil, fmt.Errorf("权限标识符已存在")
			}
		}
		updates["permission_code"] = req.PermissionCode
	}

	// 更新其他字段
	if req.MenuType != "" {
		updates["menu_type"] = req.MenuType
	}
	if req.RoutePath != existingMenu.RoutePath {
		updates["route_path"] = req.RoutePath
	}
	if req.Component != existingMenu.Component {
		updates["component"] = req.Component
	}
	if req.Icon != existingMenu.Icon {
		updates["icon"] = req.Icon
	}
	if req.SortOrder != 0 {
		updates["sort_order"] = req.SortOrder
	}
	if req.Visible != existingMenu.Visible {
		updates["visible"] = req.Visible
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}

	// 如果没有更新内容，直接返回
	if len(updates) == 0 {
		return s.convertToMenuInfo(existingMenu), nil
	}

	// 执行更新
	err = s.menuRepo.UpdateMenu(ctx, menuID, updates)
	if err != nil {
		logger.Error("更新菜单失败", err)
		return nil, fmt.Errorf("更新菜单失败: %w", err)
	}

	// 记录业务日志
	logger.BusinessLog("菜单管理", "更新菜单", "", fmt.Sprintf("更新菜单成功: %s", existingMenu.MenuName))

	// 返回更新后的菜单信息
	return s.GetMenuByID(ctx, menuID)
}

// DeleteMenu 删除菜单
func (s *menuService) DeleteMenu(ctx context.Context, menuID string) error {
	// 验证菜单是否存在
	menu, err := s.menuRepo.GetMenuByID(ctx, menuID)
	if err != nil {
		return err
	}

	// 检查是否有子菜单
	hasChildren, err := s.menuRepo.HasChildMenus(ctx, menuID)
	if err != nil {
		return fmt.Errorf("检查子菜单失败: %w", err)
	}
	if hasChildren {
		return fmt.Errorf("该菜单下还有子菜单，请先删除子菜单")
	}

	// 删除菜单
	err = s.menuRepo.DeleteMenu(ctx, menuID)
	if err != nil {
		logger.Error("删除菜单失败", err)
		return fmt.Errorf("删除菜单失败: %w", err)
	}

	// 记录业务日志
	logger.BusinessLog("菜单管理", "删除菜单", "", fmt.Sprintf("删除菜单成功: %s", menu.MenuName))

	return nil
}

// GetMenuList 获取菜单列表
func (s *menuService) GetMenuList(ctx context.Context, req *model.MenuQueryRequest) (*model.MenuListResponse, error) {
	// 构建查询过滤器
	filter := bson.M{}

	if req.MenuName != "" {
		filter["menu_name"] = bson.M{"$regex": req.MenuName, "$options": "i"}
	}
	if req.MenuType != "" {
		filter["menu_type"] = req.MenuType
	}
	if req.Status != "" {
		filter["status"] = req.Status
	}
	if req.Visible != nil {
		filter["visible"] = *req.Visible
	}
	if req.PermissionCode != "" {
		filter["permission_code"] = bson.M{"$regex": req.PermissionCode, "$options": "i"}
	}

	// 查询菜单列表
	menus, err := s.menuRepo.GetMenuList(ctx, filter)
	if err != nil {
		logger.Error("查询菜单列表失败", err)
		return nil, fmt.Errorf("查询菜单列表失败: %w", err)
	}

	// 构建树形结构
	menuTree := s.buildMenuTree(menus, "")

	logger.Infof("菜单列表查询成功: Total=%d", len(menus))

	return &model.MenuListResponse{
		Menus: menuTree,
		Total: int64(len(menus)),
	}, nil
}

// GetMenuTree 获取菜单树
func (s *menuService) GetMenuTree(ctx context.Context, req *model.MenuQueryRequest) ([]model.MenuInfo, error) {
	response, err := s.GetMenuList(ctx, req)
	if err != nil {
		return nil, err
	}

	return response.Menus, nil
}

// GetUserMenus 获取用户菜单（用于前端渲染）
func (s *menuService) GetUserMenus(ctx context.Context, menuIDs []string) ([]model.UserMenuResponse, error) {
	if len(menuIDs) == 0 {
		return []model.UserMenuResponse{}, nil
	}

	// 获取用户有权限的菜单
	menus, err := s.menuRepo.GetMenusByIDs(ctx, menuIDs)
	if err != nil {
		return nil, fmt.Errorf("获取用户菜单失败: %w", err)
	}

	// 过滤启用且可见的菜单
	var visibleMenus []model.Menu
	for _, menu := range menus {
		if menu.Status == "enable" && menu.Visible && menu.MenuType != "button" {
			visibleMenus = append(visibleMenus, menu)
		}
	}

	// 构建用户菜单树
	userMenus := s.buildUserMenuTree(visibleMenus, "")

	return userMenus, nil
}

// GetUserMenusByRoles 根据角色ID获取用户菜单
func (s *menuService) GetUserMenusByRoles(ctx context.Context, roleIDs []string) ([]model.UserMenuResponse, error) {
	if len(roleIDs) == 0 {
		return []model.UserMenuResponse{}, nil
	}

	// 获取角色关联的菜单ID
	menuIDs, err := s.menuRepo.GetMenuIDsByRoleIDs(ctx, roleIDs)
	if err != nil {
		return nil, fmt.Errorf("获取角色菜单失败: %w", err)
	}

	// 获取这些菜单的详细信息
	menus, err := s.menuRepo.GetMenusByIDs(ctx, menuIDs)
	if err != nil {
		return nil, fmt.Errorf("获取菜单详细信息失败: %w", err)
	}

	// 过滤启用且可见的菜单
	var visibleMenus []model.Menu
	for _, menu := range menus {
		if menu.Status == "enable" && menu.Visible && menu.MenuType != "button" {
			visibleMenus = append(visibleMenus, menu)
		}
	}

	// 构建用户菜单树
	userMenus := s.buildUserMenuTree(visibleMenus, "")

	return userMenus, nil
}

// BatchUpdateMenuStatus 批量更新菜单状态
func (s *menuService) BatchUpdateMenuStatus(ctx context.Context, req *model.BatchUpdateMenuStatusRequest) error {
	if len(req.MenuIDs) == 0 {
		return fmt.Errorf("菜单ID列表不能为空")
	}

	// 执行批量更新
	err := s.menuRepo.BatchUpdateMenuStatus(ctx, req.MenuIDs, req.Status)
	if err != nil {
		logger.Error("批量更新菜单状态失败", err)
		return fmt.Errorf("批量更新菜单状态失败: %w", err)
	}

	// 记录业务日志
	statusDesc := "启用"
	if req.Status == "disable" {
		statusDesc = "禁用"
	}

	logger.BusinessLog("菜单管理", "批量更新菜单状态", "",
		fmt.Sprintf("批量%s菜单成功: %d个菜单", statusDesc, len(req.MenuIDs)))

	return nil
}

// GetMenuStats 获取菜单统计信息
func (s *menuService) GetMenuStats(ctx context.Context) (*model.MenuStatsResponse, error) {
	stats, err := s.menuRepo.GetMenuStats(ctx)
	if err != nil {
		logger.Error("获取菜单统计信息失败", err)
		return nil, fmt.Errorf("获取菜单统计信息失败: %w", err)
	}

	return stats, nil
}

// 辅助方法

// convertToMenuInfo 转换为菜单信息
func (s *menuService) convertToMenuInfo(menu *model.Menu) *model.MenuInfo {
	return &model.MenuInfo{
		ID:             menu.ID.Hex(),
		MenuID:         menu.MenuID,
		ParentID:       menu.ParentID,
		MenuName:       menu.MenuName,
		MenuType:       menu.MenuType,
		RoutePath:      menu.RoutePath,
		Component:      menu.Component,
		PermissionCode: menu.PermissionCode,
		Icon:           menu.Icon,
		SortOrder:      menu.SortOrder,
		Visible:        menu.Visible,
		Status:         menu.Status,
		CreatedAt:      menu.CreatedAt,
		UpdatedAt:      menu.UpdatedAt,
		Children:       []model.MenuInfo{},
	}
}

// buildMenuTree 构建菜单树
func (s *menuService) buildMenuTree(menus []model.Menu, parentID string) []model.MenuInfo {
	var tree []model.MenuInfo

	for _, menu := range menus {
		if menu.ParentID == parentID {
			menuInfo := s.convertToMenuInfo(&menu)
			menuInfo.Children = s.buildMenuTree(menus, menu.MenuID)
			tree = append(tree, *menuInfo)
		}
	}

	// 按排序号排序
	sort.Slice(tree, func(i, j int) bool {
		return tree[i].SortOrder < tree[j].SortOrder
	})

	return tree
}

// buildUserMenuTree 构建用户菜单树
func (s *menuService) buildUserMenuTree(menus []model.Menu, parentID string) []model.UserMenuResponse {
	var tree []model.UserMenuResponse

	for _, menu := range menus {
		if menu.ParentID == parentID {
			userMenu := model.UserMenuResponse{
				MenuID:    menu.MenuID,
				MenuName:  menu.MenuName,
				RoutePath: menu.RoutePath,
				Component: menu.Component,
				Icon:      menu.Icon,
				SortOrder: menu.SortOrder,
				Children:  s.buildUserMenuTree(menus, menu.MenuID),
			}
			tree = append(tree, userMenu)
		}
	}

	// 按排序号排序
	sort.Slice(tree, func(i, j int) bool {
		return tree[i].SortOrder < tree[j].SortOrder
	})

	return tree
}

// isChildMenu 检查是否为子菜单（递归检查）
func (s *menuService) isChildMenu(ctx context.Context, parentMenuID, checkMenuID string) (bool, error) {
	children, err := s.menuRepo.GetChildMenus(ctx, parentMenuID)
	if err != nil {
		return false, err
	}

	for _, child := range children {
		if child.MenuID == checkMenuID {
			return true, nil
		}

		// 递归检查子菜单
		isChild, err := s.isChildMenu(ctx, child.MenuID, checkMenuID)
		if err != nil {
			return false, err
		}
		if isChild {
			return true, nil
		}
	}

	return false, nil
}
