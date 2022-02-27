package csv

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/element"
	"github.com/Breeze0806/go-etl/storage/stream/file"
)

func init() {
	var opener Opener
	file.RegisterOpener("csv", &opener)
}

type Opener struct {
}

func (o *Opener) Open(filename string) (file.Stream, error) {
	return NewStream(filename)
}

type Stream struct {
	file *os.File
}

func NewStream(filename string) (file.Stream, error) {
	stream := &Stream{}
	var err error
	stream.file, err = os.Open(filename)
	if err != nil {
		return nil, err
	}
	return stream, nil
}

func (s *Stream) Rows(conf *config.JSON) (rows file.Rows, err error) {
	return NewRows(s.file, conf)
}

func (s *Stream) Close() (err error) {
	return s.file.Close()
}

type Rows struct {
	columns map[int]Column
	reader  *csv.Reader
	record  []string
	err     error
}

func NewRows(f *os.File, c *config.JSON) (file.Rows, error) {
	var conf *Config
	var err error
	if conf, err = NewConfig(c); err != nil {
		return nil, err
	}
	rows := &Rows{
		columns: make(map[int]Column),
	}
	rows.reader = csv.NewReader(f)
	rows.reader.Comma = []rune(conf.Delimiter)[0]
	for _, v := range conf.Columns {
		rows.columns[v.index()] = v
	}
	return rows, err
}

func (r *Rows) Next() bool {
	if r.record, r.err = r.reader.Read(); r.err != nil {
		if r.err == io.EOF {
			r.err = nil
		}
		return false
	}
	return true
}

func (r *Rows) Scan() (columns []element.Column, err error) {
	for i, v := range r.record {
		var c element.Column
		c, err = r.getColum(i, v)
		if err != nil {
			return nil, err
		}
		columns = append(columns, c)
	}
	return
}

func (r *Rows) Error() error {
	return r.err
}

func (r *Rows) Close() error {
	return nil
}

func (r *Rows) getColum(index int, s string) (element.Column, error) {
	c, ok := r.columns[index]
	if ok && element.ColumnType(c.Type) == element.TypeTime {
		layout := c.layout()
		t, err := time.Parse(layout, s)
		if err != nil {
			return nil, err
		}
		return element.NewDefaultColumn(element.NewTimeColumnValueWithDecoder(t,
			element.NewStringTimeDecoder(layout)),
			strconv.Itoa(index), 0), nil
	}
	return element.NewDefaultColumn(element.NewStringColumnValue(s), strconv.Itoa(index), 0), nil
}
