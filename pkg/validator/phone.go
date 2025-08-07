package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

// PhonePattern 全球电话号码正则表达式
// 支持格式：
// - +86 138-1234-5678
// - +1 (555) 123-4567
// - 0755-12345678
// - 13812345678
// - +852 2345 6789
// - (010) 1234-5678
const PhonePattern = `^(\+\d{1,3}[\s\-]?)?(\(?\d{1,4}\)?[\s\-]?)?[\d\s\-\(\)\.]{6,18}$`

var phoneRegex = regexp.MustCompile(PhonePattern)

// ValidatePhone 电话号码验证函数
func ValidatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()

	// 基本长度检查 (去除所有非数字字符后应该在7-15位之间)
	digitOnly := regexp.MustCompile(`\d`).FindAllString(phone, -1)
	if len(digitOnly) < 7 || len(digitOnly) > 15 {
		return false
	}

	// 正则表达式验证
	return phoneRegex.MatchString(phone)
}
