package response

type UserResponse struct {
	Id       int32  `json:"id"`
	Email   string `json:"email"`
	NickName string `json:"nickName"`
	Birthday string `json:"birthday"`
	Gender   string `json:"gender"`
	Role     uint32 `json:"role"`
}
