package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/graphsentinel/graphsentinel/internal/analyzers"
	"github.com/graphsentinel/graphsentinel/internal/api"
	"github.com/graphsentinel/graphsentinel/internal/config"
	"github.com/graphsentinel/graphsentinel/internal/store"
	"github.com/graphsentinel/graphsentinel/internal/workers"
)

func main() {
	cfg := config.Load()
	jobs := store.NewMemory()

	pool := workers.NewPool(cfg.WorkerCount, 256, jobs, analyzers.Analyze)
	pool.Start()

	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           api.NewRouter(jobs, pool.Submit),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("graphsentinel listening on %s", cfg.HTTPAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server: %v", err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	log.Println("shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown: %v", err)
	}

	pool.Close()
	pool.Wait()
}
