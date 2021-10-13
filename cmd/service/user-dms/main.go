package main

import (
	"flag"
	"fmt"
	"ktrain/cmd/repository"
	"ktrain/cmd/service/user-dms/handlers"
	"ktrain/pkg/config"
	"ktrain/pkg/storage"
	"ktrain/proto/pb"
	"log"
	"net"
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
	listen, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()

	psqlDB, err := storage.NewPSQLManager()
	if err != nil {
		log.Fatalf("Error when connecting database, err: %v", err)
		return
	}
	defer psqlDB.Close()
	userRepository := repository.NewUserRepository(psqlDB)
	h, err := handlers.NewUserHandler(userRepository)
	if err != nil {
		panic(err)
	}
	reflection.Register(s)
	pb.RegisterUserDMSServiceServer(s,h)
	fmt.Println("listen port 9000")
	s.Serve(listen)
}