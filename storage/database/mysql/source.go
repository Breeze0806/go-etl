package mysql

import (
	"github.com/Breeze0806/go-etl/storage/database"
)

func init() {
	var d Dialect
	database.RegisterDialect(d.Name(), d)
}

type Dialect struct{}

func (d Dialect) Source(bs *database.BaseSource) (database.Source, error) {
	return NewSource(bs)
}

func (d Dialect) Name() string {
	return "mysql"
}

type Source struct {
	*database.BaseSource
	dsn string
}

func NewSource(bs *database.BaseSource) (s database.Source, err error) {
	source := &Source{
		BaseSource: bs,
	}
	var c *Config
	if c, err = NewConfig(source.Config()); err != nil {
		return
	}

	if source.dsn, err = c.FormatDSN(); err != nil {
		return
	}
	return source, nil
}

func (s *Source) DriverName() string {
	return "mysql"
}

func (s *Source) ConnectName() string {
	return s.dsn
}

func (s *Source) Key() string {
	return s.dsn
}

func (s *Source) Table(b *database.BaseTable) database.Table {
	return NewTable(b)
}

func Quoted(s string) string {
	return "`" + s + "`"
}
