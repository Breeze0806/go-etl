package schedule

import (
	"math"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func Test_taskManager(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	m := NewTaskManager()
	for i := 0; i < 10000; i++ {
		m.PushRemain(&mockMappedTask{
			taskID: int64(i),
		})
	}
	var wg sync.WaitGroup
	for !m.IsEmpty() {
		task, ok := m.PopRemainAndAddRun()
		if !ok {
			continue
		}
		wg.Add(1)
		go func(task MappedTask) {
			defer wg.Done()
			x := rand.Int31n(math.MaxInt16)
			if x < math.MaxInt16/2 {
				m.RemoveRunAndPushRemain(task)
			} else {
				m.RemoveRun(task)
			}
		}(task)
	}
	wg.Wait()
	if m.Size() != 0 {
		t.Errorf("size() = %v want: 0", m.Size())
	}
}
