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

package csv

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
				Index: "x",
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
				Index:  "1",
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
				Index: "1",
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

func TestNewConfig(t *testing.T) {
	type args struct {
		conf *config.JSON
	}
	tests := []struct {
		name    string
		args    args
		wantC   *Config
		wantErr bool
	}{

		{
			name: "2",
			args: args{
				conf: testJSONFromString(`{"delimiter":"12"}`),
			},
			wantErr: true,
		},

		{
			name: "3",
			args: args{
				conf: testJSONFromString(`{"column":[{"index":""}]}`),
			},
			wantErr: true,
		},
		{
			name: "4",
			args: args{
				conf: testJSONFromString(`{"encoding":"12"}`),
			},
			wantErr: true,
		},
		{
			name: "5",
			args: args{
				conf: testJSONFromString(`{"encoding":1}`),
			},
			wantErr: true,
		},
		{
			name: "6",
			args: args{
				conf: testJSONFromString(`{"encoding":"utf-8","column":[{"index":"1","type":"bool"}]}`),
			},
			wantC: &Config{
				Encoding: "utf-8",
				Columns: []Column{
					{
						Index: "1",
						Type:  "bool",
					},
				},
				Delimiter: ",",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotC, err := NewConfig(tt.args.conf)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotC, tt.wantC) {
				t.Errorf("NewConfig() = %v, want %v", gotC, tt.wantC)
			}
		})
	}
}
