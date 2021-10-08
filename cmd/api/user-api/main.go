package main

import (
	"flag"
	"fmt"
	"ktrain/cmd/api/user-api/handler"

	middleware2 "ktrain/cmd/api/user-api/middleware"
	"ktrain/cmd/repository"
	"ktrain/pkg/config"
	"ktrain/pkg/storage"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

	psqlDB, err := storage.NewPSQLManager()
	if err != nil {
		log.Fatalf("Error when connecting database, err: %v", err)
		return
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	r.Route("/api", func(r chi.Router) {
		r.Use(middleware.SetHeader("Content-Type", "application/json"))

		userRepository := repository.NewUserRepository(psqlDB)

		//Authenticate
		r.Use(middleware2.NewDBTokenAuth(userRepository).Handle())

		//API handlers
		userHandler := handler.NewUserHandler(userRepository)
		r.Get("/me", userHandler.GetMyProfile)
		r.Put("/users/{id}", userHandler.UpdateUser)
		r.Delete("/users/{id}", userHandler.DeleteUser)
		r.Get("/users", userHandler.GetListUsers)
		r.Get("/users/{id}", userHandler.GetInformationUser)
		r.Post("/users", userHandler.PostNewUser)

	})
	fmt.Println("Listen at port: 8080")
	http.ListenAndServe(":8080", r)

}
