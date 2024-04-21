package forms

type PasswordLoginForm struct {
	Mobile string `form:"mobile" json:"mobile" binding:"required,mobile"`
	Password string `form:"password" json:"password" binding:"required,min=6,max=20"`
}
