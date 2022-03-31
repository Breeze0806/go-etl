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
	"fmt"
	"sync"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/element"
)

// Creater 创建输出流的创建器
type Creater interface {
	Create(filename string) (stream OutStream, err error) //创建名为filename的输出流
}

// OutStream 输出流
type OutStream interface {
	Writer(conf *config.JSON) (writer StreamWriter, err error) //创建写入器
	Close() (err error)                                        //关闭输出流
}

// StreamWriter 输出流写入器
type StreamWriter interface {
	Write(record element.Record) (err error) //写入记录
	Flush() (err error)                      //刷新至文件
	Close() (err error)                      //关闭输出流写入器
}

// RegisterCreater 通过创建器名称name注册输出流创建器creater
func RegisterCreater(name string, creater Creater) {
	if err := creaters.register(name, creater); err != nil {
		panic(err)
	}
}

// OutStreamer 输出流包装
type OutStreamer struct {
	stream OutStream
}

// NewOutStreamer 通过creater名称name的输出流包装，并打开名为filename的输出流
func NewOutStreamer(name string, filename string) (streamer *OutStreamer, err error) {
	creater, ok := creaters.creater(name)
	if !ok {
		err = fmt.Errorf("creater %v does not exist", name)
		return nil, err
	}
	streamer = &OutStreamer{}
	if streamer.stream, err = creater.Create(filename); err != nil {
		return nil, fmt.Errorf("create fail. err : %v", err)
	}
	return
}

// Writer 通过配置conf创建流写入器
func (s *OutStreamer) Writer(conf *config.JSON) (StreamWriter, error) {
	return s.stream.Writer(conf)
}

// Close 关闭写入包装
func (s *OutStreamer) Close() error {
	return s.stream.Close()
}

var creaters = &createrMap{
	creaters: make(map[string]Creater),
}

type createrMap struct {
	sync.RWMutex
	creaters map[string]Creater
}

func (o *createrMap) register(name string, creater Creater) error {
	if creater == nil {
		return fmt.Errorf("creater %v is nil", name)
	}

	o.Lock()
	defer o.Unlock()
	if _, ok := o.creaters[name]; ok {
		return fmt.Errorf("creater %v exists", name)
	}

	o.creaters[name] = creater
	return nil
}

func (o *createrMap) creater(name string) (creater Creater, ok bool) {
	o.RLock()
	defer o.RUnlock()
	creater, ok = o.creaters[name]
	return
}
