package message

import (
	"errors"
	"fmt"
	"sync"
)

// messageLifecycleStates is the complete, provider-neutral message state
// vocabulary. Its order is stable for callers that render state choices.
var messageLifecycleStates = [...]MessageStatus{
	MessageStatusAccepted,
	MessageStatusQueued,
	MessageStatusSending,
	MessageStatusSent,
	MessageStatusDelivered,
	MessageStatusRead,
	MessageStatusFailed,
	MessageStatusRetrying,
	MessageStatusDeadLetter,
	MessageStatusCanceled,
}

// ErrInvalidTransition identifies a rejected message state transition.
var ErrInvalidTransition = errors.New("invalid message state transition")

// InvalidTransitionError describes a rejected message state transition. It
// intentionally contains only state vocabulary and no message or provider
// identifiers, making it safe to return from diagnostics and logs.
type InvalidTransitionError struct {
	From MessageStatus
	To   MessageStatus
}

func (e *InvalidTransitionError) Error() string {
	if e == nil {
		return "<nil>"
	}
	return fmt.Sprintf("cannot transition message from %q to %q", e.From, e.To)
}

func (e *InvalidTransitionError) Unwrap() error { return ErrInvalidTransition }

// IsInvalidTransition reports whether err is a rejected message transition.
func IsInvalidTransition(err error) bool { return errors.Is(err, ErrInvalidTransition) }

// messageStateTransitions is deliberately explicit. In particular, terminal
// states do not silently reopen and cancellation is available only before a
// send starts.
var messageStateTransitions = map[MessageStatus][]MessageStatus{
	MessageStatusAccepted:  {MessageStatusQueued, MessageStatusCanceled},
	MessageStatusQueued:    {MessageStatusSending, MessageStatusCanceled},
	MessageStatusSending:   {MessageStatusSent, MessageStatusFailed},
	MessageStatusSent:      {MessageStatusDelivered},
	MessageStatusDelivered: {MessageStatusRead},
	MessageStatusFailed:    {MessageStatusRetrying},
	MessageStatusRetrying:  {MessageStatusSent, MessageStatusFailed, MessageStatusDeadLetter},
}

// SupportedMessageStatuses returns all message states in stable order.
func SupportedMessageStatuses() []MessageStatus {
	return append([]MessageStatus(nil), messageLifecycleStates[:]...)
}

// ValidTransitionsFrom returns the allowed targets from state in stable order.
func ValidTransitionsFrom(state MessageStatus) []MessageStatus {
	return append([]MessageStatus(nil), messageStateTransitions[state]...)
}

// CanTransition reports whether a message can move from one state to another.
// Self-transitions are not valid: duplicate provider observations are handled
// by provider event deduplication rather than by weakening transition rules.
func CanTransition(from, target MessageStatus) bool {
	for _, allowed := range messageStateTransitions[from] {
		if allowed == target {
			return true
		}
	}
	return false
}

// ValidateTransition validates a message transition without applying it.
func ValidateTransition(from, target MessageStatus) error {
	if !CanTransition(from, target) {
		return &InvalidTransitionError{From: from, To: target}
	}
	return nil
}

// StateMachine validates and applies message transitions. It is safe for
// concurrent readers and transitions.
type StateMachine struct {
	mu    sync.RWMutex
	state MessageStatus
}

// MessageStateMachine is the descriptive name for StateMachine.
type MessageStateMachine = StateMachine

// NewStateMachine starts a new message in platform-accepted state. An
// explicit state is useful when hydrating a machine from persisted state.
func NewStateMachine(initial ...MessageStatus) *StateMachine {
	state := MessageStatusAccepted
	if len(initial) > 0 {
		state = initial[0]
	}
	return &StateMachine{state: state}
}

// NewMessageStateMachine is an explicit message-named constructor.
func NewMessageStateMachine(initial ...MessageStatus) *StateMachine {
	return NewStateMachine(initial...)
}

// SupportedStates is a concise compatibility name for message callers.
func SupportedStates() []MessageStatus { return SupportedMessageStatuses() }

// ValidMessageTransitionsFrom is the explicit message-named transition query.
func ValidMessageTransitionsFrom(state MessageStatus) []MessageStatus {
	return ValidTransitionsFrom(state)
}

// CanTransitionMessage reports whether a message state transition is valid.
func CanTransitionMessage(from, target MessageStatus) bool {
	return CanTransition(from, target)
}

// ValidateMessageTransition validates a message state transition.
func ValidateMessageTransition(from, target MessageStatus) error {
	return ValidateTransition(from, target)
}

// State returns the current state, or unknown for a nil machine.
func (m *StateMachine) State() MessageStatus {
	if m == nil {
		return MessageStatusUnknown
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.state
}

// CanTransition reports whether the machine can move to target.
func (m *StateMachine) CanTransition(target MessageStatus) bool {
	if m == nil {
		return false
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	return CanTransition(m.state, target)
}

// ValidTransitions returns a defensive copy of allowed targets.
func (m *StateMachine) ValidTransitions() []MessageStatus {
	if m == nil {
		return nil
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	return ValidTransitionsFrom(m.state)
}

// Transition validates and applies target. Rejected transitions do not change
// the current state.
func (m *StateMachine) Transition(target MessageStatus) error {
	if m == nil {
		return &InvalidTransitionError{From: MessageStatusUnknown, To: target}
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if err := ValidateTransition(m.state, target); err != nil {
		return err
	}
	m.state = target
	return nil
}

// TransitionError is retained as a concise compatibility name.
type TransitionError = InvalidTransitionError
