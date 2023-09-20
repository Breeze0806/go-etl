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

package loader

import (
	"context"
	"fmt"
	"sync"

	"github.com/Breeze0806/go-etl/datax/common/plugin"
	"github.com/Breeze0806/go-etl/datax/common/spi"
	"github.com/Breeze0806/go-etl/datax/common/spi/reader"
	"github.com/Breeze0806/go-etl/datax/common/spi/writer"
)

var _centor = &centor{
	readers: make(map[string]spi.Reader),
	writers: make(map[string]spi.Writer),
}

// RegisterReader 注册名字为name的读取器reader
// 当name重复，reader为空或者reader的task和job为空会panic
func RegisterReader(name string, reader spi.Reader) {
	if err := _centor.registerReader(name, reader); err != nil {
		panic(err)
	}
}

// RegisterWriter 注册名字为name的写入器writer
// 当name重复，writer为空或者writer的task和job为空会panic
func RegisterWriter(name string, writer spi.Writer) {
	if err := _centor.registerWriter(name, writer); err != nil {
		panic(err)
	}
}

// UnregisterReaders 注销所有读取器
func UnregisterReaders() {
	_centor.unregisterReaders()
}

// UnregisterWriters 注销所有写入器
func UnregisterWriters() {
	_centor.unregisterWriters()
}

// LoadJobPlugin 目前未正常实现该函数，仅仅是个架子
// todo: 实现handle ?
func LoadJobPlugin(typ plugin.Type, name string) (plugin.Job, error) {
	return newdefaultJobPlugin(), nil
}

// LoadReaderJob ,根据名字name获取读取器的工作
// 如果name不存在，返回的布尔值为false
func LoadReaderJob(name string) (reader.Job, bool) {
	r, ok := _centor.reader(name)
	if !ok {
		return nil, false
	}
	return r.Job(), true
}

// LoadReaderTask ,根据名字name获取读取器的任务
// 如果name不存在，返回的布尔值为false
func LoadReaderTask(name string) (reader.Task, bool) {
	r, ok := _centor.reader(name)
	if !ok {
		return nil, false
	}
	return r.Task(), true
}

// LoadWriterJob ,根据名字name获取写入器的工作
// 如果name不存在，返回的布尔值为false
func LoadWriterJob(name string) (writer.Job, bool) {
	w, ok := _centor.writer(name)
	if !ok {
		return nil, false
	}
	return w.Job(), true
}

// LoadWriterTask ,根据名字name获取写入器的任务
// 如果name不存在，返回的布尔值为false
func LoadWriterTask(name string) (writer.Task, bool) {
	w, ok := _centor.writer(name)
	if !ok {
		return nil, false
	}
	return w.Task(), true
}

type centor struct {
	readersMu sync.RWMutex
	readers   map[string]spi.Reader

	writersMu sync.RWMutex
	writers   map[string]spi.Writer
}

func (l *centor) registerReader(name string, reader spi.Reader) error {

	l.readersMu.Lock()
	defer l.readersMu.Unlock()

	if reader == nil {
		return fmt.Errorf("datax: reader %v is nil", name)
	}

	if reader.Task() == nil || reader.Job() == nil {
		return fmt.Errorf("datax: reader %v has nil job or task", name)
	}

	if _, ok := l.readers[name]; ok {
		return fmt.Errorf("datax: reader %v has already registered", name)
	}

	l.readers[name] = reader
	return nil
}

func (l *centor) reader(name string) (reader spi.Reader, ok bool) {
	l.readersMu.RLock()
	defer l.readersMu.RUnlock()
	reader, ok = l.readers[name]
	return
}

func (l *centor) registerWriter(name string, writer spi.Writer) error {
	l.writersMu.Lock()
	defer l.writersMu.Unlock()

	if writer == nil {
		return fmt.Errorf("datax: writer %v is nil", name)
	}

	if writer.Task() == nil || writer.Job() == nil {
		return fmt.Errorf("datax: writer %v has nil job or task", name)
	}

	if _, ok := l.writers[name]; ok {
		return fmt.Errorf("datax: writer %v has already registered", name)
	}
	l.writers[name] = writer
	return nil
}

func (l *centor) writer(name string) (writer spi.Writer, ok bool) {
	l.writersMu.RLock()
	defer l.writersMu.RUnlock()
	writer, ok = l.writers[name]
	return
}

func (l *centor) unregisterReaders() {
	l.readersMu.Lock()
	defer l.readersMu.Unlock()
	for k := range l.readers {
		l.readers[k] = nil
	}
	l.readers = make(map[string]spi.Reader)
}

func (l *centor) unregisterWriters() {
	l.writersMu.Lock()
	defer l.writersMu.Unlock()
	for k := range l.writers {
		l.writers[k] = nil
	}
	l.writers = make(map[string]spi.Writer)
}

type defaultJobPlugin struct {
	*plugin.BaseJob
}

func newdefaultJobPlugin() *defaultJobPlugin {
	return &defaultJobPlugin{
		BaseJob: plugin.NewBaseJob(),
	}
}

func (d *defaultJobPlugin) Init(ctx context.Context) error {
	return nil
}

func (d *defaultJobPlugin) Destroy(ctx context.Context) error {
	return nil
}
