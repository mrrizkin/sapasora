package channel

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// CircuitState is the provider-neutral availability state of a circuit
// breaker. It is intentionally separate from ConnectionState: a provider can
// be connected while calls to it are temporarily blocked by this breaker.
type CircuitState string

const (
	CircuitStateClosed   CircuitState = "closed"
	CircuitStateOpen     CircuitState = "open"
	CircuitStateHalfOpen CircuitState = "half_open"

	// Short aliases make the state vocabulary convenient without coupling it to
	// a provider or to the channel lifecycle state machine.
	CircuitClosed   = CircuitStateClosed
	CircuitOpen     = CircuitStateOpen
	CircuitHalfOpen = CircuitStateHalfOpen
	Closed          = CircuitStateClosed
	Open            = CircuitStateOpen
	HalfOpen        = CircuitStateHalfOpen
)

// CircuitBreakerConfig controls consecutive-failure opening and recovery
// probing. Clock is injectable so state transitions can be tested without
// sleeping. A nil Clock uses time.Now. Now is retained as a descriptive alias
// for callers that prefer the name used by their clock abstraction; Clock wins
// when both are supplied.
type CircuitBreakerConfig struct {
	FailureThreshold int
	OpenCooldown     time.Duration
	Clock            func() time.Time
	Now              func() time.Time
}

// Validate checks the explicit circuit-breaker configuration values.
func (c CircuitBreakerConfig) Validate() error {
	if c.FailureThreshold <= 0 {
		return fmt.Errorf("%w: failure threshold must be positive", ErrInvalidCircuitBreakerConfig)
	}
	if c.OpenCooldown < 0 {
		return fmt.Errorf("%w: open cooldown cannot be negative", ErrInvalidCircuitBreakerConfig)
	}
	return nil
}

func (c CircuitBreakerConfig) normalized() CircuitBreakerConfig {
	if c.FailureThreshold <= 0 {
		// Keep the zero value useful, while still making the configured form
		// strict through Validate and NewCircuitBreakerWithError.
		c.FailureThreshold = 1
	}
	if c.OpenCooldown < 0 {
		c.OpenCooldown = 0
	}
	if c.Clock == nil {
		c.Clock = c.Now
	}
	if c.Clock == nil {
		c.Clock = time.Now
	}
	return c
}

// CircuitBreakerSnapshot is a race-free point-in-time view of breaker state.
type CircuitBreakerSnapshot struct {
	State               CircuitState
	ConsecutiveFailures int
	OpenedAt            time.Time
	RetryAt             time.Time
	ProbeInFlight       bool
}

// CircuitBreaker prevents repeated calls to an unhealthy provider. It is safe
// for concurrent use. It is deliberately operation-shaped rather than
// provider-shaped: integration code supplies a context-aware function at the
// provider/channel boundary and remains responsible for error normalization.
type CircuitBreaker struct {
	mu sync.Mutex

	config              CircuitBreakerConfig
	state               CircuitState
	consecutiveFailures int
	openedAt            time.Time
	retryAt             time.Time
	probeInFlight       bool
}

// NewCircuitBreaker creates a provider-neutral breaker. Non-positive
// thresholds are normalized to one so a zero-value config remains useful;
// callers that need strict configuration validation can use
// NewCircuitBreakerWithError.
func NewCircuitBreaker(config CircuitBreakerConfig) *CircuitBreaker {
	config = config.normalized()
	return &CircuitBreaker{
		config: config,
		state:  CircuitStateClosed,
	}
}

// NewCircuitBreakerWithError validates config before constructing a breaker.
func NewCircuitBreakerWithError(config CircuitBreakerConfig) (*CircuitBreaker, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}
	return NewCircuitBreaker(config), nil
}

// ProviderCircuitBreaker is a descriptive alias for CircuitBreaker.
type ProviderCircuitBreaker = CircuitBreaker

// NewProviderCircuitBreaker is an explicit provider-scoped constructor alias.
func NewProviderCircuitBreaker(config CircuitBreakerConfig) *CircuitBreaker {
	return NewCircuitBreaker(config)
}

// State returns the current breaker state. An open breaker moves to half-open
// when the next Execute call observes that its cooldown has elapsed; State
// itself does not mutate state or invoke the configured clock.
func (b *CircuitBreaker) State() CircuitState {
	if b == nil {
		return CircuitStateOpen
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state
}

// ConsecutiveFailures returns the current consecutive failure count.
func (b *CircuitBreaker) ConsecutiveFailures() int {
	if b == nil {
		return 0
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.consecutiveFailures
}

// Snapshot returns a copy of the current breaker state for diagnostics and
// tests. It does not transition an open circuit to half-open.
func (b *CircuitBreaker) Snapshot() CircuitBreakerSnapshot {
	if b == nil {
		return CircuitBreakerSnapshot{State: CircuitStateOpen}
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	return CircuitBreakerSnapshot{
		State:               b.state,
		ConsecutiveFailures: b.consecutiveFailures,
		OpenedAt:            b.openedAt,
		RetryAt:             b.retryAt,
		ProbeInFlight:       b.probeInFlight,
	}
}

// Execute runs operation when the circuit permits it. Calls are rejected
// while open until OpenCooldown elapses. Exactly one call is admitted as the
// half-open probe; concurrent callers receive a typed error. Context
// cancellation is checked before admission and the same context is passed to
// operation. A nil operation is rejected without changing breaker state.
func (b *CircuitBreaker) Execute(ctx context.Context, operation func(context.Context) error) error {
	if b == nil {
		return ErrNilCircuitBreaker
	}
	if operation == nil {
		return ErrNilCircuitOperation
	}
	if ctx == nil {
		return ErrNilCircuitContext
	}
	if err := ctx.Err(); err != nil {
		return err
	}

	probe, err := b.admit()
	if err != nil {
		return err
	}

	err = operation(ctx)
	b.complete(probe, err)
	return err
}

// Call is a concise alias for Execute.
func (b *CircuitBreaker) Call(ctx context.Context, operation func(context.Context) error) error {
	return b.Execute(ctx, operation)
}

func (b *CircuitBreaker) admit() (bool, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	switch b.state {
	case CircuitStateClosed:
		return false, nil
	case CircuitStateOpen:
		now := b.config.Clock()
		if now.Before(b.retryAt) {
			return false, &CircuitOpenError{OpenedAt: b.openedAt, RetryAt: b.retryAt}
		}
		b.state = CircuitStateHalfOpen
		b.probeInFlight = true
		return true, nil
	case CircuitStateHalfOpen:
		return false, &CircuitHalfOpenError{RetryAt: b.retryAt}
	default:
		// The state is private and only takes known values, but fail closed if
		// that invariant is ever broken during a future change.
		return false, &CircuitOpenError{OpenedAt: b.openedAt, RetryAt: b.retryAt}
	}
}

func (b *CircuitBreaker) complete(probe bool, operationErr error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if operationErr == nil {
		b.state = CircuitStateClosed
		b.consecutiveFailures = 0
		b.openedAt = time.Time{}
		b.retryAt = time.Time{}
		b.probeInFlight = false
		return
	}

	if probe {
		b.open()
		return
	}

	b.consecutiveFailures++
	if b.consecutiveFailures >= b.config.FailureThreshold {
		b.open()
	}
}

func (b *CircuitBreaker) open() {
	now := b.config.Clock()
	b.state = CircuitStateOpen
	b.openedAt = now
	b.retryAt = now.Add(b.config.OpenCooldown)
	b.probeInFlight = false
	if b.consecutiveFailures < b.config.FailureThreshold {
		b.consecutiveFailures = b.config.FailureThreshold
	}
}

var (
	// ErrCircuitOpen indicates that the cooldown has not elapsed.
	ErrCircuitOpen = errors.New("channel circuit breaker is open")
	// ErrCircuitHalfOpen indicates that another half-open probe is in flight.
	ErrCircuitHalfOpen = errors.New("channel circuit breaker half-open probe in progress")
	// ErrInvalidCircuitBreakerConfig indicates invalid explicit configuration.
	ErrInvalidCircuitBreakerConfig = errors.New("invalid channel circuit breaker configuration")
	// ErrNilCircuitOperation indicates that Execute was given no operation.
	ErrNilCircuitOperation = errors.New("channel circuit breaker operation is nil")
	// ErrNilCircuitContext indicates that Execute was given a nil context.
	ErrNilCircuitContext = errors.New("channel circuit breaker context is nil")
	// ErrNilCircuitBreaker indicates a call through a nil breaker pointer.
	ErrNilCircuitBreaker = errors.New("channel circuit breaker is nil")
)

// CircuitOpenError is returned when the open cooldown has not elapsed.
type CircuitOpenError struct {
	OpenedAt time.Time
	RetryAt  time.Time
}

func (e *CircuitOpenError) Error() string {
	if e == nil {
		return ErrCircuitOpen.Error()
	}
	return fmt.Sprintf("%s until %s", ErrCircuitOpen, e.RetryAt.Format(time.RFC3339Nano))
}

func (e *CircuitOpenError) Unwrap() error { return ErrCircuitOpen }

// CircuitHalfOpenError is returned to callers that arrive while the single
// half-open probe is already in flight.
type CircuitHalfOpenError struct {
	RetryAt time.Time
}

func (e *CircuitHalfOpenError) Error() string {
	if e == nil {
		return ErrCircuitHalfOpen.Error()
	}
	return fmt.Sprintf("%s; retry after %s", ErrCircuitHalfOpen, e.RetryAt.Format(time.RFC3339Nano))
}

func (e *CircuitHalfOpenError) Unwrap() error { return ErrCircuitHalfOpen }

// CircuitStateError reports whether a rejected call met an open or half-open
// circuit. It is useful when callers want one typed extraction point.
type CircuitStateError interface {
	error
	State() CircuitState
}

func (e *CircuitOpenError) State() CircuitState     { return CircuitStateOpen }
func (e *CircuitHalfOpenError) State() CircuitState { return CircuitStateHalfOpen }

// IsCircuitOpen reports whether err is a cooldown rejection.
func IsCircuitOpen(err error) bool { return errors.Is(err, ErrCircuitOpen) }

// IsCircuitHalfOpen reports whether err was rejected because a probe is in
// flight.
func IsCircuitHalfOpen(err error) bool { return errors.Is(err, ErrCircuitHalfOpen) }

// IsCircuitBreakerError reports whether err is a typed circuit rejection.
func IsCircuitBreakerError(err error) bool {
	return IsCircuitOpen(err) || IsCircuitHalfOpen(err)
}
