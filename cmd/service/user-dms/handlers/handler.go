package handlers

import (
	"context"
	"database/sql"
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