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
	"testing"

	"github.com/Breeze0806/go-etl/element"
)

func TestChannel_PushPop(t *testing.T) {
	ch, _ := NewChannel(context.TODO())
	defer ch.Close()
	if !ch.IsEmpty() {
		t.Errorf("IsEmpty() = %v want true", ch.IsEmpty())
	}

	if n := ch.Push(element.NewDefaultRecord()); n != 1 {
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
	ch, _ := NewChannel(context.TODO())
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
