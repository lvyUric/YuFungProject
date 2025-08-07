package utils

import (
	"crypto/rand"
	"fmt"
	"time"
)

// GenerateUserID 生成用户ID
func GenerateUserID() string {
	// 使用时间戳和随机数生成用户ID
	timestamp := time.Now().Unix()
	randomBytes := make([]byte, 4)
	rand.Read(randomBytes)
	return fmt.Sprintf("USR%d%X", timestamp, randomBytes)
}

// GenerateCompanyID 生成公司ID
func GenerateCompanyID() string {
	timestamp := time.Now().Unix()
	randomBytes := make([]byte, 4)
	rand.Read(randomBytes)
	return fmt.Sprintf("CMP%d%X", timestamp, randomBytes)
}

// GenerateRoleID 生成角色ID
func GenerateRoleID() string {
	timestamp := time.Now().Unix()
	randomBytes := make([]byte, 4)
	rand.Read(randomBytes)
	return fmt.Sprintf("ROL%d%X", timestamp, randomBytes)
}

// GenerateID 通用ID生成函数
func GenerateID(prefix string) string {
	timestamp := time.Now().Unix()
	randomBytes := make([]byte, 4)
	rand.Read(randomBytes)

	// 转换前缀为大写
	switch prefix {
	case "user":
		return fmt.Sprintf("USR%d%X", timestamp, randomBytes)
	case "company":
		return fmt.Sprintf("CMP%d%X", timestamp, randomBytes)
	case "role":
		return fmt.Sprintf("ROL%d%X", timestamp, randomBytes)
	default:
		return fmt.Sprintf("%s%d%X", prefix, timestamp, randomBytes)
	}
}
