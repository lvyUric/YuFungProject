package validator

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// InitCustomValidators 初始化自定义验证器
func InitCustomValidators() error {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册电话号码验证器
		if err := v.RegisterValidation("phone_pattern", ValidatePhone); err != nil {
			return err
		}
	}
	return nil
}
