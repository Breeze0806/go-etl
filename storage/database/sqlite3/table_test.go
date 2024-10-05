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

package sqlite3

import (
	"database/sql"
	"database/sql/driver"
	"net"
	"reflect"
	"testing"

	"github.com/Breeze0806/go-etl/storage/database"
	sqlite3 "github.com/mattn/go-sqlite3"
)

func TestTable_Quoted(t *testing.T) {
	tests := []struct {
		name string
		tr   *Table
		want string
	}{
		{
			name: "1",
			tr:   NewTable(database.NewBaseTable("", "", "table")),
			want: `"table"`,
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
			tr:   NewTable(database.NewBaseTable("", "", "table")),
			want: `"table"`,
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
			tr:   NewTable(database.NewBaseTable("", "", "table")),
			args: args{
				err: nil,
			},
		},
		{
			name: "2",
			tr:   NewTable(database.NewBaseTable("", "", "table")),
			args: args{
				err: &net.AddrError{},
			},
			want: false,
		},
		{
			name: "3",
			tr:   NewTable(database.NewBaseTable("", "", "table")),
			args: args{
				err: driver.ErrBadConn,
			},
			want: true,
		},
		{
			name: "4",
			tr:   NewTable(database.NewBaseTable("", "", "table")),
			args: args{
				err: &sqlite3.Error{},
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
			tr:   NewTable(database.NewBaseTable("", "", "table")),
			args: args{
				err: nil,
			},
		},
		{
			name: "2",
			tr:   NewTable(database.NewBaseTable("", "", "table")),
			args: args{
				err: &net.AddrError{},
			},
		},
		{
			name: "3",
			tr:   NewTable(database.NewBaseTable("", "", "table")),
			args: args{
				err: driver.ErrBadConn,
			},
		},
		{
			name: "4",
			tr:   NewTable(database.NewBaseTable("", "", "table")),
			args: args{
				err: sqlite3.Error{},
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
			tr:   NewTable(database.NewBaseTable("", "", "table")),
			args: args{
				mode:   "",
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
			tr:   NewTable(database.NewBaseTable("", "schema", "table")),
			args: args{
				baseField: database.NewBaseField(0, "f1", database.NewBaseFieldType(&sql.ColumnType{})),
			},
			want: []database.Field{
				NewField(database.NewBaseField(0, "f1", database.NewBaseFieldType(&sql.ColumnType{}))),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.tr.AddField(tt.args.baseField)
			if !reflect.DeepEqual(tt.tr.Fields(), tt.want) {
				t.Errorf("run %v Table.Fields() = %v want: %v", tt.name, tt.tr.Fields(), tt.want)
			}
		})
	}
}
