package main

import (
	"context"
	"flag"
	handler "ktrain/cmd/consumers/activity-log-aggregator/handlers"
	"ktrain/pkg/config"
	"ktrain/pkg/httputil"
	"ktrain/pkg/logger"
	"ktrain/pkg/storage"
	"ktrain/proto/pb"
	"log"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var (
	configPath = flag.String("config.file", "cmd/consumers/activity-log-aggregator/config.yaml", "Path to configuration file.")
)

func main() {
	sugarLogger := logger.InitLogger()
	defer func() {
		if err := sugarLogger.Sync(); err != nil {
			log.Fatalf("Error when release the buffer,err:%v", err)
			return
		}
	}()
	flag.Parse()
	err := config.BindDefault(*configPath)
	if err != nil {
		logger.Log().Fatalf("Error when binding config, err: %v", err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("mongodb.timeout"))
	defer cancel()
	mongDB, err := storage.NewMongoDBManager(ctx)
	if err != nil {
		logger.Log().Fatalf("Error when connecting database, err: %v", err)
		return
	}
	defer mongDB.Close(ctx)
	activityConn, err := grpc.Dial(":9001", grpc.WithInsecure())
	if err != nil {
		logger.Log().Panic(err)
	}
	activityClient := pb.NewActivityLogDMSServiceClient(activityConn)
	rabbitMq, err := handler.ConectRambbitMQ()
	if err != nil {
		httputil.FailOnError(err, "Failed to connect to RabbitMQ")
		return
	}
	defer rabbitMq.Close()
	err = rabbitMq.Consumers(ctx, activityClient)
	if err != nil {
		httputil.FailOnError(err, err.Error())
		return
	}
}
