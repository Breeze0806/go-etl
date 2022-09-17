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

package rdbm

import (
	"context"
	"database/sql"

	"github.com/Breeze0806/go-etl/storage/database"
)

//Execer 执行器
type Execer interface {
	//通过基础表信息获取具体表
	Table(*database.BaseTable) database.Table
	//检测连通性
	PingContext(ctx context.Context) error
	//通过query查询语句进行查询
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	//通过query查询语句进行查询
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	//通过参数param获取具体表
	FetchTableWithParam(ctx context.Context, param database.Parameter) (database.Table, error)
	//批量执行
	BatchExec(ctx context.Context, opts *database.ParameterOptions) (err error)
	//prepare/exec批量执行
	BatchExecStmt(ctx context.Context, opts *database.ParameterOptions) (err error)
	//事务批量执行
	BatchExecWithTx(ctx context.Context, opts *database.ParameterOptions) (err error)
	//事务prepare/exec批量执行
	BatchExecStmtWithTx(ctx context.Context, opts *database.ParameterOptions) (err error)
	//关闭
	Close() error
}
