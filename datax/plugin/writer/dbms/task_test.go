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

package dbms

import (
	"context"
	"errors"
	"testing"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/datax/common/spi/writer"
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
		want    *config.JSON
	}{
		{
			name: "1",
			t: NewTask(newMockDbHandler(func(name string, conf *config.JSON) (Execer, error) {
				return &MockExecer{
					config: conf,
				}, nil
			})),
			args: args{
				ctx: context.TODO(),
			},
			conf: testJSON(),
			jobConf: testJSONFromString(`{
				"connection":{
					"url":"breeze0806.xxx"
				},
				"username":"breeze0806",
				"password":"breeze0806",
				"job":{
					"setting":{
						"pool":{
						  "maxOpenConns":8,
						  "maxIdleConns":8,
						  "connMaxIdleTime":"40m",
						  "connMaxLifetime":"40m"
						},
						"retry":{
						  "type":"ntimes",
						  "strategy":{
							"n":3,
							"wait":"1s"
						  },
						  "ignoreOneByOneError":true
						}
					}
				}
			}`),
			want: testJSONFromString(`{
				"url":"breeze0806.xxx",
				"username":"breeze0806",
				"password":"breeze0806",
				"pool":{
					"maxOpenConns":8,
					"maxIdleConns":8,
					"connMaxIdleTime":"40m",
					"connMaxLifetime":"40m"
				  },
				  "retry":{
					"type":"ntimes",
					"strategy":{
					  "n":3,
					  "wait":"1s"
					},
					"ignoreOneByOneError":true
				  }
			}`),
		},
		{
			name: "2",
			t: NewTask(newMockDbHandler(func(name string, conf *config.JSON) (Execer, error) {
				return &MockExecer{}, nil
			})),
			args: args{
				ctx: context.TODO(),
			},
			conf:    testJSONFromString(`{}`),
			jobConf: testJSONFromString(`{}`),
			wantErr: true,
		},
		{
			name: "3",
			t: NewTask(newMockDbHandler(func(name string, conf *config.JSON) (Execer, error) {
				return &MockExecer{}, nil
			})),
			args: args{
				ctx: context.TODO(),
			},
			conf: testJSON(),
			jobConf: testJSONFromString(`{
				"username": 1		
			}`),
			wantErr: true,
		},
		{
			name: "4",
			t: NewTask(newMockDbHandler(func(name string, conf *config.JSON) (Execer, error) {
				return nil, errors.New("mock error")
			})),
			args: args{
				ctx: context.TODO(),
			},
			conf:    testJSON(),
			jobConf: testJSONFromString(`{}`),
			wantErr: true,
		},
		{
			name: "5",
			t: NewTask(newMockDbHandler(func(name string, conf *config.JSON) (Execer, error) {
				return &MockExecer{
					PingErr: errors.New("mock error"),
				}, nil
			})),
			args: args{
				ctx: context.TODO(),
			},
			conf:    testJSON(),
			jobConf: testJSONFromString(`{}`),
			wantErr: true,
		},
		{
			name: "6",
			t: NewTask(newMockDbHandler(func(name string, conf *config.JSON) (Execer, error) {
				return &MockExecer{
					FetchErr: errors.New("mock error"),
				}, nil
			})),
			args: args{
				ctx: context.TODO(),
			},
			conf:    testJSON(),
			jobConf: testJSONFromString(`{}`),
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
				return
			}
			if err != nil {
				return
			}
			if !equalConfigJSON(tt.t.Execer.(*MockExecer).config, tt.want) {
				t.Fatalf("Execer.config = %v, want = %v", tt.t.Execer.(*MockExecer).config, tt.want)
				return
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
				Execer:   &MockExecer{},
			},
			args: args{
				ctx: context.TODO(),
			},
		},
		{
			name: "2",
			t:    NewTask(nil),
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
