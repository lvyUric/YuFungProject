package service

import (
	"context"
	"fmt"
	"time"

	"YufungProject/internal/model"
	"YufungProject/internal/repository"
	"YufungProject/pkg/logger"
	"YufungProject/pkg/utils"

	"go.mongodb.org/mongo-driver/bson"
)

// RoleService 角色服务接口
type RoleService interface {
	// 基本CRUD操作
	CreateRole(ctx context.Context, req *model.RoleCreateRequest) (*model.RoleInfo, error)
	GetRoleByID(ctx context.Context, roleID string) (*model.RoleInfo, error)
	UpdateRole(ctx context.Context, roleID string, req *model.RoleUpdateRequest) (*model.RoleInfo, error)
	DeleteRole(ctx context.Context, roleID string) error

	// 查询操作
	GetRoleList(ctx context.Context, req *model.RoleQueryRequest) (*model.RoleListResponse, error)

	// 批量操作
	BatchUpdateRoleStatus(ctx context.Context, req *model.BatchUpdateRoleStatusRequest) error

	// 统计操作
	GetRoleStats(ctx context.Context, companyID string) (*model.RoleStatsResponse, error)
}

// roleService 角色服务实现
type roleService struct {
	roleRepo    repository.RoleRepository
	companyRepo repository.CompanyRepository
	rbacRepo    repository.RBACRepository
}

// NewRoleService 创建角色服务实例
func NewRoleService(roleRepo repository.RoleRepository, companyRepo repository.CompanyRepository, rbacRepo repository.RBACRepository) RoleService {
	return &roleService{
		roleRepo:    roleRepo,
		companyRepo: companyRepo,
		rbacRepo:    rbacRepo,
	}
}

// CreateRole 创建角色
func (s *roleService) CreateRole(ctx context.Context, req *model.RoleCreateRequest) (*model.RoleInfo, error) {
	// 检查角色标识符是否重复
	exists, err := s.roleRepo.CheckRoleKeyExists(ctx, req.RoleKey, "")
	if err != nil {
		logger.Error("检查角色标识符重复失败", err)
		return nil, fmt.Errorf("检查角色标识符重复失败: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("角色标识符 %s 已存在", req.RoleKey)
	}

	// 检查角色名称是否重复
	exists, err = s.roleRepo.CheckRoleNameExists(ctx, req.RoleName, req.CompanyID, "")
	if err != nil {
		logger.Error("检查角色名称重复失败", err)
		return nil, fmt.Errorf("检查角色名称重复失败: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("角色名称 %s 已存在", req.RoleName)
	}

	// 如果指定了公司ID，检查公司是否存在
	if req.CompanyID != "" {
		_, err := s.companyRepo.GetCompanyByID(ctx, req.CompanyID)
		if err != nil {
			return nil, fmt.Errorf("公司不存在: %w", err)
		}
	}

	// 创建角色实体
	role := &model.Role{
		RoleID:    utils.GenerateID("ROLE"),
		RoleName:  req.RoleName,
		RoleKey:   req.RoleKey,
		CompanyID: req.CompanyID,
		SortOrder: req.SortOrder,
		DataScope: req.DataScope,
		MenuIDs:   req.MenuIDs,
		Status:    req.Status,
		Remark:    req.Remark,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 保存角色
	if err := s.roleRepo.CreateRole(ctx, role); err != nil {
		logger.Error("创建角色失败", err)
		return nil, fmt.Errorf("创建角色失败: %w", err)
	}

	// 分配菜单权限
	if len(req.MenuIDs) > 0 {
		if err := s.rbacRepo.AssignPermissionsToRole(ctx, role.RoleID, req.MenuIDs); err != nil {
			logger.Error("分配角色权限失败", err)
			return nil, fmt.Errorf("分配角色权限失败: %w", err)
		}
	}

	// 返回角色信息
	return s.roleModelToInfo(role), nil
}

// GetRoleByID 根据ID获取角色
func (s *roleService) GetRoleByID(ctx context.Context, roleID string) (*model.RoleInfo, error) {
	role, err := s.roleRepo.GetRoleByID(ctx, roleID)
	if err != nil {
		logger.Error("获取角色失败", err)
		return nil, fmt.Errorf("获取角色失败: %w", err)
	}

	roleInfo := s.roleModelToInfo(role)

	// 获取角色权限
	menuIDs, err := s.rbacRepo.GetRolePermissions(ctx, roleID)
	if err != nil {
		logger.Error("获取角色权限失败", err)
		return nil, fmt.Errorf("获取角色权限失败: %w", err)
	}
	roleInfo.MenuIDs = menuIDs

	return roleInfo, nil
}

// UpdateRole 更新角色
func (s *roleService) UpdateRole(ctx context.Context, roleID string, req *model.RoleUpdateRequest) (*model.RoleInfo, error) {
	// 获取现有角色
	existingRole, err := s.roleRepo.GetRoleByID(ctx, roleID)
	if err != nil {
		return nil, fmt.Errorf("角色不存在: %w", err)
	}

	// 检查角色标识符是否重复（排除自己）
	if req.RoleKey != "" && req.RoleKey != existingRole.RoleKey {
		exists, err := s.roleRepo.CheckRoleKeyExists(ctx, req.RoleKey, roleID)
		if err != nil {
			logger.Error("检查角色标识符重复失败", err)
			return nil, fmt.Errorf("检查角色标识符重复失败: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("角色标识符 %s 已存在", req.RoleKey)
		}
	}

	// 检查角色名称是否重复（排除自己）
	if req.RoleName != "" && req.RoleName != existingRole.RoleName {
		exists, err := s.roleRepo.CheckRoleNameExists(ctx, req.RoleName, existingRole.CompanyID, roleID)
		if err != nil {
			logger.Error("检查角色名称重复失败", err)
			return nil, fmt.Errorf("检查角色名称重复失败: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("角色名称 %s 已存在", req.RoleName)
		}
	}

	// 构建更新字段
	updates := bson.M{
		"updated_at": time.Now(),
	}

	if req.RoleName != "" {
		updates["role_name"] = req.RoleName
	}
	if req.RoleKey != "" {
		updates["role_key"] = req.RoleKey
	}
	if req.SortOrder != 0 {
		updates["sort_order"] = req.SortOrder
	}
	if req.DataScope != "" {
		updates["data_scope"] = req.DataScope
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}
	if req.Remark != "" {
		updates["remark"] = req.Remark
	}

	// 更新菜单权限ID数组
	updates["menu_ids"] = req.MenuIDs

	// 更新角色
	if err := s.roleRepo.UpdateRole(ctx, roleID, updates); err != nil {
		logger.Error("更新角色失败", err)
		return nil, fmt.Errorf("更新角色失败: %w", err)
	}

	// 更新菜单权限关联
	if err := s.rbacRepo.AssignPermissionsToRole(ctx, roleID, req.MenuIDs); err != nil {
		logger.Error("更新角色权限失败", err)
		return nil, fmt.Errorf("更新角色权限失败: %w", err)
	}

	// 返回更新后的角色信息
	return s.GetRoleByID(ctx, roleID)
}

// DeleteRole 删除角色
func (s *roleService) DeleteRole(ctx context.Context, roleID string) error {
	// 检查角色是否存在
	_, err := s.roleRepo.GetRoleByID(ctx, roleID)
	if err != nil {
		return fmt.Errorf("角色不存在: %w", err)
	}

	// TODO: 检查是否有用户使用该角色
	// 这个功能需要用户服务配合实现，暂时跳过

	// 删除角色权限关联
	menuIDs, err := s.rbacRepo.GetRolePermissions(ctx, roleID)
	if err != nil {
		logger.Error("获取角色权限失败", err)
		return fmt.Errorf("获取角色权限失败: %w", err)
	}
	if len(menuIDs) > 0 {
		if err := s.rbacRepo.RemovePermissionsFromRole(ctx, roleID, menuIDs); err != nil {
			logger.Error("删除角色权限关联失败", err)
			return fmt.Errorf("删除角色权限关联失败: %w", err)
		}
	}

	// 删除角色
	if err := s.roleRepo.DeleteRole(ctx, roleID); err != nil {
		logger.Error("删除角色失败", err)
		return fmt.Errorf("删除角色失败: %w", err)
	}

	return nil
}

// GetRoleList 获取角色列表
func (s *roleService) GetRoleList(ctx context.Context, req *model.RoleQueryRequest) (*model.RoleListResponse, error) {
	// 构建查询条件
	filter := bson.M{}

	if req.RoleName != "" {
		filter["role_name"] = bson.M{"$regex": req.RoleName, "$options": "i"}
	}
	if req.RoleKey != "" {
		filter["role_key"] = bson.M{"$regex": req.RoleKey, "$options": "i"}
	}
	if req.CompanyID != "" {
		filter["company_id"] = req.CompanyID
	}
	if req.DataScope != "" {
		filter["data_scope"] = req.DataScope
	}
	if req.Status != "" {
		filter["status"] = req.Status
	}

	// 获取角色列表和总数
	roles, total, err := s.roleRepo.GetRoleList(ctx, filter, req.Page, req.PageSize)
	if err != nil {
		logger.Error("获取角色列表失败", err)
		return nil, fmt.Errorf("获取角色列表失败: %w", err)
	}

	// 转换为响应格式
	roleInfos := make([]model.RoleInfo, len(roles))
	for i, role := range roles {
		roleInfo := s.roleModelToInfo(&role)

		// 获取角色权限
		menuIDs, err := s.rbacRepo.GetRolePermissions(ctx, role.RoleID)
		if err != nil {
			logger.Error("获取角色权限失败", err)
			// 不中断流程，只记录错误
			menuIDs = []string{}
		}
		roleInfo.MenuIDs = menuIDs

		roleInfos[i] = *roleInfo
	}

	// 计算分页信息
	totalPages := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))

	return &model.RoleListResponse{
		Roles:      roleInfos,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// BatchUpdateRoleStatus 批量更新角色状态
func (s *roleService) BatchUpdateRoleStatus(ctx context.Context, req *model.BatchUpdateRoleStatusRequest) error {
	if len(req.RoleIDs) == 0 {
		return fmt.Errorf("角色ID列表不能为空")
	}

	err := s.roleRepo.BatchUpdateRoleStatus(ctx, req.RoleIDs, req.Status)
	if err != nil {
		logger.Error("批量更新角色状态失败", err)
		return fmt.Errorf("批量更新角色状态失败: %w", err)
	}

	return nil
}

// GetRoleStats 获取角色统计信息
func (s *roleService) GetRoleStats(ctx context.Context, companyID string) (*model.RoleStatsResponse, error) {
	stats, err := s.roleRepo.GetRoleStats(ctx, companyID)
	if err != nil {
		logger.Error("获取角色统计信息失败", err)
		return nil, fmt.Errorf("获取角色统计信息失败: %w", err)
	}

	return stats, nil
}

// roleModelToInfo 转换角色模型为信息DTO
func (s *roleService) roleModelToInfo(role *model.Role) *model.RoleInfo {
	return &model.RoleInfo{
		ID:        role.ID.Hex(),
		RoleID:    role.RoleID,
		RoleName:  role.RoleName,
		RoleKey:   role.RoleKey,
		CompanyID: role.CompanyID,
		SortOrder: role.SortOrder,
		DataScope: role.DataScope,
		MenuIDs:   role.MenuIDs,
		Status:    role.Status,
		Remark:    role.Remark,
		CreatedAt: role.CreatedAt,
		UpdatedAt: role.UpdatedAt,
	}
}
