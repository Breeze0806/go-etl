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

package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	mylog "github.com/Breeze0806/go/log"
)

//go:generate go run main.go
var log mylog.Logger = mylog.NewDefaultLogger(os.Stdout, mylog.ErrorLevel, "[datax]")

const (
	readerCode = `package %v

import (
	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/datax/plugin/reader"
)

var _pluginConfig string

func init() {
	var err error
	maker := &maker{}
	if _pluginConfig, err = reader.RegisterReader(maker); err != nil {
		panic(err)
	}
}

var pluginConfig = %v

//NewReaderFromFile 创建读取器
func NewReaderFromFile(filename string) (rd reader.Reader, err error) {
	r := &Reader{}
	r.pluginConf, err = config.NewJSONFromFile(filename)
	if err != nil {
		return nil, err
	}
	rd = r
	return
}

//NewReaderFromString 创建读取器
func NewReaderFromString(filename string) (rd reader.Reader, err error) {
	r := &Reader{}
	r.pluginConf, err = config.NewJSONFromString(filename)
	if err != nil {
		return nil, err
	}
	rd = r
	return
}

type maker struct{}

func (m *maker) FromFile(filename string) (reader.Reader, error) {
	return NewReaderFromFile(filename)
}

func (m *maker) Default() (reader.Reader, error) {
	return NewReaderFromString(pluginConfig)
}
`
	writerCode = `package %v

import (
	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/datax/plugin/writer"
)

var _pluginConfig string

func init() {
	var err error
	maker := &maker{}
	if _pluginConfig, err = writer.RegisterWriter(maker); err != nil {
		panic(err)
	}
}

var pluginConfig = %v

//NewWriterFromFile 创建写入器
func NewWriterFromFile(filename string) (wr writer.Writer, err error) {
	w := &Writer{}
	w.pluginConf, err = config.NewJSONFromFile(filename)
	if err != nil {
		return nil, err
	}
	wr = w
	return
}

//NewWriterFromString 创建写入器
func NewWriterFromString(filename string) (wr writer.Writer, err error) {
	w := &Writer{}
	w.pluginConf, err = config.NewJSONFromString(filename)
	if err != nil {
		return nil, err
	}
	wr = w
	return
}

type maker struct{}

func (m *maker) FromFile(filename string) (writer.Writer, error) {
	return NewWriterFromFile(filename)
}

func (m *maker) Default() (writer.Writer, error) {
	return NewWriterFromString(pluginConfig)
}`
	sourcePath = "../../../datax/"
)

func main() {
	var imports []string
	parser := pluginParser{}
	if err := parser.readPackages(sourcePath + "plugin/reader"); err != nil {
		log.Errorf("readPackages %v", err)
		return
	}
	for _, info := range parser.infos {
		if err := info.genFile(sourcePath+"plugin/reader", readerCode); err != nil {
			log.Errorf("genFile %v", err)
			return
		}
		imports = append(imports, info.genImport("reader"))
	}

	imports = append(imports, "")
	parser.infos = nil
	if err := parser.readPackages(sourcePath + "plugin/writer"); err != nil {
		log.Errorf("readPackages %v", err)
		return
	}
	for _, info := range parser.infos {
		if err := info.genFile(sourcePath+"plugin/writer", writerCode); err != nil {
			log.Errorf("genFile %v", err)
			return
		}
		imports = append(imports, info.genImport("writer"))
	}

	if err := writeAllPlugins(imports); err != nil {
		log.Errorf("writeAllPlugins fail. err: %v", err)
		return
	}
}

type pluginParser struct {
	infos []pluginInfo
}

func (p *pluginParser) readPackages(path string) (err error) {
	var list []os.FileInfo
	list, err = ioutil.ReadDir(path)
	if err != nil {
		return err
	}
	var data []byte
	for _, v := range list {
		if v.IsDir() {
			data, err = ioutil.ReadFile(filepath.Join(path, v.Name(), "resources", "plugin.json"))
			if err != nil {
				err = nil
				continue
			}
			info := pluginInfo{
				shotPackage:  v.Name(),
				pluginConfig: "`" + string(data) + "`",
			}
			p.infos = append(p.infos, info)
		}
	}
	return
}

type pluginInfo struct {
	shotPackage  string
	pluginConfig string
}

func (p *pluginInfo) genImport(typ string) string {
	return fmt.Sprintf(`	_ "github.com/Breeze0806/go-etl/datax/plugin/%s/%s"`, typ, p.shotPackage)
}

func (p *pluginInfo) genFile(path string, code string) (err error) {
	var f *os.File
	f, err = os.Create(filepath.Join(path, p.shotPackage, "plugin.go"))
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = fmt.Fprintf(f, code, p.shotPackage, p.pluginConfig)
	return
}

func writeAllPlugins(imports []string) (err error) {
	var f *os.File
	f, err = os.Create(sourcePath + "plugin.go")
	if err != nil {
		return
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	defer w.Flush()
	f.WriteString(`package datax

import (
`)
	for _, v := range imports {
		f.WriteString(v)
		f.WriteString("\n")
	}
	f.WriteString(")\n")
	return
}
