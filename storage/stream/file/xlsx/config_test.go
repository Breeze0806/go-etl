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

package xlsx

import (
	"reflect"
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

func TestColumn_validate(t *testing.T) {
	tests := []struct {
		name    string
		c       *Column
		wantErr bool
	}{
		{
			name: "1",
			c: &Column{
				Type:  "",
				Index: "1",
			},
			wantErr: true,
		},
		{
			name: "2",
			c: &Column{
				Type:  string(element.TypeTime),
				Index: "1",
			},
			wantErr: true,
		},
		{
			name: "3",
			c: &Column{
				Type:  string(element.TypeBigInt),
				Index: "1",
			},
			wantErr: true,
		},
		{
			name: "4",
			c: &Column{
				Type:  string(element.TypeBigInt),
				Index: "0",
			},
			wantErr: true,
		},
		{
			name: "5",
			c: &Column{
				Type:   string(element.TypeTime),
				Format: "yyyy-MM-dd",
				Index:  "A",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.validate(); (err != nil) != tt.wantErr {
				t.Errorf("Column.validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestColumn_index(t *testing.T) {
	tests := []struct {
		name  string
		c     *Column
		wantI int
	}{
		{
			name: "1",
			c: &Column{
				Index: "A",
			},
			wantI: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotI := tt.c.index(); gotI != tt.wantI {
				t.Errorf("Column.index() = %v, want %v", gotI, tt.wantI)
				return
			}
			if gotI := tt.c.index(); gotI != tt.wantI {
				t.Errorf("Column.index() = %v, want %v", gotI, tt.wantI)
			}
		})
	}
}

func TestColumn_layout(t *testing.T) {
	tests := []struct {
		name string
		c    *Column
		want string
	}{
		{
			name: "1",
			c: &Column{
				Format: "yyyy-MM-dd",
			},
			want: "2006-01-02",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.layout(); got != tt.want {
				t.Errorf("Column.layout() = %v, want %v", got, tt.want)
				return
			}

			if got := tt.c.layout(); got != tt.want {
				t.Errorf("Column.layout() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewInConfig(t *testing.T) {
	type args struct {
		conf *config.JSON
	}
	tests := []struct {
		name    string
		args    args
		wantC   *InConfig
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				conf: testJSONFromString("{}"),
			},
			wantErr: true,
		},

		{
			name: "3",
			args: args{
				conf: testJSONFromString(`{"sheet":"12","column":[{"index":""}]}`),
			},
			wantErr: true,
		},
		{
			name: "4",
			args: args{
				conf: testJSONFromString(`{"sheet":"sheet1","column":[{"index":"A","type":"bool"}]}`),
			},
			wantC: &InConfig{
				Sheet: "sheet1",
				Columns: []Column{
					{
						Index: "A",
						Type:  "bool",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotC, err := NewInConfig(tt.args.conf)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewInConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotC, tt.wantC) {
				t.Errorf("NewInConfig() = %v, want %v", gotC, tt.wantC)
			}
		})
	}
}

func TestNewOutConfig(t *testing.T) {
	type args struct {
		conf *config.JSON
	}
	tests := []struct {
		name    string
		args    args
		wantC   *OutConfig
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				conf: testJSONFromString("{}"),
			},
			wantErr: true,
		},

		{
			name: "2",
			args: args{
				conf: testJSONFromString(`{"sheets":["sheet1"]}`),
			},
			wantErr: true,
		},

		{
			name: "3",
			args: args{
				conf: testJSONFromString(`{"sheets":["sheet1"],"column":[{"index":""}]}`),
			},
			wantErr: true,
		},

		{
			name: "4",
			args: args{
				conf: testJSONFromString(`{"sheets":["sheet1"],"column":[{"index":"A","type":"bool"}]}`),
			},
			wantC: &OutConfig{
				Sheets: []string{"sheet1"},
				Columns: []Column{
					{
						Index: "A",
						Type:  "bool",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotC, err := NewOutConfig(tt.args.conf)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewOutConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotC, tt.wantC) {
				t.Errorf("NewOutConfig() = %v, want %v", gotC, tt.wantC)
			}
		})
	}
}
