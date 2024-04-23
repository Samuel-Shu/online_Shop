package forms

type PasswordLoginForm struct {
	Email     string `form:"email" json:"email" binding:"required,email"`
	Password  string `form:"password" json:"password" binding:"required,min=6,max=20"`
	Captcha   string `form:"captcha" json:"captcha" binding:"required,min=5,max=5"`
	CaptchaId string `form:"captcha_id" json:"captcha_id" binding:"required"`
}

type RegisterForm struct {
	Email    string `form:"email" json:"email" binding:"required,email"`
	Password string `form:"password" json:"password" binding:"required,min=6,max=20"`
	Code     string `form:"code" json:"code" binding:"required,min=6,max=6"`
}
