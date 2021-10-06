package handler

import (
	"encoding/json"
	"fmt"
	"ktrain/cmd/api/user-api/mapper"
	"ktrain/cmd/api/user-api/resquest"
	"ktrain/cmd/repository"
	"ktrain/pkg/httputil"
	"net/http"
	"strconv"

	//"github.com/go-playground/validator"
	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

type userHandler struct {
	userRepository repository.IUserRepository
}

func NewUserHandler(userRepository repository.IUserRepository) *userHandler {
	return &userHandler{
		userRepository: userRepository,
	}
}
func (h *userHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	validate = validator.New()
	res := resquest.UserResquest{}
	json.NewDecoder(r.Body).Decode(&res)
	err := validate.Struct(res)
	if err != nil {
		httputil.RespondError(w, http.StatusBadRequest, "User does not exist")
		// if validateErrors, ok := err.(validator.ValidationErrors); ok {
		// 	errMessages := make([]string, 0)
		// 	for _, v := range validateErrors {
		// 		errMessages = append(errMessages, v.Error())
		// 	}
		// 	httputil.RespondJSONValidator(w, http.StatusBadRequest, errMessages)
		// }
		return
	}
	_, err = h.userRepository.GetUserByID(res.Id)
	if err != nil {
		httputil.RespondError(w, http.StatusInternalServerError, "User does not exist")
		return
	}
	user := mapper.ToUserResquest(&res)
	resp, err := h.userRepository.UpdateUser(user)
	if err != nil {
		httputil.RespondError(w, http.StatusInternalServerError, "Error when update user")
		return
	}
	httputil.RespondSuccessWithData(w, http.StatusOK, mapper.ToUserResponse(resp))
}
func (h *userHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// res := resquest.UserIdRequest{}
	// json.NewDecoder(r.Body).Decode(&res)
	ID, _ := strconv.Atoi(chi.URLParam(r, "id"))
	err := h.userRepository.DeleteUser(int64(ID))

	fmt.Println(ID)
	if err != nil {
		httputil.RespondError(w, http.StatusInternalServerError, "Error when delete user")
		return
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
