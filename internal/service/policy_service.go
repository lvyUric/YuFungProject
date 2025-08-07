package service

import (
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"mime/multipart"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson"

	"YufungProject/internal/model"
	"YufungProject/internal/repository"
	"math"
)

type PolicyService struct {
	policyRepo          *repository.PolicyRepository
	changeRecordService *ChangeRecordService
}

func NewPolicyService(policyRepo *repository.PolicyRepository, changeRecordService *ChangeRecordService) *PolicyService {
	return &PolicyService{
		policyRepo:          policyRepo,
		changeRecordService: changeRecordService,
	}
}

// CreatePolicy 创建保单
func (s *PolicyService) CreatePolicy(ctx context.Context, req *model.PolicyCreateRequest, userID, companyID string) (*model.PolicyResponse, error) {
	// 检查重复保单
	isDuplicate, err := s.policyRepo.CheckDuplicatePolicy(ctx, req.AccountNumber, req.ProposalNumber, companyID, "")
	if err != nil {
		return nil, err
	}
	if isDuplicate {
		return nil, fmt.Errorf("账户号或投保单号已存在")
	}

	// 构建保单模型
	policy := &model.Policy{
		AccountNumber:     req.AccountNumber,
		CustomerNumber:    req.CustomerNumber,
		CustomerNameCN:    req.CustomerNameCN,
		CustomerNameEN:    req.CustomerNameEN,
		ProposalNumber:    req.ProposalNumber,
		PolicyCurrency:    req.PolicyCurrency,
		Partner:           req.Partner,
		ReferralCode:      req.ReferralCode,
		HKManager:         req.HKManager,
		ReferralPM:        req.ReferralPM,
		ReferralBranch:    req.ReferralBranch,
		ReferralSubBranch: req.ReferralSubBranch,
		ReferralDate:      req.ReferralDate,
		IsSurrendered:     req.IsSurrendered,
		PaymentDate:       req.PaymentDate,
		EffectiveDate:     req.EffectiveDate,
		PaymentMethod:     req.PaymentMethod,
		PaymentYears:      req.PaymentYears,
		PaymentPeriods:    req.PaymentPeriods,
		ActualPremium:     req.ActualPremium,
		AUM:               req.AUM,
		PastCoolingPeriod: req.PastCoolingPeriod,
		IsPaidCommission:  req.IsPaidCommission,
		IsEmployee:        req.IsEmployee,
		ReferralRate:      req.ReferralRate,
		ExchangeRate:      req.ExchangeRate,
		ExpectedFee:       req.ExpectedFee,
		PaymentPayDate:    req.PaymentPayDate,
		InsuranceCompany:  req.InsuranceCompany,
		ProductName:       req.ProductName,
		ProductType:       req.ProductType,
		Remark:            req.Remark,
		CompanyID:         companyID,
		CreatedBy:         userID,
		UpdatedBy:         userID,
	}

	// 创建保单
	err = s.policyRepo.CreatePolicy(ctx, policy)
	if err != nil {
		return nil, err
	}

	return &model.PolicyResponse{Policy: policy}, nil
}

// GetPolicyByID 获取保单详情
func (s *PolicyService) GetPolicyByID(ctx context.Context, policyID, companyID string) (*model.PolicyResponse, error) {
	policy, err := s.policyRepo.GetPolicyByID(ctx, policyID)
	if err != nil {
		return nil, err
	}
	if policy == nil {
		return nil, fmt.Errorf("保单不存在")
	}

	// 检查权限：确保保单属于当前公司
	if policy.CompanyID != companyID {
		return nil, fmt.Errorf("无权访问该保单")
	}

	return &model.PolicyResponse{Policy: policy}, nil
}

// UpdatePolicy 更新保单
func (s *PolicyService) UpdatePolicy(ctx context.Context, policyID string, req *model.PolicyUpdateRequest, userID, companyID, ipAddress, userAgent string) (*model.PolicyResponse, error) {
	// 检查保单是否存在
	policy, err := s.policyRepo.GetPolicyByID(ctx, policyID)
	if err != nil {
		return nil, err
	}
	if policy == nil {
		return nil, fmt.Errorf("保单不存在")
	}

	// 检查权限
	if policy.CompanyID != companyID {
		return nil, fmt.Errorf("无权修改该保单")
	}

	// 保存原始数据用于变更记录
	oldPolicy := *policy

	// 构建更新字段
	updates := bson.M{}
	updates["updated_by"] = userID

	// 使用反射来设置非空字段
	reqValue := reflect.ValueOf(req).Elem()
	reqType := reflect.TypeOf(req).Elem()

	for i := 0; i < reqValue.NumField(); i++ {
		field := reqValue.Field(i)
		fieldType := reqType.Field(i)

		// 获取bson标签名
		bsonTag := fieldType.Tag.Get("json")
		if bsonTag == "" || bsonTag == "-" {
			continue
		}

		// 处理不同类型的字段
		switch field.Kind() {
		case reflect.String:
			if field.String() != "" {
				updates[bsonTag] = field.String()
			}
		case reflect.Ptr:
			if !field.IsNil() {
				switch field.Elem().Kind() {
				case reflect.Bool:
					updates[bsonTag] = field.Elem().Bool()
				case reflect.Int:
					updates[bsonTag] = field.Elem().Int()
				case reflect.Float64:
					updates[bsonTag] = field.Elem().Float()
				case reflect.Struct:
					// 处理时间类型
					if field.Elem().Type().String() == "time.Time" {
						updates[bsonTag] = field.Interface()
					}
				}
			}
		}
	}

	// 执行更新
	err = s.policyRepo.UpdatePolicy(ctx, policyID, updates)
	if err != nil {
		return nil, err
	}

	// 获取更新后的保单数据
	updatedPolicy, err := s.policyRepo.GetPolicyByID(ctx, policyID)
	if err != nil {
		return nil, err
	}

	// 记录变更日志（异步处理，不影响主流程）
	go func() {
		if s.changeRecordService != nil {
			// 创建独立的上下文以避免原上下文被取消
			changeCtx := context.Background()
			err := s.changeRecordService.RecordChange(
				changeCtx,
				"policies",
				policyID,
				userID,
				companyID,
				"update",
				&oldPolicy,
				updatedPolicy,
				"",
				ipAddress,
				userAgent,
			)
			if err != nil {
				// 记录错误但不影响主流程
				fmt.Printf("Failed to record change: %v\n", err)
			}
		}
	}()

	// 返回更新后的保单
	return &model.PolicyResponse{Policy: updatedPolicy}, nil
}

// DeletePolicy 删除保单
func (s *PolicyService) DeletePolicy(ctx context.Context, policyID, companyID string) error {
	// 检查保单是否存在
	policy, err := s.policyRepo.GetPolicyByID(ctx, policyID)
	if err != nil {
		return err
	}
	if policy == nil {
		return fmt.Errorf("保单不存在")
	}

	// 检查权限
	if policy.CompanyID != companyID {
		return fmt.Errorf("无权删除该保单")
	}

	return s.policyRepo.DeletePolicy(ctx, policyID)
}

// ListPolicies 获取保单列表
func (s *PolicyService) ListPolicies(ctx context.Context, req *model.PolicyQueryRequest, companyID string) (*model.PolicyListResponse, error) {
	return s.policyRepo.ListPolicies(ctx, req, companyID)
}

// GetPolicyStatistics 获取保单统计
func (s *PolicyService) GetPolicyStatistics(ctx context.Context, companyID string) (*model.PolicyStatistics, error) {
	return s.policyRepo.GetPolicyStatistics(ctx, companyID)
}

// BatchUpdatePolicyStatus 批量更新保单状态
func (s *PolicyService) BatchUpdatePolicyStatus(ctx context.Context, req *model.BatchUpdatePolicyStatusRequest, userID, companyID string) error {
	// 验证保单ID属于当前公司
	policies, err := s.policyRepo.GetPoliciesByIDs(ctx, req.PolicyIDs)
	if err != nil {
		return err
	}

	for _, policy := range policies {
		if policy.CompanyID != companyID {
			return fmt.Errorf("保单 %s 不属于当前公司", policy.PolicyID)
		}
	}

	// 构建更新字段
	updates := bson.M{
		"updated_by": userID,
	}

	if req.IsSurrendered != nil {
		updates["is_surrendered"] = *req.IsSurrendered
	}
	if req.PastCoolingPeriod != nil {
		updates["past_cooling_period"] = *req.PastCoolingPeriod
	}
	if req.IsPaidCommission != nil {
		updates["is_paid_commission"] = *req.IsPaidCommission
	}

	return s.policyRepo.BatchUpdatePolicyStatus(ctx, req.PolicyIDs, updates)
}

// ImportPolicies 批量导入保单
func (s *PolicyService) ImportPolicies(ctx context.Context, req *model.PolicyImportRequest, userID, companyID string) ([]string, []string, error) {
	var successIDs []string
	var errors []string

	for i, policyReq := range req.Data {
		// 检查重复
		isDuplicate, err := s.policyRepo.CheckDuplicatePolicy(ctx, policyReq.AccountNumber, policyReq.ProposalNumber, companyID, "")
		if err != nil {
			errors = append(errors, fmt.Sprintf("第%d行: 检查重复时出错: %v", i+1, err))
			continue
		}
		if isDuplicate {
			errors = append(errors, fmt.Sprintf("第%d行: 账户号或投保单号已存在", i+1))
			continue
		}

		// 创建保单
		_, err = s.CreatePolicy(ctx, &policyReq, userID, companyID)
		if err != nil {
			errors = append(errors, fmt.Sprintf("第%d行: %v", i+1, err))
			continue
		}

		successIDs = append(successIDs, fmt.Sprintf("第%d行", i+1))
	}

	return successIDs, errors, nil
}

// ExportPolicies 导出保单
func (s *PolicyService) ExportPolicies(ctx context.Context, req *model.PolicyExportRequest, companyID string) ([]model.Policy, error) {
	if len(req.PolicyIDs) > 0 {
		// 导出指定保单
		policies, err := s.policyRepo.GetPoliciesByIDs(ctx, req.PolicyIDs)
		if err != nil {
			return nil, err
		}

		// 检查权限
		for _, policy := range policies {
			if policy.CompanyID != companyID {
				return nil, fmt.Errorf("保单 %s 不属于当前公司", policy.PolicyID)
			}
		}

		return policies, nil
	} else {
		// 导出全部保单
		queryReq := &model.PolicyQueryRequest{
			Page:     1,
			PageSize: 10000, // 设置一个较大的值
		}

		response, err := s.policyRepo.ListPolicies(ctx, queryReq, companyID)
		if err != nil {
			return nil, err
		}

		var policies []model.Policy
		for _, policyResp := range response.List {
			policies = append(policies, *policyResp.Policy)
		}

		return policies, nil
	}
}

// GeneratePolicyTemplate 生成保单导入模板
func (s *PolicyService) GeneratePolicyTemplate(ctx context.Context, format string) ([]byte, string, error) {
	headers := []string{
		"序号", "账户号", "客户号", "客户中文名", "客户英文名", "投保单号",
		"保单币种（USD/HKD/CNY）", "合作伙伴", "转介编号", "港分客户经理", "转介理财经理",
		"转介分行", "转介支行", "转介日期", "签单后是否退保", "缴费日期", "生效日期",
		"缴费方式（期缴、趸缴、预缴）", "缴费年期", "期缴期数", "实际缴纳保费", "AUM",
		"是否已过冷静期", "是否支付佣金", "转介费率", "汇率", "预计转介费", "支付日期",
		"是否员工", "承保公司", "保险产品名称", "产品类型", "备注说明",
	}

	var fileData []byte
	var fileName string
	var err error

	switch format {
	case "xlsx":
		fileData, err = s.generatePolicyExcelTemplate(headers)
		fileName = fmt.Sprintf("policy_template_%s.xlsx", time.Now().Format("20060102150405"))
	case "csv":
		fileData, err = s.generatePolicyCSVTemplate(headers)
		fileName = fmt.Sprintf("policy_template_%s.csv", time.Now().Format("20060102150405"))
	default:
		return nil, "", errors.New("不支持的文件格式")
	}

	return fileData, fileName, err
}

// ExportPoliciesToFile 导出保单为文件
func (s *PolicyService) ExportPoliciesToFile(ctx context.Context, req *model.PolicyExportRequest, companyID, format string) ([]byte, string, error) {
	// 获取保单数据
	policies, err := s.ExportPolicies(ctx, req, companyID)
	if err != nil {
		return nil, "", err
	}

	var fileData []byte
	var fileName string

	switch format {
	case "xlsx":
		fileData, err = s.generatePolicyExcelData(policies)
		fileName = fmt.Sprintf("policies_export_%s.xlsx", time.Now().Format("20060102150405"))
	case "csv":
		fileData, err = s.generatePolicyCSVData(policies)
		fileName = fmt.Sprintf("policies_export_%s.csv", time.Now().Format("20060102150405"))
	default:
		return nil, "", errors.New("不支持的文件格式")
	}

	return fileData, fileName, err
}

// PreviewPolicyImport 预览保单导入数据
func (s *PolicyService) PreviewPolicyImport(ctx context.Context, file multipart.File, header *multipart.FileHeader, req *model.PolicyImportFileRequest) (*model.PolicyImportResponse, error) {
	return s.processPolicyImport(ctx, file, header, req, true)
}

// ImportPoliciesFromFile 从文件导入保单
func (s *PolicyService) ImportPoliciesFromFile(ctx context.Context, file multipart.File, header *multipart.FileHeader, req *model.PolicyImportFileRequest, userID, companyID string) (*model.PolicyImportResponse, error) {
	req.UserID = userID
	req.CompanyID = companyID
	return s.processPolicyImport(ctx, file, header, req, false)
}

// processPolicyImport 处理保单导入（预览或实际导入）
func (s *PolicyService) processPolicyImport(ctx context.Context, file multipart.File, header *multipart.FileHeader, req *model.PolicyImportFileRequest, preview bool) (*model.PolicyImportResponse, error) {
	// 解析文件
	var records [][]string
	var err error

	fileName := strings.ToLower(header.Filename)
	if strings.HasSuffix(fileName, ".xlsx") || strings.HasSuffix(fileName, ".xls") {
		records, err = s.parsePolicyExcelFile(file)
	} else if strings.HasSuffix(fileName, ".csv") {
		records, err = s.parsePolicyCSVFile(file)
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

	response := &model.PolicyImportResponse{
		TotalCount: len(records),
		Errors:     []model.PolicyImportError{},
	}

	var policies []model.PolicyCreateRequest
	successCount := 0

	for i, record := range records {
		rowNum := i + 1
		if req.SkipHeader {
			rowNum = i + 2 // 考虑表头行
		}

		policy, errors := s.validateAndConvertPolicyRecord(record, rowNum)
		if len(errors) > 0 {
			response.Errors = append(response.Errors, model.PolicyImportError{
				Row:    rowNum,
				Errors: errors,
				Data:   record,
			})
			continue
		}

		if !preview {
			// 实际导入 - 检查重复
			isDuplicate, err := s.policyRepo.CheckDuplicatePolicy(ctx, policy.AccountNumber, policy.ProposalNumber, req.CompanyID, "")
			if err != nil {
				response.Errors = append(response.Errors, model.PolicyImportError{
					Row:    rowNum,
					Errors: []string{fmt.Sprintf("检查重复时出错: %v", err)},
					Data:   record,
				})
				continue
			}
			if isDuplicate && !req.UpdateExisting {
				response.Errors = append(response.Errors, model.PolicyImportError{
					Row:    rowNum,
					Errors: []string{"账户号或投保单号已存在"},
					Data:   record,
				})
				continue
			}

			// 创建保单
			_, err = s.CreatePolicy(ctx, policy, req.UserID, req.CompanyID)
			if err != nil {
				response.Errors = append(response.Errors, model.PolicyImportError{
					Row:    rowNum,
					Errors: []string{err.Error()},
					Data:   record,
				})
				continue
			}
		}

		policies = append(policies, *policy)
		successCount++
	}

	response.SuccessCount = successCount
	response.ErrorCount = len(response.Errors)

	if preview {
		// 预览时返回前10条数据
		previewCount := 10
		if len(policies) > previewCount {
			policies = policies[:previewCount]
		}
		response.Preview = policies
	}

	return response, nil
}

// 辅助方法实现

func (s *PolicyService) generatePolicyExcelTemplate(headers []string) ([]byte, error) {
	f := excelize.NewFile()
	sheetName := "Sheet1"

	// 写入表头
	for i, header := range headers {
		cell := s.getExcelColumnName(i) + "1"
		f.SetCellValue(sheetName, cell, header)
	}

	// 设置列宽
	for i := range headers {
		col := s.getExcelColumnName(i)
		f.SetColWidth(sheetName, col, col, 15)
	}

	// 添加示例数据行
	exampleData := []string{
		"1", "ACC001", "CUST001", "张三", "Zhang San", "PROP001",
		"USD", "合作伙伴A", "REF001", "经理A", "理财经理A",
		"分行A", "支行A", "2024-01-15", "否", "2024-01-20", "2024-02-01",
		"期缴", "10", "12", "10000.00", "50000.00",
		"是", "是", "2.50", "6.9000", "2500.00", "2024-02-15",
		"否", "保险公司A", "产品A", "寿险", "备注信息",
	}

	for i, data := range exampleData {
		if i < len(headers) {
			cell := s.getExcelColumnName(i) + "2"
			f.SetCellValue(sheetName, cell, data)
		}
	}

	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (s *PolicyService) generatePolicyCSVTemplate(headers []string) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// 写入表头
	if err := writer.Write(headers); err != nil {
		return nil, err
	}

	// 写入示例数据
	exampleData := []string{
		"1", "ACC001", "CUST001", "张三", "Zhang San", "PROP001",
		"USD", "合作伙伴A", "REF001", "经理A", "理财经理A",
		"分行A", "支行A", "2024-01-15", "否", "2024-01-20", "2024-02-01",
		"期缴", "10", "12", "10000.00", "50000.00",
		"是", "是", "2.50", "6.9000", "2500.00", "2024-02-15",
		"否", "保险公司A", "产品A", "寿险", "备注信息",
	}

	if err := writer.Write(exampleData); err != nil {
		return nil, err
	}

	writer.Flush()
	return buf.Bytes(), nil
}

func (s *PolicyService) generatePolicyExcelData(policies []model.Policy) ([]byte, error) {
	f := excelize.NewFile()
	sheetName := "Sheet1"

	// 表头
	headers := []string{
		"序号", "账户号", "客户号", "客户中文名", "客户英文名", "投保单号",
		"保单币种（USD/HKD/CNY）", "合作伙伴", "转介编号", "港分客户经理", "转介理财经理",
		"转介分行", "转介支行", "转介日期", "签单后是否退保", "缴费日期", "生效日期",
		"缴费方式（期缴、趸缴、预缴）", "缴费年期", "期缴期数", "实际缴纳保费", "AUM",
		"是否已过冷静期", "是否支付佣金", "转介费率", "汇率", "预计转介费", "支付日期",
		"是否员工", "承保公司", "保险产品名称", "产品类型", "备注说明",
	}

	// 写入表头
	for i, header := range headers {
		cell := s.getExcelColumnName(i) + "1"
		f.SetCellValue(sheetName, cell, header)
	}

	// 写入数据
	for i, policy := range policies {
		row := i + 2
		data := []interface{}{
			policy.SerialNumber,
			policy.AccountNumber,
			policy.CustomerNumber,
			policy.CustomerNameCN,
			policy.CustomerNameEN,
			policy.ProposalNumber,
			policy.PolicyCurrency,
			policy.Partner,
			policy.ReferralCode,
			policy.HKManager,
			policy.ReferralPM,
			policy.ReferralBranch,
			policy.ReferralSubBranch,
			s.formatDate(policy.ReferralDate),
			s.formatBool(policy.IsSurrendered),
			s.formatDate(policy.PaymentDate),
			s.formatDate(policy.EffectiveDate),
			policy.PaymentMethod,
			policy.PaymentYears,
			policy.PaymentPeriods,
			policy.ActualPremium,
			policy.AUM,
			s.formatBool(policy.PastCoolingPeriod),
			s.formatBool(policy.IsPaidCommission),
			policy.ReferralRate,
			policy.ExchangeRate,
			policy.ExpectedFee,
			s.formatDate(policy.PaymentPayDate),
			s.formatBool(policy.IsEmployee),
			policy.InsuranceCompany,
			policy.ProductName,
			policy.ProductType,
			policy.Remark,
		}

		for j, value := range data {
			cell := s.getExcelColumnName(j) + fmt.Sprintf("%d", row)
			f.SetCellValue(sheetName, cell, value)
		}
	}

	// 设置列宽
	for i := range headers {
		col := s.getExcelColumnName(i)
		f.SetColWidth(sheetName, col, col, 15)
	}

	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (s *PolicyService) generatePolicyCSVData(policies []model.Policy) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// 表头
	headers := []string{
		"序号", "账户号", "客户号", "客户中文名", "客户英文名", "投保单号",
		"保单币种（USD/HKD/CNY）", "合作伙伴", "转介编号", "港分客户经理", "转介理财经理",
		"转介分行", "转介支行", "转介日期", "签单后是否退保", "缴费日期", "生效日期",
		"缴费方式（期缴、趸缴、预缴）", "缴费年期", "期缴期数", "实际缴纳保费", "AUM",
		"是否已过冷静期", "是否支付佣金", "转介费率", "汇率", "预计转介费", "支付日期",
		"是否员工", "承保公司", "保险产品名称", "产品类型", "备注说明",
	}

	if err := writer.Write(headers); err != nil {
		return nil, err
	}

	// 写入数据
	for _, policy := range policies {
		record := []string{
			fmt.Sprintf("%d", policy.SerialNumber),
			policy.AccountNumber,
			policy.CustomerNumber,
			policy.CustomerNameCN,
			policy.CustomerNameEN,
			policy.ProposalNumber,
			policy.PolicyCurrency,
			policy.Partner,
			policy.ReferralCode,
			policy.HKManager,
			policy.ReferralPM,
			policy.ReferralBranch,
			policy.ReferralSubBranch,
			s.formatDate(policy.ReferralDate),
			s.formatBool(policy.IsSurrendered),
			s.formatDate(policy.PaymentDate),
			s.formatDate(policy.EffectiveDate),
			policy.PaymentMethod,
			fmt.Sprintf("%d", policy.PaymentYears),
			fmt.Sprintf("%d", policy.PaymentPeriods),
			fmt.Sprintf("%.2f", policy.ActualPremium),
			fmt.Sprintf("%.2f", policy.AUM),
			s.formatBool(policy.PastCoolingPeriod),
			s.formatBool(policy.IsPaidCommission),
			fmt.Sprintf("%.2f", policy.ReferralRate),
			fmt.Sprintf("%.4f", policy.ExchangeRate),
			fmt.Sprintf("%.2f", policy.ExpectedFee),
			s.formatDate(policy.PaymentPayDate),
			s.formatBool(policy.IsEmployee),
			policy.InsuranceCompany,
			policy.ProductName,
			policy.ProductType,
			policy.Remark,
		}

		if err := writer.Write(record); err != nil {
			return nil, err
		}
	}

	writer.Flush()
	return buf.Bytes(), nil
}

func (s *PolicyService) parsePolicyExcelFile(file multipart.File) ([][]string, error) {
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

func (s *PolicyService) parsePolicyCSVFile(file multipart.File) ([][]string, error) {
	reader := csv.NewReader(file)
	return reader.ReadAll()
}

func (s *PolicyService) validateAndConvertPolicyRecord(record []string, rowNum int) (*model.PolicyCreateRequest, []string) {
	var errors []string

	if len(record) < 33 {
		errors = append(errors, "数据列数不足")
		return nil, errors
	}

	// 验证必填字段
	if strings.TrimSpace(record[5]) == "" { // 投保单号
		errors = append(errors, "投保单号不能为空")
	}

	// 转换数据类型
	paymentYears, err := strconv.Atoi(strings.TrimSpace(record[18]))
	if err != nil && strings.TrimSpace(record[18]) != "" {
		errors = append(errors, "缴费年期格式错误")
		paymentYears = 0
	}

	paymentPeriods, err := strconv.Atoi(strings.TrimSpace(record[19]))
	if err != nil && strings.TrimSpace(record[19]) != "" {
		errors = append(errors, "期缴期数格式错误")
		paymentPeriods = 0
	}

	actualPremium, err := strconv.ParseFloat(strings.TrimSpace(record[20]), 64)
	if err != nil && strings.TrimSpace(record[20]) != "" {
		errors = append(errors, "实际缴纳保费格式错误")
		actualPremium = 0
	}

	aum, err := strconv.ParseFloat(strings.TrimSpace(record[21]), 64)
	if err != nil && strings.TrimSpace(record[21]) != "" {
		errors = append(errors, "AUM格式错误")
		aum = 0
	}

	referralRate, err := strconv.ParseFloat(strings.TrimSpace(record[24]), 64)
	if err != nil && strings.TrimSpace(record[24]) != "" {
		errors = append(errors, "转介费率格式错误")
		referralRate = 0
	}

	exchangeRate, err := strconv.ParseFloat(strings.TrimSpace(record[25]), 64)
	if err != nil && strings.TrimSpace(record[25]) != "" {
		errors = append(errors, "汇率格式错误")
		exchangeRate = 0
	}

	expectedFee, err := strconv.ParseFloat(strings.TrimSpace(record[26]), 64)
	if err != nil && strings.TrimSpace(record[26]) != "" {
		errors = append(errors, "预计转介费格式错误")
		expectedFee = 0
	}

	if len(errors) > 0 {
		return nil, errors
	}

	// 创建保单请求对象
	policy := &model.PolicyCreateRequest{
		AccountNumber:     strings.TrimSpace(record[1]),
		CustomerNumber:    strings.TrimSpace(record[2]),
		CustomerNameCN:    strings.TrimSpace(record[3]),
		CustomerNameEN:    strings.TrimSpace(record[4]),
		ProposalNumber:    strings.TrimSpace(record[5]),
		PolicyCurrency:    strings.TrimSpace(record[6]),
		Partner:           strings.TrimSpace(record[7]),
		ReferralCode:      strings.TrimSpace(record[8]),
		HKManager:         strings.TrimSpace(record[9]),
		ReferralPM:        strings.TrimSpace(record[10]),
		ReferralBranch:    strings.TrimSpace(record[11]),
		ReferralSubBranch: strings.TrimSpace(record[12]),
		PaymentMethod:     strings.TrimSpace(record[17]),
		PaymentYears:      paymentYears,
		PaymentPeriods:    paymentPeriods,
		ActualPremium:     actualPremium,
		AUM:               aum,
		ReferralRate:      referralRate,
		ExchangeRate:      math.Round(exchangeRate*10000) / 10000, // 保留4位小数
		ExpectedFee:       expectedFee,
		IsSurrendered:     strings.TrimSpace(record[14]) == "是" || strings.TrimSpace(record[14]) == "true",
		PastCoolingPeriod: strings.TrimSpace(record[22]) == "是" || strings.TrimSpace(record[22]) == "true",
		IsPaidCommission:  strings.TrimSpace(record[23]) == "是" || strings.TrimSpace(record[23]) == "true",
		IsEmployee:        strings.TrimSpace(record[28]) == "是" || strings.TrimSpace(record[28]) == "true",
		InsuranceCompany:  strings.TrimSpace(record[29]),
		ProductName:       strings.TrimSpace(record[30]),
		ProductType:       strings.TrimSpace(record[31]),
		Remark:            strings.TrimSpace(record[32]),
	}

	// 处理日期字段
	if strings.TrimSpace(record[13]) != "" {
		if date, err := time.Parse("2006-01-02", strings.TrimSpace(record[13])); err == nil {
			policy.ReferralDate = &date
		}
	}
	if strings.TrimSpace(record[15]) != "" {
		if date, err := time.Parse("2006-01-02", strings.TrimSpace(record[15])); err == nil {
			policy.PaymentDate = &date
		}
	}
	if strings.TrimSpace(record[16]) != "" {
		if date, err := time.Parse("2006-01-02", strings.TrimSpace(record[16])); err == nil {
			policy.EffectiveDate = &date
		}
	}
	if strings.TrimSpace(record[27]) != "" {
		if date, err := time.Parse("2006-01-02", strings.TrimSpace(record[27])); err == nil {
			policy.PaymentPayDate = &date
		}
	}

	return policy, nil
}

// 辅助函数
func (s *PolicyService) formatDate(date *time.Time) string {
	if date == nil {
		return ""
	}
	return date.Format("2006-01-02")
}

func (s *PolicyService) formatBool(b bool) string {
	if b {
		return "是"
	}
	return "否"
}

// getExcelColumnName 生成Excel列名 (A, B, ..., Z, AA, AB, ...)
func (s *PolicyService) getExcelColumnName(index int) string {
	result := ""
	for index >= 0 {
		result = string(rune('A'+index%26)) + result
		index = index/26 - 1
	}
	return result
}
