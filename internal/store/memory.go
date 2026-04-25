package store

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/graphsentinel/graphsentinel/pkg/models"
)

var (
	// ErrJobNotFound is returned when mutating a job id that does not exist.
	ErrJobNotFound = errors.New("job not found")
	// ErrJobNotRunning is returned when completing or failing a job that is not running.
	ErrJobNotRunning = errors.New("job is not running")
)

// Memory is an in-process JobStore for development and tests.
type Memory struct {
	mu   sync.RWMutex
	jobs map[string]*models.AnalysisJob
}

// NewMemory returns an empty in-memory store.
func NewMemory() *Memory {
	return &Memory{jobs: make(map[string]*models.AnalysisJob)}
}

// Enqueue assigns an ID, records the job as queued, and returns the stored job pointer.
// The submission response must not read mutable fields from this pointer after scheduling workers.
func (m *Memory) Enqueue(req models.AnalyzeRequest) (*models.AnalysisJob, error) {
	id, err := newAnalysisID()
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	job := &models.AnalysisJob{
		ID:        id,
		Status:    models.StatusQueued,
		Request:   req,
		CreatedAt: now,
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	if m.jobs == nil {
		m.jobs = make(map[string]*models.AnalysisJob)
	}
	m.jobs[id] = job
	return job, nil
}

// Get returns a job by id if present.
func (m *Memory) Get(id string) (*models.AnalysisJob, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	j, ok := m.jobs[id]
	if !ok || j == nil {
		return nil, false
	}
	return j, true
}

// Snapshot returns a copy of the job suitable for API responses and offline processing.
func (m *Memory) Snapshot(id string) (models.AnalysisJob, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	j, ok := m.jobs[id]
	if !ok || j == nil {
		return models.AnalysisJob{}, false
	}
	return *j, true
}

// TryClaim transitions a queued job to running. It returns false if the job is missing or not queued.
func (m *Memory) TryClaim(id string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	j := m.jobs[id]
	if j == nil || j.Status != models.StatusQueued {
		return false
	}
	now := time.Now().UTC()
	j.Status = models.StatusRunning
	j.StartedAt = &now
	return true
}

// Complete marks a running job successful and attaches the report.
func (m *Memory) Complete(id string, report *models.AnalysisReport) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	j := m.jobs[id]
	if j == nil {
		return ErrJobNotFound
	}
	if j.Status != models.StatusRunning {
		return ErrJobNotRunning
	}
	now := time.Now().UTC()
	j.Status = models.StatusCompleted
	j.FinishedAt = &now
	j.Report = report
	j.Error = ""
	return nil
}

// Fail marks a running job as failed with a message.
func (m *Memory) Fail(id string, msg string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	j := m.jobs[id]
	if j == nil {
		return ErrJobNotFound
	}
	if j.Status != models.StatusRunning {
		return ErrJobNotRunning
	}
	now := time.Now().UTC()
	j.Status = models.StatusFailed
	j.FinishedAt = &now
	j.Error = msg
	j.Report = nil
	return nil
}

var randomReader = rand.Read

func newAnalysisID() (string, error) {
	var b [16]byte
	if _, err := randomReader(b[:]); err != nil {
		return "", fmt.Errorf("generate analysis id: %w", err)
	}
	return hex.EncodeToString(b[:]), nil
}
