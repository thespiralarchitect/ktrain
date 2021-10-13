package main

import (
	"context"
	"flag"
	handler "ktrain/cmd/consumers/activity-log-aggregator/handlers"
	"ktrain/cmd/repository"
	"ktrain/pkg/config"
	"ktrain/pkg/httputil"
	"ktrain/pkg/storage"
	"log"

	"github.com/spf13/viper"
)

var (
	configPath = flag.String("config.file", "cmd/consumers/activity-log-aggregator/config.yaml", "Path to configuration file.")
)

func main() {
	flag.Parse()
	err := config.BindDefault(*configPath)
	if err != nil {
		log.Fatalf("Error when binding config, err: %v", err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("mongodb.timeout"))
	defer cancel()
	mongDB, err := storage.NewMongoDBManager(ctx)
	if err != nil {
		log.Fatalf("Error when connecting database, err: %v", err)
		return
	}
	defer mongDB.Close(ctx)
	activityLogRepository := repository.NewActivityLogRepository(mongDB)
	rabbitMq, err := handler.ConectRambbitMQ(activityLogRepository)
	if err != nil {
		httputil.FailOnError(err, "Failed to connect to RabbitMQ")
	}
	defer rabbitMq.Close()
	err = rabbitMq.Consumers(ctx)
	if err != nil {
		httputil.FailOnError(err, err.Error())
	}
}
