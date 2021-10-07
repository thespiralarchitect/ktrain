package httpreq

import "time"

type UserRequest struct {
	Fullname string    `json:"fullname" validate:"required"`
	Username string    `json:"username" validate:"required"`
	Gender   string    `json:"gender" validate:"required"`
	Birthday time.Time `json:"birthday" validate:"required"`
}
