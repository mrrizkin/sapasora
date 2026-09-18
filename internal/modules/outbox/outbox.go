package outbox

import (
	"context"
	"errors"
	"time"
)

// Service exposes the safe write/publication boundary to message modules.
// Callers append an event first and publish it explicitly; send jobs may only
// be enqueued after publication has succeeded.
type Service struct {
	store EventStore
}

func NewService(store EventStore) (*Service, error) {
	if store == nil {
		return nil, errors.New("outbox: store is required")
	}
	return &Service{store: store}, nil
}

// AppendAndPublish leaves a pending event behind if publication fails. A relay
// can safely call Publish again because publication is idempotent.
func (s *Service) AppendAndPublish(ctx context.Context, event OutboxEvent) error {
	if err := s.store.AppendEvent(ctx, event); err != nil {
		return err
	}
	scope := eventScope(event)
	return s.store.PublishEvent(ctx, scope, event.ID)
}

func (s *Service) AppendEvent(ctx context.Context, event OutboxEvent) error {
	return s.store.AppendEvent(ctx, event)
}

func (s *Service) PublishEvent(ctx context.Context, scope Scope, eventID string) error {
	return s.store.PublishEvent(ctx, scope, eventID)
}

func (s *Service) EnqueueSendJob(ctx context.Context, job SendJob) error {
	return s.store.EnqueueJob(ctx, job)
}

func (s *Service) CancelJob(ctx context.Context, scope Scope, jobID string, now time.Time) error {
	return s.store.CancelJob(ctx, scope, jobID, now)
}

func (s *Service) ReplayDeadLetter(ctx context.Context, scope Scope, jobID string, now time.Time) error {
	return s.store.ReplayDeadLetter(ctx, scope, jobID, now)
}
