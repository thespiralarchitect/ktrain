package dto

import "time"

type UserResponse struct {
	ID       int64  `json:"id"`
	Fullname string `json:"fullname"`
	Username string `json:"username"`
	Gender   string `json:"gender"`
	Birthday string `json:"birthday"`
}

type InformationUserResponse struct {
	ID        int64     `json:"id"`
	Fullname  string    `json:"fullname"`
	Username  string    `json:"username"`
	Gender    string    `json:"gender"`
	Birthday  string    `json:"birthday"`
	CreatedAt time.Time `json:"created"`
	UpdatedAt time.Time `json:"updated"`
}
