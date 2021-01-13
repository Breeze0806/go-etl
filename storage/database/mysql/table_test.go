package mysql

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/Breeze0806/go-etl/storage/database"
)

func TestTable_Quoted(t *testing.T) {
	tests := []struct {
		name string
		t    *Table
		want string
	}{
		{
			name: "1",
			t:    NewTable(database.NewBaseTable("db", "schema", "table")),
			want: "`db`.`table`",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.Quoted(); got != tt.want {
				t.Errorf("Table.Quoted() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTable_String(t *testing.T) {
	tests := []struct {
		name string
		t    *Table
		want string
	}{
		{
			name: "1",
			t:    NewTable(database.NewBaseTable("db", "schema", "table")),
			want: "`db`.`table`",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.String(); got != tt.want {
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
		t    *Table
		args args
		want []database.Field
	}{
		{
			name: "1",
			t:    table,
			args: args{
				baseField: database.NewBaseField("f1", &sql.ColumnType{}),
			},
			want: []database.Field{
				NewField(database.NewBaseField("f1", &sql.ColumnType{})),
			},
		},
		{
			name: "2",
			t:    table,
			args: args{
				baseField: database.NewBaseField("f2", &sql.ColumnType{}),
			},
			want: []database.Field{
				NewField(database.NewBaseField("f1", &sql.ColumnType{})),
				NewField(database.NewBaseField("f2", &sql.ColumnType{})),
			},
		},
		{
			name: "3",
			t:    table,
			args: args{
				baseField: database.NewBaseField("f3", &sql.ColumnType{}),
			},
			want: []database.Field{
				NewField(database.NewBaseField("f1", &sql.ColumnType{})),
				NewField(database.NewBaseField("f2", &sql.ColumnType{})),
				NewField(database.NewBaseField("f3", &sql.ColumnType{})),
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
