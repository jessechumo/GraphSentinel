package store

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"

	"github.com/graphsentinel/graphsentinel/pkg/models"
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

// Enqueue assigns an ID, records the job as queued, and returns a snapshot for the caller.
func (m *Memory) Enqueue(req models.AnalyzeRequest) (*models.AnalysisJob, error) {
	id := newAnalysisID()
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

func newAnalysisID() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		// crypto/rand only fails if the kernel refuses entropy; treat as fatal.
		panic("graphsentinel: cannot generate analysis id: " + err.Error())
	}
	return hex.EncodeToString(b[:])
}
