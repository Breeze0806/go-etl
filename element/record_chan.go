package element

import "sync"

//RecordChan 记录通道 修复内存溢出
type RecordChan struct {
	lock sync.Mutex
	ch   chan Record

	closed bool
}

const defaultRequestChanBuffer = 1024

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
		ch: make(chan Record, n),
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
	c.ch <- r
	return c.Buffered()
}

//PopFront 在头部弹出记录r，并且返回是否还有值
func (c *RecordChan) PopFront() (Record, bool) {
	r, ok := <-c.ch
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
