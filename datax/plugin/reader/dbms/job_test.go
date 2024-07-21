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

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/datax/common/plugin"
)

func newMockDbHandler(newQuerier func(name string, conf *config.JSON) (Querier, error)) DbHandler {
	return NewBaseDbHandler(newQuerier, nil)
}

func TestJob_Init(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		j       *Job
		conf    *config.JSON
		jobConf *config.JSON
		args    args
		wantErr bool
		want    *config.JSON
	}{
		{
			name: "1",
			j: NewJob(newMockDbHandler(func(name string, conf *config.JSON) (Querier, error) {
				return &MockQuerier{
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
			j: NewJob(newMockDbHandler(func(name string, conf *config.JSON) (Querier, error) {
				return &MockQuerier{}, nil
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
			j: NewJob(newMockDbHandler(func(name string, conf *config.JSON) (Querier, error) {
				return &MockQuerier{}, nil
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
			j: NewJob(newMockDbHandler(func(name string, conf *config.JSON) (Querier, error) {
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
			j: NewJob(newMockDbHandler(func(name string, conf *config.JSON) (Querier, error) {
				return &MockQuerier{
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
			if !equalConfigJSON(tt.j.Querier.(*MockQuerier).config, tt.want) {
				t.Fatalf("Querier.config = %v, want = %v", tt.j.Querier.(*MockQuerier).config, tt.want)
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
				Querier: &MockQuerier{},
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
				Config:  &BaseConfig{},
				Querier: &MockQuerier{},
				handler: newMockDbHandler(func(name string, conf *config.JSON) (Querier, error) {
					return &MockQuerier{}, nil
				}),
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
		{
			name: "2",
			j: &Job{
				BaseJob: plugin.NewBaseJob(),
				Config: &BaseConfig{
					Split: SplitConfig{
						Key: "f1",
					},
				},
				Querier: &MockQuerier{},
				handler: newMockDbHandler(func(name string, conf *config.JSON) (Querier, error) {
					return &MockQuerier{}, nil
				}),
			},
			args: args{
				ctx:    context.TODO(),
				number: 2,
			},
			jobConf: testJSONFromString(`{}`),
			want: []*config.JSON{
				testJSONFromString(`{"split":{"range":{"type":"bigInt","layout":"","left":"10000","right":"20000"}},"where":"f1 >= $1 and f1 < $2"}`),
				testJSONFromString(`{"split":{"range":{"type":"bigInt","layout":"","left":"20000","right":"30000"}},"where":"f1 >= $1 and f1 <= $2"}`),
			},
		},
		{
			name: "3",
			j: &Job{
				BaseJob: plugin.NewBaseJob(),
				Config: &BaseConfig{
					Where: "a < 1",
					Split: SplitConfig{
						Key: "f1",
					},
				},
				Querier: &MockQuerier{},
				handler: newMockDbHandler(func(name string, conf *config.JSON) (Querier, error) {
					return &MockQuerier{}, nil
				}),
			},
			args: args{
				ctx:    context.TODO(),
				number: 2,
			},
			jobConf: testJSONFromString(`{"where":"a < 1"}`),
			want: []*config.JSON{
				testJSONFromString(`{"where":"(a < 1) and (f1 >= $1 and f1 < $2)","split":{"range":{"type":"bigInt","layout":"","left":"10000","right":"20000"}}}`),
				testJSONFromString(`{"where":"(a < 1) and (f1 >= $1 and f1 <= $2)","split":{"range":{"type":"bigInt","layout":"","left":"20000","right":"30000"}}}`),
			},
		},
		{
			name: "4",
			j: &Job{
				BaseJob: plugin.NewBaseJob(),
				Config: &BaseConfig{

					Split: SplitConfig{
						Key: "f1",
					},
				},
				Querier: &MockQuerier{
					FetchErr: errors.New("mock error"),
				},
				handler: newMockDbHandler(func(name string, conf *config.JSON) (Querier, error) {
					return &MockQuerier{}, nil
				}),
			},
			args: args{
				ctx:    context.TODO(),
				number: 2,
			},
			jobConf: testJSONFromString(`{}`),
			wantErr: true,
		},
		{
			name: "5",
			j: &Job{
				BaseJob: plugin.NewBaseJob(),
				Config: &BaseConfig{
					Split: SplitConfig{
						Key: "f1",
					},
				},
				Querier: &MockQuerier{
					FetchMinErr: errors.New("mock error"),
				},
				handler: newMockDbHandler(func(name string, conf *config.JSON) (Querier, error) {
					return &MockQuerier{}, nil
				}),
			},
			args: args{
				ctx:    context.TODO(),
				number: 2,
			},
			jobConf: testJSONFromString(`{"where":"a < 1"}`),
			wantErr: true,
		},
		{
			name: "6",
			j: &Job{
				BaseJob: plugin.NewBaseJob(),
				Config: &BaseConfig{
					Split: SplitConfig{
						Key: "f1",
					},
				},
				Querier: &MockQuerier{
					FetchMaxErr: errors.New("mock error"),
				},
				handler: newMockDbHandler(func(name string, conf *config.JSON) (Querier, error) {
					return &MockQuerier{}, nil
				}),
			},
			args: args{
				ctx:    context.TODO(),
				number: 2,
			},
			jobConf: testJSONFromString(`{"where":"a < 1"}`),
			wantErr: true,
		},
		{
			name: "7",
			j: &Job{
				BaseJob: plugin.NewBaseJob(),
				Config: &BaseConfig{
					Split: SplitConfig{
						Key: "f1",
					},
				},
				Querier: &MockQuerier{},
				handler: newMockDbHandler(func(name string, conf *config.JSON) (Querier, error) {
					return &MockQuerier{}, nil
				}),
			},
			args: args{
				ctx:    context.TODO(),
				number: 0,
			},
			jobConf: testJSONFromString(`{}`),
			wantErr: true,
		},
		{
			name: "8",
			j: &Job{
				BaseJob: plugin.NewBaseJob(),
				Config:  &BaseConfig{},
				Querier: &MockQuerier{},
				handler: newMockDbHandler(func(name string, conf *config.JSON) (Querier, error) {
					return &MockQuerier{}, nil
				}),
			},
			args: args{
				ctx:    context.TODO(),
				number: 2,
			},
			jobConf: testJSONFromString(`{}`),
			want: []*config.JSON{
				testJSONFromString(`{}`),
			},
		},
		{
			name: "9",
			j: &Job{
				BaseJob: plugin.NewBaseJob(),
				Config: &BaseConfig{
					Split: SplitConfig{
						Key:          "f1",
						TimeAccuracy: "day",
						Range: SplitRange{
							Type:  "time",
							Left:  "2023-05-01",
							Right: "2023-05-03",
						},
					},
				},
				Querier: &MockQuerier{isTime: true},
				handler: newMockDbHandler(func(name string, conf *config.JSON) (Querier, error) {
					return &MockQuerier{isTime: true}, nil
				}),
			},
			args: args{
				ctx:    context.TODO(),
				number: 2,
			},
			jobConf: testJSONFromString(`{}`),
			want: []*config.JSON{
				testJSONFromString(`{"split":{"range":{"type":"time","layout":"2006-01-02","left":"2023-05-01","right":"2023-05-02"}},"where":"f1 >= $1 and f1 < $2"}`),
				testJSONFromString(`{"split":{"range":{"type":"time","layout":"2006-01-02","left":"2023-05-02","right":"2023-05-03"}},"where":"f1 >= $1 and f1 <= $2"}`),
			},
		},
		{
			name: "10",
			j: &Job{
				BaseJob: plugin.NewBaseJob(),
				Config: &BaseConfig{
					QuerySQL: []string{
						"select a,b from table_a join table_b on table_a.id = table_b.id",
						"select a,b,c from table_a join table_b on table_a.id = table_b.id",
					},
				},
				Querier: &MockQuerier{isTime: true},
				handler: newMockDbHandler(func(name string, conf *config.JSON) (Querier, error) {
					return &MockQuerier{isTime: true}, nil
				}),
			},
			args: args{
				ctx:    context.TODO(),
				number: 2,
			},
			jobConf: testJSONFromString(`{}`),
			want: []*config.JSON{
				testJSONFromString(`{"querySql":["select a,b from table_a join table_b on table_a.id = table_b.id"]}`),
				testJSONFromString(`{"querySql":["select a,b,c from table_a join table_b on table_a.id = table_b.id"]}`),
			},
		},
		{
			name: "11",
			j: &Job{
				BaseJob: plugin.NewBaseJob(),
				Config: &BaseConfig{
					Where: "a < 1",
					Split: SplitConfig{
						Key: "f1",
					},
				},
				Querier: &MockQuerier{},
				handler: newMockDbHandler(func(name string, conf *config.JSON) (Querier, error) {
					return &MockQuerier{}, nil
				}),
			},
			args: args{
				ctx:    context.TODO(),
				number: 1,
			},
			jobConf: testJSONFromString(`{"where":"a < 1"}`),
			want: []*config.JSON{
				testJSONFromString(`{"where":"(a < 1) and (f1 >= $1 and f1 <= $2)","split":{"range":{"type":"bigInt","layout":"","left":"10000","right":"30000"}}}`),
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
