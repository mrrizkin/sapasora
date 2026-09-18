package channel

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestReconnectBackoffDelayCapsExponentWithoutOverflow(t *testing.T) {
	policy := ReconnectBackoffPolicy{
		MinDelay: 10 * time.Millisecond,
		MaxDelay: 25 * time.Millisecond,
	}
	want := []time.Duration{10 * time.Millisecond, 20 * time.Millisecond, 25 * time.Millisecond, 25 * time.Millisecond}
	for retry, expected := range want {
		if got := policy.Delay(retry); got != expected {
			t.Fatalf("Delay(%d) = %s, want %s", retry, got, expected)
		}
	}

	policy = ReconnectBackoffPolicy{
		MinDelay: time.Duration(1 << 62),
		MaxDelay: time.Duration(1<<63 - 1),
	}
	if got, want := policy.Delay(2), time.Duration(1<<63-1); got != want {
		t.Fatalf("Delay(2) = %d, want capped duration %d", got, want)
	}
	if got := policy.Delay(int(^uint(0) >> 1)); got != time.Duration(1<<63-1) {
		t.Fatalf("Delay(max int) = %d, want capped duration", got)
	}
}

func TestReconnectBackoffRetryUsesMaxAttemptsAndDoesNotSleepAfterFinalAttempt(t *testing.T) {
	var sleeps []time.Duration
	policy := ReconnectBackoffPolicy{
		MinDelay:    time.Second,
		MaxDelay:    3 * time.Second,
		MaxAttempts: 3,
		Sleep: func(_ context.Context, delay time.Duration) error {
			sleeps = append(sleeps, delay)
			return nil
		},
	}

	attempts := 0
	wantErr := errors.New("unavailable")
	err := policy.Retry(context.Background(), func(context.Context) error {
		attempts++
		return wantErr
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("Retry() error = %v, want %v", err, wantErr)
	}
	if attempts != 3 {
		t.Fatalf("attempts = %d, want 3", attempts)
	}
	if !reflect.DeepEqual(sleeps, []time.Duration{time.Second, 2 * time.Second}) {
		t.Fatalf("sleeps = %v, want [1s 2s]", sleeps)
	}
}

func TestReconnectBackoffRetryStopsOnCancellationDuringBackoff(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	attempts := 0
	policy := ReconnectBackoffPolicy{
		MinDelay:    time.Hour,
		MaxDelay:    time.Hour,
		MaxAttempts: 3,
	}
	err := policy.Retry(ctx, func(context.Context) error {
		attempts++
		cancel()
		return errors.New("unavailable")
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Retry() error = %v, want context.Canceled", err)
	}
	if attempts != 1 {
		t.Fatalf("attempts = %d, want 1", attempts)
	}
}

func TestReconnectBackoffJitterIsInjectableAndBounded(t *testing.T) {
	policy := ReconnectBackoffPolicy{
		MinDelay: 10 * time.Millisecond,
		MaxDelay: 100 * time.Millisecond,
		Jitter: func(time.Duration) time.Duration {
			return 500 * time.Millisecond
		},
	}
	if got := policy.Delay(0); got != 100*time.Millisecond {
		t.Fatalf("Delay(0) with oversized jitter = %s, want 100ms", got)
	}

	policy.Jitter = func(time.Duration) time.Duration { return 0 }
	if got := policy.Delay(1); got != 10*time.Millisecond {
		t.Fatalf("Delay(1) with undersized jitter = %s, want 10ms", got)
	}
}
