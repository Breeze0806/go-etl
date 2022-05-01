package db2

import (
	"context"
	"testing"

	"github.com/Breeze0806/go-etl/datax/plugin/reader/rdbm"
	"github.com/Breeze0806/go-etl/element"
	"github.com/Breeze0806/go-etl/storage/database"

	//db2 storage
	_ "github.com/Breeze0806/go-etl/storage/database/db2"
)

type testPrama struct {
	*database.BaseParam
}

func (t *testPrama) Query(_ []element.Record) (s string, err error) {
	s = "select * from " + t.Table().Quoted()
	return
}

func (t *testPrama) Agrs(_ []element.Record) ([]interface{}, error) {
	return nil, nil
}

func Test_Querier_FetchRecord(t *testing.T) {
	conf := rdbm.TestJSONFromString(`{
		"url":"HOSTNAME=192.168.15.130;PORT=50000;DATABASE=testdb",
		"username":"db2inst1",
		"password":"12345678"
}`)
	var records []element.Record
	q, err := database.Open("db2", conf)
	if err != nil {
		t.Fatalf("open fail. error: %v", err)
	}
	defer q.Close()
	table, err := q.FetchTable(context.TODO(), database.NewBaseTable(
		"", "TEST", "TEST"))
	if err != nil {
		t.Fatalf("fetchTable fail. error: %v", err)
	}
	err = q.FetchRecord(context.TODO(), &testPrama{
		BaseParam: database.NewBaseParam(table, nil),
	}, database.NewBaseFetchHandler(func() (element.Record, error) {
		return element.NewDefaultRecord(), nil
	}, func(r element.Record) error {
		records = append(records, r)
		return nil
	}))
	if err != nil {
		t.Fatalf("fetchRecord fail. error: %v", err)
	}
	opts := &database.ParameterOptions{
		Table:     table,
		Mode:      database.WriteModeInsert,
		TxOptions: nil,
		Records:   records,
	}
	if err := q.BatchExec(context.TODO(), opts); err != nil {
		t.Fatalf("BatchExec fail. error: %v", err)
	}
}
