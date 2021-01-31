package plugin

import (
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
		taskID int
	}
	tests := []struct {
		name string
		b    *BaseTask
		args args
		want int
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
		taskGroupID int
	}
	tests := []struct {
		name string
		b    *BaseTask
		args args
		want int
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
