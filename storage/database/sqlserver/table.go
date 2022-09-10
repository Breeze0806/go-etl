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

	"github.com/Breeze0806/go-etl/storage/database"
)

//Table mssql表
type Table struct {
	*database.BaseTable
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

//ExecParam 获取执行参数，其中replace into的参数方式以及被注册
func (t *Table) ExecParam(mode string, txOpts *sql.TxOptions) (database.Parameter, bool) {
	return nil, false
}
