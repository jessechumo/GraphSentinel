package workers

import (
	"context"
	"sync"

	"github.com/graphsentinel/graphsentinel/internal/store"
	"github.com/graphsentinel/graphsentinel/pkg/models"
)

// Processor runs a single analysis job and returns a report or error.
type Processor func(ctx context.Context, job *models.AnalysisJob) (*models.AnalysisReport, error)

// Pool is a fixed worker count queue that claims jobs from the store and runs the processor.
type Pool struct {
	work    chan string
	wg      sync.WaitGroup
	once    sync.Once
	st      store.WorkerStore
	proc    Processor
	workers int
}

// NewPool constructs a pool. workers must be at least 1; qsize is the buffered work queue depth.
func NewPool(workers, qsize int, st store.WorkerStore, proc Processor) *Pool {
	if workers < 1 {
		workers = 1
	}
	if qsize < 1 {
		qsize = 64
	}
	if proc == nil {
		proc = func(ctx context.Context, job *models.AnalysisJob) (*models.AnalysisReport, error) {
			return nil, nil
		}
	}
	return &Pool{
		work:    make(chan string, qsize),
		st:      st,
		proc:    proc,
		workers: workers,
	}
}

// Start launches worker goroutines. Call Close then Wait during shutdown.
func (p *Pool) Start() {
	for range p.workers {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			for id := range p.work {
				p.runJob(context.Background(), id)
			}
		}()
	}
}

func (p *Pool) runJob(ctx context.Context, id string) {
	if !p.st.TryClaim(id) {
		return
	}
	job, ok := p.st.Get(id)
	if !ok {
		return
	}
	report, err := p.proc(ctx, job)
	if err != nil {
		_ = p.st.Fail(id, err.Error())
		return
	}
	if report == nil {
		_ = p.st.Fail(id, "processor returned no report")
		return
	}
	_ = p.st.Complete(id, report)
}

// Submit enqueues a job id for processing. It blocks if the queue is full.
func (p *Pool) Submit(id string) {
	p.work <- id
}

// Close stops accepting new work and signals workers to exit after draining the queue.
func (p *Pool) Close() {
	p.once.Do(func() { close(p.work) })
}

// Wait blocks until all workers exit after Close.
func (p *Pool) Wait() {
	p.wg.Wait()
}
