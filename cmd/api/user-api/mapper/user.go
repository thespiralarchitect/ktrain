package mapper

import (
	"ktrain/cmd/api/user-api/dto"
	"ktrain/cmd/model"
	"time"
)

func ToUserResponse(user *model.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:       user.ID,
		Fullname: user.Fullname,
		Username: user.Username,
		Gender:   user.Gender,
		Birthday: user.Birthday.Format("02/01/2006"),
	}
}
func ToUserModel(user *dto.UserResquest) *model.User {
	Birthday, _ := time.Parse("02/01/2006", user.Birthday)
	pReq := &model.User{
		ID:         user.Id,
		Fullname:   user.Fullname,
		Username:   user.Username,
		Gender:     user.Gender,
		Birthday:   Birthday,
		AuthTokens: []model.AuthToken{},
		CreatedAt:  time.Time{},
		UpdatedAt:  time.Time{},
	}
	return pReq
}
