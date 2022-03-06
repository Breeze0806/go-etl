package reader

import (
	"context"
	"os"
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

func newMockReader(filename string) (r *mockReader, err error) {
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

type mockReaderMaker1 struct {
	err error
}

func (m *mockReaderMaker1) FromFile(path string) (Reader, error) {
	return newMockReader(path)
}

func (m *mockReaderMaker1) Default() (Reader, error) {
	return nil, nil
}

type mockReaderMaker2 struct {
	path string
	err  error
}

func (m *mockReaderMaker2) FromFile(path string) (Reader, error) {
	m.path = path
	return nil, os.ErrNotExist
}

func (m *mockReaderMaker2) Default() (Reader, error) {
	r, err := newMockReader(m.path)
	r.pluginConf.Set("name", "reader2")
	return r, err
}

func TestRegisterReader(t *testing.T) {
	type args struct {
		maker ReaderMaker
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
				maker: &mockReaderMaker1{},
			},
			want: "github.com\\Breeze0806\\go-etl\\datax\\plugin\\reader\\resources\\plugin.json",
		},

		{
			name: "2",
			args: args{
				maker: &mockReaderMaker2{},
			},
			want: "github.com\\Breeze0806\\go-etl\\datax\\plugin\\reader\\resources\\plugin.json",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RegisterReader(tt.args.maker)
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
