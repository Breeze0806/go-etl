package file

import (
	"context"
	"testing"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/element"
)

func testJSONFromString(json string) *config.JSON {
	conf, err := config.NewJSONFromString(json)
	if err != nil {
		panic(err)
	}
	return conf
}

type mockFetchHandler struct {
	record element.Record
}

func (m *mockFetchHandler) CreateRecord() (element.Record, error) {
	return element.NewDefaultRecord(), nil
}

func (m *mockFetchHandler) OnRecord(record element.Record) error {
	m.record = record
	return nil
}

type mockRows struct {
	n int
}

func (m *mockRows) Next() bool {
	m.n++
	return m.n <= 1
}

func (m *mockRows) Scan() (columns []element.Column, err error) {
	columns = append(columns, element.NewDefaultColumn(element.NewStringColumnValue("mock"),
		"mock", 0))
	return
}

func (m *mockRows) Error() error {
	return nil
}

func (m *mockRows) Close() error {
	return nil
}

type mockInStream struct {
}

func (m *mockInStream) Rows(conf *config.JSON) (rows Rows, err error) {
	return &mockRows{}, nil
}

func (m *mockInStream) Close() (err error) {
	return
}

type mockOpener struct {
}

func (m *mockOpener) Open(filename string) (stream InStream, err error) {
	return &mockInStream{}, nil
}

func TestInStreamer_Read(t *testing.T) {
	RegisterOpener("mock", &mockOpener{})

	s, err := NewInStreamer("mock", "")
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		ctx     context.Context
		conf    *config.JSON
		handler FetchHandler
	}
	tests := []struct {
		name    string
		s       *InStreamer
		args    args
		wantErr bool
	}{
		{
			name: "1",
			s:    s,
			args: args{
				ctx:     context.TODO(),
				conf:    testJSONFromString("{}"),
				handler: &mockFetchHandler{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer tt.s.Close()
			if err := tt.s.Read(tt.args.ctx, tt.args.conf, tt.args.handler); (err != nil) != tt.wantErr {
				t.Errorf("InStreamer.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if c, _ := tt.args.handler.(*mockFetchHandler).record.GetByName("mock"); c.String() != "mock" {
				t.Errorf("InStreamer.Read() fail")
			}
		})
	}
}
