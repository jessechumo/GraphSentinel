package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds process-wide settings loaded from the environment.
type Config struct {
	HTTPAddr        string
	ShutdownTimeout time.Duration
	WorkerCount     int
}

// Load reads configuration from environment variables with sensible defaults.
func Load() Config {
	cfg := Config{
		HTTPAddr:        getenv("HTTP_ADDR", ":8080"),
		ShutdownTimeout: 15 * time.Second,
		WorkerCount:     2,
	}
	if v := os.Getenv("SHUTDOWN_TIMEOUT_SEC"); v != "" {
		if sec, err := strconv.Atoi(v); err == nil && sec > 0 {
			cfg.ShutdownTimeout = time.Duration(sec) * time.Second
		}
	}
	if v := os.Getenv("WORKER_COUNT"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			cfg.WorkerCount = n
		}
	}
	return cfg
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
