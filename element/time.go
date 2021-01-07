package element

import (
	"fmt"
	"time"
)

var defaultTimeFormat = time.RFC3339Nano

type TimeDecoder interface {
	TimeDecode(t time.Time) (interface{}, error)
}

type TimeEncoder interface {
	TimeEncode(i interface{}) (time.Time, error)
}

type StringTimeEncoder struct {
	layout string
}

func NewStringTimeEncoder(layout string) TimeEncoder {
	return &StringTimeEncoder{
		layout: layout,
	}
}

func (e *StringTimeEncoder) TimeEncode(i interface{}) (time.Time, error) {
	s, ok := i.(string)
	if !ok {
		return time.Time{}, fmt.Errorf("%v is %T, not string", i, i)
	}
	return time.Parse(e.layout, s)
}

type StringTimeDecoder struct {
	layout string
}

func NewStringTimeDecoder(layout string) TimeDecoder {
	return &StringTimeDecoder{
		layout: layout,
	}
}

func (d *StringTimeDecoder) TimeDecode(t time.Time) (interface{}, error) {
	return t.Format(d.layout), nil
}
