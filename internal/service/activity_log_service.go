package service

import (
	"YufungProject/internal/model"
	"YufungProject/internal/repository"
	"context"
)

type ActivityLogService struct {
	repo *repository.ActivityLogRepository
}

func NewActivityLogService() *ActivityLogService {
	return &ActivityLogService{
		repo: repository.NewActivityLogRepository(),
	}
}

// CreateActivityLog 创建活动记录
func (s *ActivityLogService) CreateActivityLog(ctx context.Context, log *model.ActivityLog) error {
	return s.repo.Create(ctx, log)
}

// GetActivityLogList 获取活动记录列表
func (s *ActivityLogService) GetActivityLogList(ctx context.Context, query *model.ActivityLogQuery) (*model.ActivityLogResponse, error) {
	return s.repo.GetList(ctx, query)
}

// GetRecentActivityLogs 获取最近的活动记录（用于仪表盘）
func (s *ActivityLogService) GetRecentActivityLogs(ctx context.Context, companyID string, limit int) ([]model.ActivityLog, error) {
	return s.repo.GetRecentLogs(ctx, companyID, limit)
}

// GetActivityLogByID 根据ID获取活动记录
func (s *ActivityLogService) GetActivityLogByID(ctx context.Context, id string) (*model.ActivityLog, error) {
	return s.repo.GetByID(ctx, id)
}

// GetActivityLogStatistics 获取活动记录统计
func (s *ActivityLogService) GetActivityLogStatistics(ctx context.Context, companyID string, days int) (map[string]interface{}, error) {
	return s.repo.GetStatistics(ctx, companyID, days)
}

// DeleteActivityLogsByCompanyID 删除指定公司的所有活动记录
func (s *ActivityLogService) DeleteActivityLogsByCompanyID(ctx context.Context, companyID string) error {
	return s.repo.DeleteByCompanyID(ctx, companyID)
}

// LogUserActivity 记录用户活动（便捷方法）
func (s *ActivityLogService) LogUserActivity(ctx context.Context, userID, username, companyID, companyName, operationType, moduleName, operationDesc, requestURL, requestMethod, ipAddress, userAgent string, requestParams interface{}, executionTime int64, resultStatus string) error {
	log := &model.ActivityLog{
		UserID:        userID,
		Username:      username,
		CompanyID:     companyID,
		CompanyName:   companyName,
		OperationType: operationType,
		ModuleName:    moduleName,
		OperationDesc: operationDesc,
		RequestURL:    requestURL,
		RequestMethod: requestMethod,
		RequestParams: requestParams,
		IPAddress:     ipAddress,
		UserAgent:     userAgent,
		ExecutionTime: executionTime,
		ResultStatus:  resultStatus,
	}
	return s.CreateActivityLog(ctx, log)
}

// LogCreateActivity 记录创建操作
func (s *ActivityLogService) LogCreateActivity(ctx context.Context, userID, username, companyID, companyName, moduleName, targetName, targetID string) error {
	return s.LogUserActivity(ctx, userID, username, companyID, companyName, model.OperationTypeCreate, moduleName,
		"新增了"+targetName, "", "", "", "", nil, 0, "success")
}

// LogUpdateActivity 记录更新操作
func (s *ActivityLogService) LogUpdateActivity(ctx context.Context, userID, username, companyID, companyName, moduleName, targetName, targetID string) error {
	return s.LogUserActivity(ctx, userID, username, companyID, companyName, model.OperationTypeUpdate, moduleName,
		"更新了"+targetName, "", "", "", "", nil, 0, "success")
}

// LogDeleteActivity 记录删除操作
func (s *ActivityLogService) LogDeleteActivity(ctx context.Context, userID, username, companyID, companyName, moduleName, targetName, targetID string) error {
	return s.LogUserActivity(ctx, userID, username, companyID, companyName, model.OperationTypeDelete, moduleName,
		"删除了"+targetName, "", "", "", "", nil, 0, "success")
}

// LogViewActivity 记录查看操作
func (s *ActivityLogService) LogViewActivity(ctx context.Context, userID, username, companyID, companyName, moduleName, targetName string) error {
	return s.LogUserActivity(ctx, userID, username, companyID, companyName, model.OperationTypeView, moduleName,
		"查看了"+targetName, "", "", "", "", nil, 0, "success")
}

// LogExportActivity 记录导出操作
func (s *ActivityLogService) LogExportActivity(ctx context.Context, userID, username, companyID, companyName, moduleName, targetName string) error {
	return s.LogUserActivity(ctx, userID, username, companyID, companyName, model.OperationTypeExport, moduleName,
		"导出了"+targetName, "", "", "", "", nil, 0, "success")
}

// LogImportActivity 记录导入操作
func (s *ActivityLogService) LogImportActivity(ctx context.Context, userID, username, companyID, companyName, moduleName, targetName string) error {
	return s.LogUserActivity(ctx, userID, username, companyID, companyName, model.OperationTypeImport, moduleName,
		"导入了"+targetName, "", "", "", "", nil, 0, "success")
}

// LogLoginActivity 记录登录操作
func (s *ActivityLogService) LogLoginActivity(ctx context.Context, userID, username, companyID, companyName, ipAddress, userAgent string) error {
	return s.LogUserActivity(ctx, userID, username, companyID, companyName, model.OperationTypeLogin, model.ModuleAuth,
		"用户登录", "", "", ipAddress, userAgent, nil, 0, "success")
}

// LogLogoutActivity 记录登出操作
func (s *ActivityLogService) LogLogoutActivity(ctx context.Context, userID, username, companyID, companyName, ipAddress, userAgent string) error {
	return s.LogUserActivity(ctx, userID, username, companyID, companyName, model.OperationTypeLogout, model.ModuleAuth,
		"用户登出", "", "", ipAddress, userAgent, nil, 0, "success")
}
