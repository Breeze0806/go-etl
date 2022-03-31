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
	"fmt"

	"github.com/Breeze0806/go-etl/element"
	"github.com/Breeze0806/go-etl/storage/database"
)

//Offset 位移
type Offset interface {
	fmt.Stringer
	//比较,1 代表 > right, 0 代表 == right, -1 代表 < right，当比较失败时会返回错误
	Cmp(right Offset) (int, error)
	PrimaryKey() database.PrimaryKey     //表主键
	PrimaryKeyColumns() []element.Column //主键列值
}

//BaseOffset 基本位移，用于写自己的位移主键
type BaseOffset struct {
	columns []element.Column
	key     database.PrimaryKey
}

//NewBaseOffset 根据表主键key 和主键列columns 生成基本位移
func NewBaseOffset(key database.PrimaryKey, columns []element.Column) *BaseOffset {
	return &BaseOffset{
		columns: columns,
		key:     key,
	}
}

//Cmp 与位移right比较大小， 当列大小不一致时会返回错误
//1代表 b > offset, 0 代表 b == offset, -1代表 b < right
func (b *BaseOffset) Cmp(right Offset) (cmp int, err error) {
	columns := right.PrimaryKeyColumns()
	if len(columns) != len(b.columns) {
		return 0, fmt.Errorf("the length of columns is not equaled(%v,%v)", b, right)
	}
	for i := range b.columns {
		cmp, err = b.columns[i].Cmp(columns[i])
		if err != nil {
			return 0, err
		}
		if cmp != 0 {
			return cmp, nil
		}
	}
	return 0, nil
}

//PrimaryKey 表主键
func (b *BaseOffset) PrimaryKey() database.PrimaryKey {
	return b.key
}

//PrimaryKeyColumns 表主键列值
func (b *BaseOffset) PrimaryKeyColumns() []element.Column {
	return b.columns
}
