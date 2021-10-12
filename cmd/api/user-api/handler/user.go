package handler

import (
	"ktrain/cmd/api/user-api/dto"
	"ktrain/cmd/api/user-api/mapper"
	"ktrain/cmd/model"
	"ktrain/cmd/repository"
	"ktrain/pkg/errors"
	"ktrain/pkg/httputil"
	"ktrain/proto/pb"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type userHandler struct {
	userClient            pb.UserDMSServiceClient
	activityLogRepository repository.ActivityLogRepository
}

func NewUserHandler(userClient pb.UserDMSServiceClient, activityLogRepository repository.ActivityLogRepository) *userHandler {
	return &userHandler{
		userClient:            userClient,
		activityLogRepository: activityLogRepository,
	}
}
func (h *userHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	var validate *validator.Validate
	validate = validator.New()
	req := dto.UserRequest{}
	var binder httputil.JsonBinder
	if err := binder.BindRequest(&req, r); err != nil {
		if err.Error() == "Error reading body request" {
			httputil.RespondError(w, http.StatusInternalServerError, "Error reading body request")
			return
		} else {
			httputil.RespondError(w, http.StatusInternalServerError, "Error unmarshal body request")
		}
		return
	}
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
	// _, err = h.userRepository.GetUserByID(int64(id))
	pbReq := &pb.GetUserByIDRequest{
		Id: int64(id),
	}
	ppUser, err := h.userClient.GetUserByID(r.Context(),pbReq)
	if err != nil {
		if errors.IsDataNotFound(err) {
			httputil.RespondError(w, http.StatusNotFound, "User not found in database")
			return
		}
		httputil.RespondError(w, http.StatusInternalServerError, "Error when getting user ")
		return
	}
	// user := mapper.ToUserModel(&req)
	// user.ID = int64(id)
	// resp, err := h.userRepository.UpdateUser(user)
	resp, err := h.userClient.UpdateUser(r.Context(),(*pb.UpdateUserRequest)(ppUser))
	if err != nil {
		httputil.RespondError(w, http.StatusInternalServerError, "Error when update user")
		return
	}
	userResponse := &model.User{
		ID:         resp.User.Id,
		Fullname:   resp.User.Fullname,
		Username:   resp.User.Username,
		Gender:     resp.User.Gender,
		Birthday:   resp.User.Birthday.AsTime(),
	}
	httputil.RespondSuccessWithData(w, http.StatusOK, mapper.ToUserResponse(userResponse))
}

func (h *userHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	_, err := h.activityLogRepository.CreateAction(r.Context(), ctx.Value("userID").(int64), "Delete user")
	if err != nil {
		httputil.RespondError(w, http.StatusInternalServerError, "Error when creating new action")
		return
	}
	ID, _ := strconv.Atoi(chi.URLParam(r, "id"))
	// err = h.userRepository.DeleteUser(int64(ID))
	del := &pb.DeleteUserRequest{
		Id: int64(ID),
	}
	_, err = h.userClient.DeleteUser(r.Context(), del)
	if err != nil {
		httputil.RespondError(w, http.StatusInternalServerError, "Error when delete user")
		return
	}
}

func (h *userHandler) GetMyProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	_, err := h.activityLogRepository.CreateAction(r.Context(), ctx.Value("userID").(int64), "Get my profile user")
	if err != nil {
		httputil.RespondError(w, http.StatusInternalServerError, "Error when creating new action ")
		return
	}
	// user, err := h.userRepository.GetUserByID(ctx.Value("userID").(int64))
	req := &pb.GetUserByIDRequest{
		Id: ctx.Value("userID").(int64),
	}
	user, err := h.userClient.GetUserByID(r.Context(),req)
	if err != nil {
		if errors.IsDataNotFound(err) {
			httputil.RespondError(w, http.StatusNotFound, "Your profile not found")
			return
		}
		httputil.RespondError(w, http.StatusInternalServerError, "Error when getting user profile")
		return
	}
	userResponse := &model.User{
		ID:         user.User.Id,
		Fullname:   user.User.Fullname,
		Username:   user.User.Username,
		Gender:     user.User.Gender,
		Birthday:   user.User.Birthday.AsTime(),
	}
	httputil.RespondSuccessWithData(w, http.StatusOK, mapper.ToUserResponse(userResponse))
}

func (h *userHandler) GetListUsers(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	var ids []int64
	if values["ids"] != nil {
		req := dto.UserQuery{}
		var binder httputil.QueryURLBinder
		if err := binder.BindRequest(&req, r); err != nil {
			httputil.RespondError(w, http.StatusInternalServerError, "Error when query users list")
			return
		}
		for _, v := range req.Ids {
			id, _ := strconv.Atoi(v)
			ids = append(ids, int64(id))
		}
	}
	// ctx := r.Context()
	// _, err := h.activityLogRepository.CreateAction(r.Context(), ctx.Value("userID").(int64), "Get list user")
	// if err != nil {
	// 	httputil.RespondError(w, http.StatusInternalServerError, "Error when creating action ")
	// 	return
	// }
	// users, err := h.userRepository.GetListUser(ids)
	userIds := &pb.GetListUserRequest{
		Ids: ids,
	}
	users, err := h.userClient.GetListUser(r.Context(),userIds)
	if err != nil {
		httputil.RespondError(w, http.StatusInternalServerError, "Error when getting users list")
		return
	}
	var usersResponse []*model.User
	for _, v := range users.Users {
		userRes := &model.User{
			ID:         v.Id,
			Fullname:   v.Fullname,
			Username:   v.Username,
			Gender:     v.Gender,
			Birthday:   v.Birthday.AsTime(),
		}
		usersResponse = append(usersResponse, userRes)
	}
	httputil.RespondSuccessWithData(w, http.StatusOK, mapper.ToListUsersResponse(usersResponse))
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

	// user, err := h.userRepository.GetUserByID(int64(userID))
	userId := &pb.GetUserByIDRequest{
		Id: int64(userID),
	}
	user, err := h.userClient.GetUserByID(r.Context(),userId)
	if err != nil {
		if errors.IsDataNotFound(err) {
			httputil.RespondError(w, http.StatusNotFound, "User not found")
			return
		}
		httputil.RespondError(w, http.StatusInternalServerError, "Error when getting user profile")
		return
	}
	userResponse := &model.User{
		ID:         user.User.Id,
		Fullname:   user.User.Fullname,
		Username:   user.User.Username,
		Gender:     user.User.Gender,
		Birthday:   user.User.Birthday.AsTime(),
	}
	httputil.RespondSuccessWithData(w, http.StatusOK, mapper.ToUserResponse(userResponse))
}

func (h *userHandler) PostNewUser(w http.ResponseWriter, r *http.Request) {
	u := dto.CreateUserRequest{}
	var binder httputil.JsonBinder
	if err := binder.BindRequest(&u, r); err != nil {
		if err.Error() == "Error reading body request" {
			httputil.RespondError(w, http.StatusInternalServerError, "Error reading body request")
			return
		} else {
			httputil.RespondError(w, http.StatusInternalServerError, "Error unmarshal body request")
		}
		return
	}
	validate := validator.New()
	err := validate.Struct(u)
	if err != nil {
		httputil.RespondError(w, http.StatusBadRequest, "Error when validate request")
		return
	}
	birthday, _ := time.Parse("2006-01-02", u.Birthday)
	// User := &model.User{
	// 	Fullname: u.Fullname,
	// 	Username: u.Username,
	// 	Gender:   u.Gender,
	// 	Birthday: birthday,
	// }
	ctx := r.Context()
	_, err = h.activityLogRepository.CreateAction(r.Context(), ctx.Value("userID").(int64), "Create new user ")
	if err != nil {
		httputil.RespondError(w, http.StatusInternalServerError, "Error when creating action ")
		return
	}
	// newUser, err := h.userRepository.CreateUser(User)
	user := &pb.CreateUserRequest{
		User: &pb.User{
			Fullname:  u.Fullname,
			Username:  u.Username,
			Gender:    u.Gender,
			Birthday:  &timestamppb.Timestamp{
				Seconds: birthday.Unix(),
			},

		},
	}
	newUser, err := h.userClient.CreateUser(r.Context(), user)
	if err != nil {
		httputil.RespondError(w, http.StatusInternalServerError, "Error when creating new user")
		return
	}
	userResponse := &model.User{
		ID:         newUser.User.Id,
		Fullname:   newUser.User.Fullname,
		Username:   newUser.User.Username,
		Gender:     newUser.User.Gender,
		Birthday:   newUser.User.Birthday.AsTime(),
	}
	httputil.RespondSuccessWithData(w, http.StatusOK, mapper.ToUserResponse(userResponse))
}
