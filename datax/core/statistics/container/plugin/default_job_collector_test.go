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

	"github.com/Breeze0806/go-etl/datax/core/statistics/container"
)

type testStruct struct {
	Path string `json:"path"`
}

func TestDefaultJobCollector_JSON(t *testing.T) {
	m := container.NewMetrics()
	m.Set("test", testStruct{Path: "value"})
	tests := []struct {
		name string
		d    *DefaultJobCollector
		want string
	}{
		{
			name: "1",
			d:    NewDefaultJobCollector(m).(*DefaultJobCollector),
			want: `{"test":{"path":"value"}}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.JSON().String(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DefaultJobCollector.JSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultJobCollector_JSONByKey(t *testing.T) {
	m := container.NewMetrics()
	m.Set("test", testStruct{Path: "value"})
	type args struct {
		key string
	}
	tests := []struct {
		name string
		d    *DefaultJobCollector
		args args
		want string
	}{
		{
			name: "1",
			d:    NewDefaultJobCollector(m).(*DefaultJobCollector),
			args: args{
				key: "test",
			},
			want: `{"path":"value"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.JSONByKey(tt.args.key).String(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DefaultJobCollector.JSONByKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
