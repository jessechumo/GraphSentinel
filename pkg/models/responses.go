package models

// SubmitAnalysisResponse is returned immediately after a successful POST /analyze.
type SubmitAnalysisResponse struct {
	Status     AnalysisStatus `json:"status"`
	AnalysisID string         `json:"analysis_id"`
}

// ErrorResponse is a consistent JSON shape for client and validation errors (used by handlers later).
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
}
