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

package db2

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/Breeze0806/go-etl/config"
)

//Config 配置
//ibm.com/docs/en/db2/11.1?topic=file-data-server-driver-configuration-keywords
type Config struct {
	URL      string `json:"url"`      //数据库url，包含数据库地址，数据库其他参数
	Username string `json:"username"` //用户名
	Password string `json:"password"` //密码
}

//NewConfig 创建mysql配置，如果格式不符合要求，就会报错
func NewConfig(conf *config.JSON) (c *Config, err error) {
	c = &Config{}
	err = json.Unmarshal([]byte(conf.String()), c)
	if err != nil {
		return nil, err
	}
	return
}

//FormatDSN 生成数据源连接信息，url有错会报错
func (c *Config) FormatDSN() (dsn string, err error) {
	o := make(values)
	if err = parseOpts(c.URL, o); err != nil {
		return
	}

	if _, ok := o["HOSTNAME"]; !ok {
		err = errors.New("HOSTNAME does not exist")
		return
	}

	if _, ok := o["PORT"]; !ok {
		err = errors.New("PORT does not exist")
		return
	}

	if _, ok := o["DATABASE"]; !ok {
		err = errors.New("DATABASE does not exist")
		return
	}

	o["UID"] = c.Username
	o["PWD"] = c.Password
	dsn = formatDSN(o)
	return
}

func formatDSN(o values) string {
	var kvs []string
	escaper := strings.NewReplacer(`'`, `\'`, `\`, `\\`)
	accrue := func(k, v string) {
		if v != "" {
			kvs = append(kvs, k+"="+escaper.Replace(v))
		}
	}
	for k, v := range o {
		accrue(k, v)
	}
	sort.Sort(sort.StringSlice(kvs))
	return strings.Join(kvs, ";")
}

type values map[string]string

func parseOpts(name string, o values) error {
	params := strings.Split(name, ";")
	for i, v := range params {
		ss := strings.Split(v, "=")
		if len(ss) != 2 {
			return fmt.Errorf("param %v(%v) has '='", i, v)
		}
		key := strings.TrimSpace(ss[0])
		value := strings.TrimSpace(ss[1])
		o[key] = value
	}
	return nil
}
