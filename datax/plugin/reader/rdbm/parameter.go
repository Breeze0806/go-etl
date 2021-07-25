package rdbm

import (
	"bytes"
	"database/sql"
	"errors"

	"github.com/Breeze0806/go-etl/element"
	"github.com/Breeze0806/go-etl/storage/database"
)

type TableParamConfig interface {
	GetColumns() []Column
	GetBaseTable() *database.BaseTable
}

type TableParamTable interface {
	Table(*database.BaseTable) database.Table
}

type TableParam struct {
	*database.BaseParam

	Config TableParamConfig
}

func NewTableParam(config TableParamConfig, table TableParamTable, opts *sql.TxOptions) *TableParam {
	return &TableParam{
		BaseParam: database.NewBaseParam(table.Table(config.GetBaseTable()), opts),

		Config: config,
	}
}

func (t *TableParam) Query(_ []element.Record) (string, error) {
	buf := bytes.NewBufferString("select ")
	if len(t.Config.GetColumns()) == 0 {
		return "", errors.New("column is empty")
	}
	for i, v := range t.Config.GetColumns() {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(v.GetName())
	}
	buf.WriteString(" from ")
	buf.WriteString(t.Table().Quoted())
	buf.WriteString(" where 1 = 2")
	return buf.String(), nil
}

func (t *TableParam) Agrs(_ []element.Record) ([]interface{}, error) {
	return nil, nil
}

type QueryParam struct {
	*database.BaseParam

	Config Config
}

func NewQueryParam(config Config, table database.Table, opts *sql.TxOptions) *QueryParam {
	return &QueryParam{
		BaseParam: database.NewBaseParam(table, opts),

		Config: config,
	}
}

func (q *QueryParam) Query(_ []element.Record) (string, error) {
	buf := bytes.NewBufferString("select ")
	if len(q.Table().Fields()) == 0 {
		return "", errors.New("column is empty")
	}
	for i, v := range q.Table().Fields() {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(v.Quoted())
	}
	buf.WriteString(" from ")
	buf.WriteString(q.Table().Quoted())
	if q.Config.GetWhere() != "" {
		buf.WriteString(" where ")
		buf.WriteString(q.Config.GetWhere())
	}
	return buf.String(), nil
}

func (q *QueryParam) Agrs(_ []element.Record) ([]interface{}, error) {
	return nil, nil
}
