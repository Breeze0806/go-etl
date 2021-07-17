package rdbm

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Breeze0806/go-etl/datax/common/plugin"
	"github.com/Breeze0806/go-etl/datax/core/transport/exchange"
	"github.com/Breeze0806/go-etl/storage/database"
)

type mockBatchWriter struct {
	batchSize    int
	batchTimeout time.Duration
	execer       *MockExecer
	opts         *database.ParameterOptions
}

func newMockBatchWriter(batchSize int, batchTimeout time.Duration,
	execer *MockExecer, opts *database.ParameterOptions) *mockBatchWriter {
	return &mockBatchWriter{
		batchSize:    batchSize,
		batchTimeout: batchTimeout,
		execer:       execer,
		opts:         opts,
	}
}

func (m *mockBatchWriter) JobID() int64 {
	return 0
}

func (m *mockBatchWriter) TaskGroupID() int64 {
	return 0
}

func (m *mockBatchWriter) TaskID() int64 {
	return 0
}

func (m *mockBatchWriter) BatchSize() int {
	return m.batchSize
}

func (m *mockBatchWriter) BatchTimeout() time.Duration {
	return m.batchTimeout
}

func (m *mockBatchWriter) BatchWrite(ctx context.Context) error {
	return m.execer.BatchExec(ctx, m.opts)
}

func (m *mockBatchWriter) Options() *database.ParameterOptions {
	return m.opts
}

func TestStartWrite(t *testing.T) {
	type args struct {
		ctx      context.Context
		writer   BatchWriter
		receiver plugin.RecordReceiver
	}
	tests := []struct {
		name    string
		args    args
		waitCtx time.Duration
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				ctx:      context.TODO(),
				receiver: NewMockReceiver(1000, exchange.ErrTerminate, 1*time.Millisecond),
				writer:   newMockBatchWriter(1000, 1*time.Second, &MockExecer{}, &database.ParameterOptions{}),
			},
		},
		{
			name: "2",
			args: args{
				ctx:      context.TODO(),
				receiver: NewMockReceiverWithoutWait(10000, exchange.ErrTerminate),
				writer:   newMockBatchWriter(1000, 1*time.Second, &MockExecer{}, &database.ParameterOptions{}),
			},
		},
		{
			name: "3",
			args: args{
				ctx:      context.TODO(),
				receiver: NewMockReceiverWithoutWait(10000, errors.New("mock error")),
				writer:   newMockBatchWriter(1000, 1*time.Second, &MockExecer{}, &database.ParameterOptions{}),
			},
			wantErr: true,
		},
		{
			name: "4",
			args: args{
				ctx:      context.TODO(),
				receiver: NewMockReceiverWithoutWait(10000, errors.New("mock error")),
				writer:   newMockBatchWriter(1000, 1*time.Second, &MockExecer{}, &database.ParameterOptions{}),
			},
			waitCtx: 5 * time.Microsecond,
			wantErr: false,
		},
		{
			name: "5",
			args: args{
				ctx:      context.TODO(),
				receiver: NewMockReceiver(1000, exchange.ErrTerminate, 1*time.Millisecond),
				writer: newMockBatchWriter(1000, 1*time.Second, &MockExecer{
					BatchErr: errors.New("mock error"),
					BatchN:   1,
				}, &database.ParameterOptions{}),
			},
			wantErr: true,
		},
		{
			name: "6",
			args: args{
				ctx:      context.TODO(),
				receiver: NewMockReceiverWithoutWait(10000, exchange.ErrTerminate),
				writer: newMockBatchWriter(1000, 1*time.Second, &MockExecer{
					BatchErr: errors.New("mock error"),
					BatchN:   1,
				}, &database.ParameterOptions{}),
			},
			wantErr: true,
		},
		{
			name: "7",
			args: args{
				ctx:      context.TODO(),
				receiver: NewMockReceiver(2, exchange.ErrTerminate, 1*time.Millisecond),
				writer: newMockBatchWriter(1000, 1*time.Second, &MockExecer{
					BatchErr: errors.New("mock error"),
					BatchN:   0,
				}, &database.ParameterOptions{}),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(tt.args.ctx)
			defer cancel()
			if tt.waitCtx != 0 {
				go func() {
					<-time.After(tt.waitCtx)
					cancel()
				}()
			}

			if err := StartWrite(ctx, tt.args.writer, tt.args.receiver); (err != nil) != tt.wantErr {
				t.Errorf("StartWrite() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
