package dto

type UserResponse struct {
	ID       int64  `json:"id"`
	Fullname string `json:"fullname"`
	Username string `json:"username"`
	Gender   string `json:"gender"`
	Birthday string `json:"birthday"`
}

type UserResquest struct {
	Id       int64  `json:"id"`
	Fullname string `json:"fullname" validate:"required"`
	Username string `json:"username" validate:"required"`
	Gender   string `json:"gender" validate:"required"`
	Birthday string `json:"birthday" validate:"required"`
}
