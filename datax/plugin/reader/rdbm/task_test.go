package rdbm

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/datax/common/plugin"
)

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
				Handler: newMockDbHandler(func(name string, conf *config.JSON) (Querier, error) {
					return &MockQuerier{}, nil
				}),
			},
			args: args{
				ctx: context.TODO(),
			},
			conf:    TestJSONFromFile(filepath.Join("resources", "plugin.json")),
			jobConf: TestJSONFromString(`{}`),
		},
		{
			name: "2",
			t: &Task{
				BaseTask: plugin.NewBaseTask(),
				Handler: newMockDbHandler(func(name string, conf *config.JSON) (Querier, error) {
					return &MockQuerier{}, nil
				}),
			},
			args: args{
				ctx: context.TODO(),
			},
			conf:    TestJSONFromString(`{}`),
			jobConf: TestJSONFromString(`{}`),
			wantErr: true,
		},
		{
			name: "3",
			t: &Task{
				BaseTask: plugin.NewBaseTask(),
				Handler: newMockDbHandler(func(name string, conf *config.JSON) (Querier, error) {
					return &MockQuerier{}, nil
				}),
			},
			args: args{
				ctx: context.TODO(),
			},
			conf: TestJSONFromFile(filepath.Join("resources", "plugin.json")),
			jobConf: TestJSONFromString(`{
				"username": 1		
			}`),
			wantErr: true,
		},
		{
			name: "4",
			t: &Task{
				BaseTask: plugin.NewBaseTask(),
				Handler: newMockDbHandler(func(name string, conf *config.JSON) (Querier, error) {
					return nil, errors.New("mock error")
				}),
			},
			args: args{
				ctx: context.TODO(),
			},
			conf:    TestJSONFromFile(filepath.Join("resources", "plugin.json")),
			jobConf: TestJSONFromString(`{}`),
			wantErr: true,
		},
		{
			name: "5",
			t: &Task{
				BaseTask: plugin.NewBaseTask(),
				Handler: newMockDbHandler(func(name string, conf *config.JSON) (Querier, error) {
					return &MockQuerier{
						PingErr: errors.New("mock error"),
					}, nil
				}),
			},
			args: args{
				ctx: context.TODO(),
			},
			conf:    TestJSONFromFile(filepath.Join("resources", "plugin.json")),
			jobConf: TestJSONFromString(`{}`),
			wantErr: true,
		},
		{
			name: "6",
			t: &Task{
				BaseTask: plugin.NewBaseTask(),
				Handler: newMockDbHandler(func(name string, conf *config.JSON) (Querier, error) {
					return &MockQuerier{
						FetchErr: errors.New("mock error"),
					}, nil
				}),
			},
			args: args{
				ctx: context.TODO(),
			},
			conf:    TestJSONFromFile(filepath.Join("resources", "plugin.json")),
			jobConf: TestJSONFromString(`{}`),
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
				Querier: &MockQuerier{},
			},
			args: args{
				ctx: context.TODO(),
			},
		},
		{
			name: "2",
			t: &Task{
				Querier: nil,
			},
			args: args{
				ctx: context.TODO(),
			},
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

func TestStartRead(t *testing.T) {
	type args struct {
		ctx    context.Context
		reader BatchReader
		sender plugin.RecordSender
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				ctx: context.TODO(),
				reader: NewBaseBatchReader(&Task{
					BaseTask: plugin.NewBaseTask(),
					Querier:  &MockQuerier{},
					Config:   &BaseConfig{},
				}, "", nil),
				sender: &MockSender{},
			},
		},

		{
			name: "2",
			args: args{
				ctx: context.TODO(),
				reader: NewBaseBatchReader(&Task{
					BaseTask: plugin.NewBaseTask(),
					Querier:  &MockQuerier{},
					Config:   &BaseConfig{},
				}, "Tx", nil),
				sender: &MockSender{
					SendErr: errors.New("mock error"),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := StartRead(tt.args.ctx, tt.args.reader, tt.args.sender); (err != nil) != tt.wantErr {
				t.Errorf("StartRead() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
