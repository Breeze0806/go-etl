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
	close        chan struct{}
	destroy      sync.Once
	closed       int32
	size         *atomic.Int32
}

func NewTaskSchduler(workerNumer, chanCap int) *TaskSchduler {
	t := &TaskSchduler{
		taskWrappers: make(chan *taskWrapper, chanCap),
		close:        make(chan struct{}),
		closed:       0,
		size:         atomic.NewInt32(0),
	}
	t.wg.Add(1)
	for i := 0; i < workerNumer; i++ {
		go func() {
			defer t.wg.Done()
			t.processTask()
		}()
	}

	return t
}

func (t *TaskSchduler) Push(task Task) (<-chan error, error) {
	if goatomic.CompareAndSwapInt32(&t.closed, 1, 1) {
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
	case <-t.close:
		return nil, ErrClose
	}
}

func (t *TaskSchduler) Size() int32 {
	return t.size.Load()
}

func (t *TaskSchduler) Stop() {
	if !goatomic.CompareAndSwapInt32(&t.closed, 0, 1) {
		return
	}
	close(t.close)
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
		case <-t.close:
			return
		}
	}
}
