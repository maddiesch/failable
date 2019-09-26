package failable_test

import (
	"context"
	"errors"
	"testing"
	"time"

	. "github.com/maddiesch/failable"
	"github.com/stretchr/testify/assert"
)

func assertCompleted(t *testing.T, done Completed, failed Failed) {
	select {
	case err := <-failed:
		t.Error(err)
	case <-done:
		// Okay
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Test Timeout")
	}
}

func assertFailed(t *testing.T, done Completed, failed Failed) error {
	select {
	case err := <-failed:
		return err
	case <-done:
		t.Error("Run completed without an error")
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Test Timeout")
	}
	return nil
}

func TestFailable(t *testing.T) {
	t.Run("RunWithContext", func(t *testing.T) {
		t.Run("completion", func(t *testing.T) {
			t.Parallel()

			done, failed := RunWithContext(context.Background(), func(ctx context.Context, fail FailFunc) {
				time.Sleep(5 * time.Millisecond)
			})

			assertCompleted(t, done, failed)
		})

		t.Run("fail with nil error", func(t *testing.T) {
			t.Parallel()

			done, failed := RunWithContext(context.Background(), func(ctx context.Context, fail FailFunc) {
				time.Sleep(5 * time.Millisecond)
				fail(nil)
				// Fail should exit the goroutine this sleep will ensure that
				// the timeout fails.
				time.Sleep(5 * time.Second)
			})

			err := assertFailed(t, done, failed)

			assert.Equal(t, ErrNilFailure, err)
		})

		t.Run("fail with nil error", func(t *testing.T) {
			t.Parallel()

			failure := errors.New("testing error failure")

			done, failed := RunWithContext(context.Background(), func(ctx context.Context, fail FailFunc) {
				time.Sleep(5 * time.Millisecond)
				fail(failure)
				// Fail should exit the goroutine this sleep will ensure that
				// the timeout fails.
				time.Sleep(5 * time.Second)
			})

			err := assertFailed(t, done, failed)

			assert.Equal(t, failure, err)
		})
	})
}
