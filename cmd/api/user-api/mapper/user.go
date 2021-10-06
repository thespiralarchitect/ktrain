package mapper

import (
	"ktrain/cmd/api/user-api/dto"
	"ktrain/cmd/api/user-api/resquest"
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
func ToUserResquest(user *resquest.UserResquest) *model.User {
	layout := "02/01/2006"
	Birthday, _ := time.Parse(layout, user.Birthday)
	pReq := &model.User{
		ID:         user.Id,
		Fullname:   user.Fullname,
		Username:   user.Username,
		Gender:     user.Gender,
		Birthday:   Birthday,
		AuthTokens: make([]model.AuthToken, 0),
		CreatedAt:  time.Time{},
		UpdatedAt:  time.Time{},
	}
	for _, v := range user.AuthToken {
		pReq.AuthTokens = append(pReq.AuthTokens, model.AuthToken{
			ID:        v.Id,
			UserID:    user.Id,
			Token:     v.Tocken,
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
		})
	}
	return pReq
}
