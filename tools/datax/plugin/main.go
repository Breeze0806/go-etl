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
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	mylog "github.com/Breeze0806/go/log"
)

var log mylog.Logger = mylog.NewDefaultLogger(os.Stdout, mylog.ErrorLevel, "[datax]")

const (
	sourcePluginPath = "../../../datax/plugin"

	normalFile = `package %v
`

	resourceFile = `{
	"name" : "%v%v",
	"developer":"",
	"description":""
}`
	resourceJobFile = `{
	"name": "%v%v",
	"parameter": {

	}
}`
	readmeFile = `# %v%v
`

	readerFile = `package %v
import (
	"github.com/Breeze0806/go-etl/config"
	spireader "github.com/Breeze0806/go-etl/datax/common/spi/reader"
)

//Reader 读取器
type Reader struct {
	pluginConf *config.JSON
}

//ResourcesConfig 插件资源配置
func (r *Reader) ResourcesConfig() *config.JSON {
	return r.pluginConf
}

//Job 工作
func (r *Reader) Job() spireader.Job {
	// todo like below
	//job := NewJob()
	//job.SetPluginConf(r.pluginConf)
	//return job
}

//Task 任务
func (r *Reader) Task() spireader.Task {
	// todo like below
	//task := fNewTask()
	//task.SetPluginConf(r.pluginConf)
	//return task
}
`
	writerFile = `package %v

import (
	"github.com/Breeze0806/go-etl/config"
	spiwriter "github.com/Breeze0806/go-etl/datax/common/spi/writer"
)
	

//Writer 写入器
type Writer struct {
	pluginConf *config.JSON
}

//ResourcesConfig 插件资源配置
func (w *Writer) ResourcesConfig() *config.JSON {
	return w.pluginConf
}

//Job 工作
func (w *Writer) Job() spiwriter.Job {
	// todo like below
	//job := NewJob()
	//job.SetPluginConf(w.pluginConf)
	//return job
}

//Task 任务
func (w *Writer) Task() spiwriter.Task {
    // todo like below
	//task := NewTask()
	//task.SetPluginConf(w.pluginConf)
	//return task
}
`
)

type file struct {
	filename string
	content  string
}

func (fi *file) create() (err error) {
	var f *os.File
	f, err = os.Create(fi.filename)
	if err != nil {
		return
	}
	defer f.Close()
	_, err = f.WriteString(fi.content)
	return
}

func main() {
	typ := flag.String("t", "", "set plugin type as reader or writer")
	pack := flag.String("p", "", "set plugin package name")
	delete := flag.Bool("d", false, "delete plugin package")
	flag.Parse()
	if *typ != "reader" && *typ != "writer" {
		log.Errorf("-t %v is neither reader nor writer", *typ)
		return
	}

	if *pack == "" {
		log.Errorf("-p %v is empty", *typ)
		return
	}

	p := strings.ToLower(*pack)

	packPath := filepath.Join(sourcePluginPath, *typ, p)
	if *delete {
		os.RemoveAll(packPath)
	}

	if err := os.Mkdir(packPath, 0664); err != nil {
		log.Errorf("mkdir %v fail. err: %v", packPath, err)
		return
	}

	resourcePath := filepath.Join(packPath, "resources")
	if err := os.Mkdir(resourcePath, 0664); err != nil {
		log.Errorf("mkdir %v fail. err: %v", packPath, err)
		return
	}

	var files []file

	files = append(files, file{
		filename: filepath.Join(packPath, "task.go"),
		content:  fmt.Sprintf(normalFile, p),
	})

	files = append(files, file{
		filename: filepath.Join(packPath, "job.go"),
		content:  fmt.Sprintf(normalFile, p),
	})

	files = append(files, file{
		filename: filepath.Join(packPath, "README.md"),
		content:  fmt.Sprintf(readmeFile, *pack, strings.Title(*typ)),
	})

	switch *typ {
	case "reader":
		files = append(files, file{
			filename: filepath.Join(packPath, "reader.go"),
			content:  fmt.Sprintf(readerFile, p),
		})
	case "writer":
		files = append(files, file{
			filename: filepath.Join(packPath, "writer.go"),
			content:  fmt.Sprintf(writerFile, p),
		})
	}

	files = append(files, file{
		filename: filepath.Join(resourcePath, "plugin.json"),
		content:  fmt.Sprintf(resourceFile, p, *typ),
	})

	files = append(files, file{
		filename: filepath.Join(resourcePath, "plugin_job_template.json"),
		content:  fmt.Sprintf(resourceJobFile, p, *typ),
	})
	for _, v := range files {
		if err := v.create(); err != nil {
			log.Errorf("create %+v fail. err: %v", v, err)
			return
		}
	}
}
