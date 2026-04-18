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
			writeAPIError(w, r, http.StatusRequestEntityTooLarge, models.ErrCodeBodyTooLarge, "request body too large")
			return
		}
		if errors.Is(err, io.EOF) {
			writeAPIError(w, r, http.StatusBadRequest, models.ErrCodeEmptyBody, "empty body")
			return
		}
		writeAPIError(w, r, http.StatusBadRequest, models.ErrCodeInvalidJSON, "invalid JSON")
		return
	}

	if err := dec.Decode(&struct{}{}); err != io.EOF {
		writeAPIError(w, r, http.StatusBadRequest, models.ErrCodeTrailingJSON, "trailing JSON not allowed")
		return
	}

	if err := req.Validate(); err != nil {
		writeAPIError(w, r, http.StatusBadRequest, models.ErrCodeValidation, err.Error())
		return
	}

	norm := req.Normalized()
	job, err := h.store.Enqueue(norm)
	if err != nil {
		writeAPIError(w, r, http.StatusInternalServerError, models.ErrCodeEnqueueFailed, "could not queue analysis")
		return
	}

	id := job.ID
	if h.submit != nil {
		h.submit(id)
	}

	// Response status is always queued here; read job.Status only via Snapshot after enqueue.
	writeJSON(w, http.StatusAccepted, models.SubmitAnalysisResponse{
		Status:     models.StatusQueued,
		AnalysisID: id,
	})
}

func isMaxBytesError(err error) bool {
	var me *http.MaxBytesError
	return errors.As(err, &me)
}
