package main

import (
	"flag"
	"fmt"
	"ktrain/cmd/api/user-api/handler"
	"log"

	middleware2 "ktrain/cmd/api/user-api/middleware"
	"ktrain/proto/pb"
	"ktrain/rambbitmq"

	"ktrain/pkg/config"
	"ktrain/pkg/httputil"
	"ktrain/pkg/logger"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"google.golang.org/grpc"
)

var (
	configPath = flag.String("config.file", "cmd/api/user-api/config.yaml", "Path to configuration file.")
)

func main() {
	sugarLogger := logger.InitLogger()
	defer func() {
		if err := sugarLogger.Sync(); err != nil {
			log.Fatalf("Error when release the buffer,err:%v", err)
			return
		}
	}()
	// parse command-line flags
	flag.Parse()
	err := config.BindDefault(*configPath)
	if err != nil {
		logger.Log().Fatalf("Error when binding config, err: %v", err)
		return
	}
	userConn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		logger.Log().Panic(err)
	}
	rabbitMq, err := rambbitmq.ConectRambbitMQ()
	if err != nil {
		logger.Log().Fatalf("Failed to connect to RabbitMQ, err: %v", err)
		return
	}

	defer rabbitMq.Close()
	userClient := pb.NewUserDMSServiceClient(userConn)
	activityConn, err := grpc.Dial(":9001", grpc.WithInsecure())
	if err != nil {
		logger.Log().Panic(err)
	}
	activityClient := pb.NewActivityLogDMSServiceClient(activityConn)
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("welcome"))
		if err != nil {
			logger.Log().Errorw("Error to write response,err:%v", err)
			httputil.RespondError(w, http.StatusInternalServerError, err.Error())
		}
	})
	r.Route("/api", func(r chi.Router) {
		r.Use(middleware.SetHeader("Content-Type", "application/json"))
		//Authenticate
		r.Use(middleware2.NewDBTokenAuth(userClient).Handle())
		//API handlers
		userHandler := handler.NewUserHandler(rabbitMq, userClient)
		monngoHandler := handler.NewActivityLogHandler(activityClient)
		r.Get("/users/{id}/activities", monngoHandler.GetActivity)
		r.Get("/me", userHandler.GetMyProfile)
		r.Get("/users", userHandler.GetListUsers)
		r.Get("/users/{id}", userHandler.GetInformationUser)
		r.Route("/", func(r chi.Router) {
			r.Use(middleware2.NewDBTokenAuth(userClient).HandleAdmin())
			r.Post("/users", userHandler.PostNewUser)
			r.Put("/users/{id}", userHandler.UpdateUser)
			r.Delete("/users/{id}", userHandler.DeleteUser)
		})
	})
	fmt.Println("Listen at port: 8080")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		logger.Log().Fatalf("ListenAndServe(): %v", err)
		return
	}
}
