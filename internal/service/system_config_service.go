package service

import (
	"context"
	"fmt"
	"strings"

	"YufungProject/internal/model"
	"YufungProject/internal/repository"
)

// SystemConfigService 系统配置服务接口
type SystemConfigService interface {
	CreateSystemConfig(ctx context.Context, req *model.SystemConfigCreateRequest, userID string) (*model.SystemConfigResponse, error)
	GetSystemConfigByID(ctx context.Context, configID string) (*model.SystemConfigResponse, error)
	UpdateSystemConfig(ctx context.Context, configID string, req *model.SystemConfigUpdateRequest, userID string) (*model.SystemConfigResponse, error)
	DeleteSystemConfig(ctx context.Context, configID string) error
	ListSystemConfigs(ctx context.Context, req *model.SystemConfigQueryRequest) (*model.SystemConfigListResponse, error)
	GetConfigsByType(ctx context.Context, configType string) ([]model.SystemConfigResponse, error)
}

type systemConfigService struct {
	systemConfigRepo repository.SystemConfigRepository
}

// NewSystemConfigService 创建系统配置服务实例
func NewSystemConfigService(systemConfigRepo repository.SystemConfigRepository) SystemConfigService {
	return &systemConfigService{
		systemConfigRepo: systemConfigRepo,
	}
}

// CreateSystemConfig 创建系统配置
func (s *systemConfigService) CreateSystemConfig(ctx context.Context, req *model.SystemConfigCreateRequest, userID string) (*model.SystemConfigResponse, error) {
	// 验证配置键是否已存在
	exists, err := s.systemConfigRepo.CheckKeyExists(ctx, req.ConfigType, req.ConfigKey, "", "")
	if err != nil {
		return nil, fmt.Errorf("检查配置键失败: %v", err)
	}
	if exists {
		return nil, fmt.Errorf("配置键已存在")
	}

	// 设置默认值
	if req.Status == "" {
		req.Status = "enable"
	}

	// 创建系统配置模型
	config := &model.SystemConfig{
		ConfigType:  req.ConfigType,
		ConfigKey:   strings.TrimSpace(req.ConfigKey),
		ConfigValue: strings.TrimSpace(req.ConfigValue),
		DisplayName: strings.TrimSpace(req.DisplayName),
		CompanyID:   "", // 移除公司ID限制
		SortOrder:   req.SortOrder,
		Status:      req.Status,
		Remark:      req.Remark,
		CreatedBy:   userID,
		UpdatedBy:   userID,
	}

	// 保存到数据库
	if err := s.systemConfigRepo.Create(ctx, config); err != nil {
		return nil, fmt.Errorf("创建系统配置失败: %v", err)
	}

	return s.GetSystemConfigByID(ctx, config.ConfigID)
}

// GetSystemConfigByID 根据ID获取系统配置
func (s *systemConfigService) GetSystemConfigByID(ctx context.Context, configID string) (*model.SystemConfigResponse, error) {
	config, err := s.systemConfigRepo.GetByID(ctx, configID)
	if err != nil {
		return nil, err
	}

	// 移除公司权限检查
	// if config.CompanyID != companyID {
	// 	return nil, fmt.Errorf("无权限访问该系统配置")
	// }

	return s.convertToResponse(config), nil
}

// UpdateSystemConfig 更新系统配置
func (s *systemConfigService) UpdateSystemConfig(ctx context.Context, configID string, req *model.SystemConfigUpdateRequest, userID string) (*model.SystemConfigResponse, error) {
	// 获取现有配置
	existingConfig, err := s.systemConfigRepo.GetByID(ctx, configID)
	if err != nil {
		return nil, err
	}

	// 移除公司权限检查
	// if existingConfig.CompanyID != companyID {
	// 	return nil, fmt.Errorf("无权限访问该系统配置")
	// }

	// 更新字段
	if req.ConfigValue != "" {
		existingConfig.ConfigValue = strings.TrimSpace(req.ConfigValue)
	}
	if req.DisplayName != "" {
		existingConfig.DisplayName = strings.TrimSpace(req.DisplayName)
	}
	if req.SortOrder > 0 {
		existingConfig.SortOrder = req.SortOrder
	}
	if req.Status != "" {
		existingConfig.Status = req.Status
	}
	if req.Remark != "" {
		existingConfig.Remark = req.Remark
	}
	existingConfig.UpdatedBy = userID

	// 保存更新
	if err := s.systemConfigRepo.Update(ctx, configID, existingConfig); err != nil {
		return nil, fmt.Errorf("更新系统配置失败: %v", err)
	}

	return s.GetSystemConfigByID(ctx, configID)
}

// DeleteSystemConfig 删除系统配置
func (s *systemConfigService) DeleteSystemConfig(ctx context.Context, configID string) error {
	// 检查配置是否存在
	_, err := s.systemConfigRepo.GetByID(ctx, configID)
	if err != nil {
		return err
	}

	// 移除公司权限检查
	// if config.CompanyID != companyID {
	// 	return fmt.Errorf("无权限删除该系统配置")
	// }

	return s.systemConfigRepo.Delete(ctx, configID)
}

// ListSystemConfigs 获取系统配置列表
func (s *systemConfigService) ListSystemConfigs(ctx context.Context, req *model.SystemConfigQueryRequest) (*model.SystemConfigListResponse, error) {
	result, err := s.systemConfigRepo.List(ctx, req, "") // 移除公司ID限制
	if err != nil {
		return nil, err
	}

	// 转换为响应格式
	var responseList []model.SystemConfigResponse
	for _, config := range result.List {
		responseList = append(responseList, *s.convertToResponse(&config))
	}

	return &model.SystemConfigListResponse{
		List:  result.List,
		Total: result.Total,
	}, nil
}

// GetConfigsByType 根据配置类型获取配置选项
func (s *systemConfigService) GetConfigsByType(ctx context.Context, configType string) ([]model.SystemConfigResponse, error) {
	configs, err := s.systemConfigRepo.GetByType(ctx, configType, "") // 移除公司ID限制
	if err != nil {
		return nil, err
	}

	var responses []model.SystemConfigResponse
	for _, config := range configs {
		responses = append(responses, *s.convertToResponse(&config))
	}

	return responses, nil
}

// convertToResponse 转换为响应格式
func (s *systemConfigService) convertToResponse(config *model.SystemConfig) *model.SystemConfigResponse {
	statusText := "启用"
	if config.Status == "disable" {
		statusText = "禁用"
	}

	return &model.SystemConfigResponse{
		ID:          config.ID.Hex(),
		ConfigID:    config.ConfigID,
		ConfigType:  config.ConfigType,
		ConfigKey:   config.ConfigKey,
		ConfigValue: config.ConfigValue,
		DisplayName: config.DisplayName,
		CompanyID:   config.CompanyID,
		SortOrder:   config.SortOrder,
		Status:      config.Status,
		StatusText:  statusText,
		Remark:      config.Remark,
		CreatedBy:   config.CreatedBy,
		UpdatedBy:   config.UpdatedBy,
		CreatedAt:   config.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   config.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
