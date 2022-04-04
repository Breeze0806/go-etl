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

package file

import (
	"reflect"
	"testing"
	"time"

	"github.com/Breeze0806/go-etl/config"
)

func testJSONFromString(json string) *config.JSON {
	conf, err := config.NewJSONFromString(json)
	if err != nil {
		panic(err)
	}
	return conf
}

func testBaseConfig(json string) *BaseConfig {
	conf, err := NewBaseConfig(testJSONFromString(json))
	if err != nil {
		panic(err)
	}
	return conf
}

func TestBaseConfig_GetBatchTimeout(t *testing.T) {
	tests := []struct {
		name string
		b    *BaseConfig
		want time.Duration
	}{
		{
			name: "1",
			b:    testBaseConfig(`{}`),
			want: defalutBatchTimeout,
		},
		{
			name: "2",
			b:    testBaseConfig(`{"batchTimeout":"2s"}`),
			want: 2 * time.Second,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.GetBatchTimeout(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BaseConfig.GetBatchTimeout() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseConfig_GetBatchSize(t *testing.T) {
	tests := []struct {
		name string
		b    *BaseConfig
		want int
	}{
		{
			name: "1",
			b:    testBaseConfig(`{}`),
			want: defalutBatchSize,
		},
		{
			name: "2",
			b:    testBaseConfig(`{"batchSize":30000}`),
			want: 30000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.GetBatchSize(); got != tt.want {
				t.Errorf("BaseConfig.GetBatchSize() = %v, want %v", got, tt.want)
			}
		})
	}
}
