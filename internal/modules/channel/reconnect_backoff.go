package channel

import (
	"context"
	"errors"
	"time"
)

// ErrNilOperation is returned when a retry policy is asked to run a nil
// operation.
var ErrNilOperation = errors.New("channel reconnect operation is nil")

// JitterFunc adjusts a calculated backoff delay. A nil JitterFunc disables
// jitter, which keeps the default policy deterministic. The returned delay is
// bounded by the policy's effective minimum and maximum delays.
type JitterFunc func(time.Duration) time.Duration

// SleepFunc waits for a backoff delay. It is optional; a nil SleepFunc uses a
// context-aware timer. Supplying one is useful for an application scheduler or
// deterministic tests, and implementations should honor ctx cancellation.
type SleepFunc func(context.Context, time.Duration) error

// ReconnectBackoffPolicy describes bounded retries for a provider-neutral
// channel connection. MaxAttempts counts operation calls, not waits. A wait is
// made only between failed attempts, never after the final attempt.
//
// MinDelay and MaxDelay are normalized when the policy is used: negative
// values disable that part of the delay, MaxDelay <= 0 disables waiting, and a
// MinDelay greater than MaxDelay is reduced to MaxDelay. MaxAttempts <= 0 is
// treated as one attempt so an invalid zero value still makes one operation
// call.
type ReconnectBackoffPolicy struct {
	MinDelay    time.Duration
	MaxDelay    time.Duration
	MaxAttempts int
	Jitter      JitterFunc
	Sleep       SleepFunc
}

// BackoffPolicy is the concise name for ReconnectBackoffPolicy.
type BackoffPolicy = ReconnectBackoffPolicy

// NewReconnectBackoffPolicy creates a policy with deterministic delays unless
// an optional jitter function is supplied. The policy remains value-like and
// can be further customized by setting its exported fields.
func NewReconnectBackoffPolicy(minDelay, maxDelay time.Duration, maxAttempts int, jitter ...JitterFunc) ReconnectBackoffPolicy {
	policy := ReconnectBackoffPolicy{
		MinDelay:    minDelay,
		MaxDelay:    maxDelay,
		MaxAttempts: maxAttempts,
	}
	if len(jitter) > 0 {
		policy.Jitter = jitter[0]
	}
	return policy
}

// Delay returns the delay before the attempt following retry. retry is
// zero-based: the first failed attempt uses MinDelay, the next uses twice the
// delay, and so on. The calculation stops as soon as the maximum is reached,
// avoiding duration overflow even for very large retry values.
func (p ReconnectBackoffPolicy) Delay(retry int) time.Duration {
	minimum, maximum := p.delayBounds()
	if retry < 0 || minimum <= 0 || maximum <= 0 {
		return 0
	}

	delay := minimum
	for retry > 0 && delay < maximum {
		if delay > maximum/2 {
			delay = maximum
			break
		}
		delay *= 2
		retry--
	}

	if p.Jitter != nil {
		delay = p.Jitter(delay)
		if delay < minimum {
			delay = minimum
		}
		if delay > maximum {
			delay = maximum
		}
	}
	return delay
}

// Retry executes operation until it succeeds, the attempt limit is reached,
// or ctx is canceled. The context is checked before every operation and while
// waiting between attempts. If all attempts fail, the final operation error is
// returned.
func (p ReconnectBackoffPolicy) Retry(ctx context.Context, operation func(context.Context) error) error {
	if operation == nil {
		return ErrNilOperation
	}

	attempts := p.MaxAttempts
	if attempts <= 0 {
		attempts = 1
	}

	var lastErr error
	for attempt := 0; attempt < attempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}

		lastErr = operation(ctx)
		if lastErr == nil {
			return nil
		}
		if attempt == attempts-1 {
			return lastErr
		}

		if delay := p.Delay(attempt); delay > 0 {
			if err := p.sleep(ctx, delay); err != nil {
				return err
			}
		}
	}
	return lastErr
}

// RetryReconnect runs a provider-neutral reconnect operation with policy. It
// is equivalent to policy.Retry and is convenient for callers that keep the
// policy separate from the operation.
func RetryReconnect(ctx context.Context, policy ReconnectBackoffPolicy, operation func(context.Context) error) error {
	return policy.Retry(ctx, operation)
}

func (p ReconnectBackoffPolicy) delayBounds() (time.Duration, time.Duration) {
	maximum := p.MaxDelay
	if maximum <= 0 {
		return 0, 0
	}
	minimum := p.MinDelay
	if minimum <= 0 {
		minimum = 0
	}
	if minimum > maximum {
		minimum = maximum
	}
	return minimum, maximum
}

func (p ReconnectBackoffPolicy) sleep(ctx context.Context, delay time.Duration) error {
	if p.Sleep != nil {
		return p.Sleep(ctx, delay)
	}

	timer := time.NewTimer(delay)
	defer func() {
		if !timer.Stop() {
			select {
			case <-timer.C:
			default:
			}
		}
	}()

	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
