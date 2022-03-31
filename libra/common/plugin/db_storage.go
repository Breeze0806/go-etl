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

package plugin

import (
	"context"

	"github.com/Breeze0806/go-etl/storage/database"
)

//StorageTable 数据库表
type StorageTable interface {
	Table() database.Table        //数据库表
	ReadPage(ctx context.Context, //读取一页数据库
		param PageParam) (Page, error)
	Close() error //关闭
}

//DBStorage 数据库
type DBStorage interface {
	AllTable(ctx context.Context) ([]*database.BaseTable, error) //查询所有表
	MasterTable(*database.BaseTable) (MasterTable, error)        //主数据库表
	SlaveTable(*database.BaseTable) (SlaveTable, error)          //从数据库表
}
