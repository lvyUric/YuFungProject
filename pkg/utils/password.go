package utils

import (
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword 加密密码
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// CheckPassword 验证密码
func CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// ValidatePassword 验证密码强度
func ValidatePassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	// 至少包含一个数字
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	// 至少包含一个小写字母
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	// 至少包含一个大写字母或特殊字符
	hasUpperOrSpecial := regexp.MustCompile(`[A-Z]|[^a-zA-Z0-9]`).MatchString(password)

	return hasNumber && hasLower && hasUpperOrSpecial
}
