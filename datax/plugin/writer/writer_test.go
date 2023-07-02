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

package writer

import (
	"context"
	"testing"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/datax/common/plugin"
	"github.com/Breeze0806/go-etl/datax/common/plugin/loader"
	"github.com/Breeze0806/go-etl/datax/common/spi/writer"
	"github.com/pingcap/errors"
)

type mockMaker struct {
	err  error
	conf *config.JSON
}

func newMockMaker(err error, conf *config.JSON) *mockMaker {
	return &mockMaker{
		err:  err,
		conf: conf,
	}
}

func (m *mockMaker) Default() (Writer, error) {
	return &mockWriter{conf: m.conf}, m.err
}

type mockWriter struct {
	conf *config.JSON
}

func (m *mockWriter) Job() writer.Job {
	return newMockJob()
}

func (m *mockWriter) Task() writer.Task {
	return newMockTask()
}

func (m *mockWriter) ResourcesConfig() *config.JSON {
	return m.conf
}

type mockJob struct {
	*plugin.BaseJob
}

func newMockJob() *mockJob {
	return &mockJob{
		BaseJob: plugin.NewBaseJob(),
	}
}

func (m *mockJob) Split(ctx context.Context, number int) ([]*config.JSON, error) {
	return nil, nil
}

func (m *mockJob) Init(ctx context.Context) error {
	return nil
}

func (m *mockJob) Destroy(ctx context.Context) error {
	return nil
}

type mockTask struct {
	*writer.BaseTask
}

func newMockTask() *mockTask {
	return &mockTask{
		BaseTask: writer.NewBaseTask(),
	}
}

func (m *mockTask) Init(ctx context.Context) error {
	return nil
}

func (m *mockTask) Destroy(ctx context.Context) error {
	return nil
}

func (m *mockTask) StartWrite(ctx context.Context, receiver plugin.RecordReceiver) error {
	return nil
}

func testJSONFromString(s string) *config.JSON {
	j, err := config.NewJSONFromString(s)
	if err != nil {
		panic(err)
	}
	return j
}

func TestRegisterWriter(t *testing.T) {
	type args struct {
		maker Maker
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    bool
	}{
		{
			name: "1",
			args: args{
				maker: newMockMaker(nil, testJSONFromString(`{
					"name" : "mockwriter",
					"developer":"Breeze0806",
					"dialect":"mock",
					"description":""
				}`)),
			},
			want: true,
		},
		{
			name: "2",
			args: args{
				maker: newMockMaker(errors.New("mock error"), testJSONFromString(`{}`)),
			},
			wantErr: true,
		},
		{
			name: "3",
			args: args{
				maker: newMockMaker(nil, testJSONFromString(`{}`)),
			},
			wantErr: true,
		},
		{
			name: "4",
			args: args{
				maker: newMockMaker(nil, testJSONFromString(`{
					"name":""
				}`)),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := RegisterWriter(tt.args.maker)

			if (err != nil) != tt.wantErr {
				t.Errorf("RegisterWriter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return
			}

			if _, got := loader.LoadWriterTask("mockwriter"); got != tt.want {
				t.Errorf("mockwriter() got = %v, want %v", got, tt.want)
			}
		})
	}
}
