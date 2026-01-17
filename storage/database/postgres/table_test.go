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

package postgres

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"net"
	"reflect"
	"testing"

	"github.com/Breeze0806/go-etl/element"
	"github.com/Breeze0806/go-etl/storage/database"
	"github.com/lib/pq"
	"github.com/lib/pq/oid"
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

func TestCopyInParam_Query(t *testing.T) {
	type input struct {
		t      *Table
		fields []database.Field
		txOps  *sql.TxOptions
	}

	type args struct {
		in0 []element.Record
	}
	tests := []struct {
		name      string
		input     input
		args      args
		wantQuery string
		wantErr   bool
	}{
		{
			name: "1",
			input: input{
				t: NewTable(database.NewBaseTable("db", "schema", "table")),
				fields: []database.Field{
					NewField(database.NewBaseField(0,
						"f1", NewFieldType(newMockColumnType(oid.TypeName[oid.T_int8])))),
					NewField(database.NewBaseField(0,
						"f2", NewFieldType(newMockColumnType(oid.TypeName[oid.T_numeric])))),
					NewField(database.NewBaseField(0,
						"f3", NewFieldType(newMockColumnType(oid.TypeName[oid.T_varchar])))),
				},
				txOps: nil,
			},

			args: args{
				in0: nil,
			},

			wantQuery: pq.CopyInSchema("schema", "table", "f1", "f2", "f3"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, v := range tt.input.fields {
				tt.input.t.AppendField(v)
			}
			ci := NewCopyInParam(tt.input.t, tt.input.txOps)
			gotQuery, err := ci.Query(tt.args.in0)
			if (err != nil) != tt.wantErr {
				t.Errorf("CopyInParam.Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotQuery != tt.wantQuery {
				t.Errorf("CopyInParam.Query() = %v, want %v", gotQuery, tt.wantQuery)
			}
		})
	}
}

func TestCopyInParam_Agrs(t *testing.T) {
	type input struct {
		t      *Table
		fields []*database.BaseField
		txOps  *sql.TxOptions
	}

	type args struct {
		records []element.Record
		columns [][]element.Column
	}
	tests := []struct {
		name        string
		input       input
		args        args
		wantValuers []any
		wantErr     bool
	}{
		{
			name: "1",
			input: input{
				t: NewTable(database.NewBaseTable("db", "schema", "table")),
				fields: []*database.BaseField{
					database.NewBaseField(0,
						"f1", NewFieldType(newMockColumnType(oid.TypeName[oid.T_int8]))),
					database.NewBaseField(0,
						"f2", NewFieldType(newMockColumnType(oid.TypeName[oid.T_numeric]))),
					database.NewBaseField(0,
						"f3", NewFieldType(newMockColumnType(oid.TypeName[oid.T_varchar]))),
				},
				txOps: nil,
			},

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
			},
			wantValuers: []any{
				int64(1), "2", "3",
				int64(5), "4", "6",
				int64(9), "7", "8",
			},
		},

		{
			name: "2",
			input: input{
				t: NewTable(database.NewBaseTable("db", "schema", "table")),
				fields: []*database.BaseField{
					database.NewBaseField(0,
						"f1", NewFieldType(newMockColumnType(oid.TypeName[oid.T_int8]))),
					database.NewBaseField(0,
						"f2", NewFieldType(newMockColumnType(oid.TypeName[oid.T_numeric]))),
					database.NewBaseField(0,
						"f3", NewFieldType(newMockColumnType(oid.TypeName[oid.T_varchar]))),
				},
				txOps: nil,
			},

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
						element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(7), "f4", 0),
						element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(8), "f2", 0),
					},
				},
			},
			wantErr: true,
		},

		{
			name: "3",
			input: input{
				t: NewTable(database.NewBaseTable("db", "schema", "table")),
				fields: []*database.BaseField{
					database.NewBaseField(0,
						"f1", NewFieldType(newMockColumnType(oid.TypeName[oid.T__bool]))),
					database.NewBaseField(0,
						"f2", NewFieldType(newMockColumnType(oid.TypeName[oid.T_numeric]))),
					database.NewBaseField(0,
						"f3", NewFieldType(newMockColumnType(oid.TypeName[oid.T_varchar]))),
				},
				txOps: nil,
			},

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
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, v := range tt.input.fields {
				tt.input.t.AddField(v)
			}

			for i, v := range tt.args.records {
				for _, c := range tt.args.columns[i] {
					v.Add(c)
				}
			}
			ci := NewCopyInParam(tt.input.t, tt.input.txOps)
			gotValuers, err := ci.Agrs(tt.args.records)
			if (err != nil) != tt.wantErr {
				t.Errorf("CopyInParam.Agrs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(gotValuers, tt.wantValuers) {
				t.Errorf("CopyInParam.Agrs() = %v, want %v", gotValuers, tt.wantValuers)
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
			name: WriteModeCopyIn,
			tr:   NewTable(database.NewBaseTable("db", "schema", "table")),
			args: args{
				mode:   WriteModeCopyIn,
				txOpts: nil,
			},
			want:  NewCopyInParam(NewTable(database.NewBaseTable("db", "schema", "table")), nil),
			want1: true,
		},

		{
			name: "INSERT",
			tr:   NewTable(database.NewBaseTable("db", "schema", "table")),
			args: args{
				mode:   "INSERT",
				txOpts: nil,
			},
			want:  nil,
			want1: false,
		},
		{
			name: WriteModeUpsert,
			tr:   NewTable(database.NewBaseTable("db", "schema", "table")),
			args: args{
				mode:   WriteModeUpsert,
				txOpts: nil,
			},
			want:  NewUpsetParam(NewTable(database.NewBaseTable("db", "schema", "table")), nil),
			want1: true,
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
				err: &net.AddrError{},
			},
			want: true,
		},
		{
			name: "3",
			tr:   NewTable(database.NewBaseTable("db", "schema", "table")),
			args: args{
				err: driver.ErrBadConn,
			},
			want: true,
		},
		{
			name: "4",
			tr:   NewTable(database.NewBaseTable("db", "schema", "table")),
			args: args{
				err: &pq.Error{},
			},
			want: false,
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
				err: &net.AddrError{},
			},
		},
		{
			name: "3",
			tr:   NewTable(database.NewBaseTable("db", "schema", "table")),
			args: args{
				err: driver.ErrBadConn,
			},
		},
		{
			name: "4",
			tr:   NewTable(database.NewBaseTable("db", "schema", "table")),
			args: args{
				err: &pq.Error{},
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

func TestUpsetParam_Query(t *testing.T) {
	type input struct {
		t      *Table
		fields []database.Field
		txOpts *sql.TxOptions
		config map[string]interface{} // 配置upsertSQL
	}

	type args struct {
		records []element.Record
	}

	tests := []struct {
		name      string
		input     input
		args      args
		wantQuery string
		wantErr   bool
	}{
		{
			name: "upsert with on conflict",
			input: input{
				t: NewTable(database.NewBaseTable("db", "schema", "table")),
				fields: []database.Field{
					NewField(database.NewBaseField(0,
						"id", NewFieldType(newMockColumnType(oid.TypeName[oid.T_int8])))),
					NewField(database.NewBaseField(0,
						"name", NewFieldType(newMockColumnType(oid.TypeName[oid.T_varchar])))),
				},
				config: map[string]interface{}{
					"upsertSql": "ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name",
				},
				txOpts: nil,
			},
			args: args{
				records: nil, // records参数不会影响Query方法的结果
			},
			wantQuery: `insert into "schema"."table"("id","name") values ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name`,
			wantErr:   false,
		},
		{
			name: "no upsertSql config - should error",
			input: input{
				t: NewTable(database.NewBaseTable("db", "schema", "table")),
				fields: []database.Field{
					NewField(database.NewBaseField(0,
						"id", NewFieldType(newMockColumnType(oid.TypeName[oid.T_int8])))),
					NewField(database.NewBaseField(0,
						"name", NewFieldType(newMockColumnType(oid.TypeName[oid.T_varchar])))),
				},
				config: map[string]interface{}{}, // 空配置，没有upsertSQL
				txOpts: nil,
			},
			args: args{
				records: nil,
			},
			wantQuery: "",
			wantErr:   true, // 因为配置中没有upsertSQL会导致错误
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置字段
			for _, v := range tt.input.fields {
				tt.input.t.AppendField(v)
			}

			// 设置配置
			if tt.input.config != nil {
				// 将map转换为JSON字符串并设置到表的配置中
				jsonConfig, err := json.Marshal(tt.input.config)
				if err != nil {
					t.Fatalf("Failed to create JSON config: %v", err)
				}
				tt.input.t.SetConfig(testJSONFromString(string(jsonConfig)))
			}

			ci := NewUpsetParam(tt.input.t, tt.input.txOpts)
			got, gotErr := ci.Query(tt.args.records)
			if (gotErr != nil) != tt.wantErr {
				t.Errorf("UpsetParam.Query() error = %v, wantErr %v", gotErr, tt.wantErr)
				return
			}
			if got != tt.wantQuery {
				t.Errorf("UpsetParam.Query() = %v, want %v", got, tt.wantQuery)
			}
		})
	}
}
