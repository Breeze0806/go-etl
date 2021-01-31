package taskgroup

import "sync"

//taskManager 任务管理器
//toto 不知道为什么len(remain) + len(run) 无法实时任务数,其中主要是len(run)不准确
type taskManager struct {
	sync.Mutex

	remain []*taskExecer          //待执行队列
	run    map[string]*taskExecer //运行队列
	num    int                    //任务数
}

//newTaskManager 创建任务管理器
func newTaskManager() *taskManager {
	return &taskManager{
		run: make(map[string]*taskExecer),
	}
}

//isEmpty 任务管理器是否为空
func (t *taskManager) isEmpty() bool {
	t.Lock()
	defer t.Unlock()
	return t.lockedSize() == 0
}

//size 任务数，包含待执行和运行任务
func (t *taskManager) size() int {
	t.Lock()
	defer t.Unlock()
	return t.lockedSize()
}

//lockedSize 未加锁的任务数
func (t *taskManager) lockedSize() int {
	return t.num
}

//removeRunAndPushRemain 从运行队列移动到待执行队列
func (t *taskManager) removeRunAndPushRemain(te *taskExecer) {
	t.Lock()
	defer t.Unlock()
	t.lockedRemoveRun(te)
	t.lockedPushRemain(te)
}

//pushRemain 把任务加入待执行队列
func (t *taskManager) pushRemain(te *taskExecer) {
	t.Lock()
	defer t.Unlock()
	t.lockedPushRemain(te)

}

//removeRun 从运行队列移除出任务
func (t *taskManager) removeRun(te *taskExecer) {
	t.Lock()
	defer t.Unlock()
	t.lockedRemoveRun(te)
}

//popRemainAndAddRun 从待执行队列移到运行队列中
func (t *taskManager) popRemainAndAddRun() (te *taskExecer, ok bool) {
	t.Lock()
	defer t.Unlock()
	te, ok = t.lockedPopRemain()
	if ok {
		t.lockedAddRun(te)
	}
	return
}

//lockedRemoveRun 从运行队列移除任务
func (t *taskManager) lockedRemoveRun(te *taskExecer) {
	t.run[te.Key()] = nil
	delete(t.run, te.Key())
	t.num--
}

//lockedPushRemain 将任务加入到待执行队列
func (t *taskManager) lockedPushRemain(te *taskExecer) {
	t.remain = append(t.remain, te)
	t.num++
}

//lockedPushRemain 将任务加入到运行队列
func (t *taskManager) lockedAddRun(te *taskExecer) {
	t.run[te.Key()] = te
}

//lockedPopRemain 从待执行队列带出任务te，当代执行队列中没有值时，返回false
func (t *taskManager) lockedPopRemain() (te *taskExecer, ok bool) {
	if len(t.remain) == 0 {
		return nil, false
	}
	te, ok = t.remain[0], true
	t.remain, t.remain[0] = t.remain[1:], nil
	return
}
