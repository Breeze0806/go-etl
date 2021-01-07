package taskgroup

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/Breeze0806/go-etl/config"
	coreconst "github.com/Breeze0806/go-etl/datax/common/config/core"
	"github.com/Breeze0806/go-etl/datax/core"
	"github.com/Breeze0806/go-etl/schedule"
)

type Container struct {
	*core.BaseCotainer
	jobId         int64
	taskGroupId   int64
	scheduler     *schedule.TaskSchduler
	wg            sync.WaitGroup
	tasks         *taskManager
	ctx           context.Context
	sleepInterval time.Duration
	retryInterval time.Duration
	retryMaxCount int32
}

func NewContainer(ctx context.Context, conf *config.Json) (c *Container, err error) {
	c = &Container{
		BaseCotainer: core.NewBaseCotainer(),
		tasks:        newTaskManager(),
		ctx:          ctx,
	}
	c.SetConfig(conf)
	c.jobId, err = c.Config().GetInt64(coreconst.DataxCoreContainerJobId)
	if err != nil {
		return nil, err
	}
	c.taskGroupId, err = c.Config().GetInt64(coreconst.DataxCoreContainerTaskGroupId)
	if err != nil {
		return nil, err
	}

	c.sleepInterval = time.Duration(
		c.Config().GetInt64OrDefaullt(coreconst.DataxCoreContainerJobSleepinterval, 100)) * time.Millisecond
	c.retryInterval = time.Duration(
		c.Config().GetInt64OrDefaullt(coreconst.DataxCoreContainerTaskFailoverMaxretrytimes, 10000)) * time.Millisecond
	c.retryMaxCount = int32(c.Config().GetInt64OrDefaullt(coreconst.DataxCoreContainerTaskFailoverMaxretrytimes, 1))
	log.Infof("datax job(%v) taskgruop(%v) sleepInterval: %v retryInterval: %v retryMaxCount: %v",
		c.jobId, c.taskGroupId, c.sleepInterval, c.retryInterval, c.retryMaxCount)
	return
}

func (c *Container) JobId() int64 {
	return c.jobId
}

func (c *Container) TaskGroupId() int64 {
	return c.taskGroupId
}

func (c *Container) Do() error {
	return c.Start()
}

func (c *Container) Start() (err error) {
	log.Infof("datax job(%v) taskgruop(%v)  start", c.jobId, c.taskGroupId)
	defer log.Infof("datax job(%v) taskgruop(%v)  end", c.jobId, c.taskGroupId)
	var taskConfigs []*config.Json
	if taskConfigs, err = c.Config().GetConfigArray(coreconst.DataxJobContent); err != nil {
		return err
	}
	c.scheduler = schedule.NewTaskSchduler(
		int(c.Config().GetInt64OrDefaullt(coreconst.DataxCoreContainerTaskGroupMaxWorkerNumber, 4)), len(taskConfigs))
	defer c.scheduler.Stop()
	prefixKey := strconv.FormatInt(c.jobId, 10) + "-" + strconv.FormatInt(c.taskGroupId, 10)
	log.Infof("datax job(%v) taskgruop(%v) manager config", c.jobId, c.taskGroupId)
	for i := range taskConfigs {
		var taskExecer *taskExecer

		taskExecer, err = newTaskExecer(c.ctx, taskConfigs[i], prefixKey, 0)
		if err != nil {
			return err
		}
		c.tasks.pushRemain(taskExecer)
	}
	log.Infof("datax job(%v) taskgruop(%v) start tasks", c.jobId, c.taskGroupId)
	for i := 0; i < len(taskConfigs); i++ {
		te, ok := c.tasks.popRemainAndAddRun()
		if !ok {
			continue
		}
		if err = c.startTaskExecer(te); err != nil {
			return
		}
	}
	log.Infof("datax job(%v) taskgruop(%v) manage tasks", c.jobId, c.taskGroupId)
	ticker := time.NewTicker(c.sleepInterval)
	defer ticker.Stop()
QueueLoop:
	for !c.tasks.isEmpty() {
		select {
		case <-c.ctx.Done():
			break QueueLoop
		default:
		}
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
	log.Infof("datax job(%v) taskgruop(%v) wait tasks end", c.jobId, c.taskGroupId)
	c.wg.Wait()
	if c.ctx.Err() != nil {
		return c.ctx.Err()
	}

	return nil
}

func (c *Container) startTaskExecer(te *taskExecer) (err error) {
	log.Debugf("datax job(%v) taskgruop(%v) task(%v) push", c.jobId, c.taskGroupId, te.Key())
	c.wg.Add(1)
	var errChan <-chan error
	errChan, err = c.scheduler.Push(te)
	if err != nil {
		c.tasks.removeRun(te)
		c.wg.Done()
		return err
	}
	log.Debugf("datax job(%v) taskgruop(%v) task(%v) start", c.jobId, c.taskGroupId, te.Key())
	go func(te *taskExecer) {
		defer c.wg.Done()
		timer := time.NewTimer(c.retryInterval)
		defer timer.Stop()
		select {
		case err := <-errChan:
			if err != nil && te.WriterSuportFailOverport() && te.AttemptCount() <= c.retryMaxCount {
				log.Debugf("datax job(%v) taskgruop(%v) task(%v) shutdown and retry. attemptCount: %v err: %v",
					c.jobId, c.taskGroupId, te.Key(), te.AttemptCount(), err)
				te.Shutdown()
				select {
				case <-timer.C:
				case <-c.ctx.Done():
				}
				c.tasks.removeRunAndPushRemain(te)
			} else {
				log.Debugf("datax job(%v) taskgruop(%v) task(%v) end", c.jobId, c.taskGroupId, te.Key())
				c.tasks.removeRun(te)
			}
		case <-c.ctx.Done():
		}
	}(te)
	return
}
