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

package mysql

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"fmt"

	"github.com/Breeze0806/go-etl/element"
	"github.com/Breeze0806/go-etl/storage/database"
)

const WriteModeReplace = "replace"

//Table mysql表
type Table struct {
	*database.BaseTable
}

//NewTable 创建mysql表，注意此时BaseTable中的schema参数为空，instance为数据库名，而name是表明
func NewTable(b *database.BaseTable) *Table {
	return &Table{
		BaseTable: b,
	}
}

//Quoted 表引用全名
func (t *Table) Quoted() string {
	return Quoted(t.Instance()) + "." + Quoted(t.Name())
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
	case "replace":
		return NewReplaceParam(t, txOpts), true
	}
	return nil, false
}

//ReplaceParam Replace into 参数
type ReplaceParam struct {
	*database.BaseParam
}

//NewReplaceParam  通过表table和事务参数txOps插入参数
func NewReplaceParam(t database.Table, txOps *sql.TxOptions) *ReplaceParam {
	return &ReplaceParam{
		BaseParam: database.NewBaseParam(t, txOps),
	}
}

//Query 通过多条记录 records生成批量Replace into插入sql语句
func (rp *ReplaceParam) Query(records []element.Record) (query string, err error) {
	buf := bytes.NewBufferString("replace into ")
	buf.WriteString(rp.Table().Quoted())
	buf.WriteString("(")
	for fi, f := range rp.Table().Fields() {
		if fi > 0 {
			buf.WriteString(",")
		}
		_, err = buf.WriteString(f.Quoted())
	}
	buf.WriteString(") values")

	for ri := range records {
		if ri > 0 {
			buf.WriteString(",")
		}
		buf.WriteString("(")
		for fi, f := range rp.Table().Fields() {
			if fi > 0 {
				buf.WriteString(",")
			}
			buf.WriteString(
				f.BindVar(ri*len(rp.Table().Fields()) + fi + 1))
		}
		buf.WriteString(")")
	}
	return buf.String(), nil
}

//Agrs 通过多条记录 records生成批量Replace into参数
func (rp *ReplaceParam) Agrs(records []element.Record) (valuers []interface{}, err error) {
	for _, r := range records {
		for fi, f := range rp.Table().Fields() {
			var c element.Column
			if c, err = r.GetByIndex(fi); err != nil {
				return nil, fmt.Errorf("GetByIndex(%v) err: %v", fi, err)
			}
			var v driver.Value
			if v, err = f.Valuer(c).Value(); err != nil {
				return nil, err
			}

			valuers = append(valuers, interface{}(v))
		}
	}
	return
}
