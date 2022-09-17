package sqlserver

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/element"
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
			tr:   NewTable(database.NewBaseTable("db", "schema", "table")),
			args: args{
				mode:   "insert",
				txOpts: nil,
			},
			want:  nil,
			want1: false,
		},
		{
			name: "2",
			tr:   NewTable(database.NewBaseTable("db", "schema", "table")),
			args: args{
				mode:   WriteModeCopyIn,
				txOpts: nil,
			},
			want: NewCopyInParam(NewTable(database.NewBaseTable("db",
				"schema", "table")), nil),
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

func TestTable_SetConfig(t *testing.T) {
	type args struct {
		conf *config.JSON
	}
	tests := []struct {
		name string
		tr   *Table
		args args
		want *config.JSON
	}{
		{
			name: "1",
			tr:   &Table{},
			args: args{
				conf: testJSONFromString(`{
					"username": "",
					"password": "",
					"writeMode": "",
					"column": [],
					"preSql": [],
					"connection": {
						"url": "",
						"table": {
							"schema":"",
							"name":""
						}
					},
					"batchTimeout": "1s",
					"batchSize":1000,
					"bulkOption":{}
				}`),
			},
			want: testJSONFromString(`{
					"username": "",
					"password": "",
					"writeMode": "",
					"column": [],
					"preSql": [],
					"connection": {
						"url": "",
						"table": {
							"schema":"",
							"name":""
						}
					},
					"batchTimeout": "1s",
					"batchSize":1000,
					"bulkOption":{}
				}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.tr.SetConfig(tt.args.conf)
			if !reflect.DeepEqual(tt.tr.conf, tt.want) {
				t.Errorf("got: %v want: %v", tt.tr.conf, tt.want)
				return
			}
		})
	}
}

func TestCopyInParam_Query(t *testing.T) {
	type args struct {
		in0 []element.Record
	}

	new := func(t *database.BaseTable, conf *config.JSON) *Table {
		table := NewTable(t)
		table.SetConfig(conf)
		table.AddField(database.NewBaseField(1, "f1", newMockFieldType("int")))
		table.AddField(database.NewBaseField(2, "f2", newMockFieldType("bit")))
		table.AddField(database.NewBaseField(3, "f3", newMockFieldType("varchar")))
		return table
	}

	tests := []struct {
		name      string
		ci        *CopyInParam
		args      args
		wantQuery string
		wantErr   bool
	}{
		{
			name: "1",
			ci: NewCopyInParam(new(database.NewBaseTable("db", "schema", "table"),
				testJSONFromString(`{}`)), nil),
			wantQuery: `INSERTBULK {"TableName":"[db].[schema].[table]","ColumnsName":["f1","f2","f3"],"Options":{"CheckConstraints":false,"FireTriggers":false,"KeepNulls":false,"KilobytesPerBatch":0,"RowsPerBatch":0,"Order":null,"Tablock":false}}`,
		},
		{
			name: "2",
			ci: NewCopyInParam(new(database.NewBaseTable("db", "schema", "table"),
				testJSONFromString(`{"bulkOption":{"CheckConstraints":true,"FireTriggers":true,"KeepNulls":true,"KilobytesPerBatch":1000,"RowsPerBatch":1000,"Order":["f1","f2"],"Tablock":true}}`)), nil),
			wantQuery: `INSERTBULK {"TableName":"[db].[schema].[table]","ColumnsName":["f1","f2","f3"],"Options":{"CheckConstraints":true,"FireTriggers":true,"KeepNulls":true,"KilobytesPerBatch":1000,"RowsPerBatch":1000,"Order":["f1","f2"],"Tablock":true}}`,
		},
		{
			name: "3",
			ci: NewCopyInParam(new(database.NewBaseTable("db", "schema", "table"),
				testJSONFromString(`{"bulkOption":{"CheckConstraints":"true","FireTriggers":true,"KeepNulls":true,"KilobytesPerBatch":1000,"RowsPerBatch":1000,"Order":["f1","f2"],"Tablock":true}}`)), nil),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotQuery, err := tt.ci.Query(tt.args.in0)
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
		wantValuers []interface{}
		wantErr     bool
	}{
		{
			name: "1",
			input: input{
				t: NewTable(database.NewBaseTable("db", "schema", "table")),
				fields: []*database.BaseField{
					database.NewBaseField(0,
						"f1", newMockFieldType("INT")),
					database.NewBaseField(0,
						"f2", newMockFieldType("DECIMAL")),
					database.NewBaseField(0,
						"f3", newMockFieldType("VARCHAR")),
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
			wantValuers: []interface{}{
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
						"f1", newMockFieldType("INT")),
					database.NewBaseField(0,
						"f2", newMockFieldType("DECIMAL")),
					database.NewBaseField(0,
						"f3", newMockFieldType("VARCHAR")),
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
						"f1", newMockFieldType("INT")),
					database.NewBaseField(0,
						"f2", newMockFieldType("DECIMAL")),
					database.NewBaseField(0,
						"f3", newMockFieldType("DATE")),
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
