package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Policy 保单表模型
type Policy struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`    // MongoDB主键ID
	PolicyID string             `bson:"policy_id" json:"policy_id"` // 保单唯一标识，业务主键

	// 基本信息
	SerialNumber   int    `bson:"serial_number" json:"serial_number"`       // 序号
	AccountNumber  string `bson:"account_number" json:"account_number"`     // 账户号
	CustomerNumber string `bson:"customer_number" json:"customer_number"`   // 客户号
	CustomerNameCN string `bson:"customer_name_cn" json:"customer_name_cn"` // 客户中文名
	CustomerNameEN string `bson:"customer_name_en" json:"customer_name_en"` // 客户英文名
	ProposalNumber string `bson:"proposal_number" json:"proposal_number"`   // 投保单号

	// 保单信息
	PolicyCurrency    string     `bson:"policy_currency" json:"policy_currency"`         // 保单币种(USD/HKD/CNY)
	Partner           string     `bson:"partner" json:"partner"`                         // 合作伙伴
	ReferralCode      string     `bson:"referral_code" json:"referral_code"`             // 转介编号
	HKManager         string     `bson:"hk_manager" json:"hk_manager"`                   // 港分客户经理
	ReferralPM        string     `bson:"referral_pm" json:"referral_pm"`                 // 转介理财经理
	ReferralBranch    string     `bson:"referral_branch" json:"referral_branch"`         // 转介分行
	ReferralSubBranch string     `bson:"referral_sub_branch" json:"referral_sub_branch"` // 转介支行
	ReferralDate      *time.Time `bson:"referral_date" json:"referral_date"`             // 转介日期

	// 退保信息
	IsSurrendered bool `bson:"is_surrendered" json:"is_surrendered"` // 签单后是否退保

	// 缴费信息
	PaymentDate    *time.Time `bson:"payment_date" json:"payment_date"`       // 缴费日期
	EffectiveDate  *time.Time `bson:"effective_date" json:"effective_date"`   // 生效日期
	PaymentMethod  string     `bson:"payment_method" json:"payment_method"`   // 缴费方式(期缴、趸缴、预缴)
	PaymentYears   int        `bson:"payment_years" json:"payment_years"`     // 缴费年期
	PaymentPeriods int        `bson:"payment_periods" json:"payment_periods"` // 期缴期数
	ActualPremium  float64    `bson:"actual_premium" json:"actual_premium"`   // 实际缴纳保费
	AUM            float64    `bson:"aum" json:"aum"`                         // AUM

	// 状态信息
	PastCoolingPeriod bool `bson:"past_cooling_period" json:"past_cooling_period"` // 是否已过冷静期
	IsPaidCommission  bool `bson:"is_paid_commission" json:"is_paid_commission"`   // 是否支付佣金
	IsEmployee        bool `bson:"is_employee" json:"is_employee"`                 // 是否员工

	// 费用信息
	ReferralRate   float64    `bson:"referral_rate" json:"referral_rate"`       // 转介费率
	ExchangeRate   float64    `bson:"exchange_rate" json:"exchange_rate"`       // 汇率
	ExpectedFee    float64    `bson:"expected_fee" json:"expected_fee"`         // 预计转介费
	PaymentPayDate *time.Time `bson:"payment_pay_date" json:"payment_pay_date"` // 支付日期

	// 产品信息
	InsuranceCompany string `bson:"insurance_company" json:"insurance_company"` // 承保公司
	ProductName      string `bson:"product_name" json:"product_name"`           // 保险产品名称
	ProductType      string `bson:"product_type" json:"product_type"`           // 产品类型

	// 其他信息
	Remark    string `bson:"remark" json:"remark"`         // 备注说明
	CompanyID string `bson:"company_id" json:"company_id"` // 所属公司ID（多租户隔离）

	// 系统字段
	CreatedBy string    `bson:"created_by" json:"created_by"` // 创建人
	UpdatedBy string    `bson:"updated_by" json:"updated_by"` // 更新人
	CreatedAt time.Time `bson:"created_at" json:"created_at"` // 创建时间
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"` // 更新时间
}

// PolicyCreateRequest 创建保单请求
type PolicyCreateRequest struct {
	AccountNumber     string     `json:"account_number" label:"账户号"` // 改为非必填
	CustomerNumber    string     `json:"customer_number" binding:"required" label:"客户号"`
	CustomerNameCN    string     `json:"customer_name_cn" binding:"required" label:"客户中文名"`
	CustomerNameEN    string     `json:"customer_name_en" label:"客户英文名"`
	ProposalNumber    string     `json:"proposal_number" binding:"required" label:"投保单号"`
	PolicyCurrency    string     `json:"policy_currency" binding:"required,oneof=USD HKD CNY" label:"保单币种"`
	Partner           string     `json:"partner" label:"合作伙伴"`
	ReferralCode      string     `json:"referral_code" label:"转介编号"`
	HKManager         string     `json:"hk_manager" label:"港分客户经理"`
	ReferralPM        string     `json:"referral_pm" label:"转介理财经理"`
	ReferralBranch    string     `json:"referral_branch" label:"转介分行"`
	ReferralSubBranch string     `json:"referral_sub_branch" label:"转介支行"`
	ReferralDate      *time.Time `json:"referral_date" label:"转介日期"`
	IsSurrendered     bool       `json:"is_surrendered" label:"签单后是否退保"`
	PaymentDate       *time.Time `json:"payment_date" label:"缴费日期"`
	EffectiveDate     *time.Time `json:"effective_date" label:"生效日期"`
	PaymentMethod     string     `json:"payment_method" binding:"oneof=期缴 趸缴 预缴" label:"缴费方式"`
	PaymentYears      int        `json:"payment_years" label:"缴费年期"`
	PaymentPeriods    int        `json:"payment_periods" label:"期缴期数"`
	ActualPremium     float64    `json:"actual_premium" binding:"min=0" label:"实际缴纳保费"`
	AUM               float64    `json:"aum" binding:"min=0" label:"AUM"`
	PastCoolingPeriod bool       `json:"past_cooling_period" label:"是否已过冷静期"`
	IsPaidCommission  bool       `json:"is_paid_commission" label:"是否支付佣金"`
	IsEmployee        bool       `json:"is_employee" label:"是否员工"`
	ReferralRate      float64    `json:"referral_rate" binding:"min=0,max=100" label:"转介费率"`
	ExchangeRate      float64    `json:"exchange_rate" binding:"min=0" label:"汇率"` // 汇率字段，保留4位小数
	ExpectedFee       float64    `json:"expected_fee" binding:"min=0" label:"预计转介费"`
	PaymentPayDate    *time.Time `json:"payment_pay_date" label:"支付日期"`
	InsuranceCompany  string     `json:"insurance_company" binding:"required" label:"承保公司"`
	ProductName       string     `json:"product_name" binding:"required" label:"保险产品名称"`
	ProductType       string     `json:"product_type" binding:"required" label:"产品类型"`
	Remark            string     `json:"remark" label:"备注说明"`
}

// PolicyUpdateRequest 更新保单请求
type PolicyUpdateRequest struct {
	CustomerNameCN    string     `json:"customer_name_cn" label:"客户中文名"`
	CustomerNameEN    string     `json:"customer_name_en" label:"客户英文名"`
	PolicyCurrency    string     `json:"policy_currency" binding:"omitempty,oneof=USD HKD CNY" label:"保单币种"`
	Partner           string     `json:"partner" label:"合作伙伴"`
	ReferralCode      string     `json:"referral_code" label:"转介编号"`
	HKManager         string     `json:"hk_manager" label:"港分客户经理"`
	ReferralPM        string     `json:"referral_pm" label:"转介理财经理"`
	ReferralBranch    string     `json:"referral_branch" label:"转介分行"`
	ReferralSubBranch string     `json:"referral_sub_branch" label:"转介支行"`
	ReferralDate      *time.Time `json:"referral_date" label:"转介日期"`
	IsSurrendered     *bool      `json:"is_surrendered" label:"签单后是否退保"`
	PaymentDate       *time.Time `json:"payment_date" label:"缴费日期"`
	EffectiveDate     *time.Time `json:"effective_date" label:"生效日期"`
	PaymentMethod     string     `json:"payment_method" binding:"omitempty,oneof=期缴 趸缴 预缴" label:"缴费方式"`
	PaymentYears      *int       `json:"payment_years" label:"缴费年期"`
	PaymentPeriods    *int       `json:"payment_periods" label:"期缴期数"`
	ActualPremium     *float64   `json:"actual_premium" binding:"omitempty,min=0" label:"实际缴纳保费"`
	AUM               *float64   `json:"aum" binding:"omitempty,min=0" label:"AUM"`
	PastCoolingPeriod *bool      `json:"past_cooling_period" label:"是否已过冷静期"`
	IsPaidCommission  *bool      `json:"is_paid_commission" label:"是否支付佣金"`
	IsEmployee        *bool      `json:"is_employee" label:"是否员工"`
	ReferralRate      *float64   `json:"referral_rate" binding:"omitempty,min=0,max=100" label:"转介费率"`
	ExchangeRate      *float64   `json:"exchange_rate" binding:"omitempty,min=0" label:"汇率"`
	ExpectedFee       *float64   `json:"expected_fee" binding:"omitempty,min=0" label:"预计转介费"`
	PaymentPayDate    *time.Time `json:"payment_pay_date" label:"支付日期"`
	InsuranceCompany  string     `json:"insurance_company" label:"承保公司"`
	ProductName       string     `json:"product_name" label:"保险产品名称"`
	ProductType       string     `json:"product_type" label:"产品类型"`
	Remark            string     `json:"remark" label:"备注说明"`
}

// PolicyQueryRequest 查询保单请求
type PolicyQueryRequest struct {
	Page               int        `form:"page" binding:"min=1" label:"页码"`
	PageSize           int        `form:"page_size" binding:"min=1,max=100" label:"每页数量"`
	AccountNumber      string     `form:"account_number" label:"账户号"`
	CustomerNumber     string     `form:"customer_number" label:"客户号"`
	CustomerNameCN     string     `form:"customer_name_cn" label:"客户中文名"`
	CustomerNameEN     string     `form:"customer_name_en" label:"客户英文名"`
	ProposalNumber     string     `form:"proposal_number" label:"投保单号"`
	PolicyCurrency     string     `form:"policy_currency" binding:"omitempty,oneof=USD HKD CNY" label:"保单币种"`
	Partner            string     `form:"partner" label:"合作伙伴"`
	ReferralCode       string     `form:"referral_code" label:"转介编号"`
	HKManager          string     `form:"hk_manager" label:"港分客户经理"`
	ReferralPM         string     `form:"referral_pm" label:"转介理财经理"`
	ReferralBranch     string     `form:"referral_branch" label:"转介分行"`
	ReferralSubBranch  string     `form:"referral_sub_branch" label:"转介支行"`
	PaymentMethod      string     `form:"payment_method" binding:"omitempty,oneof=期缴 趸缴 预缴" label:"缴费方式"`
	InsuranceCompany   string     `form:"insurance_company" label:"承保公司"`
	ProductName        string     `form:"product_name" label:"保险产品名称"`
	ProductType        string     `form:"product_type" label:"产品类型"`
	IsSurrendered      *bool      `form:"is_surrendered" label:"是否退保"`
	PastCoolingPeriod  *bool      `form:"past_cooling_period" label:"是否已过冷静期"`
	IsPaidCommission   *bool      `form:"is_paid_commission" label:"是否支付佣金"`
	IsEmployee         *bool      `form:"is_employee" label:"是否员工"`
	ReferralDateStart  *time.Time `form:"referral_date_start" label:"转介日期开始"`
	ReferralDateEnd    *time.Time `form:"referral_date_end" label:"转介日期结束"`
	PaymentDateStart   *time.Time `form:"payment_date_start" label:"缴费日期开始"`
	PaymentDateEnd     *time.Time `form:"payment_date_end" label:"缴费日期结束"`
	EffectiveDateStart *time.Time `form:"effective_date_start" label:"生效日期开始"`
	EffectiveDateEnd   *time.Time `form:"effective_date_end" label:"生效日期结束"`
	SortBy             string     `form:"sort_by" label:"排序字段"`
	SortOrder          string     `form:"sort_order" binding:"omitempty,oneof=asc desc" label:"排序方向"`
}

// PolicyResponse 保单响应
type PolicyResponse struct {
	*Policy
}

// PolicyListResponse 保单列表响应
type PolicyListResponse struct {
	List       []PolicyResponse `json:"list"`        // 保单列表
	Total      int64            `json:"total"`       // 总数
	Page       int              `json:"page"`        // 当前页
	PageSize   int              `json:"page_size"`   // 每页数量
	TotalPages int              `json:"total_pages"` // 总页数
}

// PolicyStatistics 保单统计
type PolicyStatistics struct {
	TotalPolicies       int64   `json:"total_policies"`        // 总保单数
	TotalPremium        float64 `json:"total_premium"`         // 总保费
	TotalAUM            float64 `json:"total_aum"`             // 总AUM
	TotalExpectedFee    float64 `json:"total_expected_fee"`    // 总预计转介费
	SurrenderedCount    int64   `json:"surrendered_count"`     // 退保数量
	EmployeeCount       int64   `json:"employee_count"`        // 员工保单数量
	CoolingPeriodCount  int64   `json:"cooling_period_count"`  // 已过冷静期数量
	PaidCommissionCount int64   `json:"paid_commission_count"` // 已支付佣金数量
}

// BatchUpdatePolicyStatusRequest 批量更新保单状态请求
type BatchUpdatePolicyStatusRequest struct {
	PolicyIDs         []string `json:"policy_ids" binding:"required,min=1"` // 保单ID数组
	IsSurrendered     *bool    `json:"is_surrendered" label:"是否退保"`
	PastCoolingPeriod *bool    `json:"past_cooling_period" label:"是否已过冷静期"`
	IsPaidCommission  *bool    `json:"is_paid_commission" label:"是否支付佣金"`
}

// PolicyImportRequest 保单导入请求
type PolicyImportRequest struct {
	Data []PolicyCreateRequest `json:"data" binding:"required,min=1" label:"保单数据"`
}

// PolicyExportRequest 保单导出请求
type PolicyExportRequest struct {
	PolicyIDs  []string `json:"policy_ids" label:"保单ID数组"` // 为空则导出所有
	ExportType string   `json:"export_type" binding:"required,oneof=xlsx csv" label:"导出类型"`
}

// PolicyImportFileRequest 保单文件导入请求
type PolicyImportFileRequest struct {
	SkipHeader     bool   `form:"skip_header"`     // 是否跳过表头行
	UpdateExisting bool   `form:"update_existing"` // 是否更新已存在的数据
	UserID         string `form:"-"`               // 由中间件设置
	CompanyID      string `form:"-"`               // 由中间件设置
}

// PolicyImportResponse 保单导入响应
type PolicyImportResponse struct {
	SuccessCount int                   `json:"success_count"` // 成功导入数量
	ErrorCount   int                   `json:"error_count"`   // 错误数量
	TotalCount   int                   `json:"total_count"`   // 总数量
	Errors       []PolicyImportError   `json:"errors"`        // 错误详情
	Preview      []PolicyCreateRequest `json:"preview"`       // 预览数据（仅预览时返回）
}

// PolicyImportError 保单导入错误
type PolicyImportError struct {
	Row    int      `json:"row"`    // 错误行号
	Errors []string `json:"errors"` // 错误信息列表
	Data   any      `json:"data"`   // 错误数据
}
