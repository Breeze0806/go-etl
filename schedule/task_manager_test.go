// Copyright 2020 the go-etl Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
