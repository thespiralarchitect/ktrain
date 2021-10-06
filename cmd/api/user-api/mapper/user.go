package mapper

import (
	"ktrain/cmd/api/user-api/dto"
	"ktrain/cmd/api/user-api/resquest"
	"ktrain/cmd/model"
	"time"

	"gorm.io/gorm"
)

func ToUserResponse(user *model.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:       user.Id,
		Fullname: user.Fullname,
		Username: user.Username,
		Gender:   user.Gender,
		Birthday: user.Birthday,
	}
}
func ToUserResquest(user *resquest.UserResquest) *model.User {
	pReq := &model.User{
		Id:          user.Id,
		Fullname:    user.Fullname,
		Username:    user.Username,
		Gender:      user.Gender,
		Birthday:    user.Birthday,
		CreatedAt:   time.Time{},
		UpdateAt:    time.Time{},
		DeletedAt:   gorm.DeletedAt{},
		Auth_tokens: make([]*model.Auth_token, 0),
	}
	for _, v := range user.Auth_tokens {
		pReq.Auth_tokens = append(pReq.Auth_tokens, &model.Auth_token{
			Id:        v.Id,
			UserID:    user.Id,
			Tocken:    v.Tocken,
			CreatedAt: time.Time{},
			UpdateAt:  time.Time{},
			DeletedAt: gorm.DeletedAt{},
		})
	}
	return pReq
}
