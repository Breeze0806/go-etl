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

	"github.com/Breeze0806/go-etl/element"
	"github.com/Breeze0806/go-etl/storage/database"
)

func Test_tableParam_Query(t *testing.T) {
	type args struct {
		in0 []element.Record
	}
	tests := []struct {
		name    string
		t       *TableParam
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "1",
			t:    NewTableParam(&BaseConfig{}, &MockQuerier{}, nil),
			args: args{
				in0: nil,
			},
			wantErr: true,
		},
		{
			name: "2",
			t: NewTableParam(&BaseConfig{
				Column: []string{
					"f1", "f2", "f3",
				},
				Connection: ConnConfig{
					Table: TableConfig{
						Db:     "db",
						Schema: "schema",
						Name:   "table",
					},
				},
			}, &MockQuerier{}, nil),
			args: args{
				in0: nil,
			},
			want: "select f1,f2,f3 from db.schema.table where 1 = 2",
		},
		{
			name: "3",
			t: NewTableParam(&BaseConfig{
				Column: []string{
					"f1",
				},
				Connection: ConnConfig{
					Table: TableConfig{
						Db:     "db",
						Schema: "schema",
						Name:   "table",
					},
				},
			}, &MockQuerier{}, nil),
			args: args{
				in0: nil,
			},
			want: "select f1 from db.schema.table where 1 = 2",
		},
		{
			name: "4",
			t: NewTableParam(&BaseConfig{
				Column: []string{
					"*",
				},
				Connection: ConnConfig{
					Table: TableConfig{
						Db:     "db",
						Schema: "schema",
						Name:   "table",
					},
				},
			}, &MockQuerier{}, nil),
			args: args{
				in0: nil,
			},
			want: "select * from db.schema.table where 1 = 2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.t.Query(tt.args.in0)
			if (err != nil) != tt.wantErr {
				t.Errorf("tableParam.Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("tableParam.Query() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_tableParam_Agrs(t *testing.T) {
	type args struct {
		in0 []element.Record
	}
	tests := []struct {
		name    string
		t       *TableParam
		args    args
		want    []interface{}
		wantErr bool
	}{
		{
			name: "1",
			t:    NewTableParam(&BaseConfig{}, &MockQuerier{}, nil),
			args: args{
				in0: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.t.Agrs(tt.args.in0)
			if (err != nil) != tt.wantErr {
				t.Errorf("tableParam.Agrs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("tableParam.Agrs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_queryParam_Query(t *testing.T) {
	type args struct {
		in0 []element.Record
	}
	tests := []struct {
		name    string
		t       *MockTable
		config  *BaseConfig
		args    args
		want    string
		wantErr bool
	}{
		{
			name:   "1",
			t:      NewMockTable(database.NewBaseTable("db", "schema", "table")),
			config: &BaseConfig{},
			args: args{
				in0: nil,
			},
			wantErr: true,
		},
		{
			name: "2",
			t:    NewMockTable(database.NewBaseTable("db", "schema", "table")),
			config: &BaseConfig{
				Column: []string{
					"f1", "f2", "f3",
				},
			},
			args: args{
				in0: nil,
			},
			want: "select f1,f2,f3 from db.schema.table",
		},
		{
			name: "3",
			t:    NewMockTable(database.NewBaseTable("db", "schema", "table")),
			config: &BaseConfig{
				Column: []string{
					"f1",
				},
			},
			args: args{
				in0: nil,
			},
			want: "select f1 from db.schema.table",
		},
		{
			name: "4",
			t:    NewMockTable(database.NewBaseTable("db", "schema", "table")),
			config: &BaseConfig{
				Column: []string{
					"f1",
				},
				Connection: ConnConfig{
					Table: TableConfig{
						Db:   "db",
						Name: "table",
					},
				},
				Where: "a <> 1",
			},
			args: args{
				in0: nil,
			},
			want: "select f1 from db.schema.table where a <> 1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i, v := range tt.config.Column {
				tt.t.AddField(database.NewBaseField(i, v, NewMockFieldType(database.GoTypeBool)))
			}
			q := NewQueryParam(tt.config, tt.t, nil)
			got, err := q.Query(tt.args.in0)
			if (err != nil) != tt.wantErr {
				t.Errorf("queryParam.Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("queryParam.Query() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_queryParam_Agrs(t *testing.T) {
	type args struct {
		in0 []element.Record
	}
	tests := []struct {
		name    string
		q       *QueryParam
		args    args
		want    []interface{}
		wantErr bool
	}{
		{
			name: "1",
			q:    NewQueryParam(&BaseConfig{}, NewMockTable(database.NewBaseTable("db", "schema", "table")), nil),
			args: args{
				in0: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.q.Agrs(tt.args.in0)
			if (err != nil) != tt.wantErr {
				t.Errorf("queryParam.Agrs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("queryParam.Agrs() = %v, want %v", got, tt.want)
			}
		})
	}
}
