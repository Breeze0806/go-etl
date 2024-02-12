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
	"compress/gzip"
	"fmt"
	"io"
	"os"
)

// Type represents the compression type
type Type string

// CompressionTypeEnum enumerates the compression types
const (
	TypeNone    Type = ""
	TypeTarGzip Type = "targz"
	TypeTar     Type = "tar"
	TypeZip     Type = "zip"
	TypeGzip    Type = "gz"
)

// ReadCloser retrieves a read closer
func (c Type) ReadCloser(f *os.File) (r io.ReadCloser, err error) {
	switch c {
	case TypeNone:
		r = NewNoneReadCloser(f)
		return
	case TypeZip:
		r, err = NewZipReadCloser(f)
		return
	case TypeGzip:
		r, err = NewGzipReadCloser(f)
		return
	}
	err = fmt.Errorf("unsupported type %v", c)
	return
}

// WriteCloser retrieves a write closer
func (c Type) WriteCloser(f *os.File) (w io.WriteCloser, err error) {
	switch c {
	case TypeNone:
		w = NewNoneWriter(f)
		return
	case TypeZip:
		w, err = NewZipWriter(f)
		return
	case TypeGzip:
		w = NewGzipWriter(f)
		return
	}
	err = fmt.Errorf("unsupported type %v", c)
	return
}

// ReadCloser is a read closer
type ReadCloser struct {
	io.Reader
}

// Read reads 'p'
func (r *ReadCloser) Read(p []byte) (n int, err error) {
	return r.Reader.Read(p)
}

// Close closes the connection
func (r *ReadCloser) Close() error {
	return nil
}

// NewNoneReadCloser retrieves a non-compression read closer
func NewNoneReadCloser(f *os.File) *ReadCloser {
	return &ReadCloser{
		Reader: f,
	}
}

// NewZipReadCloser retrieves a zip compression read closer
func NewZipReadCloser(f *os.File) (r *ReadCloser, err error) {
	r = &ReadCloser{}
	if r.Reader, err = NewZipReader(f); err != nil {
		return nil, err
	}
	return
}

// NewGzipReadCloser retrieves a gzip compression read closer
func NewGzipReadCloser(f *os.File) (r *ReadCloser, err error) {
	r = &ReadCloser{}

	if r.Reader, err = gzip.NewReader(f); err != nil {
		return nil, err
	}
	return
}

// NoneWriter is a non-compression writer
type NoneWriter struct {
	file *os.File
}

// NewNoneWriter creates a non-compression writer
func NewNoneWriter(f *os.File) (nw *NoneWriter) {
	return &NoneWriter{
		file: f,
	}
}

// Write writes 'p'
func (nw *NoneWriter) Write(p []byte) (n int, err error) {
	return nw.file.Write(p)
}

// Close closes the writer
func (nw *NoneWriter) Close() error {
	return nil
}

// GzipWriter is a gzip compression writer
type GzipWriter struct {
	writer *gzip.Writer
}

// NewGzipWriter creates a gzip compression writer
func NewGzipWriter(f *os.File) (gw *GzipWriter) {
	return &GzipWriter{
		writer: gzip.NewWriter(f),
	}
}

// Write writes 'p'
func (g *GzipWriter) Write(p []byte) (n int, err error) {
	defer g.writer.Flush()
	return g.writer.Write(p)
}

// Close closes the writer
func (g *GzipWriter) Close() error {
	return g.writer.Close()
}
