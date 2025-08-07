package repository

import (
	"context"
	"fmt"
	"time"

	"YufungProject/internal/model"
	"YufungProject/pkg/logger"
	"YufungProject/pkg/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// RoleRepository 角色数据访问接口
type RoleRepository interface {
	// 基本CRUD操作
	CreateRole(ctx context.Context, role *model.Role) error
	GetRoleByID(ctx context.Context, roleID string) (*model.Role, error)
	GetRoleByKey(ctx context.Context, roleKey string) (*model.Role, error)
	UpdateRole(ctx context.Context, roleID string, updates bson.M) error
	DeleteRole(ctx context.Context, roleID string) error

	// 查询操作
	GetRoleList(ctx context.Context, filter bson.M, page, pageSize int) ([]model.Role, int64, error)
	GetRolesByCompanyID(ctx context.Context, companyID string) ([]model.Role, error)
	GetRolesByIDs(ctx context.Context, roleIDs []string) ([]model.Role, error)

	// 批量操作
	BatchUpdateRoleStatus(ctx context.Context, roleIDs []string, status string) error

	// 统计操作
	GetRoleStats(ctx context.Context, companyID string) (*model.RoleStatsResponse, error)
	CountRolesByCompanyID(ctx context.Context, companyID string) (int64, error)

	// 验证操作
	CheckRoleKeyExists(ctx context.Context, roleKey string, excludeRoleID string) (bool, error)
	CheckRoleNameExists(ctx context.Context, roleName string, companyID string, excludeRoleID string) (bool, error)
}

// roleRepository 角色数据访问实现
type roleRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

// NewRoleRepository 创建角色数据访问实例
func NewRoleRepository(db *mongo.Database) RoleRepository {
	repo := &roleRepository{
		db:         db,
		collection: db.Collection("roles"),
	}

	// 创建索引
	repo.createIndexes()
	return repo
}

// createIndexes 创建索引
func (r *roleRepository) createIndexes() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "role_id", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("idx_role_id_unique"),
		},
		{
			Keys:    bson.D{{Key: "role_key", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("idx_role_key_unique"),
		},
		{
			Keys:    bson.D{{Key: "company_id", Value: 1}},
			Options: options.Index().SetName("idx_company_id"),
		},
		{
			Keys:    bson.D{{Key: "status", Value: 1}},
			Options: options.Index().SetName("idx_status"),
		},
		{
			Keys:    bson.D{{Key: "created_at", Value: -1}},
			Options: options.Index().SetName("idx_created_at"),
		},
	}

	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		logger.Error("创建角色表索引失败", err)
	} else {
		logger.Debug("角色表索引创建成功")
	}
}

// CreateRole 创建角色
func (r *roleRepository) CreateRole(ctx context.Context, role *model.Role) error {
	startTime := time.Now()

	// 设置创建时间
	role.CreatedAt = time.Now()
	role.UpdatedAt = time.Now()

	// 生成角色ID
	if role.RoleID == "" {
		role.RoleID = utils.GenerateID("ROLE")
	}

	_, err := r.collection.InsertOne(ctx, role)

	// 记录数据库操作日志
	logger.DBLog("CREATE", "roles", bson.M{"role_id": role.RoleID}, time.Since(startTime))

	if err != nil {
		logger.Error("创建角色失败", err)
		return fmt.Errorf("创建角色失败: %w", err)
	}

	return nil
}

// GetRoleByID 根据ID获取角色
func (r *roleRepository) GetRoleByID(ctx context.Context, roleID string) (*model.Role, error) {
	startTime := time.Now()
	filter := bson.M{"role_id": roleID}

	var role model.Role
	err := r.collection.FindOne(ctx, filter).Decode(&role)

	// 记录数据库操作日志
	logger.DBLog("FIND", "roles", filter, time.Since(startTime))

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("角色不存在")
		}
		logger.Error("查询角色失败", err)
		return nil, fmt.Errorf("查询角色失败: %w", err)
	}

	return &role, nil
}

// GetRoleByKey 根据角色标识符获取角色
func (r *roleRepository) GetRoleByKey(ctx context.Context, roleKey string) (*model.Role, error) {
	startTime := time.Now()
	filter := bson.M{"role_key": roleKey}

	var role model.Role
	err := r.collection.FindOne(ctx, filter).Decode(&role)

	// 记录数据库操作日志
	logger.DBLog("FIND", "roles", filter, time.Since(startTime))

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("角色不存在")
		}
		logger.Error("查询角色失败", err)
		return nil, fmt.Errorf("查询角色失败: %w", err)
	}

	return &role, nil
}

// UpdateRole 更新角色
func (r *roleRepository) UpdateRole(ctx context.Context, roleID string, updates bson.M) error {
	startTime := time.Now()
	filter := bson.M{"role_id": roleID}

	// 添加更新时间
	updates["updated_at"] = time.Now()

	update := bson.M{"$set": updates}

	result, err := r.collection.UpdateOne(ctx, filter, update)

	// 记录数据库操作日志
	logger.DBLog("UPDATE", "roles", filter, time.Since(startTime))

	if err != nil {
		logger.Error("更新角色失败", err)
		return fmt.Errorf("更新角色失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("角色不存在")
	}

	return nil
}

// DeleteRole 删除角色
func (r *roleRepository) DeleteRole(ctx context.Context, roleID string) error {
	startTime := time.Now()
	filter := bson.M{"role_id": roleID}

	result, err := r.collection.DeleteOne(ctx, filter)

	// 记录数据库操作日志
	logger.DBLog("DELETE", "roles", filter, time.Since(startTime))

	if err != nil {
		logger.Error("删除角色失败", err)
		return fmt.Errorf("删除角色失败: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("角色不存在")
	}

	return nil
}

// GetRoleList 获取角色列表
func (r *roleRepository) GetRoleList(ctx context.Context, filter bson.M, page, pageSize int) ([]model.Role, int64, error) {
	startTime := time.Now()

	// 计算跳过的记录数
	skip := (page - 1) * pageSize

	// 设置查询选项
	findOptions := options.Find().
		SetSort(bson.D{{Key: "sort_order", Value: 1}, {Key: "created_at", Value: -1}}).
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize))

	// 执行查询
	cursor, err := r.collection.Find(ctx, filter, findOptions)

	// 记录数据库操作日志
	logger.DBLog("FIND_LIST", "roles", filter, time.Since(startTime))

	if err != nil {
		logger.Error("查询角色列表失败", err)
		return nil, 0, fmt.Errorf("查询角色列表失败: %w", err)
	}
	defer cursor.Close(ctx)

	var roles []model.Role
	if err = cursor.All(ctx, &roles); err != nil {
		logger.Error("解析角色列表失败", err)
		return nil, 0, fmt.Errorf("解析角色列表失败: %w", err)
	}

	// 获取总数
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		logger.Error("统计角色总数失败", err)
		return nil, 0, fmt.Errorf("统计角色总数失败: %w", err)
	}

	return roles, total, nil
}

// GetRolesByCompanyID 根据公司ID获取角色列表
func (r *roleRepository) GetRolesByCompanyID(ctx context.Context, companyID string) ([]model.Role, error) {
	startTime := time.Now()

	// 查询条件：公司角色 + 平台角色
	filter := bson.M{
		"$or": []bson.M{
			{"company_id": companyID},
			{"company_id": ""},
		},
		"status": "enable",
	}

	cursor, err := r.collection.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "sort_order", Value: 1}}))

	// 记录数据库操作日志
	logger.DBLog("FIND_LIST", "roles", filter, time.Since(startTime))

	if err != nil {
		logger.Error("查询公司角色列表失败", err)
		return nil, fmt.Errorf("查询公司角色列表失败: %w", err)
	}
	defer cursor.Close(ctx)

	var roles []model.Role
	if err = cursor.All(ctx, &roles); err != nil {
		logger.Error("解析公司角色列表失败", err)
		return nil, fmt.Errorf("解析公司角色列表失败: %w", err)
	}

	return roles, nil
}

// GetRolesByIDs 根据ID列表获取角色
func (r *roleRepository) GetRolesByIDs(ctx context.Context, roleIDs []string) ([]model.Role, error) {
	startTime := time.Now()
	filter := bson.M{"role_id": bson.M{"$in": roleIDs}}

	cursor, err := r.collection.Find(ctx, filter)

	// 记录数据库操作日志
	logger.DBLog("FIND_LIST", "roles", filter, time.Since(startTime))

	if err != nil {
		logger.Error("查询角色列表失败", err)
		return nil, fmt.Errorf("查询角色列表失败: %w", err)
	}
	defer cursor.Close(ctx)

	var roles []model.Role
	if err = cursor.All(ctx, &roles); err != nil {
		logger.Error("解析角色列表失败", err)
		return nil, fmt.Errorf("解析角色列表失败: %w", err)
	}

	return roles, nil
}

// BatchUpdateRoleStatus 批量更新角色状态
func (r *roleRepository) BatchUpdateRoleStatus(ctx context.Context, roleIDs []string, status string) error {
	startTime := time.Now()
	filter := bson.M{"role_id": bson.M{"$in": roleIDs}}
	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}

	result, err := r.collection.UpdateMany(ctx, filter, update)

	// 记录数据库操作日志
	logger.DBLog("BATCH_UPDATE", "roles", filter, time.Since(startTime))

	if err != nil {
		logger.Error("批量更新角色状态失败", err)
		return fmt.Errorf("批量更新角色状态失败: %w", err)
	}

	logger.Infof("批量更新角色状态成功: 更新了 %d 个角色", result.ModifiedCount)
	return nil
}

// GetRoleStats 获取角色统计信息
func (r *roleRepository) GetRoleStats(ctx context.Context, companyID string) (*model.RoleStatsResponse, error) {
	startTime := time.Now()

	var filter bson.M
	if companyID != "" {
		// 公司管理员只能看本公司数据
		filter = bson.M{"company_id": companyID}
	} else {
		// 平台管理员可以看所有数据
		filter = bson.M{}
	}

	// 总角色数
	totalRoles, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("统计总角色数失败: %w", err)
	}

	// 启用角色数
	enabledFilter := bson.M{"status": "enable"}
	if companyID != "" {
		enabledFilter["company_id"] = companyID
	}
	enabledRoles, err := r.collection.CountDocuments(ctx, enabledFilter)
	if err != nil {
		return nil, fmt.Errorf("统计启用角色数失败: %w", err)
	}

	// 禁用角色数
	disabledFilter := bson.M{"status": "disable"}
	if companyID != "" {
		disabledFilter["company_id"] = companyID
	}
	disabledRoles, err := r.collection.CountDocuments(ctx, disabledFilter)
	if err != nil {
		return nil, fmt.Errorf("统计禁用角色数失败: %w", err)
	}

	var platformRoles, companyRoles int64

	if companyID == "" {
		// 只有平台管理员才统计平台角色和公司角色的区分
		platformRoles, err = r.collection.CountDocuments(ctx, bson.M{"company_id": ""})
		if err != nil {
			return nil, fmt.Errorf("统计平台角色数失败: %w", err)
		}

		companyRoles, err = r.collection.CountDocuments(ctx, bson.M{"company_id": bson.M{"$ne": ""}})
		if err != nil {
			return nil, fmt.Errorf("统计公司角色数失败: %w", err)
		}
	}

	// 记录数据库操作日志
	logger.DBLog("STATS", "roles", filter, time.Since(startTime))

	return &model.RoleStatsResponse{
		TotalRoles:    totalRoles,
		EnabledRoles:  enabledRoles,
		DisabledRoles: disabledRoles,
		PlatformRoles: platformRoles,
		CompanyRoles:  companyRoles,
	}, nil
}

// CountRolesByCompanyID 统计公司角色数量
func (r *roleRepository) CountRolesByCompanyID(ctx context.Context, companyID string) (int64, error) {
	startTime := time.Now()
	filter := bson.M{"company_id": companyID}

	count, err := r.collection.CountDocuments(ctx, filter)

	// 记录数据库操作日志
	logger.DBLog("COUNT", "roles", filter, time.Since(startTime))

	if err != nil {
		logger.Error("统计公司角色数量失败", err)
		return 0, fmt.Errorf("统计公司角色数量失败: %w", err)
	}

	return count, nil
}

// CheckRoleKeyExists 检查角色标识符是否存在
func (r *roleRepository) CheckRoleKeyExists(ctx context.Context, roleKey string, excludeRoleID string) (bool, error) {
	startTime := time.Now()
	filter := bson.M{"role_key": roleKey}

	if excludeRoleID != "" {
		filter["role_id"] = bson.M{"$ne": excludeRoleID}
	}

	count, err := r.collection.CountDocuments(ctx, filter)

	// 记录数据库操作日志
	logger.DBLog("CHECK_EXISTS", "roles", filter, time.Since(startTime))

	if err != nil {
		logger.Error("检查角色标识符是否存在失败", err)
		return false, fmt.Errorf("检查角色标识符是否存在失败: %w", err)
	}

	return count > 0, nil
}

// CheckRoleNameExists 检查角色名称是否存在
func (r *roleRepository) CheckRoleNameExists(ctx context.Context, roleName string, companyID string, excludeRoleID string) (bool, error) {
	startTime := time.Now()
	filter := bson.M{
		"role_name":  roleName,
		"company_id": companyID,
	}

	if excludeRoleID != "" {
		filter["role_id"] = bson.M{"$ne": excludeRoleID}
	}

	count, err := r.collection.CountDocuments(ctx, filter)

	// 记录数据库操作日志
	logger.DBLog("CHECK_EXISTS", "roles", filter, time.Since(startTime))

	if err != nil {
		logger.Error("检查角色名称是否存在失败", err)
		return false, fmt.Errorf("检查角色名称是否存在失败: %w", err)
	}

	return count > 0, nil
}
