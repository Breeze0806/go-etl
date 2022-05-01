package db2

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
			tr:   NewTable(database.NewBaseTable("", "schema", "table")),
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
			tr:   NewTable(database.NewBaseTable("", "schema", "table")),
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
			args: args{},
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
