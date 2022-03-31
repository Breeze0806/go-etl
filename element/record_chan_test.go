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

package element

import (
	"context"
	"fmt"
	"sync"
	"testing"
)

type mockRecord struct {
	i int
}

func (m *mockRecord) Add(Column) error {
	return nil
}

func (m *mockRecord) GetByIndex(i int) (Column, error) {
	return nil, nil
}

func (m *mockRecord) GetByName(name string) (Column, error) {
	return nil, nil
}

func (m *mockRecord) Set(i int, c Column) error {
	return nil
}

func (m *mockRecord) ColumnNumber() int {
	return 0
}

func (m *mockRecord) ByteSize() int64 {
	return 0
}

func (m *mockRecord) MemorySize() int64 {
	return 0
}
func (m *mockRecord) String() string {
	return ""
}

func TestRecordChan_MutilPushBackPopFront(t *testing.T) {
	c := NewRecordChan(context.TODO())
	defer c.Close()
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 1500; i++ {
			c.PushBack(&mockRecord{i: i})
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 1500; i++ {
			c.PushBack(&mockRecord{i: i})
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; {
			if _, ok := c.PopFront(); ok {
				i++
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; {
			if _, ok := c.PopFront(); ok {
				i++
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; {
			if _, ok := c.PopFront(); ok {
				i++
			}
		}
	}()
	wg.Wait()
	if c.Buffered() != 0 {
		t.Errorf("Buffered() = %v want 0", c.Buffered())
	}
}

func TestRecordChan_MutilPushBackAllPopFrontAll(t *testing.T) {
	c := NewRecordChan(context.TODO())
	defer c.Close()
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		i := 0
		c.PushBackAll(func() (Record, error) {
			if i == 1500 {
				return nil, fmt.Errorf("test over")
			}
			i++
			return &mockRecord{i: i}, nil
		})
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		i := 0
		c.PushBackAll(func() (Record, error) {
			if i == 1500 {
				return nil, fmt.Errorf("test over")
			}
			i++
			return &mockRecord{i: i}, nil
		})
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		i := 0
		c.PopFrontAll(func(Record) error {
			if i == 999 {
				return fmt.Errorf("test over")
			}
			i++
			return nil
		})
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		i := 0
		c.PopFrontAll(func(Record) error {
			if i == 999 {
				return fmt.Errorf("test over")
			}
			i++
			return nil
		})

	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		i := 0
		c.PopFrontAll(func(Record) error {
			if i == 999 {
				return fmt.Errorf("test over")
			}
			i++
			return nil
		})
	}()
	wg.Wait()
	if c.Buffered() != 0 {
		t.Errorf("Buffered() = %v want 0", c.Buffered())
	}
}

func TestRecordChan_Close(t *testing.T) {
	c := NewRecordChan(context.TODO())
	var wg sync.WaitGroup
	wg.Add(1)
	ok := true
	go func() {
		defer wg.Done()
		_, ok = c.PopFront()
	}()
	wg.Add(1)
	var err error
	go func() {
		defer wg.Done()
		err = c.PopFrontAll(func(Record) error {
			return nil
		})
	}()
	c.Close()
	wg.Wait()

	if ok {
		t.Errorf("PopFront %v want: false", ok)
	}

	if err != nil {
		t.Errorf("PopFrontAll %v wantErr : false", err)
	}

	defer func() {
		myerr := recover()
		if myerr == nil {
			t.Errorf("pushback have no error")
		}
		t.Log(myerr)
	}()
	c.PushBack(NewDefaultRecord())
}

func TestRecordChan_Order(t *testing.T) {
	c := NewRecordChan(context.TODO())
	defer c.Close()
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		i := 0
		c.PushBackAll(func() (Record, error) {
			i++
			if i == 1001 {
				return nil, fmt.Errorf("test over")
			}
			return &mockRecord{i: i}, nil
		})
	}()
	var err error
	wg.Add(1)
	go func() {
		defer wg.Done()
		i := 0
		c.PopFrontAll(func(r Record) error {
			i++
			if r.(*mockRecord).i != i {
				err = fmt.Errorf("got ：%v want: %v", r.(*mockRecord).i, i)
			}
			if i == 1000 {
				return fmt.Errorf("test over")
			}
			return nil
		})
	}()
	wg.Wait()
	if err != nil {
		t.Error(err)
	}
}
