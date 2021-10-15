package main

import (
	"context"
	"flag"
	"fmt"
	"ktrain/cmd/repository"
	"ktrain/cmd/service/activity-log-dms/handlers"
	"ktrain/pkg/config"
	"ktrain/pkg/logger"
	"ktrain/pkg/storage"
	"ktrain/proto/pb"
	"net"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	configPath = flag.String("config.file", "cmd/service/activity-log-dms/config.yaml", "Path to configuration file.")
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
		return
	}
	listen, err := net.Listen("tcp", ":9001")
	if err != nil {
		logger.InitLogger().Panic(err)
	}
	s := grpc.NewServer()
	ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("mongodb.timeout"))
	defer cancel()
	mongDB, err := storage.NewMongoDBManager(ctx)
	if err != nil {
		logger.InitLogger().Fatalf("Error when connecting database, err: %v", err)
		return
	}
	defer mongDB.Close(ctx)
	activityLogRepository := repository.NewActivityLogRepository(mongDB)
	h, err := handlers.NewActivityLogHandler(activityLogRepository)
	if err != nil {
		logger.InitLogger().Fatalf("Error when creating new activity handler, err: %v", err)
		return
	}
	reflection.Register(s)
	pb.RegisterActivityLogDMSServiceServer(s, h)
	fmt.Println("listen port 9001")
	err = s.Serve(listen)
	if err != nil {
		logger.InitLogger().Panic(err)
	}
}
