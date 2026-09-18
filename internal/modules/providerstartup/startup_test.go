package providerstartup

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"sapasora/internal/modules/device"
)

func TestRetryUsesBoundedAttempts(t *testing.T) {
	attempts := 0
	err := Retry(context.Background(), 4, 0, 0, func() error {
		attempts++
		return errors.New("provider unavailable")
	})
	if err == nil {
		t.Fatal("Retry() error = nil, want final operation error")
	}
	if attempts != 4 {
		t.Fatalf("attempts = %d, want 4", attempts)
	}
}

func TestRetryDelayUsesExponentialCap(t *testing.T) {
	want := []time.Duration{10 * time.Millisecond, 20 * time.Millisecond, 25 * time.Millisecond}
	for retry, expected := range want {
		if got := retryDelay(10*time.Millisecond, 25*time.Millisecond, retry); got != expected {
			t.Fatalf("retryDelay(%d) = %s, want %s", retry, got, expected)
		}
	}
}

func TestRetryStopsBeforeNextAttemptWhenContextCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	attempts := 0

	err := Retry(ctx, 3, 0, 0, func() error {
		attempts++
		return errors.New("should not run")
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Retry() error = %v, want context.Canceled", err)
	}
	if attempts != 0 {
		t.Fatalf("attempts = %d, want 0", attempts)
	}
}

func TestRetryInterruptsBackoffOnContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	attempts := 0

	err := Retry(ctx, 3, time.Hour, time.Hour, func() error {
		attempts++
		cancel()
		return errors.New("provider unavailable")
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Retry() error = %v, want context.Canceled", err)
	}
	if attempts != 1 {
		t.Fatalf("attempts = %d, want 1", attempts)
	}
}

func TestRunIsolatesOneProviderFailure(t *testing.T) {
	devices := []*device.Device{{PublicID: "failed"}, {PublicID: "healthy"}}
	var mu sync.Mutex
	started := make(map[string]bool)
	var startupErrors []string

	Run(context.Background(), devices, 1,
		func(_ context.Context, d *device.Device) error {
			mu.Lock()
			started[d.PublicID] = true
			mu.Unlock()
			if d.PublicID == "failed" {
				return errors.New("provider unavailable")
			}
			return nil
		},
		func(d *device.Device, _ error) {
			mu.Lock()
			startupErrors = append(startupErrors, d.PublicID)
			mu.Unlock()
		},
	)

	if !started["failed"] || !started["healthy"] {
		t.Fatalf("startup tasks = %#v, want both providers attempted", started)
	}
	if len(startupErrors) != 1 || startupErrors[0] != "failed" {
		t.Fatalf("startup errors = %#v, want only failed provider", startupErrors)
	}
}

func TestRunBoundsConcurrentStartup(t *testing.T) {
	devices := make([]*device.Device, 8)
	for i := range devices {
		devices[i] = &device.Device{PublicID: string(rune('a' + i))}
	}

	const limit = 2
	var mu sync.Mutex
	active, maxActive := 0, 0
	Run(context.Background(), devices, limit,
		func(_ context.Context, _ *device.Device) error {
			mu.Lock()
			active++
			if active > maxActive {
				maxActive = active
			}
			mu.Unlock()
			time.Sleep(5 * time.Millisecond)
			mu.Lock()
			active--
			mu.Unlock()
			return nil
		}, nil)

	if maxActive > limit {
		t.Fatalf("max concurrent startup = %d, want <= %d", maxActive, limit)
	}
}
