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
	"github.com/Breeze0806/go-etl/element"
	"github.com/Breeze0806/go-etl/storage/database"
)

//PageParam 页查询参数
type PageParam struct {
	Start Offset             //页面开始位移
	End   Offset             //页面结束位移
	Param database.Parameter //页面查询语句
}

//Page 页查询结果
type Page struct {
	Min     Offset                //页查询结果最小位移
	Max     Offset                //页查询结果最大位移
	Records map[string]PageRecord //页查询结果映射
}

//PageRecord 页记录
type PageRecord struct {
	Now    Offset         //当前页记录位移
	Record element.Record //页记录查询结果
}
