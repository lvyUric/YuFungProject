package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ActivityLog 系统活动记录
type ActivityLog struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID        string             `bson:"user_id" json:"user_id"`
	Username      string             `bson:"username" json:"username"`
	CompanyID     string             `bson:"company_id" json:"company_id"`
	CompanyName   string             `bson:"company_name" json:"company_name"`
	OperationType string             `bson:"operation_type" json:"operation_type"` // create/update/delete/view/export
	ModuleName    string             `bson:"module_name" json:"module_name"`
	OperationDesc string             `bson:"operation_desc" json:"operation_desc"`
	RequestURL    string             `bson:"request_url" json:"request_url"`
	RequestMethod string             `bson:"request_method" json:"request_method"`
	RequestParams interface{}        `bson:"request_params" json:"request_params"`
	IPAddress     string             `bson:"ip_address" json:"ip_address"`
	UserAgent     string             `bson:"user_agent" json:"user_agent"`
	OperationTime time.Time          `bson:"operation_time" json:"operation_time"`
	ExecutionTime int64              `bson:"execution_time" json:"execution_time"`     // 执行耗时(ms)
	ResultStatus  string             `bson:"result_status" json:"result_status"`       // success/failure
	TargetID      string             `bson:"target_id" json:"target_id,omitempty"`     // 操作目标ID
	TargetName    string             `bson:"target_name" json:"target_name,omitempty"` // 操作目标名称
}

// ActivityLogQuery 活动记录查询参数
type ActivityLogQuery struct {
	Page          int    `json:"page" form:"page"`
	PageSize      int    `json:"page_size" form:"page_size"`
	CompanyID     string `json:"company_id" form:"company_id"`
	UserID        string `json:"user_id" form:"user_id"`
	OperationType string `json:"operation_type" form:"operation_type"`
	ModuleName    string `json:"module_name" form:"module_name"`
	StartTime     string `json:"start_time" form:"start_time"`
	EndTime       string `json:"end_time" form:"end_time"`
}

// ActivityLogResponse 活动记录响应
type ActivityLogResponse struct {
	Total int64         `json:"total"`
	List  []ActivityLog `json:"list"`
}

// 操作类型常量
const (
	OperationTypeCreate = "create"
	OperationTypeUpdate = "update"
	OperationTypeDelete = "delete"
	OperationTypeView   = "view"
	OperationTypeExport = "export"
	OperationTypeImport = "import"
	OperationTypeLogin  = "login"
	OperationTypeLogout = "logout"
)

// 模块名称常量
const (
	ModuleUser     = "用户管理"
	ModuleRole     = "角色管理"
	ModuleMenu     = "菜单管理"
	ModuleCompany  = "公司管理"
	ModulePolicy   = "保单管理"
	ModuleCustomer = "客户管理"
	ModuleSystem   = "系统管理"
	ModuleAuth     = "认证授权"
)
