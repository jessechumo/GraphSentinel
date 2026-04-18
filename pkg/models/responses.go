package models

// Stable API error codes returned in ErrorResponse.Code.
const (
	ErrCodeInvalidJSON       = "INVALID_JSON"
	ErrCodeEmptyBody         = "EMPTY_BODY"
	ErrCodeTrailingJSON      = "TRAILING_JSON"
	ErrCodeBodyTooLarge      = "BODY_TOO_LARGE"
	ErrCodeValidation        = "VALIDATION_ERROR"
	ErrCodeEnqueueFailed     = "ENQUEUE_FAILED"
	ErrCodeMissingAnalysisID = "MISSING_ANALYSIS_ID"
	ErrCodeNotFound          = "NOT_FOUND"
)

// SubmitAnalysisResponse is returned immediately after a successful POST /analyze.
type SubmitAnalysisResponse struct {
	Status     AnalysisStatus `json:"status"`
	AnalysisID string         `json:"analysis_id"`
}

// ErrorResponse is the standard JSON shape for API errors.
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
}
