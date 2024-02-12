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
	"time"
)

// RetryTask retry task
type RetryTask struct {
	ctx      context.Context
	task     Task
	strategy RetryStrategy
}

// NewRetryTask generates retry task based on context relationship ctx
func NewRetryTask(ctx context.Context, strategy RetryStrategy, task Task) *RetryTask {
	return &RetryTask{
		ctx:      ctx,
		strategy: strategy,
		task:     task,
	}
}

// Do synchronous execution
func (r *RetryTask) Do() (err error) {
	ticker := time.NewTicker(1)
	defer ticker.Stop()
	var before time.Duration
	for i := 1; ; i++ {
		select {
		case <-r.ctx.Done():
			if err == nil {
				err = r.ctx.Err()
			}
			return
		default:
		}

		err = r.task.Do()

		retry, wait := r.strategy.Next(err, i)
		if !retry {
			return
		}

		if wait != before {
			ticker.Reset(wait)
			before = wait
		}

		select {
		case <-ticker.C:
		case <-r.ctx.Done():
			return
		}
	}
}
