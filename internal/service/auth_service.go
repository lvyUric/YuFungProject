package service

import (
	"context"
	"errors"
	"time"

	"YufungProject/configs"
	"YufungProject/internal/model"
	"YufungProject/internal/repository"
	"YufungProject/pkg/logger"
	"YufungProject/pkg/utils"
)

// AuthService 认证服务接口
type AuthService interface {
	Login(ctx context.Context, req *model.LoginRequest) (*model.LoginResponse, error)
	Register(ctx context.Context, req *model.RegisterRequest) (*model.UserInfo, error)
	ChangePassword(ctx context.Context, userID string, req *model.ChangePasswordRequest) error
	RefreshToken(ctx context.Context, refreshToken string) (*model.LoginResponse, error)
	Logout(ctx context.Context, userID string, token string) error
	ResetPassword(ctx context.Context, req *model.ResetPasswordRequest) error
	GetUserInfo(ctx context.Context, userID string) (*model.UserInfo, error)
}

type authService struct {
	userRepo repository.UserRepository
	config   *configs.Config
	jwtUtil  *utils.JWTUtil
}

// NewAuthService 创建认证服务实例
func NewAuthService(userRepo repository.UserRepository, config *configs.Config) AuthService {
	// 解析时间配置
	expiresIn, _ := time.ParseDuration(config.JWT.ExpiresIn)
	refreshExpiresIn, _ := time.ParseDuration(config.JWT.RefreshExpiresIn)

	jwtUtil := utils.NewJWTUtil(config.JWT.Secret, expiresIn, refreshExpiresIn)

	return &authService{
		userRepo: userRepo,
		config:   config,
		jwtUtil:  jwtUtil,
	}
}

// Login 用户登录
func (s *authService) Login(ctx context.Context, req *model.LoginRequest) (*model.LoginResponse, error) {
	// 查找用户
	user, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		logger.Warnf("用户登录失败 - 用户不存在: %s", req.Username)
		return nil, errors.New("用户不存在")
	}

	// 检查账户状态
	if user.Status == "inactive" {
		logger.Warnf("用户登录失败 - 账户已禁用: %s", req.Username)
		return nil, errors.New("账户已被禁用")
	}

	// 检查账户是否被锁定
	if user.LockedUntil != nil && time.Now().Before(*user.LockedUntil) {
		logger.Warnf("用户登录失败 - 账户被锁定: %s", req.Username)
		return nil, errors.New("账户已被锁定")
	}

	// 验证密码
	if !utils.CheckPassword(user.PasswordHash, req.Password) {
		// 增加登录失败次数
		newAttempts := user.LoginAttempts + 1
		var lockedUntil *time.Time

		// 检查是否需要锁定账户
		if newAttempts >= s.config.Security.MaxLoginAttempts {
			lockDuration, _ := time.ParseDuration(s.config.Security.LockoutDuration)
			lockTime := time.Now().Add(lockDuration)
			lockedUntil = &lockTime
			logger.Warnf("用户账户已锁定: %s, 锁定至: %v", req.Username, lockTime)
		}

		// 更新登录尝试次数
		s.userRepo.UpdateLoginAttempts(ctx, user.UserID, newAttempts, lockedUntil)

		logger.Warnf("用户登录失败 - 密码错误: %s, 尝试次数: %d", req.Username, newAttempts)
		return nil, errors.New("密码错误")
	}

	// 登录成功，重置登录尝试次数和更新最后登录时间
	now := time.Now()
	s.userRepo.UpdateLoginAttempts(ctx, user.UserID, 0, nil)
	s.userRepo.UpdateLastLoginTime(ctx, user.UserID, now)

	// 生成JWT令牌
	token, expiresAt, err := s.jwtUtil.GenerateToken(user.UserID, user.Username, user.CompanyID, user.RoleIDs)
	if err != nil {
		logger.Errorf("生成JWT令牌失败: %v", err)
		return nil, errors.New("生成令牌失败")
	}

	// 生成刷新令牌
	refreshToken, err := s.jwtUtil.GenerateRefreshToken(user.UserID)
	if err != nil {
		logger.Errorf("生成刷新令牌失败: %v", err)
		return nil, errors.New("生成刷新令牌失败")
	}

	// 构建用户信息
	userInfo := &model.UserInfo{
		ID:          user.ID.Hex(),
		UserID:      user.UserID,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		CompanyID:   user.CompanyID,
		RoleIDs:     user.RoleIDs,
		Status:      user.Status,
		Email:       user.Email,
		Phone:       user.Phone,
		LastLogin:   &now,
	}

	logger.Infof("用户登录成功: %s", req.Username)

	return &model.LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		User:         *userInfo,
	}, nil
}

// Register 用户注册
func (s *authService) Register(ctx context.Context, req *model.RegisterRequest) (*model.UserInfo, error) {
	// 检查用户名是否已存在
	exists, err := s.userRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		logger.Errorf("检查用户名是否存在失败: %v", err)
		return nil, errors.New("注册失败")
	}
	if exists {
		logger.Warnf("注册失败 - 用户名已存在: %s", req.Username)
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在（如果提供了邮箱）
	if req.Email != "" {
		exists, err := s.userRepo.ExistsByEmail(ctx, req.Email)
		if err != nil {
			logger.Errorf("检查邮箱是否存在失败: %v", err)
			return nil, errors.New("注册失败")
		}
		if exists {
			logger.Warnf("注册失败 - 邮箱已存在: %s", req.Email)
			return nil, errors.New("邮箱已存在")
		}
	}

	// 验证密码强度
	if len(req.Password) < s.config.Security.PasswordMinLength || !utils.ValidatePassword(req.Password) {
		logger.Warnf("注册失败 - 密码强度不足: %s", req.Username)
		return nil, errors.New("密码强度不足")
	}

	// 加密密码
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		logger.Errorf("密码加密失败: %v", err)
		return nil, errors.New("注册失败")
	}

	// 生成用户ID
	userID := utils.GenerateUserID()

	// 创建用户对象
	now := time.Now()
	user := &model.User{
		UserID:        userID,
		Username:      req.Username,
		DisplayName:   req.DisplayName,
		CompanyID:     "CMP_PLATFORM_001",          // 默认平台公司
		RoleIDs:       []string{"ROL_NORMAL_USER"}, // 默认普通用户角色
		Status:        "active",
		PasswordHash:  passwordHash,
		Email:         req.Email,
		Phone:         req.Phone,
		LoginAttempts: 0,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	// 保存用户
	err = s.userRepo.Create(ctx, user)
	if err != nil {
		logger.Errorf("保存用户失败: %v", err)
		return nil, errors.New("注册失败")
	}

	logger.Infof("用户注册成功: %s", req.Username)

	// 返回用户信息
	return &model.UserInfo{
		ID:          user.ID.Hex(),
		UserID:      user.UserID,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		CompanyID:   user.CompanyID,
		RoleIDs:     user.RoleIDs,
		Status:      user.Status,
		Email:       user.Email,
		Phone:       user.Phone,
	}, nil
}

// ChangePassword 修改密码
func (s *authService) ChangePassword(ctx context.Context, userID string, req *model.ChangePasswordRequest) error {
	// 获取用户信息
	user, err := s.userRepo.GetByUserID(ctx, userID)
	if err != nil {
		logger.Errorf("获取用户信息失败: %v, UserID: %s", err, userID)
		return errors.New("用户不存在")
	}

	// 验证当前密码
	if !utils.CheckPassword(user.PasswordHash, req.OldPassword) {
		logger.Warnf("修改密码失败 - 当前密码错误: %s", userID)
		return errors.New("当前密码错误")
	}

	// 验证新密码强度
	if len(req.NewPassword) < s.config.Security.PasswordMinLength || !utils.ValidatePassword(req.NewPassword) {
		logger.Warnf("修改密码失败 - 新密码强度不足: %s", userID)
		return errors.New("新密码强度不足")
	}

	// 加密新密码
	newPasswordHash, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		logger.Errorf("新密码加密失败: %v", err)
		return errors.New("密码修改失败")
	}

	// 更新密码
	err = s.userRepo.UpdatePassword(ctx, userID, newPasswordHash)
	if err != nil {
		logger.Errorf("更新密码失败: %v, UserID: %s", err, userID)
		return errors.New("密码修改失败")
	}

	logger.Infof("用户密码修改成功: %s", userID)
	return nil
}

// RefreshToken 刷新令牌
func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*model.LoginResponse, error) {
	// 解析刷新令牌
	userID, err := s.jwtUtil.ParseRefreshToken(refreshToken)
	if err != nil {
		logger.Warnf("刷新令牌解析失败: %v", err)
		return nil, errors.New("刷新令牌无效")
	}

	// 获取用户信息
	user, err := s.userRepo.GetByUserID(ctx, userID)
	if err != nil {
		logger.Errorf("获取用户信息失败: %v, UserID: %s", err, userID)
		return nil, errors.New("用户不存在")
	}

	// 检查用户状态
	if user.Status != "active" {
		logger.Warnf("刷新令牌失败 - 用户状态异常: %s, Status: %s", userID, user.Status)
		return nil, errors.New("用户状态异常")
	}

	// 生成新的访问令牌
	token, expiresAt, err := s.jwtUtil.GenerateToken(user.UserID, user.Username, user.CompanyID, user.RoleIDs)
	if err != nil {
		logger.Errorf("生成新JWT令牌失败: %v", err)
		return nil, errors.New("生成令牌失败")
	}

	// 生成新的刷新令牌
	newRefreshToken, err := s.jwtUtil.GenerateRefreshToken(user.UserID)
	if err != nil {
		logger.Errorf("生成新刷新令牌失败: %v", err)
		return nil, errors.New("生成刷新令牌失败")
	}

	// 构建用户信息
	userInfo := &model.UserInfo{
		ID:          user.ID.Hex(),
		UserID:      user.UserID,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		CompanyID:   user.CompanyID,
		RoleIDs:     user.RoleIDs,
		Status:      user.Status,
		Email:       user.Email,
		Phone:       user.Phone,
		LastLogin:   user.LastLoginTime,
	}

	logger.Infof("令牌刷新成功: %s", userID)

	return &model.LoginResponse{
		Token:        token,
		RefreshToken: newRefreshToken,
		ExpiresAt:    expiresAt,
		User:         *userInfo,
	}, nil
}

// Logout 用户登出
func (s *authService) Logout(ctx context.Context, userID string, token string) error {
	// 这里可以实现令牌黑名单逻辑
	// 目前简单记录日志
	logger.Infof("用户登出: %s", userID)
	return nil
}

// ResetPassword 重置密码（管理员功能）
func (s *authService) ResetPassword(ctx context.Context, req *model.ResetPasswordRequest) error {
	// 验证新密码强度
	if len(req.NewPassword) < s.config.Security.PasswordMinLength || !utils.ValidatePassword(req.NewPassword) {
		logger.Warnf("重置密码失败 - 新密码强度不足: %s", req.UserID)
		return errors.New("新密码强度不足")
	}

	// 加密新密码
	newPasswordHash, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		logger.Errorf("新密码加密失败: %v", err)
		return errors.New("密码重置失败")
	}

	// 更新密码
	err = s.userRepo.UpdatePassword(ctx, req.UserID, newPasswordHash)
	if err != nil {
		logger.Errorf("重置密码失败: %v, UserID: %s", err, req.UserID)
		return errors.New("密码重置失败")
	}

	logger.Infof("管理员重置用户密码成功: %s", req.UserID)
	return nil
}

// GetUserInfo 获取用户信息
func (s *authService) GetUserInfo(ctx context.Context, userID string) (*model.UserInfo, error) {
	user, err := s.userRepo.GetByUserID(ctx, userID)
	if err != nil {
		logger.Errorf("获取用户信息失败: %v, UserID: %s", err, userID)
		return nil, errors.New("用户不存在")
	}

	return &model.UserInfo{
		ID:          user.ID.Hex(),
		UserID:      user.UserID,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		CompanyID:   user.CompanyID,
		RoleIDs:     user.RoleIDs,
		Status:      user.Status,
		Email:       user.Email,
		Phone:       user.Phone,
		LastLogin:   user.LastLoginTime,
	}, nil
}
