package csv

import (
	"path/filepath"
	"testing"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/element"
)

func Test_ReadWrite(t *testing.T) {
	record := element.NewDefaultRecord()
	record.Add(element.NewDefaultColumn(element.NewStringColumnValueWithEncoder("20220101", element.NewStringTimeEncoder("20060102")),
		"1", 0))
	record.Add(element.NewDefaultColumn(element.NewStringColumnValue("abc"),
		"2", 0))
	type args struct {
		record   element.Record
		conf     *config.JSON
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
				record:   record,
				conf:     testJSONFromString(`{"column":[{"index":"1","type":"time","format":"yyyy-MM-dd"}]}`),
				filename: filepath.Join(t.TempDir(), "a.csv"),
			},
			wantStr: "0=2022-01-01T00:00:00Z 1=abc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wFunc := func() {
				var creater Creater
				out, err := creater.Create(tt.args.filename)
				if err != nil {
					t.Fatal(err)
				}
				defer out.Close()
				w, err := out.Writer(tt.args.conf)
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
				rows, err := in.Rows(tt.args.conf)
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
