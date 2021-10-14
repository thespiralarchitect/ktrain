package handlers

import (
	"context"
	"ktrain/cmd/repository"
	"ktrain/pkg/config"
	"ktrain/pkg/storage"
	"ktrain/proto/pb"
	"log"
	"testing"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func bindConfig() {
	err := config.BindDefault("./config.yaml")
	if err != nil {
		log.Fatalf("Error when binding config, err: %v", err)
		return
	}
}
func initTest() *storage.PSQLManager {
	bindConfig()
	psqlDB, err := storage.NewPSQLManager()
	if err != nil {
		log.Fatalf("Error when connecting database, err: %v", err)
		return nil
	}
	return psqlDB
}
func TestUserHandler_DeleteUser(t *testing.T) {
	psqlDB := initTest()
	userRepository := repository.NewUserRepository(psqlDB)
	h, err := NewUserHandler(userRepository)
	if err != nil {
		log.Fatalf("Error when creating new user handler, err: %v", err)
		return
	}
	_, err = h.DeleteUser(context.Background(), &pb.DeleteUserRequest{
		Id: 1,
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("success")
}
func TestUserHandler_ListUser(t *testing.T) {
	psqlDB := initTest()
	userRepository := repository.NewUserRepository(psqlDB)
	h, err := NewUserHandler(userRepository)
	if err != nil {
		log.Fatalf("Error when creating new user handler, err: %v", err)
		return
	}
	_, err = h.GetListUser(context.Background(), &pb.GetListUserRequest{
		Ids: []int64{1},
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("success")
}
func TestUserHandler_GetUser(t *testing.T) {
	psqlDB := initTest()
	userRepository := repository.NewUserRepository(psqlDB)
	h, err := NewUserHandler(userRepository)
	if err != nil {
		log.Fatalf("Error when creating new user handler, err: %v", err)
		return
	}
	_, err = h.GetUserByID(context.Background(), &pb.GetUserByIDRequest{
		Id: 1,
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("success")
}
func TestUserHandler_UpdateUser(t *testing.T) {
	psqlDB := initTest()
	userRepository := repository.NewUserRepository(psqlDB)
	h, err := NewUserHandler(userRepository)
	if err != nil {
		log.Fatalf("Error when creating new user handler, err: %v", err)
		return
	}
	_, err = h.UpdateUser(context.Background(), &pb.UpdateUserRequest{
		User: &pb.User{
			IsAdmin:   false,
			Id:        1,
			Fullname:  "",
			Username:  "",
			Gender:    "",
			Birthday:  &timestamppb.Timestamp{},
			AuthToken: []*pb.AuthToken{},
		},
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("success")
}
