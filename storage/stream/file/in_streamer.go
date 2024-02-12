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

// FetchHandler - Acquires the record handler
type FetchHandler interface {
	OnRecord(element.Record) error         // Process Record - Handles the record
	CreateRecord() (element.Record, error) // Create Empty Record - Creates an empty record
}

// Opener - An opener used to open an input stream
type Opener interface {
	Open(filename string) (stream InStream, err error) // Open Input Stream - Opens an input stream for the file named 'filename'
}

// InStream - Input stream
type InStream interface {
	Rows(conf *config.JSON) (rows Rows, err error) // Get Line Reader - Acquires a line reader
	Close() (err error)                            // Close Input Stream - Closes the input stream
}

// Rows - Line reader
type Rows interface {
	Next() bool                                  // Get Next Line - Returns true if there is a next line, false otherwise
	Scan() (columns []element.Column, err error) // Scan Columns - Scans the columns of each line
	Error() error                                // Get Error of Next Line - Gets the error of the next line
	Close() error                                // Close Line Reader - Closes the line reader
}

// RegisterOpener - Registers an input stream opener with the given name 'name'
func RegisterOpener(name string, opener Opener) {
	if err := openers.register(name, opener); err != nil {
		panic(err)
	}
}

// UnregisterAllOpener - Unregisters all file openers
func UnregisterAllOpener() {
	openers.unregisterAll()
}

// InStreamer - Input stream wrapper
type InStreamer struct {
	stream InStream
}

// NewInStreamer - Opens an input stream named 'filename' using the input stream opener with the given 'name'
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

// Read - Reads data using the record handler 'handler', context 'ctx', and configuration file 'conf'
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

// Close - Closes the input stream
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
