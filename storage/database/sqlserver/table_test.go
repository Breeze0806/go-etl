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

package sqlserver

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/Breeze0806/go-etl/storage/database"
)

func TestTable_Quoted(t *testing.T) {
	tests := []struct {
		name string
		tr   *Table
		want string
	}{
		{
			name: "1",
			tr:   NewTable(database.NewBaseTable("db", "schema", "table")),
			want: `[db].[schema].[table]`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.Quoted(); got != tt.want {
				t.Errorf("Table.Quoted() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTable_String(t *testing.T) {
	tests := []struct {
		name string
		tr   *Table
		want string
	}{
		{
			name: "1",
			tr:   NewTable(database.NewBaseTable("db", "schema", "table")),
			want: `[db].[schema].[table]`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.String(); got != tt.want {
				t.Errorf("Table.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTable_ExecParam(t *testing.T) {
	type args struct {
		mode   string
		txOpts *sql.TxOptions
	}
	tests := []struct {
		name  string
		tr    *Table
		args  args
		want  database.Parameter
		want1 bool
	}{
		{
			name: "1",
			tr:   NewTable(database.NewBaseTable("", "schema", "table")),
			args: args{
				mode:   "insert",
				txOpts: nil,
			},
			want:  nil,
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.tr.ExecParam(tt.args.mode, tt.args.txOpts)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Table.ExecParam() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Table.ExecParam() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestTable_AddField(t *testing.T) {
	table := NewTable(database.NewBaseTable("db", "schema", "table"))
	type args struct {
		baseField *database.BaseField
	}
	tests := []struct {
		name string
		t    *Table
		args args
		want []database.Field
	}{
		{
			name: "1",
			t:    table,
			args: args{
				baseField: database.NewBaseField(0, "f1", database.NewBaseFieldType(&sql.ColumnType{})),
			},
			want: []database.Field{
				NewField(database.NewBaseField(0, "f1", database.NewBaseFieldType(&sql.ColumnType{}))),
			},
		},
		{
			name: "2",
			t:    table,
			args: args{
				baseField: database.NewBaseField(1, "f2", database.NewBaseFieldType(&sql.ColumnType{})),
			},
			want: []database.Field{
				NewField(database.NewBaseField(0, "f1", database.NewBaseFieldType(&sql.ColumnType{}))),
				NewField(database.NewBaseField(1, "f2", database.NewBaseFieldType(&sql.ColumnType{}))),
			},
		},
		{
			name: "3",
			t:    table,
			args: args{
				baseField: database.NewBaseField(2, "f3", database.NewBaseFieldType(&sql.ColumnType{})),
			},
			want: []database.Field{
				NewField(database.NewBaseField(0, "f1", database.NewBaseFieldType(&sql.ColumnType{}))),
				NewField(database.NewBaseField(1, "f2", database.NewBaseFieldType(&sql.ColumnType{}))),
				NewField(database.NewBaseField(2, "f3", database.NewBaseFieldType(&sql.ColumnType{}))),
			},
		},
	}
	for _, tt := range tests {
		tt.t.AddField(tt.args.baseField)
		if !reflect.DeepEqual(tt.t.Fields(), tt.want) {
			t.Errorf("run %v Table.Fields() = %v want: %v", tt.name, tt.t.Fields(), tt.want)
		}
	}
}
