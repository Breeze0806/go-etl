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

package rdbm

import (
	"reflect"
	"testing"
	"time"

	"github.com/Breeze0806/go-etl/config"
	rdbmreader "github.com/Breeze0806/go-etl/datax/plugin/reader/rdbm"
)

func testBaseConfig(conf *config.JSON) (bc *BaseConfig) {
	var err error
	bc, err = NewBaseConfig(conf)
	if err != nil {
		panic(err)
	}
	return bc
}

func TestBaseConfig_GetColumns(t *testing.T) {
	tests := []struct {
		name        string
		b           *BaseConfig
		wantColumns []rdbmreader.Column
	}{
		{
			name: "1",
			b: &BaseConfig{
				Column: []string{"f1", "f2", "f3", "f4"},
			},
			wantColumns: []rdbmreader.Column{
				&rdbmreader.BaseColumn{
					Name: "f1",
				},
				&rdbmreader.BaseColumn{
					Name: "f2",
				},
				&rdbmreader.BaseColumn{
					Name: "f3",
				},
				&rdbmreader.BaseColumn{
					Name: "f4",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotColumns := tt.b.GetColumns(); !reflect.DeepEqual(gotColumns, tt.wantColumns) {
				t.Errorf("BaseConfig.GetColumns() = %v, want %v", gotColumns, tt.wantColumns)
			}
		})
	}
}

func TestBaseConfig_GetBatchTimeout(t *testing.T) {
	tests := []struct {
		name string
		b    *BaseConfig
		want time.Duration
	}{
		{
			name: "1",
			b:    testBaseConfig(testJSONFromString("{}")),
			want: defalutBatchTimeout,
		},
		{
			name: "2",
			b:    testBaseConfig(testJSONFromString(`{"batchTimeout":"100ms"}`)),
			want: 100 * time.Millisecond,
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
			b:    testBaseConfig(testJSONFromString("{}")),
			want: defalutBatchSize,
		},

		{
			name: "2",
			b:    testBaseConfig(testJSONFromString(`{"batchSize":30000}`)),
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
