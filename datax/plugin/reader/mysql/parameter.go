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

func newParameter(paramConfig *paramConfig, querier Querier) *parameter {
	p := &parameter{
		BaseParam: database.NewBaseParam(querier.Table(database.NewBaseTable(
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

type queryParam struct {
	*parameter
}

func newQueryParam(p *parameter) *queryParam {
	return &queryParam{
		parameter: p,
	}
}

func (q *queryParam) Query(_ []element.Record) (string, error) {
	buf := bytes.NewBufferString("select ")
	if len(q.paramConfig.Column) == 0 {
		return "", errors.New("column is empty")
	}
	for i, v := range q.paramConfig.Column {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(v)
	}
	buf.WriteString(" from ")
	buf.WriteString(q.Table().Quoted())
	if q.paramConfig.Where != "" {
		buf.WriteString(" where ")
		buf.WriteString(q.paramConfig.Where)
	}
	return buf.String(), nil
}

func (q *queryParam) Agrs(_ []element.Record) ([]interface{}, error) {
	return nil, nil
}
