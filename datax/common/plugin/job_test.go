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

package plugin

import (
	"reflect"
	"testing"

	"github.com/Breeze0806/go/encoding"
)

type mockJobCollector struct {
}

func (m *mockJobCollector) JSON() *encoding.JSON {
	return nil
}

func (m *mockJobCollector) JSONByKey(key string) *encoding.JSON {
	return nil
}

func TestBaseJob_SetCollector(t *testing.T) {
	type args struct {
		collector JobCollector
	}
	tests := []struct {
		name string
		b    *BaseJob
		args args
		want JobCollector
	}{
		{
			name: "1",
			b:    NewBaseJob(),
			args: args{
				collector: &mockJobCollector{},
			},
			want: &mockJobCollector{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.SetCollector(tt.args.collector)
			if !reflect.DeepEqual(tt.b.Collector(), tt.want) {
				t.Errorf("Collector() = %p want %p", tt.b.Collector(), tt.want)
			}
		})
	}
}

func TestBaseJob_SetJobID(t *testing.T) {
	type args struct {
		jobID int64
	}
	tests := []struct {
		name string
		b    *BaseJob
		args args
		want int64
	}{
		{
			name: "1",
			b:    &BaseJob{},
			args: args{
				jobID: 10,
			},
			want: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.SetJobID(tt.args.jobID)
			if tt.b.JobID() != tt.want {
				t.Errorf("JobID() = %v want %v", tt.b.JobID(), tt.want)
			}
		})
	}
}
