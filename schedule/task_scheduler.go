package schedule

import (
	"errors"
	"sync"
	goatomic "sync/atomic"

	"go.uber.org/atomic"
)

//已关闭错误
var (
	ErrClosed = errors.New("task schduler closed")
)

type taskWrapper struct {
	task   Task
	result chan error
}

//TaskSchduler 任务调度器
type TaskSchduler struct {
	taskWrappers chan *taskWrapper //待执行任务队列
	wg           sync.WaitGroup
	stop         chan struct{} //停止信号
	stopped      int32         //停止标识
	size         *atomic.Int32 //待执行队列大小
}

//NewTaskSchduler 根据执行者数workerNumer，待执行队列容量cao生成任务调度器
func NewTaskSchduler(workerNumer, cap int) *TaskSchduler {
	t := &TaskSchduler{
		taskWrappers: make(chan *taskWrapper, cap),
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

//Push 将任务task加入队列，获得执行结果通知信道，在已关闭时报错
func (t *TaskSchduler) Push(task Task) (<-chan error, error) {
	if goatomic.CompareAndSwapInt32(&t.stopped, 1, 1) {
		return nil, ErrClosed
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
		return nil, ErrClosed
	}
}

//Size 待执行队列大小
func (t *TaskSchduler) Size() int32 {
	return t.size.Load()
}

//Stop 停止任务调度器
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
