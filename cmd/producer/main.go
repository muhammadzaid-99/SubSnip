package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/muhammadzaid-99/SubSnip/internal/queue"
	"github.com/muhammadzaid-99/SubSnip/internal/status"

	"github.com/muhammadzaid-99/SubSnip/internal/api"
)

func init() {
	status.Init()
}

func main() {

	r := chi.NewRouter()
	r.Post("/submit-task", api.SubmitTaskHandler)
	r.Get("/task-status/{taskID}", api.TaskStatusHandler)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Println("Server API running on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down...")
	queue.Close() // closing the connection

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(shutdownCtx)
}
