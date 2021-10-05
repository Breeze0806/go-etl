package writer

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Breeze0806/go-etl/datax/common/plugin"
	"github.com/Breeze0806/go-etl/datax/core/transport/exchange"
	"github.com/Breeze0806/go-etl/element"
)

type mockBatchWriter struct {
	err          error
	n            int
	batchSize    int
	batchTimeout time.Duration
}

func (m *mockBatchWriter) JobID() int64 {
	return 1
}
func (m *mockBatchWriter) TaskGroupID() int64 {
	return 1
}
func (m *mockBatchWriter) TaskID() int64 {
	return 1
}
func (m *mockBatchWriter) BatchSize() int {
	return m.batchSize
}
func (m *mockBatchWriter) BatchTimeout() time.Duration {
	return m.batchTimeout
}
func (m *mockBatchWriter) BatchWrite(ctx context.Context, records []element.Record) error {
	m.n--
	if m.n <= 0 {
		return m.err
	}
	return nil
}

//MockReceiver 模拟接受器
type MockReceiver struct {
	err    error
	n      int
	ticker *time.Ticker
}

//NewMockReceiver 新建等待模拟接受器
func NewMockReceiver(n int, err error, wait time.Duration) *MockReceiver {
	return &MockReceiver{
		err:    err,
		n:      n,
		ticker: time.NewTicker(wait),
	}
}

//NewMockReceiverWithoutWait 新建无等待模拟接受器
func NewMockReceiverWithoutWait(n int, err error) *MockReceiver {
	return &MockReceiver{
		err: err,
		n:   n,
	}
}

//GetFromReader 从读取器获取记录
func (m *MockReceiver) GetFromReader() (element.Record, error) {
	m.n--
	if m.n <= 0 {
		return nil, m.err
	}
	if m.ticker != nil {
		select {
		case <-m.ticker.C:
			return element.NewDefaultRecord(), nil
		}
	}
	return element.NewDefaultRecord(), nil
}

//Shutdown 关闭
func (m *MockReceiver) Shutdown() error {
	m.ticker.Stop()
	return nil
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
				writer: &mockBatchWriter{
					batchSize:    1000,
					batchTimeout: 1 * time.Second,
				},
			},
		},
		{
			name: "2",
			args: args{
				ctx:      context.TODO(),
				receiver: NewMockReceiverWithoutWait(10000, errors.New("mock error")),
				writer: &mockBatchWriter{
					batchSize:    1000,
					batchTimeout: 1 * time.Second,
				},
			},
			waitCtx: 5 * time.Microsecond,
			wantErr: false,
		},
		{
			name: "3",
			args: args{
				ctx:      context.TODO(),
				receiver: NewMockReceiver(1000, exchange.ErrTerminate, 1*time.Millisecond),
				writer: &mockBatchWriter{
					err:          errors.New("mock error"),
					n:            1,
					batchSize:    1000,
					batchTimeout: 1 * time.Second,
				},
			},
			wantErr: true,
		},
		{
			name: "4",
			args: args{
				ctx:      context.TODO(),
				receiver: NewMockReceiverWithoutWait(10000, exchange.ErrTerminate),
				writer: &mockBatchWriter{
					err:          errors.New("mock error"),
					n:            1,
					batchSize:    1000,
					batchTimeout: 1 * time.Second,
				},
			},
			wantErr: true,
		},
		{
			name: "5",
			args: args{
				ctx:      context.TODO(),
				receiver: NewMockReceiver(2, exchange.ErrTerminate, 1*time.Millisecond),
				writer: &mockBatchWriter{
					err:          errors.New("mock error"),
					n:            0,
					batchSize:    1000,
					batchTimeout: 1 * time.Second,
				},
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
