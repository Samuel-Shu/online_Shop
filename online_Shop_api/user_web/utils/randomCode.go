package utils

import (
	"math/rand"
	"time"
)

// RandomCode 以时间为种子，生成随机六位数验证码，用于邮箱认证
func RandomCode() string {
	rand.Seed(time.Now().UnixNano()) // 设置随机数种子为当前时间纳秒值

	var code = make([]byte, 6)
	for i := range code {
		code[i] = byte(rand.Intn(10) + '0') // 生成0到9之间的随机整数并转换为字符
	}

	return string(code)
}
