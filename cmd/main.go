package main

import (
	"context"

	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/apetsko/guml/config"
	"github.com/apetsko/guml/server"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cancel()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	c, err := config.New()
	if err != nil {
		logger.Error(err.Error())
		return
	}

	if _, err := server.Run(c.Host, logger); err != nil {
		log.Fatal("HTTP server failed: " + err.Error())
	}

	<-ctx.Done()
	logger.Info("Shutting down servers...")

}
