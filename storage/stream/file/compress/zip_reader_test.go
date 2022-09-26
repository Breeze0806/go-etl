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

package compress

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestZipReader_Read(t *testing.T) {
	type args struct {
		filename string
		p        []byte
	}
	tests := []struct {
		name    string
		args    args
		wantN   int
		wantP   []byte
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				filename: "a.zip",
				p:        make([]byte, 36),
			},
			wantN: 36,
			wantP: []byte("abcdefghijklmnopqrstuvwxyz1234567890"),
		},
		{
			name: "2",
			args: args{
				filename: "a.zip",
				p:        make([]byte, 33),
			},
			wantN: 33,
			wantP: []byte("abcdefghijklmnopqrstuvwxyz1234567"),
		},
		{
			name: "3",
			args: args{
				filename: "a.zip",
				p:        make([]byte, 37),
			},
			wantN:   36,
			wantP:   []byte("abcdefghijklmnopqrstuvwxyz1234567890"),
			wantErr: true,
		},
		{
			name: "4",
			args: args{
				filename: "a.zip",
				p:        make([]byte, 6),
			},
			wantN: 6,
			wantP: []byte("abcdef"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := os.Open(filepath.Join("testdata", tt.args.filename))
			if err != nil {
				t.Errorf("Open fail. err: %v", err)
				return
			}
			defer f.Close()
			z, err := NewZipReader(f)
			if err != nil {
				t.Errorf("NewZipReader fail. err: %v", err)
				return
			}
			gotN, err := z.Read(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZipReader.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotN != tt.wantN {
				t.Errorf("ZipReader.Read() = %v, want %v", gotN, tt.wantN)
				return
			}
			if !reflect.DeepEqual(tt.args.p[:gotN], tt.wantP) {
				t.Errorf("ZipReader.Read() = %v, want %v", tt.args.p, tt.wantP)
			}
		})
	}
}
