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

// DefaultTimeFormat 默认时间格式
var DefaultTimeFormat = "2006-01-02 15:04:05.999999999Z07:00"

// TimeDecoder 时间解码器
type TimeDecoder interface {
	TimeDecode(t time.Time) (interface{}, error)
	Layout() string
}

// TimeEncoder 时间编码器
type TimeEncoder interface {
	TimeEncode(i interface{}) (time.Time, error)
}

// StringTimeEncoder 字符串时间编码器
type StringTimeEncoder struct {
	layout string //go时间格式
}

// NewStringTimeEncoder 根据go时间格式layout的字符串时间编码器
func NewStringTimeEncoder(layout string) TimeEncoder {
	return &StringTimeEncoder{
		layout: layout,
	}
}

// TimeEncode 编码成时间，若i不是string或者不是layout格式，会报错
func (e *StringTimeEncoder) TimeEncode(i interface{}) (time.Time, error) {
	s, ok := i.(string)
	if !ok {
		return time.Time{}, fmt.Errorf("%v is %T, not string", i, i)
	}
	return time.Parse(e.layout[:len(s)], s)
}

// StringTimeDecoder 字符串时间编码器
type StringTimeDecoder struct {
	layout string //go时间格式
}

// NewStringTimeDecoder 根据go时间格式layout的字符串时间编码器
func NewStringTimeDecoder(layout string) TimeDecoder {
	return &StringTimeDecoder{
		layout: layout,
	}
}

// TimeDecode 根据go时间格式layout的字符串时间编码成string
func (d *StringTimeDecoder) TimeDecode(t time.Time) (interface{}, error) {
	return t.Format(d.layout), nil
}

// Layout 时间格式
func (d *StringTimeDecoder) Layout() string {
	return d.layout
}
