package handler

import (
	"ktrain/cmd/api/user-api/mapper"
	"ktrain/cmd/repository"
	"ktrain/pkg/httputil"
	"net/http"
)

type userHandler struct {
	userRepository repository.IUserRepository
}

func NewUserHandler(userRepository repository.IUserRepository) *userHandler {
	return &userHandler{
		userRepository: userRepository,
	}
}

func (h *userHandler) GetMyProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, err := h.userRepository.GetUserByID(ctx.Value("userID").(int64))
	if err != nil {
		httputil.RespondError(w, http.StatusInternalServerError, "Error when getting user profile")
		return
	}

	httputil.RespondSuccessWithData(w, http.StatusOK, mapper.ToUserResponse(user))
}
