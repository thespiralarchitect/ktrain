package handler

import (
	"encoding/json"
	"ktrain/cmd/api/user-api/dto"
	"ktrain/cmd/api/user-api/mapper"
	"ktrain/cmd/repository"
	"ktrain/pkg/errors"
	"ktrain/pkg/httputil"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type userHandler struct {
	userRepository repository.IUserRepository
}

func NewUserHandler(userRepository repository.IUserRepository) *userHandler {
	return &userHandler{
		userRepository: userRepository,
	}
}
func (h *userHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var validate *validator.Validate
	validate = validator.New()
	req := dto.UserResquest{}
	json.NewDecoder(r.Body).Decode(&req)
	err := validate.Struct(req)
	if err != nil {
		httputil.RespondError(w, http.StatusBadRequest, "Error when validate request")
		return
	}
	_, err = h.userRepository.GetUserByID(req.Id)
	if err != nil {
		if errors.IsDataNotFound(err) {
			httputil.RespondError(w, http.StatusNotFound, "User not found in database")
			return
		}
		httputil.RespondError(w, http.StatusInternalServerError, "Error when getting user ")
		return
	}
	user := mapper.ToUserModel(&req)
	resp, err := h.userRepository.UpdateUser(user)
	if err != nil {
		httputil.RespondError(w, http.StatusInternalServerError, "Error when update user")
		return
	}
	httputil.RespondSuccessWithData(w, http.StatusOK, mapper.ToUserResponse(resp))
}
func (h *userHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ID, _ := strconv.Atoi(chi.URLParam(r, "id"))
	err := h.userRepository.DeleteUser(int64(ID))
	if err != nil {
		httputil.RespondError(w, http.StatusInternalServerError, "Error when delete user")
		return
	}
}
func (h *userHandler) GetMyProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, err := h.userRepository.GetUserByID(ctx.Value("userID").(int64))
	if err != nil {
		if errors.IsDataNotFound(err) {
			httputil.RespondError(w, http.StatusNotFound, "Your profile not found")
			return
		}
		httputil.RespondError(w, http.StatusInternalServerError, "Error when getting user profile")
		return
	}
	httputil.RespondSuccessWithData(w, http.StatusOK, mapper.ToUserResponse(user))
}
