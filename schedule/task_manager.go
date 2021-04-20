package schedule

import "sync"

//MappedTaskManager 任务管理器
//toto 不知道为什么len(remain) + len(run) 无法实时任务数,其中主要是len(run)不准确
type MappedTaskManager struct {
	sync.Mutex

	remain []MappedTask          //待执行队列
	run    map[string]MappedTask //运行队列
	num    int                   //任务数
}

//NewTaskManager 创建任务管理器
func NewTaskManager() *MappedTaskManager {
	return &MappedTaskManager{
		run: make(map[string]MappedTask),
	}
}

//IsEmpty 任务管理器是否为空
func (t *MappedTaskManager) IsEmpty() bool {
	t.Lock()
	defer t.Unlock()
	return t.lockedSize() == 0
}

//Size 任务数，包含待执行和运行任务
func (t *MappedTaskManager) Size() int {
	t.Lock()
	defer t.Unlock()
	return t.lockedSize()
}

//lockedSize 未加锁的任务数
func (t *MappedTaskManager) lockedSize() int {
	return t.num
}

//RemoveRunAndPushRemain 从运行队列移动到待执行队列
func (t *MappedTaskManager) RemoveRunAndPushRemain(task MappedTask) {
	t.Lock()
	defer t.Unlock()
	t.lockedRemoveRun(task)
	t.lockedPushRemain(task)
}

//PushRemain 把任务加入待执行队列
func (t *MappedTaskManager) PushRemain(task MappedTask) {
	t.Lock()
	defer t.Unlock()
	t.lockedPushRemain(task)
}

//RemoveRun 从运行队列移除出任务
func (t *MappedTaskManager) RemoveRun(task MappedTask) {
	t.Lock()
	defer t.Unlock()
	t.lockedRemoveRun(task)
}

//PopRemainAndAddRun 从待执行队列移到运行队列中
func (t *MappedTaskManager) PopRemainAndAddRun() (task MappedTask, ok bool) {
	t.Lock()
	defer t.Unlock()
	task, ok = t.lockedPopRemain()
	if ok {
		t.lockedAddRun(task)
	}
	return
}

//lockedRemoveRun 从运行队列移除任务
func (t *MappedTaskManager) lockedRemoveRun(task MappedTask) {
	t.run[task.Key()] = nil
	delete(t.run, task.Key())
	t.num--
}

//lockedPushRemain 将任务加入到待执行队列
func (t *MappedTaskManager) lockedPushRemain(task MappedTask) {
	t.remain = append(t.remain, task)
	t.num++
}

//lockedPushRemain 将任务加入到运行队列
func (t *MappedTaskManager) lockedAddRun(task MappedTask) {
	t.run[task.Key()] = task
}

//lockedPopRemain 从待执行队列带出任务te，当代执行队列中没有值时，返回false
func (t *MappedTaskManager) lockedPopRemain() (task MappedTask, ok bool) {
	if len(t.remain) == 0 {
		return nil, false
	}
	task, ok = t.remain[0], true
	t.remain, t.remain[0] = t.remain[1:], nil
	return
}
