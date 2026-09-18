package outbox

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// WorkerConfig controls bounded execution and lease behavior.
// ConcurrencyLimiter is a small reusable bounded-concurrency helper for
// future queue loops. Release must be called once for every successful Acquire.
type ConcurrencyLimiter struct {
	slots chan struct{}
}

func NewConcurrencyLimiter(limit int) *ConcurrencyLimiter {
	if limit <= 0 {
		limit = 1
	}
	return &ConcurrencyLimiter{slots: make(chan struct{}, limit)}
}

func (l *ConcurrencyLimiter) Acquire(ctx context.Context) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	select {
	case l.slots <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (l *ConcurrencyLimiter) Release() {
	select {
	case <-l.slots:
	default:
		panic("outbox: concurrency limiter released without acquire")
	}
}

type WorkerConfig struct {
	Scope         Scope
	WorkerID      string
	Concurrency   int
	LeaseDuration time.Duration
	PollInterval  time.Duration
	Backoff       BackoffPolicy
	Clock         func() time.Time
}

func (c WorkerConfig) withDefaults() WorkerConfig {
	if c.WorkerID == "" {
		c.WorkerID = newID("worker")
	}
	if c.Concurrency <= 0 {
		c.Concurrency = 1
	}
	if c.LeaseDuration <= 0 {
		c.LeaseDuration = time.Minute
	}
	if c.PollInterval <= 0 {
		c.PollInterval = 25 * time.Millisecond
	}
	if c.Clock == nil {
		c.Clock = time.Now
	}
	if c.Backoff.Base <= 0 {
		c.Backoff.Base = time.Second
	}
	if c.Backoff.Max <= 0 {
		c.Backoff.Max = 5 * time.Minute
	}
	return c
}

// Worker claims jobs, limits in-flight provider calls, and waits for active
// calls during graceful shutdown. It never logs or exposes job payloads.
type Worker struct {
	store  EventStore
	sender Sender
	config WorkerConfig

	stateMu       sync.Mutex
	running       bool
	stopRequested bool
	stop          chan struct{}
	done          chan struct{}
}

func NewWorker(store EventStore, sender Sender, config WorkerConfig) *Worker {
	return &Worker{store: store, sender: sender, config: config.withDefaults()}
}

// ProcessOnce claims and handles at most one job. It is convenient for a
// deterministic scheduler and for tests.
func (w *Worker) ProcessOnce(ctx context.Context) (bool, error) {
	if w.store == nil || w.sender == nil {
		return false, errors.New("outbox: worker requires store and sender")
	}
	if err := contextErr(ctx); err != nil {
		return false, err
	}
	job, err := w.store.Claim(ctx, w.config.Scope, w.config.WorkerID, w.config.Clock(), w.config.LeaseDuration)
	if errors.Is(err, ErrNoWork) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, w.handle(ctx, job)
}

func (w *Worker) handle(ctx context.Context, job SendJob) error {
	// Re-read the claim before crossing the provider boundary. This closes the
	// common cancellation-before-send race; cancellation that arrives after
	// this check is resolved by the provider's Accepted result below.
	current, err := w.store.GetJob(context.Background(), job.Scope, job.ID)
	if err != nil {
		return err
	}
	if current.Status == JobCancelled {
		return context.Canceled
	}
	if current.CancelRequested {
		now := w.config.Clock()
		if err := w.store.RetryJob(context.Background(), job.Scope, job.ID, w.config.WorkerID, now, now, context.Canceled); err != nil {
			return err
		}
		return context.Canceled
	}

	// Check cancellation before crossing the provider boundary. Once Send is
	// called, only the provider's Accepted result closes that boundary.
	if err := contextErr(ctx); err != nil {
		now := w.config.Clock()
		_ = w.store.CancelJob(context.Background(), job.Scope, job.ID, now)
		_ = w.store.RetryJob(context.Background(), job.Scope, job.ID, w.config.WorkerID, now, now, err)
		return err
	}

	result, sendErr := w.sender.Send(ctx, cloneJob(job))
	now := w.config.Clock()
	if result.Accepted {
		return w.store.AcceptJob(context.Background(), job.Scope, job.ID, w.config.WorkerID, now)
	}
	if sendErr == nil {
		sendErr = errors.New("outbox: provider did not accept job")
	}
	if ctx.Err() != nil || errors.Is(sendErr, context.Canceled) || errors.Is(sendErr, context.DeadlineExceeded) {
		// Cancellation before acceptance is terminal and must not be retried.
		// Record the request even when the provider returns a different error
		// after observing cancellation; otherwise RetryJob would leave the job
		// retrying because CancelRequested was never set.
		_ = w.store.CancelJob(context.Background(), job.Scope, job.ID, now)
		_ = w.store.RetryJob(context.Background(), job.Scope, job.ID, w.config.WorkerID, now, now, sendErr)
		return sendErr
	}
	if Classify(sendErr) == FailureRetryable && job.Attempt < job.MaxAttempts {
		next := now.Add(w.config.Backoff.Delay(job.Attempt))
		return w.store.RetryJob(context.Background(), job.Scope, job.ID, w.config.WorkerID, now, next, sendErr)
	}
	return w.store.DeadLetterJob(context.Background(), job.Scope, job.ID, w.config.WorkerID, now, sendErr)
}

// Run starts a loop and returns after ctx is cancelled and all in-flight jobs
// finish. Existing provider calls receive a non-cancelled child context so
// shutdown is graceful; an operator can cancel a specific job via CancelJob.
func (w *Worker) Run(ctx context.Context) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	if w.store == nil || w.sender == nil {
		return errors.New("outbox: worker requires store and sender")
	}
	w.stateMu.Lock()
	if w.running {
		w.stateMu.Unlock()
		return errors.New("outbox: worker already running")
	}
	w.running = true
	w.stopRequested = false
	w.stop = make(chan struct{})
	w.done = make(chan struct{})
	stop := w.stop
	done := w.done
	w.stateMu.Unlock()
	defer func() {
		w.stateMu.Lock()
		w.running = false
		close(done)
		w.stateMu.Unlock()
	}()

	semaphore := make(chan struct{}, w.config.Concurrency)
	var active sync.WaitGroup
	for {
		select {
		case <-ctx.Done():
			active.Wait()
			return nil
		case <-stop:
			active.Wait()
			return nil
		default:
		}

		select {
		case semaphore <- struct{}{}:
		case <-ctx.Done():
			active.Wait()
			return nil
		case <-stop:
			active.Wait()
			return nil
		}

		job, err := w.store.Claim(ctx, w.config.Scope, w.config.WorkerID, w.config.Clock(), w.config.LeaseDuration)
		if errors.Is(err, ErrNoWork) {
			<-semaphore
			if !waitOrStop(ctx, stop, w.config.PollInterval) {
				active.Wait()
				return nil
			}
			continue
		}
		if err != nil {
			<-semaphore
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				active.Wait()
				return nil
			}
			active.Wait()
			return fmt.Errorf("outbox: claim: %w", err)
		}

		active.Add(1)
		go func() {
			defer active.Done()
			defer func() { <-semaphore }()
			// Keep values from the caller's context, but do not abort a
			// provider call merely because intake has begun shutting down.
			jobContext := context.WithoutCancel(ctx)
			_ = w.handle(jobContext, job)
		}()
	}
}

func waitOrStop(ctx context.Context, stop <-chan struct{}, delay time.Duration) bool {
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-timer.C:
		return true
	case <-ctx.Done():
		return false
	case <-stop:
		return false
	}
}

// Shutdown stops claiming new jobs and waits for in-flight calls. If its
// context expires, it returns that error while the worker finishes in the
// background and keeps its lease state intact.
func (w *Worker) Shutdown(ctx context.Context) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	w.stateMu.Lock()
	if !w.running {
		w.stateMu.Unlock()
		return nil
	}
	if !w.stopRequested {
		close(w.stop)
		w.stopRequested = true
	}
	done := w.done
	w.stateMu.Unlock()
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
