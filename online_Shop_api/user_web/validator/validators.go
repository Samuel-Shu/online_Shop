package validator

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

func ValidateEmail(f1 validator.FieldLevel) bool {
	email := f1.Field().String()
	//使用正则表达式验证手机号是否合法
	ok, _ := regexp.MatchString(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`, email)
	if !ok {
		return false
	}
	return true
}
