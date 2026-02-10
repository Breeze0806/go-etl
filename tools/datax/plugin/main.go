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

var log mylog.Logger = mylog.NewDefaultLogger(os.Stdout, mylog.ErrorLevel, "")

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
	readmeFile = `# %v%v Plugin Documentation
`

	readerFile = `package %v
import (
	"github.com/Breeze0806/go-etl/config"
	spireader "github.com/Breeze0806/go-etl/datax/common/spi/reader"
)

//A reader is uesed to extract data from data source 
type Reader struct {
	pluginConf *config.JSON
}

//ResourcesConfig returns the configuration of the data source to initiate the reader.
func (r *Reader) ResourcesConfig() *config.JSON {
	return r.pluginConf
}

//Job returns a description of how the reader extracts data from the data source.
func (r *Reader) Job() spireader.Job {
	// todo like below
	//job := NewJob()
	//job.SetPluginConf(r.pluginConf)
	//return job
}

//Task returns the smallest execution unit obtained by maximizing the split of a Job
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
	

//A writer is uesed to load data into data source 
type Writer struct {
	pluginConf *config.JSON
}

//ResourcesConfig returns the configuration of the data source to initiate the writer.
func (w *Writer) ResourcesConfig() *config.JSON {
	return w.pluginConf
}

//Job returns a description of how the reader extracts data from the data source.
func (w *Writer) Job() spiwriter.Job {
	// todo like below
	//job := NewJob()
	//job.SetPluginConf(w.pluginConf)
	//return job
}

//Task returns the smallest execution unit obtained by maximizing the split of a Job
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

	if err := os.Mkdir(packPath, 0755); err != nil {
		log.Errorf("mkdir %v fail. err: %v", packPath, err)
		return
	}

	resourcePath := filepath.Join(packPath, "resources")
	if err := os.Mkdir(resourcePath, 0755); err != nil {
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
	// To assist in development, a MySQL reader template is automatically generated in datax/plugin/reader
	//     reader---mysql--+-----resources--+--plugin.json
	//                     |--job.go        |--plugin_job_template.json
	//                     |--reader.go
	//                     |--README.md
	//                     |--task.go
	case "reader":
		files = append(files, file{
			filename: filepath.Join(packPath, "reader.go"),
			content:  fmt.Sprintf(readerFile, p),
		})
	// To assist in development,  a MySQL writer template is automatically generated in datax/plugin/writer
	// 	writer--mysql---+-----resources--+--plugin.json
	// 					|--job.go        |--plugin_job_template.json
	// 					|--README.md
	// 					|--task.go
	// 					|--writer.go
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
