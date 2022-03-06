package file

import (
	"fmt"
	"sync"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/element"
)

//FetchHandler 获取记录句柄
type FetchHandler interface {
	OnRecord(element.Record) error
	CreateRecord() (element.Record, error)
}

type Opener interface {
	Open(filename string) (stream Stream, err error)
}

type Stream interface {
	Rows(conf *config.JSON) (rows Rows, err error)
	Close() (err error)
}

type Rows interface {
	Next() bool
	Scan() (columns []element.Column, err error)
	Error() error
	Close() error
}

func RegisterOpener(name string, opener Opener) {
	if err := openers.register(name, opener); err != nil {
		panic(err)
	}
}

type Streamer struct {
	stream Stream
}

func NewStreamer(name string, filename string) (streamer *Streamer, err error) {
	streamer = &Streamer{}
	opener, ok := openers.opener(name)
	if !ok {
		err = fmt.Errorf("opener %v does not exist", name)
		return
	}
	if streamer.stream, err = opener.Open(filename); err != nil {
		return nil, fmt.Errorf("open fail. err : %v", err)
	}
	return
}

func (s *Streamer) Read(conf *config.JSON, handler FetchHandler) (err error) {
	var rows Rows
	rows, err = s.stream.Rows(conf)
	if err != nil {
		return fmt.Errorf("rows fail. error: %v, config： %v", err, conf.String())
	}
	defer rows.Close()
	for rows.Next() {
		var columns []element.Column

		if columns, err = rows.Scan(); err != nil {
			return fmt.Errorf("Scan fail. error: %v", err)
		}
		var r element.Record

		if r, err = handler.CreateRecord(); err != nil {
			return fmt.Errorf("CreateRecord fail. error: %v", err)
		}

		for _, v := range columns {
			if err = r.Add(v); err != nil {
				return fmt.Errorf("Add fail. error: %v", err)
			}
		}
		if err = handler.OnRecord(r); err != nil {
			return fmt.Errorf("OnRecord fail. error: %v", err)
		}
	}
	if err = rows.Error(); err != nil {
		return fmt.Errorf("Error. error: %v", err)
	}
	return
}

func (s *Streamer) Close() error {
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
