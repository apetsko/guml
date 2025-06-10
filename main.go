package main

import (
	"context"
	"d2/handlers"
	"d2/server"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cancel()

	logger := slog.NewJSONHandler(os.Stdout, nil)

	if _, err := server.Run(logger); err != nil {
		log.Fatal("HTTP server failed: " + err.Error())
	}

	<-ctx.Done()
	logger.Info("Shutting down servers...")

}
