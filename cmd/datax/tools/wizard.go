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
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Breeze0806/go-etl/config"
	coreconst "github.com/Breeze0806/go-etl/datax/common/config/core"
)

//Wizard 配置向导工具
type Wizard struct {
	dataSourceFile string
	csvFile        string
}

//NewWizard 根据数据源文件dataSourceFile，源目的文件csvFile生成配置向导工具
func NewWizard(dataSourceFile, csvFile string) (w *Wizard) {
	w = &Wizard{
		dataSourceFile: dataSourceFile,
		csvFile:        csvFile,
	}
	return
}

//GenerateConfigs 生成配置文件集
func (w *Wizard) GenerateConfigs() (err error) {
	var dataSource *config.JSON
	if dataSource, err = config.NewJSONFromFile(w.dataSourceFile); err != nil {
		return err
	}

	dataSourceAbsFile := ""
	if dataSourceAbsFile, err = filepath.Abs(w.dataSourceFile); err != nil {
		return err
	}

	dataSourceExt := filepath.Ext(dataSourceAbsFile)
	os.MkdirAll(filepath.Join(filepath.Dir(dataSourceAbsFile), "config"), 0644)
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
	for {
		line++
		if record, err = r.Read(); err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}
		if len(record) != 2 {
			return fmt.Errorf("源目的文件的第%d行不是两列", line)
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
			"postgresreader", "sqlserverreader":
			if err = cloneDataSource.Set(coreconst.DataxJobContentReaderParameter+
				".connection.table.name", record[0]); err != nil {
				return err
			}
		default:
			return fmt.Errorf("reader name(%v) does not support", readerName)
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
			"postgreswriter", "sqlserverwriter":
			if err = cloneDataSource.Set(coreconst.DataxJobContentWriterParameter+
				".connection.table.name", record[1]); err != nil {
				return err
			}
		default:
			return fmt.Errorf("writer name(%v) does not support", writerName)
		}

		err = ioutil.WriteFile(dataSourcePrefix+"_"+strconv.Itoa(line)+dataSourceExt,
			[]byte(cloneDataSource.String()), 0644)
		if err != nil {
			return err
		}
	}
}
