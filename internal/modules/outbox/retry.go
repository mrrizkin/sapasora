package outbox

import (
	"context"
	"errors"
	"math/rand"
	"net"
	"sync"
	"time"
)

// FailureClass determines whether a failed provider call should be retried.
type FailureClass uint8

const (
	FailureUnknown FailureClass = iota
	FailureRetryable
	FailurePermanent
)

// RetryableError and PermanentError let provider adapters classify failures
// without coupling the outbox to a provider SDK.
type RetryableError struct{ Err error }

func (e RetryableError) Error() string {
	if e.Err == nil {
		return "outbox: retryable error"
	}
	return e.Err.Error()
}
func (e RetryableError) Unwrap() error   { return e.Err }
func (e RetryableError) Retryable() bool { return true }

// PermanentError explicitly prevents retries.
type PermanentError struct{ Err error }

func (e PermanentError) Error() string {
	if e.Err == nil {
		return "outbox: permanent error"
	}
	return e.Err.Error()
}
func (e PermanentError) Unwrap() error   { return e.Err }
func (e PermanentError) Retryable() bool { return false }

func AsRetryable(err error) error {
	if err == nil {
		return nil
	}
	return RetryableError{Err: err}
}

func AsPermanent(err error) error {
	if err == nil {
		return nil
	}
	return PermanentError{Err: err}
}

// Classify applies conservative defaults: explicit retryability and temporary
// network errors retry; cancellation, deadlines, and unknown errors do not.
func Classify(err error) FailureClass {
	if err == nil {
		return FailureUnknown
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return FailurePermanent
	}
	var marked interface{ Retryable() bool }
	if errors.As(err, &marked) {
		if marked.Retryable() {
			return FailureRetryable
		}
		return FailurePermanent
	}
	var networkErr net.Error
	if errors.As(err, &networkErr) && (networkErr.Timeout() || networkErr.Temporary()) {
		return FailureRetryable
	}
	return FailurePermanent
}

// JitterFunc is injected into BackoffPolicy so tests can be deterministic and
// deployments can choose a randomization strategy appropriate to their queues.
type JitterFunc func(time.Duration) time.Duration

// BackoffPolicy calculates the delay before retry attempt n (n starts at 1).
type BackoffPolicy struct {
	Base   time.Duration
	Max    time.Duration
	Jitter JitterFunc
}

func (p BackoffPolicy) Delay(attempt int) time.Duration {
	if attempt < 1 {
		attempt = 1
	}
	base := p.Base
	if base <= 0 {
		base = time.Second
	}
	max := p.Max
	if max <= 0 {
		max = 5 * time.Minute
	}
	delay := base
	for i := 1; i < attempt && delay < max; i++ {
		if delay > max/2 {
			delay = max
			break
		}
		delay *= 2
	}
	if delay > max {
		delay = max
	}
	if p.Jitter != nil {
		delay = p.Jitter(delay)
		if delay < 0 {
			delay = 0
		}
		if delay > max {
			delay = max
		}
	}
	return delay
}

// FullJitter returns a production-friendly jitter function. Tests should pass
// their own function instead of using this random helper.
func FullJitter(source *rand.Rand) JitterFunc {
	if source == nil {
		source = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	// rand.Rand is not safe for concurrent use. A single policy is commonly
	// shared by all worker goroutines, so serialize access to the injected
	// source rather than making callers coordinate it themselves.
	var mu sync.Mutex
	return func(max time.Duration) time.Duration {
		if max <= 0 {
			return 0
		}
		mu.Lock()
		defer mu.Unlock()
		if max == time.Duration(^uint64(0)>>1) {
			return time.Duration(source.Int63())
		}
		return time.Duration(source.Int63n(int64(max) + 1))
	}
}

// SafeDuration guards callers that derive policy values from untrusted config.
func SafeDuration(value time.Duration) time.Duration {
	if value < 0 {
		return 0
	}
	return value
}
