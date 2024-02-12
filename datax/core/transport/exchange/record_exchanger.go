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

// Error Enumeration - An enumeration of possible errors
var (
	ErrTerminate = errors.New("reader is terminated")
	ErrEmpty     = errors.New("chan is empty")
	ErrShutdown  = errors.New("exchange is shutdowned")
)

// RecordExchanger - A component responsible for exchanging records
type RecordExchanger struct {
	tran       transform.Transformer
	ch         *channel.Channel
	isShutdown bool
}

// NewRecordExchangerWithoutTransformer - Creates a new instance of a RecordExchanger without a transformer
func NewRecordExchangerWithoutTransformer(ch *channel.Channel) *RecordExchanger {
	return NewRecordExchanger(ch, &transform.NilTransformer{})
}

// NewRecordExchanger - Creates a new instance of a RecordExchanger based on a channel (ch) and a transformer (tran)
func NewRecordExchanger(ch *channel.Channel, tran transform.Transformer) *RecordExchanger {
	return &RecordExchanger{
		tran: tran,
		ch:   ch,
	}
}

// GetFromReader - Retrieves records from the Reader
// An error will be reported if the exchanger is closed, the channel is empty, or a termination message is received
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

// Shutdown - Closes the exchanger
func (r *RecordExchanger) Shutdown() error {
	r.isShutdown = true
	return nil
}

// CreateRecord - Creates a new record
func (r *RecordExchanger) CreateRecord() (element.Record, error) {
	return element.NewDefaultRecord(), nil
}

// SendWriter - Writes a record (recode) to the writer, potentially transforming it through a transformer
// An error will be reported if the transformation fails or if the channel is closed
func (r *RecordExchanger) SendWriter(record element.Record) (err error) {
	if r.isShutdown {
		return ErrShutdown
	}

	r.ch.Push(record)

	return
}

// Flush - Flushes any pending data, but is an empty method in this context
func (r *RecordExchanger) Flush() error {
	return nil
}

// Terminate - Terminates the record exchange process
func (r *RecordExchanger) Terminate() error {
	r.ch.PushTerminate()
	return nil
}
