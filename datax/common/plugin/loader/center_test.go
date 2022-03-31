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

package loader

import (
	"context"
	"reflect"
	"testing"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/datax/common/plugin"
	"github.com/Breeze0806/go-etl/datax/common/spi"
	"github.com/Breeze0806/go-etl/datax/common/spi/reader"
	"github.com/Breeze0806/go-etl/datax/common/spi/writer"
)

type mockTaskPlugin struct {
}

func (m *mockTaskPlugin) Init(ctx context.Context) error {
	return nil
}

func (m *mockTaskPlugin) Destroy(ctx context.Context) error {
	return nil
}

type mockReaderJob struct {
	*defaultJobPlugin
}

func (m *mockReaderJob) Split(ctx context.Context, number int) ([]*config.JSON, error) {
	return nil, nil
}

type mockReaderTask struct {
	*plugin.BaseTask
	*mockTaskPlugin
}

func (m *mockReaderTask) StartRead(ctx context.Context, sender plugin.RecordSender) error {
	return nil
}

type mockWriterJob struct {
	*defaultJobPlugin
}

func (m *mockWriterJob) Split(ctx context.Context, number int) ([]*config.JSON, error) {
	return nil, nil
}

type mockWriterTask struct {
	*mockTaskPlugin
	*writer.BaseTask
}

func (m *mockWriterTask) StartWrite(ctx context.Context, receiver plugin.RecordReceiver) error {
	return nil
}

type mockNilReader struct{}

func (m *mockNilReader) Job() reader.Job {
	return nil
}

func (m *mockNilReader) Task() reader.Task {
	return nil
}

type mockNilWriter struct{}

func (m *mockNilWriter) Job() writer.Job {
	return nil
}

func (m *mockNilWriter) Task() writer.Task {
	return nil
}

type mockReader struct{}

func (m *mockReader) Job() reader.Job {
	return &mockReaderJob{}
}

func (m *mockReader) Task() reader.Task {
	return &mockReaderTask{}
}

type mockWriter struct{}

func (m *mockWriter) Job() writer.Job {
	return &mockWriterJob{}
}

func (m *mockWriter) Task() writer.Task {
	return &mockWriterTask{}
}

func TestRegisterReader(t *testing.T) {
	_centor = &centor{
		readers: map[string]spi.Reader{
			"mock": &mockReader{},
		},
	}

	type args struct {
		name   string
		reader spi.Reader
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "1",
			args:    args{},
			wantErr: true,
		},
		{
			name: "2",
			args: args{
				name:   "mock",
				reader: &mockNilReader{},
			},
			wantErr: true,
		},
		{
			name: "3",
			args: args{
				name:   "mock1",
				reader: &mockReader{},
			},
			wantErr: false,
		},
		{
			name: "4",
			args: args{
				name:   "mock",
				reader: &mockReader{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if err := recover(); (err != nil) != tt.wantErr {
					t.Errorf("panic err: %v wantErr %v", err, tt.wantErr)
				}
			}()

			RegisterReader(tt.args.name, tt.args.reader)
		})
	}
}

func TestRegisterWriter(t *testing.T) {
	_centor = &centor{
		writers: map[string]spi.Writer{
			"mock": &mockWriter{},
		},
	}

	type args struct {
		name   string
		writer spi.Writer
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "1",
			args:    args{},
			wantErr: true,
		},
		{
			name: "2",
			args: args{
				name:   "mock",
				writer: &mockNilWriter{},
			},
			wantErr: true,
		},
		{
			name: "3",
			args: args{
				name:   "mock1",
				writer: &mockWriter{},
			},
			wantErr: false,
		},
		{
			name: "4",
			args: args{
				name:   "mock",
				writer: &mockWriter{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if err := recover(); (err != nil) != tt.wantErr {
					t.Errorf("panic err: %v wantErr %v", err, tt.wantErr)
				}
			}()

			RegisterWriter(tt.args.name, tt.args.writer)
		})
	}
}

func TestLoadJobPlugin(t *testing.T) {
	type args struct {
		typ  plugin.Type
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    plugin.Job
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				typ:  plugin.Writer,
				name: "1111",
			},
			want: newdefaultJobPlugin(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadJobPlugin(tt.args.typ, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadJobPlugin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadJobPlugin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadReaderJob(t *testing.T) {
	_centor = &centor{
		readers: map[string]spi.Reader{
			"mock": &mockReader{},
		},
	}

	type args struct {
		name string
	}
	tests := []struct {
		name  string
		args  args
		want  reader.Job
		want1 bool
	}{
		{
			name: "1",
			args: args{
				name: "mock",
			},
			want:  &mockReaderJob{},
			want1: true,
		},
		{
			name: "2",
			args: args{
				name: "mock1",
			},
			want:  nil,
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := LoadReaderJob(tt.args.name)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadReaderJob() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("LoadReaderJob() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestLoadReaderTask(t *testing.T) {
	_centor = &centor{
		readers: map[string]spi.Reader{
			"mock": &mockReader{},
		},
	}
	type args struct {
		name string
	}
	tests := []struct {
		name  string
		args  args
		want  reader.Task
		want1 bool
	}{
		{
			name: "1",
			args: args{
				name: "mock",
			},
			want:  &mockReaderTask{},
			want1: true,
		},
		{
			name: "2",
			args: args{
				name: "mock1",
			},
			want:  nil,
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := LoadReaderTask(tt.args.name)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadReaderTask() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("LoadReaderTask() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestLoadWriterJob(t *testing.T) {
	_centor = &centor{
		writers: map[string]spi.Writer{
			"mock": &mockWriter{},
		},
	}
	type args struct {
		name string
	}
	tests := []struct {
		name  string
		args  args
		want  writer.Job
		want1 bool
	}{
		{
			name: "1",
			args: args{
				name: "mock",
			},
			want:  &mockWriterJob{},
			want1: true,
		},
		{
			name: "2",
			args: args{
				name: "mock1",
			},
			want:  nil,
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := LoadWriterJob(tt.args.name)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadWriterJob() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("LoadWriterJob() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestLoadWriterTask(t *testing.T) {
	_centor = &centor{
		writers: map[string]spi.Writer{
			"mock": &mockWriter{},
		},
	}
	type args struct {
		name string
	}
	tests := []struct {
		name  string
		args  args
		want  writer.Task
		want1 bool
	}{
		{
			name: "1",
			args: args{
				name: "mock",
			},
			want:  &mockWriterTask{},
			want1: true,
		},
		{
			name: "2",
			args: args{
				name: "mock1",
			},
			want:  nil,
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := LoadWriterTask(tt.args.name)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadWriterTask() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("LoadWriterTask() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_defaultJobPlugin(t *testing.T) {
	d := &defaultJobPlugin{}
	if err := d.Init(context.Background()); err != nil {
		t.Errorf("Init() err = %v", err)
	}
	if err := d.Destroy(context.Background()); err != nil {
		t.Errorf("Destroy() err = %v", err)
	}
}

func TestUnregisterReaders(t *testing.T) {
	_centor = &centor{
		readers: map[string]spi.Reader{
			"mock": &mockReader{},
		},
	}
	tests := []struct {
		name string
	}{
		{
			name: "1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			UnregisterReaders()
		})
	}
}

func TestUnregisterWriters(t *testing.T) {
	_centor = &centor{
		writers: map[string]spi.Writer{
			"mock": &mockWriter{},
		},
	}
	tests := []struct {
		name string
	}{
		{
			name: "1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			UnregisterWriters()
		})
	}
}
