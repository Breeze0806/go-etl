package channel

import "github.com/Breeze0806/go-etl/datax/common/element"

type Channel struct {
	records *RecordChan
}

func NewChannel() (*Channel, error) {
	return &Channel{
		records: NewRecordChan(),
	}, nil
}

func (c *Channel) Size() int {
	return c.records.Buffered()
}

func (c *Channel) IsEmpty() bool {
	return c.Size() == 0
}

func (c *Channel) Push(r element.Record) int {
	return c.records.PushBack(r)
}

func (c *Channel) Pop() (r element.Record, ok bool) {
	return c.records.PopFront()
}

func (c *Channel) PushAll(fetchRecord func() (element.Record, error)) error {
	return c.records.PushBackAll(fetchRecord)
}

func (c *Channel) PopAll(onRecord func(element.Record) error) error {
	return c.records.PopFrontAll(onRecord)
}

func (c *Channel) Close() {
	c.records.Close()
}

func (c *Channel) PushTerminate() int {
	return c.Push(element.GetTerminateRecord())
}
