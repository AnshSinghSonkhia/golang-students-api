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
)

func main() {
	// load config

	cfg := config.MustLoad()

	// TODO: database setup

	// setup router

	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the Golang Students API!"))
	})

	// setup server

	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	slog.Info("Server started %s", slog.String("address", cfg.Addr)) // log the server address

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

	err := server.Shutdown(ctx) // shutdown the server gracefully
	if err != nil {
		slog.Error("Failed to shutdown server", slog.String("error", err.Error()))
	} // log any error that occurs during shutdown

	slog.Info("Server shutdown successfully")
}
