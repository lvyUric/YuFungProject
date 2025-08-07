package model

// 响应状态码
const (
	// 通用状态码
	CodeSuccess       = 200 // 成功
	CodeInvalidParams = 400 // 请求参数错误
	CodeUnauthorized  = 401 // 未授权
	CodeForbidden     = 403 // 禁止访问
	CodeNotFound      = 404 // 资源不存在
	CodeConflict      = 409 // 资源冲突
	CodeAccountLocked = 423 // 账户被锁定
	CodeServerError   = 500 // 服务器内部错误

	// 认证相关错误码
	CodeAuthFailed      = 1001 // 认证失败
	CodeTokenInvalid    = 1002 // 令牌无效
	CodeTokenExpired    = 1003 // 令牌过期
	CodeAccountDisabled = 1004 // 账户被禁用

	// 用户相关错误码
	CodeUserExists    = 2001 // 用户已存在
	CodeUserNotExists = 2002 // 用户不存在
	CodeEmailExists   = 2003 // 邮箱已存在
	CodePhoneExists   = 2004 // 手机号已存在

	// 公司相关错误码
	CodeCompanyExists     = 3001 // 公司已存在
	CodeCompanyNotExists  = 3002 // 公司不存在
	CodeCompanyExpired    = 3003 // 公司已过期
	CodeUserQuotaExceeded = 3004 // 用户配额已满

	// 角色权限相关错误码
	CodeRoleExists     = 4001 // 角色已存在
	CodeRoleNotExists  = 4002 // 角色不存在
	CodePermissionDeny = 4003 // 权限不足

	// 业务相关错误码
	CodePolicyExists    = 5001 // 保单已存在
	CodePolicyNotExists = 5002 // 保单不存在
	CodeTableNotExists  = 5003 // 表结构不存在
	CodeFieldNotExists  = 5004 // 字段不存在
)

// 响应消息
const (
	MsgSuccess = "操作成功"
	MsgError   = "操作失败"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`               // 状态码
	Message string      `json:"message"`            // 响应消息
	Data    interface{} `json:"data,omitempty"`     // 响应数据
	TraceID string      `json:"trace_id,omitempty"` // 追踪ID
}

// SuccessResponse 成功响应
func SuccessResponse(message string, data interface{}) *Response {
	if message == "" {
		message = MsgSuccess
	}
	return &Response{
		Code:    CodeSuccess,
		Message: message,
		Data:    data,
	}
}

// ErrorResponse 错误响应
func ErrorResponse(code int, message string, data interface{}) *Response {
	if message == "" {
		message = MsgError
	}
	return &Response{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

// Success 成功响应（简化版）
func Success(data interface{}) *Response {
	return SuccessResponse(MsgSuccess, data)
}

// Error 错误响应（简化版）
func Error(code int, message string) *Response {
	return ErrorResponse(code, message, nil)
}

// ValidationError 参数验证错误响应
func ValidationError(err error) *Response {
	return ErrorResponse(CodeInvalidParams, "参数验证失败", err.Error())
}

// ServerError 服务器错误响应
func ServerError(message string) *Response {
	return ErrorResponse(CodeServerError, message, nil)
}

// UnauthorizedError 未授权错误响应
func UnauthorizedError(message string) *Response {
	if message == "" {
		message = "未授权访问"
	}
	return ErrorResponse(CodeUnauthorized, message, nil)
}

// ForbiddenError 禁止访问错误响应
func ForbiddenError(message string) *Response {
	if message == "" {
		message = "权限不足"
	}
	return ErrorResponse(CodeForbidden, message, nil)
}

// NotFoundError 资源不存在错误响应
func NotFoundError(message string) *Response {
	if message == "" {
		message = "资源不存在"
	}
	return ErrorResponse(CodeNotFound, message, nil)
}

// ConflictError 资源冲突错误响应
func ConflictError(message string) *Response {
	if message == "" {
		message = "资源冲突"
	}
	return ErrorResponse(CodeConflict, message, nil)
}
