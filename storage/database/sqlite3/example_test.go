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

package sqlite3_test

import (
	"context"
	"fmt"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/element"
	"github.com/Breeze0806/go-etl/storage/database"
	"github.com/Breeze0806/go-etl/storage/database/sqlite3"
)

func testJSONFromString(s string) *config.JSON {
	json, err := config.NewJSONFromString(s)
	if err != nil {
		panic(err)
	}
	return json
}

type TableParam struct {
	*database.BaseParam
}

func NewTableParam() *TableParam {
	return &TableParam{
		BaseParam: database.NewBaseParam(sqlite3.NewTable(database.NewBaseTable("", "", "test")), nil),
	}
}

func (t *TableParam) Query(_ []element.Record) (string, error) {
	return "select * from test", nil
}

func (t *TableParam) Agrs(_ []element.Record) ([]any, error) {
	return nil, nil
}

type FetchHandler struct {
}

func (f *FetchHandler) OnRecord(r element.Record) error {
	fmt.Println(r)
	return nil
}

func (f *FetchHandler) CreateRecord() (element.Record, error) {
	return element.NewDefaultRecord(), nil
}

func Example_sqlite3() {
	fmt.Println("strat")
	db, err := database.Open("sqlite3", testJSONFromString(`{"url":"E:\\projects\\sqlite3\\test.db"}`))
	if err != nil {
		fmt.Printf("open fail. err: %v", err)
		return
	}
	defer db.Close()
	err = db.FetchRecord(context.TODO(), NewTableParam(), &FetchHandler{})
	if err != nil {
		fmt.Printf("fetchRecord fail. err: %v", err)
		return
	}
}
