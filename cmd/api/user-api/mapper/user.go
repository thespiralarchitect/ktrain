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
func ToListUsersResponse(users []*model.User) []*dto.UserResponse {
	listUsersResponse := []*dto.UserResponse{}
	for _, user := range users {
		userResponse := &dto.UserResponse{
			ID:       user.ID,
			Fullname: user.Fullname,
			Username: user.Username,
			Gender:   user.Gender,
			Birthday: user.Birthday.Format("02/01/2006"),
		}
		listUsersResponse = append(listUsersResponse, userResponse)
	}
	return listUsersResponse
}
func ToInformationUserResponse(user *model.User) *dto.InformationUserResponse {
	return &dto.InformationUserResponse{
		ID:        user.ID,
		Fullname:  user.Fullname,
		Username:  user.Username,
		Gender:    user.Gender,
		Birthday:  user.CreatedAt.Format("02/01/2006"),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
