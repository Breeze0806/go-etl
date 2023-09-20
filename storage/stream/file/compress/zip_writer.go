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

	"github.com/google/uuid"
)

// ZipWriter zip压缩写入器
type ZipWriter struct {
	writer    *zip.Writer
	nowWriter io.Writer
}

// NewZipWriter 创建zip压缩写入器
func NewZipWriter(f *os.File) (zw *ZipWriter, err error) {
	var fi fs.FileInfo

	if fi, err = f.Stat(); err != nil {
		return nil, err
	}

	zw = &ZipWriter{
		writer: zip.NewWriter(f),
	}

	if zw.nowWriter, err = zw.writer.Create(fi.Name() + uuid.New().String()); err != nil {
		return nil, err
	}
	return
}

// Write 写入p
func (z *ZipWriter) Write(p []byte) (n int, err error) {
	defer z.writer.Flush()
	return z.nowWriter.Write(p)
}

// Close 关闭
func (z *ZipWriter) Close() error {
	return z.writer.Close()
}
