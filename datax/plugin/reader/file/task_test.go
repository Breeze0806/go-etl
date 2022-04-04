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
	"testing"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/datax/common/plugin"
	"github.com/Breeze0806/go-etl/element"
	"github.com/Breeze0806/go-etl/storage/stream/file"
)

func testJSONFromString(json string) *config.JSON {
	conf, err := config.NewJSONFromString(json)
	if err != nil {
		panic(err)
	}
	return conf
}

type mockRows struct {
	n int
}

func (m *mockRows) Next() bool {
	m.n++
	return m.n <= 1
}

func (m *mockRows) Scan() (columns []element.Column, err error) {
	columns = append(columns, element.NewDefaultColumn(element.NewStringColumnValue("mock"),
		"mock", 0))
	return
}

func (m *mockRows) Error() error {
	return nil
}

func (m *mockRows) Close() error {
	return nil
}

type mockInStream struct {
}

func (m *mockInStream) Rows(conf *config.JSON) (rows file.Rows, err error) {
	return &mockRows{}, nil
}

func (m *mockInStream) Close() (err error) {
	return
}

type mockOpener struct {
}

func (m *mockOpener) Open(filename string) (stream file.InStream, err error) {
	return &mockInStream{}, nil
}

//MockSender 模拟发送器
type MockSender struct {
	record    element.Record
	CreateErr error
	SendErr   error
}

//CreateRecord 创建记录
func (m *MockSender) CreateRecord() (element.Record, error) {
	return element.NewDefaultRecord(), m.CreateErr
}

//SendWriter 发往写入器
func (m *MockSender) SendWriter(record element.Record) error {
	m.record = record
	return m.SendErr
}

//Flush 刷新至写入器
func (m *MockSender) Flush() error {
	return nil
}

//Terminate 终止发送数据
func (m *MockSender) Terminate() error {
	return nil
}

//Shutdown 关闭
func (m *MockSender) Shutdown() error {
	return nil
}
func TestTask_StartRead(t *testing.T) {
	file.RegisterOpener("mock", &mockOpener{})
	type args struct {
		ctx    context.Context
		sender plugin.RecordSender
	}
	tests := []struct {
		name    string
		conf    *config.JSON
		jobConf *config.JSON
		args    args
		wantErr bool
	}{
		{
			name:    "1",
			conf:    testJSONFromString(`{"opener":"mock"}`),
			jobConf: testJSONFromString(`{"path":"mockfile","content":[{},{}]}`),
			args: args{
				ctx:    context.TODO(),
				sender: &MockSender{},
			},
		},

		{
			name:    "2",
			conf:    testJSONFromString(`{"opener":"mock"}`),
			jobConf: testJSONFromString(`{"path":"mockfile"}`),
			args: args{
				ctx:    context.TODO(),
				sender: &MockSender{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := NewTask()
			defer task.Destroy(tt.args.ctx)
			task.SetPluginConf(tt.conf)
			task.SetPluginJobConf(tt.jobConf)
			if err := task.Init(tt.args.ctx); err != nil {
				t.Errorf("init fail. err: %v", err)
				return
			}
			if err := task.StartRead(tt.args.ctx, tt.args.sender); (err != nil) != tt.wantErr {
				t.Errorf("StartRead fail. err: %v  wantErr: %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if c, _ := tt.args.sender.(*MockSender).record.GetByName("mock"); c.String() != "mock" {
					t.Errorf("StartRead() fail")
				}
			}
		})
	}
}

func TestTask_Init(t *testing.T) {
	file.RegisterOpener("mock1", &mockOpener{})
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		conf    *config.JSON
		jobConf *config.JSON
		args    args
		wantErr bool
	}{
		{
			name:    "1",
			conf:    testJSONFromString(`{"opener1":"mock2"}`),
			jobConf: testJSONFromString(`{"path":"mockfile","content":[{},{}]}`),
			args: args{
				ctx: context.TODO(),
			},
			wantErr: true,
		},
		{
			name:    "2",
			conf:    testJSONFromString(`{"opener":"mock1"}`),
			jobConf: testJSONFromString(`{"path1":"mockfile","content":[{},{}]}`),
			args: args{
				ctx: context.TODO(),
			},
			wantErr: true,
		},
		{
			name:    "3",
			conf:    testJSONFromString(`{"opener":"mock2"}`),
			jobConf: testJSONFromString(`{"path":"mockfile","content":[{},{}]}`),
			args: args{
				ctx: context.TODO(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := NewTask()
			defer task.Destroy(tt.args.ctx)
			task.SetPluginConf(tt.conf)
			task.SetPluginJobConf(tt.jobConf)
			if err := task.Init(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Task.Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
