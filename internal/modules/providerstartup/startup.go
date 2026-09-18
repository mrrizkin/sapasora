// Package providerstartup contains small lifecycle helpers shared by provider
// startup implementations.
package providerstartup

import (
	"context"
	"sync"

	"sapasora/internal/modules/device"
)

// StartFunc starts one provider device. A failure is reported per device and
// does not stop the remaining startup tasks.
type StartFunc func(context.Context, *device.Device) error

// ErrorFunc receives an error from one startup task.
type ErrorFunc func(*device.Device, error)

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
