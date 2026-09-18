# Sidekiq-like Background Job Processing System

This package provides a comprehensive background job processing system inspired by Ruby's Sidekiq. It includes multiple job queues, a worker pool, persistence options, and a dashboard for monitoring.

## Features

- **Job Queues**: Multiple named queues for different types of jobs
- **Worker Pool**: Configurable number of concurrent workers
- **Job Persistence**: Store jobs in database (PostgreSQL, SQLite) or memory
- **Retry Logic**: Automatic retry with exponential backoff
- **Dashboard**: Web UI to monitor queue status and job progress
- **Job Builder**: Fluent API for creating and enqueueing jobs
- **Delayed Jobs**: Schedule jobs to run at a later time

## Architecture

```
[Job Enqueuer] -> [Job Queue] -> [Worker Pool] -> [Job Handlers]
      |              |              |
   [Dashboard]  [Persistence]  [Statistics]
```

## Usage

### Basic Job Enqueuing

```go
import "sapasora/platform/scheduler"

// Create a function job handler
fn := func(ctx context.Context, args ...interface{}) error {
    // Do work here
    fmt.Printf("Processing with args: %v\n", args)
    return nil
}

// Enqueue the job
err := advancedScheduler.PerformAsync(ctx, fn, "arg1", "arg2")
if err != nil {
    // handle error
}
```

### Custom Job Handler

```go
type EmailJob struct {
    To      string
    Subject string
    Body    string
}

func (e *EmailJob) Perform(ctx context.Context, args ...interface{}) error {
    // Send email
    return sendEmail(e.To, e.Subject, e.Body)
}

// Create and enqueue custom job
job := &EmailJob{
    To: "user@example.com",
    Subject: "Welcome!",
    Body: "Welcome to our service",
}

builder := scheduler.NewJobBuilder()
jobObj := builder.
    ToQueue("email").
    WithArgs("user@example.com", "Welcome!", "Welcome to our service").
    Build(&EmailJob{})

err := advancedScheduler.Enqueue(ctx, jobObj.JobHandler, "email", "user@example.com", "Welcome!", "Welcome to our service")
```

### Delayed Jobs

```go
// Enqueue job to run 1 hour from now
err := advancedScheduler.EnqueueIn(ctx, time.Hour, jobHandler, "queue-name", args...)

// Enqueue job to run at specific time
atTime := time.Now().Add(2 * time.Hour)
err := advancedScheduler.EnqueueAt(ctx, atTime, jobHandler, "queue-name", args...)
```

## Configuration

The scheduler can be configured via environment variables:

```env
SCHEDULER_WORKERS=5                    # Number of worker goroutines
SCHEDULER_QUEUE=default                # Default queue name
SCHEDULER_MAX_RETRIES=25               # Maximum retry attempts
SCHEDULER_TIMEOUT=30m                  # Job processing timeout
SCHEDULER_DB_STORAGE=true              # Enable database persistence
SCHEDULER_DB_DRIVER=postgres           # Database driver (postgres, sqlite)
SCHEDULER_DB_DSN=postgresql://...      # Database connection string
```

## Dashboard API Endpoints

- `GET /admin/jobs` - Main jobs dashboard
- `GET /admin/jobs/queues` - Queue statistics
- `GET /admin/jobs/workers` - Worker statistics
- `GET /admin/jobs/stats` - JSON stats endpoint
- `POST /admin/jobs/retry` - Retry a failed job
- `POST /admin/jobs/delete` - Delete a job

## Integration with Sapasora

The scheduler is integrated with Sapasora's dependency injection system using Uber's FX:

```go
app := fx.New(
    // ... other modules
    scheduler.Module,  // Provides both traditional cron and advanced scheduler
    fx.Invoke(useScheduler),
)

func useScheduler(advancedScheduler *scheduler.AdvancedScheduler) {
    // Use the scheduler
}
```

## Persistence Options

### Memory Storage (Default)
- Jobs are stored in memory
- Suitable for development and simple use cases
- Jobs are lost when application restarts

### Database Storage
- Jobs are persisted to PostgreSQL or SQLite
- Jobs survive application restarts
- Supports job history and monitoring

## Worker Pool Configuration

The worker pool can be dynamically scaled:

```go
// Add a worker
err := advancedScheduler.WorkerPool.AddWorker()

// Remove a worker
err := advancedScheduler.WorkerPool.RemoveWorker()

// Get statistics
stats := advancedScheduler.GetWorkerStats()
```

## Error Handling

The scheduler includes comprehensive error handling:

- Failed jobs are retried automatically (up to max_retries)
- Error details are stored with each job
- Dead job processing after all retries are exhausted
- Detailed logging for debugging

## Performance Considerations

- Use appropriate number of workers based on available CPU cores and I/O patterns
- Monitor queue lengths to detect bottlenecks
- Configure timeouts to prevent hanging jobs
- Use database persistence for production environments
