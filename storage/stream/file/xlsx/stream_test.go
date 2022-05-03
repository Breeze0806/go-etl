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

package xlsx

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/element"
)

func Test_WriteRead(t *testing.T) {
	type args struct {
		columns  []element.Column
		inConf   *config.JSON
		outConf  *config.JSON
		filename string
	}
	tests := []struct {
		name    string
		args    args
		wantStr string
	}{
		{
			name: "1",
			args: args{
				columns: []element.Column{
					element.NewDefaultColumn(element.NewStringColumnValueWithEncoder(
						"20220101", element.NewStringTimeEncoder("20060102")), "1", 0),
					element.NewDefaultColumn(element.NewStringColumnValue("abc"),
						"2", 0),
				},
				inConf:   testJSONFromString(`{"sheet":"where","column":[{"index":"A","type":"time","format":"yyyy-MM-dd"}]}`),
				outConf:  testJSONFromString(`{"sheets":["where"],"column":[{"index":"A","type":"time","format":"yyyy-MM-dd"}]}`),
				filename: filepath.Join(t.TempDir(), "a.xlsx"),
			},
			wantStr: "0=2022-01-01 00:00:00Z 1=abc",
		},
		{
			name: "2",
			args: args{
				columns: []element.Column{
					element.NewDefaultColumn(element.NewStringColumnValueWithEncoder(
						"20220101", element.NewStringTimeEncoder("20060102")), "1", 0),
					element.NewDefaultColumn(element.NewNilStringColumnValue(),
						"2", 0),
				},
				inConf:   testJSONFromString(`{"sheet":"where","nullFormat":"\\N","column":[{"index":"A","type":"time","format":"yyyy-MM-dd"}]}`),
				outConf:  testJSONFromString(`{"sheets":["where"],"nullFormat":"\\N","column":[{"index":"A","type":"time","format":"yyyy-MM-dd"}]}`),
				filename: filepath.Join(t.TempDir(), "a.xlsx"),
			},
			wantStr: "0=2022-01-01 00:00:00Z 1=<nil>",
		},
		{
			name: "3",
			args: args{
				columns: []element.Column{
					element.NewDefaultColumn(element.NewNilStringColumnValue(), "1", 0),
					element.NewDefaultColumn(element.NewStringColumnValue("abc"),
						"2", 0),
				},
				inConf:   testJSONFromString(`{"sheet":"where","nullFormat":"\\N","column":[{"index":"A","type":"time","format":"yyyy-MM-dd"}]}`),
				outConf:  testJSONFromString(`{"sheets":["where"],"nullFormat":"\\N","column":[{"index":"A","type":"time","format":"yyyy-MM-dd"}]}`),
				filename: filepath.Join(t.TempDir(), "a.xlsx"),
			},
			wantStr: "0=<nil> 1=abc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer os.Remove(tt.args.filename)
			record := element.NewDefaultRecord()
			for _, c := range tt.args.columns {
				record.Add(c)
			}
			wFunc := func() {
				var creator Creator
				out, err := creator.Create(tt.args.filename)
				if err != nil {
					t.Fatal(err)
				}
				defer out.Close()
				w, err := out.Writer(tt.args.outConf)
				if err != nil {
					t.Fatal(err)
				}
				defer w.Close()
				defer w.Flush()
				if err = w.Write(record); err != nil {
					t.Fatal(err)
				}
			}

			var got []element.Record
			rFunc := func() {
				var opener Opener
				in, err := opener.Open(tt.args.filename)
				if err != nil {
					t.Fatal(err)
				}
				defer in.Close()
				rows, err := in.Rows(tt.args.inConf)
				if err != nil {
					t.Fatal(err)
				}
				defer rows.Close()
				for rows.Next() {
					r := element.NewDefaultRecord()
					cols, err := rows.Scan()
					if err != nil {
						t.Fatal(err)
					}
					for _, v := range cols {
						r.Add(v)
					}
					got = append(got, r)
				}
				if err = rows.Error(); err != nil {
					t.Fatal(err)
				}
			}
			wFunc()
			rFunc()
			if len(got) != 1 {
				t.Fatal("len is not 1")
			}
			if got[0].String() != tt.wantStr {
				t.Fatalf("got: %v want: %v", got[0].String(), tt.wantStr)
			}
		})
	}
}

func TestWriter_Write(t *testing.T) {
	rotateLine = 1
	type args struct {
		records  []element.Record
		outConf  *config.JSON
		filename string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "1",
			args: args{
				records: []element.Record{
					element.NewDefaultRecord(),
					element.NewDefaultRecord(),
				},
				outConf:  testJSONFromString(`{"sheets":["where","where1"]}`),
				filename: filepath.Join(t.TempDir(), "a.xlsx"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer os.Remove(tt.args.filename)

			var creator Creator
			out, err := creator.Create(tt.args.filename)
			if err != nil {
				t.Fatal(err)
			}
			defer out.Close()
			w, err := out.Writer(tt.args.outConf)
			if err != nil {
				t.Fatal(err)
			}
			defer w.Close()
			defer w.Flush()
			for _, r := range tt.args.records {
				if err = w.Write(r); err != nil {
					t.Fatal(err)
				}
			}

		})
	}
}
