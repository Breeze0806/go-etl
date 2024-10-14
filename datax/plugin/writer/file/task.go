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
	"sync"
	"time"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/datax/common/plugin"
	"github.com/Breeze0806/go-etl/datax/common/spi/writer"
	"github.com/Breeze0806/go-etl/datax/core/transport/exchange"
	"github.com/Breeze0806/go-etl/element"
	"github.com/Breeze0806/go-etl/storage/stream/file"
)

// Task normal file task
type Task struct {
	*writer.BaseTask

	streamer  *file.OutStreamer
	conf      Config
	newConfig func(conf *config.JSON) (Config, error)
	content   *config.JSON
}

// NewTask Create a task by obtaining the configuration newConfig
func NewTask(newConfig func(conf *config.JSON) (Config, error)) *Task {
	return &Task{
		BaseTask:  writer.NewBaseTask(),
		newConfig: newConfig,
	}
}

// Init Initialization
func (t *Task) Init(ctx context.Context) (err error) {
	var name string
	if name, err = t.PluginConf().GetString("creator"); err != nil {
		return t.Wrapf(err, "GetString fail")
	}
	var filename string
	if filename, err = t.PluginJobConf().GetString("path"); err != nil {
		return t.Wrapf(err, "GetString fail")
	}

	if t.content, err = t.PluginJobConf().GetConfig("content"); err != nil {
		return t.Wrapf(err, "GetString fail")
	}

	if t.conf, err = t.newConfig(t.content); err != nil {
		return t.Wrapf(err, "newConfig fail")
	}

	if t.streamer, err = file.NewOutStreamer(name, filename); err != nil {
		return t.Wrapf(err, "NewOutStreamer fail")
	}

	return
}

// Destroy Destruction
func (t *Task) Destroy(ctx context.Context) (err error) {
	if t.streamer != nil {
		err = t.streamer.Close()
	}
	return t.Wrapf(err, "Close fail")
}

// StartWrite Start writing
func (t *Task) StartWrite(ctx context.Context, receiver plugin.RecordReceiver) (err error) {
	var sw file.StreamWriter
	if sw, err = t.streamer.Writer(t.content); err != nil {
		return t.Wrapf(err, "Writer fail")
	}

	recordChan := make(chan element.Record)
	var rerr error
	afterCtx, cancel := context.WithCancel(ctx)
	var wg sync.WaitGroup
	wg.Add(1)
	// Read records from the receiver and place them into the recordChan channel
	go func() {
		defer func() {
			wg.Done()
			// Close the recordChan channel
			close(recordChan)
			log.Debugf(t.Format("get records end"))
		}()
		log.Debugf(t.Format("start to get records"))
		for {
			select {
			case <-afterCtx.Done():
				return
			default:
			}
			var record element.Record
			record, rerr = receiver.GetFromReader()
			if rerr != nil && rerr != exchange.ErrEmpty {
				return
			}

			// When the receiver returns a non-empty error, write it to the recordChan
			if rerr != exchange.ErrEmpty {
				select {
				// Prevent not writing to the recordChan when the context ctx is closed
				case <-afterCtx.Done():
					return
				case recordChan <- record:
				}

			}
		}
	}()
	ticker := time.NewTicker(t.conf.GetBatchTimeout())
	defer ticker.Stop()
	cnt := 0
	log.Debugf(t.Format("start to write"))
	for {
		select {
		case record, ok := <-recordChan:
			if !ok {
				// When writing ends, write the remaining records to the database
				if cnt > 0 {
					if err = sw.Flush(); err != nil {
						log.Errorf(t.Format("Flush error: %v"), err)
					}
				}
				if err == nil {
					err = rerr
				}
				goto End
			}

			// Write to file
			if err = sw.Write(record); err != nil {
				log.Errorf(t.Format("Write error: %v"), err)
				goto End
			}
			cnt++
			// When the data volume exceeds the single batch size, write to the file
			if cnt >= t.conf.GetBatchSize() {
				if err = sw.Flush(); err != nil {
					log.Errorf(t.Format("Flush error: %v"), err)
					goto End
				}
				cnt = 0
			}
		// When the written data does not reach the single batch size, write even if the timeout occurs
		case <-ticker.C:
			if cnt > 0 {
				if err = sw.Flush(); err != nil {
					log.Errorf(t.Format("Flush error: %v"), err)
					goto End
				}
			}
			cnt = 0
		}
	}
End:
	if cerr := sw.Close(); cerr != nil {
		log.Errorf(t.Format("Close error: %v"), cerr)
	}
	cancel()
	log.Debugf(t.Format("wait all goroutine"))
	// Wait for the goroutine to finish
	wg.Wait()
	log.Debugf(t.Format(" wait all goroutine end"))
	switch {
	// Starting to write is not an error when externally canceled
	case ctx.Err() != nil:
		return nil
	// When the error is a stop, it is also not an error
	case err == exchange.ErrTerminate:
		return nil
	}
	return t.Wrapf(err, "")
}
