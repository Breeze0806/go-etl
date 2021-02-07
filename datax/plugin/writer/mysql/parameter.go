package mysql

import (
	"bytes"
	"errors"

	"github.com/Breeze0806/go-etl/element"
	"github.com/Breeze0806/go-etl/storage/database"
)

type parameter struct {
	*database.BaseParam

	paramConfig *paramConfig
}

func newParameter(paramConfig *paramConfig, execer Execer) *parameter {
	p := &parameter{
		BaseParam: database.NewBaseParam(execer.Table(database.NewBaseTable(
			paramConfig.Connection.Table.Db, "", paramConfig.Connection.Table.Name)), nil),
		paramConfig: paramConfig,
	}

	return p
}

type tableParam struct {
	*parameter
}

func newTableParam(p *parameter) *tableParam {
	return &tableParam{
		parameter: p,
	}
}

func (t *tableParam) Query(_ []element.Record) (string, error) {
	buf := bytes.NewBufferString("select ")
	if len(t.paramConfig.Column) == 0 {
		return "", errors.New("column is empty")
	}
	for i, v := range t.paramConfig.Column {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(v)
	}
	buf.WriteString(" from ")
	buf.WriteString(t.Table().Quoted())
	buf.WriteString(" where 1 = 2")
	return buf.String(), nil
}

func (t *tableParam) Agrs(_ []element.Record) ([]interface{}, error) {
	return nil, nil
}
