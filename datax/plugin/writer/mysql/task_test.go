package mysql

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/datax/common/plugin"
	"github.com/Breeze0806/go-etl/datax/common/spi/writer"
	"github.com/Breeze0806/go-etl/datax/core/transport/exchange"
	"github.com/Breeze0806/go-etl/datax/plugin/writer/rdbm"
	"github.com/Breeze0806/go-etl/element"
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
func newMockReceiverWithoutWait(n int, err error) *MockReceiver {
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
		select {
		case <-m.ticker.C:
			return element.NewDefaultRecord(), nil
		}
	}
	return element.NewDefaultRecord(), nil
}

func (m *MockReceiver) Shutdown() error {
	m.ticker.Stop()
	return nil
}

func TestTask_Init(t *testing.T) {
	type args struct {
		ctx context.Context
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
			t: &Task{
				BaseTask: writer.NewBaseTask(),
				newExecer: func(name string, conf *config.JSON) (rdbm.Execer, error) {
					return &rdbm.MockExecer{}, nil
				},
			},
			args: args{
				ctx: context.TODO(),
			},
			conf:    rdbm.TestJSONFromFile(filepath.Join("resources", "plugin.json")),
			jobConf: rdbm.TestJSONFromString(`{}`),
		},
		{
			name: "2",
			t: &Task{
				BaseTask: writer.NewBaseTask(),
				newExecer: func(name string, conf *config.JSON) (rdbm.Execer, error) {
					return &rdbm.MockExecer{}, nil
				},
			},
			args: args{
				ctx: context.TODO(),
			},
			conf:    rdbm.TestJSONFromString(`{}`),
			jobConf: rdbm.TestJSONFromString(`{}`),
			wantErr: true,
		},
		{
			name: "3",
			t: &Task{
				BaseTask: writer.NewBaseTask(),
				newExecer: func(name string, conf *config.JSON) (rdbm.Execer, error) {
					return &rdbm.MockExecer{}, nil
				},
			},
			args: args{
				ctx: context.TODO(),
			},
			conf: rdbm.TestJSONFromFile(filepath.Join("resources", "plugin.json")),
			jobConf: rdbm.TestJSONFromString(`{
				"username": 1
			}`),
			wantErr: true,
		},
		{
			name: "4",
			t: &Task{
				BaseTask: writer.NewBaseTask(),
				newExecer: func(name string, conf *config.JSON) (rdbm.Execer, error) {
					return nil, errors.New("mock error")
				},
			},
			args: args{
				ctx: context.TODO(),
			},
			conf:    rdbm.TestJSONFromFile(filepath.Join("resources", "plugin.json")),
			jobConf: rdbm.TestJSONFromString(`{}`),
			wantErr: true,
		},
		{
			name: "5",
			t: &Task{
				BaseTask: writer.NewBaseTask(),
				newExecer: func(name string, conf *config.JSON) (rdbm.Execer, error) {
					return &rdbm.MockExecer{
						QueryErr: errors.New("mock error"),
					}, nil
				},
			},
			args: args{
				ctx: context.TODO(),
			},
			conf:    rdbm.TestJSONFromFile(filepath.Join("resources", "plugin.json")),
			jobConf: rdbm.TestJSONFromString(`{}`),
			wantErr: true,
		},
		{
			name: "6",
			t: &Task{
				BaseTask: writer.NewBaseTask(),
				newExecer: func(name string, conf *config.JSON) (rdbm.Execer, error) {
					return &rdbm.MockExecer{
						FetchErr: errors.New("mock error"),
					}, nil
				},
			},
			args: args{
				ctx: context.TODO(),
			},
			conf:    rdbm.TestJSONFromFile(filepath.Join("resources", "plugin.json")),
			jobConf: rdbm.TestJSONFromString(`{}`),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.t.SetPluginConf(tt.conf)
			tt.t.SetPluginJobConf(tt.jobConf)
			err := tt.t.Init(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Task.Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTask_Destroy(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		t       *Task
		args    args
		wantErr bool
	}{
		{
			name: "1",
			t: &Task{
				BaseTask: writer.NewBaseTask(),
				execer:   &rdbm.MockExecer{},
			},
			args: args{
				ctx: context.TODO(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.t.Destroy(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Task.Destroy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTask_StartWrite(t *testing.T) {
	type args struct {
		ctx      context.Context
		receiver plugin.RecordReceiver
	}
	tests := []struct {
		name    string
		t       *Task
		args    args
		wantErr bool
	}{
		{
			name: "1",
			t: &Task{
				BaseTask: writer.NewBaseTask(),
				execer:   &rdbm.MockExecer{},
				param:    newParameter(&paramConfig{}, &rdbm.MockExecer{}),
			},
			args: args{
				ctx:      context.TODO(),
				receiver: newMockReceiverWithoutWait(10000, exchange.ErrTerminate),
			},
		},
		{
			name: "2",
			t: &Task{
				BaseTask: writer.NewBaseTask(),
				execer:   &rdbm.MockExecer{},
				param:    newParameter(&paramConfig{}, &rdbm.MockExecer{}),
			},
			args: args{
				ctx:      context.TODO(),
				receiver: newMockReceiverWithoutWait(10000, errors.New("mock error")),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.t.StartWrite(tt.args.ctx, tt.args.receiver); (err != nil) != tt.wantErr {
				t.Errorf("Task.StartWrite() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
