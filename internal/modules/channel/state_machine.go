package channel

import (
	"errors"
	"fmt"
	"sync"
)

// lifecycleStates is the complete set of lifecycle states supported by the
// channel state machine. Its order is stable for callers that need to render
// or validate the state vocabulary.
var lifecycleStates = [...]ConnectionState{
	StatePending,
	StateConnecting,
	StateConnected,
	StateDegraded,
	StateDisconnected,
	StateExpired,
	StateError,
}

// ErrInvalidTransition identifies an attempted lifecycle transition that is
// not allowed by the channel state machine.
var ErrInvalidTransition = errors.New("invalid channel state transition")

// InvalidTransitionError describes a rejected lifecycle transition. It wraps
// ErrInvalidTransition so callers can use either errors.Is or errors.As
// without depending on an error string.
type InvalidTransitionError struct {
	From ConnectionState
	To   ConnectionState
}

func (e *InvalidTransitionError) Error() string {
	if e == nil {
		return "<nil>"
	}
	return fmt.Sprintf("cannot transition channel from %q to %q", e.From, e.To)
}

func (e *InvalidTransitionError) Unwrap() error {
	return ErrInvalidTransition
}

// IsInvalidTransition reports whether err is a rejected channel lifecycle
// transition.
func IsInvalidTransition(err error) bool {
	return errors.Is(err, ErrInvalidTransition)
}

// lifecycleStateTransitions is the deterministic lifecycle transition table.
// A state is allowed to remain unchanged; Transition treats that as an
// idempotent success even though it is not listed in the table.
var lifecycleStateTransitions = map[ConnectionState][]ConnectionState{
	StatePending: {
		StateConnecting,
		StateExpired,
		StateError,
	},
	StateConnecting: {
		StateConnected,
		StateDegraded,
		StateDisconnected,
		StateExpired,
		StateError,
	},
	StateConnected: {
		StateDegraded,
		StateDisconnected,
		StateExpired,
		StateError,
	},
	StateDegraded: {
		StateConnecting,
		StateConnected,
		StateDisconnected,
		StateExpired,
		StateError,
	},
	StateDisconnected: {
		StatePending,
	},
	StateExpired: {
		StatePending,
	},
	StateError: {
		StatePending,
	},
}

// SupportedLifecycleStates returns the supported lifecycle states in stable
// order. It returns a copy so callers cannot alter the state vocabulary.
func SupportedLifecycleStates() []ConnectionState {
	return append([]ConnectionState(nil), lifecycleStates[:]...)
}

// ValidTransitionsFrom returns the allowed next states from state in stable
// order. It returns a copy so callers cannot alter the transition rules.
func ValidTransitionsFrom(state ConnectionState) []ConnectionState {
	return validTransitions(state)
}

// CanTransition reports whether a lifecycle state can move to target. Staying
// in the same valid state is accepted as an idempotent operation.
func CanTransition(from, target ConnectionState) bool {
	return canTransition(from, target)
}

// ValidateTransition validates a lifecycle transition without applying it.
func ValidateTransition(from, target ConnectionState) error {
	if !canTransition(from, target) {
		return &InvalidTransitionError{From: from, To: target}
	}
	return nil
}

// StateMachine validates and applies provider-neutral channel lifecycle
// transitions. It is safe for concurrent readers and transitions; each
// transition observes and updates one current state atomically.
type StateMachine struct {
	mu    sync.RWMutex
	state ConnectionState
}

// LifecycleStateMachine is a descriptive alias for StateMachine.
type LifecycleStateMachine = StateMachine

// NewStateMachine creates a state machine. With no initial state it starts in
// pending, which is the lifecycle entry state for a new channel. An explicit
// state is useful when hydrating a machine from persisted channel state.
func NewStateMachine(initial ...ConnectionState) *StateMachine {
	state := StatePending
	if len(initial) > 0 {
		state = initial[0]
	}
	return &StateMachine{state: state}
}

// NewLifecycleStateMachine is an explicit lifecycle-named constructor.
func NewLifecycleStateMachine(initial ...ConnectionState) *StateMachine {
	return NewStateMachine(initial...)
}

// State returns the machine's current lifecycle state.
func (m *StateMachine) State() ConnectionState {
	if m == nil {
		return StateUnknown
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.state
}

// CanTransition reports whether the machine can move to target. Remaining in
// the current valid state is always allowed as an idempotent operation.
func (m *StateMachine) CanTransition(target ConnectionState) bool {
	if m == nil {
		return false
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	return canTransition(m.state, target)
}

// ValidTransitions returns the allowed next states from the machine's current
// state in deterministic order. It returns a new slice on every call.
func (m *StateMachine) ValidTransitions() []ConnectionState {
	if m == nil {
		return nil
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	return validTransitions(m.state)
}

// Transition validates and applies target as the next state. Invalid
// transitions leave the current state unchanged and return a typed error.
func (m *StateMachine) Transition(target ConnectionState) error {
	if m == nil {
		return &InvalidTransitionError{From: StateUnknown, To: target}
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if err := ValidateTransition(m.state, target); err != nil {
		return err
	}
	m.state = target
	return nil
}

// TransitionError is retained as a concise name for callers that prefer it.
type TransitionError = InvalidTransitionError

func canTransition(from, to ConnectionState) bool {
	if from == to && isLifecycleState(from) {
		return true
	}
	for _, allowed := range lifecycleStateTransitions[from] {
		if allowed == to {
			return true
		}
	}
	return false
}

func validTransitions(from ConnectionState) []ConnectionState {
	return append([]ConnectionState(nil), lifecycleStateTransitions[from]...)
}

func isLifecycleState(state ConnectionState) bool {
	for _, candidate := range lifecycleStates {
		if candidate == state {
			return true
		}
	}
	return false
}
