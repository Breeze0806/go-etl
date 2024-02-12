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
	"github.com/pingcap/errors"
)

// Creator - The creator that generates the output stream.
type Creator interface {
	Create(filename string) (stream OutStream, err error) // Create an output stream named 'filename'.
}

// OutStream - Represents the output stream.
type OutStream interface {
	Writer(conf *config.JSON) (writer StreamWriter, err error) // Create a writer for writing to the output stream.
	Close() (err error)                                        // Close the output stream.
}

// StreamWriter - A writer for writing to the output stream.
type StreamWriter interface {
	Write(record element.Record) (err error) // Write a record to the output stream.
	Flush() (err error)                      // Flush the data to the file.
	Close() (err error)                      // Close the output stream writer.
}

// RegisterCreator - Registers an output stream creator with the given name 'name'.
func RegisterCreator(name string, creator Creator) {
	if err := creators.register(name, creator); err != nil {
		panic(err)
	}
}

// UnregisterAllCreater - Unregister all file openers.
func UnregisterAllCreater() {
	creators.unregisterAll()
}

// OutStreamer - A wrapper for the output stream.
type OutStreamer struct {
	stream OutStream
}

// NewOutStreamer - Opens an output stream named 'filename' using the creator with the given name 'name'.
func NewOutStreamer(name string, filename string) (streamer *OutStreamer, err error) {
	creator, ok := creators.creator(name)
	if !ok {
		err = errors.Errorf("creator %v does not exist", name)
		return nil, err
	}
	streamer = &OutStreamer{}
	if streamer.stream, err = creator.Create(filename); err != nil {
		return nil, errors.Wrapf(err, "create fail")
	}
	return
}

// Writer - Creates a stream writer based on the configuration 'conf'.
func (s *OutStreamer) Writer(conf *config.JSON) (StreamWriter, error) {
	return s.stream.Writer(conf)
}

// Close - Closes the writing wrapper.
func (s *OutStreamer) Close() error {
	return s.stream.Close()
}

var creators = &creatorMap{
	creators: make(map[string]Creator),
}

type creatorMap struct {
	sync.RWMutex
	creators map[string]Creator
}

func (o *creatorMap) register(name string, creator Creator) error {
	if creator == nil {
		return fmt.Errorf("creator %v is nil", name)
	}

	o.Lock()
	defer o.Unlock()
	if _, ok := o.creators[name]; ok {
		return fmt.Errorf("creator %v exists", name)
	}

	o.creators[name] = creator
	return nil
}

func (o *creatorMap) creator(name string) (creator Creator, ok bool) {
	o.RLock()
	defer o.RUnlock()
	creator, ok = o.creators[name]
	return
}

func (o *creatorMap) unregisterAll() {
	o.Lock()
	defer o.Unlock()
	o.creators = make(map[string]Creator)
}
