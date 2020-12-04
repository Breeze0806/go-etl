package plugin

import "github.com/Breeze0806/go-etl/datax/common/config"

type Defalut interface {
	Pluggable
	PreCheck()
	Prepare()
	Post()
	PreHandler(conf *config.Json)
	PostHandler(conf *config.Json)
}


