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

package dbms

import (
	"bytes"
	"database/sql"

	"github.com/Breeze0806/go-etl/element"
	"github.com/Breeze0806/go-etl/storage/database"
	"github.com/pingcap/errors"
)

// TableParamConfig Table parameter configuration
type TableParamConfig interface {
	GetColumns() []Column              // Get column information
	GetBaseTable() *database.BaseTable // Get table information
}

// TableParamTable Get the table of the corresponding database through table parameters
type TableParamTable interface {
	Table(*database.BaseTable) database.Table // Get the table of the corresponding database through table parameters
}

// TableParam Table parameters
type TableParam struct {
	*database.BaseParam

	Config TableParamConfig
}

// NewTableParam Get table parameter configuration config, get table parameters through table parameters of the corresponding database table and transaction options opts
func NewTableParam(config TableParamConfig, table TableParamTable, opts *sql.TxOptions) *TableParam {
	return &TableParam{
		BaseParam: database.NewBaseParam(table.Table(config.GetBaseTable()), opts),

		Config: config,
	}
}

// Query Get the query statement
func (t *TableParam) Query(_ []element.Record) (string, error) {
	buf := bytes.NewBufferString("select ")
	if len(t.Config.GetColumns()) == 0 {
		return "", errors.NewNoStackError("column is empty")
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

// Agrs Get query parameters
func (t *TableParam) Agrs(_ []element.Record) ([]any, error) {
	return nil, nil
}

// QueryParam Query parameters
type QueryParam struct {
	*database.BaseParam

	Config Config
}

// NewQueryParam Get query parameters through relational database input configuration config, corresponding database table table, and transaction options opts
func NewQueryParam(config Config, table database.Table, opts *sql.TxOptions) *QueryParam {
	return &QueryParam{
		BaseParam: database.NewBaseParam(table, opts),

		Config: config,
	}
}

// Query Get the query statement
func (q *QueryParam) Query(_ []element.Record) (string, error) {
	if len(q.Config.GetQuerySQL()) > 1 {
		return "", errors.NewNoStackError("too much querySQL")
	}

	if len(q.Config.GetQuerySQL()) == 1 {
		return q.Config.GetQuerySQL()[0], nil
	}

	buf := bytes.NewBufferString("select ")
	if len(q.Table().Fields()) == 0 {
		return "", errors.NewNoStackError("column is empty")
	}

	canUseConfig := false
	if len(q.Table().Fields()) == len(q.Config.GetColumns()) {
		canUseConfig = true
	}

	for i, v := range q.Table().Fields() {
		if i > 0 {
			buf.WriteString(",")
		}

		col := v.Select()
		if canUseConfig && v.Name() != q.Config.GetColumns()[i].GetName() {
			col = q.Config.GetColumns()[i].GetName()
		}
		buf.WriteString(col)
	}
	buf.WriteString(" from ")
	buf.WriteString(q.Table().Quoted())
	if q.Config.GetWhere() != "" {
		buf.WriteString(" where ")
		buf.WriteString(q.Config.GetWhere())
	}
	return buf.String(), nil
}

// Agrs Get query parameters
func (q *QueryParam) Agrs(_ []element.Record) (a []any, err error) {
	if len(q.Config.GetQuerySQL()) > 0 {
		return nil, nil
	}

	if q.Config.GetSplitConfig().Key != "" {
		for _, v := range q.Table().Fields() {
			if q.Config.GetSplitConfig().Key == v.Name() {
				var left, right element.Column
				if left, err = q.Config.GetSplitConfig().Range.leftColumn(v.Name()); err != nil {
					return
				}
				if right, err = q.Config.GetSplitConfig().Range.rightColumn(v.Name()); err != nil {
					return
				}
				var li, ri any
				if li, err = v.Valuer(left).Value(); err != nil {
					return
				}
				if ri, err = v.Valuer(right).Value(); err != nil {
					return
				}
				a = append(a, li, ri)
				return
			}
		}
	}
	return nil, nil
}

// SplitParam Splitting parameters
type SplitParam struct {
	*database.BaseParam

	Config Config
}

// NewSplitParam Get table parameter configuration config, get split table parameters through table parameters of the corresponding database table and transaction options opts
func NewSplitParam(config Config, table TableParamTable, opts *sql.TxOptions) *SplitParam {
	return &SplitParam{
		BaseParam: database.NewBaseParam(table.Table(config.GetBaseTable()), opts),

		Config: config,
	}
}

// Query Get the query statement
func (s *SplitParam) Query(_ []element.Record) (string, error) {
	buf := bytes.NewBufferString("select ")

	buf.WriteString(s.Config.GetSplitConfig().Key)
	buf.WriteString(" from ")
	buf.WriteString(s.Table().Quoted())
	buf.WriteString(" where 1 = 2")

	return buf.String(), nil
}

// Agrs Get query parameters
func (s *SplitParam) Agrs(_ []element.Record) ([]any, error) {
	return nil, nil
}

// MinParam Minimum value parameter
type MinParam struct {
	*database.BaseParam

	Config Config
}

// NewMinParam Get the minimum value parameter through relational database input configuration config, corresponding database table table, and transaction options opts
func NewMinParam(config Config, table database.Table, opts *sql.TxOptions) *MinParam {
	return &MinParam{
		BaseParam: database.NewBaseParam(table, opts),

		Config: config,
	}
}

// Query Get the query statement
func (m *MinParam) Query(_ []element.Record) (string, error) {
	buf := bytes.NewBufferString("select min(")
	buf.WriteString(m.Config.GetSplitConfig().Key)

	buf.WriteString(") from ")
	buf.WriteString(m.Table().Quoted())
	if m.Config.GetWhere() != "" {
		buf.WriteString(" where ")
		buf.WriteString(m.Config.GetWhere())
	}
	return buf.String(), nil
}

// Agrs Get query parameters
func (m *MinParam) Agrs(_ []element.Record) ([]any, error) {
	return nil, nil
}

// MaxParam Maximum value parameter
type MaxParam struct {
	*database.BaseParam

	Config Config
}

// NewMaxParam Get query parameters through relational database input configuration config, corresponding database table table, and transaction options opts
func NewMaxParam(config Config, table database.Table, opts *sql.TxOptions) *MaxParam {
	return &MaxParam{
		BaseParam: database.NewBaseParam(table, opts),

		Config: config,
	}
}

// Query Get the query statement
func (m *MaxParam) Query(_ []element.Record) (string, error) {
	buf := bytes.NewBufferString("select max(")
	buf.WriteString(m.Config.GetSplitConfig().Key)

	buf.WriteString(") from ")
	buf.WriteString(m.Table().Quoted())
	if m.Config.GetWhere() != "" {
		buf.WriteString(" where ")
		buf.WriteString(m.Config.GetWhere())
	}
	return buf.String(), nil
}

// Agrs Get query parameters
func (m *MaxParam) Agrs(_ []element.Record) ([]any, error) {
	return nil, nil
}
