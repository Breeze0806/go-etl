package database

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"reflect"
	"sync"
	"testing"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/element"
)

var once sync.Once

func registerMock() {
	once.Do(func() {
		RegisterDialect("mock", &mockDialect{
			name: "mock",
		})
		RegisterDialect("mockErr", &mockDialect{
			name: "",
			err:  errors.New("mock error"),
		})
		RegisterDialect("test", &mockDialect{
			name: "test",
		})
		sql.Register("mock", &mockDriver{
			rows: &mockRows{
				columns: []string{
					"f1", "f2", "f3", "f4",
				},
				types: []*mockFieldType{
					newMockFieldType(GoTypeBool),
					newMockFieldType(GoTypeInt64),
					newMockFieldType(GoTypeFloat64),
					newMockFieldType(GoTypeString),
				},
				columnValues: [][]driver.Value{
					{false, int64(1), float64(1), string("1")},
					{true, int64(2), float64(2), string("2")},
				},
			},
		})
	})
}

func TestDB(t *testing.T) {
	registerMock()
	db, err := Open("mock", testJsonFromString("{}"))
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

	for i, v := range gotTable.Fields() {
		if !reflect.DeepEqual(v.Name(), wantTable.Fields()[i].Name()) {
			t.Errorf("%v got.name: %v want.name: %v", i, v.Name(), wantTable.Fields()[i].Name())
			return
		}
		if !reflect.DeepEqual(v.Name(), wantTable.Fields()[i].Name()) {
			t.Errorf("%v got.type: %v want.type: %v", i, v.Type().DatabaseTypeName(), wantTable.Fields()[i].Type().DatabaseTypeName())
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

func TestOpen(t *testing.T) {
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
				name: "test?",
				conf: testJsonFromString("{}"),
			},
			wantErr: true,
		},
		{
			name: "3",
			args: args{
				name: "mockErr",
				conf: testJsonFromString("{}"),
			},
			wantErr: true,
		},
		{
			name: "4",
			args: args{
				name: "mock",
				conf: testJsonFromString(`{"connMaxIdleTime":"1","connMaxLifetime":"1"}`),
			},
			wantErr: true,
		},
		{
			name: "5",
			args: args{
				name: "mock",
				conf: testJsonFromString(`{"connMaxIdleTime":"1s","connMaxLifetime":"1s"}`),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Open(tt.args.name, tt.args.conf)
			t.Log(err)
			if (err != nil) != tt.wantErr {
				t.Errorf("Open() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
