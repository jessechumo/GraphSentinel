package config

import (
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds process-wide settings loaded from the environment.
type Config struct {
	HTTPAddr        string
	ShutdownTimeout time.Duration
	WorkerCount     int
	WorkerQueueSize int
	LogLevel        slog.Level

	ReadHeaderTimeout time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
}

// Load reads configuration from environment variables with sensible defaults.
func Load() Config {
	cfg := Config{
		HTTPAddr:          getenv("HTTP_ADDR", ":8080"),
		ShutdownTimeout:   15 * time.Second,
		WorkerCount:       2,
		WorkerQueueSize:   256,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
		LogLevel:          slog.LevelInfo,
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
	if v := os.Getenv("WORKER_QUEUE_SIZE"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			cfg.WorkerQueueSize = n
		}
	}
	if v := os.Getenv("READ_TIMEOUT_SEC"); v != "" {
		if sec, err := strconv.Atoi(v); err == nil && sec > 0 {
			cfg.ReadTimeout = time.Duration(sec) * time.Second
		}
	}
	if v := os.Getenv("WRITE_TIMEOUT_SEC"); v != "" {
		if sec, err := strconv.Atoi(v); err == nil && sec > 0 {
			cfg.WriteTimeout = time.Duration(sec) * time.Second
		}
	}
	if v := os.Getenv("IDLE_TIMEOUT_SEC"); v != "" {
		if sec, err := strconv.Atoi(v); err == nil && sec > 0 {
			cfg.IdleTimeout = time.Duration(sec) * time.Second
		}
	}
	if v := os.Getenv("LOG_LEVEL"); v != "" {
		if parsed, ok := parseLogLevel(v); ok {
			cfg.LogLevel = parsed
		}
	}
	return cfg
}

func parseLogLevel(s string) (slog.Level, bool) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "debug":
		return slog.LevelDebug, true
	case "info", "":
		return slog.LevelInfo, true
	case "warn", "warning":
		return slog.LevelWarn, true
	case "error":
		return slog.LevelError, true
	default:
		return slog.LevelInfo, false
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
