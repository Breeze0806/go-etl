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

package oracle

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"fmt"

	"github.com/Breeze0806/go-etl/element"
	"github.com/Breeze0806/go-etl/storage/database"
	"github.com/godror/godror"
	"github.com/pingcap/errors"
)

//WriteModeInsert intert into 写入方式
const WriteModeInsert = "insert"

//Table oracle表
type Table struct {
	*database.BaseTable
}

//NewTable 创建oracle表，注意此时BaseTable中的schema参数为空，instance为数据库名，而name是表明
func NewTable(b *database.BaseTable) *Table {
	return &Table{
		BaseTable: b,
	}
}

//Quoted 表引用全名
func (t *Table) Quoted() string {
	return Quoted(t.Schema()) + "." + Quoted(t.Name())
}

func (t *Table) String() string {
	return t.Quoted()
}

//AddField 新增列
func (t *Table) AddField(baseField *database.BaseField) {
	t.AppendField(NewField(baseField))
}

//ExecParam 获取执行参数，其中replace into的参数方式以及被注册
func (t *Table) ExecParam(mode string, txOpts *sql.TxOptions) (database.Parameter, bool) {
	switch mode {
	case WriteModeInsert:
		return NewInsertParam(t, txOpts), true
	}
	return nil, false
}

//ShouldRetry 重试
func (t *Table) ShouldRetry(err error) bool {
	return godror.IsBadConn(errors.Cause(err))
}

//ShouldOneByOne 单个重试
func (t *Table) ShouldOneByOne(err error) bool {
	_, ok := errors.Cause(err).(*godror.OraErr)
	return ok && !godror.IsBadConn(err)
}

//InsertParam Insert into 参数
type InsertParam struct {
	*database.BaseParam
}

//NewInsertParam 通过表table和事务参数txOpts插入参数
func NewInsertParam(t database.Table, txOpts *sql.TxOptions) *InsertParam {
	return &InsertParam{
		BaseParam: database.NewBaseParam(t, txOpts),
	}
}

//Query 通过多条记录 records生成批量insert into插入sql语句
func (ip *InsertParam) Query(_ []element.Record) (query string, err error) {
	buf := bytes.NewBufferString("insert into ")
	buf.WriteString(ip.Table().Quoted())
	buf.WriteString("(")
	for fi, f := range ip.Table().Fields() {
		if fi > 0 {
			buf.WriteString(",")
		}
		_, err = buf.WriteString(f.Quoted())
	}
	buf.WriteString(") values (")

	for fi, f := range ip.Table().Fields() {
		if fi > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(
			f.BindVar(fi + 1))
	}
	buf.WriteString(")")

	return buf.String(), nil
}

//Agrs 通过多条记录 records生成批量insert into参数
func (ip *InsertParam) Agrs(records []element.Record) (valuers []interface{}, err error) {
	for fi, f := range ip.Table().Fields() {
		var ba [][]byte
		var sa []string
		for _, r := range records {
			var c element.Column
			if c, err = r.GetByIndex(fi); err != nil {
				return nil, fmt.Errorf("GetByIndex(%v) err: %v", fi, err)
			}
			var v driver.Value
			if v, err = f.Valuer(c).Value(); err != nil {
				return nil, err
			}
			switch data := v.(type) {
			case nil:
				ba = append(ba, nil)
			case []byte:
				ba = append(ba, data)
			case string:
				sa = append(sa, data)
			}
		}
		var a interface{}

		if len(ba) > 0 {
			a = ba
		}
		if len(sa) > 0 {
			a = sa
		}
		valuers = append(valuers, a)
	}
	return
}
