package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AnshSinghSonkhia/golang-students-api/internal/config"
	"github.com/AnshSinghSonkhia/golang-students-api/internal/http/handlers/student"
	"github.com/AnshSinghSonkhia/golang-students-api/internal/storage/sqlite"
)

func main() {
	// load config

	cfg := config.MustLoad()

	// database setup

	storage, err := sqlite.New(cfg) // initialize the SQLite database with the configuration
	if err != nil {
		log.Fatalf("Failed to initialize storage: %s", err.Error()) // log an error if the storage initialization fails
	}

	// log the storage path
	slog.Info("Storage initialized", slog.String("storage_path", cfg.StoragePath), slog.String("env", cfg.Env), slog.String("version", "1.0.0"))

	// setup router

	router := http.NewServeMux()

	// register the student handler for POST requests to /api/students
	router.HandleFunc("POST /api/students", student.New(storage))

	// register the student handler for GET requests to /api/students/{id}
	router.HandleFunc("GET /api/students/{id}", student.GetByID(storage))

	// register the student handler for GET requests to /api/students
	router.HandleFunc("GET /api/students", student.GetList(storage))

	// setup server

	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	slog.Info("Server started", slog.String("address", cfg.Addr)) // log the server address

	fmt.Printf("Server is started and runnin %s\n", cfg.HTTPServer.Addr)

	// gracefully shutdown server on interrupt signal
	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM) // catch interrupt signals

	go func() { // run server in a goroutine
		err := server.ListenAndServe()
		if err != nil {
			log.Fatalf("Failed to start server: %s", err.Error())
		}
	}()

	// Wait for interrupt signal
	<-done // block until an interrupt signal is received

	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // create a context with a timeout for graceful shutdown

	defer cancel() // ensure the context is cancelled after use

	err = server.Shutdown(ctx) // shutdown the server gracefully
	if err != nil {
		slog.Error("Failed to shutdown server", slog.String("error", err.Error()))
	} // log any error that occurs during shutdown

	slog.Info("Server shutdown successfully")
}
