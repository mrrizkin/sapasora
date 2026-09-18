package scheduler

import (
	"context"
	"fmt"
	"sync"
	"time"

	"sapasora/platform/logger"
)

// MemoryJobQueue is an in-memory implementation of JobQueue
type MemoryJobQueue struct {
	jobs       map[string]*Job
	queue      map[string]chan *Job
	mutex      sync.RWMutex
	log        *logger.Logger
	workers    int
	workerPool *MemoryWorkerPool
}

// NewMemoryJobQueue creates a new in-memory job queue
func NewMemoryJobQueue(log *logger.Logger, workers int) *MemoryJobQueue {
	if workers <= 0 {
		workers = 1
	}

	q := &MemoryJobQueue{
		jobs:    make(map[string]*Job),
		queue:   make(map[string]chan *Job),
		log:     log,
		workers: workers,
	}

	// Initialize default queue
	q.queue["default"] = make(chan *Job, 100)

	return q
}

// Enqueue adds a job to the queue
func (q *MemoryJobQueue) Enqueue(ctx context.Context, job *Job) error {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	// Create the queue channel if it doesn't exist
	if _, exists := q.queue[job.Queue]; !exists {
		q.queue[job.Queue] = make(chan *Job, 100)
	}

	// Store the job
	q.jobs[job.ID] = job
	enqueuedAt := time.Now()
	job.EnqueuedAt = &enqueuedAt

	// Add to the queue channel
	select {
	case q.queue[job.Queue] <- job:
		q.log.Info("Job enqueued", "job_id", job.ID, "queue", job.Queue)
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Dequeue retrieves a job from the queue
func (q *MemoryJobQueue) Dequeue(ctx context.Context, queue string) (*Job, error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	queueChan, exists := q.queue[queue]
	if !exists {
		return nil, fmt.Errorf("queue %s does not exist", queue)
	}

	select {
	case job := <-queueChan:
		startedAt := time.Now()
		job.StartedAt = &startedAt
		q.log.Info("Job dequeued", "job_id", job.ID, "queue", queue)
		return job, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// Retry a job that failed
func (q *MemoryJobQueue) Retry(ctx context.Context, jobID string) error {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	job, exists := q.jobs[jobID]
	if !exists {
		return fmt.Errorf("job %s not found", jobID)
	}

	if job.RetryCount >= job.MaxRetries {
		now := time.Now()
		job.FailedAt = &now
		q.log.Info("Job failed permanently", "job_id", jobID, "retries", job.RetryCount)
		return nil
	}

	// Increase retry count
	job.RetryCount++

	// Remove from any active tracking
	if job.StartedAt != nil && job.FinishedAt == nil {
		now := time.Now()
		job.FinishedAt = &now
	}

	// Reset for retry
	job.StartedAt = nil
	job.FinishedAt = nil
	job.Error = ""

	// Add back to queue
	enqueuedAt := time.Now()
	job.EnqueuedAt = &enqueuedAt

	select {
	case q.queue[job.Queue] <- job:
		q.log.Info("Job retried", "job_id", jobID, "attempt", job.RetryCount)
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Delete removes a job from the queue
func (q *MemoryJobQueue) Delete(ctx context.Context, jobID string) error {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	job, exists := q.jobs[jobID]
	if !exists {
		return fmt.Errorf("job %s not found", jobID)
	}

	// Mark as finished
	finishedAt := time.Now()
	job.FinishedAt = &finishedAt

	delete(q.jobs, jobID)
	q.log.Info("Job deleted", "job_id", jobID)

	return nil
}

// Get retrieves a job by ID
func (q *MemoryJobQueue) Get(ctx context.Context, jobID string) (*Job, error) {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	job, exists := q.jobs[jobID]
	if !exists {
		return nil, fmt.Errorf("job %s not found", jobID)
	}

	return job, nil
}

// Stats returns queue statistics
func (q *MemoryJobQueue) Stats(ctx context.Context) (*QueueStats, error) {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	stats := &QueueStats{
		Queues:     make(map[string]int64),
		RetryStats: make(map[string]int64),
		Processes:  make(map[string]string),
		Workers:    q.workers,
	}

	// Count jobs in each queue
	for queueName, queueChan := range q.queue {
		stats.Queues[queueName] = int64(len(queueChan))
	}

	// Count processed, failed, and retry jobs
	for _, job := range q.jobs {
		if job.FinishedAt != nil {
			if job.Error != "" {
				stats.Failed++
				stats.RetryStats[job.Queue]++
			} else {
				stats.Processed++
			}
		} else if job.StartedAt != nil {
			// Running jobs
			stats.Retry++
		} else {
			// Scheduled/queued jobs
			stats.Scheduled++
		}
	}

	return stats, nil
}

// CleanUp removes old finished jobs
func (q *MemoryJobQueue) CleanUp(ctx context.Context, olderThan time.Duration) error {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	cutoffTime := time.Now().Add(-olderThan)

	for id, job := range q.jobs {
		if job.FinishedAt != nil && job.FinishedAt.Before(cutoffTime) {
			delete(q.jobs, id)
			q.log.Info("Cleaned up old job", "job_id", id)
		}
	}

	return nil
}

// JobQueueWithStorage adds database persistence to the job queue
type JobQueueWithStorage struct {
	baseQueue JobQueue
	storage   JobStorage
	log       *logger.Logger
}

// JobStorage interface for job persistence
type JobStorage interface {
	StoreJob(ctx context.Context, job *Job) error
	GetJob(ctx context.Context, jobID string) (*Job, error)
	DeleteJob(ctx context.Context, jobID string) error
	UpdateJob(ctx context.Context, job *Job) error
	GetAllJobs(ctx context.Context) ([]*Job, error)
	GetJobsByQueue(ctx context.Context, queue string) ([]*Job, error)
	GetFailedJobs(ctx context.Context) ([]*Job, error)
	CleanOldJobs(ctx context.Context, olderThan time.Duration) error
}

// NewJobQueueWithStorage creates a job queue with persistence
func NewJobQueueWithStorage(
	baseQueue JobQueue,
	storage JobStorage,
	log *logger.Logger,
) *JobQueueWithStorage {
	return &JobQueueWithStorage{
		baseQueue: baseQueue,
		storage:   storage,
		log:       log,
	}
}

// Enqueue adds a job to the queue and stores it in persistent storage
func (q *JobQueueWithStorage) Enqueue(ctx context.Context, job *Job) error {
	// Store in persistent storage first
	if err := q.storage.StoreJob(ctx, job); err != nil {
		q.log.Error("Failed to store job", "job_id", job.ID, "error", err)
		return err
	}

	// Then add to the base queue
	return q.baseQueue.Enqueue(ctx, job)
}

// Dequeue retrieves a job from the queue
func (q *JobQueueWithStorage) Dequeue(ctx context.Context, queue string) (*Job, error) {
	// First try base queue
	job, err := q.baseQueue.Dequeue(ctx, queue)
	if err != nil {
		return nil, err
	}

	// Update in storage
	if err := q.storage.UpdateJob(ctx, job); err != nil {
		q.log.Error("Failed to update job in storage", "job_id", job.ID, "error", err)
		// Continue anyway, as the job was already dequeued from memory
	}

	return job, nil
}

// Retry a job that failed
func (q *JobQueueWithStorage) Retry(ctx context.Context, jobID string) error {
	job, err := q.baseQueue.Get(ctx, jobID)
	if err != nil {
		return err
	}

	if err := q.baseQueue.Retry(ctx, jobID); err != nil {
		return err
	}

	// Update in storage
	if err := q.storage.UpdateJob(ctx, job); err != nil {
		q.log.Error("Failed to update job in storage", "job_id", jobID, "error", err)
	}

	return nil
}

// Delete removes a job from the queue and from persistent storage
func (q *JobQueueWithStorage) Delete(ctx context.Context, jobID string) error {
	// Delete from base queue
	if err := q.baseQueue.Delete(ctx, jobID); err != nil {
		return err
	}

	// Then delete from storage
	return q.storage.DeleteJob(ctx, jobID)
}

// Get retrieves a job by ID
func (q *JobQueueWithStorage) Get(ctx context.Context, jobID string) (*Job, error) {
	// First check in memory
	job, err := q.baseQueue.Get(ctx, jobID)
	if err == nil {
		return job, nil
	}

	// If not found in memory, check persistent storage
	return q.storage.GetJob(ctx, jobID)
}

// Stats returns queue statistics
func (q *JobQueueWithStorage) Stats(ctx context.Context) (*QueueStats, error) {
	return q.baseQueue.Stats(ctx)
}

// CleanUp removes old finished jobs
func (q *JobQueueWithStorage) CleanUp(ctx context.Context, olderThan time.Duration) error {
	// Clean up in memory
	if err := q.baseQueue.CleanUp(ctx, olderThan); err != nil {
		return err
	}

	// Clean up in storage
	return q.storage.CleanOldJobs(ctx, olderThan)
}

// MemoryJobStorage is an in-memory implementation of JobStorage
type MemoryJobStorage struct {
	jobs  map[string]*Job
	mutex sync.RWMutex
}

// NewMemoryJobStorage creates a new in-memory job storage
func NewMemoryJobStorage() *MemoryJobStorage {
	return &MemoryJobStorage{
		jobs: make(map[string]*Job),
	}
}

// StoreJob stores a job in memory
func (s *MemoryJobStorage) StoreJob(ctx context.Context, job *Job) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	jobCopy := *job
	// We don't store the JobHandler as it's not serializable
	jobCopy.JobHandler = nil
	s.jobs[job.ID] = &jobCopy
	return nil
}

// GetJob retrieves a job by ID
func (s *MemoryJobStorage) GetJob(ctx context.Context, jobID string) (*Job, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	job, exists := s.jobs[jobID]
	if !exists {
		return nil, fmt.Errorf("job %s not found", jobID)
	}

	return job, nil
}

// DeleteJob deletes a job by ID
func (s *MemoryJobStorage) DeleteJob(ctx context.Context, jobID string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.jobs, jobID)
	return nil
}

// UpdateJob updates an existing job
func (s *MemoryJobStorage) UpdateJob(ctx context.Context, job *Job) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Store a copy without the handler
	jobCopy := *job
	jobCopy.JobHandler = nil
	s.jobs[job.ID] = &jobCopy
	return nil
}

// GetAllJobs returns all jobs
func (s *MemoryJobStorage) GetAllJobs(ctx context.Context) ([]*Job, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	jobs := make([]*Job, 0, len(s.jobs))
	for _, job := range s.jobs {
		jobs = append(jobs, job)
	}

	return jobs, nil
}

// GetJobsByQueue returns all jobs for a specific queue
func (s *MemoryJobStorage) GetJobsByQueue(ctx context.Context, queue string) ([]*Job, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var jobs []*Job
	for _, job := range s.jobs {
		if job.Queue == queue {
			jobs = append(jobs, job)
		}
	}

	return jobs, nil
}

// GetFailedJobs returns all failed jobs
func (s *MemoryJobStorage) GetFailedJobs(ctx context.Context) ([]*Job, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var jobs []*Job
	for _, job := range s.jobs {
		if job.Error != "" || job.FailedAt != nil {
			jobs = append(jobs, job)
		}
	}

	return jobs, nil
}

// CleanOldJobs removes jobs older than the specified duration
func (s *MemoryJobStorage) CleanOldJobs(ctx context.Context, olderThan time.Duration) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	cutoffTime := time.Now().Add(-olderThan)

	for id, job := range s.jobs {
		if job.FinishedAt != nil && job.FinishedAt.Before(cutoffTime) {
			delete(s.jobs, id)
		}
	}

	return nil
}
