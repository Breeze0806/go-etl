package database

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/Breeze0806/go-etl/element"
)

func TestBaseTable_Instance(t *testing.T) {
	tests := []struct {
		name string
		b    *BaseTable
		want string
	}{
		{
			name: "1",
			b:    NewBaseTable("db", "schema", "table"),
			want: "db",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.Instance(); got != tt.want {
				t.Errorf("BaseTable.Instance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseTable_Schema(t *testing.T) {
	tests := []struct {
		name string
		b    *BaseTable
		want string
	}{
		{
			name: "1",
			b:    NewBaseTable("db", "schema", "table"),
			want: "schema",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.Schema(); got != tt.want {
				t.Errorf("BaseTable.Schema() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseTable_Name(t *testing.T) {
	tests := []struct {
		name string
		b    *BaseTable
		want string
	}{
		{
			name: "1",
			b:    NewBaseTable("db", "schema", "table"),
			want: "table",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.Name(); got != tt.want {
				t.Errorf("BaseTable.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseTable_Fields(t *testing.T) {
	b := NewBaseTable("", "", "")
	type args struct {
		f Field
	}

	tests := []struct {
		name string
		b    *BaseTable
		args args
		want []Field
	}{
		{
			name: "1",
			args: args{
				f: newMockField(NewBaseField(1, "1", NewBaseFieldType(&sql.ColumnType{})), NewBaseFieldType(&sql.ColumnType{})),
			},
			want: []Field{
				newMockField(NewBaseField(1, "1", NewBaseFieldType(&sql.ColumnType{})), NewBaseFieldType(&sql.ColumnType{})),
			},
		},
		{
			name: "2",
			args: args{
				f: newMockField(NewBaseField(1, "2", NewBaseFieldType(&sql.ColumnType{})), NewBaseFieldType(&sql.ColumnType{})),
			},
			want: []Field{
				newMockField(NewBaseField(1, "1", NewBaseFieldType(&sql.ColumnType{})), NewBaseFieldType(&sql.ColumnType{})),
				newMockField(NewBaseField(1, "2", NewBaseFieldType(&sql.ColumnType{})), NewBaseFieldType(&sql.ColumnType{})),
			},
		},
		{
			name: "3",
			args: args{
				f: newMockField(NewBaseField(1, "3", NewBaseFieldType(&sql.ColumnType{})), NewBaseFieldType(&sql.ColumnType{})),
			},
			want: []Field{
				newMockField(NewBaseField(1, "1", NewBaseFieldType(&sql.ColumnType{})), NewBaseFieldType(&sql.ColumnType{})),
				newMockField(NewBaseField(1, "2", NewBaseFieldType(&sql.ColumnType{})), NewBaseFieldType(&sql.ColumnType{})),
				newMockField(NewBaseField(1, "3", NewBaseFieldType(&sql.ColumnType{})), NewBaseFieldType(&sql.ColumnType{})),
			},
		},
	}
	for _, tt := range tests {
		b.AppendField(tt.args.f)
		if got := b.Fields(); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("run %v BaseTable.Fields() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestBaseParam_Table(t *testing.T) {
	tests := []struct {
		name string
		b    *BaseParam
		want Table
	}{
		{
			name: "1",
			b:    NewBaseParam(newMockTable(NewBaseTable("db", "schema", "table")), nil),
			want: newMockTable(NewBaseTable("db", "schema", "table")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.Table(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BaseParam.Table() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInsertParam_TxOptions(t *testing.T) {
	tests := []struct {
		name string
		i    *InsertParam
		want *sql.TxOptions
	}{
		{
			name: "1",
			i: NewInsertParam(nil, &sql.TxOptions{
				Isolation: sql.LevelRepeatableRead,
				ReadOnly:  true,
			}),
			want: &sql.TxOptions{
				Isolation: sql.LevelRepeatableRead,
				ReadOnly:  true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.TxOptions(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InsertParam.TxOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInsertParam_Query(t *testing.T) {
	type args struct {
		records []element.Record
		columns [][]element.Column
		fields  []Field
		t       *BaseTable
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
				fields: []Field{
					newMockField(NewBaseField(0, "f1", nil), newMockFieldType(GoTypeInt64)),
					newMockField(NewBaseField(1, "f2", nil), newMockFieldType(GoTypeFloat64)),
					newMockField(NewBaseField(2, "f3", nil), newMockFieldType(GoTypeString)),
				},
				t: NewBaseTable("db", "schema", "table"),
			},
			wantQuery: "insert into db.schema.table(f1,f2,f3) values($1,$2,$3),($4,$5,$6),($7,$8,$9)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, v := range tt.args.fields {
				tt.args.t.AppendField(v)
			}
			table := newMockTable(tt.args.t)
			for i, r := range tt.args.records {
				for _, c := range tt.args.columns[i] {
					r.Add(c)
				}
			}

			insertParam := NewInsertParam(table, nil)
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
		fields  []Field
		t       *BaseTable
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
				fields: []Field{
					newMockField(NewBaseField(0, "f1", nil), newMockFieldType(GoTypeInt64)),
					newMockField(NewBaseField(1, "f2", nil), newMockFieldType(GoTypeFloat64)),
					newMockField(NewBaseField(2, "f3", nil), newMockFieldType(GoTypeString)),
				},
				t: NewBaseTable("db", "schema", "table"),
			},
			wantValuers: []interface{}{
				int64(1), float64(2.0), "3",
				int64(4), float64(5.0), "6",
				int64(7), float64(8.0), "9",
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
						element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(6), "f4", 0),
					},
					{
						element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(9), "f3", 0),
						element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(7), "f1", 0),
						element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(8), "f2", 0),
					},
				},
				fields: []Field{
					newMockField(NewBaseField(0, "f1", nil), newMockFieldType(GoTypeInt64)),
					newMockField(NewBaseField(1, "f2", nil), newMockFieldType(GoTypeFloat64)),
					newMockField(NewBaseField(2, "f3", nil), newMockFieldType(GoTypeString)),
				},
				t: NewBaseTable("db", "schema", "table"),
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
				fields: []Field{
					newMockField(NewBaseField(0, "f1", nil), newMockFieldType(GoTypeInt64)),
					newMockField(NewBaseField(1, "f2", nil), newMockFieldType(GoTypeFloat64)),
					newMockField(NewBaseField(2, "f3", nil), newMockFieldType(GoTypeTime)),
				},
				t: NewBaseTable("db", "schema", "table"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, v := range tt.args.fields {
				tt.args.t.AppendField(v)
			}
			table := newMockTable(tt.args.t)
			for i, r := range tt.args.records {
				for _, c := range tt.args.columns[i] {
					r.Add(c)
				}
			}

			insertParam := NewInsertParam(table, nil)
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

func TestTableQueryParam_TxOptions(t *testing.T) {
	tests := []struct {
		name string
		t    *TableQueryParam
		want *sql.TxOptions
	}{
		{
			name: "1",
			t:    NewTableQueryParam(newMockTable(NewBaseTable("db", "schema", "table"))),
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.TxOptions(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TableQueryParam.TxOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTableQueryParam_Query(t *testing.T) {
	type args struct {
		records []element.Record
	}
	tests := []struct {
		name    string
		t       *TableQueryParam
		args    args
		wantS   string
		wantErr bool
	}{
		{
			name: "1",
			t:    NewTableQueryParam(newMockTable(NewBaseTable("db", "schema", "table"))),
			args: args{
				records: nil,
			},
			wantS: "select * from db.schema.table where 1 = 2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotS, err := tt.t.Query(tt.args.records)
			if (err != nil) != tt.wantErr {
				t.Errorf("TableQueryParam.Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotS != tt.wantS {
				t.Errorf("TableQueryParam.Query() = %v, want %v", gotS, tt.wantS)
			}
		})
	}
}

func TestTableQueryParam_Agrs(t *testing.T) {
	type args struct {
		records []element.Record
	}
	tests := []struct {
		name    string
		t       *TableQueryParam
		args    args
		wantA   []interface{}
		wantErr bool
	}{
		{
			name: "1",
			t:    NewTableQueryParam(newMockTable(NewBaseTable("db", "schema", "table"))),
			args: args{
				records: nil,
			},
			wantA: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotA, err := tt.t.Agrs(tt.args.records)
			if (err != nil) != tt.wantErr {
				t.Errorf("TableQueryParam.Agrs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotA, tt.wantA) {
				t.Errorf("TableQueryParam.Agrs() = %v, want %v", gotA, tt.wantA)
			}
		})
	}
}

func TestBaseParam_SetTable(t *testing.T) {
	type args struct {
		t Table
	}
	tests := []struct {
		name string
		b    *BaseParam
		args args
		want Table
	}{
		{
			name: "1",
			b:    NewBaseParam(nil, nil),
			args: args{
				t: newMockTable(NewBaseTable("db", "schema", "table")),
			},
			want: newMockTable(NewBaseTable("db", "schema", "table")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.SetTable(tt.args.t)
			if !reflect.DeepEqual(tt.b.Table(), tt.want) {
				t.Errorf("got: %v want: %v", tt.b.Table(), tt.want)
			}
		})
	}
}

func TestBaseParam_SettxOps(t *testing.T) {
	type args struct {
		txOps *sql.TxOptions
	}
	tests := []struct {
		name string
		b    *BaseParam
		args args
		want *sql.TxOptions
	}{
		{
			name: "1",
			b:    NewBaseParam(nil, nil),
			args: args{
				&sql.TxOptions{
					Isolation: sql.LevelRepeatableRead,
					ReadOnly:  true,
				},
			},
			want: &sql.TxOptions{
				Isolation: sql.LevelRepeatableRead,
				ReadOnly:  true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.SetTxOps(tt.args.txOps)
			if !reflect.DeepEqual(tt.b.TxOptions(), tt.want) {
				t.Errorf("got: %v want: %v", tt.b.TxOptions(), tt.want)
			}
		})
	}
}
