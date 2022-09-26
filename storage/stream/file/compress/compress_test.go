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
	"io"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestType_ReadWrite(t *testing.T) {
	type args struct {
		filename string
		p        []byte
	}
	tests := []struct {
		name  string
		c     Type
		args  args
		wantP []byte
	}{
		{
			name: "1",
			c:    TypeNone,
			args: args{
				filename: "a",
				p:        make([]byte, 36),
			},
			wantP: []byte("abcdefghijklmnopqrstuvwxyz1234567890"),
		},
		{
			name: "2",
			c:    TypeZip,
			args: args{
				filename: "a.zip",
				p:        make([]byte, 36),
			},
			wantP: []byte("abcdefghijklmnopqrstuvwxyz1234567890"),
		},
		{
			name: "3s",
			c:    TypeGzip,
			args: args{
				filename: "a.gz",
				p:        make([]byte, 36),
			},
			wantP: []byte("abcdefghijklmnopqrstuvwxyz1234567890"),
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
				w, err := tt.c.WriteCloser(f)
				defer w.Close()
				if err != nil {
					t.Errorf("WriteCloser fail. err: %v", err)
					return
				}
				w.Write(tt.wantP)
			}

			read := func() {
				f, err := os.Open(filename)
				if err != nil {
					t.Errorf("Open fail. err: %v", err)
					return
				}
				defer f.Close()
				r, err := tt.c.ReadCloser(f)
				defer r.Close()
				if err != nil {
					t.Errorf("ReadCloser fail. err: %v", err)
					return
				}
				gotN, _ := r.Read(tt.args.p)
				if !reflect.DeepEqual(tt.args.p[:gotN], tt.wantP) {
					t.Errorf("Read() = %v, want %v", tt.args.p, tt.wantP)
				}
			}
			write()
			read()
		})
	}
}

func TestType_ReadCloser(t *testing.T) {
	type args struct {
		filename string
		p        []byte
	}
	tests := []struct {
		name    string
		c       Type
		args    args
		wantP   []byte
		wantErr bool
	}{
		{
			name: "1",
			c:    Type("7z"),
			args: args{
				filename: "a.tar",
				p:        make([]byte, 36),
			},
			wantErr: true,
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
			r, err := tt.c.ReadCloser(f)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadCloser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			defer r.Close()
			gotN, _ := r.Read(tt.args.p)
			if !reflect.DeepEqual(tt.args.p[:gotN], tt.wantP) {
				t.Errorf("Read() = %v, want %v", tt.args.p, tt.wantP)
			}
		})
	}
}

func TestType_WriteCloser(t *testing.T) {
	type args struct {
		f *os.File
	}
	tests := []struct {
		name    string
		c       Type
		args    args
		wantW   io.WriteCloser
		wantErr bool
	}{
		{
			name:    "1",
			c:       Type("7z"),
			args:    args{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotW, err := tt.c.WriteCloser(tt.args.f)
			if (err != nil) != tt.wantErr {
				t.Errorf("Type.WriteCloser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotW, tt.wantW) {
				t.Errorf("Type.WriteCloser() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
