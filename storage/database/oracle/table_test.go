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

package oracle

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/Breeze0806/go-etl/element"
	"github.com/Breeze0806/go-etl/storage/database"
	"github.com/godror/godror"
)

func TestNewTable(t *testing.T) {
	type args struct {
		b *database.BaseTable
	}
	tests := []struct {
		name string
		args args
		want *Table
	}{
		{
			name: "1",
			args: args{
				b: database.NewBaseTable("db", "schema", "table"),
			},
			want: NewTable(database.NewBaseTable("db", "schema", "table")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTable(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTable_Quoted(t *testing.T) {
	tests := []struct {
		name string
		tr   *Table
		want string
	}{
		{
			name: "1",
			tr:   NewTable(database.NewBaseTable("db", "schema", "table")),
			want: `"schema"."table"`,
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
			want: `"schema"."table"`,
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

func TestTable_AddField(t *testing.T) {
	table := NewTable(database.NewBaseTable("db", "schema", "table"))
	type args struct {
		baseField *database.BaseField
	}
	tests := []struct {
		name string
		tr   *Table
		args args
		want []database.Field
	}{
		{
			name: "1",
			tr:   table,
			args: args{
				baseField: database.NewBaseField(0, "f1", database.NewBaseFieldType(&sql.ColumnType{})),
			},
			want: []database.Field{
				NewField(database.NewBaseField(0, "f1", database.NewBaseFieldType(&sql.ColumnType{}))),
			},
		},
		{
			name: "2",
			tr:   table,
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
			tr:   table,
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
		tt.tr.AddField(tt.args.baseField)
		if !reflect.DeepEqual(tt.tr.Fields(), tt.want) {
			t.Errorf("run %v Table.Fields() = %v want: %v", tt.name, tt.tr.Fields(), tt.want)
		}
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
			tr:   NewTable(database.NewBaseTable("db", "schema", "table")),
			args: args{
				mode:   WriteModeInsert,
				txOpts: nil,
			},
			want:  NewInsertParam(NewTable(database.NewBaseTable("db", "schema", "table")), nil),
			want1: true,
		},
		{
			name: "2",
			tr:   NewTable(database.NewBaseTable("db", "schema", "table")),
			args: args{
				mode:   "copyIn",
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

func TestInsertParam_Query(t *testing.T) {
	type args struct {
		records []element.Record
		columns [][]element.Column
		fields  []database.Field
		t       database.Table
	}
	tests := []struct {
		name      string
		args      args
		wantQuery string
		wantErr   bool
	}{
		{
			name: "1",
			args: args{
				records: []element.Record{
					element.NewDefaultRecord(),
					element.NewDefaultRecord(),
					element.NewDefaultRecord(),
				},
				columns: [][]element.Column{
					{
						element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(1), "f1", 0),
						element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(2), "f2", 0),
						element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(3), "f3", 0),
					},
					{
						element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(5), "f2", 0),
						element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(4), "f1", 0),
						element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(6), "f3", 0),
					},
					{
						element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(9), "f3", 0),
						element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(7), "f1", 0),
						element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(8), "f2", 0),
					},
				},
				fields: []database.Field{
					NewField(database.NewBaseField(1, "f1", newMockColumnType("VARCHAR2"))),
					NewField(database.NewBaseField(2, "f2", newMockColumnType("VARCHAR2"))),
					NewField(database.NewBaseField(3, "f3", newMockColumnType("VARCHAR2"))),
				},
				t: NewTable(database.NewBaseTable("db", "schema", "table")),
			},
			wantQuery: `insert into "schema"."table"("f1","f2","f3") values (:1,:2,:3)`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, v := range tt.args.fields {
				tt.args.t.(*Table).AppendField(v)
			}

			for i, r := range tt.args.records {
				for _, c := range tt.args.columns[i] {
					r.Add(c)
				}
			}

			insertParam := NewInsertParam(tt.args.t, nil)
			gotQuery, err := insertParam.Query(tt.args.records)
			if (err != nil) != tt.wantErr {
				t.Errorf("InsertParam.Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotQuery != tt.wantQuery {
				t.Errorf("InsertParam.Query() = %v, want %v", gotQuery, tt.wantQuery)
			}
		})
	}
}

func TestInsertParam_Agrs(t *testing.T) {
	type args struct {
		records []element.Record
		columns [][]element.Column
		fields  []database.Field
		t       database.Table
	}
	tests := []struct {
		name        string
		args        args
		wantValuers []interface{}
		wantErr     bool
	}{
		{
			name: "1",
			args: args{
				records: []element.Record{
					element.NewDefaultRecord(),
					element.NewDefaultRecord(),
					element.NewDefaultRecord(),
				},
				columns: [][]element.Column{
					{
						element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(1), "f1", 0),
						element.NewDefaultColumn(element.NewNilBigIntColumnValue(), "f2", 0),
						element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(3), "f3", 0),
					},
					{
						element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(5), "f2", 0),
						element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(4), "f1", 0),
						element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(6), "f3", 0),
					},
					{
						element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(9), "f3", 0),
						element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(7), "f1", 0),
						element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(8), "f2", 0),
					},
				},
				fields: []database.Field{
					NewField(database.NewBaseField(1, "f1", newMockColumnType("VARCHAR2"))),
					NewField(database.NewBaseField(2, "f2", newMockColumnType("LONG"))),
					NewField(database.NewBaseField(3, "f3", newMockColumnType("VARCHAR2"))),
				},
				t: NewTable(database.NewBaseTable("db", "schema", "table")),
			},
			wantValuers: []interface{}{
				[]string{"1", "5", "9"},
				[][]byte{nil, []byte("4"), []byte("7")},
				[]string{"3", "6", "8"},
			},
		},
		{
			name: "2",
			args: args{
				records: []element.Record{
					element.NewDefaultRecord(),
					element.NewDefaultRecord(),
					element.NewDefaultRecord(),
				},
				columns: [][]element.Column{
					{
						element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(1), "f1", 0),
						element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(2), "f2", 0),
						element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(3), "f3", 0),
					},
					{
						element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(5), "f2", 0),
						element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(4), "f1", 0),
					},
					{
						element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(9), "f3", 0),
						element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(7), "f1", 0),
						element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(8), "f2", 0),
					},
				},
				fields: []database.Field{
					NewField(database.NewBaseField(1, "f1", newMockColumnType("VARCHAR2"))),
					NewField(database.NewBaseField(2, "f2", newMockColumnType("VARCHAR2"))),
					NewField(database.NewBaseField(3, "f3", newMockColumnType("VARCHAR2"))),
				},
				t: NewTable(database.NewBaseTable("db", "schema", "table")),
			},
			wantErr: true,
		},
		{
			name: "3",
			args: args{
				records: []element.Record{
					element.NewDefaultRecord(),
					element.NewDefaultRecord(),
					element.NewDefaultRecord(),
				},
				columns: [][]element.Column{
					{
						element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(1), "f1", 0),
						element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(2), "f2", 0),
						element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(3), "f3", 0),
					},
					{
						element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(5), "f2", 0),
						element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(4), "f1", 0),
						element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(6), "f3", 0),
					},
					{
						element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(9), "f3", 0),
						element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(7), "f1", 0),
						element.NewDefaultColumn(element.NewStringColumnValue("we"), "f2", 0),
					},
				},
				fields: []database.Field{
					NewField(database.NewBaseField(1, "f1", newMockColumnType("DATE"))),
					NewField(database.NewBaseField(2, "f2", newMockColumnType("DATE"))),
					NewField(database.NewBaseField(3, "f3", newMockColumnType("BOOLEAN"))),
				},
				t: NewTable(database.NewBaseTable("db", "schema", "table")),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, v := range tt.args.fields {
				tt.args.t.(*Table).AppendField(v)
			}

			for i, r := range tt.args.records {
				for _, c := range tt.args.columns[i] {
					r.Add(c)
				}
			}

			insertParam := NewInsertParam(tt.args.t, nil)
			gotValuers, err := insertParam.Agrs(tt.args.records)
			if (err != nil) != tt.wantErr {
				t.Errorf("InsertParam.Agrs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotValuers, tt.wantValuers) {
				t.Errorf("InsertParam.Agrs() = %v, want %v", gotValuers, tt.wantValuers)
			}
		})
	}
}

func TestTable_ShouldRetry(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		tr   *Table
		args args
		want bool
	}{
		{
			name: "1",
			tr:   NewTable(database.NewBaseTable("db", "schema", "table")),
			args: args{
				err: nil,
			},
		},
		{
			name: "2",
			tr:   NewTable(database.NewBaseTable("db", "schema", "table")),
			args: args{
				err: &godror.OraErr{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.ShouldRetry(tt.args.err); got != tt.want {
				t.Errorf("Table.ShouldRetry() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTable_ShouldOneByOne(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		tr   *Table
		args args
		want bool
	}{
		{
			name: "1",
			tr:   NewTable(database.NewBaseTable("db", "schema", "table")),
			args: args{
				err: nil,
			},
		},
		{
			name: "2",
			tr:   NewTable(database.NewBaseTable("db", "schema", "table")),
			args: args{
				err: &godror.OraErr{},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.ShouldOneByOne(tt.args.err); got != tt.want {
				t.Errorf("Table.ShouldOneByOne() = %v, want %v", got, tt.want)
			}
		})
	}
}
