package exchange

import (
	"errors"

	"github.com/Breeze0806/go-etl/datax/common/element"
	"github.com/Breeze0806/go-etl/datax/core/transport/channel"
	"github.com/Breeze0806/go-etl/datax/transform"
)

var (
	ErrTerminate = errors.New("reader is terminated")
	ErrClose     = errors.New("chan close")
	ErrShutdown  = errors.New("exchange is shutdowned")
)

type RecordExchanger struct {
	tran       transform.Transformer
	ch         *channel.Channel
	isShutdown bool
}

func NewRecordExchangerWithoutTransformer(ch *channel.Channel) *RecordExchanger {
	return NewRecordExchanger(ch, &transform.NilTransformer{})
}

func NewRecordExchanger(ch *channel.Channel, tran transform.Transformer) *RecordExchanger {
	return &RecordExchanger{
		tran: tran,
		ch:   ch,
	}
}

func (r *RecordExchanger) GetFromReader() (element.Record, error) {
	if r.isShutdown {
		return nil, ErrShutdown
	}
	record, ok := r.ch.Pop()
	if !ok {
		return nil, ErrClose
	}

	switch record.(type) {
	case *element.TerminateRecord:
		return nil, ErrTerminate
	default:
		return record, nil
	}
}

func (r *RecordExchanger) Shutdown() error {
	r.isShutdown = true
	return nil
}

func (r *RecordExchanger) CreateRecord() (element.Record, error) {
	return element.NewDefaultRecord(), nil
}

func (r *RecordExchanger) SendWriter(record element.Record) (err error) {
	if r.isShutdown {
		return ErrShutdown
	}
	var newRecord element.Record
	if newRecord, err = r.tran.DoTransform(record); err == nil {
		r.ch.Push(newRecord)
	}
	return
}

func (r *RecordExchanger) Flush() error {
	return nil
}

func (r *RecordExchanger) Terminate() error {
	r.ch.PushTerminate()
	return nil
}
