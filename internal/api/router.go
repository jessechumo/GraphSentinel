package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/graphsentinel/graphsentinel/internal/store"
)

// NewRouter returns the HTTP handler tree for the GraphSentinel API.
// submit is called with a new job id after a successful enqueue; nil skips background scheduling (tests).
func NewRouter(js store.JobStore, submit func(id string)) http.Handler {
	h := &handler{store: js, submit: submit}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(requestLog())
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	r.Get("/health", health)
	r.Post("/analyze", h.analyze)
	r.Get("/analysis/{id}", h.getAnalysis)

	return r
}
