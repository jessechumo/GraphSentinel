package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/graphsentinel/graphsentinel/pkg/models"
)

func (h *handler) getAnalysis(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeAPIError(w, http.StatusBadRequest, "missing analysis id")
		return
	}

	job, ok := h.store.Snapshot(id)
	if !ok {
		writeAPIError(w, http.StatusNotFound, "analysis not found")
		return
	}

	writeJSON(w, http.StatusOK, models.NewGetAnalysisResponse(&job))
}
