package resquest

type UserResquest struct {
	Id        int64        `json:"id `
	Fullname  string       `json:"fullname" validate:"required"`
	Username  string       `json:"username" validate:"required"`
	Gender    string       `json:"gender" validate:"required"`
	Birthday  string       `json:"birthday" validate:"required"`
	AuthToken []*AuthToken `json:"authtokens,omitempty"`
}
type AuthToken struct {
	Id     int64  `json:"id"`
	Tocken string `json:"tocken"`
}
type UserIdRequest struct {
	Id int64 `json:"id "`
}
