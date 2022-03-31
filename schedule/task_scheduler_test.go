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
	"sync"
	"testing"
	"time"
)

func TestTaskSchduler_Once(t *testing.T) {
	schduler := NewTaskSchduler(2, 0)
	wait := make(chan struct{})
	waited := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			if i == 100 {
				close(wait)
				<-waited
			}
			schduler.Push(&mockTask{})
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-wait
		schduler.Stop()
		close(waited)
	}()
	wg.Wait()
}

func TestTaskSchduler_Multi(t *testing.T) {
	schduler := NewTaskSchduler(2, 0)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			schduler.Push(&mockTask{})
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			schduler.Push(&mockTask{})
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			schduler.Push(&mockTask{})
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			schduler.Stop()
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			schduler.Stop()
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			schduler.Size()
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			schduler.Size()
		}
	}()
	wg.Wait()
}

func TestTaskSchduler_Size(t *testing.T) {
	schduler := NewTaskSchduler(1, 0)
	schduler.Push(&mockTask{100 * time.Millisecond})
	if schduler.Size() != 1 {
		t.Errorf("Size() = %v want: 1", schduler.Size())
	}
	time.Sleep(1 * time.Second)
	if schduler.Size() != 0 {
		t.Errorf("Size() = %v want: 0", schduler.Size())
	}
}

func TestTaskSchduler_Stop(t *testing.T) {
	schduler := NewTaskSchduler(1, 0)
	schduler.Stop()
	_, err := schduler.Push(&mockTask{})
	if err != ErrClosed {
		t.Errorf("Push() = %v want: %v", err, ErrClosed)
	}
}
