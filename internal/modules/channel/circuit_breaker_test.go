package channel

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

type breakerTestClock struct {
	nanoseconds atomic.Int64
}

func newBreakerTestClock(at time.Time) *breakerTestClock {
	clock := &breakerTestClock{}
	clock.nanoseconds.Store(at.UnixNano())
	return clock
}

func (c *breakerTestClock) Now() time.Time {
	return time.Unix(0, c.nanoseconds.Load())
}

func (c *breakerTestClock) Advance(duration time.Duration) {
	c.nanoseconds.Add(int64(duration))
}

func TestCircuitBreakerCountsConsecutiveFailuresAndSuccessResets(t *testing.T) {
	clock := newBreakerTestClock(time.Unix(100, 0))
	breaker := NewCircuitBreaker(CircuitBreakerConfig{
		FailureThreshold: 2,
		OpenCooldown:     time.Minute,
		Clock:            clock.Now,
	})
	failure := errors.New("provider unavailable")
	calls := 0

	operation := func(context.Context) error {
		calls++
		return failure
	}
	if err := breaker.Execute(context.Background(), operation); !errors.Is(err, failure) {
		t.Fatalf("first Execute() error = %v, want %v", err, failure)
	}
	if got := breaker.State(); got != CircuitStateClosed {
		t.Fatalf("state after one failure = %q, want closed", got)
	}
	if got := breaker.ConsecutiveFailures(); got != 1 {
		t.Fatalf("failure count = %d, want 1", got)
	}

	if err := breaker.Execute(context.Background(), func(context.Context) error { calls++; return nil }); err != nil {
		t.Fatalf("successful Execute() error = %v", err)
	}
	if got := breaker.ConsecutiveFailures(); got != 0 {
		t.Fatalf("failure count after success = %d, want 0", got)
	}

	if err := breaker.Execute(context.Background(), operation); err == nil {
		t.Fatal("second failure Execute() returned nil")
	}
	if err := breaker.Execute(context.Background(), operation); err == nil {
		t.Fatal("threshold failure Execute() returned nil")
	}
	if got := breaker.State(); got != CircuitStateOpen {
		t.Fatalf("state at threshold = %q, want open", got)
	}
	if calls != 4 {
		t.Fatalf("operation calls = %d, want 4", calls)
	}
}

func TestCircuitBreakerCooldownAdmitsOneProbeAndSuccessCloses(t *testing.T) {
	clock := newBreakerTestClock(time.Unix(200, 0))
	breaker := NewCircuitBreaker(CircuitBreakerConfig{
		FailureThreshold: 1,
		OpenCooldown:     10 * time.Second,
		Clock:            clock.Now,
	})
	failure := errors.New("temporary failure")
	if err := breaker.Execute(context.Background(), func(context.Context) error { return failure }); !errors.Is(err, failure) {
		t.Fatalf("initial failure = %v", err)
	}

	var rejectedCalls atomic.Int64
	if err := breaker.Execute(context.Background(), func(context.Context) error {
		rejectedCalls.Add(1)
		return nil
	}); !IsCircuitOpen(err) {
		t.Fatalf("early call error = %v, want circuit open", err)
	}
	if rejectedCalls.Load() != 0 {
		t.Fatal("open circuit invoked operation before cooldown")
	}

	clock.Advance(10 * time.Second)
	if err := breaker.Execute(context.Background(), func(context.Context) error {
		rejectedCalls.Add(1)
		return nil
	}); err != nil {
		t.Fatalf("half-open probe error = %v", err)
	}
	if got := breaker.State(); got != CircuitStateClosed {
		t.Fatalf("state after successful probe = %q, want closed", got)
	}
	if got := breaker.ConsecutiveFailures(); got != 0 {
		t.Fatalf("failure count after successful probe = %d, want 0", got)
	}
	if rejectedCalls.Load() != 1 {
		t.Fatalf("successful calls = %d, want 1", rejectedCalls.Load())
	}
}

func TestCircuitBreakerHalfOpenProbeIsSingleUnderConcurrency(t *testing.T) {
	clock := newBreakerTestClock(time.Unix(300, 0))
	breaker := NewCircuitBreaker(CircuitBreakerConfig{
		FailureThreshold: 1,
		OpenCooldown:     time.Second,
		Clock:            clock.Now,
	})
	if err := breaker.Execute(context.Background(), func(context.Context) error { return errors.New("down") }); err == nil {
		t.Fatal("initial failure returned nil")
	}
	clock.Advance(time.Second)

	const workers = 16
	entered := make(chan struct{})
	release := make(chan struct{})
	var operationCalls atomic.Int64
	results := make(chan error, workers)
	for range workers {
		go func() {
			results <- breaker.Execute(context.Background(), func(context.Context) error {
				if operationCalls.Add(1) == 1 {
					close(entered)
				}
				<-release
				return nil
			})
		}()
	}

	select {
	case <-entered:
	case <-time.After(time.Second):
		t.Fatal("half-open probe did not start")
	}
	for i := 0; i < workers-1; i++ {
		select {
		case err := <-results:
			if !IsCircuitHalfOpen(err) {
				t.Fatalf("concurrent call %d error = %v, want half-open rejection", i, err)
			}
		case <-time.After(time.Second):
			t.Fatalf("timed out waiting for concurrent half-open rejection %d", i)
		}
	}
	if got := operationCalls.Load(); got != 1 {
		t.Fatalf("provider operation calls while probe blocked = %d, want 1", got)
	}

	close(release)
	select {
	case err := <-results:
		if err != nil {
			t.Fatalf("probe result = %v, want nil", err)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for probe result")
	}
	if got := breaker.State(); got != CircuitStateClosed {
		t.Fatalf("state after probe = %q, want closed", got)
	}
}

func TestCircuitBreakerContextAndTypedErrors(t *testing.T) {
	clock := newBreakerTestClock(time.Unix(400, 0))
	breaker := NewCircuitBreaker(CircuitBreakerConfig{FailureThreshold: 1, OpenCooldown: time.Minute, Clock: clock.Now})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	called := false
	if err := breaker.Execute(ctx, func(context.Context) error {
		called = true
		return nil
	}); !errors.Is(err, context.Canceled) {
		t.Fatalf("canceled Execute() error = %v, want context.Canceled", err)
	}
	if called {
		t.Fatal("canceled context invoked operation")
	}

	if err := breaker.Execute(context.Background(), func(context.Context) error { return errors.New("down") }); err == nil {
		t.Fatal("failure returned nil")
	}
	var openErr *CircuitOpenError
	if err := breaker.Execute(context.Background(), func(context.Context) error { t.Fatal("open operation invoked"); return nil }); !errors.As(err, &openErr) || !errors.Is(err, ErrCircuitOpen) {
		t.Fatalf("open error = %T %v, want CircuitOpenError and ErrCircuitOpen", err, err)
	}
	if err := breaker.Execute(context.Background(), nil); !errors.Is(err, ErrNilCircuitOperation) {
		t.Fatalf("nil operation error = %v, want ErrNilCircuitOperation", err)
	}

	if _, err := NewCircuitBreakerWithError(CircuitBreakerConfig{FailureThreshold: 0}); !errors.Is(err, ErrInvalidCircuitBreakerConfig) {
		t.Fatalf("invalid config error = %v, want ErrInvalidCircuitBreakerConfig", err)
	}
}
