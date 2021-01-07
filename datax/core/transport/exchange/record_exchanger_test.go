package exchange

import (
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

func (m *mockRecord) ColumnNumber() int {
	return 0
}

func (m *mockRecord) ByteSize() int64 {
	return 0
}

func (m *mockRecord) MemorySize() int64 {
	return 0
}
func TestRecordExchanger(t *testing.T) {
	ch, _ := channel.NewChannel()
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
