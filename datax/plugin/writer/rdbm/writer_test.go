package rdbm

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/datax/common/plugin"
	"github.com/Breeze0806/go-etl/datax/common/spi/writer"
)

type mockJob struct {
	*plugin.BaseJob
}

func (m *mockJob) Init(ctx context.Context) (err error) {
	return
}

func (m *mockJob) Destroy(ctx context.Context) (err error) {
	return
}

func (m *mockJob) Split(ctx context.Context, number int) ([]*config.JSON, error) {
	return nil, nil
}

type mockTask struct {
	*writer.BaseTask
}

func (m *mockTask) Init(ctx context.Context) (err error) {
	return
}

func (m *mockTask) Destroy(ctx context.Context) (err error) {
	return
}

func (m *mockTask) StartWrite(ctx context.Context, receiver plugin.RecordReceiver) (err error) {
	return
}

type mockWriter struct {
	pluginConf *config.JSON
}

func newMockWriter(filename string) (w *mockWriter, err error) {
	w = &mockWriter{}
	w.pluginConf, err = config.NewJSONFromFile(filename)
	if err != nil {
		return nil, err
	}
	return
}

func (w *mockWriter) ResourcesConfig() *config.JSON {
	return w.pluginConf
}

func (w *mockWriter) Job() writer.Job {
	return &mockJob{}
}

func (w *mockWriter) Task() writer.Task {
	return &mockTask{}
}

func TestRegisterWriter(t *testing.T) {
	type args struct {
		new func(string) (Writer, error)
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				new: func(path string) (Writer, error) {
					return newMockWriter(path)
				},
			},
			want: "github.com\\Breeze0806\\go-etl\\datax\\plugin\\writer\\rdbm\\resources\\plugin.json",
		},
		{
			name: "2",
			args: args{
				new: func(path string) (Writer, error) {
					return nil, errors.New("mock error")
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RegisterWriter(tt.args.new)
			if (err != nil) != tt.wantErr {
				t.Errorf("RegisterWriter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !strings.Contains(got, tt.want) {
				t.Errorf("RegisterWriter() = %v, want %v", got, tt.want)
			}
		})
	}
}
