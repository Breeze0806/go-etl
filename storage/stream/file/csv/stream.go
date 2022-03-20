package csv

import (
	"encoding/csv"
	"fmt"
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
	var creater Creater
	file.RegisterCreater("csv", &creater)
}

type Opener struct {
}

func (o *Opener) Open(filename string) (file.InStream, error) {
	return NewInStream(filename)
}

type Creater struct {
}

func (c *Creater) Create(filename string) (file.OutStream, error) {
	return NewOutStream(filename)
}

type Stream struct {
	file *os.File
}

func NewInStream(filename string) (file.InStream, error) {
	stream := &Stream{}
	var err error
	stream.file, err = os.Open(filename)
	if err != nil {
		return nil, err
	}
	return stream, nil
}

func NewOutStream(filename string) (file.OutStream, error) {
	stream := &Stream{}
	var err error
	stream.file, err = os.Create(filename)
	if err != nil {
		return nil, err
	}
	return stream, nil
}

func (s *Stream) Writer(conf *config.JSON) (file.StreamWriter, error) {
	return NewWriter(s.file, conf)
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
	return rows, nil
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
			return nil, fmt.Errorf("layout: %v error: %v", layout, err)
		}
		return element.NewDefaultColumn(element.NewTimeColumnValueWithDecoder(t,
			element.NewStringTimeDecoder(layout)),
			strconv.Itoa(index), 0), nil
	}
	return element.NewDefaultColumn(element.NewStringColumnValue(s), strconv.Itoa(index), 0), nil
}

type Writer struct {
	writer  *csv.Writer
	columns map[int]Column
}

func NewWriter(f *os.File, c *config.JSON) (file.StreamWriter, error) {
	var conf *Config
	var err error
	if conf, err = NewConfig(c); err != nil {
		return nil, err
	}
	w := &Writer{
		writer:  csv.NewWriter(f),
		columns: make(map[int]Column),
	}
	w.writer.Comma = []rune(conf.Delimiter)[0]
	for _, v := range conf.Columns {
		w.columns[v.index()] = v
	}
	return w, nil
}

func (w *Writer) Flush() (err error) {
	w.writer.Flush()
	return
}

func (w *Writer) Close() (err error) {
	w.writer.Flush()
	return
}

func (w *Writer) Write(record element.Record) (err error) {
	var records []string
	for i := 0; i < record.ColumnNumber(); i++ {
		var col element.Column
		if col, err = record.GetByIndex(i); err != nil {
			return
		}
		var s string
		if s, err = w.getRecord(col, i); err != nil {
			return
		}
		records = append(records, s)
	}
	return w.writer.Write(records)
}

func (w *Writer) getRecord(col element.Column, i int) (s string, err error) {
	if c, ok := w.columns[i]; ok && element.ColumnType(c.Type) == element.TypeTime {
		var t time.Time
		if t, err = col.AsTime(); err != nil {
			return
		}
		s = t.Format(c.layout())
		return
	}
	s = col.String()
	return
}
