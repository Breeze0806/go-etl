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

package file

import (
	"context"
	"errors"
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
	record    element.Record
	createErr error
	onErr     error
}

func (m *mockFetchHandler) CreateRecord() (element.Record, error) {
	return element.NewDefaultRecord(), m.createErr
}

func (m *mockFetchHandler) OnRecord(record element.Record) error {
	m.record = record
	return m.onErr
}

type mockRows struct {
	n       int
	scanErr error
	err     error
}

func (m *mockRows) Next() bool {
	m.n++
	return m.n <= 1
}

func (m *mockRows) Scan() (columns []element.Column, err error) {
	columns = append(columns, element.NewDefaultColumn(element.NewStringColumnValue("mock"),
		"mock", 0))
	err = m.scanErr
	return
}

func (m *mockRows) Error() error {
	return m.err
}

func (m *mockRows) Close() error {
	return nil
}

type mockInStream struct {
	rows    Rows
	rowsErr error
}

func (m *mockInStream) Rows(conf *config.JSON) (rows Rows, err error) {
	return m.rows, m.rowsErr
}

func (m *mockInStream) Close() (err error) {
	return
}

type mockOpener struct {
	inStream InStream
	openErr  error
}

func (m *mockOpener) Open(filename string) (stream InStream, err error) {
	return m.inStream, m.openErr
}

func TestInStreamer_Read(t *testing.T) {
	UnregisterAllOpener()
	RegisterOpener("mock", &mockOpener{
		inStream: &mockInStream{
			rows: &mockRows{},
		},
	})

	type args struct {
		name    string
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
			args: args{
				name:    "mock",
				ctx:     context.TODO(),
				conf:    testJSONFromString("{}"),
				handler: &mockFetchHandler{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewInStreamer(tt.args.name, "")
			if err != nil {
				t.Fatal(err)
			}
			defer s.Close()
			if err := s.Read(tt.args.ctx, tt.args.conf, tt.args.handler); (err != nil) != tt.wantErr {
				t.Errorf("InStreamer.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if c, _ := tt.args.handler.(*mockFetchHandler).record.GetByName("mock"); c.String() != "mock" {
				t.Errorf("InStreamer.Read() fail")
			}
		})
	}
}

func TestInStreamer_ReadErr(t *testing.T) {
	UnregisterAllOpener()
	RegisterOpener("mockCreateErr", &mockOpener{
		inStream: &mockInStream{
			rows: &mockRows{},
		},
	})
	RegisterOpener("mockOnErr", &mockOpener{
		inStream: &mockInStream{
			rows: &mockRows{},
		},
	})
	RegisterOpener("mockRowsErr", &mockOpener{
		inStream: &mockInStream{
			rows:    &mockRows{},
			rowsErr: errors.New("mock error"),
		},
	})
	RegisterOpener("mockScanErr", &mockOpener{
		inStream: &mockInStream{
			rows: &mockRows{
				scanErr: errors.New("mock error"),
			},
		},
	})

	RegisterOpener("mockErr", &mockOpener{
		inStream: &mockInStream{
			rows: &mockRows{
				err: errors.New("mock error"),
			},
		},
	})
	type args struct {
		name    string
		ctx     context.Context
		conf    *config.JSON
		handler FetchHandler
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				name:    "mockRowsErr",
				ctx:     context.TODO(),
				conf:    testJSONFromString("{}"),
				handler: &mockFetchHandler{},
			},
			wantErr: true,
		},
		{
			name: "2",
			args: args{
				name:    "mockScanErr",
				ctx:     context.TODO(),
				conf:    testJSONFromString("{}"),
				handler: &mockFetchHandler{},
			},
			wantErr: true,
		},
		{
			name: "3",
			args: args{
				name:    "mockErr",
				ctx:     context.TODO(),
				conf:    testJSONFromString("{}"),
				handler: &mockFetchHandler{},
			},
			wantErr: true,
		},
		{
			name: "4",
			args: args{
				name: "mockCreateErr",
				ctx:  context.TODO(),
				conf: testJSONFromString("{}"),
				handler: &mockFetchHandler{
					createErr: errors.New("mock error"),
				},
			},
			wantErr: true,
		},
		{
			name: "5",
			args: args{
				name: "mockOnErr",
				ctx:  context.TODO(),
				conf: testJSONFromString("{}"),
				handler: &mockFetchHandler{
					onErr: errors.New("mock error"),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewInStreamer(tt.args.name, "")
			if err != nil {
				t.Fatal(err)
			}
			defer s.Close()
			if err := s.Read(tt.args.ctx, tt.args.conf, tt.args.handler); (err != nil) != tt.wantErr {
				t.Errorf("InStreamer.Read() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewInStreamerErr(t *testing.T) {
	UnregisterAllOpener()
	RegisterOpener("mockOpenErr", &mockOpener{
		inStream: nil,
		openErr:  errors.New("mock err"),
	})
	type args struct {
		name     string
		filename string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				name:     "mock",
				filename: "",
			},
			wantErr: true,
		},
		{
			name: "2",
			args: args{
				name:     "mockOpenErr",
				filename: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewInStreamer(tt.args.name, tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewInStreamer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_openerMap_registerErr(t *testing.T) {
	type args struct {
		name   string
		opener Opener
	}
	tests := []struct {
		name    string
		o       *openerMap
		args    args
		wantErr bool
	}{
		{
			name: "1",
			o: &openerMap{
				openers: make(map[string]Opener),
			},
			args: args{
				name:   "mock",
				opener: nil,
			},
			wantErr: true,
		},
		{
			name: "1",
			o: &openerMap{
				openers: map[string]Opener{
					"mock": &mockOpener{},
				},
			},
			args: args{
				name:   "mock",
				opener: &mockOpener{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.o.register(tt.args.name, tt.args.opener); (err != nil) != tt.wantErr {
				t.Errorf("openerMap.register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
