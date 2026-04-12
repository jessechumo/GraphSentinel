package models

// AnalysisStatus is the lifecycle state of an analysis job.
type AnalysisStatus string

const (
	StatusQueued    AnalysisStatus = "queued"
	StatusRunning   AnalysisStatus = "running"
	StatusCompleted AnalysisStatus = "completed"
	StatusFailed    AnalysisStatus = "failed"
)

// Terminal returns true if no further state transitions are expected.
func (s AnalysisStatus) Terminal() bool {
	return s == StatusCompleted || s == StatusFailed
}
