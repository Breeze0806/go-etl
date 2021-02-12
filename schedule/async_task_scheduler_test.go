package schedule

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestAsyncTaskScheduler_PostErrorLast(t *testing.T) {
	ctx := context.Background()
	scheduler := NewAsyncTaskScheduler(ctx, 1, 0)
	defer scheduler.Close()
	n := 1000
	for i := 0; i < n; i++ {
		var postErr error
		if i == n-1 {
			postErr = errors.New("mock error")
		}

		if err := scheduler.Push(newMockAsyncTask(
			time.Duration(0), []error{nil, postErr})); err != nil {
			t.Fatalf("Push error= %v", err)
		}
	}

	select {
	case <-scheduler.Errors():
	case <-ctx.Done():
	}

	if got := scheduler.Size(); got != 0 {
		t.Fatalf("Size() = %v  want: 0 ", got)
	}
}

func TestAsyncTaskScheduler_PostErrorMid(t *testing.T) {
	ctx := context.Background()
	scheduler := NewAsyncTaskScheduler(ctx, 1, 0)
	defer scheduler.Close()
	n := 1000
	var err error
	for i := 0; i < n*10; i++ {
		var postErr error
		if i == n-1 {
			postErr = errors.New("mock error")
		}

		if err = scheduler.Push(newMockAsyncTask(
			time.Duration(0), []error{nil, postErr})); err != nil {
			break
		}
	}
	if err == nil {
		t.Fatal("Push Last has no error")
	}
	select {
	case <-scheduler.Errors():
	case <-ctx.Done():
	}

	if got := scheduler.Size(); got != 0 {
		t.Fatalf("Size() = %v  want: 0 ", got)
	}
}
func TestAsyncTaskScheduler_DoError(t *testing.T) {
	ctx := context.Background()
	scheduler := NewAsyncTaskScheduler(ctx, 1, 0)
	defer scheduler.Close()
	n := 1000
	var err error
	for i := 0; i < n*10; i++ {
		var doErr error
		if i == n-1 {
			doErr = errors.New("mock error")
		}
		if err = scheduler.Push(newMockAsyncTask(
			time.Duration(0), []error{doErr, nil})); err != nil {
			break
		}
	}
	if err == nil {
		t.Fatal("Push Last has no error")
	}

	select {
	case <-scheduler.Errors():
	case <-ctx.Done():
	}
}

func TestAsyncTaskScheduler_Close(t *testing.T) {
	ctx := context.Background()
	scheduler := NewAsyncTaskScheduler(ctx, 1, 0)
	scheduler.Close()
	n := 1000
	var err error
	for i := 0; i < n*10; i++ {
		var doErr error
		if i == n-1 {
			doErr = errors.New("mock error")
		}
		if err = scheduler.Push(newMockAsyncTask(
			time.Duration(0), []error{doErr, nil})); err != nil {
			break
		}
	}
	scheduler.Close()
	if err == nil {
		t.Fatal("Push Last has no error")
	}

	select {
	case <-scheduler.Errors():
	case <-ctx.Done():
	}
}

func TestAsyncTaskScheduler_CancelOne(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	scheduler := NewAsyncTaskScheduler(ctx, 8, 2)
	n := 10000
	go func() {
		time.Sleep(1 * time.Second)
		cancel()
	}()
	for i := 0; i < n; i++ {
		scheduler.Push(newMockAsyncTask(
			1*time.Microsecond, []error{nil, nil}))
	}
	select {
	case <-scheduler.Errors():
	case <-ctx.Done():
	}
	scheduler.Close()
}

func TestAsyncTaskScheduler_CancelMutil(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	scheduler := NewAsyncTaskScheduler(ctx, 1, 0)
	n := 10000
	go func() {
		time.Sleep(1 * time.Second)
		cancel()
	}()
	for i := 0; i < n; i++ {
		scheduler.Push(newMockAsyncTask(
			1*time.Microsecond, []error{nil, nil}))
	}
	select {
	case <-scheduler.Errors():
	case <-ctx.Done():
	}
	scheduler.Close()
}
