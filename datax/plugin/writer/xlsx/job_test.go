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

package xlsx

import (
	"context"
	"reflect"
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

func TestJob_Split(t *testing.T) {
	type args struct {
		ctx    context.Context
		number int
	}
	tests := []struct {
		name    string
		j       *Job
		jobConf *config.JSON

		args        args
		wantConfigs []*config.JSON
		wantErr     bool
	}{
		{
			name:    "1",
			j:       NewJob(),
			jobConf: testJSONFromString(`{"xlsxs":[{"path":"file1","sheets":["Sheet1","Sheet2"]},{"path":"file2","sheets":["Sheet1"]}],"batchSize":1000,"batchTimeout":"2s","column":[],"nullFormat":"(null)","hasHeader":true,"header":["1"],"sheetRow":2}`),
			args: args{
				ctx: context.TODO(),
			},
			wantConfigs: []*config.JSON{
				testJSONFromString(`{"path":"file1","content":{"batchSize":1000,"batchTimeout":"2s","column":[],"sheets":["Sheet1","Sheet2"],"nullFormat":"(null)","hasHeader":true,"header":["1"],"sheetRow":2}}`),
				testJSONFromString(`{"path":"file2","content":{"batchSize":1000,"batchTimeout":"2s","column":[],"sheets":["Sheet1"],"nullFormat":"(null)","hasHeader":true,"header":["1"],"sheetRow":2}}`),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.j.SetPluginJobConf(tt.jobConf)
			if err := tt.j.Init(tt.args.ctx); err != nil {
				t.Errorf("Job.Init() error = %v", err)
				return
			}
			gotConfigs, err := tt.j.Split(tt.args.ctx, tt.args.number)
			if (err != nil) != tt.wantErr {
				t.Errorf("Job.Split() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotConfigs, tt.wantConfigs) {
				t.Errorf("Job.Split() = %v, want %v", gotConfigs, tt.wantConfigs)
			}
		})
	}
}
