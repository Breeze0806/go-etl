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
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/Breeze0806/go-etl/config"
	dbmsreader "github.com/Breeze0806/go-etl/datax/plugin/reader/dbms"
	"github.com/Breeze0806/go-etl/schedule"
)

func testBaseConfig(conf *config.JSON) (bc *BaseConfig) {
	var err error
	bc, err = NewBaseConfig(conf)
	if err != nil {
		panic(err)
	}
	return bc
}

type mockNTimeTask struct {
	err error
	n   int
}

func (m *mockNTimeTask) Do() error {
	m.n--
	if m.n == 0 {
		return m.err
	}
	return errors.New("mock error")
}

type mockRetryJudger struct{}

func (m *mockRetryJudger) ShouldRetry(err error) bool {
	return err != nil
}

func TestBaseConfig_GetColumns(t *testing.T) {
	tests := []struct {
		name        string
		b           *BaseConfig
		wantColumns []dbmsreader.Column
	}{
		{
			name: "1",
			b: &BaseConfig{
				Column: []string{"f1", "f2", "f3", "f4"},
			},
			wantColumns: []dbmsreader.Column{
				&dbmsreader.BaseColumn{
					Name: "f1",
				},
				&dbmsreader.BaseColumn{
					Name: "f2",
				},
				&dbmsreader.BaseColumn{
					Name: "f3",
				},
				&dbmsreader.BaseColumn{
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

func TestBaseConfig_GetRetryStrategy(t *testing.T) {
	type args struct {
		j schedule.RetryJudger
	}
	tests := []struct {
		name    string
		b       *BaseConfig
		args    args
		want    schedule.RetryStrategy
		wantErr bool
	}{
		{
			name: "1",
			b:    testBaseConfig(testJSONFromString(`{"retry":{"type":"ntimes","strategy":{"wait":"1s","n":3}}}`)),
			args: args{
				j: &mockRetryJudger{},
			},
			want: schedule.NewNTimesRetryStrategy(&mockRetryJudger{}, 3, 1*time.Second),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.GetRetryStrategy(tt.args.j)
			if (err != nil) != tt.wantErr {
				t.Errorf("BaseConfig.GetRetryStrategy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BaseConfig.GetRetryStrategy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseConfig_IgnoreOneByOneError(t *testing.T) {
	tests := []struct {
		name string
		b    *BaseConfig
		want bool
	}{
		{
			name: "1",
			b:    testBaseConfig(testJSONFromString(`{"retry":{"ignoreOneByOneError":true}}`)),
			want: true,
		},
		{
			name: "2",
			b:    testBaseConfig(testJSONFromString(`{"retry":{"ignoreOneByOneError":"true"}}`)),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.IgnoreOneByOneError(); got != tt.want {
				t.Errorf("BaseConfig.IgnoreOneByOneError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseConfig_GetPreSQL(t *testing.T) {
	tests := []struct {
		name string
		b    *BaseConfig
		want []string
	}{
		{
			name: "1",
			b: &BaseConfig{
				PreSQL: []string{"delete from a", "create table a"},
			},
			want: []string{"delete from a", "create table a"},
		},
		{
			name: "2",
			b: &BaseConfig{
				PreSQL: []string{"", "delete from a", "", "create table a", ""},
			},
			want: []string{"delete from a", "create table a"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.GetPreSQL(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BaseConfig.GetPreSQL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseConfig_GetPostSQL(t *testing.T) {
	tests := []struct {
		name string
		b    *BaseConfig
		want []string
	}{
		{
			name: "1",
			b: &BaseConfig{
				PostSQL: []string{"delete from a", "create table a"},
			},
			want: []string{"delete from a", "create table a"},
		},
		{
			name: "2",
			b: &BaseConfig{
				PostSQL: []string{"", "delete from a", "", "create table a", ""},
			},
			want: []string{"delete from a", "create table a"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.GetPostSQL(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BaseConfig.GetPostSQL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewBaseConfig(t *testing.T) {
	type args struct {
		conf *config.JSON
	}
	tests := []struct {
		name    string
		args    args
		wantC   *BaseConfig
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				conf: testJSONFromString(`{"preSQL":["select * from a"]}`),
			},
			wantErr: true,
		},
		{
			name: "2",
			args: args{
				conf: testJSONFromString(`{"preSQL":[" select * from a"]}`),
			},
			wantErr: true,
		},
		{
			name: "3",
			args: args{
				conf: testJSONFromString(`{"preSQL":[" SELECT * from a"]}`),
			},
			wantErr: true,
		},
		{
			name: "4",
			args: args{
				conf: testJSONFromString(`{"postSQL":["select * from a"]}`),
			},
			wantErr: true,
		},
		{
			name: "5",
			args: args{
				conf: testJSONFromString(`{"postSQL":[" select * from a"]}`),
			},
			wantErr: true,
		},
		{
			name: "6",
			args: args{
				conf: testJSONFromString(`{"postSQL":[" SELECT * from a"]}`),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotC, err := NewBaseConfig(tt.args.conf)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewBaseConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotC, tt.wantC) {
				t.Errorf("NewBaseConfig() = %v, want %v", gotC, tt.wantC)
			}
		})
	}
}
