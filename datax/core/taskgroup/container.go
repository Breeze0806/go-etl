package taskgroup

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/Breeze0806/go-etl/datax/common/config"
	coreconst "github.com/Breeze0806/go-etl/datax/common/config/core"
	"github.com/Breeze0806/go-etl/datax/common/plugin/loader"
	"github.com/Breeze0806/go-etl/datax/common/schedule"
	"github.com/Breeze0806/go-etl/datax/core"
	"github.com/Breeze0806/go-etl/datax/core/statistics/communication"
	"github.com/Breeze0806/go-etl/datax/core/taskgroup/runner"
	"github.com/Breeze0806/go-etl/datax/core/transport/channel"
	"github.com/Breeze0806/go-etl/datax/core/transport/exchange"
)

type Container struct {
	*core.BaseCotainer
	jobId       int64
	taskGroupId int64
	scheduler   *schedule.TaskSchduler
	wg          sync.WaitGroup
	tasks       *taskMap
	ctx         context.Context
}

func NewContainer(ctx context.Context, conf *config.Json) (c *Container, err error) {
	c = &Container{
		tasks: newTaskMap(),
		ctx:   ctx,
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
	var taskConfigs []*config.Json
	if taskConfigs, err = c.Config().GetArray(coreconst.DataxJobContent); err != nil {
		return err
	}
	c.scheduler = schedule.NewTaskSchduler(len(taskConfigs), len(taskConfigs))
	prefixKey := strconv.FormatInt(c.jobId, 10) + "-" + strconv.FormatInt(c.taskGroupId, 10)
	for i := range taskConfigs {
		var taskExecer *taskExecer

		taskExecer, err = newTaskExecer(c.ctx, taskConfigs[i], prefixKey, 0)
		if err != nil {
			return err
		}
		c.tasks.pushRemain(taskExecer)
	}
	sleepIntervalInMillSec := time.Duration(
		c.Config().GetInt64OrDefaullt(coreconst.DataxCoreContainerJobSleepinterval, 100)) * time.Millisecond
	for {
		if c.tasks.isEmpty() {
			break
		}
		te, ok := c.tasks.popRemainAndAddRun()
		if !ok {
			continue
		}
		c.wg.Add(1)
		var errChan <-chan error
		errChan, err = c.scheduler.Push(te)
		if err != nil {
			c.tasks.removeRun(te)
			c.wg.Done()
			return err
		}

		go func(te *taskExecer) {
			defer c.wg.Done()
			select {
			case err := <-errChan:
				if err != nil {
					c.tasks.removeRunAndPushRemain(te)
				} else {
					c.tasks.removeRun(te)
				}
			case <-c.ctx.Done():
			}
		}(te)

		select {
		case <-time.After(sleepIntervalInMillSec):
		case <-c.ctx.Done():
		}
	}
	c.wg.Wait()
	c.scheduler.Stop()
	return nil
}

type taskExecer struct {
	taskConf     *config.Json
	taskId       int64
	ctx          context.Context
	cancel       context.CancelFunc
	channel      *channel.Channel
	writerRunner runner.Runner
	readerRunner runner.Runner
	wg           sync.WaitGroup
	errors       chan error
	//todo: 初始化
	taskCommunication communication.Communication
	destroy           sync.Once
	key               string
}

func newTaskExecer(ctx context.Context, taskConf *config.Json, prefixKey string, attemptCount int) (t *taskExecer, err error) {
	t = &taskExecer{
		taskConf: taskConf,
		errors:   make(chan error),
	}
	t.ctx, t.cancel = context.WithCancel(ctx)
	t.channel, err = channel.NewChannel()
	if err != nil {
		return nil, err
	}

	t.taskId, err = taskConf.GetInt64(coreconst.TaskId)
	if err != nil {
		return nil, err
	}
	t.key = prefixKey + "-" + strconv.FormatInt(t.taskId, 10)
	name := ""
	name, err = taskConf.GetString(coreconst.JobReaderName)
	if err != nil {
		return nil, err
	}

	readTask, ok := loader.LoadReaderTask(name)
	if !ok {
		return nil, fmt.Errorf("reader task name (%v) does not exist", name)
	}
	exchanger := exchange.NewRecordExchangerWithoutTransformer(t.channel)
	t.readerRunner = runner.NewReader(ctx, readTask, exchanger)

	name, err = taskConf.GetString(coreconst.JobWriterName)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	writeTask, ok := loader.LoadWriterTask(name)
	if !ok {
		return nil, fmt.Errorf("writer task name (%v) does not exist", name)
	}
	t.writerRunner = runner.NewWriter(ctx, writeTask, exchanger)

	return
}

func (t *taskExecer) Start() {
	t.wg.Add(1)
	var writerWg sync.WaitGroup
	writerWg.Add(1)
	go func() {
		defer t.wg.Done()
		writerWg.Done()
		if err := t.writerRunner.Run(); err != nil {
			t.errors <- fmt.Errorf("writer task(%v) fail, err: %v", t.Key(), err)
		}
	}()
	writerWg.Wait()
	var readWg sync.WaitGroup
	t.wg.Add(1)
	readWg.Add(1)
	go func() {
		defer t.wg.Done()
		readWg.Done()
		if err := t.readerRunner.Run(); err != nil {
			t.errors <- fmt.Errorf("reader task(%v) fail, err: %v", t.Key(), err)
		}
	}()
	readWg.Wait()
}

func (t *taskExecer) Do() error {
	defer t.Shutdown()
	t.Start()
	var errors []error
RunLoop:
	for {
		select {
		case err := <-t.errors:
			errors = append(errors, err)
		case <-t.ctx.Done():
			break RunLoop
		}
	}

ErrorLoop:
	for {
		select {
		case err := <-t.errors:
			errors = append(errors, err)
		default:
			break ErrorLoop
		}
	}

	s := ""
	for i, v := range errors {
		if i > 0 {
			s += " "
		}
		s += v.Error()
	}
	if s != "" {
		return fmt.Errorf("%v", s)
	}
	return nil
}

func (t *taskExecer) Key() string {
	return t.key
}

func (t *taskExecer) Shutdown() {
	t.destroy.Do(func() {
		t.readerRunner.Shutdown()
		t.writerRunner.Shutdown()
		t.cancel()
		t.wg.Wait()
	})
}

type taskMap struct {
	sync.Mutex
	remain []*taskExecer
	run    map[string]*taskExecer
}

func newTaskMap() *taskMap {
	return &taskMap{
		run: make(map[string]*taskExecer),
	}
}

func (t *taskMap) isEmpty() bool {
	t.Lock()
	defer t.Unlock()
	return len(t.remain)+len(t.run) == 0
}

func (t *taskMap) removeRunAndPushRemain(te *taskExecer) {
	t.Lock()
	defer t.Unlock()
	t.lockedRemoveRun(te)
	t.lockedPushRemain(te)
}

func (t *taskMap) pushRemain(te *taskExecer) {
	t.Lock()
	defer t.Unlock()
	t.lockedPushRemain(te)
}

func (t *taskMap) removeRun(te *taskExecer) {
	t.Lock()
	defer t.Unlock()
	t.lockedRemoveRun(te)
}

func (t *taskMap) popRemainAndAddRun() (te *taskExecer, ok bool) {
	t.Lock()
	defer t.Unlock()
	te, ok = t.lockedPopRemain()
	if ok {
		t.lockedAddRun(te)
	}
	return
}

func (t *taskMap) lockedRemoveRun(te *taskExecer) {
	t.run[te.Key()] = nil
	delete(t.run, te.Key())
}

func (t *taskMap) lockedPushRemain(te *taskExecer) {
	t.remain = append(t.remain, te)
}

func (t *taskMap) lockedAddRun(te *taskExecer) {
	t.run[te.Key()] = te
}

func (t *taskMap) lockedPopRemain() (te *taskExecer, ok bool) {
	if len(t.remain) == 0 {
		return nil, false
	}
	te, ok = t.remain[0], true
	t.remain, t.remain[0] = t.remain[1:], nil
	return
}
