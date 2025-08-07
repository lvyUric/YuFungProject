package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User 用户表模型
type User struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`                // MongoDB主键ID
	UserID        string             `bson:"user_id" json:"user_id"`                 // 用户唯一标识，业务主键
	Username      string             `bson:"username" json:"username"`               // 登录用户名，唯一
	DisplayName   string             `bson:"display_name" json:"display_name"`       // 用户显示名称（中文名等）
	CompanyID     string             `bson:"company_id" json:"company_id"`           // 所属保险经纪公司ID
	RoleIDs       []string           `bson:"role_ids" json:"role_ids"`               // 用户角色ID数组，支持多角色
	Status        string             `bson:"status" json:"status"`                   // 用户状态：active=激活, inactive=禁用, locked=锁定
	LastLoginTime *time.Time         `bson:"last_login_time" json:"last_login_time"` // 最后登录时间，可为空
	PasswordHash  string             `bson:"password_hash" json:"-"`                 // 密码哈希值，不返回给前端
	Email         string             `bson:"email" json:"email"`                     // 邮箱地址，可选
	Phone         string             `bson:"phone" json:"phone"`                     // 手机号码，可选
	Remark        string             `bson:"remark" json:"remark"`                   // 备注信息
	LoginAttempts int                `bson:"login_attempts" json:"login_attempts"`   // 登录失败次数，用于防暴力破解
	LockedUntil   *time.Time         `bson:"locked_until" json:"locked_until"`       // 账户锁定截止时间，可为空
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`           // 创建时间
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`           // 更新时间
}

// Company 保险经纪公司表模型
type Company struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`          // MongoDB主键ID
	CompanyID   string             `bson:"company_id" json:"company_id"`     // 公司唯一标识，业务主键
	CompanyName string             `bson:"company_name" json:"company_name"` // 公司名称
	CompanyCode string             `bson:"company_code" json:"company_code"` // 内部公司代码

	// 负责人信息
	ContactPerson string `bson:"contact_person" json:"contact_person"` // 联络人

	// 联系方式
	TelNo        string `bson:"tel_no" json:"tel_no"`               // 固定电话
	Mobile       string `bson:"mobile" json:"mobile"`               // 移动电话
	ContactPhone string `bson:"contact_phone" json:"contact_phone"` // 联系电话（保留兼容）
	Email        string `bson:"email" json:"email"`                 // 邮箱地址

	// 中文地址信息
	AddressCNProvince string `bson:"address_cn_province" json:"address_cn_province"` // 中文地址-省/自治区/直辖市
	AddressCNCity     string `bson:"address_cn_city" json:"address_cn_city"`         // 中文地址-市
	AddressCNDistrict string `bson:"address_cn_district" json:"address_cn_district"` // 中文地址-县/区
	AddressCNDetail   string `bson:"address_cn_detail" json:"address_cn_detail"`     // 中文地址-详细地址

	// 英文地址信息
	AddressENProvince string `bson:"address_en_province" json:"address_en_province"` // 英文地址-省/自治区/直辖市
	AddressENCity     string `bson:"address_en_city" json:"address_en_city"`         // 英文地址-市
	AddressENDistrict string `bson:"address_en_district" json:"address_en_district"` // 英文地址-县/区
	AddressENDetail   string `bson:"address_en_detail" json:"address_en_detail"`     // 英文地址-详细地址

	Address string `bson:"address" json:"address"` // 原有地址字段（保留兼容）

	// 业务信息
	BrokerCode string `bson:"broker_code" json:"broker_code"` // 经纪人代码
	Link       string `bson:"link" json:"link"`               // 相关链接

	// 登录信息
	Username     string `bson:"username" json:"username"` // 用户名
	PasswordHash string `bson:"password_hash" json:"-"`   // 密码哈希值（不返回前端）

	// 系统字段
	ValidStartDate   time.Time `bson:"valid_start_date" json:"valid_start_date"`     // 有效期开始日期
	ValidEndDate     time.Time `bson:"valid_end_date" json:"valid_end_date"`         // 有效期结束日期
	UserQuota        int       `bson:"user_quota" json:"user_quota"`                 // 允许创建的用户数量配额
	CurrentUserCount int       `bson:"current_user_count" json:"current_user_count"` // 当前已创建的用户数量
	Status           string    `bson:"status" json:"status"`                         // 状态：active=有效, inactive=停用, expired=过期
	Remark           string    `bson:"remark" json:"remark"`                         // 备注信息（保留兼容）
	SubmittedBy      string    `bson:"submitted_by" json:"submitted_by"`             // 提交人
	CreatedAt        time.Time `bson:"created_at" json:"created_at"`                 // 创建时间（提交时间）
	UpdatedAt        time.Time `bson:"updated_at" json:"updated_at"`                 // 更新时间
}

// Role 角色表模型
type Role struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`      // MongoDB主键ID
	RoleID    string             `bson:"role_id" json:"role_id"`       // 角色唯一标识，业务主键
	RoleName  string             `bson:"role_name" json:"role_name"`   // 角色名称
	RoleKey   string             `bson:"role_key" json:"role_key"`     // 角色标识符，用于权限判断
	CompanyID string             `bson:"company_id" json:"company_id"` // 所属公司ID，空表示平台级角色
	SortOrder int                `bson:"sort_order" json:"sort_order"` // 排序号
	DataScope string             `bson:"data_scope" json:"data_scope"` // 数据权限范围：all=全部, company=本公司, self=个人
	MenuIDs   []string           `bson:"menu_ids" json:"menu_ids"`     // 菜单权限ID数组
	Status    string             `bson:"status" json:"status"`         // 状态：enable=启用, disable=禁用
	Remark    string             `bson:"remark" json:"remark"`         // 备注信息
	CreatedAt time.Time          `bson:"created_at" json:"created_at"` // 创建时间
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"` // 更新时间
}

// Menu 菜单表模型
type Menu struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`                // MongoDB主键ID
	MenuID         string             `bson:"menu_id" json:"menu_id"`                 // 菜单唯一标识，业务主键
	ParentID       string             `bson:"parent_id" json:"parent_id"`             // 父菜单ID，根菜单为空
	MenuName       string             `bson:"menu_name" json:"menu_name"`             // 菜单名称
	MenuType       string             `bson:"menu_type" json:"menu_type"`             // 菜单类型：directory=目录, menu=菜单, button=按钮
	RoutePath      string             `bson:"route_path" json:"route_path"`           // 路由路径
	Component      string             `bson:"component" json:"component"`             // 组件路径
	PermissionCode string             `bson:"permission_code" json:"permission_code"` // 权限标识符
	Icon           string             `bson:"icon" json:"icon"`                       // 菜单图标
	SortOrder      int                `bson:"sort_order" json:"sort_order"`           // 排序号
	Visible        bool               `bson:"visible" json:"visible"`                 // 是否在菜单中显示
	Status         string             `bson:"status" json:"status"`                   // 状态：enable=启用, disable=禁用
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`           // 创建时间
	UpdatedAt      time.Time          `bson:"updated_at" json:"updated_at"`           // 更新时间
}

// TableStructure 动态表结构定义表模型
type TableStructure struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`          // MongoDB主键ID
	TableID     string             `bson:"table_id" json:"table_id"`         // 表结构唯一标识，业务主键
	TableName   string             `bson:"table_name" json:"table_name"`     // 表名（英文）
	DisplayName string             `bson:"display_name" json:"display_name"` // 表显示名称（中文）
	TableType   string             `bson:"table_type" json:"table_type"`     // 表类型：system=系统表, custom=自定义表
	CompanyID   string             `bson:"company_id" json:"company_id"`     // 所属公司ID，空表示平台级表
	Description string             `bson:"description" json:"description"`   // 表描述
	Status      string             `bson:"status" json:"status"`             // 状态：active=启用, inactive=禁用
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`     // 创建时间
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`     // 更新时间
}

// FieldDefinition 字段定义表模型
type FieldDefinition struct {
	ID              primitive.ObjectID     `bson:"_id,omitempty" json:"id"`                  // MongoDB主键ID
	FieldID         string                 `bson:"field_id" json:"field_id"`                 // 字段唯一标识，业务主键
	TableID         string                 `bson:"table_id" json:"table_id"`                 // 所属表结构ID
	FieldName       string                 `bson:"field_name" json:"field_name"`             // 字段名（英文）
	DisplayName     string                 `bson:"display_name" json:"display_name"`         // 字段显示名称（中文）
	FieldType       string                 `bson:"field_type" json:"field_type"`             // 字段类型：string=文本, number=数字, date=日期, boolean=布尔, enum=枚举, file=文件
	FieldLength     int                    `bson:"field_length" json:"field_length"`         // 字段长度限制
	Required        bool                   `bson:"required" json:"required"`                 // 是否必填
	DefaultValue    string                 `bson:"default_value" json:"default_value"`       // 默认值
	EnumOptions     []string               `bson:"enum_options" json:"enum_options"`         // 枚举选项（当字段类型为enum时）
	ValidationRules map[string]interface{} `bson:"validation_rules" json:"validation_rules"` // 验证规则JSON
	SortOrder       int                    `bson:"sort_order" json:"sort_order"`             // 排序号
	Visible         bool                   `bson:"visible" json:"visible"`                   // 是否在表单中显示
	CreatedAt       time.Time              `bson:"created_at" json:"created_at"`             // 创建时间
	UpdatedAt       time.Time              `bson:"updated_at" json:"updated_at"`             // 更新时间
}

// OperationLog 系统操作日志表模型
type OperationLog struct {
	ID            primitive.ObjectID     `bson:"_id,omitempty" json:"id"`              // MongoDB主键ID
	LogID         string                 `bson:"log_id" json:"log_id"`                 // 日志唯一标识，业务主键
	UserID        string                 `bson:"user_id" json:"user_id"`               // 操作用户ID
	Username      string                 `bson:"username" json:"username"`             // 操作用户名
	CompanyID     string                 `bson:"company_id" json:"company_id"`         // 所属公司ID
	OperationType string                 `bson:"operation_type" json:"operation_type"` // 操作类型：create=创建, update=更新, delete=删除, view=查看, export=导出
	ModuleName    string                 `bson:"module_name" json:"module_name"`       // 模块名称
	OperationDesc string                 `bson:"operation_desc" json:"operation_desc"` // 操作描述
	RequestURL    string                 `bson:"request_url" json:"request_url"`       // 请求URL
	RequestMethod string                 `bson:"request_method" json:"request_method"` // 请求方法：GET, POST, PUT, DELETE等
	RequestParams map[string]interface{} `bson:"request_params" json:"request_params"` // 请求参数JSON
	IPAddress     string                 `bson:"ip_address" json:"ip_address"`         // 操作者IP地址
	UserAgent     string                 `bson:"user_agent" json:"user_agent"`         // 浏览器用户代理
	OperationTime time.Time              `bson:"operation_time" json:"operation_time"` // 操作时间
	ExecutionTime int                    `bson:"execution_time" json:"execution_time"` // 执行耗时（毫秒）
	ResultStatus  string                 `bson:"result_status" json:"result_status"`   // 执行结果：success=成功, failure=失败
	ErrorMessage  string                 `bson:"error_message" json:"error_message"`   // 错误信息（失败时记录）
}

// DataChangeLog 数据变更记录表模型
type DataChangeLog struct {
	ID            primitive.ObjectID     `bson:"_id,omitempty" json:"id"`              // MongoDB主键ID
	ChangeID      string                 `bson:"change_id" json:"change_id"`           // 变更记录唯一标识，业务主键
	TableName     string                 `bson:"table_name" json:"table_name"`         // 操作的表名
	RecordID      string                 `bson:"record_id" json:"record_id"`           // 被变更记录的ID
	UserID        string                 `bson:"user_id" json:"user_id"`               // 操作用户ID
	CompanyID     string                 `bson:"company_id" json:"company_id"`         // 所属公司ID
	ChangeType    string                 `bson:"change_type" json:"change_type"`       // 变更类型：insert=新增, update=更新, delete=删除
	OldValues     map[string]interface{} `bson:"old_values" json:"old_values"`         // 变更前数据JSON
	NewValues     map[string]interface{} `bson:"new_values" json:"new_values"`         // 变更后数据JSON
	ChangedFields []string               `bson:"changed_fields" json:"changed_fields"` // 变更字段列表
	ChangeTime    time.Time              `bson:"change_time" json:"change_time"`       // 变更时间
	ChangeReason  string                 `bson:"change_reason" json:"change_reason"`   // 变更原因
}

// UserCreateRequest 用户创建请求
type UserCreateRequest struct {
	Username    string   `json:"username" binding:"required,min=3,max=50"`      // 用户名
	DisplayName string   `json:"display_name" binding:"required,min=2,max=100"` // 显示名称
	Password    string   `json:"password" binding:"required,min=8"`             // 密码
	CompanyID   string   `json:"company_id" binding:"required"`                 // 所属公司ID
	RoleIDs     []string `json:"role_ids" binding:"required"`                   // 角色ID数组
	Email       string   `json:"email" binding:"omitempty,email"`               // 邮箱地址
	Phone       string   `json:"phone" binding:"omitempty"`                     // 手机号码
	Remark      string   `json:"remark"`                                        // 备注信息
}

// UserUpdateRequest 用户更新请求
type UserUpdateRequest struct {
	DisplayName string   `json:"display_name" binding:"omitempty,min=2,max=100"`          // 显示名称
	RoleIDs     []string `json:"role_ids" binding:"omitempty"`                            // 角色ID数组
	Email       string   `json:"email" binding:"omitempty,email"`                         // 邮箱地址
	Phone       string   `json:"phone" binding:"omitempty"`                               // 手机号码
	Remark      string   `json:"remark"`                                                  // 备注信息
	Status      string   `json:"status" binding:"omitempty,oneof=active inactive locked"` // 用户状态
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"` // 用户名
	Password string `json:"password" binding:"required"` // 密码
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username    string `json:"username" binding:"required,min=3,max=50"`      // 用户名
	DisplayName string `json:"display_name" binding:"required,min=2,max=100"` // 显示名称
	Password    string `json:"password" binding:"required,min=8"`             // 密码
	Email       string `json:"email" binding:"omitempty,email"`               // 邮箱地址
	Phone       string `json:"phone" binding:"omitempty"`                     // 手机号码
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`       // 当前密码
	NewPassword string `json:"new_password" binding:"required,min=8"` // 新密码
}

// ResetPasswordRequest 重置密码请求
type ResetPasswordRequest struct {
	UserID      string `json:"user_id" binding:"required"`            // 用户ID
	NewPassword string `json:"new_password" binding:"required,min=8"` // 新密码
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token        string    `json:"token"`         // 访问令牌
	RefreshToken string    `json:"refresh_token"` // 刷新令牌
	ExpiresAt    time.Time `json:"expires_at"`    // 令牌过期时间
	User         UserInfo  `json:"user"`          // 用户信息
}

// UserInfo 用户信息
type UserInfo struct {
	ID          string     `json:"id"`           // 用户ID
	UserID      string     `json:"user_id"`      // 用户唯一标识
	Username    string     `json:"username"`     // 用户名
	DisplayName string     `json:"display_name"` // 显示名称
	CompanyID   string     `json:"company_id"`   // 所属公司ID
	RoleIDs     []string   `json:"role_ids"`     // 角色ID数组
	Status      string     `json:"status"`       // 用户状态
	Email       string     `json:"email"`        // 邮箱地址
	Phone       string     `json:"phone"`        // 手机号码
	LastLogin   *time.Time `json:"last_login"`   // 最后登录时间
}

// UserListResponse 用户列表响应
type UserListResponse struct {
	Users      []UserInfo `json:"users"`       // 用户列表
	Total      int64      `json:"total"`       // 总记录数
	Page       int        `json:"page"`        // 当前页码
	PageSize   int        `json:"page_size"`   // 每页大小
	TotalPages int        `json:"total_pages"` // 总页数
}

// RefreshTokenRequest 刷新令牌请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"` // 刷新令牌
}

// ==========================
// 公司管理相关模型
// ==========================

// CreateCompanyRequest 创建公司请求
type CreateCompanyRequest struct {
	CompanyName string `json:"company_name" binding:"required,min=2,max=100"` // 公司名称
	CompanyCode string `json:"company_code" binding:"omitempty,max=50"`       // 内部公司代码

	// 负责人信息
	ContactPerson string `json:"contact_person" binding:"omitempty,max=100"` // 联络人

	// 联系方式
	TelNo        string `json:"tel_no" binding:"omitempty"`        // 固定电话
	Mobile       string `json:"mobile" binding:"omitempty"`        // 移动电话
	ContactPhone string `json:"contact_phone" binding:"omitempty"` // 联系电话（保留兼容）
	Email        string `json:"email" binding:"required,email"`    // 邮箱地址

	// 中文地址信息
	AddressCNProvince string `json:"address_cn_province" binding:"omitempty,max=50"` // 中文地址-省/自治区/直辖市
	AddressCNCity     string `json:"address_cn_city" binding:"omitempty,max=50"`     // 中文地址-市
	AddressCNDistrict string `json:"address_cn_district" binding:"omitempty,max=50"` // 中文地址-县/区
	AddressCNDetail   string `json:"address_cn_detail" binding:"omitempty,max=200"`  // 中文地址-详细地址

	// 英文地址信息
	AddressENProvince string `json:"address_en_province" binding:"omitempty,max=100"` // 英文地址-省/自治区/直辖市
	AddressENCity     string `json:"address_en_city" binding:"omitempty,max=100"`     // 英文地址-市
	AddressENDistrict string `json:"address_en_district" binding:"omitempty,max=100"` // 英文地址-县/区
	AddressENDetail   string `json:"address_en_detail" binding:"omitempty,max=200"`   // 英文地址-详细地址

	Address string `json:"address" binding:"omitempty,max=500"` // 原有地址字段（保留兼容）

	// 业务信息
	BrokerCode string `json:"broker_code" binding:"omitempty,max=50"` // 经纪人代码
	Link       string `json:"link" binding:"omitempty,url"`           // 相关链接

	// 登录信息
	Username string `json:"username" binding:"omitempty,min=3,max=50"` // 用户名
	Password string `json:"password" binding:"omitempty,min=8"`        // 密码

	// 系统字段
	ValidStartDate string `json:"valid_start_date" binding:"required"`           // 有效期开始日期 (YYYY-MM-DD)
	ValidEndDate   string `json:"valid_end_date" binding:"required"`             // 有效期结束日期 (YYYY-MM-DD)
	UserQuota      int    `json:"user_quota" binding:"required,min=1,max=10000"` // 用户配额
	Remark         string `json:"remark" binding:"omitempty,max=500"`            // 备注信息
}

// UpdateCompanyRequest 更新公司请求
type UpdateCompanyRequest struct {
	CompanyName string `json:"company_name" binding:"omitempty,min=2,max=100"` // 公司名称
	CompanyCode string `json:"company_code" binding:"omitempty,max=50"`        // 内部公司代码

	// 负责人信息
	ContactPerson string `json:"contact_person" binding:"omitempty,max=100"` // 联络人

	// 联系方式
	TelNo        string `json:"tel_no" binding:"omitempty"`        // 固定电话
	Mobile       string `json:"mobile" binding:"omitempty"`        // 移动电话
	ContactPhone string `json:"contact_phone" binding:"omitempty"` // 联系电话（保留兼容）
	Email        string `json:"email" binding:"omitempty,email"`   // 邮箱地址

	// 中文地址信息
	AddressCNProvince string `json:"address_cn_province" binding:"omitempty,max=50"` // 中文地址-省/自治区/直辖市
	AddressCNCity     string `json:"address_cn_city" binding:"omitempty,max=50"`     // 中文地址-市
	AddressCNDistrict string `json:"address_cn_district" binding:"omitempty,max=50"` // 中文地址-县/区
	AddressCNDetail   string `json:"address_cn_detail" binding:"omitempty,max=200"`  // 中文地址-详细地址

	// 英文地址信息
	AddressENProvince string `json:"address_en_province" binding:"omitempty,max=100"` // 英文地址-省/自治区/直辖市
	AddressENCity     string `json:"address_en_city" binding:"omitempty,max=100"`     // 英文地址-市
	AddressENDistrict string `json:"address_en_district" binding:"omitempty,max=100"` // 英文地址-县/区
	AddressENDetail   string `json:"address_en_detail" binding:"omitempty,max=200"`   // 英文地址-详细地址

	Address string `json:"address" binding:"omitempty,max=500"` // 原有地址字段（保留兼容）

	// 业务信息
	BrokerCode string `json:"broker_code" binding:"omitempty,max=50"` // 经纪人代码
	Link       string `json:"link" binding:"omitempty,url"`           // 相关链接

	// 登录信息
	Username string `json:"username" binding:"omitempty,min=3,max=50"` // 用户名
	Password string `json:"password" binding:"omitempty,min=8"`        // 密码

	// 系统字段
	ValidStartDate string `json:"valid_start_date"`                                 // 有效期开始日期
	ValidEndDate   string `json:"valid_end_date"`                                   // 有效期结束日期
	UserQuota      int    `json:"user_quota" binding:"omitempty,min=1,max=10000"`   // 用户配额
	Status         string `json:"status" binding:"omitempty,oneof=active inactive"` // 状态
	Remark         string `json:"remark" binding:"omitempty,max=500"`               // 备注信息（保留兼容）
}

// CompanyQueryRequest 公司查询请求
type CompanyQueryRequest struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`                           // 页码，默认1
	PageSize int    `form:"page_size" binding:"omitempty,min=1"`                      // 每页大小，默认20
	Status   string `form:"status" binding:"omitempty,oneof=active inactive expired"` // 状态筛选
	Keyword  string `form:"keyword"`                                                  // 关键词搜索（公司名称）
}

// CompanyInfo 公司信息
type CompanyInfo struct {
	ID          string `json:"id"`           // 主键ID
	CompanyID   string `json:"company_id"`   // 公司唯一标识
	CompanyName string `json:"company_name"` // 公司名称
	CompanyCode string `json:"company_code"` // 内部公司代码

	// 负责人信息
	ContactPerson string `json:"contact_person"` // 联络人

	// 联系方式
	TelNo        string `json:"tel_no"`        // 固定电话
	Mobile       string `json:"mobile"`        // 移动电话
	ContactPhone string `json:"contact_phone"` // 联系电话（保留兼容）
	Email        string `json:"email"`         // 邮箱地址

	// 中文地址信息
	AddressCNProvince string `json:"address_cn_province"` // 中文地址-省/自治区/直辖市
	AddressCNCity     string `json:"address_cn_city"`     // 中文地址-市
	AddressCNDistrict string `json:"address_cn_district"` // 中文地址-县/区
	AddressCNDetail   string `json:"address_cn_detail"`   // 中文地址-详细地址

	// 英文地址信息
	AddressENProvince string `json:"address_en_province"` // 英文地址-省/自治区/直辖市
	AddressENCity     string `json:"address_en_city"`     // 英文地址-市
	AddressENDistrict string `json:"address_en_district"` // 英文地址-县/区
	AddressENDetail   string `json:"address_en_detail"`   // 英文地址-详细地址

	Address string `json:"address"` // 原有地址字段（保留兼容）

	// 业务信息
	BrokerCode string `json:"broker_code"` // 经纪人代码
	Link       string `json:"link"`        // 相关链接

	// 登录信息
	Username string `json:"username"` // 用户名

	// 系统字段
	ValidStartDate   string `json:"valid_start_date"`   // 有效期开始日期
	ValidEndDate     string `json:"valid_end_date"`     // 有效期结束日期
	UserQuota        int    `json:"user_quota"`         // 用户配额
	CurrentUserCount int    `json:"current_user_count"` // 当前用户数量
	Status           string `json:"status"`             // 状态
	StatusText       string `json:"status_text"`        // 状态文本
	Remark           string `json:"remark"`             // 备注信息（保留兼容）
	SubmittedBy      string `json:"submitted_by"`       // 提交人
	CreatedAt        string `json:"created_at"`         // 创建时间
	UpdatedAt        string `json:"updated_at"`         // 更新时间
}

// CompanyListResponse 公司列表响应
type CompanyListResponse struct {
	Companies  []CompanyInfo `json:"companies"`   // 公司列表
	Total      int64         `json:"total"`       // 总记录数
	Page       int           `json:"page"`        // 当前页码
	PageSize   int           `json:"page_size"`   // 每页大小
	TotalPages int           `json:"total_pages"` // 总页数
}

// CompanyStatsResponse 公司统计响应
type CompanyStatsResponse struct {
	TotalCompanies   int64 `json:"total_companies"`   // 总公司数
	ActiveCompanies  int64 `json:"active_companies"`  // 有效公司数
	ExpiredCompanies int64 `json:"expired_companies"` // 过期公司数
	TotalUsers       int64 `json:"total_users"`       // 总用户数
}

// 导入导出相关模型
// ==========================

// CompanyImportRequest 公司导入请求
type CompanyImportRequest struct {
	SkipHeader     bool `form:"skip_header"`     // 是否跳过表头行
	UpdateExisting bool `form:"update_existing"` // 是否更新已存在的公司
}

// CompanyImportError 导入错误信息
type CompanyImportError struct {
	Row    int      `json:"row"`    // 错误行号
	Errors []string `json:"errors"` // 错误信息列表
	Data   any      `json:"data"`   // 错误数据
}

// CompanyImportResponse 公司导入响应
type CompanyImportResponse struct {
	SuccessCount int                  `json:"success_count"` // 成功导入数量
	ErrorCount   int                  `json:"error_count"`   // 错误数量
	TotalCount   int                  `json:"total_count"`   // 总数量
	Errors       []CompanyImportError `json:"errors"`        // 错误详情
	Preview      []CompanyInfo        `json:"preview"`       // 预览数据（仅预览时返回）
}

// CompanyExportRequest 公司导出请求
type CompanyExportRequest struct {
	Status     string   `json:"status"`      // 状态筛选
	Keyword    string   `json:"keyword"`     // 关键词
	IDs        []string `json:"ids"`         // 指定ID列表
	ExportType string   `json:"export_type"` // 导出类型：all=全部, selected=选中, filtered=筛选
	Format     string   `json:"format"`      // 导出格式：xlsx, csv
	Template   bool     `json:"template"`    // 是否导出模板
}

// CompanyExportResponse 公司导出响应
type CompanyExportResponse struct {
	FileURL       string `json:"file_url"`       // 文件下载URL
	FileName      string `json:"file_name"`      // 文件名
	DownloadToken string `json:"download_token"` // 下载令牌
}

// BatchUpdateUserStatusRequest 批量更新用户状态请求
type BatchUpdateUserStatusRequest struct {
	UserIDs []string `json:"user_ids" binding:"required"`                            // 用户ID列表
	Status  string   `json:"status" binding:"required,oneof=active inactive locked"` // 新状态
}

// UserImportRequest 用户导入请求
type UserImportRequest struct {
	SkipHeader     bool `form:"skip_header"`     // 是否跳过表头行
	UpdateExisting bool `form:"update_existing"` // 是否更新已存在的用户
}

// UserImportError 用户导入错误信息
type UserImportError struct {
	Row    int      `json:"row"`    // 错误行号
	Errors []string `json:"errors"` // 错误信息列表
	Data   any      `json:"data"`   // 错误数据
}

// UserImportResponse 用户导入响应
type UserImportResponse struct {
	SuccessCount int               `json:"success_count"` // 成功导入数量
	ErrorCount   int               `json:"error_count"`   // 错误数量
	TotalCount   int               `json:"total_count"`   // 总数量
	Errors       []UserImportError `json:"errors"`        // 错误详情
	Preview      []UserInfo        `json:"preview"`       // 预览数据（仅预览时返回）
}

// UserExportRequest 用户导出请求
type UserExportRequest struct {
	Status     string   `json:"status"`      // 状态筛选
	CompanyID  string   `json:"company_id"`  // 公司ID筛选
	Keyword    string   `json:"keyword"`     // 关键词
	IDs        []string `json:"ids"`         // 指定ID列表
	ExportType string   `json:"export_type"` // 导出类型：all=全部, selected=选中, filtered=筛选
	Format     string   `json:"format"`      // 导出格式：xlsx, csv
	Template   bool     `json:"template"`    // 是否导出模板
}

// UserExportResponse 用户导出响应
type UserExportResponse struct {
	FileURL       string `json:"file_url"`       // 文件下载URL
	FileName      string `json:"file_name"`      // 文件名
	DownloadToken string `json:"download_token"` // 下载令牌
}

// ==========================
// 角色管理相关模型
// ==========================

// RoleCreateRequest 角色创建请求
type RoleCreateRequest struct {
	RoleName  string   `json:"role_name" binding:"required,min=2,max=50"`            // 角色名称
	RoleKey   string   `json:"role_key" binding:"required,min=2,max=50"`             // 角色标识符
	CompanyID string   `json:"company_id"`                                           // 所属公司ID（空表示平台级角色）
	SortOrder int      `json:"sort_order"`                                           // 排序号
	DataScope string   `json:"data_scope" binding:"required,oneof=all company self"` // 数据权限范围
	MenuIDs   []string `json:"menu_ids"`                                             // 菜单权限ID数组
	Status    string   `json:"status" binding:"omitempty,oneof=enable disable"`      // 状态
	Remark    string   `json:"remark"`                                               // 备注信息
}

// RoleUpdateRequest 角色更新请求
type RoleUpdateRequest struct {
	RoleName  string   `json:"role_name" binding:"omitempty,min=2,max=50"`            // 角色名称
	RoleKey   string   `json:"role_key" binding:"omitempty,min=2,max=50"`             // 角色标识符
	SortOrder int      `json:"sort_order"`                                            // 排序号
	DataScope string   `json:"data_scope" binding:"omitempty,oneof=all company self"` // 数据权限范围
	MenuIDs   []string `json:"menu_ids"`                                              // 菜单权限ID数组
	Status    string   `json:"status" binding:"omitempty,oneof=enable disable"`       // 状态
	Remark    string   `json:"remark"`                                                // 备注信息
}

// RoleInfo 角色信息
type RoleInfo struct {
	ID        string    `json:"id"`         // 角色ID
	RoleID    string    `json:"role_id"`    // 角色唯一标识
	RoleName  string    `json:"role_name"`  // 角色名称
	RoleKey   string    `json:"role_key"`   // 角色标识符
	CompanyID string    `json:"company_id"` // 所属公司ID
	SortOrder int       `json:"sort_order"` // 排序号
	DataScope string    `json:"data_scope"` // 数据权限范围
	MenuIDs   []string  `json:"menu_ids"`   // 菜单权限ID数组
	Status    string    `json:"status"`     // 状态
	Remark    string    `json:"remark"`     // 备注信息
	CreatedAt time.Time `json:"created_at"` // 创建时间
	UpdatedAt time.Time `json:"updated_at"` // 更新时间
}

// RoleListResponse 角色列表响应
type RoleListResponse struct {
	Roles      []RoleInfo `json:"roles"`       // 角色列表
	Total      int64      `json:"total"`       // 总记录数
	Page       int        `json:"page"`        // 当前页码
	PageSize   int        `json:"page_size"`   // 每页大小
	TotalPages int        `json:"total_pages"` // 总页数
}

// RoleQueryRequest 角色查询请求
type RoleQueryRequest struct {
	Page      int    `form:"page" binding:"omitempty,min=1"`              // 页码
	PageSize  int    `form:"page_size" binding:"omitempty,min=1,max=100"` // 每页大小
	RoleName  string `form:"role_name"`                                   // 角色名称搜索
	RoleKey   string `form:"role_key"`                                    // 角色标识符搜索
	CompanyID string `form:"company_id"`                                  // 公司ID筛选
	DataScope string `form:"data_scope"`                                  // 数据权限范围筛选
	Status    string `form:"status"`                                      // 状态筛选
}

// BatchUpdateRoleStatusRequest 批量更新角色状态请求
type BatchUpdateRoleStatusRequest struct {
	RoleIDs []string `json:"role_ids" binding:"required,min=1"`              // 角色ID数组
	Status  string   `json:"status" binding:"required,oneof=enable disable"` // 目标状态
}

// RoleStatsResponse 角色统计响应
type RoleStatsResponse struct {
	TotalRoles    int64 `json:"total_roles"`    // 总角色数
	EnabledRoles  int64 `json:"enabled_roles"`  // 启用角色数
	DisabledRoles int64 `json:"disabled_roles"` // 禁用角色数
	PlatformRoles int64 `json:"platform_roles"` // 平台角色数
	CompanyRoles  int64 `json:"company_roles"`  // 公司角色数
}

// ==========================
// 菜单管理相关模型
// ==========================

// MenuCreateRequest 菜单创建请求
type MenuCreateRequest struct {
	ParentID       string `json:"parent_id"`                                                // 父菜单ID，根菜单为空
	MenuName       string `json:"menu_name" binding:"required,min=1,max=50"`                // 菜单名称
	MenuType       string `json:"menu_type" binding:"required,oneof=directory menu button"` // 菜单类型
	RoutePath      string `json:"route_path"`                                               // 路由路径
	Component      string `json:"component"`                                                // 组件路径
	PermissionCode string `json:"permission_code"`                                          // 权限标识符
	Icon           string `json:"icon"`                                                     // 菜单图标
	SortOrder      int    `json:"sort_order"`                                               // 排序号
	Visible        bool   `json:"visible"`                                                  // 是否在菜单中显示
	Status         string `json:"status" binding:"omitempty,oneof=enable disable"`          // 状态
}

// MenuUpdateRequest 菜单更新请求
type MenuUpdateRequest struct {
	ParentID       string `json:"parent_id"`                                                 // 父菜单ID
	MenuName       string `json:"menu_name" binding:"omitempty,min=1,max=50"`                // 菜单名称
	MenuType       string `json:"menu_type" binding:"omitempty,oneof=directory menu button"` // 菜单类型
	RoutePath      string `json:"route_path"`                                                // 路由路径
	Component      string `json:"component"`                                                 // 组件路径
	PermissionCode string `json:"permission_code"`                                           // 权限标识符
	Icon           string `json:"icon"`                                                      // 菜单图标
	SortOrder      int    `json:"sort_order"`                                                // 排序号
	Visible        bool   `json:"visible"`                                                   // 是否在菜单中显示
	Status         string `json:"status" binding:"omitempty,oneof=enable disable"`           // 状态
}

// MenuInfo 菜单信息
type MenuInfo struct {
	ID             string     `json:"id"`              // 菜单ID
	MenuID         string     `json:"menu_id"`         // 菜单唯一标识
	ParentID       string     `json:"parent_id"`       // 父菜单ID
	MenuName       string     `json:"menu_name"`       // 菜单名称
	MenuType       string     `json:"menu_type"`       // 菜单类型
	RoutePath      string     `json:"route_path"`      // 路由路径
	Component      string     `json:"component"`       // 组件路径
	PermissionCode string     `json:"permission_code"` // 权限标识符
	Icon           string     `json:"icon"`            // 菜单图标
	SortOrder      int        `json:"sort_order"`      // 排序号
	Visible        bool       `json:"visible"`         // 是否在菜单中显示
	Status         string     `json:"status"`          // 状态
	CreatedAt      time.Time  `json:"created_at"`      // 创建时间
	UpdatedAt      time.Time  `json:"updated_at"`      // 更新时间
	Children       []MenuInfo `json:"children"`        // 子菜单列表
}

// MenuTreeNode 菜单树节点
type MenuTreeNode struct {
	MenuID         string         `json:"menu_id"`         // 菜单ID
	MenuName       string         `json:"menu_name"`       // 菜单名称
	MenuType       string         `json:"menu_type"`       // 菜单类型
	RoutePath      string         `json:"route_path"`      // 路由路径
	Component      string         `json:"component"`       // 组件路径
	PermissionCode string         `json:"permission_code"` // 权限标识符
	Icon           string         `json:"icon"`            // 菜单图标
	SortOrder      int            `json:"sort_order"`      // 排序号
	Visible        bool           `json:"visible"`         // 是否显示
	Status         string         `json:"status"`          // 状态
	Children       []MenuTreeNode `json:"children"`        // 子节点
}

// MenuQueryRequest 菜单查询请求
type MenuQueryRequest struct {
	MenuName       string `form:"menu_name"`       // 菜单名称搜索
	MenuType       string `form:"menu_type"`       // 菜单类型筛选
	Status         string `form:"status"`          // 状态筛选
	Visible        *bool  `form:"visible"`         // 是否显示筛选
	PermissionCode string `form:"permission_code"` // 权限标识符搜索
}

// MenuListResponse 菜单列表响应
type MenuListResponse struct {
	Menus []MenuInfo `json:"menus"` // 菜单列表（树形结构）
	Total int64      `json:"total"` // 总记录数
}

// BatchUpdateMenuStatusRequest 批量更新菜单状态请求
type BatchUpdateMenuStatusRequest struct {
	MenuIDs []string `json:"menu_ids" binding:"required,min=1"`              // 菜单ID数组
	Status  string   `json:"status" binding:"required,oneof=enable disable"` // 目标状态
}

// MenuStatsResponse 菜单统计响应
type MenuStatsResponse struct {
	TotalMenus     int64 `json:"total_menus"`     // 总菜单数
	EnabledMenus   int64 `json:"enabled_menus"`   // 启用菜单数
	DisabledMenus  int64 `json:"disabled_menus"`  // 禁用菜单数
	DirectoryMenus int64 `json:"directory_menus"` // 目录类型菜单数
	PageMenus      int64 `json:"page_menus"`      // 页面类型菜单数
	ButtonMenus    int64 `json:"button_menus"`    // 按钮类型菜单数
}

// UserMenuResponse 用户菜单响应（用于前端菜单渲染）
type UserMenuResponse struct {
	MenuID    string             `json:"menu_id"`    // 菜单ID
	MenuName  string             `json:"menu_name"`  // 菜单名称
	RoutePath string             `json:"route_path"` // 路由路径
	Component string             `json:"component"`  // 组件路径
	Icon      string             `json:"icon"`       // 菜单图标
	SortOrder int                `json:"sort_order"` // 排序号
	Children  []UserMenuResponse `json:"children"`   // 子菜单
}

// ==========================
// RBAC中间表模型
// ==========================

// UserRole 用户角色关联表模型
type UserRole struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`      // MongoDB主键ID
	UserID    string             `bson:"user_id" json:"user_id"`       // 用户ID
	RoleID    string             `bson:"role_id" json:"role_id"`       // 角色ID
	CreatedAt time.Time          `bson:"created_at" json:"created_at"` // 创建时间
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"` // 更新时间
}

// RolePermission 角色权限关联表模型
type RolePermission struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`                // MongoDB主键ID
	RoleID         string             `bson:"role_id" json:"role_id"`                 // 角色ID
	MenuID         string             `bson:"menu_id" json:"menu_id"`                 // 菜单ID（权限ID）
	PermissionType string             `bson:"permission_type" json:"permission_type"` // 权限类型：menu=菜单权限, button=按钮权限
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`           // 创建时间
	UpdatedAt      time.Time          `bson:"updated_at" json:"updated_at"`           // 更新时间
}

// ==========================
// 系统配置相关模型
// ==========================

// SystemConfig 系统配置表模型
type SystemConfig struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ConfigID    string             `bson:"config_id" json:"config_id"`       // 配置项ID
	ConfigType  string             `bson:"config_type" json:"config_type"`   // 配置类型：hk_manager, referral_branch, partner
	ConfigKey   string             `bson:"config_key" json:"config_key"`     // 配置键
	ConfigValue string             `bson:"config_value" json:"config_value"` // 配置值
	DisplayName string             `bson:"display_name" json:"display_name"` // 显示名称
	CompanyID   string             `bson:"company_id" json:"company_id"`     // 所属公司ID
	SortOrder   int                `bson:"sort_order" json:"sort_order"`     // 排序
	Status      string             `bson:"status" json:"status"`             // 状态：enable/disable
	Remark      string             `bson:"remark" json:"remark"`             // 备注
	CreatedBy   string             `bson:"created_by" json:"created_by"`     // 创建人
	UpdatedBy   string             `bson:"updated_by" json:"updated_by"`     // 更新人
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`     // 创建时间
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`     // 更新时间
}

// SystemConfigCreateRequest 创建系统配置请求
type SystemConfigCreateRequest struct {
	ConfigType  string `json:"config_type" binding:"required,oneof=hk_manager referral_branch partner" label:"配置类型"`
	ConfigKey   string `json:"config_key" binding:"required" label:"配置键"`
	ConfigValue string `json:"config_value" binding:"required" label:"配置值"`
	DisplayName string `json:"display_name" binding:"required" label:"显示名称"`
	SortOrder   int    `json:"sort_order" label:"排序"`
	Status      string `json:"status" binding:"oneof=enable disable" label:"状态"`
	Remark      string `json:"remark" label:"备注"`
}

// SystemConfigUpdateRequest 更新系统配置请求
type SystemConfigUpdateRequest struct {
	ConfigValue string `json:"config_value" label:"配置值"`
	DisplayName string `json:"display_name" label:"显示名称"`
	SortOrder   int    `json:"sort_order" label:"排序"`
	Status      string `json:"status" binding:"omitempty,oneof=enable disable" label:"状态"`
	Remark      string `json:"remark" label:"备注"`
}

// SystemConfigQueryRequest 查询系统配置请求
type SystemConfigQueryRequest struct {
	ConfigType string `json:"config_type" form:"config_type"`
	Status     string `json:"status" form:"status"`
	Keyword    string `json:"keyword" form:"keyword"`
	Page       int    `json:"page" form:"page"`
	PageSize   int    `json:"page_size" form:"page_size"`
}

// SystemConfigListResponse 系统配置列表响应
type SystemConfigListResponse struct {
	List  []SystemConfig `json:"list"`
	Total int64          `json:"total"`
}

// SystemConfigResponse 系统配置响应
type SystemConfigResponse struct {
	ID          string `json:"id"`
	ConfigID    string `json:"config_id"`
	ConfigType  string `json:"config_type"`
	ConfigKey   string `json:"config_key"`
	ConfigValue string `json:"config_value"`
	DisplayName string `json:"display_name"`
	CompanyID   string `json:"company_id"`
	SortOrder   int    `json:"sort_order"`
	Status      string `json:"status"`
	StatusText  string `json:"status_text"`
	Remark      string `json:"remark"`
	CreatedBy   string `json:"created_by"`
	UpdatedBy   string `json:"updated_by"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
