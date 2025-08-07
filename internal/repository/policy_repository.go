package repository

import (
	"context"
	"math"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"YufungProject/internal/model"
	"YufungProject/pkg/utils"
)

const PolicyCollection = "policies"

type PolicyRepository struct {
	db *mongo.Database
}

func NewPolicyRepository(db *mongo.Database) *PolicyRepository {
	return &PolicyRepository{db: db}
}

// CreatePolicy 创建保单
func (r *PolicyRepository) CreatePolicy(ctx context.Context, policy *model.Policy) error {
	collection := r.db.Collection(PolicyCollection)

	// 生成业务主键
	policy.PolicyID = utils.GenerateID("POL")
	policy.CreatedAt = time.Now()
	policy.UpdatedAt = time.Now()

	// 自动生成序号
	serialNumber, err := r.getNextSerialNumber(ctx, policy.CompanyID)
	if err != nil {
		return err
	}
	policy.SerialNumber = serialNumber

	_, err = collection.InsertOne(ctx, policy)
	return err
}

// GetPolicyByID 根据ID获取保单
func (r *PolicyRepository) GetPolicyByID(ctx context.Context, policyID string) (*model.Policy, error) {
	collection := r.db.Collection(PolicyCollection)

	var policy model.Policy
	err := collection.FindOne(ctx, bson.M{"policy_id": policyID}).Decode(&policy)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &policy, nil
}

// UpdatePolicy 更新保单
func (r *PolicyRepository) UpdatePolicy(ctx context.Context, policyID string, updates bson.M) error {
	collection := r.db.Collection(PolicyCollection)

	updates["updated_at"] = time.Now()

	_, err := collection.UpdateOne(
		ctx,
		bson.M{"policy_id": policyID},
		bson.M{"$set": updates},
	)
	return err
}

// DeletePolicy 删除保单
func (r *PolicyRepository) DeletePolicy(ctx context.Context, policyID string) error {
	collection := r.db.Collection(PolicyCollection)

	_, err := collection.DeleteOne(ctx, bson.M{"policy_id": policyID})
	return err
}

// ListPolicies 查询保单列表
func (r *PolicyRepository) ListPolicies(ctx context.Context, req *model.PolicyQueryRequest, companyID string) (*model.PolicyListResponse, error) {
	collection := r.db.Collection(PolicyCollection)

	// 构建查询条件
	filter := bson.M{"company_id": companyID}

	// 添加搜索条件
	if req.AccountNumber != "" {
		filter["account_number"] = bson.M{"$regex": req.AccountNumber, "$options": "i"}
	}
	if req.CustomerNumber != "" {
		filter["customer_number"] = bson.M{"$regex": req.CustomerNumber, "$options": "i"}
	}
	if req.CustomerNameCN != "" {
		filter["customer_name_cn"] = bson.M{"$regex": req.CustomerNameCN, "$options": "i"}
	}
	if req.CustomerNameEN != "" {
		filter["customer_name_en"] = bson.M{"$regex": req.CustomerNameEN, "$options": "i"}
	}
	if req.ProposalNumber != "" {
		filter["proposal_number"] = bson.M{"$regex": req.ProposalNumber, "$options": "i"}
	}
	if req.PolicyCurrency != "" {
		filter["policy_currency"] = req.PolicyCurrency
	}
	if req.Partner != "" {
		filter["partner"] = bson.M{"$regex": req.Partner, "$options": "i"}
	}
	if req.ReferralCode != "" {
		filter["referral_code"] = bson.M{"$regex": req.ReferralCode, "$options": "i"}
	}
	if req.HKManager != "" {
		filter["hk_manager"] = bson.M{"$regex": req.HKManager, "$options": "i"}
	}
	if req.ReferralPM != "" {
		filter["referral_pm"] = bson.M{"$regex": req.ReferralPM, "$options": "i"}
	}
	if req.ReferralBranch != "" {
		filter["referral_branch"] = bson.M{"$regex": req.ReferralBranch, "$options": "i"}
	}
	if req.ReferralSubBranch != "" {
		filter["referral_sub_branch"] = bson.M{"$regex": req.ReferralSubBranch, "$options": "i"}
	}
	if req.PaymentMethod != "" {
		filter["payment_method"] = req.PaymentMethod
	}
	if req.InsuranceCompany != "" {
		filter["insurance_company"] = bson.M{"$regex": req.InsuranceCompany, "$options": "i"}
	}
	if req.ProductName != "" {
		filter["product_name"] = bson.M{"$regex": req.ProductName, "$options": "i"}
	}
	if req.ProductType != "" {
		filter["product_type"] = bson.M{"$regex": req.ProductType, "$options": "i"}
	}

	// 布尔字段筛选
	if req.IsSurrendered != nil {
		filter["is_surrendered"] = *req.IsSurrendered
	}
	if req.PastCoolingPeriod != nil {
		filter["past_cooling_period"] = *req.PastCoolingPeriod
	}
	if req.IsPaidCommission != nil {
		filter["is_paid_commission"] = *req.IsPaidCommission
	}
	if req.IsEmployee != nil {
		filter["is_employee"] = *req.IsEmployee
	}

	// 日期范围筛选
	if req.ReferralDateStart != nil || req.ReferralDateEnd != nil {
		dateFilter := bson.M{}
		if req.ReferralDateStart != nil {
			dateFilter["$gte"] = *req.ReferralDateStart
		}
		if req.ReferralDateEnd != nil {
			dateFilter["$lte"] = *req.ReferralDateEnd
		}
		filter["referral_date"] = dateFilter
	}

	if req.PaymentDateStart != nil || req.PaymentDateEnd != nil {
		dateFilter := bson.M{}
		if req.PaymentDateStart != nil {
			dateFilter["$gte"] = *req.PaymentDateStart
		}
		if req.PaymentDateEnd != nil {
			dateFilter["$lte"] = *req.PaymentDateEnd
		}
		filter["payment_date"] = dateFilter
	}

	if req.EffectiveDateStart != nil || req.EffectiveDateEnd != nil {
		dateFilter := bson.M{}
		if req.EffectiveDateStart != nil {
			dateFilter["$gte"] = *req.EffectiveDateStart
		}
		if req.EffectiveDateEnd != nil {
			dateFilter["$lte"] = *req.EffectiveDateEnd
		}
		filter["effective_date"] = dateFilter
	}

	// 计算总数
	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	// 设置默认分页
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	// 计算跳过数量
	skip := (req.Page - 1) * req.PageSize

	// 构建排序
	sort := bson.D{}
	if req.SortBy != "" {
		sortOrder := 1
		if req.SortOrder == "desc" {
			sortOrder = -1
		}
		sort = append(sort, bson.E{Key: req.SortBy, Value: sortOrder})
	} else {
		sort = append(sort, bson.E{Key: "created_at", Value: -1}) // 默认按创建时间倒序
	}

	// 查询数据
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(req.PageSize)).
		SetSort(sort)

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var policies []model.Policy
	if err = cursor.All(ctx, &policies); err != nil {
		return nil, err
	}

	// 转换为响应格式
	var policyResponses []model.PolicyResponse
	for _, policy := range policies {
		policyResponses = append(policyResponses, model.PolicyResponse{Policy: &policy})
	}

	// 计算总页数
	totalPages := int(math.Ceil(float64(total) / float64(req.PageSize)))

	return &model.PolicyListResponse{
		List:       policyResponses,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// CheckDuplicatePolicy 检查重复保单
func (r *PolicyRepository) CheckDuplicatePolicy(ctx context.Context, accountNumber, proposalNumber, companyID string, excludePolicyID string) (bool, error) {
	collection := r.db.Collection(PolicyCollection)

	// 构建查询条件
	conditions := []bson.M{
		{"proposal_number": proposalNumber}, // 投保单号必须唯一
	}

	// 只有当账户号不为空时才检查账户号唯一性
	if strings.TrimSpace(accountNumber) != "" {
		conditions = append(conditions, bson.M{"account_number": accountNumber})
	}

	filter := bson.M{
		"company_id": companyID,
		"$or":        conditions,
	}

	// 排除指定的保单ID（用于更新时检查）
	if excludePolicyID != "" {
		filter["policy_id"] = bson.M{"$ne": excludePolicyID}
	}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// GetPolicyStatistics 获取保单统计信息
func (r *PolicyRepository) GetPolicyStatistics(ctx context.Context, companyID string) (*model.PolicyStatistics, error) {
	collection := r.db.Collection(PolicyCollection)

	filter := bson.M{"company_id": companyID}

	// 聚合查询统计信息
	pipeline := []bson.M{
		{"$match": filter},
		{
			"$group": bson.M{
				"_id":                   nil,
				"total_policies":        bson.M{"$sum": 1},
				"total_premium":         bson.M{"$sum": "$actual_premium"},
				"total_aum":             bson.M{"$sum": "$aum"},
				"total_expected_fee":    bson.M{"$sum": "$expected_fee"},
				"surrendered_count":     bson.M{"$sum": bson.M{"$cond": []interface{}{"$is_surrendered", 1, 0}}},
				"employee_count":        bson.M{"$sum": bson.M{"$cond": []interface{}{"$is_employee", 1, 0}}},
				"cooling_period_count":  bson.M{"$sum": bson.M{"$cond": []interface{}{"$past_cooling_period", 1, 0}}},
				"paid_commission_count": bson.M{"$sum": bson.M{"$cond": []interface{}{"$is_paid_commission", 1, 0}}},
			},
		},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return &model.PolicyStatistics{}, nil
	}

	result := results[0]

	return &model.PolicyStatistics{
		TotalPolicies:       getInt64FromInterface(result["total_policies"]),
		TotalPremium:        getFloat64FromInterface(result["total_premium"]),
		TotalAUM:            getFloat64FromInterface(result["total_aum"]),
		TotalExpectedFee:    getFloat64FromInterface(result["total_expected_fee"]),
		SurrenderedCount:    getInt64FromInterface(result["surrendered_count"]),
		EmployeeCount:       getInt64FromInterface(result["employee_count"]),
		CoolingPeriodCount:  getInt64FromInterface(result["cooling_period_count"]),
		PaidCommissionCount: getInt64FromInterface(result["paid_commission_count"]),
	}, nil
}

// BatchUpdatePolicyStatus 批量更新保单状态
func (r *PolicyRepository) BatchUpdatePolicyStatus(ctx context.Context, policyIDs []string, updates bson.M) error {
	collection := r.db.Collection(PolicyCollection)

	updates["updated_at"] = time.Now()

	filter := bson.M{"policy_id": bson.M{"$in": policyIDs}}

	_, err := collection.UpdateMany(ctx, filter, bson.M{"$set": updates})
	return err
}

// GetPoliciesByIDs 根据ID列表获取保单
func (r *PolicyRepository) GetPoliciesByIDs(ctx context.Context, policyIDs []string) ([]model.Policy, error) {
	collection := r.db.Collection(PolicyCollection)

	filter := bson.M{"policy_id": bson.M{"$in": policyIDs}}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var policies []model.Policy
	if err = cursor.All(ctx, &policies); err != nil {
		return nil, err
	}

	return policies, nil
}

// getNextSerialNumber 获取下一个序号
func (r *PolicyRepository) getNextSerialNumber(ctx context.Context, companyID string) (int, error) {
	collection := r.db.Collection(PolicyCollection)

	// 查找当前公司的最大序号
	filter := bson.M{"company_id": companyID}
	opts := options.FindOne().SetSort(bson.D{{Key: "serial_number", Value: -1}})

	var policy model.Policy
	err := collection.FindOne(ctx, filter, opts).Decode(&policy)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 1, nil // 第一个保单，序号为1
		}
		return 0, err
	}

	return policy.SerialNumber + 1, nil
}

// 辅助函数：从interface{}转换为int64
func getInt64FromInterface(v interface{}) int64 {
	switch val := v.(type) {
	case int64:
		return val
	case int32:
		return int64(val)
	case int:
		return int64(val)
	case float64:
		return int64(val)
	case string:
		if i, err := strconv.ParseInt(val, 10, 64); err == nil {
			return i
		}
	}
	return 0
}

// 辅助函数：从interface{}转换为float64
func getFloat64FromInterface(v interface{}) float64 {
	switch val := v.(type) {
	case float64:
		return val
	case float32:
		return float64(val)
	case int64:
		return float64(val)
	case int32:
		return float64(val)
	case int:
		return float64(val)
	case string:
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			return f
		}
	}
	return 0
}
