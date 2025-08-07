package repository

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"YufungProject/internal/model"
	"YufungProject/pkg/utils"
)

// SystemConfigRepository 系统配置数据访问层接口
type SystemConfigRepository interface {
	Create(ctx context.Context, config *model.SystemConfig) error
	GetByID(ctx context.Context, configID string) (*model.SystemConfig, error)
	Update(ctx context.Context, configID string, config *model.SystemConfig) error
	Delete(ctx context.Context, configID string) error
	List(ctx context.Context, req *model.SystemConfigQueryRequest, companyID string) (*model.SystemConfigListResponse, error)
	GetByType(ctx context.Context, configType, companyID string) ([]model.SystemConfig, error)
	CheckKeyExists(ctx context.Context, configType, configKey, companyID, excludeID string) (bool, error)
}

type systemConfigRepository struct {
	db *mongo.Database
}

// NewSystemConfigRepository 创建系统配置数据访问层实例
func NewSystemConfigRepository(db *mongo.Database) SystemConfigRepository {
	return &systemConfigRepository{
		db: db,
	}
}

// Create 创建系统配置
func (r *systemConfigRepository) Create(ctx context.Context, config *model.SystemConfig) error {
	collection := r.db.Collection("system_configs")

	// 生成配置ID
	config.ConfigID = utils.GenerateID("CONFIG")
	config.CreatedAt = time.Now()
	config.UpdatedAt = time.Now()

	_, err := collection.InsertOne(ctx, config)
	return err
}

// GetByID 根据ID获取系统配置
func (r *systemConfigRepository) GetByID(ctx context.Context, configID string) (*model.SystemConfig, error) {
	collection := r.db.Collection("system_configs")

	var config model.SystemConfig
	filter := bson.M{"config_id": configID}

	err := collection.FindOne(ctx, filter).Decode(&config)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("系统配置不存在")
		}
		return nil, err
	}

	return &config, nil
}

// Update 更新系统配置
func (r *systemConfigRepository) Update(ctx context.Context, configID string, config *model.SystemConfig) error {
	collection := r.db.Collection("system_configs")

	filter := bson.M{"config_id": configID}
	update := bson.M{
		"$set": bson.M{
			"config_value": config.ConfigValue,
			"display_name": config.DisplayName,
			"sort_order":   config.SortOrder,
			"status":       config.Status,
			"remark":       config.Remark,
			"updated_by":   config.UpdatedBy,
			"updated_at":   time.Now(),
		},
	}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("系统配置不存在")
	}

	return nil
}

// Delete 删除系统配置
func (r *systemConfigRepository) Delete(ctx context.Context, configID string) error {
	collection := r.db.Collection("system_configs")

	filter := bson.M{"config_id": configID}
	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("系统配置不存在")
	}

	return nil
}

// List 获取系统配置列表
func (r *systemConfigRepository) List(ctx context.Context, req *model.SystemConfigQueryRequest, companyID string) (*model.SystemConfigListResponse, error) {
	collection := r.db.Collection("system_configs")

	// 构建查询条件
	filter := bson.M{}

	// 只有当companyID不为空时才添加公司过滤条件
	if companyID != "" {
		filter["company_id"] = companyID
	}

	if req.ConfigType != "" {
		filter["config_type"] = req.ConfigType
	}

	if req.Status != "" {
		filter["status"] = req.Status
	}

	if req.Keyword != "" {
		filter["$or"] = []bson.M{
			{"config_key": primitive.Regex{Pattern: req.Keyword, Options: "i"}},
			{"config_value": primitive.Regex{Pattern: req.Keyword, Options: "i"}},
			{"display_name": primitive.Regex{Pattern: req.Keyword, Options: "i"}},
		}
	}

	// 设置分页参数
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	skip := int64((page - 1) * pageSize)
	limit := int64(pageSize)

	// 查询总数
	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	// 查询数据
	findOptions := options.Find().
		SetSkip(skip).
		SetLimit(limit).
		SetSort(bson.D{
			{Key: "sort_order", Value: 1},
			{Key: "created_at", Value: -1},
		})

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var configs []model.SystemConfig
	if err = cursor.All(ctx, &configs); err != nil {
		return nil, err
	}

	return &model.SystemConfigListResponse{
		List:  configs,
		Total: total,
	}, nil
}

// GetByType 根据配置类型获取配置列表
func (r *systemConfigRepository) GetByType(ctx context.Context, configType, companyID string) ([]model.SystemConfig, error) {
	collection := r.db.Collection("system_configs")

	filter := bson.M{
		"config_type": configType,
		"status":      "enable",
	}

	// 只有当companyID不为空时才添加公司过滤条件
	if companyID != "" {
		filter["company_id"] = companyID
	}

	findOptions := options.Find().SetSort(bson.D{
		{Key: "sort_order", Value: 1},
		{Key: "created_at", Value: -1},
	})

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var configs []model.SystemConfig
	if err = cursor.All(ctx, &configs); err != nil {
		return nil, err
	}

	return configs, nil
}

// CheckKeyExists 检查配置键是否存在
func (r *systemConfigRepository) CheckKeyExists(ctx context.Context, configType, configKey, companyID, excludeID string) (bool, error) {
	collection := r.db.Collection("system_configs")

	filter := bson.M{
		"config_type": configType,
		"config_key":  configKey,
		"company_id":  companyID,
	}

	if excludeID != "" {
		filter["config_id"] = bson.M{"$ne": excludeID}
	}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
