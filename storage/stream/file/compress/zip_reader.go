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
	"archive/zip"
	"io"
	"io/fs"
	"os"
)

// ZipReader - A reader for ZIP files.
type ZipReader struct {
	reader    *zip.Reader
	now       int
	nowReader io.ReadCloser
}

// NewZipReader - Creates a ZIP reader using the file 'f'.
func NewZipReader(f *os.File) (zr *ZipReader, err error) {
	var fi fs.FileInfo
	fi, err = f.Stat()
	if err != nil {
		return nil, err
	}
	zr = &ZipReader{}
	zr.reader, err = zip.NewReader(f, fi.Size())
	if err != nil {
		return nil, err
	}
	return zr, nil
}

// Read - Reads the content passed in 'p'.
func (z *ZipReader) Read(p []byte) (n int, err error) {
	if z.now == 0 {
		if err = z.getNowReader(); err != nil {
			return
		}
	}

	for readBytes := 0; n < len(p); {
		readBytes, err = z.nowReader.Read(p[n:])
		n += readBytes
		if err == io.EOF {
			if err = z.nowReader.Close(); err != nil {
				return
			}
			if err = z.getNowReader(); err != nil {
				return
			}
		}

		if err != nil {
			return
		}
	}
	return
}

func (z *ZipReader) getNowReader() (err error) {
	if z.now >= len(z.reader.File) {
		return io.EOF
	}

	z.nowReader, err = z.reader.File[z.now].Open()
	if err != nil {
		return err
	}
	z.now++
	return
}
