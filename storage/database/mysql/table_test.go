package mysql

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/Breeze0806/go-etl/element"
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

func TestTable_ExecParam(t *testing.T) {
	type args struct {
		mode   string
		txOpts *sql.TxOptions
	}
	tests := []struct {
		name  string
		t     *Table
		args  args
		want  database.Parameter
		want1 bool
	}{
		{
			name: "1",
			t:    NewTable(database.NewBaseTable("db", "", "table")),
			args: args{
				mode:   WriteModeReplace,
				txOpts: nil,
			},
			want:  NewReplaceParam(NewTable(database.NewBaseTable("db", "", "table")), nil),
			want1: true,
		},
		{
			name: "2",
			t:    NewTable(database.NewBaseTable("db", "", "table")),
			args: args{
				mode:   database.WriteModeInsert,
				txOpts: nil,
			},
			want:  nil,
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.t.ExecParam(tt.args.mode, tt.args.txOpts)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Table.ExecParam() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Table.ExecParam() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestReplaceParam_Query(t *testing.T) {
	type args struct {
		records []element.Record
		columns [][]element.Column
		fields  []database.Field
		t       *database.BaseTable
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
					NewField(database.NewBaseField(0, "f1", newMockFieldType("BIGINT"))),
					NewField(database.NewBaseField(1, "f2", newMockFieldType("DECIMAL"))),
					NewField(database.NewBaseField(2, "f3", newMockFieldType("STRING"))),
				},
				t: database.NewBaseTable("db", "", "table"),
			},
			wantQuery: "replace into `db`.`table`(`f1`,`f2`,`f3`) values(?,?,?),(?,?,?),(?,?,?)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, v := range tt.args.fields {
				tt.args.t.AppendField(v)
			}
			table := NewTable(tt.args.t)
			for i, r := range tt.args.records {
				for _, c := range tt.args.columns[i] {
					r.Add(c)
				}
			}

			rp, _ := table.ExecParam("replace", nil)
			gotQuery, err := rp.Query(tt.args.records)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReplaceParam.Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotQuery != tt.wantQuery {
				t.Errorf("ReplaceParam.Query() = %v, want %v", gotQuery, tt.wantQuery)
			}
		})
	}
}

func TestReplaceParam_Agrs(t *testing.T) {
	type args struct {
		records []element.Record
		columns [][]element.Column
		fields  []database.Field
		t       *database.BaseTable
	}
	tests := []struct {
		name        string
		rp          *ReplaceParam
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
					NewField(database.NewBaseField(0, "f1", newMockFieldType("BIGINT"))),
					NewField(database.NewBaseField(1, "f2", newMockFieldType("DECIMAL"))),
					NewField(database.NewBaseField(2, "f3", newMockFieldType("MEDIUMINT"))),
				},
				t: database.NewBaseTable("db", "", "table"),
			},
			wantValuers: []interface{}{
				"1", "2", "3",
				"4", "5", "6",
				"7", "8", "9",
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
						element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(6), "f3", 0),
					},
					{
						element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(9), "f3", 0),
						element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(7), "f1", 0),
						element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(8), "f2", 0),
					},
				},
				fields: []database.Field{
					NewField(database.NewBaseField(0, "f1", newMockFieldType("BIGINT"))),
					NewField(database.NewBaseField(1, "f2", newMockFieldType("DECIMAL"))),
					NewField(database.NewBaseField(2, "f3", newMockFieldType("STRING"))),
				},
				t: database.NewBaseTable("db", "", "table"),
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
						element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(8), "f2", 0),
					},
				},
				fields: []database.Field{
					NewField(database.NewBaseField(0, "f1", newMockFieldType("BIGINT"))),
					NewField(database.NewBaseField(1, "f2", newMockFieldType("DECIMAL"))),
					NewField(database.NewBaseField(2, "f4", newMockFieldType("MEDIUMINT"))),
				},
				t: database.NewBaseTable("db", "", "table"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, v := range tt.args.fields {
				tt.args.t.AppendField(v)
			}
			table := NewTable(tt.args.t)
			for i, r := range tt.args.records {
				for _, c := range tt.args.columns[i] {
					r.Add(c)
				}
			}

			rp, _ := table.ExecParam("replace", nil)
			gotValuers, err := rp.Agrs(tt.args.records)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReplaceParam.Agrs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotValuers, tt.wantValuers) {
				t.Errorf("ReplaceParam.Agrs() = %v, want %v", gotValuers, tt.wantValuers)
			}
		})
	}
}
