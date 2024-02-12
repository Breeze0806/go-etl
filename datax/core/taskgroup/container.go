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
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Breeze0806/go-etl/config"
	coreconst "github.com/Breeze0806/go-etl/datax/common/config/core"
	"github.com/Breeze0806/go-etl/datax/core"
	"github.com/Breeze0806/go-etl/datax/core/statistics/container"
	"github.com/Breeze0806/go-etl/schedule"
	"github.com/pingcap/errors"
)

// Container represents the environment for a group of tasks.
type Container struct {
	*core.BaseCotainer

	Err error

	jobID          int64
	taskGroupID    int64
	scheduler      *schedule.TaskSchduler
	wg             sync.WaitGroup
	tasks          *taskManager
	ctx            context.Context
	reportInterval time.Duration
	retryInterval  time.Duration
	sleepInterval  time.Duration
	retryMaxCount  int32
}

// NewContainer creates a task group container based on the JSON configuration conf.
// If jobID or taskGroupID is invalid, an error will be reported.
func NewContainer(ctx context.Context, conf *config.JSON) (c *Container, err error) {
	c = &Container{
		BaseCotainer: core.NewBaseCotainer(),
		tasks:        newTaskManager(),
		ctx:          ctx,
	}
	c.SetConfig(conf)
	c.SetMetrics(container.NewMetrics())
	c.jobID, err = c.Config().GetInt64(coreconst.DataxCoreContainerJobID)
	if err != nil {
		return nil, err
	}
	c.taskGroupID, err = c.Config().GetInt64(coreconst.DataxCoreContainerTaskGroupID)
	if err != nil {
		return nil, err
	}
	c.Metrics().Set("taskGroupID", c.taskGroupID)
	c.reportInterval = time.Duration(
		c.Config().GetInt64OrDefaullt(coreconst.DataxCoreContainerTaskGroupReportinterval, 1)) * time.Second
	c.sleepInterval = time.Duration(
		c.Config().GetInt64OrDefaullt(coreconst.DataxCoreContainerTaskGroupSleepinterval, 1)) * time.Second
	c.retryInterval = time.Duration(
		c.Config().GetInt64OrDefaullt(coreconst.DataxCoreContainerTaskFailoverRetryintervalinmsec, 1000)) * time.Millisecond
	c.retryMaxCount = int32(c.Config().GetInt64OrDefaullt(coreconst.DataxCoreContainerTaskFailoverMaxretrytimes, 1))
	log.Infof("datax job(%v) taskgruop(%v) reportInterval: %v retryInterval: %v retryMaxCount: %v config: %v",
		c.jobID, c.taskGroupID, c.reportInterval, c.retryInterval, c.retryMaxCount, conf)
	return
}

// JobID refers to the unique identifier for a job.
func (c *Container) JobID() int64 {
	return c.jobID
}

// TaskGroupID refers to the unique identifier for a group of tasks.
func (c *Container) TaskGroupID() int64 {
	return c.taskGroupID
}

// Do represents the execution action.
func (c *Container) Do() error {
	return c.Start()
}

// Start initiates the execution of tasks using the task scheduler based on the JSON configurations.
func (c *Container) Start() (err error) {
	log.Infof("datax job(%v) taskgruop(%v)  start", c.jobID, c.taskGroupID)
	defer log.Infof("datax job(%v) taskgruop(%v)  end", c.jobID, c.taskGroupID)
	var taskConfigs []*config.JSON
	if taskConfigs, err = c.Config().GetConfigArray(coreconst.DataxJobContent); err != nil {
		return err
	}
	c.scheduler = schedule.NewTaskSchduler(
		int(c.Config().GetInt64OrDefaullt(coreconst.DataxCoreContainerTaskGroupMaxWorkerNumber, 4)), len(taskConfigs))
	defer c.scheduler.Stop()
	log.Infof("datax job(%v) taskgruop(%v) manager config", c.jobID, c.taskGroupID)
	for i := range taskConfigs {
		var taskExecer *taskExecer

		taskExecer, err = newTaskExecer(c.ctx, taskConfigs[i], c.jobID, c.taskGroupID, 0)
		if err != nil {
			return err
		}
		// Add the task executor to the pending execution queue.
		c.tasks.pushRemain(taskExecer)
	}
	log.Infof("datax job(%v) taskgruop(%v) start tasks", c.jobID, c.taskGroupID)
	var tasks []*taskExecer
	for i := 0; i < len(taskConfigs); i++ {
		// Move tasks from the pending execution queue to the running queue.
		te, ok := c.tasks.popRemainAndAddRun()
		if !ok {
			continue
		}
		tasks = append(tasks, te)
		// Begin the execution of tasks.
		if err = c.startTaskExecer(te); err != nil {
			return
		}
	}
	log.Infof("datax job(%v) taskgruop(%v) manage tasks", c.jobID, c.taskGroupID)
	ticker := time.NewTicker(c.sleepInterval)
	defer ticker.Stop()
QueueLoop:
	// The task queue is not empty.
	for !c.tasks.isRunsEmpty() {
		for !c.tasks.isEmpty() {
			select {
			case <-c.ctx.Done():
				break QueueLoop
			default:
			}
			// Move tasks from the pending execution queue to the running queue.
			te, ok := c.tasks.popRemainAndAddRun()
			if !ok {
				select {
				case <-ticker.C:
				case <-c.ctx.Done():
					break QueueLoop
				}
				continue
			}
			if err = c.startTaskExecer(te); err != nil {
				return
			}
		}
		select {
		case <-c.ctx.Done():
			break QueueLoop
		case <-ticker.C:
		}
	}
	log.Infof("datax job(%v) taskgruop(%v) wait tasks end", c.jobID, c.taskGroupID)
	// Wait for all tasks to complete.
	c.wg.Wait()
	if c.ctx.Err() != nil {
		return c.ctx.Err()
	}

	b := &strings.Builder{}
	for _, t := range tasks {
		if t.Err != nil {
			b.WriteString(t.Err.Error())
			b.WriteByte(' ')
		}
	}

	if b.Len() != 0 {
		return errors.NewNoStackError(b.String())
	}

	return nil
}

// startTaskExecer initiates a task.
func (c *Container) startTaskExecer(te *taskExecer) (err error) {
	log.Debugf("datax job(%v) taskgruop(%v) task(%v) push", c.jobID, c.taskGroupID, te.Key())
	c.wg.Add(1)
	var errChan <-chan error
	// Add the task to the scheduler.
	errChan, err = c.scheduler.Push(te)
	if err != nil {
		// Remove a task from the running queue when an error occurs.
		c.tasks.removeRun(te)
		c.wg.Done()
		return err
	}
	go func(te *taskExecer) {
		log.Debugf("datax job(%v) taskgruop(%v) task(%v) start", c.jobID, c.taskGroupID, te.Key())
		defer func() {
			log.Debugf("datax job(%v) taskgruop(%v) task(%v) end", c.jobID, c.taskGroupID, te.Key())
			c.wg.Done()
		}()
		statsTimer := time.NewTicker(c.reportInterval)
		defer statsTimer.Stop()
		for {
			select {
			case te.Err = <-errChan:
				// If a task fails and the number of retries does not exceed the maximum retry limit, and the task supports failure retries, then decide whether to retry the task.
				if te.Err != nil && te.WriterSuportFailOverport() && te.AttemptCount() <= c.retryMaxCount {
					log.Debugf("datax job(%v) taskgruop(%v) task(%v) shutdown and retry. attemptCount: %v err: %v",
						c.jobID, c.taskGroupID, te.Key(), te.AttemptCount(), te.Err)
					// Close a task.
					te.Shutdown()
					timer := time.NewTimer(c.retryInterval)
					defer timer.Stop()
					select {
					case <-timer.C:
					case <-c.ctx.Done():
						return
					}
					// Move a task from the running queue back to the pending execution queue.
					c.tasks.removeRunAndPushRemain(te)
				} else {
					// Remove a task from the task scheduler.
					c.tasks.removeRun(te)
					c.setStats(te)
				}
				return
			case <-c.ctx.Done():
				return
			case <-statsTimer.C:
				c.setStats(te)
			}
		}
	}(te)
	return
}

func (c *Container) setStats(te *taskExecer) {
	key := "metrics." + strconv.FormatInt(te.taskID, 10)
	stats := te.Stats()

	c.Metrics().Set(key, stats)
}
