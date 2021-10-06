package dto

type UserResponse struct {
	ID       int64  `json:"id"`
	FullName string `json:"fullname"`
	Username string `json:"username"`
	Gender   string `json:"gender"`
	Birthday string `json:"birthday"`
}
