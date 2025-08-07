package service

import (
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"mime/multipart"
	"strconv"
	"strings"
	"time"

	"YufungProject/internal/model"
	"YufungProject/internal/repository"
	"YufungProject/pkg/logger"
	"YufungProject/pkg/utils"

	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// CompanyService 公司服务接口
type CompanyService interface {
	// 创建公司
	CreateCompany(ctx context.Context, req *model.CreateCompanyRequest) (*model.CompanyInfo, error)
	// 获取公司详情
	GetCompanyByID(ctx context.Context, companyID string) (*model.CompanyInfo, error)
	// 获取公司列表
	GetCompanyList(ctx context.Context, req *model.CompanyQueryRequest) (*model.CompanyListResponse, error)
	// 更新公司
	UpdateCompany(ctx context.Context, companyID string, req *model.UpdateCompanyRequest) (*model.CompanyInfo, error)
	// 删除公司
	DeleteCompany(ctx context.Context, companyID string) error
	// 获取公司统计
	GetCompanyStats(ctx context.Context) (*model.CompanyStatsResponse, error)

	// 导入导出功能
	// 导出公司数据
	ExportCompany(ctx context.Context, req *model.CompanyExportRequest) (*model.CompanyExportResponse, error)
	// 生成导入模板
	GenerateTemplate(ctx context.Context, format string) ([]byte, string, error)
	// 预览导入数据
	PreviewImport(ctx context.Context, file multipart.File, header *multipart.FileHeader, req *model.CompanyImportRequest) (*model.CompanyImportResponse, error)
	// 导入公司数据
	ImportCompany(ctx context.Context, file multipart.File, header *multipart.FileHeader, req *model.CompanyImportRequest) (*model.CompanyImportResponse, error)
}

// companyService 公司服务实现
type companyService struct {
	companyRepo repository.CompanyRepository
	userRepo    repository.UserRepository
}

// NewCompanyService 创建公司服务实例
func NewCompanyService(companyRepo repository.CompanyRepository, userRepo repository.UserRepository) CompanyService {
	return &companyService{
		companyRepo: companyRepo,
		userRepo:    userRepo,
	}
}

// CreateCompany 创建公司
func (s *companyService) CreateCompany(ctx context.Context, req *model.CreateCompanyRequest) (*model.CompanyInfo, error) {
	// 检查公司名称是否已存在
	exists, err := s.companyRepo.ExistsCompanyName(ctx, req.CompanyName, "")
	if err != nil {
		return nil, fmt.Errorf("检查公司名称失败: %v", err)
	}
	if exists {
		return nil, fmt.Errorf("公司名称已存在")
	}

	// 解析日期
	validStartDate, err := time.Parse("2006-01-02", req.ValidStartDate)
	if err != nil {
		return nil, fmt.Errorf("有效期开始日期格式错误")
	}

	validEndDate, err := time.Parse("2006-01-02", req.ValidEndDate)
	if err != nil {
		return nil, fmt.Errorf("有效期结束日期格式错误")
	}

	// 验证日期合理性
	if validEndDate.Before(validStartDate) {
		return nil, fmt.Errorf("有效期结束日期不能早于开始日期")
	}

	// 创建公司模型
	company := &model.Company{
		CompanyName:   req.CompanyName,
		CompanyCode:   req.CompanyCode,
		ContactPerson: req.ContactPerson,
		TelNo:         req.TelNo,
		Mobile:        req.Mobile,
		ContactPhone:  req.ContactPhone,
		Email:         req.Email,
		// 中文地址信息
		AddressCNProvince: req.AddressCNProvince,
		AddressCNCity:     req.AddressCNCity,
		AddressCNDistrict: req.AddressCNDistrict,
		AddressCNDetail:   req.AddressCNDetail,
		// 英文地址信息
		AddressENProvince: req.AddressENProvince,
		AddressENCity:     req.AddressENCity,
		AddressENDistrict: req.AddressENDistrict,
		AddressENDetail:   req.AddressENDetail,
		Address:           req.Address, // 保留兼容
		BrokerCode:        req.BrokerCode,
		Link:              req.Link,
		Username:          req.Username,
		ValidStartDate:    validStartDate,
		ValidEndDate:      validEndDate,
		UserQuota:         req.UserQuota,
		CurrentUserCount:  0,
		Status:            "active",
		Remark:            req.Remark,
	}

	// 保存到数据库
	if err := s.companyRepo.CreateCompany(ctx, company); err != nil {
		return nil, fmt.Errorf("创建公司失败: %v", err)
	}

	// 记录业务日志
	logger.BusinessLog("公司管理", "创建公司", company.CompanyID, fmt.Sprintf("新建公司: %s", company.CompanyName))

	// 转换为响应格式
	return s.convertToCompanyInfo(company), nil
}

// GetCompanyByID 获取公司详情
func (s *companyService) GetCompanyByID(ctx context.Context, companyID string) (*model.CompanyInfo, error) {
	company, err := s.companyRepo.GetCompanyByID(ctx, companyID)
	if err != nil {
		return nil, fmt.Errorf("查询公司失败: %v", err)
	}
	if company == nil {
		return nil, fmt.Errorf("公司不存在")
	}

	// 获取当前用户数量
	userCount, _ := s.companyRepo.GetCompanyUserStats(ctx, companyID)
	company.CurrentUserCount = int(userCount)

	return s.convertToCompanyInfo(company), nil
}

// GetCompanyList 获取公司列表
func (s *companyService) GetCompanyList(ctx context.Context, req *model.CompanyQueryRequest) (*model.CompanyListResponse, error) {
	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	// 查询公司列表
	companies, total, err := s.companyRepo.GetCompanyList(ctx, req.Page, req.PageSize, req.Status)
	if err != nil {
		return nil, fmt.Errorf("查询公司列表失败: %v", err)
	}

	// 转换为响应格式
	companyInfos := make([]model.CompanyInfo, 0, len(companies))
	for _, company := range companies {
		companyInfos = append(companyInfos, *s.convertToCompanyInfo(company))
	}

	// 计算总页数
	totalPages := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))

	return &model.CompanyListResponse{
		Companies:  companyInfos,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// UpdateCompany 更新公司
func (s *companyService) UpdateCompany(ctx context.Context, companyID string, req *model.UpdateCompanyRequest) (*model.CompanyInfo, error) {
	// 检查公司是否存在
	existingCompany, err := s.companyRepo.GetCompanyByID(ctx, companyID)
	if err != nil {
		return nil, fmt.Errorf("查询公司失败: %v", err)
	}
	if existingCompany == nil {
		return nil, fmt.Errorf("公司不存在")
	}

	// 构建更新字段
	updates := bson.M{}

	if req.CompanyName != "" && req.CompanyName != existingCompany.CompanyName {
		// 检查新名称是否已存在
		exists, err := s.companyRepo.ExistsCompanyName(ctx, req.CompanyName, companyID)
		if err != nil {
			return nil, fmt.Errorf("检查公司名称失败: %v", err)
		}
		if exists {
			return nil, fmt.Errorf("公司名称已存在")
		}
		updates["company_name"] = req.CompanyName
	}

	// 基本信息字段
	if req.CompanyCode != "" {
		updates["company_code"] = req.CompanyCode
	}
	if req.ContactPerson != "" {
		updates["contact_person"] = req.ContactPerson
	}
	if req.TelNo != "" {
		updates["tel_no"] = req.TelNo
	}
	if req.Mobile != "" {
		updates["mobile"] = req.Mobile
	}
	if req.ContactPhone != "" {
		updates["contact_phone"] = req.ContactPhone
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}

	// 中文地址字段
	if req.AddressCNProvince != "" {
		updates["address_cn_province"] = req.AddressCNProvince
	}
	if req.AddressCNCity != "" {
		updates["address_cn_city"] = req.AddressCNCity
	}
	if req.AddressCNDistrict != "" {
		updates["address_cn_district"] = req.AddressCNDistrict
	}
	if req.AddressCNDetail != "" {
		updates["address_cn_detail"] = req.AddressCNDetail
	}

	// 英文地址字段
	if req.AddressENProvince != "" {
		updates["address_en_province"] = req.AddressENProvince
	}
	if req.AddressENCity != "" {
		updates["address_en_city"] = req.AddressENCity
	}
	if req.AddressENDistrict != "" {
		updates["address_en_district"] = req.AddressENDistrict
	}
	if req.AddressENDetail != "" {
		updates["address_en_detail"] = req.AddressENDetail
	}

	// 兼容性地址字段
	if req.Address != "" {
		updates["address"] = req.Address
	}

	// 业务信息字段
	if req.BrokerCode != "" {
		updates["broker_code"] = req.BrokerCode
	}
	if req.Link != "" {
		updates["link"] = req.Link
	}
	if req.Username != "" {
		updates["username"] = req.Username
	}

	// 验证日期合理性
	if req.ValidStartDate != "" {
		validStartDate, err := time.Parse("2006-01-02", req.ValidStartDate)
		if err != nil {
			return nil, fmt.Errorf("有效期开始日期格式错误")
		}
		updates["valid_start_date"] = validStartDate
	}

	if req.ValidEndDate != "" {
		validEndDate, err := time.Parse("2006-01-02", req.ValidEndDate)
		if err != nil {
			return nil, fmt.Errorf("有效期结束日期格式错误")
		}
		updates["valid_end_date"] = validEndDate
	}

	if req.UserQuota > 0 {
		updates["user_quota"] = req.UserQuota
	}

	if req.Status != "" {
		updates["status"] = req.Status
	}

	if req.Remark != "" {
		updates["remark"] = req.Remark
	}

	// 执行更新
	if len(updates) > 0 {
		if err := s.companyRepo.UpdateCompany(ctx, companyID, updates); err != nil {
			return nil, fmt.Errorf("更新公司失败: %v", err)
		}

		// 记录业务日志
		logger.BusinessLog("公司管理", "更新公司", companyID, "公司信息更新")
	}

	// 返回更新后的公司信息
	return s.GetCompanyByID(ctx, companyID)
}

// DeleteCompany 删除公司
func (s *companyService) DeleteCompany(ctx context.Context, companyID string) error {
	// 检查公司是否存在
	company, err := s.companyRepo.GetCompanyByID(ctx, companyID)
	if err != nil {
		return fmt.Errorf("查询公司失败: %v", err)
	}
	if company == nil {
		return fmt.Errorf("公司不存在")
	}

	// 检查是否有关联用户
	userCount, err := s.companyRepo.GetCompanyUserStats(ctx, companyID)
	if err != nil {
		return fmt.Errorf("检查公司用户失败: %v", err)
	}
	if userCount > 0 {
		return fmt.Errorf("公司下还有 %d 个用户，无法删除", userCount)
	}

	// 执行删除
	if err := s.companyRepo.DeleteCompany(ctx, companyID); err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("公司不存在")
		}
		return fmt.Errorf("删除公司失败: %v", err)
	}

	// 记录业务日志
	logger.BusinessLog("公司管理", "删除公司", companyID, fmt.Sprintf("删除公司: %s", company.CompanyName))

	return nil
}

// GetCompanyStats 获取公司统计
func (s *companyService) GetCompanyStats(ctx context.Context) (*model.CompanyStatsResponse, error) {
	// 获取所有公司
	companies, total, err := s.companyRepo.GetCompanyList(ctx, 1, 1000, "")
	if err != nil {
		return nil, fmt.Errorf("查询公司统计失败: %v", err)
	}

	stats := &model.CompanyStatsResponse{
		TotalCompanies: total,
	}

	// 统计各状态公司数量
	now := time.Now()
	for _, company := range companies {
		switch company.Status {
		case "active":
			// 检查是否过期
			if company.ValidEndDate.Before(now) {
				stats.ExpiredCompanies++
			} else {
				stats.ActiveCompanies++
			}
		case "inactive":
			// 不计入有效或过期
		}

		stats.TotalUsers += int64(company.CurrentUserCount)
	}

	return stats, nil
}

// convertToCompanyInfo 转换为公司信息响应格式
func (s *companyService) convertToCompanyInfo(company *model.Company) *model.CompanyInfo {
	info := &model.CompanyInfo{
		ID:            company.ID.Hex(),
		CompanyID:     company.CompanyID,
		CompanyName:   company.CompanyName,
		CompanyCode:   company.CompanyCode,
		ContactPerson: company.ContactPerson,
		TelNo:         company.TelNo,
		Mobile:        company.Mobile,
		ContactPhone:  company.ContactPhone,
		Email:         company.Email,
		// 中文地址信息
		AddressCNProvince: company.AddressCNProvince,
		AddressCNCity:     company.AddressCNCity,
		AddressCNDistrict: company.AddressCNDistrict,
		AddressCNDetail:   company.AddressCNDetail,
		// 英文地址信息
		AddressENProvince: company.AddressENProvince,
		AddressENCity:     company.AddressENCity,
		AddressENDistrict: company.AddressENDistrict,
		AddressENDetail:   company.AddressENDetail,
		Address:           company.Address,
		BrokerCode:        company.BrokerCode,
		Link:              company.Link,
		Username:          company.Username,
		ValidStartDate:    company.ValidStartDate.Format("2006-01-02"),
		ValidEndDate:      company.ValidEndDate.Format("2006-01-02"),
		UserQuota:         company.UserQuota,
		CurrentUserCount:  company.CurrentUserCount,
		Status:            company.Status,
		Remark:            company.Remark,
		SubmittedBy:       company.SubmittedBy,
		CreatedAt:         company.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:         company.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	// 设置状态文本
	switch company.Status {
	case "active":
		if time.Now().After(company.ValidEndDate) {
			info.Status = "expired"
			info.StatusText = "已过期"
		} else {
			info.StatusText = "正常"
		}
	case "inactive":
		info.StatusText = "停用"
	default:
		info.StatusText = "未知"
	}

	return info
}

// ExportCompany 导出公司数据
func (s *companyService) ExportCompany(ctx context.Context, req *model.CompanyExportRequest) (*model.CompanyExportResponse, error) {
	var companies []model.CompanyInfo
	var err error

	// 根据导出类型获取数据
	switch req.ExportType {
	case "all":
		// 导出全部数据
		queryReq := &model.CompanyQueryRequest{Page: 1, PageSize: 10000}
		response, err := s.GetCompanyList(ctx, queryReq)
		if err != nil {
			return nil, err
		}
		companies = response.Companies
	case "selected":
		// 导出选中数据
		if len(req.IDs) == 0 {
			return nil, errors.New("请选择要导出的数据")
		}
		for _, id := range req.IDs {
			company, err := s.GetCompanyByID(ctx, id)
			if err != nil {
				continue // 忽略不存在的公司
			}
			companies = append(companies, *company)
		}
	case "filtered":
		// 导出筛选结果
		queryReq := &model.CompanyQueryRequest{
			Page:     1,
			PageSize: 10000,
			Status:   req.Status,
			Keyword:  req.Keyword,
		}
		response, err := s.GetCompanyList(ctx, queryReq)
		if err != nil {
			return nil, err
		}
		companies = response.Companies
	default:
		return nil, errors.New("不支持的导出类型")
	}

	// 生成文件
	var fileData []byte
	var fileName string

	if req.Template {
		// 生成模板文件
		fileData, fileName, err = s.GenerateTemplate(ctx, req.Format)
	} else {
		// 生成数据文件
		fileData, fileName, err = s.generateDataFile(companies, req.Format)
	}

	if err != nil {
		return nil, err
	}

	// 这里简化处理，实际项目中可能需要将文件保存到对象存储
	// 并返回下载链接
	response := &model.CompanyExportResponse{
		FileURL:  "/api/download/" + fileName, // 简化的下载链接
		FileName: fileName,
	}

	// 记录文件大小用于日志
	logger.Infof("生成导出文件成功: %s, 大小: %d bytes", fileName, len(fileData))

	return response, nil
}

// GenerateTemplate 生成导入模板
func (s *companyService) GenerateTemplate(ctx context.Context, format string) ([]byte, string, error) {
	headers := []string{
		"公司名称", "公司代码", "负责人中文名", "负责人英文名", "联络人",
		"固定电话", "移动电话", "邮箱地址", "中文地址省份", "中文地址城市",
		"中文地址区县", "中文地址详细", "英文地址省份", "英文地址城市",
		"英文地址区县", "英文地址详细", "经纪人代码", "相关链接",
		"用户名", "备注信息", "有效期开始", "有效期结束", "用户配额",
	}

	var fileData []byte
	var fileName string
	var err error

	switch format {
	case "xlsx":
		fileData, err = s.generateExcelTemplate(headers)
		fileName = fmt.Sprintf("company_template_%s.xlsx", time.Now().Format("20060102150405"))
	case "csv":
		fileData, err = s.generateCSVTemplate(headers)
		fileName = fmt.Sprintf("company_template_%s.csv", time.Now().Format("20060102150405"))
	default:
		return nil, "", errors.New("不支持的文件格式")
	}

	return fileData, fileName, err
}

// PreviewImport 预览导入数据
func (s *companyService) PreviewImport(ctx context.Context, file multipart.File, header *multipart.FileHeader, req *model.CompanyImportRequest) (*model.CompanyImportResponse, error) {
	return s.processImport(ctx, file, header, req, true)
}

// ImportCompany 导入公司数据
func (s *companyService) ImportCompany(ctx context.Context, file multipart.File, header *multipart.FileHeader, req *model.CompanyImportRequest) (*model.CompanyImportResponse, error) {
	return s.processImport(ctx, file, header, req, false)
}

// processImport 处理导入（预览或实际导入）
func (s *companyService) processImport(ctx context.Context, file multipart.File, header *multipart.FileHeader, req *model.CompanyImportRequest, preview bool) (*model.CompanyImportResponse, error) {
	// 解析文件
	var records [][]string
	var err error

	fileName := strings.ToLower(header.Filename)
	if strings.HasSuffix(fileName, ".xlsx") || strings.HasSuffix(fileName, ".xls") {
		records, err = s.parseExcelFile(file)
	} else if strings.HasSuffix(fileName, ".csv") {
		records, err = s.parseCSVFile(file)
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

	response := &model.CompanyImportResponse{
		TotalCount: len(records),
		Errors:     []model.CompanyImportError{},
	}

	var companies []model.CompanyInfo
	successCount := 0

	for i, record := range records {
		rowNum := i + 1
		if req.SkipHeader {
			rowNum = i + 2 // 考虑表头行
		}

		company, errors := s.validateAndConvertRecord(record, rowNum)
		if len(errors) > 0 {
			response.Errors = append(response.Errors, model.CompanyImportError{
				Row:    rowNum,
				Errors: errors,
				Data:   record,
			})
			continue
		}

		// 检查是否存在
		if !preview && req.UpdateExisting {
			exists, _ := s.companyRepo.ExistsCompanyName(ctx, company.CompanyName, "")
			if exists {
				// 更新现有公司逻辑...
			}
		}

		if !preview {
			// 实际导入
			companyModel := s.convertToCompanyModel(company)
			if err := s.companyRepo.CreateCompany(ctx, companyModel); err != nil {
				response.Errors = append(response.Errors, model.CompanyImportError{
					Row:    rowNum,
					Errors: []string{err.Error()},
					Data:   record,
				})
				continue
			}
		}

		companies = append(companies, *company)
		successCount++
	}

	response.SuccessCount = successCount
	response.ErrorCount = len(response.Errors)

	if preview {
		// 预览时返回前10条数据
		previewCount := 10
		if len(companies) < previewCount {
			previewCount = len(companies)
		}
		response.Preview = companies[:previewCount]
	}

	return response, nil
}

// 辅助方法
func (s *companyService) generateExcelTemplate(headers []string) ([]byte, error) {
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

func (s *companyService) generateCSVTemplate(headers []string) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	if err := writer.Write(headers); err != nil {
		return nil, err
	}

	writer.Flush()
	return buf.Bytes(), nil
}

func (s *companyService) generateDataFile(companies []model.CompanyInfo, format string) ([]byte, string, error) {
	var fileData []byte
	var fileName string
	var err error

	switch format {
	case "xlsx":
		fileData, err = s.generateExcelData(companies)
		fileName = fmt.Sprintf("companies_export_%s.xlsx", time.Now().Format("20060102150405"))
	case "csv":
		fileData, err = s.generateCSVData(companies)
		fileName = fmt.Sprintf("companies_export_%s.csv", time.Now().Format("20060102150405"))
	default:
		return nil, "", errors.New("不支持的文件格式")
	}

	return fileData, fileName, err
}

func (s *companyService) generateExcelData(companies []model.CompanyInfo) ([]byte, error) {
	f := excelize.NewFile()
	sheetName := "Sheet1"

	// 表头
	headers := []string{
		"公司名称", "公司代码", "负责人中文名", "负责人英文名", "联络人",
		"固定电话", "移动电话", "邮箱地址", "中文地址省份", "中文地址城市",
		"中文地址区县", "中文地址详细", "英文地址省份", "英文地址城市",
		"英文地址区县", "英文地址详细", "经纪人代码", "用户配额", "状态",
		"有效期开始", "有效期结束", "创建时间",
	}

	// 写入表头
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
	}

	// 写入数据
	for rowIndex, company := range companies {
		row := rowIndex + 2
		values := []interface{}{
			company.CompanyName, company.CompanyCode, company.ContactPerson,
			company.TelNo, company.Mobile, company.Email, company.AddressCNDetail,
			company.AddressENDetail, company.BrokerCode, company.Link,
			company.Username, company.UserQuota, company.StatusText,
			company.ValidStartDate, company.ValidEndDate, company.CreatedAt,
		}

		for colIndex, value := range values {
			cell := fmt.Sprintf("%c%d", 'A'+colIndex, row)
			f.SetCellValue(sheetName, cell, value)
		}
	}

	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (s *companyService) generateCSVData(companies []model.CompanyInfo) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// 写入表头
	headers := []string{
		"公司名称", "公司代码", "联络人", "固定电话", "移动电话", "邮箱地址",
		"中文地址详细", "英文地址详细", "经纪人代码", "相关链接", "用户名",
		"用户配额", "状态", "有效期开始", "有效期结束", "创建时间",
	}

	if err := writer.Write(headers); err != nil {
		return nil, err
	}

	// 写入数据
	for _, company := range companies {
		record := []string{
			company.CompanyName, company.CompanyCode, company.ContactPerson,
			company.TelNo, company.Mobile, company.Email, company.AddressCNDetail,
			company.AddressENDetail, company.BrokerCode, company.Link,
			company.Username, strconv.Itoa(company.UserQuota), company.StatusText,
			company.ValidStartDate, company.ValidEndDate, company.CreatedAt,
		}

		if err := writer.Write(record); err != nil {
			return nil, err
		}
	}

	writer.Flush()
	return buf.Bytes(), nil
}

func (s *companyService) parseExcelFile(file multipart.File) ([][]string, error) {
	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// 获取第一个工作表
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, errors.New("Excel文件中没有工作表")
	}

	rows, err := f.GetRows(sheets[0])
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (s *companyService) parseCSVFile(file multipart.File) ([][]string, error) {
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil
}

func (s *companyService) validateAndConvertRecord(record []string, rowNum int) (*model.CompanyInfo, []string) {
	var errors []string

	// 确保记录有足够的字段
	for len(record) < 12 {
		record = append(record, "")
	}

	// 验证必填字段
	if strings.TrimSpace(record[0]) == "" {
		errors = append(errors, "公司名称不能为空")
	}

	if strings.TrimSpace(record[5]) == "" {
		errors = append(errors, "邮箱地址不能为空")
	}

	// 验证用户配额
	userQuota := 1
	if strings.TrimSpace(record[11]) != "" {
		if quota, err := strconv.Atoi(strings.TrimSpace(record[11])); err != nil {
			errors = append(errors, "用户配额必须是数字")
		} else if quota < 1 {
			errors = append(errors, "用户配额必须大于0")
		} else {
			userQuota = quota
		}
	}

	// 验证日期格式
	validStartDate := strings.TrimSpace(record[9])
	validEndDate := strings.TrimSpace(record[10])

	if validStartDate == "" {
		validStartDate = time.Now().Format("2006-01-02")
	}
	if validEndDate == "" {
		validEndDate = time.Now().AddDate(1, 0, 0).Format("2006-01-02")
	}

	company := &model.CompanyInfo{
		CompanyName:     strings.TrimSpace(record[0]),
		CompanyCode:     strings.TrimSpace(record[1]),
		ContactPerson:   strings.TrimSpace(record[2]),
		TelNo:           strings.TrimSpace(record[3]),
		Mobile:          strings.TrimSpace(record[4]),
		Email:           strings.TrimSpace(record[5]),
		AddressCNDetail: strings.TrimSpace(record[6]),
		AddressENDetail: strings.TrimSpace(record[7]),
		BrokerCode:      strings.TrimSpace(record[8]),
		Link:            strings.TrimSpace(record[9]),
		Username:        strings.TrimSpace(record[10]),
		Remark:          strings.TrimSpace(record[11]),
		ValidStartDate:  validStartDate,
		ValidEndDate:    validEndDate,
		UserQuota:       userQuota,
		Status:          "active",
		StatusText:      "有效",
	}

	return company, errors
}

func (s *companyService) convertToCompanyModel(info *model.CompanyInfo) *model.Company {
	now := time.Now()
	validStart, _ := time.Parse("2006-01-02", info.ValidStartDate)
	validEnd, _ := time.Parse("2006-01-02", info.ValidEndDate)

	return &model.Company{
		CompanyID:        utils.GenerateCompanyID(),
		CompanyName:      info.CompanyName,
		CompanyCode:      info.CompanyCode,
		ContactPerson:    info.ContactPerson,
		TelNo:            info.TelNo,
		Mobile:           info.Mobile,
		ContactPhone:     info.Mobile, // Assuming ContactPhone is Mobile for now
		Email:            info.Email,
		AddressCNDetail:  info.AddressCNDetail,
		AddressENDetail:  info.AddressENDetail,
		Address:          info.AddressCNDetail, // Use CN detail as default address
		BrokerCode:       info.BrokerCode,
		Link:             info.Link,
		Username:         info.Username,
		ValidStartDate:   validStart,
		ValidEndDate:     validEnd,
		UserQuota:        info.UserQuota,
		CurrentUserCount: 0,
		Status:           info.Status,
		Remark:           info.Remark,
		SubmittedBy:      "", // Will be set by service layer
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}
