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
	"reflect"
	"testing"
	"time"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/datax/common/plugin"
)

func newMockDbHandler(newExecer func(name string, conf *config.JSON) (Execer, error)) DbHandler {
	return NewBaseDbHandler(newExecer, nil)
}

func TestJob_Init(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		j       *Job
		args    args
		conf    *config.JSON
		jobConf *config.JSON
		wantErr bool
		want    *config.JSON
	}{
		{
			name: "1",
			j: NewJob(newMockDbHandler(func(name string, conf *config.JSON) (Execer, error) {
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
			j: NewJob(newMockDbHandler(func(name string, conf *config.JSON) (Execer, error) {
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
			j: NewJob(newMockDbHandler(func(name string, conf *config.JSON) (Execer, error) {
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
			j: NewJob(newMockDbHandler(func(name string, conf *config.JSON) (Execer, error) {
				return nil, errors.New("mock error")
			})),
			args: args{
				ctx: context.TODO(),
			},
			conf: testJSON(),
			jobConf: testJSONFromString(`{
			}`),
			wantErr: true,
		},
		{
			name: "5",
			j: NewJob(newMockDbHandler(func(name string, conf *config.JSON) (Execer, error) {
				return &MockExecer{
					PingErr: errors.New("mock error"),
				}, nil
			})),
			args: args{
				ctx: context.TODO(),
			},
			conf: testJSON(),
			jobConf: testJSONFromString(`{
			}`),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.j.SetPluginConf(tt.conf)
			tt.j.SetPluginJobConf(tt.jobConf)
			err := tt.j.Init(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Job.Init() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if !equalConfigJSON(tt.j.Execer.(*MockExecer).config, tt.want) {
				t.Fatalf("Execer.config = %v, want = %v", tt.j.Execer.(*MockExecer).config, tt.want)
				return
			}
		})
	}
}

func TestJob_Destroy(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		j       *Job
		args    args
		wantErr bool
	}{
		{
			name: "1",
			j: &Job{
				Execer: &MockExecer{},
			},
			args: args{
				ctx: context.TODO(),
			},
		},
		{
			name: "2",
			j:    &Job{},
			args: args{
				ctx: context.TODO(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.j.Destroy(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Job.Destroy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestJob_Split(t *testing.T) {
	type args struct {
		ctx    context.Context
		number int
	}
	tests := []struct {
		name    string
		j       *Job
		args    args
		jobConf *config.JSON
		want    []*config.JSON
		wantErr bool
	}{
		{
			name: "1",
			j: &Job{
				BaseJob: plugin.NewBaseJob(),
			},
			args: args{
				ctx:    context.TODO(),
				number: 1,
			},
			jobConf: testJSONFromString(`{}`),
			want: []*config.JSON{
				testJSONFromString(`{}`),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.j.SetPluginJobConf(tt.jobConf)
			got, err := tt.j.Split(tt.args.ctx, tt.args.number)
			if (err != nil) != tt.wantErr {
				t.Errorf("Job.Split() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Job.Split() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJob_Prepare(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		j       *Job
		args    args
		timeout time.Duration
		wantErr bool
	}{
		{
			name: "1",
			j: &Job{
				conf: &BaseConfig{
					PreSQL: []string{
						"delate",
						"drop",
						"create",
					},
				},
				Execer: &MockExecer{},
			},
			args: args{
				ctx: context.TODO(),
			},
		},
		{
			name: "2",
			j: &Job{
				conf: &BaseConfig{
					PreSQL: []string{
						"wait",
						"drop",
						"create",
					},
				},
				Execer: &MockExecer{},
			},
			args: args{
				ctx: context.TODO(),
			},
			timeout: 10 * time.Millisecond,
			wantErr: true,
		},
		{
			name: "3",
			j: &Job{
				conf: &BaseConfig{
					PreSQL: []string{
						"drop",
						"create",
					},
				},
				Execer: &MockExecer{
					ExecErr: errors.New("mock error"),
				},
			},
			args: args{
				ctx: context.TODO(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(tt.args.ctx)
			go func() {
				if tt.timeout != 0 {
					<-time.After(tt.timeout)
				}
				cancel()
			}()
			if err := tt.j.Prepare(ctx); (err != nil) != tt.wantErr {
				t.Errorf("Job.Prepare() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestJob_Post(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		j       *Job
		args    args
		timeout time.Duration
		wantErr bool
	}{
		{
			name: "1",
			j: &Job{
				conf: &BaseConfig{
					PostSQL: []string{
						"delate",
						"drop",
						"create",
					},
				},
				Execer: &MockExecer{},
			},
			args: args{
				ctx: context.TODO(),
			},
		},
		{
			name: "2",
			j: &Job{
				conf: &BaseConfig{
					PostSQL: []string{
						"wait",
						"drop",
						"create",
					},
				},
				Execer: &MockExecer{},
			},
			args: args{
				ctx: context.TODO(),
			},
			timeout: 10 * time.Millisecond,
			wantErr: true,
		},
		{
			name: "3",
			j: &Job{
				conf: &BaseConfig{
					PostSQL: []string{
						"drop",
						"create",
					},
				},
				Execer: &MockExecer{ExecErr: errors.New("mock error")},
			},
			args: args{
				ctx: context.TODO(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(tt.args.ctx)
			go func() {
				if tt.timeout != 0 {
					<-time.After(tt.timeout)
				}
				cancel()
			}()
			if err := tt.j.Post(ctx); (err != nil) != tt.wantErr {
				t.Errorf("Job.Post() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
