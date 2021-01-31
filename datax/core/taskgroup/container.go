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

//Container 任务组容器环境
type Container struct {
	*core.BaseCotainer

	jobID         int64
	taskGroupID   int64
	scheduler     *schedule.TaskSchduler
	wg            sync.WaitGroup
	tasks         *taskManager
	ctx           context.Context
	sleepInterval time.Duration
	retryInterval time.Duration
	retryMaxCount int32
}

//NewContainer 根据JSON配置conf创建任务组容器
//当jobID 和 taskGroupID非法就会报错
func NewContainer(ctx context.Context, conf *config.JSON) (c *Container, err error) {
	c = &Container{
		BaseCotainer: core.NewBaseCotainer(),
		tasks:        newTaskManager(),
		ctx:          ctx,
	}
	c.SetConfig(conf)
	c.jobID, err = c.Config().GetInt64(coreconst.DataxCoreContainerJobID)
	if err != nil {
		return nil, err
	}
	c.taskGroupID, err = c.Config().GetInt64(coreconst.DataxCoreContainerTaskGroupID)
	if err != nil {
		return nil, err
	}

	c.sleepInterval = time.Duration(
		c.Config().GetInt64OrDefaullt(coreconst.DataxCoreContainerJobSleepinterval, 100)) * time.Millisecond
	c.retryInterval = time.Duration(
		c.Config().GetInt64OrDefaullt(coreconst.DataxCoreContainerTaskFailoverMaxretrytimes, 10000)) * time.Millisecond
	c.retryMaxCount = int32(c.Config().GetInt64OrDefaullt(coreconst.DataxCoreContainerTaskFailoverMaxretrytimes, 1))
	log.Infof("datax job(%v) taskgruop(%v) sleepInterval: %v retryInterval: %v retryMaxCount: %v",
		c.jobID, c.taskGroupID, c.sleepInterval, c.retryInterval, c.retryMaxCount)
	return
}

//JobID 工作编号
func (c *Container) JobID() int64 {
	return c.jobID
}

//TaskGroupID 任务组编号
func (c *Container) TaskGroupID() int64 {
	return c.taskGroupID
}

//Do 执行
func (c *Container) Do() error {
	return c.Start()
}

//Start 开始运行，使用任务调度器执行这些JSON配置
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
	prefixKey := strconv.FormatInt(c.jobID, 10) + "-" + strconv.FormatInt(c.taskGroupID, 10)
	log.Infof("datax job(%v) taskgruop(%v) manager config", c.jobID, c.taskGroupID)
	for i := range taskConfigs {
		var taskExecer *taskExecer

		taskExecer, err = newTaskExecer(c.ctx, taskConfigs[i], prefixKey, 0)
		if err != nil {
			return err
		}
		//将任务执行器加入到待执行队列
		c.tasks.pushRemain(taskExecer)
	}
	log.Infof("datax job(%v) taskgruop(%v) start tasks", c.jobID, c.taskGroupID)
	for i := 0; i < len(taskConfigs); i++ {
		//从待执行队列加入运行队列
		te, ok := c.tasks.popRemainAndAddRun()
		if !ok {
			continue
		}
		//开始运行
		if err = c.startTaskExecer(te); err != nil {
			return
		}
	}
	log.Infof("datax job(%v) taskgruop(%v) manage tasks", c.jobID, c.taskGroupID)
	ticker := time.NewTicker(c.sleepInterval)
	defer ticker.Stop()
QueueLoop:
	//任务队列不为空
	for !c.tasks.isEmpty() {
		select {
		case <-c.ctx.Done():
			break QueueLoop
		default:
		}
		//从待执行队列加入运行队列
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
	log.Infof("datax job(%v) taskgruop(%v) wait tasks end", c.jobID, c.taskGroupID)
	// 等待所有任务携程结束
	c.wg.Wait()
	if c.ctx.Err() != nil {
		return c.ctx.Err()
	}

	return nil
}

//startTaskExecer 开始任务
func (c *Container) startTaskExecer(te *taskExecer) (err error) {
	log.Debugf("datax job(%v) taskgruop(%v) task(%v) push", c.jobID, c.taskGroupID, te.Key())
	c.wg.Add(1)
	var errChan <-chan error
	//将任务加入到调度器
	errChan, err = c.scheduler.Push(te)
	if err != nil {
		//错误发生时，从运行队列移除任务
		c.tasks.removeRun(te)
		c.wg.Done()
		return err
	}
	log.Debugf("datax job(%v) taskgruop(%v) task(%v) start", c.jobID, c.taskGroupID, te.Key())
	go func(te *taskExecer) {
		defer c.wg.Done()
		timer := time.NewTimer(c.retryInterval)
		defer timer.Stop()
		select {
		case err := <-errChan:
			//当失败时，重试次数不超过最大重试次数，写入任务是否支持失败冲时，这些决定写入任务是否冲时
			if err != nil && te.WriterSuportFailOverport() && te.AttemptCount() <= c.retryMaxCount {
				log.Debugf("datax job(%v) taskgruop(%v) task(%v) shutdown and retry. attemptCount: %v err: %v",
					c.jobID, c.taskGroupID, te.Key(), te.AttemptCount(), err)
				//关闭任务
				te.Shutdown()
				select {
				case <-timer.C:
				case <-c.ctx.Done():
				}
				//从运行队列移到待执行队列
				c.tasks.removeRunAndPushRemain(te)
			} else {
				log.Debugf("datax job(%v) taskgruop(%v) task(%v) end", c.jobID, c.taskGroupID, te.Key())
				//从任务调度器移除
				c.tasks.removeRun(te)
			}
		case <-c.ctx.Done():
		}
	}(te)
	return
}
