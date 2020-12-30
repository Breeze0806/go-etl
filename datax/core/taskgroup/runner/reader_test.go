package runner

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/Breeze0806/go-etl/datax/common/plugin"
)

var errMockTest = errors.New("mock test error")

func TestReader_Run(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		r       *Reader
		args    args
		wantErr bool
	}{
		{
			name: "1",
			r: NewReader(newMockReaderTask([]error{
				nil, nil, nil, nil, nil,
			}), &mockRecordSender{}),
			args: args{
				ctx: context.TODO(),
			},
			wantErr: false,
		},
		{
			name: "2",
			r: NewReader(newMockReaderTask([]error{
				errMockTest, nil, nil, nil, nil,
			}), &mockRecordSender{}),
			args: args{
				ctx: context.TODO(),
			},
			wantErr: true,
		},
		{
			name: "3",
			r: NewReader(newMockReaderTask([]error{
				nil, errMockTest, nil, nil, nil,
			}), &mockRecordSender{}),
			args: args{
				ctx: context.TODO(),
			},
			wantErr: true,
		},
		{
			name: "4",
			r: NewReader(newMockReaderTask([]error{
				nil, nil, errMockTest, nil, nil,
			}), &mockRecordSender{}),
			args: args{
				ctx: context.TODO(),
			},
			wantErr: true,
		},
		{
			name: "5",
			r: NewReader(newMockReaderTask([]error{
				nil, nil, nil, errMockTest, nil,
			}), &mockRecordSender{}),
			args: args{
				ctx: context.TODO(),
			},
			wantErr: true,
		},
		{
			name: "6",
			r: NewReader(newMockReaderTask([]error{
				nil, nil, nil, nil, errMockTest,
			}), &mockRecordSender{}),
			args: args{
				ctx: context.TODO(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.Run(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Reader.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestReader_Plugin(t *testing.T) {
	tests := []struct {
		name string
		r    *Reader
		want plugin.Task
	}{
		{
			name: "1",
			r: NewReader(newMockReaderTask([]error{
				nil, nil, nil, nil, nil,
			}), &mockRecordSender{}),
			want: newMockReaderTask([]error{
				nil, nil, nil, nil, nil,
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.Plugin(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reader.Plugin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReader_Shutdown(t *testing.T) {
	tests := []struct {
		name    string
		r       *Reader
		wantErr bool
	}{
		{
			name: "1",
			r: NewReader(newMockReaderTask([]error{
				nil, nil, nil, nil, nil,
			}), &mockRecordSender{}),

			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.Shutdown(); (err != nil) != tt.wantErr {
				t.Errorf("Reader.Shutdown() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
