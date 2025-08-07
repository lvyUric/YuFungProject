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

// MenuRepository 菜单数据访问接口
type MenuRepository interface {
	// 基本CRUD操作
	CreateMenu(ctx context.Context, menu *model.Menu) error
	GetMenuByID(ctx context.Context, menuID string) (*model.Menu, error)
	UpdateMenu(ctx context.Context, menuID string, updates bson.M) error
	DeleteMenu(ctx context.Context, menuID string) error

	// 查询操作
	GetMenuList(ctx context.Context, filter bson.M) ([]model.Menu, error)
	GetMenuTree(ctx context.Context, filter bson.M) ([]model.Menu, error)
	GetChildMenus(ctx context.Context, parentID string) ([]model.Menu, error)
	GetMenusByIDs(ctx context.Context, menuIDs []string) ([]model.Menu, error)
	GetMenuIDsByRoleIDs(ctx context.Context, roleIDs []string) ([]string, error)

	// 批量操作
	BatchUpdateMenuStatus(ctx context.Context, menuIDs []string, status string) error

	// 统计操作
	GetMenuStats(ctx context.Context) (*model.MenuStatsResponse, error)

	// 验证操作
	CheckMenuNameExists(ctx context.Context, menuName string, parentID string, excludeMenuID string) (bool, error)
	CheckPermissionCodeExists(ctx context.Context, permissionCode string, excludeMenuID string) (bool, error)
	HasChildMenus(ctx context.Context, menuID string) (bool, error)
}

// menuRepository 菜单数据访问实现
type menuRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

// NewMenuRepository 创建菜单数据访问实例
func NewMenuRepository(db *mongo.Database) MenuRepository {
	repo := &menuRepository{
		db:         db,
		collection: db.Collection("menus"),
	}

	// 创建索引
	repo.createIndexes()
	return repo
}

// createIndexes 创建索引
func (r *menuRepository) createIndexes() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "menu_id", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("idx_menu_id_unique"),
		},
		{
			Keys:    bson.D{{Key: "parent_id", Value: 1}},
			Options: options.Index().SetName("idx_parent_id"),
		},
		{
			Keys:    bson.D{{Key: "menu_type", Value: 1}},
			Options: options.Index().SetName("idx_menu_type"),
		},
		{
			Keys:    bson.D{{Key: "status", Value: 1}},
			Options: options.Index().SetName("idx_status"),
		},
		{
			Keys:    bson.D{{Key: "sort_order", Value: 1}},
			Options: options.Index().SetName("idx_sort_order"),
		},
		{
			Keys:    bson.D{{Key: "permission_code", Value: 1}},
			Options: options.Index().SetSparse(true).SetName("idx_permission_code"),
		},
	}

	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		logger.Error("创建菜单表索引失败", err)
	} else {
		logger.Debug("菜单表索引创建成功")
	}
}

// CreateMenu 创建菜单
func (r *menuRepository) CreateMenu(ctx context.Context, menu *model.Menu) error {
	startTime := time.Now()

	// 设置创建时间
	menu.CreatedAt = time.Now()
	menu.UpdatedAt = time.Now()

	// 生成菜单ID
	if menu.MenuID == "" {
		menu.MenuID = utils.GenerateID("MENU")
	}

	_, err := r.collection.InsertOne(ctx, menu)

	// 记录数据库操作日志
	logger.DBLog("CREATE", "menus", bson.M{"menu_id": menu.MenuID}, time.Since(startTime))

	if err != nil {
		logger.Error("创建菜单失败", err)
		return fmt.Errorf("创建菜单失败: %w", err)
	}

	return nil
}

// GetMenuByID 根据ID获取菜单
func (r *menuRepository) GetMenuByID(ctx context.Context, menuID string) (*model.Menu, error) {
	startTime := time.Now()
	filter := bson.M{"menu_id": menuID}

	var menu model.Menu
	err := r.collection.FindOne(ctx, filter).Decode(&menu)

	// 记录数据库操作日志
	logger.DBLog("FIND", "menus", filter, time.Since(startTime))

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("菜单不存在")
		}
		logger.Error("查询菜单失败", err)
		return nil, fmt.Errorf("查询菜单失败: %w", err)
	}

	return &menu, nil
}

// UpdateMenu 更新菜单
func (r *menuRepository) UpdateMenu(ctx context.Context, menuID string, updates bson.M) error {
	startTime := time.Now()
	filter := bson.M{"menu_id": menuID}

	// 添加更新时间
	updates["updated_at"] = time.Now()

	update := bson.M{"$set": updates}

	result, err := r.collection.UpdateOne(ctx, filter, update)

	// 记录数据库操作日志
	logger.DBLog("UPDATE", "menus", filter, time.Since(startTime))

	if err != nil {
		logger.Error("更新菜单失败", err)
		return fmt.Errorf("更新菜单失败: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("菜单不存在")
	}

	return nil
}

// DeleteMenu 删除菜单
func (r *menuRepository) DeleteMenu(ctx context.Context, menuID string) error {
	startTime := time.Now()
	filter := bson.M{"menu_id": menuID}

	result, err := r.collection.DeleteOne(ctx, filter)

	// 记录数据库操作日志
	logger.DBLog("DELETE", "menus", filter, time.Since(startTime))

	if err != nil {
		logger.Error("删除菜单失败", err)
		return fmt.Errorf("删除菜单失败: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("菜单不存在")
	}

	return nil
}

// GetMenuList 获取菜单列表
func (r *menuRepository) GetMenuList(ctx context.Context, filter bson.M) ([]model.Menu, error) {
	startTime := time.Now()

	// 设置查询选项：按层级和排序号排序
	findOptions := options.Find().SetSort(bson.D{
		{Key: "parent_id", Value: 1},
		{Key: "sort_order", Value: 1},
		{Key: "created_at", Value: 1},
	})

	cursor, err := r.collection.Find(ctx, filter, findOptions)

	// 记录数据库操作日志
	logger.DBLog("FIND_LIST", "menus", filter, time.Since(startTime))

	if err != nil {
		logger.Error("查询菜单列表失败", err)
		return nil, fmt.Errorf("查询菜单列表失败: %w", err)
	}
	defer cursor.Close(ctx)

	var menus []model.Menu
	if err = cursor.All(ctx, &menus); err != nil {
		logger.Error("解析菜单列表失败", err)
		return nil, fmt.Errorf("解析菜单列表失败: %w", err)
	}

	return menus, nil
}

// GetMenuTree 获取菜单树
func (r *menuRepository) GetMenuTree(ctx context.Context, filter bson.M) ([]model.Menu, error) {
	startTime := time.Now()

	// 获取所有菜单
	menus, err := r.GetMenuList(ctx, filter)
	if err != nil {
		return nil, err
	}

	// 记录数据库操作日志
	logger.DBLog("FIND_TREE", "menus", filter, time.Since(startTime))

	return menus, nil
}

// GetChildMenus 获取子菜单
func (r *menuRepository) GetChildMenus(ctx context.Context, parentID string) ([]model.Menu, error) {
	startTime := time.Now()
	filter := bson.M{"parent_id": parentID}

	findOptions := options.Find().SetSort(bson.D{{Key: "sort_order", Value: 1}})
	cursor, err := r.collection.Find(ctx, filter, findOptions)

	// 记录数据库操作日志
	logger.DBLog("FIND_CHILDREN", "menus", filter, time.Since(startTime))

	if err != nil {
		logger.Error("查询子菜单失败", err)
		return nil, fmt.Errorf("查询子菜单失败: %w", err)
	}
	defer cursor.Close(ctx)

	var menus []model.Menu
	if err = cursor.All(ctx, &menus); err != nil {
		logger.Error("解析子菜单失败", err)
		return nil, fmt.Errorf("解析子菜单失败: %w", err)
	}

	return menus, nil
}

// GetMenusByIDs 根据ID列表获取菜单
func (r *menuRepository) GetMenusByIDs(ctx context.Context, menuIDs []string) ([]model.Menu, error) {
	startTime := time.Now()
	filter := bson.M{"menu_id": bson.M{"$in": menuIDs}}

	cursor, err := r.collection.Find(ctx, filter)

	// 记录数据库操作日志
	logger.DBLog("FIND_LIST", "menus", filter, time.Since(startTime))

	if err != nil {
		logger.Error("查询菜单列表失败", err)
		return nil, fmt.Errorf("查询菜单列表失败: %w", err)
	}
	defer cursor.Close(ctx)

	var menus []model.Menu
	if err = cursor.All(ctx, &menus); err != nil {
		logger.Error("解析菜单列表失败", err)
		return nil, fmt.Errorf("解析菜单列表失败: %w", err)
	}

	return menus, nil
}

// BatchUpdateMenuStatus 批量更新菜单状态
func (r *menuRepository) BatchUpdateMenuStatus(ctx context.Context, menuIDs []string, status string) error {
	startTime := time.Now()
	filter := bson.M{"menu_id": bson.M{"$in": menuIDs}}
	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}

	result, err := r.collection.UpdateMany(ctx, filter, update)

	// 记录数据库操作日志
	logger.DBLog("BATCH_UPDATE", "menus", filter, time.Since(startTime))

	if err != nil {
		logger.Error("批量更新菜单状态失败", err)
		return fmt.Errorf("批量更新菜单状态失败: %w", err)
	}

	logger.Infof("批量更新菜单状态成功: 更新了 %d 个菜单", result.ModifiedCount)
	return nil
}

// GetMenuStats 获取菜单统计信息
func (r *menuRepository) GetMenuStats(ctx context.Context) (*model.MenuStatsResponse, error) {
	startTime := time.Now()

	// 总菜单数
	totalMenus, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("统计总菜单数失败: %w", err)
	}

	// 启用菜单数
	enabledMenus, err := r.collection.CountDocuments(ctx, bson.M{"status": "enable"})
	if err != nil {
		return nil, fmt.Errorf("统计启用菜单数失败: %w", err)
	}

	// 禁用菜单数
	disabledMenus, err := r.collection.CountDocuments(ctx, bson.M{"status": "disable"})
	if err != nil {
		return nil, fmt.Errorf("统计禁用菜单数失败: %w", err)
	}

	// 目录类型菜单数
	directoryMenus, err := r.collection.CountDocuments(ctx, bson.M{"menu_type": "directory"})
	if err != nil {
		return nil, fmt.Errorf("统计目录菜单数失败: %w", err)
	}

	// 页面类型菜单数
	pageMenus, err := r.collection.CountDocuments(ctx, bson.M{"menu_type": "menu"})
	if err != nil {
		return nil, fmt.Errorf("统计页面菜单数失败: %w", err)
	}

	// 按钮类型菜单数
	buttonMenus, err := r.collection.CountDocuments(ctx, bson.M{"menu_type": "button"})
	if err != nil {
		return nil, fmt.Errorf("统计按钮菜单数失败: %w", err)
	}

	// 记录数据库操作日志
	logger.DBLog("STATS", "menus", bson.M{}, time.Since(startTime))

	return &model.MenuStatsResponse{
		TotalMenus:     totalMenus,
		EnabledMenus:   enabledMenus,
		DisabledMenus:  disabledMenus,
		DirectoryMenus: directoryMenus,
		PageMenus:      pageMenus,
		ButtonMenus:    buttonMenus,
	}, nil
}

// CheckMenuNameExists 检查菜单名称是否存在
func (r *menuRepository) CheckMenuNameExists(ctx context.Context, menuName string, parentID string, excludeMenuID string) (bool, error) {
	startTime := time.Now()
	filter := bson.M{
		"menu_name": menuName,
		"parent_id": parentID,
	}

	if excludeMenuID != "" {
		filter["menu_id"] = bson.M{"$ne": excludeMenuID}
	}

	count, err := r.collection.CountDocuments(ctx, filter)

	// 记录数据库操作日志
	logger.DBLog("CHECK_EXISTS", "menus", filter, time.Since(startTime))

	if err != nil {
		logger.Error("检查菜单名称是否存在失败", err)
		return false, fmt.Errorf("检查菜单名称是否存在失败: %w", err)
	}

	return count > 0, nil
}

// CheckPermissionCodeExists 检查权限标识符是否存在
func (r *menuRepository) CheckPermissionCodeExists(ctx context.Context, permissionCode string, excludeMenuID string) (bool, error) {
	startTime := time.Now()
	filter := bson.M{"permission_code": permissionCode}

	if excludeMenuID != "" {
		filter["menu_id"] = bson.M{"$ne": excludeMenuID}
	}

	count, err := r.collection.CountDocuments(ctx, filter)

	// 记录数据库操作日志
	logger.DBLog("CHECK_EXISTS", "menus", filter, time.Since(startTime))

	if err != nil {
		logger.Error("检查权限标识符是否存在失败", err)
		return false, fmt.Errorf("检查权限标识符是否存在失败: %w", err)
	}

	return count > 0, nil
}

// HasChildMenus 检查是否有子菜单
func (r *menuRepository) HasChildMenus(ctx context.Context, menuID string) (bool, error) {
	startTime := time.Now()
	filter := bson.M{"parent_id": menuID}

	count, err := r.collection.CountDocuments(ctx, filter)

	// 记录数据库操作日志
	logger.DBLog("COUNT_CHILDREN", "menus", filter, time.Since(startTime))

	if err != nil {
		logger.Error("检查子菜单失败", err)
		return false, fmt.Errorf("检查子菜单失败: %w", err)
	}

	return count > 0, nil
}

// GetMenuIDsByRoleIDs 根据角色ID列表获取菜单ID列表
func (r *menuRepository) GetMenuIDsByRoleIDs(ctx context.Context, roleIDs []string) ([]string, error) {
	startTime := time.Now()

	// 参数验证
	if len(roleIDs) == 0 {
		logger.Debug("角色ID列表为空，返回空菜单ID列表")
		return []string{}, nil
	}

	logger.Debugf("查询角色菜单权限，roleIDs: %+v", roleIDs)

	// 查询角色表获取这些角色的菜单权限
	rolesCollection := r.db.Collection("roles")
	filter := bson.M{"role_id": bson.M{"$in": roleIDs}, "status": "enable"}

	logger.Debugf("MongoDB查询filter: %+v", filter)

	cursor, err := rolesCollection.Find(ctx, filter)

	// 记录数据库操作日志
	logger.DBLog("FIND_ROLE_MENUS", "roles", filter, time.Since(startTime))

	if err != nil {
		logger.Error("查询角色菜单权限失败", err)
		return nil, fmt.Errorf("查询角色菜单权限失败: %w", err)
	}
	defer cursor.Close(ctx)

	// 收集所有角色的菜单ID
	menuIDSet := make(map[string]bool)

	for cursor.Next(ctx) {
		var role model.Role
		if err := cursor.Decode(&role); err != nil {
			logger.Error("解析角色数据失败", err)
			continue
		}

		logger.Debugf("角色 %s 的菜单权限: %+v", role.RoleID, role.MenuIDs)

		// 添加该角色的所有菜单ID到集合中
		for _, menuID := range role.MenuIDs {
			if menuID != "" {
				menuIDSet[menuID] = true
			}
		}
	}

	// 转换为切片
	var menuIDs []string
	for menuID := range menuIDSet {
		menuIDs = append(menuIDs, menuID)
	}

	logger.Debugf("获取到的菜单ID列表: %+v", menuIDs)

	return menuIDs, nil
}
