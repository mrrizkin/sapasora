package scheduler

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"sapasora/platform/logger"

	_ "github.com/lib/pq"           // PostgreSQL driver
	_ "github.com/mattn/go-sqlite3" // SQLite driver
	// Add other database drivers as needed
)

// DBJobStorage is a database implementation of JobStorage
type DBJobStorage struct {
	db  *sql.DB
	log *logger.Logger
}

// NewDBJobStorage creates a new database-based job storage
func NewDBJobStorage(db *sql.DB, log *logger.Logger) *DBJobStorage {
	storage := &DBJobStorage{
		db:  db,
		log: log,
	}

	// Ensure the table exists
	storage.createTable()

	return storage
}

// createTable creates the jobs table if it doesn't exist
func (s *DBJobStorage) createTable() {
	query := `
	CREATE TABLE IF NOT EXISTS scheduler_jobs (
		id TEXT PRIMARY KEY,
		queue TEXT NOT NULL,
		class TEXT NOT NULL,
		args TEXT,
		created_at TIMESTAMP NOT NULL,
		enqueued_at TIMESTAMP,
		started_at TIMESTAMP,
		finished_at TIMESTAMP,
		failed_at TIMESTAMP,
		error TEXT,
		retry_count INTEGER DEFAULT 0,
		max_retries INTEGER DEFAULT 25,
		backtrace TEXT,
		payload TEXT
	);
	`

	_, err := s.db.Exec(query)
	if err != nil {
		s.log.Error("Failed to create scheduler_jobs table", "error", err)
	}
}

// StoreJob stores a job in the database
func (s *DBJobStorage) StoreJob(ctx context.Context, job *Job) error {
	argsBytes, err := json.Marshal(job.Args)
	if err != nil {
		return fmt.Errorf("failed to marshal job args: %w", err)
	}

	payloadBytes, err := json.Marshal(job.Payload)
	if err != nil {
		return fmt.Errorf("failed to marshal job payload: %w", err)
	}

	var backtraceBytes []byte
	if job.Backtrace != nil {
		backtraceBytes, err = json.Marshal(job.Backtrace)
		if err != nil {
			return fmt.Errorf("failed to marshal job backtrace: %w", err)
		}
	}

	query := `
	INSERT INTO scheduler_jobs (
		id, queue, class, args, created_at, enqueued_at, started_at,
		finished_at, failed_at, error, retry_count, max_retries, backtrace, payload
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`

	_, err = s.db.Exec(query,
		job.ID,
		job.Queue,
		job.Class,
		string(argsBytes),
		job.CreatedAt,
		job.EnqueuedAt,
		job.StartedAt,
		job.FinishedAt,
		job.FailedAt,
		job.Error,
		job.RetryCount,
		job.MaxRetries,
		string(backtraceBytes),
		string(payloadBytes),
	)

	if err != nil {
		s.log.Error("Failed to store job", "job_id", job.ID, "error", err)
		return fmt.Errorf("failed to store job: %w", err)
	}

	s.log.Info("Job stored in database", "job_id", job.ID)
	return nil
}

// GetJob retrieves a job by ID from the database
func (s *DBJobStorage) GetJob(ctx context.Context, jobID string) (*Job, error) {
	query := `
	SELECT id, queue, class, args, created_at, enqueued_at, started_at,
		   finished_at, failed_at, error, retry_count, max_retries, backtrace, payload
	FROM scheduler_jobs
	WHERE id = $1
	`

	row := s.db.QueryRow(query, jobID)

	var job Job
	var argsBytes, payloadBytes, backtraceBytes sql.NullString
	var enqueuedAt, startedAt, finishedAt, failedAt sql.NullTime

	err := row.Scan(
		&job.ID,
		&job.Queue,
		&job.Class,
		&argsBytes,
		&job.CreatedAt,
		&enqueuedAt,
		&startedAt,
		&finishedAt,
		&failedAt,
		&job.Error,
		&job.RetryCount,
		&job.MaxRetries,
		&backtraceBytes,
		&payloadBytes,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("job %s not found", jobID)
		}
		return nil, fmt.Errorf("failed to get job: %w", err)
	}

	// Set nullable fields
	if enqueuedAt.Valid {
		job.EnqueuedAt = &enqueuedAt.Time
	}
	if startedAt.Valid {
		job.StartedAt = &startedAt.Time
	}
	if finishedAt.Valid {
		job.FinishedAt = &finishedAt.Time
	}
	if failedAt.Valid {
		job.FailedAt = &failedAt.Time
	}

	// Unmarshal JSON fields
	if argsBytes.Valid {
		if err := json.Unmarshal([]byte(argsBytes.String), &job.Args); err != nil {
			return nil, fmt.Errorf("failed to unmarshal job args: %w", err)
		}
	} else {
		job.Args = []any{} // Default to empty slice
	}

	if payloadBytes.Valid {
		if err := json.Unmarshal([]byte(payloadBytes.String), &job.Payload); err != nil {
			return nil, fmt.Errorf("failed to unmarshal job payload: %w", err)
		}
	} else {
		job.Payload = make(map[string]any) // Default to empty map
	}

	if backtraceBytes.Valid {
		var backtrace []string
		if err := json.Unmarshal([]byte(backtraceBytes.String), &backtrace); err != nil {
			return nil, fmt.Errorf("failed to unmarshal job backtrace: %w", err)
		}
		job.Backtrace = backtrace
	}

	return &job, nil
}

// DeleteJob deletes a job by ID from the database
func (s *DBJobStorage) DeleteJob(ctx context.Context, jobID string) error {
	query := "DELETE FROM scheduler_jobs WHERE id = $1"
	result, err := s.db.Exec(query, jobID)
	if err != nil {
		s.log.Error("Failed to delete job", "job_id", jobID, "error", err)
		return fmt.Errorf("failed to delete job: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("job %s not found", jobID)
	}

	s.log.Info("Job deleted from database", "job_id", jobID)
	return nil
}

// UpdateJob updates an existing job in the database
func (s *DBJobStorage) UpdateJob(ctx context.Context, job *Job) error {
	argsBytes, err := json.Marshal(job.Args)
	if err != nil {
		return fmt.Errorf("failed to marshal job args: %w", err)
	}

	payloadBytes, err := json.Marshal(job.Payload)
	if err != nil {
		return fmt.Errorf("failed to marshal job payload: %w", err)
	}

	var backtraceBytes []byte
	if job.Backtrace != nil {
		backtraceBytes, err = json.Marshal(job.Backtrace)
		if err != nil {
			return fmt.Errorf("failed to marshal job backtrace: %w", err)
		}
	}

	query := `
	UPDATE scheduler_jobs
	SET queue = $2, class = $3, args = $4, enqueued_at = $5, started_at = $6,
		finished_at = $7, failed_at = $8, error = $9, retry_count = $10,
		max_retries = $11, backtrace = $12, payload = $13
	WHERE id = $1
	`

	_, err = s.db.Exec(query,
		job.ID,
		job.Queue,
		job.Class,
		string(argsBytes),
		job.EnqueuedAt,
		job.StartedAt,
		job.FinishedAt,
		job.FailedAt,
		job.Error,
		job.RetryCount,
		job.MaxRetries,
		string(backtraceBytes),
		string(payloadBytes),
	)

	if err != nil {
		s.log.Error("Failed to update job", "job_id", job.ID, "error", err)
		return fmt.Errorf("failed to update job: %w", err)
	}

	s.log.Info("Job updated in database", "job_id", job.ID)
	return nil
}

// GetAllJobs returns all jobs from the database
func (s *DBJobStorage) GetAllJobs(ctx context.Context) ([]*Job, error) {
	query := `
	SELECT id, queue, class, args, created_at, enqueued_at, started_at,
		   finished_at, failed_at, error, retry_count, max_retries, backtrace, payload
	FROM scheduler_jobs
	ORDER BY created_at DESC
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all jobs: %w", err)
	}
	defer rows.Close()

	var jobs []*Job
	for rows.Next() {
		var job Job
		var argsBytes, payloadBytes, backtraceBytes sql.NullString
		var enqueuedAt, startedAt, finishedAt, failedAt sql.NullTime

		err := rows.Scan(
			&job.ID,
			&job.Queue,
			&job.Class,
			&argsBytes,
			&job.CreatedAt,
			&enqueuedAt,
			&startedAt,
			&finishedAt,
			&failedAt,
			&job.Error,
			&job.RetryCount,
			&job.MaxRetries,
			&backtraceBytes,
			&payloadBytes,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan job: %w", err)
		}

		// Set nullable fields
		if enqueuedAt.Valid {
			job.EnqueuedAt = &enqueuedAt.Time
		}
		if startedAt.Valid {
			job.StartedAt = &startedAt.Time
		}
		if finishedAt.Valid {
			job.FinishedAt = &finishedAt.Time
		}
		if failedAt.Valid {
			job.FailedAt = &failedAt.Time
		}

		// Unmarshal JSON fields
		if argsBytes.Valid {
			if err := json.Unmarshal([]byte(argsBytes.String), &job.Args); err != nil {
				s.log.Error("Failed to unmarshal job args", "job_id", job.ID, "error", err)
				continue // Skip this job, but continue with others
			}
		} else {
			job.Args = []any{} // Default to empty slice
		}

		if payloadBytes.Valid {
			if err := json.Unmarshal([]byte(payloadBytes.String), &job.Payload); err != nil {
				s.log.Error("Failed to unmarshal job payload", "job_id", job.ID, "error", err)
				continue // Skip this job, but continue with others
			}
		} else {
			job.Payload = make(map[string]any) // Default to empty map
		}

		if backtraceBytes.Valid {
			var backtrace []string
			if err := json.Unmarshal([]byte(backtraceBytes.String), &backtrace); err != nil {
				s.log.Error("Failed to unmarshal job backtrace", "job_id", job.ID, "error", err)
				continue // Skip this job, but continue with others
			}
			job.Backtrace = backtrace
		}

		jobs = append(jobs, &job)
	}

	return jobs, nil
}

// GetJobsByQueue returns all jobs for a specific queue
func (s *DBJobStorage) GetJobsByQueue(ctx context.Context, queue string) ([]*Job, error) {
	query := `
	SELECT id, queue, class, args, created_at, enqueued_at, started_at,
		   finished_at, failed_at, error, retry_count, max_retries, backtrace, payload
	FROM scheduler_jobs
	WHERE queue = $1
	ORDER BY created_at DESC
	`

	rows, err := s.db.Query(query, queue)
	if err != nil {
		return nil, fmt.Errorf("failed to get jobs by queue: %w", err)
	}
	defer rows.Close()

	var jobs []*Job
	for rows.Next() {
		var job Job
		var argsBytes, payloadBytes, backtraceBytes sql.NullString
		var enqueuedAt, startedAt, finishedAt, failedAt sql.NullTime

		err := rows.Scan(
			&job.ID,
			&job.Queue,
			&job.Class,
			&argsBytes,
			&job.CreatedAt,
			&enqueuedAt,
			&startedAt,
			&finishedAt,
			&failedAt,
			&job.Error,
			&job.RetryCount,
			&job.MaxRetries,
			&backtraceBytes,
			&payloadBytes,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan job: %w", err)
		}

		// Set nullable fields
		if enqueuedAt.Valid {
			job.EnqueuedAt = &enqueuedAt.Time
		}
		if startedAt.Valid {
			job.StartedAt = &startedAt.Time
		}
		if finishedAt.Valid {
			job.FinishedAt = &finishedAt.Time
		}
		if failedAt.Valid {
			job.FailedAt = &failedAt.Time
		}

		// Unmarshal JSON fields
		if argsBytes.Valid {
			if err := json.Unmarshal([]byte(argsBytes.String), &job.Args); err != nil {
				s.log.Error("Failed to unmarshal job args", "job_id", job.ID, "error", err)
				continue // Skip this job, but continue with others
			}
		} else {
			job.Args = []any{} // Default to empty slice
		}

		if payloadBytes.Valid {
			if err := json.Unmarshal([]byte(payloadBytes.String), &job.Payload); err != nil {
				s.log.Error("Failed to unmarshal job payload", "job_id", job.ID, "error", err)
				continue // Skip this job, but continue with others
			}
		} else {
			job.Payload = make(map[string]any) // Default to empty map
		}

		if backtraceBytes.Valid {
			var backtrace []string
			if err := json.Unmarshal([]byte(backtraceBytes.String), &backtrace); err != nil {
				s.log.Error("Failed to unmarshal job backtrace", "job_id", job.ID, "error", err)
				continue // Skip this job, but continue with others
			}
			job.Backtrace = backtrace
		}

		jobs = append(jobs, &job)
	}

	return jobs, nil
}

// GetFailedJobs returns all failed jobs
func (s *DBJobStorage) GetFailedJobs(ctx context.Context) ([]*Job, error) {
	query := `
	SELECT id, queue, class, args, created_at, enqueued_at, started_at,
		   finished_at, failed_at, error, retry_count, max_retries, backtrace, payload
	FROM scheduler_jobs
	WHERE error IS NOT NULL OR failed_at IS NOT NULL
	ORDER BY created_at DESC
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get failed jobs: %w", err)
	}
	defer rows.Close()

	var jobs []*Job
	for rows.Next() {
		var job Job
		var argsBytes, payloadBytes, backtraceBytes sql.NullString
		var enqueuedAt, startedAt, finishedAt, failedAt sql.NullTime

		err := rows.Scan(
			&job.ID,
			&job.Queue,
			&job.Class,
			&argsBytes,
			&job.CreatedAt,
			&enqueuedAt,
			&startedAt,
			&finishedAt,
			&failedAt,
			&job.Error,
			&job.RetryCount,
			&job.MaxRetries,
			&backtraceBytes,
			&payloadBytes,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan job: %w", err)
		}

		// Set nullable fields
		if enqueuedAt.Valid {
			job.EnqueuedAt = &enqueuedAt.Time
		}
		if startedAt.Valid {
			job.StartedAt = &startedAt.Time
		}
		if finishedAt.Valid {
			job.FinishedAt = &finishedAt.Time
		}
		if failedAt.Valid {
			job.FailedAt = &failedAt.Time
		}

		// Unmarshal JSON fields
		if argsBytes.Valid {
			if err := json.Unmarshal([]byte(argsBytes.String), &job.Args); err != nil {
				s.log.Error("Failed to unmarshal job args", "job_id", job.ID, "error", err)
				continue // Skip this job, but continue with others
			}
		} else {
			job.Args = []any{} // Default to empty slice
		}

		if payloadBytes.Valid {
			if err := json.Unmarshal([]byte(payloadBytes.String), &job.Payload); err != nil {
				s.log.Error("Failed to unmarshal job payload", "job_id", job.ID, "error", err)
				continue // Skip this job, but continue with others
			}
		} else {
			job.Payload = make(map[string]any) // Default to empty map
		}

		if backtraceBytes.Valid {
			var backtrace []string
			if err := json.Unmarshal([]byte(backtraceBytes.String), &backtrace); err != nil {
				s.log.Error("Failed to unmarshal job backtrace", "job_id", job.ID, "error", err)
				continue // Skip this job, but continue with others
			}
			job.Backtrace = backtrace
		}

		jobs = append(jobs, &job)
	}

	return jobs, nil
}

// CleanOldJobs removes jobs older than the specified duration
func (s *DBJobStorage) CleanOldJobs(ctx context.Context, olderThan time.Duration) error {
	cutoffTime := time.Now().Add(-olderThan)
	query := "DELETE FROM scheduler_jobs WHERE finished_at < $1"

	result, err := s.db.Exec(query, cutoffTime)
	if err != nil {
		s.log.Error("Failed to clean old jobs", "error", err)
		return fmt.Errorf("failed to clean old jobs: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	s.log.Info("Cleaned up old jobs", "count", rowsAffected)
	return nil
}

// RedisJobStorage is a Redis-based implementation of JobStorage
// This would require importing a Redis client like go-redis/redis
type RedisJobStorage struct {
	// redisClient *redis.Client
	log *logger.Logger
}

// NewRedisJobStorage creates a new Redis-based job storage
// func NewRedisJobStorage(redisClient *redis.Client, log *logger.Logger) *RedisJobStorage {
// 	return &RedisJobStorage{
// 		redisClient: redisClient,
// 		log:         log,
// 	}
// }

// For now, we'll keep the interface and implementation separate
// The Redis implementation would go here when Redis support is added to the project
