package taskgroup

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/Breeze0806/go-etl/datax/common/config"
	coreconst "github.com/Breeze0806/go-etl/datax/common/config/core"
	"github.com/Breeze0806/go-etl/datax/common/plugin/loader"
	"github.com/Breeze0806/go-etl/datax/common/spi/writer"
	"github.com/Breeze0806/go-etl/datax/core/statistics/communication"
	"github.com/Breeze0806/go-etl/datax/core/taskgroup/runner"
	"github.com/Breeze0806/go-etl/datax/core/transport/channel"
	"github.com/Breeze0806/go-etl/datax/core/transport/exchange"
	"go.uber.org/atomic"
)

type taskExecer struct {
	taskConf     *config.Json
	taskId       int64
	ctx          context.Context
	channel      *channel.Channel
	writerRunner runner.Runner
	readerRunner runner.Runner
	wg           sync.WaitGroup
	errors       chan error
	//todo: 初始化
	taskCommunication communication.Communication
	destroy           sync.Once
	key               string

	cancalMutex  sync.Mutex
	cancel       context.CancelFunc
	attemptCount *atomic.Int32
}

func newTaskExecer(ctx context.Context, taskConf *config.Json, prefixKey string, attemptCount int) (t *taskExecer, err error) {
	t = &taskExecer{
		taskConf:     taskConf,
		errors:       make(chan error, 2),
		ctx:          ctx,
		attemptCount: atomic.NewInt32(int32(attemptCount)),
	}
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
	t.readerRunner = runner.NewReader(readTask, exchanger, t.key)

	name, err = taskConf.GetString(coreconst.JobWriterName)
	if err != nil {
		return nil, err
	}

	writeTask, ok := loader.LoadWriterTask(name)
	if !ok {
		return nil, fmt.Errorf("writer task name (%v) does not exist", name)
	}
	t.writerRunner = runner.NewWriter(writeTask, exchanger, t.key)

	return
}

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
			t.errors <- fmt.Errorf("writer task(%v) fail, err: %v", t.Key(), err)
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
			t.errors <- fmt.Errorf("reader task(%v) fail, err: %v", t.Key(), err)
		}
	}()
	readerWg.Wait()
}

func (t *taskExecer) AttemptCount() int32 {
	return t.attemptCount.Load()
}

func (t *taskExecer) Do() error {
	log.Debugf("taskExecer %v start to do", t.key)
	defer func() {
		t.attemptCount.Inc()
		log.Debugf("taskExecer %v end to do", t.key)
	}()
	t.Start()
	log.Debugf("taskExecer %v do wait runner stop", t.key)
	t.wg.Wait()

	var errs []error
ErrorLoop:
	for {
		select {
		case err := <-t.errors:
			errs = append(errs, err)
		default:
			break ErrorLoop
		}
	}

	s := ""
	for i, v := range errs {
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

func (t *taskExecer) WriterSuportFailOverport() bool {
	return t.writerRunner.Plugin().(writer.Task).SupportFailOver()
}

func (t *taskExecer) Shutdown() {
	log.Debugf("taskExecer %v starts to shutdown", t.key)
	defer log.Debugf("taskExecer %v ends to shutdown", t.key)
	t.cancalMutex.Lock()
	if t.cancel != nil {
		t.cancel()
	}

	t.cancalMutex.Unlock()
	log.Debugf("taskExecer %v shutdown wait runner stop", t.key)
	t.wg.Wait()

	t.readerRunner.Shutdown()
	t.writerRunner.Shutdown()

}
