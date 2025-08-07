package service

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"YufungProject/internal/model"
	"YufungProject/internal/repository"
	"YufungProject/pkg/logger"
)

type ChangeRecordService struct {
	changeRecordRepo *repository.ChangeRecordRepository
	userRepo         repository.UserRepository
}

func NewChangeRecordService(changeRecordRepo *repository.ChangeRecordRepository, userRepo repository.UserRepository) *ChangeRecordService {
	return &ChangeRecordService{
		changeRecordRepo: changeRecordRepo,
		userRepo:         userRepo,
	}
}

// RecordChange 记录数据变更
func (s *ChangeRecordService) RecordChange(ctx context.Context, tableName, recordID, userID, companyID, changeType string, oldData, newData interface{}, changeReason, ipAddress, userAgent string) error {
	// 获取用户信息
	user, err := s.userRepo.GetByUserID(ctx, userID)
	if err != nil {
		logger.Error("Failed to get user info", "user_id", userID, "error", err)
		// 即使获取用户信息失败，也要记录变更，使用默认用户名
	}

	username := "unknown"
	if user != nil {
		username = user.Username
	}

	// 比较数据变更
	changedFields, oldValues, newValues := s.compareData(oldData, newData, changeType)

	// 如果没有变更字段且不是插入或删除操作，则不记录
	if len(changedFields) == 0 && changeType == "update" {
		return nil
	}

	record := &model.ChangeRecord{
		TableName:     tableName,
		RecordID:      recordID,
		UserID:        userID,
		Username:      username,
		CompanyID:     companyID,
		ChangeType:    changeType,
		OldValues:     oldValues,
		NewValues:     newValues,
		ChangedFields: changedFields,
		ChangeReason:  changeReason,
		IPAddress:     ipAddress,
		UserAgent:     userAgent,
	}

	return s.changeRecordRepo.CreateChangeRecord(ctx, record)
}

// GetChangeRecordsByPolicy 获取保单的变更记录
func (s *ChangeRecordService) GetChangeRecordsByPolicy(ctx context.Context, policyID string, days int, page, pageSize int) ([]*model.ChangeRecordResponse, int64, error) {
	records, total, err := s.changeRecordRepo.GetChangeRecordsByTableAndRecord(ctx, "policies", policyID, days, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	var responses []*model.ChangeRecordResponse
	for _, record := range records {
		response := s.convertToResponse(record)
		responses = append(responses, response)
	}

	return responses, total, nil
}

// GetChangeRecordsList 获取变更记录列表
func (s *ChangeRecordService) GetChangeRecordsList(ctx context.Context, params *model.ChangeRecordListParams) ([]*model.ChangeRecordResponse, int64, error) {
	records, total, err := s.changeRecordRepo.GetChangeRecordsList(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	var responses []*model.ChangeRecordResponse
	for _, record := range records {
		response := s.convertToResponse(record)
		responses = append(responses, response)
	}

	return responses, total, nil
}

// convertToResponse 转换为响应格式
func (s *ChangeRecordService) convertToResponse(record *model.ChangeRecord) *model.ChangeRecordResponse {
	response := &model.ChangeRecordResponse{
		ID:                  record.ID.Hex(),
		ChangeID:            record.ChangeID,
		TableName:           record.TableName,
		RecordID:            record.RecordID,
		UserID:              record.UserID,
		Username:            record.Username,
		CompanyID:           record.CompanyID,
		ChangeType:          record.ChangeType,
		OldValues:           record.OldValues,
		NewValues:           record.NewValues,
		ChangedFields:       record.ChangedFields,
		ChangeTime:          record.ChangeTime,
		ChangeReason:        record.ChangeReason,
		IPAddress:           record.IPAddress,
		UserAgent:           record.UserAgent,
		ChangeTimeFormatted: record.ChangeTime.Format("2006-01-02 15:04:05"),
	}

	// 生成变更详情
	response.ChangeDetails = s.generateChangeDetails(record.TableName, record.ChangedFields, record.OldValues, record.NewValues, record.ChangeType)

	return response
}

// generateChangeDetails 生成变更详情
func (s *ChangeRecordService) generateChangeDetails(tableName string, changedFields []string, oldValues, newValues map[string]interface{}, changeType string) []model.ChangeDetail {
	var details []model.ChangeDetail

	for _, fieldName := range changedFields {
		detail := model.ChangeDetail{
			FieldName:  fieldName,
			FieldLabel: model.GetFieldLabel(tableName, fieldName),
		}

		// 获取旧值和新值
		if oldValues != nil {
			detail.OldValue = oldValues[fieldName]
		}
		if newValues != nil {
			detail.NewValue = newValues[fieldName]
		}

		// 格式化显示文本
		detail.OldValueText = s.formatValue(detail.OldValue, fieldName)
		detail.NewValueText = s.formatValue(detail.NewValue, fieldName)

		details = append(details, detail)
	}

	return details
}

// formatValue 格式化值为显示文本
func (s *ChangeRecordService) formatValue(value interface{}, fieldName string) string {
	if value == nil {
		return "无"
	}

	switch v := value.(type) {
	case bool:
		if v {
			if strings.Contains(fieldName, "surrendered") {
				return "是"
			}
			return "是"
		}
		if strings.Contains(fieldName, "surrendered") {
			return "否"
		}
		return "否"
	case float64:
		// 对于金额字段，添加货币符号
		if strings.Contains(fieldName, "premium") || strings.Contains(fieldName, "aum") {
			return fmt.Sprintf("%.2f", v)
		}
		return fmt.Sprintf("%.2f", v)
	case string:
		if v == "" {
			return "无"
		}
		return v
	case time.Time:
		return v.Format("2006-01-02")
	default:
		return fmt.Sprintf("%v", v)
	}
}

// compareData 比较数据变更
func (s *ChangeRecordService) compareData(oldData, newData interface{}, changeType string) ([]string, map[string]interface{}, map[string]interface{}) {
	var changedFields []string
	oldValues := make(map[string]interface{})
	newValues := make(map[string]interface{})

	switch changeType {
	case "insert":
		// 插入操作，记录所有新值
		if newData != nil {
			newMap := s.structToMap(newData)
			for k, v := range newMap {
				if s.shouldTrackField(k) {
					changedFields = append(changedFields, k)
					newValues[k] = v
				}
			}
		}
	case "delete":
		// 删除操作，记录所有旧值
		if oldData != nil {
			oldMap := s.structToMap(oldData)
			for k, v := range oldMap {
				if s.shouldTrackField(k) {
					changedFields = append(changedFields, k)
					oldValues[k] = v
				}
			}
		}
	case "update":
		// 更新操作，比较新旧值
		if oldData != nil && newData != nil {
			oldMap := s.structToMap(oldData)
			newMap := s.structToMap(newData)

			for k, newVal := range newMap {
				if !s.shouldTrackField(k) {
					continue
				}

				oldVal, exists := oldMap[k]
				if !exists || !s.valuesEqual(oldVal, newVal) {
					changedFields = append(changedFields, k)
					if exists {
						oldValues[k] = oldVal
					}
					newValues[k] = newVal
				}
			}
		}
	}

	return changedFields, oldValues, newValues
}

// structToMap 将结构体转换为map
func (s *ChangeRecordService) structToMap(data interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	if data == nil {
		return result
	}

	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return result
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// 获取json tag作为字段名
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		// 处理json tag中的选项（如omitempty）
		jsonName := strings.Split(jsonTag, ",")[0]
		if jsonName == "" {
			jsonName = field.Name
		}

		// 跳过私有字段
		if !value.CanInterface() {
			continue
		}

		result[jsonName] = value.Interface()
	}

	return result
}

// valuesEqual 比较两个值是否相等
func (s *ChangeRecordService) valuesEqual(a, b interface{}) bool {
	// 处理时间类型的比较
	if timeA, okA := a.(time.Time); okA {
		if timeB, okB := b.(time.Time); okB {
			return timeA.Equal(timeB)
		}
	}

	// 处理指针类型的时间比较
	if ptrA, okA := a.(*time.Time); okA {
		if ptrB, okB := b.(*time.Time); okB {
			if ptrA == nil && ptrB == nil {
				return true
			}
			if ptrA == nil || ptrB == nil {
				return false
			}
			return ptrA.Equal(*ptrB)
		}
	}

	return reflect.DeepEqual(a, b)
}

// shouldTrackField 判断是否应该跟踪该字段的变更
func (s *ChangeRecordService) shouldTrackField(fieldName string) bool {
	// 跳过系统字段和敏感字段
	skipFields := map[string]bool{
		"id":            true,
		"_id":           true,
		"created_at":    true,
		"updated_at":    true,
		"created_by":    true,
		"updated_by":    true,
		"password":      true,
		"password_hash": true,
	}

	return !skipFields[fieldName]
}
