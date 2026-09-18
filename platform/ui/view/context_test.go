package view

import (
	"context"
	"testing"
	"time"
)

func TestCombineContextsWaitsForCancellation(t *testing.T) {
	first, cancelFirst := context.WithCancel(context.Background())
	defer cancelFirst()
	second, cancelSecond := context.WithCancel(context.Background())
	defer cancelSecond()

	combined := CombineContexts(first, second)
	select {
	case <-combined.Done():
		t.Fatal("combined context canceled before a source context")
	case <-time.After(20 * time.Millisecond):
	}

	cancelSecond()
	select {
	case <-combined.Done():
	case <-time.After(time.Second):
		t.Fatal("combined context did not cancel")
	}
	if combined.Err() != context.Canceled {
		t.Fatalf("combined error = %v, want %v", combined.Err(), context.Canceled)
	}
}

func TestCombineContextsUsesFirstCancellationError(t *testing.T) {
	first, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	second := context.Background()

	combined := CombineContexts(first, second)
	select {
	case <-combined.Done():
	case <-time.After(time.Second):
		t.Fatal("combined context did not cancel after deadline")
	}
	if combined.Err() != context.DeadlineExceeded {
		t.Fatalf("combined error = %v, want %v", combined.Err(), context.DeadlineExceeded)
	}
}
