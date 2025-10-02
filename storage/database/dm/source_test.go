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

package dm

import (
	"testing"

	"github.com/Breeze0806/go-etl/storage/database"
)

func TestDialect_Source(t *testing.T) {
	type args struct {
		bs *database.BaseSource
	}
	tests := []struct {
		name    string
		d       Dialect
		args    args
		wantS   database.Source
		wantErr bool
	}{
		{
			name: "1",
			d:    Dialect{},
			args: args{
				bs: database.NewBaseSource(testJSON()),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.d.Source(tt.args.bs)
			if (err != nil) != tt.wantErr {
				t.Errorf("Dialect.Source() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestSource_DriverName(t *testing.T) {
	s, err := NewSource(database.NewBaseSource(testJSON()))
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name string
		s    *Source
		want string
	}{
		{
			name: "1",
			s:    s.(*Source),
			want: "dm",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.DriverName(); got != tt.want {
				t.Errorf("Source.DriverName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSource_ConnectName(t *testing.T) {
	s, err := NewSource(database.NewBaseSource(testJSON()))
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name string
		s    *Source
		want string
	}{
		{
			name: "1",
			s:    s.(*Source),
			want: "dm://username:password@ip:port",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.ConnectName(); got != tt.want {
				t.Errorf("Source.ConnectName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSource_Key(t *testing.T) {
	s, err := NewSource(database.NewBaseSource(testJSON()))
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name string
		s    *Source
		want string
	}{
		{
			name: "1",
			s:    s.(*Source),
			want: "dm://username:password@ip:port",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Key(); got != tt.want {
				t.Errorf("Source.Key() = %v, want %v", got, tt.want)
			}
		})
	}
}
