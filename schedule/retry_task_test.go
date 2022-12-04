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
	"errors"
	"testing"
	"time"
)

type mockNTimeTask struct {
	err error
	n   int
}

func (m *mockNTimeTask) Do() error {
	m.n--
	if m.n == 0 {
		return m.err
	}
	return errors.New("mock error")
}

func TestRetryTask_Do(t *testing.T) {
	type args struct {
		ctx      context.Context
		strategy RetryStrategy
		task     Task
		timeout  time.Duration
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				ctx:      context.TODO(),
				strategy: NewNTimesRetryStrategy(&mockRetryJudger{}, 10, 1*time.Millisecond),
				task: &mockNTimeTask{
					n: 2,
				},
			},
			wantErr: false,
		},
		{
			name: "2",
			args: args{
				ctx:      context.TODO(),
				strategy: NewNTimesRetryStrategy(&mockRetryJudger{}, 10, 1*time.Millisecond),
				task: &mockNTimeTask{
					n: 11,
				},
			},
			wantErr: true,
		},
		{
			name: "2",
			args: args{
				ctx:      context.TODO(),
				strategy: NewNTimesRetryStrategy(&mockRetryJudger{}, 10, 1*time.Millisecond),
				task: &mockNTimeTask{
					n: 11,
				},
			},
			wantErr: true,
		},
		{
			name: "3",
			args: args{
				ctx:      context.TODO(),
				strategy: NewNTimesRetryStrategy(&mockRetryJudger{}, 10, 2*time.Millisecond),
				task: &mockNTimeTask{
					n: 11,
				},
				timeout: 1 * time.Nanosecond,
			},
			wantErr: true,
		},
		{
			name: "4",
			args: args{
				ctx:      context.TODO(),
				strategy: NewNTimesRetryStrategy(&mockRetryJudger{}, 10, 2*time.Millisecond),
				task: &mockNTimeTask{
					n: 11,
				},
				timeout: 2 * time.Millisecond,
			},
			wantErr: true,
		},
		{
			name: "4",
			args: args{
				ctx:      context.TODO(),
				strategy: NewNTimesRetryStrategy(&mockRetryJudger{}, 10, 2*time.Millisecond),
				task: &mockNTimeTask{
					n: 11,
				},
				timeout: 2*time.Millisecond + 1*time.Nanosecond,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(tt.args.ctx)
			defer func() {
				if tt.args.timeout == 0 {
					cancel()
				}
			}()
			go func() {
				if tt.args.timeout != 0 {
					<-time.After(tt.args.timeout)
					cancel()
				}
			}()
			r := NewRetryTask(ctx, tt.args.strategy, tt.args.task)
			if err := r.Do(); (err != nil) != tt.wantErr {
				t.Errorf("RetryTask.Do() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
