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
	"encoding/csv"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	mylog "github.com/Breeze0806/go/log"
)

var log mylog.Logger = mylog.NewDefaultLogger(os.Stdout, mylog.ErrorLevel, "")

type translate struct {
	line int
	en   string
	chn  string
}

func main() {
	translate := flag.Bool("t", false, "translate")
	flag.Parse()

	packages := []string{"schedule"}
	var codeFiles []string //= []string{
	// 	"storage/database/mysql/field.go",
	// 	"storage/database/oracle/field.go",
	// 	"storage/database/mysql/table.go",
	// 	"storage/database/db2/field.go",
	// 	"storage/database/oracle/table.go",
	// 	"storage/database/postgres/table.go",
	// 	"storage/database/sqlserver/field.go",
	// 	"storage/database/sqlserver/table.go",
	// }
	for _, v := range packages {
		if err := filepath.Walk(v, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				filename := info.Name()
				if filepath.Ext(filename) == ".go" {
					if !strings.HasSuffix(filename, "_test.go") {
						codeFiles = append(codeFiles, path)
					}
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

			if *translate {
				log.Infof("translateComment %v", filename)
				if err := translateComment(filename); err != nil {
					log.Errorf("translateComment %v fail. err : %v", filename, err)
					return
				}

			} else {
				log.Infof("fetchComment %v", filename)
				if err := fetchComment(filename); err != nil {
					log.Errorf("fetchComment %v fail. err : %v", filename, err)
					return
				}
			}
		}(v)
	}
	wg.Wait()
}

// func readPackages(path string) (packages []string, err error) {
// 	var list []os.FileInfo
// 	list, err = ioutil.ReadDir(path)
// 	if err != nil {
// 		return
// 	}

// 	for _, v := range list {
// 		if v.IsDir() {
// 			switch v.Name() {
// 			case "vendor", ".vscode", ".git":
// 			default:
// 				packages = append(packages, v.Name())
// 			}
// 		}
// 	}
// 	return
// }

func fetchComment(filename string) (err error) {
	fset := token.NewFileSet()
	var astFile *ast.File
	astFile, err = parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return
	}

	chnFile := filename + ".chn"
	var f, fen *os.File
	f, err = os.Create(chnFile)
	if err != nil {
		return
	}
	defer f.Close()

	//if _, err = os.Stat(filename + ".en"); err != nil {
	fen, err = os.Create(filename + ".en")
	if err != nil {
		log.Errorf("create %v.en fail. err : %v", filename, err)
		return
	}
	defer fen.Close()
	//}

	var w, wen *csv.Writer
	w = csv.NewWriter(f)
	w.Comma = '^'
	defer w.Flush()
	w.Write([]string{
		"请将中文翻译成英文", "",
	})
	l := 1

	if fen != nil {
		wen = csv.NewWriter(fen)
		wen.Comma = '^'
		defer wen.Flush()
		wen.Write([]string{
			"请将中文翻译成英文", "",
		})
	}
	for i, commentGroup := range astFile.Comments {
		if i == 0 {
			continue
		}
		for _, comment := range commentGroup.List {
			pos := comment.Pos()
			line, _ := fset.Position(pos).Line, fset.Position(pos).Column
			l++
			record := []string{
				strconv.Itoa(line),
				comment.Text,
			}
			w.Write(record)
			if fen != nil {
				for i := range record {
					record[i] = strings.ReplaceAll(record[i], ":", " -")
				}
				wen.Write(record)
			}
			if l%1000 == 0 {
				if fen != nil {
					wen.Flush()
				}
				w.Flush()
			}
		}
	}
	return
}

func translateComment(filename string) (err error) {

	var tm map[int]*translate
	if tm, err = mapComment(filename); err != nil {
		return fmt.Errorf("mapComment %w", err)
	}
	return replaceComment(filename, tm)
}

func mapComment(filename string) (tm map[int]*translate, err error) {
	chnFileName := filename + ".chn"
	enFilename := filename + ".en"
	var chnFile, enFile *os.File
	if chnFile, err = os.Open(chnFileName); err != nil {
		return
	}
	defer chnFile.Close()

	if enFile, err = os.Open(enFilename); err != nil {
		return
	}
	defer enFile.Close()

	chnReader := csv.NewReader(chnFile)
	chnReader.Comma = '^'
	tm = make(map[int]*translate)

	if _, err = chnReader.Read(); err != nil {
		return
	}
	var record []string

	for {
		if record, err = chnReader.Read(); err != nil {
			if err == io.EOF {
				err = nil
				break
			}
			err = fmt.Errorf("chn read %w", err)
			return
		}

		line, _ := strconv.Atoi(record[0])

		tm[line] = &translate{
			line: line,
			chn:  strings.TrimSpace(record[1]),
		}
	}
	enReader := csv.NewReader(enFile)
	enReader.Comma = '^'
	for {
		enReader.FieldsPerRecord = 0
		if record, err = enReader.Read(); err != nil {
			if err == io.EOF {
				err = nil
				return
			}
			err = fmt.Errorf("en read %w", err)

			return
		}

		if len(record) == 0 {
			continue
		}
		line, _ := strconv.Atoi(record[0])
		if t, ok := tm[line]; ok {
			t.en = strings.TrimSpace(record[1])
		}
	}
}

func replaceComment(filename string, tm map[int]*translate) (err error) {

	rlFilename := filename + ".rl"
	defer os.Rename(rlFilename, filename)

	var f, rf *os.File
	if f, err = os.Open(filename); err != nil {
		return
	}
	defer f.Close()

	if rf, err = os.Create(rlFilename); err != nil {
		return
	}
	defer rf.Close()

	r := bufio.NewReaderSize(f, 1024*32)
	w := bufio.NewWriter(rf)
	defer w.Flush()
	var line []byte
	l := 0
	for {
		if line, _, err = r.ReadLine(); err != nil {
			if err == io.EOF {
				err = nil
				return
			}
			return
		}
		l++
		s := string(line)

		if t, ok := tm[l]; ok {
			s = strings.ReplaceAll(s, t.chn, t.en)
		}
		w.WriteString(s)
		w.WriteString("\r\n")
		if l%1000 == 0 {
			w.Flush()
		}
	}

}
