package postgres

import (
	"database/sql"
	"database/sql/driver"
	"fmt"

	"github.com/Breeze0806/go-etl/element"
	"github.com/Breeze0806/go-etl/storage/database"
	"github.com/lib/pq"
)

//Table postgres表
type Table struct {
	*database.BaseTable
}

//NewTable 创建postgres表，注意此时BaseTable中的schema参数为架构名，instance为数据库名，而name是表明
func NewTable(b *database.BaseTable) *Table {
	return &Table{
		BaseTable: b,
	}
}

//Quoted 引用全名
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

//ExecParam 获取执行参数，其中copy in的参数方式以及被注册
func (t *Table) ExecParam(mode string, txOpts *sql.TxOptions) (database.Parameter, bool) {
	switch mode {
	case "copyIn":
		return NewCopyInParam(t, txOpts), true
	}
	return nil, false
}

//CopyInParam copy in 参数
type CopyInParam struct {
	*database.BaseParam
}

//NewCopyInParam  通过表table和事务参数txOps插入参数
func NewCopyInParam(t database.Table, txOps *sql.TxOptions) *CopyInParam {
	return &CopyInParam{
		BaseParam: database.NewBaseParam(t, txOps),
	}
}

//Query 批量copy in插入sql语句
func (ci *CopyInParam) Query(_ []element.Record) (query string, err error) {
	var columns []string
	for _, f := range ci.Table().Fields() {
		columns = append(columns, f.Name())
	}
	return pq.CopyInSchema(ci.Table().Schema(), ci.Table().Name(),
		columns...), nil
}

//Agrs 通过多条记录 records生成批量copy in参数
func (ci *CopyInParam) Agrs(records []element.Record) (valuers []interface{}, err error) {
	for _, r := range records {
		for _, f := range ci.Table().Fields() {
			var c element.Column
			if c, err = r.GetByName(f.Name()); err != nil {
				return nil, fmt.Errorf("GetByName(%v) err: %v", f.Name(), err)
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
