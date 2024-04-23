package forms

type SendEmail struct {
	Email string `form:"email" json:"email" binding:"required,email"`
	Type string `form:"type" json:"type" binding:"required,oneof=register login"`
}
