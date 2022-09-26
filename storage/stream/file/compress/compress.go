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

//Type 压缩类型
type Type string

//压缩类型枚举
const (
	TypeNone    Type = ""
	TypeTarGzip Type = "targz"
	TypeTar     Type = "tar"
	TypeZip     Type = "zip"
	TypeGzip    Type = "gz"
)

//ReadCloser 获取读取关闭器
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

//WriteCloser 获取写入关闭器
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

//ReadCloser 读取关闭器
type ReadCloser struct {
	io.Reader
}

//Read 读取p
func (r *ReadCloser) Read(p []byte) (n int, err error) {
	return r.Reader.Read(p)
}

//Close 关闭
func (r *ReadCloser) Close() error {
	return nil
}

//NewNoneReadCloser 获取无压缩读取关闭器
func NewNoneReadCloser(f *os.File) *ReadCloser {
	return &ReadCloser{
		Reader: f,
	}
}

//NewZipReadCloser 获取zip压缩读取关闭器
func NewZipReadCloser(f *os.File) (r *ReadCloser, err error) {
	r = &ReadCloser{}
	if r.Reader, err = NewZipReader(f); err != nil {
		return nil, err
	}
	return
}

//NewGzipReadCloser 获取gzip压缩读取关闭器
func NewGzipReadCloser(f *os.File) (r *ReadCloser, err error) {
	r = &ReadCloser{}

	if r.Reader, err = gzip.NewReader(f); err != nil {
		return nil, err
	}
	return
}

//NoneWriter 无压缩写入器
type NoneWriter struct {
	file *os.File
}

//NewNoneWriter 创建无压缩写入器
func NewNoneWriter(f *os.File) (nw *NoneWriter) {
	return &NoneWriter{
		file: f,
	}
}

//Write 写入p
func (nw *NoneWriter) Write(p []byte) (n int, err error) {
	return nw.file.Write(p)
}

//Close 关闭
func (nw *NoneWriter) Close() error {
	return nil
}

//GzipWriter Gzip压缩写入器
type GzipWriter struct {
	writer *gzip.Writer
}

//NewGzipWriter 创建gzip压缩写入器
func NewGzipWriter(f *os.File) (gw *GzipWriter) {
	return &GzipWriter{
		writer: gzip.NewWriter(f),
	}
}

//Write 写入p
func (g *GzipWriter) Write(p []byte) (n int, err error) {
	defer g.writer.Flush()
	return g.writer.Write(p)
}

//Close 关闭
func (g *GzipWriter) Close() error {
	return g.writer.Close()
}
