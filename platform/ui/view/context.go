package view

import (
	"context"
	"sync"
	"time"
)

func NewViewCtx() *ViewCtx {
	return &ViewCtx{
		userData: make(userData, 0),
	}
}

type ViewCtx struct {
	userData userData
}

func (v *ViewCtx) Deadline() (deadline time.Time, ok bool) {
	return
}

func (v *ViewCtx) Done() <-chan struct{} {
	return nil
}

func (v *ViewCtx) Err() error {
	return nil
}

func (v *ViewCtx) Value(key any) any {
	return v.userData.Get(key)
}

func (v *ViewCtx) SetUserValue(key any, value any) {
	v.userData.Set(key, value)
}

func (v *ViewCtx) SetUserValueBytes(key []byte, value any) {
	v.userData.SetBytes(key, value)
}

func (v *ViewCtx) GetUserValue(key any) any {
	return v.userData.Get(key)
}

func (v *ViewCtx) GetUserValueBytes(key []byte) any {
	return v.userData.GetBytes(key)
}

func (v *ViewCtx) Clone() *ViewCtx {
	return &ViewCtx{
		userData: v.userData.Clone(),
	}
}

// CombinedContext merges multiple contexts with proper cancellation and deadline handling
type CombinedContext struct {
	contexts []context.Context
	done     chan struct{}
	err      error
	mu       sync.RWMutex
	cancel   context.CancelFunc
}

func (c *CombinedContext) Deadline() (deadline time.Time, ok bool) {
	var earliest time.Time
	var hasDeadline bool

	for _, ctx := range c.contexts {
		if d, ok := ctx.Deadline(); ok {
			if !hasDeadline || d.Before(earliest) {
				earliest = d
				hasDeadline = true
			}
		}
	}

	return earliest, hasDeadline
}

func (c *CombinedContext) Done() <-chan struct{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.done
}

func (c *CombinedContext) Err() error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.err
}

func (c *CombinedContext) Value(key any) any {
	for _, ctx := range c.contexts {
		if v := ctx.Value(key); v != nil {
			return v
		}
	}
	return nil
}

// startWatching monitors all contexts for cancellation
func (c *CombinedContext) startWatching() {
	if len(c.contexts) == 0 {
		return
	}

	// Create a goroutine to watch for cancellation from any context
	go func() {
		defer func() {
			c.mu.Lock()
			close(c.done)
			c.mu.Unlock()
		}()

		// Use select to wait for any context to be done
		cases := make([]any, len(c.contexts))
		for i, ctx := range c.contexts {
			cases[i] = ctx.Done()
		}

		// Wait for any context to be cancelled
		for _, ctx := range c.contexts {
			if ctx.Done() != nil {
				go func(ctx context.Context) {
					<-ctx.Done()
					c.mu.Lock()
					if c.err == nil {
						c.err = ctx.Err()
					}
					c.mu.Unlock()
					if c.cancel != nil {
						c.cancel()
					}
				}(ctx)
			}
		}
	}()
}

func CombineContexts(contexts ...context.Context) context.Context {
	if len(contexts) == 0 {
		return context.Background()
	}

	if len(contexts) == 1 {
		return contexts[0]
	}

	// Filter out nil contexts
	validContexts := make([]context.Context, 0, len(contexts))
	for _, ctx := range contexts {
		if ctx != nil {
			validContexts = append(validContexts, ctx)
		}
	}

	if len(validContexts) == 0 {
		return context.Background()
	}

	if len(validContexts) == 1 {
		return validContexts[0]
	}

	combined := &CombinedContext{
		contexts: validContexts,
		done:     make(chan struct{}),
	}

	// Set up cancellation watching
	combined.startWatching()

	return combined
}

func (v *ViewCtx) WithContext(ctxs ...context.Context) context.Context {
	allContexts := make([]context.Context, 0, len(ctxs)+1)
	allContexts = append(allContexts, v)
	allContexts = append(allContexts, ctxs...)
	return CombineContexts(allContexts...)
}
