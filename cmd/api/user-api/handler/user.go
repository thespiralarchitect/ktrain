package handler

import (
	"ktrain/cmd/api/user-api/dto"
	"ktrain/cmd/api/user-api/mapper"
	"ktrain/cmd/model"
	"ktrain/pkg/errors"
	"ktrain/pkg/httputil"
	"ktrain/pkg/logger"
	"ktrain/proto/pb"
	"ktrain/rambbitmq"
	"net/http"
	"strconv"
	"time"

	middleware2 "ktrain/cmd/api/user-api/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type userHandler struct {
	userClient pb.UserDMSServiceClient
	rabbitmq   *rambbitmq.RabbitMqManager
	validator  *validator.Validate
}

func NewUserHandler(rabbitmq *rambbitmq.RabbitMqManager, userClient pb.UserDMSServiceClient) *userHandler {
	return &userHandler{
		userClient: userClient,
		rabbitmq:   rabbitmq,
		validator:  validator.New(),
	}
}
func (h *userHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	req := dto.UserRequest{}
	var binder httputil.JsonBinder
	if err := binder.BindRequest(&req, r); err != nil {
		if err.Error() == "Error reading body request" {
			logger.InitLogger().Errorw("Error reading body request", "error", err)
			//h.logger.Errorw("Error reading body request", "error", err)
			httputil.RespondError(w, http.StatusInternalServerError, "Error reading body request")
			return
		} else {
			logger.InitLogger().Errorw("Error unmarshal body request", "error", err)
			//h.logger.Errorw("Error unmarshal body request", "error", err)
			httputil.RespondError(w, http.StatusInternalServerError, "Error unmarshal body request")
		}
		return
	}
	err := h.validator.Struct(req)
	if err != nil {
		logger.InitLogger().Errorw("Error when validate request", "error", err)
		//h.logger.Errorw("Error when validate request", "error", err)
		httputil.RespondError(w, http.StatusBadRequest, "Error when validate request")
		return
	}
	ctx := r.Context()
	body := dto.UserActivityLogMessage{
		ID:  ctx.Value(middleware2.ContextKey("userID")).(int64),
		Log: "Update user",
	}

	err = h.rabbitmq.Publish(body)
	if err != nil {
		httputil.FailOnError(err, err.Error())
		return
	}
	pbReq := &pb.GetUserByIDRequest{
		Id: int64(id),
	}
	ppUser, err := h.userClient.GetUserByID(r.Context(), pbReq)
	if err != nil {
		if errors.IsDataNotFound(err) {
			logger.InitLogger().Errorw("User not found in database", "error", err)
			//h.logger.Errorw("User not found in database", "error", err)
			httputil.RespondError(w, http.StatusNotFound, "User not found in database")
			return
		}
		logger.InitLogger().Errorw("Error when getting user", "error", err)
		//h.logger.Errorw("Error when getting user", "error", err)
		httputil.RespondError(w, http.StatusInternalServerError, "Error when getting user ")
		return
	}
	resp, err := h.userClient.UpdateUser(r.Context(), (*pb.UpdateUserRequest)(ppUser))
	if err != nil {
		logger.InitLogger().Errorw("Error when update user", "error", err)
		//h.logger.Errorw("Error when update user", "error", err)
		httputil.RespondError(w, http.StatusInternalServerError, "Error when update user")
		return
	}
	userResponse := &model.User{
		ID:       resp.User.Id,
		Fullname: resp.User.Fullname,
		Username: resp.User.Username,
		Gender:   resp.User.Gender,
		Birthday: resp.User.Birthday.AsTime(),
	}
	httputil.RespondSuccessWithData(w, http.StatusOK, mapper.ToUserResponse(userResponse))
}

func (h *userHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	body := dto.UserActivityLogMessage{
		ID:  ctx.Value(middleware2.ContextKey("userID")).(int64),
		Log: "Delete user",
	}
	err := h.rabbitmq.Publish(body)
	if err != nil {
		httputil.FailOnError(err, err.Error())
		return
	}
	ID, _ := strconv.Atoi(chi.URLParam(r, "id"))
	del := &pb.DeleteUserRequest{
		Id: int64(ID),
	}
	_, err = h.userClient.DeleteUser(r.Context(), del)
	if err != nil {
		logger.InitLogger().Errorw("Error when delete user", "error", err)
		//h.logger.Errorw("Error when delete user", "error", err)
		httputil.RespondError(w, http.StatusInternalServerError, "Error when delete user")
		return
	}
}

func (h *userHandler) GetMyProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	body := dto.UserActivityLogMessage{
		ID:  ctx.Value(middleware2.ContextKey("userID")).(int64),
		Log: "Get my profile ",
	}
	err := h.rabbitmq.Publish(body)
	if err != nil {
		httputil.FailOnError(err, err.Error())
		return
	}
	req := &pb.GetUserByIDRequest{
		Id: ctx.Value("userID").(int64),
	}
	user, err := h.userClient.GetUserByID(r.Context(), req)
	if err != nil {
		if errors.IsDataNotFound(err) {
			logger.InitLogger().Errorw("Your profile not found", "error", err)
			//h.logger.Errorw("Your profile not found", "error", err)
			httputil.RespondError(w, http.StatusNotFound, "Your profile not found")
			return
		}
		logger.InitLogger().Errorw("Error when getting user profile,err:%v", err)
		//h.logger.Errorw("Error when getting user profile,err:%v", err)
		httputil.RespondError(w, http.StatusInternalServerError, "Error when getting user profile")
		return
	}
	userResponse := &model.User{
		ID:       user.User.Id,
		Fullname: user.User.Fullname,
		Username: user.User.Username,
		Gender:   user.User.Gender,
		Birthday: user.User.Birthday.AsTime(),
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
			logger.InitLogger().Errorw("Error when query users list", "error", err)
			//h.logger.Errorw("Error when query users list", "error", err)
			httputil.RespondError(w, http.StatusInternalServerError, "Error when query users list")
			return
		}
		for _, v := range req.Ids {
			id, _ := strconv.Atoi(v)
			ids = append(ids, int64(id))
		}
	}
	ctx := r.Context()
	body := dto.UserActivityLogMessage{
		ID:  ctx.Value(middleware2.ContextKey("userID")).(int64),
		Log: "Get list user",
	}
	err := h.rabbitmq.Publish(body)
	if err != nil {
		httputil.FailOnError(err, err.Error())
		return
	}
	userIds := &pb.GetListUserRequest{
		Ids: ids,
	}
	users, err := h.userClient.GetListUser(r.Context(), userIds)
	if err != nil {
		logger.InitLogger().Errorw("Error when getting users list", "error", err)
		//h.logger.Errorw("Error when getting users list", "error", err)
		httputil.RespondError(w, http.StatusInternalServerError, "Error when getting users list")
		return
	}
	var usersResponse []*model.User
	for _, v := range users.Users {
		userRes := &model.User{
			ID:       v.Id,
			Fullname: v.Fullname,
			Username: v.Username,
			Gender:   v.Gender,
			Birthday: v.Birthday.AsTime(),
		}
		usersResponse = append(usersResponse, userRes)
	}
	httputil.RespondSuccessWithData(w, http.StatusOK, mapper.ToListUsersResponse(usersResponse))
}

func (h *userHandler) GetInformationUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		logger.InitLogger().Errorw("Request invalid", "error", err)
		httputil.RespondError(w, http.StatusBadRequest, "Request invalid")
		return
	}
	ctx := r.Context()
	body := dto.UserActivityLogMessage{
		ID:  ctx.Value(middleware2.ContextKey("userID")).(int64),
		Log: "Get infor user",
	}
	err = h.rabbitmq.Publish(body)
	if err != nil {
		httputil.FailOnError(err, err.Error())
		return
	}
	userId := &pb.GetUserByIDRequest{
		Id: int64(userID),
	}
	user, err := h.userClient.GetUserByID(r.Context(), userId)
	if err != nil {
		if errors.IsDataNotFound(err) {
			logger.InitLogger().Errorw("User not found,err:%v", err)
			//h.logger.Errorw("User not found,err:%v", err)
			httputil.RespondError(w, http.StatusNotFound, "User not found")
			return
		}
		logger.InitLogger().Errorw("Error when getting user profile", "error", err)
		//h.logger.Errorw("Error when getting user profile", "error", err)
		httputil.RespondError(w, http.StatusInternalServerError, "Error when getting user profile")
		return
	}
	userResponse := &model.User{
		ID:       user.User.Id,
		Fullname: user.User.Fullname,
		Username: user.User.Username,
		Gender:   user.User.Gender,
		Birthday: user.User.Birthday.AsTime(),
	}
	httputil.RespondSuccessWithData(w, http.StatusOK, mapper.ToUserResponse(userResponse))
}

func (h *userHandler) PostNewUser(w http.ResponseWriter, r *http.Request) {
	u := dto.CreateUserRequest{}
	var binder httputil.JsonBinder
	if err := binder.BindRequest(&u, r); err != nil {
		if err.Error() == "Error reading body request" {
			logger.InitLogger().Errorw("Error reading body request", "error", err)
			//h.logger.Errorw("Error reading body request", "error", err)
			httputil.RespondError(w, http.StatusInternalServerError, "Error reading body request")
			return
		} else {
			logger.InitLogger().Errorw("Error unmarshal body request", "error", err)
			//h.logger.Errorw("Error unmarshal body request", "error", err)
			httputil.RespondError(w, http.StatusInternalServerError, "Error unmarshal body request")
		}
		return
	}
	err := h.validator.Struct(u)
	if err != nil {
		logger.InitLogger().Errorw("Error when validate request", "error", err)
		//h.logger.Errorw("Error when validate request", "error", err)
		httputil.RespondError(w, http.StatusBadRequest, "Error when validate request")
		return
	}
	birthday, _ := time.Parse("2006-01-02", u.Birthday)
	ctx := r.Context()
	body := dto.UserActivityLogMessage{
		ID:  ctx.Value(middleware2.ContextKey("userID")).(int64),
		Log: "Create user",
	}
	err = h.rabbitmq.Publish(body)
	if err != nil {
		httputil.FailOnError(err, err.Error())
		return
	}
	user := &pb.CreateUserRequest{
		User: &pb.User{
			Fullname: u.Fullname,
			Username: u.Username,
			Gender:   u.Gender,
			Birthday: &timestamppb.Timestamp{
				Seconds: birthday.Unix(),
			},
		},
	}
	newUser, err := h.userClient.CreateUser(r.Context(), user)
	if err != nil {
		logger.InitLogger().Errorw("Error when creating new user", "error", err)
		//h.logger.Errorw("Error when creating new user", "error", err)
		httputil.RespondError(w, http.StatusInternalServerError, "Error when creating new user")
		return
	}
	userResponse := &model.User{
		ID:       newUser.User.Id,
		Fullname: newUser.User.Fullname,
		Username: newUser.User.Username,
		Gender:   newUser.User.Gender,
		Birthday: newUser.User.Birthday.AsTime(),
	}
	httputil.RespondSuccessWithData(w, http.StatusOK, mapper.ToUserResponse(userResponse))
}
