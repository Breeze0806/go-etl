package database

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"testing"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/element"
)

func TestDB(t *testing.T) {
	registerMock()
	db, err := testDB("mock", testJsonFromString("{}"))
	if err != nil {
		t.Errorf("Open mock error %v", err)
		return
	}
	defer db.Close()
	gotTable, err := db.FetchTable(context.TODO(), NewBaseTable("db", "schema", "table"))
	if err != nil {
		t.Errorf("FetchTable error %v", err)
		return
	}

	wantTable := &mockTable{
		BaseTable: &BaseTable{
			instance: "db",
			schema:   "schema",
			name:     "table",
			fields: []Field{
				newMockField(NewBaseField("f1", newMockFieldType(GoTypeBool)), newMockFieldType(GoTypeBool)),
				newMockField(NewBaseField("f2", newMockFieldType(GoTypeInt64)), newMockFieldType(GoTypeInt64)),
				newMockField(NewBaseField("f3", newMockFieldType(GoTypeFloat64)), newMockFieldType(GoTypeFloat64)),
				newMockField(NewBaseField("f4", newMockFieldType(GoTypeString)), newMockFieldType(GoTypeString)),
			},
		},
	}

	if gotTable.String() != wantTable.String() {
		t.Errorf("got: %v want: %v", gotTable.String(), wantTable.String())
		return
	}

	if len(gotTable.Fields()) != len(wantTable.Fields()) {
		t.Errorf("got.field: %v want.fields: %v", len(gotTable.Fields()), len(wantTable.Fields()))
		return
	}

	for i, v := range gotTable.Fields() {
		if !reflect.DeepEqual(v.Name(), wantTable.Fields()[i].Name()) {
			t.Errorf("field %v got.name: %v want.name: %v", i, v.Name(), wantTable.Fields()[i].Name())
			return
		}
		if !reflect.DeepEqual(v.Name(), wantTable.Fields()[i].Name()) {
			t.Errorf("field %v got.type: %v want.type: %v", i, v.Type().DatabaseTypeName(), wantTable.Fields()[i].Type().DatabaseTypeName())
			return
		}
	}

	var gotRecords []element.Record

	if err = db.FetchRecord(context.TODO(), NewTableQueryParam(gotTable), func(r element.Record) error {
		gotRecords = append(gotRecords, r)
		return nil
	}); err != nil {
		t.Errorf("FetchTable error %v", err)
		return
	}
	columns := [][]element.Column{
		{
			element.NewDefaultColumn(element.NewBoolColumnValue(false), "f1", 0),
			element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(1), "f2", 0),
			element.NewDefaultColumn(element.NewDecimalColumnValueFromFloat(1.0), "f3", 0),
			element.NewDefaultColumn(element.NewStringColumnValue("1"), "f4", 0),
		},
		{
			element.NewDefaultColumn(element.NewBoolColumnValue(true), "f1", 0),
			element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(2), "f2", 0),
			element.NewDefaultColumn(element.NewDecimalColumnValueFromFloat(2.0), "f3", 0),
			element.NewDefaultColumn(element.NewStringColumnValue("2"), "f4", 0),
		},
	}
	var wantRecords []element.Record
	for _, row := range columns {
		record := element.NewDefaultRecord()
		for _, c := range row {
			record.Add(c)
		}
		wantRecords = append(wantRecords, record)
	}

	if !reflect.DeepEqual(gotRecords, wantRecords) {
		t.Errorf("got: %v want: %v", gotRecords, wantRecords)
	}

	gotRecords = nil
	if err = db.FetchRecordWithTx(context.TODO(), NewTableQueryParam(gotTable), func(r element.Record) error {
		gotRecords = append(gotRecords, r)
		return nil
	}); err != nil {
		t.Errorf("FetchTable error %v", err)
		return
	}
	if !reflect.DeepEqual(gotRecords, wantRecords) {
		t.Errorf("got: %v want: %v", gotRecords, wantRecords)
	}

	if err = db.BatchExec(context.TODO(), &ParameterOptions{
		Table:     gotTable,
		TxOptions: nil,
		Mode:      "insert",
		Records:   wantRecords,
	}); err != nil {
		t.Errorf("FetchTable error %v", err)
		return
	}

	if err = db.BatchExecWithTx(context.TODO(), &ParameterOptions{
		Table:     gotTable,
		TxOptions: nil,
		Mode:      "insert",
		Records:   wantRecords,
	}); err != nil {
		t.Errorf("FetchTable error %v", err)
		return
	}

	if err = db.BatchExecStmtWithTx(context.TODO(), &ParameterOptions{
		Table:     gotTable,
		TxOptions: nil,
		Mode:      "insert",
		Records:   wantRecords,
	}); err != nil {
		t.Errorf("FetchTable error %v", err)
		return
	}
}

func TestNewDB(t *testing.T) {
	registerMock()
	type args struct {
		name string
		conf *config.Json
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				name: "test",
				conf: testJsonFromString("{}"),
			},
			wantErr: true,
		},
		{
			name: "2",
			args: args{
				name: "mock",
				conf: testJsonFromString(`{"pool":{"connMaxIdleTime":"1","connMaxLifetime":"1"}}`),
			},
			wantErr: true,
		},
		{
			name: "3",
			args: args{
				name: "mock",
				conf: testJsonFromString(`{"pool":{"connMaxIdleTime":"1s","connMaxLifetime":"1s"}}`),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := testDB(tt.args.name, tt.args.conf)
			defer func() {
				if db != nil {
					db.Close()
				}
			}()
			if (err != nil) != tt.wantErr {
				t.Errorf("Open() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestDB_FetchTableWithParam(t *testing.T) {
	registerMock()
	type args struct {
		ctx   context.Context
		param Parameter
	}
	tests := []struct {
		name    string
		d       *DB
		args    args
		want    Table
		wantErr bool
	}{
		{
			name: "1",
			d:    testMustDB("mock", testJsonFromString("{}")),
			args: args{
				ctx: context.TODO(),
				param: NewTableQueryParam(&mockTableWithOther{
					mockTable: &mockTable{
						BaseTable: NewBaseTable("db", "schema", "table"),
					},
				}),
			},
			want: &mockTable{
				BaseTable: &BaseTable{
					instance: "db",
					schema:   "schema",
					name:     "table",
					fields: []Field{
						newMockField(NewBaseField("f1", newMockFieldType(GoTypeBool)), newMockFieldType(GoTypeBool)),
						newMockField(NewBaseField("f2", newMockFieldType(GoTypeInt64)), newMockFieldType(GoTypeInt64)),
						newMockField(NewBaseField("f3", newMockFieldType(GoTypeFloat64)), newMockFieldType(GoTypeFloat64)),
						newMockField(NewBaseField("f4", newMockFieldType(GoTypeString)), newMockFieldType(GoTypeString)),
					},
				},
			},
		},
		{
			name: "2",
			d:    testMustDB("mock", testJsonFromString("{}")),
			args: args{
				ctx: context.TODO(),
				param: NewTableQueryParam(&mockTableWithOther{
					mockTable: &mockTable{
						BaseTable: NewBaseTable("db", "schema", "table"),
					},
					err: errors.New("mock error"),
				}),
			},
			wantErr: true,
		},
		{
			name: "3",
			d:    testMustDB("mock", testJsonFromString("{}")),
			args: args{
				ctx: context.TODO(),
				param: NewTableQueryParam(&mockTableWithNoAdder{
					BaseTable: NewBaseTable("db", "schema", "table"),
				}),
			},
			wantErr: true,
		},
		{
			name: "4",
			d:    testMustDB("mock", testJsonFromString("{}")),
			args: args{
				ctx: context.TODO(),
				param: &mockParameter{
					BaseParam: NewBaseParam(newMockTable(NewBaseTable("db", "schema", "table")), nil),
					queryErr:  errors.New("mock error"),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer tt.d.Close()
			got, err := tt.d.FetchTableWithParam(tt.args.ctx, tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("DB.FetchTableWithParam() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got == nil && tt.want == nil {
				return
			}

			if got.String() != tt.want.String() {
				t.Errorf("got: %v want: %v", got.String(), tt.want.String())
				return
			}

			if len(got.Fields()) != len(tt.want.Fields()) {
				t.Errorf("got.field: %v want.fields: %v", len(got.Fields()), len(tt.want.Fields()))
				return
			}
			for i, v := range got.Fields() {
				if !reflect.DeepEqual(v.Name(), tt.want.Fields()[i].Name()) {
					t.Errorf("field %v got.name: %v want.name: %v", i, v.Name(), tt.want.Fields()[i].Name())
					return
				}
				if !reflect.DeepEqual(v.Name(), tt.want.Fields()[i].Name()) {
					t.Errorf("field %v got.type: %v want.type: %v", i, v.Type().DatabaseTypeName(), tt.want.Fields()[i].Type().DatabaseTypeName())
					return
				}
			}
		})
	}
}

func Test_getQueryAndAgrs(t *testing.T) {
	registerMock()
	type args struct {
		param   Parameter
		records []element.Record
	}
	tests := []struct {
		name      string
		args      args
		wantQuery string
		wantAgrs  []interface{}
		wantErr   bool
	}{
		{
			name: "1",
			args: args{
				param: &mockParameter{
					agrsErr: errors.New("mock error"),
				},
				records: nil,
			},
			wantErr: true,
		},
		{
			name: "2",
			args: args{
				param: &mockParameter{
					queryErr: errors.New("mock error"),
				},
				records: nil,
			},
			wantErr: true,
		},
		{
			name: "3",
			args: args{
				param:   NewTableQueryParam(newMockTable(NewBaseTable("db", "schema", "table"))),
				records: nil,
			},
			wantQuery: "select * from db.schema.table where 1 = 2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotQuery, gotAgrs, err := getQueryAndAgrs(tt.args.param, tt.args.records)
			if (err != nil) != tt.wantErr {
				t.Errorf("getQueryAndAgrs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotQuery != tt.wantQuery {
				t.Errorf("getQueryAndAgrs() gotQuery = %v, want %v", gotQuery, tt.wantQuery)
			}
			if !reflect.DeepEqual(gotAgrs, tt.wantAgrs) {
				t.Errorf("getQueryAndAgrs() gotAgrs = %v, want %v", gotAgrs, tt.wantAgrs)
			}
		})
	}
}

func Test_execParam(t *testing.T) {
	registerMock()
	type args struct {
		opts *ParameterOptions
	}
	tests := []struct {
		name      string
		args      args
		wantParam Parameter
		wantErr   bool
	}{
		{
			name: "1",
			args: args{
				opts: &ParameterOptions{
					Table: &mockTable{
						BaseTable: NewBaseTable("db", "schema", "table"),
					},
					Mode: "insert",
				},
			},
			wantParam: NewInsertParam(&mockTable{
				BaseTable: NewBaseTable("db", "schema", "table"),
			}, nil),
		},
		{
			name: "2",
			args: args{
				opts: &ParameterOptions{
					Table: &mockTable{
						BaseTable: NewBaseTable("db", "schema", "table"),
					},
					Mode: "copy in",
				},
			},
			wantErr: true,
		},
		{
			name: "3",
			args: args{
				opts: &ParameterOptions{
					Table: &mockTableWithOther{
						mockTable: &mockTable{
							BaseTable: NewBaseTable("db", "schema", "table"),
						},
					},
					Mode: "copy in",
				},
			},
			wantErr: true,
		},
		{
			name: "4",
			args: args{
				opts: &ParameterOptions{
					Table: &mockTableWithOther{
						mockTable: &mockTable{
							BaseTable: NewBaseTable("db", "schema", "table"),
						},
					},
					Mode: "insert",
				},
			},
			wantParam: NewInsertParam(&mockTable{
				BaseTable: NewBaseTable("db", "schema", "table"),
			}, nil),
		},
		{
			name: "5",
			args: args{
				opts: &ParameterOptions{
					Table: &mockTableWithOther{
						mockTable: &mockTable{
							BaseTable: NewBaseTable("db", "schema", "table"),
						},
						execParams: map[string]func(t Table, txOpts *sql.TxOptions) Parameter{
							"mock": func(t Table, txOpts *sql.TxOptions) Parameter {
								return &mockParameter{
									BaseParam: NewBaseParam(t, txOpts),
								}
							},
						},
					},
					Mode: "mock",
				},
			},
			wantParam: &mockParameter{
				BaseParam: NewBaseParam(&mockTableWithOther{
					mockTable: &mockTable{
						BaseTable: NewBaseTable("db", "schema", "table"),
					},
					execParams: map[string]func(t Table, txOpts *sql.TxOptions) Parameter{
						"mock": func(t Table, txOpts *sql.TxOptions) Parameter {
							return &mockParameter{
								BaseParam: NewBaseParam(t, txOpts),
							}
						},
					},
				}, nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotParam, err := execParam(tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("execParam() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotParam == nil && tt.wantParam == nil {
				return
			}
			if reflect.ValueOf(gotParam).Type() != reflect.ValueOf(tt.wantParam).Type() {
				t.Errorf("execParam() = %T, want %T", gotParam, tt.wantParam)
			}

			if !reflect.DeepEqual(gotParam.Table().String(), tt.wantParam.Table().String()) {
				t.Errorf("execParam() = %v, want %v", gotParam, tt.wantParam)
			}
		})
	}
}

func TestDB_BatchExec(t *testing.T) {
	registerMock()
	type args struct {
		ctx  context.Context
		opts *ParameterOptions
	}
	tests := []struct {
		name    string
		d       *DB
		args    args
		wantErr bool
	}{
		{
			name: "1",
			d:    testMustDB("mock", testJsonFromString("{}")),
			args: args{
				ctx: context.TODO(),
				opts: &ParameterOptions{
					Table: &mockTableWithOther{
						mockTable: &mockTable{
							BaseTable: NewBaseTable("db", "schema", "table"),
						},
						execParams: map[string]func(t Table, txOpts *sql.TxOptions) Parameter{
							"mock": func(t Table, txOpts *sql.TxOptions) Parameter {
								return &mockParameter{
									BaseParam: NewBaseParam(t, txOpts),
									queryErr:  errors.New("mock error"),
								}
							},
						},
					},
					Mode: "mock",
				},
			},
			wantErr: true,
		},
		{
			name: "2",
			d:    testMustDB("mock", testJsonFromString("{}")),
			args: args{
				ctx: context.TODO(),
				opts: &ParameterOptions{
					Table: &mockTableWithOther{
						mockTable: &mockTable{
							BaseTable: NewBaseTable("db", "schema", "table"),
						},
						execParams: map[string]func(t Table, txOpts *sql.TxOptions) Parameter{
							"mock": func(t Table, txOpts *sql.TxOptions) Parameter {
								return &mockParameter{
									BaseParam: NewBaseParam(t, txOpts),
									queryErr:  errors.New("mock error"),
								}
							},
						},
					},
					Mode: "mock1",
				},
			},
			wantErr: true,
		},
		{
			name: "3",
			d:    testMustDB("mock", testJsonFromString("{}")),
			args: args{
				ctx: context.TODO(),
				opts: &ParameterOptions{
					Table: &mockTableWithOther{
						mockTable: &mockTable{
							BaseTable: NewBaseTable("db", "schema", "table"),
						},
						execParams: map[string]func(t Table, txOpts *sql.TxOptions) Parameter{
							"mock": func(t Table, txOpts *sql.TxOptions) Parameter {
								return &mockParameter{
									BaseParam: NewBaseParam(t, txOpts),
								}
							},
						},
					},
					Mode: "mock",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer tt.d.Close()
			if err := tt.d.BatchExec(tt.args.ctx, tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("DB.BatchExec() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDB_BatchExecWithTx(t *testing.T) {
	registerMock()
	type args struct {
		ctx  context.Context
		opts *ParameterOptions
	}
	tests := []struct {
		name    string
		d       *DB
		args    args
		wantErr bool
	}{
		{
			name: "1",
			d:    testMustDB("mock", testJsonFromString("{}")),
			args: args{
				ctx: context.TODO(),
				opts: &ParameterOptions{
					Table: &mockTableWithOther{
						mockTable: &mockTable{
							BaseTable: NewBaseTable("db", "schema", "table"),
						},
						execParams: map[string]func(t Table, txOpts *sql.TxOptions) Parameter{
							"mock": func(t Table, txOpts *sql.TxOptions) Parameter {
								return &mockParameter{
									BaseParam: NewBaseParam(t, txOpts),
									queryErr:  errors.New("mock error"),
								}
							},
						},
					},
					Mode: "mock",
				},
			},
			wantErr: true,
		},
		{
			name: "2",
			d:    testMustDB("mock", testJsonFromString("{}")),
			args: args{
				ctx: context.TODO(),
				opts: &ParameterOptions{
					Table: &mockTableWithOther{
						mockTable: &mockTable{
							BaseTable: NewBaseTable("db", "schema", "table"),
						},
						execParams: map[string]func(t Table, txOpts *sql.TxOptions) Parameter{
							"mock": func(t Table, txOpts *sql.TxOptions) Parameter {
								return &mockParameter{
									BaseParam: NewBaseParam(t, txOpts),
									queryErr:  errors.New("mock error"),
								}
							},
						},
					},
					Mode: "mock1",
				},
			},
			wantErr: true,
		},
		{
			name: "3",
			d:    testMustDB("mock", testJsonFromString("{}")),
			args: args{
				ctx: context.TODO(),
				opts: &ParameterOptions{
					Table: &mockTableWithOther{
						mockTable: &mockTable{
							BaseTable: NewBaseTable("db", "schema", "table"),
						},
						execParams: map[string]func(t Table, txOpts *sql.TxOptions) Parameter{
							"mock": func(t Table, txOpts *sql.TxOptions) Parameter {
								return &mockParameter{
									BaseParam: NewBaseParam(t, txOpts),
								}
							},
						},
					},
					Mode: "mock",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer tt.d.Close()
			if err := tt.d.BatchExecWithTx(tt.args.ctx, tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("DB.BatchExecWithTx() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDB_BatchExecStmtWithTx(t *testing.T) {
	registerMock()
	type args struct {
		ctx  context.Context
		opts *ParameterOptions
	}
	tests := []struct {
		name    string
		d       *DB
		args    args
		wantErr bool
	}{
		{
			name: "1",
			d:    testMustDB("mock", testJsonFromString("{}")),
			args: args{
				ctx: context.TODO(),
				opts: &ParameterOptions{
					Table: &mockTableWithOther{
						mockTable: &mockTable{
							BaseTable: NewBaseTable("db", "schema", "table"),
						},
						execParams: map[string]func(t Table, txOpts *sql.TxOptions) Parameter{
							"mock": func(t Table, txOpts *sql.TxOptions) Parameter {
								return &mockParameter{
									BaseParam: NewBaseParam(t, txOpts),
									queryErr:  errors.New("mock error"),
								}
							},
						},
					},
					Mode: "mock",
				},
			},
			wantErr: true,
		},
		{
			name: "2",
			d:    testMustDB("mock", testJsonFromString("{}")),
			args: args{
				ctx: context.TODO(),
				opts: &ParameterOptions{
					Table: &mockTableWithOther{
						mockTable: &mockTable{
							BaseTable: NewBaseTable("db", "schema", "table"),
						},
						execParams: map[string]func(t Table, txOpts *sql.TxOptions) Parameter{
							"mock": func(t Table, txOpts *sql.TxOptions) Parameter {
								return &mockParameter{
									BaseParam: NewBaseParam(t, txOpts),
									queryErr:  errors.New("mock error"),
								}
							},
						},
					},
					Mode: "mock1",
				},
			},
			wantErr: true,
		},
		{
			name: "3",
			d:    testMustDB("mock", testJsonFromString("{}")),
			args: args{
				ctx: context.TODO(),
				opts: &ParameterOptions{
					Table: &mockTableWithOther{
						mockTable: &mockTable{
							BaseTable: NewBaseTable("db", "schema", "table"),
						},
						execParams: map[string]func(t Table, txOpts *sql.TxOptions) Parameter{
							"mock": func(t Table, txOpts *sql.TxOptions) Parameter {
								return &mockParameter{
									BaseParam: NewBaseParam(t, txOpts),
								}
							},
						},
					},
					Mode: "mock",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer tt.d.Close()
			if err := tt.d.BatchExecStmtWithTx(tt.args.ctx, tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("DB.BatchExecStmtWithTx() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDB_FetchRecord(t *testing.T) {
	registerMock()
	db := testMustDB("mock", testJsonFromString("{}"))
	defer db.Close()
	table, _ := db.FetchTable(context.TODO(), NewBaseTable("db", "schema", "table"))
	type args struct {
		ctx      context.Context
		param    Parameter
		onRecord func(element.Record) error
	}
	tests := []struct {
		name    string
		d       *DB
		args    args
		wantErr bool
	}{
		{
			name: "1",
			d:    db,
			args: args{
				ctx: context.TODO(),
				param: &mockParameter{
					BaseParam: NewBaseParam(table, nil),
					queryErr:  errors.New("mock error"),
				},
				onRecord: func(r element.Record) error {
					return nil
				},
			},
			wantErr: true,
		},
		{
			name: "2",
			d:    db,
			args: args{
				ctx: context.TODO(),
				param: &mockParameter{
					BaseParam: NewBaseParam(table, nil),
				},
				onRecord: func(r element.Record) error {
					return errors.New("mock error")
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.d.FetchRecord(tt.args.ctx, tt.args.param, tt.args.onRecord); (err != nil) != tt.wantErr {
				t.Errorf("DB.FetchRecord() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDB_FetchRecordWithTx(t *testing.T) {
	registerMock()
	db := testMustDB("mock", testJsonFromString("{}"))
	defer db.Close()
	table, _ := db.FetchTable(context.TODO(), NewBaseTable("db", "schema", "table"))
	type args struct {
		ctx      context.Context
		param    Parameter
		onRecord func(element.Record) error
	}
	tests := []struct {
		name    string
		d       *DB
		args    args
		wantErr bool
	}{
		{
			name: "1",
			d:    db,
			args: args{
				ctx: context.TODO(),
				param: &mockParameter{
					BaseParam: NewBaseParam(table, nil),
					queryErr:  errors.New("mock error"),
				},
				onRecord: func(r element.Record) error {
					return nil
				},
			},
			wantErr: true,
		},
		{
			name: "2",
			d:    db,
			args: args{
				ctx: context.TODO(),
				param: &mockParameter{
					BaseParam: NewBaseParam(table, nil),
				},
				onRecord: func(r element.Record) error {
					return errors.New("mock error")
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.d.FetchRecordWithTx(tt.args.ctx, tt.args.param, tt.args.onRecord)
			if (err != nil) != tt.wantErr {
				t.Errorf("DB.FetchRecordWithTx() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDB_Close(t *testing.T) {
	tests := []struct {
		name    string
		d       *DB
		wantErr bool
	}{
		{
			name:    "1",
			d:       &DB{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.d.Close(); (err != nil) != tt.wantErr {
				t.Errorf("DB.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
