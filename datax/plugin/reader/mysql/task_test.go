package mysql

import (
	"context"
	"errors"
	"testing"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/datax/common/plugin"
	"github.com/Breeze0806/go-etl/datax/plugin/reader/rdbm"
	"github.com/Breeze0806/go-etl/element"
)

type mockSender struct {
	createErr error
	sendErr   error
}

func (m *mockSender) CreateRecord() (element.Record, error) {
	return element.NewDefaultRecord(), m.createErr
}

func (m *mockSender) SendWriter(record element.Record) error {
	return m.sendErr
}

func (m *mockSender) Flush() error {
	return nil
}

func (m *mockSender) Terminate() error {
	return nil
}

func (m *mockSender) Shutdown() error {
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
				BaseTask: plugin.NewBaseTask(),
				newQuerier: func(name string, conf *config.JSON) (rdbm.Querier, error) {
					return &rdbm.MockQuerier{}, nil
				},
			},
			args: args{
				ctx: context.TODO(),
			},
			conf:    rdbm.TestJSONFromFile(_pluginConfig),
			jobConf: rdbm.TestJSONFromString(`{}`),
		},
		{
			name: "2",
			t: &Task{
				BaseTask: plugin.NewBaseTask(),
				newQuerier: func(name string, conf *config.JSON) (rdbm.Querier, error) {
					return &rdbm.MockQuerier{}, nil
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
				BaseTask: plugin.NewBaseTask(),
				newQuerier: func(name string, conf *config.JSON) (rdbm.Querier, error) {
					return &rdbm.MockQuerier{}, nil
				},
			},
			args: args{
				ctx: context.TODO(),
			},
			conf: rdbm.TestJSONFromFile(_pluginConfig),
			jobConf: rdbm.TestJSONFromString(`{
				"username": 1		
			}`),
			wantErr: true,
		},
		{
			name: "4",
			t: &Task{
				BaseTask: plugin.NewBaseTask(),
				newQuerier: func(name string, conf *config.JSON) (rdbm.Querier, error) {
					return nil, errors.New("mock error")
				},
			},
			args: args{
				ctx: context.TODO(),
			},
			conf:    rdbm.TestJSONFromFile(_pluginConfig),
			jobConf: rdbm.TestJSONFromString(`{}`),
			wantErr: true,
		},
		{
			name: "5",
			t: &Task{
				BaseTask: plugin.NewBaseTask(),
				newQuerier: func(name string, conf *config.JSON) (rdbm.Querier, error) {
					return &rdbm.MockQuerier{
						QueryErr: errors.New("mock error"),
					}, nil
				},
			},
			args: args{
				ctx: context.TODO(),
			},
			conf:    rdbm.TestJSONFromFile(_pluginConfig),
			jobConf: rdbm.TestJSONFromString(`{}`),
			wantErr: true,
		},
		{
			name: "6",
			t: &Task{
				BaseTask: plugin.NewBaseTask(),
				newQuerier: func(name string, conf *config.JSON) (rdbm.Querier, error) {
					return &rdbm.MockQuerier{
						FetchErr: errors.New("mock error"),
					}, nil
				},
			},
			args: args{
				ctx: context.TODO(),
			},
			conf:    rdbm.TestJSONFromFile(_pluginConfig),
			jobConf: rdbm.TestJSONFromString(`{}`),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.t.SetPluginConf(tt.conf)
			tt.t.SetPluginJobConf(tt.jobConf)
			if err := tt.t.Init(tt.args.ctx); (err != nil) != tt.wantErr {
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
				BaseTask: plugin.NewBaseTask(),
				querier:  &rdbm.MockQuerier{},
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

func TestTask_StartRead(t *testing.T) {
	type args struct {
		ctx    context.Context
		sender plugin.RecordSender
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
				BaseTask: plugin.NewBaseTask(),
				querier:  &rdbm.MockQuerier{},
			},
			args: args{
				ctx:    context.TODO(),
				sender: &mockSender{},
			},
			wantErr: false,
		},
		{
			name: "2",
			t: &Task{
				BaseTask: plugin.NewBaseTask(),
				querier:  &rdbm.MockQuerier{},
			},
			args: args{
				ctx: context.TODO(),
				sender: &mockSender{
					createErr: errors.New("mock error"),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.t.StartRead(tt.args.ctx, tt.args.sender); (err != nil) != tt.wantErr {
				t.Errorf("Task.StartRead() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
