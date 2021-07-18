package postgres

import (
	"github.com/Breeze0806/go-etl/storage/database"
	"github.com/lib/pq"

	//postgres driver 可以指定网络访问的超时时间
	_ "github.com/Breeze0806/go/database/pqto"
)

func init() {
	var d Dialect
	database.RegisterDialect(d.Name(), d)
}

//Dialect postgres数据库方言
type Dialect struct{}

//Source 生产数据源
func (d Dialect) Source(bs *database.BaseSource) (database.Source, error) {
	return NewSource(bs)
}

//Name 数据库方言的注册名
func (d Dialect) Name() string {
	return "postgres"
}

//Source postgres数据源
type Source struct {
	*database.BaseSource //基础数据源

	dsn string
}

//NewSource 生成postgres数据源，在配置文件错误时会报错
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

//DriverName github.com/Breeze0806/go/database/pqto的驱动名
func (s *Source) DriverName() string {
	return "pgTimeout"
}

//ConnectName github.com/Breeze0806/go/database/pqto的数据源连接信息
func (s *Source) ConnectName() string {
	return s.dsn
}

//Key 数据源的关键字，用于DBWrapper的复用
func (s *Source) Key() string {
	return s.dsn
}

//Table 生成mysql的表
func (s *Source) Table(b *database.BaseTable) database.Table {
	return NewTable(b)
}

//Quoted postgres引用函数
func Quoted(s string) string {
	return pq.QuoteIdentifier(s)
}
