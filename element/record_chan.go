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

package element

import (
	"context"
	"sync"
)

//RecordChan 记录通道 修复内存溢出
type RecordChan struct {
	lock   sync.Mutex
	ch     chan Record
	ctx    context.Context
	closed bool
}

const defaultRequestChanBuffer = 1024

//NewRecordChan 创建记录通道
func NewRecordChan(ctx context.Context) *RecordChan {
	return NewRecordChanBuffer(ctx, 0)
}

//NewRecordChanBuffer 创建容量n的记录通道
func NewRecordChanBuffer(ctx context.Context, n int) *RecordChan {
	if n <= 0 {
		n = defaultRequestChanBuffer
	}
	var ch = &RecordChan{
		ctx: ctx,
		ch:  make(chan Record, n),
	}
	return ch
}

//Close 关闭
func (c *RecordChan) Close() {
	c.lock.Lock()
	if !c.closed {
		c.closed = true
		close(c.ch)
	}
	c.lock.Unlock()
}

//Buffered 记录通道内的元素数量
func (c *RecordChan) Buffered() int {
	return len(c.ch)
}

//PushBack 在尾部追加记录r，并且返回队列大小
func (c *RecordChan) PushBack(r Record) int {
	select {
	case c.ch <- r:
	case <-c.ctx.Done():
	}
	return c.Buffered()
}

//PopFront 在头部弹出记录r，并且返回是否还有值
func (c *RecordChan) PopFront() (r Record, ok bool) {
	select {
	case r, ok = <-c.ch:
	case <-c.ctx.Done():
		r, ok = nil, false
	}
	return r, ok
}

//PushBackAll 通过函数fetchRecord获取多个记录，在尾部追加
func (c *RecordChan) PushBackAll(fetchRecord func() (Record, error)) error {
	for {
		r, err := fetchRecord()
		if err != nil {
			return err
		}
		c.PushBack(r)
	}
}

//PopFrontAll 通过函数onRecord从头部弹出所有记录
func (c *RecordChan) PopFrontAll(onRecord func(Record) error) error {
	for {
		r, ok := c.PopFront()
		if ok {
			if err := onRecord(r); err != nil {
				return err
			}
		} else {
			return nil
		}
	}
}
