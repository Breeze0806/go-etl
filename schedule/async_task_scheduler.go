// Copyright 2020 the go-etl Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package schedule

import (
	"context"
	"sync"
	goatomic "sync/atomic"

	"go.uber.org/atomic"
)

type asyncTaskWrapper struct {
	task   AsyncTask
	result chan error
}

type asyncTaskResult struct {
	task   AsyncTask
	result chan error
}

// AsyncTaskScheduler: Asynchronous task scheduler
type AsyncTaskScheduler struct {
	tasks   chan *asyncTaskWrapper
	results chan *asyncTaskResult
	errors  chan error
	wg      sync.WaitGroup
	ctx     context.Context
	cancel  context.CancelFunc
	size    *atomic.Int32
	closed  int32
}

// NewAsyncTaskScheduler: Create an asynchronous task scheduler using the context ctx, the number of parallel workers numWorker, and the task channel size chanSize.
// Create an asynchronous task scheduler with the specified context
func NewAsyncTaskScheduler(ctx context.Context,
	numWorker, chanSize int) *AsyncTaskScheduler {
	a := &AsyncTaskScheduler{
		tasks:   make(chan *asyncTaskWrapper, chanSize),
		results: make(chan *asyncTaskResult, chanSize),
		errors:  make(chan error, 1),
		size:    atomic.NewInt32(0),
		closed:  0,
	}

	a.ctx, a.cancel = context.WithCancel(ctx)

	for i := 0; i < numWorker; i++ {
		a.wg.Add(1)
		go func() {
			defer a.wg.Done()
			a.processTask()
		}()
	}

	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		a.errors <- a.processResult()
		close(a.errors)
	}()

	return a
}

// Push: Asynchronously execute a task.
func (a *AsyncTaskScheduler) Push(task AsyncTask) (err error) {
	if goatomic.CompareAndSwapInt32(&a.closed, 1, 1) {
		return ErrClosed
	}

	result := make(chan error, 1)
	select {
	case <-a.ctx.Done():
		return a.ctx.Err()
	case err = <-a.errors:
		return
	case a.tasks <- &asyncTaskWrapper{
		task:   task,
		result: result,
	}:
	}

	select {
	case <-a.ctx.Done():
		return a.ctx.Err()
	case err = <-a.errors:
		return
	case a.results <- &asyncTaskResult{
		task:   task,
		result: result,
	}:
	}

	a.size.Inc()
	return nil
}

// Size: The number of tasks currently in the asynchronous task scheduler.
func (a *AsyncTaskScheduler) Size() int32 {
	return a.size.Load()
}

// Errors: Error listener for the asynchronous task scheduler.
func (a *AsyncTaskScheduler) Errors() <-chan error {
	return a.errors
}

// Close: Close the asynchronous task scheduler.
func (a *AsyncTaskScheduler) Close() error {
	if !goatomic.CompareAndSwapInt32(&a.closed, 0, 1) {
		return ErrClosed
	}
	a.cancel()
	a.wg.Wait()

	return nil
}

func (a *AsyncTaskScheduler) processTask() {
	for {
		select {
		case <-a.ctx.Done():
			return
		case task, ok := <-a.tasks:
			if !ok {
				return
			}
			select {
			case <-a.ctx.Done():
				return
			case task.result <- task.task.Do():
				close(task.result)
			}
		}
	}
}

func (a *AsyncTaskScheduler) processResult() error {
	for {
		select {
		case <-a.ctx.Done():
			return a.ctx.Err()
		case res, ok := <-a.results:
			if !ok {
				return nil
			}
			select {
			case <-a.ctx.Done():
				return a.ctx.Err()
			case doRes, ok := <-res.result:
				if !ok {
					return nil
				}
				defer a.size.Dec()
				if doRes != nil {
					return doRes
				}

				if err := res.task.Post(); err != nil {
					return err
				}

			}
		}
	}
}
