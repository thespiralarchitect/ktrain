package resquest

type UserResquest struct {
	Id          int64         `json:"id `
	Fullname    string        `json:"fullname" validate:"required"`
	Username    string        `json:"username" validate:"required"`
	Gender      string        `json:"gender" validate:"required"`
	Birthday    string        `json:"birthday" validate:"required"`
	Auth_tokens []*Auth_token `json:"authtokens,omitempty"`
}
type Auth_token struct {
	Id     int64  `json:"id"`
	Tocken string `json:"tocken"`
}
type UserIdRequest struct {
	Id int64 `json:"id "`
}
