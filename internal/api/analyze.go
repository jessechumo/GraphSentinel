package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/graphsentinel/graphsentinel/pkg/models"
)

// POST body limit: code may be up to models' max plus JSON envelope and escaping headroom.
const maxAnalyzeBodyBytes = 512 << 10

func (h *handler) analyze(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxAnalyzeBodyBytes)

	var req models.AnalyzeRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		if isMaxBytesError(err) {
			writeAPIError(w, http.StatusRequestEntityTooLarge, "request body too large")
			return
		}
		if errors.Is(err, io.EOF) {
			writeAPIError(w, http.StatusBadRequest, "empty body")
			return
		}
		writeAPIError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	if err := dec.Decode(&struct{}{}); err != io.EOF {
		writeAPIError(w, http.StatusBadRequest, "trailing JSON not allowed")
		return
	}

	if err := req.Validate(); err != nil {
		writeAPIError(w, http.StatusBadRequest, err.Error())
		return
	}

	norm := req.Normalized()
	job, err := h.store.Enqueue(norm)
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, "could not queue analysis")
		return
	}

	writeJSON(w, http.StatusAccepted, models.SubmitAnalysisResponse{
		Status:     job.Status,
		AnalysisID: job.ID,
	})
}

func isMaxBytesError(err error) bool {
	var me *http.MaxBytesError
	return errors.As(err, &me)
}
