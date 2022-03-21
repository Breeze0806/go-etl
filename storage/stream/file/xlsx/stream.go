package xlsx

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/element"
	"github.com/Breeze0806/go-etl/storage/stream/file"
	"github.com/xuri/excelize/v2"
)

func init() {
	var opener Opener
	file.RegisterOpener("xlsx", &opener)
	var creater Creater
	file.RegisterCreater("xlsx", &creater)
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
	file     *excelize.File
	filename string
}

func NewInStream(filename string) (file.InStream, error) {
	stream := &Stream{}
	var err error
	stream.file, err = excelize.OpenFile(filename)
	if err != nil {
		return nil, err
	}
	return stream, nil
}

func NewOutStream(filename string) (file.OutStream, error) {
	stream := &Stream{
		filename: filename,
	}
	stream.file = excelize.NewFile()
	return stream, nil
}

func (s *Stream) Rows(conf *config.JSON) (file.Rows, error) {
	return NewRows(s.file, conf)
}

func (s *Stream) Writer(conf *config.JSON) (file.StreamWriter, error) {
	return NewWriter(s.file, conf)
}

func (s *Stream) Close() (err error) {
	if s.filename != "" {
		err = s.file.SaveAs(s.filename)
		return
	}
	return s.file.Close()
}

type Rows struct {
	*excelize.Rows

	columns map[int]Column
}

func NewRows(f *excelize.File, c *config.JSON) (rows *Rows, err error) {
	var conf *InConfig
	if conf, err = NewInConfig(c); err != nil {
		return
	}
	rows = &Rows{
		columns: make(map[int]Column),
	}
	for _, v := range conf.Columns {
		rows.columns[v.index()] = v
	}
	rows.Rows, err = f.Rows(conf.Sheet)
	return
}

func (r *Rows) Scan() (columns []element.Column, err error) {
	var record []string
	record, err = r.Columns()
	if err != nil {
		return nil, err
	}

	for i, v := range record {
		var c element.Column
		c, err = r.getColum(i, v)
		if err != nil {
			return nil, err
		}
		columns = append(columns, c)
	}
	return
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

type Writer struct {
	file       *excelize.File
	writer     *excelize.StreamWriter
	conf       *OutConfig
	row        int
	sheetIndex int
	columns    map[int]Column
}

func NewWriter(f *excelize.File, c *config.JSON) (file.StreamWriter, error) {
	w := &Writer{
		file:    f,
		columns: make(map[int]Column),
	}
	var err error
	w.conf, err = NewOutConfig(c)
	if err != nil {
		return nil, err
	}
	for _, v := range w.conf.Columns {
		w.columns[v.index()] = v
	}
	var name string
	name, err = w.getSheetName()
	if err != nil {
		return nil, err
	}
	w.writer, err = w.file.NewStreamWriter(name)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (w *Writer) Write(record element.Record) (err error) {
	w.row++
	if w.row == excelize.TotalRows {
		if err = w.writer.Flush(); err != nil {
			return
		}
		w.row = 1
		var name string
		name, err = w.getSheetName()
		if err != nil {
			return
		}
		w.writer, err = w.file.NewStreamWriter(name)
		if err != nil {
			return
		}
	}

	var records []interface{}
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
	axis, _ := excelize.CoordinatesToCellName(1, w.row)
	return w.writer.SetRow(axis, records)
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

func (w *Writer) Flush() (err error) {
	return
}

func (w *Writer) Close() (err error) {
	err = w.writer.Flush()
	return
}

func (w *Writer) getSheetName() (string, error) {
	if w.sheetIndex < len(w.conf.Sheets) {
		w.sheetIndex++
		return w.conf.Sheets[w.sheetIndex-1], nil
	}
	return "", fmt.Errorf("index out of range in sheets")
}
