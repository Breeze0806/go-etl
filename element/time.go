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
	"fmt"
	"time"
)

// DefaultTimeFormat - Default time format
var DefaultTimeFormat = "2006-01-02 15:04:05.999999999Z07:00"

// TimeDecoder - Time decoder
type TimeDecoder interface {
	TimeDecode(t time.Time) (interface{}, error)
	Layout() string
}

// TimeEncoder - Time encoder
type TimeEncoder interface {
	TimeEncode(i interface{}) (time.Time, error)
}

// StringTimeEncoder - String time encoder
type StringTimeEncoder struct {
	layout string // go time format
}

// NewStringTimeEncoder - A string time encoder based on the layout of the go time format
func NewStringTimeEncoder(layout string) TimeEncoder {
	return &StringTimeEncoder{
		layout: layout,
	}
}

// TimeEncode - Encode to time. If 'i' is not a string or not in the layout format, an error will be reported.
func (e *StringTimeEncoder) TimeEncode(i interface{}) (time.Time, error) {
	s, ok := i.(string)
	if !ok {
		return time.Time{}, fmt.Errorf("%v is %T, not string", i, i)
	}
	return time.Parse(e.layout[:len(s)], s)
}

// StringTimeDecoder - String time decoder
type StringTimeDecoder struct {
	layout string // go time format
}

// NewStringTimeDecoder - A string time decoder based on the layout of the go time format
func NewStringTimeDecoder(layout string) TimeDecoder {
	return &StringTimeDecoder{
		layout: layout,
	}
}

// TimeDecode - Decode a string time based on the layout of the go time format into a string
func (d *StringTimeDecoder) TimeDecode(t time.Time) (interface{}, error) {
	return t.Format(d.layout), nil
}

// Layout - Time format
func (d *StringTimeDecoder) Layout() string {
	return d.layout
}
