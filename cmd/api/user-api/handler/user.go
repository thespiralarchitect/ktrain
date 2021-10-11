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
	userRepository        repository.IUserRepository
	activityLogRepository repository.ActivityLogRepository
}

func NewUserHandler(userRepository repository.IUserRepository, activityLogRepository repository.ActivityLogRepository) *userHandler {
	return &userHandler{
		userRepository:        userRepository,
		activityLogRepository: activityLogRepository,
	}
}
func (h *userHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	var validate *validator.Validate
	validate = validator.New()
	req := dto.UserRequest{}
	json.NewDecoder(r.Body).Decode(&req)
	err := validate.Struct(req)
	if err != nil {
		httputil.RespondError(w, http.StatusBadRequest, "Error when validate request")
		return
	}
	ctx := r.Context()
	_, err = h.activityLogRepository.CreateAction(r.Context(), ctx.Value("userID").(int64), "Update user")
	if err != nil {
		httputil.RespondError(w, http.StatusInternalServerError, "Error when creating new action ")
		return
	}
	_, err = h.userRepository.GetUserByID(int64(id))
	if err != nil {
		if errors.IsDataNotFound(err) {
			httputil.RespondError(w, http.StatusNotFound, "User not found in database")
			return
		}
		httputil.RespondError(w, http.StatusInternalServerError, "Error when getting user ")
		return
	}
	user := mapper.ToUserModel(&req)
	user.ID = int64(id)
	resp, err := h.userRepository.UpdateUser(user)
	if err != nil {
		httputil.RespondError(w, http.StatusInternalServerError, "Error when update user")
		return
	}
	httputil.RespondSuccessWithData(w, http.StatusOK, mapper.ToUserResponse(resp))
}
func (h *userHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	_, err := h.activityLogRepository.CreateAction(r.Context(), ctx.Value("userID").(int64), "Delete user")
	if err != nil {
		httputil.RespondError(w, http.StatusInternalServerError, "Error when creating new action")
		return
	}
	ID, _ := strconv.Atoi(chi.URLParam(r, "id"))
	err = h.userRepository.DeleteUser(int64(ID))
	if err != nil {
		httputil.RespondError(w, http.StatusInternalServerError, "Error when delete user")
		return
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
	_, err := h.activityLogRepository.CreateAction(r.Context(), ctx.Value("userID").(int64), "Get my profile user")
	if err != nil {
		httputil.RespondError(w, http.StatusInternalServerError, "Error when creating new action ")
		return
	}
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
	ctx := r.Context()
	_, err := h.activityLogRepository.CreateAction(r.Context(), ctx.Value("userID").(int64), "Get list user")
	if err != nil {
		httputil.RespondError(w, http.StatusInternalServerError, "Error when creating action ")
		return
	}
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
		httputil.RespondError(w, http.StatusBadRequest, "Request invalid")
		return
	}
	ctx := r.Context()
	_, err = h.activityLogRepository.CreateAction(r.Context(), ctx.Value("userID").(int64), "Get infor user")
	if err != nil {
		httputil.RespondError(w, http.StatusInternalServerError, "Error when creating action ")
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
	ctx := r.Context()
	_, err := h.activityLogRepository.CreateAction(r.Context(), ctx.Value("userID").(int64), "Create new user ")
	if err != nil {
		httputil.RespondError(w, http.StatusInternalServerError, "Error when creating action ")
		return
	}
	newUser, err := h.userRepository.CreateUser(User)
	if err != nil {
		httputil.RespondError(w, http.StatusInternalServerError, "Error when creating new user")
		return
	}
	httputil.RespondSuccessWithData(w, http.StatusOK, mapper.ToUserResponse(newUser))
}
