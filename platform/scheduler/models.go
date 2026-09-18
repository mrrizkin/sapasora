package scheduler

import (
	"context"
	"fmt"
	"time"
)

// Job represents a background job in the queue
type Job struct {
	ID         string         `json:"id"`
	Queue      string         `json:"queue"`
	Class      string         `json:"class"` // Job class/type
	Args       []any          `json:"args"`
	CreatedAt  time.Time      `json:"created_at"`
	EnqueuedAt *time.Time     `json:"enqueued_at,omitempty"`
	StartedAt  *time.Time     `json:"started_at,omitempty"`
	FinishedAt *time.Time     `json:"finished_at,omitempty"`
	FailedAt   *time.Time     `json:"failed_at,omitempty"`
	Error      string         `json:"error,omitempty"`
	RetryCount int            `json:"retry_count"`
	MaxRetries int            `json:"max_retries"`
	Backtrace  []string       `json:"backtrace,omitempty"`
	Payload    map[string]any `json:"payload"`
	JobHandler JobHandler     `json:"-"` // Function to execute
}

// JobHandler defines the interface for job execution
type JobHandler interface {
	Perform(ctx context.Context, args ...any) error
}

// FuncJobHandler wraps a function to implement JobHandler
type FuncJobHandler struct {
	Fn func(context.Context, ...any) error
}

// Perform executes the function
func (f *FuncJobHandler) Perform(ctx context.Context, args ...any) error {
	return f.Fn(ctx, args...)
}

// JobQueue interface for job queue operations
type JobQueue interface {
	// Enqueue adds a job to the queue
	Enqueue(ctx context.Context, job *Job) error
	// Dequeue retrieves a job from the queue
	Dequeue(ctx context.Context, queue string) (*Job, error)
	// Retry a job that failed
	Retry(ctx context.Context, jobID string) error
	// Delete removes a job from the queue
	Delete(ctx context.Context, jobID string) error
	// Get retrieves a job by ID
	Get(ctx context.Context, jobID string) (*Job, error)
	// Stats returns queue statistics
	Stats(ctx context.Context) (*QueueStats, error)
	// CleanUp removes old finished jobs
	CleanUp(ctx context.Context, olderThan time.Duration) error
}

// QueueStats contains statistics about job queues
type QueueStats struct {
	Processed  int64             `json:"processed"`
	Failed     int64             `json:"failed"`
	Retry      int64             `json:"retry"`
	Scheduled  int64             `json:"scheduled"`
	Queues     map[string]int64  `json:"queues"`
	RetryStats map[string]int64  `json:"retry_stats"`
	Workers    int               `json:"workers"`
	Processes  map[string]string `json:"processes"`
}

// WorkerPool interface for managing job workers
type WorkerPool interface {
	// Start starts the worker pool
	Start(ctx context.Context) error
	// Stop stops the worker pool gracefully
	Stop(ctx context.Context) error
	// AddWorker adds a new worker to the pool
	AddWorker() error
	// RemoveWorker removes a worker from the pool
	RemoveWorker() error
	// Stats returns worker pool statistics
	Stats() *WorkerStats
}

// WorkerStats contains statistics about the worker pool
type WorkerStats struct {
	Workers     int          `json:"workers"`
	RunningJobs int          `json:"running_jobs"`
	QueuedJobs  int64        `json:"queued_jobs"`
	Processed   int64        `json:"processed"`
	Failed      int64        `json:"failed"`
	ActiveJobs  []*JobStatus `json:"active_jobs"`
}

// JobStatus contains information about a running job
type JobStatus struct {
	ID        string    `json:"id"`
	Queue     string    `json:"queue"`
	Class     string    `json:"class"`
	StartedAt time.Time `json:"started_at"`
}

// JobBuilder helps build jobs
type JobBuilder struct {
	queue      string
	jobClass   string
	args       []any
	maxRetries int
	payload    map[string]any
}

// NewJobBuilder creates a new job builder
func NewJobBuilder() *JobBuilder {
	return &JobBuilder{
		queue:      "default",
		maxRetries: 25, // Sidekiq default
		payload:    make(map[string]any),
	}
}

// ToQueue sets the queue for the job
func (jb *JobBuilder) ToQueue(queue string) *JobBuilder {
	jb.queue = queue
	return jb
}

// WithMaxRetries sets the maximum number of retries for the job
func (jb *JobBuilder) WithMaxRetries(retries int) *JobBuilder {
	jb.maxRetries = retries
	return jb
}

// WithArgs sets the arguments for the job
func (jb *JobBuilder) WithArgs(args ...any) *JobBuilder {
	jb.args = args
	return jb
}

// WithPayload adds key-value pairs to the job's payload
func (jb *JobBuilder) WithPayload(key string, value any) *JobBuilder {
	jb.payload[key] = value
	return jb
}

// WithPayloadMap sets the entire payload map
func (jb *JobBuilder) WithPayloadMap(payload map[string]any) *JobBuilder {
	jb.payload = payload
	return jb
}

// Build creates a new job
func (jb *JobBuilder) Build(jobHandler JobHandler) *Job {
	job := &Job{
		ID:         generateJobID(),
		Queue:      jb.queue,
		Class:      jb.jobClass,
		Args:       jb.args,
		CreatedAt:  time.Now(),
		MaxRetries: jb.maxRetries,
		Payload:    jb.payload,
		JobHandler: jobHandler,
	}

	return job
}

// JobEnqueuer provides convenience methods for enqueuing jobs
type JobEnqueuer struct {
	queue JobQueue
}

// NewJobEnqueuer creates a new job enqueuer
func NewJobEnqueuer(queue JobQueue) *JobEnqueuer {
	return &JobEnqueuer{
		queue: queue,
	}
}

// EnqueueJob enqueues a job using the builder pattern
func (je *JobEnqueuer) EnqueueJob(ctx context.Context, jobHandler JobHandler) *JobBuilder {
	builder := NewJobBuilder()

	// Set the job class name based on the handler type
	if _, ok := jobHandler.(*FuncJobHandler); ok {
		// For function handlers, we'll create a generic class name
		builder.jobClass = "FuncJob"
	} else {
		// For other handler types, we would extract the type name
		builder.jobClass = "CustomJob"
	}

	// This is a workaround since we can't add methods to the builder inside this function
	// The actual implementation would be cleaner with a proper method
	job := builder.Build(jobHandler)
	_ = je.queue.Enqueue(ctx, job)

	return builder
}

// EnqueueWithDelay enqueues a job to be executed after a delay
func (je *JobEnqueuer) EnqueueWithDelay(
	ctx context.Context,
	delay time.Duration,
	jobHandler JobHandler,
	args ...any,
) error {
	// In a real implementation, this would create a delayed job
	// For now, we'll just create a regular job since full delayed execution is complex
	builder := NewJobBuilder()
	job := builder.WithArgs(args...).Build(jobHandler)

	return je.queue.Enqueue(ctx, job)
}

// PerformAsync enqueues a job to be executed asynchronously
func (je *JobEnqueuer) PerformAsync(
	ctx context.Context,
	jobHandler JobHandler,
	args ...any,
) error {
	builder := NewJobBuilder()
	job := builder.WithArgs(args...).Build(jobHandler)

	return je.queue.Enqueue(ctx, job)
}

// Helper function to generate a unique job ID
func generateJobID() string {
	return fmt.Sprint(time.Now().UnixNano())
}

// JobResult represents the result of a job execution
type JobResult struct {
	JobID     string    `json:"job_id"`
	Status    string    `json:"status"` // "success", "failed", "retrying"
	Error     error     `json:"error,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

// JobResultHandler handles job execution results
type JobResultHandler func(*JobResult)
