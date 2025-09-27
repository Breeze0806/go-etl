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

package dm

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/Breeze0806/go-etl/storage/database"
)

func TestTable_Quoted(t *testing.T) {
	type args struct {
		schema string
		name   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1",
			args: args{
				schema: "schema",
				name:   "table",
			},
			want: `"schema"."table"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table := NewTable(database.NewBaseTable(tt.args.schema, "", tt.args.name))
			if got := table.Quoted(); got != tt.want {
				t.Errorf("Table.Quoted() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTable_String(t *testing.T) {
	table := NewTable(database.NewBaseTable("schema", "", "table"))
	if got, want := table.String(), `"schema"."table"`; got != want {
		t.Errorf("Table.String() = %v, want %v", got, want)
	}
}

func TestTable_AddField(t *testing.T) {
	type args struct {
		baseField *database.BaseField
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "1",
			args: args{
				baseField: database.NewBaseField(0, "f", database.NewBaseFieldType(&sql.ColumnType{})),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table := NewTable(database.NewBaseTable("schema", "", "table"))
			table.AddField(tt.args.baseField)
		})
	}
}

func TestTable_ExecParam(t *testing.T) {
	type args struct {
		mode   string
		txOpts *sql.TxOptions
	}
	tests := []struct {
		name   string
		args   args
		want   database.Parameter
		wantOk bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table := NewTable(database.NewBaseTable("schema", "", "table"))
			got, gotOk := table.ExecParam(tt.args.mode, tt.args.txOpts)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Table.ExecParam() got = %v, want %v", got, tt.want)
			}
			if gotOk != tt.wantOk {
				t.Errorf("Table.ExecParam() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}
