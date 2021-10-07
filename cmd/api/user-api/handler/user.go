package handler

import (
	"encoding/json"
	"io/ioutil"
	"ktrain/cmd/api/user-api/dto"
	"ktrain/cmd/api/user-api/mapper"
	"ktrain/cmd/model"
	"ktrain/cmd/repository"
	"ktrain/pkg/errors"
	"ktrain/pkg/httputil"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
)

type userHandler struct {
	userRepository repository.IUserRepository
}

func NewUserHandler(userRepository repository.IUserRepository) *userHandler {
	return &userHandler{
		userRepository: userRepository,
	}
}
func (h *userHandler) readBodyRequest(w http.ResponseWriter, r *http.Request, u *dto.CreateUserRequest) bool {
	var validate *validator.Validate
	validate = validator.New()
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		httputil.RespondError(w, http.StatusInternalServerError, "Error read body request")
		return false
	}
	err = json.Unmarshal(b, &u)
	if err != nil {
		httputil.RespondError(w, http.StatusInternalServerError, "Error unmarshal body request")
		return false
	}

	err = validate.Struct(u)
	if err != nil {
		httputil.RespondError(w, http.StatusBadRequest, "Validation error")
		return false
	}
	return true
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

func (h *userHandler) GetListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userRepository.GetListUser()
	if err != nil {
		httputil.RespondError(w, http.StatusInternalServerError, "Error when getting users list")
		return
	}
	httputil.RespondSuccessWithData(w, http.StatusOK, mapper.ToListUsersResponse(users))
}

func (h *userHandler) GetInformationUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		httputil.RespondError(w, http.StatusForbidden, "Error happend with user id ")
		return
	}
	user, err := h.userRepository.GetUserByID(int64(userID))
	if err != nil {
		if errors.IsDataNotFound(err) {
			httputil.RespondError(w, http.StatusNotFound, "User not found")
			return
		}
		httputil.RespondError(w, http.StatusInternalServerError, "Error when getting user profile")
		return
	}
	httputil.RespondSuccessWithData(w, http.StatusOK, mapper.ToUserResponse(user))
}

func (h *userHandler) PostNewUser(w http.ResponseWriter, r *http.Request) {
	var u dto.CreateUserRequest
	if ok := h.readBodyRequest(w, r, &u); !ok {
		return
	}
	birthday, _ := time.Parse("2006-01-02", u.Birthday)
	User := &model.User{
		Fullname: u.Fullname,
		Username: u.Username,
		Gender:   u.Gender,
		Birthday: birthday,
	}
	newUser, err := h.userRepository.CreateUser(User)
	if err != nil {
		httputil.RespondError(w, http.StatusInternalServerError, "Error when creating new user")
		return
	}
	httputil.RespondSuccessWithData(w, http.StatusOK, mapper.ToUserResponse(newUser))
}
