package repository

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"YufungProject/internal/model"
	"YufungProject/pkg/utils"
)

const ChangeRecordCollection = "change_records"

type ChangeRecordRepository struct {
	db *mongo.Database
}

func NewChangeRecordRepository(db *mongo.Database) *ChangeRecordRepository {
	return &ChangeRecordRepository{db: db}
}

// CreateChangeRecord 创建变更记录
func (r *ChangeRecordRepository) CreateChangeRecord(ctx context.Context, record *model.ChangeRecord) error {
	collection := r.db.Collection(ChangeRecordCollection)

	// 生成唯一标识
	record.ChangeID = utils.GenerateID("CHG")
	record.ChangeTime = time.Now()

	_, err := collection.InsertOne(ctx, record)
	return err
}

// GetChangeRecordsByTableAndRecord 获取指定表和记录的变更记录
func (r *ChangeRecordRepository) GetChangeRecordsByTableAndRecord(ctx context.Context, tableName, recordID string, days int, page, pageSize int) ([]*model.ChangeRecord, int64, error) {
	collection := r.db.Collection(ChangeRecordCollection)

	// 构建查询条件
	filter := bson.M{
		"table_name": tableName,
		"record_id":  recordID,
	}

	// 添加时间范围条件（默认查询最近指定天数的记录）
	if days > 0 {
		startTime := time.Now().AddDate(0, 0, -days)
		filter["change_time"] = bson.M{"$gte": startTime}
	}

	// 计算总数
	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	skip := (page - 1) * pageSize
	opts := options.Find().
		SetSort(bson.D{{Key: "change_time", Value: -1}}). // 按时间倒序
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize))

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var records []*model.ChangeRecord
	for cursor.Next(ctx) {
		var record model.ChangeRecord
		if err := cursor.Decode(&record); err != nil {
			return nil, 0, err
		}
		records = append(records, &record)
	}

	return records, total, nil
}

// GetChangeRecordsList 获取变更记录列表（支持多种筛选条件）
func (r *ChangeRecordRepository) GetChangeRecordsList(ctx context.Context, params *model.ChangeRecordListParams) ([]*model.ChangeRecord, int64, error) {
	collection := r.db.Collection(ChangeRecordCollection)

	// 构建查询条件
	filter := bson.M{}

	if params.TableName != "" {
		filter["table_name"] = params.TableName
	}
	if params.RecordID != "" {
		filter["record_id"] = params.RecordID
	}
	if params.UserID != "" {
		filter["user_id"] = params.UserID
	}
	if params.CompanyID != "" {
		filter["company_id"] = params.CompanyID
	}
	if params.ChangeType != "" {
		filter["change_type"] = params.ChangeType
	}

	// 时间范围查询
	if params.StartTime != "" || params.EndTime != "" {
		timeFilter := bson.M{}
		if params.StartTime != "" {
			if startTime, err := time.Parse("2006-01-02", params.StartTime); err == nil {
				timeFilter["$gte"] = startTime
			}
		}
		if params.EndTime != "" {
			if endTime, err := time.Parse("2006-01-02", params.EndTime); err == nil {
				// 结束时间设置为当天的23:59:59
				endTime = endTime.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
				timeFilter["$lte"] = endTime
			}
		}
		if len(timeFilter) > 0 {
			filter["change_time"] = timeFilter
		}
	}

	// 计算总数
	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	page := params.Page
	pageSize := params.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	skip := (page - 1) * pageSize
	opts := options.Find().
		SetSort(bson.D{{Key: "change_time", Value: -1}}).
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize))

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var records []*model.ChangeRecord
	for cursor.Next(ctx) {
		var record model.ChangeRecord
		if err := cursor.Decode(&record); err != nil {
			return nil, 0, err
		}
		records = append(records, &record)
	}

	return records, total, nil
}

// DeleteChangeRecordsByTableAndRecord 删除指定表和记录的变更记录（当记录被删除时使用）
func (r *ChangeRecordRepository) DeleteChangeRecordsByTableAndRecord(ctx context.Context, tableName, recordID string) error {
	collection := r.db.Collection(ChangeRecordCollection)

	filter := bson.M{
		"table_name": tableName,
		"record_id":  recordID,
	}

	_, err := collection.DeleteMany(ctx, filter)
	return err
}

// CleanupOldRecords 清理旧的变更记录（可以通过定时任务调用）
func (r *ChangeRecordRepository) CleanupOldRecords(ctx context.Context, beforeDays int) error {
	collection := r.db.Collection(ChangeRecordCollection)

	cutoffTime := time.Now().AddDate(0, 0, -beforeDays)
	filter := bson.M{
		"change_time": bson.M{"$lt": cutoffTime},
	}

	result, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}

	fmt.Printf("Cleaned up %d old change records\n", result.DeletedCount)
	return nil
}

// CreateIndexes 创建索引以优化查询性能
func (r *ChangeRecordRepository) CreateIndexes(ctx context.Context) error {
	collection := r.db.Collection(ChangeRecordCollection)

	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "table_name", Value: 1},
				{Key: "record_id", Value: 1},
				{Key: "change_time", Value: -1},
			},
		},
		{
			Keys: bson.D{
				{Key: "company_id", Value: 1},
				{Key: "change_time", Value: -1},
			},
		},
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "change_time", Value: -1},
			},
		},
		{
			Keys: bson.D{{Key: "change_time", Value: -1}},
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	return err
}
