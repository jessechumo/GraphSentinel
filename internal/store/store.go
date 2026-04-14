// Package store abstracts persistence for analysis jobs and results.
package store

import "github.com/graphsentinel/graphsentinel/pkg/models"

// JobStore persists queued analysis work. Implementations must be safe for concurrent use.
type JobStore interface {
	Enqueue(req models.AnalyzeRequest) (*models.AnalysisJob, error)
	Get(id string) (*models.AnalysisJob, bool)
}

// WorkerStore is the store surface used by background workers (claim + finish).
type WorkerStore interface {
	JobStore
	TryClaim(id string) bool
	Complete(id string, report *models.AnalysisReport) error
	Fail(id string, msg string) error
}
