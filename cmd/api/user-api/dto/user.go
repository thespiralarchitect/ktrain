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
	Password string `json:"password" validate:"required"`
}

type UserRequest struct {
	Fullname string `json:"fullname" validate:"required"`
	Username string `json:"username" validate:"required"`
	Gender   string `json:"gender" validate:"required"`
	Birthday string `json:"birthday" validate:"required"`
}

type UserQuery struct {
	Ids []string
}

type ActionRequest struct {
	ID     int64
	Action string
}
type ActionResponse struct {
	Action []string `json:"action"`
}
type UserActivityLogMessage struct {
	ID  int64  `json:"id"`
	Log string `json:"log"`
}
