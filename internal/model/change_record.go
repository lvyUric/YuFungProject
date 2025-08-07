package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ChangeRecord 数据变更记录模型
type ChangeRecord struct {
	ID            primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	ChangeID      string                 `bson:"change_id" json:"change_id"`                             // 变更记录唯一标识
	TableName     string                 `bson:"table_name" json:"table_name"`                           // 表名
	RecordID      string                 `bson:"record_id" json:"record_id"`                             // 记录ID
	UserID        string                 `bson:"user_id" json:"user_id"`                                 // 操作用户ID
	Username      string                 `bson:"username" json:"username"`                               // 用户名
	CompanyID     string                 `bson:"company_id" json:"company_id"`                           // 所属公司ID
	ChangeType    string                 `bson:"change_type" json:"change_type"`                         // 变更类型：insert/update/delete
	OldValues     map[string]interface{} `bson:"old_values,omitempty" json:"old_values,omitempty"`       // 变更前数据
	NewValues     map[string]interface{} `bson:"new_values,omitempty" json:"new_values,omitempty"`       // 变更后数据
	ChangedFields []string               `bson:"changed_fields" json:"changed_fields"`                   // 变更字段列表
	ChangeTime    time.Time              `bson:"change_time" json:"change_time"`                         // 变更时间
	ChangeReason  string                 `bson:"change_reason,omitempty" json:"change_reason,omitempty"` // 变更原因
	IPAddress     string                 `bson:"ip_address,omitempty" json:"ip_address,omitempty"`       // IP地址
	UserAgent     string                 `bson:"user_agent,omitempty" json:"user_agent,omitempty"`       // 浏览器信息
}

// ChangeRecordListParams 变更记录查询参数
type ChangeRecordListParams struct {
	TableName  string `json:"table_name" form:"table_name"`   // 表名
	RecordID   string `json:"record_id" form:"record_id"`     // 记录ID
	UserID     string `json:"user_id" form:"user_id"`         // 用户ID
	CompanyID  string `json:"company_id" form:"company_id"`   // 公司ID
	ChangeType string `json:"change_type" form:"change_type"` // 变更类型
	StartTime  string `json:"start_time" form:"start_time"`   // 开始时间
	EndTime    string `json:"end_time" form:"end_time"`       // 结束时间
	Page       int    `json:"page" form:"page"`               // 页码
	PageSize   int    `json:"page_size" form:"page_size"`     // 每页数量
}

// ChangeRecordResponse 变更记录响应
type ChangeRecordResponse struct {
	ID            string                 `json:"id"`
	ChangeID      string                 `json:"change_id"`
	TableName     string                 `json:"table_name"`
	RecordID      string                 `json:"record_id"`
	UserID        string                 `json:"user_id"`
	Username      string                 `json:"username"`
	CompanyID     string                 `json:"company_id"`
	ChangeType    string                 `json:"change_type"`
	OldValues     map[string]interface{} `json:"old_values,omitempty"`
	NewValues     map[string]interface{} `json:"new_values,omitempty"`
	ChangedFields []string               `json:"changed_fields"`
	ChangeTime    time.Time              `json:"change_time"`
	ChangeReason  string                 `json:"change_reason,omitempty"`
	IPAddress     string                 `json:"ip_address,omitempty"`
	UserAgent     string                 `json:"user_agent,omitempty"`
	// 格式化显示字段
	ChangeTimeFormatted string         `json:"change_time_formatted"`
	ChangeDetails       []ChangeDetail `json:"change_details"`
}

// ChangeDetail 变更详情
type ChangeDetail struct {
	FieldName    string      `json:"field_name"`     // 字段名
	FieldLabel   string      `json:"field_label"`    // 字段显示名
	OldValue     interface{} `json:"old_value"`      // 旧值
	NewValue     interface{} `json:"new_value"`      // 新值
	OldValueText string      `json:"old_value_text"` // 旧值显示文本
	NewValueText string      `json:"new_value_text"` // 新值显示文本
}

// FieldLabel 字段标签映射（用于显示友好的字段名）
var PolicyFieldLabels = map[string]string{
	"account_number":      "账户号",
	"customer_number":     "客户号",
	"customer_name_cn":    "客户中文名",
	"customer_name_en":    "客户英文名",
	"proposal_number":     "投保单号",
	"policy_currency":     "保单币种",
	"partner":             "合作伙伴",
	"referral_code":       "转介编号",
	"hk_manager":          "港分客户经理",
	"referral_pm":         "转介理财经理",
	"referral_branch":     "转介分行",
	"referral_sub_branch": "转介支行",
	"referral_date":       "转介日期",
	"is_surrendered":      "签单后是否退保",
	"payment_date":        "缴费日期",
	"effective_date":      "生效日期",
	"payment_method":      "缴费方式",
	"payment_years":       "缴费年期",
	"payment_periods":     "期缴期数",
	"actual_premium":      "实际缴纳保费",
	"aum":                 "AUM",
	"commission_rate":     "佣金比例",
	"policy_year":         "保单年度",
	"remark":              "备注",
	"status":              "状态",
}

// GetFieldLabel 获取字段显示标签
func GetFieldLabel(tableName, fieldName string) string {
	switch tableName {
	case "policies":
		if label, exists := PolicyFieldLabels[fieldName]; exists {
			return label
		}
	}
	return fieldName
}
