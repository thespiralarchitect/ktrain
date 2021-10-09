package dto

import "time"

type UserResponse struct {
	ID        int64     `json:"id"`
	Fullname  string    `json:"fullname"`
	Username  string    `json:"username"`
	Gender    string    `json:"gender"`
	Birthday  string    `json:"birthday"`
	CreatedAt time.Time `json:"created"`
	UpdatedAt time.Time `json:"updated"`
}

type CreateUserRequest struct {
	Fullname string `json:"fullname" validate:"required"`
	Username string `json:"username" validate:"required"`
	Gender   string `json:"gender" validate:"required"`
	Birthday string `json:"birthday" validate:"required"`
}

type UserResquest struct {
	Id       int64  `json:"id"`
	Fullname string `json:"fullname" validate:"required"`
	Username string `json:"username" validate:"required"`
	Gender   string `json:"gender" validate:"required"`
	Birthday string `json:"birthday" validate:"required"`
}

type UserQuery struct {
	Id []string
}