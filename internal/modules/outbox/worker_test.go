package outbox

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestConcurrencyLimiterBoundsAndHonorsCancellation(t *testing.T) {
	limiter := NewConcurrencyLimiter(1)
	if err := limiter.Acquire(context.Background()); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := limiter.Acquire(ctx); err == nil {
		t.Fatal("acquire succeeded after cancellation")
	}
	limiter.Release()
}

func TestWorkerGracefulShutdownAndBoundedConcurrency(t *testing.T) {
	store := NewMemoryStore()
	scope := testScope()
	now := time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < 8; i++ {
		publishedJob(t, store, scope, now, 2)
	}

	entered := make(chan struct{}, 8)
	release := make(chan struct{})
	var active atomic.Int32
	var maximum atomic.Int32
	sender := senderFunc(func(context.Context, SendJob) (SendResult, error) {
		current := active.Add(1)
		for {
			old := maximum.Load()
			if current <= old || maximum.CompareAndSwap(old, current) {
				break
			}
		}
		entered <- struct{}{}
		<-release
		active.Add(-1)
		return SendResult{Accepted: true}, nil
	})
	worker := NewWorker(store, sender, WorkerConfig{
		Scope:         scope,
		WorkerID:      "bounded-worker",
		Concurrency:   2,
		LeaseDuration: time.Minute,
		PollInterval:  time.Millisecond,
		Clock:         func() time.Time { return now },
	})
	runDone := make(chan error, 1)
	go func() { runDone <- worker.Run(context.Background()) }()
	for i := 0; i < 2; i++ {
		select {
		case <-entered:
		case <-time.After(time.Second):
			t.Fatal("worker did not fill its concurrency bound")
		}
	}
	if got := maximum.Load(); got != 2 {
		t.Fatalf("maximum concurrency = %d, want 2", got)
	}
	shutdownDone := make(chan error, 1)
	go func() { shutdownDone <- worker.Shutdown(context.Background()) }()
	close(release)
	if err := <-shutdownDone; err != nil {
		t.Fatal(err)
	}
	if err := <-runDone; err != nil {
		t.Fatal(err)
	}
	if active.Load() != 0 {
		t.Fatal("worker shutdown returned with active provider calls")
	}
}

func TestBackoffIsDeterministicWithInjectedJitter(t *testing.T) {
	var seen []time.Duration
	var mu sync.Mutex
	policy := BackoffPolicy{Base: 2 * time.Second, Max: 10 * time.Second, Jitter: func(value time.Duration) time.Duration {
		mu.Lock()
		seen = append(seen, value)
		mu.Unlock()
		return value / 2
	}}
	if got := policy.Delay(3); got != 4*time.Second {
		t.Fatalf("delay = %s, want 4s", got)
	}
	if len(seen) != 1 || seen[0] != 8*time.Second {
		t.Fatalf("jitter input = %v", seen)
	}
}
