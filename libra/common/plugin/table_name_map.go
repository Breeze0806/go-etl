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

import "github.com/Breeze0806/go-etl/storage/database"

//TableNameMap 表名映射
type TableNameMap interface {
	//通过主数据库表名获取master获取对应的从数据库表名slave,有错误时返回err
	SlaveTableName(master *database.BaseTable) (slave *database.BaseTable, err error)
	Close() error
}

//TableNameMapMaker 表名映射生成器
type TableNameMapMaker interface {
	TableNameMap() TableNameMap
}
