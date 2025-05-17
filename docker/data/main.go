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
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Breeze0806/go-etl/element"
)

func main() {
	write(`split.csv`, 0)
	write(`split1.csv`, 10000000)
}

func write(filename string, start int) {
	f, err := os.Create(filename)
	if err != nil {
		fmt.Println("crete file fail. err:", err)
		return
	}
	defer f.Close()

	w := csv.NewWriter(f)
	for i := start; i < start+10000000; i++ {
		record := []string{strconv.Itoa(i),
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, i/1000).Format(element.DefaultTimeFormat[:10]),
			base64.StdEncoding.EncodeToString([]byte{byte(i / 100 / 100), byte((i / 100) % 100), byte(i % 100)}),
		}
		w.Write(record)
		if (i+1)%1000 == 0 {
			w.Flush()
		}
	}
	w.Flush()
}
