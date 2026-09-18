package message

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"sync"
	"testing"
)

func TestMessageStateMachineUsesAcceptedAsDefaultAndStrictTransitions(t *testing.T) {
	machine := NewMessageStateMachine()
	if got := machine.State(); got != MessageStatusAccepted {
		t.Fatalf("State() = %q, want %q", got, MessageStatusAccepted)
	}
	want := []MessageStatus{MessageStatusQueued, MessageStatusCanceled}
	if got := machine.ValidTransitions(); !reflect.DeepEqual(got, want) {
		t.Fatalf("ValidTransitions() = %v, want %v", got, want)
	}
	if err := machine.Transition(MessageStatusSent); err == nil || !errors.Is(err, ErrInvalidTransition) {
		t.Fatalf("Transition(sent) error = %v, want ErrInvalidTransition", err)
	}
	if got := machine.State(); got != MessageStatusAccepted {
		t.Fatalf("state after rejected transition = %q", got)
	}
	if err := machine.Transition(MessageStatusQueued); err != nil {
		t.Fatal(err)
	}
	if err := machine.Transition(MessageStatusQueued); err == nil {
		t.Fatal("self-transition was accepted")
	}
}

func TestMessageStateMachineSupportsDocumentedLifecycle(t *testing.T) {
	transitions := map[MessageStatus][]MessageStatus{
		MessageStatusAccepted:  {MessageStatusQueued, MessageStatusCanceled},
		MessageStatusQueued:    {MessageStatusSending, MessageStatusCanceled},
		MessageStatusSending:   {MessageStatusSent, MessageStatusFailed},
		MessageStatusSent:      {MessageStatusDelivered},
		MessageStatusDelivered: {MessageStatusRead},
		MessageStatusFailed:    {MessageStatusRetrying},
		MessageStatusRetrying:  {MessageStatusSent, MessageStatusFailed, MessageStatusDeadLetter},
	}
	for from, targets := range transitions {
		for _, target := range targets {
			if !CanTransition(from, target) {
				t.Errorf("CanTransition(%q, %q) = false", from, target)
			}
		}
	}
	for _, terminal := range []MessageStatus{MessageStatusRead, MessageStatusDeadLetter, MessageStatusCanceled} {
		if got := ValidTransitionsFrom(terminal); len(got) != 0 {
			t.Errorf("ValidTransitionsFrom(%q) = %v, want no transitions", terminal, got)
		}
	}
}

func TestMessageRepositoryTransitionEventsAreAtomicAndProviderAcceptanceIsDistinct(t *testing.T) {
	repository := NewInMemoryRepository()
	value := testMessage("tenant-a", "message-state")
	if err := repository.CreateMessage(context.Background(), &value); err != nil {
		t.Fatal(err)
	}

	platformEvent, err := repository.TransitionMessage(context.Background(), "tenant-a", value.PublicID, MessageTransition{Target: MessageStatusQueued, RequestID: "request-safe"})
	if err != nil {
		t.Fatal(err)
	}
	if platformEvent.Source != MessageEventSourcePlatform || platformEvent.Status != MessageStatusQueued {
		t.Fatalf("platform transition event = %+v", platformEvent)
	}
	if _, err := repository.TransitionMessage(context.Background(), "tenant-a", value.PublicID, MessageTransition{Target: MessageStatusSent}); !errors.Is(err, ErrInvalidTransition) {
		t.Fatalf("invalid transition error = %v", err)
	}
	updated, err := repository.GetMessageByPublicID(context.Background(), "tenant-a", value.PublicID)
	if err != nil {
		t.Fatal(err)
	}
	updated.Status = MessageStatusSent
	if err := repository.UpdateMessage(context.Background(), updated); err == nil {
		t.Fatal("UpdateMessage allowed a lifecycle status mutation without an event")
	}
	updated, err = repository.GetMessageByPublicID(context.Background(), "tenant-a", value.PublicID)
	if err != nil || updated.Status != MessageStatusQueued {
		t.Fatalf("status after rejected UpdateMessage = %q, %v", updated.Status, err)
	}

	providerAccepted, err := repository.RecordProviderStatusEvent(context.Background(), "tenant-a", value.PublicID, ProviderStatusEvent{Status: MessageStatusAccepted, EventID: "provider-accepted-1", ProviderMessageID: "opaque-provider-id"})
	if err != nil {
		t.Fatal(err)
	}
	if providerAccepted.Source != MessageEventSourceProvider || providerAccepted.Status != MessageStatusAccepted {
		t.Fatalf("provider acceptance event = %+v", providerAccepted)
	}
	encoded, err := json.Marshal(providerAccepted)
	if err != nil {
		t.Fatal(err)
	}
	for _, opaque := range []string{"provider-accepted-1", "opaque-provider-id"} {
		if strings.Contains(string(encoded), opaque) {
			t.Fatalf("provider diagnostic leaked %q: %s", opaque, encoded)
		}
	}
	stored, err := repository.GetMessageByPublicID(context.Background(), "tenant-a", value.PublicID)
	if err != nil {
		t.Fatal(err)
	}
	if stored.Status != MessageStatusQueued {
		t.Fatalf("provider acceptance changed platform state to %q", stored.Status)
	}

	duplicate, err := repository.RecordProviderStatusEvent(context.Background(), "tenant-a", value.PublicID, ProviderStatusEvent{Status: MessageStatusAccepted, EventID: "provider-accepted-1"})
	if err != nil || duplicate.ID != providerAccepted.ID {
		t.Fatalf("duplicate provider event = %+v, %v", duplicate, err)
	}
	events, err := repository.ListMessageEvents(context.Background(), "tenant-a", value.PublicID)
	if err != nil || len(events) != 2 {
		t.Fatalf("events after deduplication = %d, %v", len(events), err)
	}
}

func TestMessageRepositoryDeduplicatesProviderStatusEventsUnderRace(t *testing.T) {
	repository := NewInMemoryRepository()
	value := testMessage("tenant-race", "message-race")
	if err := repository.CreateMessage(context.Background(), &value); err != nil {
		t.Fatal(err)
	}
	const workers = 32
	var wait sync.WaitGroup
	errorsCh := make(chan error, workers)
	ids := make(chan uint64, workers)
	for i := 0; i < workers; i++ {
		wait.Add(1)
		go func() {
			defer wait.Done()
			event, err := repository.RecordProviderStatusEvent(context.Background(), "tenant-race", value.PublicID, ProviderStatusEvent{Status: MessageStatusAccepted, ProviderEventID: "provider-event-race"})
			if err != nil {
				errorsCh <- err
				return
			}
			ids <- event.ID
		}()
	}
	wait.Wait()
	close(errorsCh)
	close(ids)
	for err := range errorsCh {
		t.Fatal(err)
	}
	var first uint64
	for id := range ids {
		if first == 0 {
			first = id
		}
		if id != first {
			t.Fatalf("deduplicated provider event IDs differ: %d and %d", first, id)
		}
	}
	events, err := repository.ListMessageEvents(context.Background(), "tenant-race", value.PublicID)
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 1 {
		t.Fatalf("provider event count under race = %d, want 1", len(events))
	}
}
