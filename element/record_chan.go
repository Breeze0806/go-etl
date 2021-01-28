package element

import (
	"sync"
)

//RecordChan 记录通道
type RecordChan struct {
	lock sync.Mutex
	cond *sync.Cond

	data []Record
	buff []Record

	waits  int
	closed bool
}

const defaultRequestChanBuffer = 128

//NewRecordChan 创建记录通道
func NewRecordChan() *RecordChan {
	return NewRecordChanBuffer(0)
}

//NewRecordChanBuffer 创建容量n的记录通道
func NewRecordChanBuffer(n int) *RecordChan {
	if n <= 0 {
		n = defaultRequestChanBuffer
	}
	var ch = &RecordChan{
		buff: make([]Record, n),
	}
	ch.cond = sync.NewCond(&ch.lock)
	return ch
}

//Close 关闭
func (c *RecordChan) Close() {
	c.lock.Lock()
	if !c.closed {
		c.closed = true
		c.cond.Broadcast()
	}
	c.lock.Unlock()
}

//Buffered 记录通道内的元素数量
func (c *RecordChan) Buffered() int {
	c.lock.Lock()
	n := len(c.data)
	c.lock.Unlock()
	return n
}

//PushBack 在尾部追加记录r，并且返回队列大小
func (c *RecordChan) PushBack(r Record) int {
	c.lock.Lock()
	n := c.lockedPushBack(r)
	c.lock.Unlock()
	return n
}

//PopFront 在头部弹出记录r，并且返回是否还有值
func (c *RecordChan) PopFront() (Record, bool) {
	c.lock.Lock()
	r, ok := c.lockedPopFront()
	c.lock.Unlock()
	return r, ok
}

func (c *RecordChan) lockedPushBack(r Record) int {
	if c.closed {
		panic("send on closed chan")
	}
	if c.waits != 0 {
		c.cond.Signal()
	}
	c.data = append(c.data, r)
	return len(c.data)
}

func (c *RecordChan) lockedPopFront() (Record, bool) {
	for len(c.data) == 0 {
		if c.closed {
			return nil, false
		}
		c.data = c.buff[:0]
		c.waits++
		c.cond.Wait()
		c.waits--
	}
	var r = c.data[0]
	c.data, c.data[0] = c.data[1:], nil
	return r, true
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
