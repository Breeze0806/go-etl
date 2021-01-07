package schedule

import (
	"errors"
	"sync"
	goatomic "sync/atomic"

	"go.uber.org/atomic"
)

var (
	ErrClose = errors.New("task schduler closed")
)

type taskWrapper struct {
	task   Task
	result chan error
}

type TaskSchduler struct {
	taskWrappers chan *taskWrapper
	wg           sync.WaitGroup
	stop         chan struct{}
	stopped      int32
	size         *atomic.Int32
}

func NewTaskSchduler(workerNumer, chanCap int) *TaskSchduler {
	t := &TaskSchduler{
		taskWrappers: make(chan *taskWrapper, chanCap),
		stop:         make(chan struct{}),
		stopped:      0,
		size:         atomic.NewInt32(0),
	}

	for i := 0; i < workerNumer; i++ {
		t.wg.Add(1)
		go func() {
			defer t.wg.Done()
			t.processTask()
		}()
	}

	return t
}

func (t *TaskSchduler) Push(task Task) (<-chan error, error) {
	if goatomic.CompareAndSwapInt32(&t.stopped, 1, 1) {
		return nil, ErrClose
	}
	tw := &taskWrapper{
		task:   task,
		result: make(chan error, 1),
	}

	select {
	case t.taskWrappers <- tw:
		t.size.Inc()
		return tw.result, nil
	case <-t.stop:
		return nil, ErrClose
	}
}

func (t *TaskSchduler) Size() int32 {
	return t.size.Load()
}

func (t *TaskSchduler) Stop() {
	if !goatomic.CompareAndSwapInt32(&t.stopped, 0, 1) {
		return
	}
	close(t.stop)
	t.wg.Wait()
}

func (t *TaskSchduler) processTask() {
	for {
		select {
		case tw, ok := <-t.taskWrappers:
			if !ok {
				return
			}
			tw.result <- tw.task.Do()
			t.size.Dec()
		case <-t.stop:
			return
		}
	}
}
