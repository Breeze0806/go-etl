package file

import (
	"fmt"
	"sync"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/element"
)

type Creater interface {
	Create(filename string) (stream OutStream, err error)
}

type OutStream interface {
	Writer(conf *config.JSON) (writer StreamWriter, err error)
	Close() (err error)
}

type StreamWriter interface {
	Write(record element.Record) (err error)
	Flush() (err error)
	Close() (err error)
}

func RegisterCreater(name string, creater Creater) {
	if err := creaters.register(name, creater); err != nil {
		panic(err)
	}
}

type OutStreamer struct {
	stream OutStream
}

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

func (s *OutStreamer) Writer(conf *config.JSON) (StreamWriter, error) {
	return s.stream.Writer(conf)
}

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

func (o *createrMap) unregisterAll() {
	o.Lock()
	defer o.Unlock()
	o.creaters = make(map[string]Creater)
}
