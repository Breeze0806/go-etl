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

package tools

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Breeze0806/go-etl/config"
	coreconst "github.com/Breeze0806/go-etl/datax/common/config/core"
)

// Wizard is a configuration wizard tool
type Wizard struct {
	dataSourceFile string
	csvFile        string
}

// NewWizard generates a configuration wizard tool based on the data source file dataSourceFile and the source-destination file csvFile
func NewWizard(dataSourceFile, csvFile string) (w *Wizard) {
	w = &Wizard{
		dataSourceFile: dataSourceFile,
		csvFile:        csvFile,
	}
	return
}

// GenerateConfigsAndScripts generates a set of configuration files and execution scripts
func (w *Wizard) GenerateConfigsAndScripts() (err error) {
	var dataSource *config.JSON
	if dataSource, err = config.NewJSONFromFile(w.dataSourceFile); err != nil {
		return err
	}

	dataSourceAbsFile := ""
	if dataSourceAbsFile, err = filepath.Abs(w.dataSourceFile); err != nil {
		return err
	}

	dataSourceExt := filepath.Ext(dataSourceAbsFile)
	if err = os.MkdirAll(filepath.Join(filepath.Dir(dataSourceAbsFile), "config"), 0o755); err != nil {
		return err
	}
	dataSourcePrefix := filepath.Join(filepath.Dir(dataSourceAbsFile), "config",
		filepath.Base(dataSourceAbsFile)[:len(filepath.Base(dataSourceAbsFile))-len(dataSourceExt)])

	var f *os.File
	if f, err = os.Open(w.csvFile); err != nil {
		return err
	}
	defer f.Close()

	r := csv.NewReader(f)

	line := 0
	var record []string
	var scripts []string
	for {
		line++
		if record, err = r.Read(); err == io.EOF {
			break
		}

		if err != nil {
			return err
		}
		if len(record) != 2 {
			return fmt.Errorf("the source-destination file is not two columns on line %d", line)
		}

		var readerName, writerName string
		readerName, err = dataSource.GetString(coreconst.DataxJobContentReaderName)
		if err != nil {
			return err
		}

		writerName, err = dataSource.GetString(coreconst.DataxJobContentWriterName)
		if err != nil {
			return err
		}

		cloneDataSource := dataSource.CloneConfig()

		switch readerName {
		case "xlsxreader":
			if err = cloneDataSource.Set(coreconst.DataxJobContentReaderParameter+
				".xlsxs.0.path", record[0]); err != nil {
				return err
			}
		case "csvreader":
			if err = cloneDataSource.Set(coreconst.DataxJobContentReaderParameter+
				".path.0", record[0]); err != nil {
				return err
			}
		case "db2reader", "mysqlreader", "oraclereader",
			"postgresreader", "sqlserverreader", "sqlite3reader":
			if err = cloneDataSource.Set(coreconst.DataxJobContentReaderParameter+
				".connection.table.name", record[0]); err != nil {
				return err
			}
		default:
			return fmt.Errorf("reader name(%v) is not supported", readerName)
		}

		switch writerName {
		case "xlsxwriter":
			if err = cloneDataSource.Set(coreconst.DataxJobContentWriterParameter+
				".xlsxs.0.path", record[0]); err != nil {
				return err
			}
		case "csvwriter":
			if err = cloneDataSource.Set(coreconst.DataxJobContentWriterParameter+
				".path.0", record[1]); err != nil {
				return err
			}
		case "db2writer", "mysqlwriter", "oraclewriter",
			"postgreswriter", "sqlserverwriter", "sqlite3writer":
			if err = cloneDataSource.Set(coreconst.DataxJobContentWriterParameter+
				".connection.table.name", record[1]); err != nil {
				return err
			}
		default:
			return fmt.Errorf("writer name(%v) is not supported", writerName)
		}

		filename := dataSourcePrefix + "_" + strconv.Itoa(line) + dataSourceExt
		err = os.WriteFile(filename, []byte(cloneDataSource.String()), fs.FileMode(0o644))
		if err != nil {
			return err
		}
		scripts = append(scripts, generateScript(filename))
	}

	err = os.WriteFile("run"+ext(), []byte(strings.Join(scripts, "\n")), fs.FileMode(0o644))
	if err != nil {
		return err
	}
	return
}
