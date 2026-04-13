package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/graphsentinel/graphsentinel/internal/store"
)

// NewRouter returns the HTTP handler tree for the GraphSentinel API.
func NewRouter(js store.JobStore) http.Handler {
	h := &handler{store: js}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	r.Get("/health", health)
	r.Post("/analyze", h.analyze)

	return r
}
