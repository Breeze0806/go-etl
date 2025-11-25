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

package parquet

import (
	"context"
	"testing"

	"github.com/Breeze0806/go-etl/config"
)

func testJSONFromString(json string) *config.JSON {
	conf, err := config.NewJSONFromString(json)
	if err != nil {
		panic(err)
	}
	return conf
}

func TestNewJob(t *testing.T) {
	tests := []struct {
		name string
		want *Job
	}{
		{
			name: "Create new job",
			want: &Job{
				Job: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewJob(); got == nil {
				t.Errorf("NewJob() = %v, want non-nil", got)
			}
		})
	}
}

func TestJob_Init(t *testing.T) {
	type fields struct {
		Job  *Job
		conf *Config
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Initialize with valid config",
			fields: fields{
				Job: NewJob(),
			},
			args: args{
				ctx: context.TODO(),
			},
			wantErr: false,
		},
		{
			name: "Initialize with invalid config",
			fields: fields{
				Job: NewJob(),
			},
			args: args{
				ctx: context.TODO(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := tt.fields.Job
			if !tt.wantErr {
				j.SetPluginJobConf(testJSONFromString(`{"path":["/tmp/test.parquet"],"column":[{"name":"id","type":"INT64"}]}`))
			} else {
				j.SetPluginJobConf(testJSONFromString(`{"path":"/tmp/test.parquet"}`))
			}

			if err := j.Init(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Job.Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestJob_Destroy(t *testing.T) {
	type fields struct {
		Job  *Job
		conf *Config
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Destroy job",
			fields: fields{
				Job: NewJob(),
			},
			args: args{
				ctx: context.TODO(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := tt.fields.Job
			if err := j.Destroy(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Job.Destroy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestJob_Split(t *testing.T) {
	type fields struct {
		Job  *Job
		conf *Config
	}
	type args struct {
		ctx    context.Context
		number int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*config.JSON
		wantErr bool
	}{
		{
			name: "Split with single path",
			fields: fields{
				Job: NewJob(),
			},
			args: args{
				ctx:    context.TODO(),
				number: 1,
			},
			want: []*config.JSON{
				testJSONFromString(`{"path":"/tmp/test.parquet","content":{"column":[{"name":"id","type":"INT64"}]}}`),
			},
			wantErr: false,
		},
		{
			name: "Split with multiple paths",
			fields: fields{
				Job: NewJob(),
			},
			args: args{
				ctx:    context.TODO(),
				number: 1,
			},
			want: []*config.JSON{
				testJSONFromString(`{"path":"/tmp/test1.parquet","content":{"column":[{"name":"id","type":"INT64"}]}}`),
				testJSONFromString(`{"path":"/tmp/test2.parquet","content":{"column":[{"name":"id","type":"INT64"}]}}`),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := tt.fields.Job

			if tt.name == "Split with single path" {
				j.SetPluginJobConf(testJSONFromString(`{"path":["/tmp/test.parquet"],"column":[{"name":"id","type":"INT64"}]}`))
				j.conf = &Config{
					Path: []string{"/tmp/test.parquet"},
				}
			} else if tt.name == "Split with multiple paths" {
				j.SetPluginJobConf(testJSONFromString(`{"path":["/tmp/test1.parquet","/tmp/test2.parquet"],"column":[{"name":"id","type":"INT64"}]}`))
				j.conf = &Config{
					Path: []string{"/tmp/test1.parquet", "/tmp/test2.parquet"},
				}
			}

			got, err := j.Split(tt.args.ctx, tt.args.number)
			if (err != nil) != tt.wantErr {
				t.Errorf("Job.Split() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(got) != len(tt.want) {
					t.Errorf("Job.Split() length = %v, want %v", len(got), len(tt.want))
				}

				for i := range got {
					if i < len(tt.want) {
						// Check if paths match
						gotPath, _ := got[i].GetString("path")
						wantPath, _ := tt.want[i].GetString("path")
						if gotPath != wantPath {
							t.Errorf("Job.Split()[%d] path = %v, want %v", i, gotPath, wantPath)
						}
					}
				}
			}
		})
	}
}
