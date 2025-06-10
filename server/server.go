package server

import (
	"context"

	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/apetsko/guml/handlers"
	"github.com/go-chi/chi/v5"
	"golang.org/x/sync/errgroup"
)

func Run(addr string, logger *slog.Logger) (*http.Server, error) {
	srv := &http.Server{
		Addr:              addr,
		Handler:           router(),
		ReadHeaderTimeout: 3 * time.Second,
	}

	g, ctx := errgroup.WithContext(context.Background())

	logger.Info("Running server on " + addr)
	g.Go(func() error {
		<-ctx.Done()
		five := 5 * time.Second
		shutdownCtx, cancel := context.WithTimeout(context.Background(), five)
		defer cancel()
		return srv.Shutdown(shutdownCtx)
	})

	g.Go(func() error {
		return srv.ListenAndServe()
	})

	go func() {
		if err := g.Wait(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("HTTP server error", err)
		}
	}()

	return srv, nil
}

func router() http.Handler {
	r := chi.NewRouter()
	r.Route("/uml", func(r chi.Router) {
		r.Get("/index", handlers.Index)
		r.Post("/upload", handlers.Upload)
		r.Get("/", handlers.Link)
	})
	return r
}
