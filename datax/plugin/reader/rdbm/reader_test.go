package rdbm

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/datax/common/plugin"
	"github.com/Breeze0806/go-etl/datax/common/spi/reader"
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
	*plugin.BaseTask
}

func (m *mockTask) Init(ctx context.Context) (err error) {
	return
}

func (m *mockTask) Destroy(ctx context.Context) (err error) {
	return
}

func (m *mockTask) StartRead(ctx context.Context, sender plugin.RecordSender) (err error) {
	return
}

type mockReader struct {
	pluginConf *config.JSON
}

func NewMockReader(filename string) (r *mockReader, err error) {
	r = &mockReader{}
	r.pluginConf, err = config.NewJSONFromFile(filename)
	if err != nil {
		return nil, err
	}
	return
}

func (r *mockReader) ResourcesConfig() *config.JSON {
	return r.pluginConf
}

func (r *mockReader) Job() reader.Job {
	return &mockJob{}
}
func (r *mockReader) Task() reader.Task {
	return &mockTask{}
}

func TestRegisterReader(t *testing.T) {
	type args struct {
		new func(string) (Reader, error)
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
				new: func(path string) (Reader, error) {
					return NewMockReader(path)
				},
			},
			want: "github.com\\Breeze0806\\go-etl\\datax\\plugin\\reader\\rdbm\\resources\\plugin.json",
		},
		{
			name: "2",
			args: args{
				new: func(path string) (Reader, error) {
					return nil, errors.New("mock error")
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RegisterReader(tt.args.new)
			if (err != nil) != tt.wantErr {
				t.Errorf("RegisterReader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !strings.Contains(got, tt.want) {
				t.Errorf("RegisterReader() = %v, want %v", got, tt.want)
			}
		})
	}
}
