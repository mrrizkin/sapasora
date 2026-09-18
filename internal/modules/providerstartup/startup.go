// Package providerstartup contains small lifecycle helpers shared by provider
// startup implementations.
package providerstartup

import (
	"context"
	"sync"
	"time"

	"sapasora/internal/modules/device"
)

const DefaultStartupConcurrency = 4

// StartFunc starts one provider device. A failure is reported per device and
// does not stop the remaining startup tasks.
type StartFunc func(context.Context, *device.Device) error

// ErrorFunc receives an error from one startup task.
type ErrorFunc func(*device.Device, error)

// Retry runs one device connection attempt with bounded exponential backoff.
// Each caller owns its retry loop, so devices do not delay one another. The
// operation is never started after ctx is canceled, and waiting for the next
// attempt is interruptible by ctx.
func Retry(
	ctx context.Context,
	maxAttempts int,
	initialDelay time.Duration,
	maxDelay time.Duration,
	operation func() error,
) error {
	if maxAttempts < 1 {
		maxAttempts = 1
	}
	if maxDelay < 0 {
		maxDelay = 0
	}

	var lastErr error
	for attempt := 0; attempt < maxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}

		lastErr = operation()
		if lastErr == nil {
			return nil
		}
		if attempt == maxAttempts-1 {
			return lastErr
		}

		delay := retryDelay(initialDelay, maxDelay, attempt)
		if delay <= 0 {
			continue
		}

		timer := time.NewTimer(delay)
		select {
		case <-timer.C:
		case <-ctx.Done():
			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}
			return ctx.Err()
		}
	}

	return lastErr
}

func retryDelay(initialDelay, maxDelay time.Duration, retry int) time.Duration {
	if initialDelay <= 0 || maxDelay == 0 {
		return 0
	}
	if initialDelay > maxDelay {
		initialDelay = maxDelay
	}

	delay := initialDelay
	for range retry {
		if delay >= maxDelay-delay {
			return maxDelay
		}
		delay *= 2
	}
	return delay
}

// Run starts tasks with a bounded number of workers. It waits for the startup
// callbacks to return, but deliberately does not return a task error: one
// unavailable provider must not prevent the application from starting.
func Run(
	ctx context.Context,
	devices []*device.Device,
	limit int,
	start StartFunc,
	onError ErrorFunc,
) {
	if len(devices) == 0 {
		return
	}
	if limit < 1 {
		limit = 1
	}
	if limit > len(devices) {
		limit = len(devices)
	}

	jobs := make(chan *device.Device)
	var workers sync.WaitGroup
	workers.Add(limit)
	for range limit {
		go func() {
			defer workers.Done()
			for d := range jobs {
				if ctx.Err() != nil {
					return
				}
				if err := start(ctx, d); err != nil && onError != nil {
					onError(d, err)
				}
			}
		}()
	}

	for _, d := range devices {
		select {
		case jobs <- d:
		case <-ctx.Done():
			close(jobs)
			workers.Wait()
			return
		}
	}
	close(jobs)
	workers.Wait()
}
