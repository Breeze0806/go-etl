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

type mockStreamWriter struct {
	record element.Record
}

func (m *mockStreamWriter) Write(record element.Record) (err error) {
	m.record = record
	return
}

func (m *mockStreamWriter) Flush() (err error) {
	return
}

func (m *mockStreamWriter) Close() (err error) {
	return
}

type mockOutStream struct {
}

func (m *mockOutStream) Writer(conf *config.JSON) (writer file.StreamWriter, err error) {
	return &mockStreamWriter{}, nil
}

func (m *mockOutStream) Close() (err error) {
	return
}

type mockCreater struct {
}

func (m *mockCreater) Create(filename string) (stream file.OutStream, err error) {
	return &mockOutStream{}, nil
}
func TestTask_StartWrite(t *testing.T) {
	file.RegisterCreater("mock", &mockCreater{})

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
			conf:    testJSONFromString(`{"creater":"mock"}`),
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
			conf:    testJSONFromString(`{"creater":"mock"}`),
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
				receiver: NewMockReceiverWithoutWait(10000, errors.New("mock error")),
			},
			conf:    testJSONFromString(`{"creater":"mock"}`),
			jobConf: testJSONFromString(`{"path":"file1","content":{}}`),
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
