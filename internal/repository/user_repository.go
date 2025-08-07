package repository

import (
	"context"
	"time"

	"YufungProject/internal/model"
	"YufungProject/pkg/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UserRepository 用户数据访问接口
type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*model.User, error)
	GetByUserID(ctx context.Context, userID string) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Update(ctx context.Context, userID string, update bson.M) error
	UpdatePassword(ctx context.Context, userID string, passwordHash string) error
	UpdateLoginAttempts(ctx context.Context, userID string, attempts int, lockedUntil *time.Time) error
	UpdateLastLoginTime(ctx context.Context, userID string, loginTime time.Time) error
	List(ctx context.Context, filter bson.M, page, pageSize int) ([]*model.User, int64, error)
	Delete(ctx context.Context, userID string) error
	ExistsByUsername(ctx context.Context, username string) (bool, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}

type userRepository struct {
	collection *mongo.Collection
}

// NewUserRepository 创建用户仓库实例
func NewUserRepository(db *mongo.Database) UserRepository {
	collection := db.Collection("users")

	// 创建索引
	repo := &userRepository{collection: collection}
	repo.createIndexes()

	return repo
}

// createIndexes 创建索引
func (r *userRepository) createIndexes() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "user_id", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("idx_user_id"),
		},
		{
			Keys:    bson.D{{Key: "username", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("idx_username"),
		},
		{
			Keys:    bson.D{{Key: "company_id", Value: 1}},
			Options: options.Index().SetName("idx_company_id"),
		},
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetSparse(true).SetName("idx_email"),
		},
		{
			Keys:    bson.D{{Key: "status", Value: 1}},
			Options: options.Index().SetName("idx_status"),
		},
		{
			Keys:    bson.D{{Key: "created_at", Value: 1}},
			Options: options.Index().SetName("idx_created_at"),
		},
		{
			Keys: bson.D{
				{Key: "company_id", Value: 1},
				{Key: "status", Value: 1},
			},
			Options: options.Index().SetName("idx_company_status"),
		},
	}

	start := time.Now()
	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	duration := time.Since(start)

	if err != nil {
		logger.Errorf("创建用户表索引失败: %v", err)
	} else {
		logger.DBLog("CREATE_INDEX", "users", "all_indexes", duration)
	}
}

// Create 创建用户
func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	start := time.Now()
	_, err := r.collection.InsertOne(ctx, user)
	duration := time.Since(start)

	if err != nil {
		logger.DBLog("INSERT_ERROR", "users", bson.M{"user_id": user.UserID}, duration)
		logger.Errorf("创建用户失败: %v, UserID: %s", err, user.UserID)
		return err
	}

	logger.DBLog("INSERT", "users", bson.M{"user_id": user.UserID}, duration)
	logger.Infof("用户创建成功: UserID=%s, Username=%s", user.UserID, user.Username)
	return nil
}

// GetByID 根据ObjectID获取用户
func (r *userRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*model.User, error) {
	start := time.Now()
	filter := bson.M{"_id": id}

	var user model.User
	err := r.collection.FindOne(ctx, filter).Decode(&user)
	duration := time.Since(start)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			logger.DBLog("FIND_NOT_FOUND", "users", filter, duration)
		} else {
			logger.DBLog("FIND_ERROR", "users", filter, duration)
			logger.Errorf("根据ID查询用户失败: %v, ID: %s", err, id.Hex())
		}
		return nil, err
	}

	logger.DBLog("FIND", "users", filter, duration)
	return &user, nil
}

// GetByUserID 根据用户ID获取用户
func (r *userRepository) GetByUserID(ctx context.Context, userID string) (*model.User, error) {
	start := time.Now()
	filter := bson.M{"user_id": userID}

	var user model.User
	err := r.collection.FindOne(ctx, filter).Decode(&user)
	duration := time.Since(start)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			logger.DBLog("FIND_NOT_FOUND", "users", filter, duration)
		} else {
			logger.DBLog("FIND_ERROR", "users", filter, duration)
			logger.Errorf("根据UserID查询用户失败: %v, UserID: %s", err, userID)
		}
		return nil, err
	}

	logger.DBLog("FIND", "users", filter, duration)
	return &user, nil
}

// GetByUsername 根据用户名获取用户
func (r *userRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	start := time.Now()
	filter := bson.M{"username": username}

	var user model.User
	err := r.collection.FindOne(ctx, filter).Decode(&user)
	duration := time.Since(start)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			logger.DBLog("FIND_NOT_FOUND", "users", filter, duration)
		} else {
			logger.DBLog("FIND_ERROR", "users", filter, duration)
			logger.Errorf("根据用户名查询用户失败: %v, Username: %s", err, username)
		}
		return nil, err
	}

	logger.DBLog("FIND", "users", filter, duration)
	return &user, nil
}

// GetByEmail 根据邮箱获取用户
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	start := time.Now()
	filter := bson.M{"email": email}

	var user model.User
	err := r.collection.FindOne(ctx, filter).Decode(&user)
	duration := time.Since(start)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			logger.DBLog("FIND_NOT_FOUND", "users", filter, duration)
		} else {
			logger.DBLog("FIND_ERROR", "users", filter, duration)
			logger.Errorf("根据邮箱查询用户失败: %v, Email: %s", err, email)
		}
		return nil, err
	}

	logger.DBLog("FIND", "users", filter, duration)
	return &user, nil
}

// Update 更新用户信息
func (r *userRepository) Update(ctx context.Context, userID string, update bson.M) error {
	start := time.Now()
	filter := bson.M{"user_id": userID}
	update["updated_at"] = time.Now()

	result, err := r.collection.UpdateOne(ctx, filter, bson.M{"$set": update})
	duration := time.Since(start)

	if err != nil {
		logger.DBLog("UPDATE_ERROR", "users", filter, duration)
		logger.Errorf("更新用户失败: %v, UserID: %s", err, userID)
		return err
	}

	logger.DBLog("UPDATE", "users", filter, duration)
	logger.Infof("用户更新成功: UserID=%s, ModifiedCount=%d", userID, result.ModifiedCount)
	return nil
}

// UpdatePassword 更新用户密码
func (r *userRepository) UpdatePassword(ctx context.Context, userID string, passwordHash string) error {
	start := time.Now()
	filter := bson.M{"user_id": userID}
	update := bson.M{
		"password_hash": passwordHash,
		"updated_at":    time.Now(),
	}

	result, err := r.collection.UpdateOne(ctx, filter, bson.M{"$set": update})
	duration := time.Since(start)

	if err != nil {
		logger.DBLog("UPDATE_PASSWORD_ERROR", "users", filter, duration)
		logger.Errorf("更新用户密码失败: %v, UserID: %s", err, userID)
		return err
	}

	logger.DBLog("UPDATE_PASSWORD", "users", filter, duration)
	logger.Infof("用户密码更新成功: UserID=%s, ModifiedCount=%d", userID, result.ModifiedCount)
	return nil
}

// UpdateLoginAttempts 更新登录尝试次数
func (r *userRepository) UpdateLoginAttempts(ctx context.Context, userID string, attempts int, lockedUntil *time.Time) error {
	start := time.Now()
	filter := bson.M{"user_id": userID}
	update := bson.M{
		"login_attempts": attempts,
		"locked_until":   lockedUntil,
		"updated_at":     time.Now(),
	}

	result, err := r.collection.UpdateOne(ctx, filter, bson.M{"$set": update})
	duration := time.Since(start)

	if err != nil {
		logger.DBLog("UPDATE_LOGIN_ATTEMPTS_ERROR", "users", filter, duration)
		logger.Errorf("更新登录尝试次数失败: %v, UserID: %s", err, userID)
		return err
	}

	logger.DBLog("UPDATE_LOGIN_ATTEMPTS", "users", filter, duration)
	logger.Infof("登录尝试次数更新成功: UserID=%s, Attempts=%d, ModifiedCount=%d", userID, attempts, result.ModifiedCount)
	return nil
}

// UpdateLastLoginTime 更新最后登录时间
func (r *userRepository) UpdateLastLoginTime(ctx context.Context, userID string, loginTime time.Time) error {
	start := time.Now()
	filter := bson.M{"user_id": userID}
	update := bson.M{
		"last_login_time": loginTime,
		"updated_at":      time.Now(),
	}

	result, err := r.collection.UpdateOne(ctx, filter, bson.M{"$set": update})
	duration := time.Since(start)

	if err != nil {
		logger.DBLog("UPDATE_LAST_LOGIN_ERROR", "users", filter, duration)
		logger.Errorf("更新最后登录时间失败: %v, UserID: %s", err, userID)
		return err
	}

	logger.DBLog("UPDATE_LAST_LOGIN", "users", filter, duration)
	logger.Infof("最后登录时间更新成功: UserID=%s, ModifiedCount=%d", userID, result.ModifiedCount)
	return nil
}

// List 分页查询用户列表
func (r *userRepository) List(ctx context.Context, filter bson.M, page, pageSize int) ([]*model.User, int64, error) {
	start := time.Now()

	// 计算总数
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		duration := time.Since(start)
		logger.DBLog("COUNT_ERROR", "users", filter, duration)
		logger.Errorf("查询用户总数失败: %v", err)
		return nil, 0, err
	}

	// 查询数据
	skip := (page - 1) * pageSize
	findOptions := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize))

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		duration := time.Since(start)
		logger.DBLog("FIND_LIST_ERROR", "users", filter, duration)
		logger.Errorf("查询用户列表失败: %v", err)
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var users []*model.User
	if err = cursor.All(ctx, &users); err != nil {
		duration := time.Since(start)
		logger.DBLog("DECODE_LIST_ERROR", "users", filter, duration)
		logger.Errorf("解析用户列表失败: %v", err)
		return nil, 0, err
	}

	duration := time.Since(start)
	logger.DBLog("FIND_LIST", "users", filter, duration)
	logger.Infof("用户列表查询成功: Total=%d, Page=%d, PageSize=%d, ResultCount=%d", total, page, pageSize, len(users))

	return users, total, nil
}

// Delete 删除用户（软删除）
func (r *userRepository) Delete(ctx context.Context, userID string) error {
	start := time.Now()
	filter := bson.M{"user_id": userID}
	update := bson.M{
		"status":     "inactive",
		"updated_at": time.Now(),
	}

	result, err := r.collection.UpdateOne(ctx, filter, bson.M{"$set": update})
	duration := time.Since(start)

	if err != nil {
		logger.DBLog("DELETE_ERROR", "users", filter, duration)
		logger.Errorf("删除用户失败: %v, UserID: %s", err, userID)
		return err
	}

	logger.DBLog("DELETE", "users", filter, duration)
	logger.Infof("用户删除成功: UserID=%s, ModifiedCount=%d", userID, result.ModifiedCount)
	return nil
}

// ExistsByUsername 检查用户名是否存在
func (r *userRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	start := time.Now()
	filter := bson.M{"username": username}

	count, err := r.collection.CountDocuments(ctx, filter)
	duration := time.Since(start)

	if err != nil {
		logger.DBLog("EXISTS_CHECK_ERROR", "users", filter, duration)
		logger.Errorf("检查用户名是否存在失败: %v, Username: %s", err, username)
		return false, err
	}

	logger.DBLog("EXISTS_CHECK", "users", filter, duration)
	exists := count > 0
	logger.Debugf("用户名存在检查: Username=%s, Exists=%v", username, exists)

	return exists, nil
}

// ExistsByEmail 检查邮箱是否存在
func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	start := time.Now()
	filter := bson.M{"email": email}

	count, err := r.collection.CountDocuments(ctx, filter)
	duration := time.Since(start)

	if err != nil {
		logger.DBLog("EXISTS_CHECK_ERROR", "users", filter, duration)
		logger.Errorf("检查邮箱是否存在失败: %v, Email: %s", err, email)
		return false, err
	}

	logger.DBLog("EXISTS_CHECK", "users", filter, duration)
	exists := count > 0
	logger.Debugf("邮箱存在检查: Email=%s, Exists=%v", email, exists)

	return exists, nil
}
