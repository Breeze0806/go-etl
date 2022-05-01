package db2

import (
	"database/sql"

	"github.com/Breeze0806/go-etl/storage/database"
)

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
	return Quoted(t.Schema()) + "." + Quoted(t.Name())
}

func (t *Table) String() string {
	return t.Quoted()
}

//AddField 新增列
func (t *Table) AddField(baseField *database.BaseField) {
	t.AppendField(NewField(baseField))
}

//ExecParam 获取执行参数
func (t *Table) ExecParam(mode string, txOpts *sql.TxOptions) (database.Parameter, bool) {
	return nil, false
}
