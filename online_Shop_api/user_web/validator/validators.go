package validator

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

func ValidateMobile(f1 validator.FieldLevel) bool {
	mobile := f1.Field().String()
	//使用正则表达式验证手机号是否合法
	ok, _ := regexp.MatchString(`^(13[0-9]|14[579]|15[0-3,5-9]|16[6]|17[0135678]|18[0-9]|19[89])\d{8}$`, mobile)
	if !ok {
		return false
	}
	return true
}
