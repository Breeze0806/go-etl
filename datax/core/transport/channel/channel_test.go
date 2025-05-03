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

package channel

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/element"
)

type mockRecord struct {
	*element.DefaultRecord

	n int64
}

func (m *mockRecord) ByteSize() int64 {
	return m.n
}
func TestChannel_PushPop(t *testing.T) {
	ch := NewChannel(context.TODO(), nil)
	defer ch.Close()
	if !ch.IsEmpty() {
		t.Errorf("IsEmpty() = %v want true", ch.IsEmpty())
	}

	if n, _ := ch.Push(element.NewDefaultRecord()); n != 1 {
		t.Errorf("Push() = %v want 1", n)
	}
	if n := ch.PushTerminate(); n != 2 {
		t.Errorf("Push() = %v want 2", n)
	}
	if _, ok := ch.Pop(); !ok {
		t.Errorf("Pop() = %v want true", ok)
	}
}

func TestChannel_PushAllPopAll(t *testing.T) {
	ch := NewChannel(context.TODO(), nil)
	defer ch.Close()
	if !ch.IsEmpty() {
		t.Errorf("IsEmpty() = %v want true", ch.IsEmpty())
	}
	i := 0
	if err := ch.PushAll(func() (element.Record, error) {
		if i == 2 {
			return nil, fmt.Errorf("test over")
		}
		i++
		return element.NewDefaultRecord(), nil
	}); err == nil {
		t.Errorf("PushAll() = %v want not nil", err)
	}
	ch.Close()
	if err := ch.PopAll(func(element.Record) error {
		return nil
	}); err != nil {
		t.Errorf("PopAll() = %v want nil", err)
	}
}

func TestChannelWithRateLimit(t *testing.T) {
	conf, _ := config.NewJSONFromString(`{
		"core":{
			"transport":{
				"channel":{
					"speed":{
						"byte":10000,
						"record":10
					}
				}
			}
		}
	}`)
	want := 1000
	b := 100
	ch := NewChannel(context.TODO(), conf)
	defer ch.Close()
	var wg sync.WaitGroup
	wg.Add(1)
	n := 0
	go func() {
		defer wg.Done()
		for {
			r, _ := ch.Pop()
			switch r.(type) {
			case *element.TerminateRecord:
				return
			}
			n++
		}
	}()
	for i := 0; i < want; i++ {
		ch.Push(&mockRecord{
			DefaultRecord: element.NewDefaultRecord(),
			n:             int64(b),
		})
	}
	ch.PushTerminate()
	wg.Wait()

	if n != want {
		t.Errorf("want:%v n:%v", want, n)
	}

	if ch.StatsJSON().TotalByte != int64(b*want) {
		t.Errorf("TotalByte:%v want:%v", ch.StatsJSON().TotalByte, b*want)
	}

	if ch.StatsJSON().TotalRecord != int64(want) {
		t.Errorf("TotalRecord:%v want:%v", ch.StatsJSON().TotalRecord, want+1)
	}
}

func TestChannelWithRateLimit_Err(t *testing.T) {
	conf, _ := config.NewJSONFromString(`{
		"core":{
			"transport":{
				"channel":{
					"speed":{
						"byte":10000,
						"record":10
					}
				}
			}
		}
	}`)
	want := 1000
	ch := NewChannel(context.TODO(), conf)
	defer ch.Close()
	for i := 0; i < want; i++ {
		_, err := ch.Push(&mockRecord{
			DefaultRecord: element.NewDefaultRecord(),
			n:             int64(100000),
		})
		if err == nil {
			t.Fatal("want error back")
		}
	}
}
