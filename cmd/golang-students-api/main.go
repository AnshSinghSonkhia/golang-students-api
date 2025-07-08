package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/AnshSinghSonkhia/golang-students-api/internal/config"
)

func main() {
	// TODO: load config

	cfg := config.MustLoad()

	// TODO: database setup

	// TODO: setup router

	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the Golang Students API!"))
	})

	// TODO: setup server

	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	fmt.Printf("Server is started and runnin %s\n", cfg.HTTPServer.Addr)

	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Failed to start server: %s", err.Error())
	}

}
