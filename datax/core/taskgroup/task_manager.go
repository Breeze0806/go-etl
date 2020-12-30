package taskgroup

import "sync"

type taskManager struct {
	sync.Mutex
	remain []*taskExecer
	run    map[string]*taskExecer
	num    int
}

func newTaskManager() *taskManager {
	return &taskManager{
		run: make(map[string]*taskExecer),
	}
}

func (t *taskManager) isEmpty() bool {
	t.Lock()
	defer t.Unlock()
	return t.lockedSize() == 0
}

func (t *taskManager) size() int {
	t.Lock()
	defer t.Unlock()
	return t.lockedSize()
}

func (t *taskManager) lockedSize() int {
	return t.num
}

func (t *taskManager) removeRunAndPushRemain(te *taskExecer) {
	t.Lock()
	defer t.Unlock()
	t.lockedRemoveRun(te)
	t.lockedPushRemain(te)
}

func (t *taskManager) pushRemain(te *taskExecer) {
	t.Lock()
	defer t.Unlock()
	t.lockedPushRemain(te)

}

func (t *taskManager) removeRun(te *taskExecer) {
	t.Lock()
	defer t.Unlock()
	t.lockedRemoveRun(te)
}

func (t *taskManager) popRemainAndAddRun() (te *taskExecer, ok bool) {
	t.Lock()
	defer t.Unlock()
	te, ok = t.lockedPopRemain()
	if ok {
		t.lockedAddRun(te)
	}
	return
}

func (t *taskManager) lockedRemoveRun(te *taskExecer) {
	t.run[te.Key()] = nil
	delete(t.run, te.Key())
	t.num--
}

func (t *taskManager) lockedPushRemain(te *taskExecer) {
	t.remain = append(t.remain, te)
	t.num++
}

func (t *taskManager) lockedAddRun(te *taskExecer) {
	t.run[te.Key()] = te
}

func (t *taskManager) lockedPopRemain() (te *taskExecer, ok bool) {
	if len(t.remain) == 0 {
		return nil, false
	}
	te, ok = t.remain[0], true
	t.remain, t.remain[0] = t.remain[1:], nil
	return
}
