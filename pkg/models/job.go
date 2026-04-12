package models

import "time"

// AnalysisJob is a queued or in-flight analysis unit. Persistence details live in store implementations.
type AnalysisJob struct {
	ID         string          `json:"analysis_id"`
	Status     AnalysisStatus  `json:"status"`
	Request    AnalyzeRequest  `json:"request"`
	CreatedAt  time.Time       `json:"created_at"`
	StartedAt  *time.Time      `json:"started_at,omitempty"`
	FinishedAt *time.Time      `json:"finished_at,omitempty"`
	Report     *AnalysisReport `json:"report,omitempty"`
	Error      string          `json:"error,omitempty"`
}
