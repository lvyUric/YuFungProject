package repository

import (
	"YufungProject/internal/model"
	"YufungProject/pkg/database"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ActivityLogRepository struct {
	collection *mongo.Collection
}

func NewActivityLogRepository() *ActivityLogRepository {
	return &ActivityLogRepository{
		collection: database.MongoDB.Collection("activity_logs"),
	}
}

// Create 创建活动记录
func (r *ActivityLogRepository) Create(ctx context.Context, log *model.ActivityLog) error {
	log.OperationTime = time.Now()
	_, err := r.collection.InsertOne(ctx, log)
	return err
}

// GetList 获取活动记录列表
func (r *ActivityLogRepository) GetList(ctx context.Context, query *model.ActivityLogQuery) (*model.ActivityLogResponse, error) {
	filter := bson.M{}

	// 公司过滤
	if query.CompanyID != "" {
		filter["company_id"] = query.CompanyID
	}

	// 用户过滤
	if query.UserID != "" {
		filter["user_id"] = query.UserID
	}

	// 操作类型过滤
	if query.OperationType != "" {
		filter["operation_type"] = query.OperationType
	}

	// 模块名称过滤
	if query.ModuleName != "" {
		filter["module_name"] = query.ModuleName
	}

	// 时间范围过滤
	if query.StartTime != "" || query.EndTime != "" {
		timeFilter := bson.M{}
		if query.StartTime != "" {
			startTime, _ := time.Parse("2006-01-02 15:04:05", query.StartTime)
			timeFilter["$gte"] = startTime
		}
		if query.EndTime != "" {
			endTime, _ := time.Parse("2006-01-02 15:04:05", query.EndTime)
			timeFilter["$lte"] = endTime
		}
		if len(timeFilter) > 0 {
			filter["operation_time"] = timeFilter
		}
	}

	// 分页参数
	page := query.Page
	if page <= 0 {
		page = 1
	}
	pageSize := query.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	// 计算总数
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	// 查询数据
	opts := options.Find().
		SetSort(bson.D{{Key: "operation_time", Value: -1}}).
		SetSkip(int64((page - 1) * pageSize)).
		SetLimit(int64(pageSize))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []model.ActivityLog
	if err = cursor.All(ctx, &logs); err != nil {
		return nil, err
	}

	return &model.ActivityLogResponse{
		Total: total,
		List:  logs,
	}, nil
}

// GetRecentLogs 获取最近的活动记录（用于仪表盘）
func (r *ActivityLogRepository) GetRecentLogs(ctx context.Context, companyID string, limit int) ([]model.ActivityLog, error) {
	filter := bson.M{}
	if companyID != "" {
		filter["company_id"] = companyID
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "operation_time", Value: -1}}).
		SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []model.ActivityLog
	if err = cursor.All(ctx, &logs); err != nil {
		return nil, err
	}

	return logs, nil
}

// GetByID 根据ID获取活动记录
func (r *ActivityLogRepository) GetByID(ctx context.Context, id string) (*model.ActivityLog, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var log model.ActivityLog
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&log)
	if err != nil {
		return nil, err
	}

	return &log, nil
}

// DeleteByCompanyID 删除指定公司的所有活动记录
func (r *ActivityLogRepository) DeleteByCompanyID(ctx context.Context, companyID string) error {
	_, err := r.collection.DeleteMany(ctx, bson.M{"company_id": companyID})
	return err
}

// GetStatistics 获取活动记录统计
func (r *ActivityLogRepository) GetStatistics(ctx context.Context, companyID string, days int) (map[string]interface{}, error) {
	// 计算时间范围
	endTime := time.Now()
	startTime := endTime.AddDate(0, 0, -days)

	filter := bson.M{
		"operation_time": bson.M{
			"$gte": startTime,
			"$lte": endTime,
		},
	}

	if companyID != "" {
		filter["company_id"] = companyID
	}

	// 按操作类型统计
	pipeline := []bson.M{
		{"$match": filter},
		{"$group": bson.M{
			"_id":   "$operation_type",
			"count": bson.M{"$sum": 1},
		}},
		{"$sort": bson.M{"count": -1}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	// 按模块统计
	pipeline2 := []bson.M{
		{"$match": filter},
		{"$group": bson.M{
			"_id":   "$module_name",
			"count": bson.M{"$sum": 1},
		}},
		{"$sort": bson.M{"count": -1}},
	}

	cursor2, err := r.collection.Aggregate(ctx, pipeline2)
	if err != nil {
		return nil, err
	}
	defer cursor2.Close(ctx)

	var moduleResults []bson.M
	if err = cursor2.All(ctx, &moduleResults); err != nil {
		return nil, err
	}

	// 总记录数
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total":           total,
		"operation_types": results,
		"modules":         moduleResults,
	}, nil
}
