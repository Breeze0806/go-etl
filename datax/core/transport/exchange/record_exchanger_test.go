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

package exchange

import (
	"context"
	"sync"
	"testing"

	"github.com/Breeze0806/go-etl/datax/core/transport/channel"
	"github.com/Breeze0806/go-etl/element"
)

type mockRecord struct {
	i int
}

func (m *mockRecord) Add(element.Column) error {
	return nil
}

func (m *mockRecord) GetByIndex(i int) (element.Column, error) {
	return nil, nil
}

func (m *mockRecord) GetByName(name string) (element.Column, error) {
	return nil, nil
}

func (m *mockRecord) Set(i int, c element.Column) error {
	return nil
}

func (m *mockRecord) Put(c element.Column) error {
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
func TestRecordExchanger(t *testing.T) {
	ch := channel.NewChannel(context.TODO(), nil)
	defer ch.Close()
	re := NewRecordExchangerWithoutTransformer(ch)
	defer re.Shutdown()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 1; i <= 1000; i++ {
			re.SendWriter(&mockRecord{
				i: i,
			})
		}
	}()

	for i := 1; i <= 1000; i++ {
		r, _ := re.GetFromReader()
		if r.(*mockRecord).i != i {
			t.Errorf("GetFromReader() = %v  want %v", r.(*mockRecord).i, i)
		}
	}
	wg.Wait()
	re.Flush()
	re.Terminate()
	_, err := re.GetFromReader()
	if err != ErrTerminate {
		t.Errorf("GetFromReader() err = %v  want %v", err, ErrTerminate)
	}

	re.Shutdown()

	r, _ := re.CreateRecord()
	err = re.SendWriter(r)
	if err != ErrShutdown {
		t.Errorf("GetFromReader() err = %v  want %v", err, ErrTerminate)
	}
	_, err = re.GetFromReader()
	if err != ErrShutdown {
		t.Errorf("GetFromReader() err = %v  want %v", err, ErrTerminate)
	}
}
