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

func TestZipWriter_WriteRead(t *testing.T) {
	type args struct {
		filename string
		p        []byte
	}
	tests := []struct {
		name  string
		args  args
		wantN int
		wantP []byte
		wErr  bool
		rErr  bool
	}{
		{
			name: "1",
			args: args{
				filename: "test.zip",
				p:        make([]byte, 36),
			},
			wantP: []byte("abcdefghijklmnopqrstuvwxyz1234567890"),
			wantN: 36,
			rErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filename := filepath.Join(os.TempDir(), tt.args.filename)
			os.Remove(filename)
			defer os.Remove(filename)
			write := func() {
				f, err := os.Create(filename)
				if err != nil {
					t.Errorf("Open fail. err: %v", err)
					return
				}
				defer f.Close()
				z, err := NewZipWriter(f)
				defer z.Close()
				if err != nil {
					t.Errorf("Open fail. err: %v", err)
					return
				}
				gotN, err := z.Write(tt.wantP)
				if (err != nil) != tt.wErr {
					t.Errorf("ZipWriter.Write() error = %v, wantErr %v", err, tt.wErr)
					return
				}
				if gotN != tt.wantN {
					t.Errorf("ZipWriter.Write() = %v, want %v", gotN, tt.wantN)
				}
			}

			read := func() {
				f, err := os.Open(filename)
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
				if (err != nil) != tt.rErr {
					t.Errorf("ZipReader.Read() error = %v, wantErr %v", err, tt.rErr)
					return
				}
				if gotN != tt.wantN {
					t.Errorf("ZipReader.Read() = %v, want %v", gotN, tt.wantN)
					return
				}
				if !reflect.DeepEqual(tt.args.p[:gotN], tt.wantP) {
					t.Errorf("ZipReader.Read() = %v, want %v", tt.args.p, tt.wantP)
				}
			}
			write()
			read()
		})
	}
}
