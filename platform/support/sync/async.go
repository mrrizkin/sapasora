// Package sync provides a facade for async operations.
package sync

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"sync/atomic"
	"time"

	"sapasora/platform/support/console"
)

// Future represents an asynchronous computation that may or may not have finished.
type Future[T any] struct {
	await  func() (T, error)
	cancel context.CancelFunc
	done   <-chan struct{}
	state  *int32 // atomic: 0=running, 1=completed, 2=cancelled
}

// Await waits for the async operation to complete and returns the result.
func (f *Future[T]) Await() (T, error) {
	return f.await()
}

// Cancel cancels the async operation if it's still running.
func (f *Future[T]) Cancel() {
	if f.cancel != nil {
		f.cancel()
		if f.state != nil {
			atomic.CompareAndSwapInt32(f.state, 0, 2) // running -> cancelled
		}
	}
}

// IsDone returns true if the async operation has completed (either successfully, with error, or cancelled).
func (f *Future[T]) IsDone() bool {
	select {
	case <-f.done:
		return true
	default:
		return false
	}
}

// IsCancelled returns true if the async operation was cancelled.
func (f *Future[T]) IsCancelled() bool {
	if f.state != nil {
		return atomic.LoadInt32(f.state) == 2
	}
	return false
}

// WithTimeout returns a new Future that will timeout after the specified duration.
func (f *Future[T]) WithTimeout(timeout time.Duration) *Future[T] {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	newFuture := &Future[T]{
		done:  f.done,
		state: f.state,
	}

	newFuture.await = func() (T, error) {
		select {
		case <-f.done:
			return f.Await()
		case <-ctx.Done():
			var zero T
			return zero, ctx.Err()
		}
	}

	newFuture.cancel = func() {
		cancel()
		f.Cancel()
	}

	return newFuture
}

// Async starts an asynchronous function and returns a Future.
//
//	future := Async(func() (string, error) {
//	    // Simulate a long-running operation
//	    time.Sleep(time.Second * 5)
//	    return "Hello, world!", nil
//	})
//
//	result, err := future.Await()
//	if err != nil {
//	    fmt.Println("Error:", err)
//	} else {
//	    fmt.Println("Result:", result)
//	}
func Async[T any](fn func() (T, error)) *Future[T] {
	return AsyncWithContext(context.Background(), fn)
}

// AsyncWithContext starts an asynchronous function with a context and returns a Future.
//
//	future := AsyncWithContext(context.Background(), func() (string, error) {
//	    // Simulate a long-running operation
//	    time.Sleep(time.Second * 5)
//	    return "Hello, world!", nil
//	})
//
//	result, err := future.Await()
//	if err != nil {
//	    fmt.Println("Error:", err)
//	} else {
//	    fmt.Println("Result:", result)
//	}
func AsyncWithContext[T any](ctx context.Context, fn func() (T, error)) *Future[T] {
	var result T
	var err error

	done := make(chan struct{})
	ctx, cancel := context.WithCancel(ctx)
	state := int32(0) // 0=running, 1=completed, 2=cancelled

	go func() {
		defer close(done)
		defer func() {
			if p := recover(); p != nil {
				err = convertPanicToError(p)
				// Print stack trace for debugging
				buf := make([]byte, 64<<10)
				n := runtime.Stack(buf, false)
				console.Error("Panic in Async operation: %v\n%s", p, buf[:n])
			}
			atomic.StoreInt32(&state, 1) // completed
		}()

		select {
		case <-ctx.Done():
			// Context was cancelled or timed out
			err = ctx.Err()
			atomic.CompareAndSwapInt32(&state, 0, 2) // running -> cancelled
		default:
			result, err = fn()
		}
	}()

	return &Future[T]{
		await: func() (T, error) {
			select {
			case <-done:
				return result, err
			case <-ctx.Done():
				var zero T
				return zero, ctx.Err()
			}
		},
		cancel: cancel,
		done:   done,
		state:  &state,
	}
}

// AsyncWithTimeout starts an asynchronous function with a timeout and returns a Future.
//
//	future := AsyncWithTimeout(func() (string, error) {
//	    // Simulate a long-running operation
//	    time.Sleep(time.Second * 5)
//	    return "Hello, world!", nil
//	}, time.Second * 3)
//
//	result, err := future.Await()
//	if err != nil {
//	    fmt.Println("Error:", err)
//	} else {
//	    fmt.Println("Result:", result)
//	}
func AsyncWithTimeout[T any](fn func() (T, error), timeout time.Duration) *Future[T] {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	future := AsyncWithContext(ctx, fn)

	// Override the cancel function to use the timeout context's cancel
	originalCancel := future.cancel
	future.cancel = func() {
		cancel()
		if originalCancel != nil {
			originalCancel()
		}
	}

	return future
}

func convertPanicToError(p any) error {
	var err error
	switch v := p.(type) {
	case error:
		err = v
	case string:
		err = errors.New(v)
	default:
		err = fmt.Errorf("%v", v)
	}
	return fmt.Errorf("panic in async operation: %w", err)
}
