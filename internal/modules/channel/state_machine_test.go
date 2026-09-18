package channel

import (
	"errors"
	"reflect"
	"testing"
)

func TestStateMachineUsesPendingAsDefault(t *testing.T) {
	machine := NewStateMachine()
	if got := machine.State(); got != StatePending {
		t.Fatalf("State() = %q, want %q", got, StatePending)
	}
	if got := machine.ValidTransitions(); !reflect.DeepEqual(got, []ConnectionState{StateConnecting, StateExpired, StateError}) {
		t.Fatalf("ValidTransitions() = %v", got)
	}
}

func TestStateMachineAcceptsEveryDocumentedTransition(t *testing.T) {
	tests := []struct {
		name string
		from ConnectionState
		to   []ConnectionState
	}{
		{name: "pending", from: StatePending, to: []ConnectionState{StateConnecting, StateExpired, StateError}},
		{name: "connecting", from: StateConnecting, to: []ConnectionState{StateConnected, StateDegraded, StateDisconnected, StateExpired, StateError}},
		{name: "connected", from: StateConnected, to: []ConnectionState{StateDegraded, StateDisconnected, StateExpired, StateError}},
		{name: "degraded", from: StateDegraded, to: []ConnectionState{StateConnecting, StateConnected, StateDisconnected, StateExpired, StateError}},
		{name: "disconnected", from: StateDisconnected, to: []ConnectionState{StatePending}},
		{name: "expired", from: StateExpired, to: []ConnectionState{StatePending}},
		{name: "error", from: StateError, to: []ConnectionState{StatePending}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			machine := NewStateMachine(tt.from)
			if got := machine.ValidTransitions(); !reflect.DeepEqual(got, tt.to) {
				t.Fatalf("ValidTransitions() = %v, want %v", got, tt.to)
			}
			for _, target := range tt.to {
				if err := machine.Transition(target); err != nil {
					t.Fatalf("Transition(%q) error = %v", target, err)
				}
				if got := machine.State(); got != target {
					t.Fatalf("State() = %q, want %q", got, target)
				}
				machine = NewStateMachine(tt.from)
			}
		})
	}
}

func TestStateMachineRejectsInvalidTransitionWithTypedError(t *testing.T) {
	machine := NewStateMachine(StatePending)
	err := machine.Transition(StateConnected)
	if err == nil {
		t.Fatal("Transition() returned nil for an invalid transition")
	}
	if !errors.Is(err, ErrInvalidTransition) || !IsInvalidTransition(err) {
		t.Fatalf("Transition() error = %v, want ErrInvalidTransition", err)
	}
	var transitionErr *InvalidTransitionError
	if !errors.As(err, &transitionErr) {
		t.Fatalf("Transition() error type = %T, want *InvalidTransitionError", err)
	}
	if transitionErr.From != StatePending || transitionErr.To != StateConnected {
		t.Fatalf("InvalidTransitionError = %+v", transitionErr)
	}
	if got := machine.State(); got != StatePending {
		t.Fatalf("State() = %q after rejection, want %q", got, StatePending)
	}
}

func TestStateMachineRejectsEveryUndocumentedTransition(t *testing.T) {
	states := SupportedLifecycleStates()
	for _, from := range states {
		for _, to := range states {
			if CanTransition(from, to) {
				continue
			}

			machine := NewStateMachine(from)
			err := machine.Transition(to)
			if err == nil {
				t.Errorf("Transition(%q -> %q) returned nil", from, to)
				continue
			}
			var transitionErr *InvalidTransitionError
			if !errors.As(err, &transitionErr) || transitionErr.From != from || transitionErr.To != to {
				t.Errorf("Transition(%q -> %q) error = %T %+v, want matching InvalidTransitionError", from, to, err, transitionErr)
			}
			if got := machine.State(); got != from {
				t.Errorf("State() = %q after rejecting %q -> %q, want %q", got, from, to, from)
			}
		}
	}
}

func TestStateMachineReturnsDefensiveTransitionCopies(t *testing.T) {
	states := SupportedLifecycleStates()
	states[0] = StateUnknown
	if got := SupportedLifecycleStates()[0]; got != StatePending {
		t.Fatalf("SupportedLifecycleStates() exposed mutable storage: first state = %q", got)
	}

	transitions := ValidTransitionsFrom(StatePending)
	transitions[0] = StateUnknown
	if got := ValidTransitionsFrom(StatePending)[0]; got != StateConnecting {
		t.Fatalf("ValidTransitionsFrom() exposed mutable storage: first transition = %q", got)
	}
}

func TestStateMachineAllowsIdempotentTransitions(t *testing.T) {
	for _, state := range SupportedLifecycleStates() {
		machine := NewStateMachine(state)
		if err := machine.Transition(state); err != nil {
			t.Errorf("Transition(%q) error = %v", state, err)
		}
		if got := machine.State(); got != state {
			t.Errorf("State() = %q, want %q", got, state)
		}
	}
}

func TestStateMachineRejectsUnknownState(t *testing.T) {
	machine := NewStateMachine(StateUnknown)
	if machine.CanTransition(StatePending) {
		t.Fatal("CanTransition() accepted a transition from unknown")
	}
	if err := machine.Transition(StatePending); err == nil {
		t.Fatal("Transition() accepted a transition from unknown")
	}

	machine = NewStateMachine(StatePending)
	if err := machine.Transition(StateUnknown); err == nil {
		t.Fatal("Transition() accepted unknown target state")
	}
}
