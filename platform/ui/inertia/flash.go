package inertia

import (
	"context"
	"fmt"

	"github.com/romsar/gonertia"
)

// Flash handling for Inertia.js
type flash struct {
	errors       map[string]gonertia.ValidationErrors
	clearHistory map[string]bool
}

func newFlash() gonertia.FlashProvider {
	return &flash{
		errors:       make(map[string]gonertia.ValidationErrors),
		clearHistory: make(map[string]bool),
	}
}

func (f *flash) FlashErrors(ctx context.Context, errors gonertia.ValidationErrors) error {
	sessionID := f.getSessionIDFromContext(ctx)
	f.errors[sessionID] = errors
	return nil
}

func (f *flash) GetErrors(ctx context.Context) (gonertia.ValidationErrors, error) {
	sessionID := f.getSessionIDFromContext(ctx)
	errors, ok := f.errors[sessionID]
	if !ok {
		return nil, fmt.Errorf("history doesn't exist with that session id")
	}
	delete(f.errors, sessionID) // Clean up after retrieving
	return errors, nil
}

func (f *flash) FlashClearHistory(ctx context.Context) error {
	sessionID := f.getSessionIDFromContext(ctx)
	f.clearHistory[sessionID] = true
	return nil
}

func (f *flash) ShouldClearHistory(ctx context.Context) (bool, error) {
	sessionID := f.getSessionIDFromContext(ctx)
	clearHistory, ok := f.clearHistory[sessionID]
	if !ok {
		return false, fmt.Errorf("history doesn't exist with that session id")
	}
	delete(f.clearHistory, sessionID) // Clean up after retrieving
	return clearHistory, nil
}

func (f *flash) getSessionIDFromContext(ctx context.Context) string {
	sessionID, ok := ctx.Value("session_id").(string)
	if !ok {
		return ""
	}
	return sessionID
}
