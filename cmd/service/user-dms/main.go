package main

import (
	"flag"
	"fmt"
	"ktrain/cmd/repository"
	"ktrain/cmd/service/user-dms/handlers"
	"ktrain/pkg/config"
	"ktrain/pkg/logger"
	"ktrain/pkg/storage"
	"ktrain/proto/pb"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	configPath = flag.String("config.file", "cmd/service/user-dms/config.yaml", "Path to configuration file.")
)

func main() {
	// logger := middleware2.InitLogger()
	// defer func() {
	// 	if err := logger.Sync(); err != nil {
	// 		log.Fatalf("Error when release the buffer,err:%v", err)
	// 	}
	// }()
	flag.Parse()
	err := config.BindDefault(*configPath)
	if err != nil {
		logger.InitLogger().Fatalf("Error when binding config, err: %v", err)
		//logger.Fatalf("Error when binding config, err: %v", err)
		return
	}
	listen, err := net.Listen("tcp", ":9000")
	if err != nil {
		logger.InitLogger().Panic(err)
	}
	s := grpc.NewServer()

	psqlDB, err := storage.NewPSQLManager()
	if err != nil {
		logger.InitLogger().Fatalf("Error when connecting database, err: %v", err)
		//logger.Fatalf("Error when connecting database, err: %v", err)
		return
	}
	defer psqlDB.Close()
	userRepository := repository.NewUserRepository(psqlDB)
	h, err := handlers.NewUserHandler(userRepository)
	if err != nil {
		logger.InitLogger().Fatalf("Error when creating new user handler, err: %v", err)
		//logger.Fatalf("Error when creating new user handler, err: %v", err)
		return
	}
	reflection.Register(s)
	pb.RegisterUserDMSServiceServer(s, h)
	fmt.Println("listen port 9000")
	err = s.Serve(listen)
	if err != nil {
		logger.InitLogger().Panic(err)
	}
}
