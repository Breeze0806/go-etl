package xlsx

import (
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
}

type Opener struct {
}

func (o *Opener) Open(filename string) (file.InStream, error) {
	return NewInStream(filename)
}

type Stream struct {
	file *excelize.File
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

func (s *Stream) Rows(conf *config.JSON) (file.Rows, error) {
	return NewRows(s.file, conf)
}

func (s *Stream) Close() (err error) {
	return s.file.Close()
}

type Rows struct {
	*excelize.Rows

	columns map[int]Column
}

func NewRows(f *excelize.File, c *config.JSON) (rows *Rows, err error) {
	var conf *Config
	if conf, err = NewConfig(c); err != nil {
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
