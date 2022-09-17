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

package file

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/datax/common/plugin"
	"github.com/Breeze0806/go-etl/datax/core/transport/exchange"
	"github.com/Breeze0806/go-etl/element"
	"github.com/Breeze0806/go-etl/storage/stream/file"
)

type MockReceiver struct {
	err    error
	n      int
	ticker *time.Ticker
}

func NewMockReceiver(n int, err error, wait time.Duration) *MockReceiver {
	return &MockReceiver{
		err:    err,
		n:      n,
		ticker: time.NewTicker(wait),
	}
}

func NewMockReceiverWithoutWait(n int, err error) *MockReceiver {
	return &MockReceiver{
		err: err,
		n:   n,
	}
}

func (m *MockReceiver) GetFromReader() (element.Record, error) {
	m.n--
	if m.n <= 0 {
		return nil, m.err
	}
	if m.ticker != nil {
		<-m.ticker.C
		return element.NewDefaultRecord(), nil
	}
	return element.NewDefaultRecord(), nil
}

func (m *MockReceiver) Shutdown() error {
	m.ticker.Stop()
	return nil
}

type mockStreamWriter struct {
	record   element.Record
	writeErr error
	flushErr error
	closeErr error
}

func (m *mockStreamWriter) Write(record element.Record) (err error) {
	m.record = record
	return m.writeErr
}

func (m *mockStreamWriter) Flush() (err error) {
	return m.flushErr
}

func (m *mockStreamWriter) Close() (err error) {
	return m.closeErr
}

type mockOutStream struct {
	writer    file.StreamWriter
	writerErr error
}

func (m *mockOutStream) Writer(conf *config.JSON) (writer file.StreamWriter, err error) {
	return m.writer, m.writerErr
}

func (m *mockOutStream) Close() (err error) {
	return nil
}

type mockCreater struct {
	stream file.OutStream
}

func (m *mockCreater) Create(filename string) (stream file.OutStream, err error) {
	return m.stream, nil
}

func TestTask_StartWrite(t *testing.T) {
	file.UnregisterAllCreater()
	file.RegisterCreator("mock", &mockCreater{
		stream: &mockOutStream{
			writer: &mockStreamWriter{},
		},
	})

	file.RegisterCreator("mockWriterErr", &mockCreater{
		stream: &mockOutStream{
			writer:    &mockStreamWriter{},
			writerErr: errors.New("mock error"),
		},
	})
	file.RegisterCreator("mockCloseErr", &mockCreater{
		stream: &mockOutStream{
			writer: &mockStreamWriter{
				closeErr: errors.New("mock error"),
			},
		},
	})
	file.RegisterCreator("mockWriteErr", &mockCreater{
		stream: &mockOutStream{
			writer: &mockStreamWriter{
				writeErr: errors.New("mock error"),
			},
		},
	})
	file.RegisterCreator("mockFlushErr", &mockCreater{
		stream: &mockOutStream{
			writer: &mockStreamWriter{
				flushErr: errors.New("flush error"),
			},
		},
	})
	type args struct {
		ctx      context.Context
		receiver plugin.RecordReceiver
	}
	tests := []struct {
		name    string
		t       *Task
		args    args
		conf    *config.JSON
		jobConf *config.JSON
		wantErr bool
	}{
		{
			name: "1",
			t: NewTask(func(conf *config.JSON) (Config, error) {
				c, err := NewBaseConfig(conf)
				if err != nil {
					return nil, err
				}
				return c, nil
			}),
			args: args{
				ctx:      context.TODO(),
				receiver: NewMockReceiver(1000, exchange.ErrTerminate, 1*time.Millisecond),
			},
			conf:    testJSONFromString(`{"creator":"mock"}`),
			jobConf: testJSONFromString(`{"path":"file1","content":{}}`),
		},
		{
			name: "2",
			t: NewTask(func(conf *config.JSON) (Config, error) {
				c, err := NewBaseConfig(conf)
				if err != nil {
					return nil, err
				}
				return c, nil
			}),
			args: args{
				ctx:      context.TODO(),
				receiver: NewMockReceiverWithoutWait(10000, exchange.ErrTerminate),
			},
			conf:    testJSONFromString(`{"creator":"mock"}`),
			jobConf: testJSONFromString(`{"path":"file1","content":{}}`),
		},
		{
			name: "3",
			t: NewTask(func(conf *config.JSON) (Config, error) {
				c, err := NewBaseConfig(conf)
				if err != nil {
					return nil, err
				}
				return c, nil
			}),
			args: args{
				ctx:      context.TODO(),
				receiver: NewMockReceiverWithoutWait(10, errors.New("mock error")),
			},
			conf:    testJSONFromString(`{"creator":"mock"}`),
			jobConf: testJSONFromString(`{"path":"file1","content":{}}`),
			wantErr: true,
		},
		{
			name: "4",
			t: NewTask(func(conf *config.JSON) (Config, error) {
				c, err := NewBaseConfig(conf)
				if err != nil {
					return nil, err
				}
				return c, nil
			}),
			args: args{
				ctx:      context.TODO(),
				receiver: NewMockReceiver(1, exchange.ErrTerminate, 1*time.Millisecond),
			},
			conf:    testJSONFromString(`{"creator":"mockWriterErr"}`),
			jobConf: testJSONFromString(`{"path":"file1","content":{}}`),
			wantErr: true,
		},
		{
			name: "5",
			t: NewTask(func(conf *config.JSON) (Config, error) {
				c, err := NewBaseConfig(conf)
				if err != nil {
					return nil, err
				}
				return c, nil
			}),
			args: args{
				ctx:      context.TODO(),
				receiver: NewMockReceiver(1, exchange.ErrTerminate, 1*time.Millisecond),
			},
			conf:    testJSONFromString(`{"creator":"mockCloseErr"}`),
			jobConf: testJSONFromString(`{"path":"file1","content":{}}`),
			wantErr: false,
		},
		{
			name: "6",
			t: NewTask(func(conf *config.JSON) (Config, error) {
				c, err := NewBaseConfig(conf)
				if err != nil {
					return nil, err
				}
				return c, nil
			}),
			args: args{
				ctx:      context.TODO(),
				receiver: NewMockReceiver(2, exchange.ErrTerminate, 1*time.Millisecond),
			},
			conf:    testJSONFromString(`{"creator":"mockWriteErr"}`),
			jobConf: testJSONFromString(`{"path":"file1","content":{}}`),
			wantErr: true,
		},
		{
			name: "7",
			t: NewTask(func(conf *config.JSON) (Config, error) {
				c, err := NewBaseConfig(conf)
				if err != nil {
					return nil, err
				}
				return c, nil
			}),
			args: args{
				ctx:      context.TODO(),
				receiver: NewMockReceiver(2, exchange.ErrTerminate, 1*time.Millisecond),
			},
			conf:    testJSONFromString(`{"creator":"mockFlushErr"}`),
			jobConf: testJSONFromString(`{"path":"file1","content":{}}`),
			wantErr: true,
		},
		{
			name: "8",
			t: NewTask(func(conf *config.JSON) (Config, error) {
				c, err := NewBaseConfig(conf)
				if err != nil {
					return nil, err
				}
				return c, nil
			}),
			args: args{
				ctx:      context.TODO(),
				receiver: NewMockReceiverWithoutWait(1001, exchange.ErrTerminate),
			},
			conf:    testJSONFromString(`{"creator":"mockFlushErr"}`),
			jobConf: testJSONFromString(`{"path":"file1","content":{}}`),
			wantErr: true,
		},
		{
			name: "9",
			t: NewTask(func(conf *config.JSON) (Config, error) {
				c, err := NewBaseConfig(conf)
				if err != nil {
					return nil, err
				}
				return c, nil
			}),
			args: args{
				ctx:      context.TODO(),
				receiver: NewMockReceiver(10, exchange.ErrTerminate, 500*time.Millisecond),
			},
			conf:    testJSONFromString(`{"creator":"mockFlushErr"}`),
			jobConf: testJSONFromString(`{"path":"file1","content":{}}`),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.t.SetPluginConf(tt.conf)
			tt.t.SetPluginJobConf(tt.jobConf)
			defer tt.t.Destroy(tt.args.ctx)
			if err := tt.t.Init(tt.args.ctx); err != nil {
				t.Errorf("Task.Init() error = %v", err)
				return
			}
			if err := tt.t.StartWrite(tt.args.ctx, tt.args.receiver); (err != nil) != tt.wantErr {
				t.Errorf("Task.StartWrite() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTask_Init(t *testing.T) {
	file.UnregisterAllCreater()
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		tr      *Task
		args    args
		conf    *config.JSON
		jobConf *config.JSON
		wantErr bool
	}{
		{
			name: "1",
			tr: NewTask(func(conf *config.JSON) (Config, error) {
				c, err := NewBaseConfig(conf)
				if err != nil {
					return nil, err
				}
				return c, nil
			}),
			args: args{
				ctx: context.TODO(),
			},
			conf:    testJSONFromString(`{"creater":"mock"}`),
			jobConf: testJSONFromString(`{"path":"file1","content":{}}`),
			wantErr: true,
		},
		{
			name: "2",
			tr: NewTask(func(conf *config.JSON) (Config, error) {
				c, err := NewBaseConfig(conf)
				if err != nil {
					return nil, err
				}
				return c, nil
			}),
			args: args{
				ctx: context.TODO(),
			},
			conf:    testJSONFromString(`{"creator":"mock"}`),
			jobConf: testJSONFromString(`{"path1":"file1","content":{}}`),
			wantErr: true,
		},
		{
			name: "3",
			tr: NewTask(func(conf *config.JSON) (Config, error) {
				c, err := NewBaseConfig(conf)
				if err != nil {
					return nil, err
				}
				return c, nil
			}),
			args: args{
				ctx: context.TODO(),
			},
			conf:    testJSONFromString(`{"creator":"mock"}`),
			jobConf: testJSONFromString(`{"path":"file1","content1":{}}`),
			wantErr: true,
		},
		{
			name: "4",
			tr: NewTask(func(conf *config.JSON) (Config, error) {
				c, err := NewBaseConfig(conf)
				if err != nil {
					return nil, err
				}
				return c, nil
			}),
			args: args{
				ctx: context.TODO(),
			},
			conf:    testJSONFromString(`{"creator":"mock"}`),
			jobConf: testJSONFromString(`{"path":"file1","content":{"batchTimeout":"2"}}`),
			wantErr: true,
		},
		{
			name: "5",
			tr: NewTask(func(conf *config.JSON) (Config, error) {
				c, err := NewBaseConfig(conf)
				if err != nil {
					return nil, err
				}
				return c, nil
			}),
			args: args{
				ctx: context.TODO(),
			},
			conf:    testJSONFromString(`{"creator":"mock"}`),
			jobConf: testJSONFromString(`{"path":"file1","content":{}}`),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.tr.SetPluginConf(tt.conf)
			tt.tr.SetPluginJobConf(tt.jobConf)
			defer tt.tr.Destroy(tt.args.ctx)
			if err := tt.tr.Init(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Task.Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTask_StartWriteCancel(t *testing.T) {
	file.UnregisterAllCreater()
	file.RegisterCreator("mock", &mockCreater{
		stream: &mockOutStream{
			writer: &mockStreamWriter{},
		},
	})

	type args struct {
		ctx      context.Context
		receiver plugin.RecordReceiver
	}
	tests := []struct {
		name    string
		t       *Task
		args    args
		conf    *config.JSON
		jobConf *config.JSON
		wantErr bool
	}{
		{
			name: "1",
			t: NewTask(func(conf *config.JSON) (Config, error) {
				c, err := NewBaseConfig(conf)
				if err != nil {
					return nil, err
				}
				return c, nil
			}),
			args: args{
				ctx:      context.Background(),
				receiver: NewMockReceiver(1000, exchange.ErrTerminate, 10*time.Millisecond),
			},
			conf:    testJSONFromString(`{"creator":"mock"}`),
			jobConf: testJSONFromString(`{"path":"file1","content":{}}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.t.SetPluginConf(tt.conf)
			tt.t.SetPluginJobConf(tt.jobConf)
			ctx, cancel := context.WithCancel(tt.args.ctx)
			defer tt.t.Destroy(ctx)
			go func() {
				<-time.After(100 * time.Millisecond)
				cancel()
			}()
			if err := tt.t.Init(ctx); err != nil {
				t.Errorf("Task.Init() error = %v", err)
				return
			}
			if err := tt.t.StartWrite(ctx, tt.args.receiver); (err != nil) != tt.wantErr {
				t.Errorf("Task.StartWrite() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
