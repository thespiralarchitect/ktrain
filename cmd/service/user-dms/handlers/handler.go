package handlers

import (
	"context"
	"database/sql"
	"ktrain/cmd/model"
	"ktrain/cmd/repository"
	"ktrain/proto/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserHandler struct {
	pb.UnimplementedUserDMSServiceServer
	userRepository repository.IUserRepository
}

func NewUserHandler( userRepository repository.IUserRepository)(*UserHandler,error){
	return &UserHandler{
		userRepository:        userRepository,
	},nil
}
func (h *UserHandler) GetUserByID(ctx context.Context, in *pb.GetUserByIDRequest)(*pb.GetUserByIDResponse, error){
	user, err := h.userRepository.GetUserByID(in.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, err
	}
	uResponse := &pb.GetUserByIDResponse{
		User: &pb.User{
			Id:        user.ID,
			Fullname:  user.Fullname,
			Username:  user.Username,
			Gender:    user.Gender,
			Birthday:  &timestamppb.Timestamp{
				Seconds: user.Birthday.Unix(),
			},
		},
	}
	return uResponse,nil
}
func(h *UserHandler) GetAuthToken(ctx context.Context, in *pb.GetAuthTokenRequest)(*pb.GetAuthTokenResponse, error){
	auth, err := h.userRepository.GetAuthToken(in.Token)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, err
	}
	uResponse := &pb.GetAuthTokenResponse{
		AuthToken: &pb.AuthToken{
			Id:     auth.ID,
			UserId: auth.UserID,
			Token:  auth.Token,
		},
	}
	return uResponse, nil
}
func(h *UserHandler) UpdateUser(ctx context.Context, in *pb.UpdateUserRequest)(*pb.UpdateUserResponse, error){
	uReq := &model.User{
		ID:         in.User.Id,
		Fullname:   in.User.Fullname,
		Username:   in.User.Username,
		Gender:     in.User.Gender,
		Birthday:   in.User.Birthday.AsTime(),
	}
	upUser, err := h.userRepository.UpdateUser(uReq)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, err
	}
	uResponse := &pb.UpdateUserResponse{
		User: &pb.User{
			Id:        upUser.ID,
			Fullname:  upUser.Fullname,
			Username:  upUser.Username,
			Gender:    upUser.Gender,
			Birthday:  &timestamppb.Timestamp{
				Seconds: upUser.Birthday.Unix(),
			},
		},
	}
	return uResponse, nil
}
func(h *UserHandler) DeleteUser(ctx context.Context, in *pb.DeleteUserRequest)(*pb.Empty,error)  {
	if err:= h.userRepository.DeleteUser(in.Id); err != nil {
		if err == sql.ErrNoRows{
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, err
	}
	return &pb.Empty{}, nil
}
func(h *UserHandler) GetListUser(ctx context.Context, in *pb.GetListUserRequest)(*pb.GetListUserResponse, error){
	users, err := h.userRepository.GetListUser(in.Ids)
	if err != nil {
		if err == sql.ErrNoRows{
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, err
	}
	var uRes []*pb.User
	for _, v := range users {
		user := &pb.User{
			Id:        v.ID,
			Fullname:  v.Fullname,
			Username:  v.Username,
			Gender:    v.Gender,
			Birthday:  &timestamppb.Timestamp{
				Seconds: v.Birthday.Unix(),
			},
		}
		uRes = append(uRes,user)
	}
	uRespose := &pb.GetListUserResponse{
		Users: uRes,
	}
	return uRespose, nil
}
func(h *UserHandler) CreateUser(ctx context.Context, in *pb.CreateUserRequest)(*pb.CreateUserResponse, error){
	userReq := &model.User{
		Fullname:   in.User.Fullname,
		Username:   in.User.Username,
		Gender:     in.User.Gender,
		Birthday:   in.User.Birthday.AsTime(),
	}
	newUser, err := h.userRepository.CreateUser(userReq)
	if err != nil{
		return nil, err
	}
	uRespose := &pb.CreateUserResponse{
		User: &pb.User{
			Id:        newUser.ID,
			Fullname:  newUser.Fullname,
			Username:  newUser.Username,
			Gender:    newUser.Gender,
			Birthday:  &timestamppb.Timestamp{
				Seconds: newUser.Birthday.Unix(),
			},
		},
	}
	return uRespose, nil
}