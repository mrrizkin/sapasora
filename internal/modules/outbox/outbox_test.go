package outbox

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func testScope() Scope { return Scope{TenantID: "tenant-1", WorkspaceID: "workspace-1"} }

func publishedJob(t *testing.T, store *MemoryStore, scope Scope, now time.Time, max int) SendJob {
	t.Helper()
	event := NewEvent(scope, "message.created", "message-1", []byte(`{"body":"private"}`), now)
	if err := store.AppendEvent(context.Background(), event); err != nil {
		t.Fatal(err)
	}
	if err := store.PublishEvent(context.Background(), event.Scope, event.ID); err != nil {
		t.Fatal(err)
	}
	job := NewSendJob(scope, event.ID, "future-provider", "sms", []byte(`private message`), max, now)
	if err := store.EnqueueJob(context.Background(), job); err != nil {
		t.Fatal(err)
	}
	return job
}

func TestPublicationBoundaryAndRedactedDiagnostics(t *testing.T) {
	store := NewMemoryStore()
	scope := testScope()
	now := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	event := NewEvent(scope, "message.created", "message-1", []byte(`{"body":"do-not-log"}`), now)
	if err := store.AppendEvent(context.Background(), event); err != nil {
		t.Fatal(err)
	}
	job := NewSendJob(scope, event.ID, "provider", "channel", []byte("do-not-log"), 2, now)
	if err := store.EnqueueJob(context.Background(), job); !errors.Is(err, ErrNotPublished) {
		t.Fatalf("enqueue before publication: got %v", err)
	}
	if err := store.PublishEvent(context.Background(), event.Scope, event.ID); err != nil {
		t.Fatal(err)
	}
	if err := store.PublishEvent(context.Background(), event.Scope, event.ID); err != nil {
		t.Fatalf("idempotent publication: %v", err)
	}
	stored, err := store.GetEvent(context.Background(), scope, event.ID)
	if err != nil || stored.Status != EventPublished {
		t.Fatalf("stored event = %#v, err = %v", stored, err)
	}
	if err := store.EnqueueJob(context.Background(), job); err != nil {
		t.Fatal(err)
	}
	encoded, err := json.Marshal(stored)
	if err != nil {
		t.Fatal(err)
	}
	if string(encoded) == "" || contains(string(encoded), "do-not-log") {
		t.Fatalf("event payload leaked in JSON: %s", encoded)
	}
}

func TestDiagnosticsAndCancellationBeforeSendAreRedacted(t *testing.T) {
	store := NewMemoryStore()
	scope := testScope()
	now := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	job := publishedJob(t, store, scope, now, 2)
	claimed, err := store.Claim(context.Background(), scope, "worker", now, time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.CancelJob(context.Background(), scope, job.ID, now); err != nil {
		t.Fatal(err)
	}
	var calls atomic.Int32
	worker := NewWorker(store, senderFunc(func(context.Context, SendJob) (SendResult, error) {
		calls.Add(1)
		return SendResult{Accepted: true}, nil
	}), WorkerConfig{WorkerID: "worker", Clock: func() time.Time { return now }})
	if err := worker.handle(context.Background(), claimed); !errors.Is(err, context.Canceled) {
		t.Fatalf("pre-send cancellation = %v", err)
	}
	if calls.Load() != 0 {
		t.Fatal("provider was called after cancellation")
	}
	stored, err := store.GetJob(context.Background(), scope, job.ID)
	if err != nil || stored.Status != JobCancelled {
		t.Fatalf("cancelled job = %#v, err = %v", stored, err)
	}
	if strings := fmt.Sprint(stored); strings == "" || contains(strings, "private message") {
		t.Fatalf("job payload leaked in String: %s", strings)
	}
}

func TestClaimIsExclusiveAndLeaseExpiryReclaims(t *testing.T) {
	store := NewMemoryStore()
	scope := testScope()
	now := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)
	job := publishedJob(t, store, scope, now, 3)
	otherScope := Scope{TenantID: "tenant-2", WorkspaceID: "workspace-2"}
	otherJob := publishedJob(t, store, otherScope, now, 3)
	if _, err := store.GetJob(context.Background(), scope, otherJob.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("cross-scope get: %v", err)
	}
	if scoped, err := store.Claim(context.Background(), otherScope, "scope-worker", now, time.Minute); err != nil || scoped.ID != otherJob.ID {
		t.Fatalf("scoped claim = %#v, err = %v", scoped, err)
	}

	var claimed atomic.Int32
	var wg sync.WaitGroup
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := store.Claim(context.Background(), scope, "worker", now, time.Minute); err == nil {
				claimed.Add(1)
			}
		}()
	}
	wg.Wait()
	if got := claimed.Load(); got != 1 {
		t.Fatalf("claims = %d, want exactly one", got)
	}
	if _, err := store.Claim(context.Background(), scope, "worker-2", now.Add(30*time.Second), time.Minute); !errors.Is(err, ErrNoWork) {
		t.Fatalf("claim before expiry: %v", err)
	}
	reclaimed, err := store.Claim(context.Background(), scope, "worker-2", now.Add(time.Minute), time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	if reclaimed.ID != job.ID || reclaimed.Attempt != 2 {
		t.Fatalf("reclaimed = %#v", reclaimed)
	}
	if err := store.CancelJob(context.Background(), scope, job.ID, now.Add(time.Minute)); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Claim(context.Background(), scope, "worker-3", now.Add(2*time.Minute), time.Minute); !errors.Is(err, ErrNoWork) {
		t.Fatalf("expired cancelled job claim: %v", err)
	}
	cancelled, err := store.GetJob(context.Background(), scope, job.ID)
	if err != nil || cancelled.Status != JobCancelled {
		t.Fatalf("expired cancelled job = %#v, err = %v", cancelled, err)
	}
}

func TestRetryBackoffDeadLetterAndReplay(t *testing.T) {
	store := NewMemoryStore()
	scope := testScope()
	now := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)
	job := publishedJob(t, store, scope, now, 2)
	clock := now
	var calls atomic.Int32
	sender := senderFunc(func(context.Context, SendJob) (SendResult, error) {
		calls.Add(1)
		return SendResult{}, AsRetryable(errors.New("provider unavailable"))
	})
	worker := NewWorker(store, sender, WorkerConfig{Scope: scope, WorkerID: "worker", Clock: func() time.Time { return clock }, Backoff: BackoffPolicy{
		Base: time.Second, Max: time.Minute, Jitter: func(d time.Duration) time.Duration { return d + 500*time.Millisecond },
	}})
	worked, err := worker.ProcessOnce(context.Background())
	if !worked || err != nil {
		t.Fatalf("first process: worked=%v err=%v", worked, err)
	}
	stored, _ := store.GetJob(context.Background(), scope, job.ID)
	if stored.Status != JobRetrying || !stored.NextAttemptAt.Equal(now.Add(1500*time.Millisecond)) {
		t.Fatalf("retry state = %#v", stored)
	}
	clock = stored.NextAttemptAt
	worked, err = worker.ProcessOnce(context.Background())
	if !worked || err != nil {
		t.Fatalf("second process: worked=%v err=%v", worked, err)
	}
	stored, _ = store.GetJob(context.Background(), scope, job.ID)
	if stored.Status != JobDead || stored.LastError == "" || contains(stored.LastError, "provider unavailable") {
		t.Fatalf("dead-letter state = %#v", stored)
	}
	if calls.Load() != 2 {
		t.Fatalf("provider calls = %d, want 2", calls.Load())
	}
	if err := store.ReplayDeadLetter(context.Background(), scope, job.ID, clock); err != nil {
		t.Fatal(err)
	}
	stored, _ = store.GetJob(context.Background(), scope, job.ID)
	if stored.Status != JobPending || stored.Attempt != 0 || stored.ReplayCount != 1 {
		t.Fatalf("replay state = %#v", stored)
	}
}

func TestContextCancellationBeforeAcceptanceIsTerminal(t *testing.T) {
	store := NewMemoryStore()
	scope := testScope()
	now := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	job := publishedJob(t, store, scope, now, 3)
	entered := make(chan struct{})
	release := make(chan struct{})
	sender := senderFunc(func(context.Context, SendJob) (SendResult, error) {
		close(entered)
		<-release
		return SendResult{}, errors.New("provider stopped after cancellation")
	})
	worker := NewWorker(store, sender, WorkerConfig{Scope: scope, WorkerID: "worker", Clock: func() time.Time { return now }})
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { _, err := worker.ProcessOnce(ctx); done <- err }()
	<-entered
	cancel()
	close(release)
	<-done
	stored, _ := store.GetJob(context.Background(), scope, job.ID)
	if stored.Status != JobCancelled {
		t.Fatalf("context-cancelled state = %#v", stored)
	}
}

func TestCancellationBeforeAcceptance(t *testing.T) {
	store := NewMemoryStore()
	scope := testScope()
	now := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)
	job := publishedJob(t, store, scope, now, 3)
	entered := make(chan struct{})
	release := make(chan struct{})
	sender := senderFunc(func(context.Context, SendJob) (SendResult, error) {
		close(entered)
		<-release
		return SendResult{}, context.Canceled
	})
	worker := NewWorker(store, sender, WorkerConfig{Scope: scope, WorkerID: "worker", Clock: func() time.Time { return now }})
	done := make(chan error, 1)
	go func() { _, err := worker.ProcessOnce(context.Background()); done <- err }()
	<-entered
	if err := store.CancelJob(context.Background(), scope, job.ID, now); err != nil {
		t.Fatal(err)
	}
	close(release)
	<-done
	stored, _ := store.GetJob(context.Background(), scope, job.ID)
	if stored.Status != JobCancelled {
		t.Fatalf("cancelled state = %#v", stored)
	}
}

type senderFunc func(context.Context, SendJob) (SendResult, error)

func (f senderFunc) Send(ctx context.Context, job SendJob) (SendResult, error) { return f(ctx, job) }

func contains(value, part string) bool {
	for i := 0; i+len(part) <= len(value); i++ {
		if value[i:i+len(part)] == part {
			return true
		}
	}
	return false
}
