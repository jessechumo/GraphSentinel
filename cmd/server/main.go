package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/graphsentinel/graphsentinel/internal/analyzers"
	"github.com/graphsentinel/graphsentinel/internal/api"
	"github.com/graphsentinel/graphsentinel/internal/config"
	"github.com/graphsentinel/graphsentinel/internal/store"
	"github.com/graphsentinel/graphsentinel/internal/workers"
)

func main() {
	cfg := config.Load()

	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.LogLevel})
	slog.SetDefault(slog.New(logHandler))

	jobs := store.NewMemory()

	pool := workers.NewPool(cfg.WorkerCount, cfg.WorkerQueueSize, jobs, analyzers.Analyze)
	pool.Start()

	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           api.NewRouter(jobs, pool.Submit),
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		ReadTimeout:       cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
	}

	go func() {
		slog.Info("graphsentinel listening", slog.String("addr", cfg.HTTPAddr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server exited", slog.Any("err", err))
			os.Exit(1)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	slog.Info("shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Warn("graceful shutdown", slog.Any("err", err))
	}

	pool.Close()
	pool.Wait()
}
