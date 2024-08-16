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

package plugin

import (
	"errors"
	"testing"

	"github.com/Breeze0806/go-etl/element"
)

type mockTaskCollector struct {
}

func (m *mockTaskCollector) CollectDirtyRecordWithError(record element.Record, err error) {
	return
}
func (m *mockTaskCollector) CollectDirtyRecordWithMsg(record element.Record, msgErr string) {
	return
}
func (m *mockTaskCollector) CollectDirtyRecord(record element.Record, err error, msgErr string) {
	return
}
func (m *mockTaskCollector) CollectMessage(key string, value string) {
	return
}
func TestBaseTask_SetTaskID(t *testing.T) {
	type args struct {
		taskID int64
	}
	tests := []struct {
		name string
		b    *BaseTask
		args args
		want int64
	}{
		{
			name: "1",
			b:    NewBaseTask(),
			args: args{
				taskID: 1,
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.SetTaskID(tt.args.taskID)
			if tt.b.TaskID() != tt.want {
				t.Errorf("TaskId() = %v want %v", tt.b.TaskID(), tt.want)
			}
		})
	}
}

func TestBaseTask_SetTaskGroupID(t *testing.T) {
	type args struct {
		taskGroupID int64
	}
	tests := []struct {
		name string
		b    *BaseTask
		args args
		want int64
	}{
		{
			name: "1",
			b:    NewBaseTask(),
			args: args{
				taskGroupID: 1,
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.SetTaskGroupID(tt.args.taskGroupID)
			if tt.b.TaskGroupID() != tt.want {
				t.Errorf("TaskGroupId() = %v want %v", tt.b.TaskGroupID(), tt.want)
			}
		})
	}
}

func TestBaseTask_SetTaskCollector(t *testing.T) {
	type args struct {
		collector TaskCollector
	}
	tests := []struct {
		name string
		b    *BaseTask
		args args
		want TaskCollector
	}{
		{
			name: "1",
			b:    NewBaseTask(),
			args: args{
				collector: &mockTaskCollector{},
			},
			want: &mockTaskCollector{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.SetTaskCollector(tt.args.collector)
			if tt.b.TaskCollector() != tt.want {
				t.Errorf("TaskCollector() = %p want %p", tt.b.TaskCollector(), tt.want)
			}
		})
	}
}

func TestBaseTask_SetJobID(t *testing.T) {
	type fields struct {
		BasePlugin  *BasePlugin
		jobID       int64
		taskID      int64
		taskGroupID int64
		collector   TaskCollector
	}
	type args struct {
		jobID int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int64
	}{
		{
			name:   "1",
			fields: fields{},
			args: args{
				jobID: 1,
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BaseTask{
				BasePlugin:  tt.fields.BasePlugin,
				jobID:       tt.fields.jobID,
				taskID:      tt.fields.taskID,
				taskGroupID: tt.fields.taskGroupID,
				collector:   tt.fields.collector,
			}
			b.SetJobID(tt.args.jobID)
			if b.JobID() != tt.want {
				t.Errorf("JobID() = %v want %v", b.JobID(), tt.want)
			}
		})
	}
}

func TestBaseTask_Wrapf(t *testing.T) {
	type args struct {
		err    error
		format string
		args   []any
	}
	tests := []struct {
		name    string
		b       *BaseTask
		args    args
		wantErr bool
	}{
		{
			name: "1",
			b:    &BaseTask{},
			args: args{
				err:    nil,
				format: "mock",
			},
		},
		{
			name: "1",
			b:    &BaseTask{},
			args: args{
				err:    errors.New("mock error"),
				format: "mock",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.b.Wrapf(tt.args.err, tt.args.format, tt.args.args...); (err != nil) != tt.wantErr {
				t.Errorf("BaseTask.Wrapf() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
