// Package failable provides a function that does work with failable chanels already setup.
package failable

import (
	"context"
	"errors"
	"runtime"
)

var (
	// ErrNilFailure sent if the error passed to the fail function is nil
	ErrNilFailure = errors.New("nil failure")
)

// FailFunc can be called to fail an operation immediately.
//
// If the failure function is passed a nil error the failed channel will be signaled with ErrCompleted
type FailFunc func(error)

// HandlerFunc is the main function that is called with a context and the failure function.
type HandlerFunc func(context.Context, FailFunc)

// SimpleHandlerFunc allows you to directly return an error instead of calling the failure function.
type SimpleHandlerFunc func(context.Context) error

// Completed is the channel this is signaled when the function successfully completes.
type Completed <-chan struct{}

// Failed is the channel this is signaled with the failing error.
type Failed <-chan error

// Do performs a run and handles the channel's appropriately.
func Do(fn func(FailFunc)) error {
	return DoWithContext(context.Background(), func(ctx context.Context, fail FailFunc) {
		fn(fail)
	})
}

// DoWithContext performs a run and handles the channel's appropriately
func DoWithContext(ctx context.Context, fn HandlerFunc) error {
	done, fail := RunWithContext(ctx, fn)

	select {
	case <-done:
		return nil
	case err := <-fail:
		return err
	}
}

// Run performs the work in a new goroutine
//
// If the fail function is called the goroutine will exit and the failed channel will be signaled with the error.
func Run(fn func(FailFunc)) (Completed, Failed) {
	return RunWithContext(context.Background(), func(ctx context.Context, fail FailFunc) {
		fn(fail)
	})
}

// RunWithContext performs the work in a new goroutine passing the context into the handler function.
func RunWithContext(ctx context.Context, fn HandlerFunc) (Completed, Failed) {
	completed := make(chan struct{})
	failed := make(chan error)

	fail := func(err error) {
		if err == nil {
			failed <- ErrNilFailure
		} else {
			failed <- err
		}
		runtime.Goexit()
	}

	go func() {
		fn(ctx, fail)

		completed <- struct{}{}
	}()

	return completed, failed
}

// RunSimple performs the function in a new goroutine and will fail with returned error if there is one.
func RunSimple(fn func() error) (Completed, Failed) {
	return RunSimpleWithContext(context.Background(), func(ctx context.Context) error {
		return fn()
	})
}

// RunSimpleWithContext performs the function in a new goroutine and will fail with returned error if there is one.
//
// The context is passed to the function.
func RunSimpleWithContext(ctx context.Context, fn SimpleHandlerFunc) (Completed, Failed) {
	return RunWithContext(ctx, func(ctx context.Context, fail FailFunc) {
		err := fn(ctx)

		if err != nil {
			fail(err)
		}
	})
}
