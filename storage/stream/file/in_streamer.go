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

package file

import (
	"context"
	"fmt"
	"sync"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/element"
	"github.com/pingcap/errors"
)

// FetchHandler 获取记录句柄
type FetchHandler interface {
	OnRecord(element.Record) error         //处理记录
	CreateRecord() (element.Record, error) //创建空记录
}

// Opener 用于打开一个输入流的打开器
type Opener interface {
	Open(filename string) (stream InStream, err error) //打开文件名filename的输入流
}

// InStream 输入流
type InStream interface {
	Rows(conf *config.JSON) (rows Rows, err error) //获取行读取器
	Close() (err error)                            //关闭输入流
}

// Rows 行读取器
type Rows interface {
	Next() bool                                  //获取下一行，如果没有返回false，有返回true
	Scan() (columns []element.Column, err error) //扫描出每一行的列
	Error() error                                //获取下一行的错误
	Close() error                                //关闭行读取器
}

// RegisterOpener 通过打开器名称name注册输入流打开器opener
func RegisterOpener(name string, opener Opener) {
	if err := openers.register(name, opener); err != nil {
		panic(err)
	}
}

// UnregisterAllOpener 注销所有文件打开器
func UnregisterAllOpener() {
	openers.unregisterAll()
}

// InStreamer 输入流包装
type InStreamer struct {
	stream InStream
}

// NewInStreamer 通过opener名称name的输入流打开器，并打开名为filename的输入流
func NewInStreamer(name string, filename string) (streamer *InStreamer, err error) {
	opener, ok := openers.opener(name)
	if !ok {
		err = errors.Errorf("opener %v does not exist", name)
		return
	}
	streamer = &InStreamer{}
	if streamer.stream, err = opener.Open(filename); err != nil {
		return nil, errors.Wrapf(err, "open(%v) fail", filename)
	}
	return
}

// Read 使用获取记录句柄handler，传入上下文ctx和配置文件conf获取对应数据
func (s *InStreamer) Read(ctx context.Context, conf *config.JSON, handler FetchHandler) (err error) {
	var rows Rows
	rows, err = s.stream.Rows(conf)
	if err != nil {
		return errors.Wrapf(err, "rows fail. config: %v", conf)
	}
	defer rows.Close()
	for rows.Next() {
		var columns []element.Column
		select {
		case <-ctx.Done():
			return nil
		default:
		}
		if columns, err = rows.Scan(); err != nil {
			return errors.Wrapf(err, "Scan fail")
		}
		if len(columns) > 0 {
			var r element.Record

			if r, err = handler.CreateRecord(); err != nil {
				return errors.Wrapf(err, "CreateRecord fail")
			}

			for _, v := range columns {
				if err = r.Add(v); err != nil {
					return errors.Wrapf(err, "Add fail")
				}
			}
			if err = handler.OnRecord(r); err != nil {
				return errors.Wrapf(err, "OnRecord fail")
			}
		}
	}
	if err = rows.Error(); err != nil {
		return errors.Wrapf(err, "Error")
	}
	return
}

// Close 关闭输入流
func (s *InStreamer) Close() error {
	return s.stream.Close()
}

var openers = &openerMap{
	openers: make(map[string]Opener),
}

type openerMap struct {
	sync.RWMutex
	openers map[string]Opener
}

func (o *openerMap) register(name string, opener Opener) error {
	if opener == nil {
		return fmt.Errorf("opener %v is nil", name)
	}

	o.Lock()
	defer o.Unlock()
	if _, ok := o.openers[name]; ok {
		return fmt.Errorf("opener %v exists", name)
	}

	o.openers[name] = opener
	return nil
}

func (o *openerMap) opener(name string) (opener Opener, ok bool) {
	o.RLock()
	defer o.RUnlock()
	opener, ok = o.openers[name]
	return
}

func (o *openerMap) unregisterAll() {
	o.Lock()
	defer o.Unlock()
	o.openers = make(map[string]Opener)
}
