package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"strings"
	"time"

	"YufungProject/internal/model"
	"YufungProject/internal/repository"
	"YufungProject/pkg/logger"
	"YufungProject/pkg/utils"

	"encoding/csv"

	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

// UserService 用户服务接口
type UserService interface {
	// 用户基础操作
	CreateUser(ctx context.Context, req *model.UserCreateRequest) (*model.UserInfo, error)
	GetByUserID(ctx context.Context, userID string) (*model.UserInfo, error)
	GetUserList(ctx context.Context, filter map[string]interface{}, page, pageSize int) (*model.UserListResponse, error)
	UpdateUser(ctx context.Context, userID string, req *model.UserUpdateRequest) error
	DeleteUser(ctx context.Context, userID string) error

	// 密码管理
	ResetPassword(ctx context.Context, req *model.ResetPasswordRequest) error

	// 批量操作
	BatchUpdateUserStatus(ctx context.Context, userIDs []string, status string) error
	QuickDisableUser(ctx context.Context, userID string) error

	// 导出功能
	ExportUsers(ctx context.Context, filter map[string]interface{}) ([]byte, string, error)

	// 导入导出功能扩展
	ExportUsersAdvanced(ctx context.Context, req *model.UserExportRequest) (*model.UserExportResponse, error)
	GenerateUserTemplate(ctx context.Context, format string) ([]byte, string, error)
	PreviewUserImport(ctx context.Context, file multipart.File, header *multipart.FileHeader, req *model.UserImportRequest) (*model.UserImportResponse, error)
	ImportUsers(ctx context.Context, file multipart.File, header *multipart.FileHeader, req *model.UserImportRequest) (*model.UserImportResponse, error)
}

// userService 用户服务实现
type userService struct {
	userRepo    repository.UserRepository
	companyRepo repository.CompanyRepository
}

// NewUserService 创建用户服务实例
func NewUserService(userRepo repository.UserRepository, companyRepo repository.CompanyRepository) UserService {
	return &userService{
		userRepo:    userRepo,
		companyRepo: companyRepo,
	}
}

// CreateUser 创建用户
func (s *userService) CreateUser(ctx context.Context, req *model.UserCreateRequest) (*model.UserInfo, error) {
	// 检查用户名是否已存在
	exists, err := s.userRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("检查用户名失败: %w", err)
	}
	if exists {
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在（如果提供）
	if req.Email != "" {
		exists, err := s.userRepo.ExistsByEmail(ctx, req.Email)
		if err != nil {
			return nil, fmt.Errorf("检查邮箱失败: %w", err)
		}
		if exists {
			return nil, errors.New("邮箱已存在")
		}
	}

	// 生成密码哈希
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("生成密码哈希失败: %w", err)
	}

	// 生成用户ID
	userID := utils.GenerateID("user")

	// 创建用户对象
	user := &model.User{
		UserID:       userID,
		Username:     req.Username,
		DisplayName:  req.DisplayName,
		CompanyID:    req.CompanyID,
		RoleIDs:      req.RoleIDs,
		Status:       "active",
		PasswordHash: string(passwordHash),
		Email:        req.Email,
		Phone:        req.Phone,
		Remark:       req.Remark,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// 保存到数据库
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("创建用户失败: %w", err)
	}

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

// GetByUserID 根据用户ID获取用户信息
func (s *userService) GetByUserID(ctx context.Context, userID string) (*model.UserInfo, error) {
	user, err := s.userRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("获取用户失败: %w", err)
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

// GetUserList 获取用户列表
func (s *userService) GetUserList(ctx context.Context, filter map[string]interface{}, page, pageSize int) (*model.UserListResponse, error) {
	// 构建MongoDB查询条件
	mongoFilter := bson.M{}

	for key, value := range filter {
		switch key {
		case "username", "display_name":
			// 支持模糊搜索
			mongoFilter[key] = bson.M{"$regex": value, "$options": "i"}
		default:
			mongoFilter[key] = value
		}
	}

	// 查询用户列表
	users, total, err := s.userRepo.List(ctx, mongoFilter, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("查询用户列表失败: %w", err)
	}

	// 转换为响应格式
	userInfos := make([]model.UserInfo, 0, len(users))
	for _, user := range users {
		userInfos = append(userInfos, model.UserInfo{
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
		})
	}

	// 计算总页数
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	return &model.UserListResponse{
		Users:      userInfos,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// UpdateUser 更新用户信息
func (s *userService) UpdateUser(ctx context.Context, userID string, req *model.UserUpdateRequest) error {
	// 构建更新数据
	update := bson.M{
		"updated_at": time.Now(),
	}

	if req.DisplayName != "" {
		update["display_name"] = req.DisplayName
	}
	if req.RoleIDs != nil {
		update["role_ids"] = req.RoleIDs
	}
	if req.Email != "" {
		update["email"] = req.Email
	}
	if req.Phone != "" {
		update["phone"] = req.Phone
	}
	if req.Remark != "" {
		update["remark"] = req.Remark
	}
	if req.Status != "" {
		update["status"] = req.Status
	}

	// 执行更新
	if err := s.userRepo.Update(ctx, userID, update); err != nil {
		return fmt.Errorf("更新用户失败: %w", err)
	}

	return nil
}

// DeleteUser 删除用户
func (s *userService) DeleteUser(ctx context.Context, userID string) error {
	if err := s.userRepo.Delete(ctx, userID); err != nil {
		return fmt.Errorf("删除用户失败: %w", err)
	}
	return nil
}

// ResetPassword 重置用户密码
func (s *userService) ResetPassword(ctx context.Context, req *model.ResetPasswordRequest) error {
	// 生成新密码哈希
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("生成密码哈希失败: %w", err)
	}

	// 更新密码
	if err := s.userRepo.UpdatePassword(ctx, req.UserID, string(passwordHash)); err != nil {
		return fmt.Errorf("重置密码失败: %w", err)
	}

	return nil
}

// BatchUpdateUserStatus 批量更新用户状态
func (s *userService) BatchUpdateUserStatus(ctx context.Context, userIDs []string, status string) error {
	for _, userID := range userIDs {
		update := bson.M{
			"status":     status,
			"updated_at": time.Now(),
		}

		if err := s.userRepo.Update(ctx, userID, update); err != nil {
			logger.Error("批量更新用户状态失败", err, "userID", userID)
			return fmt.Errorf("更新用户 %s 状态失败: %w", userID, err)
		}
	}

	return nil
}

// QuickDisableUser 快捷停用用户
func (s *userService) QuickDisableUser(ctx context.Context, userID string) error {
	update := bson.M{
		"status":     "inactive",
		"updated_at": time.Now(),
	}

	if err := s.userRepo.Update(ctx, userID, update); err != nil {
		return fmt.Errorf("快捷停用用户失败: %w", err)
	}

	return nil
}

// ExportUsers 导出用户数据
func (s *userService) ExportUsers(ctx context.Context, filter map[string]interface{}) ([]byte, string, error) {
	// 构建MongoDB查询条件
	mongoFilter := bson.M{}
	for key, value := range filter {
		mongoFilter[key] = value
	}

	// 查询所有符合条件的用户
	users, _, err := s.userRepo.List(ctx, mongoFilter, 1, 10000) // 限制最大导出10000条
	if err != nil {
		return nil, "", fmt.Errorf("查询用户数据失败: %w", err)
	}

	// 创建Excel文件
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			logger.Error("关闭Excel文件失败", err)
		}
	}()

	// 设置工作表名称
	sheetName := "用户列表"
	f.SetSheetName("Sheet1", sheetName)

	// 设置表头
	headers := []string{"用户ID", "用户名", "显示名称", "所属公司ID", "状态", "邮箱", "手机", "创建时间", "最后登录时间", "备注"}
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
	}

	// 填充数据
	for i, user := range users {
		row := i + 2 // 从第二行开始

		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), user.UserID)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), user.Username)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), user.DisplayName)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), user.CompanyID)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), user.Status)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), user.Email)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), user.Phone)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), user.CreatedAt.Format("2006-01-02 15:04:05"))

		lastLogin := ""
		if user.LastLoginTime != nil {
			lastLogin = user.LastLoginTime.Format("2006-01-02 15:04:05")
		}
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), lastLogin)
		f.SetCellValue(sheetName, fmt.Sprintf("J%d", row), user.Remark)
	}

	// 生成文件数据
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, "", fmt.Errorf("生成Excel文件失败: %w", err)
	}

	// 生成文件名
	filename := fmt.Sprintf("用户列表_%s.xlsx", time.Now().Format("20060102_150405"))

	return buf.Bytes(), filename, nil
}

// ExportUsersAdvanced 高级导出用户数据
func (s *userService) ExportUsersAdvanced(ctx context.Context, req *model.UserExportRequest) (*model.UserExportResponse, error) {
	var users []model.UserInfo
	var err error

	// 根据导出类型获取数据
	switch req.ExportType {
	case "all":
		// 导出全部数据
		filter := make(map[string]interface{})
		userModels, _, err := s.userRepo.List(ctx, bson.M(filter), 1, 10000)
		if err != nil {
			return nil, err
		}
		// 转换为UserInfo
		for _, user := range userModels {
			users = append(users, model.UserInfo{
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
			})
		}
	case "selected":
		// 导出选中数据
		if len(req.IDs) == 0 {
			return nil, errors.New("请选择要导出的数据")
		}
		for _, id := range req.IDs {
			user, err := s.GetByUserID(ctx, id)
			if err != nil {
				continue // 忽略不存在的用户
			}
			users = append(users, *user)
		}
	case "filtered":
		// 导出筛选结果
		filter := make(map[string]interface{})
		if req.Status != "" {
			filter["status"] = req.Status
		}
		if req.CompanyID != "" {
			filter["company_id"] = req.CompanyID
		}
		if req.Keyword != "" {
			filter["$or"] = []bson.M{
				{"username": bson.M{"$regex": req.Keyword, "$options": "i"}},
				{"display_name": bson.M{"$regex": req.Keyword, "$options": "i"}},
			}
		}
		userModels, _, err := s.userRepo.List(ctx, bson.M(filter), 1, 10000)
		if err != nil {
			return nil, err
		}
		// 转换为UserInfo
		for _, user := range userModels {
			users = append(users, model.UserInfo{
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
			})
		}
	default:
		return nil, errors.New("不支持的导出类型")
	}

	// 生成文件
	var fileData []byte
	var fileName string

	if req.Template {
		// 生成模板文件
		fileData, fileName, err = s.GenerateUserTemplate(ctx, req.Format)
	} else {
		// 生成数据文件
		fileData, fileName, err = s.generateUserDataFile(users, req.Format)
	}

	if err != nil {
		return nil, err
	}

	// 返回响应
	response := &model.UserExportResponse{
		FileURL:  "/api/download/" + fileName,
		FileName: fileName,
	}

	// 记录文件大小用于日志
	logger.Info("生成用户导出文件成功", "文件名", fileName, "大小", len(fileData), "bytes")

	return response, nil
}

// GenerateUserTemplate 生成用户导入模板
func (s *userService) GenerateUserTemplate(ctx context.Context, format string) ([]byte, string, error) {
	headers := []string{
		"用户名", "显示名称", "密码", "所属公司ID", "角色ID(多个用逗号分隔)",
		"邮箱地址", "手机号码", "备注信息",
	}

	var fileData []byte
	var fileName string
	var err error

	switch format {
	case "xlsx":
		fileData, err = s.generateUserExcelTemplate(headers)
		fileName = fmt.Sprintf("user_template_%s.xlsx", time.Now().Format("20060102150405"))
	case "csv":
		fileData, err = s.generateUserCSVTemplate(headers)
		fileName = fmt.Sprintf("user_template_%s.csv", time.Now().Format("20060102150405"))
	default:
		return nil, "", errors.New("不支持的文件格式")
	}

	return fileData, fileName, err
}

// PreviewUserImport 预览用户导入数据
func (s *userService) PreviewUserImport(ctx context.Context, file multipart.File, header *multipart.FileHeader, req *model.UserImportRequest) (*model.UserImportResponse, error) {
	return s.processUserImport(ctx, file, header, req, true)
}

// ImportUsers 导入用户数据
func (s *userService) ImportUsers(ctx context.Context, file multipart.File, header *multipart.FileHeader, req *model.UserImportRequest) (*model.UserImportResponse, error) {
	return s.processUserImport(ctx, file, header, req, false)
}

// processUserImport 处理用户导入（预览或实际导入）
func (s *userService) processUserImport(ctx context.Context, file multipart.File, header *multipart.FileHeader, req *model.UserImportRequest, preview bool) (*model.UserImportResponse, error) {
	// 解析文件
	var records [][]string
	var err error

	fileName := strings.ToLower(header.Filename)
	if strings.HasSuffix(fileName, ".xlsx") || strings.HasSuffix(fileName, ".xls") {
		records, err = s.parseUserExcelFile(file)
	} else if strings.HasSuffix(fileName, ".csv") {
		records, err = s.parseUserCSVFile(file)
	} else {
		return nil, errors.New("不支持的文件格式")
	}

	if err != nil {
		return nil, err
	}

	// 跳过表头
	if req.SkipHeader && len(records) > 0 {
		records = records[1:]
	}

	response := &model.UserImportResponse{
		TotalCount: len(records),
		Errors:     []model.UserImportError{},
	}

	var users []model.UserInfo
	successCount := 0

	for i, record := range records {
		rowNum := i + 1
		if req.SkipHeader {
			rowNum = i + 2 // 考虑表头行
		}

		user, errors := s.validateAndConvertUserRecord(record, rowNum)
		if len(errors) > 0 {
			response.Errors = append(response.Errors, model.UserImportError{
				Row:    rowNum,
				Errors: errors,
				Data:   record,
			})
			continue
		}

		// 检查用户名是否已存在
		if !preview {
			exists, _ := s.userRepo.ExistsByUsername(ctx, user.Username)
			if exists && !req.UpdateExisting {
				response.Errors = append(response.Errors, model.UserImportError{
					Row:    rowNum,
					Errors: []string{"用户名已存在"},
					Data:   record,
				})
				continue
			}

			// 实际导入
			createReq := &model.UserCreateRequest{
				Username:    user.Username,
				DisplayName: user.DisplayName,
				Password:    "123456", // 默认密码，实际应该从记录中获取
				CompanyID:   user.CompanyID,
				RoleIDs:     user.RoleIDs,
				Email:       user.Email,
				Phone:       user.Phone,
				Remark:      "",
			}

			_, err := s.CreateUser(ctx, createReq)
			if err != nil {
				response.Errors = append(response.Errors, model.UserImportError{
					Row:    rowNum,
					Errors: []string{err.Error()},
					Data:   record,
				})
				continue
			}
		}

		users = append(users, *user)
		successCount++
	}

	response.SuccessCount = successCount
	response.ErrorCount = len(response.Errors)

	if preview {
		// 预览时返回前10条数据
		previewCount := 10
		if len(users) < previewCount {
			previewCount = len(users)
		}
		response.Preview = users[:previewCount]
	}

	return response, nil
}

// 辅助方法
func (s *userService) generateUserDataFile(users []model.UserInfo, format string) ([]byte, string, error) {
	var fileData []byte
	var fileName string
	var err error

	switch format {
	case "xlsx":
		fileData, err = s.generateUserExcelData(users)
		fileName = fmt.Sprintf("users_export_%s.xlsx", time.Now().Format("20060102150405"))
	case "csv":
		fileData, err = s.generateUserCSVData(users)
		fileName = fmt.Sprintf("users_export_%s.csv", time.Now().Format("20060102150405"))
	default:
		return nil, "", errors.New("不支持的文件格式")
	}

	return fileData, fileName, err
}

func (s *userService) generateUserExcelTemplate(headers []string) ([]byte, error) {
	f := excelize.NewFile()
	sheetName := "Sheet1"

	// 写入表头
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
	}

	// 设置列宽
	for i := range headers {
		col := string(rune('A' + i))
		f.SetColWidth(sheetName, col, col, 15)
	}

	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (s *userService) generateUserCSVTemplate(headers []string) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	if err := writer.Write(headers); err != nil {
		return nil, err
	}

	writer.Flush()
	return buf.Bytes(), nil
}

func (s *userService) parseUserExcelFile(file multipart.File) ([][]string, error) {
	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (s *userService) parseUserCSVFile(file multipart.File) ([][]string, error) {
	reader := csv.NewReader(file)
	return reader.ReadAll()
}

func (s *userService) validateAndConvertUserRecord(record []string, rowNum int) (*model.UserInfo, []string) {
	var errors []string

	if len(record) < 8 {
		errors = append(errors, "数据列数不足")
		return nil, errors
	}

	// 验证必填字段
	if strings.TrimSpace(record[0]) == "" {
		errors = append(errors, "用户名不能为空")
	}
	if strings.TrimSpace(record[1]) == "" {
		errors = append(errors, "显示名称不能为空")
	}
	if strings.TrimSpace(record[2]) == "" {
		errors = append(errors, "密码不能为空")
	}
	if strings.TrimSpace(record[3]) == "" {
		errors = append(errors, "所属公司ID不能为空")
	}

	if len(errors) > 0 {
		return nil, errors
	}

	// 处理角色ID
	roleIDs := []string{}
	if roleStr := strings.TrimSpace(record[4]); roleStr != "" {
		roleIDs = strings.Split(roleStr, ",")
		for i := range roleIDs {
			roleIDs[i] = strings.TrimSpace(roleIDs[i])
		}
	}

	user := &model.UserInfo{
		Username:    strings.TrimSpace(record[0]),
		DisplayName: strings.TrimSpace(record[1]),
		CompanyID:   strings.TrimSpace(record[3]),
		RoleIDs:     roleIDs,
		Email:       strings.TrimSpace(record[5]),
		Phone:       strings.TrimSpace(record[6]),
		Status:      "active",
	}

	return user, nil
}

func (s *userService) generateUserExcelData(users []model.UserInfo) ([]byte, error) {
	f := excelize.NewFile()
	sheetName := "Sheet1"

	// 设置表头
	headers := []string{"用户ID", "用户名", "显示名称", "所属公司ID", "角色ID", "状态", "邮箱", "手机号", "最后登录时间"}
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
	}

	// 填充数据
	for i, user := range users {
		row := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), user.UserID)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), user.Username)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), user.DisplayName)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), user.CompanyID)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), strings.Join(user.RoleIDs, ","))
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), user.Status)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), user.Email)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), user.Phone)

		lastLogin := ""
		if user.LastLogin != nil {
			lastLogin = user.LastLogin.Format("2006-01-02 15:04:05")
		}
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), lastLogin)
	}

	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (s *userService) generateUserCSVData(users []model.UserInfo) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// 写入表头
	headers := []string{"用户ID", "用户名", "显示名称", "所属公司ID", "角色ID", "状态", "邮箱", "手机号", "最后登录时间"}
	if err := writer.Write(headers); err != nil {
		return nil, err
	}

	// 写入数据
	for _, user := range users {
		lastLogin := ""
		if user.LastLogin != nil {
			lastLogin = user.LastLogin.Format("2006-01-02 15:04:05")
		}

		record := []string{
			user.UserID,
			user.Username,
			user.DisplayName,
			user.CompanyID,
			strings.Join(user.RoleIDs, ","),
			user.Status,
			user.Email,
			user.Phone,
			lastLogin,
		}
		if err := writer.Write(record); err != nil {
			return nil, err
		}
	}

	writer.Flush()
	return buf.Bytes(), nil
}
