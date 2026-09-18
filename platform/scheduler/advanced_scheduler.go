package scheduler

import (
	"context"
	"fmt"
	"time"

	"sapasora/platform/config"
	"sapasora/platform/database"
	"sapasora/platform/logger"
)

// AdvancedScheduler is the main Sidekiq-like scheduler that integrates all components
type AdvancedScheduler struct {
	queue       JobQueue
	workerPool  WorkerPool
	log         *logger.Logger
	db          *database.Database
	dashboard   *DashboardHandler
	jobEnqueuer *JobEnqueuer
	config      SchedulerConfig
}

// SchedulerConfig holds configuration for the advanced scheduler
type SchedulerConfig struct {
	Workers    int           `name:"workers"     env:"SCHEDULER_WORKERS,default=3"`
	Queue      string        `name:"queue"       env:"SCHEDULER_QUEUE,default=default"`
	MaxRetries int           `name:"max_retries" env:"SCHEDULER_MAX_RETRIES,default=25"`
	Timeout    time.Duration `name:"timeout"     env:"SCHEDULER_TIMEOUT,default=30m"`
	DBStorage  bool          `name:"db_storage"  env:"SCHEDULER_DB_STORAGE,default=false"`
	DBDriver   string        `name:"db_driver"   env:"SCHEDULER_DB_DRIVER"`
	DBDSN      string        `name:"db_dsn"      env:"SCHEDULER_DB_DSN"`
}

// NewAdvancedScheduler creates a new advanced scheduler
func NewAdvancedScheduler(
	cfg config.Config,
	log *logger.Logger,
	db *database.Database,
) (*AdvancedScheduler, error) {
	// Load configuration
	schedulerConfig := SchedulerConfig{
		Workers:    cfg.GetInt("scheduler.workers", 3),
		Queue:      cfg.GetString("scheduler.queue", "default"),
		MaxRetries: cfg.GetInt("scheduler.max_retries", 25),
		Timeout:    cfg.GetDuration("scheduler.timeout", 30*time.Minute),
		DBStorage:  cfg.GetBool("scheduler.db_storage", false),
		DBDriver:   cfg.GetString("scheduler.db_driver", ""),
		DBDSN:      cfg.GetString("scheduler.db_dsn", ""),
	}

	// Create the base in-memory queue
	baseQueue := NewMemoryJobQueue(log, schedulerConfig.Workers)

	var persistentQueue JobQueue = baseQueue

	// Add database storage if configured
	if schedulerConfig.DBStorage {
		if db != nil && db.DB != nil {
			// Get the underlying *sql.DB from GORM using the DB() method or raw connection
			// In GORM, we can use db.DB.DB (the underlying sql.DB)
			sqlDB, err := db.DB.DB()
			if err != nil {
				log.Error("Failed to get underlying SQL DB", "error", err)
				return nil, fmt.Errorf("failed to get underlying SQL DB: %w", err)
			}
			storage := NewDBJobStorage(sqlDB, log)
			persistentQueue = NewJobQueueWithStorage(baseQueue, storage, log)
		} else {
			log.Warn("Database storage requested but no database provided, using memory only")
		}
	}

	// Create worker pool
	workerPool := NewMemoryWorkerPool(persistentQueue, log, schedulerConfig.Workers)

	// Create dashboard handler
	dashboard := NewDashboardHandler(persistentQueue, workerPool, log)

	// Create job enqueuer
	jobEnqueuer := NewJobEnqueuer(persistentQueue)

	scheduler := &AdvancedScheduler{
		queue:       persistentQueue,
		workerPool:  workerPool,
		log:         log,
		db:          db,
		dashboard:   dashboard,
		jobEnqueuer: jobEnqueuer,
		config:      schedulerConfig,
	}

	return scheduler, nil
}

// Start starts the scheduler and its components
func (s *AdvancedScheduler) Start(ctx context.Context) error {
	s.log.Info("Starting advanced scheduler", "workers", s.config.Workers)

	// Start the worker pool
	if err := s.workerPool.Start(ctx); err != nil {
		return fmt.Errorf("failed to start worker pool: %w", err)
	}

	s.log.Info("Advanced scheduler started successfully")
	return nil
}

// Stop stops the scheduler and its components
func (s *AdvancedScheduler) Stop(ctx context.Context) error {
	s.log.Info("Stopping advanced scheduler")

	// Stop the worker pool
	if err := s.workerPool.Stop(ctx); err != nil {
		return fmt.Errorf("failed to stop worker pool: %w", err)
	}

	s.log.Info("Advanced scheduler stopped")
	return nil
}

// Enqueue adds a job to the queue
func (s *AdvancedScheduler) Enqueue(
	ctx context.Context,
	jobHandler JobHandler,
	queue string,
	args ...any,
) error {
	if queue == "" {
		queue = s.config.Queue
	}

	job := NewJobBuilder().
		ToQueue(queue).
		WithMaxRetries(s.config.MaxRetries).
		WithArgs(args...).
		Build(jobHandler)

	return s.queue.Enqueue(ctx, job)
}

// EnqueueIn enqueues a job to run after a specified delay
func (s *AdvancedScheduler) EnqueueIn(
	ctx context.Context,
	delay time.Duration,
	jobHandler JobHandler,
	queue string,
	args ...any,
) error {
	// In a real implementation, this would schedule the job to run after the delay
	// For now, we'll use a simple approach with a goroutine

	go func() {
		time.Sleep(delay)
		_ = s.Enqueue(ctx, jobHandler, queue, args...)
	}()

	return nil
}

// EnqueueAt enqueues a job to run at a specific time
func (s *AdvancedScheduler) EnqueueAt(
	ctx context.Context,
	at time.Time,
	jobHandler JobHandler,
	queue string,
	args ...any,
) error {
	// Calculate the delay until the specified time
	delay := time.Until(at)
	if delay < 0 {
		// Time is in the past, run immediately
		return s.Enqueue(ctx, jobHandler, queue, args...)
	}

	return s.EnqueueIn(ctx, delay, jobHandler, queue, args...)
}

// GetJob retrieves a job by ID
func (s *AdvancedScheduler) GetJob(ctx context.Context, jobID string) (*Job, error) {
	return s.queue.Get(ctx, jobID)
}

// RetryJob retries a failed job
func (s *AdvancedScheduler) RetryJob(ctx context.Context, jobID string) error {
	return s.queue.Retry(ctx, jobID)
}

// DeleteJob removes a job
func (s *AdvancedScheduler) DeleteJob(ctx context.Context, jobID string) error {
	return s.queue.Delete(ctx, jobID)
}

// GetStats returns scheduler statistics
func (s *AdvancedScheduler) GetStats(ctx context.Context) (*QueueStats, error) {
	return s.queue.Stats(ctx)
}

// GetWorkerStats returns worker pool statistics
func (s *AdvancedScheduler) GetWorkerStats() *WorkerStats {
	return s.workerPool.Stats()
}

// GetJobEnqueuer returns the job enqueuer for convenience
func (s *AdvancedScheduler) GetJobEnqueuer() *JobEnqueuer {
	return s.jobEnqueuer
}

// RegisterDashboardRoutes registers the dashboard routes with the provided router
func (s *AdvancedScheduler) RegisterDashboardRoutes(router any) {
	// The router type would need to be more specific based on the actual router used
	// This is a placeholder implementation
	if mux, ok := router.(*any); ok {
		// In a real implementation, this would register the routes
		_ = mux
	}
}

// PerformAsync is a convenience method to enqueue a function to run asynchronously
func (s *AdvancedScheduler) PerformAsync(
	ctx context.Context,
	fn func(context.Context, ...any) error,
	args ...any,
) error {
	handler := &FuncJobHandler{Fn: fn}
	return s.Enqueue(ctx, handler, "", args...)
}
