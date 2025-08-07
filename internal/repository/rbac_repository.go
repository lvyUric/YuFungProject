package repository

import (
	"context"
	"fmt"
	"time"

	"YufungProject/internal/model"
	"YufungProject/pkg/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// RBACRepository RBAC数据访问接口
type RBACRepository interface {
	// 用户角色关联
	AssignRolesToUser(ctx context.Context, userID string, roleIDs []string) error
	RemoveRolesFromUser(ctx context.Context, userID string, roleIDs []string) error
	GetUserRoles(ctx context.Context, userID string) ([]string, error)
	GetRoleUsers(ctx context.Context, roleID string) ([]string, error)

	// 角色权限关联
	AssignPermissionsToRole(ctx context.Context, roleID string, menuIDs []string) error
	RemovePermissionsFromRole(ctx context.Context, roleID string, menuIDs []string) error
	GetRolePermissions(ctx context.Context, roleID string) ([]string, error)
	GetPermissionRoles(ctx context.Context, menuID string) ([]string, error)

	// 用户权限查询
	GetUserPermissions(ctx context.Context, userID string) ([]string, error)
	CheckUserPermission(ctx context.Context, userID string, permissionCode string) (bool, error)
}

// rbacRepository RBAC数据访问实现
type rbacRepository struct {
	db                       *mongo.Database
	userRoleCollection       *mongo.Collection
	rolePermissionCollection *mongo.Collection
}

// NewRBACRepository 创建RBAC数据访问实例
func NewRBACRepository(db *mongo.Database) RBACRepository {
	repo := &rbacRepository{
		db:                       db,
		userRoleCollection:       db.Collection("user_roles"),
		rolePermissionCollection: db.Collection("role_permissions"),
	}

	// 创建索引
	repo.createIndexes()
	return repo
}

// createIndexes 创建索引
func (r *rbacRepository) createIndexes() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 用户角色表索引
	userRoleIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "user_id", Value: 1}, {Key: "role_id", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("idx_user_role_unique"),
		},
		{
			Keys:    bson.D{{Key: "user_id", Value: 1}},
			Options: options.Index().SetName("idx_user_id"),
		},
		{
			Keys:    bson.D{{Key: "role_id", Value: 1}},
			Options: options.Index().SetName("idx_role_id"),
		},
	}

	_, err := r.userRoleCollection.Indexes().CreateMany(ctx, userRoleIndexes)
	if err != nil {
		logger.Error("创建用户角色表索引失败", err)
	} else {
		logger.Debug("用户角色表索引创建成功")
	}

	// 角色权限表索引
	rolePermissionIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "role_id", Value: 1}, {Key: "menu_id", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("idx_role_permission_unique"),
		},
		{
			Keys:    bson.D{{Key: "role_id", Value: 1}},
			Options: options.Index().SetName("idx_role_id_perm"),
		},
		{
			Keys:    bson.D{{Key: "menu_id", Value: 1}},
			Options: options.Index().SetName("idx_menu_id_perm"),
		},
	}

	_, err = r.rolePermissionCollection.Indexes().CreateMany(ctx, rolePermissionIndexes)
	if err != nil {
		logger.Error("创建角色权限表索引失败", err)
	} else {
		logger.Debug("角色权限表索引创建成功")
	}
}

// AssignRolesToUser 为用户分配角色
func (r *rbacRepository) AssignRolesToUser(ctx context.Context, userID string, roleIDs []string) error {
	startTime := time.Now()

	// 先删除用户现有的角色
	_, err := r.userRoleCollection.DeleteMany(ctx, bson.M{"user_id": userID})
	if err != nil {
		logger.Error("删除用户现有角色失败", err)
		return fmt.Errorf("删除用户现有角色失败: %w", err)
	}

	// 添加新的角色关联
	if len(roleIDs) > 0 {
		var userRoles []interface{}
		now := time.Now()

		for _, roleID := range roleIDs {
			userRole := model.UserRole{
				UserID:    userID,
				RoleID:    roleID,
				CreatedAt: now,
				UpdatedAt: now,
			}
			userRoles = append(userRoles, userRole)
		}

		_, err = r.userRoleCollection.InsertMany(ctx, userRoles)
		if err != nil {
			logger.Error("分配用户角色失败", err)
			return fmt.Errorf("分配用户角色失败: %w", err)
		}
	}

	// 记录数据库操作日志
	logger.DBLog("ASSIGN_ROLES", "user_roles", bson.M{"user_id": userID, "role_ids": roleIDs}, time.Since(startTime))

	return nil
}

// RemoveRolesFromUser 移除用户角色
func (r *rbacRepository) RemoveRolesFromUser(ctx context.Context, userID string, roleIDs []string) error {
	startTime := time.Now()
	filter := bson.M{
		"user_id": userID,
		"role_id": bson.M{"$in": roleIDs},
	}

	_, err := r.userRoleCollection.DeleteMany(ctx, filter)

	// 记录数据库操作日志
	logger.DBLog("REMOVE_ROLES", "user_roles", filter, time.Since(startTime))

	if err != nil {
		logger.Error("移除用户角色失败", err)
		return fmt.Errorf("移除用户角色失败: %w", err)
	}

	return nil
}

// GetUserRoles 获取用户角色
func (r *rbacRepository) GetUserRoles(ctx context.Context, userID string) ([]string, error) {
	startTime := time.Now()
	filter := bson.M{"user_id": userID}

	cursor, err := r.userRoleCollection.Find(ctx, filter)

	// 记录数据库操作日志
	logger.DBLog("GET_USER_ROLES", "user_roles", filter, time.Since(startTime))

	if err != nil {
		logger.Error("查询用户角色失败", err)
		return nil, fmt.Errorf("查询用户角色失败: %w", err)
	}
	defer cursor.Close(ctx)

	var userRoles []model.UserRole
	if err = cursor.All(ctx, &userRoles); err != nil {
		logger.Error("解析用户角色失败", err)
		return nil, fmt.Errorf("解析用户角色失败: %w", err)
	}

	roleIDs := make([]string, len(userRoles))
	for i, ur := range userRoles {
		roleIDs[i] = ur.RoleID
	}

	return roleIDs, nil
}

// GetRoleUsers 获取角色用户
func (r *rbacRepository) GetRoleUsers(ctx context.Context, roleID string) ([]string, error) {
	startTime := time.Now()
	filter := bson.M{"role_id": roleID}

	cursor, err := r.userRoleCollection.Find(ctx, filter)

	// 记录数据库操作日志
	logger.DBLog("GET_ROLE_USERS", "user_roles", filter, time.Since(startTime))

	if err != nil {
		logger.Error("查询角色用户失败", err)
		return nil, fmt.Errorf("查询角色用户失败: %w", err)
	}
	defer cursor.Close(ctx)

	var userRoles []model.UserRole
	if err = cursor.All(ctx, &userRoles); err != nil {
		logger.Error("解析角色用户失败", err)
		return nil, fmt.Errorf("解析角色用户失败: %w", err)
	}

	userIDs := make([]string, len(userRoles))
	for i, ur := range userRoles {
		userIDs[i] = ur.UserID
	}

	return userIDs, nil
}

// AssignPermissionsToRole 为角色分配权限
func (r *rbacRepository) AssignPermissionsToRole(ctx context.Context, roleID string, menuIDs []string) error {
	startTime := time.Now()

	// 先删除角色现有的权限
	_, err := r.rolePermissionCollection.DeleteMany(ctx, bson.M{"role_id": roleID})
	if err != nil {
		logger.Error("删除角色现有权限失败", err)
		return fmt.Errorf("删除角色现有权限失败: %w", err)
	}

	// 添加新的权限关联
	if len(menuIDs) > 0 {
		var rolePermissions []interface{}
		now := time.Now()

		for _, menuID := range menuIDs {
			rolePermission := model.RolePermission{
				RoleID:         roleID,
				MenuID:         menuID,
				PermissionType: "menu", // 默认为菜单权限
				CreatedAt:      now,
				UpdatedAt:      now,
			}
			rolePermissions = append(rolePermissions, rolePermission)
		}

		_, err = r.rolePermissionCollection.InsertMany(ctx, rolePermissions)
		if err != nil {
			logger.Error("分配角色权限失败", err)
			return fmt.Errorf("分配角色权限失败: %w", err)
		}
	}

	// 记录数据库操作日志
	logger.DBLog("ASSIGN_PERMISSIONS", "role_permissions", bson.M{"role_id": roleID, "menu_ids": menuIDs}, time.Since(startTime))

	return nil
}

// RemovePermissionsFromRole 移除角色权限
func (r *rbacRepository) RemovePermissionsFromRole(ctx context.Context, roleID string, menuIDs []string) error {
	startTime := time.Now()
	filter := bson.M{
		"role_id": roleID,
		"menu_id": bson.M{"$in": menuIDs},
	}

	_, err := r.rolePermissionCollection.DeleteMany(ctx, filter)

	// 记录数据库操作日志
	logger.DBLog("REMOVE_PERMISSIONS", "role_permissions", filter, time.Since(startTime))

	if err != nil {
		logger.Error("移除角色权限失败", err)
		return fmt.Errorf("移除角色权限失败: %w", err)
	}

	return nil
}

// GetRolePermissions 获取角色权限
func (r *rbacRepository) GetRolePermissions(ctx context.Context, roleID string) ([]string, error) {
	startTime := time.Now()
	filter := bson.M{"role_id": roleID}

	cursor, err := r.rolePermissionCollection.Find(ctx, filter)

	// 记录数据库操作日志
	logger.DBLog("GET_ROLE_PERMISSIONS", "role_permissions", filter, time.Since(startTime))

	if err != nil {
		logger.Error("查询角色权限失败", err)
		return nil, fmt.Errorf("查询角色权限失败: %w", err)
	}
	defer cursor.Close(ctx)

	var rolePermissions []model.RolePermission
	if err = cursor.All(ctx, &rolePermissions); err != nil {
		logger.Error("解析角色权限失败", err)
		return nil, fmt.Errorf("解析角色权限失败: %w", err)
	}

	menuIDs := make([]string, len(rolePermissions))
	for i, rp := range rolePermissions {
		menuIDs[i] = rp.MenuID
	}

	return menuIDs, nil
}

// GetPermissionRoles 获取权限角色
func (r *rbacRepository) GetPermissionRoles(ctx context.Context, menuID string) ([]string, error) {
	startTime := time.Now()
	filter := bson.M{"menu_id": menuID}

	cursor, err := r.rolePermissionCollection.Find(ctx, filter)

	// 记录数据库操作日志
	logger.DBLog("GET_PERMISSION_ROLES", "role_permissions", filter, time.Since(startTime))

	if err != nil {
		logger.Error("查询权限角色失败", err)
		return nil, fmt.Errorf("查询权限角色失败: %w", err)
	}
	defer cursor.Close(ctx)

	var rolePermissions []model.RolePermission
	if err = cursor.All(ctx, &rolePermissions); err != nil {
		logger.Error("解析权限角色失败", err)
		return nil, fmt.Errorf("解析权限角色失败: %w", err)
	}

	roleIDs := make([]string, len(rolePermissions))
	for i, rp := range rolePermissions {
		roleIDs[i] = rp.RoleID
	}

	return roleIDs, nil
}

// GetUserPermissions 获取用户权限
func (r *rbacRepository) GetUserPermissions(ctx context.Context, userID string) ([]string, error) {
	startTime := time.Now()

	// 聚合查询：用户 -> 角色 -> 权限
	pipeline := []bson.M{
		{
			"$match": bson.M{"user_id": userID},
		},
		{
			"$lookup": bson.M{
				"from":         "role_permissions",
				"localField":   "role_id",
				"foreignField": "role_id",
				"as":           "permissions",
			},
		},
		{
			"$unwind": "$permissions",
		},
		{
			"$group": bson.M{
				"_id": "$permissions.menu_id",
			},
		},
	}

	cursor, err := r.userRoleCollection.Aggregate(ctx, pipeline)

	// 记录数据库操作日志
	logger.DBLog("GET_USER_PERMISSIONS", "user_roles", bson.M{"user_id": userID}, time.Since(startTime))

	if err != nil {
		logger.Error("查询用户权限失败", err)
		return nil, fmt.Errorf("查询用户权限失败: %w", err)
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		logger.Error("解析用户权限失败", err)
		return nil, fmt.Errorf("解析用户权限失败: %w", err)
	}

	menuIDs := make([]string, len(results))
	for i, result := range results {
		menuIDs[i] = result["_id"].(string)
	}

	return menuIDs, nil
}

// CheckUserPermission 检查用户权限
func (r *rbacRepository) CheckUserPermission(ctx context.Context, userID string, permissionCode string) (bool, error) {
	startTime := time.Now()

	// 聚合查询：检查用户是否有特定权限
	pipeline := []bson.M{
		{
			"$match": bson.M{"user_id": userID},
		},
		{
			"$lookup": bson.M{
				"from":         "role_permissions",
				"localField":   "role_id",
				"foreignField": "role_id",
				"as":           "permissions",
			},
		},
		{
			"$unwind": "$permissions",
		},
		{
			"$lookup": bson.M{
				"from":         "menus",
				"localField":   "permissions.menu_id",
				"foreignField": "menu_id",
				"as":           "menu",
			},
		},
		{
			"$unwind": "$menu",
		},
		{
			"$match": bson.M{"menu.permission_code": permissionCode},
		},
		{
			"$limit": 1,
		},
	}

	cursor, err := r.userRoleCollection.Aggregate(ctx, pipeline)

	// 记录数据库操作日志
	logger.DBLog("CHECK_USER_PERMISSION", "user_roles", bson.M{"user_id": userID, "permission_code": permissionCode}, time.Since(startTime))

	if err != nil {
		logger.Error("检查用户权限失败", err)
		return false, fmt.Errorf("检查用户权限失败: %w", err)
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		logger.Error("解析权限检查结果失败", err)
		return false, fmt.Errorf("解析权限检查结果失败: %w", err)
	}

	return len(results) > 0, nil
}
