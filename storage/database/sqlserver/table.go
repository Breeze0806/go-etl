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

package sqlserver

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/element"
	"github.com/Breeze0806/go-etl/storage/database"
	mssql "github.com/denisenkom/go-mssqldb"
)

//WriteModeCopyIn copy in写入方式
const WriteModeCopyIn = "copyIn"

//Table mssql表
type Table struct {
	*database.BaseTable

	conf *config.JSON
}

//NewTable 创建mssql表，注意此时BaseTable中的schema参数为空，instance为数据库名，而name是表明
func NewTable(b *database.BaseTable) *Table {
	return &Table{
		BaseTable: b,
	}
}

//Quoted 表引用全名
func (t *Table) Quoted() string {
	return Quoted(t.Instance()) + "." + Quoted(t.Schema()) + "." + Quoted(t.Name())
}

func (t *Table) String() string {
	return t.Quoted()
}

//AddField 新增列
func (t *Table) AddField(baseField *database.BaseField) {
	t.AppendField(NewField(baseField))
}

//SetConfig 设置配置
func (t *Table) SetConfig(conf *config.JSON) {
	t.conf = conf
}

//ExecParam 获取执行参数，其中replace into的参数方式以及被注册
func (t *Table) ExecParam(mode string, txOpts *sql.TxOptions) (database.Parameter, bool) {
	switch mode {
	case WriteModeCopyIn:
		return NewCopyInParam(t, txOpts), true
	}
	return nil, false
}

//CopyInParam copy in 参数
type CopyInParam struct {
	*database.BaseParam
}

//NewCopyInParam  通过表table和事务参数txOpts插入参数
func NewCopyInParam(t database.Table, txOpts *sql.TxOptions) *CopyInParam {
	return &CopyInParam{
		BaseParam: database.NewBaseParam(t, txOpts),
	}
}

//Query 批量copy in插入sql语句
func (ci *CopyInParam) Query(_ []element.Record) (query string, err error) {
	var conf *config.JSON
	conf, err = ci.Table().(*Table).conf.GetConfig("bulkOption")
	if err != nil {
		err = nil
		conf, _ = config.NewJSONFromString("{}")
	}

	opt := mssql.BulkOptions{}
	err = json.Unmarshal([]byte(conf.String()), &opt)
	if err != nil {
		return
	}

	var columns []string
	for _, f := range ci.Table().Fields() {
		columns = append(columns, f.Name())
	}
	return mssql.CopyIn(ci.Table().Quoted(), opt,
		columns...), nil
}

//Agrs 通过多条记录 records生成批量copy in参数
func (ci *CopyInParam) Agrs(records []element.Record) (valuers []interface{}, err error) {
	for _, r := range records {
		for fi, f := range ci.Table().Fields() {
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
