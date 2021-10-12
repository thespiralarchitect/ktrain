package main

import (
	"context"
	"flag"
	"fmt"
	"ktrain/cmd/repository"
	"ktrain/cmd/service/user-dms/handlers"
	"ktrain/pkg/config"
	"ktrain/pkg/storage"
	"ktrain/proto/pb"
	"log"
	"net"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	configPath = flag.String("config.file", "cmd/api/user-api/config.yaml", "Path to configuration file.")
)
func main() {
	flag.Parse()
	err := config.BindDefault(*configPath)
	if err != nil {
		log.Fatalf("Error when binding config, err: %v", err)
		return
	}
	userConn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	userClient := pb.NewUserDMSServiceClient(userConn)
	listen, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()

	ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("mongodb.timeout"))
	defer cancel()
	mongDB, err := storage.NewMongoDBManager(ctx)
	if err != nil {
		log.Fatalf("Error when connecting database, err: %v", err)
		return
	}
	defer mongDB.Close(ctx)
	psqlDB, err := storage.NewPSQLManager()
	if err != nil {
		log.Fatalf("Error when connecting database, err: %v", err)
		return
	}
	defer psqlDB.Close()
	userRepository := repository.NewUserRepository(psqlDB)
	activityLogRepository := repository.NewActivityLogRepository(mongDB)
	h, err := handlers.NewUserHandler(userClient,userRepository,activityLogRepository)
	if err != nil {
		panic(err)
	}
	reflection.Register(s)
	pb.RegisterUserDMSServiceServer(s,h)
	fmt.Println("listen port 9000")
	s.Serve(listen)
}