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
	"encoding/json"
	"fmt"
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
				conf: testJSONFromString(`{"delimiter":"12"}`),
			},
			wantErr: true,
		},

		{
			name: "2",
			args: args{
				conf: testJSONFromString(`{"column":[{"index":""}]}`),
			},
			wantErr: true,
		},
		{
			name: "3",
			args: args{
				conf: testJSONFromString(`{"encoding":"12"}`),
			},
			wantErr: true,
		},
		{
			name: "4",
			args: args{
				conf: testJSONFromString(`{"encoding":1}`),
			},
			wantErr: true,
		},
		{
			name: "5",
			args: args{
				conf: testJSONFromString(`{"startRow":-1}`),
			},
			wantErr: true,
		},
		{
			name: "6",
			args: args{
				conf: testJSONFromString(`{"comment":"as"}`),
			},
			wantErr: true,
		},
		{
			name: "7",
			args: args{
				conf: testJSONFromString(`{"encoding":"utf-8","column":[{"index":"1","type":"bool"}],"delimiter":"\u0010"}`),
			},
			wantC: &InConfig{
				Encoding: "utf-8",
				Columns: []Column{
					{
						Index: "1",
						Type:  "bool",
					},
				},
				Delimiter: "\u0010",
			},
		},
		{
			name: "8",
			args: args{
				conf: testJSONFromString(`{"encoding":"gbk","column":[{"index":"1","type":"bool"}]}`),
			},
			wantC: &InConfig{
				Encoding: "gbk",
				Columns: []Column{
					{
						Index: "1",
						Type:  "bool",
					},
				},
			},
		},
		{
			name: "9",
			args: args{
				conf: testJSONFromString(`{"column":[{"index":"1","type":"bool"}]}`),
			},
			wantC: &InConfig{
				Columns: []Column{
					{
						Index: "1",
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
				t.Errorf("NewConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotC, tt.wantC) {
				t.Errorf("NewConfig() = %v, want %v", gotC, tt.wantC)
			}
		})
	}
}

func TestMarshalInConfig(t *testing.T) {
	type args struct {
		conf *InConfig
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "1",
			args: args{
				conf: &InConfig{
					Delimiter: string([]byte{0x10}),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.args.conf)
			if err != nil {
				t.Fatalf("Marshal fail. err: %v", err)
			}
			fmt.Println(string(data))
		})
	}
}

func TestConfig_startLine(t *testing.T) {
	tests := []struct {
		name string
		c    *InConfig
		want int
	}{
		{
			name: "1",
			c:    &InConfig{},
			want: 1,
		},
		{
			name: "2",
			c: &InConfig{
				StartRow: 2,
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.startRow(); got != tt.want {
				t.Errorf("Config.startLine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_encoding(t *testing.T) {
	tests := []struct {
		name string
		c    *InConfig
		want string
	}{
		{
			name: "1",
			c:    &InConfig{},
			want: "utf-8",
		},
		{
			name: "2",
			c: &InConfig{
				Encoding: "gbk",
			},
			want: "gbk",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.encoding(); got != tt.want {
				t.Errorf("Config.encoding() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_comment(t *testing.T) {
	tests := []struct {
		name string
		c    *InConfig
		want rune
	}{
		{
			name: "1",
			c:    &InConfig{},
		},
		{
			name: "2",
			c: &InConfig{
				Comment: "#",
			},
			want: rune('#'),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.comment(); got != tt.want {
				t.Errorf("Config.comment() = %v, want %v", got, tt.want)
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
				conf: testJSONFromString(`{"delimiter":"12"}`),
			},
			wantErr: true,
		},

		{
			name: "2",
			args: args{
				conf: testJSONFromString(`{"column":[{"index":""}]}`),
			},
			wantErr: true,
		},
		{
			name: "3",
			args: args{
				conf: testJSONFromString(`{"encoding":"12"}`),
			},
			wantErr: true,
		},
		{
			name: "4",
			args: args{
				conf: testJSONFromString(`{"encoding":1}`),
			},
			wantErr: true,
		},
		{
			name: "7",
			args: args{
				conf: testJSONFromString(`{"encoding":"utf-8","column":[{"index":"1","type":"bool"}],"delimiter":"\u0010"}`),
			},
			wantC: &OutConfig{
				Encoding: "utf-8",
				Columns: []Column{
					{
						Index: "1",
						Type:  "bool",
					},
				},
				Delimiter: "\u0010",
			},
		},
		{
			name: "8",
			args: args{
				conf: testJSONFromString(`{"encoding":"gbk","column":[{"index":"1","type":"bool"}]}`),
			},
			wantC: &OutConfig{
				Encoding: "gbk",
				Columns: []Column{
					{
						Index: "1",
						Type:  "bool",
					},
				},
			},
		},
		{
			name: "9",
			args: args{
				conf: testJSONFromString(`{"column":[{"index":"1","type":"bool"}]}`),
			},
			wantC: &OutConfig{
				Columns: []Column{
					{
						Index: "1",
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

func TestOutConfig_encoding(t *testing.T) {
	tests := []struct {
		name string
		c    *OutConfig
		want string
	}{
		{
			name: "1",
			c:    &OutConfig{},
			want: "utf-8",
		},
		{
			name: "2",
			c: &OutConfig{
				Encoding: "gbk",
			},
			want: "gbk",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.encoding(); got != tt.want {
				t.Errorf("OutConfig.encoding() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOutConfig_comma(t *testing.T) {
	tests := []struct {
		name string
		c    *OutConfig
		want rune
	}{
		{
			name: "1",
			c:    &OutConfig{},
			want: rune(','),
		},
		{
			name: "2",
			c: &OutConfig{
				Delimiter: "\u0010",
			},
			want: rune(0x0010),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.comma(); got != tt.want {
				t.Errorf("OutConfig.comma() = %v, want %v", got, tt.want)
			}
		})
	}
}
