package main

import (
	"context"
	"flag"
	"fmt"
	"ktrain/cmd/api/user-api/handler"
	middleware2 "ktrain/cmd/api/user-api/middleware"
	"ktrain/rambbitmq"

	"ktrain/cmd/repository"
	"ktrain/pkg/config"
	"ktrain/pkg/httputil"
	"ktrain/pkg/storage"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/spf13/viper"
)

var (
	configPath = flag.String("config.file", "cmd/api/user-api/config.yaml", "Path to configuration file.")
)

func main() {
	// parse command-line flags
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

	psqlDB, err := storage.NewPSQLManager()
	if err != nil {
		log.Fatalf("Error when connecting database, err: %v", err)
		return
	}
	defer psqlDB.Close()
	rabbitMq, err := rambbitmq.ConectRambbitMQ()
	defer rabbitMq.Close()
	if err != nil {
		httputil.FailOnError(err, "Failed to connect to RabbitMQ")
	}
	rambbitmqRepository := repository.NewActivityLogRambbitMqRepository(rabbitMq)
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	r.Route("/api", func(r chi.Router) {
		r.Use(middleware.SetHeader("Content-Type", "application/json"))
		userRepository := repository.NewUserRepository(psqlDB)
		activityLogRepository := repository.NewActivityLogRepository(mongDB)
		//Authenticate
		r.Use(middleware2.NewDBTokenAuth(userRepository).Handle())
		//API handlers
		userHandler := handler.NewUserHandler(rambbitmqRepository, userRepository, activityLogRepository)
		monngoHandler := handler.NewActivityLogHandler(activityLogRepository)
		r.Get("/users/{id}/activities", monngoHandler.GetActivity)
		r.Get("/me", userHandler.GetMyProfile)
		r.Get("/users", userHandler.GetListUsers)
		r.Get("/users/{id}", userHandler.GetInformationUser)
		r.Route("/", func(r chi.Router) {
			r.Use(middleware2.NewDBTokenAuth(userRepository).HandleAdmin())
			r.Post("/users", userHandler.PostNewUser)
			r.Put("/users/{id}", userHandler.UpdateUser)
			r.Delete("/users/{id}", userHandler.DeleteUser)
		})
	})
	fmt.Println("Listen at port: 8080")
	http.ListenAndServe(":8080", r)
}
