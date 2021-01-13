package mysql

import (
	"encoding/json"
	"fmt"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/storage/database"
	"github.com/go-sql-driver/mysql"
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

func (s *Source) Table(b *database.BaseTable) database.Table {
	return NewTable(b)
}

type Config struct {
	URL      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewConfig(conf *config.Json) (c *Config, err error) {
	c = &Config{}
	err = json.Unmarshal([]byte(conf.String()), c)
	if err != nil {
		return nil, err
	}
	return
}

func (c *Config) FormatDSN() (dsn string, err error) {
	var mysqlConf *mysql.Config
	if mysqlConf, err = mysql.ParseDSN(c.URL); err != nil {
		return
	}
	fmt.Printf("%+v", mysqlConf)
	mysqlConf.User = c.Username
	mysqlConf.Passwd = c.Password
	mysqlConf.ParseTime = true
	dsn = mysqlConf.FormatDSN()
	return
}

func Quoted(s string) string {
	return "`" + s + "`"
}
