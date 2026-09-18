package scheduler

import (
	"context"
	"errors"
	"sync"
	"time"

	"sapasora/platform/logger"
)

// MemoryWorkerPool is an in-memory implementation of WorkerPool
type MemoryWorkerPool struct {
	queue      JobQueue
	workers    []*Worker
	workerWg   sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
	log        *logger.Logger
	running    bool
	workerID   int
	workerIDMu sync.Mutex
	stats      *WorkerStats
	statsMu    sync.RWMutex
}

// Worker represents a single worker in the pool
type Worker struct {
	ID       int
	pool     *MemoryWorkerPool
	log      *logger.Logger
	running  bool
	jobQueue JobQueue
}

// NewMemoryWorkerPool creates a new worker pool
func NewMemoryWorkerPool(queue JobQueue, log *logger.Logger, numWorkers int) *MemoryWorkerPool {
	ctx, cancel := context.WithCancel(context.Background())

	pool := &MemoryWorkerPool{
		queue:   queue,
		workers: make([]*Worker, 0),
		log:     log,
		ctx:     ctx,
		cancel:  cancel,
		stats: &WorkerStats{
			Workers:     0,
			RunningJobs: 0,
			QueuedJobs:  0,
			Processed:   0,
			Failed:      0,
			ActiveJobs:  make([]*JobStatus, 0),
		},
	}

	// Add initial workers
	for range numWorkers {
		pool.AddWorker()
	}

	return pool
}

// Start starts the worker pool
func (p *MemoryWorkerPool) Start(ctx context.Context) error {
	if p.running {
		return errors.New("worker pool already running")
	}

	p.running = true
	p.log.Info("Starting worker pool", "workers", len(p.workers))

	// Start each worker
	for _, worker := range p.workers {
		p.workerWg.Add(1)
		go worker.Work(p.ctx)
	}

	p.log.Info("Worker pool started", "workers", len(p.workers))
	return nil
}

// Stop stops the worker pool gracefully
func (p *MemoryWorkerPool) Stop(ctx context.Context) error {
	if !p.running {
		return errors.New("worker pool not running")
	}

	p.log.Info("Stopping worker pool", "workers", len(p.workers))

	// Cancel all workers
	p.cancel()

	// Wait for all workers to finish
	done := make(chan struct{})
	go func() {
		p.workerWg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// All workers have stopped
	case <-ctx.Done():
		return ctx.Err()
	}

	p.running = false
	p.log.Info("Worker pool stopped")
	return nil
}

// AddWorker adds a new worker to the pool
func (p *MemoryWorkerPool) AddWorker() error {
	p.workerIDMu.Lock()
	p.workerID++
	workerID := p.workerID
	p.workerIDMu.Unlock()

	worker := &Worker{
		ID:       workerID,
		pool:     p,
		log:      p.log,
		jobQueue: p.queue,
	}

	// Add to the pool
	p.workers = append(p.workers, worker)

	// Update stats
	p.statsMu.Lock()
	p.stats.Workers++
	p.statsMu.Unlock()

	p.log.Info("Worker added to pool", "worker_id", workerID, "total_workers", len(p.workers))

	return nil
}

// RemoveWorker removes a worker from the pool
func (p *MemoryWorkerPool) RemoveWorker() error {
	if len(p.workers) <= 1 {
		return errors.New("cannot remove worker: at least one worker required")
	}

	// Remove the last worker
	p.workers = p.workers[:len(p.workers)-1]

	// Update stats
	p.statsMu.Lock()
	p.stats.Workers--
	p.statsMu.Unlock()

	p.log.Info("Worker removed from pool", "total_workers", len(p.workers))

	return nil
}

// Stats returns worker pool statistics
func (p *MemoryWorkerPool) Stats() *WorkerStats {
	p.statsMu.RLock()
	defer p.statsMu.RUnlock()

	// Create a copy to avoid race conditions
	statsCopy := *p.stats

	// Create a copy of the active jobs slice
	activeJobsCopy := make([]*JobStatus, len(statsCopy.ActiveJobs))
	copy(activeJobsCopy, statsCopy.ActiveJobs)

	statsCopy.ActiveJobs = activeJobsCopy

	return &statsCopy
}

// Work is the main work loop for a worker
func (w *Worker) Work(ctx context.Context) {
	defer w.pool.workerWg.Done()

	// Keep track of which jobs this worker is currently running
	runningJobs := make(map[string]*Job)

	for {
		select {
		case <-ctx.Done():
			// Context cancelled, worker should stop
			w.log.Info("Worker shutting down", "worker_id", w.ID)

			// Process any running jobs if needed (e.g., requeue if gracefully stopping)
			for _, job := range runningJobs {
				w.log.Info(
					"Worker stopped while job was running",
					"job_id",
					job.ID,
					"worker_id",
					w.ID,
				)
			}
			return

		default:
			// Try to get a job from the queue
			job, err := w.jobQueue.Dequeue(
				ctx,
				"default",
			) // In a real implementation, we'd support multiple queues
			if err != nil {
				// No jobs available, wait a bit before trying again
				time.Sleep(100 * time.Millisecond)
				continue
			}

			// Track this job as running
			runningJobs[job.ID] = job

			// Update stats
			w.pool.statsMu.Lock()
			w.pool.stats.RunningJobs++
			w.pool.stats.ActiveJobs = append(w.pool.stats.ActiveJobs, &JobStatus{
				ID:        job.ID,
				Queue:     job.Queue,
				Class:     job.Class,
				StartedAt: *job.StartedAt,
			})
			w.pool.statsMu.Unlock()

			// Execute the job
			jobErr := job.JobHandler.Perform(ctx, job.Args...)

			// Mark job as finished
			finishedAt := time.Now()
			job.FinishedAt = &finishedAt

			if jobErr != nil {
				// Job failed
				job.Error = jobErr.Error()
				w.pool.statsMu.Lock()
				w.pool.stats.Failed++
				w.pool.statsMu.Unlock()

				// Retry the job if possible
				if job.RetryCount < job.MaxRetries {
					job.RetryCount++
					err := w.jobQueue.Retry(ctx, job.ID)
					if err != nil {
						w.log.Error("Failed to retry job", "job_id", job.ID, "error", err)
					} else {
						w.log.Info("Job retried", "job_id", job.ID, "attempt", job.RetryCount)
					}
				} else {
					w.log.Error("Job failed permanently", "job_id", job.ID, "error", jobErr)
				}
			} else {
				// Job succeeded
				w.pool.statsMu.Lock()
				w.pool.stats.Processed++
				w.pool.statsMu.Unlock()

				w.log.Info("Job completed successfully", "job_id", job.ID)
			}

			// Remove from running jobs
			delete(runningJobs, job.ID)

			// Update stats
			w.pool.statsMu.Lock()
			w.pool.stats.RunningJobs--
			// Remove from active jobs list
			activeJobs := make([]*JobStatus, 0, len(w.pool.stats.ActiveJobs))
			for _, status := range w.pool.stats.ActiveJobs {
				if status.ID != job.ID {
					activeJobs = append(activeJobs, status)
				}
			}
			w.pool.stats.ActiveJobs = activeJobs
			w.pool.statsMu.Unlock()

			// Update the job in the queue
			_ = w.jobQueue.Delete(ctx, job.ID) // In a real implementation, this would be an update
		}
	}
}

// QueueWorkerPool is a more advanced queue-based worker pool implementation
type QueueWorkerPool struct {
	queues        map[string]JobQueue
	workers       []*QueueWorker
	workerWg      sync.WaitGroup
	ctx           context.Context
	cancel        context.CancelFunc
	log           *logger.Logger
	running       bool
	workerID      int
	workerIDMu    sync.Mutex
	mu            sync.RWMutex
	stats         *WorkerStats
	statsMu       sync.RWMutex
	queuesToWatch []string
}

// QueueWorker represents a worker that can watch multiple queues
type QueueWorker struct {
	ID         int
	pool       *QueueWorkerPool
	log        *logger.Logger
	running    bool
	queues     []string // Queues this worker watches
	queuesMap  map[string]JobQueue
	queueIndex int
}

// NewQueueWorkerPool creates a new queue-based worker pool
func NewQueueWorkerPool(
	queues map[string]JobQueue,
	log *logger.Logger,
	numWorkers int,
	queuesToWatch []string,
) *QueueWorkerPool {
	ctx, cancel := context.WithCancel(context.Background())

	if len(queuesToWatch) == 0 {
		queuesToWatch = []string{"default"}
	}

	pool := &QueueWorkerPool{
		queues:        queues,
		workers:       make([]*QueueWorker, 0),
		log:           log,
		ctx:           ctx,
		cancel:        cancel,
		queuesToWatch: queuesToWatch,
		stats: &WorkerStats{
			Workers:     0,
			RunningJobs: 0,
			QueuedJobs:  0,
			Processed:   0,
			Failed:      0,
			ActiveJobs:  make([]*JobStatus, 0),
		},
	}

	// Add initial workers
	for range numWorkers {
		pool.AddWorker()
	}

	return pool
}

// Start starts the queue worker pool
func (p *QueueWorkerPool) Start(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.running {
		return errors.New("queue worker pool already running")
	}

	p.running = true
	p.log.Info("Starting queue worker pool", "workers", len(p.workers), "queues", p.queuesToWatch)

	// Start each worker
	for _, worker := range p.workers {
		p.workerWg.Add(1)
		go worker.Work(p.ctx)
	}

	p.log.Info("Queue worker pool started", "workers", len(p.workers), "queues", p.queuesToWatch)
	return nil
}

// Stop stops the queue worker pool gracefully
func (p *QueueWorkerPool) Stop(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.running {
		return errors.New("queue worker pool not running")
	}

	p.log.Info("Stopping queue worker pool", "workers", len(p.workers))

	// Cancel all workers
	p.cancel()

	// Wait for all workers to finish
	done := make(chan struct{})
	go func() {
		p.workerWg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// All workers have stopped
	case <-ctx.Done():
		return ctx.Err()
	}

	p.running = false
	p.log.Info("Queue worker pool stopped")
	return nil
}

// AddWorker adds a new worker to the pool
func (p *QueueWorkerPool) AddWorker() error {
	p.workerIDMu.Lock()
	p.workerID++
	workerID := p.workerID
	p.workerIDMu.Unlock()

	worker := &QueueWorker{
		ID:        workerID,
		pool:      p,
		log:       p.log,
		queues:    append([]string(nil), p.queuesToWatch...), // copy the slice
		queuesMap: make(map[string]JobQueue),
	}

	// Create a mapping of queues this worker will watch
	for _, qName := range p.queuesToWatch {
		if q, exists := p.queues[qName]; exists {
			worker.queuesMap[qName] = q
		} else {
			// Create a default queue if it doesn't exist
			defaultQueue := NewMemoryJobQueue(p.log, 1)
			p.queues[qName] = defaultQueue
			worker.queuesMap[qName] = defaultQueue
		}
	}

	// Add to the pool
	p.workers = append(p.workers, worker)

	// Update stats
	p.statsMu.Lock()
	p.stats.Workers++
	p.statsMu.Unlock()

	p.log.Info("Queue worker added to pool", "worker_id", workerID, "total_workers", len(p.workers))

	return nil
}

// RemoveWorker removes a worker from the pool
func (p *QueueWorkerPool) RemoveWorker() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.workers) <= 1 {
		return errors.New("cannot remove worker: at least one worker required")
	}

	// Remove the last worker
	p.workers = p.workers[:len(p.workers)-1]

	// Update stats
	p.statsMu.Lock()
	p.stats.Workers--
	p.statsMu.Unlock()

	p.log.Info("Queue worker removed from pool", "total_workers", len(p.workers))

	return nil
}

// Stats returns worker pool statistics
func (p *QueueWorkerPool) Stats() *WorkerStats {
	p.statsMu.RLock()
	defer p.statsMu.RUnlock()

	// Create a copy to avoid race conditions
	statsCopy := *p.stats

	// Create a copy of the active jobs slice
	activeJobsCopy := make([]*JobStatus, len(statsCopy.ActiveJobs))
	copy(activeJobsCopy, statsCopy.ActiveJobs)

	statsCopy.ActiveJobs = activeJobsCopy

	return &statsCopy
}

// Work is the main work loop for a queue worker
func (w *QueueWorker) Work(ctx context.Context) {
	defer w.pool.workerWg.Done()

	// Keep track of which jobs this worker is currently running
	runningJobs := make(map[string]*Job)

	for {
		select {
		case <-ctx.Done():
			// Context cancelled, worker should stop
			w.log.Info("Queue worker shutting down", "worker_id", w.ID)

			// Process any running jobs if needed (e.g., requeue if gracefully stopping)
			for _, job := range runningJobs {
				w.log.Info(
					"Queue worker stopped while job was running",
					"job_id",
					job.ID,
					"worker_id",
					w.ID,
				)
			}
			return

		default:
			// Try to get a job from one of the watched queues
			job, queueName, err := w.getNextJob(ctx)
			if err != nil {
				// No jobs available across all queues, wait a bit before trying again
				time.Sleep(100 * time.Millisecond)
				continue
			}

			// Track this job as running
			runningJobs[job.ID] = job

			// Update stats
			w.pool.statsMu.Lock()
			w.pool.stats.RunningJobs++
			w.pool.stats.ActiveJobs = append(w.pool.stats.ActiveJobs, &JobStatus{
				ID:        job.ID,
				Queue:     job.Queue,
				Class:     job.Class,
				StartedAt: *job.StartedAt,
			})
			w.pool.statsMu.Unlock()

			// Execute the job
			jobErr := job.JobHandler.Perform(ctx, job.Args...)

			// Mark job as finished
			finishedAt := time.Now()
			job.FinishedAt = &finishedAt

			if jobErr != nil {
				// Job failed
				job.Error = jobErr.Error()
				w.pool.statsMu.Lock()
				w.pool.stats.Failed++
				w.pool.statsMu.Unlock()

				// Retry the job if possible
				if job.RetryCount < job.MaxRetries {
					job.RetryCount++
					err := w.pool.queues[queueName].Retry(ctx, job.ID)
					if err != nil {
						w.log.Error("Failed to retry job", "job_id", job.ID, "error", err)
					} else {
						w.log.Info("Job retried", "job_id", job.ID, "attempt", job.RetryCount)
					}
				} else {
					w.log.Error("Job failed permanently", "job_id", job.ID, "error", jobErr)
				}
			} else {
				// Job succeeded
				w.pool.statsMu.Lock()
				w.pool.stats.Processed++
				w.pool.statsMu.Unlock()

				w.log.Info("Job completed successfully", "job_id", job.ID)
			}

			// Remove from running jobs
			delete(runningJobs, job.ID)

			// Update stats
			w.pool.statsMu.Lock()
			w.pool.stats.RunningJobs--
			// Remove from active jobs list
			activeJobs := make([]*JobStatus, 0, len(w.pool.stats.ActiveJobs))
			for _, status := range w.pool.stats.ActiveJobs {
				if status.ID != job.ID {
					activeJobs = append(activeJobs, status)
				}
			}
			w.pool.stats.ActiveJobs = activeJobs
			w.pool.statsMu.Unlock()

			// Update the job in the queue
			_ = w.pool.queues[queueName].Delete(ctx, job.ID)
		}
	}
}

// getNextJob tries to get the next available job from any of the watched queues
func (w *QueueWorker) getNextJob(ctx context.Context) (*Job, string, error) {
	// Round-robin through the queues
	for i := 0; i < len(w.queues); i++ {
		queueName := w.queues[w.queueIndex%len(w.queues)]
		w.queueIndex++

		queue, exists := w.queuesMap[queueName]
		if !exists {
			continue
		}

		job, err := queue.Dequeue(ctx, queueName)
		if err == nil {
			return job, queueName, nil
		}
	}

	// No jobs found in any queue
	return nil, "", errors.New("no jobs available")
}
