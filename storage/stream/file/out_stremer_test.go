package file

import (
	"testing"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/element"
)

type mockStreamWriter struct {
	record element.Record
}

func (m *mockStreamWriter) Write(record element.Record) (err error) {
	m.record = record
	return
}

func (m *mockStreamWriter) Flush() (err error) {
	return
}

func (m *mockStreamWriter) Close() (err error) {
	return
}

type mockOutStream struct {
}

func (m *mockOutStream) Writer(conf *config.JSON) (writer StreamWriter, err error) {
	return &mockStreamWriter{}, nil
}

func (m *mockOutStream) Close() (err error) {
	return
}

type mockCreater struct {
}

func (m *mockCreater) Create(filename string) (stream OutStream, err error) {
	return &mockOutStream{}, nil
}

func TestOutStreamer_Write(t *testing.T) {
	RegisterCreater("mock", &mockCreater{})

	s, err := NewOutStreamer("mock", "")
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		conf *config.JSON
	}
	tests := []struct {
		name    string
		s       *OutStreamer
		args    args
		wantErr bool
	}{
		{
			name: "1",
			s:    s,
			args: args{
				conf: testJSONFromString("{}"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer tt.s.Close()
			got, err := tt.s.Writer(tt.args.conf)
			if (err != nil) != tt.wantErr {
				t.Errorf("OutStreamer.Writer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			defer got.Close()
			record := element.NewDefaultRecord()
			record.Add(element.NewDefaultColumn(element.NewStringColumnValue("mock"),
				"mock", 0))
			defer got.Flush()
			if err = got.Write(record); err != nil {
				t.Errorf("Write() error = %v", err)
				return
			}
			if c, _ := got.(*mockStreamWriter).record.GetByName("mock"); c.String() != "mock" {
				t.Errorf("InStreamer.write() fail")
				return
			}
		})
	}
}
