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

// RecordChan: Record channel. Fixes memory overflow issues.
type RecordChan struct {
	lock   sync.Mutex
	ch     chan Record
	ctx    context.Context
	closed bool
}

const defaultRequestChanBuffer = 1024

// NewRecordChan: Create a new record channel.
func NewRecordChan(ctx context.Context) *RecordChan {
	return NewRecordChanBuffer(ctx, 0)
}

// NewRecordChanBuffer: Create a new record channel with a capacity of n.
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

// Close: Close the channel.
func (c *RecordChan) Close() {
	c.lock.Lock()
	if !c.closed {
		c.closed = true
		close(c.ch)
	}
	c.lock.Unlock()
}

// Buffered: Number of elements currently buffered in the record channel.
func (c *RecordChan) Buffered() int {
	return len(c.ch)
}

// PushBack: Append a record r to the end of the channel and return the current size of the queue.
func (c *RecordChan) PushBack(r Record) int {
	select {
	case c.ch <- r:
	case <-c.ctx.Done():
	}
	return c.Buffered()
}

// PopFront: Remove and return the record at the front of the channel, and indicate whether there are more values remaining.
func (c *RecordChan) PopFront() (r Record, ok bool) {
	select {
	case r, ok = <-c.ch:
	case <-c.ctx.Done():
		r, ok = nil, false
	}
	return r, ok
}

// PushBackAll: Append multiple records obtained through the fetchRecord function to the end of the channel.
func (c *RecordChan) PushBackAll(fetchRecord func() (Record, error)) error {
	for {
		r, err := fetchRecord()
		if err != nil {
			return err
		}
		c.PushBack(r)
	}
}

// PopFrontAll: Remove and return all records from the front of the channel using the onRecord function.
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
