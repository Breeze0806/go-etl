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
	"bytes"
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	mylog "github.com/Breeze0806/go/log"
)

var log mylog.Logger = mylog.NewDefaultLogger(os.Stdout, mylog.ErrorLevel, "[datax]")
var licenseHeader = `// Copyright 2020 the go-etl Authors.
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

`

func main() {
	check := flag.Bool("c", false, "check licenseHeader")
	flag.Parse()
	packages, err := readPackages("./")
	if err != nil {
		log.Errorf("readPackages fail. err : %v", err)
		return
	}
	log.Infof("packages: %v", packages)
	var codeFiles []string
	for _, v := range packages {
		if err := filepath.Walk(v, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				filename := info.Name()
				if filepath.Ext(filename) == ".go" {
					codeFiles = append(codeFiles, path)
				}
			}
			return nil
		}); err != nil {
			log.Errorf("Walk %v fail. err : %v", v, err)
			return
		}
	}
	log.Infof("codeFiles: %v", codeFiles)
	c := make(chan struct{}, 20)
	var wg sync.WaitGroup
	for _, v := range codeFiles {
		c <- struct{}{}
		wg.Add(1)
		go func(filename string) {
			defer func() {
				<-c
				wg.Done()
			}()

			if *check {
				log.Infof("checkLicenseHeader %v", filename)
				if err = checkLicenseHeader(filename); err != nil {
					log.Errorf("checkLicenseHeader %v fail. err : %v", filename, err)
					os.Exit(1)
				}
			} else {
				log.Infof("addLicenseHeader %v", filename)
				if err = addLicenseHeader(filename); err != nil {
					log.Errorf("addLicenseHeader %v fail. err : %v", filename, err)
				}
			}
		}(v)
	}
	wg.Wait()
}

//添加许可证
func readPackages(path string) (packages []string, err error) {
	var list []os.FileInfo
	list, err = ioutil.ReadDir(path)
	if err != nil {
		return
	}

	for _, v := range list {
		if v.IsDir() {
			switch v.Name() {
			case "vendor", ".vscode", ".git":
			default:
				packages = append(packages, v.Name())
			}
		}
	}
	return
}

//检查许可证
func addLicenseHeader(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	data = bytes.ReplaceAll(data, []byte("\r"), []byte(""))
	if bytes.HasPrefix(data, bytes.ReplaceAll([]byte(licenseHeader), []byte("\r"), []byte(""))) {
		return nil
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(licenseHeader)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func checkLicenseHeader(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	data = bytes.ReplaceAll(data, []byte("\r"), []byte(""))
	if bytes.HasPrefix(data, bytes.ReplaceAll([]byte(licenseHeader), []byte("\r"), []byte(""))) {
		return nil
	}
	return fmt.Errorf("has no license header")
}
