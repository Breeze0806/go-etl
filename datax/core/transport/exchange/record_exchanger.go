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

package exchange

import (
	"errors"

	"github.com/Breeze0806/go-etl/datax/core/transport/channel"
	"github.com/Breeze0806/go-etl/datax/transform"
	"github.com/Breeze0806/go-etl/element"
)

// 错误枚举
var (
	ErrTerminate = errors.New("reader is terminated")
	ErrEmpty     = errors.New("chan is empty")
	ErrShutdown  = errors.New("exchange is shutdowned")
)

// RecordExchanger 记录交换器
type RecordExchanger struct {
	tran       transform.Transformer
	ch         *channel.Channel
	isShutdown bool
}

// NewRecordExchangerWithoutTransformer 生成不带转化器的记录交换器
func NewRecordExchangerWithoutTransformer(ch *channel.Channel) *RecordExchanger {
	return NewRecordExchanger(ch, &transform.NilTransformer{})
}

// NewRecordExchanger 根据通道ch和转化器tran生成的记录交换器
func NewRecordExchanger(ch *channel.Channel, tran transform.Transformer) *RecordExchanger {
	return &RecordExchanger{
		tran: tran,
		ch:   ch,
	}
}

// GetFromReader 从Reader中获取记录
// 当交换器关闭，通道为空或者收到终止消息也会报错
func (r *RecordExchanger) GetFromReader() (newRecord element.Record, err error) {
	if r.isShutdown {
		return nil, ErrShutdown
	}
	record, ok := r.ch.Pop()
	if !ok {
		return nil, ErrEmpty
	}

	switch record.(type) {
	case *element.TerminateRecord:
		return nil, ErrTerminate
	default:
		if newRecord, err = r.tran.DoTransform(record); err != nil {
			return nil, err
		}
		return
	}
}

// Shutdown 关闭
func (r *RecordExchanger) Shutdown() error {
	r.isShutdown = true
	return nil
}

// CreateRecord 创建记录
func (r *RecordExchanger) CreateRecord() (element.Record, error) {
	return element.NewDefaultRecord(), nil
}

// SendWriter 向写入器写入记录recode,其中还会通过转化器的转化
// 当转化失败或者通道已关闭时就会报错
func (r *RecordExchanger) SendWriter(record element.Record) (err error) {
	if r.isShutdown {
		return ErrShutdown
	}

	r.ch.Push(record)

	return
}

// Flush 刷新，空方法
func (r *RecordExchanger) Flush() error {
	return nil
}

// Terminate 终止记录交换
func (r *RecordExchanger) Terminate() error {
	r.ch.PushTerminate()
	return nil
}
