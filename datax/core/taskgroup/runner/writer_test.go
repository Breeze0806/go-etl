package runner

import (
	"context"
	"reflect"
	"testing"

	"github.com/Breeze0806/go-etl/datax/common/plugin"
)

func TestWriter_Run(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		w       *Writer
		args    args
		wantErr bool
	}{
		{
			name: "1",
			w: NewWriter(newMockWriterTask([]error{
				nil, nil, nil, nil, nil,
			}), &mockRecordReceiver{}, "mock"),
			args: args{
				ctx: context.TODO(),
			},
			wantErr: false,
		},
		{
			name: "2",
			w: NewWriter(newMockWriterTask([]error{
				errMockTest, nil, nil, nil, nil,
			}), &mockRecordReceiver{}, "mock"),
			args: args{
				ctx: context.TODO(),
			},
			wantErr: true,
		},
		{
			name: "3",
			w: NewWriter(newMockWriterTask([]error{
				nil, errMockTest, nil, nil, nil,
			}), &mockRecordReceiver{}, "mock"),
			args: args{
				ctx: context.TODO(),
			},
			wantErr: true,
		},
		{
			name: "4",
			w: NewWriter(newMockWriterTask([]error{
				nil, nil, errMockTest, nil, nil,
			}), &mockRecordReceiver{}, "mock"),
			args: args{
				ctx: context.TODO(),
			},
			wantErr: true,
		},
		{
			name: "5",
			w: NewWriter(newMockWriterTask([]error{
				nil, nil, nil, errMockTest, nil,
			}), &mockRecordReceiver{}, "mock"),
			args: args{
				ctx: context.TODO(),
			},
			wantErr: true,
		},
		{
			name: "6",
			w: NewWriter(newMockWriterTask([]error{
				nil, nil, nil, nil, errMockTest,
			}), &mockRecordReceiver{}, "mock"),
			args: args{
				ctx: context.TODO(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.w.Run(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Writer.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWriter_Plugin(t *testing.T) {
	tests := []struct {
		name string
		w    *Writer
		want plugin.Task
	}{
		{
			name: "1",
			w: NewWriter(newMockWriterTask([]error{
				nil, nil, nil, nil, nil,
			}), &mockRecordReceiver{}, "mock"),
			want: newMockWriterTask([]error{
				nil, nil, nil, nil, nil,
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.w.Plugin(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Writer.Plugin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWriter_Shutdown(t *testing.T) {
	tests := []struct {
		name    string
		w       *Writer
		wantErr bool
	}{
		{
			name: "1",
			w: NewWriter(newMockWriterTask([]error{
				nil, nil, nil, nil, nil,
			}), &mockRecordReceiver{}, "mock"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.w.Shutdown(); (err != nil) != tt.wantErr {
				t.Errorf("Writer.Shutdown() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
