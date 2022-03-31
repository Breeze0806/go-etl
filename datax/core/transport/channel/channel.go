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

package channel

import (
	"context"

	"github.com/Breeze0806/go-etl/element"
)

//Channel 通道
type Channel struct {
	records *element.RecordChan
}

//NewChannel 创建通道
func NewChannel(ctx context.Context) (*Channel, error) {
	return &Channel{
		records: element.NewRecordChan(ctx),
	}, nil
}

//Size 通道记录大小
func (c *Channel) Size() int {
	return c.records.Buffered()
}

//IsEmpty 通道是否为空
func (c *Channel) IsEmpty() bool {
	return c.Size() == 0
}

//Push 将记录r加入通道
func (c *Channel) Push(r element.Record) int {
	return c.records.PushBack(r)
}

//Pop 将记录弹出，当通道中不存在记录，就会返回false
func (c *Channel) Pop() (r element.Record, ok bool) {
	return c.records.PopFront()
}

//PushAll 通过fetchRecord函数加入多条记录
func (c *Channel) PushAll(fetchRecord func() (element.Record, error)) error {
	return c.records.PushBackAll(fetchRecord)
}

//PopAll 通过onRecord函数弹出多条记录
func (c *Channel) PopAll(onRecord func(element.Record) error) error {
	return c.records.PopFrontAll(onRecord)
}

//Close 关闭
func (c *Channel) Close() {
	c.records.Close()
}

//PushTerminate 加入终止记录
func (c *Channel) PushTerminate() int {
	return c.Push(element.GetTerminateRecord())
}
