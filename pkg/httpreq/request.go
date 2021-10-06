package httpreq

import "time"

type UserRequest struct {
	Fullname string    `json:"fullname"`
	Username string    `json:"username"`
	Gender   string    `json:"gender"`
	Birthday time.Time `json:"birthday"`
}
