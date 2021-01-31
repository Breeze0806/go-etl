package taskgroup

import (
	"math"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func Test_taskManager(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	m := newTaskManager()
	for i := 0; i < 10000; i++ {
		m.pushRemain(&taskExecer{
			taskID: int64(i),
		})
	}
	var wg sync.WaitGroup
	for !m.isEmpty() {
		te, ok := m.popRemainAndAddRun()
		if !ok {
			continue
		}
		wg.Add(1)
		go func(te *taskExecer) {
			defer wg.Done()
			x := rand.Int31n(math.MaxInt16)
			if x < math.MaxInt16/2 {
				m.removeRunAndPushRemain(te)
			} else {
				m.removeRun(te)
			}
		}(te)
	}
	wg.Wait()
	if m.size() != 0 {
		t.Errorf("size() = %v want: 0", m.size())
	}
}
