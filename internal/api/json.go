package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/graphsentinel/graphsentinel/pkg/models"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(true)
	_ = enc.Encode(v)
}

// writeAPIError writes a standard ErrorResponse and logs server-side.
func writeAPIError(w http.ResponseWriter, r *http.Request, status int, code, msg string) {
	if r != nil {
		logAPIError(r, status, code, msg)
	}
	writeJSON(w, status, models.ErrorResponse{Error: msg, Code: code})
}

func logAPIError(r *http.Request, status int, code, msg string) {
	lvl := slog.LevelInfo
	if status >= 500 {
		lvl = slog.LevelError
	}
	slog.LogAttrs(r.Context(), lvl, "api_error",
		slog.Int("http_status", status),
		slog.String("err_code", code),
		slog.String("err_msg", msg),
		slog.String("request_id", middleware.GetReqID(r.Context())),
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
	)
}
