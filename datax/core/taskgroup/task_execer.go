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

package taskgroup

import (
	"context"
	"fmt"
	"sync"

	"github.com/Breeze0806/go-etl/config"
	coreconst "github.com/Breeze0806/go-etl/datax/common/config/core"
	"github.com/Breeze0806/go-etl/datax/common/plugin/loader"
	"github.com/Breeze0806/go-etl/datax/common/spi/writer"
	"github.com/Breeze0806/go-etl/datax/core/taskgroup/runner"
	"github.com/Breeze0806/go-etl/datax/core/transport/channel"
	"github.com/Breeze0806/go-etl/datax/core/transport/exchange"
	"github.com/pingcap/errors"
	"go.uber.org/atomic"
)

type taskExecer struct {
	Err error

	taskConf     *config.JSON // Task JSON Configuration - Configuration for a task in JSON format
	taskID       int64        // Task ID - Unique identifier for the task
	ctx          context.Context
	channel      *channel.Channel // Log Channel - Channel used for logging task-related information
	writerRunner runner.Runner    // Write Runner - Component responsible for writing operations in the task
	readerRunner runner.Runner    // Execute Runner - Component responsible for executing the main operations of the task
	wg           sync.WaitGroup
	errors       chan error

	destroy      sync.Once
	key          string
	exchanger    *exchange.RecordExchanger
	cancalMutex  sync.Mutex         // Since the cancel function can be called by multiple threads, locking is required
	cancel       context.CancelFunc // Cancel Function - Function used to cancel the task's execution
	attemptCount *atomic.Int32      // Execution Count - Number of times the task has been executed
}

// newTaskExecer - Creates a new task executor based on the context (ctx), task configuration (taskConf), prefix key (prefixKey),
// and attempt count (attemptCount). It will report an error if the task ID does not exist, or if the configured writer or reader does not exist.

func newTaskExecer(ctx context.Context, taskConf *config.JSON,
	jobID, taskGroupID int64, attemptCount int) (t *taskExecer, err error) {
	t = &taskExecer{
		taskConf:     taskConf,
		errors:       make(chan error, 2),
		ctx:          ctx,
		attemptCount: atomic.NewInt32(int32(attemptCount)),
	}
	t.channel = channel.NewChannel(ctx, taskConf)

	t.taskID, err = taskConf.GetInt64(coreconst.TaskID)
	if err != nil {
		return nil, err
	}
	t.key = fmt.Sprintf("%v-%v-%v", jobID, taskGroupID, t.taskID)
	readName, writeName := "", ""
	readName, err = taskConf.GetString(coreconst.JobReaderName)
	if err != nil {
		return nil, err
	}

	writeName, err = taskConf.GetString(coreconst.JobWriterName)
	if err != nil {
		return nil, err
	}

	var readConf, writeConf *config.JSON
	readConf, err = taskConf.GetConfig(coreconst.JobReaderParameter)
	if err != nil {
		return nil, err
	}

	writeConf, err = taskConf.GetConfig(coreconst.JobWriterParameter)
	if err != nil {
		return nil, err
	}

	readTask, ok := loader.LoadReaderTask(readName)
	if !ok {
		return nil, errors.Errorf("reader task name (%v) does not exist", readName)
	}
	readTask.SetJobID(jobID)
	readTask.SetTaskGroupID(taskGroupID)
	readTask.SetTaskID(t.taskID)
	readTask.SetPluginJobConf(readConf)
	readTask.SetPeerPluginName(writeName)
	readTask.SetPeerPluginJobConf(writeConf)
	t.exchanger = exchange.NewRecordExchangerWithoutTransformer(t.channel)
	t.readerRunner = runner.NewReader(readTask, t.exchanger, t.key)

	writeTask, ok := loader.LoadWriterTask(writeName)
	if !ok {
		return nil, errors.Errorf("writer task name (%v) does not exist", writeName)
	}
	writeTask.SetJobID(jobID)
	writeTask.SetTaskGroupID(taskGroupID)
	writeTask.SetTaskID(t.taskID)
	writeTask.SetPluginJobConf(writeConf)
	writeTask.SetPeerPluginName(readName)
	writeTask.SetPeerPluginJobConf(readConf)
	t.writerRunner = runner.NewWriter(writeTask, t.exchanger, t.key)

	return
}

// Start - Executes the read runner and write runner in separate goroutines
func (t *taskExecer) Start() {
	var ctx context.Context
	t.cancalMutex.Lock()
	ctx, t.cancel = context.WithCancel(t.ctx)
	t.cancalMutex.Unlock()

	log.Debugf("taskExecer %v start to run writer", t.key)
	t.wg.Add(1)
	var writerWg sync.WaitGroup
	writerWg.Add(1)
	go func() {
		defer t.wg.Done()
		writerWg.Done()
		if err := t.writerRunner.Run(ctx); err != nil {
			log.Errorf("writer task(%v) fail, err: %v", t.Key(), err)
			t.errors <- err
		} else {
			t.errors <- nil
		}
	}()
	writerWg.Wait()

	log.Debugf("taskExecer %v start to run reader", t.key)
	var readerWg sync.WaitGroup
	t.wg.Add(1)
	readerWg.Add(1)
	go func() {
		defer t.wg.Done()
		readerWg.Done()
		if err := t.readerRunner.Run(ctx); err != nil {
			log.Errorf("reader task(%v) fail, err: %v", t.Key(), err)
			t.errors <- err
		} else {
			t.errors <- nil
		}
	}()
	readerWg.Wait()
}

// AttemptCount - Number of attempts made to execute the task
func (t *taskExecer) AttemptCount() int32 {
	return t.attemptCount.Load()
}

// Do - The execution function for the task
func (t *taskExecer) Do() (err error) {
	log.Debugf("taskExecer %v start to do", t.key)
	defer func() {
		t.attemptCount.Inc()
		log.Debugf("taskExecer %v end to do", t.key)
	}()
	// Execute Read-Write Runner - Executes the runner responsible for both reading and writing operations
	t.Start()
	log.Debugf("taskExecer %v do wait runner stop", t.key)
	cnt := 0
	for {
		select {
		case err := <-t.errors:
			if err != nil {
				return err
			}
			cnt++
			if cnt == 2 {
				return nil
			}
		case <-t.ctx.Done():
		}
	}
}

// Key - A keyword or identifier used in the context of the task
func (t *taskExecer) Key() string {
	return t.key
}

// WriterSupportsFailover - Indicates whether the writer supports failover or error retry mechanisms
func (t *taskExecer) WriterSuportFailOverport() bool {
	task, ok := t.writerRunner.Plugin().(writer.Task)
	if !ok {
		return false
	}
	return task.SupportFailOver()
}

// Shutdown - Stops the writer through cancellation and closes both the reader and writer
func (t *taskExecer) Shutdown() {
	log.Debugf("taskExecer %v starts to shutdown", t.key)
	defer log.Debugf("taskExecer %v ends to shutdown", t.key)
	t.wg.Add(1)
	go func() {
		t.wg.Done()
		for {
			var rerr error
			_, rerr = t.exchanger.GetFromReader()
			if rerr != nil && rerr != exchange.ErrEmpty {
				return
			}
		}
	}()

	t.cancalMutex.Lock()
	if t.cancel != nil {
		t.cancel()
	}
	t.cancalMutex.Unlock()

	log.Debugf("taskExecer %v shutdown wait runner stop", t.key)
	t.wg.Wait()

	log.Debugf("taskExecer %v shutdown reader", t.key)

Loop:
	for {
		select {
		case <-t.errors:
		default:
			break Loop
		}
	}

	t.readerRunner.Shutdown()

	log.Debugf("taskExecer %v shutdown writer", t.key)
	t.writerRunner.Shutdown()
}

// Stats - Information or metrics related to the task's execution
type Stats struct {
	TaskID  int64             `json:"taskID"`
	Channel channel.StatsJSON `json:"channel"`
}

// Get Stats - Retrieves the statistical information related to the task's execution
func (t *taskExecer) Stats() Stats {
	return Stats{
		TaskID:  t.taskID,
		Channel: t.channel.StatsJSON(),
	}
}
