package mapper

import (
	"ktrain/cmd/api/user-api/dto"
	"ktrain/cmd/model"
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
