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
