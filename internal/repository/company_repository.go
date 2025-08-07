package repository

import (
	"context"
	"time"

	"YufungProject/internal/model"
	"YufungProject/pkg/logger"
	"YufungProject/pkg/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CompanyRepository 公司仓库接口
type CompanyRepository interface {
	// 创建公司
	CreateCompany(ctx context.Context, company *model.Company) error
	// 根据ID获取公司
	GetCompanyByID(ctx context.Context, companyID string) (*model.Company, error)
	// 获取公司列表
	GetCompanyList(ctx context.Context, page, pageSize int, status string) ([]*model.Company, int64, error)
	// 更新公司
	UpdateCompany(ctx context.Context, companyID string, updates bson.M) error
	// 删除公司
	DeleteCompany(ctx context.Context, companyID string) error
	// 检查公司名称是否存在
	ExistsCompanyName(ctx context.Context, companyName, excludeID string) (bool, error)
	// 获取公司用户统计
	GetCompanyUserStats(ctx context.Context, companyID string) (int64, error)
}

// companyRepository 公司仓库实现
type companyRepository struct {
	collection *mongo.Collection
	userRepo   UserRepository
}

// NewCompanyRepository 创建公司仓库实例
func NewCompanyRepository(db *mongo.Database, userRepo UserRepository) CompanyRepository {
	repo := &companyRepository{
		collection: db.Collection("companies"),
		userRepo:   userRepo,
	}

	// 创建索引
	repo.createIndexes()

	return repo
}

// createIndexes 创建索引
func (r *companyRepository) createIndexes() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "company_id", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("idx_company_id"),
		},
		{
			Keys:    bson.D{{Key: "company_name", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("idx_company_name"),
		},
		{
			Keys:    bson.D{{Key: "status", Value: 1}},
			Options: options.Index().SetName("idx_status"),
		},
		{
			Keys:    bson.D{{Key: "created_at", Value: 1}},
			Options: options.Index().SetName("idx_created_at"),
		},
	}

	if _, err := r.collection.Indexes().CreateMany(ctx, indexes); err != nil {
		logger.Errorf("创建公司表索引失败: %v", err)
	}
}

// CreateCompany 创建公司
func (r *companyRepository) CreateCompany(ctx context.Context, company *model.Company) error {
	// 生成公司ID
	company.CompanyID = utils.GenerateCompanyID()
	company.ID = primitive.NewObjectID()
	company.CurrentUserCount = 0
	company.CreatedAt = time.Now()
	company.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, company)
	if err != nil {
		logger.Errorf("创建公司失败: %v", err)
		return err
	}

	logger.Infof("公司创建成功: %s - %s", company.CompanyID, company.CompanyName)
	return nil
}

// GetCompanyByID 根据ID获取公司
func (r *companyRepository) GetCompanyByID(ctx context.Context, companyID string) (*model.Company, error) {
	var company model.Company
	err := r.collection.FindOne(ctx, bson.M{"company_id": companyID}).Decode(&company)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		logger.Errorf("查询公司失败: %v", err)
		return nil, err
	}

	return &company, nil
}

// GetCompanyList 获取公司列表
func (r *companyRepository) GetCompanyList(ctx context.Context, page, pageSize int, status string) ([]*model.Company, int64, error) {
	// 构建查询条件
	filter := bson.M{}
	if status != "" {
		filter["status"] = status
	}

	// 计算总数
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		logger.Errorf("统计公司数量失败: %v", err)
		return nil, 0, err
	}

	// 计算分页
	skip := (page - 1) * pageSize

	// 查询选项
	opts := options.Find()
	opts.SetSkip(int64(skip))
	opts.SetLimit(int64(pageSize))
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}})

	// 执行查询
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		logger.Errorf("查询公司列表失败: %v", err)
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var companies []*model.Company
	if err = cursor.All(ctx, &companies); err != nil {
		logger.Errorf("解析公司列表失败: %v", err)
		return nil, 0, err
	}

	// 获取每个公司的用户统计
	for _, company := range companies {
		userCount, _ := r.GetCompanyUserStats(ctx, company.CompanyID)
		company.CurrentUserCount = int(userCount)
	}

	return companies, total, nil
}

// UpdateCompany 更新公司
func (r *companyRepository) UpdateCompany(ctx context.Context, companyID string, updates bson.M) error {
	updates["updated_at"] = time.Now()

	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"company_id": companyID},
		bson.M{"$set": updates},
	)
	if err != nil {
		logger.Errorf("更新公司失败: %v", err)
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	logger.Infof("公司更新成功: %s", companyID)
	return nil
}

// DeleteCompany 删除公司
func (r *companyRepository) DeleteCompany(ctx context.Context, companyID string) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"company_id": companyID})
	if err != nil {
		logger.Errorf("删除公司失败: %v", err)
		return err
	}

	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	logger.Infof("公司删除成功: %s", companyID)
	return nil
}

// ExistsCompanyName 检查公司名称是否存在
func (r *companyRepository) ExistsCompanyName(ctx context.Context, companyName, excludeID string) (bool, error) {
	filter := bson.M{"company_name": companyName}
	if excludeID != "" {
		filter["company_id"] = bson.M{"$ne": excludeID}
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		logger.Errorf("检查公司名称是否存在失败: %v", err)
		return false, err
	}

	return count > 0, nil
}

// GetCompanyUserStats 获取公司用户统计
func (r *companyRepository) GetCompanyUserStats(ctx context.Context, companyID string) (int64, error) {
	if r.userRepo == nil {
		return 0, nil
	}

	// 这里应该调用用户仓库的方法来统计用户数量
	// 由于我们还没有实现用户仓库的这个方法，先返回0
	return 0, nil
}
